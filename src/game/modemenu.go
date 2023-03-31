package game

import (
	"github.com/bunyk/fasolasi/src/config"
	"github.com/bunyk/fasolasi/src/ui"
	"github.com/faiface/pixel/pixelgl"
)

type ModeMenu struct {
	SongID int
}

func (mm *ModeMenu) Loop(win *pixelgl.Window) ui.Scene {
	win.Clear(config.BackgroundColor)
	renderFingering(win)
	ui.Prepare()
	defer ui.Finish(win)

	choice := ui.Menu(win, win.Bounds(), []string{
		"Training",
		"30 bpm challenge",
		"60 bpm challenge",
		"90 bpm challenge",
		"← back to songs",
	})
	if win.Pressed(pixelgl.KeyEscape) {
		choice = 4
	}
	switch choice {
	case 0:
		return NewSession(mm.SongID, "training", 60)
	case 1:
		return NewSession(mm.SongID, "challenge", 30)
	case 2:
		return NewSession(mm.SongID, "challenge", 60)
	case 3:
		return NewSession(mm.SongID, "challenge", 90)
	case 4:
		return &SongMenu{}
	}
	return mm
}
