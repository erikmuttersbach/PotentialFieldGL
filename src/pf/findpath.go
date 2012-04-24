package pf

import (
	"common/pq"
	"container/heap"
	"container/list"
	"fmt"
	"graph"
	"geo"
)

func FindPath(s *Sim, start, end *geo.Vec) *list.List {
	startNode := s.nav.NodeAtPoint(start)
	endNode := s.nav.NodeAtPoint(end)
	path1 := findPath1(startNode, endNode)
	return findPath2(s, start, end, path1)
}

func findPath1(start, end graph.Node) *list.List {
	openlist := pq.New(1000) // TODO This is shit
	closedlist := make(map[graph.Node]bool)
	pre := make(map[graph.Node]graph.Node)
	g := make(map[graph.Node]float64)
	
	// Initialisierung der Open List, die Closed List ist noch leer
    // (die Priorität bzw. der f Wert des Startknotens ist unerheblich)
    heap.Push(&openlist, &pq.Item{Value: start, Priority: 0})
    
    // diese Schleife wird durchlaufen bis entweder
    // - die optimale Lösung gefunden wurde oder
    // - feststeht, dass keine Lösung existiert
    for ; openlist.Len() > 0 ; {
    
        // Knoten mit dem geringsten f Wert aus der Open List entfernen
        currentItem := heap.Pop(&openlist).(*pq.Item)
        currentCoord := currentItem.Value.(graph.Node)
        
        //fmt.Println("current coord is", currentCoord)
        
        // Wurde das Ziel gefunden?
        if currentCoord == end {
        	path := list.New()
        	for p := end; p != start; p = pre[p] {
        		path.PushFront(p)
        	}
        	path.PushFront(start)
            
            return path
        }
            
        // Wenn das Ziel noch nicht gefunden wurde: Nachfolgeknoten
        // des aktuellen Knotens auf die Open List setzen
        for _, succCoord := range (currentCoord).GetLinks() {
        
        	//fmt.Println(" Checking", succCoord)
		
			// TODO use IntCoord here
			// wenn der Nachfolgeknoten bereits auf der Closed List ist - tue nichts
	        if closedlist[succCoord] {
	            continue
			}
			
	        // g Wert für den neuen Weg berechnen: g Wert des Vorgängers plus
	        // die Kosten der gerade benutzten Kante
	        tentative_g := g[currentCoord] + (currentCoord).C(succCoord)
	        
	        // wenn der Nachfolgeknoten bereits auf der Open List ist,
	        // aber der neue Weg nicht besser ist als der alte - tue nichts
	        if openlist.Contains(succCoord) && tentative_g >= g[succCoord] {
	            continue
	        }
	        
	        // Vorgängerzeiger setzen und g Wert merken
	        pre[succCoord] = currentCoord
	        g[succCoord] = tentative_g
	        
	        // f Wert des Knotens in der Open List aktualisieren
	        // bzw. Knoten mit f Wert in die Open List einfügen
	        h := (succCoord).C(end)
	        //h = 0
	        f := tentative_g + h
	        if succItem := openlist.GetItem(succCoord); succItem != nil {
	            openlist.ChangePriority(succItem, f)
	        } else {
	        	heap.Push(&openlist, &pq.Item{Value: succCoord, Priority: f})
	        }
	        
	        /*fmt.Print(" openlist is ")
	        for _, item := range openlist {
	        	fmt.Print(item.Value," [",item.Priority,"] ", c(currentCoord, succCoord))
	        }
	        fmt.Println()*/
	        _ = fmt.Println
		}
        
        // der aktuelle Knoten ist nun abschließend untersucht
        closedlist[currentCoord] = true
    }
    
    // die Open List ist leer, es existiert kein Pfad zum Ziel
	return nil
}

func findPath2(s *Sim, start, end *geo.Vec, path1 *list.List) *list.List {
	fmt.Println("Route from", start, end)
	
	la := list.New()
	ra := list.New()
	
	la.PushBack(start)
	ra.PushBack(start)
	
	// Is there a direct connection between start and end?
	if !s.IntersectLine(start, end.Sub(start)) {
		la.PushBack(end)
		return la
	}
	
	for e := path1.Front(); e != nil; e = e.Next() {
		node := e.Value.(*NavNode)
		fmt.Println("Looking at Poly", node)
		
		if e.Next() == nil {
			continue
		}
		
		var l, r *geo.Vec
		
		// Find the "port" to the next poly
		for i, link :=range  node.links {
			if link == e.Next().Value {
				l = node.node.Points[i]
				r = node.node.Points[(i+1)%node.node.Len()]
			}
		}
		
		_ = r
		
		if la.Back().Value != end {
			fmt.Println(" Adding left port", l)
			la.PushBack(l)	
			
			if !s.IntersectLine(l, end.Sub(l)) {
				fmt.Println(" Found end (left)")
				la.PushBack(end)
			}
		}
		
		if ra.Back().Value != end {
			fmt.Println(" Adding right port", r)
			ra.PushBack(r)	
			
			if !s.IntersectLine(r, end.Sub(r)) {
				fmt.Println(" Found end (right)")
				ra.PushBack(end)
			}
		}
	}
	
	lengthl := 0.0
	fmt.Println("left")
	for e := la.Front(); e != nil; e = e.Next() {
		fmt.Println(" ", e.Value)
		
		if e.Next() != nil {
			this := e.Value.(*geo.Vec)
			next := e.Next().Value.(*geo.Vec)
			lengthl += this.Distance(next)
		}
	}
	fmt.Println(" length:", lengthl)
	
	lengthr := 0.0
	fmt.Println("right")
	for e := ra.Front(); e != nil; e = e.Next() {
		fmt.Println(" ", e.Value)
		
		if e.Next() != nil {
			this := e.Value.(*geo.Vec)
			next := e.Next().Value.(*geo.Vec)
			lengthr += this.Distance(next)
		}
	}
	fmt.Println(" length:", lengthr)
	
	if lengthr > lengthl {
		return ra
	}
	
	return la
}