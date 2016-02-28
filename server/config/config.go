package config

import (
  "bitbucket.org/classroomsystems/ini"
  "strconv"
  "errors"
)



var conf, err = ini.LoadFile("config.ini")

func GetString(section, key string) (value string) {
  value, success := conf.Get(section, key)
  if (success) {
    return value
  } else {
    return ""
  }
}

func GetInt(section, key string) (int, error) {
  str, success := conf.Get(section, key)
  if (success) {
    return strconv.Atoi(str)
  } else {
    return 0, errors.New("Value for " + key + " not found.")
  }
}

func GetFloat64(section, key string) (float64, error) {
  str, success := conf.Get(section, key)
  if (success) {
    return strconv.ParseFloat(str, 64)
  } else {
    return 0.0, errors.New("Value for " + key + " not found.")
  }
}

func GetBool(section, key string) (bool, error) {
  str, success := conf.Get(section, key)
  if (success) {
    return strconv.ParseBool(str)
  } else {
    return false, errors.New("Value for " + key + " not found.")
  }
}


func Config() (*ini.Config) {
  return &conf
}
