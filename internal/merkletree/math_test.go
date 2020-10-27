package merkletree

// This file tests basic functions of the MMR data structure
// independent of any specific implementation of MMRs.

import (
	"encoding/binary"
	"testing"
)

// {pos, height}
var heightTable = [][]int{
	{0, 0},
	{1, 0},
	{2, 1},
	{3, 0},
	{4, 0},
	{5, 1},
	{6, 2},
	{7, 0},
	{8, 0},
	{9, 1},
	{10, 0},
	{11, 0},
	{12, 1},
	{13, 2},
	{14, 3},
}

func TestHeight(t *testing.T) {
	for _, c := range heightTable {
		pos, out := c[0], c[1]
		result := height(pos)
		if result != out {
			t.Errorf("height(%d): %d (expected %d)", pos, result, out)
		}
	}
}

func TestPeaks(t *testing.T) {
	// {size, peaks...}
	table := [][]int{
		{1, 0},
		{2, 0, 1},
		{3, 2},
		{4, 2, 3},
		{5, 2, 3, 4},
		{6, 2, 5},
		{7, 6},
		{8, 6, 7},
		{9, 6, 7, 8},
		{10, 6, 9},

		{20, 14, 17, 18, 19},
		{40, 30, 37, 38, 39},
		{80, 62, 77, 78, 79},
		{160, 126, 157, 158, 159},
		{202, 126, 189, 196, 199, 200, 201},
		{217, 126, 189, 204, 211, 214, 215, 216},
		{248, 126, 189, 220, 235, 242, 245, 246, 247},
		{255, 254},
		{345, 254, 317, 332, 339, 342, 343, 344},
		{481, 254, 381, 444, 475, 478, 479, 480},
		{511, 510},
	}
	for _, c := range table {
		in, out := c[0], c[1:]
		if result := peaks(in); !intSliceEqual(out, result) {
			t.Errorf("peaks(%d): expected '%v', got '%v'", in, out, result)
			continue
		}
	}
}

func TestLeftChild(t *testing.T) {
	// pos, height, first child
	table := [][]int{
		{2, 1, 0},
		{5, 1, 3},
		{6, 2, 2},
	}
	for _, vals := range table {
		pos, h, fc := vals[0], vals[1], vals[2]
		if out := leftChild(pos, h); out != fc {
			t.Errorf("leftChild(%d, %d) is %d, expected %d", pos, h, out, fc)
		}
	}
}

func TestChildren(t *testing.T) {
	// pos, height, children...
	table := [][]int{
		// leaves
		{0, 0},
		{1, 0},

		// non-leaves
		{2, 1, 0, 1},
		{5, 1, 3, 4},
		{9, 1, 7, 8},
		{12, 1, 10, 11},
		{6, 2, 2, 5},
		{13, 2, 9, 12},
		{14, 3, 6, 13},
	}
	for _, vals := range table {
		pos, h, expected := vals[0], vals[1], vals[2:]
		if out := children(pos, h); !intSliceEqual(expected, out) {
			t.Errorf("children(%d, %d): expected '%v', got '%v'", pos, h, expected, out)
		}
	}

}

func TestPath(t *testing.T) {
	// pos, pathEntry's...
	_ = map[int][]pathEntry{
		0: {},
		1: {{}, {}},
	}
	// compare arrays of pathEntry's: p, expected
}

func TestProof(t *testing.T) {
	t.Errorf("%v", proof(3, 19))
}

func TestSize(t *testing.T) {
	var i uint64
	t.Errorf("%d", binary.Size(i))
}

type testcase struct {
	a        []int
	b        []int
	expected bool
}

func TestIntSliceEqual(t *testing.T) {
	cases := []testcase{
		{[]int{1, 1, 2, 3}, []int{1, 1, 2, 3}, true},
		{[]int{1, 1, 2, 3}, []int{1, 1, 2}, false},
		{[]int{1, 1, 2, 3}, []int{}, false},
		{[]int{1, 1, 2, 3}, nil, false},
		{nil, nil, true},
	}

	for _, c := range cases {
		o := intSliceEqual(c.a, c.b)
		if o != c.expected {
			t.Errorf("intSliceEqual(%v, %v): %v, expected %v", c.a, c.b, o, c.expected)
		}
	}
}
