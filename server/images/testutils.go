package images


import (
	"github.com/phzfi/RIC/server/logging"
	"github.com/joonazan/imagick/imagick"
	"path/filepath"
	"errors"
)


func CheckDistortion(blob ImageBlob, reffn string, tol float64, resfn string) (err error) {
	
	ref := NewImage()
	defer ref.Destroy()
	err = ref.FromFile(reffn)
	if err != nil {
		logging.Debug("Could not load image")
		return
	}

	img := NewImage()
	defer img.Destroy()
	err = img.FromBlob(blob)
	if err != nil {
		logging.Debug("Could not load image")
		return
	}

	trash, d := img.CompareImages(ref.MagickWand, imagick.METRIC_MEAN_SQUARED_ERROR)
	trash.Destroy()


	err = img.WriteImage(filepath.FromSlash(resfn))
	if err != nil {
		logging.Debug("Could not write result file:%v", resfn)
		err = nil
	}
	
	if d > tol {
		logging.Debug("Too much distortion. res:%v, ref:%v, tol:%v, dist:%v", resfn, reffn, tol, d)
		err = errors.New("Too much distrotion.")
	}

	return
}
