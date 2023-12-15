//go:build !debug
// +build !debug

package logging

func Debug(v ...interface{})            {}
func Debugf(s string, v ...interface{}) {}
