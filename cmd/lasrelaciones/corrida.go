package main

// Fase XI - arithmetic in the RELATIONS. Her work order, in her order:
//   1. open the black box: freeze the exact S_distancia(k) that gave PR/N = 0.281
//   2. three layers: A distance-pure / B distance-arithmetic / C permuted
//   3. the arithmetic catalogue is FROZEN HERE, before any spectrum: Lambda, mu,
//      Liouville lambda. Nothing else may be added after looking.
//   4. match density, per-k structure, live levels; both rails from Phase X
//   5. only then, the echo: operator -> spectrum -> echo -> log p
//   6. and her section 13: replicate the q=0.05 border with many seeds

import (
	"fmt"
	"math"
)

const (
	NR    = 400
	T0R   = 100.0
	KR    = 120
	FR    = 30.0
	SEMR  = 6
	PRokR = 0.090 // Phase X's declared COMPARABLE band
)

// --- frozen arithmetic, integers only ---------------------------------------

func esPotenciaPrima(k int) bool {
	if k < 2 {
		return false
	}
	for q := 2; q*q <= k; q++ {
		if k%q == 0 {
			for k%q == 0 {
				k /= q
			}
			return k == 1
		}
	}
	return true
}

// mu is the Moebius function; liou is Liouville (-1)^Omega.
func mu(k int) int {
	if k == 1 {
		return 1
	}
	m, cnt := 1, 0
	for q := 2; q*q <= k; q++ {
		if k%q == 0 {
			k /= q
			cnt++
			if k%q == 0 {
				return 0
			}
			m = -m
		}
	}
	if k > 1 {
		m = -m
		cnt++
	}
	_ = cnt
	return m
}

func liou(k int) int {
	om := 0
	for q := 2; q*q <= k; q++ {
		for k%q == 0 {
			k /= q
			om++
		}
	}
	if k > 1 {
		om++
	}
	if om%2 == 1 {
		return -1
	}
	return 1
}

// --- distance sign machinery -------------------------------------------------

// signoDist wraps a per-distance sign table s[1..KR] as a bond sign field:
// the sign of bond (i,j) is s[|i-j|]. Toeplitz by construction.
func signoDist(s []float64) func(int, int) float64 {
	return func(i, j int) float64 {
		k := j - i
		if k < 0 {
			k = -k
		}
		if k < 1 || k > KR {
			return 1
		}
		return s[k]
	}
}

func tablaAzar(q float64, d *dado) []float64 {
	s := make([]float64, KR+1)
	for k := 1; k <= KR; k++ {
		s[k] = 1
		if d.u() < q {
			s[k] = -1
		}
	}
	return s
}

func permutar(s []float64, d *dado) []float64 {
	t := append([]float64(nil), s...)
	for k := KR; k > 1; k-- {
		j := 1 + int(d.u()*float64(k))
		t[k], t[j] = t[j], t[k]
	}
	return t
}

func densidadTabla(s []float64) float64 {
	n := 0
	for k := 1; k <= KR; k++ {
		if s[k] < 0 {
			n++
		}
	}
	return float64(n) / float64(KR)
}

// signoEnlace is Phase X's per-bond family, for the border replication.
func signoEnlace(n, kmax int, q float64, d *dado) func(int, int) float64 {
	t := make([][]float64, n)
	for i := range t {
		t[i] = make([]float64, kmax+1)
		for k := range t[i] {
			t[i][k] = 1
		}
	}
	for i := 0; i < n; i++ {
		for k := 1; k <= kmax && i+k < n; k++ {
			if d.u() < q {
				t[i][k] = -1
			}
		}
	}
	return func(i, j int) float64 {
		if i > j {
			i, j = j, i
		}
		if j-i > kmax {
			return 1
		}
		return t[i][j-i]
	}
}

// --- fichas ------------------------------------------------------------------

type ficha struct {
	nom                    string
	s5, s10, s20, alfa, pr float64
	dens, fr               float64
	vivos                  int
	des                    float64
	n                      int
	niv                    []float64 // kept for the echo step
}

