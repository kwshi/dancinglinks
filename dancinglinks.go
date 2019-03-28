package dancinglinks

//import "fmt"

type Step struct {
	option  int
	choices []int
}


type itemNode struct {
	// Linked list neighbors.
	left  *itemNode
	right *itemNode

	// Blank "anchor" entry node whose `down` points to the
	// first/top-most entry covering this item.
	head *entryNode

	// Number of (remaining) entries that cover the item.  At each
	// iteration the dancing links algorithm chooses the item with the
	// fewest options covering it.
	choices int
}

type entryNode struct {
	// The item covered by this entry.
	item *itemNode

	// The index of the option this entry belongs to.
	option int

	// Linked list neighbors.
	up   *entryNode
	down *entryNode
}

type DancingLinks struct {
	// A list of options.  Each "option" is a list of pointers to the
	// entries provided by that option.
	options [][]*entryNode

	// Blank "anchor" item node, whose `right` points to the
	// first/left-most item to be covered.
	itemHead *itemNode

	// Pre-selected options and associated deleted option indices.
	selected []int
	deleted  []int
}

func New(itemCount int, options [][]int) *DancingLinks {
	dl := &DancingLinks{
		options:  make([][]*entryNode, len(options)),
		itemHead: &itemNode{},
		selected: []int{},
		deleted:  []int{},
	}

	// Construct item list.
	items := make([]*itemNode, itemCount)
	lastItem := dl.itemHead
	for index := range items {
		newItem := &itemNode{
			left: lastItem,
			head: &entryNode{option: -1},
		}

		// Add item to slice.
		items[index] = newItem

		// Append to linked list.
		lastItem.right = newItem
		lastItem = newItem
	}

	// Make linked list cyclic to reduce edge cases.
	lastItem.right = dl.itemHead
	dl.itemHead.left = lastItem

	// Keep track of bottom-most node for each column (item).
	lastEntries := make([]*entryNode, itemCount)
	for itemIndex, item := range items {
		lastEntries[itemIndex] = item.head
	}

	// Create and append entry nodes.
	for option, optionItems := range options {
		for _, itemIndex := range optionItems {
			newEntry := &entryNode{
				item:   items[itemIndex],
				option: option,
				up:     lastEntries[itemIndex],
			}

			newEntry.item.choices++

			// Add entry to corresponding row (option) record.
			dl.options[option] = append(dl.options[option], newEntry)

			// Append to column-specific linked list.
			lastEntries[itemIndex].down = newEntry
			lastEntries[itemIndex] = newEntry
		}
	}

	// Make column lists cyclic to reduce edge cases.
	for index, item := range items {
		lastEntries[index].down = item.head
		item.head.up = lastEntries[index]
	}

	return dl
}

func FromMatrix(matrix [][]bool) *DancingLinks {
	itemCount := 0
	options := make([][]int, len(matrix))

	for i, row := range matrix {
		if len(row) > itemCount {
			itemCount = len(row)
		}

		option := []int{}
		for j, cell := range row {
			if cell {
				option = append(option, j)
			}
		}

		options[i] = option
	}

	return New(itemCount, options)
}

func (dl *DancingLinks) ToMatrix() [][]bool {
	items := map[*itemNode]int{}
	index := 0
	for item := dl.itemHead.right; item != dl.itemHead; item = item.right {
		items[item] = index
		index++
	}

	mat := make([][]bool, len(dl.options))

	for i, option := range dl.options {
		row := make([]bool, len(items))
		for _, entry := range option {
			row[items[entry.item]] = true
		}
		mat[i] = row
	}

	return mat
}

// Solves an exact cover problem with Donald Knuth's dancing links
// algorithm.  A solution to an exact cover problem is a subset (the
// cover) of options, such that every item in items is contained in
// exactly one option in the cover.
//
// In the matrix representation of the exact cover problem, each item
// in items corresponds to a column in the matrix, and each option in
// options corresponds to a row in the matrix.

func (dl *DancingLinks) AllSolutions() [][]Step {
	covers := [][]Step{}
	dl.GenerateSolutions(func(cover []Step) bool {
		covers = append(covers, cover)
		return true
	})
	return covers
}

