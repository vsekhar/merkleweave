// Package merkletree provides a simple Merkle tree implementation.
//
// Under the hood it uses a Merkle Mountain Range (MMR).
package merkletree

import (
	"bytes"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/sha3"
)

// HashLength is the number of bytes to read from the Shake hash.
const HashLength = 64

// MerkleTree is a simple Merkle Tree data structure.
type MerkleTree struct {
	data  [][]byte           // from users
	nodes [][HashLength]byte // hashes of data and children
}

// New returns a new empty MerkleTree.
func New() *MerkleTree {
	return &MerkleTree{
		data:  make([][]byte, 0),
		nodes: make([][HashLength]byte, 0),
	}
}

// Len returns the number of entries in the MerkleTree.
func (m *MerkleTree) Len() int {
	return len(m.data)
}

// At returns the entry at pos in the MerkleTree. If pos does not exist in the
// MerkleTree, At panics.
func (m *MerkleTree) At(pos int) []byte {
	return m.data[pos]
}

// Append adds an entry to the MerkleTree.
func (m *MerkleTree) Append(b []byte) {
	pos := len(m.data)
	h := height(pos)
	shaker := sha3.NewShake256()

	// Hash left child and write child (if not a leaf).
	cs := children(pos, h)
	if cs != nil {
		left := m.nodes[cs[0]][:]
		right := m.nodes[cs[1]][:]

		if _, err := shaker.Write(left); err != nil {
			panic(err)
		}
		if _, err := shaker.Write(right); err != nil {
			panic(err)
		}
	}

	// Hash the current node's data.
	_, err := shaker.Write(b[:])
	if err != nil {
		panic(err)
	}

	// Store.
	var node [HashLength]byte
	_, err = shaker.Read(node[:])
	if err != nil {
		panic(err)
	}
	m.data = append(m.data, b)
	m.nodes = append(m.nodes, node)
}

// Summary is a summary of a tree.
type Summary struct {
	N       int
	Summary [HashLength]byte
}

// Equals returns true if the summaries are equal.
func (s Summary) Equals(s2 Summary) bool {
	if s.N != s2.N {
		return false
	}
	if bytes.Compare(s.Summary[:], s2.Summary[:]) != 0 {
		return false
	}
	return true
}

// ShortString returns a short string representation of a Summary.
//
// Hashes are truncated to 8 base64-encoded characters.
func (s Summary) String() string {
	hash := base64.RawURLEncoding.EncodeToString(s.Summary[:])
	return fmt.Sprintf("%d:%s", s.N, hash)
}

// Summary returns the length and hash of the Merkle tree.
func (m *MerkleTree) Summary() Summary {
	s := Summary{}
	s.N = m.Len()
	ps := peaks(s.N)
	shaker := sha3.NewShake256()
	for _, pos := range ps {
		if _, err := shaker.Write(m.nodes[pos][:]); err != nil {
			panic(err)
		}
	}
	if _, err := shaker.Read(s.Summary[:]); err != nil {
		panic(err)
	}
	return s
}

// EmptyTreeSummary is the fixed summary of an empty Merkle tree.
var EmptyTreeSummary Summary = Summary{
	N:       0,
	Summary: [HashLength]byte{0x46, 0xb9, 0xdd, 0x2b, 0xb, 0xa8, 0x8d, 0x13, 0x23, 0x3b, 0x3f, 0xeb, 0x74, 0x3e, 0xeb, 0x24, 0x3f, 0xcd, 0x52, 0xea, 0x62, 0xb8, 0x1b, 0x82, 0xb5, 0xc, 0x27, 0x64, 0x6e, 0xd5, 0x76, 0x2f, 0xd7, 0x5d, 0xc4, 0xdd, 0xd8, 0xc0, 0xf2, 0x0, 0xcb, 0x5, 0x1, 0x9d, 0x67, 0xb5, 0x92, 0xf6, 0xfc, 0x82, 0x1c, 0x49, 0x47, 0x9a, 0xb4, 0x86, 0x40, 0x29, 0x2e, 0xac, 0xb3, 0xb7, 0xc4, 0xbe},
}

// TODO: ProveEntry
// TODO: ProveSummary
