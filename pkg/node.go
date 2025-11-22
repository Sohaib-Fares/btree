package pkg

import (
	"bytes"
)

type item struct {
	key []byte
	val []byte
}

type node struct {
	items       [maxItems]*item
	children    [maxChildren]*node
	nbrChildren int
	nbrItems    int
}

func (n *node) isLeaf() bool {
	return n.nbrChildren == 0
}

func (n *node) Search(key []byte) (int, bool) {
	low, high := 0, n.nbrItems
	var mid int

	for low < high {
		mid = (low + high) / 2
		cmp := bytes.Compare(key, n.items[mid].key)

		switch cmp {
		case -1:
			high = mid
		case 1:
			low = mid + 1
		case 0:
			return mid, true
		}
	}
	return low, false
}

func (n *node) insertItemAt(pos int, i *item) {
	if pos <= n.nbrItems {
		copy(n.items[pos+1:n.nbrItems+1], n.items[pos:n.nbrItems])
		n.items[pos] = i
		n.nbrItems++
	}
}

func (n *node) insertChildAt(pos int, c *node) {
	if pos <= n.nbrChildren {
		copy(n.children[pos+1:n.nbrChildren+1], n.children[pos:n.nbrChildren])
	}
	n.children[pos] = c
	n.nbrChildren++
}

func (n *node) Split() (*item, *node) {
	mid := minItems
	midItem := n.items[mid]

	newNode := &node{}
	newNode.nbrItems = copy(newNode.items[:], n.items[mid+1:])

	if !n.isLeaf() {
		copy(newNode.children[:], n.children[mid+1:])
		newNode.nbrChildren = newNode.nbrItems + 1
	}

	for i, l := mid, n.nbrItems; i < l; i++ {
		n.items[i] = nil
		n.nbrItems--

		if !n.isLeaf() {
			n.children[i+1] = nil
			n.nbrChildren--
		}
	}
	return midItem, newNode
}

func (n *node) insert(item *item) bool {

	pos, found := n.Search(item.key)

	if found {
		n.items[pos] = item
		return false
	}

	if n.isLeaf() {
		n.insertItemAt(pos, item)
		return true
	}

	if n.children[pos].nbrItems >= maxItems {

		midItem, newNode := n.children[pos].Split()

		n.insertItemAt(pos, midItem)
		n.insertChildAt(pos+1, newNode)

		switch cmp := bytes.Compare(item.key, n.items[pos].key); {
		case cmp < 0:
			// do nothing, keep the same path
		case cmp > 0:
			pos++

		case cmp == 0:

			n.items[pos] = item
			return false // This is an update, not a new insertion so we return false

		}

	}

	return n.children[pos].insert(item)
}

func (n *node) removeItemAt(pos int) *item {
	removedItem := n.items[pos]
	n.items[pos] = nil

	if lastPos := n.nbrItems - 1; pos < lastPos {
		copy(n.items[pos:lastPos], n.items[pos+1:lastPos+1])
		n.items[lastPos] = nil
	}
	n.nbrItems--

	return removedItem
}

func (n *node) removeChildAt(pos int) *node {
	removedCHild := n.children[pos]
	n.children[pos] = nil

	if lastPos := n.nbrChildren - 1; pos < lastPos {
		copy(n.children[pos:lastPos], n.children[pos+1:lastPos+1])
		n.children[lastPos] = nil
	}
	n.nbrChildren--

	return removedCHild
}

func (n *node) fillChildAt(pos int) {

	switch {

	case pos > 0 && n.children[pos-1].nbrItems > minItems:

		left, right := n.children[pos-1], n.children[pos]
		copy(right.items[1:right.nbrItems+1], right.items[:right.nbrItems])
		right.items[0] = n.items[pos-1]
		right.nbrItems++
		if !left.isLeaf() {
			right.insertChildAt(0, left.removeChildAt(left.nbrChildren-1))
		}
		n.items[pos-1] = left.removeItemAt(left.nbrItems - 1)

	case pos < n.nbrChildren-1 && n.children[pos+1].nbrItems > minItems:

		left, right := n.children[pos], n.children[pos+1]
		left.items[left.nbrItems] = n.items[pos]
		left.nbrItems++
		if !right.isLeaf() {
			left.insertChildAt(left.nbrChildren, right.removeChildAt(0))
		}
		n.items[pos] = right.removeItemAt(0)

	default:

		if pos >= n.nbrItems {
			pos = n.nbrItems - 1
		}

		left, right := n.children[pos], n.children[pos+1]
		left.items[left.nbrItems] = n.removeItemAt(pos)
		left.nbrItems++
		copy(left.items[left.nbrItems:], right.items[:right.nbrItems])
		left.nbrItems += right.nbrItems
		if !left.isLeaf() {
			copy(left.children[left.nbrChildren:], right.children[:right.nbrChildren])
			left.nbrChildren += right.nbrChildren
		}
		n.removeChildAt(pos + 1)
		right = nil
	}
}

func (n *node) delete(key []byte, isSeekingSuccessor bool) *item {
	pos, found := n.Search(key)

	var next *node

	if found {
		if n.isLeaf() {
			return n.removeItemAt(pos)
		}
		next, isSeekingSuccessor = n.children[pos+1], true
	} else {
		next = n.children[pos]
	}
	if n.isLeaf() && isSeekingSuccessor {
		return n.removeItemAt(0)
	}
	if next == nil {
		return nil
	}

	deletedItem := next.delete(key, isSeekingSuccessor)

	if found && isSeekingSuccessor {
		n.items[pos] = deletedItem
	}

	if next.nbrItems < minItems {
		if found && isSeekingSuccessor {
			n.fillChildAt(pos + 1)
		} else {
			n.fillChildAt(pos)
		}
	}
	return deletedItem
}
