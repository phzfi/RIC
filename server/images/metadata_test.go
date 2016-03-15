package images

import (
  "testing"
  "fmt"
  "bytes"
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
      if before[i] != after[i] {
        fmt.Println(before[i] + " - " + after[i])
        t.Fatal("Image metadata does not match! Different values.")
      }
    }
  } else {
    t.Fatal("Image metadata does not match! Different length.")
  }
}

func TestPreserveMetadataJpgResize(t *testing.T) {
  imgBefore := NewImage()
  imgAfter := NewImage()
  defer imgBefore.Destroy()
  defer imgAfter.Destroy()

  imgBefore.FromFile("../testimages/metadata/IPTC-Photometadata.jpg")
  //EXIFbefore := readImageProfiles(imgBefore, "exif")
  EXIFbefore := []byte(imgBefore.GetImageProfile("exif"))
  //IPTCbefore := readImageProfiles(imgBefore, "iptc")
  //XMP_before := readImageProfiles(imgBefore, "xmp")

  imgBefore.Resize(100, 100)
  imgAfter.FromBlob(imgBefore.Blob())
  EXIFafter := []byte(imgAfter.GetImageProfile("exif"))
  //EXIFafter := readImageProfiles(imgAfter, "exif")
  //IPTCafter := readImageProfiles(imgAfter, "iptc")
  //XMP_after := readImageProfiles(imgAfter, "xmp")

  fmt.Println(bytes.Equal(EXIFbefore, EXIFafter))
  //compare(t, EXIFbefore, EXIFafter)
  //compare(t, IPTCbefore, IPTCafter)
  //compare(t, XMP_before, XMP_after)
}

func TestPreserveMetadataJpgToPNG(t *testing.T) {
  imgBefore := NewImage()
  imgAfter := NewImage()
  defer imgBefore.Destroy()
  defer imgAfter.Destroy()

  imgBefore.FromFile("../testimages/metadata/IPTC-Photometadata.jpg")
  EXIFbefore := readImageProfiles(imgBefore, "exif")
  IPTCbefore := readImageProfiles(imgBefore, "iptc")
  XMP_before := readImageProfiles(imgBefore, "xmp")

  imgBefore.Convert("PNG")
  imgAfter.FromBlob(imgBefore.Blob())
  EXIFafter := readImageProfiles(imgAfter, "exif")
  IPTCafter := readImageProfiles(imgAfter, "iptc")
  XMP_after := readImageProfiles(imgAfter, "xmp")

  compare(t, EXIFbefore, EXIFafter)
  compare(t, IPTCbefore, IPTCafter)
  compare(t, XMP_before, XMP_after)
}

func TestPreserveMetadataJpgToTiff(t *testing.T) {
  imgBefore := NewImage()
  imgAfter := NewImage()
  defer imgBefore.Destroy()
  defer imgAfter.Destroy()

  imgBefore.FromFile("../testimages/metadata/IPTC-Photometadata.jpg")
  //EXIFbefore := readImageProfiles(imgBefore, "exif")
  IPTCbefore := imgBefore.GetImageProfile("iptc")
  //XMP_before := imgBefore.GetImageProfile("xmp")
  fmt.Println(imgBefore.GetImageProfiles("*"))

  imgBefore.Convert("TIFF")
  imgBefore.SetImageProfile("tiff", []byte(imgBefore.GetImageProfile("exif")))
  fmt.Println(imgBefore.GetImageProfiles("*"))
  imgAfter.FromBlob(imgBefore.Blob())
  fmt.Println(imgAfter.GetImageProfiles("*"))
  fmt.Println(imgAfter.GetImageProperties("*"))
  //EXIFafter := readImageProfiles(imgAfter, "exif")
  IPTCafter := imgAfter.GetImageProfile("iptc")
  //XMP_after := imgAfter.GetImageProfile("xmp")


  //compare(t, EXIFbefore, EXIFafter)
  fmt.Println(bytes.Equal([]byte(IPTCbefore), []byte(IPTCafter)))
  //compare(t, IPTCbefore, IPTCafter)
  //compare(t, XMP_before, XMP_after)
}
