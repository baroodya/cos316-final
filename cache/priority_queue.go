// This example demonstrates a priority queue built using the heap interface.
// This is adapted from https://golang.org/pkg/container/heap/#example__priorityQueue
package cache

import (
	"container/heap"
	"fmt"
)

// An Item is something we manage in a priority queue.
type Item struct {
	key    string // The value of the item; arbitrary.
	priority int    // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index    int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the item with the lowest priority so we use less than here.
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func (pq *PriorityQueue) Remove(item *Item) {
	for i := 0; i < len(*pq); i++ {
		if (*pq)[i] == item {
			heap.Remove(pq, i)
			return
		}
	}
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) Update(item *Item, priority int) {
	item.priority = priority

	// // Print the priority queue.
	// fmt.Println("Priority queue:")
	// for i := 0; i < pq.Len(); i++ {
	// 	fmt.Printf("%d -> %.2d:%s\n", (*pq)[i].index, (*pq)[i].priority, (*pq)[i].key)
	// }
	// fmt.Println()
	heap.Fix(pq, item.index)

	// fmt.Println("Updated priority queue:")
	// // Print the priority queue.
	// for i := 0; i < pq.Len(); i++ {
	// 	fmt.Printf("%d -> %.2d:%s\n", (*pq)[i].index, (*pq)[i].priority, (*pq)[i].key)
	// }
	// fmt.Println()
}

// This example creates a PriorityQueue with some items, adds and manipulates an item,
// and then removes the items in priority order.
func main() {
	// Some items and their priorities.
	items := map[string]int{
		"banana": 3, "apple": 2, "pear": 4,
	}

	// Create a priority queue, put the items in it, and
	// establish the priority queue (heap) invariants.
	pq := make(PriorityQueue, len(items))
	i := 0
	for value, priority := range items {
		pq[i] = &Item{
			key: value,
			priority: priority,
			index:    i,
		}
		i++
	}
	heap.Init(&pq)

	// Insert a new item and then modify its priority.
	item := &Item{
		key: "orange",
		priority: 1,
	}
	heap.Push(&pq, item)
	pq.Update(item, 5)

	// Take the items out; they arrive in decreasing priority order.
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		fmt.Printf("%.2d:%s ", item.priority, item.key)
	}
}