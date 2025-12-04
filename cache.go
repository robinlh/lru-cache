package main

import "time"

type Cache[K comparable, V any] interface {
	Get(key K) (V, bool)
	Put(key K, value V, ttl time.Duration)
	Size() int
}

type Node[S comparable, T any] struct {
	key       S
	data      T
	next      *Node[S, T]
	prev      *Node[S, T]
	expiresAt time.Time
}

type DoublyLinkedList[S comparable, T any] struct {
	head   *Node[S, T]
	tail   *Node[S, T]
	length int
}

type LRUCache[S comparable, T any] struct {
	capacity         int
	hashMap          map[S]*Node[S, T]
	doublyLinkedList *DoublyLinkedList[S, T]
}

func (cache *LRUCache[S, T]) Get(key S) (T, bool) {
	node, ok := cache.hashMap[key]
	if !ok {
		var zero T
		return zero, false
	}
	if node.expiresAt.Before(time.Now()) {
		cache.doublyLinkedList.RemoveNode(node)
		delete(cache.hashMap, key)
		var zero T
		return zero, false
	} else {
		cache.doublyLinkedList.MoveToFront(node)
		return node.data, true
	}
}

func (cache *LRUCache[S, T]) Put(key S, nodeData T, ttl time.Duration) {
	nodeIn, ok := cache.hashMap[key]
	if ok {
		nodeIn.data = nodeData
		nodeIn.expiresAt = time.Now().Add(ttl)
		cache.doublyLinkedList.MoveToFront(nodeIn)
	} else {
		newNode := &Node[S, T]{
			key:       key,
			data:      nodeData,
			expiresAt: time.Now().Add(ttl),
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

func (cache *LRUCache[S, T]) Size() int {
	return cache.doublyLinkedList.length
}

func (list *DoublyLinkedList[S, T]) AddToFront(n *Node[S, T]) {
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

func (list *DoublyLinkedList[S, T]) AddToEnd(n *Node[S, T]) {
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

func (list *DoublyLinkedList[S, T]) RemoveNode(n *Node[S, T]) {
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

func (list *DoublyLinkedList[S, T]) MoveToFront(n *Node[S, T]) {
	if n == nil || n == list.head {
		return
	}

	list.RemoveNode(n)
	list.AddToFront(n)
}

func (list *DoublyLinkedList[S, T]) PopTail() *Node[S, T] {
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
