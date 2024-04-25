package system

import (
	"mob/pkg/component"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sedyh/mizu/pkg/engine"
)

type ShopItemSytem struct{}

func (s *ShopItemSytem) Update(w engine.World) {
	items := w.View(component.ShopItem{}, component.Render{})
	gw, gh := ebiten.WindowSize()
	var sep float64 = float64(gw) / float64(len(items.Filter())+1)

	i := 0
	items.Each(func(e engine.Entity) {
		var shopItem *component.ShopItem
		var render *component.Render
		e.Get(&shopItem, &render)

		// arrange in grid
		render.X = sep*float64(i) - render.W/2 + sep
		render.Y = float64(gh/2) - render.H/2

		i++
	})
}
