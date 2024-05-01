package system

import (
	"mob/pkg/component"
	"mob/pkg/font"
	"slices"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/sedyh/mizu/pkg/engine"
	"golang.org/x/exp/shiny/materialdesign/colornames"
)

type UIGridLayout struct {
	*component.UIGrid
	*component.Render
}

func (ui *UIGridLayout) Update(w engine.World) {
	uiW, uiH := ui.Render.GetSize()

	var columns, rows int

	items := w.View(component.UIChild{})
	var children []engine.Entity

	// belongs to this grid?
	for _, e := range items.Filter() {
		var child *component.UIChild
		e.Get(&child)
		if child.Parent == ui.UIGrid.ID {
			children = append(children, e)
		}
	}

	// calculate rows if only columns is given
	columns = max(1, ui.UIGrid.Columns)
	autoArrange := ui.UIGrid.Rows == 0
	if autoArrange {
		rows = max(1, len(children)/columns)
	} else {
		rows = ui.UIGrid.Rows
	}

	cellW := float64(uiW / columns)
	cellH := float64(uiH / rows)

	for i, e := range children {
		var child *component.UIChild
		var render *component.Render
		var list *component.UIList
		e.Get(&child, &render, &list)
		if list != nil && !list.FitContents {
			render.Resize(int(cellW), int(cellH))
		}
		x, y := child.X, child.Y
		if autoArrange {
			x = int(i % columns)
			y = int(i / columns)
		}
		render.X = ui.Render.X + float64(x)*cellW
		render.Y = ui.Render.Y + float64(y)*cellH
		// alignment
		dx, dy := render.GetSize()
		if ui.UIGrid.Justify == component.CENTER {
			render.X += float64(cellW)/2 - float64(dx)/2
		}
		if ui.UIGrid.Justify == component.END {
			render.X += float64(cellW) - float64(dx)
		}
		if ui.UIGrid.Align == component.CENTER {
			render.Y += float64(cellH)/2 - float64(dy)/2
		}
		if ui.UIGrid.Align == component.END {
			render.Y += float64(cellH) - float64(dy)
		}
		// TODO use render.GeoM.Apply instead? or just set these values to 0?
		render.X += render.OX
		render.Y += render.OY
	}
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
	var totalW, totalH int
	for _, e := range items {
		var child *component.UIChild
		var render *component.Render
		var list *component.UIList
		e.Get(&child, &render, &list)
		// belongs to this list?
		if child.Parent != ui.UIList.ID {
			continue
		}
		childW, childH := render.GetSize()
		totalW += childW
		totalH += childH
		// arrange nodes in list
		render.X = ui.Render.X + x
		render.Y = ui.Render.Y + y

		dx, dy := float64(ui.Render.Image.Bounds().Dx()), float64(ui.Render.Image.Bounds().Dy())
		if ui.UIList.Direction == component.HORIZONTAL {
			x += float64(childW)
			// alignment
			if ui.UIList.Justify == component.CENTER {
				render.X += dx/2 - float64(childW)/2
			}
			if ui.UIList.Justify == component.END {
				render.X += dx - float64(childW)
			}
			if ui.UIList.Align == component.CENTER {
				render.Y += dy/2 - float64(childH)/2
			}
			if ui.UIList.Align == component.END {
				render.Y += dy - float64(childH)
			}
		}

		if ui.UIList.Direction == component.VERTICAL {
			y += float64(childH)
			// alignment
			if ui.UIList.Align == component.CENTER {
				render.X += dx/2 - float64(childW)/2
			}
			if ui.UIList.Align == component.END {
				render.X += dx - float64(childW)
			}
			if ui.UIList.Justify == component.CENTER {
				render.Y += dy/2 - float64(childH)/2
			}
			if ui.UIList.Justify == component.END {
				render.Y += dy - float64(childH)
			}
		}
		// TODO use render.GeoM.Apply instead? or just set these values to 0?
		render.X += render.OX
		render.Y += render.OY
	}
	// resize to fit contents
	if ui.UIList.FitContents {
		ui.Render.Resize(totalW, totalH)
	}
}

type UIRenderLabel struct {
	*component.Render
	*component.UILabel
}

func (t *UIRenderLabel) Update(world engine.World) {
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
		} else {
			// draw text
			for c, char := range txt.Text {
				// draw options
				op := &text.DrawOptions{}
				op.GeoM.Translate(x, y)
				op.LayoutOptions.PrimaryAlign = t.UILabel.HAlign
				op.LayoutOptions.SecondaryAlign = t.UILabel.VAlign
				color := txt.Color
				if color == nil {
					color = colornames.White
				}
				op.ColorScale.ScaleWithColor(color)
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

type UIRenderChild struct {
	*component.Render
	*component.UIChild
}

func (c *UIRenderChild) Update(w engine.World) {

}
