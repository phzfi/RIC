// +build debug

package logging

import "log"

func Debug(v ...interface{}) {
	log.Println(v...)
}
