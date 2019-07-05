package old

import "testing"
import _ "github.com/mattn/go-sqlite3"

// TestLoadFrom ...
func TestLoadFrom(t *testing.T) {
	t.Logf("%+v", LoadFrom("seed.db"))
}
