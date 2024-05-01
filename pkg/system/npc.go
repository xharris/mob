package system

import (
	"log/slog"
	"math/rand"
	"mob/pkg/component"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sedyh/mizu/pkg/engine"
	"golang.org/x/exp/shiny/materialdesign/colornames"
)

const debug = true

type NPC struct{}

func (*NPC) Update(w engine.World) {
	v := w.View(component.NPC{}, component.Render{})
	entities := v.Filter()

	var allies []engine.Entity
	for _, e := range entities {
		var npc *component.NPC
		e.Get(&npc)
		if npc.Type == component.Ally {
			allies = append(allies, e)
		}
	}
	var enemies []engine.Entity
	for _, e := range entities {
		var npc *component.NPC
		e.Get(&npc)
		if npc.Type == component.Enemy {
			enemies = append(enemies, e)
		}
	}

	v = w.View(component.Follow{}, component.NPC{}, component.Render{})
	for _, e := range v.Filter() {
		var f *component.Follow
		var npc *component.NPC
		var render *component.Render
		e.Get(&f, &npc, &render)

		// reset
		npc.Target = component.EntityNone

		// set aggro range
		if npc.Type == component.Ally {
			npc.AggroRange = -1
		}
		if npc.Type == component.Enemy {
			npc.AggroRange = npc.GetAttackRange() * 2
		}

		// get list of targets
		var targets []engine.Entity
		if npc.Type == component.Enemy {
			targets = allies
		}
		if npc.Type == component.Ally {
			targets = enemies
		}

		if len(targets) > 0 && npc.Target == component.EntityNone {
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
			slog.Info("aggro?", "atk range", npc.GetAttackRange(), "range", npc.AggroRange, "dist", render.Distance(*tRender))
			if npc.AggroRange < 0 || render.Distance(*tRender) < float64(npc.AggroRange) {
				npc.Target = target.ID()
			}
		}

		targetEntity, hasTarget := w.GetEntity(npc.Target)
		if hasTarget {
			var tRender *component.Render
			var tNPC *component.NPC
			targetEntity.Get(&tRender, &tNPC)

			f.Radius = float64(max(render.GetSize()))
			f.Target = targetEntity.ID()
		}
	}
}

func (*NPC) Draw(w engine.World, screen *ebiten.Image) {
	if !debug {
		return
	}
	v := w.View(component.Render{}, component.NPC{})

	for _, e := range v.Filter() {
		var render *component.Render
		var npc *component.NPC
		e.Get(&render, &npc)

		x, y := render.Apply(0, 0)
		vector.StrokeCircle(screen, float32(x), float32(y), float32(npc.AggroRange), 1, colornames.Red200, true)
	}
}
