package main

import (
	"sync"
	"time"
)

type Cache[K comparable, V any] interface {
	Get(key K) (V, bool)
	Put(key K, value V)
	Size() int
}

type Node[K comparable, V any] struct {
	key       K
	data      V
	next      *Node[K, V]
	prev      *Node[K, V]
	expiresAt time.Time
}

type DoublyLinkedList[K comparable, V any] struct {
	head   *Node[K, V]
	tail   *Node[K, V]
	length int
}

type LRUCache[K comparable, V any] struct {
	capacity         int
	hashMap          map[K]*Node[K, V]
	doublyLinkedList *DoublyLinkedList[K, V]
	defaultTTL       time.Duration
	mu               sync.RWMutex
}

func (cache *LRUCache[K, V]) Get(key K) (V, bool) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	node, ok := cache.hashMap[key]
	if !ok {
		var zero V
		return zero, false
	}
	if node.expiresAt.Before(time.Now()) {
		cache.doublyLinkedList.RemoveNode(node)
		delete(cache.hashMap, key)
		var zero V
		return zero, false
	} else {
		cache.doublyLinkedList.MoveToFront(node)
		return node.data, true
	}
}

func (cache *LRUCache[K, V]) Put(key K, nodeData V) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	nodeIn, ok := cache.hashMap[key]
	if ok {
		nodeIn.data = nodeData
		nodeIn.expiresAt = time.Now().Add(cache.defaultTTL)
		cache.doublyLinkedList.MoveToFront(nodeIn)
	} else {
		newNode := &Node[K, V]{
			key:       key,
			data:      nodeData,
			expiresAt: time.Now().Add(cache.defaultTTL),
		}
		cache.hashMap[key] = newNode
		cache.doublyLinkedList.AddToFront(newNode)

		if cache.doublyLinkedList.length > cache.capacity {
			lruNode := cache.doublyLinkedList.PopTail()
			if lruNode != nil {
				delete(cache.hashMap, lruNode.key)
			}
		}
	}
}

func (cache *LRUCache[K, V]) Size() int {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	return cache.doublyLinkedList.length
}

func (list *DoublyLinkedList[K, V]) AddToFront(n *Node[K, V]) {
	if n == nil {
		return
	}

	n.prev = nil
	n.next = nil

	if list.head == nil {
		list.head = n
		list.tail = n
	} else {
		list.head.prev = n
		n.next = list.head
		list.head = n
	}
	list.length++
}

func (list *DoublyLinkedList[K, V]) AddToEnd(n *Node[K, V]) {
	if n == nil {
		return
	}

	n.prev = nil
	n.next = nil

	if list.head == nil {
		list.head = n
		list.tail = n
	} else {
		n.prev = list.tail
		list.tail.next = n
		list.tail = n
	}
	list.length++
}

func (list *DoublyLinkedList[K, V]) RemoveNode(n *Node[K, V]) {
	if n == nil {
		return
	}

	if n == list.head {
		list.head = n.next
	}

	if n == list.tail {
		list.tail = n.prev
	}

	if n.prev != nil {
		n.prev.next = n.next
	}

	if n.next != nil {
		n.next.prev = n.prev
	}

	n.prev = nil
	n.next = nil

	list.length--
}

func (list *DoublyLinkedList[K, V]) MoveToFront(n *Node[K, V]) {
	if n == nil || n == list.head {
		return
	}

	list.RemoveNode(n)
	list.AddToFront(n)
}

func (list *DoublyLinkedList[K, V]) PopTail() *Node[K, V] {
	if list.tail == nil {
		return nil
	}

	node := list.tail
	list.tail = list.tail.prev

	if list.tail != nil {
		list.tail.next = nil
	} else {
		list.head = nil
	}

	node.prev = nil
	node.next = nil

	list.length--
	return node
}

func NewLRUCache[K comparable, V any](capacity int, defaultTTL time.Duration) *LRUCache[K, V] {
	return &LRUCache[K, V]{
		capacity:         capacity,
		hashMap:          make(map[K]*Node[K, V]),
		doublyLinkedList: &DoublyLinkedList[K, V]{},
		defaultTTL:       defaultTTL,
	}
}
