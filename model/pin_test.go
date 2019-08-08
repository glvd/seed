package model

import (
	"testing"

	"github.com/ipfs/go-cid"
)

type dummyResolved struct {
}

// Cid ...
func (d dummyResolved) Cid() cid.Cid {
	panic("implement me")
}

// Root ...
func (d dummyResolved) Root() cid.Cid {
	panic("implement me")
}

// Remainder ...
func (d dummyResolved) Remainder() string {
	panic("implement me")
}

// String ...
func (d dummyResolved) String() string {
	return "/ipld/QmS4ustL54uo8FzR9455qaxZwuMiUhyvMcX9Ba8nUH4uVv"
}

// Namespace ...
func (d dummyResolved) Namespace() string {
	panic("implement me")
}

// Mutable ...
func (d dummyResolved) Mutable() bool {
	panic("implement me")
}

// IsValid ...
func (d dummyResolved) IsValid() error {
	panic("implement me")
}

// TestPinHash ...
func TestPinHash(t *testing.T) {
	hash := PinHash(dummyResolved{})
	if hash != "QmS4ustL54uo8FzR9455qaxZwuMiUhyvMcX9Ba8nUH4uVv" {
		t.Error(hash, "QmS4ustL54uo8FzR9455qaxZwuMiUhyvMcX9Ba8nUH4uVv")
	}
}
