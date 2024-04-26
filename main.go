package main

import (
	"log"
	"mob/pkg/font"
	"mob/pkg/scene"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sedyh/mizu/pkg/engine"
)

func main() {
	font.Init()

	ebiten.SetWindowSize(600, 400)
	ebiten.SetWindowTitle("pawns in dungeon game")

	if err := ebiten.RunGame(engine.NewGame(&scene.Setup{})); err != nil {
		log.Fatal(err)
	}
}
