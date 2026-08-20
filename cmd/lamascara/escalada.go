package main

// escalada.go - the blessed dial, turned to full power. N was declared a prior
// decision in the Phase XII acta; the auditor's boxed question is whether
//   M_gemelos - M_permutada
// separates systematically from zero as N = 400 -> 800 -> 1600 -> 3200.
// Everything else stays frozen: kmax = 120, F = 30, the long tail, the mask
// rule, the k-matched permuted control with 6 seeds. Seeds run in parallel.

import (
	"fmt"
	"math"
	"runtime"
	"sync"
	"time"
)

type puntoN struct {
	N, m, vivos0, vivosG int
	s0, gem              float64
	permM, permS         float64
	prG, prP             float64
	z                    float64
}

func escalada() {
	fmt.Println("🧵🚀 LA ESCALADA — la perilla N a fondo, con la bendición de la auditora")
	fmt.Println("     pregunta encajonada: ¿M_gemelos − M_permutada se separa de cero al crecer N?")
	fmt.Printf("     todo lo demás congelado: kmax=%d · F=%.0f · máscara g=2 · permutada por k, %d semillas\n",
		KM, FM, SEMM)
	fmt.Printf("     núcleos disponibles: %d — las semillas corren en paralelo\n\n", runtime.NumCPU())

	var puntos []puntoN
	for _, N := range []int{400, 800, 1600, 3200} {
		t0 := time.Now()
		ms := medioN(4000, N, T0M)
		if len(ms) < N {
			fmt.Printf("     N=%d: el mar sólo tiene %d modos — se corre con esos\n", N, len(ms))
		}
		h := normalizar1(homogenea(len(ms)))
		c0 := colaLarga(0.5)
		Mg := mascaraGapN(ms, 2, len(ms))

		var s0F, gemF ficha
		var wg sync.WaitGroup
		wg.Add(2)
		go func() { defer wg.Done(); s0F = medirMN(ms, h, c0, nil, len(ms)) }()
		go func() { defer wg.Done(); gemF = medirMN(ms, h, c0, Mg, len(ms)) }()

		permS10 := make([]float64, SEMM)
		permPR := make([]float64, SEMM)
		sem := make(chan struct{}, runtime.NumCPU())
		for s := 0; s < SEMM; s++ {
			wg.Add(1)
			go func(s int) {
				defer wg.Done()
				sem <- struct{}{}
				defer func() { <-sem }()
				ds := &dado{s: 90000 + uint64(N)*10 + uint64(s)}
				Mp := mascaraPermutadaN(len(ms), Mg.porK, ds)
				f := medirMN(ms, h, c0, Mp, len(ms))
				permS10[s], permPR[s] = f.s10, f.pr
			}(s)
		}
		wg.Wait()

		p := puntoN{N: len(ms), m: Mg.m, vivos0: s0F.vivos, vivosG: gemF.vivos,
			s0: s0F.s10, gem: gemF.s10, permM: media(permS10), permS: desvio(permS10),
			prG: gemF.pr, prP: media(permPR)}
		p.z = (p.permM - p.gem) / math.Max(p.permS, 1e-9)
		puntos = append(puntos, p)
		fmt.Printf("     N=%4d · gemelos %5d enlaces · vivos %d/%d · Σ² S0 %7.3f · gemela %7.3f · permutada %7.3f ± %6.3f · PR %5.3f/%5.3f · Δ = %+6.2fσ  (%s)\n",
			p.N, p.m, p.vivosG, p.vivos0, p.s0, p.gem, p.permM, p.permS, p.prG, p.prP, p.z,
			time.Since(t0).Round(time.Second))
	}

	fmt.Println("\n§ LA LECTURA, con la regla de antes: señal = Δ creciendo sistemáticamente")
	crece := true
	for i := 1; i < len(puntos); i++ {
		if puntos[i].z <= puntos[i-1].z {
			crece = false
		}
	}
	ultimo := puntos[len(puntos)-1]
	switch {
	case crece && ultimo.z > 2:
		fmt.Println("     ⟹ Δ CRECE con N y el último punto supera 2σ: la correspondencia exacta")
		fmt.Println("       par↔gemelo empieza a importar. ESO sería cosa nueva — replicar con otras")
		fmt.Println("       semillas y más N antes de nombrarlo.")
	case ultimo.z > 2:
		fmt.Println("     ⟹ el último punto supera 2σ pero la serie no es monótona: sugestivo, no")
		fmt.Println("       sistemático. Más semillas en el N grande antes de decir nada.")
	default:
		fmt.Println("     ⟹ Δ NO se separa de cero al crecer N: la ventaja de los gemelos sigue")
		fmt.Println("       siendo su distribución por distancia, a toda escala medida. El canal de")
		fmt.Println("       la correspondencia exacta queda cerrado con datos, no con opinión.")
	}
}

// --- N-generic copies of the fixed-N helpers (NM stays 400 for the base run) --

func medioN(topeP, n int, t0 float64) []modo { return medio(topeP, n, t0) }

func mascaraGapN(ms []modo, g, n int) *mascara {
	M := &mascara{neg: map[int]bool{}, porK: map[int]int{}, n: n}
	for i := 0; i < len(ms); i++ {
		for j := i + 1; j < len(ms) && j-i <= KM; j++ {
			d := ms[i].p - ms[j].p
			if d < 0 {
				d = -d
			}
			if d == g {
				M.neg[i*n+j] = true
				M.porK[j-i]++
				M.m++
			}
		}
	}
	return M
}

func mascaraPermutadaN(n int, porK map[int]int, d *dado) *mascara {
	M := &mascara{neg: map[int]bool{}, n: n}
	for k, c := range porK {
		puestos := 0
		for puestos < c {
			i := int(d.u() * float64(n-k))
			id := i*n + (i + k)
			if !M.neg[id] {
				M.neg[id] = true
				M.m++
				puestos++
			}
		}
	}
	return M
}

func medirMN(ms []modo, h []float64, c func(int) float64, M *mascara, n int) ficha {
	var sg func(int, int) float64
	if M != nil {
		sg = func(i, j int) float64 {
			if i > j {
				i, j = j, i
			}
			if M.neg[i*n+j] {
				return -1
			}
			return 1
		}
	}
	r := espectro(ms, h, c, KM, FM, sg)
	o := medir(r)
	f := ficha{s10: o.s10, pr: o.pr, vivos: o.vivos, n: 1}
	if M != nil {
		f.m = M.m
	}
	return f
}
