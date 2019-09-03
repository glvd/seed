package seed_test

import (
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/glvd/seed"
	"github.com/glvd/seed/model"
	"github.com/godcong/go-trait"
)

// TestCheck ...
func TestCheck(t *testing.T) {
	data := seed.NewDatabase(model.MustDatabase(model.InitSQLite3("test.db")), seed.DatabaseShowSQLArg())

	check := seed.NewCheck(seed.CheckPinTypeArg("recursive"))

	s := seed.NewSeed(data)

	s.Start()

	s.RunTask(seed.NewTask(check))

	s.Wait()

}

// TestSliceBenchmark ...
func BenchmarkSliceBenchmark(b *testing.B) {
	var vals []interface{}
	for i := 0; i < 250000; i++ {
		vals = append(vals, trait.GenerateRandomString(32))
		vals = append(vals, rand.Intn(time.Now().Nanosecond()))
	}
	b.StartTimer()
	count := 0
	for i := 0; i < 100; i++ {
		if checkSlice(trait.GenerateRandomString(3), vals...) {
			count++
		}
	}
	b.StopTimer()
	b.Log(count)
}

// TestSliceBenchmark ...
func BenchmarkArrayBenchmark(b *testing.B) {
	var vals []interface{}
	for i := 0; i < 250000; i++ {
		vals = append(vals, trait.GenerateRandomString(32))
		vals = append(vals, rand.Intn(time.Now().Nanosecond()))
	}
	b.StartTimer()
	count := 0
	for i := 0; i < 100; i++ {
		if checkArray(trait.GenerateRandomString(3), vals) {
			count++
		}
	}
	b.StopTimer()
	b.Log(count)
}

func checkSlice(s string, v ...interface{}) bool {
	for i := range v {
		if vv, b := (v[i]).(string); b {
			if strings.Index(vv, s) > 0 {
				return true
			}
		}
	}
	return false
}

func checkArray(s string, v []interface{}) bool {
	if v == nil {
		return false
	}
	for i := range v {
		if vv, b := (v[i]).(string); b {
			if strings.Index(vv, s) > 0 {
				return true
			}
		}
	}
	return false
}
