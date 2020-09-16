// Package merkleweave provides in in-memory write-optimized Merkle tree-like
// data structure. Unlike Merkle trees, a Merkle weave supports concurrent
// writes.
package merkleweave

import (
	"bytes"
	"fmt"
	"sort"
	"sync"

	"github.com/vsekhar/merkleweave/internal/merkletree"
)

const prefixBytes = 2
const numTrees = 1 << (prefixBytes * 8)
const numCrossTrees = 2
const minDataLen = prefixBytes * numCrossTrees

type prefix [prefixBytes]byte

func (p *prefix) Less(p2 prefix) bool {
	r := bytes.Compare(p[:], p2[:])
	return r == -1
}

type prefixes []prefix

func (ps prefixes) Len() int {
	return len(ps)
}

func (ps prefixes) Less(i, j int) bool {
	return ps[i].Less(ps[j])
}

func (ps prefixes) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}

func fromInt(n int) prefix {
	r := prefix{}
	for i := 0; i < prefixBytes; i++ {
		r[i] = byte(n % (1 << 8))
		n /= (1 << 8)
	}
	if n > 0 {
		panic("overflow")
	}
	return r
}

func prefixesOf(b []byte) [numCrossTrees]prefix {
	var r [numCrossTrees]prefix
	for i := 0; i < numCrossTrees; i++ {
		p := prefix{}
		s := p[:]
		if n := copy(s, b); n != prefixBytes {
			panic("too few bytes for a prefix")
		}
		r[i] = p
		b = b[prefixBytes:]
	}
	return r
}

// SaltLen is the number of bytes of random data generated when adding an entry.
const SaltLen = 64

type tree struct {
	m *sync.Mutex
	t *merkletree.MerkleTree
}

type treeMap map[prefix]tree

// MerkleWeave is a write-optimized Merkle tree-like data structure.
type MerkleWeave struct {
	ts treeMap
}

// New returns a new MerkleWeave.
func New() *MerkleWeave {
	ret := &MerkleWeave{ts: make(treeMap)}
	for i := 0; i < numTrees; i++ {
		t := tree{
			m: new(sync.Mutex),
			t: merkletree.New(),
		}
		ret.ts[fromInt(i)] = t
	}
	return ret
}

// Append adds an entry to a MerkleWeave.
func (m *MerkleWeave) Append(b []byte) error {
	if len(b) < minDataLen {
		return fmt.Errorf("at least %d bytes needed, got %d bytes", minDataLen, len(b))
	}
	ps := prefixesOf(b)

	// sort and dedupe to prevent deadlock
	deduped := make(map[prefix]struct{})
	for _, v := range ps {
		deduped[v] = struct{}{}
	}
	sorted := make(prefixes, 0, len(ps))
	for k := range deduped {
		sorted = append(sorted, k)
	}
	sort.Sort(sorted)
	for _, p := range sorted {
		l := m.ts[p].m
		l.Lock()
		defer l.Unlock()
	}

	for i := 0; i < numCrossTrees; i++ {
		t := m.ts[ps[i]].t
		if err := t.Append(b); err != nil {
			// If an error occurs, spurious entries in other trees may be left
			// behind but this is ok.
			return err
		}
	}

	return nil
}
