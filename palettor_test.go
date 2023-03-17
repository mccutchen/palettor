package palettor

import (
	"bytes"
	"encoding/base64"
	"image/png"
	"testing"

	"github.com/stretchr/testify/assert"
)

// base64-encoded 4x4 png, w/ black, white, red, & blue pixels
var testImageData = []byte("iVBORw0KGgoAAAANSUhEUgAAAAIAAAACCAIAAAD91JpzAAAAE0lEQVQIHWMAgv///zP8ZwCC/wAh7AT8vKm73AAAAABJRU5ErkJggg==")

func TestExtract(t *testing.T) {
	decoder := base64.NewDecoder(base64.StdEncoding, bytes.NewReader(testImageData))
	img, err := png.Decode(decoder)
	if err != nil {
		t.Fatalf("invalid test image: %s", err)
	}

	_, err = Extract(5, 100, img)
	assert.Error(t, err, "should error when k is too large")

	palette, _ := Extract(4, 100, img)
	assert.Equal(t, 4, palette.Count())
}
