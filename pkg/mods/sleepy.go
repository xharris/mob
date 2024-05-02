package mods

import (
	"log/slog"
	"mob/pkg/component"
)

type sleepy struct {
	ticked     bool
	isSleeping bool
}

func Sleepy() component.Mod {
	return component.Mod{
		Name:   "Sleepy",
		Desc:   "Might take a nap",
		Type:   component.Debuff,
		Target: component.TargetSelf,
		Move:   &sleepy{},
		Order:  component.BeforeAll,
	}
}

func (s *sleepy) Tick(mt component.MoveTick) {
	var combat *component.Combat
	var velocity *component.Velocity
	mt.Self.Get(&combat, &velocity)

	if s.isSleeping {
		combat.MoveSpeed = 0
	}
	if s.ticked {
		return
	}
	s.ticked = true
	if velocity == nil {
		slog.Warn("could not sleep, missing velocity")
		return
	}
	// skip the remaining actions
	combat.ReplaceMoves(1, Sleepy())
}
