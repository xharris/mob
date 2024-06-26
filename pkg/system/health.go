package system

import (
	"mob/pkg/component"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sedyh/mizu/pkg/engine"
	"golang.org/x/exp/shiny/materialdesign/colornames"
)

type Health struct {
	*component.Render
	*component.Health
}

func (h *Health) Draw(w engine.World, _ *ebiten.Image) {
	texture := h.GetTexture(h.Health, 30, 5)
	texture.Y = -20
	// draw health bar
	ratio := float32(h.Health.Remaining / h.Health.Total)
	texture.RenderGeometry.X = -15
	texture.RenderGeometry.Y = -20
	vector.DrawFilledRect(texture.Image, 0, 0, 30*ratio, 5, colornames.White, false)
}
