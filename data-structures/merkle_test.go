package datastructures

import (
	"bytes"
	"crypto/sha256"
	"testing"
)

// ========== MerkleTree Tests ==========

func TestMerkleTree_NewEmpty(t *testing.T) {
	tree := NewMerkleTree([][]byte{})

	if tree.Root() != nil {
		t.Error("Empty tree should have nil root")
	}
}

func TestMerkleTree_SingleElement(t *testing.T) {
	data := [][]byte{[]byte("hello")}
	tree := NewMerkleTree(data)

	expectedHash := sha256.Sum256([]byte("hello"))
	if !bytes.Equal(tree.Root(), expectedHash[:]) {
		t.Error("Single element tree root should equal hash of element")
	}
}

func TestMerkleTree_MultipleElements(t *testing.T) {
	data := [][]byte{
		[]byte("a"),
		[]byte("b"),
		[]byte("c"),
		[]byte("d"),
	}
	tree := NewMerkleTree(data)

	if tree.Root() == nil {
		t.Fatal("Tree root should not be nil")
	}
	if len(tree.Root()) != 32 {
		t.Errorf("Root hash length = %d, want 32", len(tree.Root()))
	}
}

func TestMerkleTree_OddElements(t *testing.T) {
	data := [][]byte{
		[]byte("a"),
		[]byte("b"),
		[]byte("c"),
	}
	tree := NewMerkleTree(data)

	if tree.Root() == nil {
		t.Fatal("Tree root should not be nil")
	}
	if len(tree.Root()) != 32 {
		t.Errorf("Root hash length = %d, want 32", len(tree.Root()))
	}
}

func TestMerkleTree_GenerateProof(t *testing.T) {
	data := [][]byte{
		[]byte("a"),
		[]byte("b"),
		[]byte("c"),
		[]byte("d"),
	}
	tree := NewMerkleTree(data)

	// Generate proof for index 2 ("c")
	proof, err := tree.GenerateProof(2)
	if err != nil {
		t.Fatalf("GenerateProof failed: %v", err)
	}

	if proof == nil {
		t.Fatal("Proof should not be nil")
	}

	expectedLeafHash := sha256.Sum256([]byte("c"))
	if !bytes.Equal(proof.LeafHash, expectedLeafHash[:]) {
		t.Error("Proof leaf hash should match hash of 'c'")
	}

	if !bytes.Equal(proof.RootHash, tree.Root()) {
		t.Error("Proof root hash should match tree root")
	}

	// For 4 elements, tree depth is 2, so we need 2 siblings
	if len(proof.Siblings) != 2 {
		t.Errorf("Siblings count = %d, want 2", len(proof.Siblings))
	}
}

func TestMerkleTree_GenerateProofOutOfBounds(t *testing.T) {
	data := [][]byte{[]byte("a"), []byte("b")}
	tree := NewMerkleTree(data)

	_, err := tree.GenerateProof(-1)
	if err == nil {
		t.Error("GenerateProof with negative index should return error")
	}

	_, err = tree.GenerateProof(5)
	if err == nil {
		t.Error("GenerateProof with out of bounds index should return error")
	}
}

func TestMerkleTree_VerifyProof(t *testing.T) {
	data := [][]byte{
		[]byte("a"),
		[]byte("b"),
		[]byte("c"),
		[]byte("d"),
	}
	tree := NewMerkleTree(data)

	// Test each leaf
	for i := 0; i < len(data); i++ {
		proof, err := tree.GenerateProof(i)
		if err != nil {
			t.Fatalf("GenerateProof(%d) failed: %v", i, err)
		}

		if !VerifyProof(proof) {
			t.Errorf("VerifyProof failed for index %d", i)
		}
	}
}

func TestMerkleTree_VerifyProofTampered(t *testing.T) {
	data := [][]byte{
		[]byte("a"),
		[]byte("b"),
		[]byte("c"),
		[]byte("d"),
	}
	tree := NewMerkleTree(data)

	proof, _ := tree.GenerateProof(0)

	// Tamper with the leaf hash
	fakeHash := sha256.Sum256([]byte("fake"))
	tamperedProof := &MerkleProof{
		LeafHash: fakeHash[:],
		RootHash: proof.RootHash,
		Siblings: proof.Siblings,
		PathBits: proof.PathBits,
	}

	if VerifyProof(tamperedProof) {
		t.Error("VerifyProof should fail for tampered leaf hash")
	}
}

func TestMerkleTree_VerifyProofNil(t *testing.T) {
	if VerifyProof(nil) {
		t.Error("VerifyProof(nil) should return false")
	}
}

// ========== SparseMerkleTree Tests ==========

