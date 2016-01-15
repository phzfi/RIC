package cache

type DenyUpscale struct {
	Resizer
}

func (d DenyUpscale) GetImage(filename string, xsize uint, ysize uint)(images.ImageBlob, error) {
	x, y, err := ImageSize(filename)
	if x < xsize {
		xsize = x
	}
	if y < ysize {
		ysize = y
	}
	return d.Resizer.GetImage(filename, xsize, ysize)
}
