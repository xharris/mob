package component

import (
	"mob/pkg/ally"

	"golang.org/x/exp/shiny/materialdesign/colornames"
)

type ShopItem struct {
	AddMods     []ally.Mod
	RemoveMods  []ally.Mod
	UpgradeMods []ally.Mod
}

func (s *ShopItem) GetUIText() (ui UILabel) {
	for _, addMod := range s.AddMods {
		ui.Text = append(ui.Text,
			UILabelText{Text: addMod.Name + "\n", Color: colornames.Red300},
			UILabelText{Text: addMod.Desc, Color: colornames.Grey100},
		)
	}
	return
}
