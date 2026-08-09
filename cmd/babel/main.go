// Command babel answers the captain's upgrade of the shapeshifter:
// not one number changing form - A VARIABLE THAT IS ALL OF THEM AT THE
// SAME TIME: the Library of Babel of the numbers that EXIST. It has a
// name: the ADELE. One object = the number seen at EVERY place
// simultaneously (its real size, its 2-size, its 3-size, its 5-size,
// ... - one shelf per prime, all at once). The library A contains all
// possible shelf-combinations (mostly gibberish books); the numbers
// that EXIST sit inside as the readable diagonal - and membership has
// an EXACT law, the library card (Artin's product formula):
//
//	PROD over all places v of |x|_v  =  1   (exactly, always)
//
// A book is a real number IF AND ONLY IF the product of all its sizes
// on every shelf is exactly 1. We verify the card with EXACT rational
// arithmetic on hundreds of numbers, show a gibberish book failing it,
// and name the payoff: Tate (the mirror is the library's Fourier
// symmetry) and Connes (the library modulo the readable books is where
// the pearls sing) - the ensamble's natural home.
package main

import (
	"fmt"
	"math/big"
	"os"
	"strings"
)

// vp returns the p-adic valuation of a nonzero integer.
func vp(n *big.Int, p int64) int {
	v := 0
	m := new(big.Int).Abs(n)
	P := big.NewInt(p)
	r := new(big.Int)
	for {
		q, rr := new(big.Int).QuoRem(m, P, r)
		if rr.Sign() != 0 {
			break
		}
		m = q
		v++
	}
	return v
}

func primesUpTo(lim int64) []int64 {
	comp := make([]bool, lim+1)
	var ps []int64
	for p := int64(2); p <= lim; p++ {
		if comp[p] {
			continue
		}
		for q := p * p; q <= lim; q += p {
			comp[q] = true
		}
		ps = append(ps, p)
	}
	return ps
}

