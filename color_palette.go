package palettor

import (
	"image/color"
)

// A ColorPalette represents the dominant colors extracted from an image.
type ColorPalette struct {
	colorWeights map[color.Color]float64
	iterations   int
}

// Colors returns the colors in a color palette.
func (p *ColorPalette) Colors() []color.Color {
	colors := make([]color.Color, len(p.colorWeights))
	i := 0
	for color := range p.colorWeights {
		colors[i] = color
		i++
	}
	return colors
}

// Count returns the number of colors in a color palette.
func (p *ColorPalette) Count() int {
	return len(p.colorWeights)
}

// Iterations returns the number of iterations required to extract the colors
// of a palette.
func (p *ColorPalette) Iterations() int {
	return p.iterations
}

// Weight returns the weight of a color in a palette as a float in the range
// [0, 1], or 0 if a given color is not in a palette.
func (p *ColorPalette) Weight(c color.Color) float64 {
	weight, _ := p.colorWeights[c]
	return weight
}
