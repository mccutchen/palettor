package palettor

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDistanceSquared(t *testing.T) {
	a := forceHCL(color.RGBA{0, 0, 0, 255})
	b := forceHCL(color.RGBA{255, 255, 255, 255})
	assert.InDelta(t, 1, a.distanceSquared(b), .0001, "distance should be square of Euclidean distance")

	a = forceHCL(color.RGBA{0, 0, 0, 1})
	b = forceHCL(color.RGBA{0, 0, 0, 255})
	assert.Equal(t, 0.00, a.distanceSquared(b), "alpha channel should be ignored for the purpose of distance")

	c := forceHCL(randomColor())
	assert.Equal(t, 0.00, c.distanceSquared(c), "distance from between identical colors should be 0")
}

func TestColor(t *testing.T) {
	assert.Implements(t, (*color.Color)(nil), new(hcl))

	input := color.RGBA{123, 123, 123, 255}
	inputR, inputG, inputB, inputA := input.RGBA()

	c, err := toHCL(input)
	assert.NoError(t, err)

	r, g, b, a := c.RGBA()
	assert.Equal(t, inputR, r)
	assert.Equal(t, inputG, g)
	assert.Equal(t, inputB, b)
	assert.Equal(t, inputA, a)
}
