package game

import (
	"mob/pkg/component"
	"mob/pkg/level"
)

type State struct {
	Allies         []component.NPC
	Level          level.Level
	LevelX, LevelY int
}

func NewState() (s State) {
	s.Level = level.Level{
		Level: 1,
	}
	return
}
