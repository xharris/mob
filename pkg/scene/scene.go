package scene

import (
	"mob/pkg/game"
	"mob/pkg/save"
)

type Scene struct {
	Save  save.Save
	State game.State
}
