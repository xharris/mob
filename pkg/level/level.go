package level

import (
	"fmt"
	"log/slog"
	"math/rand"
	"mob/pkg/component"
)

type Level struct {
	W, H  int
	Rooms [][]component.Room
	Level int
}

func (l *Level) GetFinish() (int, int, bool) {
	for lx := range l.W {
		for ly := range l.H {
			if l.Rooms[lx][ly].Type == component.Finish {
				return lx, ly, true
			}
		}
	}
	return -1, -1, false
}

func (l *Level) generateRoom(x, y int) {
	// generate room
	if l.Rooms[y][x].Type == component.Empty {
		// TODO make it random
		room := component.Room{
			Type:    component.Fight,
			Enemies: []component.NPC{},
		}
		for range rand.Intn(l.Level+2) + 1 {
			enemy := component.NewNPC(
				component.WithMod(component.Mod{
					Name:   "Slash",
					Type:   component.Attack,
					Target: component.TargetAlly,
					Range:  component.Melee,
				}),
				component.WithName("goblin"),
				component.WithType(component.Enemy),
			)
			room.Enemies = append(room.Enemies, enemy)
		}
		l.Rooms[y][x] = room
	}
	// generate neighbors
	for nx := x - 1; nx <= x+1; nx += 1 {
		for ny := y - 1; ny <= y+1; ny += 1 {
			// within bounds?
			neighbor := l.GetRoom(nx, ny)
			if neighbor != nil && neighbor.Type == component.Empty {
				l.generateRoom(nx, ny)
			}
		}
	}
}

func Generate(level int) (l Level, startX, startY int) {
	l.Level = level
	// level size
	size := max(3, level)
	randOffset := rand.Intn(size)
	signs := []int{-1, 1}
	randSign := signs[rand.Intn(len(signs))]
	l.W = level + 2 + (randOffset * randSign)
	l.H = level + 2 + (randOffset * randSign * -1)
	// starting room
	cornerX := []int{0, l.W - 1}
	cornerY := []int{0, l.H - 1}
	startX = cornerX[rand.Intn(2)]
	startY = cornerY[rand.Intn(2)]
	// ending room
	finishX := cornerX[0]
	finishY := cornerY[0]
	if finishX == startX {
		finishX = cornerX[1]
	}
	if finishY == startY {
		finishY = cornerY[1]
	}
	// fill rooms array and set starting room
	l.Rooms = make([][]component.Room, l.H)
	for y := range l.H {
		l.Rooms[y] = make([]component.Room, l.W)
		for x := range l.W {
			room := component.Room{X: y, Y: x, Type: component.Empty}
			if x == startX && y == startY {
				room.Type = component.Start
			}
			if x == finishX && y == finishY {
				room.Type = component.Finish
			}
			l.Rooms[y][x] = room
		}
	}
	// flood fill random room types from start
	l.generateRoom(startX, startY)
	slog.Info("generated level", "w", l.W, "h", l.H)

	for y := range l.H {
		for x := range l.W {
			room := l.GetRoom(x, y)
			if room == nil {
				fmt.Print(' ')
			} else {
				fmt.Print(room.Type)
			}
		}
		fmt.Println()
	}
	return
}

func (l *Level) IsOutsideBounds(x, y int) bool {
	return x < 0 || y < 0 || y >= l.H || x >= l.W
}

func (l *Level) GetRoom(x, y int) *component.Room {
	if !l.IsOutsideBounds(x, y) {
		return &l.Rooms[y][x]
	}
	return nil
}

func (l *Level) IterRooms() (rooms []*component.Room) {
	for y := range l.H {
		for x := range l.W {
			rooms = append(rooms, &l.Rooms[y][x])
		}
	}
	return
}
