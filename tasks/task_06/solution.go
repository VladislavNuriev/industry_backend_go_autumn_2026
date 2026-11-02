package main

import "container/list"

type entry[K comparable, V any] struct {
	key   K
	value V
}

type LRUCache[K comparable, V any] struct {
	capacity int
	ll       list.List
	items    map[K]*list.Element
}

func NewLRUCache[K comparable, V any](capacity int) *LRUCache[K, V] {
	c := &LRUCache[K, V]{capacity: capacity}
	if capacity > 0 {
		c.items = make(map[K]*list.Element, capacity)
	}
	return c
}

func (c *LRUCache[K, V]) Get(key K) (value V, ok bool) {
	if c.capacity <= 0 {
		return value, false
	}
	elem, ok := c.items[key]
	if !ok {
		return value, false
	}
	c.ll.MoveToFront(elem)
	return elem.Value.(*entry[K, V]).value, true
}

func (c *LRUCache[K, V]) Set(key K, value V) {
	if c.capacity <= 0 {
		return
	}
	if elem, ok := c.items[key]; ok {
		// Существующий ключ: меняем только значение, давность не трогаем.
		elem.Value.(*entry[K, V]).value = value
		return
	}
	if c.ll.Len() >= c.capacity {
		// Вытесняем самый давний (back).
		back := c.ll.Back()
		if back != nil {
			old := back.Value.(*entry[K, V])
			delete(c.items, old.key)
			c.ll.Remove(back)
		}
	}
	e := &entry[K, V]{key: key, value: value}
	elem := c.ll.PushFront(e)
	c.items[key] = elem
}

type LRU[K comparable, V any] interface {
	Get(K) (V, bool)
	Set(K, V)
}
