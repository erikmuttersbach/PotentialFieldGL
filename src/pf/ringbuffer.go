package pf

import (
	
)

type Ringbuffer struct {
	buffer []interface{}
	front int
}

func NewRingbuffer(l int) *Ringbuffer {
	return &Ringbuffer{
		buffer: make([]interface{}, l),
		front: l-1,
	}
}

func (r *Ringbuffer) AddToFront(val interface{}) {
	insertAt := (r.front+1) % len(r.buffer)
	
	r.buffer[insertAt] = val
	r.front = insertAt
}

func (r *Ringbuffer) Front() interface{} {
	return r.buffer[r.front]
}

func (r* Ringbuffer) Elements() []interface{} {
	elms := make([]interface{}, len(r.buffer))
	for i := 0; i<len(r.buffer); i++ {
		val := r.buffer[(len(r.buffer)+r.front-i)%len(r.buffer)]
		elms[i] = val
	}
	return elms
}