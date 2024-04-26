package component

import (
	"golang.org/x/exp/shiny/materialdesign/colornames"
)

type Tooltip struct {
	Parent UI_ID
	Text   []UILabelText
	Shown  bool
}

// add text from ShopItem
func (t *Tooltip) UseShopItem(si ShopItem) {
	for _, addMod := range si.AddMods {
		t.Text = append(t.Text,
			UILabelText{Text: addMod.Name, Color: colornames.Red300, Newline: true},
			UILabelText{Text: addMod.Desc, Color: colornames.Grey100, Newline: true},
		)
	}
}
