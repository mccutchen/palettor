package palettor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPalette(t *testing.T) {
	iterations := 1
	converged := true
	palette := &Palette{
		converged:  converged,
		iterations: iterations,
	}
	palette.add(black, 0.75)
	palette.add(white, 0.25)

	assert.Equal(t, 2, palette.Count())
	assert.Equal(t, converged, palette.Converged())
	assert.Equal(t, iterations, palette.Iterations())

	assert.Equal(t, 0.75, palette.Weight(black), "wrong weight for black")
	assert.Equal(t, 0.00, palette.Weight(red), "wrong weight for unknown color")

	for _, color := range palette.Colors() {
		assert.Contains(t, palette.entries, asKey(color))
	}

	// ensure entries are sorted by weight
	expectedEntries := []Entry{
		{white, 0.25},
		{black, 0.75},
	}
	assert.Equal(t, expectedEntries, palette.Entries())
}
