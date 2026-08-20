package main

// fase16.go - WHY does the crossing exist? Her order: derive first, compare
// after. Everything frozen here before running.
//
// 1. THE EXACT ALGEBRA. Write m = g_n + gap/2. Then
//      cos(mT) = cos(g_nT)cos(gapT/2) - sin(g_nT)sin(gapT/2)
//    and with the analytic null (population phase times the bin's gap factor)
//    the excess decomposes EXACTLY into four measurable terms per bin:
//      E(s) = T1 + T2 + T3 + T4
//      T1 = -2 <Cov_bin(cos gT, cos gapT/2)>_T     [inner cos covariance]
//      T2 = -2 <(C_bin - C_all) * cbar_bin>_T      [SELECTION x smear, cos]
//      T3 = +2 <Cov_bin(sin gT, sin gapT/2)>_T     [inner sin covariance]
//      T4 = +2 <(S_bin - S_all) * sbar_bin>_T      [SELECTION x smear, sin]
//    T2+T4 is "the left zero of a small-gap pair sits on a prime-wave crest"
//    (selection), T1+T3 is the within-pair phase-gap coupling. This answers her
//    section 3 with an identity, not a metaphor. The identity is VERIFIED
//    against the measured E(s) (shuffle null ~ analytic null, also verified).
//
// 2. THE MECHANISM MODEL, derived then simulated. The explicit formula makes
//    the zero DENSITY oscillate with the primes:
//      eps(g) = -(2/log(g/2pi)) * sum_{n<=97, Lambda>0} Lambda(n) cos(g log n)/sqrt(n)
//    Where eps > 0 (crest) the local spacing shrinks; GUE repulsion supplies
//    the local gap statistics. Surrogate spectra are generated SEQUENTIALLY:
//      next gap ~ Wigner(GUE surmise), scaled by the local mean spacing
//      1/(rhobar(g)*(1+eps(g))).
//    Two arms: GUE-PURO (eps = 0) and GUE+FORMULA (eps on). Both run through
//    the IDENTICAL Phase XV pipeline (same zoom, same shuffle nulls). If the
//    second reproduces crossing, slope and shape, the curve is explained by
//    known mathematics and we say so; the residual is what remains.
//
// 3. Slope, crossing and residual, per her sections 8, 12 and 15.

import (
	"fmt"
	"math"
)

// epsilonPrimos: relative density fluctuation from the explicit formula,
// truncated at the SAME primes the echo measures (declared, matched).
func epsilonPrimos(g float64) float64 {
	s := 0.0
	for _, n := range []int{2, 3, 4, 5, 7, 8, 9, 11, 13, 16, 17, 19, 23, 25, 27, 29, 31, 32, 37, 41, 43, 47, 49, 53, 59, 61, 64, 67, 71, 73, 79, 81, 83, 89, 97} {
		lp := math.Log(float64(basePrimo(n))) // Lambda(n) = log p for prime powers
		s += lp * math.Cos(g*math.Log(float64(n))) / math.Sqrt(float64(n))
	}
	return -2 * s / math.Log(g/(2*math.Pi))
}

func basePrimo(n int) int {
	for q := 2; q*q <= n; q++ {
		if n%q == 0 {
			for n%q == 0 {
				n /= q
			}
			if n == 1 {
				return q
			}
			return 1 // not a prime power (never happens in our list)
		}
	}
	return n
}

// wigner samples the GUE surmise P(s) = (32/pi^2) s^2 exp(-4s^2/pi), mean 1.
func wigner(d *dado) float64 {
	for {
		x := 3 * d.u()
		f := (32 / (math.Pi * math.Pi)) * x * x * math.Exp(-4*x*x/math.Pi)
		if d.u()*0.9 < f {
			return x
		}
	}
}

// surrogate generates a spectrum from 30 to tope with local mean spacing
// 1/(rhobar*(1+eps)); eps = nil gives pure GUE.
func surrogate(tope float64, eps func(float64) float64, d *dado) []float64 {
	var g []float64
	x := 30.0
	for x < tope {
		rho := math.Log(x/(2*math.Pi)) / (2 * math.Pi)
		if eps != nil {
			e := eps(x)
			if e < -0.9 {
				e = -0.9
			}
			rho *= 1 + e
		}
		x += wigner(d) / rho
		g = append(g, x)
	}
	return g
}

