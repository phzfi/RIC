package images

import (
  "testing"
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

func compare(t *testing.T, before []string, after []string) {
  if len(before) == len(after) {
    for i, _ := range before {
      if before[i] != before[i] {
        t.Fatal("Image metadata does not match! Different values.")
      }
    }
  } else {
    t.Fatal("Image metadata does not match! Different length.")
  }
}

func TestPreserveEXIFJpgResize(t *testing.T) {
  imgBefore := NewImage()
  imgAfter := NewImage()
  defer imgBefore.Destroy()
  defer imgAfter.Destroy()

  imgBefore.FromFile("../testimages/resize/toresize.jpg")
  EXIFbefore := readImageProperties(imgBefore, "exif:*")

  imgBefore.Resize(100, 100)
  imgAfter.FromBlob(imgBefore.Blob())
  EXIFafter := readImageProperties(imgAfter, "exif:*")

  compare(t, EXIFbefore, EXIFafter)
}

func TestPreserveEXIFJpgToPNG(t *testing.T) {
  imgBefore := NewImage()
  imgAfter := NewImage()
  defer imgBefore.Destroy()
  defer imgAfter.Destroy()

  imgBefore.FromFile("../testimages/resize/toresize.jpg")
  EXIFbefore := readImageProperties(imgBefore, "exif:*")

  imgBefore.Convert("PNG")
  imgAfter.FromBlob(imgBefore.Blob())
  EXIFafter := readImageProperties(imgAfter, "exif:*")

  compare(t, EXIFbefore, EXIFafter)
}

func TestPreserveEXIFJpgToTiff(t *testing.T) {
  imgBefore := NewImage()
  imgAfter := NewImage()
  defer imgBefore.Destroy()
  defer imgAfter.Destroy()

  imgBefore.FromFile("../testimages/resize/toresize.jpg")
  EXIFbefore := readImageProperties(imgBefore, "exif:*")
  fmt.Println(imgBefore.GetImageProfiles("*"))

  imgBefore.Convert("TIFF")
  imgAfter.FromBlob(imgBefore.Blob())
  EXIFafter := readImageProperties(imgAfter, "exif:*")
  fmt.Println(imgAfter.GetImageProfiles("*"))

  compare(t, EXIFbefore, EXIFafter)
}
