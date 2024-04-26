package scene

import (
	"fmt"
	"mob/pkg/allymod"
	"mob/pkg/component"
	"mob/pkg/font"
	"mob/pkg/system"

	"github.com/sedyh/mizu/pkg/engine"
	"golang.org/x/exp/shiny/materialdesign/colornames"
)

type Shop struct {
	Free         bool
	AllyModCount int
}

type ShopAlly struct {
	component.Render
	component.ShopItem
	component.Tooltip
	component.Rect
	component.Clickable
	component.UIList
}

type List struct {
	component.UIList
	component.Render
}

type ShopCostLabel struct {
	component.Render
	component.UIChild
	component.UILabel
}

type Button struct {
	component.Render
	component.Clickable
	component.UILabel
	component.UIChild
}

func (s *Shop) Setup(w engine.World) {
	w.AddComponents(
		component.Render{}, component.ShopItem{}, component.Tooltip{}, component.Rect{},
		component.UIList{}, component.UILabel{}, component.UIChild{}, component.Clickable{},
	)
	w.AddSystems(
		&system.RenderSystem{}, &system.ShopItemSytem{}, &system.TooltipSystem{},
		&system.RenderRect{}, &system.UIRenderLabel{}, &system.UIListLayout{},
		&system.Clickable{},
	)

	if s.AllyModCount <= 0 {
		s.AllyModCount = 3
	}

	var purchasedAllies []ShopAlly

	// buyable ally mods
	for i := range s.AllyModCount {
		shopItemID := component.UI_ID(fmt.Sprintf("shop-item-%d", i))
		shopitem := ShopAlly{
			Render: component.NewRender(component.WRenderSize(16, 16)),
			ShopItem: component.ShopItem{
				AddMods: []allymod.Mod{
					{Name: "Slash", Desc: "Swing my sword", Type: allymod.GOOD},
					{Name: "Block", Desc: "Block an attack", Type: allymod.GOOD},
				},
			},
			Tooltip: component.Tooltip{
				Parent: "shop-item-tooltip",
			},
			Rect: component.Rect{
				Color: colornames.Blue500,
			},
			Clickable: component.Clickable{},
			UIList: component.UIList{
				ID:        shopItemID,
				Direction: component.VERTICAL,
				Reverse:   true,
			},
		}
		shopitem.Clickable.Click = func() {
			purchasedAllies = append(purchasedAllies, shopitem)
		}
		// cost
		if s.Free {
			shopitem.ShopItem.Cost = 0
		} else {
			shopitem.ShopItem.Cost = 3
		}
		costLabel := ShopCostLabel{
			Render: component.NewRender(),
			UIChild: component.UIChild{
				Parent: shopItemID,
			},
			UILabel: component.UILabel{
				Text: []component.UILabelText{
					{Text: fmt.Sprintf("$%d", shopitem.ShopItem.Cost), Color: colornames.Yellow500},
				},
			},
		}
		shopitem.Tooltip.UseShopItem(shopitem.ShopItem)
		w.AddEntities(&shopitem, &costLabel)
	}
	// tooltip container
	b := w.Bounds().Max
	mainTTArea := List{
		Render: component.NewRender(component.WRenderSize(b.X, b.Y)),
		UIList: component.UIList{
			ID:        "shop-item-tooltip",
			Direction: component.VERTICAL,
			Reverse:   true,
			// Align:     component.END,
		},
	}
	mainTTArea.Render.X = 10
	mainTTArea.Render.Y = -10
	w.AddEntities(&mainTTArea)
	// continue button
	actions := List{
		Render: component.NewRender(component.WRenderSize(b.X, b.Y)),
		UIList: component.UIList{
			ID:        "actions",
			Direction: component.VERTICAL,
			Reverse:   true,
			Align:     component.END,
		},
	}
	fx, fy := font.DefaultFont.Measure("Continue")
	continueButton := Button{
		Render: component.NewRender(component.WRenderDebug()),
		Clickable: component.Clickable{
			Click: func() {
				// go to strategy scene
			},
		},
		UILabel: component.UILabel{
			Text: []component.UILabelText{
				{Text: "Continue", Color: colornames.Red300},
			},
		},
		UIChild: component.UIChild{
			Parent: "actions",
			W:      int(fx),
			H:      int(fy),
		},
	}
	w.AddEntities(&actions, &continueButton)
}
