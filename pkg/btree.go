package pkg

import "errors"

const (
	degree      = 5
	maxChildren = 2 * degree
	maxItems    = maxChildren - 1
	minItems    = degree - 1
)

type BTree struct {
	root *node
}

func NewBtree() *BTree {
	return &BTree{}
}

func (t *BTree) Find(key []byte) ([]byte, error) {
	for next := t.root; next != nil; {
		pos, found := next.Search(key)
		if found {
			return next.items[pos].val, nil
		}
		next = next.children[pos]
	}
	return nil, errors.New("key not found")
}
