package palettor

import (
	"image/color"
	"math/rand"
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
	return &color.RGBA{
		R: uint8(r.Intn(255)),
		G: uint8(r.Intn(255)),
		B: uint8(r.Intn(255)),
		A: uint8(r.Intn(255)),
	}
}

func newColor(r, g, b, a uint8) color.Color {
	return &color.RGBA{
		R: r,
		G: g,
		B: b,
		A: a,
	}
}

func TestDistance(t *testing.T) {
	a := newColor(0, 0, 0, 0)
	b := newColor(255, 255, 255, 0)
	// Distance actually calculated based on uint16 channels, not uint8, due to
	// the RGBA() method on the Color type:
	// http://golang.org/pkg/image/color/#Color
	expected := (0xFFFF * 0xFFFF) + (0xFFFF * 0xFFFF) + (0xFFFF * 0xFFFF)
	if distance(a, b) != expected {
		t.Fatalf("distance should be square of Euclidean distance; %d != %d", distance(a, b), expected)
	}

	a = newColor(0, 0, 0, 0)
	b = newColor(0, 0, 0, 255)
	if distance(a, b) != 0 {
		t.Fatalf("alpha channel is ignored for the purpose of distance")
	}

	c := randomColor()
	if distance(c, c) != 0 {
		t.Fatalf("distance from between identical colors should be 0")
	}
}

func TestNearest(t *testing.T) {
	var haystack []color.Color = []color.Color{black, white, red, green, blue}

	if nearest(black, haystack) != black {
		t.Fatalf("nearest color to self should be self")
	}
	if nearest(darkGray, haystack) != black {
		t.Fatalf("dark gray should be nearest to black")
	}
	if nearest(mostlyRed, haystack) != red {
		t.Fatalf("mostly-red should be nearest to red")
	}
}

func TestFindCentroid(t *testing.T) {
	var cluster []color.Color = []color.Color{black, white, red, mostlyRed}
	centroid := findCentroid(cluster)
	found := false
	for _, c := range cluster {
		if c == centroid {
			found = true
		}
	}
	if !found {
		t.Fatalf("centroid should be a member of the cluster")
	}
}

func TestCluster(t *testing.T) {
	var observations []color.Color = []color.Color{black, white, red}

	k := 4
	clusters, iterations, err := Cluster(k, 100, observations)
	if clusters != nil || iterations != 0 || err == nil {
		t.Fatalf("too few observations should result in an error")
	}

	k = 3
	clusters, _, _ = Cluster(k, 100, observations)
	if len(clusters) != k {
		t.Fatalf("expected %d clusters, got %d; %v", k, len(clusters), clusters)
	}

	k = 2
	observations = []color.Color{black, white}
	clusters, _, _ = Cluster(k, 100, observations)
	weight, _ := clusters[black]
	if weight != 0.5 {
		t.Fatalf("expected weight of black cluster to be 0.5")
	}
	weight, _ = clusters[white]
	if weight != 0.5 {
		t.Fatalf("expected weight of white cluster to be 0.5")
	}
}

func BenchmarkCluster(b *testing.B) {
	// We generally expect an input image to have been thumbnailed down to a
	// manageable size (e.g. 200x200 pixels) before its colors are extracted.
	observationCount := 200 * 200
	observations := make([]color.Color, observationCount)
	for i := 0; i < observationCount; i++ {
		observations[i] = randomColor()
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Cluster(4, 100, observations)
	}
}
