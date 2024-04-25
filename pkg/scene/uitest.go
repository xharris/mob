package scene

import (
	"fmt"
	"mob/pkg/component"
	"mob/pkg/system"

	"github.com/sedyh/mizu/pkg/engine"
	"golang.org/x/exp/shiny/materialdesign/colornames"
)

type UITest struct{}

type TestContainer struct {
	component.Render
	component.UIGrid
}

type TestLabel struct {
	component.Render
	component.UILabel
	component.UIChild
}

func (u *UITest) Setup(w engine.World) {
	w.AddComponents(component.UIGrid{}, component.UILabel{}, component.UIChild{}, component.Render{})
	w.AddSystems(&system.UILayout{}, &system.RenderSystem{}, &system.UIRenderText{})

	grid := TestContainer{
		Render: component.NewRender(component.WGameSize()),
		UIGrid: component.UIGrid{
			ID:      "test",
			Rows:    3,
			Columns: 3,
		},
	}
	w.AddEntities(&grid)

	for x := range 3 {
		for y := range 3 {
			txt := TestLabel{
				Render: component.NewRender(),
				UILabel: component.UILabel{
					Text: []component.UILabelText{
						{Text: fmt.Sprintf("row %d", x), Color: colornames.White},
						{Text: fmt.Sprintf(" col %d", y), Color: colornames.White},
					},
				},
				UIChild: component.UIChild{
					Parent: "test",
					X:      x,
					Y:      y,
				},
			}
			w.AddEntities(&txt)
		}
	}
}
