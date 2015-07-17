package kmeans

import (
	"log"
	"math/rand"
	"time"
)

// FindClusters finds clusters
func FindClusters(k int, xs []int, maxIterations int) map[int][]int {
	centroids := initializeCentroids(k, xs)

	var clusters map[int][]int

	for i := 0; i < maxIterations; i++ {
		clusters = make(map[int][]int)

		log.Println()
		log.Printf("iteration: %d", i)
		log.Printf("centroids: %v", centroids)
		log.Printf("clusters:  %v", clusters)

		// assign observations to centroids
		log.Printf("=== assignment")
		for _, x := range xs {
			log.Printf("... x == %v", x)
			centroid := nearest(x, centroids)
			clusters[centroid] = append(clusters[centroid], x)
		}

		// update centroids
		converged := true
		newCentroids := make([]int, k)
		j := 0
		log.Printf("=== update")
		log.Printf("... len(clusters) == %v", len(clusters))
		for centroid, observations := range clusters {
			newCentroid := findCentroid(observations)
			log.Printf("... observations == %v", observations)
			log.Printf("... old centroid == %v", centroid)
			log.Printf("... new centroid == %v", newCentroid)
			log.Printf("... j == %v", j)
			if newCentroid != centroid {
				log.Printf("!!! NOT CONVERGED")
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
	return clusters
}

// FindCentroids finds centroids
func FindCentroids(k int, xs []int, maxIterations int) []int {
	clusters := FindClusters(k, xs, maxIterations)
	centroids := make([]int, len(clusters))
	i := 0
	for key := range clusters {
		centroids[i] = key
		i++
	}
	return centroids
}

func nearest(x int, xs []int) int {
	var minDist int
	var nearestX int
	for i, xn := range xs {
		dist := distance(x, xn)
		if i == 0 || dist < minDist {
			minDist = dist
			nearestX = xn
		}
	}
	return nearestX
}

func distance(a, b int) int {
	delta := a - b
	return delta * delta
}

func initializeCentroids(k int, xs []int) []int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	total := len(xs)
	centroids := make([]int, k)
	for i := 0; i < k; i++ {
		centroids[i] = xs[r.Intn(total)]
	}
	return centroids
}

func findCentroid(xs []int) int {
	total := 0
	for _, x := range xs {
		total += x
	}
	center := total / len(xs)
	return nearest(center, xs)
}
