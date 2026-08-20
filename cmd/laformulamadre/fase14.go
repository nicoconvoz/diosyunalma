package main

// fase14.go - THE SWEEP. Her work order: abandon the two boxes, measure the
// whole curve signal(s) with s = gap over MEAN SPACING, empirical nulls from
// MANY shuffles, several depths, threshold-stability, and no post-hoc cuts.
//
// FROZEN BEFORE RUNNING (her sections 7, 10, 14):
//   - unfolding: s_n = (g_{n+1}-g_n) * log(g_n/2pi) / (2pi)  - the LOCAL mean
//     spacing, which fixes a soft spot of F365: low zeros have wide gaps by
//     HEIGHT, not by looseness. Declared here as an improvement, not a tune.
//   - bins of s: {0-0.3, 0.3-0.5, 0.5-0.7, 0.7-0.9, 0.9-1.1, 1.1-1.3,
//     1.3-1.6, 1.6-2.0, 2.0+}  - nine, symmetric around 1, frozen.
//   - periods: T = log p for the 23 primes 5..97, frozen (F365's set).
//   - statistic per bin: mean over periods of -E on the bin's midpoints.
//   - null: 200 shuffles of the gap pool; each midpoint rebuilt as
//     g_n + gap'/2 and re-binned by the foreign gap unfolded AT ITS OWN
//     height - same gap set, pairing destroyed, identical protocol.
//   - depths: zeros up to gamma = 1000, 2000, 4000 (scan step 0.02 so close
//     pairs at height are not skipped).
//   - decision rules: her section 21, verbatim.

import (
	"fmt"
	"math"
)

const nBar = 200

var cortes = []float64{0, 0.3, 0.5, 0.7, 0.9, 1.1, 1.3, 1.6, 2.0, 99}

func cerosPaso(t1, paso float64) []float64 {
	var g []float64
	a, za := 10.0, zetaZ(10.0)
	for b := 10 + paso; b <= t1; b += paso {
		zb := zetaZ(b)
		if (za < 0) != (zb < 0) {
			lo, hi, zlo := a, b, za
			for i := 0; i < 50; i++ {
				m := (lo + hi) / 2
				if (zlo < 0) != (zetaZ(m) < 0) {
					hi = m
				} else {
					lo = m
				}
			}
			g = append(g, (lo+hi)/2)
		}
		a, za = b, zb
	}
	return g
}

func binDe(s float64) int {
	for k := 0; k+1 < len(cortes); k++ {
		if s >= cortes[k] && s < cortes[k+1] {
			return k
		}
	}
	return len(cortes) - 2
}

// senalBins: the per-bin statistic for one assignment of gaps to bases.
func senalBins(base, gaps []float64, Tp []float64) ([]float64, []int) {
	nb := len(cortes) - 1
	sc := make([]float64, nb) // sum of cos over (pairs in bin, periods)
	ss := make([]float64, nb)
	cnt := make([]int, nb)
	_ = ss
	for i := range base {
		s := gaps[i] * math.Log(base[i]/(2*math.Pi)) / (2 * math.Pi)
		b := binDe(s)
		cnt[b]++
		m := base[i] + gaps[i]/2
		for _, T := range Tp {
			sc[b] += math.Cos(m * T)
		}
	}
	out := make([]float64, nb)
	for b := 0; b < nb; b++ {
		if cnt[b] > 0 {
			out[b] = -2 * sc[b] / float64(cnt[b]*len(Tp)) // mean -E per bin
		}
	}
	return out, cnt
}

func fase14() {
	fmt.Println("🪞📈 FASE XIV — EL BARRIDO: la curva que los datos quieran mostrar")
	fmt.Println("   malla congelada: s ∈ {0;0,3;0,5;0,7;0,9;1,1;1,3;1,6;2,0;∞} · 23 primos 5..97")
	fmt.Println("   desdoblado LOCAL declarado: s = gap·log(γ/2π)/2π — corrige que los ceros bajos")
	fmt.Println("   tengan gaps anchos por ALTURA; F365 usaba el espaciado global (mejora, no ajuste)")
	fmt.Printf("   nulo empírico: %d barajados del pool de gaps, protocolo idéntico al real\n\n", nBar)

	var Tp []float64
	for _, p := range []int{5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97} {
		Tp = append(Tp, math.Log(float64(p)))
	}
	gTodos := cerosPaso(4000, 0.02)
	fmt.Printf("   ceros calculados hasta γ = 4000: %d (paso 0,02 para no saltear pares pegados)\n", len(gTodos))

	d := &dado{s: 20260821}
	type filaB struct {
		tope, sLo, sHi, real, mu, sd float64
		n                            int
	}
	var tabla []filaB
	for _, tope := range []float64{1000, 2000, 4000} {
		var g []float64
		for _, x := range gTodos {
			if x <= tope {
				g = append(g, x)
			}
		}
		base := g[:len(g)-1]
		gaps := diffs(g)

		real, cnt := senalBins(base, gaps, Tp)

		// empirical null: many shuffles, same statistic, same binning protocol
		nb := len(cortes) - 1
		acum := make([][]float64, nb)
		for r := 0; r < nBar; r++ {
			bar := append([]float64(nil), gaps...)
			for i := len(bar) - 1; i > 0; i-- {
				j := int(d.u() * float64(i+1))
				bar[i], bar[j] = bar[j], bar[i]
			}
			v, _ := senalBins(base, bar, Tp)
			for b := 0; b < nb; b++ {
				acum[b] = append(acum[b], v[b])
			}
		}

		fmt.Printf("\n   ═══ M = %d ceros (γ ≤ %.0f) ═══\n", len(g), tope)
		fmt.Printf("   %-11s %6s %9s %9s %8s %8s %9s %8s\n",
			"bin s", "pares", "real", "nulo med", "nulo σ", "Δ", "z_emp", "p_emp")
		for b := 0; b < nb; b++ {
			if cnt[b] < 15 {
				fmt.Printf("   %4.1f–%-4.1f %6d   (menos de 15 pares: no se reporta)\n", cortes[b], cortes[b+1], cnt[b])
				continue
			}
			mu, sd := media(acum[b]), desvio(acum[b])
			zb := (real[b] - mu) / math.Max(sd, 1e-12)
			ext := 0
			for _, v := range acum[b] {
				if math.Abs(v-mu) >= math.Abs(real[b]-mu) {
					ext++
				}
			}
			marca := ""
			if math.Abs(zb) >= 3 {
				marca = "  ⚡"
			}
			fmt.Printf("   %4.1f–%-4.1f %6d %+9.4f %+9.4f %8.4f %+8.4f %+9.2f %8.3f%s\n",
				cortes[b], cortes[b+1], cnt[b], real[b], mu, sd, real[b]-mu, zb, float64(ext)/nBar, marca)
			tabla = append(tabla, filaB{tope, cortes[b], cortes[b+1], real[b], mu, sd, cnt[b]})
		}
	}

	fmt.Println("\n   § F365 REPRODUCIDO (las dos cajas viejas, espaciado global) — referencia:")
	fmt.Println("     649: +12,4σ/−8,0σ · 1517: +23,5σ/−20,6σ · barajado ~⅓ (ver corrida madre)")
	fmt.Println("\n   § LECTURA con sus reglas del §21: buscar región AMPLIA y replicada entre")
	fmt.Println("     profundidades, y el cruce de signo estable; nada de elegir el bin más lindo.")
	var F []([7]float64)
	for _, f := range tabla {
		F = append(F, [7]float64{f.tope, f.sLo, f.sHi, f.real, f.mu, f.sd, float64(f.n)})
	}
	dibujar14(F)
}
