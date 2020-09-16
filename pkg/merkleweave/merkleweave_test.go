package merkleweave_test

import (
	"crypto/rand"
	"sync"
	"testing"

	"github.com/vsekhar/merkleweave/internal/merkletree"
	"github.com/vsekhar/merkleweave/pkg/merkleweave"
)

func TestDuplicatePrefixes(t *testing.T) {
	m := merkleweave.New()
	b1 := []byte{1, 2, 1, 2}
	if err := m.Append(b1); err != nil {
		t.Error(err)
	}
}

const records = 1 << 10 // 1024
const recordLen = 64

func testData() [records][recordLen]byte {
	var r [records][recordLen]byte
	for i := 0; i < records; i++ {
		rand.Read(r[i][:])
	}
	return r
}

type appendFunc func([]byte) error

func benchmarkAppend(b *testing.B, a appendFunc, d [records][recordLen]byte) {
	procs := 256
	N := b.N / procs
	b.ResetTimer()
	wg := sync.WaitGroup{}
	wg.Add(procs)
	for i := 0; i < procs; i++ {
		go func() {
			for j := 0; j < N; j++ {
				k := j % len(d)
				if err := a(d[k][:]); err != nil {
					b.Error(err)
					return
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkMerkleWeave(b *testing.B) {
	m := merkleweave.New()
	d := testData()
	benchmarkAppend(b, func(b []byte) error { return m.Append(b) }, d)
}

func BenchmarkMerkleTree(b *testing.B) {
	m := merkletree.New()
	mu := sync.Mutex{}
	appendFunc := func(b []byte) error {
		mu.Lock()
		defer mu.Unlock()
		return m.Append(b)
	}
	d := testData()
	benchmarkAppend(b, appendFunc, d)
}
