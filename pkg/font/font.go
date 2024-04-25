package font

import (
	"bytes"
	"os"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var faceSource *text.GoTextFaceSource

func GetTextFaceSource() (*text.GoTextFaceSource, error) {
	if faceSource != nil {
		return faceSource, nil
	}
	data, err := os.ReadFile("asset/font/Retro Gaming.ttf")
	if err != nil {
		return nil, err
	}
	faceSource, err = text.NewGoTextFaceSource(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return faceSource, nil
}
