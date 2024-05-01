package component

import (
	"time"
)

type PathAgent struct {
	TargetX, TargetY float64
	TargetChanged    bool
	Path             [][]float64
	Cooldown         time.Time
	Pathing          bool
	NextPoint        [][]float64
}

func (p *PathAgent) GetNextPoint() (x float64, y float64, done bool) {
	if len(p.Path) == 0 {
		done = true
		return
	}
	x, y = p.Path[0][0], p.Path[0][1]
	return
}

func (p *PathAgent) SetTarget(x, y float64) {
	p.TargetX = x
	p.TargetY = y
	p.TargetChanged = true
}

func (p *PathAgent) NeedNewPath() bool {
	return p.TargetChanged && (p.Cooldown.IsZero() || p.Cooldown.Before(time.Now()))
}

func (p *PathAgent) SetPath(path [][]float64) {
	p.Path = path
	p.TargetChanged = false
	p.Cooldown = time.Now().Add(time.Second * time.Duration(10))
}

type PathObstacle struct {
	Solid  bool
	Weight int
}

func (p *PathObstacle) GetWeight() int {
	if p.Solid {
		return 999999
	}
	return p.Weight
}
