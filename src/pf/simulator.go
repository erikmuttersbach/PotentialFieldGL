package pf

import (
	"fmt"
	"math"
	"time"
	"math2"
	"github.com/banthar/gl"
	"github.com/banthar/Go-SDL/sdl"
)


type unit struct {
	id int
	pos Pos
	end Pos
	
	trail *Ringbuffer
	color *Color
}

type Sim struct {
	units map[int]*unit
	buildings map[int]*building
	static [][]float64
	nav *NavMesh	// List of Polygons
	
	ui *UI
	
	run bool
}

type Color struct {
	r,g,b,a float32
}

func NewColor(r,g,b,a float32) *Color {
	return &Color{r,g,b,a}
}

type building struct {
	id int
	x, y, w, h float64
}

func NewSim() *Sim {	
	_ = fmt.Println
	
	s := &Sim{
		units: make(map[int]*unit),
		buildings: make(map[int]*building),
		ui: InitUI(),
		run: false,
		nav: NewNavMesh(),
	}
	
	// Init static potential field
	s.static = make([][]float64, 300)
	for x := 0; x<300; x++ {
		s.static[x] = make([]float64, 300)
		for y := 0; y<300; y++ {
			s.static[x][y] = 0.0
		}
	}
	
	return s
}

func (s *Sim) Init2() {
	/*poly1 := []*Pos{
		&Pos{0, 50},
		&Pos{50, 0}, 
		&Pos{50, 150},
		&Pos{0, 100},
	}
	node1 := NewNavNode(&Polygon{poly1})
	s.nav.nodes.PushBack(node1)
	
	poly2 := []*Pos{
		&Pos{200, 0},
		&Pos{200, 150}, 
		&Pos{50, 150},
		&Pos{50, 0},
	}
	node2 := NewNavNode(&Polygon{poly2})
	e2 := s.nav.nodes.PushBack(node2)
	
	node1.Merge(node2)
	s.nav.nodes.Remove(e2)
	
	fmt.Println(node1, node1.node.IsConvex())*/
	
	poly1 := []*Pos{
		&Pos{0, 50},
		&Pos{50, 0}, 
		&Pos{100, 0},
		&Pos{150, 0},
		&Pos{50, 150},
		&Pos{0, 100},
	}
	node1 := NewNavNode(&Polygon{poly1})
	s.nav.nodes.PushBack(node1)
	
	fmt.Println(node1.node.IsConvex())
	
}

func (s *Sim) Init() {
	/*s.nav.PushFront(&Polygon{
		points: []Pos{
			Pos{x: 1, y: 1},
			Pos{x: 50, y: 1},
			Pos{x: 50, y: 50},
			Pos{x: 1, y: 50},
		},
	})*/
	
	outers := []*Pos{
		&Pos{0, 0},
		&Pos{299, 0},
		&Pos{299, 299},
		&Pos{0, 299},
	}
	
	for _, building := range s.buildings {
		inners := []*Pos{
			&Pos{building.x, building.y},
			&Pos{building.x+building.w, building.y},
			&Pos{building.x+building.w, building.y+building.h},
			&Pos{building.x, building.y+building.h},
		}
		
		var prev *NavNode
		for i_inner, inner := range inners {
		
			tri1 := []*Pos{
				outers[i_inner],
				outers[(i_inner+1)%len(outers)], 
				inner, 
			}
			node1 := NewNavNode(&Polygon{tri1})
			
			tri2 := []*Pos{
				outers[(i_inner+1)%len(outers)],
				inners[(i_inner+1)%len(inners)],
				inner, 
			}
			node2 := NewNavNode(&Polygon{tri2})
			
			// Connect the links
			if(prev != nil) {
				node1.links[2] = prev
				prev.links[0] = node1
			}
			node1.links[1] = node2
			node2.links[2] = node1
			
			s.nav.nodes.PushBack(node1)
			s.nav.nodes.PushBack(node2)
			
			prev = node2
		}	
		
		firstNode := s.nav.nodes.Front().Value.(*NavNode)
		firstNode.links[2] = prev
		prev.links[0] = firstNode
	}
	
	s.nav.Reduce()
	fmt.Println("\n", s.nav.nodes.Back().Value)
	/*
	
	
	node1 := s.nav.nodes.Front().Value.(*NavNode)
	node2 := s.nav.nodes.Front().Next().Value.(*NavNode)
	//fmt.Println(node1)
	//fmt.Println(node2)
	
	fmt.Println()
	new1 := node1.Merge(node2)
	
	s.nav.nodes.Front().Value = new1
	s.nav.nodes.Remove(s.nav.nodes.Front().Next())
	//fmt.Println(new1)
	
	fmt.Println("\nSECOND MERGE")
	node1 = s.nav.nodes.Front().Value.(*NavNode)
	node2 = s.nav.nodes.Back().Value.(*NavNode)
	
	fmt.Println(node1)
	fmt.Println(node2)
	fmt.Println()
	
	new2 := node1.Merge(node2)
	fmt.Println(new2)
	
	s.nav.nodes.Front().Value = new2
	s.nav.nodes.Remove(s.nav.nodes.Back())
	fmt.Println(new2.node.IsConvex())*/
	
	// TWO-A
	/*node2 = s.nav.nodes.Front().Next().Value.(*NavNode)
	node3 := s.nav.nodes.Front().Next().Next().Value.(*NavNode)
	node2.Link(node3)
	s.nav.nodes.Remove(s.nav.nodes.Front().Next().Next())
	
	fmt.Println(node2.node)
	
	// TWO-A
	node2 = s.nav.nodes.Front().Next().Value.(*NavNode)
	node3 = s.nav.nodes.Front().Next().Next().Value.(*NavNode)
	node2.Link(node3)
	s.nav.nodes.Remove(s.nav.nodes.Front().Next().Next())
	
	fmt.Println(node2.node)	*/
}

