package main

import (
	"fmt"
	"os"

	io "github.com/faideww/aoc-2024/lib"
)

type TrieNode struct {
	value      int
	isTerminal bool
	children   map[byte]*TrieNode
}

func NewTrieNode() *TrieNode {
	return &TrieNode{
		children: make(map[byte]*TrieNode),
	}
}

func main() {
	input := io.ReadInputFile(os.Args[1])

	trie, words := parseInput(input)

	sum := 0
	for _, word := range words {
		isPossible := isPatternPossible(trie, word)
		if isPossible {
			sum++
		}
	}

	fmt.Printf("possibles: %d\n", sum)

	sum2 := 0
	for _, word := range words {
		sum2 += countPossiblePatterns(trie, word)
	}

	fmt.Printf("all possibles: %d\n", sum2)
}

func parseInput(input string) (*TrieNode, []string) {
	components := io.TrimAndSplitBy(input, "\n\n")

	stripes := io.TrimAndSplitBy(components[0], ", ")
	words := io.TrimAndSplit(components[1])

	// build trie
	root := NewTrieNode()

	for _, stripe := range stripes {
		current := root
		for i := 0; i < len(stripe); i++ {
			if _, ok := current.children[stripe[i]]; !ok {
				current.children[stripe[i]] = NewTrieNode()
			}
			current = current.children[stripe[i]]
		}
		current.value = len(stripe)
		current.isTerminal = true
	}

	return root, words
}

func isPatternPossible(trieRoot *TrieNode, pattern string) bool {
	// modified search: we progressively walk the trie to find all
	// matching substrings. when we do, for each substring we move
	// the cursor forward by the length of the match, until:
	// 1. we find a match that ends exactly at the pattern's length (match found)
	// 2. we fail to find a match at any point
	// 3. we reach the end of the string and are not at a terminal node

	cache := make(map[string]bool)

	return isSubpatternPossible(trieRoot, pattern, cache)
}

func countPossiblePatterns(trieRoot *TrieNode, pattern string) int {
	cache := make(map[string]int)
	return countSubpatterns(trieRoot, pattern, cache)
}

func isSubpatternPossible(trieRoot *TrieNode, pattern string, cache map[string]bool) bool {
	if val, ok := cache[pattern]; ok {
		return val
	}

	if len(pattern) == 0 {
		return true
	}

	substrings := findAllMatches(trieRoot, pattern)
	canMatch := false
	for _, substring := range substrings {
		if isSubpatternPossible(trieRoot, pattern[substring:], cache) {
			canMatch = true
			break
		}
	}

	cache[pattern] = canMatch
	return canMatch
}

func countSubpatterns(trieRoot *TrieNode, pattern string, cache map[string]int) int {
	if val, ok := cache[pattern]; ok {
		return val
	}

	if len(pattern) == 0 {
		return 1
	}

	substrings := findAllMatches(trieRoot, pattern)
	matches := 0
	for _, substring := range substrings {
		result := countSubpatterns(trieRoot, pattern[substring:], cache)
		if result > 0 {
			matches += result
		}
	}

	cache[pattern] = matches
	return matches
}

func findLongestMatch(trieRoot *TrieNode, pattern string) int {
	current := trieRoot
	var i int
	for i = 0; i < len(pattern); i++ {
		if _, ok := current.children[pattern[i]]; !ok {
			break
		}
		current = current.children[pattern[i]]
	}

	if current.isTerminal {
		fmt.Printf("match: %s (value: %d)\n", pattern[:current.value], current.value)
		return current.value
	} else {
		fmt.Printf("no match found; reached %s (%d)\n", pattern[:i], i)
		return -1
	}
}

func findAllMatches(trieRoot *TrieNode, pattern string) []int {
	current := trieRoot
	is := []int{}
	var i int
	for i = 0; i < len(pattern); i++ {
		if _, ok := current.children[pattern[i]]; !ok {
			break
		}
		current = current.children[pattern[i]]
		if current.isTerminal {
			is = append(is, i+1)
		}
	}

	return is
}

func printTrie(t *TrieNode, leftPad int) {

	for char, c := range t.children {
		for i := 0; i < leftPad; i++ {
			fmt.Printf(" ")
		}
		fmt.Printf("-%c", char)
		if c.isTerminal {
			fmt.Printf("*")
		}
		fmt.Println()
		printTrie(c, leftPad+1)
	}
}
