package dancinglinks

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"
)

var (
	classic = New(7, [][]int{
		[]int{2, 4},
		[]int{0, 3, 6},
		[]int{1, 2, 5},
		[]int{0, 3, 5},
		[]int{1, 6},
		[]int{3, 4, 6},
	})
	classicMatrix = [][]bool{
		[]bool{false, false, true, false, true, false, false},
		[]bool{true, false, false, true, false, false, true},
		[]bool{false, true, true, false, false, true, false},
		[]bool{true, false, false, true, false, true, false},
		[]bool{false, true, false, false, false, false, true},
		[]bool{false, false, false, true, true, false, true},
	}
	classicSolution = [][]int{
		[]int{0, 3, 4},
	}

	classicDuplicates = New(7, [][]int{
		[]int{2, 4},
		[]int{2, 4},
		[]int{0, 3, 6},
		[]int{1, 2, 5},
		[]int{0, 3, 5},
		[]int{0, 3, 5},
		[]int{1, 6},
		[]int{3, 4, 6},
	})
	classicDuplicatesMatrix = [][]bool{
		[]bool{false, false, true, false, true, false, false},
		[]bool{false, false, true, false, true, false, false},
		[]bool{true, false, false, true, false, false, true},
		[]bool{false, true, true, false, false, true, false},
		[]bool{true, false, false, true, false, true, false},
		[]bool{true, false, false, true, false, true, false},
		[]bool{false, true, false, false, false, false, true},
		[]bool{false, false, false, true, true, false, true},
	}
	classicDuplicatesSolution = [][]int{
		[]int{0, 4, 6},
		[]int{0, 5, 6},
		[]int{1, 4, 6},
		[]int{1, 5, 6},
	}

	impossible = New(2, [][]int{
		[]int{1},
	})
	impossibleMatrix = [][]bool{
		[]bool{false, true},
	}
	impossibleSolution = [][]int{}

	trivial         = New(0, [][]int{})
	trivialMatrix   = [][]bool{}
	trivialSolution = [][]int{
		[]int{},
	}
)

func sprintMatrix(matrix [][]bool) string {
	b := &strings.Builder{}
	for i, row := range matrix {
		b.WriteString(fmt.Sprintf("%d:", i))
		for _, cell := range row {
			b.WriteByte(' ')
			if cell {
				b.WriteByte('1')
			} else {
				b.WriteByte('0')
			}
		}
		b.WriteByte('\n')
	}
	return b.String()
}

func testToMatrix(t *testing.T, dl DancingLinks, expected [][]bool) {
	mat := dl.ToMatrix()
	if !reflect.DeepEqual(mat, expected) {
		t.Errorf(
			"matrix mismatch:\nshould be\n%s\ngot\n%s",
			sprintMatrix(expected), sprintMatrix(mat),
		)
	}
}

func TestToMatrix(t *testing.T) {
	testToMatrix(t, classic, classicMatrix)
	testToMatrix(t, classicDuplicates, classicDuplicatesMatrix)
	testToMatrix(t, impossible, impossibleMatrix)
	testToMatrix(t, trivial, trivialMatrix)
}

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

func testExample(t *testing.T, dl DancingLinks, expected [][]int) {
	covers := dl.AllSolutions()
	sortSequences(covers)

	if !reflect.DeepEqual(covers, expected) {
		t.Errorf("\nshould be:\n  %v\nsolver returned:\n  %v", expected, covers)
	}
}

func TestExamples(t *testing.T) {
	testExample(t, classic, classicSolution)
	testExample(t, classicDuplicates, classicDuplicatesSolution)
	testExample(t, impossible, impossibleSolution)
	testExample(t, trivial, trivialSolution)
}

func TestYieldBreak(t *testing.T) {
	count := 0
	classicDuplicates.GenerateSolutions(func([]int) bool {
		count++
		return count < 2
	})

	if count != 2 {
		t.Errorf("short-circuit failed: should run twice, but ran %d times", count)
	}
}