func (dl *DancingLinks) AnySolution() []Step {
	var solution []Step
	dl.GenerateSolutions(func(cover []Step) bool {
		solution = cover
		return false
	})
	return solution
}

func (dl *DancingLinks) ForceOptions(indices ...int) {
	for _, index := range indices {
		dl.selected = append(dl.selected, index)
		dl.chooseOption(index, &dl.deleted)
	}
}

func (dl *DancingLinks) chooseOption(index int, deleted *[]int) {
	// Keep track of deleted options so that (1) we don't do redundant
	// deletes, which break things, and (2) we can un-delete them in
	// reverse order.  The slice stores indices of deleted options in
	// the order they are deleted.

	// Delete each covered item.
	for _, covered := range dl.options[index] {
		item := covered.item

		// Delete covered item from linked list.
		item.left.right = item.right
		item.right.left = item.left

		// Delete all options that cover the same item, since we can
		// only cover each item once.
		for conflict := item.head.down; conflict != item.head; conflict = conflict.down {
			// We can only delete nodes once; trying to re-delete may
			// break things.  So if we've already deleted something, don't
			// try delete it again.
			if intSliceContains(*deleted, conflict.option) {
				continue
			}

			// Record deleted option.
			*deleted = append(*deleted, conflict.option)

			// To delete an option, we go through and delete each entry in
			// the option.
			for _, entry := range dl.options[conflict.option] {
				entry.up.down = entry.down
				entry.down.up = entry.up

				// Update the corresponding item's record of remaining
				// items.
				entry.item.choices--
			}
		}
	}
}


type searchStage struct {
	choice int
	choices []int
}
type cleanupStage struct {
	choice int
	deleted []int
}


type stageType interface{}


func (dl *DancingLinks) GenerateSolutions(yield func([]Step) bool) {
	stages := []stageType{
		searchStage{-1, nil},
	}

	path := []Step{}

	for len(stages) > 0 {
		stage := stages[len(stages)-1]
		stages = stages[:len(stages)-1]

		switch stage := stage.(type) {

		case cleanupStage:
			path = path[:len(path)-1]

			// Uncover items in reverse order.
			entries := dl.options[stage.choice]
			for i := range entries {
				// We deleted the items left to right (increasing index), so we
				// uncover the items right to left (decreasing index).
				entry := entries[len(entries)-1-i]
				item := entry.item

				// Uncover item.
				item.left.right = item
				item.right.left = item
			}

			// Restore conflicting options in reverse order.
			for i := range stage.deleted {
				// Retrieve index of deleted option, in reverse order.
				option := stage.deleted[len(stage.deleted)-1-i]

				// To restore the option, we restore each entry in the option.
				for _, entry := range dl.options[option] {
					entry.up.down = entry
					entry.down.up = entry

					// Update item's choices counter.
					entry.item.choices++
				}
			}

		case searchStage:
			if stage.choice != -1 {
				deleted := []int{}
				dl.chooseOption(stage.choice, &deleted)
				stages = append(stages, cleanupStage{stage.choice, deleted})
				path = append(path, Step{stage.choice, stage.choices})
			}

			choices := dl.nextChoices()

			if choices == nil {
				keepGoing := yield(append([]Step{}, path...))

				// If the yield call returned false, then _after cleaning up and
				// restoring the link states we quit.
				if !keepGoing {
					return
				}

				break
			}

			// Consider each option that covers the first item.
			for i := range choices {
				choice := choices[len(choices)-1-i]
				stages = append(stages, searchStage{choice, choices})
			}
		}
	}
}

func (dl *DancingLinks) nextChoices() []int {
	// First item to cover.  We find the item with the fewest remaining
	// choices.
	first := dl.itemHead.right
	for item := first; item != dl.itemHead; item = item.right {
		if item.choices < first.choices {
			first = item
		}
	}

	// Nothing left to cover!
	if first == dl.itemHead {
		return nil
	}

	choices := []int{}
	for choice := first.head.down; choice != first.head; choice = choice.down {
		choices = append(choices, choice.option)
	}

	return choices
}

func intSliceContains(slice []int, element int) bool {
	for _, e := range slice {
		if e == element {
			return true
		}
	}
	return false
}
