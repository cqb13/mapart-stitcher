# Mapart Stitcher

A program that can automatically combine and scale maparts made out of several maps into a single image.

## Usage

1. Export maps from Minecraft into a supported format, `png` or `jpg`.
2. Run the `stitch` command on the directory to stitch the exported images together.

```sh
mas stitch ./cool-map/

mas scale ./cool-stitched-map.png

FLAGS:
-o output-path.png
-s scale

```

_Note: Exported maps must be 128x128 in order to be stitched together_

## Export Guide

Maps must be exported with the following suffix `(row-col)`

Ex: `cool-map-(0-0).png`

The top left map must always be (0-0).

For example a 2x2 map should be exported like this:
╔═════════╦═════════╗
║map-(0-0)║map-(0-1)║
╠═════════╬═════════╣
║map-(1-0)║map-(1-1)║
╚═════════╩═════════╝
