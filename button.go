package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Button struct {
	position Vector
	size     Vector
}

func newButton(pos Vector) *Button {
	return &Button{pos, Vector{20, 20}}
}

func (src *Button) Draw(image *ebiten.Image) {
	clr := color.RGBA{255, 255, 255, 255}
	if src.Hovered() {
		clr = color.RGBA{150, 150, 150, 255}
	}
	ebitenutil.DrawRect(image, src.position.x, src.position.y, src.size.x, src.size.y, clr)
}

func (src *Button) Hovered() bool {
	a, b := ebiten.CursorPosition()
	x, y := float64(a), float64(b)
	return x > src.position.x && x < src.position.x+src.size.x && y > src.position.y && y < src.position.y+src.size.y
}

func (src *Button) Clicked() bool {
	return src.Hovered() && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
}
