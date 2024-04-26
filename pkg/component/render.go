package component

import "github.com/hajimehoshi/ebiten/v2"

type Render struct {
	Image            *ebiten.Image
	DrawImageOptions ebiten.DrawImageOptions
	w, h             int
	Z                int
	X, Y             float64
	mouseInside      bool
	Visible          bool
	Debug            bool
}

type RenderOption func(*Render)

func NewRender(opts ...RenderOption) Render {
	gw, gh := ebiten.WindowSize()
	r := Render{
		Image:            ebiten.NewImage(gw, gh),
		DrawImageOptions: ebiten.DrawImageOptions{},
		w:                0,
		h:                0,
		Z:                0,
		Visible:          true,
	}
	for _, opt := range opts {
		opt(&r)
	}
	return r
}

func WRenderSize(w, h int) RenderOption {
	return func(r *Render) {
		r.Image = ebiten.NewImage(int(w), int(h))
	}
}

func WGameSize() RenderOption {
	return func(r *Render) {
		gw, gh := ebiten.WindowSize()
		r.Image = ebiten.NewImage(gw, gh)
	}
}

func WRenderDebug() RenderOption {
	return func(r *Render) {
		r.Debug = true
	}
}

func (r *Render) Resize(w, h int) {
	if w != r.w && h != r.h {
		newImage := ebiten.NewImage(w, h)
		newImage.DrawImage(r.Image, &ebiten.DrawImageOptions{})
		r.Image = newImage
	}
}

func (r *Render) GetSize() (int, int) {
	bounds := r.Image.Bounds().Max
	return bounds.X, bounds.Y
}

// resize other to fit inside r
func (r *Render) Fit(other *Render) {
	w, h := r.GetSize()
	otherW, otherH := other.GetSize()
	other.Resize(min(w, otherW), min(h, otherH))
}

func (r *Render) MouseEntered() bool {
	geom := r.DrawImageOptions.GeoM
	geom.Invert()
	imx, imy := ebiten.CursorPosition()
	mx, my := geom.Apply(float64(imx), float64(imy))
	w, h := r.GetSize()
	if !r.mouseInside && mx > 0 && mx < float64(w) && my > 0 && my < float64(h) {
		r.mouseInside = true
		return true
	}
	return false
}

func (r *Render) MouseExited() bool {
	geom := r.DrawImageOptions.GeoM
	geom.Invert()
	imx, imy := ebiten.CursorPosition()
	mx, my := geom.Apply(float64(imx), float64(imy))
	w, h := r.GetSize()
	if r.mouseInside && !(mx > 0 && mx < float64(w) && my > 0 && my < float64(h)) {
		r.mouseInside = false
		return true
	}
	return false
}

func (r *Render) MouseInside() bool {
	return r.mouseInside
}