func medirFicha(nom string, ms []modo, h []float64, c func(int) float64, sg func(int, int) float64, d *dado) ficha {
	r := espectro(ms, h, c, KR, FR, sg)
	o := medir(r)
	f := ficha{nom: nom, s5: o.s5, s10: o.s10, s20: o.s20, alfa: o.alfa, pr: o.pr,
		vivos: o.vivos, n: 1, niv: r.niveles}
	if sg != nil {
		f.fr = frustracion(sg, len(ms), KR, d)
		tot, neg := 0, 0
		for i := 0; i < len(ms); i++ {
			for k := 1; k <= KR && i+k < len(ms); k++ {
				tot++
				if sg(i, i+k) < 0 {
					neg++
				}
			}
		}
		f.dens = float64(neg) / float64(tot)
	}
	return f
}

func familiaFicha(nom string, ms []modo, h []float64, c func(int) float64, d *dado, hacer func(*dado) func(int, int) float64) ficha {
	var s10s []float64
	acc := ficha{nom: nom}
	for s := 0; s < SEMR; s++ {
		f := medirFicha(nom, ms, h, c, hacer(d), d)
		s10s = append(s10s, f.s10)
		acc.s5 += f.s5 / SEMR
		acc.s20 += f.s20 / SEMR
		acc.alfa += f.alfa / SEMR
		acc.pr += f.pr / SEMR
		acc.dens += f.dens / SEMR
		acc.fr += f.fr / SEMR
		acc.vivos += f.vivos
	}
	acc.s10 = media(s10s)
	acc.des = desvio(s10s)
	acc.vivos /= SEMR
	acc.n = SEMR
	return acc
}

func fila11(f ficha, minViv int) string {
	sd := "        "
	if f.n > 1 {
		sd = fmt.Sprintf("±%-7.3f", f.des)
	}
	riel := "ok"
	if f.vivos < minViv {
		riel = "VACIADA"
	} else if f.pr < PRokR {
		riel = "atrapada"
	}
	return fmt.Sprintf("     %-26s %5d %8.3f %s %7.3f %7.3f %7.3f %7.3f  %s",
		f.nom, f.vivos, f.s10, sd, f.alfa, f.pr, f.dens, f.fr, riel)
}

func cab11() {
	fmt.Printf("     %-26s %5s %8s %8s %7s %7s %7s %7s  %s\n",
		"brazo", "vivos", "Σ²(10)", "±", "α", "PR/N", "dens−", "frustr", "rieles")
}

