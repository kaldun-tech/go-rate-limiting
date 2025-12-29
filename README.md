# Go Algorithm Practice

A collection of algorithm and data structure implementations in Go for interview preparation and skill development.

## Structure

```
go-algorithm-practice/
├── rate-limiting/          # Rate limiting algorithms
│   └── token-bucket/      # Token bucket implementation (complete)
├── data-structures/        # Common data structures
│   ├── lru-cache.go       # LRU Cache (TODO)
│   ├── bst.go             # Binary Search Tree (TODO)
│   ├── trie.go            # Prefix Tree (TODO)
│   ├── heap.go            # Min/Max Heap (TODO)
│   └── linkedlist.go      # Linked List problems (TODO)
├── algorithms/             # Algorithm patterns
│   ├── graph.go           # Graph traversal (BFS, DFS, etc.) (TODO)
│   ├── sorting.go         # Sorting algorithms (TODO)
│   ├── dynamic-programming.go  # DP problems (TODO)
│   ├── backtracking.go    # Backtracking patterns (TODO)
│   └── two-pointers.go    # Two pointer techniques (TODO)
└── examples/              # Usage examples

```

## What's Implemented

### Rate Limiting - Token Bucket ✅

A complete, production-ready implementation of the Token Bucket rate limiting algorithm used by Stripe, GitHub, and AWS.

**Features:**
- Thread-safe with `sync.Mutex`
- Key-based rate limiting (per-user, per-API-key, per-IP)
- Burst support
- Detailed rate limit information for HTTP headers
- Weighted costs (different operations cost different tokens)

**Quick example:**
```go
import tokenbucket "github.com/kaldun-tech/go-algorithm-practice/rate-limiting/token-bucket"

// Create limiter: 10 requests/second
limiter := tokenbucket.NewTokenBucket(10, time.Second, 0)

// Check if request allowed
if limiter.Allow("user:alice") {
    fmt.Println("Request allowed!")
}
```

See [rate-limiting/CLAUDE.md](rate-limiting/CLAUDE.md) for detailed implementation notes and interview tips.

**Run the examples:**
```bash
go run examples/main.go
```

## What's Coming (TODOs)

All other files are boilerplate with TODO implementations. Each file includes:
- Clear problem descriptions
- Common interview questions
- Implementation hints
- Time/space complexity notes
- Related problems to practice

### Data Structures

**LRU Cache** - HashMap + Doubly Linked List for O(1) operations
- Get/Put in O(1) time
- Classic interview question

**Binary Search Tree** - Insert, search, delete, validate
- Three deletion cases
- BST validation with bounds

**Trie (Prefix Tree)** - String operations
- Autocomplete
- Word search
- Prefix matching

**Min/Max Heap** - Priority queue
- K largest elements
- Merge K sorted lists

**Linked List** - Classic problems
- Reverse list
- Detect cycle (Floyd's algorithm)
- Find middle
- Merge sorted lists

### Algorithm Patterns

**Graph Algorithms**
- BFS/DFS traversal
- Cycle detection
- Topological sort
- Shortest path
- Connected components

**Sorting**
- QuickSort, MergeSort, HeapSort
- Binary Search
- QuickSelect (Kth largest)
- Merge K sorted arrays

**Dynamic Programming**
- Fibonacci, Coin Change
- Longest Increasing Subsequence
- Knapsack 0/1
- Edit Distance
- Longest Common Subsequence

**Backtracking**
- Permutations, Combinations, Subsets
- N-Queens
- Sudoku Solver
- Word Search
- Generate Parentheses

**Two Pointers / Sliding Window**
- Two Sum, Three Sum
- Container With Most Water
- Longest Substring Without Repeating
- Minimum Window Substring
- Trapping Rain Water

## Getting Started

### Installation

```bash
git clone https://github.com/kaldun-tech/go-algorithm-practice.git
cd go-algorithm-practice
```

### Running Tests

```bash
# Test everything
go test ./...

# Test specific package
go test ./rate-limiting/token-bucket

# Run with coverage
go test -cover ./...
```

### Implementation Workflow

1. **Pick a problem** from the boilerplate files
2. **Read the comments** - they explain the problem and give hints
3. **Implement the solution** - replace `// TODO: Implement`
4. **Write tests** - create `*_test.go` files
5. **Benchmark** - add benchmarks for performance-critical code

## Why Go?

This repository uses Go for algorithm practice because:
- **Simplicity**: Clean syntax, no inheritance complexity
- **Performance**: Compiled, fast, good for algorithm analysis
- **Standard library**: Excellent testing/benchmarking tools
- **Industry relevance**: Used at Google, Uber, Stripe, Cloudflare
- **Concurrency**: Built-in goroutines for concurrent algorithms

## Interview Focus

Each implementation emphasizes:
- **Clarity over cleverness** - readable, maintainable code
- **Edge cases** - handling invalid inputs, empty cases, boundaries
- **Time/space complexity** - understanding trade-offs
- **Testing** - comprehensive test cases
- **Explanation** - comments that show understanding

## Project Philosophy

**This is NOT a copy-paste solutions repository.**

All TODO implementations should be written manually to:
- Build muscle memory
- Understand edge cases deeply
- Practice explaining your thought process
- Make mistakes and learn from them

The boilerplate provides structure and hints, but YOU implement the algorithms.

## Resources

- **LeetCode** - Practice problems with test cases
- **NeetCode** - Curated problem lists and explanations
- **Go by Example** - Learn Go idioms
- **Go standard library** - `sort`, `container/heap`, `container/list` for reference

## Contributing

This is a personal practice repository, but feel free to:
- Fork for your own practice
- Suggest improvements via issues
- Share alternative approaches

## License

MIT

---

**Status:** Active development. Token Bucket is complete and production-ready. Other algorithms are structured TODOs for manual implementation.
