package graph

import (

)

type Node interface {
	GetLinks() []Node
	C(other Node) float64
}