package system

import (
	"mob/pkg/component"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sedyh/mizu/pkg/engine"
)

type RenderRect struct {
	*component.Rect
	*component.Render
}

func (r *RenderRect) Draw(w engine.World, _ *ebiten.Image) {
	texture := r.Render.GetTexture(r.Rect, r.Rect.W, r.Rect.H)
	vector.DrawFilledRect(texture.Image, 0, 0, float32(r.Rect.W), float32(r.Rect.H), r.Rect.Color, false)
}
