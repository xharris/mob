package component

import (
	"mob/pkg/pawn"
)

type ShopItem struct {
	AddMods     []pawn.Mod
	RemoveMods  []pawn.Mod
	UpgradeMods []pawn.Mod
}
