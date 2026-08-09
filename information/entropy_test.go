package information

import (
	"math"
	"testing"
)

func closeTo(t *testing.T, got, want float64, label string) {
	t.Helper()
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("%s = %v, want %v", label, got, want)
	}
}

func TestShannonOnKnownDistributions(t *testing.T) {
	tests := []struct {
		name   string
		counts []int
		want   float64
	}{
		{name: "a certain outcome carries no information", counts: []int{7}, want: 0},
		{name: "a fair coin is one bit", counts: []int{50, 50}, want: 1},
		{name: "four equal outcomes are two bits", counts: []int{3, 3, 3, 3}, want: 2},
		{name: "eight equal outcomes are three bits", counts: []int{1, 1, 1, 1, 1, 1, 1, 1}, want: 3},
		{name: "zero counts are ignored", counts: []int{5, 0, 5, 0}, want: 1},
		{name: "no observations at all", counts: []int{}, want: 0},
		{name: "all zero counts", counts: []int{0, 0}, want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			closeTo(t, Shannon(tt.counts), tt.want, "Shannon")
		})
	}
}

// TestShannonIsMaximisedByTheUniformDistribution states the bound every budget
// number is read against: no distribution over n outcomes can exceed log2(n).
func TestShannonIsMaximisedByTheUniformDistribution(t *testing.T) {
	uniform := Shannon([]int{10, 10, 10, 10, 10})
	skewed := Shannon([]int{46, 1, 1, 1, 1})

	if uniform <= skewed {
		t.Errorf("uniform entropy %v should exceed skewed %v", uniform, skewed)
	}
	if want := math.Log2(5); math.Abs(uniform-want) > 1e-9 {
		t.Errorf("uniform entropy = %v, want log2(5) = %v", uniform, want)
	}
}

func TestConditionalEntropyAtOrderZeroIsTheMarginal(t *testing.T) {
	seq := []int{0, 1, 1, 0, 1, 0, 0, 1, 1, 1}

	counts := make([]int, 2)
	for _, v := range seq {
		counts[v]++
	}

	closeTo(t, ConditionalEntropy(seq, 0), Shannon(counts), "ConditionalEntropy order 0")
}

// TestConditionalEntropyOnAPerfectlyPredictableChain is the sharp case: an
// alternating sequence is fully determined by its previous symbol, so knowing
// one step of history removes every bit of uncertainty.
func TestConditionalEntropyOnAPerfectlyPredictableChain(t *testing.T) {
	seq := make([]int, 200)
	for i := range seq {
		seq[i] = i % 2
	}

	closeTo(t, ConditionalEntropy(seq, 0), 1, "order 0 of an alternating chain")
	closeTo(t, ConditionalEntropy(seq, 1), 0, "order 1 of an alternating chain")
	closeTo(t, ConditionalEntropy(seq, 2), 0, "order 2 of an alternating chain")
}

// TestConditionalEntropyNeedsTwoStepsForAPeriodOfFour checks that the order
// reported is the order actually required: a period-4 pattern is not settled by
// one symbol of history but is by two.
func TestConditionalEntropyNeedsTwoStepsForAPeriodOfFour(t *testing.T) {
	pattern := []int{0, 0, 1, 1}
	seq := make([]int, 400)
	for i := range seq {
		seq[i] = pattern[i%4]
	}

	if got := ConditionalEntropy(seq, 1); got < 0.4 {
		t.Errorf("order 1 = %v, want appreciable uncertainty left", got)
	}
	closeTo(t, ConditionalEntropy(seq, 2), 0, "order 2 of a period-4 chain")
}

func TestConditionalEntropyOnAConstantChainIsZero(t *testing.T) {
	seq := make([]int, 50)

	for order := 0; order <= 3; order++ {
		closeTo(t, ConditionalEntropy(seq, order), 0, "constant chain")
	}
}

func TestConditionalEntropyHandlesDegenerateInput(t *testing.T) {
	closeTo(t, ConditionalEntropy(nil, 0), 0, "nil sequence")
	closeTo(t, ConditionalEntropy([]int{1, 2, 3}, 5), 0, "order beyond the sequence")
	closeTo(t, ConditionalEntropy([]int{1, 2, 3}, -1), 0, "negative order")
}

// TestMutualInformationIsTheEntropyDrop pins the relation the budget is built
// on: what the past tells you about the next symbol is exactly the uncertainty
// it removes.
func TestMutualInformationIsTheEntropyDrop(t *testing.T) {
	seq := []int{0, 1, 1, 0, 1, 0, 0, 1, 1, 1, 0, 1, 0, 0, 1, 0}

	for order := 1; order <= 3; order++ {
		want := ConditionalEntropy(seq, 0) - ConditionalEntropy(seq, order)
		closeTo(t, MutualInformation(seq, order), want, "MutualInformation")
	}
}

func TestMutualInformationIsZeroWhenThePastSaysNothing(t *testing.T) {
	seq := make([]int, 40)

	closeTo(t, MutualInformation(seq, 1), 0, "constant chain carries no news")
	closeTo(t, MutualInformation([]int{5}, 1), 0, "a single symbol")
}
