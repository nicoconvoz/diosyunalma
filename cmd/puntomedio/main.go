package main

// EL TEOREMA DEL PUNTO MEDIO - Fase XIII. Pure mathematics, no operator, no
// spectrum: the auditor forbade spectral batteries and ordered the thread
// pulled algebraically. The finding and its name are the captain's.
//
// THEOREM (del Punto Medio, Nico): for odd primes p < q,
//     (p+1)/2 = (q-1)/2  <=>  q = p + 2,
// and the shared value c is an integer with p = 2c-1, q = 2c+1.
// Proof: (p+1)/2 = (q-1)/2 <=> p+1 = q-1 <=> q = p+2. Both sides are integers
// because p, q are odd. Then c = (p+1)/2 gives 2c-1 = p and 2c+1 = q. QED.
//
// GENERAL FORM (the relation that CONTAINS the twin case):
// let p < q be odd primes, g = q-p (even), r = g/2, m = (p+q)/2. Then
//   (I)   m is an INTEGER and p = m - r, q = m + r          [pure geometry]
//   (II)  a-(q) - a+(p) = r - 1                              [anchor identity]
//         so g=2 shares the anchor (difference 0), g=4 touches (difference 1),
//         g>4 leaves exactly r-2 = (g-4)/2 integers of hueco - the verified
//         dictionary of Phase XII is the identity (II) read case by case.
//   (III) THE CENTER LAW: for p, q > 3,   3 | m  <=>  g is NOT divisible by 6.
//         Proof: p,q prime > 3 gives p,q != 0 mod 3. m = p + r.
//         If g != 0 mod 6 then r != 0 mod 3, and of the three residues of m
//         only one survives the two exclusions p != 0 and q = p+2r != 0 mod 3:
//         working mod 3, p in {1,2}; q = p + 2r. For r = 1 mod 3: p=1 gives
//         q=0 (forbidden), so p=2, m=p+1=0. For r = 2 mod 3: p=2 gives q=0
//         (forbidden), so p=1, m=p+2=0. Either way 3 | m.
//         If g = 0 mod 6 then r = 0 mod 3 and m = p != 0 mod 3. QED.
//   (IV)  THE HALVING: the anchor map T(n) = (n-1)/2 is the standard bijection
//         odd numbers -> integers, and it HALVES every gap: if q - p = g then
//         T(q) - T(p) = g/2. Twin primes are exactly the pairs whose images
//         are CONSECUTIVE integers. The +-1/2 are the two halves of that map.
//
// This program verifies (I)-(IV) exhaustively and prints the center table.

import "fmt"

func criba(tope int) []bool {
	c := make([]bool, tope+1)
	for i := 2; i <= tope; i++ {
		c[i] = true
	}
	for i := 2; i*i <= tope; i++ {
		if c[i] {
			for j := i * i; j <= tope; j += i {
				c[j] = false
			}
		}
	}
	return c
}

