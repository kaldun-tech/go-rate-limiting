# Blockchain Protocol Engineering Study Plan

## Current Progress

### Data Structures
- [x] Trie (basic prefix tree)
- [x] BST
- [x] LRU Cache
- [x] Token Bucket Rate Limiter
- [ ] Heap (priority queue)
- [ ] Merkle Trees
- [ ] Patricia Tries (MPT)
- [ ] Bloom Filters
- [ ] DAGs

### Algorithms
- [ ] Serialization (RLP, SSZ)
- [ ] Graph Traversal
- [ ] Gossip Protocols
- [ ] Consensus (PBFT, Tendermint)

---

## Study Order

### Phase 1: Heap
**File:** `data-structures/heap.go`

Priority queues are used for:
- Transaction mempool ordering (by gas price)
- Validator selection (by stake weight)
- Timer management in consensus

**Key operations to implement:**
- [ ] Insert (sift up)
- [ ] ExtractMin/ExtractMax (sift down)
- [ ] Peek
- [ ] Heapify (build heap from array)

**Notes:**
```
```

---

### Phase 2: Serialization
**File:** `algorithms/serialization.go`

Start here because RLP/SSZ are needed to understand how blockchain data is encoded.

**RLP (Recursive Length Prefix):**
- [ ] Encode string (byte array)
- [ ] Encode list
- [ ] Encode uint64
- [ ] Decode item
- [ ] Streaming decoder

**SSZ (Simple Serialize):**
- [ ] Fixed-size encoding (uint8, uint16, uint32, uint64)
- [ ] Variable-size encoding
- [ ] Merkleization (hash tree root)
- [ ] Merkle proofs

**Test yourself:**
- Encode an Ethereum transaction manually
- Verify against go-ethereum's RLP output

**Notes:**
```
```

---

### Phase 3: Graph Algorithms
**File:** `algorithms/graph.go`

Foundation for DAGs, network topology, transaction ordering.

**Core operations:**
- [ ] AddEdge
- [ ] BFS (level-order, shortest path in unweighted graph)
- [ ] DFS (cycle detection, connected components)
- [ ] HasCycle (critical for DAG validation)
- [ ] TopologicalSort (transaction ordering)
- [ ] ShortestPath
- [ ] NumConnectedComponents

**Blockchain applications:**
- UTXO graph traversal
- Block DAG ordering (Narwhal, IOTA)
- Network partition detection

**Notes:**
```
```

---

### Phase 4: Gossip Protocols
**File:** `algorithms/gossip.go`

How blocks and transactions propagate through the network.

**Core concepts:**
- [ ] Message deduplication (seen set / bloom filter)
- [ ] Broadcast (eager push to all peers)
- [ ] RandomBroadcast (probabilistic, bandwidth-efficient)
- [ ] PushPull (state synchronization)
- [ ] PeerSampling (random peer selection)

**Simulation:**
- [ ] GossipSimulator with multiple nodes
- [ ] Measure propagation time
- [ ] Simulate message loss

**Blockchain applications:**
- Transaction propagation
- Block announcements
- Peer discovery

**Notes:**
```
```

---

### Phase 5: Consensus Algorithms
**File:** `algorithms/consensus.go`

The most complex topic. Builds on everything else.

**PBFT:**
- [ ] IsPrimary (leader selection)
- [ ] OnPrePrepare
- [ ] OnPrepare
- [ ] OnCommit
- [ ] HasQuorum (2f+1 calculation)
- [ ] StartViewChange

**Tendermint-style:**
- [ ] GetProposer (weighted by stake)
- [ ] OnProposal
- [ ] OnPrevote
- [ ] OnPrecommit
- [ ] HasTwoThirdsMajority

**PoS Primitives:**
- [ ] ValidatorSet management
- [ ] SelectProposer (stake-weighted)
- [ ] CalculateVotingPower

**Finality Gadget (Casper FFG style):**
- [ ] Checkpoint tracking
- [ ] Justification
- [ ] Finalization

**Notes:**
```
```

---

## Phase 6: Blockchain Data Structures (Future)

After algorithms, circle back to implement blockchain-specific data structures:

**Merkle Trees:**
- Binary Merkle tree
- Proof generation/verification
- Sparse Merkle trees

**Patricia Tries (MPT):**
- Ethereum's Modified Patricia Trie
- Hex-prefix encoding
- State root calculation

**Bloom Filters:**
- Insert/query
- Ethereum log bloom
- Light client filtering

**DAGs:**
- UTXO graph
- Conflict detection
- DAG-based consensus structures

---

## Resources

### Specifications
- [Ethereum Yellow Paper](https://ethereum.github.io/yellowpaper/paper.pdf)
- [RLP Spec](https://ethereum.org/en/developers/docs/data-structures-and-encoding/rlp/)
- [SSZ Spec](https://github.com/ethereum/consensus-specs/blob/dev/ssz/simple-serialize.md)

### Papers
- [PBFT](https://pmg.csail.mit.edu/papers/osdi99.pdf)
- [Tendermint](https://arxiv.org/abs/1807.04938)
- [Casper FFG](https://arxiv.org/abs/1710.09437)
- [Narwhal and Tusk](https://arxiv.org/abs/2105.11827)

### Code References
- [go-ethereum](https://github.com/ethereum/go-ethereum)
- [prysm](https://github.com/prysmaticlabs/prysm)
- [tendermint](https://github.com/tendermint/tendermint)

---

## Session Log

Use this to track study sessions and insights.

| Date | Topic | Time | Notes |
|------|-------|------|-------|
| | | | |
