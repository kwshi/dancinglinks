package dancinglinks

import (
	"reflect"
	"sort"
	"testing"
)

func sortSequences(sequences [][]int) {
	// First, sort each individual cover.
	for _, seq := range sequences {
		sort.Ints(seq)
	}

	sort.Slice(sequences, func(i, j int) bool {
		otherSeq := sequences[j]

		// Lexicographically compare covers.
		for k, value := range sequences[i] {
			// If we've run out of things to compare, then the shorter
			// sequence (the other sequence) is less.
			if k == len(otherSeq) {
				return false
			}

			// Compare leading sequences.  The sequence with lower leading
			// entries is lower.  If leading entries are the same, move on
			// to the next entry.
			switch {
			case value < otherSeq[k]:
				return true
			case value > otherSeq[k]:
				return false
			}
		}

		// We've run out of things to compare; the shorter list (the
		// current sequence) is less.
		return true
	})
}

func check(t *testing.T, dl DancingLinks, expected [][]int) {
	covers := dl.CollectSolutions()
	sortSequences(covers)

	if !reflect.DeepEqual(covers, expected) {
		t.Errorf("\nshould be:\n  %v\nsolver returned:\n  %v", expected, covers)
	}
}

func TestClassic(t *testing.T) {
	check(
		t,
		New(7, [][]int{
			[]int{2, 4},
			[]int{0, 3, 4, 5, 6},
			[]int{1, 2, 5},
			[]int{0, 3, 5},
			[]int{1, 6},
			[]int{3, 4, 6},
		}),
		[][]int{
			[]int{0, 3, 4},
		},
	)
}

func TestClassicDuplicates(t *testing.T) {
	check(
		t,
		New(7, [][]int{
			[]int{2, 4},
			[]int{2, 4},
			[]int{0, 3, 4, 5, 6},
			[]int{1, 2, 5},
			[]int{0, 3, 5},
			[]int{0, 3, 5},
			[]int{1, 6},
			[]int{3, 4, 6},
		}),
		[][]int{
			[]int{0, 4, 6},
			[]int{0, 5, 6},
			[]int{1, 4, 6},
			[]int{1, 5, 6},
		},
	)
}

func TestImpossible(t *testing.T) {
	check(
		t,
		New(2, [][]int{
			[]int{1},
		}),
		[][]int{},
	)
}

func TestTrivial(t *testing.T) {
	check(
		t,
		New(0, [][]int{}),
		[][]int{
			[]int{},
		},
	)
}
