package internal

import (
	"image"
)

func ScaleImage(inputPath string, outputPath string, scale int) error {
	return nil
}

func scaleImage(img image.Image, scale int) image.Image {
	rows := img.Bounds().Max.X * scale
	cols := img.Bounds().Max.Y * scale

	resized := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{cols, rows}})

	for x := range rows {
		for y := range cols {
			xp := x / scale
			yp := y / scale
			resized.Set(x, y, img.At(xp, yp))
		}
	}
	return resized
}
