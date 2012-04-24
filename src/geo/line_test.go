package geo

import (
    "testing"
    "math"
)

func TestIntersectLinePoint(t *testing.T) {
	ap := &Vec{1,1}
	av := &Vec{1,1}
	
	// Real intersection
	ip := IntersectLinePoint(ap, av, &Vec{1.5, 1.5})
	if ip == nil || !(ip.X == 1.5 && ip.Y == 1.5) {
		t.Fatal(ip)
	}
	
	// edge intersection
	ip = IntersectLinePoint(ap, av, &Vec{1, 1})
	if ip == nil || !(ip.X == 1 && ip.Y == 1) {
		t.Fatal(ip)
	}
	
	// edge intersection
	ip = IntersectLinePoint(ap, av, &Vec{2, 2})
	if ip == nil || !(ip.X == 2 && ip.Y == 2) {
		t.Fatal(ip)
	}
	
	// edge intersection, y always 0
	ip = IntersectLinePoint(&Vec{1, 0}, &Vec{2, 0}, &Vec{1.5, 0})
	if ip == nil || !(ip.X == 1.5 && ip.Y == 0) {
		t.Fatal(ip)
	}
	
	// edge intersection, y always 0
	ip = IntersectLinePoint(&Vec{1, 0}, &Vec{2, 0}, &Vec{1, 0})
	if ip == nil || !(ip.X == 1 && ip.Y == 0) {
		t.Fatal(ip)
	}
	
	// edge intersection, x always 0
	ip = IntersectLinePoint(&Vec{0, -1}, &Vec{0, 10}, &Vec{0, 9})
	if ip == nil || !(ip.X == 0 && ip.Y == 9) {
		t.Fatal(ip)
	}
	
	// edge intersection, x always 0
	ip = IntersectLinePoint(&Vec{0, -1}, &Vec{0, 10}, &Vec{0, -1})
	if ip == nil || !(ip.X == 0 && ip.Y == -1) {
		t.Fatal(ip)
	}
	
	// no intersection
	ip = IntersectLinePoint(ap, av, &Vec{3, 3})
	if ip != nil {
		t.Fatal(ip)
	}
	
	// no intersection
	ip = IntersectLinePoint(ap, av, &Vec{0, 0})
	if ip != nil {
		t.Fatal(ip)
	}
}

func TestIntersectLines(t *testing.T) {
	ap := &Vec{1,1}
	av := &Vec{1,1}
	
	bp := &Vec{2,1}
	bv := &Vec{-1, 1}
	
	// Real intersection
	ip := IntersectLines(ap, av, bp, bv)
	if ip == nil || !(ip.X == 1.5 && ip.Y == 1.5) {
		t.Fatal(ip)
	}
	
	// intersection at base point ap
	ip = IntersectLines(ap, av, ap, bv)
	if ip == nil || !(ip.X == ap.X && ip.Y == ap.Y) {
		t.Fatal(ip)
	}
	
	// intersection with one base point on line
	ip = IntersectLines(ap, av, &Vec{1.5, 1.5}, &Vec{1, 0})
	if ip == nil || !(ip.X == 1.5 && ip.Y == 1.5) {
		t.Fatal(ip)
	}
	
	// no intersection (parallel)
	ip = IntersectLines(ap, av, bp, av)
	if ip != nil {
		t.Fatal(ip)
	}
	
	// Inf+ intersection, parallel line with distance 0
	ip = IntersectLines(&Vec{0, 0}, &Vec{2, 2}, &Vec{1, 1}, &Vec{2, 2})
	if !ip.IsInf(1) {
		t.Fatal(ip)
	}
	
	// Inf+ intersection, parallel line with distance 0
	ip = IntersectLines(&Vec{0, 0}, &Vec{2, 2}, &Vec{3, 3}, &Vec{2, 2})
	if ip != nil {
		t.Fatal(ip)
	}
	
	_ = math.E
}


func TestIntersectLines2(t *testing.T) {
	ap := &Vec{1,1}
	av := &Vec{1,1}
	
	bp := &Vec{2,1}
	bv := &Vec{-1, 1}
	
	// Real intersection
	ip := IntersectLines2(ap, av, bp, bv)
	if ip == nil || !(ip.X == 1.5 && ip.Y == 1.5) {
		t.Fatal(ip)
	}
	
	// intersection at base point ap
	ip = IntersectLines2(ap, av, ap, bv)
	if ip != nil {
		t.Fatal(ip)
	}
	
	// intersection with one base point on line
	ip = IntersectLines2(ap, av, ap, &Vec{1, 0})
	if ip != nil {
		t.Fatal(ip)
	}
	
	// intersection with target point
	ip = IntersectLines2(&Vec{1,1}, &Vec{2,1}, &Vec{7,4}, &Vec{-2,-1})
	if ip != nil {
		t.Fatal(ip)
	}
	
	// no intersection (parallel)
	ip = IntersectLines2(ap, av, bp, av)
	if ip != nil {
		t.Fatal(ip)
	}
	
	// Inf+ intersection, parallel line with distance 0
	ip = IntersectLines2(&Vec{0, 0}, &Vec{2, 2}, &Vec{1, 1}, &Vec{2, 2})
	if ip == nil || !ip.IsInf(1) {
		t.Fatal(ip)
	}
	
	// Inf+ intersection, parallel line with distance 0
	ip = IntersectLines2(&Vec{0, 0}, &Vec{2, 2}, &Vec{3, 3}, &Vec{-2, -2})
	if ip == nil || !ip.IsInf(1) {
		t.Fatal(ip)
	}
	
	// Inf+ intersection, parallel line with distance 0
	ip = IntersectLines2(&Vec{0, 0}, &Vec{2, 2}, &Vec{3, 3}, &Vec{-1, -1})
	if ip != nil {
		t.Fatal(ip)
	}
	
	// no intersection
	ip = IntersectLines2(&Vec{10, 0}, &Vec{-5, 5}, &Vec{5, 10}, &Vec{-5, -10})
	if ip != nil {
		t.Fatal(ip)
	}
}