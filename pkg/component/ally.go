package component

import "mob/pkg/allymod"

type Ally struct {
	Mods []allymod.Mod
	Name string
}

func (a *Ally) AddMod(mod allymod.Mod) {
	var mods []allymod.Mod
	for _, am := range a.Mods {
		if am.Name == mod.Name {
			mods = append(mods, mod)
		} else {
			mods = append(mods, am)
		}
	}
	a.Mods = mods
}

func (a *Ally) RemoveMod(mod allymod.Mod) {
	var mods []allymod.Mod
	for _, am := range a.Mods {
		if am.Name != mod.Name {
			mods = append(mods, mod)
		}
	}
	a.Mods = mods
}
