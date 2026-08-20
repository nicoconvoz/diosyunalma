package main

// Fase XII - the ARITHMETIC GEOMETRY OF RELATIONS. The finding is the captain's,
// found BY HAND: for odd primes p and q, (p+1)/2 = (q-1)/2 exactly when q = p+2.
// Formalised here as ANCHORS: every odd prime p owns two integer anchors
// a-(p) = (p-1)/2 and a+(p) = (p+1)/2, the two integers flanking its half.
// The geometry the plate drew then falls out as theorems:
//   g = 2  ->  a+(p) = a-(q)          shared anchor  (the captain's equality)
//   g = 4  ->  a-(q) - a+(p) = 1      adjacent anchors, touching, no overlap
//   g > 4  ->  (g-4)/2 integers strictly between them: the hueco
// Section 1 VERIFIES those rules exhaustively before anything else runs.
//
// THE FROZEN MASK RULE (written before any spectrum, per her sections 7 and 12):
//   masked bond (i,j) <=> |i-j| <= kmax AND p_i != p_j AND |p_i - p_j| = g
//                         with BOTH prime (they are, by construction of the medium)
//   masked bonds get sign -1; every other bond keeps +1. Amplitudes untouched.
//   Three declared classes, all run, all published: g = 2, g = 4, g = 6.
// Controls, also frozen: (a) random mask, SAME bond count, uniform over eligible
// bonds, 6 seeds; (b) permuted mask, SAME bond count AND same per-distance
// distribution |i-j|, 6 seeds. Rails of Phase X/XI unchanged.
//
// PREDICTION, pre-registered: after five applications of the blade the honest
// expectation is mask = controls. The one open door: twin primes give modes with
// nearly equal log p, so their harmonics sit CLOSE in the ordered spectrum - the
// mask is not uniform in k, and that geometry is genuinely arithmetic.

import (
	"fmt"
	"math"
	"os"
)

const (
	NM    = 400
	T0M   = 100.0
	KM    = 120
	FM    = 30.0
	SEMM  = 6
	PRokM = 0.090
)

type mascara struct {
	neg  map[int]bool // bond id = i*n + j (i<j)
	porK map[int]int  // distance histogram of masked bonds
	m    int
	n    int // row stride for the bond id (0 = the base NM)
}

func bondID(i, j int) int {
	if i > j {
		i, j = j, i
	}
	return i*NM + j
}

// mascaraGap is THE rule: -1 on bonds whose underlying primes differ by g.
func mascaraGap(ms []modo, g int) *mascara {
	M := &mascara{neg: map[int]bool{}, porK: map[int]int{}}
	for i := 0; i < len(ms); i++ {
		for j := i + 1; j < len(ms) && j-i <= KM; j++ {
			d := ms[i].p - ms[j].p
			if d < 0 {
				d = -d
			}
			if d == g {
				M.neg[bondID(i, j)] = true
				M.porK[j-i]++
				M.m++
			}
		}
	}
	return M
}

// mascaraAzar: same bond count, uniform over eligible bonds.
func mascaraAzar(n, m int, d *dado) *mascara {
	M := &mascara{neg: map[int]bool{}}
	for M.m < m {
		i := int(d.u() * float64(n))
		k := 1 + int(d.u()*float64(KM))
		if i+k >= n {
			continue
		}
		id := bondID(i, i+k)
		if !M.neg[id] {
			M.neg[id] = true
			M.m++
		}
	}
	return M
}

// mascaraPermutada: same bond count AND same per-distance histogram, random sites.
func mascaraPermutada(n int, porK map[int]int, d *dado) *mascara {
	M := &mascara{neg: map[int]bool{}}
	for k, c := range porK {
		puestos := 0
		for puestos < c {
			i := int(d.u() * float64(n-k))
			id := bondID(i, i+k)
			if !M.neg[id] {
				M.neg[id] = true
				M.m++
				puestos++
			}
		}
	}
	return M
}

func (M *mascara) sg() func(int, int) float64 {
	return func(i, j int) float64 {
		if M.neg[bondID(i, j)] {
			return -1
		}
		return 1
	}
}

type ficha struct {
	nom         string
	s10, pr, fr float64
	vivos, m    int
	des         float64
	n           int
	niv         []float64
}

