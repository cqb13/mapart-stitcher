package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Not enough arguments.")
		help()
		return
	}

	cmd := args[0]

	fs := flag.NewFlagSet(cmd, flag.ExitOnError)

	var outputFlag string
	var scaleFlag int

	fs.StringVar(&outputFlag, "o", "", "output image")
	fs.IntVar(&scaleFlag, "s", 1, "scale amount")

	fs.Parse(args[1:])

	switch cmd {
	case "stitch":
		fmt.Println("stitch stuff here")
	case "scale":
		fmt.Println("scale stuff here")
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
