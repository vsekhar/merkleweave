// Package merkleweave provides in in-memory write-optimized Merkle tree-like
// data structure. Unlike Merkle trees, a Merkle weave supports concurrent
// writes.
package merkleweave

import (
	"sync"

	"github.com/vsekhar/merkleweave/internal/merkletree"
)

const prefixBytes = 1
const numTrees = 1 << (prefixBytes * 8)

type tree struct {
	m sync.Mutex
	t *merkletree.MerkleTree
}

type treeMap map[byte]tree

// MerkleWeave is a write-optimized Merkle tree-like data structure.
type MerkleWeave struct {
	ts treeMap
}

// New returns a new MerkleWeave.
func New() *MerkleWeave {
	ret := &MerkleWeave{
		ts: make(map[byte]tree),
	}
	for i := 0; i < numTrees; i++ {
		ret.ts[byte(i)] = tree{
			t: merkletree.New(),
		}
	}
	return ret
}

// TODO: append
