// Command primoyescala answers the captain's follow-up: what relation does a
// prime number have with a numeration scale of that prime's own order - are
// there bridges?
//
// There are bridges, and they carry heavy traffic. But the honest answer is
// stronger than "there are bridges":
//
//	THERE IS NO BRIDGE TO BUILD. THE PRIME **IS** THE SCALE.
//
// Ostrowski's theorem says every nontrivial way of measuring size on the
// rational numbers is either the ordinary one or the p-adic one attached to
// some prime. One scale per prime, plus exactly one more. That list is
// complete - nothing else exists. So primes and scales are not two families
// with bridges between them; they are ONE family, counted twice.
//
// Below that headline, five exact bridges are measured, each of them turning
// an ARITHMETIC fact (how a prime divides something) into a WRITING fact
// (what the base-p digits look like):
//
//	THE TRAILING ZEROS   v_p(n) = how many zeros n ends with, written in base p
//	LEGENDRE             v_p(n!) = (n - s_p(n)) / (p-1), s_p = digit sum
//	KUMMER               v_p(C(m+n,m)) = number of CARRIES adding m+n in base p
//	LUCAS                C(m,n) mod p = product of C over the base-p digits
//	THE PERIOD           the repeating length of 1/q in base b IS ord_b(q),
//	                     the order of b in the group mod q  (+ Midy's halves)
//
// And then back to the book: zeta's Euler product has exactly one factor per
// prime - one per scale - and completing it adds exactly ONE more factor, the
// one of the infinite place, which is where the half lives (F242). The primes
// hand over the scales; the one scale that is not a prime's hands over the ½.
package main

import (
	"fmt"
	"math/big"
	"os"
	"strings"
)

// ------------------------------------------------------------------ helpers

func vp(n int, p int) int {
	if n == 0 {
		return 0
	}
	v := 0
	for n%p == 0 {
		n /= p
		v++
	}
	return v
}

func digitosEnBase(n, b int) []int { // least significant first
	if n == 0 {
		return []int{0}
	}
	var ds []int
	for ; n > 0; n /= b {
		ds = append(ds, n%b)
	}
	return ds
}

func cerosDelFinal(n, b int) int {
	ds := digitosEnBase(n, b)
	c := 0
	for _, d := range ds {
		if d != 0 {
			break
		}
		c++
	}
	return c
}

func sumaDigitos(n, b int) int {
	s := 0
	for _, d := range digitosEnBase(n, b) {
		s += d
	}
	return s
}

func escribirEnBase(n, b int) string {
	ds := digitosEnBase(n, b)
	var sb strings.Builder
	for i := len(ds) - 1; i >= 0; i-- {
		if b <= 36 {
			sb.WriteByte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"[ds[i]])
		} else {
			fmt.Fprintf(&sb, "%02d ", ds[i])
		}
	}
	return sb.String()
}

func vpBig(z *big.Int, p int64) int {
	if z.Sign() == 0 {
		return 0
	}
	bp := big.NewInt(p)
	t := new(big.Int).Set(z)
	m := new(big.Int)
	v := 0
	for {
		q, r := new(big.Int).QuoRem(t, bp, m)
		if r.Sign() != 0 {
			return v
		}
		t = q
		v++
	}
}

func binom(m, n int) *big.Int {
	return new(big.Int).Binomial(int64(m), int64(n))
}

func factorial(n int) *big.Int {
	return new(big.Int).MulRange(1, int64(n))
}

// acarreos counts the carries when adding m and n written in base p.
func acarreos(m, n, p int) int {
	c, llev := 0, 0
	for m > 0 || n > 0 || llev > 0 {
		s := m%p + n%p + llev
		if s >= p {
			llev = 1
			c++
		} else {
			llev = 0
		}
		m /= p
		n /= p
		if m == 0 && n == 0 && llev == 0 {
			break
		}
	}
	return c
}

// orden returns the multiplicative order of b modulo q (gcd(b,q)=1).
func orden(b, q int) int {
	x := b % q
	for k := 1; k <= q; k++ {
		if x == 1 {
			return k
		}
		x = x * b % q
	}
	return -1
}

