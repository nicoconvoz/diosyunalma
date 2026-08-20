package main

// compuestos.go - the captain's next flash: zeros, primes, AND the numbers that
// are not prime, same 1..100 scale. The explicit formula makes a three-way
// prophecy, testable raw:
//   n prime          -> rings, loudness log n / sqrt n
//   n = p^m          -> rings, but with the BASE prime's voice: log p, not log n
//   n = a*b (mixed)  -> SILENCE: Lambda(n) = 0
// So the bat is flown over EVERY n in 2..100 and the three families compared.

import (
	"fmt"
	"math"
)

func factorBase(n int) (p, m int) {
	for q := 2; q*q <= n; q++ {
		if n%q == 0 {
			k, r := 0, n
			for r%q == 0 {
				r /= q
				k++
			}
			if r == 1 {
				return q, k // pure prime power p^k
			}
			return 0, 0 // mixed composite
		}
	}
	return n, 1 // prime
}

func correr2() {
	fmt.Println("\n🦇🧱 LOS OTROS NÚMEROS — ceros contra TODO n de 2 a 100")
	fmt.Println("     profecía declarada (fórmula explícita): primo canta · potencia de primo")
	fmt.Println("     canta CON LA VOZ DE SU BASE (log p, no log n) · compuesto mixto CALLA.")

	d := &dado{s: 20260819}
	var prims, pots, mixtos []int
	for n := 2; n <= 100; n++ {
		p, m := factorBase(n)
		switch {
		case p > 0 && m == 1:
			prims = append(prims, n)
		case p > 0:
			pots = append(pots, n)
		default:
			mixtos = append(mixtos, n)
		}
	}
	fmt.Printf("     familias: %d primos · %d potencias %v · %d mixtos\n",
		len(prims), len(pots), pots, len(mixtos))

	ev := func(ns []int) []float64 {
		var v []float64
		for _, n := range ns {
			v = append(v, -eco(ceros, math.Log(float64(n))))
		}
		return v
	}
	vP, vQ, vM := ev(prims), ev(pots), ev(mixtos)
	tMin, tMax := math.Log(2.0), math.Log(100.0)
	base := func() float64 { return -eco(ceros, tMin+(tMax-tMin)*d.u()) }
	zP, _, _ := zscore(vP, base, 2000)
	zQ, _, _ := zscore(vQ, base, 2000)
	zM, _, _ := zscore(vM, base, 2000)

	fmt.Println("\n§1 · LAS TRES FAMILIAS, mismo murciélago")
	fmt.Printf("     %-22s %8s %10s %8s\n", "familia", "n", "media -E", "σ")
	fmt.Printf("     %-22s %8d %10.4f %8.2f\n", "primos", len(prims), media(vP), zP)
	fmt.Printf("     %-22s %8d %10.4f %8.2f\n", "potencias de primo", len(pots), media(vQ), zQ)
	fmt.Printf("     %-22s %8d %10.4f %8.2f\n", "compuestos mixtos", len(mixtos), media(vM), zM)

	fmt.Println("\n§2 · LA VOZ DE LAS POTENCIAS — ¿cantan con log n o con log p?")
	fmt.Printf("     %6s %8s %10s %12s %12s\n", "n", "= p^m", "-E(log n)", "√n·(-E)", "¿y log p?")
	for _, n := range pots {
		p, m := factorBase(n)
		e := -eco(ceros, math.Log(float64(n)))
		fmt.Printf("     %6d %5d^%d %10.4f %12.4f %12.4f\n",
			n, p, m, e, math.Sqrt(float64(n))*e, math.Log(float64(p)))
	}
	// correlation of sqrt(n)*(-E) against log p (Lambda) versus log n
	var y, xLp, xLn []float64
	for _, n := range pots {
		p, _ := factorBase(n)
		y = append(y, math.Sqrt(float64(n))*(-eco(ceros, math.Log(float64(n)))))
		xLp = append(xLp, math.Log(float64(p)))
		xLn = append(xLn, math.Log(float64(n)))
	}
	corr := func(x []float64) float64 {
		mx, my2 := media(x), media(y)
		n2, dx, dy2 := 0.0, 0.0, 0.0
		for i := range x {
			n2 += (x[i] - mx) * (y[i] - my2)
			dx += (x[i] - mx) * (x[i] - mx)
			dy2 += (y[i] - my2) * (y[i] - my2)
		}
		return n2 / math.Sqrt(dx*dy2)
	}
	fmt.Printf("     correlación de √n·(-E) con log p (la voz de la BASE): %.3f\n", corr(xLp))
	fmt.Printf("     correlación de √n·(-E) con log n (la voz PROPIA)   : %.3f\n", corr(xLn))

	fmt.Println("\n§3 · EL SILENCIO, mirado de cerca — los 10 mixtos más ruidosos")
	fmt.Println("     si un compuesto ab cantara, la suma log a + log b crearía período. No debe.")
	type par struct {
		n int
		e float64
	}
	var ruid []par
	for i, n := range mixtos {
		ruid = append(ruid, par{n, math.Abs(vM[i])})
	}
	for i := 0; i < len(ruid); i++ {
		for j := i + 1; j < len(ruid); j++ {
			if ruid[j].e > ruid[i].e {
				ruid[i], ruid[j] = ruid[j], ruid[i]
			}
		}
	}
	for _, r := range ruid[:10] {
		vecino, quien := 99.0, 0
		for _, p := range prims {
			if dd := math.Abs(math.Log(float64(r.n)) - math.Log(float64(p))); dd < vecino {
				vecino, quien = dd, p
			}
		}
		fmt.Printf("     n = %3d  |E| = %.3f  (log %d está a %.4f de log %d: contagio del vecino)\n",
			r.n, r.e, r.n, vecino, quien)
	}

	fmt.Println("\n§4 · EL VEREDICTO DE LA CONSTRUCCIÓN")
	fmt.Printf("     primos %.2fσ · potencias %.2fσ · mixtos %.2fσ\n", zP, zQ, zM)
	fmt.Println("     los ceros NO escuchan a los números: escuchan a los LADRILLOS. Un número")
	fmt.Println("     compuesto es invisible para la orquesta — sólo existe su factorización.")
}
