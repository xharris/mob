package test

import (
	"fmt"
	"math"
	"mob/pkg/pathfinding"
	"testing"

	"github.com/beefsack/go-astar"
)

type Tilemap [][]Tile

func (t *Tilemap) Get(x, y int) Tile {
	return (*t)[x][y]
}

type Tile struct {
	X, Y, Weight int
	Tiles        *Tilemap
}

func (c Tile) GetNeighbors() (neighbors []pathfinding.Node) {
	offsets := []int{-1, 0, 1}
	for x := range offsets {
		for y := range offsets {
			if x == 0 && y == 0 {
				continue
			}
			nx, ny := x+c.X, y+c.Y
			if c.Tiles != nil && nx >= 0 && ny >= 0 && nx < len(*c.Tiles) && ny < len((*c.Tiles)[nx]) {
				neighbors = append(neighbors, c.Tiles.Get(nx, ny))
			}
		}
	}
	return
}

func (c Tile) GetCost(to pathfinding.Node) int {
	return to.(Tile).Weight - c.Weight
}

func (c Tile) GetDistance(to pathfinding.Node) int {
	return int(math.Abs(float64(c.X-to.(Tile).X)) + math.Abs(float64(c.Y-to.(Tile).Y)))
}

func (c Tile) GetID() string {
	return fmt.Sprintf("%d,%d", c.X, c.Y)
}

func TestAStarFindPath(t *testing.T) {
	var stringMap []string = []string{
		"0000",
		"xxx0",
		"00x0",
		"0000",
	}
	var tilemap Tilemap = make(Tilemap, len(stringMap))
	for y, str := range stringMap {
		tilemap[y] = make([]Tile, len(str))
		for x, c := range str {
			tile := Tile{
				X:     x,
				Y:     y,
				Tiles: &tilemap,
			}
			if string(c) == "x" {
				tile.Weight = 10
			} else {
				tile.Weight = 0
			}
			tilemap[y][x] = tile
		}
	}
	start := Tile{X: 0, Y: 0, Weight: 0, Tiles: &tilemap}
	finish := Tile{X: 0, Y: 3, Weight: 0, Tiles: &tilemap}
	path := pathfinding.FindPath(start, finish)
	for _, p := range path {
		t.Logf("%d %d", p.(Tile).X, p.(Tile).Y)
	}
	if len(path) != 6 {
		t.Fatalf("path is incorrect size, expected=%d, got=%d", 6, len(path))
	}
}

type Node struct {
	Name string
	H    float64
	W    World
}

func (c *Node) PathNeighbors() (neighbors []astar.Pather) {
	for _, id := range c.W.Connections(c.GetID()) {
		neighbors = append(neighbors, c.W.GetNode(id))
	}
	return
}

func (c *Node) PathNeighborCost(to astar.Pather) float64 {
	return c.W.GetCost(c.GetID(), to.(*Node).GetID())
}

func (c *Node) PathEstimatedCost(to astar.Pather) float64 {
	return c.H
}

func (c Node) GetID() string {
	return c.Name
}

type World struct {
	connections map[string]map[string]float64
	nodes       []Node
}

func NewWorld() World {
	return World{
		connections: make(map[string]map[string]float64),
		nodes:       make([]Node, 0),
	}
}

func (w World) Connect(from string, to string, cost float64) {
	if _, ok := w.connections[from]; !ok {
		w.connections[from] = make(map[string]float64)
	}
	w.connections[from][to] = cost
}

func (w World) Connections(from string) (c []string) {
	if connections, ok := w.connections[from]; ok {
		for connection := range connections {
			c = append(c, connection)
		}
	}
	return
}

func (w World) GetCost(from string, to string) float64 {
	if connections, ok := w.connections[from]; ok {
		for _, cost := range connections {
			return cost
		}
	}
	return 99999
}

func (w *World) Add(node Node) {
	w.nodes = append(w.nodes, node)
}

func (w World) GetNode(id string) *Node {
	for _, node := range w.nodes {
		if node.GetID() == id {
			return &node
		}
	}
	return nil
}

func TestAStarFindPath2(t *testing.T) {
	w := NewWorld()
	w.Connect("a", "b", 4)
	w.Connect("a", "c", 3)
	w.Connect("b", "e", 12)
	w.Connect("b", "f", 5)
	w.Connect("c", "e", 10)
	w.Connect("c", "d", 7)
	w.Connect("d", "e", 2)
	w.Connect("e", "z", 5)
	nodes := []Node{
		{Name: "a", H: 14},
		{Name: "b", H: 12},
		{Name: "c", H: 11},
		{Name: "d", H: 6},
		{Name: "e", H: 4},
		{Name: "f", H: 11},
		{Name: "z", H: 0},
	}
	for _, node := range nodes {
		w.Add(node)
		node.W = w
	}
	start := nodes[0]
	goal := nodes[len(nodes)-1]

	path, _, _ := astar.Path(&start, &goal)
	for _, p := range path {
		t.Logf("%s", p.(*Node).Name)
	}
	if len(path) != 6 {
		t.Fatalf("path is incorrect size, expected=%d, got=%d", 6, len(path))
	}
}