func main() {
	fmt.Println("📐 EL TEOREMA DEL PUNTO MEDIO — Fase XIII, matemática pura")
	fmt.Println("   el hallazgo y el nombre son del capitán; acá se demuestra, se generaliza")
	fmt.Println("   y se verifica de manera exhaustiva. Ningún cero participa de nada.")

	const TOPE = 200000
	esP := criba(TOPE)
	var ps []int
	for n := 3; n <= TOPE; n += 2 {
		if esP[n] {
			ps = append(ps, n)
		}
	}

	// --- (teorema base) exhaustive over ALL odd pairs, not just primes --------
	fmt.Println("\n§1 · EL TEOREMA BASE, verificado sobre TODOS los pares impares ≤ 20000")
	fallas := 0
	for p := 3; p <= 20000; p += 2 {
		for q := p + 2; q <= 20000; q += 2 {
			izq := (p+1)/2 == (q-1)/2
			der := q == p+2
			if izq != der {
				fallas++
			}
		}
	}
	fmt.Printf("   (p+1)/2 = (q−1)/2 ⟺ q = p+2 : %d fallas en ~25 millones de pares\n", fallas)
	fmt.Println("   y NO depende de la primalidad: es identidad de los IMPARES. Lo específico")
	fmt.Println("   de los primos aparece recién en la ley del centro (§3).")

	// --- (II) anchor identity over consecutive prime pairs --------------------
	fmt.Println("\n§2 · LA IDENTIDAD DE LAS ANCLAS (la forma general): a⁻(q) − a⁺(p) = g/2 − 1")
	f2 := 0
	for i := 0; i+1 < len(ps); i++ {
		p, q := ps[i], ps[i+1]
		if (q-1)/2-(p+1)/2 != (q-p)/2-1 {
			f2++
		}
	}
	fmt.Printf("   verificada sobre %d pares consecutivos de primos ≤ %d: %d fallas\n", len(ps)-1, TOPE, f2)
	fmt.Println("   g=2 → diferencia 0 (ancla compartida: el teorema del capitán)")
	fmt.Println("   g=4 → diferencia 1 (se tocan) · g>4 → hueco de (g−4)/2 exacto")

	// --- (III) the center law over ALL prime pairs ----------------------------
	fmt.Println("\n§3 · LA LEY DEL CENTRO — lo que los PRIMOS agregan a la geometría")
	fmt.Println("   para p, q > 3 primos:  3 | m  ⟺  6 ∤ g   (con m = (p+q)/2, g = q−p)")
	var p3 []int
	for _, p := range ps {
		if p > 3 && p <= 6000 {
			p3 = append(p3, p)
		}
	}
	f3, casos := 0, 0
	for i := 0; i < len(p3); i++ {
		for j := i + 1; j < len(p3); j++ {
			m := (p3[i] + p3[j]) / 2
			g := p3[j] - p3[i]
			casos++
			if (m%3 == 0) != (g%6 != 0) {
				f3++
			}
		}
	}
	fmt.Printf("   verificada sobre TODOS los %d pares de primos en (3, 6000]: %d fallas\n", casos, f3)

	// --- the center table by gap class ----------------------------------------
	fmt.Println("\n§4 · LA TABLA DE LOS CENTROS — la firma mod 6 de cada clase de gap")
	fmt.Println("   (pares CONSECUTIVOS ≤ 200000; residuos del centro m observados)")
	fmt.Printf("   %4s %6s %8s %8s %14s %10s\n", "g", "r=g/2", "Δanclas", "hueco", "m mod 6", "pares")
	for _, g := range []int{2, 4, 6, 8, 10, 12, 14, 16, 18} {
		res := map[int]int{}
		n := 0
		for i := 0; i+1 < len(ps); i++ {
			if ps[i] > 3 && ps[i+1]-ps[i] == g {
				res[(ps[i]+ps[i+1])/2%6]++
				n++
			}
		}
		var firmas []int
		for k := 0; k < 6; k++ {
			if res[k] > 0 {
				firmas = append(firmas, k)
			}
		}
		hueco := "—"
		if g >= 4 {
			hueco = fmt.Sprintf("%d", (g-4)/2)
		}
		fmt.Printf("   %4d %6d %8d %8s %14v %10d\n", g, g/2, g/2-1, hueco, firmas, n)
	}
	fmt.Println("   lectura: la clase del gap se LEE en el residuo del centro —")
	fmt.Println("   gemelos en m ≡ 0 (mod 6) · g=4 en m ≡ 3 · g=6 esquiva los múltiplos de 3 ·")
	fmt.Println("   g=8 vuelve a 3 · g=10 vuelve a 0 · g=12 esquiva otra vez. Período 6 en g.")

	// --- (IV) the halving ------------------------------------------------------
	fmt.Println("\n§5 · LA MITAD — el papel exacto de los ±1/2")
	f5 := 0
	for i := 0; i+1 < len(ps); i++ {
		p, q := ps[i], ps[i+1]
		if (q-1)/2-(p-1)/2 != (q-p)/2 {
			f5++
		}
	}
	fmt.Printf("   T(n) = (n−1)/2 es la biyección impares→enteros y DIVIDE todo gap por 2:\n")
	fmt.Printf("   T(q) − T(p) = g/2, verificado en %d pares: %d fallas\n", len(ps)-1, f5)
	fmt.Println("   los gemelos son EXACTAMENTE los pares cuyas imágenes son enteros CONSECUTIVOS.")
	fmt.Println("   Y el ancla compartida es la mitad del punto medio: c = m/2, con 3 | c (§3).")

	fmt.Println("\n§6 · VEREDICTO DE HONESTIDAD (su §14), sin anestesia")
	fmt.Println("   · el teorema base es identidad de los IMPARES: paridad, no primalidad")
	fmt.Println("   · los ±1/2 son ESTRUCTURALES COMO COORDENADA (la biyección que divide gaps")
	fmt.Println("     por 2) y TRIVIALES COMO ARITMÉTICA: no crean estructura nueva")
	fmt.Println("   · la ley del centro (III) SÍ es específica de los primos y NO se reduce a")
	fmt.Println("     q−p=g: clasifica los centros en progresiones mod 6 disjuntas por clase de")
	fmt.Println("     gap. Es matemática CLÁSICA (criba mod 3) — verdadera, demostrada, no nueva")
	fmt.Println("     para el mundo; nueva como organización para este taller")
	fmt.Println("   · la relación general candidata: (p,q) = (m−r, m+r) con la firma mod 6 del")
	fmt.Println("     centro determinada por g mod 6 — y el caso gemelo es la clase m ≡ 0 (mod 6)")

	dibujar13()
}
