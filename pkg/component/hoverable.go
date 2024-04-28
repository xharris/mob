package component

import "github.com/sedyh/mizu/pkg/engine"

type Hoverable struct {
	Enter func(e engine.Entity)
	Exit  func(e engine.Entity)
}
