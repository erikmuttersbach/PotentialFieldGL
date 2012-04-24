package pf

import (
	"math"
	"geo"
)

type Pos struct {
	x float64
	y float64
}

type PosI struct {
	x int
	y int
}

func (p *Pos) ToScreen() (int, int) {
	return int(p.x), int(p.y)
}

func (p Pos) Distance(to Pos) float64 {
	a := math.Abs(p.x-to.x)
	b := math.Abs(p.y-to.y) 
	return math.Sqrt(a*a + b*b)
}

func P(x,y int) Pos {
	return Pos{x: float64(x), y: float64(y)}
}

func P64(x,y float64) Pos {
	return Pos{x: x, y: y}
}

func (p *Pos) V() *geo.Vec {
	return &geo.Vec{p.x, p.y}
}