func medirM(nom string, ms []modo, h []float64, c func(int) float64, M *mascara, d *dado) ficha {
	var sg func(int, int) float64
	if M != nil {
		sg = M.sg()
	}
	r := espectro(ms, h, c, KM, FM, sg)
	o := medir(r)
	f := ficha{nom: nom, s10: o.s10, pr: o.pr, vivos: o.vivos, n: 1, niv: r.niveles}
	if M != nil {
		f.m = M.m
		f.fr = frustracion(sg, len(ms), KM, d)
	}
	return f
}

func familiaM(nom string, ms []modo, h []float64, c func(int) float64, d *dado, hacer func(*dado) *mascara) ficha {
	var s10s []float64
	acc := ficha{nom: nom}
	for s := 0; s < SEMM; s++ {
		f := medirM(nom, ms, h, c, hacer(d), d)
		s10s = append(s10s, f.s10)
		acc.pr += f.pr / SEMM
		acc.fr += f.fr / SEMM
		acc.vivos += f.vivos
		acc.m = f.m
	}
	acc.s10 = media(s10s)
	acc.des = desvio(s10s)
	acc.vivos /= SEMM
	acc.n = SEMM
	return acc
}

func filaM(f ficha, minViv int) string {
	sd := "        "
	if f.n > 1 {
		sd = fmt.Sprintf("±%-7.3f", f.des)
	}
	riel := "ok"
	if f.vivos < minViv {
		riel = "VACIADA"
	} else if f.pr < PRokM {
		riel = "atrapada"
	}
	return fmt.Sprintf("     %-28s %6d %5d %8.3f %s %7.3f %7.4f  %s",
		f.nom, f.m, f.vivos, f.s10, sd, f.pr, f.fr, riel)
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "escalada" {
		escalada()
		return
	}
	fmt.Println("🧵🌌 LA MÁSCARA — Fase XII: la geometría de las relaciones entre primos,")
	fmt.Println("     convertida en máscara RALA. El hallazgo manual es del capitán:")
	fmt.Println("     (p+1)/2 = (q−1)/2  ⟺  q = p+2 — dos primos gemelos comparten el ancla.")

	// -----------------------------------------------------------------------
	fmt.Println("\n§1 · VALIDACIÓN ALGEBRAICA (su §14.1–2) — las tres reglas, verificadas a mano larga")
	fmt.Println("     formalización: cada primo impar p tiene dos ANCLAS enteras, a±(p) = (p±1)/2.")
	fmt.Println("     reglas a verificar sobre TODOS los pares de primos impares consecutivos ≤ 100000:")
	ps := primos(100000)
	casos, v2, v4, vh, fallas := 0, 0, 0, 0, 0
	for i := 1; i+1 < len(ps); i++ {
		p, q := ps[i], ps[i+1]
		g := q - p
		aMasP := (p + 1) / 2
		aMenQ := (q - 1) / 2
		casos++
		switch {
		case g == 2:
			if aMasP == aMenQ {
				v2++
			} else {
				fallas++
			}
		case g == 4:
			if aMenQ-aMasP == 1 {
				v4++
			} else {
				fallas++
			}
		default:
			if aMenQ-aMasP-1 == (g-4)/2 {
				vh++
			} else {
				fallas++
			}
		}
	}
	fmt.Printf("     %d pares verificados · g=2 comparte ancla: %d/%d · g=4 anclas adyacentes: %d\n",
		casos, v2, v2, v4)
	fmt.Printf("     g>4 hueco = (g−4)/2 enteros exactos: %d · FALLAS: %d\n", vh, fallas)
	if fallas == 0 {
		fmt.Println("     ✓ las tres reglas del diccionario geométrico son EXACTAS, sin excepción.")
		fmt.Println("       La igualdad del capitán es un teorema chico y verdadero: la geometría de")
		fmt.Println("       la lámina es identidad aritmética, no coincidencia visual.")
	} else {
		fmt.Println("     ⛔ hay fallas: el diccionario NO es exacto y nada de lo que sigue vale")
		return
	}

	// -----------------------------------------------------------------------
	ms := medio(4000, NM, T0M)
	c0 := colaLarga(0.5)
	h := normalizar1(homogenea(len(ms)))
	d := &dado{s: 20260819}

	fmt.Println("\n§2 · LA MÁSCARA EN EL OPERADOR — la regla, congelada antes de todo espectro")
	fmt.Println("     enlace marcado (i,j) ⟺ |i−j| ≤ 120 y |p_i − p_j| = g · signo −1 · amplitud intacta")
	fmt.Println("     clases declaradas y TODAS publicadas: g = 2 (gemelos), 4 (primos), 6 (sexys)")
	fmt.Println("     controles: azar a igual cantidad (6 sem.) · permutada a igual cantidad E")
	fmt.Println("     igual distribución por distancia (6 sem.) · rieles de Fase X vigentes")

	f0 := medirM("S0 · sin máscara", ms, h, c0, nil, d)
	minViv := int(0.8 * float64(f0.vivos))
	tot := 0
	for i := 0; i < len(ms); i++ {
		for j := i + 1; j < len(ms) && j-i <= KM; j++ {
			tot++
		}
	}

	fmt.Printf("\n     %-28s %6s %5s %8s %8s %7s %7s  %s\n",
		"brazo", "marc.", "vivos", "Σ²(10)", "±", "PR/N", "frustr", "rieles")
	fmt.Println(filaM(f0, minViv))

	type resg struct {
		g       int
		B, A, C ficha
		dens    float64
	}
	var res []resg
	for _, g := range []int{2, 4, 6} {
		Mg := mascaraGap(ms, g)
		if Mg.m == 0 {
			fmt.Printf("     M · g=%d: CERO enlaces marcados — la clase no existe en este medio\n", g)
			continue
		}
		B := medirM(fmt.Sprintf("M · gemela g=%d", g), ms, h, c0, Mg, d)
		A := familiaM(fmt.Sprintf("A · azar, %d enlaces", Mg.m), ms, h, c0, d,
			func(dd *dado) *mascara { return mascaraAzar(len(ms), Mg.m, dd) })
		C := familiaM("C · permutada por k", ms, h, c0, d,
			func(dd *dado) *mascara { return mascaraPermutada(len(ms), Mg.porK, dd) })
		fmt.Println(filaM(B, minViv))
		fmt.Println(filaM(A, minViv))
		fmt.Println(filaM(C, minViv))
		kmin, kmax := 999, 0
		for k := range Mg.porK {
			if k < kmin {
				kmin = k
			}
			if k > kmax {
				kmax = k
			}
		}
		fmt.Printf("       (g=%d: %d enlaces = %.4f%% de %d · distancias k ∈ [%d,%d])\n\n",
			g, Mg.m, 100*float64(Mg.m)/float64(tot), tot, kmin, kmax)
		res = append(res, resg{g, B, A, C, float64(Mg.m) / float64(tot)})
	}

	// -----------------------------------------------------------------------
	fmt.Println("§3 · EL VEREDICTO (su §12–13): señal = máscara mejor que SUS DOS controles")
	fmt.Println("     por más de 2σ, con rieles intactos")
	hayViva := false
	for _, r := range res {
		zA := (r.A.s10 - r.B.s10) / math.Max(r.A.des, 1e-9)
		zC := (r.C.s10 - r.B.s10) / math.Max(r.C.des, 1e-9)
		fmt.Printf("     g=%d: contra azar %+6.2fσ · contra permutada %+6.2fσ · ΔΣ² vs S0 = %+7.3f",
			r.g, zA, zC, r.B.s10-f0.s10)
		switch {
		case r.B.vivos < minViv:
			fmt.Println("  → VACIADA")
		case zA > 2 && zC > 2 && r.B.pr >= PRokM:
			fmt.Println("  → SEÑAL CANDIDATA")
			hayViva = true
		case math.Abs(r.B.s10-f0.s10) < 3*math.Max(r.A.des, r.C.des):
			fmt.Println("  → la máscara es demasiado RALA para mover el espectro: nulo limpio")
		default:
			fmt.Println("  → indistinguible de sus controles")
		}
	}
	if !hayViva {
		fmt.Println("\n     ⟹ ninguna clase separa de sus controles emparejados. Su §13 manda: no hay")
		fmt.Println("       señal aritmética detectable por este canal A ESTE TAMAÑO. El eco no se corre.")
		fmt.Println("       Nota honesta de escala: con N=400 el medio contiene poquísimos pares gemelos")
		fmt.Println("       (ver §2) — la máscara existe pero casi no toca al operador. Decidir si la")
		fmt.Println("       geometría del capitán opera pide un medio más grande, no otra regla.")
	}

	dibujar12()
}
