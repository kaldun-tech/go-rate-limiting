package datastructures

// PatriciaTrie implements a Modified Patricia Trie (MPT)
// https://ethereum.org/en/developers/docs/data-structures-and-encoding/patricia-merkle-trie/
// Blockchain uses:
// - Ethereum state trie
// - Storage trie
// - Transaction trie
// - Receipt trie
type PatriciaTrie struct {
	root PatriciaNode
}

// PatriciaNode represents a node in the Patricia trie
// Can be one of: EmptyNode, LeafNode, ExtensionNode, BranchNode
type PatriciaNode interface {
	Hash() []byte
}

// EmptyNode represents an empty node
type EmptyNode struct{}

// LeafNode stores a key suffix and value
type LeafNode struct {
	Path  []byte // remaining key nibbles (hex-prefix encoded)
	Value []byte
}

// ExtensionNode stores a shared path prefix
type ExtensionNode struct {
	Path  []byte       // shared nibbles (hex-prefix encoded)
	Child PatriciaNode
}

// BranchNode has 16 children (one for each hex nibble) plus an optional value
type BranchNode struct {
	Children [16]PatriciaNode
	Value    []byte // non-nil if this node also stores a value
}

// NewPatriciaTrie creates an empty Patricia trie
func NewPatriciaTrie() *PatriciaTrie {
	// TODO: implement
	panic("not implemented")
}

// Get retrieves the value for a given key
func (t *PatriciaTrie) Get(key []byte) ([]byte, bool) {
	// TODO: implement
	// Hint: Convert key to nibbles, traverse trie matching path
	panic("not implemented")
}

// Put inserts or updates a key-value pair
func (t *PatriciaTrie) Put(key, value []byte) {
	// TODO: implement
	// Hint: Convert key to nibbles, find insertion point, restructure nodes
	panic("not implemented")
}

// Delete removes a key from the trie
func (t *PatriciaTrie) Delete(key []byte) bool {
	// TODO: implement
	// Hint: May need to merge extension nodes after deletion
	panic("not implemented")
}

// RootHash returns the hash of the root node (state root)
func (t *PatriciaTrie) RootHash() []byte {
	// TODO: implement
	// Hint: Recursively hash nodes using RLP encoding
	panic("not implemented")
}

// GenerateProof creates a Merkle proof for a key
func (t *PatriciaTrie) GenerateProof(key []byte) ([][]byte, error) {
	// TODO: implement
	// Hint: Collect RLP-encoded nodes along the path
	panic("not implemented")
}

// VerifyProof verifies a Patricia trie proof
func VerifyPatriciaProof(rootHash, key []byte, proof [][]byte) ([]byte, error) {
	// TODO: implement
	panic("not implemented")
}

// Hex-prefix encoding helpers

// HexPrefixEncode encodes nibbles with a flag indicating leaf vs extension
func HexPrefixEncode(nibbles []byte, isLeaf bool) []byte {
	// TODO: implement
	// Hint: First nibble encodes flags (odd length + leaf/extension)
	panic("not implemented")
}

// HexPrefixDecode decodes hex-prefix encoded bytes to nibbles
func HexPrefixDecode(encoded []byte) (nibbles []byte, isLeaf bool) {
	// TODO: implement
	panic("not implemented")
}

// KeyToNibbles converts a byte key to nibbles (half-bytes)
func KeyToNibbles(key []byte) []byte {
	// TODO: implement
	// Hint: Each byte becomes two nibbles
	panic("not implemented")
}

// Hash implementations for node types

func (n *EmptyNode) Hash() []byte {
	// TODO: implement (empty node has a specific hash)
	panic("not implemented")
}

func (n *LeafNode) Hash() []byte {
	// TODO: implement (RLP encode then keccak256)
	panic("not implemented")
}

func (n *ExtensionNode) Hash() []byte {
	// TODO: implement
	panic("not implemented")
}

func (n *BranchNode) Hash() []byte {
	// TODO: implement
	panic("not implemented")
}
