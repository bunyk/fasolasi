package common

import "github.com/faiface/pixel/pixelgl"

// Returns self or next scene on each iteration
type Scene interface {
	Loop(w *pixelgl.Window) Scene
}
