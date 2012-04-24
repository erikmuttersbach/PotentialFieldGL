// A priority queue
package pq

import (
    "container/heap"
)

// An Item is something we manage in a priority queue.
type Item struct {
    Value    interface{} // The value of the item; arbitrary.
    Priority float64    // The priority of the item in the queue.
    // The index is needed by changePriority and is maintained by the heap.Interface methods.
    index int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
    // We want Pop to give us the highest, not lowest, priority so we use greater than here.
    return pq[i].Priority < pq[j].Priority
}

func (pq PriorityQueue) Swap(i, j int) {
    pq[i], pq[j] = pq[j], pq[i]
    pq[i].index = i
    pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
    // Push and Pop use pointer receivers because they modify the slice's length,
    // not just its contents.
    // To simplify indexing expressions in these methods, we save a copy of the
    // slice object. We could instead write (*pq)[i].
    a := *pq
    n := len(a)
    a = a[0 : n+1]
    item := x.(*Item)
    item.index = n
    a[n] = item
    *pq = a
}

func (pq *PriorityQueue) Pop() interface{} {
    a := *pq
    n := len(a)
    item := a[n-1]
    item.index = -1 // for safety
    *pq = a[0 : n-1]
    return item
}

// update is not used by the example but shows how to take the top item from
// the queue, update its priority and value, and put it back.
/*func (pq *PriorityQueue) update(value interface{}, priority float64) {
    item := heap.Pop(pq).(*Item)
    item.Value = value
    item.Priority = priority
    heap.Push(pq, item)
}*/

// changePriority is not used by the example but shows how to change the
// priority of an arbitrary item.
func (pq *PriorityQueue) ChangePriority(item *Item, priority float64) {
    heap.Remove(pq, item.index)
    item.Priority = priority
    heap.Push(pq, item)
}

// Returns the item for a value in the list
func (pq *PriorityQueue) GetItem(search interface{}) *Item {
	for _, item  := range *pq {
		if(item.Value == search) {
			return item
		}
	}
	return nil
}

func (pq *PriorityQueue) Contains(search interface{}) bool {
	return pq.GetItem(search) != nil
}

func New(len int) PriorityQueue {
	return make(PriorityQueue, 0, len)
}