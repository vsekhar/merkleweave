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
	b1 := [hashLength]byte{1, 2, 3} // rest are zeros
	if err := m.Append(b1); err != nil {
		t.Error(err)
	}
	if m.Len() != 1 {
		t.Error("expected len 1")
	}
}

func TestSummary(t *testing.T) {
	m := merkletree.New()
	b1 := [hashLength]byte{1, 2, 3} // rest are zeros
	if err := m.Append(b1); err != nil {
		t.Fatal(err)
	}
	s, err := m.Summary()
	if err != nil {
		t.Fatal(err)
	}
	str := base64.RawURLEncoding.EncodeToString(s[:])
	good1 := "q3dir6rpvXpwRwZynBycjmAmgtFVDRErpqUA9o-KKGLS9150t-smvUMFXlvd8u8URHoKaybZAwmHS2PzmiYZLg"
	if str != good1 {
		t.Errorf("expected %#v, got %#v", good1, str)
	}

	for i := 0; i < 100; i++ {
		if err := m.Append(b1); err != nil {
			t.Fatal(err)
		}
	}
	s, err = m.Summary()
	if err != nil {
		t.Fatal(err)
	}
	str = base64.RawURLEncoding.EncodeToString(s[:])
	good101 := "Tbr2bI4gBhAdwb0KLgt677KMS2-WUGpWuCdaKlT_SNlxLVEte2WjkpCwfe9HxaC6vYsuLqQ5-ac7n7HiuARkLg"
	if str != good101 {
		t.Errorf("expected %#v, got %#v", good101, str)
	}
}
