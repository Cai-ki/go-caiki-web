package lru

import (
	"container/list"
)

type Cache struct {
	maxBytes int
	nowBytes int
	list     *list.List
	cache    map[string]*list.Element

	OnEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Size() int
}

func NewCache(maxBytes int, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		list:      list.New(),
		cache:     map[string]*list.Element{},
		OnEvicted: onEvicted,
	}
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	if element, ok := c.cache[key]; ok {
		c.list.MoveToFront(element)
		kv := element.Value.(*entry)
		return kv.value, true
	}
	return
}

func (kv *entry) Size() int {
	return /*int(unsafe.Sizeof(*kv)) + */ len(kv.key) + kv.value.Size()
}

func (c *Cache) RemoveOldest() {
	element := c.list.Back()
	if element != nil {
		c.list.Remove(element)
		kv := element.Value.(*entry)
		delete(c.cache, kv.key)
		c.nowBytes -= kv.Size()
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Add(key string, value Value) {
	if element, ok := c.cache[key]; ok {
		c.list.MoveToFront(element)
		kv := element.Value.(*entry)
		c.nowBytes += value.Size() - kv.value.Size()
		kv.value = value
	} else {
		kv := &entry{key, value}
		element := c.list.PushFront(kv)
		c.cache[key] = element
		c.nowBytes += kv.Size()
	}

	for c.maxBytes != 0 && c.nowBytes > c.maxBytes {
		c.RemoveOldest()
	}
}

func (c *Cache) Len() int {
	return c.list.Len()
}
