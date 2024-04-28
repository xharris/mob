package scene

import (
	"log/slog"
	"mob/pkg/component"
	"mob/pkg/system"

	"github.com/sedyh/mizu/pkg/engine"
	"golang.org/x/exp/shiny/materialdesign/colornames"
)

type Strategy struct {
	Scene
}

func (s *Strategy) Setup(w engine.World) {
	slog.Info("strategy")
	component.AddComponents(w)
	system.AddSystems(w)

	// list allies (click to change priority, view mods, etc)
	// for _, ally := range s.State.Allies {

	// }

	// continue button
	actions := List{
		Render: component.NewRender(),
		UIList: component.UIList{
			ID:        "actions",
			Direction: component.VERTICAL,
			Align:     component.END,
			Justify:   component.END,
		},
		UIChild: component.UIChild{
			Parent: "ui",
			X:      2,
			Y:      2,
		},
	}
	w.AddEntities(&actions)
	continueButton := Button{
		Render: component.NewRender(),
		Clickable: component.Clickable{
			Click: func(engine.Entity) {
				// go to fight scene
				w.ChangeScene(&Fight{
					Scene: s.Scene,
				})
			},
		},
		UILabel: component.UILabel{
			Text: []component.UILabelText{
				{Text: "Continue", Color: colornames.Red300},
			},
		},
		UIChild: component.UIChild{
			Parent: "actions",
		},
	}
	w.AddEntities(&continueButton)
}
