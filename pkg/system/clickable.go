package system

import (
	"mob/pkg/component"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/sedyh/mizu/pkg/engine"
)

type Clickable struct{}

func (*Clickable) Update(w engine.World) {
	v := w.View(component.Clickable{}, component.Render{})

	for _, e := range v.Filter() {
		var clickable *component.Clickable
		var render *component.Render
		e.Get(&clickable, &render)
		if !clickable.Disabled && render.MouseInside() && inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
			clickable.Click(e)
		}
		// if clickable.Disabled {
		// 	render.AlphaLevel = component.AlphaMid
		// }
	}
}
