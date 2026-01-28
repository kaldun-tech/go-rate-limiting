# Sparse Merkle Tree: Bit Indexing and Node Addressing

## Overview

In a sparse Merkle tree, keys are binary data (e.g., 256-bit hashes). Each bit of the key determines the path from root to leaf: 0 = go left, 1 = go right.

This document explains how to extract individual bits from a `[]byte` key and how to address nodes in the tree.

---

## Part 1: Bit Indexing

### The Problem

Given a key as `[]byte`, we need to access individual bits by position.

Example key: `[]byte{0xA5, 0x3C}` (2 bytes = 16 bits)

```
Byte index:    [0]                 [1]
Hex:           0xA5                0x3C
Binary:        1 0 1 0 0 1 0 1     0 0 1 1 1 1 0 0
Bit position:  0 1 2 3 4 5 6 7     8 9 ...
               ↑ MSB       LSB ↑
```

**MSB (Most Significant Bit)**: The leftmost bit in a byte, has the highest value (128 in an 8-bit byte).

**LSB (Least Significant Bit)**: The rightmost bit, has the lowest value (1).

We use MSB-first ordering so bit 0 is the first decision from the root.

---

### Step 1: Find the Byte

```go
byteIdx := pos / 8
```

Integer division tells you which byte contains the bit:

| Bit Position | pos / 8 | Byte Index |
|--------------|---------|------------|
| 0-7          | 0       | 0          |
| 8-15         | 1       | 1          |
| 16-23        | 2       | 2          |

---

### Step 2: Find the Bit Within the Byte

```go
bitIdx := 7 - (pos % 8)
```

This converts from "position from left" to "position from right" (which is how bit operations work).

```
pos % 8 (position within byte, from left):  0  1  2  3  4  5  6  7
bitIdx (position from right, for shifts):   7  6  5  4  3  2  1  0
```

**Why the conversion?** Bit operations like shifts and masks work from the right (LSB), but we want to read bits from the left (MSB) for tree traversal.

---

### Step 3: Extract the Bit

Two equivalent ways:

```go
// Option A: Shift right, then mask
bit := (key[byteIdx] >> bitIdx) & 1

// Option B: Mask in place, check non-zero
isSet := (key[byteIdx] & (1 << bitIdx)) != 0
```

**Visual example** - Extract bit at pos=2 from `0xA5`:

```
pos = 2
byteIdx = 2 / 8 = 0
bitIdx = 7 - (2 % 8) = 7 - 2 = 5

key[0] = 0xA5 = 1 0 1 0 0 1 0 1
                    ↑
                  pos=2 (bitIdx=5 from right)

Mask:  1 << 5    = 0 0 1 0 0 0 0 0

       0xA5      = 1 0 1 0 0 1 0 1
   AND mask      = 0 0 1 0 0 0 0 0
       result    = 0 0 1 0 0 0 0 0  (non-zero → bit is 1)
```

---

### Step 4: Flip a Bit (for finding siblings)

```go
key[byteIdx] ^= (1 << bitIdx)
```

XOR with 1 flips the bit: 0→1, 1→0

```
       0xA5      = 1 0 1 0 0 1 0 1
   XOR mask      = 0 0 1 0 0 0 0 0   (1 << 5)
       result    = 1 0 0 0 0 1 0 1   = 0x85
                       ↑
                   bit flipped!
```

---

### Complete Helper Function

```go
// isRightChild returns true if the path goes right at the given bit position
func isRightChild(key []byte, pos int) bool {
    byteIdx := pos / 8
    bitIdx := uint(7 - (pos % 8))
    return (key[byteIdx] & (1 << bitIdx)) != 0
}
```

---

## Part 2: Tree Levels and Bit Positions

### Level Numbering

```
Level depth (root):        [ROOT]
                          /      \
Level depth-1:          [L]      [R]      ← bit 0 decides
                       /   \    /   \
Level depth-2:       [LL] [LR][RL] [RR]   ← bit 1 decides
                      ...
Level 1:              ...                  ← bit depth-2 decides
Level 0 (leaves):    [leaf nodes]          ← bit depth-1 decides
```

### Mapping Levels to Bit Positions

When walking **up** the tree (leaf to root) at level `L`:

```go
bitPos := depth - 1 - level
```

| Level | Bit Position (depth=8) | Description |
|-------|------------------------|-------------|
| 0     | 7                      | Just above leaves |
| 1     | 6                      | |
| ...   | ...                    | |
| 6     | 1                      | |
| 7     | 0                      | Root's children |

---

## Part 3: Node Addressing

### The Problem

We need to store intermediate node hashes in a map. How do we create unique keys?

### Solution: Level + Key

Use `fmt.Sprintf("%d:%x", level, key)` as the map key.

**Why this works:**
- The full key identifies which leaf path we're on
- The level identifies how far up that path
- Siblings at the same level have different keys (they differ in the bit at that level)

```
Key = 0xA5 (binary: 10100101), depth = 8

Path from root to leaf:
  Root
   └─[1]→ Right (bit 0 = 1)
       └─[0]→ Left (bit 1 = 0)
           └─[1]→ Right (bit 2 = 1)
               └─[0]→ Left (bit 3 = 0)
                   └─... down to leaf

Node addresses along this path:
  "7:a5" - level 7 (root's child)
  "6:a5" - level 6
  "5:a5" - level 5
  ...
  "0:a5" - level 0 (leaf)
```

### Finding Sibling's Address

To find the sibling at a given level:
1. Copy the key
2. Flip the bit at `depth - 1 - level`
3. Use the same level

```go
func (t *SparseMerkleTree) getSiblingHash(level int, key []byte) []byte {
    keyCopy := make([]byte, len(key))
    copy(keyCopy, key)

    bitPos := t.depth - 1 - level
    byteIdx := bitPos / 8
    bitIdx := uint(7 - (bitPos % 8))

    keyCopy[byteIdx] ^= (1 << bitIdx)  // Flip the bit

    return t.getNodeHash(level, keyCopy)
}
```

---

## Part 4: Putting It Together - The Set Operation

```go
func (t *SparseMerkleTree) Set(key, value []byte) error {
    // 1. Store the value
    t.values[string(key)] = value

    // 2. Compute leaf hash
    currentHash := sha256.Sum256(value)

    // 3. Walk from leaf (level 0) to root (level depth)
    for level := 0; level < t.depth; level++ {
        // Which bit determines left/right at this level?
        bitPos := t.depth - 1 - level
        isRight := isRightChild(key, bitPos)

        // Get sibling hash (from nodes map or default)
        siblingHash := t.getSiblingHash(level, key)

        // Combine in correct order
        var parentHash []byte
        if isRight {
            parentHash = hashPair(siblingHash, currentHash[:])
        } else {
            parentHash = hashPair(currentHash[:], siblingHash)
        }

        // Store this intermediate hash
        t.setNodeHash(level+1, key, parentHash)

        // Move up
        copy(currentHash[:], parentHash)
    }

    // 4. Update root
    t.root = currentHash[:]
    return nil
}
```

---

## Quick Reference

| Operation | Formula |
|-----------|---------|
| Byte containing bit | `byteIdx = pos / 8` |
| Bit position in byte | `bitIdx = 7 - (pos % 8)` |
| Extract bit | `(key[byteIdx] >> bitIdx) & 1` |
| Check if bit set | `(key[byteIdx] & (1 << bitIdx)) != 0` |
| Flip bit | `key[byteIdx] ^= (1 << bitIdx)` |
| Level to bit position | `bitPos = depth - 1 - level` |
| Node map key | `fmt.Sprintf("%d:%x", level, key)` |
