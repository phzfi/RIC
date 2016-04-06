package config

import (
  "bitbucket.org/classroomsystems/ini"
  "github.com/phzfi/RIC/server/logging"
  "strconv"
  "errors"
)

type Conf struct{
  conf ini.Config
}

type DefaultConf struct{
   minHeight int
   minWidth int
   maxHeight int
   maxWidth int
   addMark bool
   imgpath string
   tokens int
   vertical float64
   horizontal float64
}

func ReadConfig(path string) (config Conf, err error) {
  conf, err := ini.LoadFile(path)
  if err != nil {
    logging.Debug("Error reading config " + err.Error())
    return
  }
  config = Conf {
    conf: conf,
  }
  return
}


func (conf Conf) GetString(section, key string) (value string, err error) {
  value, success := conf.conf.Get(section, key)
  if (success) {
    return
  } else {
    return "", errors.New("Value not found for "+ key +" in "+ section)
  }
}

func (conf Conf) GetInt(section, key string) (value int, err error) {
  str, success := conf.conf.Get(section, key)
  if success {
    return strconv.Atoi(str)
  } else {
    return 0, errors.New("Value not found for "+ key +" in "+ section)
  }
}

func (conf Conf) GetUint64(section, key string) (value uint64, err error) {
  str, success := conf.conf.Get(section, key)
  if success {
    return strconv.ParseUint(str, 10, 64)
  } else {
    return 0, errors.New("Value not found for "+ key +" in "+ section)
  }
}

func (conf Conf) GetFloat64(section, key string) (value float64, err error) {
  str, success := conf.conf.Get(section, key)
  if (success) {
    return strconv.ParseFloat(str, 64)
  } else {
    return 0.0, errors.New("Value not found for "+ key +" in "+ section)
  }
}

func (conf Conf) GetBool(section, key string) (value bool, err error) {
  str, success := conf.conf.Get(section, key)
  if (success) {
    return strconv.ParseBool(str)
  } else {
    return false, errors.New("Value not found for "+ key +" in "+ section)
  }
}
