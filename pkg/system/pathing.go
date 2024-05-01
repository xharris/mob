package system

// import (
// 	"fmt"
// 	"log/slog"
// 	"math"
// 	"mob/pkg/component"
// 	"mob/pkg/priorityqueue"
// 	"slices"

// 	"github.com/sedyh/mizu/pkg/engine"
// )

// type Pathing struct{}

// const GridSize int = 32

// type Node struct {
// 	X, Y     int
// 	Obstacle component.PathObstacle
// 	G, H, F  float64
// }

// func (p *Node) Distance(to *Node) float64 {
// 	return float64(to.Obstacle.GetWeight() - p.Obstacle.GetWeight())
// }

// func (p *Node) SamePosition(other *Node) bool {
// 	return p.X == other.X && p.Y == other.Y
// }

// type Grid map[int]map[int]*Node

// func NewGrid() Grid {
// 	return make(Grid)
// }

// func (g Grid) Add(x, y float64, obstacle component.PathObstacle) *Node {
// 	gx, gy := int(math.Round(x/float64(GridSize))), int(math.Round(y/float64(GridSize)))
// 	if _, ok := g[gx]; !ok {
// 		g[gx] = make(map[int]*Node)
// 	}
// 	n := &Node{X: gx, Y: gy, Obstacle: component.PathObstacle{}}
// 	g[gx][gy] = n
// 	return n
// }

// // get obstacle at position
// func (g Grid) Get(x, y int) (n *Node, exists bool) {
// 	if _, exists = g[x]; exists {
// 		n, exists = g[x][y]
// 		return
// 	}
// 	return
// }

// func (g *Grid) FindPath(start, goal *Node) [][]float64 {
// 	path := [][]float64{
// 		{float64(goal.X * GridSize), float64(goal.Y * GridSize)},
// 	}
// 	return path

// 	open := priorityqueue.New[*Node]() // f
// 	closed := priorityqueue.New[*Node]()
// 	var parent map[*Node]*Node = make(map[*Node]*Node)

// 	open.Set(start, 0)
// 	for !open.Empty() {
// 		q := open.Pop()
// 		if q == nil {
// 			slog.Info("break early")
// 			break
// 		}
// 		// get successors
// 		var successors []*Node
// 		offsets := []int{-1, 0, 1}
// 		for x := range offsets {
// 			for y := range offsets {
// 				s, exists := g.Get(q.X+x, q.Y+y)

// 				if exists && s != q {
// 					parent[s] = q
// 					successors = append(successors, s)
// 				}
// 			}
// 		}
// 		slog.Info("q", "x", q.X, "y", q.Y, "succs", len(successors))
// 		for _, s := range successors {
// 			if s == goal {
// 				slog.Info("reached goal")
// 				goto done // stop search
// 			}
// 			s.G = q.G + s.Distance(q)
// 			s.H = s.Distance(goal)
// 			s.F = s.G + s.H

// 			// if a node with the same position as successor is in the OPEN list
// 			// which has a lower f than successor
// 			// skip this successor
// 			skip := false
// 			for _, o := range open.Values() {
// 				if o.SamePosition(s) && o.F < s.F {
// 					skip = true
// 					break
// 				}
// 			}
// 			if skip {
// 				slog.Info("skip open")
// 				continue
// 			}
// 			// if a node with the same position as successor is in the CLOSED list
// 			// which has a lower f than successor
// 			// skip this successor
// 			// otherwise, add  the node to the open list
// 			skip = false
// 			for _, c := range closed.Values() {
// 				if c.SamePosition(s) {
// 					if c.F < s.F {
// 						skip = true
// 						break
// 					} else {
// 						open.Set(s, int(s.F))
// 						closed.Remove(s)
// 					}
// 				}
// 			}
// 			if skip {
// 				slog.Info("skip closed")
// 				continue
// 			}
// 			open.Set(s, int(s.F))
// 		}
// 		closed.Set(q, 0)
// 	}
// done:
// 	var finalPath [][]float64
// 	next := parent[goal]
// 	for next != nil {
// 		finalPath = append(finalPath, []float64{float64(next.X * GridSize), float64(next.Y * GridSize)})
// 		next = parent[next]
// 	}
// 	slices.Reverse(finalPath)
// 	if len(finalPath) == 0 {
// 		slog.Info("no path")
// 	}
// 	return finalPath
// }

// func (*Pathing) Update(w engine.World) {
// 	obstacles := w.View(component.PathObstacle{}, component.Render{})
// 	agents := w.View(component.PathAgent{}, component.Render{})

// 	if len(agents.Filter()) <= 0 {
// 		return
// 	}

// 	for _, e := range agents.Filter() {
// 		var agent *component.PathAgent
// 		var render *component.Render
// 		var velocity *component.Velocity
// 		e.Get(&agent, &render, &velocity)

// 		if agent.NeedNewPath() {
// 			grid := NewGrid()
// 			// create empty grid
// 			for x := 0; x < w.Bounds().Dx(); x += GridSize {
// 				for y := 0; y < w.Bounds().Dy(); y += GridSize {
// 					grid.Add(float64(x), float64(y), component.PathObstacle{})
// 				}
// 			}
// 			// add obstacles
// 			for _, e := range obstacles.Filter() {
// 				var obs *component.PathObstacle
// 				var render *component.Render
// 				e.Get(&obs, &render)
// 				grid.Add(render.X, render.Y, *obs)
// 			}

// 			start := grid.Add(render.X, render.Y, component.PathObstacle{})
// 			target := grid.Add(agent.TargetX, agent.TargetY, component.PathObstacle{})

// 			for x := range w.Bounds().Dx() / GridSize {
// 				for y := range w.Bounds().Dy() / GridSize {
// 					n, ok := grid.Get(x, y)
// 					if !ok {
// 						fmt.Print("-")
// 						continue
// 					}
// 					if start.SamePosition(n) {
// 						fmt.Print("S")
// 					} else if target.SamePosition(n) {
// 						fmt.Print("T")
// 					} else {
// 						fmt.Print("0")
// 					}
// 				}
// 				fmt.Println()
// 			}

// 			aPath := grid.FindPath(start, target)
// 			agent.SetPath(aPath)
// 		}

// 		// move towards next poin in path
// 		if targetX, targetY, done := agent.GetNextPoint(); velocity != nil && !done {
// 			dist := math.Sqrt(math.Pow(targetX-render.X, 2) + math.Pow(targetY-render.Y, 2))
// 			velocity.X = (targetX - render.X) / dist * 5
// 			velocity.Y = (targetY - render.Y) / dist * 5
// 		}
// 	}
// }
