package ui

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aquilax/truncate"
	"github.com/bunyk/fasolasi/src/common"
	"github.com/bunyk/fasolasi/src/config"
	"github.com/bunyk/fasolasi/src/game"
	"github.com/faiface/pixel/pixelgl"
)

type SongMenu struct {
	Songs  []string
	Offset int
}

func NewSongMenu() *SongMenu {
	files, err := os.ReadDir(config.SongsFolder)
	if err != nil {
		log.Fatal(err)
	}

	songs := make([]string, len(files))
	for i, file := range files {
		songs[i] = file.Name()
	}
	return &SongMenu{Songs: songs}
}

func cleanupName(fn string) string {
	return truncate.Truncate(
		strings.ReplaceAll(fn, ".txt", ""),
		config.MenuButtonMaxChars,
		"...",
		truncate.PositionEnd,
	)
}

func (sm *SongMenu) Loop(win *pixelgl.Window) common.Scene {
	if win.Pressed(pixelgl.KeyEscape) {
		return &MainMenu{}
	}
	win.Clear(config.BackgroundColor)
	imguiPrepare()
	defer func() {
		imguiFinish(win)
		win.Update()
	}()

	fl := FlexRows(win.Bounds().Norm(), config.MenuButtonWidth, config.MenuButtonHeight, config.MenuVerticalSpacing, config.MenuMaxItems)

	haveButtons := 0
	if sm.Offset > 0 {
		if button(win, fl(haveButtons), "↑ Up") {
			sm.Offset-- // Go up
		}
		haveButtons++
	}

	// we could have no more than config.MenuMaxItems buttons in menu
	// Last button is "Back", so for songs we have MenuMaxItems - haveButtons - 1 remaining buttons
	limit := config.MenuMaxItems - haveButtons - 1
	showDown := false
	if sm.Offset+limit <= len(sm.Songs) { // we do not see end of the list, need "↓ Down" button
		showDown = true
		limit -= 1 // need one more slot for "Down button"
	} else {
		limit = len(sm.Songs) - sm.Offset
	}

	for _, song := range sm.Songs[sm.Offset : sm.Offset+limit] {
		if button(win, fl(haveButtons), cleanupName(song)) {
			return game.NewSession(song)
		}
		haveButtons++
	}

	if showDown {
		if button(win, fl(haveButtons), "↓ Down") {
			sm.Offset++ // Go down
			fmt.Println("Down, new offset:", sm.Offset)
		}
		haveButtons++
	}
	if button(win, fl(haveButtons), "← Back") {
		return &MainMenu{}
	}

	return sm
}