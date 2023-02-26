package config

import (
	"flag"
	"log"

	"golang.org/x/image/colornames"
)

// TODO: maybe later parse some yaml to allow changing this values without rebuild
// Settings section of the menu probably also should have control over some of the values

const SongsFolder = "./songs"

// Input tuning
const MicrophoneSampleRate = 188200
const MicrophoneBufferLength = 11025 / 2

const NoteSPS = 0.15         // Note speed in screens per second
const TimeLinePosition = 0.3 // Position of time line in screens from the left

// We draw it more to the right, because you should blow earlier, audio input buffer takes time to fill
// and screen is not updated instanly (just 60 or 50 times per second)

const TimeLinePositionLatency = TimeLinePosition + NoteSPS*(MicrophoneBufferLength/MicrophoneSampleRate+0.02)
const NoteRadius = 40

const BreathInterval = 0.05 // Pause between notes
const TimeBeforeFirstNote = 2.0

var BackgroundColor = colornames.Antiquewhite

var ButtonColor = colornames.White
var ButtonTextColor = colornames.Black
var ButtonShadowColor = colornames.Darkgray
var SelectionColor = colornames.Lightblue
var MenuButtonWidth = 500.0
var MenuButtonHeight = 60.0
var MenuVerticalSpacing = 30.0
var MenuButtonMaxChars = 30
var MenuMaxItems = 7

// CLI arguments, also part of config in some kind
var SongFilename string
var GameMode string

func init() {
	flag.StringVar(&SongFilename, "song", "required", "path to file of the song")
	flag.StringVar(&GameMode, "mode", "training", "mode of game training or challenge")

	flag.Parse()
	if SongFilename == "required" {
		log.Fatal("-song is required")
	}
	if !(GameMode == "training" || GameMode == "challenge") {
		log.Fatal("mode could be only training or challenge")
	}
}
