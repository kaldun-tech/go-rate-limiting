package datastructures

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
)

// MerkleTree implements a binary Merkle tree
// https://en.wikipedia.org/wiki/Merkle_tree
// Blockchain uses:
// - Transaction root in block headers
// - State verification
// - Light client proofs
type MerkleTree struct {
	root   *MerkleNode
	leaves []*MerkleNode
}

// MerkleNode represents a node in the Merkle tree
type MerkleNode struct {
	Hash   []byte
	Left   *MerkleNode
	Right  *MerkleNode
	Parent *MerkleNode
}

// MerkleProof contains the data needed to verify inclusion
type MerkleProof struct {
	Siblings [][]byte // sibling hashes along the path
	// PathBits indicates the position of the sibling
	// true -> sibling is on the right, false -> sibling is on left
	PathBits []bool
	LeafHash []byte
	RootHash []byte
}

// NewMerkleTree builds a Merkle tree from a list of data items
// Time: O(n)
func NewMerkleTree(data [][]byte) *MerkleTree {
	tree := &MerkleTree{
		leaves: []*MerkleNode{},
	}
	if len(data) == 0 {
		// Edge case: empty data -> return empty or nil
		return tree
	}

	// Hash each data item to create leaves
	for _, b := range data {
		hash := sha256.Sum256(b)
		leaf := &MerkleNode{
			Hash: hash[:], // Slice syntax converts [32]byte to []byte
		}
		tree.leaves = append(tree.leaves, leaf)
	}

	if len(tree.leaves) == 1 {
		// Edge case: single node -> set root and return
		tree.root = tree.leaves[0]
		return tree
	}

	// Pair and hash up the tree
	nextLevel := buildParentLevel(tree.leaves)
	for 1 < len(nextLevel) {
		nextLevel = buildParentLevel(nextLevel)
	}
	// Behold: one root to rule them all
	tree.root = nextLevel[0]

	return tree
}

// Builds a parent level of the tree from an input set of children
// Done when only one remains
func buildParentLevel(children []*MerkleNode) []*MerkleNode {
	parents := []*MerkleNode{}
	for i := 0; i < len(children); i += 2 {
		left := children[i]
		if len(children) <= i+1 {
			// Handle an odd number of children at the end
			// Option A: duplicate last node
			parents = append(parents, buildParent(left, left))
			// Option B: promote as-is
			// parents = append(parents, left)
		} else {
			// Normal case
			right := children[i+1]
			parents = append(parents, buildParent(left, right))
		}
	}
	return parents
}

// Creates a parent of two child nodes
func buildParent(left, right *MerkleNode) *MerkleNode {
	parent := &MerkleNode{
		Hash:  hashPair(left.Hash, right.Hash),
		Left:  left,
		Right: right,
	}
	left.Parent = parent
	right.Parent = parent
	return parent
}

// Helper function to combine hashes of siblings to build parent
func hashPair(left, right []byte) []byte {
	combined := append(left, right...) // ... syntax expands right array
	hash := sha256.Sum256(combined)
	return hash[:] // convert [32]byte to []byte
}

// Root returns the root hash of the tree
func (t *MerkleTree) Root() []byte {
	if t.root == nil {
		return nil
	}
	return t.root.Hash
}

