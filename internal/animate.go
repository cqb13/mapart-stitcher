package internal

import (
	"fmt"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/png"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func AnimateMapart(inputDir string, outputPath string, scale int, gifFrameDelay int, log bool) error {
	if log {
		fmt.Printf("Loading map images...\n")
	}

	maps, err := loadMapImages(inputDir, log)
	if err != nil {
		return err
	}

	if log {
		fmt.Printf("Loaded map images\n")
		fmt.Printf("Creating gif\n")
	}

	anim := gif.GIF{}

	for _, img := range maps {
		if scale != 1 {
			if log {
				fmt.Printf("Scaling image %dx\n", scale)
			}
			img = scaleImage(img, scale)
			if log {
				fmt.Printf("Scaled image, image size is %dx%d\n", img.Bounds().Max.X, img.Bounds().Max.Y)
			}
		}

		bounds := img.Bounds()
		palettedImg := image.NewPaletted(bounds, palette.Plan9)
		draw.FloydSteinberg.Draw(palettedImg, bounds, img, image.Point{0, 0})

		anim.Image = append(anim.Image, palettedImg)
		anim.Delay = append(anim.Delay, gifFrameDelay)
	}

	if log {
		fmt.Printf("Created gif\n")
	}

	if !strings.HasSuffix(outputPath, ".gif") {
		outputPath += ".gif"
		if log {
			fmt.Println("Automatically adding .gif file extension to output path")
		}
	}

	if log {
		fmt.Printf("Saving gif to %s...\n", outputPath)
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("Failed to create output file, %s: %w", outputPath, err)
	}
	defer outputFile.Close()

	err = gif.EncodeAll(outputFile, &anim)
	if err != nil {
		return fmt.Errorf("Failed to save gif to %s: %w", outputPath, err)
	}
	fmt.Printf("Saved gif to %s\n", outputPath)

	return nil
}

func loadMapImages(inputDir string, log bool) ([]image.Image, error) {
	_ = log

	entries, err := os.ReadDir(inputDir)
	if err != nil {
		return nil, fmt.Errorf("Failed to read input directory: %w", err)
	}

	sort.Slice(entries, func(i, j int) bool {
		getNum := func(name string) int {
			parts := strings.Split(name, "-")
			n, _ := strconv.Atoi(parts[len(parts)-1])
			return n
		}
		return getNum(entries[i].Name()) < getNum(entries[j].Name())
	})

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".png") {
			continue
		}
	}

	maps := []image.Image{}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".png") {
			continue
		}

		imgPath := filepath.Join(inputDir, entry.Name())

		file, err := os.Open(imgPath)
		if err != nil {
			fmt.Printf("Failed to open %s: %s, Skipping...\n", imgPath, err)
			continue
		}
		defer file.Close()

		img, err := png.Decode(file)
		if err != nil {
			fmt.Printf("Failed to decode png %s: %s, Skipping...", imgPath, err)
			continue
		}

		maps = append(maps, img)
	}

	return maps, nil
}
