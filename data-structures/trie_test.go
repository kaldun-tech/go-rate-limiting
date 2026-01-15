package datastructures

import (
	"reflect"
	"sort"
	"testing"
)

func TestTrie_InsertAndSearch(t *testing.T) {
	trie := NewTrie()

	// Search in empty trie
	if trie.Search("hello") {
		t.Error("Search in empty trie should return false")
	}

	// Insert and search single word
	trie.Insert("hello")
	if !trie.Search("hello") {
		t.Error("Search(hello) should return true after insert")
	}

	// Search for non-existing word
	if trie.Search("world") {
		t.Error("Search(world) should return false")
	}

	// Search for prefix (not a complete word)
	if trie.Search("hell") {
		t.Error("Search(hell) should return false - not a complete word")
	}

	// Insert overlapping words
	trie.Insert("hell")
	if !trie.Search("hell") {
		t.Error("Search(hell) should return true after insert")
	}
	if !trie.Search("hello") {
		t.Error("Search(hello) should still return true")
	}
}

func TestTrie_InsertAndSearch_SharedPrefixes(t *testing.T) {
	trie := NewTrie()

	words := []string{"car", "card", "care", "careful", "cart"}
	for _, w := range words {
		trie.Insert(w)
	}

	// All inserted words should be found
	for _, w := range words {
		if !trie.Search(w) {
			t.Errorf("Search(%s) should return true", w)
		}
	}

	// Partial prefixes should not be found as words
	if trie.Search("ca") {
		t.Error("Search(ca) should return false - not a complete word")
	}
	if trie.Search("caref") {
		t.Error("Search(caref) should return false - not a complete word")
	}
}

func TestTrie_StartsWith(t *testing.T) {
	trie := NewTrie()

	// Empty trie
	if trie.StartsWith("a") {
		t.Error("StartsWith on empty trie should return false")
	}

	trie.Insert("hello")
	trie.Insert("help")
	trie.Insert("world")

	// Valid prefixes
	if !trie.StartsWith("hel") {
		t.Error("StartsWith(hel) should return true")
	}
	if !trie.StartsWith("hello") {
		t.Error("StartsWith(hello) should return true - exact match is valid prefix")
	}
	if !trie.StartsWith("w") {
		t.Error("StartsWith(w) should return true")
	}

	// Invalid prefixes
	if trie.StartsWith("hex") {
		t.Error("StartsWith(hex) should return false")
	}
	if trie.StartsWith("helloo") {
		t.Error("StartsWith(helloo) should return false - extends beyond word")
	}
}

func TestTrie_Delete_NonExistent(t *testing.T) {
	trie := NewTrie()

	// Delete from empty trie
	if trie.Delete("hello") {
		t.Error("Delete from empty trie should return false")
	}

	trie.Insert("hello")

	// Delete non-existent word
	if trie.Delete("world") {
		t.Error("Delete(world) should return false - word not in trie")
	}

	// Original word should still exist
	if !trie.Search("hello") {
		t.Error("Search(hello) should still return true")
	}
}

func TestTrie_Delete_PrefixOfAnother(t *testing.T) {
	trie := NewTrie()
	trie.Insert("car")
	trie.Insert("cart")

	// Delete "car" which is prefix of "cart"
	if !trie.Delete("car") {
		t.Error("Delete(car) should return true")
	}

	if trie.Search("car") {
		t.Error("Search(car) should return false after deletion")
	}

	// "cart" should still exist
	if !trie.Search("cart") {
		t.Error("Search(cart) should still return true")
	}
}

func TestTrie_Delete_WordWithPrefixInTrie(t *testing.T) {
	trie := NewTrie()
	trie.Insert("car")
	trie.Insert("cart")

	// Delete "cart" which has prefix "car" in trie
	if !trie.Delete("cart") {
		t.Error("Delete(cart) should return true")
	}

	if trie.Search("cart") {
		t.Error("Search(cart) should return false after deletion")
	}

	// "car" should still exist
	if !trie.Search("car") {
		t.Error("Search(car) should still return true")
	}
}

func TestTrie_Delete_OnlyWord(t *testing.T) {
	trie := NewTrie()
	trie.Insert("hello")

	if !trie.Delete("hello") {
		t.Error("Delete(hello) should return true")
	}

	if trie.Search("hello") {
		t.Error("Search(hello) should return false after deletion")
	}

	// Trie should be empty - no prefixes should match
	if trie.StartsWith("h") {
		t.Error("StartsWith(h) should return false - trie should be empty")
	}
}

