package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

// element of cache.
type CacheItem struct {
	key   Key
	value interface{}
}

// Create a new cache item.
func NewCacheItem(key Key, value interface{}) *CacheItem {
	return &CacheItem{
		key:   key,
		value: value,
	}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	if item, ok := c.items[key]; ok {
		itemValue := item.Value.(*CacheItem)
		itemValue.value = value
		c.queue.MoveToFront(item)
		return true
	}

	newItem := NewCacheItem(key, value)
	addedItem := c.queue.PushFront(newItem)
	c.items[key] = addedItem

	if c.queue.Len() > c.capacity {
		removedItem := c.queue.Back()
		if removedItem != nil {
			c.queue.Remove(removedItem)
			delete(c.items, removedItem.Value.(*CacheItem).key)
		}
	}

	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	if item, ok := c.items[key]; ok {
		c.queue.MoveToFront(item)
		return item.Value.(*CacheItem).value, true
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}
