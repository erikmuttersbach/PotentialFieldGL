package geo

import (
	"math"	
)

type Vec struct {
	X, Y float64
}

func (v *Vec) Scale(u float64) *Vec {
	return &Vec{v.X*u, v.Y*u}
}

func (u *Vec) Add(v *Vec) *Vec {
	return &Vec{v.X+u.X, v.Y+u.Y}
}

// u - v
func (u *Vec) Sub(v *Vec) *Vec {
	return &Vec{u.X-v.X, u.Y-v.Y}
}

func (p *Vec) Distance(to *Vec) float64 {
	a := math.Abs(p.X-to.X)
	b := math.Abs(p.Y-to.Y) 
	return math.Sqrt(a*a + b*b)
}

func (p *Vec) Len() float64 {
	a := math.Abs(p.X)
	b := math.Abs(p.Y) 
	return math.Sqrt(a*a + b*b)
}

func (p *Vec) IsInf(sign int) bool {
	return math.IsInf(p.X, sign) && math.IsInf(p.Y, sign)
}

