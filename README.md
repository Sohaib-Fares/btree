# B-Tree Implementation in Go

A B-Tree key-value store library built as a learning project. Implements a B-Tree of order 5 with support for insertion, search, and deletion operations.

## Installation

```bash
go get github.com/red7-c/btree
```

## Quick Start

```go
package main

import (
	"fmt"
	"github.com/red7-c/btree/pkg"
)

func main() {
	tree := pkg.NewBtree()
	
	// Insert
	tree.Insert([]byte("name"), []byte("Alice"))
	tree.Insert([]byte("age"), []byte("30"))
	
	// Find
	val, err := tree.Find([]byte("name"))
	if err == nil {
		fmt.Println(string(val)) // Output: Alice
	}
	
	// Delete
	tree.Delete([]byte("age"))
}
```

## Features

- **Insert**: Add or update key-value pairs
- **Search**: Binary search for keys with O(log n) complexity
- **Delete**: Remove keys while maintaining B-Tree properties
- **Automatic balancing**: Tree rebalances on splits and merges

## B-Tree Properties

- Order: 5 (up to 9 items per node, up to 10 children)
- Stores byte slices for both keys and values
- Maintains sorted order for efficient range queries

## Status

Early development. Core operations are implemented and tested. Suitable for learning.

## Future Enhancements

- Iteration/traversal
- Configurable order instead of fixed one.
`