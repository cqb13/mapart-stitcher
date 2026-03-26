package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/cqb13/mapart-stitcher/internal"
)

func main() {
	args := os.Args[1:]

	if len(args) < 2 {
		fmt.Println("Not enough arguments.")
		help()
		return
	}

	cmd := args[0]
	input := args[1]

	fs := flag.NewFlagSet(cmd, flag.ExitOnError)

	var outputPath string
	var scale int

	fs.StringVar(&outputPath, "o", "map.png", "output image")
	fs.IntVar(&scale, "s", 1, "scale amount")

	fs.Parse(args[2:])

	switch cmd {
	case "stitch":
		err := internal.StitchMapart(input, outputPath, scale)
		if err != nil {
			fmt.Printf("Failed to stitch maps: %s.\n", err)
			return
		}
	case "scale":
		err := internal.ScaleImage(input, outputPath, scale)
		if err != nil {
			fmt.Printf("Failed to scale: %s.\n", err)
			return
		}
	default:
		fmt.Println("Unknown command.")
		help()
		return
	}
}

func help() {
	fmt.Println("Usage:")
	fmt.Println("\tmas stitch map-directory/")
	fmt.Println("\tmas scale cool-map.png")
	fmt.Println("Flags:")
	fmt.Println("\t-o output-path")
	fmt.Println("\t-s scale")
}
