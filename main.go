package main

import (
	"log"
	"mob/pkg/component"
	"mob/pkg/pawn"
	"mob/pkg/system"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sedyh/mizu/pkg/engine"
)

type Game struct{}

type ShopPawnItem struct {
	component.Render
	component.ShopItem
	component.Tooltip
}

type MainTooltipArea struct {
	component.RenderTooltips
	component.Render
}

func (g *Game) Setup(w engine.World) {
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

	gw, gh := ebiten.WindowSize()
	w.AddEntities(&MainTooltipArea{
		Render:         component.NewRender(float64(gw), float64(gh)),
		RenderTooltips: component.RenderTooltips{},
	})
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

}

func (g *Game) Layout(w, h int) (sw, sh int) {
	return 600, 400
}

func main() {
	ebiten.SetWindowSize(600, 400)
	ebiten.SetWindowTitle("pawns in dungeon game")

	if err := ebiten.RunGame(engine.NewGame(&Game{})); err != nil {
		log.Fatal(err)
	}
}
