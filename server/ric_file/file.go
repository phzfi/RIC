package ric_file

import (
	"encoding/base64"
	"github.com/phzfi/RIC/server/logging"
	"fmt"
	"errors"
	"crypto/md5"
	"io"
)


func DecodeFilename(filename string) (remoteUrl string, md5Filename string, err error) {
	if len(filename[1:]) == 0 {
		err = errors.New("encoded filename missing")
		return
	}
	decoded, encodeErr := base64.StdEncoding.DecodeString(filename[1:])
	if encodeErr != nil {
		logging.Debug(fmt.Sprintf("invalid request filename format: %s, %s", encodeErr, filename))
		err = encodeErr
		return
	}
	remoteUrl = string(decoded)
	if len(remoteUrl) == 0 {
		err = errors.New("no filename could be extracted")
		return
	}

	md5Hash := md5.New()
	io.WriteString(md5Hash, remoteUrl)
	md5Filename = fmt.Sprintf("%x", md5Hash.Sum(nil))

	return
}
