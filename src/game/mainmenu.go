package game

import (
	"fmt"

	"github.com/bunyk/fasolasi/src/common"
	"github.com/bunyk/fasolasi/src/config"
	"github.com/bunyk/fasolasi/src/ui"
	"github.com/faiface/pixel/pixelgl"
)

type MainMenu struct {
}

func (mm *MainMenu) Loop(win *pixelgl.Window) common.Scene {
	win.Clear(config.BackgroundColor)
	renderFingering(win)

	ui.Prepare()
	defer ui.Finish(win)

	choice := ui.Menu(win, win.Bounds(), []string{
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
