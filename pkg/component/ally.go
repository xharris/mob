package component

import "mob/pkg/mod"

type Ally struct {
	Mods []mod.Mod
	Name string
}

func (a *Ally) AddMod(m mod.Mod) {
	var mods []mod.Mod
	for _, am := range a.Mods {
		if am.Name == m.Name {
			mods = append(mods, m)
		} else {
			mods = append(mods, am)
		}
	}
	a.Mods = mods
}

func (a *Ally) RemoveMod(m mod.Mod) {
	var mods []mod.Mod
	for _, am := range a.Mods {
		if am.Name != m.Name {
			mods = append(mods, m)
		}
	}
	a.Mods = mods
}
