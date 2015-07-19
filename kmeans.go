package palettor

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"
)

// ClusterColors finds k clusters in the given colors using the "standard"
// k-means clustering algorithm. It returns a ColorPalette, after running the
// algorithm up to maxIterations times.
//
// Note: in terms of the standard algorithm[1], an observation in this
// implementation is simply a color, and we use the RGB channels as Euclidean
// coordinates for the purposes of finding the distance between two colors.
//
// [1]: https://en.wikipedia.org/wiki/K-means_clustering#Standard_algorithm
func ClusterColors(k, maxIterations int, colors []color.Color) (*ColorPalette, error) {
	colorCount := len(colors)
	if colorCount < k {
		return nil, fmt.Errorf("too few colors for k (%d < %d)", colorCount, k)
	}

	centroids := initializeCentroids(k, colors)
	var clusters map[color.Color][]color.Color

	// The algorithm isn't guaranteed to converge, so we put a limit on the
	// number of attempts we will make.
	var iterations int
	for iterations = 0; iterations < maxIterations; iterations++ {
		clusters = make(map[color.Color][]color.Color, k)

		// Assign each color to the cluster of the closest centroid.
		for _, x := range colors {
			centroid := nearest(x, centroids)
			clusters[centroid] = append(clusters[centroid], x)
		}

		// Pick new centroids from each cluster. If none of the centroids
		// change, the clusters have stabilized and we're done.
		converged := true
		newCentroids := make([]color.Color, k)
		j := 0
		for centroid, cluster := range clusters {
			newCentroid := findCentroid(cluster)
			if newCentroid != centroid {
				converged = false
			}
			newCentroids[j] = newCentroid
			j++
		}
		centroids = newCentroids
		if converged {
			break
		}
	}

	clusterWeights := make(map[color.Color]float64, k)
	for centroid, cluster := range clusters {
		clusterWeights[centroid] = float64(len(cluster)) / float64(colorCount)
	}
	return &ColorPalette{
		colorWeights: clusterWeights,
		iterations:   iterations,
	}, nil
}

// Generate the initial list of k centroids from the given list of colors.
func initializeCentroids(k int, colors []color.Color) []color.Color {
	// We just choose k random items from the list. Apparently, this is the
	// "Forgy Method". TODO: Try the Random Partition method?
	// https://en.wikipedia.org/wiki/K-means_clustering#Initialization_methods
	//
	// We take care to track the random indexes we've used to avoid picking the
	// same color for multiple centroids in the case len(colors) is
	// close to k.
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	centroids := make([]color.Color, k)
	colorCount := len(colors)
	usedIndexes := make(map[int]bool, k)
	var index int
	for i := 0; i < k; i++ {
		for {
			index = r.Intn(colorCount)
			if used, _ := usedIndexes[index]; !used {
				usedIndexes[index] = true
				break
			}
		}
		centroids[i] = colors[index]
	}
	return centroids
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
	var r, g, b, a, count uint32
	for _, x := range colors {
		r1, g1, b1, a1 := x.RGBA()
		r += r1
		g += g1
		b += b1
		a += a1
		count++
	}
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
		dist := distance(needle, candidate)
		if i == 0 || dist < minDist {
			minDist = dist
			result = candidate
		}
	}
	return result
}

// Calculate the square of the Euclidean distance between two colors.
func distance(a, b color.Color) int {
	r1, g1, b1, _ := a.RGBA()
	r2, g2, b2, _ := b.RGBA()
	dr := int(r1) - int(r2)
	dg := int(g1) - int(g2)
	db := int(b1) - int(b2)
	return dr*dr + dg*dg + db*db
}