// periodoDe returns the repeating period of 1/q in base b by long division,
// plus the repeating digits themselves.
func periodoDe(q, b int) (int, []int) {
	visto := map[int]int{}
	var digs []int
	r := 1 % q
	for i := 0; r != 0; i++ {
		if pos, ok := visto[r]; ok {
			return i - pos, digs[pos:]
		}
		visto[r] = i
		r *= b
		digs = append(digs, r/q)
		r %= q
	}
	return 0, nil
}

// --------------------------------------------------------- the p-adic scale

func vpRat(x *big.Rat, p int64) int {
	return vpBig(x.Num(), p) - vpBig(x.Denom(), p)
}

// tamanoP is |x|_p, the size of x on the scale of the prime p.
func tamanoP(x *big.Rat, p int64) *big.Rat {
	if x.Sign() == 0 {
		return new(big.Rat)
	}
	v := vpRat(x, p)
	av := v
	if av < 0 {
		av = -av
	}
	pot := new(big.Int).Exp(big.NewInt(p), big.NewInt(int64(av)), nil)
	r := new(big.Rat).SetInt64(1)
	if v > 0 {
		r.SetFrac(big.NewInt(1), pot)
	} else if v < 0 {
		r.SetInt(pot)
	}
	return r
}

func tamanoReal(x *big.Rat) *big.Rat {
	r := new(big.Rat).Set(x)
	if r.Sign() < 0 {
		r.Neg(r)
	}
	return r
}

func maxRat(a, b *big.Rat) *big.Rat {
	if a.Cmp(b) >= 0 {
		return a
	}
	return b
}

// ----------------------------------------------------------------- the run

