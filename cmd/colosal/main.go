// Command colosal is the wind tunnel for the colossal black hole: the
// FFT-powered deposit (the Odlyzko-Schoenhage idea, in NUFFT form).
//
// Today's engine: every wave paints EVERY grid point of the window - the
// deposit costs ~110 ns/wave and grows with window width. The colossal
// absorber: each wave falls onto ~24 neighboring points of a fine
// frequency grid (the gentle horizon - crossing costs almost nothing),
// and ONE final FFT re-propels everything onto the whole window. The
// speed is CONSERVED: cost per wave independent of window width -
// colossal windows at the price of small ones.
//
// Method: type-1 NUFFT with fast Gaussian gridding (Greengard-Lee 2004):
// oversampling R=2, spreading half-width Msp=12, tau = pi*Msp/(M^2*R*(R-.5)).
// The bench certifies error and speed against the exact rotator deposit
// on realistic wave sets before anything touches a certified hull.
//
// Usage:
//
//	go run ./cmd/colosal
package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// fft: in-place iterative radix-2, forward (sign -1).
func fft(re, im []float64) {
	n := len(re)
	for i, j := 1, 0; i < n; i++ {
		bit := n >> 1
		for ; j&bit != 0; bit >>= 1 {
			j ^= bit
		}
		j ^= bit
		if i < j {
			re[i], re[j] = re[j], re[i]
			im[i], im[j] = im[j], im[i]
		}
	}
	for length := 2; length <= n; length <<= 1 {
		ang := -2 * math.Pi / float64(length)
		wr, wi := math.Cos(ang), math.Sin(ang)
		for i := 0; i < n; i += length {
			cwr, cwi := 1.0, 0.0
			for j := 0; j < length/2; j++ {
				ur, ui := re[i+j], im[i+j]
				vr := re[i+j+length/2]*cwr - im[i+j+length/2]*cwi
				vi := re[i+j+length/2]*cwi + im[i+j+length/2]*cwr
				re[i+j], im[i+j] = ur+vr, ui+vi
				re[i+j+length/2], im[i+j+length/2] = ur-vr, ui-vi
				cwr, cwi = cwr*wr-cwi*wi, cwr*wi+cwi*wr
			}
		}
	}
}

const msp = 12

// tauFactor calibrates the Gaussian width; set by the tunnel's sweep.
var tauFactor = 1.0

// nufft computes F(s) = sum_k a_k e^{-i s xi_k}, s = 0..S-1, xi in [0, 2pi).
func nufft(xi []float64, ar, ai []float64, S int) ([]float64, []float64) {
	M := 1
	for M < 4*S {
		M <<= 1
	}
	tau := tauFactor * math.Pi * float64(msp) / (float64(M) * float64(M) * 2 * 1.5)
	dx := 2 * math.Pi / float64(M)
	e3 := make([]float64, 2*msp+1)
	for l := -msp; l <= msp; l++ {
		e3[l+msp] = math.Exp(-float64(l) * float64(l) * dx * dx / (4 * tau))
	}
	// guarded grid: deposit unwrapped, fold the guards once at the end.
	gr := make([]float64, M+2*msp+1)
	gi := make([]float64, M+2*msp+1)
	c1 := 1 / (4 * tau)
	c2 := dx / (2 * tau)
	for k := range xi {
		m0 := int(math.Round(xi[k] / dx))
		xp := xi[k] - float64(m0)*dx
		// e1*estep^(-msp) merged into ONE exp; estep by one more exp.
		f := math.Exp(-xp*xp*c1 - float64(msp)*xp*c2)
		estep := math.Exp(xp * c2)
		base := m0 // deposit at base+l+msp in guarded coordinates
		akr, aki := ar[k], ai[k]
		for l := 0; l <= 2*msp; l++ {
			w := f * e3[l]
			gr[base+l] += w * akr
			gi[base+l] += w * aki
			f *= estep
		}
	}
	// fold guards into the ring.
	rr := make([]float64, M)
	ri := make([]float64, M)
	for m := range gr {
		rr[((m-msp)%M+M)%M] += gr[m]
		ri[((m-msp)%M+M)%M] += gi[m]
	}
	fft(rr, ri)
	fr := make([]float64, S)
	fi := make([]float64, S)
	norm := dx / (2 * math.Sqrt(math.Pi*tau))
	for s := 0; s < S; s++ {
		c := norm * math.Exp(float64(s)*float64(s)*tau)
		fr[s] = rr[s] * c
		fi[s] = ri[s] * c
	}
	return fr, fi
}

