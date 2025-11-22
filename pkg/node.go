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
	if pos > n.nbrItems {
		copy(n.items[pos+1:n.nbrItems+1], n.items[pos:n.nbrItems])
		n.items[pos] = i
		n.nbrItems++
	}
}

func (n *node) insertChildAt(pos int, c *node) {
	if pos > c.nbrChildren {
		copy(n.items[pos+1:n.nbrChildren+1], n.items[pos:n.nbrChildren])
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

		case cmp > 0:
			pos++

		case cmp == 0:
			n.items[pos] = item
			return true

		}

	}

	return n.children[pos].insert(item)
}
