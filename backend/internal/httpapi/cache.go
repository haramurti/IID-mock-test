package httpapi

import (
	"sync"
	"time"
)

// ttlCache is a tiny in-memory cache with per-entry expiry.
// Good enough for one instance; production would use Redis (see README).
type ttlCache[V any] struct {
	mu    sync.Mutex
	ttl   time.Duration
	now   func() time.Time
	items map[string]cacheEntry[V]
}

type cacheEntry[V any] struct {
	val     V
	expires time.Time
}

func newTTLCache[V any](ttl time.Duration) *ttlCache[V] {
	return &ttlCache[V]{ttl: ttl, now: time.Now, items: map[string]cacheEntry[V]{}}
}

func (c *ttlCache[V]) Get(key string) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.items[key]
	if !ok || c.now().After(e.expires) {
		delete(c.items, key)
		var zero V
		return zero, false
	}
	return e.val, true
}

func (c *ttlCache[V]) Set(key string, val V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = cacheEntry[V]{val: val, expires: c.now().Add(c.ttl)}
}
