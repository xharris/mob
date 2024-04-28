package component

import "github.com/sedyh/mizu/pkg/engine"

func AddComponents(w engine.World) {
	w.AddComponents(
		Render{}, ShopItem{}, Rect{}, Health{},
		UIList{}, UILabel{}, UIChild{}, Clickable{},
		UIGrid{}, Hoverable{}, Room{},
	)
}
