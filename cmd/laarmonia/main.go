package main

// La Armonia - the captain's flash: take ONLY the primes below 100 and ONLY the
// Riemann zeros below 100, project a harmony between them, and let the bat
// search for the hidden relation. Small on purpose: if the alignment already
// shows at this tiny proportion, it scales.
//
// Both directions are flown:
//   zeros sing, primes listen:  E(T) = (2/M) sum_n cos(gamma_n T) at T = m log p
//   primes sing, zeros listen:  P(t) = sum_p (log p/sqrt p) cos(t log p) at t = gamma_n
// Controls: random periods / random points, 2000 resamples each.

import (
	"fmt"
	"math"
	"sort"
)

var ceros = []float64{
	14.134725, 21.022040, 25.010858, 30.424876, 32.935062,
	37.586178, 40.918719, 43.327073, 48.005151, 49.773832,
	52.970321, 56.446248, 59.347044, 60.831779, 65.112544,
	67.079811, 69.546402, 72.067158, 75.704691, 77.144840,
	79.337375, 82.910381, 84.735493, 87.425275, 88.809111,
	92.491899, 94.651344, 95.870634, 98.831194,
}

var primosCien = []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47,
	53, 59, 61, 67, 71, 73, 79, 83, 89, 97}

type dado struct{ s uint64 }

func (d *dado) u() float64 {
	d.s ^= d.s << 13
	d.s ^= d.s >> 7
	d.s ^= d.s << 17
	return float64(d.s>>11) / float64(uint64(1)<<53)
}

func media(v []float64) float64 {
	s := 0.0
	for _, x := range v {
		s += x
	}
	return s / float64(len(v))
}

func desvio(v []float64) float64 {
	m := media(v)
	s := 0.0
	for _, x := range v {
		s += (x - m) * (x - m)
	}
	return math.Sqrt(s / float64(len(v)-1))
}

// eco is the bat over the zeros: E(T) = (2/M) sum cos(gamma_n T).
func eco(zs []float64, T float64) float64 {
	s := 0.0
	for _, g := range zs {
		s += math.Cos(g * T)
	}
	return 2 * s / float64(len(zs))
}

// canto is the bat over the primes: P(t) = sum (log p/sqrt p) cos(t log p),
// normalised by the total weight so it is comparable across prime sets.
func canto(ps []int, t float64) float64 {
	s, w := 0.0, 0.0
	for _, p := range ps {
		lp := math.Log(float64(p))
		wp := lp / math.Sqrt(float64(p))
		s += wp * math.Cos(t*lp)
		w += wp
	}
	return s / w
}

func zscore(vals []float64, base func() float64, n int) (float64, float64, float64) {
	var res []float64
	for r := 0; r < n; r++ {
		res = append(res, base())
	}
	mb, sb := media(res), desvio(res)
	return (media(vals) - mb) / math.Max(sb, 1e-12), mb, sb
}

