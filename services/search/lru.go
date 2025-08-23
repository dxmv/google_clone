package main

import (
	"fmt"
	"strings"
)

type ListNode[K comparable, V any] struct {
	Key   K
	Value V
	Prev  *ListNode[K, V]
	Next  *ListNode[K, V]
}

type LRUCache[K comparable, V any] struct {
	Capacity int
	Head     *ListNode[K, V]
	Tail     *ListNode[K, V]
	Cache    map[K]*ListNode[K, V]
}

func NewLRUCache[K comparable, V any](capacity int) LRUCache[K, V] {
	lru := LRUCache[K, V]{
		Capacity: capacity,
		Head:     nil,
		Tail:     nil,
		Cache:    make(map[K]*ListNode[K, V]),
	}
	return lru
}

// to string function
func (lru *LRUCache[K, V]) String() string {
	var sb strings.Builder
	// go through string and add each value to the result
	current := lru.Head
	sb.WriteString("[")
	for current != nil {
		var str string
		if current.Next == nil {
			str = fmt.Sprintf("{%v: %v}", current.Key, current.Value)
		} else {
			str = fmt.Sprintf("{%v: %v}, ", current.Key, current.Value)
		}
		sb.WriteString(str)
		current = current.Next
	}
	sb.WriteString("]")
	return sb.String()
}

func (lru *LRUCache[K, V]) Put(key K, value V) {
	if lru.Capacity == 0 {
		return
	}
	// if the list is empty, add the value to the head
	if lru.Head == nil {
		lru.Head = &ListNode[K, V]{Key: key, Value: value}
		lru.Tail = lru.Head
		lru.Cache[key] = lru.Head
		return
	}
	// if the list is not empty, add the value to the head
	node, ok := lru.Cache[key]
	if ok {
		lru.moveNodeToHead(node)
		node.Value = value
		node.Key = key
		return
	}
	newNode := &ListNode[K, V]{Key: key, Value: value}
	newNode.Next = lru.Head
	lru.Head.Prev = newNode
	lru.Head = newNode
	lru.Cache[key] = newNode
	if len(lru.Cache) > lru.Capacity {
		lru.deleteNode()
	}
}

func (lru *LRUCache[K, V]) Get(key K) V {
	var zero V
	// if the list is empty, return -1
	if lru.Head == nil {
		return zero
	}
	// if the value is in the list, return the value
	node, ok := lru.Cache[key]
	if !ok {
		return zero
	}
	// if the value is in the list, move it to the head
	lru.moveNodeToHead(node)
	return node.Value
}

func (lru *LRUCache[K, V]) moveNodeToHead(node *ListNode[K, V]) {
	if node == lru.Head {
		return
	}
	if node == lru.Tail {
		lru.Tail = node.Prev
	}
	if node.Prev != nil {
		node.Prev.Next = node.Next
	}
	if node.Next != nil {
		node.Next.Prev = node.Prev
	}
	node.Next = lru.Head
	node.Prev = nil
	lru.Head.Prev = node
	lru.Head = node
}

func (lru *LRUCache[K, V]) deleteNode() {
	if lru.Tail == nil {
		return
	}
	tailNode := lru.Tail
	delete(lru.Cache, tailNode.Key)
	prev := tailNode.Prev
	if prev != nil {
		prev.Next = nil
	}
	lru.Tail = prev
	if lru.Tail == nil { // list became empty
		lru.Head = nil
	}

	tailNode.Prev, tailNode.Next = nil, nil
}
