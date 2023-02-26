package main

import (
	"fmt"
	"log"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/bunyk/fasolasi/src/common"
	"github.com/bunyk/fasolasi/src/ui"
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:     "FaSoLaSi",
		Bounds:    pixel.R(0, 0, 1024, 768),
		VSync:     true,
		Resizable: true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// FPS tracking
	frames := 0
	second := time.Tick(time.Second * 5)

	var currentScene common.Scene
	// currentScene := game.NewSession(song, config.GameMode)
	currentScene = &ui.MainMenu{}

	for !win.Closed() {
		currentScene = currentScene.Loop(win)

		win.Update()

		frames++
		select {
		case <-second:
			fmt.Printf("FPS: %d\n", frames/5)
			frames = 0
		default:
		}
	}
}

func main() {
	pixelgl.Run(run)
}
