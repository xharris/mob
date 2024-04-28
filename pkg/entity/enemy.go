package entity

import "mob/pkg/component"

type Enemy struct {
	component.Health
	component.Enemy
}
