# LRU Cache Interview Whiteboard Questions

Practice guide for demonstrating LRU Cache understanding in technical interviews.

## Core Implementation Questions

### 1. Design LRU Cache
**Question:** Design and implement an LRU (Least Recently Used) cache with O(1) get and put operations.

**Requirements:**
- `Get(key)` - Return value if exists, -1 otherwise
- `Put(key, value)` - Insert or update key-value pair
- When capacity is reached, evict the least recently used item
- Both operations must be O(1)

**Key Points to Cover:**
- HashMap for O(1) lookup
- Doubly linked list to track recency order
- Head = most recently used, Tail = least recently used
- Get/Put moves item to head
- Evict from tail when full
- Time: O(1) for both, Space: O(capacity)

**Follow-up:** Why do we need both a hash map AND a linked list?

---

### 2. Data Structure Choice
**Question:** Why use a doubly linked list instead of a singly linked list?

**Answer:**
- Need to remove nodes from middle of list (when moving to head)
- With singly linked list: O(n) to find previous node
- With doubly linked list: O(1) with prev pointer
- Trade-off: Extra pointer per node vs O(1) operations

**Follow-up:** Could you use an array instead of a linked list?

---

### 3. Get Operation
**Question:** Implement the Get operation.

**Key Points to Cover:**
```go
func (c *LRUCache) Get(key int) int {
    n, exists := c.cache[key]
    if !exists {
        return -1
    }
    c.moveToHead(n)  // Mark as recently used!
    return n.value
}
```

**Important:** Don't forget to update recency on Get!

**Follow-up:** What if we only wanted to update recency on Put, not Get?

---

### 4. Put Operation
**Question:** Implement the Put operation with eviction.

**Key Points to Cover:**
```go
func (c *LRUCache) Put(key int, value int) {
    if exists {
        // Update existing: change value + move to head
        node.value = value
        c.moveToHead(node)
    } else {
        // Add new node
        node = createNode(key, value)
        c.cache[key] = node
        c.addHead(node)

        // Evict if over capacity
        if len(c.cache) > capacity {
            tail := c.popTail()
            delete(c.cache, tail.key)  // Don't forget hash map!
        }
    }
}
```

**Common mistake:** Forgetting to delete from hash map when evicting

**Follow-up:** What if capacity is 0? Should we allow it?

---

### 5. Move to Head
**Question:** Implement moveToHead - the trickiest linked list operation.

**Key Points to Cover:**
```go
func (c *LRUCache) moveToHead(n *Node) {
    if n == c.head {
        return  // Already at head
    }

    // Unlink from current position
    if n.prev != nil {
        n.prev.next = n.next
    }
    if n.next != nil {
        n.next.prev = n.prev
    }
    if n == c.tail {
        c.tail = n.prev  // Update tail if moving tail
    }

    // Add to head
    n.prev = nil
    n.next = c.head
    if c.head != nil {
        c.head.prev = n
    }
    c.head = n
}
```

**Edge cases:**
- Node is already head
- Node is tail
- Only one node in list

**Follow-up:** Draw the pointer updates on the whiteboard.

---

### 6. Pop Tail (Eviction)
**Question:** Implement tail removal for eviction.

**Key Points to Cover:**
```go
func (c *LRUCache) popTail() *Node {
    if c.tail == nil {
        return nil
    }

    node := c.tail

    if c.tail.prev != nil {
        c.tail.prev.next = nil
        c.tail = c.tail.prev
    } else {
        // Last node - clear both head and tail
        c.head = nil
        c.tail = nil
    }

    return node
}
```

**Edge cases:**
- Empty cache
- Single node
- Need to update both head and tail when removing last node

---

## Conceptual Questions

### 7. Time Complexity Analysis
**Question:** Prove that LRU cache operations are O(1).

**Answer:**
- **Get:**
  - Hash map lookup: O(1)
  - Move to head (linked list): O(1) - have direct pointer
  - Total: O(1)

- **Put:**
  - Hash map insert/update: O(1)
  - Add to head: O(1)
  - Pop tail: O(1) - have direct tail pointer
  - Hash map delete: O(1)
  - Total: O(1)

