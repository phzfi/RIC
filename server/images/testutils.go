package images

import (
	"errors"
	"fmt"
	"github.com/joonazan/imagick/imagick"
	"path/filepath"
)

func CheckDistortion(blob ImageBlob, reffn string, tol float64, resfn string) (err error) {

	ref := imagick.NewMagickWand()
	defer ref.Destroy()
	err = ref.ReadImage(reffn)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not load ref image: %v Err: %v", reffn, err))
	}

	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	err = mw.ReadImageBlob(blob)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not load blob. Err: %v", err))
	}
	
	trash, d := mw.CompareImages(ref, imagick.METRIC_MEAN_SQUARED_ERROR)
	trash.Destroy()

	err = mw.WriteImage(filepath.FromSlash(resfn))
	if err != nil {
		return errors.New(fmt.Sprintf("Could not write result file: %v. Err: %v", resfn, err))
	}

	if d > tol {
		return errors.New(fmt.Sprintf("Too much distortion. res: %v, ref: %v, tol: %v, dist: %v", resfn, reffn, tol, d))
	}

	return
}

type TestCase struct {
	Testfn, Reffn, Resfn string
}

type FormatTestCase struct {
	TestCase
	Format string
}

type SizeTestCase struct {
	TestCase
	W, H int
}

type TestCaseAll struct {
	TestCase
	Format string
	W, H   int
}

func CheckImage(blob ImageBlob, c TestCase, tol float64, f func(Image) error) (err error) {
	img := NewImage()
	defer img.Destroy()
	img.FromBlob(blob)

	err = f(img)
	if err != nil {
		return
	}

	err = CheckDistortion(blob, c.Reffn, tol, c.Resfn)
	if err != nil {
		return
	}

	return
}

func CheckFormatFunc(c FormatTestCase) func(Image) error {
	return func(img Image) error {
		f := img.GetImageFormat()
		if f != c.Format {
			return errors.New(fmt.Sprintf("Bad image format. Requested %v, Got %v", c.Format, f))
		}
		return nil
	}
}

func CheckSizeFunc(c SizeTestCase) func(Image) error {
	return func(img Image) error {
		w := img.GetWidth()
		h := img.GetHeight()
		if w != c.W || h != c.H {
			return errors.New(fmt.Sprintf("Bad image size. Requested (%v, %v) , Got (%v, %v)", c.W, c.H, w, h))
		}
		return nil
	}
}

func CheckAllFunc(c TestCaseAll) func(Image) error {
	return func(img Image) error {
		w := img.GetWidth()
		h := img.GetHeight()
		if w != c.W || h != c.H {
			return errors.New(fmt.Sprintf("Bad image size. Requested (%v, %v) , Got (%v, %v)", c.W, c.H, w, h))
		}
		f := img.GetImageFormat()
		if f != c.Format {
			return errors.New(fmt.Sprintf("Bad image format. Requested %v, Got %v", c.Format, f))
		}
		return nil
	}
}

func FormatTest(c FormatTestCase, blob ImageBlob, tolerance float64) error {
	return CheckImage(blob, c.TestCase, tolerance, CheckFormatFunc(c))
}

func SizeTest(c SizeTestCase, blob ImageBlob, tolerance float64) error {
	return CheckImage(blob, c.TestCase, tolerance, CheckSizeFunc(c))
}

func TestAll(c TestCaseAll, blob ImageBlob, tolerance float64) error {
	return CheckImage(blob, c.TestCase, tolerance, CheckAllFunc(c))
}
