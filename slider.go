package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Slider struct {
	position Vector
	size     Vector
	value    float64
	start    int
	end      int
	name     string
	unit     string
	left     Button
	right    Button
}

func newSlider(position Vector, name, u string, min, max int) *Slider {
	return &Slider{position, Vector{400, 10}, float64(max-min) / 2, min, max, name, u, *newButton(position.Add(Vector{-20, -5})), *newButton(position.Add(Vector{400, -5}))}
}

func (src *Slider) Update() (changed bool) {
	if src.Hovered() && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, _ := ebiten.CursorPosition()
		src.value = float64(src.start) + ((float64(x)-src.position.x)/src.size.x)*float64(src.end-src.start)
		return true
	}
	if src.left.Clicked() && src.start > 0 {
		dif := src.end - src.start
		src.end = src.start
		src.start = src.end - dif
		src.value = float64(src.start) + float64(dif)/2
		return true
	}
	if src.right.Clicked() {
		dif := src.end - src.start
		src.start = src.end
		src.end = src.start + dif
		src.value = float64(src.start) + float64(dif)/2
		return true
	}
	return false
}

func (src *Slider) Draw(image *ebiten.Image) {
	ebitenutil.DrawRect(image, src.position.x, src.position.y, src.size.x, src.size.y, color.RGBA{100, 100, 100, 255})

	pos := ((src.value-float64(src.start))/float64(src.end-src.start)*src.size.x + src.position.x) - 15

	clr := color.RGBA{255, 255, 255, 255}
	if src.Hovered() {
		clr = color.RGBA{150, 150, 150, 255}
	}
	ebitenutil.DrawCircle(image, pos+15, src.position.y+5, 15, clr)
	ebitenutil.DebugPrintAt(image, fmt.Sprintf("%v: %.2f %v", src.name, src.value, src.unit), 20, int((src.position.y)-30))
	ebitenutil.DebugPrintAt(image, fmt.Sprintf("%v %v", src.start, src.unit), 40, int(src.position.y))
	ebitenutil.DebugPrintAt(image, fmt.Sprintf("%v %v", src.end, src.unit), 550, int(src.position.y))

	src.left.Draw(image)
	src.right.Draw(image)
}

func (src *Slider) Hovered() bool {
	a, b := ebiten.CursorPosition()
	x, y := float64(a), float64(b)
	return x > src.position.x && x < src.position.x+src.size.x && y > src.position.y-20 && y < src.position.y+src.size.y+20
}
