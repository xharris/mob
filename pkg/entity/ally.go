package entity

import "mob/pkg/component"

type Ally struct {
	component.Health
	component.Ally
}
