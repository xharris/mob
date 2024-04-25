package component

import "github.com/hajimehoshi/ebiten/v2"

type Render struct {
	Image            *ebiten.Image
	DrawImageOptions ebiten.DrawImageOptions
	Size             float64
	W, H             float64
	Z                int
	X, Y             float64
	mouseInside      bool
}

func NewRender(w, h float64) Render {
	return Render{
		Image:            ebiten.NewImage(int(w), int(h)),
		DrawImageOptions: ebiten.DrawImageOptions{},
		W:                w,
		H:                h,
		Z:                0,
	}
}

func (r *Render) Resize(w, h int) {
	newImage := ebiten.NewImage(w, h)
	newImage.DrawImage(r.Image, &ebiten.DrawImageOptions{})
	r.Image = newImage
}

func (r *Render) MouseEntered() bool {
	geom := r.DrawImageOptions.GeoM
	geom.Invert()
	imx, imy := ebiten.CursorPosition()
	mx, my := geom.Apply(float64(imx), float64(imy))
	if !r.mouseInside && mx > 0 && mx < r.W && my > 0 && my < r.H {
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
	if r.mouseInside && !(mx > 0 && mx < r.W && my > 0 && my < r.H) {
		r.mouseInside = false
		return true
	}
	return false
}
