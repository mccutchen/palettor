package palettor

import (
	"fmt"
	"image/color"
	"math"

	"github.com/lucasb-eyer/go-colorful"
)

func toHCL(col color.Color) (hcl, error) {
	intermediate, ok := colorful.MakeColor(col)
	if !ok {
		return hcl{}, fmt.Errorf("color has alpha channel 0: %+v", col)
	}
	h, c, l := intermediate.Hcl()
	return hcl{h, c, l}, nil
}

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
//
// Note: we may want to weight these to get greater C/L variance.
func (c hcl) distanceSquared(other hcl) float64 {
	dh := c.hueDistance(other)
	dc := c.c - other.c
	dl := c.l - other.l
	return dh*dh + dc*dc + dl*dl
}

// hueDistance calculates the angular distance between c.h and other.h. The
// arithmetic distance can be misleading: 0 and 360 have an arithmetic delta of
// 360, but they coincide (zero hue distance).
func (c hcl) hueDistance(other hcl) float64 {
	delta := math.Mod(other.h-c.h, 360)
	// Pick the shorter angular distance: 'clockwise' or 'counterclockwise'
	// around the unit circle.
	return math.Min(
		math.Abs(delta),
		math.Abs(360-delta),
	)
}

func mean(colors []hcl) hcl {
	return hcl{
		h: meanHue(colors),
		c: arithmeticMean(colors, func(c hcl) float64 { return c.c }),
		l: arithmeticMean(colors, func(c hcl) float64 { return c.l }),
	}
}

// meanHue implements a circular mean: averaging H-values can lead to visually
// improper centroids. See https://en.wikipedia.org/wiki/Circular_mean#Example
func meanHue(colors []hcl) float64 {
	meanSin := arithmeticMean(colors, func(c hcl) float64 {
		return math.Sin(radians(c.h))
	})
	meanCos := arithmeticMean(colors, func(c hcl) float64 {
		return math.Cos(radians(c.h))
	})
	return degrees(math.Atan(meanSin / meanCos))
}

func radians(degrees float64) float64 {
	return degrees * (math.Pi / 180)
}

func degrees(radians float64) float64 {
	return math.Mod(radians*(180/math.Pi), 360)
}

func arithmeticMean(colors []hcl, accessor func(hcl) float64) float64 {
	var sum float64
	for _, c := range colors {
		sum += accessor(c)
	}
	return sum / float64(len(colors))
}
