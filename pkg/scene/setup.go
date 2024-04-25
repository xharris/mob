package scene

import (
	"mob/pkg/component"
	"mob/pkg/pawn"
	"mob/pkg/system"

	"github.com/sedyh/mizu/pkg/engine"
	"golang.org/x/exp/shiny/materialdesign/colornames"
)

type ShopPawnItem struct {
	component.Render
	component.ShopItem
	component.Tooltip
	component.Rect
}

type MainTooltipArea struct {
	component.RenderTooltips
	component.Render
}

type Setup struct{}

func (s *Setup) Setup(w engine.World) {
	w.AddComponents(component.Render{}, component.ShopItem{}, component.Tooltip{}, component.RenderTooltips{}, component.Rect{})
	w.AddSystems(&system.RenderSystem{}, &system.ShopItemSytem{}, &system.TooltipSystem{}, &system.RenderTooltips{}, &system.RenderRect{})

	// 3 free pawns
	for range 3 {
		shopitem := ShopPawnItem{
			Render: component.NewRender(component.WRenderSize(16, 16)),
			ShopItem: component.ShopItem{
				AddMods: []pawn.Mod{
					{Name: "Slash", Desc: "Swing my sword", Type: pawn.MOD_GOOD},
					{Name: "Block", Desc: "Block an attack", Type: pawn.MOD_GOOD},
				},
			},
			Tooltip: component.Tooltip{},
			Rect: component.Rect{
				Color: colornames.Blue500,
			},
		}
		shopitem.Tooltip.UseShopItem(shopitem.ShopItem)
		w.AddEntities(&shopitem)
	}

	b := w.Bounds().Max
	mainTTArea := MainTooltipArea{
		Render:         component.NewRender(component.WRenderSize(b.X, b.Y)),
		RenderTooltips: component.RenderTooltips{},
	}
	mainTTArea.Render.X = 0 // float64(b.X) / 2
	mainTTArea.Render.Y = 0
	w.AddEntities(&mainTTArea)
}
