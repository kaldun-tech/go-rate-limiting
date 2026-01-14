package datastructures

// TrieNode represents a node in a Trie (Prefix Tree)
// Go strings are UTF-8 encoded
// Iterating over a string with range gives runes (int32 equivalent to Java char)
// See https://go.dev/blog/strings
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
// See https://en.wikipedia.org/wiki/Trie
// Root is a dummy node with no value, children hold the values
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

func (t *Trie) insertHelper(word string, n *TrieNode) {
	// val is of type rune (int32 ~ char)
	for _, val := range word {
		if n.children[val] == nil {
			// Missing -> create new child in position
			n.children[val] = &TrieNode{
				children: make(map[rune]*TrieNode),
			}
		}
		n = n.children[val]
	}
	// Mark the final node as end
	n.isEnd = true
}

// Insert adds a word to the trie
func (t *Trie) Insert(word string) {
	// Not helpful to run search first
	// As we may need to set isEnd in a longer word
	t.insertHelper(word, t.root)
}

// Search checks if a word exists in the trie
// Must be a complete word (isEnd == true)
func (t *Trie) Search(word string) bool {
	var n *TrieNode = t.root
	for _, val := range word {
		if n.children[val] == nil {
			// Not found
			return false
		}
		// Persist despite isEnd flag because trie may contain intermediate words
		// Example: searching for card, contains car with r marked as end
		n = n.children[val]
	}

	// Went through and found all matching characters
	// Return true only for complete word
	return n.isEnd
}

// StartsWith checks if any word in the trie starts with the given prefix
func (t *Trie) StartsWith(prefix string) bool {
	n := t.root
	for _, val := range prefix {
		if n.children[val] == nil {
			// Not found
			return false
		}
		// Continue searching
		n = n.children[val]
	}

	// Found all matching characters of the prefix
	return true
}

// Delete removes a word from the trie
// Advanced: clean up unused nodes
func (t *Trie) Delete(word string) bool {
	if !t.Search(word) {
		// Not here
		return false
	}

	// TODO
	return false
}

// FindAllWithPrefix returns all words that start with the given prefix
// Useful for autocomplete
func (t *Trie) FindAllWithPrefix(prefix string) []string {
	words := []string{}

	if !t.StartsWith(prefix) {
		// None found
		return words
	}

	// TODO DFS to find, then BFS to build all words
	return words
}
