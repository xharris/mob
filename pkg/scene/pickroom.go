package scene

import (
	"fmt"
	"log/slog"
	"mob/pkg/component"
	"mob/pkg/lang"
	"mob/pkg/system"

	"github.com/sedyh/mizu/pkg/engine"
	"golang.org/x/exp/shiny/materialdesign/colornames"
)

type PickGrid struct {
	component.UIGrid
	component.Render
}

type RoomChoice struct {
	component.Render
	component.UIChild
	component.UILabel
	component.Clickable
	component.Hoverable
}

type PickRoom struct {
	Scene
}

func (c *PickRoom) Setup(w engine.World) {
	slog.Info("pick a room")
	component.AddComponents(w)
	system.AddSystems(w)

	pickGrid := PickGrid{
		Render: component.NewRender(),
		UIGrid: component.UIGrid{
			ID:      "pick-grid",
			Rows:    3,
			Columns: 3,
			Align:   component.CENTER,
			Justify: component.CENTER,
		},
	}
	w.AddEntities(&pickGrid)

	lvlX, lvlY := c.State.LevelX, c.State.LevelY
	for x := range 3 {
		for y := range 3 {
			nx, ny := (x-1)+lvlX, (y-1)+lvlY
			slog.Info("room", "x", x, "y", y, "nx", nx, "ny", ny)
			if ((x-1) == 0 && (y-1) == 0) || (x == y) || nx < 0 || ny < 0 || nx >= c.State.Level.W || ny >= c.State.Level.H {
				continue
			}
			room := c.State.Level.Rooms[nx][ny]
			slog.Info("good", "room", lang.Get(fmt.Sprintf("room%d", room.Type)))
			roomChoice := RoomChoice{
				UIChild: component.UIChild{
					Parent: "pick-grid",
					X:      y,
					Y:      x,
				},
				UILabel: component.UILabel{
					Text: []component.UILabelText{
						{Text: lang.Get(fmt.Sprintf("room%d", room.Type)), Color: colornames.White},
					},
				},
				Render: component.NewRender(),
				Clickable: component.Clickable{
					Click: func(e engine.Entity) {
						// go to room
					},
				},
				Hoverable: component.Hoverable{},
			}
			w.AddEntities(&roomChoice)
		}
	}
}
