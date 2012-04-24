package pf

import (
	"common/pq"
	"container/heap"
	"container/list"
	"fmt"
	"graph"
)

func FindPath(start, end graph.Node) *list.List {
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

