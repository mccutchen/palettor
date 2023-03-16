package palettor

import (
	"image"
	"image/color"
	_ "image/jpeg"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/stretchr/testify/assert"
)

var (
	r         = rand.New(rand.NewSource(time.Now().UnixNano()))
	black     = newColor(0, 0, 0, 255)
	white     = newColor(255, 255, 255, 255)
	red       = newColor(255, 0, 0, 255)
	green     = newColor(0, 255, 0, 255)
	blue      = newColor(0, 0, 255, 255)
	darkGray  = newColor(1, 1, 1, 255)
	mostlyRed = newColor(200, 0, 0, 255)
)

func randomColor() colorful.Color {
	return colorful.Hcl(r.Float64()*360, r.Float64(), r.Float64())
}

func newColor(r, g, b, a int) colorful.Color {
	// Bodge: keep constants.
	color, ok := colorful.MakeColor(&color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: uint8(a),
	})
	if !ok {
		panic("Color fixtures must have nonzero A-channel")
	}

	return color
}

func TestDistanceSquared(t *testing.T) {
	a := newColor(0, 0, 0, 255)
	b := newColor(255, 255, 255, 255)

	assert.InDelta(t, 1, distanceSquared(a, b), .0001, "distance should be square of Euclidean distance")

	a = newColor(0, 0, 0, 1)
	b = newColor(0, 0, 0, 255)
	assert.Equal(t, 0.00, distanceSquared(a, b), "alpha channel should be ignored for the purpose of distance")

	c := randomColor()
	assert.Equal(t, 0.00, distanceSquared(c, c), "distance from between identical colors should be 0")
}

func TestNearest(t *testing.T) {
	var haystack = []colorful.Color{black, white, red, green, blue}

	if nearest(black, haystack) != black {
		t.Errorf("nearest color to self should be self")
	}
	if nearest(darkGray, haystack) != black {
		t.Errorf("dark gray should be nearest to black")
	}
	if nearest(mostlyRed, haystack) != red {
		t.Errorf("mostly-red should be nearest to red")
	}
}

func TestFindCentroid(t *testing.T) {
	var cluster = []colorful.Color{black, white, red, mostlyRed}
	centroid := findCentroid(cluster)
	found := false
	for _, c := range cluster {
		if c == centroid {
			found = true
		}
	}
	if !found {
		t.Errorf("centroid should be a member of the cluster")
	}
}

func TestCluster(t *testing.T) {
	var colors = []colorful.Color{black, white, red}

	k := 4
	_, err := clusterColors(k, 100, colors)
	if err == nil {
		t.Errorf("too few colors should result in an error")
	}

	k = 3
	palette, err := clusterColors(k, 100, colors)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if palette.Count() != k {
		t.Errorf("expected %d clusters, got %d", k, palette.Count())
	}

	k = 2
	colors = []colorful.Color{black, white}
	palette, _ = clusterColors(k, 100, colors)
	if palette.Weight(black) != 0.5 {
		t.Errorf("expected weight of black cluster to be 0.5")
	}
	if palette.Weight(white) != 0.5 {
		t.Errorf("expected weight of white cluster to be 0.5")
	}

	// If there are not enough unique colors to cluster, it's okay for the size
	// of the extracted palette to be < k
	k = 3
	palette, _ = clusterColors(k, 100, []colorful.Color{black, black, black, black, black, white})
	if palette.Count() > 2 {
		t.Errorf("actual palette can be smaller than k")
	}
}

func BenchmarkClusterColors200x200(b *testing.B) {
	reader, err := os.Open("testdata/resized.jpg")
	if err != nil {
		b.Fatal(err)
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		b.Fatal(err)
	}

	colors, err := getColors(img)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := clusterColors(4, 100, colors); err != nil {
			b.Error(err)
		}
	}
}
