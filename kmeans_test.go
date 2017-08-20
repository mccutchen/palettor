package palettor

import (
	"image"
	"image/color"
	_ "image/jpeg"
	"math/rand"
	"os"
	"testing"
	"time"
)

var (
	r         = rand.New(rand.NewSource(time.Now().UnixNano()))
	black     = newColor(0, 0, 0, 0)
	white     = newColor(255, 255, 255, 0)
	red       = newColor(255, 0, 0, 0)
	green     = newColor(0, 255, 0, 0)
	blue      = newColor(0, 0, 255, 0)
	darkGray  = newColor(1, 1, 1, 0)
	mostlyRed = newColor(200, 0, 0, 0)
)

func randomColor() color.Color {
	return newColor(r.Intn(255), r.Intn(255), r.Intn(255), r.Intn(255))
}

func newColor(r, g, b, a int) color.Color {
	return &color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: uint8(a),
	}
}

func TestDistanceSquared(t *testing.T) {
	a := newColor(0, 0, 0, 0)
	b := newColor(255, 255, 255, 0)
	expected := (0xFFFF * 0xFFFF) + (0xFFFF * 0xFFFF) + (0xFFFF * 0xFFFF)
	if distanceSquared(a, b) != expected {
		t.Errorf("distance should be square of Euclidean distance; %d != %d", distanceSquared(a, b), expected)
	}

	a = newColor(0, 0, 0, 0)
	b = newColor(0, 0, 0, 255)
	if distanceSquared(a, b) != 0 {
		t.Errorf("alpha channel is ignored for the purpose of distance")
	}

	c := randomColor()
	if distanceSquared(c, c) != 0 {
		t.Errorf("distance from between identical colors should be 0")
	}
}

func TestNearest(t *testing.T) {
	var haystack = []color.Color{black, white, red, green, blue}

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
	var cluster = []color.Color{black, white, red, mostlyRed}
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
	var colors = []color.Color{black, white, red}

	k := 4
	palette, err := ClusterColors(k, 100, colors)
	if err == nil {
		t.Errorf("too few colors should result in an error")
	}

	k = 3
	palette, _ = ClusterColors(k, 100, colors)
	if palette.Count() != k {
		t.Errorf("expected %d clusters, got %d", k, palette.Count())
	}

	k = 2
	colors = []color.Color{black, white}
	palette, _ = ClusterColors(k, 100, colors)
	if palette.Weight(black) != 0.5 {
		t.Errorf("expected weight of black cluster to be 0.5")
	}
	if palette.Weight(white) != 0.5 {
		t.Errorf("expected weight of white cluster to be 0.5")
	}

	// If there are not enough unique colors to cluster, it's okay for the size
	// of the extracted palette to be < k
	k = 3
	palette, _ = ClusterColors(k, 100, []color.Color{black, black, black, black, black, white})
	if palette.Count() > 2 {
		t.Errorf("actual palette can be smaller than k")
	}
}

func BenchmarkClusterColors_200x200(b *testing.B) {
	reader, err := os.Open("testdata/resized.jpg")
	if err != nil {
		b.Fatal(err)
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		b.Fatal(err)
	}

	colors := getColors(img)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ClusterColors(4, 100, colors)
	}
}
