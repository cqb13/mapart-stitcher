package internal

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const singleMapDimension = 128

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

func StitchMapart(inputDir string, outputPath string, scale int) error {
	fmt.Printf("Loading maps...\n")
	rows, cols, imgGrid, err := loadMaps(inputDir)
	if err != nil {
		return err
	}

	fmt.Printf("Loaded maps, final map size is %dx%d\n", cols, rows)

	_ = imgGrid

	return nil
}

func loadMaps(inputDir string) (int, int, [][]image.Image, error) {
	entries, err := os.ReadDir(inputDir)
	if err != nil {
		return 0, 0, nil, err
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
		return 0, 0, nil, fmt.Errorf("No valid map images founds.")
	}

	return rows, cols, imgGrid, nil
}
