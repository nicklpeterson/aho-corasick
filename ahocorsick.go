package main

import (
	"fmt"
	"io"
	"os"
	"sync"
)

const BufferSize = 256

type Processor func(word string, endIndex int, filename string)

type Automata struct {
	trie *trie
}

func NewAutomata(targets []string) Automata {
	automata := Automata{
		newTrie(targets),
	}
	return automata
}

func (a *Automata) SimpleStringSearch(text string) map[string]int {
	wordCounts := make(map[string]int)
	processor := func(word string, endIndex int, filename string) {
		wordCounts[word] += 1
	}
	a.search(processor, text, 0, "")
	return wordCounts
}

func (a *Automata) FileSearch(processor Processor, filePath string, wg *sync.WaitGroup) {
	defer wg.Done()
	var fileProcessingWg sync.WaitGroup
	channel := make(chan []byte, 10)
	fileProcessingWg.Add(2)
	go producer(&channel, filePath, &fileProcessingWg)
	go consumer(&channel, &fileProcessingWg, a, processor, filePath)
	fileProcessingWg.Wait()
}

func (a *Automata) MultipleFileSearch(processor Processor, filenames []string) {
	var wg sync.WaitGroup
	for _, path := range filenames {
		wg.Add(1)
		go a.FileSearch(processor, path, &wg)
	}
	wg.Wait()
}

func (a *Automata) StringSearch(processor Processor, text string) {
	a.search(processor, text, 0, "")
}

func (a *Automata) search(processor Processor, text string, initialState int, filename string) (state int) {
	state = initialState
	for i, c := range text {
		state = a.trie.transition(state, c)
		exitPtr := a.trie.getExitLink(state)
		for exitPtr != 0 {
			processor(a.trie.nodes[exitPtr].word, i, filename)
			exitPtr = a.trie.getExitLink(exitPtr)
		}
		if a.trie.nodes[state].leaf {
			processor(a.trie.nodes[state].word, i, filename)
		}
	}
	return
}

func consumer(channel *chan []byte, wg *sync.WaitGroup, a *Automata, processor Processor, filename string) {
	defer wg.Done()
	state := 0
	buffer, ok := <- *channel
	for ok {
		text := string(buffer)
		state = a.search(processor, text, state, filename)
		buffer, ok = <- *channel
	}
}

func producer(channel *chan []byte, filename string, wg *sync.WaitGroup)  {
	defer wg.Done()

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)

	buffer := make([]byte, BufferSize)

	bytesRead, err := file.Read(buffer)

	for err == nil{
		*channel <- buffer[:bytesRead]
		bytesRead, err = file.Read(buffer)
	}

	if err != io.EOF {
		fmt.Println(err)
	}

	close(*channel)
}
