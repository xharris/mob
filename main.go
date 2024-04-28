package main

import (
	"log"
	"mob/pkg/font"
	"mob/pkg/game"
	"mob/pkg/lang"
	"mob/pkg/save"
	"mob/pkg/scene"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sedyh/mizu/pkg/engine"
)

func main() {
	lang.Init()
	font.Init()

	ebiten.SetWindowSize(600, 400)
	ebiten.SetWindowTitle("allies in a dungeon")

	setup := &scene.Setup{
		NewGame: true,
		Scene: scene.Scene{
			Save: save.Save{
				GameName: "mob",
			},
			State: game.NewState(),
		},
	}

	if err := ebiten.RunGame(engine.NewGame(setup)); err != nil {
		log.Fatal(err)
	}
}
