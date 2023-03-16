// Package palettor provides a way to extract the color palette from an image
// using k-means clustering.
package palettor

import (
	"fmt"
	"image"

	"github.com/lucasb-eyer/go-colorful"
)

// Extract finds the k most dominant colors in the given image using the
// "standard" k-means clustering algorithm. It returns a Palette, after running
// the algorithm up to maxIterations times.
func Extract(k, maxIterations int, img image.Image) (*Palette, error) {
	imgColors, err := getColors(img)
	if err != nil {
		return nil, fmt.Errorf("error extracting colors from image: %w", err)
	}
	return clusterColors(k, maxIterations, imgColors)
}

func getColors(img image.Image) ([]colorful.Color, error) {
	bounds := img.Bounds()
	pixelCount := (bounds.Max.X - bounds.Min.X) * (bounds.Max.Y - bounds.Min.Y)
	colors := make([]colorful.Color, pixelCount)
	i := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			var ok bool
			if colors[i], ok = colorful.MakeColor(img.At(x, y)); !ok {
				// FIXME: presumably fails for images with transparency.
				return nil, fmt.Errorf("pixel at (%v, %v) has a-channel 0", x, y)
			}
			i++
		}
	}
	return colors, nil
}
