package component

import (
	"mob/pkg/allymod"
)

type ShopItem struct {
	AddMods     []allymod.Mod
	RemoveMods  []allymod.Mod
	UpgradeMods []allymod.Mod
	Cost        int
}
