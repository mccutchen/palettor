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
	r = rand.New(rand.NewSource(time.Now().UnixNano()))

	black     = forceHCL(color.RGBA{0, 0, 0, 255})
	white     = forceHCL(color.RGBA{255, 255, 255, 255})
	red       = forceHCL(color.RGBA{255, 0, 0, 255})
	green     = forceHCL(color.RGBA{0, 255, 0, 255})
	blue      = forceHCL(color.RGBA{0, 0, 255, 255})
	darkGrey  = forceHCL(color.RGBA{1, 1, 1, 255})
	mostlyRed = forceHCL(color.RGBA{200, 0, 0, 255})
)

func randomColor() colorful.Color {
	return colorful.Hcl(r.Float64()*360, r.Float64(), r.Float64())
}

func TestNearest(t *testing.T) {
	var haystack = []hcl{black, white, red, green, blue}

	assert.Equal(t, black, nearest(black, haystack), "nearest color to self should be self")
	assert.Equal(t, black, nearest(darkGrey, haystack), "dark gray should be nearest to black")
	assert.Equal(t, red, nearest(mostlyRed, haystack), "mostly-red should be nearest to red")
}

func TestFindCentroid(t *testing.T) {
	var cluster = []hcl{black, white, red, mostlyRed}
	centroid := findCentroid(cluster)

	assert.Contains(t, cluster, centroid, "centroid should be a member of the cluster")
}

func TestCluster(t *testing.T) {
	var colors = []hcl{black, white, red}

	k := 4
	_, err := clusterColors(k, 100, colors)
	assert.Error(t, err, "too few colors should result in an error")

	k = 3
	palette, err := clusterColors(k, 100, colors)
	assert.NoError(t, err)
	assert.Equal(t, k, palette.Count(), "got unexpected number of clusters")

	k = 2
	colors = []hcl{black, white}
	palette, _ = clusterColors(k, 100, colors)
	assert.Equal(t, 0.5, palette.Weight(black), "expected weight of black cluster to be 0.5")
	assert.Equal(t, 0.5, palette.Weight(white), "expected weight of white cluster to be 0.5")

	// If there are not enough unique colors to cluster, it's okay for the size
	// of the extracted palette to be < k
	k = 3
	palette, _ = clusterColors(k, 100, []hcl{black, black, black, black, black, white})
	assert.LessOrEqual(t, palette.Count(), 2, "actual palette can be smaller than k")
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

func forceHCL(c color.Color) hcl {
	out, err := toHCL(c)
	if err != nil {
		panic(err)
	}
	return out
}
