package images

import (
	"testing"
)

func compare(t *testing.T, before []string, after []string) {
	if len(before) == len(after) {
		for i, _ := range before {
			// only compare field names: Some field values should be different since
			// they are changed in resize/convert (e.g. exif:thumbnail:XResolution).
			if before[i] != after[i] {
				t.Fatal("Image metadata does not match! Different fields.")
			}
		}
	} else {
		t.Fatal("Image metadata does not match! Different number of fields.")
	}
}

func TestPreserveMetadataJpgResize(t *testing.T) {
	imgBefore := NewImage()
	imgAfter := NewImage()
	defer imgBefore.Destroy()
	defer imgAfter.Destroy()

	imgBefore.FromFile("../testimages/metadata/IPTC-Photometadata.jpg")
	exifPropertiesBefore := imgBefore.GetImageProperties("exif:*")
	iptc4xmpPropertiesBefore := imgBefore.GetImageProperties("Iptc4xmp*")
	exifProfilesBefore := imgBefore.GetImageProfiles("exif")
	iptcProfilesBefore := imgBefore.GetImageProfiles("iptc")
	xmp_ProfilesBefore := imgBefore.GetImageProfiles("xmp")

	imgBefore.Resize(100, 100)
	imgAfter.FromBlob(imgBefore.Blob())
	exifPropertiesAfter := imgAfter.GetImageProperties("exif:*")
	iptc4xmpPropertiesAfter := imgAfter.GetImageProperties("Iptc4xmp*")
	exifProfilesAfter := imgAfter.GetImageProfiles("exif")
	iptcProfilesAfter := imgAfter.GetImageProfiles("iptc")
	xmp_ProfilesAfter := imgAfter.GetImageProfiles("xmp")

	compare(t, exifPropertiesBefore, exifPropertiesAfter)
	compare(t, iptc4xmpPropertiesBefore, iptc4xmpPropertiesAfter)
	compare(t, exifProfilesBefore, exifProfilesAfter)
	compare(t, iptcProfilesBefore, iptcProfilesAfter)
	compare(t, xmp_ProfilesBefore, xmp_ProfilesAfter)
}

func TestPreserveMetadataJpgToPNG(t *testing.T) {
	imgBefore := NewImage()
	imgAfter := NewImage()
	defer imgBefore.Destroy()
	defer imgAfter.Destroy()

	imgBefore.FromFile("../testimages/metadata/IPTC-Photometadata.jpg")
	exifPropertiesBefore := imgBefore.GetImageProperties("exif:*")
	iptc4xmpPropertiesBefore := imgBefore.GetImageProperties("Iptc4xmp*")
	exifProfilesBefore := imgBefore.GetImageProfiles("exif")
	iptcProfilesBefore := imgBefore.GetImageProfiles("iptc")
	xmp_ProfilesBefore := imgBefore.GetImageProfiles("xmp")

	imgBefore.Convert("PNG")
	imgAfter.FromBlob(imgBefore.Blob())
	exifPropertiesAfter := imgAfter.GetImageProperties("exif:*")
	iptc4xmpPropertiesAfter := imgAfter.GetImageProperties("Iptc4xmp*")
	exifProfilesAfter := imgAfter.GetImageProfiles("exif")
	iptcProfilesAfter := imgAfter.GetImageProfiles("iptc")
	xmp_ProfilesAfter := imgAfter.GetImageProfiles("xmp")

	compare(t, exifPropertiesBefore, exifPropertiesAfter)
	compare(t, iptc4xmpPropertiesBefore, iptc4xmpPropertiesAfter)
	compare(t, exifProfilesBefore, exifProfilesAfter)
	compare(t, iptcProfilesBefore, iptcProfilesAfter)
	compare(t, xmp_ProfilesBefore, xmp_ProfilesAfter)
}

func TestPreserveMetadataJpgToTiff(t *testing.T) {
	imgBefore := NewImage()
	imgAfter := NewImage()
	defer imgBefore.Destroy()
	defer imgAfter.Destroy()

	imgBefore.FromFile("../testimages/metadata/IPTC-Photometadata.jpg")
	iptc4xmpPropertiesBefore := imgBefore.GetImageProperties("Iptc4xmp*")
	iptcProfilesBefore := imgBefore.GetImageProfiles("iptc")
	xmp_ProfilesBefore := imgBefore.GetImageProfiles("xmp")

	imgBefore.Convert("TIFF")
	imgAfter.FromBlob(imgBefore.Blob())
	iptc4xmpPropertiesAfter := imgAfter.GetImageProperties("Iptc4xmp*")
	iptcProfilesAfter := imgAfter.GetImageProfiles("iptc")
	xmp_ProfilesAfter := imgAfter.GetImageProfiles("xmp")

	compare(t, iptc4xmpPropertiesBefore, iptc4xmpPropertiesAfter)
	compare(t, iptcProfilesBefore, iptcProfilesAfter)
	compare(t, xmp_ProfilesBefore, xmp_ProfilesAfter)
}

func TestPreserveMetadataJpgToGif(t *testing.T) {
	imgBefore := NewImage()
	imgAfter := NewImage()
	defer imgBefore.Destroy()
	defer imgAfter.Destroy()

	imgBefore.FromFile("../testimages/metadata/IPTC-Photometadata.jpg")
	iptcProfilesBefore := imgBefore.GetImageProfiles("iptc")

	imgBefore.Convert("GIF")
	imgAfter.FromBlob(imgBefore.Blob())
	iptcProfilesAfter := imgAfter.GetImageProfiles("iptc")

	compare(t, iptcProfilesBefore, iptcProfilesAfter)
}