func TestSparseMerkleTree_New(t *testing.T) {
	tree := NewSparseMerkleTree(8)

	if tree.root == nil {
		t.Fatal("New tree should have a root")
	}
	if len(tree.root) != 32 {
		t.Errorf("Root length = %d, want 32", len(tree.root))
	}
	if tree.depth != 8 {
		t.Errorf("Depth = %d, want 8", tree.depth)
	}
	if len(tree.defaultHashes) != 9 {
		t.Errorf("DefaultHashes length = %d, want 9", len(tree.defaultHashes))
	}
}

func TestSparseMerkleTree_DefaultHashes(t *testing.T) {
	tree := NewSparseMerkleTree(4)

	// Level 0 should be hash of empty bytes
	emptyHash := sha256.Sum256([]byte{})
	if !bytes.Equal(tree.defaultHashes[0], emptyHash[:]) {
		t.Error("Default hash at level 0 should be hash of empty bytes")
	}

	// Each level should be hash of two copies of previous level
	for i := 1; i <= 4; i++ {
		expected := hashPair(tree.defaultHashes[i-1], tree.defaultHashes[i-1])
		if !bytes.Equal(tree.defaultHashes[i], expected) {
			t.Errorf("Default hash at level %d is incorrect", i)
		}
	}
}

func TestSparseMerkleTree_SetGet(t *testing.T) {
	tree := NewSparseMerkleTree(8)
	key := []byte{0x42} // 1 byte for depth 8
	value := []byte("hello")

	err := tree.Set(key, value)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	got, err := tree.Get(key)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if !bytes.Equal(got, value) {
		t.Errorf("Get = %v, want %v", got, value)
	}
}

func TestSparseMerkleTree_SetUpdatesRoot(t *testing.T) {
	tree := NewSparseMerkleTree(8)
	originalRoot := make([]byte, len(tree.root))
	copy(originalRoot, tree.root)

	key := []byte{0x00}
	value := []byte("test")

	err := tree.Set(key, value)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	if bytes.Equal(tree.root, originalRoot) {
		t.Error("Root should change after Set")
	}
}

func TestSparseMerkleTree_GetNotFound(t *testing.T) {
	tree := NewSparseMerkleTree(8)
	key := []byte{0x42}

	_, err := tree.Get(key)
	if err == nil {
		t.Error("Get on non-existent key should return error")
	}
}

func TestSparseMerkleTree_ValidateKeyNil(t *testing.T) {
	tree := NewSparseMerkleTree(8)

	err := tree.Set(nil, []byte("value"))
	if err == nil {
		t.Error("Set with nil key should return error")
	}

	_, err = tree.Get(nil)
	if err == nil {
		t.Error("Get with nil key should return error")
	}
}

func TestSparseMerkleTree_ValidateKeyLength(t *testing.T) {
	tree := NewSparseMerkleTree(8) // requires 1 byte key

	// Too short
	err := tree.Set([]byte{}, []byte("value"))
	if err == nil {
		t.Error("Set with empty key should return error")
	}

	// Too long
	err = tree.Set([]byte{0x00, 0x01}, []byte("value"))
	if err == nil {
		t.Error("Set with too long key should return error")
	}
}

func TestSparseMerkleTree_MultipleKeys(t *testing.T) {
	tree := NewSparseMerkleTree(8)

	keys := [][]byte{{0x00}, {0x01}, {0xFF}, {0x80}}
	values := [][]byte{[]byte("a"), []byte("b"), []byte("c"), []byte("d")}

	// Set all
	for i, key := range keys {
		err := tree.Set(key, values[i])
		if err != nil {
			t.Fatalf("Set key %x failed: %v", key, err)
		}
	}

	// Get all
	for i, key := range keys {
		got, err := tree.Get(key)
		if err != nil {
			t.Fatalf("Get key %x failed: %v", key, err)
		}
		if !bytes.Equal(got, values[i]) {
			t.Errorf("Get key %x = %v, want %v", key, got, values[i])
		}
	}
}

func TestSparseMerkleTree_GenerateProof(t *testing.T) {
	tree := NewSparseMerkleTree(8)
	key := []byte{0x42}
	value := []byte("hello")

	err := tree.Set(key, value)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	proof, err := tree.GenerateProof(key)
	if err != nil {
		t.Fatalf("GenerateProof failed: %v", err)
	}

	if proof == nil {
		t.Fatal("Proof should not be nil")
	}

	if !bytes.Equal(proof.RootHash, tree.root) {
		t.Error("Proof root hash should match tree root")
	}

	if len(proof.Siblings) != 8 {
		t.Errorf("Siblings count = %d, want 8", len(proof.Siblings))
	}

	if len(proof.PathBits) != 8 {
		t.Errorf("PathBits count = %d, want 8", len(proof.PathBits))
	}
}

func TestSparseMerkleTree_GenerateProofInvalidKey(t *testing.T) {
	tree := NewSparseMerkleTree(8)

	_, err := tree.GenerateProof(nil)
	if err == nil {
		t.Error("GenerateProof with nil key should return error")
	}

	_, err = tree.GenerateProof([]byte{0x00, 0x01})
	if err == nil {
		t.Error("GenerateProof with invalid key length should return error")
	}
}

