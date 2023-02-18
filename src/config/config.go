package config

import (
	"flag"
	"log"
)

// TODO: maybe later parse some yaml to allow changing this values without rebuild
// Settings section of the menu probably also should have control over some of the values

// Input tuning
const MicrophoneSampleRate = 188200
const MicrophoneBufferLength = 11025 / 2

const NoteSPS = 0.15         // Note speed in screens per second
const TimeLinePosition = 0.3 // Position of time line in screens from the left
const NoteRadius = 40

const BreathInterval = 0.02 // Pause between notes
const TimeBeforeFirstNote = 2.0

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
