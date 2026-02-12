package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"go-pacman/game"
)

func main() {
	ebiten.SetWindowSize(game.ScreenWidth*game.Scale, game.ScreenHeight*game.Scale)
	ebiten.SetWindowTitle("Go Pac-Man")
	if err := ebiten.RunGame(game.New()); err != nil {
		log.Fatal(err)
	}
}
