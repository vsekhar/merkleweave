// Package merkletree provides a simple Merkle tree implementation.
//
// Under the hood it uses a Merkle Mountain Range (MMR).
package merkletree

import (
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

// Summary returns the length and hash of the MerkleTree.
func (m *MerkleTree) Summary() (len int, r [HashLength]byte) {
	len = m.Len()
	ps := peaks(len)
	shaker := sha3.NewShake256()
	for _, pos := range ps {
		if _, err := shaker.Write(m.nodes[pos][:]); err != nil {
			panic(err)
		}
	}
	if _, err := shaker.Read(r[:]); err != nil {
		panic(err)
	}
	return
}

// TODO: ProveEntry
// TODO: ProveSummary
