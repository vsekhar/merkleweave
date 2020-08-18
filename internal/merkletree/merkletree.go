// Package merkletree provides a simple Merkle tree implementation.
//
// Under the hood it uses a Merkle Mountain Range (MMR).
package merkletree

import (
	"golang.org/x/crypto/sha3"
)

const hashLength = 64

// MerkleTree is a simple Merkle Tree data structure.
type MerkleTree struct {
	data  [][hashLength]byte // from users
	nodes [][hashLength]byte // hashes of data and children
}

// New returns a new empty MerkleTree.
func New() *MerkleTree {
	return &MerkleTree{
		data:  make([][hashLength]byte, 0),
		nodes: make([][hashLength]byte, 0),
	}
}

// Len returns the number of entries in the MerkleTree.
func (m *MerkleTree) Len() int {
	return len(m.data)
}

// At returns the entry at pos in the MerkleTree. If pos does not exist in the
// MerkleTree, At panics.
func (m *MerkleTree) At(pos int) [hashLength]byte {
	return m.data[pos]
}

// Append adds an entry to the MerkleTree.
func (m *MerkleTree) Append(b [hashLength]byte) error {
	pos := len(m.data)
	h := height(pos)
	shaker := sha3.NewShake256()

	// Hash left child and write child (if not a leaf).
	cs := children(pos, h)
	if cs != nil {
		left := m.nodes[cs[0]][:]
		right := m.nodes[cs[1]][:]

		if _, err := shaker.Write(left); err != nil {
			return err
		}
		if _, err := shaker.Write(right); err != nil {
			return err
		}
	}

	// Hash the current node's data.
	_, err := shaker.Write(b[:])
	if err != nil {
		return err
	}

	// Store.
	var node [hashLength]byte
	_, err = shaker.Read(node[:])
	if err != nil {
		return err
	}
	m.data = append(m.data, b)
	m.nodes = append(m.nodes, node)
	return nil
}

// Summary returns a hash of the MerkleTree.
func (m *MerkleTree) Summary() (r [hashLength]byte, err error) {
	ps := peaks(m.Len())
	shaker := sha3.NewShake256()
	for _, pos := range ps {
		if _, err = shaker.Write(m.nodes[pos][:]); err != nil {
			return
		}
	}
	_, err = shaker.Read(r[:])
	return
}