func findNearest(pos *Pos, search []*Pos) int {
	minDist := math.MaxFloat64
	min_i := 0
	
	for i, p := range search {
		dist := p.Distance(*pos)
		if dist < minDist {
			minDist = dist
			min_i = i
		}
	}
	
	return min_i
}

func (s *Sim) AddUnit(start, end Pos, color *Color) {
	id := len(s.units)
	unit := &unit{
		id: id,
		pos: start,
		end: end,
		trail: NewRingbuffer(5),
		color: color,
	}
	
	s.units[id] = unit
}

func (s *Sim) AddBuilding(x, y, w, h float64) {
	id := len(s.buildings)
	building := &building{
		id: id,
		x: x,
		y: y,
		w: w,
		h: h,
	}	
	s.buildings[id] = building
	
	// update the static potential field
	for ix := math2.Round64(x); ix < math2.Round64(x+w); ix++ {
		for iy := math2.Round64(y); iy < math2.Round64(y+h); iy++ {
			s.static[ix][iy] = 1			
		}
	}
}



func (s *Sim) Update() {
	speed := 20.0/float64(time.Second) // -> 10 units per second
	
	lastUpdate := time.Now()
	time.Sleep(1*time.Nanosecond)
		
	for ;; {
	
		dTime := time.Since(lastUpdate)
		start := time.Now()
		
		// Update the units
		for _, unit := range s.units {
			if(unit.pos == unit.end) {
				// TODO when to finish?
				continue
			}
		
			min, radMin := 1.0, 0.0
			for i := 0; i<16; i++ {
				// TODO We can cache all these static directions
				rad := math.Pi*2*float64(i)/16.0
				dy := math.Sin(rad)
				dx := math.Cos(rad)
				
				pot := s.potential(unit, unit.pos.x+dx, unit.pos.y+dy)
				
				if(pot < min) {
					min = pot
					radMin = rad
				}
			}
			
			d := speed*float64(dTime)
			unit.pos.x += math.Cos(radMin)*d
			unit.pos.y += math.Sin(radMin)*d
			
			if last := unit.trail.Front(); last != nil {
				lastPos := last.(Pos)
				if(lastPos.Distance(unit.pos) >= 0.25) {
					unit.trail.AddToFront(unit.pos)
				}
			} else {
				unit.trail.AddToFront(unit.pos)
			}
		}
		
		// Important: update BEFORE the render
		lastUpdate = time.Now()
		
		fps := 1/(float64(time.Since(start))/float64(time.Second))
		_ = fps
		//fmt.Printf("%.f \t\n", fps)
		//fmt.Println(fps, s.units[0].pos)
		//s.ui.fpsLabel.SetLabel(fmt.Sprintf("%f", fps))
		
		time.Sleep(1*time.Millisecond)
	}
}

