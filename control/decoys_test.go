package control

import (
	"math/rand"
	"testing"
)

func TestOddDecoyShape(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	got := OddDecoy(100, 25, rng)

	if len(got) != 25 {
		t.Fatalf("len = %d, want 25", len(got))
	}
	if got[0] != 2 {
		t.Errorf("first element = %d, want 2 (the sole even prime is always kept)", got[0])
	}
	for i := 1; i < len(got); i++ {
		if got[i]%2 == 0 {
			t.Errorf("element %d = %d, want odd", i, got[i])
		}
		if got[i] <= got[i-1] {
			t.Errorf("not strictly increasing at %d: %d then %d", i, got[i-1], got[i])
		}
		if got[i] > 100 {
			t.Errorf("element %d = %d exceeds the limit", i, got[i])
		}
	}
}

func TestOddDecoyIsReproducibleFromASeed(t *testing.T) {
	a := OddDecoy(1000, 50, rand.New(rand.NewSource(9)))
	b := OddDecoy(1000, 50, rand.New(rand.NewSource(9)))

	for i := range a {
		if a[i] != b[i] {
			t.Fatalf("same seed diverged at %d: %d vs %d", i, a[i], b[i])
		}
	}
}

// TestCramerDecoyDensity checks the model the decoy implements: including n
// with probability 1/ln n reproduces the prime-counting density, so the count
// below 10^5 must land near pi(10^5) = 9592.
func TestCramerDecoyDensity(t *testing.T) {
	got := CramerDecoy(100_000, rand.New(rand.NewSource(7)))

	if got[0] != 2 {
		t.Errorf("first element = %d, want 2", got[0])
	}
	for i := 1; i < len(got); i++ {
		if got[i] <= got[i-1] {
			t.Fatalf("not strictly increasing at %d", i)
		}
	}
	if n := len(got); n < 8_800 || n > 10_400 {
		t.Errorf("count = %d, want near pi(1e5) = 9592", n)
	}
}

func TestDecoysHandleDegenerateInput(t *testing.T) {
	rng := rand.New(rand.NewSource(1))

	if got := OddDecoy(100, 0, rng); len(got) != 0 {
		t.Errorf("count 0 = %v, want empty", got)
	}
	if got := CramerDecoy(2, rng); len(got) != 1 || got[0] != 2 {
		t.Errorf("limit 2 = %v, want [2]", got)
	}
}
