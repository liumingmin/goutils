package utils

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"os"
)

func ReadImage(filePath string) (image.Image, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	ext := FileExt(filePath)

	if ext == ".jpg" {
		img, err := jpeg.Decode(file)
		if err != nil {
			return nil, err
		}
		file.Close()

		return img, nil
	} else if ext == ".png" {
		img, err := png.Decode(file)
		if err != nil {
			return nil, err
		}
		file.Close()

		return img, nil
	}

	return nil, errors.New(ext)
}

func WriteImage(img image.Image, filePath string) error {
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	return jpeg.Encode(out, img, nil)
}
