package images

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// Returns Image from file.
func LoadImage(filename string) (img Image, err error) {

	reader, err := os.Open(filename)
	if err != nil {
		return
	}
	// Remember to free resources after you're done
	defer reader.Close()

	buffer := bytes.NewBuffer([]byte{})
	_, err = buffer.ReadFrom(reader)
	if err != nil {
		return
	}
	blob := ImageBlob(buffer.Bytes())

	img = NewImage()
	err = img.ReadImageBlob(blob)
	return
}

// Return binary ImageBlob of an image from web.
func LoadImageWeb(url string) (image Image, err error) {

	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("Couldn't load image. Server returned %i", resp.StatusCode))
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	blob := body

	image = NewImage()
	err = image.ReadImageBlob(blob)

	return
}


func ImageFromBlob(blob ImageBlob) (img Image, err error) {
	img = NewImage()
	err = img.ReadImageBlob(blob)
	return
}
