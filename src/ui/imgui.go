package ui

// written with the help of http://iki.fi/sol/imgui/

import (
	"fmt"

	"github.com/bunyk/fasolasi/src/common"
	"github.com/bunyk/fasolasi/src/config"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
)

type UiState struct {
	hotitem    int
	activeitem int
	maxitem    int
}

var uistate = UiState{}
var imd = imdraw.New(nil)

// Prepare for IMGUI code
func Prepare() {
	uistate.hotitem = 0
	uistate.maxitem = 0
}

// Finish up after IMGUI code, end frame
func Finish(win *pixelgl.Window) {
	if win.Pressed(pixelgl.MouseButtonLeft) {
		if uistate.activeitem == 0 {
			// If the mouse is clicked, but no widget is active, we need to mark the active item unavailable so that we won't activate the next widget we drag the cursor onto.
			uistate.activeitem = -1
		}
	} else {
		// If mouse isn't down, we need to clear the active item in order not to make the widgets confused on the active state (and to enable the next clicked widget to become active).
		uistate.activeitem = 0
	}
	win.Update()
}

func nextID() int {
	uistate.maxitem++
	return uistate.maxitem
}

func Button(win *pixelgl.Window, location pixel.Rect, label string) bool {
	txt := text.New(location.Center(), common.TextAtlas)
	txt.Color = config.ButtonTextColor
	txt.Dot.X -= txt.BoundsOf(label).W() / 2
	txt.Dot.Y -= txt.BoundsOf(label).H() / 2
	fmt.Fprintln(txt, label)

	id := nextID()
	if location.Contains(win.MousePosition()) {
		uistate.hotitem = id
		if uistate.activeitem == 0 && win.Pressed(pixelgl.MouseButtonLeft) {
			uistate.activeitem = id
		}
	}
	imd.Clear()
	imd.Color = config.ButtonShadowColor
	imd.Push(
		location.Min.Add(pixel.V(8, -8)),
		location.Max.Add(pixel.V(8, -8)),
	)
	imd.Rectangle(0)

	if uistate.hotitem == id {
		imd.Color = config.SelectionColor
		if uistate.activeitem == id {
			// Button is both 'hot' and 'active'
			imd.Push(
				location.Min.Add(pixel.V(2, -2)),
				location.Max.Add(pixel.V(2, -2)),
			)
		} else {
			// Button is merely 'hot'
			imd.Push(location.Min, location.Max)
		}
	} else {
		// button is not hot, but it may be active
		imd.Color = config.ButtonColor
		imd.Push(location.Min, location.Max)
	}
	imd.Rectangle(0)
	imd.Draw(win)
	txt.Draw(win, pixel.IM)

	// If button is hot and active, but mouse button is not
	// down, the user must have clicked the button.
	if !win.Pressed(pixelgl.MouseButtonLeft) &&
		uistate.hotitem == id &&
		uistate.activeitem == id {
		return true
	}

	// Otherwise, no clicky.
	return false
}

/*
func slider(win *pixelgl.Window, location pixel.Rect, max int, value *int) bool {
  // Check for hotness
  if location.Contains(win.MousePosition()) {
    uistate.hotitem = id;
    if (uistate.activeitem == 0 && uistate.mousedown) {
      uistate.activeitem = id
	  }
  }
}
*/

func Menu(win *pixelgl.Window, location pixel.Rect, items []string) int {
	fl := FlexRows(location, config.MenuButtonWidth, config.MenuButtonHeight, config.MenuVerticalSpacing, len(items))

	clicked := -1
	for i, label := range items {
		if Button(win, fl(i), label) {
			clicked = i
		}
	}
	return clicked
}