**Space:** O(capacity) for hash map + linked list nodes

**Follow-up:** What if we used a balanced BST instead of hash map?

---

### 8. Real-World Use Cases
**Question:** Where are LRU caches used in practice?

**Answer:**
- **CPU cache eviction** - L1/L2/L3 caches
- **Database query caches** - MySQL query cache
- **CDN edge caching** - Content delivery networks
- **DNS caching** - Operating system DNS resolver
- **Page replacement** - Operating system virtual memory
- **Browser caches** - Recently accessed web pages
- **API response caching** - Redis, Memcached

**Follow-up:** When might LRU NOT be the best eviction policy?

---

### 9. LRU vs Other Eviction Policies
**Question:** What are alternatives to LRU?

**Answer:**
- **LFU (Least Frequently Used)** - Track access count, evict lowest
- **FIFO (First In First Out)** - Queue, evict oldest insertion
- **Random** - Evict random entry (surprisingly effective!)
- **TTL (Time To Live)** - Expire after fixed time
- **ARC (Adaptive Replacement Cache)** - Balances recency and frequency

**Trade-offs:**
- LRU: Good for temporal locality, can be fooled by scans
- LFU: Good for frequency-based, but slow to adapt
- FIFO: Simple but ignores usage patterns
- Random: No overhead, no worst case

**Follow-up:** How would you implement LFU?

---

### 10. Thread Safety
**Question:** How would you make this LRU cache thread-safe?

**Answer:**
```go
type LRUCache struct {
    mu       sync.Mutex
    capacity int
    cache    map[int]*Node
    head     *Node
    tail     *Node
}

func (c *LRUCache) Get(key int) int {
    c.mu.Lock()
    defer c.mu.Unlock()
    // ... implementation
}
```

**Options:**
1. Single mutex (simple, coarse-grained locking)
2. Read-write mutex (if reads >> writes)
3. Sharding (multiple caches, hash key to shard)
4. Lock-free algorithms (complex, high performance)

**Follow-up:** How would sharding work?

---

## Coding Challenges

### 11. LRU Cache with Expiration
**Question:** Add time-based expiration to LRU cache.

**Hint:** Store `expiresAt time.Time` in each node. Check on Get/Put.

**Key Points:**
- Check expiration before returning in Get
- Remove expired entries (may need background cleanup)
- Two eviction triggers: capacity OR expiration

---

### 12. LFU Cache
**Question:** Implement Least Frequently Used cache with O(1) operations.

**Hint:**
- Track access count for each key
- Maintain frequency buckets (doubly linked lists)
- Track minimum frequency

**Much harder than LRU!** This is a follow-up for senior roles.

---

### 13. Size-Based Eviction
**Question:** Items have different sizes. Cache has byte limit, not item limit.

**Hint:**
- Track `currentSize` and `maxSize`
- Store size in each node
- Evict until `currentSize + newItemSize <= maxSize`

**Key Points:**
- May need to evict multiple items for one large insertion
- Update size tracking on all operations

---

### 14. K-Way LRU
**Question:** Design cache that tracks K different access patterns separately.

**Example:** Separate LRU lists for reads vs writes.

**Hint:** Multiple linked lists, single hash map with type tag.

---

### 15. Distributed LRU
**Question:** How would you implement LRU across multiple servers?

**Approaches:**
1. **Consistent hashing** - Each server owns subset of keys
2. **Shared state** - Redis/Memcached cluster
3. **Local caches** - Each server has local LRU (may have duplicates)
4. **Two-tier** - Local L1 + shared L2 cache

**Trade-offs:** Consistency, network overhead, cache hit rate

---

## Common Pitfalls

### Mistake 1: Forgetting to Update Hash Map on Eviction
```go
// WRONG
func (c *LRUCache) Put(key, value int) {
    // ... add node ...
    if len(c.cache) > capacity {
        c.popTail()  // Only removes from list!
    }
}

// RIGHT
func (c *LRUCache) Put(key, value int) {
    // ... add node ...
    if len(c.cache) > capacity {
        tail := c.popTail()
        delete(c.cache, tail.key)  // Must delete from map!
    }
}
```

