package scene

import (
	"log/slog"
	"mob/pkg/component"
	"mob/pkg/system"

	"github.com/sedyh/mizu/pkg/engine"
)

type Fight struct {
	Scene
}

func (f *Fight) Setup(w engine.World) {
	slog.Info("fight")
	component.AddComponents(w)
	system.AddSystems(w)
}
