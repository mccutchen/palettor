# palettor

Yet another way to extract dominant colors from an image using [k-means clustering][1].

[![Build Status](https://travis-ci.org/mccutchen/palettor.svg?branch=master)](http://travis-ci.org/mccutchen/palettor)


## Example

```go
package main

import (
    "image"
    _ "image/gif"
    _ "image/jpeg"
    _ "image/png"
    "log"
    "os"

    "github.com/mccutchen/palettor"
    "github.com/nfnt/resize"
)

func main() {
    // Read an image from STDIN
    originalImg, _, err := image.Decode(os.Stdin)
    if err != nil {
        log.Fatal(err)
    }

    // Reduce it to a manageable size
    img := resize.Thumbnail(200, 200, originalImg, resize.Lanczos3)

    // Extract the 3 most dominant colors
    k := 3

    // Stop the clustering algorithm after 100 iterations if the clusters have
    // not yet converged
    maxIterations := 100

    palette, err := palettor.DominantColors(k, maxIterations, img)

    // The only possible error is if k is larger than the number of pixels in
    // the input image
    if err != nil {
        log.Fatalf("image too small")
    }

    // Palette is a mapping from color to the weight of that color's cluster,
    // which can be used as an approximation for that color's relative
    // dominance
    for color, weight := range palette {
        log.Printf("color: %v; weight: %v", color, weight)
    }

    // Example output:
    // 2015/07/19 10:27:52 color: {44 120 135}; weight: 0.17482142857142857
    // 2015/07/19 10:27:52 color: {140 103 150}; weight: 0.39558035714285716
    // 2015/07/19 10:27:52 color: {189 144 118}; weight: 0.42959821428571426
}
```


[1]: https://en.wikipedia.org/wiki/K-means_clustering#Standard_algorithm
