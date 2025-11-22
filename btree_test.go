package btree

import (
	"bytes"
	"fmt"
	"testing"

)

// TestBTreeInsertAndFind tests basic insertion and retrieval
func TestBTreeInsertAndFind(t *testing.T) {
	tree := NewBtree()

	// Table-driven test: define test cases in a slice
	testCases := []struct {
		name string
		key  []byte
		val  []byte
	}{
		{"insert hello", []byte("hello"), []byte("world")},
		{"insert foo", []byte("foo"), []byte("bar")},
		{"insert number", []byte("123"), []byte("456")},
		{"insert empty value", []byte("empty"), []byte("")},
	}

	// Insert all test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tree.Insert(tc.key, tc.val)
		})
	}

	// Verify all inserts by finding them
	for _, tc := range testCases {
		t.Run("find_"+tc.name, func(t *testing.T) {
			val, err := tree.Find(tc.key)
			if err != nil {
				t.Errorf("expected to find key %s, got error: %v", tc.key, err)
			}
			if !bytes.Equal(val, tc.val) {
				t.Errorf("expected value %s, got %s", tc.val, val)
			}
		})
	}
}

// TestBTreeFindNonExistent tests finding keys that don't exist
func TestBTreeFindNonExistent(t *testing.T) {
	tree := NewBtree()
	tree.Insert([]byte("exists"), []byte("value"))

	testCases := []struct {
		name string
		key  []byte
	}{
		{"find non-existent", []byte("does-not-exist")},
		{"find empty key", []byte("")},
		{"find nil key", nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tree.Find(tc.key)
			if err == nil {
				t.Errorf("expected error when finding non-existent key %s", tc.key)
			}
		})
	}
}

// TestBTreeInsertDuplicate tests inserting duplicate keys (should update)
func TestBTreeInsertDuplicate(t *testing.T) {
	tree := NewBtree()

	key := []byte("duplicate")
	tree.Insert(key, []byte("first"))
	tree.Insert(key, []byte("second"))

	val, err := tree.Find(key)
	if err != nil {
		t.Fatalf("expected to find key, got error: %v", err)
	}
	if !bytes.Equal(val, []byte("second")) {
		t.Errorf("expected value 'second', got %s", val)
	}
}

// TestBTreeLargeInsertion tests inserting many items to trigger splits
func TestBTreeLargeInsertion(t *testing.T) {
	tree := NewBtree()
	numItems := 100

	// Insert many items
	for i := 0; i < numItems; i++ {
		key := []byte(fmt.Sprintf("key_%03d", i))
		val := []byte(fmt.Sprintf("val_%03d", i))
		tree.Insert(key, val)
	}

	// Verify all items can be found
	for i := 0; i < numItems; i++ {
		key := []byte(fmt.Sprintf("key_%03d", i))
		expectedVal := []byte(fmt.Sprintf("val_%03d", i))

		val, err := tree.Find(key)
		if err != nil {
			t.Errorf("expected to find key %s, got error: %v", key, err)
		}
		if !bytes.Equal(val, expectedVal) {
			t.Errorf("expected value %s, got %s", expectedVal, val)
		}
	}
}

// TestBTreeDelete tests deletion of keys
func TestBTreeDelete(t *testing.T) {
	tree := NewBtree()

	// Insert test data
	keys := [][]byte{
		[]byte("key1"),
		[]byte("key2"),
		[]byte("key3"),
		[]byte("key4"),
		[]byte("key5"),
	}

	for i, key := range keys {
		tree.Insert(key, []byte(fmt.Sprintf("val%d", i+1)))
	}

	testCases := []struct {
		name      string
		key       []byte
		shouldDel bool
	}{
		{"delete existing key", []byte("key3"), true},
		{"delete another key", []byte("key1"), true},
		{"delete non-existent", []byte("not-there"), false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			deleted := tree.Delete(tc.key)
			if deleted != tc.shouldDel {
				t.Errorf("expected Delete to return %v, got %v", tc.shouldDel, deleted)
			}

			// Verify key is gone
			if tc.shouldDel {
				_, err := tree.Find(tc.key)
				if err == nil {
					t.Errorf("expected key %s to be deleted, but still found", tc.key)
				}
			}
		})
	}
}

// TestBTreeDeleteFromEmpty tests deleting from empty tree
func TestBTreeDeleteFromEmpty(t *testing.T) {
	tree := NewBtree()

	deleted := tree.Delete([]byte("anything"))
	if deleted {
		t.Error("expected Delete to return false on empty tree")
	}
}

// TestBTreeDeleteAll tests deleting all items
func TestBTreeDeleteAll(t *testing.T) {
	tree := NewBtree()

	// Insert items
	keys := [][]byte{
		[]byte("a"), []byte("b"), []byte("c"),
		[]byte("d"), []byte("e"),
	}

	for _, key := range keys {
		tree.Insert(key, []byte("value"))
	}

	// Delete all items
	for _, key := range keys {
		deleted := tree.Delete(key)
		if !deleted {
			t.Errorf("expected to delete key %s", key)
		}
	}

	// Verify all are gone
	for _, key := range keys {
		_, err := tree.Find(key)
		if err == nil {
			t.Errorf("key %s should be deleted but still found", key)
		}
	}
}

// TestBTreeStressTest tests mixed operations
func TestBTreeStressTest(t *testing.T) {
	tree := NewBtree()

	// Insert 50 items
	for i := 0; i < 50; i++ {
		key := []byte(fmt.Sprintf("key_%02d", i))
		tree.Insert(key, []byte(fmt.Sprintf("val_%02d", i)))
	}

	// Delete every other item
	for i := 0; i < 50; i += 2 {
		key := []byte(fmt.Sprintf("key_%02d", i))
		tree.Delete(key)
	}

	// Verify remaining items
	for i := 1; i < 50; i += 2 {
		key := []byte(fmt.Sprintf("key_%02d", i))
		_, err := tree.Find(key)
		if err != nil {
			t.Errorf("expected to find key %s", key)
		}
	}

	// Verify deleted items are gone
	for i := 0; i < 50; i += 2 {
		key := []byte(fmt.Sprintf("key_%02d", i))
		_, err := tree.Find(key)
		if err == nil {
			t.Errorf("key %s should be deleted", key)
		}
	}
}
