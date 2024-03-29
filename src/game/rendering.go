package game

import (
	"fmt"
	"image/color"

	"github.com/bunyk/fasolasi/src/config"
	"github.com/bunyk/fasolasi/src/notes"
	"github.com/bunyk/fasolasi/src/ui"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
)

func time2X(width, time, currentTime float64) float64 {
	return width * (config.TimeLinePosition - (currentTime-time)*config.NoteSPS)
}

func renderNotes(win *pixelgl.Window, song []notes.SongNote, played []playedNote, time float64) {
	width := win.Bounds().W()
	ybase := float64(win.Bounds().H()/2 - config.NoteRadius*4)

	for _, note := range song {
		renderNote(win, time, width, ybase, false, note)
	}
	for _, note := range played {
		renderNote(win, time, width, ybase, note.Correct, note.SongNote)
	}
}

var rainbow = []color.Color{
	colornames.Red,
	colornames.Orange,
	colornames.Yellow,
	colornames.Green,
	colornames.Blue,
	colornames.Violet,
}

func renderNote(win *pixelgl.Window, time, width, ybase float64, colorful bool, note notes.SongNote) {
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

	imd := imdraw.New(nil)
	imd.EndShape = imdraw.SharpEndShape
	if note.Pitch.HasAdditionalLine() {
		imd.Color = colornames.Black
		imd.Push(
			pixel.V(startX-config.NoteRadius*2, ycenter),
			pixel.V(endX+config.NoteRadius*2, ycenter),
		)
		imd.Line(1)
	}
	border := 0.0
	var textColor color.Color = colornames.White
	if colorful {
		imd.Color = rainbow[(10+int(note.Pitch.Bottom*4))%len(rainbow)]
		textColor = oppositeColor(imd.Color)
	} else {
		if !note.Pitch.IsHalf { // white key
			imd.Color = colornames.White
			textColor = colornames.Black
			imd.Push(
				pixel.V(startX+1, ycenter-config.WhiteNoteWidth*config.NoteRadius+1),
				pixel.V(endX-1, ycenter+config.WhiteNoteWidth*config.NoteRadius-1),
			)
			imd.Rectangle(0)
			border = 2.0
		}
		imd.Color = colornames.Black
	}

	height := config.WhiteNoteWidth
	if note.Pitch.IsHalf {
		height = config.BlackNoteWidth
	}
	corner1 := pixel.V(startX+1, ycenter-height*config.NoteRadius+1)
	corner2 := pixel.V(endX-1, ycenter+height*config.NoteRadius-1)
	imd.Push(corner1, corner2)
	imd.Rectangle(border)
	imd.Draw(win)

	ui.Label(win, pixel.Rect{Min: corner1, Max: corner2}, note.Pitch.Title(), textColor)
}

func oppositeColor(c color.Color) color.Color {
	r, g, b, _ := c.RGBA()
	return color.RGBA{uint8(255 - r), uint8(255 - g), uint8(255 - b), 255}
}

func renderProgress(win *pixelgl.Window, progress float64) {
	if progress > 1.0 {
		progress = 1.0
	}
	if progress < 0.0 {
		progress = 0.0
	}
	imd := imdraw.New(nil)
	width := win.Bounds().W()

	imd.Color = colornames.Black
	imd.Push(
		pixel.V(2, 2),
		pixel.V(width-2, 10),
	)
	imd.Rectangle(1)

	imd.Color = colornames.Red
	imd.Push(
		pixel.V(2, 2),
		pixel.V(2+(width-5)*progress, 9),
	)
	imd.Rectangle(0)

	imd.Draw(win)
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
		pixel.V(width*config.TimeLinePositionLatency, 0),
		pixel.V(width*config.TimeLinePositionLatency, height),
	)
	imd.Line(3)

	imd.Draw(win)
}

func renderScore(win *pixelgl.Window, score int, big bool) {
	scoreTxt := text.New(pixel.ZV, ui.TextAtlas)
	scoreTxt.Color = colornames.Black
	scoreTxt.Clear()
	fmt.Fprintf(scoreTxt, "%d", score)
	var pos pixel.Matrix
	if big {
		pos = pixel.IM.Moved(win.Bounds().Center()).Scaled(win.Bounds().Center(), 5)
	} else {
		pos = pixel.IM.Moved(pixel.V(0, win.Bounds().H()-40))
	}
	scoreTxt.Draw(win, pos)
}

func hightLightNote(win *pixelgl.Window, color color.Color, note notes.Pitch) {
	imd := imdraw.New(nil)
	width := win.Bounds().W()
	ybase := float64(win.Bounds().H()/2 - config.NoteRadius*4)
	if note.Name == "p" {
		return
	}
	ycenter := ybase + note.Bottom*config.NoteRadius*2
	imd.Color = color

	height := config.WhiteNoteWidth
	if note.IsHalf {
		height = config.BlackNoteWidth
	}
	imd.Push(
		pixel.V(0, ycenter-height*config.NoteRadius+1),
		pixel.V(width, ycenter+height*config.NoteRadius-1),
	)
	imd.Rectangle(0)
	imd.Draw(win)
}

func soundVisualization(win *pixelgl.Window, col color.Color, data [][2]float64) {
	imd := imdraw.New(nil)
	imd.Color = col
	width := win.Bounds().W()
	height := win.Bounds().H()
	every := 1.0
	if width > 0 {
		every = float64(len(data)) / width
	}
	for i := 0.0; i < width; i += 2.0 {
		imd.Push(pixel.V(float64(i), height*(0.1+data[int(i*every)][0]*0.2)))
	}
	imd.Line(1)
	imd.Draw(win)
}

var recorderSprite *pixel.Sprite

func init() {
	pic, err := ui.LoadPicture("sprites/recorder.png")
	if err != nil {
		panic(err)
	}

	recorderSprite = pixel.NewSprite(pic, pic.Bounds())
}

func renderFingering(win *pixelgl.Window) {
	scale := win.Bounds().H() / recorderSprite.Frame().H()
	recorderSprite.Draw(win, pixel.IM.
		Scaled(pixel.ZV, scale).
		Moved(pixel.V(50, win.Bounds().H()/2.0)),
	)
}

func renderMessage(win *pixelgl.Window, msg string) {
	txt := text.New(pixel.ZV, ui.TextAtlas)
	txt.Color = colornames.Black
	fmt.Fprintln(txt, msg)
	var pos pixel.Matrix
	pos = pixel.IM.Moved(win.Bounds().Center().Sub(txt.Bounds().Center()))
	txt.Draw(win, pos)
}
