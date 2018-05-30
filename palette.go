package palettor

import (
	"image/color"
	"sort"
)

// A Palette represents the dominant colors extracted from an image, as a
// mapping from color to the weight of that color's cluster. The weight can be
// used as an approximation for that color's relative dominance in an image.
type Palette struct {
	colorWeights map[color.Color]float64
	converged    bool
	iterations   int
}

// Entry is a color and its weight in a Palette
type Entry struct {
	Color  color.Color `json:"color"`
	Weight float64     `json:"weight"`
}

// Entries returns a slice of Entry structs, sorted by weight
func (p *Palette) Entries() []Entry {
	entries := make([]Entry, p.Count())
	i := 0
	for color, weight := range p.colorWeights {
		entries[i] = Entry{color, weight}
		i++
	}
	sort.Sort(byWeight(entries))
	return entries
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
	return p.colorWeights[c]
}

// implement sort.Interface
type byWeight []Entry

func (a byWeight) Len() int           { return len(a) }
func (a byWeight) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byWeight) Less(i, j int) bool { return a[i].Weight < a[j].Weight }
