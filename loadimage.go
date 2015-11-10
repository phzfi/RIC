package main

import (
	"bytes"
	"os"
        "net/http"
        "io/ioutil"
        "errors"
        "fmt"
)



// Imageblob is just an image file dumped, byte by byte to an byte array.
type ImageBlob []byte

// Returns binary ImageBlob of an image from local filesystem.
func LoadImage(filename string) (blob ImageBlob, err error) {

	reader, err := os.Open(filename)
	if err != nil {
		return
	}

	buffer := bytes.NewBuffer([]byte{})
	buffer.ReadFrom(reader)
	blob = ImageBlob(buffer.Bytes())
	return
}

// Return binary ImageBlob of an image from web.
func LoadImageWeb(url string) (blob ImageBlob, err error){
    
    resp, err := http.Get(url);
    defer resp.Body.Close()
    if err != nil {
	return
    }
    
    if resp.StatusCode != 200 {
        err = errors.New(fmt.Sprintf("Couldn't load image. Server returned %i", resp.StatusCode))
        return
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return 
    }

    blob = body
    return
}

