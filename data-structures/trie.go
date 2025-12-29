package datastructures

// TrieNode represents a node in a Trie (Prefix Tree)
type TrieNode struct {
	children map[rune]*TrieNode
	isEnd    bool
}

// Trie implements a prefix tree for efficient string operations
// Common interview applications:
// - Autocomplete
// - Spell checker
// - IP routing tables
// - Word search games
type Trie struct {
	root *TrieNode
}

// NewTrie creates an empty Trie
func NewTrie() *Trie {
	return &Trie{
		root: &TrieNode{
			children: make(map[rune]*TrieNode),
		},
	}
}

// Insert adds a word to the trie
func (t *Trie) Insert(word string) {
	// TODO: Implement
	// Hint: Traverse character by character, creating nodes as needed
}

// Search checks if a word exists in the trie
// Must be a complete word (isEnd == true)
func (t *Trie) Search(word string) bool {
	// TODO: Implement
	return false
}

// StartsWith checks if any word in the trie starts with the given prefix
func (t *Trie) StartsWith(prefix string) bool {
	// TODO: Implement
	return false
}

// Delete removes a word from the trie
// Advanced: clean up unused nodes
func (t *Trie) Delete(word string) bool {
	// TODO: Implement
	return false
}

// FindAllWithPrefix returns all words that start with the given prefix
// Useful for autocomplete
func (t *Trie) FindAllWithPrefix(prefix string) []string {
	// TODO: Implement
	return nil
}
