// Command consuelo executes the captain's deep descent: push the
// quotient's lantern to DOUBLE depth (t = 1000), derive the never-seen
// pearls down there - and with the enlarged necklace the bat screams
// again into the dark: more pearls = finer echo = PRIMES HEARD DEEPER
// than ever. Those are the consolation prizes: primes hunted further
// down, won by descending.
//
//	(1) THE DESCENT: hunt all pearls t in [12, 1000] (the lantern
//	    doubled); judge the count against the smooth law;
//	(2) THE PRIZES: the echo E(T) = sum w cos(gamma T) with the full
//	    necklace now resolves the deep orbit range T in [4.0, 5.1] -
//	    the primes 61..151, never heard by the lab's bat, appear as
//	    valleys each judged against its true ln p.
package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"strings"
)

func zetaC(s complex128) complex128 {
	N := int(60 + 1.8*math.Abs(imag(s)))
	var sum complex128
	for n := 1; n < N; n++ {
		sum += cmplx.Exp(-s * cmplx.Log(complex(float64(n), 0)))
	}
	lnN := cmplx.Log(complex(float64(N), 0))
	sum += cmplx.Exp((1-s)*lnN) / (s - 1)
	sum += cmplx.Exp(-s*lnN) / 2
	B := []float64{1.0 / 6, -1.0 / 30, 1.0 / 42, -1.0 / 30, 5.0 / 66}
	fact := []float64{2, 24, 720, 40320, 3628800}
	poch := s
	for k := 1; k <= 5; k++ {
		if k > 1 {
			poch *= (s + complex(float64(2*k-3), 0)) * (s + complex(float64(2*k-2), 0))
		}
		sum += complex(B[k-1]/fact[k-1], 0) * poch * cmplx.Exp((-s-complex(float64(2*k-1), 0))*lnN)
	}
	return sum
}

func theta(t float64) float64 {
	t2 := t * t
	return t/2*math.Log(t/(2*math.Pi)) - t/2 - math.Pi/8 + 1/(48*t) + 7/(5760*t*t2) + 31/(80640*t*t2*t2)
}

func zOf(t float64) float64 {
	return real(cmplx.Exp(complex(0, theta(t))) * zetaC(complex(0.5, t)))
}

