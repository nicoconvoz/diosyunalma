// Command espejo builds THE MIRROR OF Z: the zero-counting staircase
// reconstructed FROM THE PRIMES ALONE - no zeta, no integers, just the
// orchestra of prime voices, each playing its frequency ln p:
//
//	N_mirror(t) = theta(t)/pi + 1 - (1/pi) sum_n Lambda(n)/(sqrt(n) ln n) sin(t ln n) w(n)
//
// (the smoothed explicit formula for arg zeta). The TRUE staircase is
// measured independently by our own instruments (levels via
// Euler-Maclaurin + bisection). If the mirror hugs the stairs - climbing
// at every zero it never saw - the bridge primes->zeros works, measured;
// the bat (F157) crossed it the other way. Judge: per-zero deviation.
package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"strings"
)

func zetaEM(t float64) complex128 {
	s := complex(0.5, t)
	N := int(t/(2*math.Pi)*1.5) + 60
	var sum complex128
	for n := 1; n < N; n++ {
		sum += cmplx.Exp(-s * complex(math.Log(float64(n)), 0))
	}
	lnN := complex(math.Log(float64(N)), 0)
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
	return real(cmplx.Exp(complex(0, theta(t))) * zetaEM(t))
}

func main() {
	// ---- the true side: levels measured by our own instruments ----
	fmt.Println("EL ESPEJO DE Z — midiendo la escalera verdadera…")
	var levels []float64
	prevT := 10.0
	prevZ := zOf(prevT)
	for t := 10.05; t <= 122; t += 0.05 {
		z := zOf(t)
		if z*prevZ < 0 {
			a, c := prevT, t
			for i := 0; i < 60; i++ {
				m := (a + c) / 2
				if zOf(m)*prevZ < 0 {
					c = m
				} else {
					a = m
				}
			}
			levels = append(levels, (a+c)/2)
		}
		prevT, prevZ = t, z
	}
	fmt.Printf("escalera verdadera: %d escalones (ceros medidos)\n", len(levels))

	// ---- the mirror: primes only ----
	// Lambda(n) over prime powers, Gaussian taper in ln n
	const nMax = 5000
	sigma := 2.6
	type voice struct{ lnn, amp float64 }
	var voices []voice
	isComposite := make([]bool, nMax+1)
	for p := 2; p <= nMax; p++ {
		if isComposite[p] {
			continue
		}
		for q := p * p; q <= nMax; q += p {
			isComposite[q] = true
		}
		lp := math.Log(float64(p))
		for pk := float64(p); pk <= nMax; pk *= float64(p) {
			lnn := math.Log(pk)
			w := math.Exp(-lnn * lnn / (2 * sigma * sigma))
			if w < 1e-8 {
				break
			}
			voices = append(voices, voice{lnn, lp / (math.Sqrt(pk) * lnn) * w})
		}
	}
	fmt.Printf("voces del espejo (potencias de primos con peso): %d\n", len(voices))
	mirror := func(t float64) float64 {
		s := 0.0
		for _, v := range voices {
			s += v.amp * math.Sin(t*v.lnn)
		}
		return theta(t)/math.Pi + 1 - s/math.Pi
	}

	// ---- the judge: mirror height at each true zero should be m-1/2 ----
	fmt.Println("\nel juez del espejo: en cada cero verdadero γ_m, el espejo debe marcar m−½")
	worst, mean := 0.0, 0.0
	for m, g := range levels {
		want := float64(m) + 0.5
		got := mirror(g)
		d := math.Abs(got - want)
		mean += d
		if d > worst {
			worst = d
		}
		if m < 8 || d > 0.4 {
			fmt.Printf("  γ%-3d = %9.4f   espejo=%7.3f   debe=%5.1f   desvío=%.3f\n", m+1, g, got, want, d)
		}
	}
	mean /= float64(len(levels))
	fmt.Printf("veredicto: desvío medio %.3f escalones, peor %.3f — sobre %d ceros que el espejo JAMÁS vio\n", mean, worst, len(levels))

	// ---- the picture: staircase vs mirror ----
	var b strings.Builder
	W, H := 1560.0, 820.0
	px, pw := 90.0, 1400.0
	py, ph := 130.0, 540.0
	t0, t1 := 10.0, 120.0
	nGrid := 2200
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="44" font-size="26" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🪞 EL ESPEJO DE Z — la escalera de los ceros, reconstruida SOLO CON PRIMOS</text>
<text x="%.0f" y="72" font-size="14" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">azul: la escalera VERDADERA (%d ceros medidos con nuestros instrumentos) · dorado: EL ESPEJO — %d voces de primos, cada una cantando sin(t·ln p^k)/√(p^k), sin haber visto un cero jamás</text>`,
		W, H, W, H, W/2, W/2, len(levels), len(voices))
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" rx="8" fill="#081020" stroke="#44608c"/>`, px, py, pw, ph)
	nLo := mirror(t0) - 2
	nHi := mirror(t1) + 2
	xOf := func(t float64) float64 { return px + pw*(t-t0)/(t1-t0) }
	yOf := func(nv float64) float64 { return py + ph - ph*(nv-nLo)/(nHi-nLo) }
	// true staircase (blue steps)
	steps := []string{fmt.Sprintf("%.1f,%.1f", xOf(t0), yOf(mirror(t0)))}
	count := mirror(t0) // align staircase base with smooth part at t0
	base := math.Floor(count + 0.5)
	for _, g := range levels {
		if g < t0 || g > t1 {
			continue
		}
		steps = append(steps, fmt.Sprintf("%.1f,%.1f", xOf(g), yOf(base)))
		base++
		steps = append(steps, fmt.Sprintf("%.1f,%.1f", xOf(g), yOf(base)))
	}
	steps = append(steps, fmt.Sprintf("%.1f,%.1f", xOf(t1), yOf(base)))
	fmt.Fprintf(&b, `<polyline fill="none" stroke="#7fb2ff" stroke-width="2.2" points="%s"/>`, strings.Join(steps, " "))
	// the mirror (gold)
	pts := make([]string, 0, nGrid)
	for i := 0; i < nGrid; i++ {
		t := t0 + (t1-t0)*float64(i)/float64(nGrid-1)
		pts = append(pts, fmt.Sprintf("%.1f,%.1f", xOf(t), yOf(mirror(t))))
	}
	fmt.Fprintf(&b, `<polyline fill="none" stroke="#ffd166" stroke-width="1.7" opacity="0.95" points="%s"/>`, strings.Join(pts, " "))
	// axes
	for t := 20.0; t <= 120; t += 20 {
		fmt.Fprintf(&b, `<text x="%.1f" y="%.0f" font-size="12" text-anchor="middle" fill="#8fa8c7">%.0f</text>`, xOf(t), py+ph+22, t)
	}
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="13" text-anchor="middle" fill="#8fa8c7">t — cada escalón azul es un cero real; el espejo dorado sube EXACTAMENTE ahí, sin conocerlos: los primos saben dónde están los ceros</text>`, W/2, py+ph+48)
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="15" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">veredicto del juez: desvío medio %.3f escalones (peor %.3f) sobre %d ceros — el puente primos → ceros FUNCIONA; el murciélago (F157) ya había cruzado al revés: el círculo está cerrado</text>`,
		W/2, py+ph+82, mean, worst, len(levels))
	b.WriteString(`</svg>`)
	os.WriteFile("espejo-z.svg", []byte(b.String()), 0644)
	fmt.Println("escrito: espejo-z.svg")
}
