// +build debug

package logging

import (
	"log"
	"fmt"
)

func Debug(v ...interface{}) {
	log.Println(v...)
}

func Debugf(s string, v ...interface{}) {
	Debug(fmt.Sprintf(s, v))
}