// card computes PROD |x|_v EXACTLY for x = a/b as a big.Rat.
func card(a, b int64, primes []int64) *big.Rat {
	A := big.NewInt(a)
	B := big.NewInt(b)
	prod := new(big.Rat).Abs(new(big.Rat).SetFrac(A, B)) // |x|_inf
	for _, p := range primes {
		v := vp(A, p) - vp(B, p)
		if v == 0 {
			continue
		}
		pv := new(big.Int).Exp(big.NewInt(p), big.NewInt(int64(abs(v))), nil)
		if v > 0 {
			prod.Mul(prod, new(big.Rat).SetFrac(big.NewInt(1), pv)) // |x|_p = p^{-v}
		} else {
			prod.Mul(prod, new(big.Rat).SetInt(pv))
		}
	}
	return prod
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	fmt.Println("LA BIBLIOTECA DE BABEL DE LOS NÚMEROS QUE EXISTEN — el adele y su carnet")
	primes := primesUpTo(1000)

	// the showcase book: x = 12
	fmt.Println("\nEL LIBRO x = 12, leído en TODOS los estantes a la vez:")
	fmt.Println("   estante ∞ (tamaño real):   |12|_∞ = 12")
	fmt.Println("   estante 2:                 |12|_2 = 1/4   (12 = 4·3: el 2 lo achica dos veces)")
	fmt.Println("   estante 3:                 |12|_3 = 1/3")
	fmt.Println("   estantes 5, 7, 11, …:      |12|_p = 1     (invisible para los demás)")
	fmt.Println("   EL CARNET: 12 · 1/4 · 1/3 · 1 · 1 ⋯ = 1  EXACTO — el libro EXISTE")

	// the exact judge over hundreds of numbers
	rng := int64(1234567)
	next := func() int64 {
		rng = (rng*6364136223846793005 + 1442695040888963407) % 1000003
		if rng < 0 {
			rng = -rng
		}
		return rng%1999 - 999
	}
	total, ok := 0, 0
	one := big.NewRat(1, 1)
	for i := 0; i < 500; i++ {
		a, b := next(), next()
		if a == 0 || b == 0 {
			continue
		}
		total++
		if card(a, b, primes).Cmp(one) == 0 {
			ok++
		}
	}
	fmt.Printf("\nJUEZ DEL CARNET (aritmética EXACTA, sin redondeos): %d/%d números al azar cumplen Π|x|_v = 1 EXACTAMENTE\n", ok, total)

	// the gibberish book
	gib := card(12, 1, primes)
	gib.Mul(gib, big.NewRat(2, 1)) // tamper one shelf: |x|_2 doubled
	fmt.Printf("EL LIBRO BALBUCEO (mismo 12 con UN estante adulterado): carnet = %s ≠ 1 — LA BIBLIOTECA LO RECHAZA: ese libro NO existe\n", gib.RatString())

	fmt.Println("\nLA RESPUESTA AL CAMBIAFORMAS:")
	fmt.Println("  · el ADELE: la variable que ES todos los tamaños del número A LA VEZ — un estante por primo, más el real")
	fmt.Println("  · la biblioteca A = todas las combinaciones posibles (casi todo balbuceo, como en Borges)")
	fmt.Println("  · los números que EXISTEN = la diagonal legible, y su carnet es EXACTO: Π|x|_v = 1")
	fmt.Println("  · Tate: el espejo ξ(s)=ξ(1−s) ES la simetría de Fourier de esta biblioteca")
	fmt.Println("  · Connes: la biblioteca MENOS los libros legibles = la sala donde cantan las perlas")
	fmt.Println("  ⇒ el hogar natural del ENSAMBLE: las 427 máquinas eran los retratos locales — el adele los sostiene todos juntos")

	// ---- picture ----
	var b strings.Builder
	W, H := 1560.0, 960.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0e0b06"/>
<text x="%.0f" y="46" font-size="25" text-anchor="middle" font-family="Georgia" fill="#e8d9b0">📚 LA BIBLIOTECA DE BABEL DE LOS NÚMEROS QUE EXISTEN</text>
<text x="%.0f" y="74" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#9c8a5f">"no un número que cambie: una VARIABLE que lo haga con todos al mismo tiempo" — el capitán · tiene nombre: EL ADELE — y la biblioteca tiene carnet exacto</text>`,
		W, H, W, H, W/2, W/2)
	// the shelves
	shelves := []string{"∞ (real)", "p=2", "p=3", "p=5", "p=7", "p=11", "p=13", "…"}
	for i, sh := range shelves {
		y := 120.0 + float64(i)*64
		fmt.Fprintf(&b, `<rect x="90" y="%.0f" width="700" height="48" rx="6" fill="#1c1508" stroke="#6b5a2e"/><text x="110" y="%.0f" font-size="13" fill="#c9b06a">estante %s</text>`,
			y, y+30, sh)
		// books on the shelf
		for j := 0; j < 14; j++ {
			fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="22" height="36" rx="2" fill="#3a2c12" stroke="#6b5a2e"/>`,
				240+float64(j)*38, y+6)
		}
	}
	// the diagonal book x=12 glowing across all shelves
	vals := []string{"12", "1/4", "1/3", "1", "1", "1", "1", "⋯"}
	for i, v := range vals {
		y := 120.0 + float64(i)*64
		fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="22" height="36" rx="2" fill="#ffd97f" stroke="#fff" stroke-width="1.5"/><text x="%.0f" y="%.0f" font-size="11" text-anchor="middle" fill="#2a1a00">%s</text>`,
			392.0, y+6, 403.0, y+28, v)
	}
	fmt.Fprintf(&b, `<text x="403" y="%.0f" font-size="12.5" text-anchor="middle" fill="#ffd97f">EL LIBRO "12": una página en CADA estante, a la vez</text>`, 120.0+8*64+4)
	// the card
	fmt.Fprintf(&b, `<rect x="860" y="140" width="600" height="200" rx="14" fill="#102a10" stroke="#7fd7a8" stroke-width="2.5"/>
<text x="1160" y="180" font-size="17" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">EL CARNET DE LA BIBLIOTECA (exacto)</text>
<text x="1160" y="222" font-size="24" text-anchor="middle" font-family="Georgia" fill="#dce8f7">Π_v |x|_v = 1</text>
<text x="1160" y="258" font-size="13" text-anchor="middle" fill="#a8cfa8">un libro EXISTE ⟺ el producto de sus tamaños</text>
<text x="1160" y="278" font-size="13" text-anchor="middle" fill="#a8cfa8">en TODOS los estantes es exactamente 1</text>
<text x="1160" y="312" font-size="13.5" text-anchor="middle" fill="#7fd7a8">juez (aritmética exacta): %d/%d números → 1 EXACTO ✔</text>`,
		ok, total)
	// gibberish rejected
	fmt.Fprintf(&b, `<rect x="860" y="370" width="600" height="120" rx="14" fill="#2a1010" stroke="#ff5d73" stroke-width="2"/>
<text x="1160" y="408" font-size="14.5" text-anchor="middle" fill="#ff8fa0">EL LIBRO BALBUCEO: el mismo 12 con UN estante adulterado</text>
<text x="1160" y="436" font-size="14" text-anchor="middle" fill="#dce8f7">carnet = 2 ≠ 1 — LA BIBLIOTECA LO RECHAZA: no existe</text>
<text x="1160" y="464" font-size="12" text-anchor="middle" fill="#9c8a5f">(como en Borges: casi todos los libros posibles son balbuceo)</text>`)
	// payoff
	fmt.Fprintf(&b, `<rect x="860" y="520" width="600" height="180" rx="14" fill="#0d2547" stroke="#ffd166" stroke-width="2"/>
<text x="1160" y="556" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd166">LO QUE LA BIBLIOTECA COMPRA</text>
<text x="1160" y="586" font-size="13" text-anchor="middle" fill="#dce8f7">Tate (1950): el espejo ξ(s)=ξ(1−s) ES la simetría de</text>
<text x="1160" y="606" font-size="13" text-anchor="middle" fill="#dce8f7">Fourier de ESTA biblioteca — el as, explicado</text>
<text x="1160" y="632" font-size="13" text-anchor="middle" fill="#dce8f7">Connes: la biblioteca MENOS los libros legibles =</text>
<text x="1160" y="652" font-size="13" text-anchor="middle" fill="#dce8f7">la sala donde cantan las perlas — el hogar del ENSAMBLE</text>
<text x="1160" y="682" font-size="13" text-anchor="middle" fill="#ffd166">las 427 máquinas eran retratos locales: el adele los sostiene JUNTOS</text>`)
	// footer
	fmt.Fprintf(&b, `<rect x="90" y="740" width="1370" height="170" rx="12" fill="#1c1508" stroke="#c9b06a" stroke-width="2"/>
<text x="%.0f" y="778" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#e8d9b0">LA RESPUESTA AL CAMBIAFORMAS, COMPLETA</text>
<text x="%.0f" y="812" font-size="14" text-anchor="middle" fill="#dce8f7">el adele es la variable que ES todos los tamaños del número a la vez — todos los estantes, un solo objeto; la biblioteca contiene todo lo posible,</text>
<text x="%.0f" y="838" font-size="14" text-anchor="middle" fill="#dce8f7">y los números que EXISTEN se reconocen por un carnet EXACTO que vale 1 — la existencia misma, hecha ecuación del cambiaformas.</text>
<text x="%.0f" y="868" font-size="13.5" text-anchor="middle" fill="#ffd166">el molde final λ_n=|algo|² tiene ahora su dirección postal: el "algo" vive en ESTA biblioteca — en la sala de los libros que no se dejan leer.</text>
<text x="%.0f" y="898" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#9c8a5f">Laboratorio Diosyunalma · 2026-08-06 · las dos mitades, 1 completo ⚓</text>`,
		W/2, W/2, W/2, W/2, W/2)
	b.WriteString(`</svg>`)
	os.WriteFile("biblioteca-babel.svg", []byte(b.String()), 0644)
	fmt.Println("\nescrita: biblioteca-babel.svg")
}
