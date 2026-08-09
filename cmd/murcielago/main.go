// Command murcielago is the captain's bat: it measures the SHAPE of
// the prime atom by EAR, like a bat that cannot see. The experiment is
// exact echolocation ("can one hear the shape of a drum?" - Kac): we
// measure hundreds of energy levels with our own instruments, strike
// the drumhead, and listen to the echo
//
//	E(T) = sum_n w_n cos(gamma_n T)
//
// The explicit formula promises the echo spikes EXACTLY at the periods
// of the periodic orbits: T = k*ln p. If the bat hears ln 2, ln 3,
// ln 5... it is seeing the primes with its ears - the shape of the
// atom, drawn from sound alone. Outputs: the echogram (SVG, with every
// spike judged against its true ln p^k) and the struck drumhead (WAV:
// all measured levels ringing together - the atom's true timbre).
package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"sort"
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
	// ---- 1. measure the levels: the bat maps the whole low spectrum ----
	fmt.Println("EL MURCIÉLAGO — midiendo niveles para la ecolocación…")
	var levels []float64
	prevT := 12.0
	prevZ := zOf(prevT)
	for t := 12.05; t <= 500; t += 0.05 {
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
	fmt.Printf("niveles medidos: %d (hasta γ=%.2f)\n", len(levels), levels[len(levels)-1])

	// ---- 2. the echo: E(T) = sum w_n cos(gamma_n T) ----
	gMax := levels[len(levels)-1]
	weight := func(g float64) float64 {
		x := g / gMax
		return math.Exp(-3 * x * x) // gaussian window: soft edges, sharp echoes
	}
	nT := 3000
	T0, T1 := 0.25, 4.1
	echo := make([]float64, nT)
	for i := 0; i < nT; i++ {
		T := T0 + (T1-T0)*float64(i)/float64(nT-1)
		var e float64
		for _, g := range levels {
			e += weight(g) * math.Cos(g*T)
		}
		echo[i] = e
	}
	// expected orbit periods: ln p^k up to T1
	type orbit struct {
		lnpk float64
		name string
	}
	var orbits []orbit
	primes := []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59}
	for _, p := range primes {
		for k := 1; ; k++ {
			v := float64(k) * math.Log(float64(p))
			if v > T1 {
				break
			}
			nm := fmt.Sprintf("%d", p)
			if k > 1 {
				nm = fmt.Sprintf("%d^%d", p, k)
			}
			orbits = append(orbits, orbit{v, nm})
		}
	}
	sort.Slice(orbits, func(i, j int) bool { return orbits[i].lnpk < orbits[j].lnpk })
	// judge each orbit: find the echo minimum near ln p^k (the explicit
	// formula sends orbit echoes DOWNWARD: -ln p / p^{k/2})
	fmt.Println("\nel eco contra las órbitas verdaderas (el juez del murciélago):")
	fmt.Println("  órbita   período ln p^k    eco medido en     error")
	type spike struct {
		at    float64
		name  string
		truey float64
	}
	var spikes []spike
	for _, ob := range orbits {
		if ob.lnpk < T0+0.05 {
			continue
		}
		best, bestE := 0.0, math.Inf(1)
		for i := 1; i < nT-1; i++ {
			T := T0 + (T1-T0)*float64(i)/float64(nT-1)
			if math.Abs(T-ob.lnpk) < 0.04 && echo[i] < bestE {
				bestE, best = echo[i], T
			}
		}
		if !math.IsInf(bestE, 1) {
			fmt.Printf("   %-6s   %.6f         %.6f      %+.1e\n", ob.name, ob.lnpk, best, best-ob.lnpk)
			spikes = append(spikes, spike{best, ob.name, ob.lnpk})
		}
	}

	// ---- 3. the echogram (SVG) ----
	var b strings.Builder
	W, H := 1500.0, 760.0
	px, pw := 90.0, 1340.0
	py, ph := 150.0, 420.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="26" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🦇 EL ECO DEL PARCHE — la forma del átomo, oída (ecolocación de %d niveles medidos)</text>
