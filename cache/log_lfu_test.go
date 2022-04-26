/******************************************************************************
 * lfu_test.go
 * Author:
 * Usage:    `go test`  or  `go test -v`
 * Description:
 *    An incomplete unit testing suite for lfu.go. You are welcome to change
 *    anything in this file however you would like. You are strongly encouraged
 *    to create additional tests for your implementation, as the ones provided
 *    here are extremely basic, and intended only to demonstrate how to test
 *    your program.
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

func TestLogGet(t *testing.T) {
	capacity := 64
	lfu := NewLogLfu(capacity)
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
}

func TestLogRemove(t * testing.T) {
	capacity := 64
	lfu := NewLogLfu(capacity)
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

func TestLogEvictSimple(t *testing.T) {
	capacity := 100
	lfu := NewLogLfu(capacity)
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

func TestLogEvict(t *testing.T) {
	capacity := 100
	lfu := NewLogLfu(capacity)
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
