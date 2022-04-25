/******************************************************************************
 * lru_test.go
 * Author:
 * Usage:    `go test`  or  `go test -v`
 * Description:
 *    An incomplete unit testing suite for lru.go. You are welcome to change
 *    anything in this file however you would like. You are strongly encouraged
 *    to create additional tests for your implementation, as the ones provided
 *    here are extremely basic, and intended only to demonstrate how to test
 *    your program.
 ******************************************************************************/

package cache

import (
	"fmt"
	"testing"
	"time"
)

/******************************************************************************/
/*                                Constants                                   */
/******************************************************************************/
// Constants can go here

/******************************************************************************/
/*                                  Tests                                     */
/******************************************************************************/

func TestLFU(t *testing.T) {
	capacity := 64
	lru := NewLru(capacity)
	checkCapacity(t, lru, capacity)

	for i := 0; i < 4; i++ {
		key := fmt.Sprintf("key%d", i)
		val := []byte(key)
		ok := lru.Set(key, val)
		if !ok {
			t.Errorf("Failed to add binding with key: %s", key)
			t.FailNow()
		}

		res, _ := lru.Get(key)
		if !bytesEqual(res, val) {
			t.Errorf("Wrong value %s for binding with key: %s", res, key)
			t.FailNow()
		}
	}
}

func TestLFU53(t *testing.T) {
	capacity := 100
	lru := NewLru(capacity)
	checkCapacity(t, lru, capacity)

	// sets 0 thru 9
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("____%d", i)
		val := []byte(key)
		ok := lru.Set(key, val)
		if !ok {
			t.Errorf("Failed to add binding with key: %s", key)
			t.FailNow()
		}
	}

	// get 0
	key := fmt.Sprintf("____0")
	res, found := lru.Get(key)
	if !found {
		t.Errorf("Could not find %s as binding with key: %s", res, key)
		t.FailNow()
	}

	// set 10
	key = fmt.Sprintf("___10")
	val := []byte("____a")
	ok := lru.Set(key, val)
	if !ok {
		t.Errorf("Failed to add binding with key: %s", key)
		t.FailNow()
	}

	// gets 0 thru 10, 1 should not get cache hit
	for i := 0; i <= 10; i++ {
		key = fmt.Sprintf("____%d", i)
		if i == 10 {
			key = fmt.Sprintf("___10")
		}
		res, found := lru.Get(key)
		if found && i == 1 {
			t.Errorf("Found %s as binding with key: %s", res, key)
			t.FailNow()
		} else if !found && i != 1 {
			t.Errorf("Could not find %s as binding with key: %s", res, key)
			t.FailNow()
		}
	}
}

func TestLFU59(t *testing.T) {
	capacity := 20
	lru := NewLru(capacity)
	checkCapacity(t, lru, capacity)

	// set 1111 -> aaaa
	key := fmt.Sprintf("1111")
	val := []byte("aaaa")
	ok := lru.Set(key, val)
	if !ok {
		t.Errorf("Failed to add binding with key: %s", key)
		t.FailNow()
	}

	// set 2222 -> bbbb
	key = fmt.Sprintf("2222")
	val = []byte("bbbb")
	ok = lru.Set(key, val)
	if !ok {
		t.Errorf("Failed to add binding with key: %s", key)
		t.FailNow()
	}

	// get
	key = fmt.Sprintf("1111")
	res, found := lru.Get(key)
	if !found {
		t.Errorf("Could not find %s as binding with key: %s", res, key)
		t.FailNow()
	}

	// get
	key = fmt.Sprintf("2222")
	res, found = lru.Get(key)
	if !found {
		t.Errorf("Could not find %s as binding with key: %s", res, key)
		t.FailNow()
	}

	// set 3333 -> cccc
	key = fmt.Sprintf("3333")
	val = []byte("cccc")
	ok = lru.Set(key, val)
	if !ok {
		t.Errorf("Failed to add binding with key: %s", key)
		t.FailNow()
	}

	// get
	key = fmt.Sprintf("1111")
	res, found = lru.Get(key)
	if found {
		t.Errorf("Found %s as binding with key: %s", res, key)
		t.FailNow()
	}

	// remove
	key = fmt.Sprintf("2222")
	_, ok = lru.Remove(key)
	if !ok {
		t.Errorf("Failed to remove binding with key: %s", key)
		t.FailNow()
	}

	// get
	key = fmt.Sprintf("2222")
	res, found = lru.Get(key)
	if found {
		t.Errorf("Found %s as binding with key: %s", res, key)
		t.FailNow()
	}

	// get
	key = fmt.Sprintf("3333")
	res, found = lru.Get(key)
	if !found {
		t.Errorf("Could not find %s as binding with key: %s", res, key)
		t.FailNow()
	}

}
