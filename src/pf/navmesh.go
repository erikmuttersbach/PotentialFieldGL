package pf

import (
	"container/list"	
	"fmt"
)

type NavMesh struct {
	nodes *list.List
}

type NavNode struct {
	node *Polygon
	links map[int]*NavNode
}

func NewNavMesh() *NavMesh {
	return &NavMesh{
		list.New(),
	}
}

func NewNavNode(node *Polygon) *NavNode {
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

func (nn *NavNode) HasCorner(pos *Pos) int {
	for i, p := range nn.node.points {
		if *p == *pos {
			return i
		}
	}
	
	return -1
}

func (nn *NavNode) Merge(other *NavNode) *NavNode {
	_ = fmt.Println
	
	for this_i, this_corner := range nn.node.points {
		if other_i := other.HasCorner(this_corner); other_i >= 0 {
			
			newPoly := &Polygon{
				make([]*Pos, nn.node.Len()+other.node.Len()-2),
			}
			
			newLinks := make(map[int]*NavNode)
			
			i := 0
			for ; i<other.node.Len(); i++ {
				offset := (other_i+i)%other.node.Len()
				newPoly.points[i] = other.node.points[offset]
				
				if i<other.node.Len()-1 && other.links[offset] != nil {
					newLinks[i] = other.links[offset]
				}
			}
			for ; i < nn.node.Len()+other.node.Len()-2; i++ {
				offset := (i-other.node.Len()+this_i+2)%nn.node.Len()
				linkOffset := (i-other.node.Len()+this_i+2-1)%nn.node.Len()
				newPoly.points[i] = nn.node.points[offset]
				
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
	fmt.Println("\nReducing", nn.node)
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

