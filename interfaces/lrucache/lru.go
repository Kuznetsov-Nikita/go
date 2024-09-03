//go:build !solution

package lrucache

import "container/list"

type data struct {
	key   int
	value int
}

type LruCache struct {
	storage *list.List
	cap     int
}

func (c LruCache) Get(key int) (int, bool) {
	for e := c.storage.Front(); e != nil; e = e.Next() {
		if e.Value.(data).key == key {
			c.storage.MoveToBack(e)
			return e.Value.(data).value, true
		}
	}

	return 0, false
}

func (c LruCache) Set(key, value int) {
	if c.cap == 0 {
		return
	}

	for e := c.storage.Front(); e != nil; e = e.Next() {
		if e.Value.(data).key == key {
			e.Value = data{key: key, value: value}
			c.storage.MoveToBack(e)
			return
		}
	}
	if c.cap > c.storage.Len() {
		c.storage.PushBack(data{key: key, value: value})
	} else {
		c.storage.Front().Value = data{key: key, value: value}
		c.storage.MoveToBack(c.storage.Front())
	}
}

func (c LruCache) Range(f func(key, value int) bool) {
	for e := c.storage.Front(); e != nil; e = e.Next() {
		if !f(e.Value.(data).key, e.Value.(data).value) {
			return
		}
	}
}

func (c LruCache) Clear() {
	c.storage = list.New()
}

func New(cap int) Cache {
	return LruCache{storage: list.New(), cap: cap}
}
