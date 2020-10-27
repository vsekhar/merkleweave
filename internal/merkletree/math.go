package merkletree

//  Algos: https://github.com/mimblewimble/grin/blob/master/core/src/core/pmmr/pmmr.rs

import (
	"math/bits"
)

func pow2(x int) int {
	return 2 << (x - 1)
}

func log2(x int) int {
	return bits.Len(uint(x)) - 1
}

func leftChild(pos, h int) int {
	return pos - pow2(h)
}

func rightChild(pos, h int) int {
	return pos - 1
}

// FYI: sibling delta = pow2(h+1) - 1

func children(pos, h int) []int {
	if h == 0 {
		return nil
	}
	return []int{leftChild(pos, h), rightChild(pos, h)}
}

// peaks returns the index of peaks in an MMR of size n.
//
// Source: https://github.com/mimblewimble/grin/blob/78220febeda94595159ece675e77e26986a3c11d/core/src/core/pmmr/pmmr.rs#L402
func peaks(n int) (peaks []int) {
	p := 0 // partition (advances as we bag peaks)
	for n-p > 0 {
		nextPeak := pow2(log2(n-p+1)) - 1
		peaks = append(peaks, p+nextPeak-1)
		p += nextPeak
	}
	return peaks
}

// height returns the height (counting from 0) of the node at index n.
func height(pos int) int {
	switch {
	case pos < 0:
		panic("index cannot be negative")
	case pos == 0:
		return 0
	case pos < 2:
		return 0
	default:
		var upos = uint(pos)
		const allOnes = (1 << bits.UintSize) - 1
		var peakSize uint = allOnes >> bits.LeadingZeros(upos)
		for peakSize != 0 {
			if upos >= peakSize {
				upos -= peakSize
			}
			peakSize >>= 1
		}
		return int(upos)
	}
}

type pathEntry struct {
	pre  []int
	post []int
}

// path returns a sequence of nodes whose hashes comprise a proof of the node
// at pos.
func path(pos int) (p []pathEntry) {
	h := height(pos)
	cs := children(pos, h)
	if cs != nil {
		cpe := pathEntry{pre: make([]int, 0, len(cs))}
		for _, c := range cs {
			cpe.pre = append(cpe.pre, c)
		}
		p = append(p, cpe)
	}
	panic("not fully implemented")
}

// proof returns a set of node indexes needed to prove the inclusion of a
// summary of MMR of size from to a summary of MMR of size to.
func proof(from, to int) (nodes []int) {
	if to < from {
		panic("reverse proof requested")
	}
	fromPeaks := peaks(from)
	toPeaks := peaks(to)

	// skip matching peaks
	for len(fromPeaks) != 0 && len(toPeaks) != 0 {
		if fromPeaks[0] == toPeaks[0] {
			fromPeaks, toPeaks = fromPeaks[1:], toPeaks[1:]
		}
	}

	if len(fromPeaks) == 0 {
		// nothing to prove, that is fromPeaks is a subset of toPeaks
		return
	}

	// remaining fromPeaks should be fully contained under the first remaining
	// toPeak.
	toPeak := toPeaks[0]
	// c := toPeak
	for fromPeaks[0] != toPeak {
		// walk down left children until we find one, then walk down right children
		// i think...
	}
	return
}

func intSliceEqual(a, b []int) bool {
	switch {
	case a == nil && b == nil:
		return true
	case a == nil && len(b) == 0:
		return true
	case len(a) == 0 && b == nil:
		return true
	case a == nil || b == nil:
		return false
	case len(a) != len(b):
		return false
	default:
		for i := range a {
			if a[i] != b[i] {
				return false
			}
		}
		return true
	}
}
