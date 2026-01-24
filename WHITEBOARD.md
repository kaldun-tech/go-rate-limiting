# Blockchain Protocol Engineering - Whiteboard Exercises

Quick-reference for interview prep. Each topic: whiteboard prompt, blockchain significance, complexity, key insight.

---

## Data Structures

### Heap / Priority Queue

**Whiteboard:** Draw a min-heap after inserting [5, 3, 8, 1, 2]. Then extract-min twice. Show the sift-up and sift-down operations.

**Blockchain Significance:**
- Mempool ordering: prioritize transactions by gas price / fee
- Validator selection: weighted random by stake
- Consensus timeouts: manage pending timers efficiently

**Complexity:**
| Operation | Time | Space |
|-----------|------|-------|
| Insert | O(log n) | O(1) |
| Extract | O(log n) | O(1) |
| Peek | O(1) | O(1) |
| Heapify | O(n) | O(1) |

**Key Insight:** Heaps give you the "best" item in O(log n) without sorting everything. Mempools don't need full ordering—just the highest-fee tx next.

---

### Trie (Prefix Tree)

**Whiteboard:** Insert "eth", "ethereum", "etc" into a trie. Show the structure. How would you find all words starting with "et"?

**Blockchain Significance:**
- Foundation for Patricia Tries (Ethereum state)
- Efficient prefix lookups for address/key matching
- Deterministic ordering for state roots

**Complexity:**
| Operation | Time | Space |
|-----------|------|-------|
| Insert | O(k) | O(k) |
| Search | O(k) | O(1) |
| Prefix query | O(k + results) | O(1) |

*k = key length*

**Key Insight:** Tries give O(key length) lookups regardless of how many keys exist. Hash tables are O(1) amortized but can't do prefix queries.

---

### Merkle Tree

**Whiteboard:** Given 4 transactions [A, B, C, D], draw the Merkle tree. Generate a proof that B is included. How does a light client verify?

**Blockchain Significance:**
- Block transaction roots: verify tx inclusion without full block
- State roots: commit to entire state in 32 bytes
- Light clients: SPV proofs with O(log n) data

**Complexity:**
| Operation | Time | Space |
|-----------|------|-------|
| Build tree | O(n) | O(n) |
| Generate proof | O(log n) | O(log n) |
| Verify proof | O(log n) | O(1) |

**Key Insight:** Verification is O(log n) with O(log n) proof size. A light client can verify a tx in a block of 10,000 txs with only ~14 hashes.

---

### Patricia Trie (MPT)

**Whiteboard:** Explain the difference between a regular trie and a Patricia trie. Why does Ethereum use path compression?

**Blockchain Significance:**
- Ethereum state trie: accounts, storage, receipts
- Path compression reduces depth for sparse keyspaces
- Enables state proofs for light clients and bridges

**Complexity:**
| Operation | Time | Space |
|-----------|------|-------|
| Get/Put | O(k) | O(k) |
| Proof generation | O(k) | O(k) |
| State root | O(1)* | O(1) |

*After updates propagate up*

**Key Insight:** 256-bit addresses would make a regular trie 256 levels deep. Patricia compression makes it proportional to the number of accounts, not key length.

---

### Bloom Filter

**Whiteboard:** Insert "tx1", "tx2" into a Bloom filter with 3 hash functions. Check for "tx3" (not inserted). Can you get false negatives?

**Blockchain Significance:**
- Ethereum log blooms: filter blocks by topic without scanning
- Light client sync: skip blocks that can't contain relevant events
- Mempool deduplication: fast "probably seen" check

**Complexity:**
| Operation | Time | Space |
|-----------|------|-------|
| Insert | O(k) | O(1) |
| Query | O(k) | O(1) |
| Space | - | O(m) bits |

*k = hash functions, m = filter size*

**Key Insight:** No false negatives, only false positives. If Bloom says "no," it's definitely no. If "yes," check the actual data. Tunable tradeoff between size and false positive rate.

---

### DAG (Directed Acyclic Graph)

**Whiteboard:** Draw a UTXO graph for: A sends to B, B sends to C and D. How do you detect double-spend attempts?

**Blockchain Significance:**
- UTXO model: Bitcoin, Cardano transaction graphs
- DAG-based consensus: IOTA Tangle, Hedera Hashgraph
- Block DAGs: parallel block production (Narwhal)

