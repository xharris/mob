package font

import (
	"bytes"
	"os"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	ebitext "github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Font struct {
	faceSource *ebitext.GoTextFaceSource
	Size       float64
}

type FontOption func(*Font)

var DefaultFont, _ = NewFont("asset/font/Retro Gaming.ttf", WFontSize(14))

func NewFont(name string, opts ...FontOption) (*Font, error) {
	f := &Font{}
	// load face source
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}
	f.faceSource, err = ebitext.NewGoTextFaceSource(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	// options
	for _, opt := range opts {
		opt(f)
	}
	return f, nil
}

func WFontSize(size float64) FontOption {
	return func(f *Font) {
		f.Size = size
	}
}

func (f *Font) Measure(text string) (float64, float64) {
	return ebitext.Measure(text, f.Face(), 0)
}

func (f *Font) Face() text.Face {
	return &text.GoTextFace{Source: f.faceSource, Size: f.Size}
}
