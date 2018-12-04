package ops

import (
	"testing"
)

func TestImageSize(t *testing.T) {
	s := MakeImageSource()
	s.AddRoot("https://upload.wikimedia.org/wikipedia/commons/")
	s.AddRoot("../")

        w, h, err := s.ImageSize("b/b4/JPEG_example_JPG_RIP_100.jpg")
        if err != nil {
            t.Fatal(err)
        }
        if w != 313 || h != 234 {
            t.Fatal("Wrong image size returned. got w: %v, h:%v. expected w:%v, h:%v", w, h, 313, 234)
        }

        w, h, err = s.ImageSize("testimages/resize/toresize.jpg")
        if err != nil {
            t.Fatal(err)
        }
        if w != 1200 || h != 900 {
            t.Fatal("Wrong image size returned. got w: %v, h:%v. expected w:%v, h:%v", w, h, 1200, 900)
        }

        w, h, err = s.ImageSize("IMAGETHATDOESNOTEXIST")
        if err == nil {
            t.Fatal("No error returned when trying to get size of non existing image from fs.")
        }

        w, h, err = s.ImageSize("https://IMAGETHATDOESNOTEXIST")
        if err == nil {
            t.Fatal("No error returned when trying to get size of non existing image from fs.")
        }

}
