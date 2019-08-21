package seed_test

import (
	"testing"
	"time"

	"github.com/glvd/seed"
)

// TestNewSeed ...
func TestNewSeed(t *testing.T) {
	seed := seed.NewSeed()
	seed.Start()
	time.Sleep(3 * time.Second)
	seed.Wait()
	//"level":"info","ts":1566365331.8032362,"caller":"seed/seed.go:136","msg":"Seed starting"}
}
