package palettor

import (
	"encoding/json"
	"fmt"
	"image/color"
	"sort"

	"github.com/lucasb-eyer/go-colorful"
)

type rgbaKey [4]uint32

func asKey(c color.Color) rgbaKey {
	r, g, b, a := c.RGBA()
	return rgbaKey{r, g, b, a}
}

// A Palette represents the dominant colors extracted from an image, as a
// mapping from color to the weight of that color's cluster. The weight can be
// used as an approximation for that color's relative dominance in an image.
type Palette struct {
	entries    map[rgbaKey]Entry
	converged  bool
	iterations int
}

func (p *Palette) add(c color.Color, weight float64) {
	if p.entries == nil {
		p.entries = make(map[rgbaKey]Entry)
	}
	p.entries[asKey(c)] = Entry{Color: c, Weight: weight}
}

// Entry is a color and its weight in a Palette
type Entry struct {
	Color  color.Color `json:"color"`
	Weight float64     `json:"weight"`
}

// MarshalJSON turns e into a more usefully readable JSON structure, with a hex
// value and RGB values in the 0-255 interval.
func (e Entry) MarshalJSON() ([]byte, error) {
	type Alias Entry
	// Bodge: convert to colorful.Color for easier representation.
	c, ok := colorful.MakeColor(e.Color)
	if !ok {
		return nil, fmt.Errorf("colorful can't handle color: %+v", e.Color)
	}
	r, g, b := c.RGB255()
	return json.Marshal(&struct {
		Color color.RGBA `json:"color"`
		Hex   string     `json:"hex"`
		Alias
	}{
		Color: color.RGBA{r, g, b, 255},
		Hex:   c.Hex(),
		Alias: (Alias)(e),
	})
}

// Entries returns a slice of Entry structs, sorted by weight
func (p *Palette) Entries() []Entry {
	entries := make([]Entry, p.Count())
	i := 0
	for _, entry := range p.entries {
		entries[i] = entry
		i++
	}
	sort.Sort(byWeight(entries))
	return entries
}

// Colors returns a slice of the colors that comprise a Palette.
func (p *Palette) Colors() []color.Color {
	var colors []color.Color
	for _, entry := range p.entries {
		colors = append(colors, entry.Color)
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
	return len(p.entries)
}

// Iterations returns the number of iterations required to extract the colors
// of a Palette.
func (p *Palette) Iterations() int {
	return p.iterations
}

// Weight returns the weight of a color in a Palette as a float in the range
// [0, 1], or 0 if a given color is not found.
func (p *Palette) Weight(c color.Color) float64 {
	entry, ok := p.entries[asKey(c)]
	if !ok {
		return 0
	}
	return entry.Weight
}

// implement sort.Interface
type byWeight []Entry

func (a byWeight) Len() int           { return len(a) }
func (a byWeight) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byWeight) Less(i, j int) bool { return a[i].Weight < a[j].Weight }
