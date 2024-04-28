package component

import (
	"mob/pkg/mod"
)

type ShopItem struct {
	AddMods     []mod.Mod
	RemoveMods  []mod.Mod
	UpgradeMods []mod.Mod
	Cost        int
	Name        string
	Purchased   bool
}
