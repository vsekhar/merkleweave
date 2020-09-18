package merkletree_test

import (
	"encoding/base64"
	"testing"

	"github.com/vsekhar/merkleweave/internal/merkletree"
)

// Keep this in sync with merkletree.hashLength (though we don't want to export
// that one and we don't want to do in-package tests here).
const hashLength = 64

func TestAppend(t *testing.T) {
	m := merkletree.New()
	if m.Len() != 0 {
		t.Fatal("expected empty merkletree")
	}
	b1 := []byte{1, 2, 3}
	m.Append(b1)
	if m.Len() != 1 {
		t.Error("expected len 1")
	}
}

func TestSummary(t *testing.T) {
	m := merkletree.New()
	b1 := []byte{1, 2, 3}
	m.Append(b1)
	n, s := m.Summary()
	if n != 1 {
		t.Errorf("expected length of %d, got %d", 1, n)
	}
	str := base64.RawURLEncoding.EncodeToString(s[:])
	good1 := "wixj8lSLFdmnBhbGlJxYaCiN1SwMNmV-G7h3g3Yox6ZV8vFHlgRl61rf1y_2XVx7YPTYgSaQaOc1uAk4-P7b4A"
	if str != good1 {
		t.Errorf("expected %#v, got %#v", good1, str)
	}

	times := 100
	for i := 0; i < times; i++ {
		m.Append(b1)
	}
	n, s = m.Summary()
	if n != times+1 {
		t.Errorf("expected length of %d, got %d", times*len(b1), n)
	}
	str = base64.RawURLEncoding.EncodeToString(s[:])
	good101 := "feJxwpLst4bh-4prEMa-Xcy6R6Tdk9w7sbmseq-goqzvJ_1PkmE5EjadvOD1L4SrY04nYyPM7yyWMRkZkumUWw"
	if str != good101 {
		t.Errorf("expected %#v, got %#v", good101, str)
	}
}
