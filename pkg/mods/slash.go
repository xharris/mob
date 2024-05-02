package mods

import (
	"log/slog"
	"mob/pkg/component"
)

type slash struct{}

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
	if mt.Timer.Ratio() == 0 {
		slog.Info("start slash")
	}
	if mt.Timer.Ratio() == 1 {
		slog.Info("done slash")
	}
}
