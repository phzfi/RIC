package main

import (
	"errors"
	"fmt"
	"github.com/phzfi/RIC/server/config"
	"github.com/phzfi/RIC/server/logging"
	"github.com/phzfi/RIC/server/ops"
	"github.com/phzfi/RIC/server/ric_file"
	"github.com/valyala/fasthttp"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func HandleReceiveFile(uri *fasthttp.URI, source ops.ImageSource, serverWhitelistConfigPath string) (filename string, err error) {

	filePath, md5Filename, _, filePathErr := getFilePathInSystem(uri, source)
	if filePathErr != nil {
		err = filePathErr
		return
	}

	if !fileExists(filePath) {
		err = errors.New(fmt.Sprintf("file does not exist: %v", filePath))
		HandleRequestExternalFile(uri, source, serverWhitelistConfigPath)

		return
	}

	filename = md5Filename

	return
}

func HandleRequestExternalFile(uri *fasthttp.URI, source ops.ImageSource, serverWhitelistConfigPath string) (filename string, err error) {

	filePath, md5Filename, remoteUrl, filePathErr := getFilePathInSystem(uri, source)
	if filePathErr != nil {
		err = filePathErr
		return
	}

	if !fileExists(filePath) {
		requestUrl, uriErr := url.ParseRequestURI(remoteUrl)
		if uriErr != nil {
			logging.Debugf("Invalid url given as parameter: %s", remoteUrl)
			err = uriErr
			return
		}

		if !isPathAllowed(serverWhitelistConfigPath, requestUrl.Host) {
			logging.Debugf("unauthorized url given: %v", remoteUrl)
			err = errors.New("unauthorized url given")
			return
		}
		getRemoteFile(remoteUrl, filePath)
	}

	filename = md5Filename

	return
}

func getRemoteFile(url string, filePath string) (err error) {

	go func(url string) {
		fmt.Printf("Fetching %s \n", url)
		response, httpErr := http.Get(url)

		if httpErr != nil {
			//log.Fatal(httpErr)
			logging.Debug(fmt.Sprintf("failed to retrieve external image: %s, %s", url, httpErr))
			err = httpErr
			return
		}
		defer response.Body.Close()

		if response.StatusCode != 200 || !strings.Contains(response.Header.Get("Content-Type"), "image/") {
			logging.Debug(fmt.Sprintf("invalid response received: status code: %v , Content-type:%v", response.StatusCode, response.Header.Get("Content-Type")))
			err = errors.New("invalid response received")
			return
		}

		file, copyErr := os.Create(filePath)
		defer file.Close()

		_, copyErr = io.Copy(file, response.Body)
		if copyErr != nil {
			err = copyErr
			return
		}
	}(url)

	return
}

func getFilePathInSystem(uri *fasthttp.URI, source ops.ImageSource) (filePath string, md5Filename string, remoteUrl string, err error) {
	rawFilename := string(uri.Path())

	remoteUrl, md5Filename, decodeErr := ric_file.DecodeFilename(rawFilename)
	if decodeErr != nil {
		logging.Debug(decodeErr)
		err = decodeErr
		return
	}
	rootDir, rootErr := source.GetDefaultRoot()
	if rootErr != nil {
		logging.Debug(err)
		err = rootErr
		return
	}
	//TODO: Check that the domain/url is allowed (we don't want to work as a proxy)
	filePath = rootDir + "/" + md5Filename

	return
}

func GetFileSize(uri *fasthttp.URI, source ops.ImageSource) (fileSize int64, err error) {

	filePath, _, _, filePathErr := getFilePathInSystem(uri, source)
	if filePathErr != nil {
		err = filePathErr
		return
	}

	stat, fileErr := os.Stat(filePath)

	if fileErr != nil {
		err = fileErr
		return
	}

	fileSize = stat.Size()

	return
}

func CreateOperations(filename string, uri *fasthttp.URI, source ops.ImageSource, marker ops.Watermarker) (operations []ops.Operation, format string, err, invalid error) {

	width, height, cropX, cropY, mode, format, requestUrl, invalid := getParams(uri.QueryArgs())

	if invalid != nil {
		logging.Debug(invalid)
		return
	}

	if requestUrl != "" {
		source.AddRoot(requestUrl)
	}

	ow, oh, err := source.ImageSize(filename)
	if err != nil {
		return
	}

	operations = []ops.Operation{source.LoadImageOp(filename)}

	adjustWidth := func() {
		width = roundedIntegerDivision(height*ow, oh)
	}

	adjustHeight := func() {
		height = roundedIntegerDivision(width*oh, ow)
	}

	adjustSize := func() {
		if height == 0 && width != 0 {
			adjustHeight()
		} else if height != 0 && width == 0 {
			adjustWidth()
		} else if width == 0 && height == 0 {
			width, height = ow, oh
		}
	}

	denyUpscale := func() {
		h0 := height
		w0 := width
		if width > ow {
			height = roundedIntegerDivision(ow*h0, w0)
			width = ow
		}
		if height > oh || height > h0 {
			width = roundedIntegerDivision(oh*w0, h0)
			height = oh
		}
	}

	resize := func() {
		denyUpscale()
		adjustSize()
		operations = append(operations, ops.Resize{width, height})
	}

	liquid := func() {
		denyUpscale()
		adjustSize()
		operations = append(operations, ops.LiquidRescale{width, height})
	}

	crop := func() {
		if width == 0 {
			width = ow
		}
		if height == 0 {
			height = oh
		}
		operations = append(operations, ops.Crop{width, height, cropX, cropY})
	}

	cropmid := func() {
		if width == 0 || width > ow {
			width = ow
		}
		if height == 0 || height > oh {
			height = oh
		}
		midW := roundedIntegerDivision(ow, 2)
		midH := roundedIntegerDivision(oh, 2)
		cropx := midW - roundedIntegerDivision(width, 2)
		cropy := midH - roundedIntegerDivision(height, 2)
		operations = append(operations, ops.Crop{width, height, cropx, cropy})
	}

	fit := func() {
		if width > ow {
			width = ow
		}
		if height > oh {
			height = oh
		}
		if width != 0 && height != 0 {
			if ow*height > width*oh {
				adjustHeight()
			} else {
				adjustWidth()
			}
			operations = append(operations, ops.Resize{width, height})
		} else {
			resize()
		}
	}

	watermark := func() {
		heightOK := height > marker.MinHeight && height < marker.MaxHeight
		widthOK := width > marker.MinWidth && width < marker.MaxWidth
		if marker.AddMark && heightOK && widthOK {
			logging.Debug("Adding watermarkOp")
			operations = append(operations, ops.WatermarkOp(marker.WatermarkImage, marker.Horizontal, marker.Vertical))
		}
	}

	switch mode {
	case resizeMode:
		resize()
	case fitMode:
		fit()
	case liquidMode:
		liquid()
	case cropMode:
		crop()
	case cropmidMode:
		cropmid()
	}

	if true == false {
		watermark()
	}

	operations = append(operations, ops.Convert{format})

	return
}

func DeleteFile(uri *fasthttp.URI, source ops.ImageSource) error {
	filename := string(uri.Path())

	decodedPath, md5Filename, decodeErr := ric_file.DecodeFilename(filename)
	if decodeErr != nil {
		logging.Debug(decodeErr)
		return decodeErr
	}
	logging.Debug(fmt.Sprintf("Attempting to delete file: %s (%s)", md5Filename, decodedPath))
	rootDir, rootErr := source.GetDefaultRoot()
	if rootErr != nil {
		logging.Debug(rootErr)
		return rootErr
	}

	filePath := rootDir + "/" + md5Filename

	if fileExists(filePath) {
		removeErr := os.Remove(filePath)
		if removeErr != nil {
			logging.Debug(rootErr)
			return errors.New("failed to delete file")
		}

		logging.Debugf("File deleted: %s (%s)", md5Filename, decodedPath)
		return nil
	} else {
		return errors.New("file does not exist")
	}
}

func roundedIntegerDivision(n, m int) int {
	if (n < 0) == (m < 0) {
		return (n + m/2) / m
	} else { // -5 / 6 should round to -1
		return (n - m/2) / m
	}
}

var stringToMode = map[string]mode{
	"":        resizeMode,
	"resize":  resizeMode,
	"fit":     fitMode,
	"crop":    cropMode,
	"cropmid": cropmidMode,
	"liquid":  liquidMode,
}

var supportedFormats = map[string]string{
	"":     "jpeg",
	"jpg":  "jpeg",
	"jpeg": "jpeg",
	"gif":  "gif",
	"webp": "webp",
	"bmp":  "bmp",
	"png":  "png",
	"tiff": "tiff",
}

type mode int

const (
	fitMode = mode(1 + iota)
	cropMode
	cropmidMode
	liquidMode
	resizeMode

	widthParam  = "width"
	heightParam = "height"
	modeParam   = "mode"
	formatParam = "format"
	cropxParam  = "cropx"
	cropyParam  = "cropy"
	urlParam    = "url"
)

// returns validated parameters from request and error if invalid
func getParams(a *fasthttp.Args) (width int, height int, cropX int, cropY int, mode mode, format, url string, err error) {

	if strings.Contains(a.String(), "%3F") { // %3F = ?
		err = errors.New("Invalid characters in request!")
		return
	}

	defer func() {
		if msg := recover(); msg != nil {
			err = msg.(error)
		}
	}()

	width = getUint(a, widthParam)
	height = getUint(a, heightParam)

	cropX = getUint(a, cropxParam)
	cropY = getUint(a, cropyParam)

	mode = stringToMode[string(a.Peek(modeParam))]
	if mode == 0 {
		err = errors.New("Invalid mode!")
		return
	}

	format, formatFound := supportedFormats[strings.ToLower(string(a.Peek(formatParam)))]
	if !formatFound {
		err = errors.New("Invalid format '" + string(a.Peek(formatParam)) + "'!")
		return
	}
	// TODO: verify that the format is one we support.
	// We do not want to support TXT, for instance

	url = string(a.Peek(urlParam))

	a.Del(widthParam)
	a.Del(heightParam)
	a.Del(modeParam)
	a.Del(formatParam)
	a.Del(cropxParam)
	a.Del(cropyParam)
	a.Del(urlParam)

	if a.Len() != 0 {
		err = errors.New("Invalid parameter " + a.String())
		return
	}

	err = nil
	return
}

func getUint(a *fasthttp.Args, param string) int {
	v, err := a.GetUint(param)
	if isParseError(err) {
		panic(err)
	}
	if v == -1 {
		v = 0
	}
	return v
}

func isParseError(err error) bool {
	return err != nil && err != fasthttp.ErrNoArgValue
}

// Exists reports whether the named file or directory exists.
func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func isPathAllowed(configPath string, host string) (allowed bool) {
	allowedHosts, err := config.ReadHostWhitelist(configPath)

	if err != nil {
		return false
	}
	for _, allowedHost := range allowedHosts {
		if allowedHost == host {

			logging.Debugf("Permission for host %v found", host)
			return true
		}
	}
	return false
}
