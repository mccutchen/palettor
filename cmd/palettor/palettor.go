package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"math"
	"os"

	"github.com/mccutchen/palettor"
	"github.com/nfnt/resize"
	"github.com/pkg/profile"
)

func main() {
	var (
		k          = flag.Int("k", 3, "Palette size")
		maxIters   = flag.Int("max", 500, "Maximum k-means iterations")
		jsonOutput = flag.Bool("json", false, "Output color palette in JSON format")
		noResize   = flag.Bool("no-resize", false, "Do not resize input image before processing")
		doProfile  = flag.Bool("profile", false, "Capture profile")
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

	img, format, err := loadImage(input)
	if err != nil {
		log.Fatalf("Error decoding image: %s", err)
	}

	// Get the image down to a more manageable size
	if !*noResize {
		img = resize.Thumbnail(200, 200, img, resize.NearestNeighbor)
	}

	// Only start profiling after the image is loaded
	if *doProfile {
		defer profile.Start().Stop()
	}

	palette, err := palettor.Extract(*k, *maxIters, img)
	if err != nil {
		log.Fatalf("Error extracting color palette: %s", err)
	}

	if *jsonOutput {
		if err := json.NewEncoder(os.Stdout).Encode(palette.Entries()); err != nil {
			log.Fatalf("Error encoding JSON: %s", err)
		}
		return
	}

	if err := drawPalette(os.Stdout, img, palette, format); err != nil {
		log.Fatalf("Error encoding palette: %s", err)
	}
}

func loadImage(src io.Reader) (image.Image, string, error) {
	img, format, err := image.Decode(src)
	if err != nil {
		return nil, "", err
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
		draw.Draw(img2, img2.Bounds(), img, image.Point{}, draw.Src)
		img = img2
	}

	return img, format, nil
}

// Draw a palette over the bottom 10% of an image
func drawPalette(dst io.Writer, img image.Image, palette *palettor.Palette, format string) error {
	drawImg := img.(draw.Image)

	imgWidth := img.Bounds().Dx()
	imgHeight := img.Bounds().Dy()

	paletteHeight := int(math.Ceil(float64(imgHeight) * 0.1))
	yOffset := imgHeight - paletteHeight
	xOffset := 0

	for _, entry := range palette.Entries() {
		colorWidth := int(math.Ceil(float64(imgWidth) * entry.Weight))
		bounds := image.Rect(xOffset, yOffset, xOffset+colorWidth, yOffset+paletteHeight)
		draw.Draw(drawImg, bounds, &image.Uniform{entry.Color}, image.Point{}, draw.Src)
		xOffset += colorWidth
	}

	switch format {
	case "jpeg":
		return jpeg.Encode(dst, drawImg, nil)
	case "gif":
		return gif.Encode(dst, drawImg, nil)
	default:
		return png.Encode(dst, drawImg)
	}
}
