package internal

import (
	"fmt"
	"image"
	"image/png"
	"os"
)

func ScaleImage(inputPath string, outputPath string, scale int) error {
	fileInfo, err := os.Stat(inputPath)
	if err != nil {
		return err
	}

	massScale := fileInfo.IsDir()

	if !massScale {
		err := loadScaleSave(inputPath, outputPath, scale)
		if err != nil {
			return err
		}
	}

	return nil
}

func loadScaleSave(imgPath string, outputPath string, scale int) error {
	img, err := loadImage(imgPath)
	if err != nil {
		return err
	}

	scaledImg := scaleImage(img, scale)

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("Failed to create output file, %s: %w", outputPath, err)
	}

	err = png.Encode(outputFile, scaledImg)
	if err != nil {
		return fmt.Errorf("Failed to save image to %s: %w", outputPath, err)
	}

	return nil
}

func loadImage(imgPath string) (image.Image, error) {
	file, err := os.Open(imgPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to open %s: %s", imgPath, err)
	}

	img, err := png.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("Failed to decode png %s: %s", imgPath, err)
	}

	return img, nil
}

func scaleImage(img image.Image, scale int) image.Image {
	rows := img.Bounds().Max.X * scale
	cols := img.Bounds().Max.Y * scale

	scaledImg := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{cols, rows}})

	for x := range rows {
		for y := range cols {
			xp := x / scale
			yp := y / scale
			scaledImg.Set(x, y, img.At(xp, yp))
		}
	}
	return scaledImg
}