func TestTrie_Delete_SharedPrefixes(t *testing.T) {
	trie := NewTrie()
	trie.Insert("cat")
	trie.Insert("car")

	// Delete "cat", "car" should remain
	if !trie.Delete("cat") {
		t.Error("Delete(cat) should return true")
	}

	if trie.Search("cat") {
		t.Error("Search(cat) should return false after deletion")
	}

	if !trie.Search("car") {
		t.Error("Search(car) should still return true")
	}

	// "ca" prefix should still exist
	if !trie.StartsWith("ca") {
		t.Error("StartsWith(ca) should return true")
	}
}

func TestTrie_Delete_PrefixNotAWord(t *testing.T) {
	trie := NewTrie()
	trie.Insert("cart")

	// Try to delete "car" which exists as prefix but not as word
	if trie.Delete("car") {
		t.Error("Delete(car) should return false - not a complete word in trie")
	}

	// "cart" should still exist
	if !trie.Search("cart") {
		t.Error("Search(cart) should still return true")
	}
}

func TestTrie_FindAllWithPrefix(t *testing.T) {
	trie := NewTrie()

	// Empty trie
	if got := trie.FindAllWithPrefix("a"); len(got) != 0 {
		t.Errorf("FindAllWithPrefix on empty trie = %v, want []", got)
	}

	trie.Insert("car")
	trie.Insert("card")
	trie.Insert("care")
	trie.Insert("cart")
	trie.Insert("cat")

	// Find all words with prefix "car"
	got := trie.FindAllWithPrefix("car")
	want := []string{"car", "card", "care", "cart"}
	sort.Strings(got)
	sort.Strings(want)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("FindAllWithPrefix(car) = %v, want %v", got, want)
	}

	// Find all words with prefix "ca"
	got = trie.FindAllWithPrefix("ca")
	want = []string{"car", "card", "care", "cart", "cat"}
	sort.Strings(got)
	sort.Strings(want)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("FindAllWithPrefix(ca) = %v, want %v", got, want)
	}

	// Exact match prefix
	got = trie.FindAllWithPrefix("cat")
	want = []string{"cat"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("FindAllWithPrefix(cat) = %v, want %v", got, want)
	}

	// Non-existent prefix
	got = trie.FindAllWithPrefix("dog")
	if len(got) != 0 {
		t.Errorf("FindAllWithPrefix(dog) = %v, want []", got)
	}
}

func TestTrie_FindAllWithPrefix_EmptyPrefix(t *testing.T) {
	trie := NewTrie()
	trie.Insert("a")
	trie.Insert("ab")
	trie.Insert("b")

	// Empty prefix should return all words
	got := trie.FindAllWithPrefix("")
	want := []string{"a", "ab", "b"}
	sort.Strings(got)
	sort.Strings(want)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("FindAllWithPrefix('') = %v, want %v", got, want)
	}
}

func TestTrie_Unicode(t *testing.T) {
	trie := NewTrie()

	// Test with unicode characters
	trie.Insert("café")
	trie.Insert("日本語")
	trie.Insert("emoji😀")

	if !trie.Search("café") {
		t.Error("Search(café) should return true")
	}
	if !trie.Search("日本語") {
		t.Error("Search(日本語) should return true")
	}
	if !trie.Search("emoji😀") {
		t.Error("Search(emoji😀) should return true")
	}

	if !trie.StartsWith("caf") {
		t.Error("StartsWith(caf) should return true")
	}
	if !trie.StartsWith("日本") {
		t.Error("StartsWith(日本) should return true")
	}
}

func TestTrie_EmptyString(t *testing.T) {
	trie := NewTrie()

	// Insert empty string
	trie.Insert("")

	if !trie.Search("") {
		t.Error("Search('') should return true after inserting empty string")
	}

	// Other words should still work
	trie.Insert("hello")
	if !trie.Search("hello") {
		t.Error("Search(hello) should return true")
	}
	if !trie.Search("") {
		t.Error("Search('') should still return true")
	}
}

// Benchmarks

func BenchmarkTrie_Insert(b *testing.B) {
	words := []string{"apple", "application", "apply", "banana", "band", "bandana"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		trie := NewTrie()
		for _, w := range words {
			trie.Insert(w)
		}
	}
}

func BenchmarkTrie_Search(b *testing.B) {
	trie := NewTrie()
	words := []string{"apple", "application", "apply", "banana", "band", "bandana"}
	for _, w := range words {
		trie.Insert(w)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		trie.Search(words[i%len(words)])
	}
}

func BenchmarkTrie_Delete(b *testing.B) {
	words := []string{"apple", "application", "apply", "banana", "band", "bandana"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		trie := NewTrie()
		for _, w := range words {
			trie.Insert(w)
		}
		b.StartTimer()
		trie.Delete(words[i%len(words)])
	}
}

func BenchmarkTrie_FindAllWithPrefix(b *testing.B) {
	trie := NewTrie()
	words := []string{"apple", "application", "apply", "app", "banana", "band", "bandana"}
	for _, w := range words {
		trie.Insert(w)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		trie.FindAllWithPrefix("app")
	}
}
