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

func TestGet(t *testing.T) {
	capacity := 64
	lfu := NewLfu(capacity)
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

func TestRemove(t * testing.T) {
	capacity := 64
	lfu := NewLfu(capacity)
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
