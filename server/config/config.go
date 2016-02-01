package config

import (
  "bitbucket.org/classroomsystems/ini"
  "strconv"
)

var conf, err = ini.LoadFile("config.ini")

func Watermark(key string) (string) {
  value, success := conf.Get("watermark", key)
  if (success) {
    return value
  } else {
    return ""
  }
}

func WatermarkInt(key string) (value int, err error) {
  str, success := conf.Get("watermark", key)
  if success {
    return strconv.Atoi(str)
  } else {
    return 0, nil
  }
}

func WatermarkFloat64(key string) (value float64, err error) {
  str, success := conf.Get("watermark", key)
  if (success) {
    return strconv.ParseFloat(str, 64)
  } else {
    return 0.0, nil
  }
}


func Server(key string) (string) {
  value, success := conf.Get("server", key)
  if (success) {
    return value
  } else {
    return ""
  }
}

func ServerInt(key string) (value int, err error) {
  str, success := conf.Get("server", key)
  if (success) {
    return strconv.Atoi(str)
  } else {
    return 0, nil
  }
}


func Config() (*ini.Config) {
  return &conf
}
