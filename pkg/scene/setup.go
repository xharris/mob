package scene

import (
	"mob/pkg/component"
	"mob/pkg/pawn"
	"mob/pkg/system"

	"github.com/sedyh/mizu/pkg/engine"
)

type ShopPawnItem struct {
	component.Render
	component.ShopItem
	component.Tooltip
}

type MainTooltipArea struct {
	component.RenderTooltips
	component.Render
}

type Setup struct{}

func (s *Setup) Setup(w engine.World) {
	w.AddComponents(component.Render{}, component.ShopItem{}, component.Tooltip{}, component.RenderTooltips{})
	w.AddSystems(&system.RenderSystem{}, &system.ShopItemSytem{}, &system.TooltipSystem{}, &system.RenderTooltips{})

	// 3 free pawns
	for range 3 {
		shopitem := ShopPawnItem{
			Render: component.NewRender(16, 16),
			ShopItem: component.ShopItem{
				AddMods: []pawn.Mod{
					{Name: "Slash", Desc: "Swing my sword", Type: pawn.MOD_GOOD},
					{Name: "Block", Desc: "Block an attack", Type: pawn.MOD_GOOD},
				},
			},
			Tooltip: component.Tooltip{},
		}
		shopitem.Tooltip.UseShopItem(shopitem.ShopItem)
		w.AddEntities(&shopitem)
	}

	b := w.Bounds().Max
	mainTTArea := MainTooltipArea{
		Render:         component.NewRender(float64(b.X), float64(b.Y)),
		RenderTooltips: component.RenderTooltips{},
	}
	mainTTArea.Render.X = float64(b.X) / 2
	mainTTArea.Render.Y = 0
	w.AddEntities(&mainTTArea)
}
