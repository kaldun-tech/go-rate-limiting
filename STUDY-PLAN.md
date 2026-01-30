# Blockchain Protocol Engineering Study Plan

## Current Progress

### Data Structures
- [x] Trie (basic prefix tree)
- [x] BST
- [x] LRU Cache
- [x] Token Bucket Rate Limiter
- [x] Heap (priority queue)
- [x] Merkle Trees
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
- [x] Insert (sift up)
- [x] ExtractMin/ExtractMax (sift down)
- [x] Peek
- [x] Heapify (build heap from array)

**Benchmarking goals:**
- [ ] Heapify 1M elements in <100ms
- [ ] Push/Pop operations <1μs each
- [ ] Memory overhead <2x the raw array size

**Interview questions:**
1. "Why is Heapify O(n) instead of O(n log n)? Walk me through the math."
2. "How would you implement a heap that supports efficient decrease-key for Dijkstra's algorithm?"
3. "Design a data structure that returns the median in O(1) time." (Hint: two heaps)

**Notes:**
```
Peek operation is identical across Min and Max heaps.
The internal siftDown operation is identical and can be done on the base class using a comparator
Pop (ExtractMin/Max) and Heapify then become identical using the siftDown with appropriate comparator (< or >)
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

**Benchmarking goals:**
- [ ] RLP encoder should match go-ethereum's performance within 2x
- [ ] Encode/decode 10,000 transactions in <50ms
- [ ] SSZ Merkleization of 1M leaves in <500ms
- [ ] Zero-allocation encode path for fixed-size types

**Test yourself:**
- Encode an Ethereum transaction manually
- Verify against go-ethereum's RLP output

**Protobuf comparison (from Hedera experience):**
```
At Hedera, I optimized protobuf serialization 5.6x. Key techniques:
- Pre-allocated buffers to avoid GC pressure
- Unsafe string/byte conversions where safe
- Batched encoding for repeated fields

Apply same thinking to RLP/SSZ:
- RLP: Length-prefix overhead vs protobuf varint - when is each better?
- SSZ: Fixed offsets enable parallel decoding - protobuf can't do this
- SSZ: Merkleization enables partial proofs - protobuf has no equivalent

Interview story: "I improved serialization throughput from X to Y by [technique].
The same principle applies to RLP/SSZ because..."
```

**Interview questions:**
1. "Why does Ethereum use RLP instead of protobuf or JSON? What are the tradeoffs?"
2. "How does SSZ enable light clients that protobuf cannot?"
3. "A block has 50,000 transactions. How would you optimize RLP encoding for this?"
4. "Explain how SSZ Merkleization works and why it's useful for state proofs."

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

**Benchmarking goals:**
- [ ] Topological sort of 100K nodes in <50ms
- [ ] Cycle detection in <10ms for 10K edges
- [ ] BFS/DFS should be O(V+E) with minimal constant factor

**Blockchain applications:**
- UTXO graph traversal
- Block DAG ordering (Narwhal, IOTA)
- Network partition detection

**Interview questions:**
1. "How would you detect a double-spend in a UTXO graph?"
2. "Explain how topological sort is used in DAG-based consensus (Narwhal/Tusk)."
3. "A transaction references 5 inputs from different blocks. How do you order execution?"
4. "How would you detect network partitions in a P2P gossip network?"

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

**Benchmarking goals:**
- [ ] 1000-node simulation completes in <1s
- [ ] Message reaches 99% of nodes within O(log n) rounds
- [ ] Bloom filter false positive rate <1% with 100K seen messages

**Blockchain applications:**
- Transaction propagation
- Block announcements
- Peer discovery

**Failure scenarios:**
- What happens if 30% of messages are dropped? (partition tolerance)
- What if a malicious node floods the network? (DoS resistance)
- What if clocks are skewed by 5 seconds? (timing assumptions)

**Interview questions:**
1. "How do you prevent a malicious node from flooding the gossip network?"
2. "Explain the tradeoff between eager push and lazy pull in gossip protocols."
3. "How would you implement message deduplication with bounded memory?"
4. "A transaction takes 30 seconds to propagate. How would you debug this?"

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

**Benchmarking goals:**
- [ ] Process 10,000 votes in <100ms
- [ ] View change completes in <3 round trips
- [ ] Validator set updates in O(log n) time

**Failure scenarios (critical for interviews):**

*PBFT failure cases:*
- What if you receive 2f+1 prepare messages but the pre-prepare was invalid?
- What if the primary sends different pre-prepares to different replicas?
- What happens during a view change if the new primary is also Byzantine?
- What if network partitions during the commit phase?

*Tendermint failure cases:*
- What if proposer is offline? (timeout → prevote nil → next round)
- What if you see conflicting proposals in the same round?
- What if 1/3 of validators go offline after prevote but before precommit?
- How do you handle equivocation (validator signing two blocks at same height)?

*Finality failure cases:*
- What if a checkpoint is justified but never finalized?
- How do you handle a finality reversion (theoretically impossible but...)?
- What's the "nothing at stake" problem and how does slashing address it?

**Interview questions:**
1. "Walk me through what happens when a PBFT primary is Byzantine."
2. "Why does Tendermint need both prevote and precommit phases?"
3. "Explain the 2/3 threshold - why not 1/2 or 3/4?"
4. "How does Casper FFG achieve finality? What are the slashing conditions?"
5. "A validator sees two conflicting blocks at the same height. What should it do?"

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

**Benchmarking goals:**
- [ ] Proof generation for 10,000 leaves in <100ms
- [ ] Proof verification in <1ms
- [ ] Sparse tree update in O(log n) time

**Interview questions:**
1. "How would you optimize proof generation for a block with 50,000 transactions?"
2. "Explain how Ethereum's state proofs work for a light client querying an account balance."
3. "What's the difference between a binary Merkle tree and a sparse Merkle tree? When would you use each?"

---

**Patricia Tries (MPT):**
- Ethereum's Modified Patricia Trie
- Hex-prefix encoding
- State root calculation

**Benchmarking goals:**
- [ ] 10,000 key-value insertions in <500ms
- [ ] State root calculation in <100ms
- [ ] Proof generation in <10ms

**Interview questions:**
1. "Why does Ethereum use a Patricia trie instead of a regular Merkle tree?"
2. "Explain hex-prefix encoding and why it's needed."
3. "How would you implement efficient batch updates to the state trie?"

---

**Bloom Filters:**
- Insert/query
- Ethereum log bloom
- Light client filtering

**Benchmarking goals:**
- [ ] 1M insertions in <100ms
- [ ] False positive rate <1% with optimal parameters
- [ ] Query time <100ns

**Interview questions:**
1. "How does Ethereum use bloom filters in block headers?"
2. "A light client wants to find all Transfer events for an address. Explain the process."
3. "How do you choose optimal bloom filter parameters?"

---

**DAGs:**
- UTXO graph
- Conflict detection
- DAG-based consensus structures

**Benchmarking goals:**
- [ ] Conflict detection in <1ms for 1000 pending transactions
- [ ] Topological ordering of 10K transactions in <50ms

**Interview questions:**
1. "How does a DAG-based consensus (like Narwhal) differ from chain-based?"
2. "Explain how IOTA's Tangle handles conflicting transactions."
3. "What are the tradeoffs of DAG vs chain for transaction throughput?"

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
