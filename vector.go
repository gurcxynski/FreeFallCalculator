package main

import "math"

type Vector struct {
	x float64
	y float64
}

func (src *Vector) Length() float64 {
	return math.Sqrt(src.x + src.y)
}

func (src *Vector) Add(arg Vector) Vector {
	return Vector{src.x + arg.x, src.y + arg.y}
}
