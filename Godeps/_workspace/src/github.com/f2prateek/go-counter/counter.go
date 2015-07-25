package counter

import (
	"sync"
	"sync/atomic"
)

type Counter struct {
	values map[string]*int64
	sync.Mutex
}

// Create a new Counter.
func New() *Counter {
	m := make(map[string]*int64)
	return &Counter{values: m}
}

// Increment the value for the given key.
func (c *Counter) Increment(key string) {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.values[key]; !ok {
		var n int64 = 0
		c.values[key] = &n
	}

	atomic.AddInt64(c.values[key], 1)
}

// Return a copy of the counter values.
func (c *Counter) Values() map[string]int64 {
	c.Lock()
	defer c.Unlock()

	v := make(map[string]int64)
	for key, value := range c.values {
		v[key] = *value
	}
	return v
}
