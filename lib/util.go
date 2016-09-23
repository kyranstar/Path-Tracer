package lib

import (
	"image"
	"image/png"
	"os"
)

const EPS = 1e-9

func WritePng(filename string, img image.Image) (err error) {
	file, err := os.Create(filename)
	if err != nil {
		return
	}

	defer func() {
		if cerr := file.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	err = png.Encode(file, img)
	return
}
