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

	fl := FlexRows(win.Bounds().Norm(), config.MenuButtonWidth, config.MenuButtonHeight, config.MenuVerticalSpacing, 3)

	if button(win, fl(0), "Training") {
		return game.NewSession(mm.Song, "training")
	}
	if button(win, fl(1), "Challenge") {
		return game.NewSession(mm.Song, "challenge")
	} // TODO: add different complexities (configure BPM)
	if button(win, fl(2), "← back to songs") || win.Pressed(pixelgl.KeyEscape) {
		return &SongMenu{}
	}

	return mm
}
