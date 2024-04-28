package component

import "github.com/sedyh/mizu/pkg/engine"

type Clickable struct {
	Click    func(e engine.Entity)
	Disabled bool
}
