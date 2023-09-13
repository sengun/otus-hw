package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	mutex    *sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (cache *lruCache) Set(key Key, value interface{}) bool {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	if listItem, exists := cache.items[key]; exists {
		cacheItem := listItem.Value.(*cacheItem)
		cacheItem.value = value
		cache.queue.MoveToFront(listItem)

		return true
	}
	listItem := cache.queue.PushFront(&cacheItem{key: key, value: value})
	cache.items[key] = listItem
	if cache.queue.Len() > cache.capacity {
		cache.DisplaceLastItem()
	}

	return false
}

func (cache *lruCache) Get(key Key) (interface{}, bool) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	if listItem, exists := cache.items[key]; exists {
		cache.queue.MoveToFront(listItem)
		cacheItem := listItem.Value.(*cacheItem)
		return cacheItem.value, true
	}

	return nil, false
}

func (cache *lruCache) Clear() {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	cache.queue = NewList()
	cache.items = map[Key]*ListItem{}
}

func (cache *lruCache) DisplaceLastItem() {
	lastListItem := cache.queue.Back()
	if lastListItem == nil {
		return
	}
	cacheItem := lastListItem.Value.(*cacheItem)
	delete(cache.items, cacheItem.key)
	cache.queue.Remove(lastListItem)
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		mutex:    &sync.Mutex{},
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
