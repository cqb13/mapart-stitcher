package internal

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"
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

		return nil
	}

	err = os.Mkdir(outputPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Failed to create output directory: %w", err)
	}

	entries, err := os.ReadDir(inputPath)
	if err != nil {
		return fmt.Errorf("Failed to read input directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".png") {
			continue
		}

		imgPath := filepath.Join(inputPath, entry.Name())
		scaledImgPath := filepath.Join(outputPath, entry.Name())

		err = loadScaleSave(imgPath, scaledImgPath, scale)
		if err != nil {
			fmt.Printf("Failed to scale %s: %s, Skipping...\n", imgPath, err)
			continue
		}
	}

	return nil
}

func loadScaleSave(imgPath string, outputPath string, scale int) error {
	if !strings.HasSuffix(imgPath, ".png") {
		return fmt.Errorf("Input path must lead to a png")
	}

	if !strings.HasSuffix(outputPath, ".png") {
		return fmt.Errorf("Output path must lead to a png")
	}

	fmt.Printf("Loading image %s...\n", imgPath)
	img, err := loadImage(imgPath)
	if err != nil {
		return err
	}
	fmt.Printf("Loaded image, image size is %dx%d\n", img.Bounds().Max.X, img.Bounds().Max.Y)

	fmt.Printf("Scaling image %dx\n", scale)
	scaledImg := scaleImage(img, scale)
	fmt.Printf("Scaled image, image size is %dx%d\n", scaledImg.Bounds().Max.X, scaledImg.Bounds().Max.Y)

	fmt.Printf("Saving image to %s...\n", outputPath)
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("Failed to create output file, %s: %w", outputPath, err)
	}

	err = png.Encode(outputFile, scaledImg)
	if err != nil {
		return fmt.Errorf("Failed to save image to %s: %w", outputPath, err)
	}
	fmt.Printf("Saved scaled image to %s\n", outputPath)

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
