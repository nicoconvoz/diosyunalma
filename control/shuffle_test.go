package control

import (
	"math/rand"
	"reflect"
	"sort"
	"testing"
)

// TestShuffleGapsPreservesTheMultiset is the property that makes the whole
// control valid. The decoy must carry the exact same distances as the real
// sequence — only their order may change. If the multiset drifted, any
// difference we later measure could be blamed on the distribution rather than
// on the arrangement, and the experiment would prove nothing.
func TestShuffleGapsPreservesTheMultiset(t *testing.T) {
	gaps := []int{1, 2, 2, 4, 2, 4, 2, 4, 6, 2, 6, 4, 2, 4, 6, 6}
	rng := rand.New(rand.NewSource(1))

	for trial := 0; trial < 50; trial++ {
		got := ShuffleGaps(gaps, rng)

		if len(got) != len(gaps) {
			t.Fatalf("len = %d, want %d", len(got), len(gaps))
		}

		gotSorted := append([]int(nil), got...)
		wantSorted := append([]int(nil), gaps...)
		sort.Ints(gotSorted)
		sort.Ints(wantSorted)

		if !reflect.DeepEqual(gotSorted, wantSorted) {
			t.Fatalf("multiset changed: got %v, want %v", gotSorted, wantSorted)
		}
	}
}

// TestShuffleGapsLeavesTheInputAlone guards against the classic in-place bug:
// a decoy that mutates the real data destroys the very thing it is compared to.
func TestShuffleGapsLeavesTheInputAlone(t *testing.T) {
	gaps := []int{1, 2, 2, 4, 2, 4, 2, 4, 6}
	original := append([]int(nil), gaps...)

	ShuffleGaps(gaps, rand.New(rand.NewSource(1)))

	if !reflect.DeepEqual(gaps, original) {
		t.Errorf("input was mutated: got %v, want %v", gaps, original)
	}
}

// TestShuffleGapsIsReproducibleFromASeed keeps experiments repeatable: the same
// seed must replay the same decoys, or a published result cannot be checked.
func TestShuffleGapsIsReproducibleFromASeed(t *testing.T) {
	gaps := []int{1, 2, 2, 4, 2, 4, 2, 4, 6, 2, 6, 4}

	first := ShuffleGaps(gaps, rand.New(rand.NewSource(99)))
	second := ShuffleGaps(gaps, rand.New(rand.NewSource(99)))

	if !reflect.DeepEqual(first, second) {
		t.Errorf("same seed gave different decoys: %v vs %v", first, second)
	}
}

func TestShuffleGapsHandlesTinyInputs(t *testing.T) {
	rng := rand.New(rand.NewSource(1))

	if got := ShuffleGaps(nil, rng); len(got) != 0 {
		t.Errorf("ShuffleGaps(nil) = %v, want empty", got)
	}
	if got := ShuffleGaps([]int{6}, rng); !reflect.DeepEqual(got, []int{6}) {
		t.Errorf("ShuffleGaps([6]) = %v, want [6]", got)
	}
}
