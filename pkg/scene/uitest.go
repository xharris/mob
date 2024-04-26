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

type TestSubContainer struct {
	component.Render
	component.UIList
	component.UIChild
}

type TestLabel struct {
	component.Render
	component.UILabel
	component.UIChild
}

func (u *UITest) Setup(w engine.World) {
	w.AddComponents(component.UIGrid{}, component.UILabel{}, component.UIChild{}, component.Render{}, component.UIList{})
	w.AddSystems(&system.UIGridLayout{}, &system.RenderSystem{}, &system.UIRenderLabel{}, &system.UIListLayout{})

	grid := TestContainer{
		Render: component.NewRender(component.WGameSize()),
		UIGrid: component.UIGrid{
			ID:      "test",
			Rows:    3,
			Columns: 3,
		},
	}
	w.AddEntities(&grid)

	aligns := []component.UIAlign{
		component.START, component.CENTER, component.END,
	}

	for x := range 3 {
		for y := range 3 {
			if x == 1 && y == 1 {
				// 3 texts in vertical list
				sublist := TestSubContainer{
					Render: component.NewRender(),
					UIList: component.UIList{
						ID:        "center",
						Direction: component.VERTICAL,
					},
					UIChild: component.UIChild{
						Parent: "test",
						X:      x,
						Y:      y,
					},
				}
				w.AddEntities(&sublist)
				for range 3 {
					txt := TestLabel{
						Render: component.NewRender(),
						UILabel: component.UILabel{
							Text: []component.UILabelText{
								{Text: fmt.Sprintf("row %d", x), Color: colornames.White},
								{Text: fmt.Sprintf(" col %d", y), Color: colornames.White},
							},
						},
						UIChild: component.UIChild{
							Parent: "center",
						},
					}
					w.AddEntities(&txt)
				}
				// 3 more text in horizontal list
				subsubgrid := TestSubContainer{
					Render: component.NewRender(),
					UIList: component.UIList{
						ID:        "center-footer",
						Direction: component.HORIZONTAL,
					},
					UIChild: component.UIChild{
						Parent: "center",
					},
				}
				w.AddEntities(&subsubgrid)
				for i := range 3 {
					txt := TestLabel{
						Render: component.NewRender(),
						UILabel: component.UILabel{
							Text: []component.UILabelText{
								{Text: fmt.Sprintf("center %d", i), Color: colornames.White},
							},
						},
						UIChild: component.UIChild{
							Parent: "center-footer",
						},
					}
					w.AddEntities(&txt)
				}
			} else {
				txtContainerID := component.UI_ID(fmt.Sprintf("text-container-%d-%d", x, y))
				txtContainer := TestSubContainer{
					Render: component.NewRender(component.WRenderDebug()),
					UIList: component.UIList{
						ID:      txtContainerID,
						Align:   aligns[x],
						Justify: aligns[y],
					},
					UIChild: component.UIChild{
						Parent: "test",
						X:      x,
						Y:      y,
					},
				}
				txt := TestLabel{
					Render: component.NewRender(component.WRenderDebug()),
					UILabel: component.UILabel{
						Text: []component.UILabelText{
							{Text: fmt.Sprintf("rowfeafew %d", x), Color: colornames.White},
							{Newline: true},
							{Text: fmt.Sprintf("col %d", y), Color: colornames.White},
						},
					},
					UIChild: component.UIChild{
						Parent: txtContainerID,
					},
				}
				w.AddEntities(&txtContainer, &txt)
			}
		}
	}
}
