package merkleweave

import (
	"crypto/rand"
	"encoding/base64"
	"sync"
	"testing"

	"github.com/vsekhar/merkleweave/internal/merkletree"
)

func fromString(s string) (r [merkletree.HashLength]byte) {
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	copy(r[:], b)
	return
}

func TestSummary(t *testing.T) {
	m := New()
	b1 := []byte{1, 2, 3, 4}
	m.Append(b1)
	s := m.Summary()
	good := newEmptySummary()
	prefixes := prefixesOf(b1)
	p1 := toInt(prefixes[0])
	p2 := toInt(prefixes[1])
	ts := merkletree.Summary{
		N:       1,
		Summary: fromString("HNYSpDRgx5Ujm7dHjpjiCu2oLrmq-htl0U4ByLazz10FhYeo9DZE1tANR9CgFdnl12TXikoQxyBmF-PS3woP0A"),
	}
	t.Log(base64.RawURLEncoding.EncodeToString(ts.Summary[:]))
	good.ss[p1] = ts
	good.ss[p2] = ts

	if !s.Equals(&good) {
		t.Errorf("expected %s, got %s", good.ShortString(), s.ShortString())
	}
}

func TestDuplicatePrefixes(t *testing.T) {
	m := New()
	b1 := []byte{1, 2, 1, 2}
	m.Append(b1)
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

type appendFunc func([]byte)

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
				a(d[k][:])
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkCompare(b *testing.B) {
	d := testData()

	t := merkletree.New()
	mu := sync.Mutex{}
	appendFunc := func(b []byte) {
		mu.Lock()
		defer mu.Unlock()
		t.Append(b)
	}
	b.Run("merkletree", func(b *testing.B) {
		benchmarkAppend(b, appendFunc, d)
	})

	w := New()
	b.Run("merkleweave", func(b *testing.B) {
		benchmarkAppend(b, w.Append, d)
	})
}
