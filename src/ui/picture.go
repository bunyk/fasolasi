package ui

import (
	"embed"
	"image"

	_ "image/png"

	"github.com/faiface/pixel"
)

//go:embed sprites/*
var f embed.FS

func LoadPicture(path string) (pixel.Picture, error) {
	file, err := f.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}
