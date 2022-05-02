/******************************************************************************
 * linear_lfu_test.go
 * Author:
 * Usage:    `go test`  or  `go test -v`
 * Description:
 *    An incomplete unit testing suite for linear_lfu.go
 ******************************************************************************/

package cache

import (
	"fmt"
	"testing"
)

/******************************************************************************/
/*                                Constants                                   */
/******************************************************************************/
// Constants can go here

/******************************************************************************/
/*                                  Tests                                     */
/******************************************************************************/

func TestLinearSetGet(t *testing.T) {
	capacity := 64
	lfu := NewLinearLfu(capacity, 1.0)
	checkCapacity(t, lfu, capacity)

	for i := 0; i < 4; i++ {
		key := fmt.Sprintf("key%d", i)
		val := []byte(key)
		ok := lfu.Set(key, val)
		if !ok {
			t.Errorf("Failed to add binding with key: %s", key)
			t.FailNow()
		}

		res, _ := lfu.Get(key)
		if !bytesEqual(res, val) {
			t.Errorf("Wrong value %s for binding with key: %s", res, key)
			t.FailNow()
		}
	}

	// allows empty string as valid key
	key := ""
	val := []byte("val")
	ok := lfu.Set(key, val)
	if !ok {
		t.Errorf("Failed to add binding with key: %s", key)
		t.FailNow()
	}

	res, _ := lfu.Get(key)
	if !bytesEqual(res, val) {
		t.Errorf("Wrong value %s for binding with key: %s", res, key)
		t.FailNow()
	}
}

func TestLinearRemove(t *testing.T) {
	capacity := 64
	lfu := NewLinearLfu(capacity, 1.0)
	checkCapacity(t, lfu, capacity)

	for i := 0; i < 4; i++ {
		key := fmt.Sprintf("key%d", i)
		val := []byte(key)
		ok := lfu.Set(key, val)
		if !ok {
			t.Errorf("Failed to add binding with key: %s", key)
			t.FailNow()
		}

		res, _ := lfu.Get(key)
		if !bytesEqual(res, val) {
			t.Errorf("Wrong value %s for binding with key: %s", res, key)
			t.FailNow()
		}
	}

	for i := 0; i < 4; i++ {
		key := fmt.Sprintf("key%d", i)
		refVal := []byte(key)
		val, ok := lfu.Remove(key)
		if !ok {
			t.Errorf("Failed to remove binding with key: %s", key)
			t.FailNow()
		}

		if !bytesEqual(val, refVal) {
			t.Errorf("Wrong value %s for binding %s with key: %s", val, key, refVal)

			t.FailNow()
		}
	}
}

func TestLinearLen(t *testing.T) {
	// length of empty
	capacity := 100
	lfu := NewLinearLfu(capacity, 1.0)
	len := lfu.Len()

	if len != 0 {
		t.Errorf("Empty LFU does not have length 0, instead has length %d", len)
		t.FailNow()
	}

	// add some and verify length
	key := "Hello"
	val := []byte("World")
	ok := lfu.Set(key, val)
	if !ok {
		t.Errorf("Failed to add binding with key: %s", key)
		t.FailNow()
	}

	len = lfu.Len()
	if len != 1 {
		t.Errorf("LFU does not have length 1, instead has length %d", len)
		t.FailNow()
	}

	// take some out and verify length
	_, ok = lfu.Remove(key)
	if !ok {
		t.Errorf("Failed to remove binding with key: %s", key)
		t.FailNow()
	}
	len = lfu.Len()

	if len != 0 {
		t.Errorf("Empty LFU does not have length 0, instead has length %d", len)
		t.FailNow()
	}
}

func TestLinearMaxStorage(t *testing.T) {
	// set max storage to 100
	capacity := 100
	lfu := NewLinearLfu(capacity, 1.0)
	checkCapacity(t, lfu, capacity)

	// set max storage to 0
	capacity = 0
	lfu = NewLinearLfu(capacity, 1.0)
	checkCapacity(t, lfu, capacity)

	// set max storage to positive val
	capacity = 1024
	lfu = NewLinearLfu(capacity, 1.0)
	checkCapacity(t, lfu, capacity)
}

