package scene

import (
	"github.com/sedyh/mizu/pkg/engine"
)

type Setup struct{}

func (s *Setup) Setup(w engine.World) {
	w.ChangeScene(&Shop{
		Free:         true,
		AllyModCount: 3,
	})
}
