package main

import (
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"

	"github.com/nfnt/resize"
	"kmeans"
)

func pixelCount(img image.Image) int {
	bounds := img.Bounds()
	return (bounds.Max.X - bounds.Min.X) * (bounds.Max.Y - bounds.Min.Y)
}

func main() {
	k := 3

	originalImg, _, err := image.Decode(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	img := resize.Thumbnail(200, 200, originalImg, resize.Lanczos3)
	colors := make([]color.Color, pixelCount(img))
	bounds := img.Bounds()
	i := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			colors[i] = img.At(x, y)
			i++
		}
	}

	log.Printf("color count:  %v", len(colors))
	log.Printf("first color: %#v %T", colors[0], colors[0])

	clusters, _ := kmeans.Cluster(k, colors, 100)
	log.Printf("clusters: %v", clusters)

	// png.Encode(os.Stdout, img)
}
