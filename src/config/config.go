package config

import (
	"flag"
	"log"
)

// TODO: maybe later parse some yaml to allow changing this values without rebuild

// Input tuning
const MicrophoneSampleRate = 188200
const MicrophoneBufferLength = 11025 / 2

const NoteSPS = 0.15         // Note speed in screens per second
const TimeLinePosition = 0.3 // Position of time line in screens from the left
const NoteRadius = 40

// CLI arguments, also part of config in some kind
var SongFilename string
var GameMode string

func init() {
	flag.StringVar(&SongFilename, "song", "required", "path to file of the song")
	flag.StringVar(&SongFilename, "mode", "training", "mode of game training or challenge")

	flag.Parse()
	if SongFilename == "required" {
		log.Fatal("-song is required")
	}
}