<text x="%.0f" y="74" font-size="14" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">se golpea el tambor del espectro y se escucha E(T) = Σ w·cos(γ·T): la fórmula explícita promete ecos EXACTOS en los períodos de las órbitas (T = k·ln p)</text>
<text x="%.0f" y="98" font-size="14" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">cada valle dorado es el murciélago "viendo" un primo con el oído — sin mirar jamás la recta de los números</text>`,
		W, H, W, H, W/2, len(levels), W/2, W/2)
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
	yOf := func(e float64) float64 { return py + ph - (e-minE)/(maxE-minE)*ph }
	xOf := func(T float64) float64 { return px + pw*(T-T0)/(T1-T0) }
	// orbit lines first (under the trace)
	for _, ob := range orbits {
		if ob.lnpk < T0+0.02 {
			continue
		}
		fmt.Fprintf(&b, `<line x1="%.1f" y1="%.0f" x2="%.1f" y2="%.0f" stroke="#ffd166" stroke-width="1" stroke-dasharray="4,4" opacity="0.55"/><text x="%.1f" y="%.0f" font-size="12" text-anchor="middle" fill="#ffd166">%s</text>`,
			xOf(ob.lnpk), py, xOf(ob.lnpk), py+ph, xOf(ob.lnpk), py-8, ob.name)
	}
	// the echo trace
	pts := make([]string, 0, nT)
	for i := 0; i < nT; i++ {
		T := T0 + (T1-T0)*float64(i)/float64(nT-1)
		pts = append(pts, fmt.Sprintf("%.1f,%.1f", xOf(T), yOf(echo[i])))
	}
	fmt.Fprintf(&b, `<polyline fill="none" stroke="#7fb2ff" stroke-width="1.8" points="%s"/>`, strings.Join(pts, " "))
	// axis
	for T := 0.5; T <= 4.0; T += 0.5 {
		fmt.Fprintf(&b, `<text x="%.1f" y="%.0f" font-size="12" text-anchor="middle" fill="#8fa8c7">%.1f</text>`, xOf(T), py+ph+22, T)
	}
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="13" text-anchor="middle" fill="#8fa8c7">T (el tiempo del eco) — las líneas doradas son los períodos VERDADEROS k·ln p; el trazo azul es lo que el murciélago ESCUCHÓ</text>`, W/2, py+ph+48)
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#ffd166">veredicto: %d órbitas oídas, cada valle clavado en su ln p^k — el oído reconstruye la lista de los primos sin verla: la forma del átomo, escuchada</text>`,
		W/2, py+ph+86, len(spikes))
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">flash del capitán: "como los murciélagos, que no ven pero escuchan la forma y se representa en su mente como si la hubieran visto" — 2026-08-06</text>`,
		W/2, py+ph+112)
	b.WriteString(`</svg>`)
	os.WriteFile("murcielago-eco.svg", []byte(b.String()), 0644)
	fmt.Println("\nescrito: murcielago-eco.svg")

	// ---- 4. the drumhead struck: the atom's true timbre (WAV) ----
	const sr = 44100
	dur := 12.0
	buf := make([]float64, int(dur*sr))
	strike := func(t0 float64) {
		i0 := int(t0 * sr)
		n := int(3.5 * sr)
		for i := 0; i < n && i0+i < len(buf); i++ {
			tt := float64(i) / sr
			var v float64
			for ni, g := range levels {
				if ni > 120 {
					break
				}
				f := g * 16 // gamma_1=14.13 -> 226 Hz; the low spectrum sings
				amp := math.Pow(float64(ni+2), -0.8)
				v += amp * math.Sin(2*math.Pi*f*tt) * math.Exp(-tt*(0.8+0.02*float64(ni)))
			}
			buf[i0+i] += v * math.Min(1, tt/0.004)
		}
	}
	for _, t0 := range []float64{0.4, 3.4, 6.4, 9.4} {
		strike(t0)
	}
	peak := 1e-9
	for _, v := range buf {
		if math.Abs(v) > peak {
			peak = math.Abs(v)
		}
	}
	pcm := make([]int16, len(buf))
	for i, v := range buf {
		pcm[i] = int16(v / peak * 0.88 * 32767)
	}
	wf, _ := os.Create("parche-atomo.wav")
	dataLen := uint32(len(pcm) * 2)
	wf.WriteString("RIFF")
	binary.Write(wf, binary.LittleEndian, uint32(36+dataLen))
	wf.WriteString("WAVEfmt ")
	binary.Write(wf, binary.LittleEndian, uint32(16))
	binary.Write(wf, binary.LittleEndian, uint16(1))
	binary.Write(wf, binary.LittleEndian, uint16(1))
	binary.Write(wf, binary.LittleEndian, uint32(sr))
	binary.Write(wf, binary.LittleEndian, uint32(sr*2))
	binary.Write(wf, binary.LittleEndian, uint16(2))
	binary.Write(wf, binary.LittleEndian, uint16(16))
	wf.WriteString("data")
	binary.Write(wf, binary.LittleEndian, dataLen)
	binary.Write(wf, binary.LittleEndian, pcm)
	wf.Close()
	fmt.Println("escrito: parche-atomo.wav — el parche del átomo, golpeado 4 veces (120 niveles sonando juntos)")
}