func TestSparseMerkleTree_ProofVerification(t *testing.T) {
	tree := NewSparseMerkleTree(8)
	key := []byte{0x42}
	value := []byte("hello")

	// Debug: what should the leaf hash be?
	expectedLeaf := sha256.Sum256(value)
	t.Logf("sha256(hello) = %x", expectedLeaf)

	err := tree.Set(key, value)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Debug: check what's stored at level 0
	t.Logf("After Set, nodes map has %d entries", len(tree.nodes))

	proof, err := tree.GenerateProof(key)
	if err != nil {
		t.Fatalf("GenerateProof failed: %v", err)
	}

	// Debug: manually verify step by step
	hash := proof.LeafHash
	t.Logf("LeafHash from proof: %x", hash)
	t.Logf("Expected root: %x", proof.RootHash)
	t.Logf("Tree root: %x", tree.root)

	for i, sib := range proof.Siblings {
		sibIsRight := proof.PathBits[i]
		t.Logf("Level %d: sibIsRight=%v, sibling=%x", i, sibIsRight, sib[:8])
		if sibIsRight {
			hash = hashPair(hash, sib)
		} else {
			hash = hashPair(sib, hash)
		}
		t.Logf("  -> hash after combine: %x", hash[:8])
	}
	t.Logf("Final computed hash: %x", hash)

	if !VerifyProof(proof) {
		t.Error("VerifyProof should return true for valid proof")
	}
}

func TestSparseMerkleTree_ProofVerificationMultipleKeys(t *testing.T) {
	tree := NewSparseMerkleTree(8)

	keys := [][]byte{{0x00}, {0x55}, {0xAA}, {0xFF}}
	values := [][]byte{[]byte("a"), []byte("b"), []byte("c"), []byte("d")}

	// Set all values
	for i, key := range keys {
		err := tree.Set(key, values[i])
		if err != nil {
			t.Fatalf("Set key %x failed: %v", key, err)
		}
	}

	// Verify proofs for all keys
	for _, key := range keys {
		proof, err := tree.GenerateProof(key)
		if err != nil {
			t.Fatalf("GenerateProof for key %x failed: %v", key, err)
		}

		if !VerifyProof(proof) {
			t.Errorf("VerifyProof failed for key %x", key)
		}
	}
}

func TestSparseMerkleTree_EmptyKeyProof(t *testing.T) {
	tree := NewSparseMerkleTree(8)
	key := []byte{0x42}

	// Generate proof for key that was never set (should use default hash)
	proof, err := tree.GenerateProof(key)
	if err != nil {
		t.Fatalf("GenerateProof failed: %v", err)
	}

	// Leaf hash should be the default hash at level 0
	if !bytes.Equal(proof.LeafHash, tree.defaultHashes[0]) {
		t.Error("LeafHash for unset key should be default hash")
	}

	if !VerifyProof(proof) {
		t.Error("VerifyProof should return true for empty key proof")
	}
}

func TestSparseMerkleTree_LargerDepth(t *testing.T) {
	tree := NewSparseMerkleTree(16) // 2 byte keys
	key := []byte{0x12, 0x34}
	value := []byte("test value")

	err := tree.Set(key, value)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	got, err := tree.Get(key)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if !bytes.Equal(got, value) {
		t.Errorf("Get = %v, want %v", got, value)
	}

	proof, err := tree.GenerateProof(key)
	if err != nil {
		t.Fatalf("GenerateProof failed: %v", err)
	}

	if len(proof.Siblings) != 16 {
		t.Errorf("Siblings count = %d, want 16", len(proof.Siblings))
	}

	if !VerifyProof(proof) {
		t.Error("VerifyProof failed for larger depth tree")
	}
}

// ========== isRightChild Tests ==========

func TestIsRightChild(t *testing.T) {
	tests := []struct {
		key      []byte
		pos      int
		expected bool
	}{
		{[]byte{0x80}, 0, true},  // 10000000 - bit 0 is 1
		{[]byte{0x80}, 1, false}, // 10000000 - bit 1 is 0
		{[]byte{0x40}, 0, false}, // 01000000 - bit 0 is 0
		{[]byte{0x40}, 1, true},  // 01000000 - bit 1 is 1
		{[]byte{0xFF}, 0, true},  // 11111111 - all bits are 1
		{[]byte{0xFF}, 7, true},
		{[]byte{0x00}, 0, false}, // 00000000 - all bits are 0
		{[]byte{0x00}, 7, false},
		{[]byte{0x01}, 7, true},  // 00000001 - only LSB is 1
		{[]byte{0x01}, 0, false},
	}

	for _, tc := range tests {
		got := isRightChild(tc.key, tc.pos)
		if got != tc.expected {
			t.Errorf("isRightChild(%x, %d) = %v, want %v", tc.key, tc.pos, got, tc.expected)
		}
	}
}
