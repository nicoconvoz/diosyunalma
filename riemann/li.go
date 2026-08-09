// Package riemann holds the numeric bridge between Li's criterion and the
// Hilbert–Pólya program.
//
// The bridge is one line of algebra: |1 − 1/ρ| = 1 exactly when Re ρ = 1/2,
// so the Möbius map z = 1 − 1/ρ carries the critical line to the unit circle.
// A Hilbert–Pólya operator — self-adjoint, spectrum on a real line — becomes,
// through the Cayley transform, a unitary operator with spectrum on that
// circle. Li's coefficients are then traces of its powers, λₙ = Tr(I − Uⁿ),
// and each conjugate pair contributes 2(1 − cos nθ): nonnegative for free.
//
// Li positivity for every n is equivalent to the Riemann Hypothesis
// (Li 1997); the spectral reading is due to Bombieri and Lagarias (1999).
// This package computes the partial sums so measured zeros can be fed in.
package riemann

import "math/cmplx"

// LiSum returns Σ_ρ Re[1 − (1 − 1/ρ)ⁿ] over the given zeros.
//
// The list should be closed under conjugation, which makes the true sum real;
// taking the real part keeps floating-point residue out of the result. Feeding
// a finite set gives a PARTIAL coefficient: useful for watching the mechanism,
// never a substitute for the full sum over all zeros.
//
// The detector property, tested rather than assumed: a zero on the critical
// line contributes a bounded, nonnegative oscillation for every n, while a
// zero off the line — completed to the quadruple the functional equation
// demands — sends some λₙ exponentially negative.
func LiSum(zeros []complex128, n int) float64 {
	if n < 1 {
		return 0
	}

	sum := 0.0
	power := complex(float64(n), 0)
	for _, rho := range zeros {
		z := 1 - 1/rho
		sum += real(1 - cmplx.Pow(z, power))
	}

	return sum
}
