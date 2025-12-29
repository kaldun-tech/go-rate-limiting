package algorithms

// Backtracking - explore all possibilities by trying options and backtracking when invalid
// Pattern: make choice, explore, undo choice

// Permutations generates all permutations of nums
// LeetCode #46
func Permutations(nums []int) [][]int {
	// TODO: Implement
	// Hint: For each position, try each unused number, backtrack
	return nil
}

// Combinations generates all k-size combinations from 1 to n
// LeetCode #77
func Combinations(n, k int) [][]int {
	// TODO: Implement
	// Hint: For each number, either include it or skip it
	return nil
}

// Subsets generates all subsets of nums
// LeetCode #78
func Subsets(nums []int) [][]int {
	// TODO: Implement
	// Hint: For each element, either include it or don't
	return nil
}

// LetterCombinations of phone number
// '2' -> "abc", '3' -> "def", etc.
// LeetCode #17
func LetterCombinations(digits string) []string {
	// TODO: Implement
	// Hint: Map digits to letters, backtrack through choices
	return nil
}

// GenerateParentheses generates all valid n pairs of parentheses
// Example: n=3 -> ["((()))","(()())","(())()","()(())","()()()"]
// LeetCode #22
func GenerateParentheses(n int) []string {
	// TODO: Implement
	// Hint: Track open/close count, add '(' if open < n, add ')' if close < open
	return nil
}

// SolveSudoku solves a 9x9 Sudoku board (modify in place)
// LeetCode #37
func SolveSudoku(board [][]byte) bool {
	// TODO: Implement
	// Hint: For each empty cell, try 1-9, check valid, backtrack if stuck
	return false
}

// NQueens solves n-queens problem
// Place n queens on n×n board so no two attack each other
// LeetCode #51
func NQueens(n int) [][]string {
	// TODO: Implement
	// Hint: Place queen row by row, check column and diagonals
	return nil
}

// WordSearch - does word exist in 2D board?
// Can move up/down/left/right, can't reuse cells
// LeetCode #79
func WordSearch(board [][]byte, word string) bool {
	// TODO: Implement
	// Hint: DFS from each cell, mark visited, backtrack
	return false
}

// PalindromePartitioning - partition string into palindromes
// LeetCode #131
func PalindromePartitioning(s string) [][]string {
	// TODO: Implement
	// Hint: Try every partition point, check if palindrome, backtrack
	return nil
}

// CombinationSum - find all combinations that sum to target
// Can reuse numbers
// LeetCode #39
func CombinationSum(candidates []int, target int) [][]int {
	// TODO: Implement
	// Hint: For each number, take it (and can take again) or skip it
	return nil
}

// Other backtracking problems to practice:
// - Combination Sum II (no reuse)
// - Combination Sum III (limited count)
// - Restore IP Addresses
// - Word Search II (with Trie)
// - Expression Add Operators
