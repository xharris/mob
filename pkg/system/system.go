package system

import (
	"mob/pkg/timer"

	"github.com/sedyh/mizu/pkg/engine"
)

func AddSystems(w engine.World) {
	w.AddSystems(
		&RenderSystem{}, &ShopItem{},
		&RenderRect{}, &UIRenderLabel{}, &UIListLayout{},
		&Clickable{}, &UIGridLayout{}, &Hoverable{},
		&Velocity{}, &Follow{}, &NPC{}, &Combat{},
		&timer.TimerTick{}, &Combat{},
	)
}
