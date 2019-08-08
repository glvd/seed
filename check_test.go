package seed

import "testing"

func TestCheck(t *testing.T) {
	seed := NewSeed(Check("recursive"))
	seed.Start()
	seed.Wait()
}
