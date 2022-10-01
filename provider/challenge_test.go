package provider

import (
	"testing"
)

func TestChallenge(t *testing.T) {
	challenge, verifier := NewChallenge(DifficultyMin, MultiplierMin)
	t.Logf("Got challenge %v", challenge)
	t.Logf("Got verifier %v", verifier)
	solution := challenge.solve()
	t.Logf("Got solution %v", solution)
	if !verifier.Verify(solution) {
		t.Fatal("Expected valid solution")
	}
	solution.Response[0]++
	if verifier.Verify(solution) {
		t.Fatal("Expected invalid solution [0]")
	}
	solution.Response[0]--
	solution.Response[3]++
	if verifier.Verify(solution) {
		t.Fatal("Expected invalid solution [1]")
	}
	solution.Response[3]--
	solution.Response = append(solution.Response, 0)
	if verifier.Verify(solution) {
		t.Fatal("Expected invalid solution [2]")
	}
	solution.Response = solution.Response[:4]
	solution.Id = "invalid"
	if verifier.Verify(solution) {
		t.Fatal("Expected invalid solution [3]")
	}
	// Hard challenge
	t.Log("Switching to hard challenge")
	challenge, verifier = NewChallenge(1, MultiplierMax)
	t.Logf("Got challenge %v", challenge)
	t.Logf("Got verifier %v", verifier)
	solution = challenge.solve()
	t.Logf("Got solution %v", solution)
	if !verifier.Verify(solution) {
		t.Fatal("Expected valid solution")
	}
	solution.Response[0]++
	if verifier.Verify(solution) {
		t.Fatal("Expected invalid solution [0]")
	}
	solution.Response[0]--
	solution.Response[3]++
	if verifier.Verify(solution) {
		t.Fatal("Expected invalid solution [1]")
	}
	solution.Response[3]--
	solution.Response = append(solution.Response, 0)
	if verifier.Verify(solution) {
		t.Fatal("Expected invalid solution [2]")
	}
	solution.Response = solution.Response[:4]
	solution.Id = "invalid"
	if verifier.Verify(solution) {
		t.Fatal("Expected invalid solution [3]")
	}
}

func BenchmarkChallenge(b *testing.B) {
	challenge, verifier := NewChallenge(10, 1000)
	solution := challenge.solve()
	b.Run("NewChallenge", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			NewChallenge(10, 1000)
		}
	})
	b.Run("solve", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			challenge.solve()
		}
	})
	b.Run("Verify", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			verifier.Verify(solution)
		}
	})
}
