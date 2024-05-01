package pathfinding

import (
	"log/slog"
	"mob/pkg/priorityqueue"
	"slices"
)

type Grid map[int]map[int]Node

type Node interface {
	GetNeighbors() []Node
	// g, distance from start
	GetCost(to Node) int
	// h, distance from goal
	GetDistance(to Node) int
	GetID() string
}

type Cost struct {
	F, G, H int
}

func NewCost(g, h int) (c Cost) {
	c.G = g
	c.H = h
	c.F = g + h
	return
}

func (c *Cost) SetH(h int) {
	c.H = h
	c.F = c.G + c.H
}

func (c *Cost) SetG(g int) {
	c.G = g
	c.F = c.G + c.H
}

func FindPath(start Node, goal Node) (fullPath []Node) {
	// var parent map[string]Node = make(map[string]Node)
	// distFromStart := priorityqueue.New[Node]() // g
	// totalDist := priorityqueue.New[Node]()     // f
	// closed := priorityqueue.New[Node]()

	// distFromStart.Set(from, 0)
	// fromCost := NewCost(0, from.GetDistance(to))
	// totalDist.Set(from, fromCost.F)

	// for !totalDist.Empty() {
	// 	current := totalDist.Pop()
	// 	if current.GetID() == to.GetID() {
	// 		break
	// 	}
	// 	for _, neighbor := range current.GetNeighbors() {
	// 		if closed.Has(neighbor) {
	// 			continue
	// 		}
	// 		parent[neighbor.GetID()] = current
	// 		cost := NewCost(neighbor.GetCost(current), neighbor.GetDistance(to))
	// 		// add previous g
	// 		if _, g, found := distFromStart.Get(neighbor); found {
	// 			cost.SetG(cost.G + g)
	// 		}
	// 		distFromStart.Set(neighbor, cost.G)
	// 	}
	// 	closed.Set(current, 0)
	// }

	frontier := priorityqueue.New[Node]()
	frontier.Set(start, 0)
	var cameFrom map[string]Node = make(map[string]Node)
	var costSoFar map[string]int = make(map[string]int) // g

	costSoFar[start.GetID()] = 0

	for !frontier.Empty() {
		current := frontier.Pop()
		if current.GetID() == goal.GetID() {
			break
		}
		for _, next := range current.GetNeighbors() {
			slog.Info("neighbors", "current", current.GetID(), "next", next.GetID())
			// new cost
			maybeG := current.GetCost(next)
			if prevG, hasPrevG := costSoFar[current.GetID()]; hasPrevG {
				maybeG += prevG
			}
			if gNext, hasNextG := costSoFar[next.GetID()]; !hasNextG || maybeG < gNext {
				costSoFar[next.GetID()] = maybeG
				cost := NewCost(maybeG, next.GetDistance(goal))
				frontier.Set(next, cost.F)
				cameFrom[next.GetID()] = current
			}
		}
	}

	slog.Info("find path", "cameFrom", cameFrom)
	next := goal
	ok := true
	for ok {
		fullPath = append(fullPath, next)
		next, ok = cameFrom[next.GetID()]
	}
	slices.Reverse(fullPath)

	return
}