func TestLinearRemainingStorage(t *testing.T) {
	// remaining storage before adding
	capacity := 10
	lfu := NewLinearLfu(capacity, 1.0)
	rem := lfu.RemainingStorage()

	if rem != capacity {
		t.Errorf("%d of remaining storage in empty cache, should be %d", rem, capacity)
		t.FailNow()
	}

	// remaining storage after adding
	key := "12345"
	val := []byte(key)
	ok := lfu.Set(key, val)
	if !ok {
		t.Errorf("Failed to add binding with key: %s", key)
		t.FailNow()
	}

	rem = lfu.RemainingStorage()
	if rem != 0 {
		t.Errorf("Remaining storage should be 0 for a full cache but is %d", rem)
		t.FailNow()
	}

	// remaining storage after removing
	_, ok = lfu.Remove(key)
	if !ok {
		t.Errorf("Failed to remove binding with key: %s", key)
		t.FailNow()
	}

	rem = lfu.RemainingStorage()
	if rem != capacity {
		t.Errorf("%d of remaining storage in empty cache, should be %d", rem, capacity)
		t.FailNow()
	}
}

func TestLinearZeroCapacity(t *testing.T) {
	capacity := 0
	lfu := NewLinearLfu(capacity, 1.0)
	checkCapacity(t, lfu, capacity)

	// check Get() returns no binding when called on empty cache
	key := "key"
	_, found := lfu.Get(key)
	if found {
		t.Errorf("Inaccurately found binding with key: %s", key)
		t.FailNow()
	}

	cacheMisses := lfu.Stats().Misses
	cacheHits := lfu.Stats().Hits
	if cacheMisses != 1 || cacheHits != 0 {
		t.Errorf("Incorrect cache stats.\n Cache Hits: %d\n Cache Misses: %d\n", cacheHits, cacheMisses)
	}

	// Set() only allows zero-size bindings in a zero-capacity cache
	key = "hello"
	value := []byte("world")
	ok := lfu.Set(key, value)
	if ok {
		t.Errorf("Should have failed to add binding with key: %s", key)
		t.FailNow()
	}

	key = ""
	value = []byte("")
	ok = lfu.Set(key, value)
	if !ok {
		t.Errorf("Failed to add binding with key: %s", key)
		t.FailNow()
	}

	_, found = lfu.Get(key)
	if !found {
		t.Errorf("Failed to find binding with key: %s", key)
		t.FailNow()
	}
}

func TestLinearTooLarge(t *testing.T) {
	capacity := 10
	lfu := NewLinearLfu(capacity, 1.0)

	// set rejects bindings too large for cache
	key := "123456"
	value := []byte(key)
	ok := lfu.Set(key, value)
	if ok {
		t.Errorf("Should have failed to add binding with key: %s", key)
		t.FailNow()
	}

	// ensure cache still empty
	len := lfu.Len()
	if len != 0 {
		t.Errorf("Cache should be empty but has length %d", len)
	}
	rem := lfu.RemainingStorage()
	if rem != capacity {
		t.Errorf("Cache should be empty but has remaining storage %d", rem)
	}
}

func TestLinearEvictSimple(t *testing.T) {
	capacity := 100
	lfu := NewLinearLfu(capacity, 1.0)
	checkCapacity(t, lfu, capacity)

	// sets 0 thru 9
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("____%d", i)
		val := []byte(key)
		ok := lfu.Set(key, val)
		if !ok {
			t.Errorf("Failed to add binding with key: %s", key)
			t.FailNow()
		}
	}

	// 0: 1
	// 1: 1
	// 2: 1
	// 3: 1
	// 4: 1
	// 5: 1
	// 6: 1
	// 7: 1
	// 8: 1
	// 9: 1

	// get 0 thru 8
	for i := 0; i < 9; i++ {
		key := fmt.Sprintf("____%d", i)
		res, found := lfu.Get(key)
		if !found {
			t.Errorf("Could not find %s as binding with key: %s", res, key)
			t.FailNow()
		}
	}

	// 0: 2
	// 1: 2
	// 2: 2
	// 3: 2
	// 4: 2
	// 5: 2
	// 6: 2
	// 7: 2
	// 8: 2
	// 9: 1

	// set 10
	key := fmt.Sprintf("___10")
	val := []byte("____a")
	ok := lfu.Set(key, val)
	if !ok {
		t.Errorf("Failed to add binding with key: %s", key)
		t.FailNow()
	}

	// 0: 2
	// 1: 2
	// 2: 2
	// 3: 2
	// 4: 2
	// 5: 2
	// 6: 2
	// 7: 2
	// 8: 2
	// 10: 1

	// gets 0 thru 10, 9 should not get cache hit
	for i := 0; i <= 10; i++ {
		key = fmt.Sprintf("____%d", i)
		if i == 10 {
			key = fmt.Sprintf("___10")
		}
		res, found := lfu.Get(key)
		// fmt.Printf("%s: %s. %t\n", key, res, found)
		if found && i == 9 {
			t.Errorf("Found %s as binding with key: %s", res, key)
			t.FailNow()
		} else if !found && i != 9 {
			t.Errorf("Could not find %s as binding with key: %s", res, key)
			t.FailNow()
		}
	}
}

