// Command echoconst applies one-half to infinity at k = 12.
//
// The speed law (Finding 81) says the gap-12 deviation is a constant core
// plus a transient dying at the square-root rate. That law is a telescope
// pointed at infinity: because the transient's death rate is known (the ½),
// the INFINITE-x value of the echo constant is computable from finite data,
// with honest error bars — the special eye, applied.
//
// This run fits dev₁₂(x) = C∞ + b·x^(−1/2) by weighted least squares to the
// four measured decades, reports C∞ ± σ, judges the closed-form ballot at
// infinity, and PRE-REGISTERS the 10¹² reading before the walker returns.
//
// Ballot status: 46/45 executed (Finding 65). 49/48 pre-registered and
// weakened. 48/47 is named HERE for the first time — post hoc, flagged as
// such, admissible only for FUTURE tests. The null always runs.
//
// Usage:
//
//	go run ./cmd/echoconst
package main

import (
	"fmt"
	"math"
)

var (
	devs   = [4]float64{1.92, 2.20, 2.21, 2.13}
	sigmas = [4]float64{0.43, 0.146, 0.05, 0.017}
	roots  = [4]float64{1e-4, 3.1623e-5, 1e-5, 3.1623e-6}
)

func main() {
	var sw, swx, swy, swxx, swxy float64
	for i := 0; i < 4; i++ {
		w := 1 / (sigmas[i] * sigmas[i])
		sw += w
		swx += w * roots[i]
		swy += w * devs[i]
		swxx += w * roots[i] * roots[i]
		swxy += w * roots[i] * devs[i]
	}
	den := sw*swxx - swx*swx
	b := (sw*swxy - swx*swy) / den
	c := (swy*swxx - swx*swxy) / den
	sigC := math.Sqrt(swxx / den)
	sigB := math.Sqrt(sw / den)

	fmt.Println("ONE-HALF APPLIED TO INFINITY AT K = 12")
	fmt.Printf("\n  the echo constant at infinity: C = %+.4f%% +/- %.4f\n", c, sigC)
	fmt.Printf("  the dying transient:           b = %+.0f +/- %.0f (gone as 1/sqrt x)\n", b, sigB)

	fmt.Println("\n  the ballot, judged at infinity:")
	type cand struct {
		name string
		v    float64
		note string
	}
	ballot := []cand{
		{"46/45", 100.0 / 45, "pre-registered; executed at 10^11"},
		{"49/48", 100.0 / 48, "pre-registered; weakened"},
		{"48/47", 100.0 / 47, "NAMED POST HOC - future tests only"},
	}
	for _, k := range ballot {
		z := (c - k.v) / sigC
		fmt.Printf("   %-6s = %+.4f%%   distance %+5.1f sigma   (%s)\n", k.name, k.v, z, k.note)
	}
	fmt.Println("   null (no simple closed form): always on the ballot")

	r12 := 1e-6
	pred := c + b*r12
	lo, hi := pred-3*sigC, pred+3*sigC
	fmt.Printf("\n  PRE-REGISTERED for the walker still crossing 10^12:\n")
	fmt.Printf("  dev12(10^12) must land in [%+.2f%%, %+.2f%%] - or the half-law dies at 12.\n", lo, hi)
	fmt.Println("\n  the transient's death rate is known, so infinity is visible from here:")
	fmt.Println("  the eye does not wait for the wave - it reads the score.")
}
