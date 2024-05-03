package component

import (
	"math"
	"reflect"

	"github.com/hajimehoshi/ebiten/v2"
)

type AlphaLevel int

const (
	AlphaNone AlphaLevel = iota
	AlphaLow
	AlphaMid
	AlphaHigh
	AlphaFull
)

type RenderGeometry struct {
	Z          int
	X, Y       float64
	OX, OY     float64
	AlphaLevel AlphaLevel
}

func NewRenderGeometry() RenderGeometry {
	return RenderGeometry{
		Z:          0,
		AlphaLevel: AlphaFull,
	}
}

func (r *RenderGeometry) GetOptions(inherit ...RenderGeometry) ebiten.DrawImageOptions {
	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(-r.OX, -r.OY)
	op.GeoM.Translate(r.X, r.Y)
	// alpha
	op.ColorScale.ScaleAlpha(float32(r.AlphaLevel) / float32(AlphaFull))
	// apply other options if given
	for _, other := range inherit {
		otherOp := other.GetOptions()
		op.GeoM.Concat(otherOp.GeoM)
		op.ColorScale.ScaleWithColorScale(otherOp.ColorScale)
	}
	return op
}

type renderTexture struct {
	RenderGeometry
	Image *ebiten.Image
	Z     int
	w, h  int
}

func (r *renderTexture) GetSize() (int, int) {
	return r.Image.Bounds().Dx(), r.Image.Bounds().Dy()
}

func (r *renderTexture) Resize(w, h int) {
	if w != r.w || h != r.h {
		newImage := ebiten.NewImage(int(w), int(h))
		newImage.DrawImage(r.Image, &ebiten.DrawImageOptions{})
		r.Image = newImage
		r.w = w
		r.h = h
	}
}

type Render struct {
	RenderGeometry
	Textures             map[string]*renderTexture
	overrideW, overrideH int
	mouseInside          bool
	Visible              bool
	Debug                bool
}

type RenderOption func(*Render)

func NewRender(opts ...RenderOption) Render {
	r := Render{
		RenderGeometry: NewRenderGeometry(),
		Textures:       make(map[string]*renderTexture),
		overrideW:      -1,
		overrideH:      -1,
		Visible:        true,
	}
	for _, opt := range opts {
		opt(&r)
	}
	return r
}

func WGameSize() RenderOption {
	return func(r *Render) {
		w, h := ebiten.WindowSize()
		r.overrideW = int(w)
		r.overrideH = int(h)
	}
}

func WRenderDebug() RenderOption {
	return func(r *Render) {
		r.Debug = true
	}
}

func WRenderPosition(x, y float64) RenderOption {
	return func(r *Render) {
		r.X = x
		r.Y = y
	}
}

func WRenderOffset(x, y float64) RenderOption {
	return func(r *Render) {
		r.OX = x
		r.OY = y
	}
}

func (r *Render) GetTexture(component interface{}, w, h int) *renderTexture {
	name := reflect.TypeOf(component).String()
	texture, ok := r.Textures[name]
	if !ok {
		// create texture
		texture = &renderTexture{
			RenderGeometry: NewRenderGeometry(),
			Image:          ebiten.NewImage(max(1, w), max(1, h)),
			w:              w,
			h:              h,
		}
		r.Textures[name] = texture
	}
	// resize?
	texture.Resize(w, h)
	return texture
}

func (r *Render) GetSize() (w, h int) {
	if r.overrideH > -1 && r.overrideW > -1 {
		w = r.overrideW
		h = r.overrideH
		return
	}
	for _, texture := range r.Textures {
		if texture.Image.Bounds().Dx() > w {
			w = texture.Image.Bounds().Dx()
		}
		if texture.Image.Bounds().Dy() > h {
			h = texture.Image.Bounds().Dy()
		}
	}
	return
}

func (r *Render) Resize(w, h int) {
	r.overrideW = w
	r.overrideH = h
}

func (r *Render) ResetSize() {
	r.overrideW = -1
	r.overrideH = -1
}

func (r *Render) MouseInside() bool {
	op := r.RenderGeometry.GetOptions()
	op.GeoM.Invert()
	imx, imy := ebiten.CursorPosition()
	mx, my := op.GeoM.Apply(float64(imx), float64(imy))
	w, h := r.GetSize()
	return mx > 0 && mx < float64(w) && my > 0 && my < float64(h)
}

func (r *Render) MouseEntered() bool {
	if !r.mouseInside && r.MouseInside() {
		r.mouseInside = true
		return true
	}
	return false
}

func (r *Render) MouseExited() bool {
	if r.mouseInside && !r.MouseInside() {
		r.mouseInside = false
		return true
	}
	return false
}

func (r *Render) Distance(other Render) float64 {
	return math.Sqrt(math.Pow(other.X-r.X, 2) + math.Pow(other.Y-r.Y, 2))
}

func (r *Render) Apply(x, y float64) (float64, float64) {
	op := r.RenderGeometry.GetOptions()
	return op.GeoM.Apply(r.OX+x, r.OY+y)
}
