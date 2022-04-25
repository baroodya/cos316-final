package cache

import (
	"container/heap"
	"math"
)

// An LogLFU is a fixed-size in-memory cache with least-frequently-used eviction
type LogLFU struct {
	// whatever fields you want here
	pq PriorityQueue
	lookup map[string]*[]byte
	items map[string]*Item
	maxSize int
	currSize int
	stats *Stats

	alpha float64
	beta float64
	cacheAccesses int
}

// NewLogLFU returns a pointer to a new LogLFU with a capacity to store limit bytes
func NewLfu(limit int) *LogLFU {
	cache := new(LogLFU)

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
	cache.alpha = 0.5
	// Constant Base for the log operation
	cache.beta = 2.0
	cache.cacheAccesses = 0
	return cache
}

// MaxStorage returns the maximum number of bytes this LogLFU can store
func (lfu *LogLFU) MaxStorage() int {
	return lfu.maxSize
}

// RemainingStorage returns the number of unused bytes available in this LogLFU
func (lfu *LogLFU) RemainingStorage() int {
	return lfu.maxSize - lfu.currSize
}

// Get returns the value associated with the given key, if it exists.
// This operation counts as a "use" for that key-value pair
// ok is true if a value was found and false otherwise.
func (lfu *LogLFU) Get(key string) (value []byte, ok bool) {
	lfu.cacheAccesses++
	valPointer := lfu.lookup[key]

	if valPointer == nil {
		lfu.stats.Misses++
		return nil, false
	}

	itemPointer := lfu.items[key]
	item := *itemPointer

	// update priority of element in priority queue
	newPriority := lfu.getLogPriority(item.accesses)
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

// priority = alpha * log_beta(cache accesses) * (key accesses)
func (lfu *LogLFU) getLogPriority(accesses float64) float64 {
	changeOfBase := math.Log1p(float64(lfu.cacheAccesses)) / math.Log1p(lfu.beta)
	return lfu.alpha * changeOfBase * accesses
}

// Remove removes and returns the value associated with the given key, if it exists.
// ok is true if a value was found and false otherwise
func (lfu *LogLFU) Remove(key string) (value []byte, ok bool) {
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
func (lfu *LogLFU) Set(key string, value []byte) bool {
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
		EvictLogLFU(lfu)
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
		lfu.pq.Update(existingItemPointer, existingItem.priority + 1)
		existingItem.priority++
		lfu.lookup[key] = existingVal
		lfu.items[key] = &existingItem
		lfu.currSize += addedSize
	// Only add to the queue if it doesn't exist yet
	} else {
		item := &Item{
			key: key, 
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
func EvictLogLFU(lfu *LogLFU) {
	item := heap.Pop(&lfu.pq).(*Item)
	// fmt.Printf("Evicting %s: %s. Accesses: %d\n", item.key, string(*lfu.lookup[item.key]), item.priority)
	key := item.key
	value := *(lfu.lookup[key])
	delete(lfu.lookup, key)
	delete(lfu.items, key)
	lfu.currSize -= len(key) + len(value)
}

// Len returns the number of bindings in the LogLFU.
func (lfu *LogLFU) Len() int {
	return lfu.pq.Len()
}

// Stats returns statistics about how many search hits and misses have occurred.
func (lfu *LogLFU) Stats() *Stats {
	return lfu.stats
}
