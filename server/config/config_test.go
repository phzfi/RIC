package config

import (
  "testing"
  "strconv"
)

// Test reading config from not existing file.
func TestReadConfigErr(t *testing.T) {
	_, err := ReadConfig("nosuchfile.ini")
	if err == nil {
		t.Fatal("Expected error")
	}
}

//Test that GetInt function returns correct value without errors.
func TestGetInt(t *testing.T) {
	conf, err := ReadConfig("testconfig.ini")
	if err != nil {
		t.Fatal(err)
	}

  res, err := conf.GetInt("watermark", "maxwidth")
  if err != nil {
		t.Fatal(err)
	}

  if (res != 1){
    t.Fatal("Wrong value. Expected 1, got: " + string(res))
  }
}

//Test that GetFloat64 function returns correct value without errors.
func TestGetFloat64(t *testing.T) {
	conf, err := ReadConfig("testconfig.ini")
	if err != nil {
		t.Fatal(err)
	}

  res, err := conf.GetFloat64("watermark", "vertical")
  if err != nil {
		t.Fatal(err)
	}

  if (res != 0.5){
    t.Fatal("Wrong value. Expected 0.5, got: " + strconv.FormatFloat(res, 'f', 3, 64))
  }
}

//Test that GetString function returns correct value without errors.
func TestGetString(t *testing.T) {
	conf, err := ReadConfig("testconfig.ini")
	if err != nil {
		t.Fatal(err)
	}

  res, err := conf.GetString("server", "cache")
  if err != nil {
		t.Fatal(err)
	}

  if (res != "lru"){
    t.Fatal("Wrong value. Expected lru, got: " + string(res))
  }
}

//Test that GetUint64 function returns correct value without errors.
func TestGetUint64(t *testing.T) {
	conf, err := ReadConfig("testconfig.ini")
	if err != nil {
		t.Fatal(err)
	}

  res, err := conf.GetUint64("watermark", "maxwidth")
  if err != nil {
		t.Fatal(err)
	}

  if (res != uint64(1)){
    t.Fatal("Wrong value. Expected 1, got: " + string(res))
  }
}

//Test that GetBool function returns correct value without errors.
func TestGetBool(t *testing.T) {
	conf, err := ReadConfig("testconfig.ini")
	if err != nil {
		t.Fatal(err)
	}

  res, err := conf.GetBool("watermark", "addmark")
  if err != nil {
		t.Fatal(err)
	}

  if (!res){
    t.Fatal("Wrong value. Expected true, got: " + strconv.FormatBool(res))
  }
}

//Test that GetInt raises errors when expected.
func TestGetIntErr(t *testing.T) {
	conf, err := ReadConfig("testconfig.ini")
	if err != nil {
		t.Fatal(err)
	}

  _, err = conf.GetInt("watermark", "addmark")
  if err == nil {
		t.Fatal("Expected error.")
  }

  _, err = conf.GetInt("nothing", "nothing")
  if err == nil {
		t.Fatal("Expected error.")
  }
}

//Test that GetString raises errors when expected.
func TestGetStringErr(t *testing.T) {
	conf, err := ReadConfig("testconfig.ini")
	if err != nil {
		t.Fatal(err)
	}

  _, err = conf.GetString("nothing", "nothing")
  if err == nil {
		t.Fatal("Expected error.")
  }
}

//Test that GetUint64 raises errors when expected.
func TestGetUint64Err(t *testing.T) {
	conf, err := ReadConfig("testconfig.ini")
	if err != nil {
		t.Fatal("Expected error.")
	}

  _, err = conf.GetUint64("watermark", "addmark")
  if err == nil {
		t.Fatal("Expected error.")
  }

  _, err = conf.GetUint64("nothing", "nothing")
  if err == nil {
		t.Fatal("Expected error.")
  }
}

//Test that GetFloat64 raises errors when expected.
func TestGetFloat64Err(t *testing.T) {
	conf, err := ReadConfig("testconfig.ini")
	if err != nil {
		t.Fatal(err)
	}

  _, err = conf.GetFloat64("watermark", "addmark")
  if err == nil {
		t.Fatal("Expected error.")
  }

  _, err = conf.GetFloat64("nothing", "nothing")
  if err == nil {
		t.Fatal("Expected error.")
  }
}

//Test that GetBool raises errors when expected.
func TestGetBoolErr(t *testing.T) {
	conf, err := ReadConfig("testconfig.ini")
	if err != nil {
		t.Fatal(err)
	}

  _, err = conf.GetBool("watermark", "vertical")
  if err == nil {
		t.Fatal("Expected error.")
  }

  _, err = conf.GetBool("nothing", "nothing")
  if err == nil {
		t.Fatal("Expected error.")
  }
}
