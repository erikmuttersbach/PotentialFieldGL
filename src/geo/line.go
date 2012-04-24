package geo

import (
	"fmt"
	"math"
)

func centerOfTri(a, b, c *Vec) Vec {
	x := (a.X + b.X + c.X) / 3
	y := (a.Y + b.Y + c.Y) / 3
	return Vec{x, y}
}

func IntersectLinePoint(ap, av, bp *Vec) *Vec {
	u := (bp.X-ap.X)/av.X
	v := (bp.Y-ap.Y)/av.Y
	s := u
	
	if math.IsInf(u, 0) || math.IsNaN(u) {
		if bp.X == ap.X + v*av.X {
			s = v
		} else {
			return nil
		}
	} else {
		if bp.Y == ap.Y + u * av.Y {
		} else {
			return nil
		}
	}
	
	if s >= 0 && s <= 1 {
		return av.Scale(s).Add(ap)
	}
	
	return nil
}

func IntersectLines(ap, av, bp, bv *Vec) *Vec {
	u := (-bp.X*bv.Y+ap.X*bv.Y+(bp.Y-ap.Y)*bv.X)/(av.Y*bv.X-av.X*bv.Y)
	v := -(av.X*(ap.Y-bp.Y)+av.Y*bp.X-ap.X*av.Y)/(av.Y*bv.X-av.X*bv.Y)
	
	//fmt.Println(ap, av, bp, bv, u, v)
	_ = fmt.Print
	if u >= 0 && u <= 1 && v >= 0 && v <= 1 {
		return ap.Add(av.Scale(u))
	} else if math.IsNaN(u) && math.IsNaN(v) {
		if ap.Distance(bp) <= av.Len() || ap.Distance(bp.Add(bv)) <= av.Len() {
			return &Vec{math.Inf(1), math.Inf(1)}
		}
	}
	
	return nil
}

func IntersectLines2(ap, av, bp, bv *Vec) *Vec {
	u := (-bp.X*bv.Y+ap.X*bv.Y+(bp.Y-ap.Y)*bv.X)/(av.Y*bv.X-av.X*bv.Y)
	v := -(av.X*(ap.Y-bp.Y)+av.Y*bp.X-ap.X*av.Y)/(av.Y*bv.X-av.X*bv.Y)
	
	//fmt.Println(u, v, ((u == 0 || u == 1) && (v > 0 && v < 1)), ((v == 0 || v == 1) && (u > 0 && u < 1)), ((u > 0 && u < 1) && (v > 0 && v < 1)))
	if ((u == 0 || u == 1) && (v > 0 && v < 1)) || ((v == 0 || v == 1) && (u > 0 && u < 1)) || ((u > 0 && u < 1) && (v > 0 && v < 1)) {
		return ap.Add(av.Scale(u))
	} else if math.IsNaN(u) && math.IsNaN(v) {
		if ap.Distance(bp) < av.Len() || ap.Distance(bp.Add(bv)) < av.Len() {
			return &Vec{math.Inf(1), math.Inf(1)}
		} 
	}
	
	return nil
}

