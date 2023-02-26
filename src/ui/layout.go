package ui

import (
	"github.com/faiface/pixel"
)

// A function that returns bounding rectangle for i-th element of flexbox (starting from 0)
type FlexMapper func(i int) pixel.Rect

// Rectangles inside parent stacked vertically, with given size and spacing, n of them
func FlexRows(parentBounds pixel.Rect, width, height, spacing float64, n int) FlexMapper {
	pc := parentBounds.Center()
	y0 := pc.Y + (float64(n)*height+float64(n-1)*spacing)/2
	return func(i int) pixel.Rect {
		return pixel.R(
			pc.X-width/2,
			y0-(height+spacing)*float64(i)-height,
			pc.X+width/2,
			y0-(height+spacing)*float64(i),
		)
	}
}

// Rectangles inside parent stacked horizontally, with given size and spacing, n of them
func FlexColumns(parentBounds pixel.Rect, width, height, spacing float64, n int) FlexMapper {
	pc := parentBounds.Center()
	x0 := pc.X - (float64(n)*width+float64(n-1)*spacing)/2
	return func(i int) pixel.Rect {
		return pixel.R(
			x0+(width+spacing)*float64(i),
			pc.Y-height/2,
			x0+(width+spacing)*float64(i)+width,
			pc.Y+height/2,
		)
	}
}
