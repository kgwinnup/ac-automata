package aca

import (
	"bytes"
	"fmt"
)

type ACAutomata[T comparable] struct {
	nodes []*Node[T]
}

// Node is the data that makes up the automata trie
type Node[T comparable] struct {
	id          int
	data        T
	children    map[T]int
	fail        int
	alternative int
	// index in the nodes slice of the match
	match int
	// length of the match, used to return a location of the substring
	matchLen int
}

// New creates a new state machine with the provided
// patterns. Patterns parameter is a lsit of patterns, where each
// individual pattern contains a sequence to match on.
func New[T comparable](patterns [][]T) *ACAutomata[T] {
	nodes := make([]*Node[T], 0)

	root := &Node[T]{
		id:          0,
		children:    make(map[T]int),
		fail:        -1,
		alternative: -1,
		match:       -1,
		matchLen:    -1,
	}

	nodes = append(nodes, root)

	cur := root
	ids := 1

	// for each pattern, build a branch for each atom in the pattern
	// starting at the root. If path already exists, traverse down it
	// until a new node can be added
	for i, pattern := range patterns {
		cur = root

		for j, atom := range pattern {

			if id, ok := cur.children[atom]; ok {
				// if this is the last node, set it as a match
				if j == len(pattern)-1 {
					nodes[id].match = i
					nodes[id].matchLen = j
				}

				cur = nodes[id]
				continue
			}

			node := &Node[T]{
				id:       ids,
				data:     atom,
				children: make(map[T]int),
				fail:     -1,
				match:    -1,
			}

			nodes = append(nodes, node)
			ids++

			// if this is the last node, set it as a match
			if j == len(pattern)-1 {
				node.match = i
				node.matchLen = j
			}

			// add this node id to the parents children
			cur.children[atom] = node.id
			cur = node
		}
	}

	root.fail = root.id

	failQueue := make([]int, 0)
	id := 0

	// set the first node of each branch fail point back to root
	for _, childId := range root.children {
		// always fail first level to root
		nodes[childId].fail = 0
		failQueue = append(failQueue, childId)
	}

	for {

		// if the initial branch nodes are all consumed
		if len(failQueue) == 0 {
			break
		}

		// pop one item off the queue
		id, failQueue = failQueue[0], failQueue[1:]
		cur = nodes[id]

		// for each child in this branch find the largest suffix
		for _, childId := range cur.children {
			child := nodes[childId]
			temp := nodes[cur.fail]

			for {
				_, ok := temp.children[child.data]

				// we're at the largest suffix so far
				if !ok && temp != root {
					temp = nodes[temp.fail]
				}

				if temp == root || ok {
					break
				}

			}

			if node, ok := temp.children[child.data]; ok {
				child.fail = node // proper suffix
			} else {
				child.fail = root.id // no suffix
			}

			// add this node to the queue for processing
			failQueue = append(failQueue, child.id)
		}

		if nodes[cur.fail].match >= 0 {
			cur.alternative = cur.fail
		} else {
			cur.alternative = nodes[cur.fail].alternative
		}
	}

	return &ACAutomata[T]{nodes: nodes}
}

