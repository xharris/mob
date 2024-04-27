package scene

import (
	"fmt"
	"log/slog"
	"mob/pkg/allymod"
	"mob/pkg/component"
	"mob/pkg/system"

	"github.com/sedyh/mizu/pkg/engine"
	"golang.org/x/exp/shiny/materialdesign/colornames"
)

type Shop struct {
	Free         bool
	AllyModCount int
}

type ShopAllyContainer struct {
	component.Render
	component.ShopItem
	component.Clickable
	component.UIList
	component.Hoverable
	component.UIChild
}

type ShopAllyRect struct {
	component.Render
	component.Rect
	component.UIChild
}

type List struct {
	component.UIChild
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

type UI struct {
	component.Render
	component.UIGrid
}

type Label struct {
	component.Render
	component.UILabel
	component.UIChild
}

func (s *Shop) Setup(w engine.World) {
	w.AddComponents(
		component.Render{}, component.ShopItem{}, component.Rect{},
		component.UIList{}, component.UILabel{}, component.UIChild{}, component.Clickable{},
		component.UIGrid{}, component.Hoverable{},
	)
	w.AddSystems(
		&system.RenderSystem{},
		&system.RenderRect{}, &system.UIRenderLabel{}, &system.UIListLayout{},
		&system.Clickable{}, &system.UIGridLayout{}, &system.Hoverable{},
	)

	if s.AllyModCount <= 0 {
		s.AllyModCount = 3
	}

	var purchased []component.ShopItem

	shopItemsUI := UI{
		Render: component.NewRender(component.WRenderDebug()),
		UIGrid: component.UIGrid{
			ID:      "shop-items-ui",
			Columns: min(5, s.AllyModCount),
			Align:   component.CENTER,
			Justify: component.CENTER,
		},
	}
	w.AddEntities(&shopItemsUI)

	// buyable ally mods
	for i := range s.AllyModCount {
		shopItemID := component.UI_ID(fmt.Sprintf("shop-item-%d", i))
		shopitem := ShopAllyContainer{
			Render: component.NewRender(component.WRenderDebug()),
			ShopItem: component.ShopItem{
				AddMods: []allymod.Mod{
					{Name: "Slash", Type: allymod.Attack, Target: allymod.Enemy},
					{Name: "Block", Type: allymod.Buff, Target: allymod.Self},
					{Name: "Sleepy", Desc: "Might take a nap", Type: allymod.Debuff, Target: allymod.Self},
				},
			},
			Clickable: component.Clickable{},
			Hoverable: component.Hoverable{},
			UIList: component.UIList{
				ID:        shopItemID,
				Direction: component.VERTICAL,
				Align:     component.CENTER,
				// Justify:     component.CENTER,
				FitContents: true,
			},
			UIChild: component.UIChild{
				Parent: "shop-items-ui",
			},
		}
		// show mods on hover
		shoptItemTooltipID := component.UI_ID(fmt.Sprintf("shop-item-tooltip-%d", i))
		shopitem.Hoverable.Enter = func() {
			label := Label{
				Render:  component.NewRender(),
				UILabel: component.UILabel{},
				UIChild: component.UIChild{
					ID:     shoptItemTooltipID,
					Parent: "shop-item-tooltip",
				},
			}
			for _, addMod := range shopitem.AddMods {
				color := colornames.Green300
				nameSuffix := "!"
				if !addMod.IsGood() {
					nameSuffix = "?"
					color = colornames.Red300
				}
				if addMod.Type == allymod.Debuff && addMod.Target == allymod.Self {
					nameSuffix = "..."
					color = colornames.Orange300
				}
				label.Text = append(label.Text,
					component.UILabelText{Text: addMod.Name, Color: color},
					component.UILabelText{Text: nameSuffix + " ", Color: color},
					component.UILabelText{Text: addMod.Desc, Color: colornames.Grey100},
					component.UILabelText{Newline: true},
				)
			}
			w.AddEntities(&label)
		}
		shopitem.Hoverable.Exit = func() {
			labels := w.View(component.UIChild{})
			for _, e := range labels.Filter() {
				var ch *component.UIChild
				e.Get(&ch)
				if ch.ID == shoptItemTooltipID {
					w.RemoveEntity(e)
				}
			}
		}
		shopitem.Clickable.Click = func() {
			purchased = append(purchased, shopitem.ShopItem)
		}
		w.AddEntities(&shopitem)
		// ally image
		allyRect := ShopAllyRect{
			Render: component.NewRender(component.WRenderSize(32, 32)),
			Rect: component.Rect{
				Color: colornames.Blue500,
			},
			UIChild: component.UIChild{
				Parent: shopItemID,
			},
		}
		w.AddEntities(&allyRect)
		// cost
		if s.Free {
			shopitem.ShopItem.Cost = i // 0
		} else {
			shopitem.ShopItem.Cost = i // 3
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
		w.AddEntities(&costLabel)
	}
	ui := UI{
		Render: component.NewRender(),
		UIGrid: component.UIGrid{
			ID:      "ui",
			Rows:    3,
			Columns: 3,
		},
	}
	w.AddEntities(&ui)
	// tooltip container
	b := w.Bounds().Max
	mainTTArea := List{
		Render: component.NewRender(),
		UIList: component.UIList{
			ID:        "shop-item-tooltip",
			Direction: component.VERTICAL,
			Justify:   component.END,
		},
		UIChild: component.UIChild{
			Parent: "ui",
			X:      0,
			Y:      2,
		},
	}
	w.AddEntities(&mainTTArea)
	// continue button
	actions := List{
		Render: component.NewRender(component.WRenderSize(b.X, b.Y)),
		UIList: component.UIList{
			ID:        "actions",
			Direction: component.VERTICAL,
			Align:     component.END,
			Justify:   component.END,
		},
		UIChild: component.UIChild{
			Parent: "ui",
			X:      2,
			Y:      2,
		},
	}
	continueButton := Button{
		Render: component.NewRender(),
		Clickable: component.Clickable{
			Click: func() {
				// go to strategy scene
				slog.Info("go to strategy scene")
			},
		},
		UILabel: component.UILabel{
			Text: []component.UILabelText{
				{Text: "Continue", Color: colornames.Red300},
			},
		},
		UIChild: component.UIChild{
			Parent: "actions",
		},
	}
	w.AddEntities(&actions, &continueButton)
}
