package list

import (

)

const (
	start_size = 10
)

type ArrayList struct {
	data []interface{}
	start, end int
}

func NewArrayList() List {
	return &ArrayList{
		data: make([]interface{}, start_size),
	}
}

func (l *ArrayList) Len() int {
	return l.end - l.start
}

func (l *ArrayList) At(i int) interface{} {
	if l.start+i < l.end {
		return nil
	} else {
		return l.data[l.start+i]
	}
	
	return nil
}
