package palettor

import (
	"reflect"
	"testing"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/stretchr/testify/assert"
)

func TestPalette(t *testing.T) {
	colorWeights := map[colorful.Color]float64{
		black: 0.75,
		white: 0.25,
	}

	iterations := 1
	converged := true
	palette := &Palette{
		colorWeights: colorWeights,
		converged:    converged,
		iterations:   iterations,
	}

	assert.Equal(t, len(colorWeights), palette.Count())
	assert.Equal(t, converged, palette.Converged())
	assert.Equal(t, iterations, palette.Iterations())

	assert.Equal(t, 0.75, palette.Weight(black), "wrong weight for black")
	assert.Equal(t, 0.00, palette.Weight(red), "wrong weight for unknown color")

	for _, color := range palette.Colors() {
		found := false
		for inputColor := range colorWeights {
			if color == inputColor {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("missing color %v from palette", color)
		}
	}

	// ensure entries are sorted by weight
	expectedEntries := []Entry{
		{white, 0.25},
		{black, 0.75},
	}
	if entries := palette.Entries(); !reflect.DeepEqual(entries, expectedEntries) {
		t.Errorf("expected entries %v, got %v", expectedEntries, entries)
	}
}
