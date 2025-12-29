package algorithms

// Graph represents an adjacency list graph
// Can be directed or undirected
type Graph struct {
	vertices int
	adjList  map[int][]int
}

// NewGraph creates a new graph with n vertices
func NewGraph(n int) *Graph {
	return &Graph{
		vertices: n,
		adjList:  make(map[int][]int),
	}
}

// AddEdge adds an edge to the graph
// For undirected graph, add edge in both directions
func (g *Graph) AddEdge(u, v int) {
	// TODO: Implement
}

// BFS performs breadth-first search starting from vertex start
// Returns the order of visited vertices
// Use case: Shortest path in unweighted graph, level-order traversal
func (g *Graph) BFS(start int) []int {
	// TODO: Implement
	// Hint: Use a queue, mark visited, process level by level
	return nil
}

// DFS performs depth-first search starting from vertex start
// Returns the order of visited vertices
// Use case: Detect cycles, topological sort, find connected components
func (g *Graph) DFS(start int) []int {
	// TODO: Implement
	// Hint: Can use recursion or stack
	return nil
}

// HasCycle detects if the graph has a cycle
// For directed graphs
func (g *Graph) HasCycle() bool {
	// TODO: Implement
	// Hint: Use DFS with three states (unvisited, visiting, visited)
	return false
}

// TopologicalSort returns vertices in topological order
// Only valid for Directed Acyclic Graphs (DAGs)
// Use case: Task scheduling, build systems, course prerequisites
func (g *Graph) TopologicalSort() ([]int, bool) {
	// TODO: Implement
	// Hint: Use DFS or Kahn's algorithm (BFS with in-degree)
	return nil, false
}

// ShortestPath finds shortest path from start to end (unweighted graph)
// Returns path and distance, or nil if no path exists
func (g *Graph) ShortestPath(start, end int) ([]int, int) {
	// TODO: Implement
	// Hint: BFS, track parent pointers to reconstruct path
	return nil, -1
}

// NumConnectedComponents counts connected components in undirected graph
func (g *Graph) NumConnectedComponents() int {
	// TODO: Implement
	// Hint: Run DFS/BFS from each unvisited vertex
	return 0
}

// IsBipartite checks if the graph can be colored with 2 colors
// Use case: Matching problems, conflict detection
func (g *Graph) IsBipartite() bool {
	// TODO: Implement
	// Hint: BFS/DFS with alternating colors, check for conflicts
	return false
}

// Common graph interview problems to also practice:
// - Clone Graph
// - Number of Islands (2D grid DFS/BFS)
// - Course Schedule (cycle detection)
// - Word Ladder (BFS)
// - Network Delay Time (Dijkstra's algorithm)
// - Minimum Spanning Tree (Kruskal's or Prim's)
