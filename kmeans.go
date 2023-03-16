package palettor

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/lucasb-eyer/go-colorful"
)

// clusterColors finds k clusters in the given colors using the "standard"
// k-means clustering algorithm. It returns a Palette, after running the
// algorithm up to maxIterations times.
//
// Note: in terms of the standard algorithm[1], an observation in this
// implementation is simply a color, and we use the RGB channels as Euclidean
// coordinates for the purposes of finding the distance between two colors.
//
// [1]: https://en.wikipedia.org/wiki/K-means_clustering#Standard_algorithm
func clusterColors(k, maxIterations int, colors []colorful.Color) (*Palette, error) {
	colorCount := len(colors)
	if colorCount < k {
		return nil, fmt.Errorf("too few colors for k (%d < %d)", colorCount, k)
	}

	centroids := initializeStep(k, colors)
	var clusters map[colorful.Color][]colorful.Color
	var converged bool

	// The algorithm isn't guaranteed to converge, so we put a limit on the
	// number of attempts we will make.
	var iterations int
	for iterations = 0; iterations < maxIterations; iterations++ {
		clusters = assignmentStep(centroids, colors)
		converged, centroids = updateStep(clusters)
		if converged {
			break
		}
	}

	clusterWeights := make(map[colorful.Color]float64, k)
	for centroid, cluster := range clusters {
		clusterWeights[centroid] = float64(len(cluster)) / float64(colorCount)
	}
	return &Palette{
		colorWeights: clusterWeights,
		iterations:   iterations,
		converged:    converged,
	}, nil
}

// Generate the initial list of k centroids from the given list of colors.
//
// TODO: Try other initialization methods?
// https://en.wikipedia.org/wiki/K-means_clustering#Initialization_methods
func initializeStep(k int, colors []colorful.Color) []colorful.Color {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	centroids := make([]colorful.Color, k)
	colorCount := len(colors)

	// Track random indexes we've used to avoid picking the same index for
	// multiple centroids in the case len(colors) is close to k.
	usedIndexes := make(map[int]struct{}, k)
	var index int
	for i := 0; i < k; i++ {
		for {
			index = r.Intn(colorCount)
			if _, used := usedIndexes[index]; !used {
				usedIndexes[index] = struct{}{}
				break
			}
		}
		centroids[i] = colors[index]
	}
	return centroids
}

// Assign each color to the cluster of the closest centroid.
func assignmentStep(centroids, colors []colorful.Color) map[colorful.Color][]colorful.Color {
	clusters := make(map[colorful.Color][]colorful.Color)
	for _, x := range colors {
		centroid := nearest(x, centroids)
		cluster, found := clusters[centroid]
		if !found {
			// allocate slice w/ maximum possible capacity to avoid possible
			// allocations per-append below
			cluster = make([]colorful.Color, 0, len(colors))
		}
		clusters[centroid] = append(cluster, x)
	}
	return clusters
}

// Pick new centroids from each cluster. If none of the centroids change, the
// clusters have stabilized and the algorithm has converged.
func updateStep(clusters map[colorful.Color][]colorful.Color) (bool, []colorful.Color) {
	converged := true
	newCentroids := make([]colorful.Color, 0, len(clusters))
	for centroid, cluster := range clusters {
		newCentroid := findCentroid(cluster)
		if newCentroid != centroid {
			converged = false
		}
		newCentroids = append(newCentroids, newCentroid)
	}
	return converged, newCentroids
}

// Find the color closest to the mean of the given colors.
//
// Note: I think this is a departure from the "standard" algorithm, which seems
// to instead use the actual mean of the given colors (which is likely
// not actually present in those colors).
func findCentroid(colors []colorful.Color) colorful.Color {
	center := meanColor(colors)
	return nearest(center, colors)
}

// Find the average color in a list of colors.
func meanColor(colors []colorful.Color) colorful.Color {
	var h, c, l float64
	for _, color := range colors {
		h1, c1, l1 := color.Hcl()
		h += h1
		c += c1
		l += l1
	}
	count := float64(len(colors))
	return colorful.Hcl(h/count, c/count, l/count)
}

// Find the item in the haystack to which the needle is closest.
func nearest(needle colorful.Color, haystack []colorful.Color) colorful.Color {
	var minDist float64
	var result colorful.Color
	for i, candidate := range haystack {
		dist := distanceSquared(needle, candidate)
		if i == 0 || dist < minDist {
			minDist = dist
			result = candidate
		}
	}
	return result
}

// Calculate the square of the Euclidean distance between two colors, ignoring
// the alpha channel.
func distanceSquared(a, b colorful.Color) float64 {

	// NOTE: is this really a decent HCL distance? H-values cover a much greater
	// range than C and L values.
	h1, c1, l1 := a.Hcl()
	h2, c2, l2 := b.Hcl()

	dh := h1 - h2
	dc := c1 - c2
	dl := l1 - l2
	return dh*dh + dc*dc + dl*dl
}
