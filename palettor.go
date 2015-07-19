package palettor

import (
	"image"
	"image/color"
)

// DominantColors uses k-means clustering to find the k most dominant colors in
// the given image.
func DominantColors(k, maxIterations int, img image.Image) (map[color.Color]float64, error) {
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
	clusters, _, err := Cluster(k, maxIterations, colors)
	return clusters, err
}
