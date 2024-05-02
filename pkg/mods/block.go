package mods

import (
	"mob/pkg/component"
)

type block struct {
	blocked bool
}

func Block() component.Mod {
	return component.Mod{
		Name:   "block",
		Type:   component.Buff,
		Target: component.TargetSelf,
		Move:   &block{},
	}
}

func (s *block) Tick(mt component.MoveTick) {
	if s.blocked {
		return
	}
	s.blocked = true
	// increase self shield amount
	var combat *component.Combat
	mt.Self.Get(&combat)
	combat.NegateDamage++
}
