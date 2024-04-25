package system

import (
	"mob/pkg/component"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sedyh/mizu/pkg/engine"
	"golang.org/x/exp/shiny/materialdesign/colornames"
)

type RenderSystem struct{}

func (r *RenderSystem) Draw(w engine.World, screen *ebiten.Image) {
	view := w.View(component.Render{})
	// sort by Z
	rs := &renderSorter{entities: view.Filter()}
	sort.Sort(rs)
	// draw all
	for _, entity := range rs.entities {
		var render *component.Render
		entity.Get(&render)
		if render.Image == nil {
			continue
		}

		// REMOVE shop item
		var shopItem *component.ShopItem
		entity.Get(&shopItem)
		if shopItem != nil {
			vector.DrawFilledRect(render.Image, 0, 0, 16, 16, colornames.Blue500, false)
		}

		// draw to screen
		render.DrawImageOptions.GeoM.Reset()
		render.DrawImageOptions.GeoM.Translate(render.X, render.Y)
		screen.DrawImage(render.Image, &render.DrawImageOptions)
	}
}

type renderSorter struct {
	entities []engine.Entity
}

func (s *renderSorter) Len() int {
	return len(s.entities)
}

func (s *renderSorter) Swap(a, b int) {
	s.entities[a], s.entities[b] = s.entities[b], s.entities[a]
}

func (s *renderSorter) Less(a, b int) bool {
	var aRender, bRender *component.Render
	s.entities[a].Get(&aRender)
	s.entities[b].Get(&bRender)
	return aRender.Z < bRender.Z
}
