package system

import (
	"image"
	"log"
	"math"
	"mob/pkg/component"
	"mob/pkg/logger"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/sedyh/mizu/pkg/engine"
)

var logNPC = logger.NewLogger(logger.WithDebug())

type NPC struct{}

var npcImages map[string]*ebiten.Image

func getNPCImage(path string) (img *ebiten.Image) {
	img, ok := npcImages[path]
	if !ok {
		var err error
		img, _, err = ebitenutil.NewImageFromFile(path)
		if err != nil {
			log.Fatalf("could not load %s (%s)", path, err)
		}
	}
	return
}

func (*NPC) Update(w engine.World) {
	v := w.View(component.NPC{}, component.Render{})
	for _, e := range v.Filter() {
		var npc *component.NPC
		var render *component.Render
		var vel *component.Velocity
		e.Get(&npc, &render, &vel)

		if vel != nil {
			if vel.X > 0 {
				render.SX = -1
			}
			if vel.X < 0 {
				render.SX = 1
			}
		}

		npc.UpdateAnimations()
	}
}

func (n *NPC) Draw(w engine.World, screen *ebiten.Image) {
	v := w.View(component.NPC{}, component.Render{})
	for _, e := range v.Filter() {
		var npc *component.NPC
		var render *component.Render
		e.Get(&npc, &render)

		texture := render.GetTexture(npc, 32, 32)
		if npc.Center {
			texture.OX = 16
			texture.OY = 16
		}
		// draw animation sprites
		var sprites []component.NPCAnimationSprite
		for _, spr := range npc.AnimationSprites {
			sprites = append(sprites, spr)
		}
		slices.SortStableFunc(sprites, func(a, b component.NPCAnimationSprite) int {
			return int(a.Part) - int(b.Part)
		})
		texture.Image.Clear()
		for _, spr := range sprites {
			img := getNPCImage(spr.Path)
			// animated image
			if spr.Frames > 1 {
				frame := int(math.Floor(float64(spr.Frames) * float64(npc.AnimationTimer.Ratio()/100.0)))
				frameWidth := img.Bounds().Dx() / spr.Frames
				img = img.SubImage(image.Rect(int(frame*frameWidth), 0, int(frame*frameWidth+frameWidth), img.Bounds().Dy())).(*ebiten.Image)
			}
			texture.Image.DrawImage(img, &ebiten.DrawImageOptions{})
		}
	}
}
