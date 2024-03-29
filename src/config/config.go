package config

import (
	_ "embed"
	"log"
	"os"

	"golang.org/x/image/colornames"

	"gopkg.in/yaml.v3"
)

// TODO: Settings section of the menu probably also should have control over some of the values

// Input tuning
const MicrophoneSampleRate = 188200
const MicrophoneBufferLength = 11025 / 2

const NoteSPS = 0.15 // Note speed in screens per second
const BlackNoteWidth = 0.3
const WhiteNoteWidth = 0.5
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

const MenuButtonWidth = 500.0
const MenuButtonHeight = 50.0
const MenuVerticalSpacing = 25.0
const MenuButtonMaxChars = 30
const MenuMaxItems = 10

const ShowFingering = false

var Songs []Song

//go:embed config.yaml
var defaultConfigData []byte

const ConfigFileName = "config.yaml"

type configFile struct {
	Songs           []Song `yaml:"songs"`
	BackgroundColor string `yaml:"background_color"`
}

func init() {
	// read config
	data, err := os.ReadFile(ConfigFileName)
	if err != nil {

		if os.IsNotExist(err) {
			data = defaultConfigData
			err := os.WriteFile(ConfigFileName, data, 0644)
			if err != nil {
				log.Printf("Failed to create default config file %s", ConfigFileName)
			}
		} else {
			log.Fatal(err)
		}
	}
	var cf configFile
	err = yaml.Unmarshal(data, &cf)
	if err != nil {
		log.Fatal(err)
	}

	// now copy values to package variables
	Songs = cf.Songs
	BackgroundColor, err = parseColor(cf.BackgroundColor)
	if err != nil {
		log.Fatalf("failed to parse background_color %s", err)
	}
}
