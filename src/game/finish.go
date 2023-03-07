package game

import (
	"fmt"

	"github.com/bunyk/fasolasi/src/config"
	"github.com/bunyk/fasolasi/src/ui"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type FinishScene struct {
	Song  string
	Mode  string
	BPM   int
	Score int
}

func (fs *FinishScene) Loop(win *pixelgl.Window) ui.Scene {
	win.Clear(config.BackgroundColor)
	renderFingering(win)
	ui.Prepare()
	defer ui.Finish(win)

	fl := ui.FlexRows(win.Bounds(), config.MenuButtonWidth, config.MenuButtonHeight, config.MenuVerticalSpacing, 4)

	ui.Label(win, fl(0), fmt.Sprintf("Score: %d", fs.Score), colornames.Black)

	if ui.Button(win, fl(1), "Retry") {
		return NewSession(fs.Song, fs.Mode, fs.BPM)
	}
	if ui.Button(win, fl(2), "Select another song") {
		return NewSongMenu()
	}
	if ui.Button(win, fl(3), "Main menu") {
		return &MainMenu{}
	}
	return fs
}
