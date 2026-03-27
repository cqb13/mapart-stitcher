package internal

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const singleMapDimension = 128

func StitchMapart(inputDir string, outputPath string, scale int) error {
	fmt.Printf("Loading map images...\n")
	rows, cols, imgGrid, err := loadMapImages(inputDir)
	if err != nil {
		return err
	}

	fmt.Printf("Loaded map images, map image size is %dx%d\n", cols, rows)
	fmt.Printf("Stitching map images...\n")
	img := stitchMapImages(imgGrid, rows, cols)
	fmt.Printf("Stitched map images, image size is %dx%d\n", img.Bounds().Max.X, img.Bounds().Max.Y)

	if scale != 1 {
		fmt.Printf("Scaling image %dx\n", scale)
		img = scaleImage(img, scale)
		fmt.Printf("Scaled image, image size is %dx%d\n", img.Bounds().Max.X, img.Bounds().Max.Y)
	}

	if !strings.HasSuffix(outputPath, ".png") {
		outputPath += ".png"
	}

	fmt.Printf("Saving image to %s...\n", outputPath)
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("Failed to create output file, %s: %w", outputPath, err)
	}

	err = png.Encode(outputFile, img)
	if err != nil {
		return fmt.Errorf("Failed to save image to %s: %w", outputPath, err)
	}
	fmt.Printf("Saved stitched image to %s\n", outputPath)

	return nil
}

func extractCoordinate(s string) (int, int, bool) {
	re := regexp.MustCompile(`-(\d+)-(\d+)\.png$`)
	matches := re.FindStringSubmatch(s)

	if matches == nil {
		return 0, 0, false
	}

	row, _ := strconv.Atoi(matches[1])
	col, _ := strconv.Atoi(matches[2])

	return row, col, true
}

func stitchMapImages(imgGrid [][]image.Image, rows int, cols int) image.Image {
	img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{cols * singleMapDimension, rows * singleMapDimension}})

	xOffset := 0
	yOffset := 0

	for _, row := range imgGrid {
		xOffset = 0
		for _, mapImage := range row {
			// copy single map over
			for y := range singleMapDimension {
				for x := range singleMapDimension {
					if mapImage == nil {
						img.Set(xOffset+x, yOffset+y, color.Transparent)
						continue
					}

					img.Set(xOffset+x, yOffset+y, mapImage.At(x, y))
				}
			}
			xOffset += singleMapDimension
		}
		yOffset += singleMapDimension
	}

	return img
}

func loadMapImages(inputDir string) (int, int, [][]image.Image, error) {
	entries, err := os.ReadDir(inputDir)
	if err != nil {
		return 0, 0, nil, fmt.Errorf("Failed to read input directory: %w", err)
	}

	imgGrid := [][]image.Image{}

	rows := 0
	cols := 0

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".png") {
			continue
		}

		row, col, matches := extractCoordinate(entry.Name())
		if !matches {
			continue
		}

		// load and parse file
		imgPath := filepath.Join(inputDir, entry.Name())

		file, err := os.Open(imgPath)
		if err != nil {
			fmt.Printf("Failed to open %s: %s, Skipping...\n", imgPath, err)
			continue
		}

		img, err := png.Decode(file)
		if err != nil {
			fmt.Printf("Failed to decode png %s: %s, Skipping...", imgPath, err)
			continue
		}

		if img.Bounds().Max.X != singleMapDimension || img.Bounds().Max.Y != singleMapDimension {
			fmt.Printf("Dimension of %s (%dx%d) don't match expected %dx%d, Skipping...", imgPath, img.Bounds().Max.X, img.Bounds().Max.Y, singleMapDimension, singleMapDimension)
			continue
		}

		// Ensure enough rows
		for len(imgGrid) <= row {
			imgGrid = append(imgGrid, make([]image.Image, cols))
			rows++
		}

		// Ensure enough columns
		if col >= cols {
			newCols := col + 1

			for i := range imgGrid {
				newRow := make([]image.Image, newCols)
				copy(newRow, imgGrid[i])
				imgGrid[i] = newRow
			}

			cols = newCols
		}

		imgGrid[row][col] = img

		fmt.Printf("%s\n", entry.Name())
	}

	if rows == 0 || cols == 0 {
		return 0, 0, nil, fmt.Errorf("No valid map images found")
	}

	return rows, cols, imgGrid, nil
}
