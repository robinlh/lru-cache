package main

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestPutAndGet(t *testing.T) {
	// given
	key := "a"
	val := 1
	cache := NewLRUCache[string, int](2, time.Hour)

	// when
	cache.Put(key, val)

	// then
	size := cache.doublyLinkedList.length
	if size != 1 {
		t.Errorf("size is %d ; want 1", size)
	}
	data, ok := cache.hashMap[key]
	if !ok {
		t.Errorf("data in cache is %d ; should be %d", data.data, val)
	}

	// when
	cacheEntry, ok := cache.Get(key)

	// then
	if !ok {
		t.Errorf("expected key to exist %s", key)
	}

	if cacheEntry != val {
		t.Errorf("expected cache entry data %d", val)
	}
}

func TestEvictionOrder(t *testing.T) {
	// given
	keyA := "a"
	valOne := 1

	keyB := "b"
	valTwo := 2

	keyC := "c"
	valThree := 3

	cache := NewLRUCache[string, int](2, time.Hour)

	// when
	cache.Put(keyA, valOne)
	cache.Put(keyB, valTwo)

	// then
	nodeMruB := cache.doublyLinkedList.head
	if nodeMruB.data != valTwo {
		t.Errorf("key %s expected first, got %d for key %s", keyB, nodeMruB.data, nodeMruB.key)
	}

	// when
	mruA, ok := cache.Get(keyA)
	nodeMruA := cache.doublyLinkedList.head

	// then
	if !ok && nodeMruA.data != valOne {
		t.Errorf("key %s expected first, got %d for key %s", keyA, mruA, keyA)
	}

	// when
	cache.Put(keyC, valThree)
	nodeMruC := cache.doublyLinkedList.head

	// then
	if nodeMruC.data != valThree {
		t.Errorf("key %s expected first, got %d for key %s", keyC, nodeMruC.data, nodeMruC.key)
	}

	// when
	evictedB, ok := cache.Get(keyB)

	// then
	if ok {
		t.Errorf("key %s should not be in cache, but was with key %s and data %d", keyB, keyB, evictedB)
	}
}

func TestTTLEviction(t *testing.T) {
	// given
	key := "a"
	val := 1

	cache := NewLRUCache[string, int](2, 50*time.Millisecond)

	// when
	cache.Put(key, val)
	time.Sleep(60 * time.Millisecond)

	// then
	evicted, ok := cache.Get(key)
	if ok {
		t.Errorf("key %s should not be in cache, but was with key %s and data %d", key, key, evicted)
	}

	size := cache.Size()
	if size != 0 {
		t.Errorf("cache should be empty after single value expired, but has size %d", size)
	}
}

func TestConcurrency(t *testing.T) {
	capacity := 10
	duration := 100 * time.Millisecond
	numWorkers := 8

	cache := NewLRUCache[string, int](capacity, 500*time.Millisecond)
	keys := []string{"a", "b", "c", "d", "e"}

	deadline := time.Now().Add(duration)

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	// ensure size never exceeds capacity
	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(1 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				size := cache.Size()
				if size > capacity {
					t.Errorf("cache size exceeded capacity: size=%d cap=%d", size, capacity)
					return
				}
			case <-done:
				return
			}
		}
	}()

	// make some workers
	for i := 0; i < numWorkers; i++ {
		go func(workerID int) {
			defer wg.Done()

			r := rand.New(rand.NewSource(time.Now().UnixNano() + int64(workerID)))

			for time.Now().Before(deadline) {
				key := keys[r.Intn(len(keys))]
				if r.Intn(2) == 0 {
					val := r.Intn(1000)
					cache.Put(key, val)
				} else {
					_, _ = cache.Get(key)
				}
			}
		}(i)
	}

	wg.Wait()
	close(done)

	if size := cache.Size(); size > capacity {
		t.Fatalf("final cache size exceeded capacity: size=%d cap=%d", size, capacity)
	}
}
