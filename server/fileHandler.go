package main

import (
	"errors"
	"github.com/phzfi/RIC/server/logging"
	"github.com/phzfi/RIC/server/ops"
	"github.com/valyala/fasthttp"
	"strings"
	"encoding/base64"
	"fmt"
	"net/http"
	"bufio"
	"os"
	"io"
	"crypto/md5"
)

func HandleReceiveFile(uri *fasthttp.URI, source ops.ImageSource, marker ops.Watermarker) (operations []ops.Operation, format string, err, invalid error) {
	filename := string(uri.Path())

	w, h, cropx, cropy, mode, format, requestUrl, invalid := getParams(uri.QueryArgs())

	if invalid != nil {
		logging.Debug(invalid)
		return
	}

	decodedPath, md5Filename := decodeFilename(filename)

	rootDir, rootErr := source.GetDefaultRoot()
	if rootErr != nil {
		logging.Debug(rootErr)
		return
	}
	//TODO: Check that the domain/url is allowed (we don't want to work as a proxy)
	filePath := rootDir + "/" + md5Filename

	if !fileExists(filePath) {

		resp, httpErr := http.Get(decodedPath)
		if httpErr != nil {
			//log.Fatal(httpErr)
			logging.Debug(fmt.Sprintf("failed to retrieve external image: %s :%s",  filename, httpErr))
			return
		}
		file, fileErr := os.OpenFile(filePath, os.O_CREATE, 0644)
		file, fileErr = os.OpenFile(filePath, os.O_WRONLY, 0644)
		if fileErr != nil {
			return
		}
		bufferedWriter := bufio.NewWriter(file)
		buffer := make([]byte, 4096)
		for {
			var bytesWritten = 0
			bytesToWrite, readErr := resp.Body.Read(buffer)
			if readErr != nil && readErr != io.EOF {
				return
			}
			writeIndex := 0
			for bytesWritten < bytesToWrite {
				bytesWrote, writeErr := bufferedWriter.Write(buffer[writeIndex:bytesToWrite])
				writeIndex = writeIndex + bytesWrote
				bytesWritten = bytesWritten + bytesWrote
				if writeErr != nil {
					return
				}
			}
			if readErr == io.EOF {break}
		}
		bufferedWriter.Flush()
		file.Close()
	}
	filename = md5Filename


	if requestUrl != "" {
		source.AddRoot(requestUrl)
	}

	ow, oh, err := source.ImageSize(filename)
	if err != nil {
		return
	}

	operations = []ops.Operation{source.LoadImageOp(filename)}


	adjustWidth := func() {
		w = roundedIntegerDivision(h*ow, oh)
	}

	adjustHeight := func() {
		h = roundedIntegerDivision(w*oh, ow)
	}

	adjustSize := func() {
		if h == 0 && w != 0 {
			adjustHeight()
		} else if h != 0 && w == 0 {
			adjustWidth()
		} else if w == 0 && h == 0 {
			w, h = ow, oh
		}
	}

	denyUpscale := func() {
		h0 := h
		w0 := w
		if w > ow {
			h = roundedIntegerDivision(ow*h0, w0)
			w = ow
		}
		if h > oh || h > h0 {
			w = roundedIntegerDivision(oh*w0, h0)
			h = oh
		}
	}

	resize := func() {
		denyUpscale()
		adjustSize()
		operations = append(operations, ops.Resize{w, h})
	}

	liquid := func() {
		denyUpscale()
		adjustSize()
		operations = append(operations, ops.LiquidRescale{w, h})
	}

	crop := func() {
		if w == 0 {
			w = ow
		}
		if h == 0 {
			h = oh
		}
		operations = append(operations, ops.Crop{w, h, cropx, cropy})
	}

	cropmid := func() {
		if w == 0 || w > ow {
			w = ow
		}
		if h == 0 || h > oh {
			h = oh
		}
		midW := roundedIntegerDivision(ow, 2)
		midH := roundedIntegerDivision(oh, 2)
		cropx := midW - roundedIntegerDivision(w, 2)
		cropy := midH - roundedIntegerDivision(h, 2)
		operations = append(operations, ops.Crop{w, h, cropx, cropy})
	}

	fit := func() {
		if w > ow {
			w = ow
		}
		if h > oh {
			h = oh
		}
		if w != 0 && h != 0 {
			if ow*h > w*oh {
				adjustHeight()
			} else {
				adjustWidth()
			}
			operations = append(operations, ops.Resize{w, h})
		} else {
			resize()
		}
	}

	watermark := func() {
		heightOK := h > marker.MinHeight && h < marker.MaxHeight
		widthOK := w > marker.MinWidth && w < marker.MaxWidth
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

func decodeFilename(filename string) (decodedPath string, md5Filename string, ) {
	// Check and download an image
	//encoded, _ := url.Parse(requestUrl)
	decoded, encodeErr := base64.StdEncoding.DecodeString(filename[1:])
	if encodeErr != nil {
		logging.Debug("invalid request filename format:", filename)
		return
	}
	decodedPath = string(decoded)
	md5Hash := md5.New()
	io.WriteString(md5Hash, decodedPath)
	md5Filename = fmt.Sprintf("%x", md5Hash.Sum(nil))

	return

}

func DeleteFile(uri *fasthttp.URI, source ops.ImageSource,) (error){
	filename := string(uri.Path())

	decodedPath, md5Filename := decodeFilename(filename)
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

		logging.Debug("File deleted: " +  decodedPath)
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
	"":					"jpeg",
	"jpg":			"jpeg",
	"jpeg":			"jpeg",
	"gif":			"gif",
	"webp":			"webp",
	"bmp":			"bmp",
	"png":			"png",
	"tiff":			"tiff",
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
func getParams(a *fasthttp.Args) (w, h, cropx, cropy int, mode mode, format, url string, err error) {

	if strings.Contains(a.String(), "%3F") { // %3F = ?
		err = errors.New("Invalid characters in request!")
		return
	}

	defer func() {
		if msg := recover(); msg != nil {
			err = msg.(error)
		}
	}()

	w = getUint(a, widthParam)
	h = getUint(a, heightParam)

	cropx = getUint(a, cropxParam)
	cropy = getUint(a, cropyParam)

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