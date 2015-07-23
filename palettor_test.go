package palettor

import (
	"bytes"
	"encoding/base64"
	"image/png"
	"testing"
)

// base64-encoded 4x4 png, w/ black, white, red, & blue pixels
var testImageData = []byte("iVBORw0KGgoAAAANSUhEUgAAAAIAAAACCAIAAAD91JpzAAAAE0lEQVQIHWMAgv///zP8ZwCC/wAh7AT8vKm73AAAAABJRU5ErkJggg==")

func TestFindPalette(t *testing.T) {
	decoder := base64.NewDecoder(base64.StdEncoding, bytes.NewReader(testImageData))
	img, err := png.Decode(decoder)
	if err != nil {
		t.Fatalf("invalid test image: %s", err)
	}

	_, err = FindPalette(5, 100, img)
	if err == nil {
		t.Errorf("k too large, expected an error")
	}

	palette, _ := FindPalette(4, 100, img)
	if palette.Count() != 4 {
		t.Errorf("expected 4 colors, got %d", palette.Count())
	}
}
