package component

import (
	"github.com/sedyh/mizu/pkg/engine"
)

type ModType int

const (
	Attack ModType = iota
	Buff
	Debuff
)

type ModTarget int

const (
	TargetAlly ModTarget = iota
	TargetEnemy
	TargetSelf
)

type ModOrder int

const (
	Normal ModOrder = iota
	BeforeAll
	AfterAll
	BeforeAttack
	AfterAttack
)

type Mod struct {
	Name   string
	Desc   string
	Type   ModType
	Target ModTarget
	Range  NPCRange
	Order  ModOrder
	Move   Move
}

func (m *Mod) IsGood() bool {
	switch m.Type {
	case Attack, Debuff:
		return m.Target == TargetEnemy
	case Buff:
		return m.Target == TargetAlly || m.Target == TargetSelf
	default:
		return false
	}
}

type NPCType int

const (
	Ally NPCType = iota
	Enemy
)

type NPCStrategy int

const (
	Closest NPCStrategy = iota
	Farthest
	HighAttack
	LowAttack
	HighHealth
	LowHealth
	Random
)

type NPCRange float64

const (
	Melee  NPCRange = 30
	Ranged NPCRange = 90
)

type NPC struct {
	Mods []Mod
	// mods gained during combat
	combatMods []Mod
	Name       string
	Type       NPCType
	// targeting strategy
	Strategy NPCStrategy
}

func NewNPC(options ...NPCOption) (n NPC) {
	for _, opt := range options {
		opt(&n)
	}
	return
}

type NPCOption func(*NPC)

func WithMod(m Mod) NPCOption {
	return func(n *NPC) {
		n.AddMod(m)
	}
}

func WithMods(mods []Mod) NPCOption {
	return func(n *NPC) {
		for _, m := range mods {
			n.AddMod(m)
		}
	}
}

func WithName(name string) NPCOption {
	return func(n *NPC) {
		n.Name = name
	}
}

func WithType(t NPCType) NPCOption {
	return func(n *NPC) {
		n.Type = t
	}
}

func (n *NPC) AddMod(m Mod) {
	var mods []Mod
	found := false
	for _, am := range n.Mods {
		if am.Name == m.Name {
			found = true
		} else {
			mods = append(mods, am)
		}
	}
	if !found {
		mods = append(mods, m)
	}
	n.Mods = mods
}

func (n *NPC) AddCombatMod(m Mod) {
	n.combatMods = append(n.combatMods, m)
}

func (n *NPC) ClearCombatMods(m Mod) {
	n.combatMods = make([]Mod, 0)
}

func (n *NPC) RemoveMod(m Mod) {
	var mods []Mod
	for _, am := range n.Mods {
		if am.Name != m.Name {
			mods = append(mods, m)
		}
	}
	n.Mods = mods
}

func GetAllNPCOfType(w engine.World, npcType NPCType) (entities []engine.Entity) {
	for _, e := range w.View(NPC{}).Filter() {
		var npc *NPC
		e.Get(&npc)
		if npc.Type == npcType {
			entities = append(entities, e)
		}
	}
	return
}
