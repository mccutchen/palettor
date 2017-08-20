// Package palettor provides a way to extract the color palette from an image
// using k-means clustering.
package palettor

import (
	"image"
	"image/color"
)

// Extract finds the k most dominant colors in the given image. It returns a
// Palette, after running the k-means algorithm up to maxIteration times. See
// ClusterColors for the actual k-means implementation.
func Extract(k, maxIterations int, img image.Image) (*Palette, error) {
	return ClusterColors(k, maxIterations, getColors(img))
}

func getColors(img image.Image) []color.Color {
	bounds := img.Bounds()
	pixelCount := (bounds.Max.X - bounds.Min.X) * (bounds.Max.Y - bounds.Min.Y)
	colors := make([]color.Color, pixelCount)
	i := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			colors[i] = img.At(x, y)
			i++
		}
	}
	return colors
}
