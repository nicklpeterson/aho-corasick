package main

type Processor func(word string, endIndex int)

type Automata struct {
	trie *trie
}

func NewAutomata(targets []string) Automata {
	automata := Automata{
		newTrie(targets),
	}
	return automata
}

func (s *Automata) SimpleStringSearch(text string) map[string]int {
	wordCounts := make(map[string]int)
	processor := func(word string, endIndex int) {
		wordCounts[word] += 1
	}
	s.StringSearch(processor, text)
	return wordCounts
}

func (s *Automata) StringSearch(processor Processor, text string) {
	state := 0
	for i, c := range text {
		state = s.trie.transition(state, c)
		exitPtr := s.trie.getExitLink(state)
		for exitPtr != 0 {
			processor(s.trie.nodes[exitPtr].word, i)
			exitPtr = s.trie.getExitLink(exitPtr)
		}
		if s.trie.nodes[state].leaf {
			processor(s.trie.nodes[state].word, i)
		}
	}
}
