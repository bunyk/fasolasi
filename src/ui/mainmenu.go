package ui

import (
	"fmt"

	"github.com/bunyk/fasolasi/src/common"
	"github.com/bunyk/fasolasi/src/config"
	"github.com/faiface/pixel/pixelgl"
)

type MainMenu struct {
}

func (mm *MainMenu) Loop(win *pixelgl.Window) common.Scene {

	if win.Pressed(pixelgl.KeyEscape) {
		win.SetClosed(true)
	}

	// Clear window and draw UI
	win.Clear(config.BackgroundColor)

	imguiPrepare()
	defer func() {
		imguiFinish(win)
		win.Update()
	}()

	fl := FlexRows(win.Bounds().Norm(), config.MenuButtonWidth, config.MenuButtonHeight, config.MenuVerticalSpacing, 4)

	if button(win, fl(0), "Play") {
		return NewSongMenu()
	}
	if button(win, fl(1), "Settings") {
		fmt.Println("TODO")
	}
	if button(win, fl(2), "Exit") {
		fmt.Println("Bye")
		win.SetClosed(true)
	}

	return mm
}
