package cache

// An LFU is a fixed-size in-memory cache with least-frequently-used eviction
type LFU struct {
	// whatever fields you want here
	maxSize int
	currSize int
	stats *Stats
}

// NewLFU returns a pointer to a new LFU with a capacity to store limit bytes
func NewLru(limit int) *LFU {
	cache := new(LFU)
	cache.maxSize = limit
	cache.currSize = 0
	cache.stats = NewStats()
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
	return nil, false
}

// Remove removes and returns the value associated with the given key, if it exists.
// ok is true if a value was found and false otherwise
func (lfu *LFU) Remove(key string) (value []byte, ok bool) {
	return nil, false
}

// Set associates the given value with the given key, possibly evicting values
// to make room. Returns true if the binding was added successfully, else false.
func (lfu *LFU) Set(key string, value []byte) bool {
	return false
}

// Evict the last element added to list
func EvictLFU(lfu *LFU, existsInQ bool) {
}

// Len returns the number of bindings in the LFU.
func (lfu *LFU) Len() int {
	return 0
}

// Stats returns statistics about how many search hits and misses have occurred.
func (lfu *LFU) Stats() *Stats {
	return lfu.stats
}
