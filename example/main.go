package main

import (
	"fmt"

	"github.com/kgwinnup/ac-automata/aca"
)

func main() {

	// input stream used with the matching patterns
	input := "pinpiiistringingting"

	// the finite automata requires a list of patterns, each patter
	// must also be a sequence that is used for matching
	initPatterns := []string{
		"i",
		"in",
		"tin",
		"pin",
		"string",
	}

	// transform the patterns into slices of slices
	patterns := make([][]rune, 0)
	for _, pattern := range initPatterns {
		temp := make([]rune, 0)
		for _, r := range pattern {
			temp = append(temp, r)
		}
		patterns = append(patterns, temp)
	}

	// create the automta
	machine := aca.New(patterns)

	// defnie a String function, this is only used for the Dot output,
	// it will label the arrows so keep them short
	toString := func(r rune) string {
		return string(r)
	}

	// print the Dot representation of the automta
	fmt.Println(machine.ToDot(patterns, toString))

	// count the pattern matches
	matches := machine.Counts(patterns, []rune(input))

	for i, pattern := range initPatterns {
		fmt.Printf("%-10v %v\n", pattern, matches[i])
	}
}
