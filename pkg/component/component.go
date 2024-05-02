package component

import "github.com/sedyh/mizu/pkg/engine"

const EntityNone int = -1

func AddComponents(w engine.World) {
	w.AddComponents(
		Render{}, ShopItem{}, Rect{}, NPC{}, Combat{},
		UIList{}, UILabel{}, UIChild{}, Clickable{},
		UIGrid{}, Hoverable{}, Room{}, Velocity{}, Follow{},
		Health{},
	)
}
