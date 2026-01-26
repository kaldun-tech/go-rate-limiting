package datastructures

import (
	"crypto/sha256"
	"errors"
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
	// TODO: implement
	// Hint: Recompute root from leaf using siblings, compare to expected root
	panic("not implemented")
}

// SparseMerkleTree implements a sparse Merkle tree
// Useful for key-value stores with a fixed key space (e.g., 256-bit keys)
// Blockchain uses:
// - Account state trees
// - Storage tries
type SparseMerkleTree struct {
	root          []byte
	depth         int
	defaultHashes [][]byte // precomputed hashes for empty subtrees
}

// NewSparseMerkleTree creates a sparse Merkle tree with the given depth
func NewSparseMerkleTree(depth int) *SparseMerkleTree {
	// TODO: implement
	// Hint: Precompute default hashes for each level
	panic("not implemented")
}

// Get retrieves the value at the given key
func (t *SparseMerkleTree) Get(key []byte) ([]byte, error) {
	// TODO: implement
	panic("not implemented")
}

// Set updates the value at the given key
func (t *SparseMerkleTree) Set(key, value []byte) error {
	// TODO: implement
	panic("not implemented")
}

// GenerateProof creates a proof of inclusion/exclusion for a key
func (t *SparseMerkleTree) GenerateProof(key []byte) (*MerkleProof, error) {
	// TODO: implement
	panic("not implemented")
}
