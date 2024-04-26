package system

import (
	"log/slog"
	"mob/pkg/component"
	"mob/pkg/font"

	"github.com/sedyh/mizu/pkg/engine"
)

type TooltipSystem struct{}

type TooltipLabel struct {
	component.Render
	component.UILabel
	component.UIChild
}

func (*TooltipSystem) Update(w engine.World) {
	view := w.View(component.Tooltip{}, component.Render{})
	f := font.DefaultFont
	for _, t := range view.Filter() {
		var r *component.Render
		var tt *component.Tooltip
		t.Get(&r, &tt)
		// show tooltip text on hover
		if r.MouseEntered() {
			slog.Info("entered")
			tt.Shown = true
			// show ui label
			var height float64
			for _, text := range tt.Text {
				if text.Font != nil {
					f = text.Font
				}
				_, h := f.Measure(text.Text)
				height += h
			}
			label := TooltipLabel{
				Render: component.NewRender(),
				UILabel: component.UILabel{
					Text: tt.Text,
				},
				UIChild: component.UIChild{
					Parent: tt.Parent,
					H:      int(height),
				},
			}
			w.AddEntities(&label)
		}
		if r.MouseExited() {
			slog.Info("exited")
			tt.Shown = false
			// remove ui label
			labels := w.View(component.UILabel{}, component.UIChild{}, component.Render{})
			labels.Each(func(e engine.Entity) {
				var ch *component.UIChild
				e.Get(&ch)
				if ch != nil && ch.Parent == tt.Parent {
					w.RemoveEntity(e)
				}
			})
		}
	}
}
