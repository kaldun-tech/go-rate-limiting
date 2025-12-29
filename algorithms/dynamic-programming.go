package algorithms

// Dynamic Programming - optimize recursive solutions with memoization or tabulation
// Classic pattern: break problem into subproblems, store results to avoid recomputation

// Fibonacci returns the nth Fibonacci number
// Classic DP intro problem
// Naive recursion: O(2^n), DP: O(n)
func Fibonacci(n int) int {
	// TODO: Implement with DP
	// Can use array/map for memoization or bottom-up tabulation
	return 0
}

// ClimbStairs - you can climb 1 or 2 steps at a time
// How many distinct ways to climb n steps?
// LeetCode #70
func ClimbStairs(n int) int {
	// TODO: Implement
	// Hint: dp[i] = dp[i-1] + dp[i-2]
	return 0
}

// CoinChange - minimum coins needed to make amount
// Given coins of different denominations, find minimum number to make target
// LeetCode #322
func CoinChange(coins []int, amount int) int {
	// TODO: Implement
	// Hint: dp[i] = min(dp[i], dp[i-coin] + 1)
	return -1
}

// LongestIncreasingSubsequence finds length of LIS
// LeetCode #300
// Time: O(n²) DP or O(n log n) with binary search
func LongestIncreasingSubsequence(nums []int) int {
	// TODO: Implement
	// Hint: dp[i] = max length ending at i
	return 0
}

// Knapsack01 solves 0/1 knapsack problem
// Given weights, values, and capacity, maximize value
// Classic DP problem
func Knapsack01(weights, values []int, capacity int) int {
	// TODO: Implement
	// Hint: dp[i][w] = max value using first i items with capacity w
	return 0
}

// LongestCommonSubsequence finds length of LCS between two strings
// LeetCode #1143
func LongestCommonSubsequence(text1, text2 string) int {
	// TODO: Implement
	// Hint: dp[i][j] based on text1[i-1] == text2[j-1]
	return 0
}

// EditDistance - minimum operations to convert word1 to word2
// Operations: insert, delete, replace
// LeetCode #72 (Levenshtein distance)
func EditDistance(word1, word2 string) int {
	// TODO: Implement
	// Hint: dp[i][j] = edit distance for word1[:i] and word2[:j]
	return 0
}

// WordBreak - can string be segmented into dictionary words?
// LeetCode #139
func WordBreak(s string, wordDict []string) bool {
	// TODO: Implement
	// Hint: dp[i] = can s[:i] be segmented?
	return false
}

// HouseRobber - rob houses to maximize money, can't rob adjacent
// LeetCode #198
func HouseRobber(nums []int) int {
	// TODO: Implement
	// Hint: dp[i] = max(dp[i-1], dp[i-2] + nums[i])
	return 0
}

// MaxSubarraySum - maximum sum of contiguous subarray
// Kadane's Algorithm - technically greedy but good DP intro
// LeetCode #53
func MaxSubarraySum(nums []int) int {
	// TODO: Implement
	// Hint: maxEndingHere = max(nums[i], maxEndingHere + nums[i])
	return 0
}

// Other important DP problems to practice:
// - Unique Paths (grid traversal)
// - Jump Game / Jump Game II
// - Decode Ways
// - Partition Equal Subset Sum
// - Longest Palindromic Substring
// - Regular Expression Matching
// - Burst Balloons
// - Distinct Subsequences
