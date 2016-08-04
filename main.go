package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/llgcode/draw2d/draw2dimg"
)

const (
	gridWidth   = 10
	gridHeight  = 7
	cellSpacing = 5
	cellSize    = 64
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

	gridImage := image.NewRGBA(image.Rect(0, 0, gridWidth*cellSize+
		(gridWidth+1)*cellSpacing, gridHeight*cellSize+
		(gridHeight+1)*cellSpacing))
	ctx := draw2dimg.NewGraphicContext(gridImage)

	for x := 0; x < gridWidth; x++ {
		for y := 0; y < gridHeight; y++ {
			sample := SampleChain(chain)
			image := SegmentImage(sample, 64)
			imgX := x*cellSize + (x+1)*cellSpacing
			imgY := y*cellSize + (y+1)*cellSpacing
			ctx.Save()
			ctx.Translate(float64(imgX), float64(imgY))
			ctx.DrawImage(image)
			ctx.Restore()
		}
	}

	outFile, err := os.Create(os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to write output:", err)
		os.Exit(1)
	}
	defer outFile.Close()

	png.Encode(outFile, gridImage)
}
