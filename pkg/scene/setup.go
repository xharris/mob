package scene

import (
	"log/slog"
	"mob/pkg/level"

	"github.com/sedyh/mizu/pkg/engine"
)

type Setup struct {
	Scene
	NewGame bool
}

func (s *Setup) Setup(w engine.World) {
	slog.Info("setup", "new game", s.NewGame)
	// new game
	if s.State.Level.Level == 1 || s.NewGame {
		// generate the level
		s.Scene.State.Level, s.Scene.State.LevelX, s.Scene.State.LevelY = level.Generate(s.State.Level.Level)
		// get one free ally from shop
		shop := &Shop{
			Scene:         s.Scene,
			AllyCount:     3,
			Free:          true,
			PurchaseLimit: 1,
			MustBuyOne:    true,
		}
		w.ChangeScene(shop)
		return
	}
	fx, fy, finishFound := s.State.Level.GetFinish()
	if !finishFound {
		slog.Warn("level does not have end")
		// TODO go to title screen
	}
	// finished level, move to next one
	if s.State.LevelX == fx && s.State.LevelY == fy {
		s.State.Level.Level++
		s.Scene.State.Level, s.Scene.State.LevelX, s.Scene.State.LevelY = level.Generate(s.State.Level.Level)
	}
	// pick next room to move to
	w.ChangeScene(&PickRoom{
		Scene: s.Scene,
	})
}
