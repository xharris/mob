package scene

import (
	"fmt"
	"log/slog"
	"mob/pkg/component"
	"mob/pkg/mods"
	"mob/pkg/system"

	"github.com/sedyh/mizu/pkg/engine"
	"golang.org/x/exp/shiny/materialdesign/colornames"
)

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

type Shop struct {
	Scene
	Free           bool
	ModCount       int
	AllyCount      int
	MustBuyOne     bool
	continueShown  bool
	purchasedCount int
	PurchaseLimit  int
}

func (s *Shop) Setup(w engine.World) {
	slog.Info("shop")
	component.AddComponents(w)
	system.AddSystems(w)

	shopItemsUI := UI{
		Render: component.NewRender(),
		UIGrid: component.UIGrid{
			ID:      "shop-items-ui",
			Columns: min(5, s.ModCount+s.AllyCount),
			Align:   component.CENTER,
			Justify: component.CENTER,
		},
	}
	w.AddEntities(&shopItemsUI)

	// buyable allies
	for i := range s.AllyCount {
		shoptItemTooltipID := component.UI_ID(fmt.Sprintf("shop-item-tooltip-%d", i))
		shopItemID := component.UI_ID(fmt.Sprintf("shop-item-%d", i))
		shopitem := ShopAllyContainer{
			Render: component.NewRender(),
			ShopItem: component.ShopItem{
				AddMods: []component.Mod{mods.Slash(), mods.Block(), mods.Sleepy()},
				Name:    "Sir Bobbington",
			},
			Clickable: component.Clickable{
				Click: func(e engine.Entity) {
					var item *component.ShopItem
					e.Get(&item)

					reachedPurchaseLimit := s.PurchaseLimit > 0 && s.purchasedCount >= s.PurchaseLimit
					if item.Purchased || reachedPurchaseLimit {
						return
					}

					// add new ally
					npc := component.NewNPC(
						component.WithMods(item.AddMods),
					)
					s.Scene.State.Allies = append(s.Scene.State.Allies, npc)

					item.Purchased = true
					s.purchasedCount++

					if s.purchasedCount > 0 || !s.MustBuyOne {
						s.showContinueButton(w)
					}
				},
			},
			// show mods on hover
			Hoverable: component.Hoverable{
				Enter: func(e engine.Entity) {
					var shopItemH *component.ShopItem
					e.Get(&shopItemH)
					label := Label{
						Render: component.NewRender(),
						UILabel: component.UILabel{
							Text: []component.UILabelText{
								{Text: shopItemH.Name, Color: colornames.BlueGrey300},
								{Newline: true},
							},
						},
						UIChild: component.UIChild{
							ID:     shoptItemTooltipID,
							Parent: "shop-item-tooltip",
						},
					}
					for _, addMod := range shopItemH.AddMods {
						color := colornames.Green300
						nameSuffix := "!"
						if !addMod.IsGood() {
							nameSuffix = "?"
							color = colornames.Red300
						}
						if addMod.Type == component.Debuff && addMod.Target == component.TargetSelf {
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
				},
				Exit: func(engine.Entity) {
					labels := w.View(component.UIChild{})
					for _, e := range labels.Filter() {
						var ch *component.UIChild
						e.Get(&ch)
						if ch.ID == shoptItemTooltipID {
							w.RemoveEntity(e)
						}
					}
				},
			},
			UIList: component.UIList{
				ID:          shopItemID,
				Direction:   component.VERTICAL,
				Align:       component.CENTER,
				FitContents: true,
			},
			UIChild: component.UIChild{
				Parent: "shop-items-ui",
			},
		}
		w.AddEntities(&shopitem)
		// ally image
		allyRect := ShopAllyRect{
			Render: component.NewRender(component.WRenderSize(32, 32)),
			Rect: component.Rect{
				Color: colornames.Green400,
			},
			UIChild: component.UIChild{
				Parent: shopItemID,
			},
		}
		w.AddEntities(&allyRect)
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
		Render: component.NewRender(),
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
	w.AddEntities(&actions)

	if !s.MustBuyOne {
		s.showContinueButton(w)
	}
}

func (s *Shop) showContinueButton(w engine.World) {
	if s.continueShown {
		return
	}
	continueButton := Button{
		Render: component.NewRender(),
		Clickable: component.Clickable{
			Click: func(engine.Entity) {
				// pick the next room
				w.ChangeScene(&PickRoom{Scene: s.Scene})
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
	w.AddEntities(&continueButton)
	s.continueShown = true
}
