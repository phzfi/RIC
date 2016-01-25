package images

import (
	"bytes"
	"errors"
	"fmt"
	//"io/ioutil"
	"github.com/valyala/fasthttp"
	"os"
)

// Returns Image from file.
func (img *Image) FromFile(filename string) error {

	reader, err := os.Open(filename)
	if err != nil {
		return err
	}
	// Remember to free resources after you're done
	defer reader.Close()

	buffer := bytes.NewBuffer([]byte{})
	_, err = buffer.ReadFrom(reader)
	if err != nil {
		return err
	}
	blob := ImageBlob(buffer.Bytes())

	return img.FromBlob(blob)
}

// Return binary ImageBlob of an image from web.
func (img *Image) FromWeb(url string) error {

	//resp, err := http.Get(url)
	statuscode, body, err := fasthttp.Get(nil, url)
	if err != nil {
		return err
	}
	//defer resp.Body.Close()

	if statuscode != 200 {
		return errors.New(fmt.Sprintf("Couldn't load image. Server returned %i", statuscode))
	}

	return img.FromBlob(body)
}

func (img *Image) FromBlob(blob ImageBlob) error {
	return img.ReadImageBlob(blob)
}
