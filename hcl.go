package palettor

import (
	"fmt"
	"image/color"

	"github.com/lucasb-eyer/go-colorful"
)

type hcl struct {
	h, c, l float64
}

// RGBA implements color.Color.
func (c hcl) RGBA() (r, g, b, a uint32) {
	// Bodge: squash floating point error to simplify testing with expected
	// output palettes.
	rFloat, gFloat, bFloat := colorful.Hcl(c.h, c.c, c.l).Clamped().RGB255()
	return color.RGBA{rFloat, gFloat, bFloat, 255}.RGBA()
}

// Calculate the square of the Euclidean distance between two colors, ignoring
// the alpha channel.
func (c hcl) distanceSquared(other hcl) float64 {
	dh := c.h - other.h
	dc := c.c - other.c
	dl := c.l - other.l
	return dh*dh + dc*dc + dl*dl
}

func toHCL(col color.Color) (hcl, error) {
	intermediate, ok := colorful.MakeColor(col)
	if !ok {
		return hcl{}, fmt.Errorf("color has alpha channel 0: %+v", col)
	}
	h, c, l := intermediate.Hcl()
	return hcl{h, c, l}, nil
}
