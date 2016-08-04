package main

import (
	"image"
	"image/color"
	"math"

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/unixpickle/markovchain"
)

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Path []Point

func (p Path) Copy() Path {
	res := make(Path, len(p))
	copy(res, p)
	return res
}

const (
	maxSegLength = 3
	segStep      = 0.5
)

// A Segment is a line segment relative to a starting
// point.
// Typically, a segment's X and Y coordinates will be
// small, e.g. below 3.
type Segment struct {
	X int
	Y int
}

// Compare first compares the X and then the Y components
// of s to those in the given state.
func (s Segment) Compare(state markovchain.State) markovchain.Comparison {
	s1 := state.(Segment)
	if s.X < s1.X {
		return markovchain.Less
	} else if s.X > s1.X {
		return markovchain.Greater
	}
	if s.Y < s1.Y {
		return markovchain.Less
	} else if s.Y > s1.Y {
		return markovchain.Greater
	}
	return markovchain.Equal
}

// A SegmentTuple is a markovchain.Sample which stores
// multiple Segments at once.
type SegmentTuple []Segment

func (s SegmentTuple) Compare(state markovchain.State) markovchain.Comparison {
	s1 := state.(SegmentTuple)
	for i, seg := range s {
		res := seg.Compare(s1[i])
		if res == markovchain.Equal {
			continue
		}
		return res
	}
	return markovchain.Equal
}

// SegmentPath turns a path into a list of short segments.
func SegmentPath(p Path) []Segment {
	if len(p) < 2 {
		return nil
	}

	var res []Segment
	curPath := p.Copy()
	curPoint := p[0]
	for len(curPath) > 1 {
		dist := distance(curPoint, curPath[0])
		if dist+segStep >= maxSegLength {
			xStep := int(curPath[0].X - curPoint.X + 0.5)
			yStep := int(curPath[0].Y - curPoint.Y + 0.5)
			curPoint.X += float64(xStep)
			curPoint.Y += float64(yStep)
			res = append(res, Segment{X: xStep, Y: yStep})
		}
		if distance(curPath[0], curPath[1]) < segStep {
			curPath = curPath[1:]
			if len(curPath) == 1 {
				xStep := int(curPath[0].X - curPoint.X + 0.5)
				yStep := int(curPath[0].Y - curPoint.Y + 0.5)
				res = append(res, Segment{X: xStep, Y: yStep})
			}
		} else {
			xDiff := curPath[1].X - curPath[0].X
			yDiff := curPath[1].Y - curPath[0].Y
			ratio := segStep / distance(curPath[0], curPath[1])
			curPath[0].X += ratio * xDiff
			curPath[0].Y += ratio * yDiff
		}
	}

	return res
}

func SegmentImage(s []Segment, imageSize int) image.Image {
	newImage := image.NewRGBA(image.Rect(0, 0, imageSize, imageSize))
	ctx := draw2dimg.NewGraphicContext(newImage)
	ctx.SetLineCap(draw2d.RoundCap)
	ctx.SetLineJoin(draw2d.RoundJoin)
	ctx.SetLineWidth(2.4)
	ctx.SetStrokeColor(color.Gray{})
	ctx.BeginPath()
	ctx.MoveTo(float64(imageSize)/2, float64(imageSize)/2)
	curX := 24.0
	curY := 24.0
	for _, s := range s {
		curX += float64(s.X)
		curY += float64(s.Y)
		ctx.LineTo(curX, curY)
	}
	ctx.Stroke()
	return newImage
}

func distance(p1, p2 Point) float64 {
	return math.Sqrt(math.Pow(p1.X-p2.X, 2) + math.Pow(p1.Y-p2.Y, 2))
}
