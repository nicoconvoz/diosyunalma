// Command costa is the coastal radar the flash demanded: if typhoons are
// the sea's most violent motion, then pure stillness - perfect agreement
// of the waves - is the sign of LAND. The zero-orchestra's waves arrive
// at every coast IN PHASE (F93), the land grows at each coast (the
// staircase of psi rises a step at every prime), and the waves are
// BOUNDED against the shore - that bound is the Riemann Hypothesis
// itself, told as geography.
//
// The stretched atlas uses Landau's lighthouse: with T the radar's
// aperture (highest zero used),
//
//	L(x) = -(2*pi/T) * sqrt(x) * sum over zeros of cos(gamma*ln(x))
//
// stands at height ~ Lambda(x) on prime powers and near sea level in
// between - so the islands carry RELIEF: each island's height is the
// logarithm of its prime. The map draws its own islands above a
// threshold, labels them, and only then checks against the true prime
// powers; misses are flagged too. Resolution blurs as x*pi/T, so the
// far bands honestly fade into radar haze.
//
// Usage:
//
//	go run ./cmd/costa                      # atlas 2..1013, 600 zeros
//	go run ./cmd/costa -hasta 62 -ceros 100 # the original small chart
package main

import (
	"flag"
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"strings"
)

var emTerms = 320

func zetaEM(s complex128, lnn []float64) complex128 {
	var sum complex128
	sig, t := real(s), imag(s)
	for n := 1; n < emTerms; n++ {
		amp := math.Exp(-sig * lnn[n])
		sn, cs := math.Sincos(t * lnn[n])
		sum += complex(amp*cs, -amp*sn)
	}
	nf := complex(float64(emTerms), 0)
	ns := cmplx.Exp(-s * complex(lnn[emTerms], 0))
	sum += ns * nf / (s - 1)
	sum += ns / 2
	sum += ns * s / nf / 12
	sum -= ns * s * (s + 1) * (s + 2) / (nf * nf * nf) / 720
	sum += ns * s * (s + 1) * (s + 2) * (s + 3) * (s + 4) /
		(nf * nf * nf * nf * nf) / 30240
	return sum
}

func theta(t float64) float64 {
	return t/2*math.Log(t/(2*math.Pi)) - t/2 - math.Pi/8 + 1/(48*t)
}

func zOf(t float64, lnn []float64) float64 {
	z := zetaEM(complex(0.5, t), lnn)
	th := theta(t)
	return real(z)*math.Cos(th) - imag(z)*math.Sin(th)
}

// primePower returns (p, k, true) when n = p^k.
func primePower(n int) (int, int, bool) {
	if n < 2 {
		return 0, 0, false
	}
	for p := 2; p*p <= n; p++ {
		if n%p == 0 {
			k := 0
			for n%p == 0 {
				n /= p
				k++
			}
			if n == 1 {
				return p, k, true
			}
			return 0, 0, false
		}
	}
	return n, 1, true
}

