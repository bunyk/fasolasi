package game

import (
	"github.com/bunyk/fasolasi/src/common"
	"github.com/bunyk/fasolasi/src/config"
	"github.com/bunyk/fasolasi/src/ui"
	"github.com/faiface/pixel/pixelgl"
)

type ModeMenu struct {
	Song string
}

func (mm *ModeMenu) Loop(win *pixelgl.Window) common.Scene {
	win.Clear(config.BackgroundColor)
	ui.Prepare()
	defer ui.Finish(win)

	choice := ui.Menu(win, win.Bounds(), []string{
		"Training",
		"Challenge",
		"← back to songs",
	})
	if win.Pressed(pixelgl.KeyEscape) {
		choice = 2
	}
	switch choice {
	case 0:
		return NewSession(mm.Song, "training")
	case 1:
		return NewSession(mm.Song, "challenge")
	case 2:
		return &SongMenu{}
	}
	return mm
}
