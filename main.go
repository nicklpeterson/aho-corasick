package main

import (
	"strconv"
)

func main() {
	automata := NewAutomata([]string{"abcde", "acb", "bc"})
	counts := automata.SimpleStringSearch("abcde")
	for key, val := range counts {
		println(key + " --> " + strconv.Itoa(val))
	}
	automata.StringSearch(processor, "abcde")
}

func processor(word string, endIndex int) {
	start := strconv.Itoa(endIndex - len(word) + 1)
	end := strconv.Itoa(endIndex)
	println(word + " at {" + start + "," + end + "}")
}
