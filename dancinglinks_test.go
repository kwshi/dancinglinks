package dancinglinks

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"
)

type example struct {
	itemCount int
	options [][]int
	matrix [][]bool
	solution [][]int
}

var (
	classic = example{
		itemCount: 7,
		options: [][]int{
			[]int{2, 4},
			[]int{0, 3, 6},
			[]int{1, 2, 5},
			[]int{0, 3, 5},
			[]int{1, 6},
			[]int{3, 4, 6},
		},
		matrix: [][]bool{
			[]bool{false, false, true, false, true, false, false},
			[]bool{true, false, false, true, false, false, true},
			[]bool{false, true, true, false, false, true, false},
			[]bool{true, false, false, true, false, true, false},
			[]bool{false, true, false, false, false, false, true},
			[]bool{false, false, false, true, true, false, true},
		},
		solution: [][]int{
			[]int{0, 3, 4},
		},
	}

	classicDuplicates = example{
		itemCount: 7,
		options: [][]int{
			[]int{2, 4},
			[]int{2, 4},
			[]int{0, 3, 6},
			[]int{1, 2, 5},
			[]int{0, 3, 5},
			[]int{0, 3, 5},
			[]int{1, 6},
			[]int{3, 4, 6},
		},
		matrix: [][]bool{
			[]bool{false, false, true, false, true, false, false},
			[]bool{false, false, true, false, true, false, false},
			[]bool{true, false, false, true, false, false, true},
			[]bool{false, true, true, false, false, true, false},
			[]bool{true, false, false, true, false, true, false},
			[]bool{true, false, false, true, false, true, false},
			[]bool{false, true, false, false, false, false, true},
			[]bool{false, false, false, true, true, false, true},
		},
		solution: [][]int{
			[]int{0, 4, 6},
			[]int{0, 5, 6},
			[]int{1, 4, 6},
			[]int{1, 5, 6},
		},
	}

	impossible = example{
		itemCount: 3,
		options: [][]int{
			[]int{0, 1},
			[]int{1, 2},
		},
		matrix: [][]bool{
			[]bool{true, true, false},
			[]bool{false, true, true},
		},
		solution: [][]int{},
	}

	trivial = example{
		itemCount: 0,
		options: [][]int{},
		matrix: [][]bool{},
		solution: [][]int{
			[]int{},
		},
	}
)

func (e example) toDancingLinks() *DancingLinks {
	return New(e.itemCount, e.options)
}

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

func testToMatrix(t *testing.T, e example) {
	mat := e.toDancingLinks().ToMatrix()
	if !reflect.DeepEqual(mat, e.matrix) {
		t.Errorf(
			"matrix mismatch:\nshould be\n%s\ngot\n%s",
			sprintMatrix(e.matrix), sprintMatrix(mat),
		)
	}
}

func TestToMatrix(t *testing.T) {
	testToMatrix(t, classic)
	testToMatrix(t, classicDuplicates)
	testToMatrix(t, impossible)
	testToMatrix(t, trivial)
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

func testExample(t *testing.T, result [][]int, correct [][]int) {
	sortSequences(result)

	if !reflect.DeepEqual(result, correct) {
		t.Errorf("\nshould be:\n  %v\nsolver returned:\n  %v", correct, result)
	}
}

func TestExamples(t *testing.T) {
	for _, e := range []example{
		classic,
		classicDuplicates,
		impossible,
		trivial,
	} {
		testExample(t, e.toDancingLinks().AllSolutions(), e.solution)
	}
}

func TestYieldBreak(t *testing.T) {
	count := 0
	classicDuplicates.toDancingLinks().GenerateSolutions(func([]int) bool {
		count++
		return count < 2
	})

	if count != 2 {
		t.Errorf("short-circuit failed: should run twice, but ran %d times", count)
	}
}

func TestForceOptions(t *testing.T) {
	dl := classicDuplicates.toDancingLinks()
	dl.ForceOptions(0)
	testExample(t, dl.AllSolutions(), [][]int{
		[]int{4, 6},
		[]int{5, 6},
	})

	dl = classicDuplicates.toDancingLinks()
	dl.ForceOptions(0, 1)
	testExample(t, dl.AllSolutions(), [][]int{
		[]int{4, 6},
		[]int{5, 6},
	})

	dl = classicDuplicates.toDancingLinks()
	dl.ForceOptions(4)
	testExample(t, dl.AllSolutions(), [][]int{
		[]int{0, 6},
		[]int{1, 6},
	})

	dl = classicDuplicates.toDancingLinks()
	dl.ForceOptions(2)
	testExample(t, dl.AllSolutions(), [][]int{})
}
