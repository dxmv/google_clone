package main

import (
	"fmt"
	"strings"
)

type ListNode struct {
	Value int
	Prev  *ListNode
	Next  *ListNode
}

type LRUCache struct {
	Capacity int
	Head     *ListNode
	Tail     *ListNode
	Cache    map[int]*ListNode
}

func NewLRUCache(capacity int) LRUCache {
	lru := LRUCache{
		Capacity: capacity,
		Head:     nil,
		Tail:     nil,
		Cache:    make(map[int]*ListNode),
	}
	return lru
}

// to string function
func (lru *LRUCache) String() string {
	var sb strings.Builder
	// go through string and add each value to the result
	current := lru.Head
	sb.WriteString("[")
	for current != nil {
		var str string
		if current.Next == nil {
			str = fmt.Sprintf("%d", current.Value)
		} else {
			str = fmt.Sprintf("%d, ", current.Value)
		}
		sb.WriteString(str)
		current = current.Next
	}
	sb.WriteString("]")
	return sb.String()
}

func (lru *LRUCache) Put(value int) {
	// if the list is empty, add the value to the head
	if lru.Head == nil {
		lru.Head = &ListNode{Value: value}
		lru.Tail = lru.Head
		lru.Cache[value] = lru.Head
		return
	}
	// if the list is not empty, add the value to the head
	node, ok := lru.Cache[value]
	if ok {
		lru.moveNodeToHead(node)
		node.Value = value
		return
	}
	newNode := &ListNode{Value: value}
	newNode.Next = lru.Head
	lru.Head.Prev = newNode
	lru.Head = newNode
	lru.Cache[value] = newNode
	if len(lru.Cache) > lru.Capacity {
		lru.deleteNode()
	}
}

func (lru *LRUCache) Get(value int) int {
	// if the list is empty, return -1
	if lru.Head == nil {
		return -1
	}
	// if the value is in the list, return the value
	node, ok := lru.Cache[value]
	if !ok {
		return -1
	}
	// if the value is in the list, move it to the head
	lru.moveNodeToHead(node)
	return node.Value
}

func (lru *LRUCache) moveNodeToHead(node *ListNode) {
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

func (lru *LRUCache) deleteNode() {
	if lru.Tail == nil {
		return
	}
	tailNode := lru.Tail
	delete(lru.Cache, tailNode.Value)
	if tailNode.Prev != nil {
		tailNode.Prev.Next = nil
	}
	tailNode.Prev = nil
	tailNode.Next = nil
	lru.Tail = tailNode.Prev
	if lru.Tail == nil {
		lru.Head = nil
	}
}
