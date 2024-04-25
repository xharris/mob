package system

import (
	"bytes"
	"log/slog"
	"mob/pkg/component"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/sedyh/mizu/pkg/engine"
)

type TooltipSystem struct {
	*component.Tooltip
	*component.Render
}

var faceSource *text.GoTextFaceSource

func loadTextFaceSource() {
	data, err := os.ReadFile("asset/font/Retro Gaming.ttf")
	if err != nil {
		slog.Warn("Could not load font", "err", err)
		return
	}
	faceSource, err = text.NewGoTextFaceSource(bytes.NewReader(data))
	if err != nil {
		slog.Warn("Could not load font", "err", err)
	}
}

func (t *TooltipSystem) Update(w engine.World) {
	// show tooltip text on hover
	if t.Render.MouseEntered() {
		slog.Info("entered")
		t.Tooltip.Shown = true
	}
	if t.Render.MouseExited() {
		slog.Info("exited")
		t.Tooltip.Shown = false
	}
}

type RenderTooltips struct {
	*component.RenderTooltips
	*component.Render
}

func (r *RenderTooltips) Update(w engine.World) {
	if faceSource == nil {
		loadTextFaceSource()
		r.Render.Image = ebiten.NewImage(ebiten.WindowSize())
	}
	tooltips := w.View(component.Tooltip{})

	r.Render.Image.Clear()
	for _, entity := range tooltips.Filter() {
		var tt *component.Tooltip
		entity.Get(&tt)
		var w, h float64
		if tt.Shown {
			for _, txt := range tt.Text {
				var txtW, txtH float64
				// draw text
				if len(txt.Text) > 0 {
					op := &text.DrawOptions{}
					op.GeoM.Translate(w, h)
					op.ColorScale.ScaleWithColor(txt.Color)
					ff := &text.GoTextFace{Source: faceSource, Size: 14}
					txtW, txtH = text.Measure(txt.Text, ff, 0)
					text.Draw(r.Render.Image, txt.Text, ff, op)
					w += txtW
				}
				// new line
				if txt.Newline {
					w = 0
					h += txtH
				}
			}
		}
	}
}
