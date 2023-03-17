// Package palettor provides a way to extract the color palette from an image
// using k-means clustering.
package palettor

import (
	"fmt"
	"image"
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

func getColors(img image.Image) ([]hcl, error) {
	bounds := img.Bounds()
	pixelCount := (bounds.Max.X - bounds.Min.X) * (bounds.Max.Y - bounds.Min.Y)
	colors := make([]hcl, pixelCount)
	i := 0
	var err error
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if colors[i], err = toHCL(img.At(x, y)); err != nil {
				return nil, fmt.Errorf("error translating pixel at (%v, %v): %w", x, y, err)
			}
			i++
		}
	}
	return colors, nil
}
