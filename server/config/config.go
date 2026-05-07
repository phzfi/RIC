package config

import (
	"log"
	"os"
	"strconv"
)

type ConfValues struct {
	Watermark Watermark
	Server    Server
}

type Server struct {
	Tokens int
	Memory uint64
}

type Watermark struct {
	ImagePath  string
	Horizontal float64
	Vertical   float64
	MaxWidth   int
	MinWidth   int
	MaxHeight  int
	MinHeight  int
	AddMark    bool
	ForceMark  string
	Text       string
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
		Tokens: 1,
		Memory: 2048 * 1024 * 1024,
	},
}

func ReadConfig() *ConfValues {
	c := defaults

	if v := os.Getenv("WATERMARK_PATH"); v != "" {
		c.Watermark.ImagePath = v
	}
	if v := os.Getenv("WATERMARK_HORIZONTAL"); v != "" {
		if val, err := strconv.ParseFloat(v, 64); err == nil {
			c.Watermark.Horizontal = val
		} else {
			log.Printf("Error parsing WATERMARK_HORIZONTAL: %s\n", err)
		}
	}
	if v := os.Getenv("WATERMARK_VERTICAL"); v != "" {
		if val, err := strconv.ParseFloat(v, 64); err == nil {
			c.Watermark.Vertical = val
		} else {
			log.Printf("Error parsing WATERMARK_VERTICAL: %s\n", err)
		}
	}
	if v := os.Getenv("WATERMARK_MAXWIDTH"); v != "" {
		if val, err := strconv.Atoi(v); err == nil {
			c.Watermark.MaxWidth = val
		} else {
			log.Printf("Error parsing WATERMARK_MAXWIDTH: %s\n", err)
		}
	}
	if v := os.Getenv("WATERMARK_MINWIDTH"); v != "" {
		if val, err := strconv.Atoi(v); err == nil {
			c.Watermark.MinWidth = val
		} else {
			log.Printf("Error parsing WATERMARK_MINWIDTH: %s\n", err)
		}
	}
	if v := os.Getenv("WATERMARK_MAXHEIGHT"); v != "" {
		if val, err := strconv.Atoi(v); err == nil {
			c.Watermark.MaxHeight = val
		} else {
			log.Printf("Error parsing WATERMARK_MAXHEIGHT: %s\n", err)
		}
	}
	if v := os.Getenv("WATERMARK_MINHEIGHT"); v != "" {
		if val, err := strconv.Atoi(v); err == nil {
			c.Watermark.MinHeight = val
		} else {
			log.Printf("Error parsing WATERMARK_MINHEIGHT: %s\n", err)
		}
	}
	if v := os.Getenv("WATERMARK_ADDMARK"); v != "" {
		if val, err := strconv.ParseBool(v); err == nil {
			c.Watermark.AddMark = val
		} else {
			log.Printf("Error parsing WATERMARK_ADDMARK: %s\n", err)
		}
	}
	if v := os.Getenv("WATERMARK_FORCEMARK"); v != "" {
		c.Watermark.ForceMark = v
	}
	if v := os.Getenv("WATERMARK_TEXT"); v != "" {
		c.Watermark.Text = v
	}
	if v := os.Getenv("SERVER_TOKENS"); v != "" {
		if val, err := strconv.Atoi(v); err == nil {
			c.Server.Tokens = val
		} else {
			log.Printf("Error parsing SERVER_TOKENS: %s\n", err)
		}
	}
	if v := os.Getenv("SERVER_MEMORY"); v != "" {
		if val, err := strconv.ParseUint(v, 10, 64); err == nil {
			c.Server.Memory = val
		} else {
			log.Printf("Error parsing SERVER_MEMORY: %s\n", err)
		}
	}

	return &c
}
