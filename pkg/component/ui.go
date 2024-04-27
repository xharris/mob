package component

import (
	"image/color"
	"mob/pkg/font"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type UI_ID string

type UIGrid struct {
	ID      UI_ID
	Rows    int
	Columns int
	Align   UIAlign
	Justify UIAlign
}

type UIListDirection int

const (
	VERTICAL UIListDirection = iota
	HORIZONTAL
)

type UIAlign int

const (
	START UIAlign = iota
	CENTER
	END
)

type UIList struct {
	ID          UI_ID
	Direction   UIListDirection
	Reverse     bool
	Align       UIAlign
	Justify     UIAlign
	FitContents bool
}

type UIChild struct {
	ID     UI_ID
	Parent UI_ID
	// grid position
	X, Y, I int
}

type UILabelText struct {
	Color   color.Color
	Text    string
	Newline bool
	// can be nil
	Font *font.Font
}

type UILabel struct {
	Text   []UILabelText
	HAlign text.Align
	VAlign text.Align
	// can be nil
	Font *font.Font
}
