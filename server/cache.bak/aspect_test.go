package cache

import (
	"fmt"
	"github.com/phzfi/RIC/server/images"
	"testing"
)

type DummyResizerForAspect [][2]uint

func (d *DummyResizerForAspect) GetImage(_ string, w, h uint) (images.ImageBlob, error) {
	*d = append(*d, [2]uint{w, h})
	return nil, nil
}

func (d DummyResizerForAspect) ImageSize(_ string) (uint, uint, error) {
	return 1000, 1000, nil
}

func (d DummyResizerForAspect) AddRoot(_ string) error {
	return nil
}

func (d DummyResizerForAspect) RemoveRoot(_ string) error {
	return nil
}

func TestAspectByWidth(t *testing.T) {
	testAspect(t, true)
}

func TestAspectByHeight(t *testing.T) {
	testAspect(t, false)
}

func TestFitImageByWidth(t *testing.T) {
	testFitImage(t, true)
}

func TestFitImageByHeight(t *testing.T) {
	testFitImage(t, false)
}

func testFitImage(t *testing.T, by_w bool) {
	dummy := DummyResizerForAspect{}
	cache := AspectPreserver{&dummy}

	if by_w {
		cache.GetImageFit("slnv", 2000, 200)
	} else {
		cache.GetImageFit("slnv", 200, 2000)
	}

	w := dummy[0][0]
	h := dummy[0][1]

	if w != 200 || h != 200 {
		t.Fatal(fmt.Sprintf("Image size was %d %d. Expected 200 200.", w, h))
	}
}

func testAspect(t *testing.T, by_w bool) {
	dummy := DummyResizerForAspect{}
	cache := AspectPreserver{&dummy}

	if by_w {
		cache.GetImageByWidth("slnv", 200)
	} else {
		cache.GetImageByHeight("slnv", 200)
	}

	w := dummy[0][0]
	h := dummy[0][1]

	if w != 200 || h != 200 {
		t.Fatal(fmt.Sprintf("Image size was %d %d. Expected 200 200.", w, h))
	}
}

func TestOriginalSize(t *testing.T) {
	dummy := DummyResizerForAspect{}
	cache := AspectPreserver{&dummy}

	cache.GetOriginalSizedImage("slt")

	w := dummy[0][0]
	h := dummy[0][1]

	if w != 1000 || h != 1000 {
		t.Fatal(fmt.Sprintf("Image size was %d %d. Expected 1000 1000.", w, h))
	}
}
