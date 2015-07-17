package main

import (
	"log"
	"math/rand"
	"time"

	"kmeans"
)

func main() {
	k := 3
	n := 50
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	observations := make([]int, n)
	for i := 0; i < n; i++ {
		observations[i] = r.Intn(n)
	}
	log.Printf("observations: %v", observations)

	centroids := kmeans.FindCentroids(k, observations, 100)
	log.Printf("centroids: %v", centroids)
}
