package cache

import (
	"container/heap"
)

// An LFU is a fixed-size in-memory cache with least-frequently-used eviction
type LFU struct {
	// whatever fields you want here
	pq PriorityQueue
	lookup map[string]*[]byte
	items map[string]*Item
	maxSize int
	currSize int
	stats *Stats
}

// NewLFU returns a pointer to a new LFU with a capacity to store limit bytes
func NewLfu(limit int) *LFU {
	cache := new(LFU)

	cache.lookup = map[string]*[]byte{}
	cache.items = map[string]*Item{}

	cache.pq = make(PriorityQueue, 0)
	heap.Init(&cache.pq)

	cache.maxSize = limit
	cache.currSize = 0
	cache.stats = new(Stats)
	cache.stats.Hits = 0
	cache.stats.Misses = 0
	return cache
}

// MaxStorage returns the maximum number of bytes this LFU can store
func (lfu *LFU) MaxStorage() int {
	return lfu.maxSize
}

// RemainingStorage returns the number of unused bytes available in this LFU
func (lfu *LFU) RemainingStorage() int {
	return lfu.maxSize - lfu.currSize
}

// Get returns the value associated with the given key, if it exists.
// This operation counts as a "use" for that key-value pair
// ok is true if a value was found and false otherwise.
func (lfu *LFU) Get(key string) (value []byte, ok bool) {
	valPointer := lfu.lookup[key]

	if valPointer == nil {
		lfu.stats.Misses++
		return nil, false
	}

	itemPointer := lfu.items[key]
	item := *itemPointer

	// update priority of element in priority queue
	lfu.pq.Update(itemPointer, item.priority + 1)
	item.priority++
	lfu.items[key] = itemPointer

	// // Print the lookup table
	// for k, v := range lfu.lookup {
	// 	fmt.Printf("%s: %s. Accesses: %d\n", k, string(*v), lfu.items[k].priority)
	// }
	// fmt.Println()

	lfu.stats.Hits++
	return *valPointer, true
}

// Remove removes and returns the value associated with the given key, if it exists.
// ok is true if a value was found and false otherwise
func (lfu *LFU) Remove(key string) (value []byte, ok bool) {
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
func (lfu *LFU) Set(key string, value []byte) bool {
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
			EvictLFU(lfu)
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
				priority: 1,
			}

			heap.Push(&lfu.pq, item)
			lfu.lookup[key] = &value
			lfu.items[key] = item
			lfu.currSize += newElSize
		}
	
		return true
}

// Evict the last element added to list
func EvictLFU(lfu *LFU) {
	item := heap.Pop(&lfu.pq).(*Item)
	// fmt.Printf("Evicting %s: %s. Accesses: %d\n", item.key, string(*lfu.lookup[item.key]), item.priority)
	key := item.key
	value := *(lfu.lookup[key])
	delete(lfu.lookup, key)
	delete(lfu.items, key)
	lfu.currSize -= len(key) + len(value)
}

// Len returns the number of bindings in the LFU.
func (lfu *LFU) Len() int {
	return lfu.pq.Len()
}

// Stats returns statistics about how many search hits and misses have occurred.
func (lfu *LFU) Stats() *Stats {
	return lfu.stats
}
