package common

import (
	"log"

	"github.com/faiface/pixel/text"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

var TextAtlas *text.Atlas

func init() {
	ttf, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}
	face := truetype.NewFace(ttf, &truetype.Options{
		Size: 30,
	})
	runeset := append([]rune("ÜÄÖüäö←↑↓→"), text.ASCII...)
	TextAtlas = text.NewAtlas(face, runeset)
}
