package component

import (
	"log/slog"
	"math/rand"
	"mob/pkg/timer"
	"slices"

	"github.com/sedyh/mizu/pkg/engine"
)

type MoveTick struct {
	W           engine.World
	Name        string
	Self        engine.Entity
	Target      engine.Entity
	Timer       timer.Timer
	WithinRange bool
}

func (m MoveTick) IsDone() bool {
	return m.Timer.Done()
}

func (m MoveTick) Valid() bool {
	return m.Self != nil && m.Target != nil
}

type Move interface {
	Tick(mt MoveTick)
}

type Combat struct {
	Mods []Mod
	// added to Mods after all mods are called
	addedMods []Mod
	Shield    int
	// stacks of damage negation
	NegateDamage int
	// targeting strategy
	Strategy NPCStrategy
	// entity id
	Target      int
	WithinRange bool
	// -1 for infinite range
	AggroRange  float64
	AttackSpeed float64
	MoveSpeed   float64
	// skip X moves
	Skip     int
	ModIndex int
	// called when executing the move
	CurrentMoveTick MoveTick
}

func NewCombat(options ...CombatOption) (c Combat) {
	c.Target = EntityNone
	c.AggroRange = -1
	c.Reset()
	for _, opt := range options {
		opt(&c)
	}
	return
}

type CombatOption func(c *Combat)

func WithCombatNPC(npc *NPC) CombatOption {
	return func(c *Combat) {
		c.Mods = npc.Mods
	}
}

func (c *Combat) GetAttackRange() (r float64) {
	for _, m := range c.Mods {
		if float64(m.Range) > r {
			r = float64(m.Range)
		}
	}
	return
}

func (c *Combat) Reset() {
	c.AttackSpeed = 1.0
	c.MoveSpeed = 1.0
	c.Skip = 0
}

func (c *Combat) sortMods() {
	slices.SortStableFunc(c.Mods, func(a Mod, b Mod) int {
		bo := b.Order
		ao := a.Order
		if ao == OrderRandom {
			ao = ModOrder(rand.Intn(int(OrderRandom)))
		}
		if bo == OrderRandom {
			bo = ModOrder(rand.Intn(int(OrderRandom)))
		}
		return int(bo) - int(ao)
	})
}

func (c *Combat) UseMove(w engine.World, self engine.Entity, target engine.Entity) (move Move, tick MoveTick, ok bool) {
	if len(c.Mods) == 0 {
		ok = false
		slog.Warn("entity does not have any mods")
		return
	}
	ok = true
	// get distance to target
	var render *Render
	self.Get(&render)
	var tRender *Render
	target.Get(&tRender)
	if render == nil || tRender == nil {
		ok = false
		slog.Warn("entity is missing render component")
		return
	}
	targetDist := render.Distance(*tRender)

	firstRun := !c.CurrentMoveTick.Valid()
	newRun := c.CurrentMoveTick.Valid() && c.ModIndex >= len(c.Mods)
	if firstRun || newRun {
		// use new mods gained during combat and sort them
		c.ModIndex = 0
		c.Mods = append(c.Mods, c.addedMods...)
		c.addedMods = make([]Mod, 0)
		c.sortMods()
		// debug
		var names []string
		for _, m := range c.Mods {
			names = append(names, m.Name)
		}
		slog.Info("moveset", "i", c.ModIndex, "current", names)
	}
	currentMod := c.Mods[c.ModIndex]
	currentMoveIsDone := c.CurrentMoveTick.Valid() && c.CurrentMoveTick.IsDone()
	if currentMoveIsDone {
		// move just finished, reset stuff
		slog.Info("finished move", "name", c.CurrentMoveTick.Name)
		if c.Skip <= 0 {
			c.Reset()
		}
		if currentMod.Once {
			// remove from mods
			var newMods []Mod
			for i := range c.Mods {
				if i != c.ModIndex {
					newMods = append(newMods, c.Mods[i])
				}
			}
			c.Mods = newMods
		} else {
			// increase index for next call
			c.ModIndex++
		}
	}
	if firstRun || newRun || currentMoveIsDone {
		c.CurrentMoveTick = MoveTick{
			W:      w,
			Timer:  timer.NewTimer(1.0 / float64(len(c.Mods))),
			Name:   currentMod.Name,
			Self:   self,
			Target: target,
		}
	}
	move = currentMod.Move
	tick = c.CurrentMoveTick
	tick.WithinRange = targetDist <= float64(currentMod.Range)

	if !ok {
		return
	}

	if c.Skip > 0 {
		c.Skip--
		slog.Info("skipped", "left", c.Skip)
		return
	}

	move.Tick(tick)

	return
}

func (c *Combat) AddMods(mods ...Mod) {
	c.addedMods = append(c.addedMods, mods...)
}

func (c *Combat) MovesLeft() int {
	return len(c.Mods) - c.ModIndex
}
