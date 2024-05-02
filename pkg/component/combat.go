package component

import (
	"log/slog"
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
	Mods   []Mod
	Shield int
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
	ModIndex    int
	// called when executing the move
	CurrentMoveTick MoveTick

	currentMoveset   []Move
	currentMoveticks []MoveTick
}

func NewCombat(options ...CombatOption) (c Combat) {
	c.Target = EntityNone
	c.AggroRange = -1
	c.AttackSpeed = 1.0
	c.MoveSpeed = 1.0
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

func (c *Combat) UseMove(w engine.World, targetDist float64) (move Move, tick MoveTick, ok bool) {
	if len(c.Mods) == 0 {
		ok = false
		return
	}
	ok = true
	if !c.CurrentMoveTick.Valid() || c.CurrentMoveTick.IsDone() {
		if c.ModIndex == 0 {
			// sort mods by order
			slices.SortStableFunc(c.Mods, func(a Mod, b Mod) int {
				return int(b.Order) - int(a.Order)
			})
		}
		// use next 'move' (mod)
		m := c.Mods[c.ModIndex]
		c.CurrentMoveTick = MoveTick{
			W:           w,
			Timer:       timer.NewTimer(float64(1 / len(c.Mods))),
			WithinRange: targetDist <= float64(m.Range),
			Name:        m.Name,
		}
		c.ModIndex++
		if c.ModIndex >= len(c.Mods) {
			c.ModIndex = 0
		}
	}
	move = c.Mods[c.ModIndex].Move
	tick = c.CurrentMoveTick
	return
}

func (c *Combat) ClearMoveset() {
	c.currentMoveset = make([]Move, 0)
	c.currentMoveticks = make([]MoveTick, 0)
}

func (c *Combat) ReplaceMoves(count int, replacement Mod) {
	for i := range c.currentMoveset {
		if i > 0 && i <= count {
			slog.Info("replace", "", c.currentMoveticks[i].Name, "->", replacement.Name)
			c.currentMoveset[i] = replacement.Move
		}
	}
}

func (c *Combat) MovesLeft() int {
	return len(c.currentMoveset)
}
