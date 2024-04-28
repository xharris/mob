package level

import (
	"fmt"
	"math/rand"
	"mob/pkg/component"
	"mob/pkg/mod"
)

type Level struct {
	W, H  int
	Rooms [][]Room
	Level int
}

type RoomType int

const (
	Empty RoomType = iota
	Start
	// move to next level
	Finish
	Fight
	// heal all allies
	Rest
	// free training for random % of allies
	Library
	// buy mods
	Shop
	// hire allies
	Recruit
	// fire an ally to pass or lose hp/gold
	Fire
	// fire a random ally to pass or lose hp/gold
	FireRandom
)

var goodRoomTypes []RoomType = []RoomType{Rest, Library, Shop, Recruit}

type Room struct {
	Enemies []component.Enemy
	// appears as 'unknown' to player
	Unknown   bool
	Completed bool
	Type      RoomType
}

func (r *RoomType) IsGood() bool {
	if *r == Start {
		return true
	}
	for _, t := range goodRoomTypes {
		if t == *r {
			return true
		}
	}
	return false
}

func (l *Level) generateRoom(x, y int) {
	// generate room
	if l.Rooms[x][y].Type == Empty {
		// TODO make it random
		room := Room{
			Type:    Fight,
			Enemies: []component.Enemy{},
		}
		for range rand.Intn(l.Level+2) + 1 {
			enemy := component.Enemy{
				Mods: []mod.Mod{
					{Name: "Slash", Type: mod.Attack, Target: mod.Ally},
				},
				Name: "goblin",
			}
			room.Enemies = append(room.Enemies, enemy)
		}
		l.Rooms[x][y] = room
	}
	// generate neighbors
	for nx := x - 1; nx <= x+1; nx += 1 {
		for ny := y - 1; ny <= y+1; ny += 1 {
			// within bounds?
			if nx >= 0 && ny >= 0 && nx < len(l.Rooms) && ny < len(l.Rooms[nx]) && l.Rooms[nx][ny].Type == Empty {
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
	l.Rooms = make([][]Room, l.W)
	for x := range l.W {
		l.Rooms[x] = make([]Room, l.H)
		for y := range l.H {
			room := Room{Type: Empty}
			if x == startX && y == startY {
				room.Type = Start
			}
			if x == finishX && y == finishY {
				room.Type = Finish
			}
			l.Rooms[x][y] = room
		}
	}
	// flood fill random room types from start
	l.generateRoom(startX, startY)
	for x := range l.W {
		for y := range l.H {
			fmt.Print(l.Rooms[x][y].Type)
		}
		fmt.Println()
	}
	return
}
