package system

import (
	"mob/pkg/component"
	"mob/pkg/font"
	"slices"

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
	v := w.View(component.UIChild{})

	var x, y float64
	items := v.Filter()
	if ui.UIList.Reverse {
		slices.Reverse(items)
	}
	for _, e := range items {
		var child *component.UIChild
		var render *component.Render
		e.Get(&child, &render)
		// belongs to this list?
		if child.Parent != ui.UIList.ID {
			continue
		}
		ui.Render.Fit(render)
		childW, childH := render.GetSize()
		// arrange nodes in list
		render.X = ui.Render.X + x
		render.Y = ui.Render.Y + y

		if ui.UIList.Direction == component.HORIZONTAL {
			x += float64(childW)
			// alignment
			if ui.UIList.Justify == component.CENTER {
				render.X += float64(ui.Render.Image.Bounds().Dx())/2 - float64(childW)/2
			}
			if ui.UIList.Justify == component.END {
				render.X += float64(ui.Render.Image.Bounds().Dx()) - float64(childW)
			}
			if ui.UIList.Align == component.CENTER {
				render.Y += float64(ui.Render.Image.Bounds().Dy())/2 - float64(childH)/2
			}
			if ui.UIList.Align == component.END {
				render.Y += float64(ui.Render.Image.Bounds().Dy()) - float64(childH)
			}
		}

		if ui.UIList.Direction == component.VERTICAL {
			y += float64(childH)
			// alignment
			if ui.UIList.Align == component.CENTER {
				render.X += float64(ui.Render.Image.Bounds().Dx())/2 - float64(childW)/2
			}
			if ui.UIList.Align == component.END {
				render.X += float64(ui.Render.Image.Bounds().Dx()) - float64(childW)
			}
			if ui.UIList.Justify == component.CENTER {
				render.Y += float64(ui.Render.Image.Bounds().Dy())/2 - float64(childH)/2
			}
			if ui.UIList.Justify == component.END {
				render.Y += float64(ui.Render.Image.Bounds().Dy()) - float64(childH)
			}
		}
	}
}

type UIRenderLabel struct {
	*component.Render
	*component.UILabel
}

func (t *UIRenderLabel) Update(world engine.World) {
	// size := t.Render.Image.Bounds().Max

	var x, y float64
	var totalW, totalH float64
	t.Render.Image.Clear()
	f := font.DefaultFont
	if t.UILabel.Font != nil {
		f = t.UILabel.Font
	}
	var lineW float64
	var lineH float64
	for i, txt := range t.UILabel.Text {
		if txt.Font != nil {
			f = txt.Font
		}
		ff := f.Face()
		if txt.Newline {
			// new line
			x = 0
			if lineH == 0 {
				_, lineH = f.Measure("|")
			}
			totalH += lineH
			y += lineH
			lineW = 0
			lineH = 0
		} else if len(txt.Text) > 0 {
			// draw text
			for c, char := range txt.Text {
				// draw options
				op := &text.DrawOptions{}
				op.GeoM.Translate(x, y)
				op.LayoutOptions.PrimaryAlign = t.UILabel.HAlign
				op.LayoutOptions.SecondaryAlign = t.UILabel.VAlign
				op.ColorScale.ScaleWithColor(txt.Color)
				// draw
				text.Draw(t.Render.Image, string(char), &ff, op)
				// measure
				txtW, txtH := f.Measure(string(char))
				// space between chars
				if char != ' ' {
					txtW -= 1
				}
				x += txtW
				// adjust total line height
				if txtH > lineH {
					lineH = txtH
				}
				lineW += txtW
				// last line, last char
				if c == len(txt.Text)-1 && i == len(t.UILabel.Text)-1 {
					totalH += lineH
				}
			}
			if lineW > totalW {
				totalW = lineW
			}
		}
	}
	t.Render.Resize(int(totalW+2), int(totalH)) // TODO why do I have to add 2?
}
