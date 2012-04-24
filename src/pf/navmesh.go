package pf

import (
	"container/list"	
	"fmt"
	"geo"
	"graph"
)

type NavMesh struct {
	nodes *list.List
}

type NavNode struct {
	node *geo.Polygon
	links map[int]*NavNode
}

type line struct {
	p1, p2 geo.Vec
}

type adjacent struct {
	a *NavNode
	a_i int
	b *NavNode
	b_i int
}

func NewNavMesh(soup *list.List) *NavMesh {
	nm := &NavMesh{
		list.New(),
	}
	
	lines := make(map[line]*adjacent)
	
	for e := soup.Front(); e != nil; e = e.Next() {
		nn := NewNavNode(e.Value.(*geo.Polygon))
		nm.nodes.PushBack(nn)
		
		for i, pt1 := range nn.node.Points {
			pt2 := nn.node.Points[(i+1)%nn.node.Len()]
			ll := line{*pt1, *pt2}
			if pt2.X < pt1.X {
				ll = line{*pt2, *pt1}
			}
			//fmt.Println("Checking ", lines[ll])
			
			if lines[ll] == nil {
				lines[ll] = &adjacent{nn, i, nil, 0}
			} else {
				lines[ll].b = nn
				lines[ll].b_i = i
			}
		}
	}
	
	
	for _, adj := range lines {
		if adj.b != nil {
			adj.a.links[adj.a_i] = adj.b
			adj.b.links[adj.b_i] = adj.a
		}
	}
	
	return nm
}

func NewNavNode(node *geo.Polygon) *NavNode {
	return &NavNode{
		node: node,
		links: make(map[int]*NavNode),
	}
}

func (nn *NavNode) String() string {
	str := fmt.Sprint(nn.node)
	str += fmt.Sprint(" links:")
	for i, link := range nn.links {
		str += fmt.Sprint(i,":")+fmt.Sprint(link.node)
	}
	return str
}

func (nn *NavNode) C(other graph.Node) float64 {
	return nn.node.Center().Distance(other.(*NavNode).node.Center())
}

func (nn *NavNode) GetLinks() []graph.Node {
	nodes := make([]graph.Node, len(nn.links))
	i:= 0
	for _, node := range nn.links {
		nodes[i] = node
		i++
	}
	return nodes;
}

func (nn *NavNode) HasCorner(pos *geo.Vec) int {
	for i, p := range nn.node.Points {
		if *p == *pos {
			return i
		}
	}
	
	return -1
}

func (nn *NavNode) Merge(other *NavNode) *NavNode {
	_ = fmt.Println
	
	for this_i, this_corner := range nn.node.Points {
		if other_i := other.HasCorner(this_corner); other_i >= 0 {
			
			newPoly := &geo.Polygon{
				make([]*geo.Vec, nn.node.Len()+other.node.Len()-2),
			}
			
			newLinks := make(map[int]*NavNode)
			
			i := 0
			for ; i<other.node.Len(); i++ {
				offset := (other_i+i)%other.node.Len()
				newPoly.Points[i] = other.node.Points[offset]
				
				if i<other.node.Len()-1 && other.links[offset] != nil {
					newLinks[i] = other.links[offset]
				}
			}
			for ; i < nn.node.Len()+other.node.Len()-2; i++ {
				offset := (i-other.node.Len()+this_i+2)%nn.node.Len()
				linkOffset := (i-other.node.Len()+this_i+2-1)%nn.node.Len()
				newPoly.Points[i] = nn.node.Points[offset]
				
				if nn.links[linkOffset] != nil {
					newLinks[i-1] = nn.links[linkOffset]
				}
			}
			
			return &NavNode{
				node: newPoly,
				links: newLinks,
			}
		}
	}
	
	return nil
}

func (nm *NavMesh) Reduce() {
	reduced := make(map[*NavNode]bool)
	nm.reduce(&reduced, nm.nodes.Front().Value.(*NavNode))
}

func (nm *NavMesh) reduce(reduced *map[*NavNode]bool, nn *NavNode) {
	// Try to reduce this node
	for canReduce := true; canReduce ; {
		for _, candidate := range nn.links {
			testNN := nn.Merge(candidate)
			if testNN.node.IsConvex() {
				canReduce = true
				
				//elm.Value = testNN
				nn.node = testNN.node
				nn.links = testNN.links
				
				// When we merge two nodes, one node (var candidate) will be
				// deleted and does not exist any more. But still, all neighbors of 
				// candidate (candidate.links) have have stored a reference to candidate:
				// candidate.links[ONE_EDGE].links[ONE_EDGE] == candidate
				// This reference must be replaced ...
				// TODO I think the while issue can be solved more easily, if we use 
				// double pointers to store the navnodes. Then we only have to make the pointer
				// of candidate, point to nn.node
				for _, link := range candidate.links {
					for link2_i, link2 := range link.links {
						if link2.node == candidate.node {
							link.links[link2_i] = nn
						}
					}
				}
				nm.RemoveNode(candidate)
				
				break
			}
		}
		canReduce = false
	}
	
	(*reduced)[nn] = true
	
	// Reduce the links, if any
	for _, link := range nn.links {
		if !(*reduced)[link] {
			nm.reduce(reduced, link)
		}
	}
}

// TODO A GOOD implementation of a list would be nice, with find, contains etc ...
func (nm *NavMesh) RemoveNode(search *NavNode) bool {
	for e := nm.nodes.Front(); e != nil; e = e.Next() {
		val := e.Value.(*NavNode)
		if val == search {
			nm.nodes.Remove(e)
			return true
		}
	}
	
	return false
}

func (nm *NavMesh) NodeAtPoint(pt *geo.Vec) *NavNode {
	for e := nm.nodes.Front(); e != nil; e = e.Next() {
		node := e.Value.(*NavNode)
		if node.node.ContainsPoint(pt) {
			return node
		}
	}
	
	return nil
}

