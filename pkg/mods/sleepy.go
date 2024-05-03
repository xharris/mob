package mods

import (
	"log/slog"
	"mob/pkg/component"
)

type sleepy struct{}

func Sleepy() component.Mod {
	return component.Mod{
		Name:   "sleepy",
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

	if mt.Timer.Ratio() == 0 {
		slog.Info("getting sleepy")
		combat.AddMods(Sleeping())
	}
}

type sleeping struct{}

func Sleeping() component.Mod {
	return component.Mod{
		Name:   "sleeping",
		Type:   component.Debuff,
		Target: component.TargetSelf,
		Move:   &sleeping{},
		Order:  component.OrderRandom,
		Once:   true,
	}
}

func (s *sleeping) Tick(mt component.MoveTick) {
	var combat *component.Combat
	mt.Self.Get(&combat)
	if mt.Timer.Ratio() == 0 {
		slog.Info("sleep")
		combat.MoveSpeed = 0
		combat.Skip += 1
	}
}
