package scene

import (
	"mob/pkg/game"
	"mob/pkg/save"

	"github.com/sedyh/mizu/pkg/engine"
)

type Setup struct {
	Scene
}

func (s *Setup) Setup(w engine.World) {
	w.ChangeScene(&Shop{
		Free:         true,
		AllyModCount: 3,
		MustBuyOne:   true,
		Scene: Scene{
			Save: save.Save{
				GameName: "mob",
			},
			State: game.State{},
		},
	})
}
