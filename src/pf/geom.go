package pf

import (
	
)

func centerOfTri(a, b, c *Pos) Pos {
	x := (a.x + b.x + c.x) / 3
	y := (a.y + b.y + c.y) / 3
	return Pos{x, y}
}
