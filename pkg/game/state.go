package game

import (
	"mob/pkg/entity"
	"mob/pkg/level"
)

type State struct {
	Allies         []entity.Ally
	Level          level.Level
	LevelX, LevelY int
}

func NewState() (s State) {
	s.Level = level.Level{
		Level: 1,
	}
	return
}
