package system

import (
	"mob/pkg/component"
	"mob/pkg/font"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/sedyh/mizu/pkg/engine"
)

type UIGridLayout struct {
	*component.UIGrid
	*component.Render
}

func (ui *UIGridLayout) Update(w engine.World) {
	uiW, uiH := ui.Render.GetSize()
	var (
		cellW float64 = float64(uiW / ui.UIGrid.Columns)
		cellH float64 = float64(uiH / ui.UIGrid.Rows)
	)

	items := w.View(component.UIChild{})

	i := 0
	items.Each(func(e engine.Entity) {
		var child *component.UIChild
		var render *component.Render
		e.Get(&child, &render)
		// belongs to this grid?
		if child.Parent != ui.UIGrid.ID {
			return
		}
		// arrange nodes in grid
		render.Resize(int(cellW), int(cellH))
		render.X = ui.Render.X + float64(child.X)*cellW
		render.Y = ui.Render.Y + float64(child.Y)*cellH
		// REMOVE? auto-arrange
		// render.X = float64(i/ui.UIGrid.Columns) * cellW
		// render.Y = float64(i%ui.UIGrid.Columns) * cellH

		i++
	})
}

type UIListLayout struct {
	*component.UIList
	*component.Render
}

func (ui *UIListLayout) Update(w engine.World) {
	items := w.View(component.UIChild{})
	size := ui.Render.Image.Bounds().Max

	i := 0
	var x, y float64
	items.Each(func(e engine.Entity) {
		var child *component.UIChild
		var render *component.Render
		e.Get(&child, &render)
		// belongs to this grid?
		if child.Parent != ui.UIList.ID {
			return
		}
		// arrange nodes in list
		render.Resize(int(child.W), int(child.H))
		render.X = ui.Render.X + x
		if ui.UIList.Reverse && ui.UIList.Direction == component.VERTICAL {
			y += float64(child.H)
			render.Y = ui.Render.Y + float64(size.Y) - y
		} else {
			render.Y = ui.Render.Y + y
			y += float64(child.H)
		}
		x += float64(child.W)
		i++
	})
}

type UIRenderLabel struct {
	*component.Render
	*component.UILabel
}

func (t *UIRenderLabel) Update(world engine.World) {
	size := t.Render.Image.Bounds().Max

	var x, y float64
	var txtW, txtH float64
	t.Render.Image.Clear()
	f := font.DefaultFont
	if t.UILabel.Font != nil {
		f = t.UILabel.Font
	}
	for _, txt := range t.UILabel.Text {
		if txt.Font != nil {
			f = txt.Font
		}
		ff := f.Face()
		// draw text
		if len(txt.Text) > 0 {
			for _, char := range txt.Text {
				// draw options
				op := &text.DrawOptions{}
				op.GeoM.Translate(x, y)
				op.LayoutOptions.PrimaryAlign = t.UILabel.VAlign
				op.LayoutOptions.SecondaryAlign = t.UILabel.VAlign
				op.ColorScale.ScaleWithColor(txt.Color)
				// draw
				text.Draw(t.Render.Image, string(char), &ff, op)
				// measure
				txtW, txtH = text.Measure(string(char), &ff, 0)
				// space between chars
				if char == ' ' {
					x += txtW
				} else {
					x += txtW - 1
				}
				// next line?
				if x > float64(size.X) {
					x = 0
				}
			}
		}
		// new line
		if txt.Newline {
			x = 0
			y += txtH
		}
	}
}