func main() {
	fmt.Println("🦇🎼 LA ARMONÍA — primos < 100 contra ceros < 100, cara a cara")
	fmt.Printf("     %d ceros · %d primos · el murciélago vuela en las DOS direcciones\n",
		len(ceros), len(primosCien))
	fmt.Println("     predicción declarada: la fórmula explícita dice que el eco en T = m·log p")
	fmt.Println("     debe ser NEGATIVO (absorción, F350: 13 de 13) con volumen log p·p^(−m/2).")

	d := &dado{s: 20260819}

	// --- direction 1: zeros sing, primes listen ------------------------------
	fmt.Println("\n§1 · LOS CEROS CANTAN, LOS PRIMOS ESCUCHAN — eco en T = log p")
	var vals []float64
	fmt.Printf("     %4s %9s %10s %10s\n", "p", "T=log p", "-E(T)", "log p/√p")
	for _, p := range primosCien {
		T := math.Log(float64(p))
		e := -eco(ceros, T)
		vals = append(vals, e)
		if p <= 13 || p >= 89 {
			fmt.Printf("     %4d %9.4f %10.4f %10.4f\n", p, T, e, T/math.Sqrt(float64(p)))
		}
	}
	tMin, tMax := math.Log(2.0), math.Log(97.0)
	z1, mb1, _ := zscore(vals, func() float64 {
		return -eco(ceros, tMin+(tMax-tMin)*d.u())
	}, 2000)
	fmt.Printf("     media de -E en los 25 primos: %.4f · en períodos al azar: %.4f\n", media(vals), mb1)
	fmt.Printf("     ⟹ los primos suenan a %.2f sigmas del azar (signo predicho: absorción)\n", z1)

	// harmonics m=2,3: the explicit formula also rings at m log p
	var vh []float64
	for _, p := range primosCien {
		for m := 2; m <= 3; m++ {
			T := float64(m) * math.Log(float64(p))
			if T <= 2*tMax {
				vh = append(vh, -eco(ceros, T))
			}
		}
	}
	z1h, _, _ := zscore(vh, func() float64 {
		return -eco(ceros, tMin+(2*tMax-tMin)*d.u())
	}, 2000)
	fmt.Printf("     y los ARMÓNICOS m·log p (potencias de primos): %.2f sigmas\n", z1h)

	// the 1/2 harmony: correlate -E(log p) with log p / p^(1/2)
	var xs, ys []float64
	for i, p := range primosCien {
		xs = append(xs, math.Log(float64(p))/math.Sqrt(float64(p)))
		ys = append(ys, vals[i])
	}
	mx, my := media(xs), media(ys)
	num, dx, dy := 0.0, 0.0, 0.0
	for i := range xs {
		num += (xs[i] - mx) * (ys[i] - my)
		dx += (xs[i] - mx) * (xs[i] - mx)
		dy += (ys[i] - my) * (ys[i] - my)
	}
	corr := num / math.Sqrt(dx*dy)
	// control: correlation with the WRONG laws p^0 and p^-1
	corrCon := func(expo float64) float64 {
		var x2 []float64
		for _, p := range primosCien {
			x2 = append(x2, math.Log(float64(p))*math.Pow(float64(p), expo))
		}
		m2 := media(x2)
		n2, d2 := 0.0, 0.0
		for i := range x2 {
			n2 += (x2[i] - m2) * (ys[i] - my)
			d2 += (x2[i] - m2) * (x2[i] - m2)
		}
		return n2 / math.Sqrt(d2*dy)
	}
	fmt.Printf("     LA RELACIÓN 1/2: correlación de -E(log p) con log p·p^(−1/2): %.3f\n", corr)
	fmt.Printf("       contra las leyes equivocadas: p^0 da %.3f · p^(−1) da %.3f\n",
		corrCon(0), corrCon(-1))

	// --- direction 2: primes sing, zeros listen ------------------------------
	fmt.Println("\n§2 · LOS PRIMOS CANTAN, LOS CEROS ESCUCHAN — P(t) evaluada en t = γₙ")
	var vz []float64
	for _, g := range ceros {
		vz = append(vz, canto(primosCien, g))
	}
	z2, mb2, _ := zscore(vz, func() float64 {
		return canto(primosCien, 10+90*d.u())
	}, 2000)
	fmt.Printf("     media de P en los 29 ceros: %.4f · en puntos al azar: %.4f\n", media(vz), mb2)
	fmt.Printf("     ⟹ los ceros suenan a %.2f sigmas del azar en el canto de los primos\n", z2)
	neg := 0
	for _, v := range vz {
		if v < 0 {
			neg++
		}
	}
	fmt.Printf("     y el SIGNO: %d de %d ceros escuchan canto NEGATIVO (absorción)\n", neg, len(vz))

	// --- scaling: does the alignment strengthen with more zeros? -------------
	fmt.Println("\n§3 · LA ESCALA — ¿la alineación crece con la proporción? (la apuesta del flash)")
	fmt.Printf("     %8s %10s\n", "ceros", "eco (σ)")
	for _, M := range []int{5, 10, 15, 20, 29} {
		zs := ceros[:M]
		var vm []float64
		for _, p := range primosCien {
			vm = append(vm, -eco(zs, math.Log(float64(p))))
		}
		zM, _, _ := zscore(vm, func() float64 {
			return -eco(zs, tMin+(tMax-tMin)*d.u())
		}, 1000)
		fmt.Printf("     %8d %10.2f\n", M, zM)
	}

	// --- the bat free-flight: where are the TOP peaks of |E|? ----------------
	fmt.Println("\n§4 · VUELO LIBRE — los 8 picos más fuertes de |E(T)|, sin decirle dónde mirar")
	type pico struct{ T, v float64 }
	var ps []pico
	prev, prev2 := 0.0, 0.0
	for T := 0.55; T <= 5.0; T += 0.0005 {
		v := math.Abs(eco(ceros, T))
		if prev > prev2 && prev > v && prev > 0.3 {
			ps = append(ps, pico{T - 0.0005, prev})
		}
		prev2, prev = prev, v
	}
	sort.Slice(ps, func(a, b int) bool { return ps[a].v > ps[b].v })
	if len(ps) > 8 {
		ps = ps[:8]
	}
	for _, pk := range ps {
		// nearest m log p with p prime, m<=4
		mejor, quien := 99.0, ""
		for _, p := range primosCien {
			for m := 1; m <= 4; m++ {
				d := math.Abs(pk.T - float64(m)*math.Log(float64(p)))
				if d < mejor {
					mejor, quien = d, fmt.Sprintf("%d·log %d", m, p)
				}
			}
		}
		marca := "  ← ¿?"
		if mejor < 0.01 {
			marca = "  ← " + quien
		}
		fmt.Printf("     T = %.4f  |E| = %.3f  (a %.4f de %s)%s\n", pk.T, pk.v, mejor, quien, marca)
	}
	fmt.Println("\n     nada de esto usa construcción nuestra: son los ceros y los primos, crudos.")

	correr2()

	correr3()
}
