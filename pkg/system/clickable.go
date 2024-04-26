package system

import (
	"mob/pkg/component"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/sedyh/mizu/pkg/engine"
)

type Clickable struct {
	*component.Clickable
	*component.Render
}

func (c *Clickable) Update(w engine.World) {
	if c.Render.MouseInside() && inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		c.Clickable.Click()
	}
}
