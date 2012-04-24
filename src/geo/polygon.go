package geo

import (
	"fmt"
	"math"
	"set"
	"container/list"
)

type Polygon struct {
	Points []*Vec
} 

func (p *Polygon) Len() int {
	return len(p.Points)
}

func (p *Polygon) String() string {
	str := "{"
	for _, pt := range p.Points {
		str += fmt.Sprint(pt)+","
	}
	str += "}"
	return str
}

func (p *Polygon) Center() (center *Vec) {
	center = &Vec{0.0,0.0}
	for i:=0; i<p.Len(); i++ {
		center.X += p.Points[i].X
		center.Y += p.Points[i].Y
	}
	
	center = center.Scale(1.0/float64(p.Len()))
	
	return center
}

// http://www.iti.fh-flensburg.de/lang/algorithmen/geo/polygon.htm
func (p *Polygon) IsConvex() bool {
	for i, pt1 := range p.Points {
		pt0 := p.Points[(len(p.Points)+i-1)%len(p.Points)]
		pt2 := p.Points[(i+1)%len(p.Points)]
		
		u := Vec{pt2.X-pt1.X, pt2.Y-pt1.Y}
		v := Vec{pt0.X-pt1.X, pt0.Y-pt1.Y}
		
		f := 0.5*(u.X*v.Y - u.Y*v.X)
		
		if f < 0 {
			return false
		}		
		
		if f == 0 {
			d := pt0.Distance(pt2)
			if pt1.Distance(pt0) > d || pt2.Distance(pt1) > d {
				return false
			} 
		}
		
		_ = math.E
		
	}
	
	return true
}

// Tests if a line, from p.points[pt_i] to v is contained
// in the polygon. 
// polygon is not required to be convex
func (p *Polygon) ContainsCornerLine(pt_i int, v *Vec) bool {
	bp := p.Points[pt_i]
	bv := v.Sub(bp)
	
	if !p.ContainsPoint(v) {
		return false
	}
	
	for i, pt := range p.Points {
		ab := pt
		av := p.Points[(i+1)%p.Len()].Sub(pt)
		
		if IntersectLines(ab, av, bp, bv) != nil {
		
			// if the intersection is on the corner point, allow it ... 
			if IntersectLinePoint(ab, av, bp) == nil {
				return false
			}
		}
	}
	
	return true
}

// Tests if a line, from p.points[pt_i] to v is contained
// in the polygon. 
// polygon is not required to be convex
func (p *Polygon) ContainsCornerCornerLine(pt1_i, pt2_i int) bool {
	bp := p.Points[pt1_i]
	bv := p.Points[pt2_i].Sub(bp)
	
	if !p.ContainsPoint(bp.Add(bv.Scale(0.5))) {
		//fmt.Println("Point ", bp.Add(bv.Scale(0.5)), "not inside -> false")
		return false
	}
	
	for i, pt := range p.Points {
		ab := pt
		av := p.Points[(i+1)%p.Len()].Sub(pt)
		
		if IntersectLines2(ab, av, bp, bv) != nil {
			//fmt.Println("intersect lines ", ab, av, "and", bp, bv)
			return false
		}
	}
	
	return true
}

func (p *Polygon) BoundingBox() (*Vec, *Vec) {
	min := *p.Points[0]
	max := *p.Points[0]
	for _, pt := range p.Points {
		if min.X > pt.X {
			min.X = pt.X
		} else if max.X < pt.X {
			max.X = pt.X
		}
		
		
		if min.Y > pt.Y {
			min.Y = pt.Y
		} else if max.Y < pt.Y {
			max.Y = pt.Y
		}
	}
	return &min, &max
}

// tests if a point is in the polygon
// if a point is on the edge, it is not in the poly ... 
func (p *Polygon) ContainsPoint(testPt *Vec) bool {
	_, max := p.BoundingBox()
	rb := &Vec{testPt.X, testPt.Y}
	rv := (&Vec{max.X+1, testPt.Y}).Sub(rb)
	
	//fmt.Println("ray", p.IntersectLine(rb, rv)) 
	
	intersections := make(map[Vec]bool)
	for i, pt := range p.Points {
		ab := pt
		av := p.Points[(i+1)%p.Len()].Sub(ab)
		//fmt.Println("Do intersect ?",  ab, av, rb, rv, ":", IntersectLines(ab, av, rb, rv))
		if ip := IntersectLines(ab, av, rb, rv); ip != nil {
			//fmt.Println(" Is point", testPt, "on line", ab,av,":", IntersectLinePoint(ab, av, testPt))
			if IntersectLinePoint(ab, av, testPt) == nil {			
				intersections[*ip] = true
			} else {
				// if the point is on the edges only once,
				// return false immediately
				return false
			}
			
		}
	}
	
	//fmt.Println("We have ", intersections)
	return len(intersections)%2 == 1
}

func (p *Polygon) IntersectLine(ap, av *Vec) *set.Set {
	ips := set.NewSet()
	for i, pt := range p.Points {
		bp := pt
		bv := p.Points[(i+1)%p.Len()].Sub(bp)
		
		
		if ip := IntersectLines(ap, av, bp, bv); ip != nil {
			ips.Put(*ip)
		}
	}
	
	return ips
} 

func (p *Polygon) SubPolygon(start, end int) *Polygon {
	var newLen int
	if end < start {
		newLen = p.Len()-start + end+1
	} else {
		newLen = end-start+1
	}
	
	//fmt.Println("start:", start, "end:", end, "len:", newLen)
	newPoly := &Polygon{make([]*Vec, newLen)}
	for i := 0; i < newLen; i++ {
		newPoly.Points[i] = p.Points[(start+i)%p.Len()]
	}
	return newPoly
}

func (p *Polygon) SplitPolygon(start, end int) (p1, p2 *Polygon) {
	return p.SubPolygon(start, end), p.SubPolygon(end, start)
}

func (p *Polygon) Triangulate() *list.List {
	if (p.Len() <= 3) {
		l := list.New()
		l.PushBack(p)
		return l
	} 
	
	l := list.New()
	
	for i, _ := range p.Points {
		if !p.ContainsCornerCornerLine(i, (i+2)%p.Len()) {
			continue
		}
		
		sub, rest := p.SplitPolygon(i, (i+2)%p.Len())
		if sub.IsConvex() {
			l.PushBack(sub)
			l.PushBackList(rest.Triangulate())
			break
		}
	}
	
	return l
}
