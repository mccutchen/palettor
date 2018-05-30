package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io"
	"log"
	"math"
	"os"

	"github.com/mccutchen/palettor"
	"github.com/nfnt/resize"
)

func main() {
	var (
		k          = flag.Int("k", 3, "Palette size")
		maxIters   = flag.Int("max", 500, "Maximum k-means iterations")
		jsonOutput = flag.Bool("json", false, "Output color palette in JSON format")
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] [INPUT]\n\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	var (
		input io.Reader
		err   error
	)
	inputPath := flag.Arg(0)
	if inputPath == "" || inputPath == "-" {
		input = os.Stdin
	} else {
		input, err = os.Open(inputPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	img, err := loadImage(input)
	if err != nil {
		log.Fatalf("Error decoding image: %s", err)
	}

	// Get the image down to a more manageable size
	thumb := resize.Thumbnail(200, 200, img, resize.NearestNeighbor)

	palette, err := palettor.Extract(*k, *maxIters, thumb)
	if err != nil {
		log.Fatalf("Error extracing color palette: %s", err)
	}

	if *jsonOutput {
		if err := json.NewEncoder(os.Stdout).Encode(palette.Entries()); err != nil {
			log.Fatalf("Error encoding JSON: %s", err)
		}
		return
	}

	drawPalette(os.Stdout, img, palette)
}

func loadImage(src io.Reader) (image.Image, error) {
	img, _, err := image.Decode(src)
	if err != nil {
		return nil, err
	}

	// Ensure we're working with RGBA data, which is necessary to a) have more
	// immediately useful JSON output and b) allow us to draw a palette back
	// onto the source image.
	//
	// In particular, JPEGs decode to *image.YCbCr, which must be converted to
	// *image.RGBA before we can draw our palette onto it.
	//
	// https://stackoverflow.com/a/47539710/151221
	if _, ok := img.(*image.RGBA); !ok {
		img2 := image.NewRGBA(img.Bounds())
		draw.Draw(img2, img2.Bounds(), img, image.ZP, draw.Src)
		img = img2
	}

	return img, nil
}

// Draw a palette over the bottom 10% of an image
func drawPalette(dst io.Writer, img image.Image, palette *palettor.Palette) {
	drawImg := img.(draw.Image)

	imgWidth := img.Bounds().Dx()
	imgHeight := img.Bounds().Dy()

	paletteHeight := int(math.Ceil(float64(imgHeight) * 0.1))
	yOffset := imgHeight - paletteHeight
	xOffset := 0

	for _, entry := range palette.Entries() {
		colorWidth := int(math.Ceil(float64(imgWidth) * entry.Weight))
		bounds := image.Rect(xOffset, yOffset, xOffset+colorWidth, yOffset+paletteHeight)
		draw.Draw(drawImg, bounds, &image.Uniform{entry.Color}, image.ZP, draw.Src)
		xOffset += colorWidth
	}

	png.Encode(dst, drawImg)
}
