package pf

import (
	"fmt"
	"math"
	"time"
	"math2"
	"geo"
	
	//"sdl/ttf"
	
	"container/list"
	
	"github.com/banthar/gl"
	"github.com/banthar/Go-SDL/sdl"
)


type unit struct {
	id int
	pos geo.Vec
	
	start, end geo.Vec
	
	trail *Ringbuffer
	color *Color
	
	path *list.List
}

type Sim struct {
	units map[int]*unit
	buildings map[int]*building
	static [][]float64
	nav *NavMesh	// List of Polygons

	// temp	
	markedNodes map[*NavNode]bool
	path *list.List
	
	ui *UI
	
	run bool
	running bool
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
		nav: nil,
		markedNodes: make(map[*NavNode]bool),
		running: false,
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

func (s *Sim) Init() {
	
	outers := []*geo.Vec{
		&geo.Vec{0, 0},
		&geo.Vec{299, 0},
		&geo.Vec{299, 299},
		&geo.Vec{0, 299},
	}
	startPoly := &geo.Polygon{outers}
	
	polySoup := list.New() 
	
	for _, building := range s.buildings {
		inners := []*geo.Vec{
			&geo.Vec{building.x+building.w, building.y+building.h},
			&geo.Vec{building.x, building.y+building.h},
			&geo.Vec{building.x, building.y},
			&geo.Vec{building.x+building.w, building.y},
		}
		innerPoly := &geo.Polygon{inners}
		
		for outer_i, outer := range outers {
			doBreak := false
			for inner_i, inner := range inners {
				//fmt.Println(outer, inner)
				pts := innerPoly.IntersectLine(outer, inner.Sub(outer)).Elements()
				if len(pts) > 0 {
					if len(pts) > 1 || pts[0].(geo.Vec) != *inner {
						break
					}
				}
			
				if startPoly.ContainsCornerLine(outer_i, inner) {
					
					newPoly := make([]*geo.Vec, len(inners)+len(outers)+2)
					i := 0
					for ; i<=len(outers) ; i++ {
						newPoly[i] = outers[(outer_i+i)%len(outers)]
					}
					
					for ii:=0; i<=len(outers) + len(inners); ii++ {
						offset := (inner_i+len(inners)-ii)%len(inners)
						newPoly[i] = inners[offset]
						i++
					}
					
					newPoly[i] = inners[0]
					
					// Triangulate the concave polygon
					tris := (&geo.Polygon{newPoly}).Triangulate()
					polySoup.PushBackList(tris)
						
					// We were able to insert the building, so break		
					doBreak = true
					break
				}
			}
			
			if doBreak {
				break
			}
		}
	}

	s.nav = NewNavMesh(polySoup)	
	s.nav.Reduce()
	
	//start := s.nav.nodes.Front().Value.(*NavNode)
	//end := s.nav.nodes.Back().Value.(*NavNode)
	//path := FindPath(start, end)
	//FindPath2(s, &geo.Vec{250, 60}, &geo.Vec{20, 60}, path)
	
	/*for e:=path.Front(); e != nil; e = e.Next() {
		nn := e.Value.(*NavNode)
		s.markedNodes[nn] = true
	}*/
	
	//s.path = FindPath(s,  &geo.Vec{250, 60}, &geo.Vec{20, 60})
}

// TODO This should go to a separate Map struct
// TODO What happens with the point-on-line intesection?
func (s *Sim) IntersectLine(ap, av *geo.Vec) bool {
	for _, building := range s.buildings {
		poly := &geo.Polygon{[]*geo.Vec{
			&geo.Vec{building.x+building.w, building.y+building.h},
			&geo.Vec{building.x, building.y+building.h},
			&geo.Vec{building.x, building.y},
			&geo.Vec{building.x+building.w, building.y},
		}}
		
		// if there is only one intersection, the line is only 
		if poly.IntersectLine(ap, av).Len() > 1 {
			return true
		}
	}
	
	return false
}

func (s *Sim) AddUnit(start, end geo.Vec, color *Color) {
	id := len(s.units)
	unit := &unit{
		id: id,
		pos: start,
		start: start,
		end: end,
		trail: NewRingbuffer(5),
		color: color,
		path: FindPath(s, &start, &end),
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
				
				pot := s.potential(unit, unit.pos.X+dx, unit.pos.Y+dy)
				
				if(pot < min) {
					min = pot
					radMin = rad
				}
			}
			
			d := speed*float64(dTime)
			unit.pos.X += math.Cos(radMin)*d
			unit.pos.Y += math.Sin(radMin)*d
			
			// Add position to trail if more than 0.25 units away
			if last := unit.trail.Front(); last != nil {
				lastPos := last.(geo.Vec)
				if(lastPos.Distance(&unit.pos) >= 0.25) {
					unit.trail.AddToFront(unit.pos)
				}
			} else {
				unit.trail.AddToFront(unit.pos)
			}
			
			// if the unit already walks at the next path piece, forget the current front one
			if unit.path.Len() > 2 {
				pt1 := unit.path.Front().Next().Value.(*geo.Vec)
				pt2 := unit.path.Front().Next().Next().Value.(*geo.Vec)
				if geo.OrthogonalVector(pt1, pt2.Sub(pt1), &unit.pos) != nil {
					unit.path.Remove(unit.path.Front())
				}
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
	//max := 300.0
	//dist := (&geo.Vec{x,y}).Distance(&unit.end)
	
	//potEnd := dist/max
	potEnd := 1.0
	pos := &geo.Vec{x,y}
	
	if unit.path != nil {
		for e := unit.path.Front(); e.Next() != nil; e = e.Next() {
			pt0 := e.Value.(*geo.Vec)
			pt1 := e.Next().Value.(*geo.Vec)
			
			av := geo.OrthogonalVector(pt0, pt1.Sub(pt0), &geo.Vec{x, y})
			if av != nil {
				potEnd = math.Min(av.Len()/10.0, 1)
			} else {
				if l := pt0.Sub(pos).Len(); l <= 10 {
					potEnd = math.Min(l/10.0, potEnd)	
				} else if l := pt1.Sub(pos).Len(); l <= 10 {
					potEnd = math.Min(l/10.0, potEnd)	
				}
			}
			
		}
	}
	
	// other units
	MIN := 5.0
	MAX := 10.0
	potDist := 0.0
	for _, oUnit := range s.units {
		if(oUnit.id != unit.id) {
			dist := oUnit.pos.Distance(&geo.Vec{x,y})
			potDist += 1-math2.MinMax(1, (dist-MIN)/(MAX-MIN), 0)
		}		
	}
	
	potDist = math2.MinMax(1, potDist, 0)
	
	// trail
	potTrail := 0.0
	elms := unit.trail.Elements();
	for _, elm := range elms {
		if(elm != nil) {
			trailPos := elm.(geo.Vec)
			dist := trailPos.Distance(&geo.Vec{x,y})
			potTrail += 1.0-math2.MinMax(1, dist/float64(1), 0)
		}		
	}
	
	potTrail = math2.MinMax(1, potTrail, 0)
	
	if unit.id == 0  && potDist > 0 {
		//fmt.Println(potDist)
	}

	pot := 0.5*potEnd+0*potDist + 0.25*potTrail
	
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
	gl.Begin(gl.POINTS)	
	for i, unit := range s.units {
		if i == 0 {
			gl.Color4f(1, 0, 0, 1)	
			fmt.Println(unit.pos.X)
		} else {
			gl.Color4f(0, 1, 0, 1)
		}
		
		gl.Vertex2f(float32(unit.pos.X), float32(unit.pos.Y))
	}
	gl.End()
	
	// Nav mesh
	/*
	for e := s.nav.nodes.Front(); e != nil; e = e.Next() {
		nn := e.Value.(*NavNode)

		gl.Color4f(0, 0, 1, 0.2)
		
		if s.markedNodes[nn] {
			gl.Color4f(1, 0, 0, 0.2)
		}
		
		gl.Begin(gl.POLYGON)
		for _, pos := range nn.node.Points {
			if pos == nil {
				continue
			}
			
			gl.Vertex2f(float32(pos.X), float32(pos.Y))
		}
		gl.End()
		
		gl.Color4f(0, 0, 1, 1)
		gl.Begin(gl.LINE_LOOP)
		for _, pos := range nn.node.Points {
			if pos == nil {
				continue
			}
			
			gl.Vertex2f(float32(pos.X), float32(pos.Y))
		}
		gl.End()
	}
	*/
	// Draw Path
	if s.units[0].path != nil {
		gl.Begin(gl.LINE_STRIP)
		gl.Color4f(1, 0, 0, 0.2)
		for e := s.units[0].path.Front(); e != nil; e = e.Next() {
			pos := e.Value.(*geo.Vec)
			gl.Vertex2f(float32(pos.X), float32(pos.Y))
		}
		gl.End()
	}
		
	// Draw Links
	/*gl.Color4f(1, 0, 0, 1)
	gl.Begin(gl.LINES)
	for i, link := range nn.links {
		pt1 := nn.node.Points[i]
		pt2 := nn.node.Points[(i+1)%nn.node.Len()]
		lineCenter := &geo.Vec{(pt1.X+pt2.X)/2, (pt1.Y+pt2.Y)/2}
		
		center := link.node.Center()
		
		gl.Vertex2d(lineCenter.X, lineCenter.Y)
		gl.Vertex2d(center.X, center.Y)
	}
	gl.End()*/
	
	
	
	
	
	
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
				} else if ev.Keysym.Sym == sdl.K_SPACE {
					if !s.running {
						s.running = true
						go s.Update()
					}
				}
			}
		}
		
		
		s.Draw()	
	
		sdl.GL_SwapBuffers()
		
		fps := 1/(float64(time.Since(start))/float64(time.Second))
		_ = fps
		//fmt.Printf("%f\n", fps)
	}
}

