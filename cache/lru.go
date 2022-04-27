package cache

import (
	"container/list"
	"log"
)

// An LRU is a fixed-size in-memory cache with least-recently-used eviction
type LRU struct {
	// whatever fields you want here
	lookup       map[string]*[]byte
	stringToNode map[string]*list.Element
	q            *list.List
	maxSize      int
	currSize     int
	stats        *Stats
}

// NewLRU returns a pointer to a new LRU with a capacity to store limit bytes
func NewLru(limit int) *LRU {
	cache := new(LRU)
	cache.lookup = map[string]*[]byte{}
	cache.stringToNode = map[string]*list.Element{}
	cache.q = list.New()
	cache.maxSize = limit
	cache.currSize = 0
	cache.stats = new(Stats)
	cache.stats.Hits = 0
	cache.stats.Misses = 0
	return cache
}

// MaxStorage returns the maximum number of bytes this LRU can store
func (lru *LRU) MaxStorage() int {
	return lru.maxSize
}

// RemainingStorage returns the number of unused bytes available in this LRU
func (lru *LRU) RemainingStorage() int {
	return lru.maxSize - lru.currSize
}

// Get returns the value associated with the given key, if it exists.
// This operation counts as a "use" for that key-value pair
// ok is true if a value was found and false otherwise.
func (lru *LRU) Get(key string) (value []byte, ok bool) {
	valPointer := lru.lookup[key]

	if valPointer == nil {
		lru.stats.Misses++
		return nil, false
	}

	// move matching element to front of queue

	currEl := lru.stringToNode[key]
	lru.q.MoveToFront(currEl)

	lru.stats.Hits++
	return *valPointer, true
}

// Remove removes and returns the value associated with the given key, if it exists.
// ok is true if a value was found and false otherwise
func (lru *LRU) Remove(key string) (value []byte, ok bool) {
	prevHits, prevMisses := lru.stats.Hits, lru.stats.Misses
	val, found := lru.Get(key)

	// ensure no changes to stats with Get
	lru.stats.Hits = prevHits
	lru.stats.Misses = prevMisses

	if !found {
		return nil, false
	}

	delete(lru.lookup, key)

	// remove matching element from queue
	currEl := lru.stringToNode[key]
	lru.q.Remove(currEl)

	delete(lru.stringToNode, key)

	lru.currSize -= len(key) + len(val)
	return val, true
}

// Set associates the given value with the given key, possibly evicting values
// to make room. Returns true if the binding was added successfully, else false.
func (lru *LRU) Set(key string, value []byte) bool {
	// Check to see if too large for cache
	newElSize := len(key) + len(value)
	if newElSize > lru.maxSize {
		return false
	}

	// Check to see if we're updating an existing key or adding a new key
	existsInQ := false
	existingVal := lru.lookup[key]
	addedSize := newElSize
	if existingVal != nil {
		existsInQ = true
		addedSize = len(value) - len(*existingVal)
	}

	// Evict until there's enough room
	for lru.currSize+addedSize > lru.maxSize {
		EvictLRU(lru, existsInQ)
	}

	// Add new key:value pair
	if existsInQ {
		currEl := lru.stringToNode[key]
		lru.q.Remove(currEl)
		lru.q.PushFront(key)
		lru.stringToNode[key] = lru.q.Front()
		lru.currSize += addedSize
		// Only add to the queue if it doesn't exist yet
	} else {
		lru.q.PushFront(key)
		lru.stringToNode[key] = lru.q.Front()
		lru.currSize += newElSize
	}

	// Add/update key in map
	lru.lookup[key] = &value

	return true
}

// Evict the last element added to list
func EvictLRU(lru *LRU, existsInQ bool) {
	backEl := lru.q.Back()
	remKey := backEl.Value.(string)

	// Get the key:value pair from the map
	valPointer := lru.lookup[remKey]

	// Bad News: We're evicting a key that doesn't exist
	if valPointer == nil {
		log.Panic()
	}

	// Remove from map
	delete(lru.lookup, remKey)

	// remove from queue only if we're setting a new value
	// Keep if we're updating a value
	if !existsInQ {
		delete(lru.stringToNode, remKey)
		lru.q.Remove(backEl)
	}

	// change size
	lru.currSize -= len(remKey) + len(*valPointer)
}

// Len returns the number of bindings in the LRU.
func (lru *LRU) Len() int {
	return lru.q.Len()
}

// Stats returns statistics about how many search hits and misses have occurred.
func (lru *LRU) Stats() *Stats {
	return lru.stats
}
