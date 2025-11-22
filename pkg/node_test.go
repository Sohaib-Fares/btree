package pkg

import (
	"bytes"
	"testing"
)

// Helper function to create a test item
func newTestItem(key, val string) *item {
	return &item{
		key: []byte(key),
		val: []byte(val),
	}
}

// Helper function to create a node with items
func newTestNode(keys ...string) *node {
	n := &node{}
	for _, key := range keys {
		n.items[n.nbrItems] = newTestItem(key, "val_"+key)
		n.nbrItems++
	}
	return n
}

// TestNodeIsLeaf tests the isLeaf method
func TestNodeIsLeaf(t *testing.T) {
	testCases := []struct {
		name        string
		nbrChildren int
		expected    bool
	}{
		{"leaf node", 0, true},
		{"internal node", 1, false},
		{"internal with many children", 5, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			n := &node{nbrChildren: tc.nbrChildren}
			if n.isLeaf() != tc.expected {
				t.Errorf("expected isLeaf() to return %v, got %v", tc.expected, n.isLeaf())
			}
		})
	}
}

// TestNodeSearch tests the binary search implementation
func TestNodeSearch(t *testing.T) {
	// Create a node with sorted keys
	n := newTestNode("apple", "banana", "cherry", "date", "elderberry")

	testCases := []struct {
		name      string
		searchKey string
		expectPos int
		expectOk  bool
	}{
		{"find first", "apple", 0, true},
		{"find middle", "cherry", 2, true},
		{"find last", "elderberry", 4, true},
		{"not found before first", "aardvark", 0, false},
		{"not found between", "blueberry", 2, false},
		{"not found after last", "fig", 5, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pos, found := n.Search([]byte(tc.searchKey))
			if found != tc.expectOk {
				t.Errorf("expected found=%v, got %v", tc.expectOk, found)
			}
			if pos != tc.expectPos {
				t.Errorf("expected position %d, got %d", tc.expectPos, pos)
			}
		})
	}
}

// TestNodeInsertItemAt tests inserting items at specific positions
func TestNodeInsertItemAt(t *testing.T) {
	n := newTestNode("a", "c", "e")

	testCases := []struct {
		name        string
		pos         int
		key         string
		expectedLen int
		checkKey    string
		checkPos    int
	}{
		{"insert at beginning", 0, "0", 4, "0", 0},
		{"insert in middle", 2, "d", 5, "d", 2},
		{"insert at end", 5, "z", 6, "z", 5},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			n.insertItemAt(tc.pos, newTestItem(tc.key, "val"))
			if n.nbrItems != tc.expectedLen {
				t.Errorf("expected %d items, got %d", tc.expectedLen, n.nbrItems)
			}
			if !bytes.Equal(n.items[tc.checkPos].key, []byte(tc.checkKey)) {
				t.Errorf("expected key %s at position %d, got %s",
					tc.checkKey, tc.checkPos, n.items[tc.checkPos].key)
			}
		})
	}
}

// TestNodeRemoveItemAt tests removing items
func TestNodeRemoveItemAt(t *testing.T) {
	testCases := []struct {
		name         string
		initialKeys  []string
		removePos    int
		expectedLen  int
		removedKey   string
		remainingKey string
		checkPos     int
	}{
		{
			name:         "remove first",
			initialKeys:  []string{"a", "b", "c"},
			removePos:    0,
			expectedLen:  2,
			removedKey:   "a",
			remainingKey: "b",
			checkPos:     0,
		},
		{
			name:         "remove middle",
			initialKeys:  []string{"a", "b", "c"},
			removePos:    1,
			expectedLen:  2,
			removedKey:   "b",
			remainingKey: "c",
			checkPos:     1,
		},
		{
			name:         "remove last",
			initialKeys:  []string{"a", "b", "c"},
			removePos:    2,
			expectedLen:  2,
			removedKey:   "c",
			remainingKey: "b",
			checkPos:     1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			n := newTestNode(tc.initialKeys...)
			removed := n.removeItemAt(tc.removePos)

			if !bytes.Equal(removed.key, []byte(tc.removedKey)) {
				t.Errorf("expected removed key %s, got %s", tc.removedKey, removed.key)
			}
			if n.nbrItems != tc.expectedLen {
				t.Errorf("expected %d items, got %d", tc.expectedLen, n.nbrItems)
			}
			if !bytes.Equal(n.items[tc.checkPos].key, []byte(tc.remainingKey)) {
				t.Errorf("expected key %s at position %d, got %s",
					tc.remainingKey, tc.checkPos, n.items[tc.checkPos].key)
			}
		})
	}
}

// TestNodeSplit tests node splitting
func TestNodeSplit(t *testing.T) {
	// Create a full node (degree=5, so minItems=4)
	n := newTestNode("a", "b", "c", "d", "e", "f", "g", "h", "i")

	midItem, rightNode := n.Split()

	// Check that middle item is returned
	if !bytes.Equal(midItem.key, []byte("e")) {
		t.Errorf("expected middle item 'e', got %s", midItem.key)
	}

	// Check left node has correct items
	if n.nbrItems != minItems {
		t.Errorf("expected left node to have %d items, got %d", minItems, n.nbrItems)
	}

	// Check right node has correct items
	expectedRightItems := 4 // items after middle
	if rightNode.nbrItems != expectedRightItems {
		t.Errorf("expected right node to have %d items, got %d", expectedRightItems, rightNode.nbrItems)
	}

	// Verify first item in right node
	if !bytes.Equal(rightNode.items[0].key, []byte("f")) {
		t.Errorf("expected first item in right node to be 'f', got %s", rightNode.items[0].key)
	}
}

// TestNodeInsert tests the insert method
func TestNodeInsert(t *testing.T) {
	testCases := []struct {
		name      string
		insertKey string
		insertVal string
		isNew     bool
		setupNode func() *node // Add function to setup node state
	}{
		{
			name:      "insert new key",
			insertKey: "newkey",
			insertVal: "newval",
			isNew:     true,
			setupNode: func() *node { return &node{} },
		},
		{
			name:      "insert duplicate",
			insertKey: "dupkey",
			insertVal: "updated",
			isNew:     false,
			setupNode: func() *node {
				n := &node{}
				n.insert(newTestItem("dupkey", "original"))
				return n
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			n := tc.setupNode()
			item := newTestItem(tc.insertKey, tc.insertVal)
			isNew := n.insert(item)

			if isNew != tc.isNew {
				t.Errorf("expected insert to return %v, got %v", tc.isNew, isNew)
			}

			// Verify item is in node
			pos, found := n.Search([]byte(tc.insertKey))
			if !found {
				t.Error("expected to find inserted key")
			}
			if !bytes.Equal(n.items[pos].val, []byte(tc.insertVal)) {
				t.Errorf("expected value %s, got %s", tc.insertVal, n.items[pos].val)
			}
		})
	}
}

// TestNodeSearchEmpty tests search on empty node
func TestNodeSearchEmpty(t *testing.T) {
	n := &node{}
	pos, found := n.Search([]byte("anything"))

	if found {
		t.Error("expected not to find key in empty node")
	}
	if pos != 0 {
		t.Errorf("expected position 0, got %d", pos)
	}
}
