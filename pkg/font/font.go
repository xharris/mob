package font

import (
	"bytes"
	"log"
	"log/slog"
	"os"

	ebitext "github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Font struct {
	faceSource *ebitext.GoTextFaceSource
	Name       string
	Size       float64
}

type FontOption func(*Font)

var DefaultFont *Font

func Init() {
	var err error
	DefaultFont, err = NewFont("asset/font/Retro Gaming.ttf", WFontSize(14))
	if err != nil {
		log.Panic("could not load font", err)
	}
}

func NewFont(name string, opts ...FontOption) (*Font, error) {
	f := &Font{Name: name}
	// load face source
	data, err := os.ReadFile(name)
	if err != nil {
		slog.Warn("could not read font", "err", err)
		return nil, err
	}
	f.faceSource, err = ebitext.NewGoTextFaceSource(bytes.NewReader(data))
	if err != nil {
		slog.Warn("could not create face source", "err", err)
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
	if !f.Valid() {
		slog.Warn("using an invalid font", "font", f)
		return 0, 0
	}
	ff := f.Face()
	return ebitext.Measure(text, &ff, 0)
}

func (f *Font) Face() ebitext.GoTextFace {
	return ebitext.GoTextFace{Source: f.faceSource, Size: f.Size}
}

func (f *Font) Valid() bool {
	return f.faceSource != nil
}
