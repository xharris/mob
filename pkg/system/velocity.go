package system

import (
	"mob/pkg/component"

	"github.com/sedyh/mizu/pkg/engine"
)

type Velocity struct {
	*component.Velocity
	*component.Render
}

func (v *Velocity) Update(w engine.World) {
	v.Render.X += v.Velocity.X
	v.Render.Y += v.Velocity.Y
}