func TestLinearEvict(t *testing.T) {
	capacity := 100
	lfu := NewLinearLfu(capacity, 1.0)
	checkCapacity(t, lfu, capacity)

	// sets 0 thru 9
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("____%d", i)
		val := []byte(key)
		ok := lfu.Set(key, val)
		if !ok {
			t.Errorf("Failed to add binding with key: %s", key)
			t.FailNow()
		}
	}

	// 0: 1
	// 1: 1
	// 2: 1
	// 3: 1
	// 4: 1
	// 5: 1
	// 6: 1
	// 7: 1
	// 8: 1
	// 9: 1

	// get 0 thru 9
	for i := 0; i < 10; i++ {
		for j := 0; j <= i; j++ {
			key := fmt.Sprintf("____%d", i)
			res, found := lfu.Get(key)
			if !found {
				t.Errorf("Could not find %s as binding with key: %s", res, key)
				t.FailNow()
			}
		}
	}

	// 0: 2
	// 1: 3
	// 2: 4
	// 3: 5
	// 4: 6
	// 5: 7
	// 6: 8
	// 7: 9
	// 8: 10
	// 9: 11

	// set 10
	key := fmt.Sprintf("___10")
	val := []byte("____a")
	ok := lfu.Set(key, val)
	if !ok {
		t.Errorf("Failed to add binding with key: %s", key)
		t.FailNow()
	}

	// 1: 3
	// 2: 4
	// 3: 5
	// 4: 6
	// 5: 7
	// 6: 8
	// 7: 9
	// 8: 10
	// 9: 11
	// 10: 1

	// gets 0 thru 10, 0 should not get cache hit
	for i := 0; i <= 10; i++ {
		key = fmt.Sprintf("____%d", i)
		if i == 10 {
			key = fmt.Sprintf("___10")
		}
		res, found := lfu.Get(key)
		// fmt.Printf("%s: %s. %t\n", key, res, found)
		if found && i == 0 {
			t.Errorf("Found %s as binding with key: %s", res, key)
			t.FailNow()
		} else if !found && i != 0 {
			t.Errorf("Could not find %s as binding with key: %s", res, key)
			t.FailNow()
		}
	}
}

func TestLinearAllMisses(t *testing.T) {
	capacity := 20
	numKeys := 3
	lfu := NewLinearLfu(capacity, 1.0)
	checkCapacity(t, lfu, capacity)

	// sets 1 thru 3
	for i := 1; i <= numKeys; i++ {
		key := fmt.Sprintf("____%d", i)
		val := []byte(key)
		ok := lfu.Set(key, val)
		if !ok {
			t.Errorf("Failed to add binding with key: %s", key)
			t.FailNow()
		}
	}

	// get 1 through 3, should all be misses
	for i := 1; i <= numKeys; i++ {
		key := fmt.Sprintf("____%d", i)
		val := []byte(key)
		_, found := lfu.Get(key)
		if !found {
			lfu.Set(key, val)
		}
	}

	cacheMisses := lfu.Stats().Misses
	if cacheMisses != numKeys {
		t.Errorf("Should have %d cache misses, only has %d", numKeys, cacheMisses)
		t.FailNow()
	}

	cacheHits := lfu.Stats().Hits
	if cacheHits != 0 {
		t.Errorf("Should have 0 cache hits, has %d", cacheHits)
		t.FailNow()

	}
}

func TestLinearAllHits(t *testing.T) {
	capacity := 30
	numKeys := 3
	lfu := NewLinearLfu(capacity, 1.0)
	checkCapacity(t, lfu, capacity)

	// sets 1 thru 3
	for i := 1; i <= numKeys; i++ {
		key := fmt.Sprintf("____%d", i)
		val := []byte(key)
		ok := lfu.Set(key, val)
		if !ok {
			t.Errorf("Failed to add binding with key: %s", key)
			t.FailNow()
		}
	}

	// get 1 through 3, should all be hits
	for i := 1; i <= numKeys; i++ {
		key := fmt.Sprintf("____%d", i)
		lfu.Get(key)
	}

	cacheMisses := lfu.Stats().Misses
	if cacheMisses != 0 {
		t.Errorf("Should have 0 cache misses, has %d", cacheMisses)
		t.FailNow()
	}

	cacheHits := lfu.Stats().Hits
	if cacheHits != numKeys {
		t.Errorf("Should have %d cache hits, only has %d", numKeys, cacheHits)
		t.FailNow()
	}
}
