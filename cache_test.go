package main

import (
	"testing"
	"time"
)

// test put and get
func TestPutAndGet(t *testing.T) {
	// given
	key := "a"
	val := 1
	ttl := time.Hour
	cache := LRUCache[string, int]{
		hashMap:          map[string]*Node[string, int]{},
		doublyLinkedList: &DoublyLinkedList[string, int]{},
		capacity:         2,
	}

	// when
	cache.Put(key, val, ttl)

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

	ttl := time.Hour

	cache := LRUCache[string, int]{
		hashMap:          map[string]*Node[string, int]{},
		doublyLinkedList: &DoublyLinkedList[string, int]{},
		capacity:         2,
	}

	// when
	cache.Put(keyA, valOne, ttl)
	cache.Put(keyB, valTwo, ttl)

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
	cache.Put(keyC, valThree, ttl)
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
	ttl := 50 * time.Millisecond

	cache := LRUCache[string, int]{
		hashMap:          map[string]*Node[string, int]{},
		doublyLinkedList: &DoublyLinkedList[string, int]{},
		capacity:         2,
	}

	// when
	cache.Put(key, val, ttl)
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