**Impact:** Memory leak - map grows unbounded!

---

### Mistake 2: Not Updating Recency on Get
```go
// WRONG
func (c *LRUCache) Get(key int) int {
    if node, exists := c.cache[key]; exists {
        return node.value  // Didn't move to head!
    }
    return -1
}

// RIGHT
func (c *LRUCache) Get(key int) int {
    if node, exists := c.cache[key]; exists {
        c.moveToHead(node)  // Mark as recently used
        return node.value
    }
    return -1
}
```

**Impact:** Get doesn't count as "use" - wrong eviction order!

---

### Mistake 3: Linked List Pointer Bugs
```go
// WRONG - lost reference to rest of list
func addHead(n *Node) {
    n.next = c.head
    c.head = n
    // Forgot: c.head.prev = n
}

// RIGHT
func addHead(n *Node) {
    n.next = c.head
    if c.head != nil {
        c.head.prev = n  // Bidirectional link!
    }
    c.head = n
    if c.tail == nil {
        c.tail = n  // First node case
    }
}
```

---

### Mistake 4: Off-by-One in Capacity Check
```go
// WRONG - allows capacity + 1 items
if len(c.cache) == capacity {
    evict()
}

// RIGHT - maintains exact capacity
if len(c.cache) > capacity {
    evict()
}
```

---

## Design Discussion Questions

### 16. Memory Overhead
**Question:** What's the memory overhead of your LRU implementation?

**Answer:** Per entry:
- Hash map: ~24 bytes (key, value, bucket pointer)
- Node: ~40 bytes (key, value, prev, next pointers)
- Total: ~64 bytes + actual key/value size

**Optimizations:**
- Use integer keys to reduce hash map overhead
- Pack node fields to reduce padding
- Use sync.Pool for node allocation

---

### 17. Cache Stampede
**Question:** Multiple threads Get() a missing key simultaneously. How to prevent all from computing the value?

**Answer:**
- **Singleflight pattern** - Only one goroutine computes, others wait
- Lock around check-and-set
- Pessimistic locking (lock key before check)

**Follow-up:** How does Redis handle this?

---

### 18. Eviction Callback
**Question:** How would you notify on eviction (e.g., write dirty data to disk)?

**Answer:**
```go
type EvictionCallback func(key, value interface{})

type LRUCache struct {
    // ...
    onEvict EvictionCallback
}

func (c *LRUCache) evict(node *Node) {
    if c.onEvict != nil {
        c.onEvict(node.key, node.value)
    }
    // ... remove from cache ...
}
```

**Use cases:** Write-back caching, logging, metrics

---

## Practice Drill

**Warm-up (5 minutes each):**
1. Draw an LRU cache with capacity 3 after: Put(1,1), Put(2,2), Get(1), Put(3,3)
2. Explain why we need O(1) operations for a cache
3. Walk through evicting tail when at capacity

**Medium (15 minutes each):**
4. Implement Get and Put on whiteboard
5. Implement moveToHead with all edge cases
6. Add thread safety with mutex

**Advanced (20+ minutes):**
7. Implement LFU cache
8. Add time-based expiration
9. Design distributed LRU cache system

---

## Key Takeaways for Interviews

1. **Two data structures** - Hash map for lookup + Doubly linked list for order
2. **Both operations O(1)** - Direct pointers enable constant time
3. **Update recency** - Both Get and Put move to head
4. **Evict from tail** - Least recently used is at tail
5. **Bidirectional links** - Don't forget prev pointers
6. **Hash map cleanup** - Always delete from map on eviction
7. **Edge cases** - Empty, single node, capacity 1
8. **Real-world context** - Used everywhere from CPU to CDN

**Interview Tips:**
- Start with data structure explanation
- Draw the state transitions
- Test with example: capacity 2, Put(1), Put(2), Get(1), Put(3)
- Discuss thread safety if time permits
- Know alternatives (LFU, FIFO, Random)

Good luck!
