package configuration

import (
  "bitbucket.org/classroomsystems/ini"
  "strconv"
  "errors"
)

var conf, err = ini.LoadFile("config.ini")

func GetString(section, key string) (string) {
  value, success := conf.Get(section, key)
  if (success) {
    return value
  } else {
    return ""
  }
}

func GetInt(section, key string) (value int, err error) {
  str, success := conf.Get(section, key)
  if success {
    return strconv.Atoi(str)
  } else {
    return 0, errors.New("Value not found for "+ key +" in "+ section)
  }
}

func GetUint64(section, key string) (value uint64, err error) {
  str, success := conf.Get(section, key)
  if success {
    return strconv.ParseUint(str, 10, 64)
  } else {
    return 0, errors.New("Value not found for "+ key +" in "+ section)
  }
}

func GetFloat64(section, key string) (value float64, err error) {
  str, success := conf.Get(section, key)
  if (success) {
    return strconv.ParseFloat(str, 64)
  } else {
    return 0.0, errors.New("Value not found for "+ key +" in "+ section)
  }
}

func GetBool(section, key string) (value bool, err error) {
  str, success := conf.Get(section, key)
  if (success) {
    return strconv.ParseBool(str)
  } else {
    return false, errors.New("Value not found for "+ key +" in "+ section)
  }
}