**Complexity:**
| Operation | Time | Space |
|-----------|------|-------|
| Topological sort | O(V + E) | O(V) |
| Cycle detection | O(V + E) | O(V) |
| Reachability | O(V + E) | O(V) |

**Key Insight:** Topological ordering gives a valid execution order respecting dependencies. DAG-based systems can process non-conflicting transactions in parallel.

---

### LRU Cache

**Whiteboard:** Design an O(1) get/put cache. What data structures do you need? Walk through: put(1), put(2), get(1), put(3) with capacity 2.

**Blockchain Significance:**
- State caching: hot accounts/storage slots
- Block/tx caching: recently verified data
- Peer connection caching: frequently contacted nodes

**Complexity:**
| Operation | Time | Space |
|-----------|------|-------|
| Get | O(1) | O(1) |
| Put | O(1) | O(1) |
| Space | - | O(capacity) |

**Key Insight:** Hash map + doubly linked list. Map gives O(1) lookup, list gives O(1) reordering. Moving to head on access maintains LRU order.

---

## Algorithms

### RLP Encoding

**Whiteboard:** Encode the list ["hello", [1, 2]] in RLP. Show the byte-level breakdown with prefixes.

**Blockchain Significance:**
- Ethereum wire protocol: all messages RLP-encoded
- Transaction serialization: signing and hashing
- MPT nodes: stored as RLP-encoded data

**Complexity:**
| Operation | Time | Space |
|-----------|------|-------|
| Encode | O(n) | O(n) |
| Decode | O(n) | O(n) |

**Key Insight:** RLP is self-describing (no schema needed) but not self-delimiting in streams. Length prefixes enable O(1) skip-ahead for nested structures.

---

### SSZ (Simple Serialize)

**Whiteboard:** How does SSZ differ from RLP? Why does Ethereum 2.0 use SSZ instead of RLP?

**Blockchain Significance:**
- Beacon chain: all consensus data SSZ-encoded
- Fixed offsets enable O(1) field access
- Native Merkleization for efficient proofs

**Complexity:**
| Operation | Time | Space |
|-----------|------|-------|
| Encode | O(n) | O(n) |
| Decode | O(n) | O(n) |
| Hash tree root | O(n) | O(log n) |
| Field access | O(1)* | O(1) |

*For fixed-size types*

**Key Insight:** SSZ trades flexibility for performance. Fixed-size offsets mean you can jump to any field without parsing everything before it. Critical for validators processing many attestations.

---

### Graph Traversal (BFS/DFS)

**Whiteboard:** Given a network of 5 peers with connections, find the shortest path (fewest hops) from peer A to peer E. Which algorithm and why?

**Blockchain Significance:**
- Block propagation: how many hops to reach all nodes?
- Network topology analysis: find bridges, partitions
- Transaction dependency graphs: execution ordering

**Complexity:**
| Operation | Time | Space |
|-----------|------|-------|
| BFS | O(V + E) | O(V) |
| DFS | O(V + E) | O(V) |
| Shortest path (unweighted) | O(V + E) | O(V) |

**Key Insight:** BFS gives shortest path in unweighted graphs (fewest hops). For block propagation analysis, BFS from a node shows how quickly messages spread.

---

### Topological Sort

