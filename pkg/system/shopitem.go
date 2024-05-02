package system

import (
	"mob/pkg/component"

	"github.com/sedyh/mizu/pkg/engine"
)

type ShopItem struct {
	*component.ShopItem
	*component.Clickable
	*component.Render
	*component.UIList
}

func (s *ShopItem) Update(w engine.World) {
	if s.ShopItem.Purchased {
		s.Clickable.Disabled = s.ShopItem.Purchased
		s.Render.AlphaLevel = component.AlphaMid

		children := w.View(component.UIChild{}, component.Render{})
		for _, c := range children.Filter() {
			var child *component.UIChild
			var render *component.Render
			c.Get(&child, &render)
			if child.Parent == s.UIList.ID {
				render.AlphaLevel = component.AlphaMid
			}
		}
	}
}
