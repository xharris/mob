package game

import (
	"mob/pkg/component"
)

type Ally struct {
	component.Health
	component.Ally
}

type State struct {
	Allies []Ally
}
