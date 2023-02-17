package game

import (
	"fmt"
	"image/color"
	"log"

	"github.com/bunyk/fasolasi/src/config"
	"github.com/bunyk/fasolasi/src/notes"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/gofont/goregular"
)

func time2X(width, time, currentTime float64) float64 {
	return width * (config.TimeLinePosition - (currentTime-time)*config.NoteSPS)
}

func renderNotes(win *pixelgl.Window, song []notes.SongNote, played []playedNote, time float64) {
	imd := imdraw.New(nil)
	imd.EndShape = imdraw.SharpEndShape
	width := win.Bounds().W()
	ybase := float64(win.Bounds().H()/2 - config.NoteRadius*4)

	for _, note := range song {
		renderNote(imd, time, width, ybase, false, note)
	}
	for _, note := range played {
		renderNote(imd, time, width, ybase, note.Correct, note.SongNote)
	}
	imd.Draw(win)
}

var rainbow = []color.Color{
	colornames.Red,
	colornames.Orange,
	colornames.Yellow,
	colornames.Green,
	colornames.Blue,
	colornames.Violet,
}

func renderNote(imd *imdraw.IMDraw, time, width, ybase float64, colorful bool, note notes.SongNote) {
	if note.Pitch.Name == "p" {
		return
	}
	end := note.End()
	if end < 0.0 {
		end = time
	}
	endX := time2X(width, end, time)
	if endX < 0 { // invisible already
		return
	}
	startX := time2X(width, note.Time, time)
	if startX > width { // still invisible
		return
	}
	ycenter := ybase + note.Pitch.Bottom*config.NoteRadius*2

	if note.Pitch.HasAdditionalLine() {
		imd.Color = colornames.Black
		imd.Push(
			pixel.V(startX-config.NoteRadius*2, ycenter),
			pixel.V(endX+config.NoteRadius*2, ycenter),
		)
		imd.Line(1)
	}
	border := 0.0
	if colorful {
		imd.Color = rainbow[(10+int(note.Pitch.Bottom*4))%len(rainbow)]
	} else {
		if note.Pitch.Height >= 0.5 { // white key
			imd.Color = colornames.White
			imd.Push(
				pixel.V(startX+1, ycenter-note.Pitch.Height*config.NoteRadius+1),
				pixel.V(endX-1, ycenter+note.Pitch.Height*config.NoteRadius-1),
			)
			imd.Rectangle(0)
			border = 2.0
		}
		imd.Color = colornames.Black
	}

	imd.Push(
		pixel.V(startX+1, ycenter-note.Pitch.Height*config.NoteRadius+1),
		pixel.V(endX-1, ycenter+note.Pitch.Height*config.NoteRadius-1),
	)
	imd.Rectangle(border)
}

func renderNoteLines(win *pixelgl.Window) {
	imd := imdraw.New(nil)
	imd.Color = colornames.Black

	width := win.Bounds().W()
	height := win.Bounds().H()

	ybase := float64(height/2 - config.NoteRadius*4)
	for i := 0; i < 5; i++ {
		imd.Push(
			pixel.V(0, ybase+float64(i)*config.NoteRadius*2),
			pixel.V(width, ybase+float64(i)*config.NoteRadius*2),
		)
		imd.Line(1)
	}
	imd.Push(
		pixel.V(width*config.TimeLinePosition, 0),
		pixel.V(width*config.TimeLinePosition, height),
	)
	imd.Line(3)

	imd.Draw(win)
}

func renderScore(win *pixelgl.Window, score int) {
	scoreTxt.Color = colornames.Black
	scoreTxt.Clear()
	fmt.Fprintf(scoreTxt, "%d", score)
	scoreTxt.Draw(win, pixel.IM.Moved(pixel.V(0, win.Bounds().H()-40)))
}

var scoreTxt *text.Text

func init() {
	ttf, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}
	face := truetype.NewFace(ttf, &truetype.Options{
		Size: 30,
	})
	atlas := text.NewAtlas(face, text.ASCII)
	scoreTxt = text.New(pixel.ZV, atlas)
}
