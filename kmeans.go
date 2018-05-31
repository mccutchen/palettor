package palettor

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"
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
func clusterColors(k, maxIterations int, colors []color.Color) (*Palette, error) {
	colorCount := len(colors)
	if colorCount < k {
		return nil, fmt.Errorf("too few colors for k (%d < %d)", colorCount, k)
	}

	centroids := initializeStep(k, colors)
	var clusters map[color.Color][]color.Color
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

	clusterWeights := make(map[color.Color]float64, k)
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
func initializeStep(k int, colors []color.Color) []color.Color {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	centroids := make([]color.Color, k)
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
func assignmentStep(centroids, colors []color.Color) map[color.Color][]color.Color {
	clusters := make(map[color.Color][]color.Color)
	for _, x := range colors {
		centroid := nearest(x, centroids)
		clusters[centroid] = append(clusters[centroid], x)
	}
	return clusters
}

// Pick new centroids from each cluster. If none of the centroids change, the
// clusters have stabilized and the algorithm has converged.
func updateStep(clusters map[color.Color][]color.Color) (bool, []color.Color) {
	converged := true
	var newCentroids []color.Color
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
func findCentroid(colors []color.Color) color.Color {
	center := meanColor(colors)
	return nearest(center, colors)
}

// Find the average color in a list of colors.
func meanColor(colors []color.Color) color.Color {
	var r, g, b, a uint32
	for _, x := range colors {
		r1, g1, b1, a1 := x.RGBA()
		r += r1
		g += g1
		b += b1
		a += a1
	}
	count := uint32(len(colors))
	return &color.RGBA64{
		R: uint16(r / count),
		G: uint16(g / count),
		B: uint16(b / count),
		A: uint16(a / count),
	}
}

// Find the item in the haystack to which the needle is closest.
func nearest(needle color.Color, haystack []color.Color) color.Color {
	var minDist int
	var result color.Color
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
func distanceSquared(a, b color.Color) int {
	r1, g1, b1, _ := a.RGBA()
	r2, g2, b2, _ := b.RGBA()
	dr := int(r1) - int(r2)
	dg := int(g1) - int(g2)
	db := int(b1) - int(b2)
	return dr*dr + dg*dg + db*db
}
