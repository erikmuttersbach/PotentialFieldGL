package pf

import (
    "testing"
)

func TestRingbuffer(t *testing.T) {
	
	r := NewRingbuffer(3)
	r.AddToFront(1)
	r.AddToFront(2)
	r.AddToFront(3)
	
	e := r.Elements()
	if(e[0] != 3 || e[1] != 2 || e[2] != 1) {
		t.Error(e)
	}
	
	r.AddToFront(4)
	e = r.Elements()
	if(e[0] != 4 || e[1] != 3 || e[2] != 2) {
		t.Error(e)
	}
	
	r.AddToFront(5)
	r.AddToFront(6)
	e = r.Elements()
	if(e[0] != 6 || e[1] != 5 || e[2] != 4) {
		t.Error(e)
	}
}