// GenerateProof creates a Merkle proof for the item at the given index
// Index refers to the position in the leaves slice
// i.e., which original data item you want to prove inclusion for
//
//	If you built the tree with:
//		data := [][]byte{itemA, itemB, itemC, itemD}
//		tree := NewMerkleTree(data)
//		Then GenerateProof(2) would generate a proof that itemC (index 2) is included in the tree
//
// Time: O(log n)
func (t *MerkleTree) GenerateProof(index int) (*MerkleProof, error) {
	if index < 0 || len(t.leaves) <= index {
		// Out of bounds
		return nil, errors.New("Index out of bounds")
	} else if t.root == nil {
		// Shouldn't happen
		return nil, errors.New("Malformed tree")
	}
	leaf := t.leaves[index]
	proof := &MerkleProof{
		LeafHash: leaf.Hash,
		RootHash: t.root.Hash,
		PathBits: []bool{},
		Siblings: [][]byte{},
	}

	// Walk from leaf to root, collecting sibling hashes
	for n := leaf; n.Parent != nil; n = n.Parent {
		var siblingIsRight bool
		var siblingHash []byte
		// Get sibling
		if n == n.Parent.Left {
			// n is left child -> sibling is on right
			siblingIsRight = true
			siblingHash = n.Parent.Right.Hash
		} else {
			// n is right child → sibling is on left
			siblingIsRight = false
			siblingHash = n.Parent.Left.Hash
		}

		proof.PathBits = append(proof.PathBits, siblingIsRight)
		proof.Siblings = append(proof.Siblings, siblingHash)
	}

	return proof, nil
}

// VerifyProof checks if a proof is valid for a given root hash
// Time: O(log n)
func VerifyProof(proof *MerkleProof) bool {
	// Check that the proof is correctly formed
	if proof == nil || len(proof.LeafHash) == 0 || len(proof.RootHash) == 0 || len(proof.Siblings) != len(proof.PathBits) {
		// Error case
		return false
	}

	// Recompute root from leaf using siblings, compare to expected root
	hash := proof.LeafHash
	for i, sib := range proof.Siblings {
		sibIsRight := proof.PathBits[i]
		if sibIsRight {
			// sibling on the right -> parent hash = hash(current || sibling)
			hash = hashPair(hash, sib)
		} else {
			// sibling on left -> parent hash = hash(sibling || current)
			hash = hashPair(sib, hash)
		}
	}

	// Final hash should match expected value at root
	return bytes.Equal(hash, proof.RootHash)
}

// SparseMerkleTree implements a sparse Merkle tree
// Useful for key-value stores with a fixed key space (e.g., 256-bit keys)
// https://medium.com/@kelvinfichter/whats-a-sparse-merkle-tree-acda70aeb837
// Blockchain uses:
// - Account state trees
// - Storage tries
type SparseMerkleTree struct {
	root          []byte
	depth         int
	defaultHashes [][]byte          // precomputed hashes for empty subtrees
	nodes         map[string][]byte // stores non-default intermediate hashes
}

// NewSparseMerkleTree creates a sparse Merkle tree with the given depth
func NewSparseMerkleTree(depth int) *SparseMerkleTree {
	tree := &SparseMerkleTree{
		depth:         depth,
		defaultHashes: make([][]byte, depth+1),
		nodes:         make(map[string][]byte),
	}

	// Precompute default hashes for each level
	// defaultHashes[i] = hash of an empty subtree of height i
	// Since all empty subtrees of the same height are identical,
	// we only need one hash per level (this is the "sparse" optimization)

	// Level 0: hash of empty leaf value
	emptyLeaf := sha256.Sum256([]byte{})
	tree.defaultHashes[0] = emptyLeaf[:]

	// Each subsequent level: hash(emptyChild || emptyChild)
	for i := 1; i <= depth; i++ {
		tree.defaultHashes[i] = hashPair(tree.defaultHashes[i-1], tree.defaultHashes[i-1])
	}

	tree.root = tree.defaultHashes[depth]
	return tree
}

// Get retrieves the value at the given key
// The key is a cryptographic hash as 256-bit value
// Go strings are immutable byte sequences so we can use the []byte as a string
func (t *SparseMerkleTree) Get(key []byte) ([]byte, error) {
	err := t.validateKey(key)
	if err != nil {
		return nil, err
	}

	val, exists := t.nodes[string(key)]
	if !exists {
		return nil, errors.New("Not found for key")
	}
	return val, nil
}

