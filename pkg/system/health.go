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

func (h *Health) Draw(w engine.World, screen *ebiten.Image) {
	opts := h.Render.DrawImageOptions
	_, rh := h.Render.GetSize()
	opts.GeoM.Translate(0, -float64(rh))
	// draw health bar
	ratio := float32(h.Health.Remaining / h.Health.Total)
	x, y := opts.GeoM.Apply(0, 0)
	vector.DrawFilledRect(screen, float32(x), float32(y), 30*ratio, 5, colornames.White, false)
}
