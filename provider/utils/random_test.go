package utils

import (
	"bytes"
	"math/rand"
	"sort"
	"strconv"
	"testing"
)

/*
*Tests
 */

func TestSeedRandom(t *testing.T) {
	// Make sure we don't error
	if err := SeedRandom(); err != nil {
		t.Fatal(err)
	}
	// Make sure we're actually seeding
	rand.Seed(0)
	n := rand.Int()
	rand.Seed(0)
	if err := SeedRandom(); err != nil {
		t.Fatal(err)
	}
	if n == rand.Int() {
		t.Fatal("We're failing to seed math/rand")
	}
	// Make sure we're generating unique seeding
	if err := SeedRandom(); err != nil {
		t.Fatal(err)
	}
	n = rand.Int()
	if err := SeedRandom(); err != nil {
		t.Fatal(err)
	}
	if n == rand.Int() {
		t.Fatal("We're failing generate unique seeding")
	}
}

func TestRandomBytes(t *testing.T) {
	// Check if we're returning a correctly sized slice
	sizes := []int{0, 1, 10, 32, 64, 1024}
	for _, s := range sizes {
		if l := len(RandomBytes(s)); l != s {
			t.Errorf("Incorrectly sized slice returned for RandomBytes(%d) (%d)", s, l)
		}
	}
	// Check if we're using the global random seed
	rand.Seed(1)
	expected := RandomBytes(1024)
	rand.Seed(1)
	actual := RandomBytes(1024)
	if !bytes.Equal(expected, actual) {
		t.Error("We're not using the global random seed for bytes generation")
	}
	// Check if we're using the slow function
	rand.Seed(1)
	slow := make([]byte, 1024)
	slowRandomBytes(slow)
	rand.Seed(1)
	fast := RandomBytes(1024)
	if bytes.Equal(slow, fast) {
		t.Fatal("We're using slowRandomBytes when calling RandomBytes")
	}
	// Make sure slowRandomBytes handles zero-length slices
	for _, s := range sizes {
		if l := len(slowRandomBytes(make([]byte, s))); l != s {
			t.Errorf("Incorrectly sized slice returned for slowRandomBytes(make([]byte, %d)) (%d)", s, l)
		}
	}
}

func TestRandomString(t *testing.T) {
	// Check if we're returning a correctly sized string
	sizes := []int{0, 1, 10, 32, 64, 1024}
	for _, s := range sizes {
		if l := len(RandomString(s)); l != s {
			t.Errorf("Incorrectly sized string returned for RandomString(%d) (%d)", s, l)
		}
	}
	// Check if we're using the global random seed
	rand.Seed(1)
	expected := RandomString(1024)
	rand.Seed(1)
	actual := RandomString(1024)
	if expected != actual {
		t.Error("We're not using the global random seed for bytes generation")
	}
}

func TestRandomIntSlice(t *testing.T) {
	if len(RandomIntSlice(10, 0, 11)) != 0 {
		t.Error("Parsed invalid")
	}
	prev := []int{}
	for i := 0; i < 200; i++ {
		maxLen := 2 + rand.Intn(6)
		test := RandomIntSlice(10, 2, maxLen)
		if len(test) < 2 || len(test) > maxLen {
			t.Fatalf("%v out of bounds with max length %d\n", test, maxLen)
		}
		sort.Slice(test, func(a, b int) bool {
			return a < b
		})
		for i := 0; i < len(test)-1; i++ {
			if test[i] < 0 || 10 <= test[i] {
				t.Fatalf("%v has out of range value at %d", test, i)
			}
			if test[i] == test[i+1] {
				t.Fatalf("%v has repeating entries", test)
			}
		}
		if len(prev) == len(test) {
			match := true
			for i := 0; i < len(test); i++ {
				if test[i] != prev[i] {
					match = false
					break
				}
			}
			if match {
				t.Fatalf("%v matches the previous test %v", test, prev)
			}
		}
		prev = test
	}
}

/*
* Benchmarks
 */

func BenchmarkSeedRandom(b *testing.B) {
	b.Run("Sequential", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if err := SeedRandom(); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("Parallel", func(b *testing.B) {
		// This might error depending on the PRNG source
		// Normally it would be okay to just retry until successful
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				if err := SeedRandom(); err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}

func BenchmarkRandomBytes(b *testing.B) {
	sizes := []int{0, 1, 64, 1024}
	for _, s := range sizes {
		sizeString := strconv.Itoa(s) + "bytes"
		// Sequential
		b.Run("Sequential"+sizeString, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				RandomBytes(s)
			}
		})
		b.Run("SlowSequential"+sizeString, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				slowRandomBytes(make([]byte, s))
			}
		})
		// Parallel
		b.Run("Parallel"+sizeString, func(b *testing.B) {
			b.RunParallel(func(p *testing.PB) {
				for p.Next() {
					RandomBytes(s)
				}
			})
		})
		b.Run("SlowParallel"+sizeString, func(b *testing.B) {
			b.RunParallel(func(p *testing.PB) {
				for p.Next() {
					slowRandomBytes(make([]byte, s))
				}
			})
		})
	}
}

func BenchmarkRandomString(b *testing.B) {
	sizes := []int{0, 1, 64, 1024}
	for _, s := range sizes {
		sizeString := strconv.Itoa(s) + "len"
		// Sequential
		b.Run("Sequential"+sizeString, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				RandomString(s)
			}
		})
		// Parallel
		b.Run("Parallel"+sizeString, func(b *testing.B) {
			b.RunParallel(func(p *testing.PB) {
				for p.Next() {
					RandomString(s)
				}
			})
		})
	}
}

func BenchmarkRandomIntSlice(b *testing.B) {
	sizes := []int{1, 3}
	for _, s := range sizes {
		sizeString := strconv.Itoa(s) + "len"
		// Sequential
		b.Run("Sequential"+sizeString, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				RandomIntSlice(10, 1, s)
			}
		})
		// Parallel
		b.Run("Parallel"+sizeString, func(b *testing.B) {
			b.RunParallel(func(p *testing.PB) {
				for p.Next() {
					RandomIntSlice(10, 1, s)
				}
			})
		})
	}
}
