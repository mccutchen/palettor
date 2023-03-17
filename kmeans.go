package palettor

import (
	"fmt"
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
func clusterColors(k, maxIterations int, colors []hcl) (*Palette, error) {
	colorCount := len(colors)
	if colorCount < k {
		return nil, fmt.Errorf("too few colors for k (%d < %d)", colorCount, k)
	}

	centroids := initializeStep(k, colors)
	var clusters map[hcl][]hcl
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

	// Build palette.
	palette := &Palette{
		iterations: iterations,
		converged:  converged,
	}
	for centroid, cluster := range clusters {
		palette.add(centroid, float64(len(cluster))/float64(colorCount))
	}
	return palette, nil
}

// Generate the initial list of k centroids from the given list of colors.
//
// TODO: Try other initialization methods?
// https://en.wikipedia.org/wiki/K-means_clustering#Initialization_methods
func initializeStep(k int, colors []hcl) []hcl {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	centroids := make([]hcl, k)
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
func assignmentStep(centroids, colors []hcl) map[hcl][]hcl {
	clusters := make(map[hcl][]hcl)
	for _, x := range colors {
		centroid := nearest(x, centroids)
		cluster, found := clusters[centroid]
		if !found {
			// allocate slice w/ maximum possible capacity to avoid possible
			// allocations per-append below
			cluster = make([]hcl, 0, len(colors))
		}
		clusters[centroid] = append(cluster, x)
	}
	return clusters
}

// Pick new centroids from each cluster. If none of the centroids change, the
// clusters have stabilized and the algorithm has converged.
func updateStep(clusters map[hcl][]hcl) (bool, []hcl) {
	converged := true
	newCentroids := make([]hcl, 0, len(clusters))
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
func findCentroid(colors []hcl) hcl {
	center := meanColor(colors)
	return nearest(center, colors)
}

// Find the average color in a list of colors.
func meanColor(colors []hcl) hcl {
	var hSum, cSum, lSum float64
	for _, color := range colors {
		hSum += color.h
		cSum += color.c
		lSum += color.l
	}
	count := float64(len(colors))
	return hcl{hSum / count, cSum / count, lSum / count}
}

// Find the item in the haystack to which the needle is closest.
func nearest(needle hcl, haystack []hcl) hcl {
	var minDist float64
	var result hcl
	for i, candidate := range haystack {
		dist := needle.distanceSquared(candidate)
		if i == 0 || dist < minDist {
			minDist = dist
			result = candidate
		}
	}
	return result
}
