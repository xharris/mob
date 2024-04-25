package component

import (
	"image/color"

	"golang.org/x/exp/shiny/materialdesign/colornames"
)

type TooltipText struct {
	Color   color.Color
	Text    string
	Newline bool
}

type Tooltip struct {
	Text  []TooltipText
	Shown bool
}

type RenderTooltips struct{}

// add text from ShopItem
func (t *Tooltip) UseShopItem(si ShopItem) {
	for _, addMod := range si.AddMods {
		t.Text = append(t.Text,
			TooltipText{Text: addMod.Name, Color: colornames.Red300, Newline: true},
			TooltipText{Text: addMod.Desc, Color: colornames.Grey100, Newline: true},
		)
	}
}
