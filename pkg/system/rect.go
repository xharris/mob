package system

import (
	"mob/pkg/component"

	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sedyh/mizu/pkg/engine"
)

type RenderRect struct {
	*component.Rect
	*component.Render
}

func (r *RenderRect) Update(w engine.World) {
	size := r.Render.Image.Bounds()
	r.Render.OX, r.Render.OY = float64(size.Dx()/2), float64(size.Dy()/2)
	vector.DrawFilledRect(r.Render.Image, 0, 0, float32(size.Dx()), float32(size.Dy()), r.Rect.Color, false)
}
