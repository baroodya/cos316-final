package cache

import (
	"container/heap"
)

// An LinearLFU is a fixed-size in-memory cache with least-frequently-used eviction
type LinearLFU struct {
	// whatever fields you want here
	pq       PriorityQueue
	lookup   map[string]*[]byte
	items    map[string]*Item
	maxSize  int
	currSize int
	stats    *Stats

	alpha         float64
	cacheAccesses int
}

// NewLinearLFU returns a pointer to a new LinearLFU with a capacity to store limit bytes
func NewLinearLfu(limit int) *LinearLFU {
	cache := new(LinearLFU)

	cache.lookup = map[string]*[]byte{}
	cache.items = map[string]*Item{}

	cache.pq = make(PriorityQueue, 0)
	heap.Init(&cache.pq)

	cache.maxSize = limit
	cache.currSize = 0
	cache.stats = new(Stats)
	cache.stats.Hits = 0
	cache.stats.Misses = 0

	// Constant multiplier for the priority of a key
	cache.alpha = 1.0
	cache.cacheAccesses = 0
	return cache
}

// MaxStorage returns the maximum number of bytes this LinearLFU can store
func (lfu *LinearLFU) MaxStorage() int {
	return lfu.maxSize
}

// RemainingStorage returns the number of unused bytes available in this LinearLFU
func (lfu *LinearLFU) RemainingStorage() int {
	return lfu.maxSize - lfu.currSize
}

// Get returns the value associated with the given key, if it exists.
// This operation counts as a "use" for that key-value pair
// ok is true if a value was found and false otherwise.
func (lfu *LinearLFU) Get(key string) (value []byte, ok bool) {
	lfu.cacheAccesses++
	valPointer := lfu.lookup[key]

	if valPointer == nil {
		lfu.stats.Misses++
		return nil, false
	}

	itemPointer := lfu.items[key]
	item := *itemPointer

	// update priority of element in priority queue
	newPriority := lfu.getLinearPriority(item.accesses)
	lfu.pq.Update(itemPointer, newPriority)
	item.priority = newPriority
	item.accesses++
	lfu.items[key] = itemPointer

	// // Print the lookup table
	// for k, v := range lfu.lookup {
	// 	fmt.Printf("%s: %s. Accesses: %d\n", k, string(*v), lfu.items[k].priority)
	// }
	// fmt.Println()

	lfu.stats.Hits++
	return *valPointer, true
}

// priority = alpha * cache accesses + key accesses
func (lfu *LinearLFU) getLinearPriority(accesses int) float64 {
	return lfu.alpha*float64(lfu.cacheAccesses) + float64(accesses)
}

// Remove removes and returns the value associated with the given key, if it exists.
// ok is true if a value was found and false otherwise
func (lfu *LinearLFU) Remove(key string) (value []byte, ok bool) {
	valPointer := lfu.lookup[key]

	if valPointer == nil {
		return nil, false
	}

	delete(lfu.lookup, key)

	// remove matching element from priority queue
	itemPointer := lfu.items[key]
	lfu.pq.Remove(itemPointer)

	delete(lfu.items, key)

	lfu.currSize -= len(key) + len(*valPointer)
	return *valPointer, true
}

// Set associates the given value with the given key, possibly evicting values
// to make room. Returns true if the binding was added successfully, else false.
func (lfu *LinearLFU) Set(key string, value []byte) bool {
	lfu.cacheAccesses++

	// Check to see if too large for cache
	newElSize := len(key) + len(value)
	if newElSize > lfu.maxSize {
		return false
	}

	// // print key value pair
	// fmt.Printf("%s: %s\n", key, value)
	// fmt.Println("------------------------------------------------------")

	// Check to see if we're updating an existing key or adding a new key
	existsInQ := false
	existingVal := lfu.lookup[key]
	addedSize := newElSize
	if existingVal != nil {
		existsInQ = true
		addedSize = len(value) - len(*existingVal)
	}

	// Evict until there's enough room
	for lfu.currSize+addedSize > lfu.maxSize {
		EvictLinearLFU(lfu)
	}

	// // Print the lookup table
	// for k, v := range lfu.lookup {
	// 	fmt.Printf("%s: %s. Accesses: %d\n", k, string(*v), lfu.items[k].priority)
	// }
	// fmt.Println()

	// Add new key:value pair
	if existsInQ {
		existingItemPointer := lfu.items[key]
		existingItem := *existingItemPointer
		lfu.pq.Update(existingItemPointer, existingItem.priority+1)
		existingItem.priority++
		lfu.lookup[key] = existingVal
		lfu.items[key] = &existingItem
		lfu.currSize += addedSize
		// Only add to the queue if it doesn't exist yet
	} else {
		item := &Item{
			key:      key,
			priority: 1.0,
			accesses: 1,
		}

		heap.Push(&lfu.pq, item)
		lfu.lookup[key] = &value
		lfu.items[key] = item
		lfu.currSize += newElSize
	}

	return true
}

// Evict the last element added to list
func EvictLinearLFU(lfu *LinearLFU) {
	item := heap.Pop(&lfu.pq).(*Item)
	// fmt.Printf("Evicting %s: %s. Accesses: %d\n", item.key, string(*lfu.lookup[item.key]), item.priority)
	key := item.key
	value := *(lfu.lookup[key])
	delete(lfu.lookup, key)
	delete(lfu.items, key)
	lfu.currSize -= len(key) + len(value)
}

// Len returns the number of bindings in the LinearLFU.
func (lfu *LinearLFU) Len() int {
	return lfu.pq.Len()
}

// Stats returns statistics about how many search hits and misses have occurred.
func (lfu *LinearLFU) Stats() *Stats {
	return lfu.stats
}
