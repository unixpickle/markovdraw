package main

import (
	"encoding/json"
	"fmt"
	"image/png"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Usage:", os.Args[0], "paths.json output.png")
		os.Exit(1)
	}

	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to read paths:", err)
		os.Exit(1)
	}

	var list []Path
	if err := json.Unmarshal(data, &list); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to unmarshal paths:", err)
		os.Exit(1)
	}

	for _, path := range list {
		for i := range path {
			path[i].X /= 5
			path[i].Y /= 5
		}
	}

	chain := BuildChain(list)
	sample := SampleChain(chain)
	image := SegmentImage(sample, 64)

	outFile, err := os.Create(os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to write output:", err)
		os.Exit(1)
	}
	defer outFile.Close()

	png.Encode(outFile, image)
}
