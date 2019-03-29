package dancinglinks

// The setup for an exact cover problem, which consists of (1) a set
// of items to cover and (2) a collection of options, i.e. subsets of
// the items.  The exact cover problem solver returns a selection of
// the option such that each item is contained in exactly one of the
// selected options.
type DancingLinks struct {
	// A list of options.  Each "option" is a list of pointers to the
	// entries provided by that option.
	options [][]*entryNode

	// Blank anchor node, whose `right` points to the first/left-most
	// item to be covered.
	itemHead *itemNode

	// Indices of required options, i.e. options that are required to be
	// in the selection.
	selected []int

	// Indices of options that were removed when selecting the
	// pre-selected/required options.
	deleted []int
}

// A decision step in the exact cover solution path.  At each step,
// the algorithm finds the lowest-index item with the fewest remaining
// options and selects one of the options that cover the item.  A Step
// records the item covered in that step, the option selected to cover
// that item, and all the remaining available options that cover the
// item.
type Step struct {
	// Index of the item to be covered by this step.
	Item int

	// Index of the option selected in this step to cover the item.
	// Option is guaranteed to be an element of Choices.
	Option int

	// All (remaining) available options that cover the item.  Choices
	// is guaranteed to contain Option.
	Choices []int
}

// A linked list node storing an item in an exact cover setup.
type itemNode struct {
	// The index of the item.
	index int

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

// A linked list node storing an entry (a 1 in the exact cover matrix)
// in the exact cover setup.
type entryNode struct {
	// The item covered by this entry.
	item *itemNode

	// The index of the option this entry belongs to.
	option int

	// Linked list neighbors.
	up   *entryNode
	down *entryNode
}

type stage struct {
	item    int
	parent  int
	deleted []int
	choices []int
	i       int
}

func New(itemCount int, options [][]int) *DancingLinks {
	dl := &DancingLinks{
		options:  make([][]*entryNode, len(options)),
		itemHead: &itemNode{index: -1},
		selected: []int{},
		deleted:  []int{},
	}

	// Construct item list.
	items := make([]*itemNode, itemCount)
	lastItem := dl.itemHead
	for index := range items {
		newItem := &itemNode{
			index: index,
			left:  lastItem,
			head:  &entryNode{option: -1},
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

func (dl *DancingLinks) ForceOptions(indices ...int) {
	for _, index := range indices {
		dl.selected = append(dl.selected, index)
		dl.chooseOption(index, &dl.deleted)
	}
}

func (dl *DancingLinks) UnforceOptions() {
	dl.restoreOptions(dl.deleted)
	dl.deleted = dl.deleted[:0]
	dl.selected = dl.selected[:0]
}

func (dl *DancingLinks) GenerateSolutions(yield func([]Step) bool) bool {

	item, choices := dl.nextChoices()
	if choices == nil {
		yield([]Step{})
		return true
	}

	stages := []*stage{
		&stage{
			item:    item,
			parent:  -1,
			deleted: nil,
			choices: choices,
			i:       0,
		},
	}

	path := []Step{}
	keepGoing := true

	for {
		s := stages[len(stages)-1]

		if s.i == len(s.choices) || !keepGoing {
			stages = stages[:len(stages)-1]

			if s.parent == -1 {
				return keepGoing
			}

			path = path[:len(path)-1]
			dl.unchooseOption(s.parent, s.deleted)
			continue
		}

		deleted := []int{}
		dl.chooseOption(s.choices[s.i], &deleted)
		path = append(path, Step{s.item, s.choices[s.i], s.choices})

		item, choices := dl.nextChoices()

		if choices == nil {
			keepGoing = yield(append([]Step{}, path...))
		}

		// Consider each option that covers the first item.
		stages = append(stages, &stage{
			item:    item,
			parent:  s.choices[s.i],
			deleted: deleted,
			choices: choices,
			i:       0,
		})

		s.i++
	}
}

func (dl *DancingLinks) GenerateCovers(yield func([]int) bool) {
	dl.GenerateSolutions(func(solution []Step) bool {
		cover := make([]int, len(solution))
		for i, step := range solution {
			cover[i] = step.Option
		}
		return yield(cover)
	})
}

func (dl *DancingLinks) AllSolutions() [][]Step {
	solutions := make([][]Step, 0)
	dl.GenerateSolutions(func(solution []Step) bool {
		solutions = append(solutions, solution)
		return true
	})
	return solutions
}

func (dl *DancingLinks) AllCovers() [][]int {
	covers := make([][]int, 0)
	dl.GenerateCovers(func(cover []int) bool {
		covers = append(covers, cover)
		return true
	})
	return covers
}

func (dl *DancingLinks) AnySolution() []Step {
	var solution []Step
	dl.GenerateSolutions(func(s []Step) bool {
		solution = s
		return false
	})
	return solution
}

func (dl *DancingLinks) AnyCover() []int {
	var cover []int
	dl.GenerateCovers(func(c []int) bool {
		cover = c
		return false
	})
	return cover
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

func (dl *DancingLinks) unchooseOption(index int, deleted []int) {
	// Uncover items in reverse order.
	entries := dl.options[index]
	for i := range entries {
		// We deleted the items left to right (increasing index), so we
		// uncover the items right to left (decreasing index).
		entry := entries[len(entries)-1-i]
		item := entry.item

		// Uncover item.
		item.left.right = item
		item.right.left = item
	}

	dl.restoreOptions(deleted)
}

func (dl *DancingLinks) restoreOptions(options []int) {
	// Restore conflicting options in reverse order.
	for i := range options {
		// To restore the option, we restore each entry in the option.
		for _, entry := range dl.options[options[len(options)-1-i]] {
			entry.up.down = entry
			entry.down.up = entry

			// Update item's choices counter.
			entry.item.choices++
		}
	}
}

func (dl *DancingLinks) nextChoices() (int, []int) {
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
		return -1, nil
	}

	choices := []int{}
	for choice := first.head.down; choice != first.head; choice = choice.down {
		choices = append(choices, choice.option)
	}

	return first.index, choices
}

func intSliceContains(slice []int, element int) bool {
	for _, e := range slice {
		if e == element {
			return true
		}
	}
	return false
}