func (a *ACAutomata[T]) ToDot(patterns [][]T, toString func(t T) string) string {

	buf := new(bytes.Buffer)

	fmt.Fprintf(buf, "digraph finite_state_machine {\n")
	fmt.Fprintf(buf, "  fontname=\"Helvetica,Arial,sans-serif\"\n")
	fmt.Fprintf(buf, "  node [fontname=\"Helvetica,Arial,sans-serif\"]\n")
	fmt.Fprintf(buf, "  edge [fontname=\"Helvetica,Arial,sans-serif\"]\n")
	fmt.Fprintf(buf, "  rankdir=LR;\n")

	fmt.Fprintf(buf, "  node [shape = doublecircle];")

	for _, node := range a.nodes {
		if node.match >= 0 {
			fmt.Fprintf(buf, " %v", node.id)
		}
	}

	fmt.Fprintf(buf, ";\n")

	fmt.Fprintf(buf, "  node [shape = circle];\n")

	for _, node := range a.nodes {
		if node.fail >= 0 {
			fmt.Fprintf(buf, "  %v -> %v [style = dashed, constraint=false];\n", node.id, a.nodes[node.fail].id)

			if node.alternative >= 0 && a.nodes[node.alternative].match >= 0 {
				fmt.Fprintf(buf, "  %v -> %v [style = dotted, constraint=false];\n", node.id, a.nodes[node.alternative].id)
			}

		}

		for _, child := range node.children {
			fmt.Fprintf(buf, "  %v -> %v [label = \"%v\"];\n", node.id, a.nodes[child].id, toString(a.nodes[child].data))
		}
	}

	fmt.Fprintf(buf, "}\n")

	return buf.String()
}

// Next method is for processing data points one at a time at any
// point in the automata. the `index` parameter defines where to start
// in the automta, the `patterns` are the list of patterns used to build
// the automta, and `c` parameter is used as the input atom to
// transition the automata.
// the return is a tuple consisting of slice of any matching nodes and
// the count for how many times that node matched, and the index of
// the next transition node.
func (a *ACAutomata[T]) Next(index int, patterns [][]T, c T) ([]int, int) {
	node := a.nodes[index]
	matches := make([]int, 0)

Start:
	if id, ok := node.children[c]; ok {
		node = a.nodes[id]

		if node.match >= 0 {
			matches = append(matches, node.match)
		}

		temp := node.alternative
		for {
			if temp < 0 {
				break
			}

			matches = append(matches, a.nodes[temp].match)
			temp = a.nodes[temp].alternative
		}

	} else {
		for {
			_, ok := node.children[c]

			if node.id == 0 {
				break
			}

			if ok {
				break
			}

			node = a.nodes[node.fail]
		}

		if _, ok := node.children[c]; ok {
			goto Start
		}
	}

	return matches, node.id
}

// Indexes will calculate at which indexes the patterns are found. The
// return value is a slice of slices, the top level slice is indexed
// the same as the patterns parameter.
func (a *ACAutomata[T]) Indexes(patterns [][]T, input []T) [][]int {
	node := a.nodes[0]
	indexes := make([][]int, len(patterns))
	for i := 0; i < len(patterns); i++ {
		indexes[i] = make([]int, 0)
	}

	for i := 0; i < len(input); i++ {
		c := input[i]

		if id, ok := node.children[c]; ok {
			node = a.nodes[id]

			if node.match >= 0 {
				indexes[node.match] = append(indexes[node.match], i-node.matchLen)
			}

			// now follow the alternative matches
			temp := node.alternative
			for {
				if temp < 0 {
					break
				}

				indexes[a.nodes[temp].match] = append(indexes[a.nodes[temp].match], i-a.nodes[temp].matchLen)
				temp = a.nodes[temp].alternative
			}

		} else {
			for {
				_, ok := node.children[c]

				if node.id == 0 {
					break
				}

				if ok {
					break
				}

				node = a.nodes[node.fail]
			}

			if _, ok := node.children[c]; ok {
				i--
			}
		}
	}

	return indexes
}

// Counts will return a slice of counts, counting how many times each
// pattern was seen in the input. The indexes of the returned slice
// will match up to the indexes of the patterns parameter.
func (a *ACAutomata[T]) Counts(patterns [][]T, input []T) []int {

	// initial starting parameters
	cur := 0
	var xs []int
	matches := make([]int, len(patterns))

	// loop over each atom in the input stream to "crank" the machine
	for _, atom := range input {
		// machine returns a list of matches and the current index
		// state for the next atom's processing location
		xs, cur = a.Next(cur, patterns, atom)

		// loop over all matches and increment the counts
		for _, x := range xs {
			matches[x]++
		}
	}

	return matches
}
