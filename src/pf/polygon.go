package pf

import (
	"fmt"
	"math"
)

type Polygon struct {
	points []*Pos
} 

func (p *Polygon) Len() int {
	return len(p.points)
}

func (p *Polygon) String() string {
	str := "{"
	for _, pt := range p.points {
		str += fmt.Sprint(pt)+","
	}
	str += "}"
	return str
}

// http://www.iti.fh-flensburg.de/lang/algorithmen/geo/polygon.htm
func (p *Polygon) IsConvex() bool {
	for i, pt1 := range p.points {
		pt0 := p.points[(len(p.points)+i-1)%len(p.points)]
		pt2 := p.points[(i+1)%len(p.points)]
		
		u := Pos{pt2.x-pt1.x, pt2.y-pt1.y}
		v := Pos{pt0.x-pt1.x, pt0.y-pt1.y}
		
		f := 0.5*(u.x*v.y - u.y*v.x)
		
		if f < 0 {
			return false
		}		
		
		if f == 0 {
			d := pt0.Distance(*pt2)
			if pt1.Distance(*pt0) > d || pt2.Distance(*pt1) > d {
				return false
			} 
		}
		
		_ = math.E
		
	}
	
	return true
}