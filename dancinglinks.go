package main

import (
	"fmt"
)

type itemNode struct {
	index int
	left  *itemNode
	right *itemNode
	head  *entryNode
}

type entryNode struct {
	itemIndex   int
	optionIndex int
	up          *entryNode
	down        *entryNode
}

type DancingLinks struct {
	// A fixed slice of pointers to item nodes.
	items    []*itemNode

	// Slice mapping index of each option to a slice of entries
	// corresponding to that option.
	entries  [][]*entryNode

	// Blank "anchor" item node, whose `right` points to the first
	// item to be covered.
	itemHead *itemNode
}

func New(itemCount int, options [][]int) DancingLinks {
	dl := DancingLinks{
		items: make([]*itemNode, itemCount),
		entries: make([][]*entryNode, len(options)),
		itemHead: &itemNode{index: -1},
	}

	// Construct item list.
	lastItem := dl.itemHead
	for index := range dl.items {
		newItem := &itemNode{
			index: index,
			left:  lastItem,
			head:  &entryNode{optionIndex: -1},
		}

		// Add item to item slice.
		dl.items[index] = newItem

		// Append to linked list.
		lastItem.right = newItem
		lastItem = newItem
	}

	// Make linked list cyclic to reduce edge cases.
	lastItem.right = dl.itemHead
	dl.itemHead.left = lastItem

	// Keep track of bottom-most node for each column (item).
	lastEntries := make([]*entryNode, itemCount)
	for itemIndex, item := range dl.items {
		lastEntries[itemIndex] = item.head
	}

	// Create and append entry nodes.
	for optionIndex, optionItems := range options {
		for _, itemIndex := range optionItems {
			newEntry := &entryNode{
				itemIndex:   itemIndex,
				optionIndex: optionIndex,
				up:          lastEntries[itemIndex],
			}

			// Add entry to corresponding row (option) record.
			dl.entries[optionIndex] = append(dl.entries[optionIndex], newEntry)

			// Append to column-specific linked list.
			lastEntries[itemIndex].down = newEntry
			lastEntries[itemIndex] = newEntry
		}
	}

	// Make column lists cyclic to reduce edge cases.
	for index, item := range dl.items {
		lastEntries[index].down = item.head
		item.head.up = lastEntries[index]
	}

	return dl
}

// Solves an exact cover problem with Donald Knuth's dancing links
// algorithm.  A solution to an exact cover problem is a subset (the
// cover) of options, such that every item in items is contained in
// exactly one option in the cover.
//
// In the matrix representation of the exact cover problem, each item
// in items corresponds to a column in the matrix, and each option in
// options corresponds to a row in the matrix.

func (dl DancingLinks) CollectSolutions() [][]int {
	covers := [][]int{}
	dl.GenerateSolutions(func(cover []int) {
		covers = append(covers, cover)
	})
	return covers
}


func (dl DancingLinks) GenerateSolutions(yield func([]int)) {
	// First item to cover.
	first := dl.itemHead.right

	// Nothing left to cover!
	if first == dl.itemHead {
		yield([]int{})
		return
	}

	// Consider each option that covers the first item.
	for candidate := first.head.down; candidate != first.head; candidate = candidate.down {

		// Keep track of deleted options so that (1) we don't do redundant
		// deletes, which break things, and (2) we can un-delete them in
		// reverse order.  The slice stores indices of deleted options in
		// the order they are deleted.
		deleted := []int{}

		// Retrieve all entries covered by the selected option.
		entries := dl.entries[candidate.optionIndex]

		// Delete each covered item.
		for _, covered := range entries {
			item := dl.items[covered.itemIndex]

			// Delete covered item from linked list.
			item.left.right = item.right
			item.right.left = item.left

			// Delete all options that cover the same item, since we can
			// only cover each item once.
			for conflict := item.head.down; conflict != item.head; conflict = conflict.down {
				// We can only delete nodes once; trying to re-delete may
				// break things.  So if we've already deleted something, don't
				// try delete it again.
				if sliceContains(deleted, conflict.optionIndex) {
					continue
				}

				// Record deleted option.
				deleted = append(deleted, conflict.optionIndex)

				// To delete an option, we go through and delete each entry in
				// the option.
				for _, entry := range dl.entries[conflict.optionIndex] {
					entry.up.down = entry.down
					entry.down.up = entry.up
				}
			}
		}

		// Recursive call.
		dl.GenerateSolutions(func(subcover []int) {
			yield(append(subcover, candidate.optionIndex))
		})

		// Uncover items in reverse order.
		for i := range entries {
			// We deleted the items left to right (increasing index), so we
			// uncover the items right to left (decreasing index).
			entry := entries[len(entries)-1-i]
			item := dl.items[entry.itemIndex]

			// Uncover item.
			item.left.right = item
			item.right.left = item
		}

		// Restore conflicting options in reverse order.
		for i := range deleted {
			// Retrieve index of deleted option, in reverse order.
			optionIndex := deleted[len(deleted)-1-i]

			// To restore the option, we restore each entry in the option.
			for _, entry := range dl.entries[optionIndex] {
				entry.up.down = entry
				entry.down.up = entry
			}
		}
	}
}

func sliceContains(slice []int, element int) bool {
	for _, e := range slice {
		if e == element {
			return true
		}
	}
	return false
}

func main() {

	dl := New(7, [][]int{
		[]int{2, 4},
		[]int{2, 4},
		[]int{0, 3, 4, 5, 6},
		[]int{1, 2, 5},
		[]int{0, 3, 5},
		[]int{1, 6},
		[]int{3, 4, 6},
	})

	fmt.Println(dl.CollectSolutions())

}
