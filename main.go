package main

import (
	"fmt"
	"log"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/bunyk/fasolasi/src/game"
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
	win.SetSmooth(true)

	// FPS tracking
	frames := 0
	second := time.Tick(time.Second * 5)

	var currentScene ui.Scene
	// currentScene = &game.MainMenu{}
	currentScene = game.NewSession("A short one.txt", "challenge", 20)
	fmt.Println(currentScene)

	for !win.Closed() {
		currentScene = currentScene.Loop(win)

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
