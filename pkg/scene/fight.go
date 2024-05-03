package scene

import (
	"log/slog"
	"mob/pkg/component"
	"mob/pkg/system"

	"github.com/sedyh/mizu/pkg/engine"
	"golang.org/x/exp/shiny/materialdesign/colornames"
)

type Fighter struct {
	component.NPC
	component.Render
	component.Rect
	component.Follow
	component.Velocity
	component.Combat
	component.Health
}

type Fight struct {
	Scene
	Room     component.Room
	PrevRoom component.Room
}

func (f *Fight) Setup(w engine.World) {
	slog.Info("fight")
	component.AddComponents(w)
	system.AddSystems(w)

	// spawn enemies
	for i, enemy := range f.Room.Enemies {
		if i > 0 {
			break // TODO REMOVE
		}
		x, y := w.Bounds().Dx()/2, w.Bounds().Dy()/2 // rand.Intn(w.Bounds().Dx()-100)+50, rand.Intn(w.Bounds().Dy()-100)+50
		fighter := Fighter{
			NPC:    enemy,
			Render: component.NewRender(component.WRenderPosition(float64(x), float64(y))),
			Rect: component.Rect{
				Color: colornames.Red400,
				W:     32,
				H:     32,
			},
			Follow:   component.Follow{Speed: 2},
			Velocity: component.Velocity{},
			Combat:   component.NewCombat(component.WithCombatNPC(&enemy)),
			Health:   component.NewHealth(),
		}
		w.AddEntities(&fighter)
	}

	// spawn allies
	for i, ally := range f.State.Allies {
		if i > 0 {
			break // TODO REMOVE
		}
		// get position of entrance
		x, y := 0, w.Bounds().Dy()/2 // rand.Intn(w.Bounds().Dx()-100)+50, rand.Intn(w.Bounds().Dy()-100)+50
		fighter := Fighter{
			NPC:    ally,
			Render: component.NewRender(component.WRenderPosition(float64(x), float64(y))),
			Rect: component.Rect{
				Color: colornames.Green400,
				W:     32,
				H:     32,
			},
			Follow:   component.Follow{Speed: 1},
			Velocity: component.Velocity{},
			Combat:   component.NewCombat(component.WithCombatNPC(&ally)),
			Health:   component.NewHealth(),
		}
		w.AddEntities(&fighter)
	}
}
