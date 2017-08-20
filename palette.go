package palettor

import (
	"image/color"
)

// A Palette represents the dominant colors extracted from an image, as a
// mapping from color to the weight of that color's cluster. The weight can be
// used as an approximation for that color's relative dominance in an image.
type Palette struct {
	colorWeights map[color.Color]float64
	converged    bool
	iterations   int
}

// Colors returns a slice of the colors that comprise a Palette.
func (p *Palette) Colors() []color.Color {
	var colors []color.Color
	for color := range p.colorWeights {
		colors = append(colors, color)
	}
	return colors
}

// Converged returns a bool indicating whether a stable set of dominant
// colors was found before the maximum number of iterations was reached.
func (p *Palette) Converged() bool {
	return p.converged
}

// Count returns the number of colors in a Palette.
func (p *Palette) Count() int {
	return len(p.colorWeights)
}

// Iterations returns the number of iterations required to extract the colors
// of a Palette.
func (p *Palette) Iterations() int {
	return p.iterations
}

// Weight returns the weight of a color in a Palette as a float in the range
// [0, 1], or 0 if a given color is not found.
func (p *Palette) Weight(c color.Color) float64 {
	weight, _ := p.colorWeights[c]
	return weight
}
