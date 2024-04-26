package scene

import (
	"fmt"
	"log/slog"
	"mob/pkg/component"
	"mob/pkg/font"
	"mob/pkg/system"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
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
	w.AddSystems(&system.UIGridLayout{}, &system.RenderSystem{}, &system.UIRenderText{}, &system.UIListLayout{})

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
			if x == 1 && y == 1 {
				textFace, err := font.GetTextFaceSource()
				if err != nil {
					slog.Error("could not render center text", "err", err)
					continue
				}
				ff := &text.GoTextFace{Source: textFace, Size: 14}
				// 3 texts in vertical list
				subgrid := TestSubContainer{
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
				w.AddEntities(&subgrid)
				for range 3 {
					_, txtH := text.Measure(fmt.Sprintf("row %d", x), ff, 0)
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
							H:      int(txtH),
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
						X:      x,
						Y:      y,
					},
				}
				w.AddEntities(&subsubgrid)
				for i := range 3 {
					txtW, _ := text.Measure(fmt.Sprintf("center %d", i), ff, 0)
					txt := TestLabel{
						Render: component.NewRender(),
						UILabel: component.UILabel{
							Text: []component.UILabelText{
								{Text: fmt.Sprintf("center %d", i), Color: colornames.White},
							},
						},
						UIChild: component.UIChild{
							Parent: "center",
							W:      int(txtW),
						},
					}
					w.AddEntities(&txt)
				}
			} else {
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
}
