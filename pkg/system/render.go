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
		if !render.Visible || render.Image == nil {
			continue
		}
		op := render.DrawImageOptions
		op.GeoM.Reset()
		op.GeoM.Translate(render.X, render.Y)
		// alpha
		op.ColorScale.Reset()
		op.ColorScale.ScaleAlpha(float32(render.AlphaLevel) / float32(component.AlphaFull))
		render.AlphaLevel = component.AlphaFull
		rw, rh := render.GetSize()
		// draw rect
		if render.Debug {
			vector.StrokeRect(
				render.Image,
				1, 1, float32(rw-1), float32(rh-1),
				1, colornames.Green200, false,
			)
			vector.StrokeLine(
				render.Image,
				0, 0, float32(rw), float32(rh),
				1, colornames.Green200, false,
			)
		}
		screen.DrawImage(render.Image, op)
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
