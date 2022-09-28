package main

import (
	_ "image/png"
	"io/ioutil"
	"log"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	textFont font.Face
)

func init() {
	fontBytes, err := ioutil.ReadFile("fonts/ChiaroStd-B.otf")
	if err != nil {
		log.Fatal(err)
	}

	tt, err := opentype.Parse(fontBytes)
	if err != nil {
		log.Fatal(err)
	}

	textFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    16,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	game := NewGame(textFont)
	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)
	ebiten.SetWindowTitle("Ocarina Player")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
