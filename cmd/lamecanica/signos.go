package main

// signos.go - THE SIGN CHANNEL, and the two gauge theorems that decide what is
// and is not reachable with a per-site state.
//
// Phase VIII's kernel was strictly positive, so the medium had no TENSION in the
// true sense: every bond pulled the same way. A bond sign field supplies it at
// exactly zero parameters, and because (sign)^2 = 1 the total force is untouched:
// only the material changes, and H stays real symmetric so R1 survives.

import "math"

// --- THEOREM 1 (gauge): a per-site sign is invisible ------------------------
//
// If sg_ij = c_i*c_j with c_i in {+1,-1}, then with D = diag(c) we have
// H' = D H D and D is real orthogonal (D = D^T = D^-1). Conjugation by an
// orthogonal matrix leaves the spectrum identical, and |v'_i| = |v_i| leaves the
// participation ratio identical. So Sigma^2, alpha, vivos and PR/N are unchanged
// BIT FOR BIT. This is a theorem, hence a unit test: signoGauge must reproduce
// signoLlano to 1e-12, and a failure means the implementation is broken.
//
// The same argument with D = diag(exp(-i*theta)) shows a per-site PHASE is also
// pure gauge. That settles Phase VIII's open question 3 in the NEGATIVE for this
// route: arithmetic cannot enter as a SITE phase, because a site phase moves no
// observable at all. The complex case is a proof rather than a run - it needs no
// solver, which is precisely why it is worth stating.
//
// --- THEOREM 2 (coboundary): a bond sign of gradient form is invisible ------
//
// A bond phase of the form theta_ij = u_i - u_j is a coboundary and is the same
// gauge transformation, so it moves nothing either. Real shear needs a bond
// connection that is NOT a coboundary, i.e. one with nonzero holonomy around a
// loop; torsion is the curl of that connection.
//
// A CORRECTION TO THE AUDITOR'S SECTION 8, offered because it cuts against what a
// lab wants to hear: torsion does NOT need a 2-D lattice. With kmax = 120 the
// coupling graph is nothing like a path - every triple i<j<k with k-i <= 120 is a
// triangle, so the graph is dense in loops. The obstruction is COBOUNDARY-NESS,
// not dimensionality. A nearest-neighbour chain would be a tree and torsion would
// then be genuinely impossible; the long-range kernel Phase VII selected removed
// that excuse. Building it means O(N*kmax) ~ 48000 bond variables instead of 400,
// so it is not a minimum mechanism and is declared out of scope, not faked.

// signoLlano is S0: Phase VIII exactly.
func signoLlano(i, j int) float64 { return 1 }

// signoGauge is S1: a pure gauge. THE UNIT TEST - must return S0.
func signoGauge(c []float64) func(int, int) float64 {
	return func(i, j int) float64 { return c[i] * c[j] }
}

// signoFrustrado is S2: sg_ij = (-1)^(u_i*u_j). This is NOT a gauge - it cannot be
// written c_i*c_j, because c_i*c_j would force a parity relation on every triple
// that (-1)^(u_i u_j) violates: with u_i = u_j = u_k = 1 the three bonds are all
// -1, whose product is -1, while any c_i*c_j product around a triangle is
// (c_i c_j)(c_j c_k)(c_k c_i) = +1. Half of the all-ones triples are frustrated.
func signoFrustrado(u []float64) func(int, int) float64 {
	return func(i, j int) float64 {
		if u[i] > 0 && u[j] > 0 {
			return -1
		}
		return 1
	}
}

// bitsAzar draws u_i in {0,1} at a GIVEN density, so the arithmetic arm can be
// compared against random signs at MATCHED density - the control that decides
// whether the sign channel is arithmetic or merely frustrated.
func bitsAzar(n int, dens float64, d *dado) []float64 {
	u := make([]float64, n)
	for i := range u {
		if d.u() < dens {
			u[i] = 1
		}
	}
	return u
}

func signosAzar(n int, d *dado) []float64 {
	c := make([]float64, n)
	for i := range c {
		c[i] = 1
		if d.u() < 0.5 {
			c[i] = -1
		}
	}
	return c
}

// bitsReciprocidad is S3: b_i = 1 iff p_i = 3 mod 4. Declared BEFORE any spectrum,
// derived from the primes alone: R6 clean. Returns the bits and their density, so
// the random control can be matched to it exactly.
func bitsReciprocidad(ms []modo) ([]float64, float64) {
	u := make([]float64, len(ms))
	k := 0
	for i, m := range ms {
		if m.p%4 == 3 {
			u[i] = 1
			k++
		}
	}
	return u, float64(k) / float64(len(ms))
}

// frustracion counts, over a sample of triangles of the coupling graph, how many
// have negative sign product. A gauge field gives exactly zero.
func frustracion(sg func(int, int) float64, n, kmax int, d *dado) float64 {
	tot, fr := 0, 0
	for t := 0; t < 20000; t++ {
		i := int(d.u() * float64(n))
		j := i + 1 + int(d.u()*float64(kmax-1))
		k := j + 1 + int(d.u()*float64(kmax-1))
		if k >= n || k-i > kmax {
			continue
		}
		tot++
		if sg(i, j)*sg(j, k)*sg(i, k) < 0 {
			fr++
		}
	}
	if tot == 0 {
		return 0
	}
	return float64(fr) / float64(tot)
}

// --- the rulers, computed and never quoted alone ----------------------------
//
// Control 16: the zeros' target must always be printed next to the random-matrix
// floors, because the target lies BELOW both. No symmetry-class engineering can
// reach it, and a model drifting toward the GOE floor has become a structureless
// random matrix - a LOSS of arithmetic content dressed as a win.

// pisoGUE is the universal number variance of the Gaussian Unitary Ensemble,
// Sigma^2(L) = (1/pi^2)(ln(2 pi L) + gamma + 1).
func pisoGUE(L float64) float64 {
	const g = 0.5772156649015329
	return (math.Log(2*math.Pi*L) + g + 1) / (math.Pi * math.Pi)
}

// pisoGOE is the same for the Gaussian Orthogonal Ensemble,
// Sigma^2(L) = (2/pi^2)(ln(2 pi L) + gamma + 1 - pi^2/8).
func pisoGOE(L float64) float64 {
	const g = 0.5772156649015329
	return 2 * (math.Log(2*math.Pi*L) + g + 1 - math.Pi*math.Pi/8) / (math.Pi * math.Pi)
}
