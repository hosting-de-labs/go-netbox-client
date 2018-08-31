package cache

import (
	"sync"
)

type Cache struct {
	data sync.Map
}

func (c *Cache) Get(key interface{}) (interface{}, bool) {
	return c.data.Load(key)
}

func (c *Cache) Set(key interface{}, val interface{}) {
	c.data.Store(key, val)
}
