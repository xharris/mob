package system

import (
	"log/slog"
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
		if !render.Visible {
			continue
		}
		for _, texture := range render.Textures {
			textureOptions := texture.RenderGeometry.GetOptions(render.RenderGeometry)
			// draw rect
			if render.Debug {
				slog.Info("render")
				tw, th := float32(texture.Image.Bounds().Dx()-1), float32(texture.Image.Bounds().Dy()-1)
				vector.StrokeRect(
					texture.Image,
					1, 1, tw-1, th-1,
					1, colornames.Green200, false,
				)
				vector.StrokeLine(
					texture.Image,
					1, 1, tw-1, th-1,
					1, colornames.Green200, false,
				)
			}
			screen.DrawImage(texture.Image, &textureOptions)
		}
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
