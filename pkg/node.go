package pkg

import "bytes"

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
