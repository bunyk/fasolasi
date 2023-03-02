package ui

import (
	"github.com/bunyk/fasolasi/src/common"
	"github.com/bunyk/fasolasi/src/config"
	"github.com/bunyk/fasolasi/src/game"
	"github.com/faiface/pixel/pixelgl"
)

type ModeMenu struct {
	Song string
}

func (mm *ModeMenu) Loop(win *pixelgl.Window) common.Scene {
	win.Clear(config.BackgroundColor)
	imguiPrepare()
	defer func() {
		imguiFinish(win)
		win.Update()
	}()

	choice := Menu(win, win.Bounds(), []string{
		"Training",
		"Challenge",
		"← back to songs",
	})
	if win.Pressed(pixelgl.KeyEscape) {
		choice = 2
	}
	switch choice {
	case 0:
		return game.NewSession(mm.Song, "training")
	case 1:
		return game.NewSession(mm.Song, "challenge")
	case 2:
		return &SongMenu{}
	}
	return mm
}
