package system

import (
	"math"
	"mob/pkg/component"

	"github.com/sedyh/mizu/pkg/engine"
)

type Follow struct{}

func (*Follow) Update(w engine.World) {
	v := w.View(component.Follow{}, component.Render{}, component.Velocity{})
	entities := v.Filter()

	for _, e := range entities {
		var f *component.Follow
		var render *component.Render
		var vel *component.Velocity
		e.Get(&f, &render, &vel)

		target, targetExists := w.GetEntity(f.Target)
		if !targetExists {
			continue
		}
		// calculate velocity
		var tRender *component.Render
		target.Get(&tRender)

		dist := math.Sqrt(math.Pow(tRender.X-render.X, 2) + math.Pow(tRender.Y-render.Y, 2))
		f.TargetReached = dist < max(10, f.Radius)
		if !f.TargetReached {
			vel.X = (tRender.X - render.X) / dist * f.Speed
			vel.Y = (tRender.Y - render.Y) / dist * f.Speed
		} else {
			vel.X = 0
			vel.Y = 0
		}
	}
}
