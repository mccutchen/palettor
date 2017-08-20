package palettor

import (
	"image/color"
	"testing"
)

func TestPalette(t *testing.T) {
	colorWeights := map[color.Color]float64{
		black: 0.5,
		white: 0.5,
	}
	iterations := 1
	converged := true
	palette := &Palette{
		colorWeights: colorWeights,
		converged:    converged,
		iterations:   iterations,
	}

	if palette.Count() != len(colorWeights) {
		t.Errorf("wrong number of colors in palette")
	}

	if palette.Converged() != converged {
		t.Errorf("wrong value for converged in palette")
	}

	if palette.Iterations() != iterations {
		t.Errorf("wrong number of iterations in palette")
	}

	if palette.Weight(black) != 0.5 {
		t.Errorf("wrong weight for black")
	}
	if palette.Weight(red) != 0 {
		t.Errorf("wrong weight for unknown color")
	}

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
}
