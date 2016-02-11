package images

import (
  "testing"
  "errors"
  "fmt"
)

func readImageProperties(img Image, prefix string) (properties []string) {
  for _, p := range img.GetImageProperties(prefix) {
    properties = append(properties, img.GetImageProperty(p))
  }
  return
}

func readImageProfiles(img Image, prefix string) (profiles []string) {
  for _, p := range img.GetImageProfiles(prefix) {
    profiles = append(profiles, img.GetImageProfile(p))
  }
  return
}

func compare(before []string, after []string) (err error) {
  if len(before) == len(after) {
    for i, _ := range before {
      if before[i] != before[i] {
        err = errors.New("Image metadata does not match! Different values.")
        return
      }
    }
  } else {
    err = errors.New("Image metadata does not match! Different length.")
  }
  return
}

func TestPreserveEXIFJpgToTif(t *testing.T) {
  imgBefore := NewImage()
  imgAfter := NewImage()
  defer imgBefore.Destroy()
  defer imgAfter.Destroy()

  imgBefore.FromFile("../testimages/resize/toresize.jpg")
  EXIFbefore := readImageProperties(imgBefore, "exif:*")
  fmt.Println(EXIFbefore)

  imgBefore.SetImageFormat("TIFF")
  imgAfter.FromBlob(imgBefore.Blob())

  EXIFafter := readImageProperties(imgAfter, "exif:*")
  fmt.Println(EXIFafter)

  err := compare(EXIFbefore, EXIFafter)
  if err != nil {
    t.Fatal(err)
  }
}
