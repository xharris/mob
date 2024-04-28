package system

import (
	"mob/pkg/component"

	"github.com/sedyh/mizu/pkg/engine"
)

type Hoverable struct {
	// *component.Hoverable
	// *component.Render
}

func (*Hoverable) Update(w engine.World) {
	v := w.View(component.Hoverable{}, component.Render{})
	for _, e := range v.Filter() {
		var render *component.Render
		var hover *component.Hoverable
		e.Get(&render, &hover)
		if render.MouseEntered() && hover.Enter != nil {
			hover.Enter(e)
		}
		if render.MouseExited() && hover.Exit != nil {
			hover.Exit(e)
		}
	}
}
