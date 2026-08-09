package pattern

import "testing"

func TestGilbreathRowsOnHandExamples(t *testing.T) {
	tests := []struct {
		name     string
		seq      []int
		maxRows  int
		wantOK   int
		wantFail int
	}{
		{
			// [2 3 5] -> [1 2] -> [1]: both rows open with 1.
			name: "tiny surviving sequence", seq: []int{2, 3, 5},
			maxRows: 10, wantOK: 2, wantFail: -1,
		},
		{
			// [2 3 7] -> [1 4] -> [3]: the second row opens with 3.
			name: "breaks in row two", seq: []int{2, 3, 7},
			maxRows: 10, wantOK: 1, wantFail: 2,
		},
		{
			// [2 4]: the first difference is 2, broken immediately.
			name: "breaks in row one", seq: []int{2, 4},
			maxRows: 10, wantOK: 0, wantFail: 1,
		},
		{
			name: "too short to difference", seq: []int{7},
			maxRows: 10, wantOK: 0, wantFail: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, fail := GilbreathRows(tt.seq, tt.maxRows)
			if ok != tt.wantOK || fail != tt.wantFail {
				t.Errorf("GilbreathRows(%v) = (%d, %d), want (%d, %d)",
					tt.seq, ok, fail, tt.wantOK, tt.wantFail)
			}
		})
	}
}

// TestGilbreathRowsOnRealPrimes records the conjecture's content on a range
// where it is verified fact: the primes below 1000 survive every row asked of
// them.
func TestGilbreathRowsOnRealPrimes(t *testing.T) {
	primes := []int{2, 3}
	for n := 5; n < 1000; n += 2 {
		isP := true
		for d := 3; d*d <= n; d += 2 {
			if n%d == 0 {
				isP = false
				break
			}
		}
		if isP {
			primes = append(primes, n)
		}
	}

	ok, fail := GilbreathRows(primes, 60)
	if fail != -1 || ok != 60 {
		t.Errorf("primes below 1000: (%d, %d), want (60, -1)", ok, fail)
	}
}

func TestGilbreathRowsRespectsMaxRows(t *testing.T) {
	ok, fail := GilbreathRows([]int{2, 3, 5, 7, 11, 13}, 2)
	if ok != 2 || fail != -1 {
		t.Errorf("capped run = (%d, %d), want (2, -1)", ok, fail)
	}
}