func (s *Sim) potential(unit *unit, x, y float64) float64 {
	// Distance to end
	max := 300.0
	dist := P64(x,y).Distance(unit.end)
	
	potEnd := dist/max
	
	// other units
	MIN := 10.0
	MAX := 15.0
	potDist := 0.0
	for _, oUnit := range s.units {
		if(oUnit.id != unit.id) {
			dist := oUnit.pos.Distance(P64(x,y))
			potDist += 1-math2.MinMax(1, (dist-MIN)/(MAX-MIN), 0)
		}		
	}
	
	potDist = math2.MinMax(1, potDist, 0)
	
	// trail
	potTrail := 0.0
	elms := unit.trail.Elements();
	for _, elm := range elms {
		if(elm != nil) {
			trailPos := elm.(Pos)
			dist := trailPos.Distance(P64(x,y))
			potTrail += 1.0-math2.MinMax(1, dist/float64(1), 0)
		}		
	}
	
	potTrail = math2.MinMax(1, potTrail, 0)
	
	if unit.id == 0  && potDist > 0 {
		//fmt.Println(potDist)
	}

	pot := 0.5*potEnd+0.1*potDist + 0.4*potTrail
	
	xd := int(x)
	yd := int(y)
	if xd >= 0 && xd < 300 && yd >= 0 && yd < 300 {
		if s.static[xd][yd] == 1 {
			pot = 1
		}
	}
	
	return pot
}

func (s *Sim) Draw() {
	

	//start := time.Now()
	
	// Init OpenGL
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.Enable(gl.BLEND)
	gl.Enable(gl.POINT_SMOOTH)
	gl.Enable(gl.LINE_SMOOTH)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	gl.Begin(gl.POINTS)
	
	gl.PointSize(1)
	for x := int(0); x < int(300); x+=1 {
		for y := int(0); y < int(300); y+=1 {
			pot := s.potential(s.units[0], float64(x), float64(y))
			
			//r := math.Min(pot*2, 1)
			//g := math.Min(pot*2, 2)-1
			gl.Color4f(1-float32(pot), 1-float32(pot), 1-float32(pot), 1)
			gl.Vertex2i(x, y)
		}
	}
	
	gl.End()
	
	// Draw Units
	gl.PointSize(1)
	gl.Begin(gl.POINTS)	
	for _, unit := range s.units {
		x, y := unit.pos.ToScreen()
		gl.Color4f(1, 0, 0, 1)
		gl.Vertex2i(x, y)
	}
	
	gl.End()
	
	// Nav mesh
	for e := s.nav.nodes.Front(); e != nil; e = e.Next() {
		poly := e.Value.(*NavNode)

		gl.Color4f(0, 0, 1, 0.2)
		gl.Begin(gl.POLYGON)
		for _, pos := range poly.node.points {
			if pos == nil {
				continue
			}
			
			gl.Vertex2f(float32(pos.x), float32(pos.y))
		}
		gl.End()
		
		gl.Color4f(0, 0, 1, 1)
		gl.Begin(gl.LINE_LOOP)
		for _, pos := range poly.node.points {
			if pos == nil {
				continue
			}
			
			gl.Vertex2f(float32(pos.x), float32(pos.y))
		}
		gl.End()
	}
	
	/*
	fps := 1/(float64(time.Since(start))/float64(time.Second))
	s.ui.fpsLabel2.SetLabel(fmt.Sprintf("%f", fps))	
	*/
}

func (s *Sim) Run() {
	for s.ui.running {
	
		start := time.Now()

		for e := sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
			switch ev := e.(type) {
			case *sdl.QuitEvent:
				s.ui.running = false
			case *sdl.KeyboardEvent:
				if ev.Keysym.Sym == sdl.K_ESCAPE {
					s.ui.running = false
				}
			/*case *sdl.MouseMotionEvent:
				if ev.State != 0 {
					pen.lineTo(Point{int(ev.X), int(ev.Y)})
				} else {
					pen.moveTo(Point{int(ev.X), int(ev.Y)})
				}*/
			}
		}
		
		s.Draw()

		sdl.GL_SwapBuffers()
		
		fps := 1/(float64(time.Since(start))/float64(time.Second))
		_ = fps
		//fmt.Printf("%f\n", fps)
	}
}

