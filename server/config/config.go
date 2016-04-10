package config

import (
	"bitbucket.org/classroomsystems/ini"
	"fmt"
	"log"
	"reflect"
	"strings"
)

type Conf struct {
	conf ini.Config
}

type ConfValues struct {
	Watermark Watermark
	Server    server
}

type server struct {
	Tokens int `ini:"concurrency"`
	Memory uint64
}

type Watermark struct {
	ImagePath  string `ini:"path"`
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
	server{
		Tokens: 1,
		Memory: 2048 * 1024 * 1024,
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
