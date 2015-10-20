package main

import (
	"bytes"
	"os"
)

// Imageblob is just an image file dumped, byte by byte to an byte array.
type ImageBlob []byte

// Returns binary ImageBlob of an image.
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
