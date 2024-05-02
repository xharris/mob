package mods

import (
	"mob/pkg/component"
)

type slash struct {
	T int
}

func Slash() component.Mod {
	return component.Mod{
		Name:   "slash",
		Type:   component.Attack,
		Target: component.TargetAlly,
		Range:  component.Melee,
		Move:   &slash{},
	}
}

func (s *slash) Tick(mt component.MoveTick) {
	if !mt.WithinRange {
		return
	}
	s.T += 1
}
