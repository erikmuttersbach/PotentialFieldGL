package geo

import (
    "testing"
)

func TestBoundingBox(t *testing.T) {
	p := &Polygon{[]*Vec{
		&Vec{0, 0},
		&Vec{10, 0},
		&Vec{10, 10},
		&Vec{0, 10},
	}}
	min, max := p.BoundingBox()
	
	if (*min != Vec{0, 0} || *max != Vec{10, 10}) {
		t.Fatal(min, max)
	}
	
	p = &Polygon{[]*Vec{
		&Vec{0, 0},
		&Vec{10, 0},
	}}
	min, max = p.BoundingBox()
	
	if (*min != Vec{0, 0} || *max != Vec{10, 0}) {
		t.Fatal(min, max)
	}
	
	p = &Polygon{[]*Vec{
		&Vec{0, 0},
	}}
	min, max = p.BoundingBox()
	
	if (*min != Vec{0, 0} || *max != Vec{0, 0}) {
		t.Fatal(min, max)
	}
}

func TestContainsPoint(t *testing.T) {
	p := &Polygon{[]*Vec{
		&Vec{10, 0},
		&Vec{10, 10},
		&Vec{0, 10},
		&Vec{0, 0},
	}}
	
	if !p.ContainsPoint(&Vec{2,2}) {
		t.FailNow()
	}
	
	if p.ContainsPoint(&Vec{0,1}) {
		t.FailNow()
	}
	
	if p.ContainsPoint(&Vec{1,0}) {
		t.FailNow()
	}
	
	if p.ContainsPoint(&Vec{10,10}) {
		t.FailNow()
	}
	
	if p.ContainsPoint(&Vec{10,0}) {
		t.FailNow()
	}
	
	if p.ContainsPoint(&Vec{0,10}) {
		t.FailNow()
	}
	
	p = &Polygon{[]*Vec{
		&Vec{0, 0},
		&Vec{10, 0},
		&Vec{5, 5},
		&Vec{5, 10},
		&Vec{0, 10},
	}}
	
	if !p.ContainsPoint(&Vec{2.5,5}) {
		t.FailNow()
	}
	
	_ = p
}


func TestContainsCornerLine(t *testing.T) {
	p := &Polygon{[]*Vec{
		&Vec{0, 0},
		&Vec{10, 0},
		&Vec{10, 10},
		&Vec{0, 10},
	}}
	
	// not contained
	if p.ContainsCornerLine(0, &Vec{-1,-1}) {
		t.Fail()
	}
	
	// not contained
	if p.ContainsCornerLine(0, &Vec{10,10}) {
		t.FailNow()
	}
	
	// contained
	if !p.ContainsCornerLine(1, &Vec{1,1}) {
		t.FailNow()
	}
	
	// not contained
	if p.ContainsCornerLine(1, &Vec{1,0}) {
		t.FailNow()
	}
	
	// not contained
	if p.ContainsCornerLine(1, &Vec{0,10}) {
		t.FailNow()
	}
	
	// not contained
	if p.ContainsCornerLine(1, &Vec{0,11}) {
		t.FailNow()
	}
		
}

func TestContainsCornerCornerLine(t *testing.T) {
	p := &Polygon{[]*Vec{
		&Vec{0, 0},
		&Vec{10, 0},
		&Vec{5, 5},
		&Vec{5, 10},
		&Vec{0, 10},
	}}
	
	// not contained
	if p.ContainsCornerCornerLine(0, 1) {
		t.Fail()
	}
	
	// not contained
	if p.ContainsCornerCornerLine(1, 3) {
		t.Fail()
	}
	
	// not contained
	if p.ContainsCornerCornerLine(3, 1) {
		t.Fail()
	}
	
	// not contained
	if p.ContainsCornerCornerLine(1, 3) {
		t.Fatalf("...")
	}
	
	// contained
	if !p.ContainsCornerCornerLine(2, 4) {
		t.Fatalf("...")
	}
	
	// contained
	if !p.ContainsCornerCornerLine(4, 2) {
		t.Fail()
	}
	
	// contained
	if !p.ContainsCornerCornerLine(3, 0) {
		t.Fail()
	}

	// contained
	if !p.ContainsCornerCornerLine(0, 3) {
		t.Fail()
	}
}

func TestSubPolygon(t *testing.T) {
	p := &Polygon{[]*Vec{
		&Vec{0, 0},
		&Vec{10, 0},
		&Vec{15, 5},
		&Vec{10, 10},
		&Vec{0, 10},
	}}
	
	pp := p.SubPolygon(0, 2)
	if pp.Points[0] != p.Points[0] || pp.Points[1] != p.Points[1] || pp.Points[2] != p.Points[2] {
		t.Fatal(pp)
	} 
	
	pp = p.SubPolygon(4, 1)
	if pp.Points[0] != p.Points[4] || pp.Points[1] != p.Points[0] || pp.Points[2] != p.Points[1] {
		t.Fatal(pp)
	} 
	
	pp = p.SubPolygon(4, 2)
	if pp.Points[0] != p.Points[4] || pp.Points[1] != p.Points[0] || pp.Points[2] != p.Points[1] || pp.Points[3] != p.Points[2] {
		t.Fatal(pp)
	} 
}

func TestSplitPolygon(t *testing.T) {
	p := &Polygon{[]*Vec{
		&Vec{0, 0},
		&Vec{10, 0},
		&Vec{15, 5},
		&Vec{10, 10},
		&Vec{0, 10},
	}}
	
	p1,p2 := p.SplitPolygon(0, 2)
	if 	p1.Points[0] != p.Points[0] || p1.Points[1] != p.Points[1] || p1.Points[2] != p.Points[2] || 
		p2.Points[0] != p.Points[2] || p2.Points[1] != p.Points[3] || p2.Points[2] != p.Points[4] || p2.Points[3] != p.Points[0] {
		t.Fatal(p1, p2)
	} 
	
	 
}
