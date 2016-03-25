package images

// Crop cuts image to requested size. Parameters x and y give offset of cropping.
func (img Image) Crop(w, h, x, y int) error {
	return img.CropImage(uint(w), uint(h), x, y)
}