func main() {
	fmt.Println("🌉 EL PRIMO Y LA ESCALA — ¿qué relación hay, y hay puentes?")
	fmt.Println("\n   la pregunta del capitán: qué relación tiene un número primo con una escala")
	fmt.Println("   de numeración del orden de ese primo. ¿Hay puentes?")
	fmt.Println("\n   SÍ, hay cinco, y son exactos. Pero la respuesta honesta es más fuerte que eso,")
	fmt.Println("   y va al final.")

	primos := []int{2, 3, 5, 7, 11, 13}

	// ---- BRIDGE 1: trailing zeros ----
	fmt.Println("\nPUENTE 1 · LOS CEROS DEL FINAL — cuántas veces p divide a n, LEÍDO EN LA ESCRITURA")
	fmt.Println("   escrito en base p, n termina en exactamente v_p(n) ceros. En ninguna otra base.")
	fmt.Println("\n      n        base p    n escrito ahí        ceros al final   v_p(n)")
	casos1 := []struct{ n, p int }{{8, 2}, {24, 2}, {96, 2}, {81, 3}, {54, 3}, {125, 5}, {350, 5}, {343, 7}, {1331, 11}}
	fallos1 := 0
	for _, c := range casos1 {
		z := cerosDelFinal(c.n, c.p)
		v := vp(c.n, c.p)
		mk := "✓"
		if z != v {
			mk = "✗"
			fallos1++
		}
		fmt.Printf("   %6d       %3d      %-16s      %3d          %3d  %s\n",
			c.n, c.p, escribirEnBase(c.n, c.p), z, v, mk)
	}
	// the contrast: base 10 is nobody's scale
	fmt.Println("\n   ¿Y en base 10, que no es la escala de ningún primo?")
	for _, n := range []int{8, 24, 350, 1000} {
		fmt.Printf("   %6d en base 10 termina en %d ceros, pero v_2=%d y v_5=%d — mira DOS primos\n",
			n, cerosDelFinal(n, 10), vp(n, 2), vp(n, 5))
	}
	fmt.Println("   y no ve a ninguno limpio: los ceros del final valen min(v₂, v₅). La escritura")
	fmt.Println("   decimal es una mezcla; la escritura en base p es EL PRIMO MISMO, hecho dígitos.")

	// ---- BRIDGE 2: Legendre ----
	fmt.Println("\nPUENTE 2 · LEGENDRE — la potencia de p en n! ES la suma de dígitos en base p")
	fmt.Println("\n        v_p(n!) = ( n − s_p(n) ) / (p−1)          s_p = suma de los dígitos en base p")
	fmt.Println("\n      n     p     n en base p        s_p(n)   (n−s_p)/(p−1)   v_p(n!) real")
	fallos2 := 0
	pruebas2 := []struct{ n, p int }{{10, 2}, {100, 2}, {1000, 2}, {50, 3}, {243, 3}, {100, 5}, {200, 7}, {150, 11}, {300, 13}}
	for _, c := range pruebas2 {
		s := sumaDigitos(c.n, c.p)
		formula := (c.n - s) / (c.p - 1)
		real := vpBig(factorial(c.n), int64(c.p))
		mk := "✓"
		if formula != real {
			mk = "✗"
			fallos2++
		}
		fmt.Printf("   %5d   %3d   %-18s   %4d        %5d          %5d  %s\n",
			c.n, c.p, escribirEnBase(c.n, c.p), s, formula, real, mk)
	}
	fmt.Println("   → un hecho puramente ARITMÉTICO (cuántas veces p divide a n factorial) es")
	fmt.Println("     exactamente un hecho de ESCRITURA (la suma de los dígitos en base p).")

	// ---- BRIDGE 3: Kummer ----
	fmt.Println("\nPUENTE 3 · KUMMER — la potencia de p en un combinatorio ES el número de ACARREOS")
	fmt.Println("   sumá m + n en base p y contá cuántas veces te llevás una. Ese número es v_p(C(m+n,m)).")
	fmt.Println("\n      m     n     p    m en base p   n en base p   acarreos   v_p(C(m+n,m))")
	fallos3 := 0
	pruebas3 := []struct{ m, n, p int }{{5, 7, 2}, {13, 19, 2}, {60, 40, 2}, {14, 17, 3}, {40, 44, 3}, {23, 31, 5}, {48, 60, 7}, {70, 85, 11}}
	for _, c := range pruebas3 {
		ac := acarreos(c.m, c.n, c.p)
		v := vpBig(binom(c.m+c.n, c.m), int64(c.p))
		mk := "✓"
		if ac != v {
			mk = "✗"
			fallos3++
		}
		fmt.Printf("   %5d %5d  %4d   %-12s  %-12s   %4d        %4d  %s\n",
			c.m, c.n, c.p, escribirEnBase(c.m, c.p), escribirEnBase(c.n, c.p), ac, v, mk)
	}
	fmt.Println("   → llevarse una al sumar en base p ES que el primo entre una vez más. La misma cosa.")

	// ---- BRIDGE 4: Lucas ----
	fmt.Println("\nPUENTE 4 · LUCAS — el combinatorio módulo p se arma DÍGITO POR DÍGITO en base p")
	fmt.Println("\n        C(m,n) ≡ Π sobre los dígitos C(mᵢ, nᵢ)   (mod p)")
	fmt.Println("\n      m     n     p    C(m,n) mod p    Π de los dígitos    ")
	fallos4 := 0
	pruebas4 := []struct{ m, n, p int }{{100, 37, 5}, {247, 88, 7}, {1000, 333, 3}, {512, 200, 2}, {365, 121, 11}, {900, 450, 13}}
	for _, c := range pruebas4 {
		directo := new(big.Int).Mod(binom(c.m, c.n), big.NewInt(int64(c.p))).Int64()
		dm := digitosEnBase(c.m, c.p)
		dn := digitosEnBase(c.n, c.p)
		prod := int64(1)
		for i := range dm {
			ni := 0
			if i < len(dn) {
				ni = dn[i]
			}
			if ni > dm[i] {
				prod = 0
				break
			}
			prod = prod * new(big.Int).Mod(binom(dm[i], ni), big.NewInt(int64(c.p))).Int64() % int64(c.p)
		}
		mk := "✓"
		if prod != directo {
			mk = "✗"
			fallos4++
		}
		fmt.Printf("   %5d %5d  %4d       %5d            %5d          %s\n", c.m, c.n, c.p, directo, prod, mk)
	}
	fmt.Println("   → el combinatorio, que es aritmética pura, se calcula mirando los DÍGITOS.")

	// ---- mass verification of the four arithmetic bridges ----
	fmt.Println("\n   ✅ VERIFICACIÓN MASIVA de los cuatro puentes aritméticos")
	fmt.Println("   (los ejemplos de arriba son para mirarlos; esto es para creerles)")
	masa := func(nombre string, n, fallos, extra int, notaExtra string) {
		mk := "✓"
		if fallos > 0 {
			mk = "✗"
		}
		fmt.Printf("   %-28s %8d casos   %4d fallos  %s   %s\n", nombre, n, fallos, mk, notaExtra)
	}
	// 1. trailing zeros
	n1, f1 := 0, 0
	for _, p := range primos {
		for n := 1; n <= 5000; n++ {
			n1++
			if cerosDelFinal(n, p) != vp(n, p) {
				f1++
			}
		}
	}
	masa("ceros del final", n1, f1, 0, "")
	// 2. Legendre
	n2, f2 := 0, 0
	for _, p := range primos {
		for n := 1; n <= 400; n++ {
			n2++
			if (n-sumaDigitos(n, p))/(p-1) != vpBig(factorial(n), int64(p)) {
				f2++
			}
		}
	}
	masa("Legendre", n2, f2, 0, "")
	// 3. Kummer
	n3, f3 := 0, 0
	for _, p := range primos {
		for m := 0; m <= 45; m++ {
			for n := 0; n <= 45; n++ {
				n3++
				if acarreos(m, n, p) != vpBig(binom(m+n, m), int64(p)) {
					f3++
				}
			}
		}
	}
	masa("Kummer", n3, f3, 0, "")
	// 4. Lucas - counting how many came out NONZERO, which is what makes it a real test
	n4, f4, nz4 := 0, 0, 0
	for _, p := range primos {
		bp := big.NewInt(int64(p))
		for m := 0; m <= 160; m++ {
			for n := 0; n <= m; n++ {
				n4++
				directo := new(big.Int).Mod(binom(m, n), bp).Int64()
				dm := digitosEnBase(m, p)
				dn := digitosEnBase(n, p)
				prod := int64(1)
				for i := range dm {
					ni := 0
					if i < len(dn) {
						ni = dn[i]
					}
					if ni > dm[i] {
						prod = 0
						break
					}
					prod = prod * new(big.Int).Mod(binom(dm[i], ni), bp).Int64() % int64(p)
				}
				if prod != directo {
					f4++
				}
				if directo != 0 {
					nz4++
				}
			}
		}
	}
	masa("Lucas", n4, f4, nz4, fmt.Sprintf("(%d con resto NO nulo — el caso que de verdad prueba)", nz4))
	fallos1 += f1
	fallos2 += f2
	fallos3 += f3
	fallos4 += f4

	// ---- BRIDGE 5: the period IS the group order ----
	fmt.Println("\nPUENTE 5 · EL PERÍODO — el largo de la repetición ES el orden del grupo")
	fmt.Println("   el período de 1/q escrito en base b es exactamente ord_b(q): el orden de b")
	fmt.Println("   en el grupo multiplicativo módulo q. La escritura VE la estructura del primo.")
	fmt.Println("\n      1/q en base b     período medido    ord_b(q)   ¿máximo (q−1)?")
	fallos5 := 0
	reptend := 0
	pruebas5 := []struct{ q, b int }{{7, 10}, {17, 10}, {19, 10}, {13, 10}, {3, 10}, {11, 10},
		{7, 2}, {11, 2}, {13, 2}, {5, 3}, {7, 3}, {31, 2}}
	for _, c := range pruebas5 {
		per, _ := periodoDe(c.q, c.b)
		o := orden(c.b, c.q)
		mk := "✓"
		if per != o {
			mk = "✗"
			fallos5++
		}
		maxi := ""
		if per == c.q-1 {
			maxi = "SÍ — b es raíz primitiva mod q"
			reptend++
		}
		fmt.Printf("   1/%-3d base %2d        %5d          %5d  %s   %s\n", c.q, c.b, per, o, mk, maxi)
	}
	fmt.Printf("   → %d de %d casos con período MÁXIMO q−1 (primos de repetición completa).\n", reptend, len(pruebas5))
	// mass verification of the period bridge
	n5m, f5m, rep5 := 0, 0, 0
	criba := func(N int) []int {
		comp := make([]bool, N+1)
		var ps []int
		for p := 2; p <= N; p++ {
			if comp[p] {
				continue
			}
			ps = append(ps, p)
			for m := p * p; m > 0 && m <= N; m += p {
				comp[m] = true
			}
		}
		return ps
	}
	for _, q := range criba(500) {
		for _, bb := range []int{2, 3, 5, 7, 10} {
			if bb%q == 0 {
				continue
			}
			n5m++
			per, _ := periodoDe(q, bb)
			if per != orden(bb, q) {
				f5m++
			}
			if per == q-1 {
				rep5++
			}
		}
	}
	fmt.Printf("   ✅ masivo: %d pares (primo q<500 × base) · %d fallos · %d con período máximo q−1\n", n5m, f5m, rep5)
	fallos5 += f5m

	fmt.Println("\n   y MIDY, de yapa: si el período es par, sus dos mitades suman puros (b−1)")
	fallosMidy := 0
	midyOK := 0
	for _, c := range []struct{ q, b int }{{7, 10}, {17, 10}, {19, 10}, {11, 2}, {13, 2}, {5, 3}} {
		per, digs := periodoDe(c.q, c.b)
		if per%2 != 0 || per == 0 {
			continue
		}
		k := per / 2
		ok := true
		var izq, der, sum []string
		for i := 0; i < k; i++ {
			s := digs[i] + digs[i+k]
			if s != c.b-1 {
				ok = false
			}
			izq = append(izq, fmt.Sprintf("%d", digs[i]))
			der = append(der, fmt.Sprintf("%d", digs[i+k]))
			sum = append(sum, fmt.Sprintf("%d", s))
		}
		mk := "✓"
		if ok {
			midyOK++
		} else {
			mk = "✗"
			fallosMidy++
		}
		fmt.Printf("   1/%-3d base %2d:  %s + %s = %s   %s\n", c.q, c.b,
			strings.Join(izq, ""), strings.Join(der, ""), strings.Join(sum, ""), mk)
	}
	fmt.Printf("   → %d de %d exactos: las dos mitades del período del primo se completan entre sí.\n", midyOK, midyOK+fallosMidy)

	// ---- BRIDGE 6: Ostrowski - there is no bridge, the prime IS the scale ----
	fmt.Println("\nPUENTE 6 · Y ACÁ SE ACABA LA METÁFORA — no hay puente: EL PRIMO **ES** LA ESCALA")
	fmt.Println("   cada primo define su propia manera de medir tamaño, |x|_p, y esa manera cumple")
	fmt.Println("   algo que la escala común NO cumple — la desigualdad ULTRAMÉTRICA:")
	fmt.Println("\n        |x+y|_p ≤ max( |x|_p , |y|_p )        ← en la escala de un primo, SIEMPRE")
	fmt.Println("        |x+y|   ≤ max( |x| , |y| )            ← en la escala común, NO")
	var xs []*big.Rat
	for _, f := range [][2]int64{{1, 2}, {3, 10}, {7, 12}, {5, 1}, {49, 6}, {1, 1024}, {121, 15}, {2, 3}, {9, 7}, {100, 21}} {
		xs = append(xs, new(big.Rat).SetFrac64(f[0], f[1]))
	}
	fmt.Println("\n      escala      pares probados   violaciones de la ultramétrica")
	violP := 0
	pares := 0
	for _, p := range primos {
		v := 0
		for _, x := range xs {
			for _, y := range xs {
				pares++
				sum := new(big.Rat).Add(x, y)
				if tamanoP(sum, int64(p)).Cmp(maxRat(tamanoP(x, int64(p)), tamanoP(y, int64(p)))) > 0 {
					v++
				}
			}
		}
		violP += v
		fmt.Printf("   la del %2d      %5d              %4d\n", p, len(xs)*len(xs), v)
	}
	violR := 0
	for _, x := range xs {
		for _, y := range xs {
			sum := new(big.Rat).Add(x, y)
			if tamanoReal(sum).Cmp(maxRat(tamanoReal(x), tamanoReal(y))) > 0 {
				violR++
			}
		}
	}
	fmt.Printf("   la común (∞)   %5d              %4d   ← la escala común NO es ultramétrica\n", len(xs)*len(xs), violR)
	fmt.Printf("\n   → %d violaciones en %d pares sobre las escalas de los primos: NINGUNA.\n", violP, pares)
	fmt.Printf("     %d violaciones sobre la escala común. Son dos familias distintas de escala.\n", violR)
	fmt.Println("\n   Y EL TEOREMA DE OSTROWSKI CIERRA LA LISTA:")
	fmt.Println("        toda manera no trivial de medir tamaño en los racionales es")
	fmt.Println("        O BIEN la común, O BIEN la de algún primo. No hay ninguna otra.")
	fmt.Println("\n   Entonces la respuesta a la pregunta del capitán:")
	fmt.Println("   NO HAY QUE CONSTRUIR NINGÚN PUENTE. Los primos y las escalas no son dos")
	fmt.Println("   familias con puentes entre sí — SON LA MISMA FAMILIA, CONTADA DOS VECES.")
	fmt.Println("   Hay exactamente una escala por primo, más una sola escala extra. Se acabó.")

	// ---- BRIDGE 7: back to the book ----
	fmt.Println("\nPUENTE 7 · Y ASÍ SE LEE EL LIBRO — un factor por escala, y el sobrante es el ½")
	fmt.Println("\n        ξ(s) = ½·s(s−1) · π^(−s/2)Γ(s/2) · Π sobre los primos 1/(1−p⁻ˢ)")
	fmt.Println("                                \\_ la escala extra _/  \\_ una por primo _/")
	fmt.Println("\n      escala            factor que aporta al libro         ¿lleva ½ adentro?")
	for _, p := range primos {
		fmt.Printf("   la del %-2d          1/(1−%d⁻ˢ)                          NO\n", p, p)
	}
	fmt.Println("   …y así una por cada primo, infinitas")
	fmt.Println("   LA EXTRA (∞)      π^(−s/2)·Γ(s/2)                    SÍ — DOS VECES")
	fmt.Println("\n   → la cuenta cierra sin sobrar nada: hay exactamente un factor de Euler por")
	fmt.Println("     escala. Los primos entregan todas las escalas menos una, y esa única escala")
	fmt.Println("     que NO es de ningún primo es justo la que trae el ½ (medido en F242).")
	fmt.Println("     Por eso los primos no saben nada del ½: el ½ es de la escala que les falta.")

	// ---- verdict ----
	fallos := fallos1 + fallos2 + fallos3 + fallos4 + fallos5 + fallosMidy + violP
	fmt.Println("\n════════ VEREDICTO ════════")
	fmt.Printf("Cinco puentes exactos, %d fallos en total sobre todas las verificaciones:\n", fallos)
	fmt.Println("  · los ceros del final ...... v_p(n) leído en la escritura")
	fmt.Println("  · Legendre ................. v_p(n!) = (n − s_p(n))/(p−1)")
	fmt.Println("  · Kummer ................... v_p del combinatorio = los acarreos")
	fmt.Println("  · Lucas .................... el combinatorio mod p, dígito por dígito")
	fmt.Println("  · el período ............... el largo de la repetición = el orden del grupo")
	fmt.Println("\nTodos hacen lo mismo: convierten un hecho ARITMÉTICO del primo en un hecho")
	fmt.Println("de ESCRITURA en base p. Y eso pasa porque, en el fondo, no son dos cosas.")
	fmt.Println("\nLA RESPUESTA CORTA, CAPITÁN: hay puentes, sí — cinco, y exactos. Pero el")
	fmt.Println("resultado bueno es que no hacían falta. EL PRIMO ES LA ESCALA. Ostrowski")
	fmt.Println("dice que la lista completa de escalas es: una por primo, más una sola más.")
	fmt.Println("Y esa única escala sobrante es la que trae el ½ que estamos persiguiendo.")
	fmt.Println("\n⚖️ LO QUE ESTO **NO** ES: ninguno de los cinco puentes toca la ubicación de los")
	fmt.Println("ceros. Son teoremas clásicos, verificados acá para tener el mapa firme y honesto.")
	fmt.Println("Lo nuevo del laboratorio es el ENSAMBLE: ver que la cuenta de escalas cierra")
	fmt.Println("exactamente con la cuenta de factores, y que el sobrante es el ½. ¿El premio?")
	fmt.Println("Todavía no.")

	escribirLamina(fallos, violP, violR, pares, reptend, len(pruebas5), midyOK, primos)
}

