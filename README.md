# Blockchain Protocol Engineering Practice

A collection of data structures and algorithms in Go focused on blockchain protocol engineering interview preparation.

## Focus Areas

### Blockchain-Relevant Data Structures

| Data Structure | Blockchain Use Case | Status |
|----------------|---------------------|--------|
| **Merkle Trees** | Block validation, transaction proofs, state roots | TODO |
| **Patricia Tries** | Ethereum state storage (MPT), account/storage tries | TODO |
| **DAGs** | UTXO models, Hedera Hashgraph, IOTA Tangle | TODO |
| **Bloom Filters** | Log filtering, light client sync, transaction lookup | TODO |
| **Hash Tables** | State management, mempool, peer tracking | Partial |

### Protocol-Relevant Algorithms

| Algorithm | Blockchain Use Case | Status |
|-----------|---------------------|--------|
| **Gossip Protocols** | P2P networking, block/tx propagation | TODO |
| **Consensus Algorithms** | BFT variants, PoS mechanisms, finality | TODO |
| **Graph Traversal** | Network topology, transaction flows, UTXO graphs | TODO |
| **Serialization** | RLP encoding, SSZ, protobuf for wire protocols | TODO |
| **Rate Limiting** | P2P DOS protection, API throttling | Done |

## Project Structure

```
go-algorithm-practice/
├── data-structures/
│   ├── merkle.go         # Merkle trees (TODO)
│   ├── patricia.go       # Patricia/MPT tries (TODO)
│   ├── dag.go            # Directed acyclic graphs (TODO)
│   ├── bloom.go          # Bloom filters (TODO)
|   ├── heap.go           # Min/Max Heap (complete)
│   ├── trie.go           # Basic trie (existing)
│   ├── bst.go            # Binary search tree (existing)
│   └── lru-cache.go      # LRU cache (existing)
├── algorithms/
│   ├── gossip/           # Gossip protocol simulation (TODO)
│   ├── consensus/        # BFT/consensus primitives (TODO)
│   ├── serialization/    # RLP, SSZ encoding (TODO)
│   └── graph.go          # Graph algorithms (existing)
├── rate-limiting/
│   └── token-bucket/     # Token bucket rate limiter (complete)
└── examples/
```

## Implementation Roadmap

### Phase 1: Core Data Structures

**Merkle Trees**
- Basic binary Merkle tree with SHA-256
- Proof generation and verification
- Sparse Merkle trees for state proofs
- Multi-proof optimization

**Patricia Tries (MPT)**
- Ethereum's Modified Patricia Trie
- Compact encoding (hex-prefix)
- State root calculation
- Proof generation for light clients

### Phase 2: Network Primitives

**Bloom Filters**
- Basic Bloom filter with configurable false positive rate
- Counting Bloom filters for deletion
- Ethereum log bloom implementation
- Light client block filtering

**Gossip Protocols**
- Epidemic broadcast simulation
- Push/pull gossip variants
- Peer sampling
- Message deduplication

### Phase 3: Consensus & Graphs

**DAGs**
- UTXO graph representation
- Topological ordering
- Conflict detection
- DAG-based consensus (Narwhal-style)

**Consensus Primitives**
- PBFT message flow simulation
- View change protocol
- Finality gadgets
- Validator set management

### Phase 4: Wire Protocols

**Serialization**
- RLP encoding/decoding (Ethereum)
- SSZ for Ethereum 2.0
- Protobuf for inter-node communication
- Benchmarking serialization performance

## What's Implemented

### Token Bucket Rate Limiter

Production-ready rate limiting for P2P and API protection.

```go
import tokenbucket "github.com/kaldun-tech/go-algorithm-practice/rate-limiting/token-bucket"

// Create limiter: 100 requests/second with burst of 20
limiter := tokenbucket.NewTokenBucket(100, time.Second, 20)

// Per-peer rate limiting
if limiter.Allow("peer:" + peerID) {
    handleMessage(msg)
}
```

**Features:**
- Thread-safe with `sync.Mutex`
- Per-key rate limiting (per-peer, per-IP)
- Burst support for traffic spikes
- Weighted costs for different message types

## Development

### Running Tests

```bash
go test ./...                    # All tests
go test -v ./rate-limiting/...   # Rate limiting tests
go test -bench=. ./...           # Benchmarks
```

### Code Quality

```bash
go fmt ./...    # Format
go vet ./...    # Static analysis
go mod tidy     # Clean deps
```

## Why These Topics?

Blockchain protocol engineering interviews focus on:

1. **Data integrity** - Merkle proofs, state roots, cryptographic commitments
2. **Distributed systems** - Consensus, gossip, network partitions
3. **Performance** - Efficient serialization, caching, batching
4. **Security** - DOS protection, rate limiting, input validation

This differs from typical leetcode prep which emphasizes:
- Array/string manipulation
- Dynamic programming patterns
- Generic tree/graph problems

The implementations here are directly applicable to building blockchain clients, validators, and infrastructure.

## Resources

### Specifications
- [Ethereum Yellow Paper](https://ethereum.github.io/yellowpaper/paper.pdf) - MPT, RLP, state model
- [Ethereum Execution Specs](https://github.com/ethereum/execution-specs) - Python reference
- [SSZ Spec](https://github.com/ethereum/consensus-specs/blob/dev/ssz/simple-serialize.md)

### Papers
- [Merkle Trees in Distributed Systems](https://www.cs.cmu.edu/~rwh/theses/okasaki.pdf)
- [PBFT](https://pmg.csail.mit.edu/papers/osdi99.pdf) - Practical Byzantine Fault Tolerance
- [Narwhal and Tusk](https://arxiv.org/abs/2105.11827) - DAG-based consensus

### Implementations
- [go-ethereum](https://github.com/ethereum/go-ethereum) - Ethereum client
- [prysm](https://github.com/prysmaticlabs/prysm) - Ethereum consensus client
- [tendermint](https://github.com/tendermint/tendermint) - BFT consensus engine

## License

MIT
