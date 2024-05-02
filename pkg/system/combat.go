package system

import (
	"math/rand"
	"mob/pkg/component"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sedyh/mizu/pkg/engine"
	"golang.org/x/exp/shiny/materialdesign/colornames"
)

const combatDebug = true

type Combat struct{}

func (*Combat) Update(w engine.World) {
	allies := component.GetAllNPCOfType(w, component.Ally)
	enemies := component.GetAllNPCOfType(w, component.Enemy)

	v := w.View(component.Follow{}, component.NPC{}, component.Render{}, component.Combat{})

	for _, e := range v.Filter() {
		var npc *component.NPC
		var render *component.Render
		var combat *component.Combat
		var f *component.Follow
		e.Get(&npc, &render, &combat, &combat, &f)

		// reset
		combat.Target = component.EntityNone

		// set aggro range
		if npc.Type == component.Ally {
			combat.AggroRange = -1
		}
		if npc.Type == component.Enemy {
			combat.AggroRange = combat.GetAttackRange() * 2
		}

		// get list of targets
		var targets []engine.Entity
		if npc.Type == component.Enemy {
			targets = allies
		}
		if npc.Type == component.Ally {
			targets = enemies
		}

		if len(targets) > 0 && combat.Target == component.EntityNone {
			var target engine.Entity
			// targeting strategy
			switch npc.Strategy {

			case component.Closest:
				slices.SortStableFunc(targets, func(a engine.Entity, b engine.Entity) int {
					var aRender *component.Render
					var bRender *component.Render
					a.Get(&aRender)
					b.Get(&bRender)
					return int(bRender.Distance(*render)) - int(aRender.Distance(*render))
				})
				target = targets[0]

			case component.Farthest:
				slices.SortStableFunc(targets, func(a engine.Entity, b engine.Entity) int {
					var aRender *component.Render
					var bRender *component.Render
					a.Get(&aRender)
					b.Get(&bRender)
					return int(aRender.Distance(*render)) - int(bRender.Distance(*render))
				})
				target = targets[0]

			case component.Random:
				target = targets[rand.Intn(len(targets))]

			default:
				target = targets[0]
			}
			// set target
			var tRender *component.Render
			target.Get(&tRender)
			if combat.AggroRange < 0 || render.Distance(*tRender) < float64(combat.AggroRange) {
				combat.Target = target.ID()
			}
		}

		targetEntity, hasTarget := w.GetEntity(combat.Target)
		if hasTarget {
			// move towards target
			var tRender *component.Render
			var tNPC *component.NPC
			targetEntity.Get(&tRender, &tNPC)

			f.Radius = combat.GetAttackRange()
			f.Target = targetEntity.ID()
			combat.WithinRange = render.Distance(*tRender) <= combat.GetAttackRange()
		}

		// has a target to fight
		target, targetExists := w.GetEntity(combat.Target)
		if !targetExists {
			continue
		}
		var tNPC *component.NPC
		var tRender *component.Render
		target.Get(&tNPC, &tRender)
		// use next move in moveset
		combat.UseMove(w, e, target)
	}
}

func (*Combat) Draw(w engine.World, screen *ebiten.Image) {
	if !combatDebug {
		return
	}
	v := w.View(component.Render{}, component.NPC{})

	for _, e := range v.Filter() {
		var render *component.Render
		var npc *component.NPC
		var combat *component.Combat
		e.Get(&render, &npc, &combat)

		x, y := render.Apply(0, 0)
		vector.StrokeCircle(screen, float32(x), float32(y), float32(combat.AggroRange), 1, colornames.Red200, true)
	}
}