func main() {
	hasta := flag.Float64("hasta", 1013, "map the coast from 2 up to here")
	factor := flag.Float64("sonar", 4, "the onboard sonar's aperture: at map point x it uses zeros up to T = sonar*x, so resolution stays constant as the ship travels")
	flag.Parse()
	Tmax := math.Max(500, *factor**hasta)
	emTerms = int(Tmax/(2*math.Pi)) + 120
	fmt.Printf("LA COSTA — atlas to %.0f, sonar aboard (T = %.0f*x, up to %.0f)\n", *hasta, *factor, Tmax)

	lnn := make([]float64, emTerms+1)
	for n := 1; n <= emTerms; n++ {
		lnn[n] = math.Log(float64(n))
	}
	var gammas []float64
	prevT, prevZ := 12.0, zOf(12.0, lnn)
	for t := 12.02; len(gammas) == 0 || gammas[len(gammas)-1] < Tmax; t += 0.02 {
		zt := zOf(t, lnn)
		if (prevZ < 0) != (zt < 0) {
			lo, hi := prevT, t
			for i := 0; i < 50 && hi-lo > 1e-9; i++ {
				mid := (lo + hi) / 2
				if (zOf(mid, lnn) < 0) == (prevZ < 0) {
					lo = mid
				} else {
					hi = mid
				}
			}
			gammas = append(gammas, (lo+hi)/2)
		}
		prevT, prevZ = t, zt
	}
	fmt.Printf("  sonar magazine: %d zeros up to gamma = %.3f; blur = pi/%.0f (constant in x)\n",
		len(gammas), gammas[len(gammas)-1], *factor)

	// the sonar travels with the ship: at map point x it listens to the
	// zeros up to T(x) = factor*x, so the blur pi*x/T(x) never grows.
	land := func(x float64) float64 {
		Tx := math.Max(500, *factor*x)
		lx := math.Log(x)
		var s float64
		var Teff float64
		for _, g := range gammas {
			if g > Tx {
				break
			}
			s += math.Cos(g * lx)
			Teff = g
		}
		return -(2 * math.Pi / Teff) * math.Sqrt(x) * s
	}

	// self-drawn land, two tiers: ISLANDS rise above half the relief
	// guide; SANDBANKS (the lowlands: prime powers stand only ln p tall)
	// rise above the noise floor but below island stature.
	type isla struct {
		x, h  float64
		n     int  // labeled integer (nearest prime power if hit)
		hit   bool // verified as prime power
		banco bool // sandbank tier
	}
	var islas []isla
	const dl = 0.0006 // ln-step: ~5 points per blur width even at x=2
	xs := []float64{}
	Ls := []float64{}
	lis := []float64{} // the li ruler (F114): the projection-free coordinate
	liAcc := 0.0
	for l := math.Log(2.0) - 0.02; l <= math.Log(*hasta)+0.02; l += dl {
		x := math.Exp(l)
		if len(xs) > 0 {
			pxv := xs[len(xs)-1]
			liAcc += (x - pxv) / math.Log((x+pxv)/2)
		}
		xs = append(xs, x)
		lis = append(lis, liAcc)
		Ls = append(Ls, land(x))
	}
	liOf := func(x float64) float64 { // interpolate the ruler at any x
		lo, hi := 0, len(xs)-1
		for hi-lo > 1 {
			m := (lo + hi) / 2
			if xs[m] <= x {
				lo = m
			} else {
				hi = m
			}
		}
		f := (x - xs[lo]) / (xs[hi] - xs[lo])
		return lis[lo] + f*(lis[hi]-lis[lo])
	}
	// the radar measures its own noise floor: windowed median of |L|
	// (+-0.3 in ln x), robust against the island spikes themselves.
	noiseAt := func(i int) float64 {
		w := int(0.3 / dl)
		lo, hi := i-w, i+w
		if lo < 0 {
			lo = 0
		}
		if hi > len(Ls) {
			hi = len(Ls)
		}
		abs := make([]float64, 0, hi-lo)
		for j := lo; j < hi; j++ {
			abs = append(abs, math.Abs(Ls[j]))
		}
		for a := 0; a < len(abs); a++ { // partial selection sort to the median
			m := a
			for b := a + 1; b < len(abs); b++ {
				if abs[b] < abs[m] {
					m = b
				}
			}
			abs[a], abs[m] = abs[m], abs[a]
			if a > len(abs)/2 {
				break
			}
		}
		return abs[len(abs)/2]
	}
	for i := 1; i < len(xs)-1; i++ {
		if Ls[i] <= Ls[i-1] || Ls[i] < Ls[i+1] {
			continue
		}
		islandBar := 0.5 * math.Log(math.Max(xs[i], 3))
		bancoBar := 3.2 * noiseAt(i)
		if Ls[i] <= math.Min(islandBar, bancoBar) {
			continue
		}
		// with the sonar aboard the blur is pi/factor, constant in x.
		r := math.Max(0.45, 0.35*math.Pi / *factor)
		best, bd := 0, r+1
		for n := int(math.Ceil(xs[i] - r)); n <= int(math.Floor(xs[i]+r)); n++ {
			if _, _, ok := primePower(n); ok && math.Abs(xs[i]-float64(n)) < bd {
				best, bd = n, math.Abs(xs[i]-float64(n))
			}
		}
		banco := Ls[i] <= islandBar
		if best > 0 {
			islas = append(islas, isla{xs[i], Ls[i], best, true, banco})
		} else {
			islas = append(islas, isla{xs[i], Ls[i], int(math.Round(xs[i])), false, banco})
		}
	}
	// a shoulder bump within the blur claims the same integer as its
	// island: keep only the tallest claimant per integer.
	bestOf := map[int]int{}
	for i, is := range islas {
		if !is.hit {
			continue
		}
		if j, ok := bestOf[is.n]; !ok || is.h > islas[j].h {
			bestOf[is.n] = i
		}
	}
	dedup := islas[:0]
	for i, is := range islas {
		if !is.hit || bestOf[is.n] == i {
			dedup = append(dedup, is)
		}
	}
	islas = dedup

	// misses: prime powers with no island claiming them.
	claimed := map[int]bool{}
	for _, is := range islas {
		if is.hit {
			claimed[is.n] = true
		}
	}
	var missed []int
	for n := 2; n <= int(*hasta); n++ {
		if _, _, ok := primePower(n); ok && !claimed[n] {
			missed = append(missed, n)
		}
	}

	// the atlas: 7 stacked bands, each a strip of coast.
	const bands = 7
	const W = 1400.0
	const bandH, topPad = 232.0, 96.0
	H := topPad + bands*bandH + 46
	scaleY := 22.0 // px per unit of ln-height
	var svg strings.Builder
	fmt.Fprintf(&svg, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="700" y="36" font-size="24" text-anchor="middle" font-family="Georgia" fill="#dce8f7">EL ATLAS DE LA COSTA — de 2 a %.0f, sonar a bordo (%d ceros)</text>
<text x="700" y="60" font-size="13" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">sonar que viaja con el barco: apertura T = %.0fx, nitidez constante; eje = regla li(x) (F114), la que NO roba visibilidad; altura de isla = ln(primo)</text>
<text x="700" y="80" font-size="13" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">verdes = islas verificadas (primos), celestes = bancos de arena (potencias, tierras bajas), rojas = falsas, triangulo gris = primo sin isla; dorado punteado = ln x</text>`,
		W, H, W, H, *hasta, len(gammas), *factor)

	liTot := liOf(*hasta)
	liBW := liTot / bands
	xAtLi := func(li float64) float64 { // invert the ruler for tick labels
		lo, hi := 0, len(lis)-1
		for hi-lo > 1 {
			m := (lo + hi) / 2
			if lis[m] <= li {
				lo = m
			} else {
				hi = m
			}
		}
		return xs[lo]
	}
	for b := 0; b < bands; b++ {
		li0 := float64(b) * liBW
		li1 := li0 + liBW
		seaY := topPad + float64(b)*bandH + bandH - 62
		px := func(li float64) float64 { return 30 + (W-60)*(li-li0)/(li1-li0) }
		fmt.Fprintf(&svg, `<rect x="30" y="%.0f" width="%.0f" height="30" fill="#12305e"/>`, seaY, W-60)
		// the land profile on the li ruler.
		fmt.Fprintf(&svg, `<polyline fill="rgba(230,165,58,0.22)" stroke="#e6a53a" stroke-width="1.3" points="30,%.0f`, seaY)
		for i := range xs {
			if lis[i] < li0 || lis[i] > li1 {
				continue
			}
			h := math.Max(Ls[i], 0)
			fmt.Fprintf(&svg, " %.1f,%.1f", px(lis[i]), seaY-h*scaleY)
		}
		fmt.Fprintf(&svg, ` %.0f,%.0f"/>`, W-30, seaY)
		// the relief guide ln x, and x ticks along the ruler.
		fmt.Fprintf(&svg, `<path d="M %.1f %.1f L %.1f %.1f" stroke="#ffd166" stroke-width="1" stroke-dasharray="5,5" opacity="0.5"/>`,
			px(li0), seaY-math.Log(math.Max(xAtLi(li0), 2))*scaleY, px(li1), seaY-math.Log(xAtLi(li1))*scaleY)
		for k := 0; k <= 6; k++ {
			li := li0 + float64(k)/6*liBW
			fmt.Fprintf(&svg, `<text x="%.1f" y="%.0f" font-size="10" fill="#5f7ba0" text-anchor="middle">%.0f</text>`,
				px(li), seaY+22, xAtLi(li))
		}
		// islands and sandbanks in this band.
		for _, is := range islas {
			li := liOf(is.x)
			if li < li0 || li > li1 {
				continue
			}
			col, rad := "#ff5d73", 3.0
			if is.hit && is.banco {
				col, rad = "#6fd3e8", 2.2
			} else if is.hit {
				col = "#7fd7a8"
			}
			fmt.Fprintf(&svg, `<circle cx="%.1f" cy="%.1f" r="%.1f" fill="%s"/>`, px(li), seaY-is.h*scaleY-6, rad, col)
			fmt.Fprintf(&svg, `<text x="%.1f" y="%.1f" font-size="10" font-weight="bold" fill="%s" text-anchor="start" transform="rotate(-90 %.1f %.1f)">%d</text>`,
				px(li)+3, seaY-is.h*scaleY-11, col, px(li)+3, seaY-is.h*scaleY-11, is.n)
		}
		// misses in this band.
		for _, n := range missed {
			li := liOf(float64(n))
			if li < li0 || li > li1 {
				continue
			}
			fmt.Fprintf(&svg, `<text x="%.1f" y="%.0f" font-size="11" fill="#8fa8c7" text-anchor="middle">&#9660;</text>`, px(li), seaY-4)
		}
	}
	fmt.Fprintf(&svg, `<text x="700" y="%.0f" font-size="13" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">Las olas tienen limite contra la tierra: ese limite, dicho exacto, es la Hipotesis de Riemann.</text>`, H-16)
	svg.WriteString(`</svg>`)
	os.WriteFile("costa-atlas.svg", []byte(svg.String()), 0644)

	hits, falses := 0, 0
	for _, is := range islas {
		if is.hit {
			hits++
		} else {
			falses++
		}
	}
	fmt.Printf("\n  islands self-drawn: %d;  verified prime powers: %d;  false: %d;  primes missed: %d of %d\n",
		len(islas), hits, falses, len(missed), hits+len(missed))
	fmt.Printf("  missed: %v\n", missed)
	fmt.Println("  written: costa-atlas.svg")
}
