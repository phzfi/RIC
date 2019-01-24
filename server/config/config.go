package config

import (
	"bitbucket.org/classroomsystems/ini"
	"errors"
	"fmt"
	"github.com/phzfi/RIC/server/logging"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type Conf struct {
	conf ini.Config
}

type ConfValues struct {
	Watermark Watermark
	Server    Server
}

type Server struct {
	Tokens              int `ini:"concurrency"`
	Memory              uint64
	ImageFolder         string
	CacheFolder         string
	HostWhitelistConfig string
	Port                int
}

type Watermark struct {
	ImagePath  string `ini:"path"`
	CachePath  string `ini:"path"`
	Horizontal float64
	Vertical   float64
	MaxWidth   int
	MinWidth   int
	MaxHeight  int
	MinHeight  int
	AddMark    bool
}

var defaults = ConfValues{
	Watermark{
		MinHeight:  200,
		MinWidth:   200,
		MaxHeight:  5000,
		MaxWidth:   5000,
		AddMark:    false,
		ImagePath:  "",
		Vertical:   0.0,
		Horizontal: 1.0,
	},
	Server{
		Tokens:              1,
		Memory:              2048 * 1024 * 1024,
		ImageFolder:         "",
		CacheFolder:         "",
		HostWhitelistConfig: "",
		Port:                8005,
	},
}

func ReadConfig(path string) (c *ConfValues) {
	copyOfDefaults := defaults
	c = &copyOfDefaults

	conf, err := ini.LoadFile(path)
	if err != nil {
		log.Println("Error reading config " + err.Error())
		return
	}

	whole := reflect.ValueOf(c).Elem()

	for categoryNo := 0; categoryNo < whole.NumField(); categoryNo++ {
		category := whole.Field(categoryNo)
		categoryName := strings.ToLower(whole.Type().Field(categoryNo).Name)

		for fieldNo := 0; fieldNo < category.NumField(); fieldNo++ {

			fieldInfo := category.Type().Field(fieldNo)
			key := fieldInfo.Tag.Get("ini")
			if key == "" {
				key = strings.ToLower(fieldInfo.Name)
			}

			field := category.Field(fieldNo)

			valueAsText, found := conf.Get(categoryName, key)
			if !found {
				// TODO: It would be better if it would only notify if a field name if written wrong instead of complaining for each missing field.
				log.Printf("%s not found in [%s], using default value of %#v.\n", key, categoryName, field.Interface())
				continue
			}
			_, err := fmt.Sscan(valueAsText, field.Addr().Interface())
			if err != nil {
				log.Printf("Error parsing %s: %s\n", valueAsText, err)
			}
		}
	}

	return
}

func ReadHostWhitelist(configPath string) (hosts []string, err error) {
	//config := ReadConfig(configPath)

	location, err := filepath.Abs(configPath)
	if _, err = os.Stat(location); os.IsNotExist(err) {
		return
	}

	configData, err := ini.LoadFile(location)
	configValue, ok := configData.Get("hosts", "allowed")

	if !ok {
		logging.Debug("Could not get server whitelist")
		err = errors.New("no server whitelist")
		fmt.Println("Failed to read server whitelist configuration file. Check 'HostWhitelistConfig'-entry in ric_config.ini")
		return
	}

	hosts = strings.Split(configValue, ",")

	return
}
