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

func (t *Trie) DeleteIterative(word string) bool {
	// Build a stack with the search path to the word including the dummy root
	n := t.root
	var searchPath = []*TrieNode{t.root}
	for _, val := range word {
		if n.children[val] == nil {
			// Not found
			return false
		}
		// Continue search
		searchPath = append(searchPath, n.children[val])
		n = n.children[val]
	}

	// Now we have a stack with nodes from the root to the parent of end n
	if !n.isEnd {
		// Last node n is not ending -> word not found
		// Example: trie contains only "cart", trying to delete "car"
		return false
	}
	// Mark as non-ending to delete the word
	n.isEnd = false

	// Iterate backwards to delete up the stack
	for i := len(word) - 1; 0 <= i; i-- {
		val := []rune(word)[i]
		parent := searchPath[i]
		if len(n.children) == 0 && !n.isEnd {
			// Childless non-ending node -> delete this and check parent
			delete(parent.children, val)
			n = parent
		} else {
			// Node has children or ends a word -> done deleting
			return true
		}
	}

	// Deleted everything going up -> empty trie
	return true
}

// Recursive helper to delete
// Returns tuple: is word found and deleted, is child node deleted
func (t *Trie) DeleteRecursive(word string, pos int, n *TrieNode) (bool, bool) {
	if len(word) < pos {
		// Index out of bounds
		return false, false
	}
	if len(word) == pos {
		// Final letter in word -> attempt deletion
		if !n.isEnd {
			// Does not end word -> not found
			// Example: trie contains only "cart", trying to delete "car"
			return false, false
		}
		n.isEnd = false

		// Delete going up the stack only if node is childless
		return true, 0 == len(n.children)
	}

	// Check children and recurse
	val := []rune(word)[pos]
	if n.children[val] == nil {
		// Word not found
		return false, false
	}
	deleteWord, deleteChild := t.DeleteRecursive(word, pos+1, n.children[val])
	if !deleteWord {
		// Word not found
		return false, false
	}

	if deleteChild {
		delete(n.children, val)
		if len(n.children) == 0 && !n.isEnd {
			// No children and non-ending -> mark this for deletion by parent
			return true, true
		}
	}
	// Word is deleted, do not delete up the stack
	return true, false
}

// Delete removes a word from the trie
// Advanced: clean up unused nodes
func (t *Trie) Delete(word string) bool {
	return t.DeleteIterative(word)
}

// Find suffixes of a prefix using DFS
func (t *Trie) findSuffixesDFS(prefix string, suffixSoFar string, words *[]string, n *TrieNode) {
	if n.isEnd {
		// Found ending of a word
		*words = append(*words, prefix+suffixSoFar)
	}
	for rune, child := range n.children {
		t.findSuffixesDFS(prefix, suffixSoFar+string(rune), words, child)
	}
}

// FindAllWithPrefix returns all words that start with the given prefix
// Useful for autocomplete
func (t *Trie) FindAllWithPrefix(prefix string) []string {
	words := []string{}

	// DFS like StartsWith to find the node containing prefix
	n := t.root
	for _, val := range prefix {
		if n.children[val] == nil {
			// Not found
			return words
		}
		// Continue searching
		n = n.children[val]
	}

	// Continue to use DFS to find all possible suffix endings (possibly only empty string)
	// Choose to build string as we go. It will be N concatenations either way for N suffixes.
	t.findSuffixesDFS(prefix, "", &words, n)

	return words
}