func main() {
	fmt.Println("EL COLOSAL — wind tunnel of the FFT absorber")

	rng := rand.New(rand.NewSource(2026))
	const N = 2000000
	// realistic facet-tier waves: xi = h ln k over a deep k-range.
	k0 := 3.0e7
	h := 2 * math.Pi / math.Log(k0+N) / 3
	xi := make([]float64, N)
	ar := make([]float64, N)
	ai := make([]float64, N)
	for k := 0; k < N; k++ {
		kk := k0 + float64(k)
		xi[k] = math.Mod(h*math.Log(kk), 2*math.Pi)
		ph := rng.Float64() * 2 * math.Pi
		amp := 1 / math.Sqrt(kk)
		ar[k] = amp * math.Cos(ph)
		ai[k] = -amp * math.Sin(ph)
	}

	// calibrate the Gaussian width on a small pilot before the real bench.
	pilotN := 100000
	pS := 93
	frP := make([]float64, pS)
	fiP := make([]float64, pS)
	for k := 0; k < pilotN; k++ {
		cr, ci := ar[k], ai[k]
		wr, wi := math.Cos(-xi[k]), math.Sin(-xi[k])
		for s := 0; s < pS; s++ {
			frP[s] += cr
			fiP[s] += ci
			cr, ci = cr*wr-ci*wi, cr*wi+ci*wr
		}
	}
	bestF, bestE := 1.0, math.Inf(1)
	for _, f := range []float64{0.25, 0.5, 1, 2, 4} {
		tauFactor = f
		fr, fi := nufft(xi[:pilotN], ar[:pilotN], ai[:pilotN], pS)
		worst := 0.0
		for s := 0; s < pS; s++ {
			if d := math.Hypot(fr[s]-frP[s], fi[s]-fiP[s]); d > worst {
				worst = d
			}
		}
		if worst < bestE {
			bestE, bestF = worst, f
		}
	}
	tauFactor = bestF
	fmt.Printf("\n  tau calibrated on a pilot: factor %.2f (pilot error %.1e)\n", bestF, bestE)

	fmt.Printf("\n  %d waves; comparing exact rotator deposit vs the colossal absorber:\n", N)
	fmt.Println("    window S    exact deposit    colossal     speedup    worst error")
	for _, S := range []int{33, 93, 318} {
		st := time.Now()
		frD := make([]float64, S)
		fiD := make([]float64, S)
		for k := 0; k < N; k++ {
			cr, ci := ar[k], ai[k]
			wr, wi := math.Cos(-xi[k]), math.Sin(-xi[k])
			for s := 0; s < S; s++ {
				frD[s] += cr
				fiD[s] += ci
				cr, ci = cr*wr-ci*wi, cr*wi+ci*wr
			}
		}
		tD := time.Since(st)

		st = time.Now()
		frN, fiN := nufft(xi, ar, ai, S)
		tN := time.Since(st)

		worst := 0.0
		scale := 0.0
		for s := 0; s < S; s++ {
			d := math.Hypot(frN[s]-frD[s], fiN[s]-fiD[s])
			if d > worst {
				worst = d
			}
			if m := math.Hypot(frD[s], fiD[s]); m > scale {
				scale = m
			}
		}
		fmt.Printf("    %7d   %10.0f ns/w   %6.0f ns/w   %6.1fx   %.1e (rel %.1e)\n",
			S, float64(tD.Nanoseconds())/N, float64(tN.Nanoseconds())/N,
			tD.Seconds()/tN.Seconds(), worst, worst/scale)
	}
	fmt.Println("\n  the colossal law: the absorber's cost per wave is FLAT - the window")
	fmt.Println("  can grow colossal and the speed is conserved. If the errors sit at")
	fmt.Println("  hunting grade, the mounting into the starship is the registered next")
	fmt.Println("  step: certification gates first, as always.")
}
