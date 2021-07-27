package main

import "sync"

type trie struct {
	nodes []node
	mutex sync.Mutex
}

type node struct {
	children   map[rune]int
	leaf       bool
	parent      int
	char       int32
	suffixLink int
	exitLink   int
	transition map[rune]int
	word       string
}

func newTrie(words []string) *trie {
	nodeSlice := make([]node, 1, 1)
	nodeSlice[0] = *newTrieNode(0,0)

	newTrie := trie{
		nodes: nodeSlice,
	}
	for _, word := range words {
		newTrie.addString(word)
	}
	return &newTrie
}

func newTrieNode(ancestor int, char int32) *node {
	trieNode := node{
		children:   make(map[rune]int),
		leaf:       false,
		parent:     ancestor,
		char:       char,
		suffixLink: -1,
		exitLink:   -1,
		transition: make(map[rune]int),
	}
	return &trieNode
}

func (i *trie) addString(word string) {
	state := 0
	for _, c := range word {
		if _, ok := i.nodes[state].children[c]; !ok {
			i.nodes[state].children[c] = len(i.nodes)
			i.nodes = append(i.nodes, *newTrieNode(state, c))
		}
		state = i.nodes[state].children[c]
	}
	i.nodes[state].leaf = true
	i.nodes[state].word = word
}

func (i *trie) getSuffixLink(node int) int {
	if i.nodes[node].suffixLink == -1 {
		if node == 0 || i.nodes[node].parent == 0 {
			i.nodes[node].suffixLink = 0
		} else {
			i.nodes[node].suffixLink = i.transition(
				i.getSuffixLink(i.nodes[node].parent),
				i.nodes[node].char)
		}
	}
	return i.nodes[node].suffixLink
}

func (i *trie) transition(node int, c rune) int {
	if _, tranSet := i.loadFromTransition(node, c); !tranSet {
		if _, childSet := i.nodes[node].children[c]; childSet {
			i.storeToTransition(node, c, i.nodes[node].children[c])
		} else if node == 0 {
			i.storeToTransition(node, c, 0)
		} else {
			i.storeToTransition(node, c, i.transition(i.getSuffixLink(node), c))
		}
	}
	result, _ := i.loadFromTransition(node, c)
	return result
}

func (i *trie) getExitLink(node int) int {
	if i.nodes[node].exitLink == -1 {
		suffixLink := i.getSuffixLink(node)
		if node == 0 || suffixLink == 0 {
			i.nodes[node].exitLink = 0
		} else if i.nodes[suffixLink].leaf {
			i.nodes[node].exitLink = suffixLink
		} else {
			i.nodes[node].exitLink = i.getExitLink(suffixLink)
		}
	}
	return i.nodes[node].exitLink
}

func (i *trie) storeToTransition(node int, c rune, val int) {
	i.mutex.Lock()
	i.nodes[node].transition[c] = val
	i.mutex.Unlock()
}

func (i *trie) loadFromTransition(node int, c rune) (int, bool) {
	i.mutex.Lock()
	val, ok := i.nodes[node].transition[c]
	i.mutex.Unlock()
	return val, ok
}