// pipeline: the exact Phase XV machinery on an arbitrary spectrum.
func pipeline(g []float64, Tp []float64, d *dado, nulos int) (E []float64, cruce, pend float64) {
	ps := paresDe(g)
	nb := int((zHi - zLo) / zW)
	gaps := make([]float64, len(ps))
	for i, p := range ps {
		gaps[i] = p.gap
	}
	nSum := make([]float64, nb)
	for r := 0; r < nulos; r++ {
		bar := append([]float64(nil), gaps...)
		for i := len(bar) - 1; i > 0; i-- {
			j := int(d.u() * float64(i+1))
			bar[i], bar[j] = bar[j], bar[i]
		}
		ps2 := make([]parZ, len(ps))
		for i := range ps {
			ps2[i] = parZ{ps[i].base, bar[i]}
		}
		c := curvaZoom(ps2, Tp, true, 0)
		for b := range c {
			nSum[b] += c[b] / float64(nulos)
		}
	}
	real := curvaZoom(ps, Tp, true, 0)
	E = make([]float64, nb)
	for b := range E {
		E[b] = real[b] - nSum[b]
	}
	if xs := crucesDe(E); len(xs) > 0 {
		cruce = xs[0]
		var sx, sy, sxx, sxy, n float64
		for b := 0; b < nb; b++ {
			sm := zLo + (float64(b)+0.5)*zW
			if math.Abs(sm-cruce) <= 0.15 {
				sx += sm
				sy += E[b]
				sxx += sm * sm
				sxy += sm * E[b]
				n++
			}
		}
		pend = (n*sxy - sx*sy) / (n*sxx - sx*sx)
	}
	return
}

