package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type entry struct {
	key   Key
	value interface{}
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	el, ok := c.items[key]
	if ok {
		el.Value.(*entry).value = value
		c.queue.MoveToFront(el)
	} else {
		e := &entry{key: key, value: value}
		li := c.queue.PushFront(e)
		c.items[key] = li
	}
	if c.queue.Len() > c.capacity {
		el1 := c.queue.Back()
		c.queue.Remove(el1)
		delete(c.items, el1.Value.(*entry).key)
	}

	return ok
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	el, ok := c.items[key]
	if ok {
		c.queue.MoveToFront(el)
		return el.Value.(*entry).value, true
	}

	return nil, false
}

func (c *lruCache) Clear() {
	c.items = make(map[Key]*ListItem, c.capacity)
	c.queue = NewList()
}
