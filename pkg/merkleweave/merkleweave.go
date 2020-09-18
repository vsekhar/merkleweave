// Package merkleweave provides in in-memory write-optimized Merkle tree-like
// data structure. Unlike Merkle trees, a Merkle weave supports concurrent
// writes.
package merkleweave

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/vsekhar/merkleweave/internal/merkletree"
)

const prefixBytes = 1
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

func toInt(p prefix) int {
	r := 0
	for i := 0; i < prefixBytes; i++ {
		r += int(p[i]) << (i * 8)
	}
	return r
}

func fromHex(s string) prefix {
	r := prefix{}
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	if len(b) != prefixBytes {
		panic(fmt.Sprintf("expected %d bytes for prefix, got %d", prefixBytes, len(b)))
	}
	copy(r[:], b)
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

// forEach runs f on each tree in parallel.
func (m *MerkleWeave) forEach(f func(i int, t *merkletree.MerkleTree)) {
	wg := sync.WaitGroup{}
	wg.Add(len(m.ts))
	for p, t := range m.ts {
		go func(p prefix, t tree) {
			t.m.Lock()
			defer t.m.Unlock()
			f(toInt(p), t.t)
			wg.Done()
		}(p, t)
	}
	wg.Wait()
}

// Append adds an entry to a MerkleWeave.
func (m *MerkleWeave) Append(b []byte) {
	if len(b) < minDataLen {
		panic(fmt.Sprintf("at least %d bytes needed, got %d bytes", minDataLen, len(b)))
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
		t.Append(b)
	}
}

// ApproxLen returns an approximate number of entries in the Merkle weave. The Merkle weave can contain spurious entries
func (m *MerkleWeave) ApproxLen() int {
	lens := [numTrees]int{}
	m.forEach(func(i int, t *merkletree.MerkleTree) {
		lens[i] = t.Len()
	})
	l := 0
	for _, i := range lens {
		l += i
	}
	return l
}

// Summary is a summary of a Merkle weave.
type Summary struct {
	ss [numTrees]merkletree.Summary
}

// Equals returns true if the Summary's are equal.
func (s *Summary) Equals(s2 *Summary) bool {
	for i, t := range s.ss {
		if !t.Equals(s2.ss[i]) {
			return false
		}
	}
	return true
}

// ShortString returns a short string representation of a Summary.
//
// Empty sub-trees (length of zero, fixed hash) are skipped. Hashes are
// truncated to 8 base64-encoded characters.
func (s *Summary) ShortString() string {
	var b strings.Builder
	for i := 0; i < numTrees; i++ {
		if s.ss[i].N == 0 {
			continue
		}
		prefix := fromInt(i)
		fmt.Fprintf(&b, "%x:%s; ", prefix, s.ss[i].String())
	}
	return b.String()
}

// Summary returns a summary of the Merkle weave.
func (m *MerkleWeave) Summary() Summary {
	r := Summary{}
	m.forEach(func(i int, t *merkletree.MerkleTree) {
		r.ss[i] = t.Summary()
		return
	})
	return r
}

// for testing
func newEmptySummary() Summary {
	s := Summary{}
	for i := range s.ss {
		s.ss[i] = merkletree.EmptyTreeSummary
	}
	return s
}

// TODO: ProveEntry
// TODO: ProveSummary
