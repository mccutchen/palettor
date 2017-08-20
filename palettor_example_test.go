package palettor

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/png"
	"log"
)

func Example() {
	// cat testdata/example.png | base64
	var exampleData = []byte("iVBORw0KGgoAAAANSUhEUgAAAAQAAAAECAIAAAAmkwkpAAAABGdBTUEAALGPC/xhBQAAACBjSFJNAAB6JgAAgIQAAPoAAACA6AAAdTAAAOpgAAA6mAAAF3CculE8AAAACXBIWXMAAAsTAAALEwEAmpwYAAACMGlUWHRYTUw6Y29tLmFkb2JlLnhtcAAAAAAAPHg6eG1wbWV0YSB4bWxuczp4PSJhZG9iZTpuczptZXRhLyIgeDp4bXB0az0iWE1QIENvcmUgNS40LjAiPgogICA8cmRmOlJERiB4bWxuczpyZGY9Imh0dHA6Ly93d3cudzMub3JnLzE5OTkvMDIvMjItcmRmLXN5bnRheC1ucyMiPgogICAgICA8cmRmOkRlc2NyaXB0aW9uIHJkZjphYm91dD0iIgogICAgICAgICAgICB4bWxuczp4bXA9Imh0dHA6Ly9ucy5hZG9iZS5jb20veGFwLzEuMC8iCiAgICAgICAgICAgIHhtbG5zOnRpZmY9Imh0dHA6Ly9ucy5hZG9iZS5jb20vdGlmZi8xLjAvIj4KICAgICAgICAgPHhtcDpDcmVhdG9yVG9vbD5BY29ybiB2ZXJzaW9uIDQuNS44PC94bXA6Q3JlYXRvclRvb2w+CiAgICAgICAgIDx0aWZmOkNvbXByZXNzaW9uPjU8L3RpZmY6Q29tcHJlc3Npb24+CiAgICAgICAgIDx0aWZmOllSZXNvbHV0aW9uPjcyPC90aWZmOllSZXNvbHV0aW9uPgogICAgICAgICA8dGlmZjpYUmVzb2x1dGlvbj43MjwvdGlmZjpYUmVzb2x1dGlvbj4KICAgICAgPC9yZGY6RGVzY3JpcHRpb24+CiAgIDwvcmRmOlJERj4KPC94OnhtcG1ldGE+CkPIGwAAAAAUSURBVAgdY/zPAAb/QTQThE2IBACKLAMB3J9gzQAAAABJRU5ErkJggg==")

	decoder := base64.NewDecoder(base64.StdEncoding, bytes.NewReader(exampleData))
	img, err := png.Decode(decoder)
	if err != nil {
		log.Fatal(err)
	}

	// For a real-world use case, it's best to use something like
	// github.com/nfnt/resize to transform images into a manageable size before
	// extracting colors:
	//
	//     img = resize.Thumbnail(200, 200, img, resize.Lanczos3)
	//
	// In this example, we're already starting from a tiny image.

	// Extract the 3 most dominant colors, halting the clustering algorithm
	// after 100 iterations if the clusters have not yet converged.
	palette, _ := Extract(3, 100, img)

	// Palette is a mapping from color to the weight of that color's cluster,
	// which can be used as an approximation for that color's relative
	// dominance
	for _, color := range palette.Colors() {
		fmt.Printf("color: %v; weight: %v\n", color, palette.Weight(color))
	}

	// Example output:
	// color: {255 0 0 255}; weight: 0.25
	// color: {255 255 255 255}; weight: 0.25
	// color: {0 0 0 255}; weight: 0.5
}
