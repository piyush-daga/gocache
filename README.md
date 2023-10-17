# gocache
A multi-threaded concurrent-safe cache implementation in go

### Supported functions
- Eviction Policy: LRU by default
- TTL: Using callbacks, instead of garbage collecting
- Safe concurrent access across multiple goroutines
- Can specify the max size of the cache

### Implementations
- sync.Map
- using mutexes
- using mutex, setting an eviction policy
