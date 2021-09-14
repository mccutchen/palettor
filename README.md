# Palettor

Yet another way to extract dominant colors from an image using [k-means clustering][1].

[![Documentation](https://pkg.go.dev/badge/github.com/mccutchen/palettor)](https://pkg.go.dev/github.com/mccutchen/palettor)
[![Build status](https://github.com/mccutchen/palettor/actions/workflows/test.yaml/badge.svg)](https://github.com/mccutchen/palettor/actions/workflows/test.yaml)
[![Code coverage](https://codecov.io/gh/mccutchen/palettor/branch/main/graph/badge.svg)](https://codecov.io/gh/mccutchen/palettor)
[![Go report card](http://goreportcard.com/badge/github.com/mccutchen/palettor)](https://goreportcard.com/report/github.com/mccutchen/palettor)


## Tests

### Unit tests

```
make test
```

### Benchmarks

```
make benchmark
```


## Usage as a library

```
go get -u github.com/mccutchen/palettor
```

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

    // Extract the 3 most dominant colors, halting the clustering algorithm
    // after 100 iterations if the clusters have not yet converged.
    k := 3
    maxIterations := 100
    palette, err := palettor.Extract(k, maxIterations, img)

    // Err will only be non-nil if k is larger than the number of pixels in the
    // input image.
    if err != nil {
        log.Fatalf("image too small")
    }

    // Palette is a mapping from color to the weight of that color's cluster,
    // which can be used as an approximation for that color's relative
    // dominance
    for _, color := range palette.Colors() {
        log.Printf("color: %v; weight: %v", color, palette.Weight(color))
    }

    // Example output:
    // 2015/07/19 10:27:52 color: {44 120 135}; weight: 0.17482142857142857
    // 2015/07/19 10:27:52 color: {140 103 150}; weight: 0.39558035714285716
    // 2015/07/19 10:27:52 color: {189 144 118}; weight: 0.42959821428571426
}
```

## The `palettor` command line application

An example command line application is provided, which reads an input image and
either a) overlays the dominant palette on the bottom of the image or b)
generates a JSON representation of the dominant color palette:

```
$ go get -u github.com/mccutchen/palettor/cmd/palettor

$ palettor -help
Usage: palettor [OPTIONS] [INPUT]

  -json
        Output color palette in JSON format
  -k int
        Palette size (default 3)
  -max int
        Maximum k-means iterations (default 500)

$ cat /Library/Desktop\ Pictures/Beach.jpg | palettor -json | jq .
[
  {
    "color": {
      "R": 70,
      "G": 134,
      "B": 154,
      "A": 255
    },
    "weight": 0.19080357142857143
  },
  {
    "color": {
      "R": 175,
      "G": 187,
      "B": 183,
      "A": 255
    },
    "weight": 0.26852678571428573
  },
  {
    "color": {
      "R": 210,
      "G": 208,
      "B": 199,
      "A": 255
    },
    "weight": 0.5406696428571428
  }
]
```


[1]: https://en.wikipedia.org/wiki/K-means_clustering#Standard_algorithm
