package scene

import (
	"fmt"
	"mob/pkg/allymod"
	"mob/pkg/component"
	"mob/pkg/game"
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
	component.Health
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
	AllyModCount   int
	MustBuyOne     bool
	continueShown  bool
	purchasedCount int
}

func (s *Shop) Setup(w engine.World) {
	w.AddComponents(
		component.Render{}, component.ShopItem{}, component.Rect{}, component.Health{},
		component.UIList{}, component.UILabel{}, component.UIChild{}, component.Clickable{},
		component.UIGrid{}, component.Hoverable{},
	)
	w.AddSystems(
		&system.RenderSystem{}, &system.ShopItem{},
		&system.RenderRect{}, &system.UIRenderLabel{}, &system.UIListLayout{},
		&system.Clickable{}, &system.UIGridLayout{}, &system.Hoverable{},
	)

	if s.AllyModCount <= 0 {
		s.AllyModCount = 3
	}

	shopItemsUI := UI{
		Render: component.NewRender(),
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
		shoptItemTooltipID := component.UI_ID(fmt.Sprintf("shop-item-tooltip-%d", i))
		shopItemID := component.UI_ID(fmt.Sprintf("shop-item-%d", i))
		shopitem := ShopAllyContainer{
			Render: component.NewRender(),
			ShopItem: component.ShopItem{
				AddMods: []allymod.Mod{
					{Name: "Slash", Type: allymod.Attack, Target: allymod.Enemy},
					{Name: "Block", Type: allymod.Buff, Target: allymod.Self},
					{Name: "Sleepy", Desc: "Might take a nap", Type: allymod.Debuff, Target: allymod.Self},
				},
				Name: "Sir Bobbington",
			},
			Health: component.Health{
				Total:     100,
				Remaining: 100,
			},
			Clickable: component.Clickable{
				Click: func(e engine.Entity) {
					var shopItem *component.ShopItem
					var health *component.Health
					e.Get(&shopItem, &health)
					s.purchaseAlly(w, shopItem, health)
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
				Color: colornames.Blue500,
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
	if s.purchasedCount > 0 || !s.MustBuyOne {
		s.showContinueButton(w)
	}
	w.AddEntities(&actions)
}

func (s *Shop) purchaseAlly(w engine.World, item *component.ShopItem, health *component.Health) {
	if s.purchasedCount > 0 || !s.MustBuyOne {
		s.showContinueButton(w)
	}
	var foundAlly *game.Ally
	for _, a := range s.Scene.State.Allies {
		if a.Name == item.Name {
			foundAlly = &a
			break
		}
	}
	if foundAlly == nil {
		// add new ally
		ally := game.Ally{
			Ally: component.Ally{
				Mods: item.AddMods,
			},
			Health: *health,
		}
		s.Scene.State.Allies = append(s.Scene.State.Allies, ally)
	} else {
		// change existing ally
		// mods
		for _, mod := range item.AddMods {
			foundAlly.AddMod(mod)
		}
		for _, mod := range item.RemoveMods {
			foundAlly.RemoveMod(mod)
		}
		// health
		foundAlly.Health = *health
	}
	item.Purchased = true
}

func (s *Shop) showContinueButton(w engine.World) {
	if s.continueShown {
		return
	}
	continueButton := Button{
		Render: component.NewRender(),
		Clickable: component.Clickable{
			Click: func(engine.Entity) {
				// go to strategy scene
				w.ChangeScene(&Strategy{
					Scene: s.Scene,
				})
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
