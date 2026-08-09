// Command speeds seeks the private speed law of the climb.
//
// Finding 66 measured every gap climbing toward zero at its own pace;
// Finding 75's impostor proved the fading valleys are shallow (a wheel of
// coins fakes them) while gap 12's invariant is deep. The hypothesis that
// unifies both, registered here: every gap's deviation is
//
//	dev(d, x) = core(d) + transient(d) · x^(−1/2)
//
// — a persistent arithmetic core plus a shallow transient dying at the
// square-root rate. The private speeds differ only through the mix.
//
// PRE-REGISTERED: the fit lands core(12) in [+1.9, +2.4] (the invariant);
// the impostor-faked valleys 30 and 54 get cores within ±1 of zero; the
// residuals stay within a few sigma of the recorded errors. The open
// question the fit answers fresh: whether 36 and 48 keep genuinely
// nonzero negative cores, or die entirely.
//
// Usage:
//
//	go run ./cmd/speeds
package main

import (
	"fmt"
	"math"
)

// deviations (%) at 10^8..10^11 (Findings 54, 63, 65) and their sigmas.
type row struct {
	d   int
	dev [4]float64
	sig [4]float64
}

var table = []row{
	{6, [4]float64{-0.97, -0.69, -0.53, -0.41}, [4]float64{0.19, 0.066, 0.02, 0.008}},
	{12, [4]float64{1.92, 2.20, 2.21, 2.13}, [4]float64{0.43, 0.146, 0.05, 0.017}},
	{18, [4]float64{-0.70, 0.14, 0.47, 0.69}, [4]float64{0.63, 0.23, 0.07, 0.024}},
	{24, [4]float64{-0.47, -0.10, 0.11, 0.57}, [4]float64{0.94, 0.33, 0.11, 0.035}},
	{30, [4]float64{-6.08, -2.10, -1.11, -0.50}, [4]float64{1.12, 0.35, 0.11, 0.034}},
	{36, [4]float64{-12.13, -12.51, -10.72, -9.39}, [4]float64{2.21, 0.64, 0.19, 0.060}},
	{42, [4]float64{1.92, -0.03, -0.37, 0.55}, [4]float64{2.66, 0.82, 0.24, 0.071}},
	{48, [4]float64{-7.12, -9.07, -6.03, -4.35}, [4]float64{4.14, 1.32, 0.37, 0.108}},
	{54, [4]float64{-18.92, -5.77, -2.34, -1.03}, [4]float64{5.85, 1.87, 0.51, 0.143}},
	{60, [4]float64{-9.24, -7.51, -4.86, -3.80}, [4]float64{6.19, 1.94, 0.50, 0.136}},
}

var roots = [4]float64{1e-4, 3.1623e-5, 1e-5, 3.1623e-6}

func main() {
	fmt.Println("THE SPEED LAW — core + transient/sqrt(x), fitted to four decades")
	fmt.Println("\n   d    core (%)   transient    worst residual / sigma")
	for _, r := range table {
		// weighted least squares on dev = core + b*root.
		var sw, swx, swy, swxx, swxy float64
		for i := 0; i < 4; i++ {
			w := 1 / (r.sig[i] * r.sig[i])
			sw += w
			swx += w * roots[i]
			swy += w * r.dev[i]
			swxx += w * roots[i] * roots[i]
			swxy += w * roots[i] * r.dev[i]
		}
		den := sw*swxx - swx*swx
		b := (sw*swxy - swx*swy) / den
		core := (swy*swxx - swx*swxy) / den
		worst := 0.0
		for i := 0; i < 4; i++ {
			res := math.Abs(r.dev[i]-core-b*roots[i]) / r.sig[i]
			if res > worst {
				worst = res
			}
		}
		fmt.Printf("  %2d   %+7.2f    %+9.0f      %.1f\n", r.d, core, b, worst)
	}
	fmt.Println("\nif the cores of the shallow valleys vanish while 12 keeps its +2,")
	fmt.Println("the whole climb reduces to one law: a deep core per gap, plus a")
	fmt.Println("wheel-shallow transient dying at the square-root rate.")
}