// Set updates the value at the given key
// Key bits indicate which direction to go at a given level: 0 for left, 1 for right
func (t *SparseMerkleTree) Set(key, value []byte) error {
	err := t.validateKey(key)
	if err != nil {
		return err
	}

	// Set the value of the node directly for Get to retrieve the raw value
	t.nodes[string(key)] = value
	// Hash the value of the leaf and store it for proof generation
	currentHash := sha256.Sum256(value)
	t.setNodeHash(0, key, currentHash[:])

	// Walk from leaf (level 0) to root (level depth). Use a new variable to not overwrite stored value
	computedHash := currentHash[:]
	for i := range t.depth {
		// Determine whether this node is left or right child from key bit at level
		pos := t.depth - 1 - i
		isRight := isRightChild(key, pos)
		// Get sibling hash from nodes map or defaultHashes[level] if sibling is empty
		siblingHash := t.getSiblingHash(i, key)

		// Combine hash of left and right
		var newHash []byte
		if isRight {
			newHash = hashPair(siblingHash, computedHash)
		} else {
			newHash = hashPair(computedHash, siblingHash)
		}

		// Store fresh slice result in nodes for that position
		t.setNodeHash(i+1, key, newHash)
		computedHash = newHash
	}

	// Update the root with the final hash after the loop
	t.root = computedHash
	return nil
}

// GenerateProof creates a proof of inclusion/exclusion for a leaf key
// Returns tuple: proof pointer, error if one occurs
func (t *SparseMerkleTree) GenerateProof(key []byte) (*MerkleProof, error) {
	err := t.validateKey(key)
	if err != nil {
		return nil, err
	}

	// Get the value or default of input leaf key
	leafHash := t.getNodeHash(0, key)

	// Build proof
	proof := &MerkleProof{
		Siblings: [][]byte{},
		PathBits: []bool{},
		RootHash: t.root,
		LeafHash: leafHash,
	}

	// Loop through levels 0 through depth - 1 collecting siblings
	for level := range t.depth {
		proof.Siblings = append(proof.Siblings, t.getSiblingHash(level, key))
		// Note relationship with isRightChild: current node isRight -> sibling on the left
		bitPos := t.depth - 1 - level
		proof.PathBits = append(proof.PathBits, !isRightChild(key, bitPos))
	}

	return proof, nil
}

func (t *SparseMerkleTree) validateKey(key []byte) error {
	if key == nil {
		return errors.New("Key cannot be nil")
	}

	// For depth d we need d bits to navigate the tree
	// Therefore key must be >= ceil(d/8) bytes long
	expectedLen := (t.depth + 7) / 8
	if len(key) != expectedLen {
		return errors.New("Invalid key length")
	}
	return nil
}

// Uses bit extraction. Returns true if bit at position is 1, false if 0
// pos 0 = first bit used from root level to initial set of children - Most significant bit (MSB)
func isRightChild(key []byte, pos int) bool {
	byteIdx := pos / 8
	bitIdx := uint(7 - (pos % 8))
	return (key[byteIdx] & (1 << bitIdx)) != 0
}

// Retrieve from nodes map, falling back to defaultHashes
func (t *SparseMerkleTree) getNodeHash(level int, key []byte) []byte {
	// For node addressing use a simple map key identifying the leaf path and level up the path
	// Works because siblings at the same level will have different keys because they differ in the bit at that level
	kStr := fmt.Sprintf("%d:%x", level, key)
	val, exists := t.nodes[kStr]
	if exists {
		return val
	}

	// Return default value
	return t.defaultHashes[level]
}

// Store in nodes map. Pretty simple
func (t *SparseMerkleTree) setNodeHash(level int, key, hash []byte) {
	kStr := fmt.Sprintf("%d:%x", level, key)
	t.nodes[kStr] = hash
}

// Get sibling hash - flip the relevant bit and look up that node
func (t *SparseMerkleTree) getSiblingHash(level int, key []byte) []byte {
	// Make defensive copy
	keyCopy := make([]byte, len(key))
	copy(keyCopy, key)
	// Which bit determines left or right at this level
	bitPos := t.depth - 1 - level
	byteIdx := bitPos / 8
	// Most significant bit is first within byte
	bitIdx := uint(7 - (bitPos % 8))
	// Flip the bit using XOR
	keyCopy[byteIdx] ^= (1 << bitIdx)

	return t.getNodeHash(level, keyCopy)
}
