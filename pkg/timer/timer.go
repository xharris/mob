package timer

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sedyh/mizu/pkg/engine"
)

type TimerTick struct{}

var ticks int = 0

func (*TimerTick) Update(w engine.World) {
	ticks += 1
}

type Timer struct {
	start int
	// seconds
	every float64
}

func NewTimer(every float64) Timer {
	return Timer{start: ticks, every: every}
}

func (t *Timer) Done() bool {
	incr := t.every * float64(ebiten.TPS())
	if incr == 0 || ticks == t.start {
		return false
	}
	if (ticks-t.start)%int(incr) == 0 {
		// slog.Info("timer", "ticks", ticks-t.start, "incr", incr, "tps", ebiten.TPS())
		return true
	}
	return false
}

// returns 0-100
func (t *Timer) Ratio() int {
	incr := t.every * float64(ebiten.TPS())
	if incr == 0 || ticks == t.start {
		return 0 // at start
	}
	if float64(ticks-t.start) == incr {
		return 100 // done
	}
	return int(float64((ticks-t.start)%int(incr)) / float64(incr) * 100)
}
