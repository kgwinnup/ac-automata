package aca

import (
	"testing"
)

func TestIndexes(t *testing.T) {
	input := "foobarfoobazbarr"

	initPatterns := []string{
		"o",
	}

	patterns := make([][]rune, 0)
	for _, pattern := range initPatterns {
		temp := make([]rune, 0)
		for _, r := range pattern {
			temp = append(temp, r)
		}
		patterns = append(patterns, temp)
	}

	machine := New(patterns)

	indexes := machine.Indexes(patterns, []rune(input))

	if len(indexes[0]) != 4 {
		t.Fatal(len(indexes[0]))
		t.Fatal("incorrect number of matches")
	}

	if indexes[0][0] != 1 {
		t.Fatal("wront index")
	}
}

func TestWithStrings(t *testing.T) {

	// input stream used with the matching patterns
	input := "foobarfoobazbarr"

	// the finite automata requires a list of patterns, each patter
	// must also be a sequence that is used for matching
	initPatterns := []string{
		"foobar",
		"oobar",
		"obar",
		"bar",
		"ar",
		"r",
	}

	patterns := make([][]rune, 0)
	for _, pattern := range initPatterns {
		temp := make([]rune, 0)
		for _, r := range pattern {
			temp = append(temp, r)
		}
		patterns = append(patterns, temp)
	}

	machine := New(patterns)

	// initial starting parameters
	cur := 0
	var xs []int
	matches := make([]int, len(initPatterns))

	// loop over each atom in the input stream to "crank" the machine
	for _, atom := range input {
		// machine returns a list of matches and the current index
		// state for the next atom's processing location
		xs, cur = machine.Next(cur, patterns, atom)

		// loop over all matches and increment the counts
		for _, x := range xs {
			matches[x]++
		}
	}

	for i, _ := range initPatterns {
		switch i {
		case 0:
			if matches[0] != 1 {
				t.Error("invalid matches for foobar")
			}
		case 1:
			if matches[1] != 1 {
				t.Error("invalid matches for oobar")
			}
		case 2:
			if matches[2] != 1 {
				t.Error("invalid matches for obar")
			}
		case 3:
			if matches[3] != 2 {
				t.Error("invalid matches for bar")
			}
		case 4:
			if matches[4] != 2 {
				t.Error("invalid matches for ar")
			}
		case 5:
			if matches[5] != 3 {
				t.Error("invalid matches for r")
			}

		}
	}

}
