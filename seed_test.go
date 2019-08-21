package seed_test

import (
	"testing"

	"github.com/glvd/seed"
)

// TestNewSeed ...
func TestNewSeed(t *testing.T) {
	seed := seed.NewSeed()
	seed.Start()
	seed.Wait()

}
