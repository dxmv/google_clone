package main

import "fmt"

type ListNode[T any] struct {
	Key   string
	Value T
	Prev  *ListNode[T]
	Next  *ListNode[T]
}

// head will be the most recently used item
// tail will be the least recently used item
type LruCache[T any] struct {
	head  *ListNode[T]
	tail  *ListNode[T]
	size  int
	cache map[string]*ListNode[T]
}

// constructor
func NewLruCache[T any](size int) *LruCache[T] {
	return &LruCache[T]{
		head:  nil,
		tail:  nil,
		size:  size,
		cache: make(map[string]*ListNode[T]),
	}
}

// get item from cache
func (c *LruCache[T]) Get(key string) (T, bool) {
	return c.head.Value, true
}

// put item in cache
func (c *LruCache[T]) Put(key string, value T) {
	// get the node from the cache
	node, ok := c.cache[key]
	// if the node is not in the cache, add it to the cache
	if !ok {
		c.addNodeToHead(key, &ListNode[T]{
			Key:   key,
			Value: value,
		})
	} else {
		// if the node is in the cache, move it to the head of the list
		c.moveNodeToHead(key, node)
		node.Value = value
	}
}

// move a given node to the head of the list
func (c *LruCache[T]) moveNodeToHead(key string, node *ListNode[T]) {

	// if the node is already at the head, return
	if node == c.head {
		return
	}
	// if the node is at the tail, add it to the head it will automatically be removed from the tail
	if node == c.tail {
		c.addNodeToHead(key, node)
		return
	}
	// if the node is in the middle, remove it from the middle
	nxt := node.Next
	prev := node.Prev
	prev.Next = nxt
	nxt.Prev = prev
	// add it to the head
	c.addNodeToHead(key, node)
}

// add a given node to the head of the list
func (c *LruCache[T]) addNodeToHead(key string, node *ListNode[T]) {
	// add the node to the cache
	c.cache[key] = node
	// if the list is empty, set the head and tail to the node
	if c.head == nil {
		c.head = node
		c.tail = node
		return
	}
	// if the list is not empty, set the head to the node
	node.Next = c.head
	c.head.Prev = node
	c.head = node
	// if the list is full, remove the least recently used item
	if len(c.cache) >= c.size {
		c.removeNodeFromTail()
	}
}

// remove the least recently used item from the tail of the list
func (c *LruCache[T]) removeNodeFromTail() {
	// edge case if the list is empty
	if c.tail == nil {
		return
	}
	// the list is length 1
	if c.tail == c.head {
		c.tail = nil
		c.head = nil
		delete(c.cache, c.tail.Key)
		return
	}
	c.tail = c.tail.Prev
	c.tail.Next = nil
	delete(c.cache, c.tail.Key)
}

func (c *LruCache[T]) printList() {
	node := c.head
	for node != nil {
		fmt.Println(node.Key, node.Value)
		node = node.Next
	}
	fmt.Println("tail", c.tail.Value)
	fmt.Println("head", c.head.Value)
}
