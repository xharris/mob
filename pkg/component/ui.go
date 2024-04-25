package component

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type UIGrid struct {
	ID      string
	Rows    int
	Columns int
}

type UIChild struct {
	Parent string
	// grid position
	X, Y int
}

type UILabelText struct {
	Color   color.Color
	Text    string
	Newline bool
}

type UILabel struct {
	Text   []UILabelText
	HAlign text.Align
	VAlign text.Align
}
