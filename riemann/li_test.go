package riemann

import (
	"math"
	"math/cmplx"
	"testing"
)

func pair(beta, gamma float64) []complex128 {
	return []complex128{complex(beta, gamma), complex(beta, -gamma)}
}

// TestLiSumFirstCoefficientByHand checks λ₁ against the closed form: for n=1
// the term is 1 − (1 − 1/ρ) = 1/ρ, so a conjugate pair contributes
// 2·Re(1/ρ) = 2β/(β²+γ²).
func TestLiSumFirstCoefficientByHand(t *testing.T) {
	const beta, gamma = 0.5, 14.134725
	want := 2 * beta / (beta*beta + gamma*gamma)

	got := LiSum(pair(beta, gamma), 1)
	if math.Abs(got-want) > 1e-14 {
		t.Errorf("LiSum(pair, 1) = %v, want %v", got, want)
	}
}

// TestLiSumOnTheCriticalLineIsACosine pins the bridge this package exists for:
// with β = 1/2 the Möbius image z = 1 − 1/ρ has |z| = 1, so a conjugate pair
// contributes exactly 2(1 − cos nθ) — nonnegative for every n, automatically.
func TestLiSumOnTheCriticalLineIsACosine(t *testing.T) {
	zeros := pair(0.5, 14.134725)
	z := 1 - 1/zeros[0]

	if r := cmplx.Abs(z); math.Abs(r-1) > 1e-15 {
		t.Fatalf("|1 − 1/ρ| = %v on the critical line, want exactly 1", r)
	}

	theta := cmplx.Phase(z)
	for n := 1; n <= 200; n++ {
		got := LiSum(zeros, n)
		want := 2 * (1 - math.Cos(float64(n)*theta))

		if math.Abs(got-want) > 1e-9 {
			t.Fatalf("n=%d: LiSum = %v, cosine form = %v", n, got, want)
		}
		if got < -1e-12 {
			t.Fatalf("n=%d: LiSum = %v is negative on the critical line", n, got)
		}
	}
}

// TestLiSumCrashesForAnOffLineZero is the detector property. A zero off the
// line, completed to the quadruple the functional equation demands, puts one
// Möbius image outside the unit circle — and (1−zⁿ) then blows up, driving
// some λₙ far negative. This is why Li positivity IS the Riemann Hypothesis.
func TestLiSumCrashesForAnOffLineZero(t *testing.T) {
	// The quadruple ρ, ρ̄, 1−ρ, 1−ρ̄ for a fictitious zero at β = 0.9, γ = 2.
	offLine := append(pair(0.9, 2), pair(0.1, 2)...)

	worst := math.Inf(1)
	for n := 1; n <= 150; n++ {
		if v := LiSum(offLine, n); v < worst {
			worst = v
		}
	}
	if worst > -5 {
		t.Errorf("min λₙ = %v for an off-line quadruple, want a crash below -5", worst)
	}

	// The same heights ON the line never go negative.
	onLine := append(pair(0.5, 2), pair(0.5, 2)...)
	for n := 1; n <= 150; n++ {
		if v := LiSum(onLine, n); v < -1e-9 {
			t.Fatalf("n=%d: on-line λₙ = %v went negative", n, v)
		}
	}
}

func TestLiSumHandlesDegenerateInput(t *testing.T) {
	if got := LiSum(nil, 5); got != 0 {
		t.Errorf("no zeros = %v, want 0", got)
	}
	if got := LiSum(pair(0.5, 14), 0); got != 0 {
		t.Errorf("n=0 = %v, want 0 (z⁰ = 1 kills every term)", got)
	}
}
