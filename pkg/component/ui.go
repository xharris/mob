package component

import (
	"image/color"
	"mob/pkg/font"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type UIGrid struct {
	ID      string
	Rows    int
	Columns int
}

type UIListDirection int

const (
	VERTICAL UIListDirection = iota
	HORIZONTAL
)

type UIList struct {
	ID        string
	Direction UIListDirection
}

type UIChild struct {
	Parent string
	// grid position
	X, Y int
	// list item size
	W, H int
}

type UILabelText struct {
	Color   color.Color
	Text    string
	Newline bool
	Font    *font.Font
}

type UILabel struct {
	Text   []UILabelText
	HAlign text.Align
	VAlign text.Align
	Font   *font.Font
}
