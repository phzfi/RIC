package images

import (
  "testing"
)

func TestPreserveMetaData(t *testing.T) {
  imgBefore := NewImage()
  imgAfter := NewImage()
  defer imgBefore.Destroy()
  defer imgAfter.Destroy()

  err := imgBefore.FromFile("../testimages/resize/toresize.jpg")
  if err != nil {
    t.Fatal(err)
  }

  var propertiesBefore []string
  var profilesBefore []string
  for _, p := range imgBefore.GetImageProperties("*") {
    propertiesBefore = append(propertiesBefore, imgBefore.GetImageProperty(p))
  }
  for _, p := range imgBefore.GetImageProfiles("*") {
    profilesBefore = append(profilesBefore, imgBefore.GetImageProfile(p))
  }

  imgBefore.Resize(100, 100)
  imgBefore.Convert("webp")
  imgAfter.FromBlob(imgBefore.Blob())

  var propertiesAfter []string
  for _, p := range imgBefore.GetImageProperties("*") {
    propertiesAfter = append(propertiesAfter, imgBefore.GetImageProperty(p))
  }
  var profilesAfter []string
  for _, p := range imgAfter.GetImageProfiles("*") {
    profilesAfter = append(profilesAfter, imgAfter.GetImageProfile(p))
  }

  for i, _ := range propertiesBefore {
    if propertiesBefore[i] != propertiesAfter[i] {
      t.Fatal("Image properties do not match!")
    }
  }
  /*
  for i, _ := range profilesBefore {
    if profilesBefore[i] != profilesAfter[i] {
      t.Fatal("Image profiles do not match!")
    }
  }
  */
}