func escribirLamina(fallos, violP, violR, pares, reptend, tot5, midyOK int, primos []int) {
	var b strings.Builder
	W, H := 1500.0, 1080.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="26" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🌉 EL PRIMO Y LA ESCALA — ¿hay puentes?</text>
<text x="%.0f" y="74" font-size="15" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">cinco puentes exactos… y el resultado de que no hacían falta</text>
`, W, H, W, H, W/2, W/2)

	// the five bridges
	puentes := []struct{ t, f, c string }{
		{"LOS CEROS DEL FINAL", "v_p(n) = ceros con que n termina en base p", "cuántas veces p divide a n, leído en la escritura"},
		{"LEGENDRE", "v_p(n!) = ( n − s_p(n) ) / (p−1)", "la potencia de p en n! ES la suma de dígitos"},
		{"KUMMER", "v_p(C(m+n,m)) = acarreos al sumar en base p", "llevarse una ES que el primo entre una vez más"},
		{"LUCAS", "C(m,n) ≡ Π C(mᵢ,nᵢ)  (mod p)", "el combinatorio se arma dígito por dígito"},
		{"EL PERÍODO", "período de 1/q en base b = ord_b(q)", "la escritura ve la estructura del grupo del primo"},
	}
	y := 108.0
	for i, p := range puentes {
		fmt.Fprintf(&b, `<rect x="40" y="%.0f" width="1420" height="76" rx="9" fill="#101f36" stroke="#26456e"/>
<text x="66" y="%.0f" font-size="15" font-family="Georgia" fill="#ffd98a">PUENTE %d · %s</text>
<text x="66" y="%.0f" font-size="16" font-family="monospace" fill="#bfe3ff">%s</text>
<text x="820" y="%.0f" font-size="14" font-family="Georgia" fill="#8fb4d9">%s</text>
<text x="1400" y="%.0f" font-size="22" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">✓</text>
`, y, y+26, i+1, p.t, y+56, p.f, y+56, p.c, y+48)
		y += 86
	}

	// Ostrowski panel
	fmt.Fprintf(&b, `<rect x="40" y="%.0f" width="1420" height="196" rx="10" fill="#161a3a" stroke="#5a4fa8"/>
<text x="750" y="%.0f" font-size="21" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">Y ACÁ SE ACABA LA METÁFORA — no hay puente: EL PRIMO **ES** LA ESCALA</text>
<text x="70" y="%.0f" font-size="16" font-family="monospace" fill="#9fd8a8">|x+y|_p ≤ max(|x|_p, |y|_p)   ultramétrica — en la escala de un primo, SIEMPRE</text>
<text x="70" y="%.0f" font-size="15" font-family="monospace" fill="#cfe6ff">      %d violaciones en %d pares probados sobre las escalas de %d primos</text>
<text x="70" y="%.0f" font-size="16" font-family="monospace" fill="#ff9e9e">|x+y|   ≤ max(|x|, |y|)       en la escala común, NO: %d violaciones</text>
<text x="70" y="%.0f" font-size="15" font-family="Georgia" fill="#c9b6ff">OSTROWSKI: toda manera no trivial de medir tamaño en los racionales es O BIEN la común, O BIEN la de algún primo.</text>
<text x="70" y="%.0f" font-size="15" font-family="Georgia" fill="#ffd98a">No hay ninguna otra. Primos y escalas no son dos familias con puentes: SON LA MISMA FAMILIA, CONTADA DOS VECES.</text>
`, y, y+32, y+68, y+94, violP, pares, len(primos), y+124, violR, y+156, y+180)
	y += 206

	// the book panel
	fmt.Fprintf(&b, `<rect x="40" y="%.0f" width="1420" height="150" rx="10" fill="#0f2b22" stroke="#2f7f63"/>
<text x="750" y="%.0f" font-size="19" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">Y ASÍ SE LEE EL LIBRO — un factor por escala, y el sobrante es el ½</text>
<text x="750" y="%.0f" font-size="19" text-anchor="middle" font-family="monospace" fill="#dce8f7">ξ(s) = ½·s(s−1) · π^(−s/2)Γ(s/2) · Π 1/(1−p⁻ˢ)</text>
<text x="750" y="%.0f" font-size="14" text-anchor="middle" font-family="monospace" fill="#c9b6ff">                la escala extra        una por cada primo</text>
<text x="750" y="%.0f" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd98a">los primos entregan TODAS las escalas menos una — y esa única que no es de ningún primo es la que trae el ½ (F242).</text>
<text x="750" y="%.0f" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">por eso los primos no saben nada del ½: el ½ es de la escala que les falta. Nada de esto toca dónde están los ceros. Todavía no.</text>
</svg>
`, y, y+30, y+62, y+82, y+112, y+136)

	if err := os.WriteFile("primo-y-escala.svg", []byte(b.String()), 0o644); err != nil {
		fmt.Println("no pude escribir la lámina:", err)
		return
	}
	fmt.Println("\n🖼️  lámina escrita: primo-y-escala.svg")
}
