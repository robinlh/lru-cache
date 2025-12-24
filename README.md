A thread-safe generic lru-cache with configurable eviction TTL.

Example usage. Create a cache with string key, int data, a capacity of 2 and an eviction TTL of 1 hour:
```
cache := NewLRUCache[string, int](2, time.Hour)
cache.Put("a", 1)
cache.Put("b", 2)
val, ok := cache.Get("a")
```
Moves key a to front (LRU)
