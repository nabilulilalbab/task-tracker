package utils

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"strings"

	"github.com/nfnt/resize"
)

func SaveResizedImage(file io.Reader, ext, savePath string, maxWidth uint) error {
	// Decode image dari io.Reader
	var img image.Image
	var err error

	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(file)
	case ".png":
		img, err = png.Decode(file)
	default:
		return errors.New("unsupported image format")
	}
	if err != nil {
		return err
	}

	// Resize image
	m := resize.Resize(maxWidth, 0, img, resize.Lanczos3)

	// Simpan ke disk
	out, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer out.Close()

	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg":
		err = jpeg.Encode(out, m, &jpeg.Options{Quality: 85})
	case ".png":
		err = png.Encode(out, m)
	}

	return err
}