**Whiteboard:** Given transactions with dependencies (tx2 depends on tx1's output), show the valid execution order. What if there's a cycle?

**Blockchain Significance:**
- Transaction ordering: respect UTXO dependencies
- Smart contract deployment: deploy dependencies first
- DAG consensus: linearize parallel blocks

**Complexity:**
| Operation | Time | Space |
|-----------|------|-------|
| Kahn's algorithm | O(V + E) | O(V) |
| DFS-based | O(V + E) | O(V) |

**Key Insight:** No valid topological order exists if there's a cycle. Cycle detection is essentially the same algorithm—if you can't complete the sort, there's a cycle (double-spend attempt in UTXO).

---

### Gossip Protocols

**Whiteboard:** A node receives a new block. It has 8 peers. Compare: (a) broadcast to all, (b) broadcast to random 3. Tradeoffs?

**Blockchain Significance:**
- Block/tx propagation: how data spreads through network
- Bandwidth efficiency: don't flood the network
- Eclipse attack resistance: diverse peer connections

**Complexity:**
| Metric | Eager (all peers) | Probabilistic |
|--------|-------------------|---------------|
| Latency | Lower | Higher |
| Bandwidth | O(degree × n) | O(fanout × n) |
| Reliability | Higher | Depends on fanout |

**Key Insight:** Probabilistic gossip trades latency for bandwidth. With fanout=3 and well-connected graph, messages still reach everyone in O(log n) rounds with high probability.

---

### PBFT Consensus

**Whiteboard:** 4 validators (A, B, C, D). A is Byzantine. Walk through a round: who proposes, what messages are exchanged, when is the block final?

**Blockchain Significance:**
- Immediate finality: no probabilistic confirmation
- Permissioned chains: Hyperledger, enterprise blockchains
- Foundation for modern BFT variants (HotStuff, Tendermint)

**Complexity:**
| Metric | Value |
|--------|-------|
| Message complexity | O(n²) per round |
| Fault tolerance | f < n/3 |
| Rounds to finality | 3 (pre-prepare, prepare, commit) |

**Key Insight:** Need 3f+1 nodes to tolerate f Byzantine faults. Why? Need 2f+1 for quorum, but f might be silent (not voting), so 2f+1 + f = 3f+1 total.

---

### Tendermint Consensus

**Whiteboard:** Compare Tendermint to PBFT. What's the propose → prevote → precommit flow? What happens if the proposer is offline?

**Blockchain Significance:**
- Cosmos ecosystem: IBC-connected chains
- Optimized for blockchain: integrated proposer selection
- Locking mechanism prevents conflicting commits

**Complexity:**
| Metric | Value |
|--------|-------|
| Message complexity | O(n²) per round |
| Fault tolerance | f < n/3 |
| Rounds to finality | 2-3 (prevote, precommit, +1 if timeout) |

**Key Insight:** Tendermint adds "locking"—once you precommit a block, you can't prevote for a different block until you see proof it's safe. Prevents safety violations during async periods.

---

### Finality Gadgets (Casper FFG)

**Whiteboard:** Explain justification vs finalization. Given checkpoints at epochs 1, 2, 3—when does epoch 1 become finalized?

**Blockchain Significance:**
- Ethereum PoS: finality on top of fork choice
- Economic finality: reverting costs slashed stake
- Bridges: safe to trust finalized state

**Complexity:**
| Metric | Value |
|--------|-------|
| Time to justification | 1 epoch (~6.4 min on Ethereum) |
| Time to finalization | 2 epochs (~12.8 min) |
| Fault tolerance | f < n/3 of stake |

**Key Insight:** Two-phase commit at epoch granularity. Justified = 2/3 voted for this checkpoint. Finalized = justified AND the next checkpoint is also justified. Provides "economic finality"—reverting requires mass slashing.

---

## Quick Complexity Reference

| Structure/Algorithm | Insert/Build | Query/Access | Space | Blockchain Use |
|---------------------|--------------|--------------|-------|----------------|
| Heap | O(log n) | O(1) peek | O(n) | Mempool |
| Trie | O(k) | O(k) | O(n×k) | State prefix |
| Merkle Tree | O(n) | O(log n) proof | O(n) | Tx inclusion |
| Patricia Trie | O(k) | O(k) | O(n) | Ethereum state |
| Bloom Filter | O(k) | O(k) | O(m) bits | Log filtering |
| LRU Cache | O(1) | O(1) | O(cap) | Hot state |
| RLP | O(n) | O(n) | O(n) | Wire protocol |
| BFS/DFS | O(V+E) | - | O(V) | Network analysis |
| PBFT/Tendermint | O(n²) msgs | - | O(n) | Block finality |

---

## Interview Tips

1. **Start with "why"** - Before diving into implementation, explain why this matters for blockchain

2. **Draw first** - Whiteboard the structure/flow before writing pseudocode

3. **State complexity early** - Mention time/space complexity as you explain, not just at the end

4. **Connect to real systems** - "Ethereum uses this for..." shows you understand practical applications

5. **Know the tradeoffs** - Every choice has costs. Bloom filters have false positives. PBFT has O(n²) messages. Show you understand the engineering decisions.
