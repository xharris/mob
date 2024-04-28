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
	component.Room
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
			if (x == y) || (nx == ny) || c.State.Level.IsOutsideBounds(nx, ny) {
				continue
			}
			room := c.State.Level.Rooms[ny][nx]
			slog.Info("good", "room", lang.Get(fmt.Sprintf("room%d", room.Type)))
			roomChoice := RoomChoice{
				UIChild: component.UIChild{
					Parent: "pick-grid",
					X:      x,
					Y:      y,
				},
				UILabel: component.UILabel{
					Text: []component.UILabelText{
						{Text: lang.Get(fmt.Sprintf("room%d", room.Type)), Color: colornames.White},
					},
				},
				Render: component.NewRender(),
				Clickable: component.Clickable{
					Click: func(e engine.Entity) {
						// TODO go to room
						var cRoom *component.Room
						e.Get(&cRoom)
						switch cRoom.Type {
						case component.Fight:
							w.ChangeScene(&Fight{
								Scene: c.Scene,
								Room:  *cRoom,
							})
						default:
							slog.Warn("room not implemented yet", "type", lang.Get(fmt.Sprintf("room%d", cRoom.Type)))
						}
					},
				},
				Hoverable: component.Hoverable{},
				Room:      room,
			}
			w.AddEntities(&roomChoice)
		}
	}
}
