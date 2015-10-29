package main

import (
	"bytes"
	"os"
)

// Imageblob is just an image file dumped, byte by byte to an byte array.
type ImageBlob []byte

// Returns binary ImageBlob of an image.
func LoadImage(filename string) (blob ImageBlob, err error) {
	blob = nil

	reader, err := os.Open(filename)
	if err != nil {
		return
	}
	// Remember to free resources after you're done
	defer reader.Close()

	buffer := bytes.NewBuffer([]byte{})

	// Remember to check for errors
	_, err = buffer.ReadFrom(reader)
	if err != nil {
		return
	}

	blob = ImageBlob(buffer.Bytes())
	return
}
