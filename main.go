package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	height   Slider
	mass     Slider
	velocity Slider
	friction Slider
	results  Results
}

type Results struct {
	t        float64
	a        float64
	v        float64
	a_values []float64
	v_values []float64
}

func count(h, m, v, b float64) *Results {
	var t, a float64
	dt := 0.0001 * h
	const g = 9.81
	avalues := []float64{}
	vvalues := []float64{}
	a = g
	for ; h > 0 || v < 0; t += dt {
		a = (m*g - b*v) / m
		h -= v * dt
		v += a * dt
		avalues = append(avalues, a)
		vvalues = append(vvalues, v)
	}
	return &Results{t, a, v, avalues, vvalues}
}

func DrawCrossedHorizontal(dst *ebiten.Image, x1, x2, y, slices float64, clr color.Color) {
	len := x2 - x1
	step := len / slices
	for x := x1; x <= x2; x += 2 * step {
		ebitenutil.DrawLine(dst, x, y, x+step, y, clr)
	}
}

func DrawGraph(screen *ebiten.Image, values []float64, start Vector, size float64, unit string, t float64) {
	zero_point := Vector{start.x, start.y + size}
	var max, min float64
	min = values[0]
	max = values[0]
	for _, v := range values {
		if v > max {
			max = v
		}
		if v < min {
			min = v
		}
	}
	y := zero_point.y
	if max > 0 {
		min = 0
		max = math.Ceil(max/10) * 10
	}
	if max <= 0 {
		max = 0
		y = zero_point.y - size
		min = math.Floor(min)
	}

	horizontal_line := &Vector{zero_point.x + size, y}
	vertical_line := &Vector{zero_point.x, zero_point.y - size}

	ebitenutil.DrawLine(screen, zero_point.x, y, horizontal_line.x, horizontal_line.y, color.White)
	ebitenutil.DrawLine(screen, zero_point.x, zero_point.y, vertical_line.x, vertical_line.y, color.White)

	arrow_length := float64(7)
	offset := float64(1)

	ebitenutil.DrawLine(screen, horizontal_line.x+offset, horizontal_line.y+offset, horizontal_line.x-arrow_length+offset, horizontal_line.y-arrow_length+offset, color.White)
	ebitenutil.DrawLine(screen, horizontal_line.x+offset, horizontal_line.y+offset, horizontal_line.x-arrow_length+offset, horizontal_line.y+arrow_length+offset, color.White)

	ebitenutil.DrawLine(screen, vertical_line.x+offset, vertical_line.y-offset, vertical_line.x+arrow_length+offset, vertical_line.y+arrow_length-offset, color.White)
	ebitenutil.DrawLine(screen, vertical_line.x+offset, vertical_line.y-offset, vertical_line.x-arrow_length+offset, vertical_line.y+arrow_length-offset, color.White)

	height := max - min
	step := size / float64(len(values))

	last := (values[len(values)-1] - min) / height * size
	DrawCrossedHorizontal(screen, start.x, start.x+size, zero_point.y-last, 20, color.White)

	for i, v := range values {
		y := (v - min) / height * size
		ebitenutil.DrawCircle(screen, zero_point.x+float64(step*float64(i)), zero_point.y-y, 2, color.White)
	}
	yzero := (values[0] - min) / height * size
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%.f", max)+unit, int(zero_point.x-60), int(start.y-5))
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%.f", min)+unit, int(zero_point.x-60), int(zero_point.y-5))
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%.2f", values[0])+unit, int(zero_point.x-60), int(zero_point.y-yzero)-5)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%.2f", values[len(values)-1])+unit, int(zero_point.x-60), int(zero_point.y-last)-5)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%.2f s", t), int(zero_point.x+size)-30, int(horizontal_line.y)+5)
}

func (game *Game) Draw(screen *ebiten.Image) {
	game.height.Draw(screen)
	game.mass.Draw(screen)
	game.velocity.Draw(screen)
	game.friction.Draw(screen)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Czas spadku: %.2f s\nPrzyspieszenie koncowe: %.2f m/s^2\nPredkosc koncowa: %.2f m/s", game.results.t, game.results.a, game.results.v), 50, 380)
	DrawGraph(screen, game.results.a_values, Vector{370, 500}, 200, " m/s^2", game.results.t)
	DrawGraph(screen, game.results.v_values, Vector{70, 500}, 200, " m/s", game.results.t)
	ebitenutil.DebugPrintAt(screen, "v(t)", 150, 470)
	ebitenutil.DebugPrintAt(screen, "a(t)", 450, 470)
}

func (game *Game) Update() error {
	if game.height.Update() || game.mass.Update() || game.velocity.Update() || game.friction.Update() {
		game.results = *count(game.height.value, game.mass.value, game.velocity.value, game.friction.value)
	}
	return nil
}

func (game *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 600, 750
}

func main() {
	ebiten.SetWindowSize(600, 750)
	ebiten.SetWindowTitle("Spadek z oporem")
	game := &Game{}
	game.height = *newSlider(Vector{100, 50}, "Wysokosc poczatkowa", "m", 0, 10)
	game.mass = *newSlider(Vector{100, 140}, "Masa", "kg", 0, 10)
	game.velocity = *newSlider(Vector{100, 230}, "Predkosc poczatkowa", "m/s", 0, 10)
	game.friction = *newSlider(Vector{100, 320}, "Wspolczynnik oporu", "", 0, 1)
	game.results = *count(game.height.value, game.mass.value, game.velocity.value, game.friction.value)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
