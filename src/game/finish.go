package game

import (
	"fmt"

	"github.com/bunyk/fasolasi/src/common"
	"github.com/bunyk/fasolasi/src/config"
	"github.com/bunyk/fasolasi/src/ui"
	"github.com/faiface/pixel/pixelgl"
)

type FinishScene struct {
	Song  string
	Mode  string
	Score int
}

func (fs *FinishScene) Loop(win *pixelgl.Window) common.Scene {
	win.Clear(config.BackgroundColor)
	renderFingering(win)
	ui.Prepare()
	defer ui.Finish(win)

	fl := ui.FlexRows(win.Bounds(), config.MenuButtonWidth, config.MenuButtonHeight, config.MenuVerticalSpacing, 4)

	ui.Label(win, fl(0), fmt.Sprintf("Score: %d", fs.Score))

	if ui.Button(win, fl(1), "Retry") {
		return NewSession(fs.Song, fs.Mode)
	}
	if ui.Button(win, fl(2), "Select another song") {
		return NewSongMenu()
	}
	if ui.Button(win, fl(3), "Main menu") {
		return &MainMenu{}
	}
	return fs
}
