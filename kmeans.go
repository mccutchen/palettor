package palettor

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"
)

// Cluster finds k clusters in the given observations and returns a mapping
// from cluster centroid to that cluster's "weight" (i.e. its size relative to
// the total number of observations) in the range [0, 1], along with the number
// of iterations taken to cluster the observations.
//
// More info: https://en.wikipedia.org/wiki/K-means_clustering#Standard_algorithm
func Cluster(k, maxIterations int, observations []color.Color) (map[color.Color]float64, int, error) {
	observationCount := len(observations)
	if observationCount < k {
		return nil, 0, fmt.Errorf("too few observations for k (%d < %d)", observationCount, k)
	}

	centroids := initializeCentroids(k, observations)
	var clusters map[color.Color][]color.Color

	// The algorithm isn't guaranteed to converge, so we put a limit on the
	// number of attempts we will make.
	var iterations int
	for iterations = 0; iterations < maxIterations; iterations++ {
		clusters = make(map[color.Color][]color.Color, k)

		// Assign each observation to the cluster of the closest centroid.
		for _, x := range observations {
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
		clusterWeights[centroid] = float64(len(cluster)) / float64(observationCount)
	}
	return clusterWeights, iterations, nil
}

func initializeCentroids(k int, observations []color.Color) []color.Color {
	// Choose k random observations as our initial centroids. Apparently, this
	// is the "Forgy Method". TODO: Try the Random Partition method?
	// https://en.wikipedia.org/wiki/K-means_clustering#Initialization_methods
	//
	// We take care to track the random indexes we've used to avoid picking the
	// same observation for multiple centroids in the case len(observations) is
	// close to k.
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	centroids := make([]color.Color, k)
	observationCount := len(observations)
	usedIndexes := make(map[int]bool, k)
	var index int
	for i := 0; i < k; i++ {
		for {
			index = r.Intn(observationCount)
			if used, _ := usedIndexes[index]; !used {
				usedIndexes[index] = true
				break
			}
		}
		centroids[i] = observations[index]
	}
	return centroids
}

// Find the observation closest to the mean of the given observations.
//
// Note: I think this is a departure from the "standard" algorithm, which seems
// to instead use the actual mean of the given observations (which is likely
// not actually present in those observations).
func findCentroid(observations []color.Color) color.Color {
	center := meanColor(observations)
	return nearest(center, observations)
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
