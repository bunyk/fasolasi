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

	// Clear window and draw UI
	win.Clear(config.BackgroundColor)

	imguiPrepare()
	defer func() {
		imguiFinish(win)
		win.Update()
	}()

	choice := Menu(win, win.Bounds(), []string{
		"Play",
		"Exit",
	})
	if win.Pressed(pixelgl.KeyEscape) {
		choice = 1
	}
	switch choice {
	case 0:
		return NewSongMenu()
	case 1:
		fmt.Println("Bye")
		win.SetClosed(true)
	}

	return mm
}
