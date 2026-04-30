import(
	"sync"
)


type Cache struct{
	mu sync.Mutex
	cache map[string]cacheEntry
}

type cacheEntry struct{
	createdAt time.Time
	val []byte
}

func (c *Cache) Add(key string, val []byte) error{
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.cache == nil{
		c.cache = make(map[string]cacheEntry)
		c.cache[key] = cacheEntry{
    		createdAt: time.Now(),
    		val:       val,
		}
		return nil
	}
	c.cache[key] = cacheEntry{
    	createdAt: time.Now(),
    	val:       val,
	}
		
	return nil
}

func (c *Cache) Get(key string) ([]byte, bool){
	c.mu.Lock()
	defer c.mu.Unlock()

	cacheEntry, ok := c.cache[key]

	return cacheEntry.val, ok
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)

	for range ticker.C {
		c.mu.Lock()

    	for key, cacheEntry := range c.cache{
			if time.Since(cacheEntry.createdAt) > interval{
				delete(c.cache, key)
			}
		}
		c.mu.Unlock()		
	}

}

func newCache() (Cache){
	c := Cache{
    	cache: map[string]cacheEntry{},
	}
	//o c.cache := make(map[string]cacheEntry)
	const interval = 5 * time.Second

	go c.reapLoop(interval)

	return c
}