package main

import (
	"fmt"
	"github.com/kwshi/dancinglinks"
)

type sudokuEntry struct {
	row, column, value int
}

func block(row, column int) int {
	return (row/3)*3 + (column / 3)
}

func solve(board [][]int) {
	options := make([][]int, 9*9*9)
	sudokuEntries := make([]sudokuEntry, 9*9*9)

	for row := 0; row < 9; row++ {
		for column := 0; column < 9; column++ {
			for value := 0; value < 9; value++ {
				entry := sudokuEntry{row, column, value}

				option := []int{
					0*9*9 + 9*value + row,
					1*9*9 + 9*value + column,
					2*9*9 + 9*value + block(row, column),
					3*9*9 + 9*row + column,
				}

				options[9*9*row + 9*column + value] = option
				sudokuEntries[9*9*row + 9*column + value] = entry
			}
		}
	}

	dl := dancinglinks.New(4*9*9, options)

	for row := 0; row < 9; row++ {
		for column := 0; column < 9; column++ {
			if board[row][column] != 0 {
				dl.ForceOptions(9*9*row + 9*column + board[row][column] - 1)
			}
		}
	}

	cover := dl.AnySolution()

	for _, option := range cover {
		entry := sudokuEntries[option]
		fmt.Printf("row %d, column %d: value %d\n", entry.row+1, entry.column+1, entry.value+1)
		board[entry.row][entry.column] = entry.value + 1
	}

	for _, row := range board {
		fmt.Println(row)
	}

	count := 0
	dl.GenerateSolutions(func([]int) bool {
		count++
		fmt.Printf("\r%d solutions found ", count)
		return true
	})
	fmt.Println()
}

func main() {
	solve([][]int{
		[]int{6, 4, 0, 0, 3, 0, 0, 0, 7},
		[]int{5, 0, 1, 0, 7, 0, 9, 0, 0},
		[]int{0, 0, 0, 0, 0, 0, 0, 1, 0},
		[]int{0, 0, 4, 9, 0, 8, 0, 6, 0},
		[]int{0, 8, 0, 0, 0, 3, 0, 2, 0},
		[]int{0, 0, 0, 4, 0, 0, 0, 0, 0},
		[]int{4, 0, 0, 1, 5, 7, 0, 3, 0},
		[]int{2, 0, 8, 3, 0, 0, 0, 4, 0},
		[]int{7, 5, 0, 0, 0, 0, 0, 9, 6},
	})
}