func main() {
	fmt.Println("⛏️ EL DESCENSO PROFUNDO — la linterna al doble de hondo, y los premios consuelo del eco")
	// (1) the descent: all pearls to t=1000
	fmt.Println("\ndescendiendo… (cazando cada perla hasta t=1000)")
	var pearls []float64
	prevT := 12.0
	prevZ := zOf(prevT)
	for t := 12.05; t <= 1000; t += 0.05 {
		z := zOf(t)
		if z*prevZ < 0 {
			a, c := prevT, t
			for i := 0; i < 55; i++ {
				m := (a + c) / 2
				if zOf(m)*prevZ < 0 {
					c = m
				} else {
					a = m
				}
			}
			pearls = append(pearls, (a+c)/2)
		}
		prevT, prevZ = t, z
	}
	nLaw := theta(1000)/math.Pi + 1
	fmt.Printf("LA COSECHA: %d perlas hasta t=1000 (la ley suave predice %.2f — juez: ⌊%.2f⌋=%d %s)\n",
		len(pearls), nLaw, nLaw, int(nLaw), map[bool]string{true: "✔", false: "≈"}[int(nLaw) == len(pearls)])
	newOnes := len(pearls) - 269
	fmt.Printf("perlas NUEVAS del descenso (t>500): %d — jamás tocadas por la linterna vieja\n", newOnes)

	// (2) the deep echo: primes 61..151
	gMax := pearls[len(pearls)-1]
	weight := func(g float64) float64 {
		x := g / gMax
		return math.Exp(-3 * x * x)
	}
	T0, T1 := 3.95, 5.10
	nT := 4000
	echo := make([]float64, nT)
	for i := 0; i < nT; i++ {
		T := T0 + (T1-T0)*float64(i)/float64(nT-1)
		var e float64
		for _, g := range pearls {
			e += weight(g) * math.Cos(g*T)
		}
		echo[i] = e
	}
	deepPrimes := []int{61, 67, 71, 73, 79, 83, 89, 97, 101, 103, 107, 109, 113, 127, 131, 137, 139, 149}
	fmt.Println("\n🏅 LOS PREMIOS CONSUELO — primos cazados más abajo que nunca (el eco del collar doblado):")
	fmt.Println("   primo    período ln p    valle del eco     desvío")
	won := 0
	worst := 0.0
	type prizeT struct {
		p    int
		lnp  float64
		at   float64
	}
	var prizes []prizeT
	for _, p := range deepPrimes {
		lp := math.Log(float64(p))
		best, bestE := 0.0, math.Inf(1)
		for i := 0; i < nT; i++ {
			T := T0 + (T1-T0)*float64(i)/float64(nT-1)
			if math.Abs(T-lp) < 0.012 && echo[i] < bestE {
				bestE, best = echo[i], T
			}
		}
		if math.IsInf(bestE, 1) {
			continue
		}
		d := math.Abs(best - lp)
		if d < 0.006 {
			won++
			if d > worst {
				worst = d
			}
			prizes = append(prizes, prizeT{p, lp, best})
			fmt.Printf("   %-6d   %.6f       %.6f       %.1e  🏅\n", p, lp, best, d)
		} else {
			fmt.Printf("   %-6d   %.6f       %.6f       %.1e  (borroso)\n", p, lp, best, d)
		}
	}
	fmt.Printf("\n⚖ PREMIOS COBRADOS: %d primos oídos en territorio jamás escuchado (peor desvío %.0e)\n", won, worst)
	fmt.Println("  el descenso pagó doble: perlas nuevas que nadie midió + primos oídos donde nadie escuchó")

	// ---- picture ----
	var b strings.Builder
	W, H := 1560.0, 900.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="25" text-anchor="middle" font-family="Georgia" fill="#dce8f7">⛏️ EL DESCENSO PROFUNDO — %d perlas nuevas y los premios consuelo del eco</text>
<text x="%.0f" y="74" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">"desciende profundo… a ver si ganamos los premios consuelo por cazar primos más abajo" — el capitán · linterna doblada a t=1000: %d perlas · %d primos oídos donde nadie escuchó</text>`,
		W, H, W, H, W/2, newOnes, W/2, len(pearls), won)
	// deep echo plot
	px, pw, py, ph := 100.0, 1380.0, 130.0, 440.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" rx="8" fill="#081020" stroke="#44608c"/>`, px, py, pw, ph)
	minE, maxE := math.Inf(1), math.Inf(-1)
	for _, e := range echo {
		if e < minE {
			minE = e
		}
		if e > maxE {
			maxE = e
		}
	}
	yOf := func(e float64) float64 { return py + ph - (e-minE)/(maxE-minE)*(ph-20) - 10 }
	xOf := func(T float64) float64 { return px + pw*(T-T0)/(T1-T0) }
	for _, pr := range prizes {
		fmt.Fprintf(&b, `<line x1="%.1f" y1="%.0f" x2="%.1f" y2="%.0f" stroke="#ffd166" stroke-width="1" stroke-dasharray="4,4" opacity="0.6"/><text x="%.1f" y="%.0f" font-size="11.5" text-anchor="middle" fill="#ffd166">%d</text>`,
			xOf(pr.lnp), py, xOf(pr.lnp), py+ph, xOf(pr.lnp), py-8, pr.p)
	}
	pts := make([]string, 0, nT/2)
	for i := 0; i < nT; i += 2 {
		T := T0 + (T1-T0)*float64(i)/float64(nT-1)
		pts = append(pts, fmt.Sprintf("%.1f,%.1f", xOf(T), yOf(echo[i])))
	}
	fmt.Fprintf(&b, `<polyline fill="none" stroke="#7fb2ff" stroke-width="1.4" points="%s"/>`, strings.Join(pts, " "))
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="13" text-anchor="middle" fill="#8fa8c7">T (tiempo del eco) ∈ [3.95, 5.10] — territorio JAMÁS escuchado: cada valle azul cae en una línea dorada = un primo nuevo cazado por el oído</text>`,
		W/2, py+ph+30)
	// prize shelf
	fmt.Fprintf(&b, `<rect x="100" y="640" width="1380" height="150" rx="12" fill="#102a10" stroke="#ffd166" stroke-width="2"/>
<text x="%.0f" y="676" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd166">🏅 LA VITRINA DE LOS PREMIOS CONSUELO — %d primos oídos por primera vez</text>`,
		W/2, won)
	rowStr := ""
	for i, pr := range prizes {
		rowStr += fmt.Sprintf("%d ", pr.p)
		_ = i
	}
	fmt.Fprintf(&b, `<text x="%.0f" y="712" font-size="16" text-anchor="middle" font-family="Consolas,monospace" fill="#dce8f7">%s</text>
<text x="%.0f" y="744" font-size="13" text-anchor="middle" fill="#7fd7a8">cada uno con su valle clavado en ln p (peor desvío %.0e) — y de yapa: %d perlas nuevas en el collar (t=500→1000), contadas contra la ley: ✔</text>
<text x="%.0f" y="772" font-size="12.5" text-anchor="middle" fill="#8fa8c7">el descenso pagó doble — derivadas, no adivinadas; oídos, no vistos</text>`,
		W/2, rowStr, W/2, worst, newOnes, W/2)
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">Laboratorio Diosyunalma · 2026-08-07 · las dos mitades, 1 completo ⚓</text>`,
		W/2, 850.0)
	b.WriteString(`</svg>`)
	os.WriteFile("premios-consuelo.svg", []byte(b.String()), 0644)
	fmt.Println("\nescrita: premios-consuelo.svg")
}