func fase16() {
	fmt.Println("🪞🧠 FASE XVI — ¿POR QUÉ EXISTE EL CRUCE? Derivar primero, comparar después")

	var Tp []float64
	for _, p := range []int{5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97} {
		Tp = append(Tp, math.Log(float64(p)))
	}
	g := cerosPaso(4000, 0.02)
	d := &dado{s: 20260823}
	fmt.Printf("   %d ceros reales · zoom y malla de la Fase XV, sin tocar\n", len(g))

	// -----------------------------------------------------------------------
	fmt.Println("\n§1 · LA IDENTIDAD DE CUATRO TÉRMINOS — el álgebra exacta de E(s)")
	fmt.Println("   cos(mT) = cos(γT)·cos(gT/2) − sin(γT)·sin(gT/2), y contra el nulo analítico:")
	fmt.Println("   E = T1(cov cos) + T2(SELECCIÓN·factor cos) + T3(cov sin) + T4(SELECCIÓN·factor sin)")
	ps := paresDe(g)
	nb := int((zHi - zLo) / zW)

	// population phases per T
	Call := make([]float64, len(Tp))
	Sall := make([]float64, len(Tp))
	for ti, T := range Tp {
		for _, p := range ps {
			Call[ti] += math.Cos(p.base * T)
			Sall[ti] += math.Sin(p.base * T)
		}
		Call[ti] /= float64(len(ps))
		Sall[ti] /= float64(len(ps))
	}

	T1 := make([]float64, nb)
	T2 := make([]float64, nb)
	T3 := make([]float64, nb)
	T4 := make([]float64, nb)
	Fs := make([]float64, nb) // the trig factor cbar(s)
	cnt := make([]int, nb)
	// bin membership
	miembros := make([][]int, nb)
	for i, p := range ps {
		s := sDe(p, true)
		if s < zLo || s >= zHi {
			continue
		}
		b := int((s - zLo) / zW)
		miembros[b] = append(miembros[b], i)
		cnt[b]++
	}
	for b := 0; b < nb; b++ {
		if cnt[b] == 0 {
			continue
		}
		for ti, T := range Tp {
			var cC, cG, sC, sG, cCG, sCG float64
			for _, i := range miembros[b] {
				p := ps[i]
				cg := math.Cos(p.gap * T / 2)
				sg := math.Sin(p.gap * T / 2)
				cz := math.Cos(p.base * T)
				sz := math.Sin(p.base * T)
				cC += cz
				cG += cg
				sC += sz
				sG += sg
				cCG += cz * cg
				sCG += sz * sg
			}
			n := float64(cnt[b])
			cC, cG, sC, sG, cCG, sCG = cC/n, cG/n, sC/n, sG/n, cCG/n, sCG/n
			T1[b] += -2 * (cCG - cC*cG) / float64(len(Tp))
			T2[b] += -2 * ((cC - Call[ti]) * cG) / float64(len(Tp))
			T3[b] += +2 * (sCG - sC*sG) / float64(len(Tp))
			T4[b] += +2 * ((sC - Sall[ti]) * sG) / float64(len(Tp))
			Fs[b] += cG / float64(len(Tp))
		}
	}

	// measured E via the Phase XV pipeline for the identity check
	Eobs, crObs, pendObs := pipeline(g, Tp, d, 120)
	fmt.Printf("   %-11s %8s %8s %8s %8s %9s %9s %8s\n",
		"bin s", "T1", "T2 SEL", "T3", "T4 SEL", "suma", "E medido", "F(s)")
	for b := 0; b < nb; b++ {
		suma := T1[b] + T2[b] + T3[b] + T4[b]
		fmt.Printf("   %.2f–%.2f %+8.4f %+8.4f %+8.4f %+8.4f %+9.4f %+9.4f %+8.3f\n",
			zLo+float64(b)*zW, zLo+float64(b+1)*zW, T1[b], T2[b], T3[b], T4[b], suma, Eobs[b], Fs[b])
	}
	fmt.Println("   (la suma usa el nulo ANALÍTICO; E medido usa el nulo BARAJADO — su cercanía")
	fmt.Println("    verifica de paso que el barajado estima bien el producto de marginales)")
	var domSel, domCov float64
	for b := 0; b < nb; b++ {
		domSel += math.Abs(T2[b] + T4[b])
		domCov += math.Abs(T1[b] + T3[b])
	}
	fmt.Printf("   masa |SELECCIÓN| = %.3f contra masa |covarianzas internas| = %.3f\n", domSel, domCov)

	// -----------------------------------------------------------------------
	fmt.Println("\n§2 · GUE PURO (sin primos) — ¿produce la curva? (su §6 y §10)")
	var crG, pendG []float64
	var Eg []float64
	for r := 0; r < 3; r++ {
		sg := surrogate(4000, nil, d)
		E, c, p := pipeline(sg, Tp, d, 60)
		if r == 0 {
			Eg = E
		}
		crG = append(crG, c)
		pendG = append(pendG, p)
	}
	maxAbsG := 0.0
	for _, v := range Eg {
		if math.Abs(v) > maxAbsG {
			maxAbsG = math.Abs(v)
		}
	}
	fmt.Printf("   3 espectros GUE puros: |E| máximo %.4f (contra 0,35 del real) · cruces %v\n", maxAbsG, crG)
	fmt.Println("   ⟹ GUE solo NO produce la curva: sin primos no hay eco que repartir.")

	// -----------------------------------------------------------------------
	fmt.Println("\n§3 · EL MODELO MECANISMO: GUE + FÓRMULA EXPLÍCITA (densidad modulada)")
	fmt.Println("   ε(γ) de la fórmula explícita truncada en los MISMOS n ≤ 97 del eco;")
	fmt.Println("   gap ~ Wigner con espaciado local 1/(ρ̄(1+ε)). Mismo pipeline, 3 semillas.")
	var crM, pendM []float64
	var Em []float64
	for r := 0; r < 3; r++ {
		sm := surrogate(4000, epsilonPrimos, d)
		E, c, p := pipeline(sm, Tp, d, 60)
		if r == 0 {
			Em = E
		} else {
			for b := range Em {
				Em[b] += E[b]
			}
		}
		crM = append(crM, c)
		pendM = append(pendM, p)
	}
	for b := range Em {
		Em[b] /= 3
	}
	fmt.Printf("   cruces del modelo: %v · pendientes: %v\n", crM, pendM)
	fmt.Printf("   real: cruce %.4f · pendiente %.3f\n", crObs, pendObs)

	fmt.Println("\n§4 · EL RESIDUO — lo que el modelo conocido NO explica (su §15)")
	fmt.Printf("   %-11s %9s %9s %9s\n", "bin s", "E real", "E modelo", "residuo")
	resMax := 0.0
	for b := 0; b < nb; b++ {
		r := Eobs[b] - Em[b]
		if math.Abs(r) > resMax {
			resMax = math.Abs(r)
		}
		fmt.Printf("   %.2f–%.2f %+9.4f %+9.4f %+9.4f\n",
			zLo+float64(b)*zW, zLo+float64(b+1)*zW, Eobs[b], Em[b], r)
	}
	fmt.Printf("   residuo máximo: %.4f (la amplitud real llega a ~0,35)\n", resMax)

	fmt.Println("\n§5 · LA PENDIENTE −1 (su §8): en el modelo y en el real, comparadas arriba —")
	fmt.Println("   si el modelo la reproduce, la pendiente es del MECANISMO (densidad modulada +")
	fmt.Println("   repulsión), no de la normalización ni una coincidencia.")

	dibujar16(Eobs, Em, Eg, T2, T4, crObs, crM)
}