func main() {
	fmt.Println("🧵🦇 LAS RELACIONES — Fase XI: ¿la aritmética puede vivir en la RELACIÓN")
	fmt.Println("     entre dos sitios, sin mirar un solo cero?")

	ms := medio(4000, NR, T0R)
	c0 := colaLarga(0.5)
	h := normalizar1(homogenea(len(ms)))
	d := &dado{s: 20260818} // Phase X's die, so the black box reopens EXACTLY

	// -----------------------------------------------------------------------
	fmt.Println("\n§1 · LA CAJA NEGRA, ABIERTA Y CONGELADA (su §4)")
	fmt.Println("     La regla exacta del brazo «por distancia» de Fase X, sin descripción verbal:")
	fmt.Println("       s[k] = −1 con probabilidad q, sorteo INDEPENDIENTE por cada k ∈ {1..120}")
	fmt.Println("       (xorshift64, semilla 20260818, el mismo dado de la Fase X)")
	fmt.Println("       signo del enlace (i,j) = s[|i−j|]  ⟹  matriz de signos TOEPLITZ")
	fmt.Println("       sin normalización extra: signo² = 1 deja la fuerza F = 30 intacta")
	fmt.Println("       dominio k = 1..120 = kmax; fuera de rango, +1")
	fmt.Println("       simetría introducida: S depende SÓLO de k — su §11 avisa del peligro")
	fmt.Println("     El PR/N = 0,281 citado corresponde a q = 0,10 con 6 semillas.")

	fmt.Println("\n§2 · EL CATÁLOGO ARITMÉTICO, CONGELADO ANTES DE MIRAR (su §12)")
	fmt.Println("     Tres reglas, sólo enteros y factorización, transformación a signo declarada:")
	fmt.Println("       S_Λ : s[k] = −1 si k es potencia de primo (Λ(k)>0), +1 si no")
	fmt.Println("       S_μ : s[k] = −1 si μ(k) = −1, +1 si μ(k) ∈ {0,+1}")
	fmt.Println("       S_λ : s[k] = λ(k) de Liouville = (−1)^Ω(k)")
	fmt.Println("     Nada se agrega después de ver resultados. Las tres se corren y se publican.")

	sLam := make([]float64, KR+1)
	sMu := make([]float64, KR+1)
	sLio := make([]float64, KR+1)
	for k := 1; k <= KR; k++ {
		sLam[k], sMu[k], sLio[k] = 1, 1, float64(liou(k))
		if esPotenciaPrima(k) {
			sLam[k] = -1
		}
		if mu(k) == -1 {
			sMu[k] = -1
		}
	}
	catalogo := []struct {
		nom string
		s   []float64
	}{
		{"B · S_Λ (potencias de primo)", sLam},
		{"B · S_μ (Moebius)", sMu},
		{"B · S_λ (Liouville)", sLio},
	}

	// -----------------------------------------------------------------------
	f0 := medirFicha("S0 · todos +1", ms, h, c0, nil, d)
	minViv := int(0.8 * float64(f0.vivos))
	fmt.Println("\n§3 · LAS TRES CAPAS — cada aritmética contra su distancia pura y su permutada")
	fmt.Printf("     rieles de Fase X, vigentes: vivos ≥ %d de %d · PR/N ≥ %.3f\n", minViv, f0.vivos, PRokR)
	cab11()
	fmt.Println(fila11(f0, minViv))

	type trio struct{ B, A, C ficha }
	var trios []trio
	for _, cand := range catalogo {
		q := densidadTabla(cand.s)
		B := medirFicha(cand.nom, ms, h, c0, signoDist(cand.s), d)
		A := familiaFicha(fmt.Sprintf("A · pura, q=%.3f", q), ms, h, c0, d,
			func(dd *dado) func(int, int) float64 { return signoDist(tablaAzar(q, dd)) })
		C := familiaFicha("C · permutada", ms, h, c0, d,
			func(dd *dado) func(int, int) float64 { return signoDist(permutar(cand.s, dd)) })
		fmt.Println(fila11(B, minViv))
		fmt.Println(fila11(A, minViv))
		fmt.Println(fila11(C, minViv))
		fmt.Println()
		trios = append(trios, trio{B, A, C})
	}

	// -----------------------------------------------------------------------
	fmt.Println("§4 · EL VEREDICTO POR CAPAS (su §7), regla fijada: señal = B mejor que A y C")
	fmt.Println("     por más de 2σ del control correspondiente, CON rieles intactos")
	for _, t := range trios {
		zA := (t.A.s10 - t.B.s10) / math.Max(t.A.des, 1e-9)
		zC := (t.C.s10 - t.B.s10) / math.Max(t.C.des, 1e-9)
		fmt.Printf("     %-26s contra pura %+6.2fσ · contra permutada %+6.2fσ", t.B.nom, zA, zC)
		switch {
		case t.B.vivos < minViv:
			fmt.Println("  → banda VACIADA: inadmisible")
		case zA > 2 && zC > 2 && t.B.pr >= PRokR:
			fmt.Println("  → SEÑAL CANDIDATA (replicar)")
		case zA > 2 && zC > 2:
			fmt.Println("  → mejora pero PR/N bajo: localización")
		case math.Abs(zA) <= 2 && math.Abs(zC) <= 2:
			fmt.Println("  → A≈B≈C: sin aporte aritmético")
		default:
			fmt.Println("  → la geometría de distancia explica lo que hay")
		}
	}

	// -----------------------------------------------------------------------
	fmt.Println("\n§5 · EL ECO SOBRE EL MEJOR BRAZO (su §10) — operador → espectro → eco → log p")
	mejor, mejorIdx := math.Inf(1), -1
	for i, t := range trios {
		if t.B.vivos >= minViv && t.B.s10 < mejor {
			mejor, mejorIdx = t.B.s10, i
		}
	}
	if mejorIdx < 0 {
		fmt.Println("     ningún brazo aritmético admisible: el eco no se corre")
	} else {
		t := trios[mejorIdx]
		fmt.Printf("     brazo: %s (Σ²(10) = %.3f)\n", t.B.nom, t.B.s10)
		var pes []int
		for _, p := range primos(200) {
			if p >= 5 {
				pes = append(pes, p)
			}
		}
		tMin, tMax := math.Log(5.0), math.Log(199.0)
		ecoDe := func(niv []float64, T float64) float64 {
			s := 0.0
			for _, l := range niv {
				s += math.Cos(l * T)
			}
			return 2 * s / float64(len(niv))
		}
		var vals []float64
		for _, p := range pes {
			vals = append(vals, math.Abs(ecoDe(t.B.niv, math.Log(float64(p)))))
		}
		var res []float64
		for r := 0; r < 300; r++ {
			m := 0.0
			for range pes {
				m += math.Abs(ecoDe(t.B.niv, tMin+(tMax-tMin)*d.u()))
			}
			res = append(res, m/float64(len(pes)))
		}
		z := (media(vals) - media(res)) / math.Max(desvio(res), 1e-12)
		fmt.Printf("     potencia en log p: %.4f · períodos al azar: %.4f ± %.4f → %.2fσ\n",
			media(vals), media(res), desvio(res), z)
		if math.Abs(z) < 2 {
			fmt.Println("     ⟹ el eco NO recupera estructura aritmética: no se fuerza la conexión")
			fmt.Println("       con La Armonía (su §16, último criterio). La banda angosta de 400")
			fmt.Println("       modos ya mostró este límite en F359.")
		} else {
			fmt.Println("     ⟹ hay potencia en log p por encima del control: replicar antes de nombrar")
		}
	}

	// -----------------------------------------------------------------------
	fmt.Println("\n§6 · EL BORDE q = 0,05 DE FASE X (su §13) — ¿estable o accidental? 12 semillas")
	var vivs, s10s, prs []float64
	for s := 0; s < 12; s++ {
		ds := &dado{s: 4242 + uint64(s)}
		f := medirFicha("borde", ms, h, c0, signoEnlace(len(ms), KR, 0.05, ds), d)
		vivs = append(vivs, float64(f.vivos))
		s10s = append(s10s, f.s10)
		prs = append(prs, f.pr)
	}
	pasa := 0
	for _, v := range vivs {
		if int(v) >= minViv {
			pasa++
		}
	}
	fmt.Printf("     vivos: %.1f ± %.1f (pasa el riel %d de 12 veces) · Σ²(10) = %.3f ± %.3f · PR/N = %.3f\n",
		media(vivs), desvio(vivs), pasa, media(s10s), desvio(s10s), media(prs))
	if pasa >= 9 {
		fmt.Println("     ⟹ el borde es ESTABLE: el brazo q=0,05 es admisible en la mayoría de semillas")
	} else if pasa >= 3 {
		fmt.Println("     ⟹ el borde OSCILA alrededor del riel: ni estable ni accidental — es un borde")
		fmt.Println("       de verdad, y decidirlo pide más modos, no más semillas")
	} else {
		fmt.Println("     ⟹ el borde era accidental: casi nunca pasa el riel")
	}

	dibujar11(f0, trios[0].B, trios[1].B, trios[2].B, minViv)
}

// keep the compiler honest about the imported helpers we reuse elsewhere
var _ = pisoGUE
var _ = bitsReciprocidad
