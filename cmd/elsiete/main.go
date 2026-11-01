// Command elsiete proves and measures the captain's seven (F422).
//
// HIS FLASH (2026-11-01): "todos dan 7". Take a prime, add it to itself, and
// whenever the result has two digits, add its digits:
//
//	2  = 2+2 = 4, 4+4 = 8, 8+8 = 16, 1+6 = 7
//	13 = 1+3 = 4, 4+4 = 8, 8+8 = 16, 1+6 = 7
//
// and when the 3 seemed to escape (3, 6, 12, 1+2 = 3, 6, 12...), he opened
// the 12 into two BRANCHES, "1+1 = 2 and 2+2 = 4, todas llegan a 7". And: "aún
// no afecta a los primos, pero puede revelarnos cosas más adelante, por eso
// necesito demostrar que es cierto".
//
// THE RULE, exactly as he plays it, in base b (his is b = 10):
//
//	a number of one digit is doubled;
//	a number of two digits or more has its digits added;
//	(branches) a number of two digits or more may instead be opened into its
//	digits, each one going on as its own branch.
//
// THE THEOREM OF THE SEVEN.
//
//	(1) Without branches, n >= 1 reaches the 7 if and only if 3 does not
//	    divide n.
//	(2) So every prime reaches the 7, except the 3.
//	(3) The 7 is not special: the same holds for 1, 2, 4, 5 and 8; the 3, 6
//	    and 9 are a prison that only multiples of 3 enter.
//	(4) With the captain's branches, EVERY n >= 1 reaches the 7.
//	(5) In base b >= 4 the clock is complete (from every one-digit unit of
//	    b-1 the rule reaches every other unit) exactly when 2 is a primitive
//	    root modulo b-1. Base 10 is one of those bases: that is why the 7
//	    appears for him.
//
// ⚡ THE BRIDGE TO THE PRIMES (what it may reveal later). When b-1 is a prime
// p, the clock is complete exactly when 2 is a primitive root modulo p. Are
// there infinitely many such primes, so that his seven works in infinitely
// many bases p+1? That is Artin's primitive root conjecture (1927), still
// open; Hooley (1967) proved it assuming the GENERALIZED RIEMANN HYPOTHESIS,
// and Artin predicted the share of such primes: his constant, the product
// over primes q of 1 - 1/(q(q-1)) = 0.3739558... The experiment measures
// that share (6).
//
// THE PROOF (on paper, for every number, not only the ones measured).
//
//	The residue modulo 9. Since 10 = 9 + 1, a number and the sum of its
//	digits leave the same remainder when divided by 9 (the old "prueba del
//	9"). Doubling multiplies that remainder by 2. So along the rule the
//	remainder r only ever changes as r -> 2r (mod 9).
//
//	A one-digit number IS its remainder (writing 9 for remainder 0). And the
//	double of a digit d >= 5 is 2d = 10 + (2d - 10), whose digits add to
//	2d - 9: the next one-digit number again equals the remainder. So once the
//	walk is in one digit, it walks exactly the clock r -> 2r modulo 9.
//
//	(1) If 3 divides n, the remainder is 0, 3 or 6, and 2*0 = 0, 2*3 = 6,
//	2*6 = 12 = 3 (mod 9): it stays in {0, 3, 6} forever and never is 7. If 3
//	does not divide n, the remainder is one of 1, 2, 4, 5, 7, 8, and doubling
//	walks them in one cycle, 1 -> 2 -> 4 -> 8 -> 7 -> 5 -> 1, because 2 is a
//	primitive root modulo 9 (its powers 2, 4, 8, 16, 32, 64 leave 2, 4, 8, 7,
//	5, 1). In at most five doublings the remainder is 7, and the one-digit
//	number is the 7.
//	(2) A prime other than 3 is not divisible by 3.
//	(3) The cycle passes through all six faces, not only through the 7.
//	(4) If 3 divides n, the walk falls into 3, 6 or 9 (n reduces to its
//	remainder). The 3 doubles to 6, the 6 to 12, the 9 to 18: in at most two
//	steps a two-digit number appears, 12 or 18, and the captain opens it. The
//	branch of the 1 (in both) is not divisible by 3, and by (1) it reaches
//	the 7. This is the part that is HIS: adding digits respects the remainder
//	by 9, so it can never get out of the prison; opening them does not, and
//	that is the key.
//	(5) The same argument with b - 1 in place of 9: digits add to a number
//	with the same remainder modulo b-1, the one-digit walk is r -> 2r modulo
//	b-1, and it visits every unit exactly when 2 generates them, that is,
//	when 2 is a primitive root modulo b-1.
//
// ⚠️ WHAT IT IS NOT (the trap this laboratory meets again, as in "el uno que
// sobra"): the net catches every prime but 3, and it also catches every
// number that is not a multiple of 3 (4, 10, 25...). It is a property of the
// 9 and of the 2 in base 10, not of the primes. His instinct to look at the
// residue and to open the branches is the right tool; whether it can say
// something about the primes is the question for later, and it starts from
// here, proven.
//
// THE EXPERIMENT. The rule is played LITERALLY, adding real digits, never
// using the residue (so the measurement does not lean on the theorem it
// checks): every n up to -hasta, every prime among them, every target from 1
// to 9, and the bases from 4 to 40.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// sumaCifras is the sum of the digits of x in base b.
func sumaCifras(x, b int) int {
	s := 0
	for x > 0 {
		s += x % b
		x /= b
	}
	return s
}

// cifras are the digits of x in base b: the captain's branches.
func cifras(x, b int) []int {
	var c []int
	for x > 0 {
		c = append(c, x%b)
		x /= b
	}
	return c
}

// paso is the rule on one number: two digits or more are added, one digit
// is doubled.
func paso(x, b int) int {
	if x >= b {
		return sumaCifras(x, b)
	}
	return 2 * x
}

// llegaSinRamas plays the rule from n and says whether it reaches meta, and
// in how many steps (0 when n is already meta).
func llegaSinRamas(n, meta, b int) (bool, int) {
	visto := map[int]bool{}
	x, pasos := n, 0
	for !visto[x] {
		if x == meta {
			return true, pasos
		}
		visto[x] = true
		x = paso(x, b)
		pasos++
	}
	return false, pasos
}

// llegaConRamas plays the rule with the captain's branches: from a number of
// two digits or more, either its digits are added or it opens into one
// branch per digit. It says whether SOME branch reaches meta, and the fewest
// steps it takes (opening counts as a step).
func llegaConRamas(n, meta, b int) (bool, int) {
	dist := map[int]int{n: 0}
	cola := []int{n}
	for len(cola) > 0 {
		x := cola[0]
		cola = cola[1:]
		if x == meta {
			return true, dist[x]
		}
		sig := []int{paso(x, b)}
		if x >= b {
			sig = append(sig, cifras(x, b)...)
		}
		for _, y := range sig {
			if _, ok := dist[y]; !ok {
				dist[y] = dist[x] + 1
				cola = append(cola, y)
			}
		}
	}
	return false, 0
}

func mcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// raizPrimitiva says whether g generates every unit modulo m: g is a unit
// and its order is the number of units.
func raizPrimitiva(g, m int) bool {
	if mcd(g, m) != 1 {
		return false
	}
	unidades := 0
	for k := 1; k < m; k++ {
		if mcd(k, m) == 1 {
			unidades++
		}
	}
	orden, x := 1, g%m
	for x != 1 {
		x = x * g % m
		orden++
	}
	return orden == unidades
}

// relojCompleto plays the rule literally in base b: from every one-digit
// unit of b-1, does it reach every other one?
func relojCompleto(b int) bool {
	for d := 1; d < b; d++ {
		for t := 1; t < b; t++ {
			if mcd(d, b-1) != 1 || mcd(t, b-1) != 1 {
				continue
			}
			if ok, _ := llegaSinRamas(d, t, b); !ok {
				return false
			}
		}
	}
	return true
}

// criba gives, for every number up to n, its smallest prime factor (a prime
// is its own).
func criba(n int) []int {
	menor := make([]int, n+1)
	for i := 2; i <= n; i++ {
		if menor[i] == 0 {
			for j := i; j <= n; j += i {
				if menor[j] == 0 {
					menor[j] = i
				}
			}
		}
	}
	return menor
}

// potencia is a^e modulo m.
func potencia(a, e, m int) int {
	r, a := 1, a%m
	for e > 0 {
		if e&1 == 1 {
			r = r * a % m
		}
		a = a * a % m
		e >>= 1
	}
	return r
}

// dosEsRaizPrimitiva says whether 2 generates every unit modulo the odd
// prime p: its order is p-1 exactly when 2^((p-1)/q) is not 1 for any prime
// q dividing p-1.
func dosEsRaizPrimitiva(p int, menor []int) bool {
	for m := p - 1; m > 1; {
		q := menor[m]
		if potencia(2, (p-1)/q, p) == 1 {
			return false
		}
		for m%q == 0 {
			m /= q
		}
	}
	return true
}

// tabla holds the literal answers for the small numbers every walk falls
// into (below 200 the walk never leaves 200), so ten million of them can be
// played in seconds without changing the rule.
type tabla struct {
	sin, con           [200]bool
	pasosSin, pasosCon [200]int
}

func nuevaTabla(meta int) *tabla {
	t := &tabla{}
	for x := 0; x < 200; x++ {
		t.sin[x], t.pasosSin[x] = llegaSinRamas(x, meta, 10)
		t.con[x], t.pasosCon[x] = llegaConRamas(x, meta, 10)
	}
	return t
}

// sinRamas is llegaSinRamas for any n: a number of 200 or more only ever has
// its digits added (the rule), down below 200.
func (t *tabla) sinRamas(n int) (bool, int) {
	pasos := 0
	for n >= 200 {
		n = sumaCifras(n, 10)
		pasos++
	}
	return t.sin[n], t.pasosSin[n] + pasos
}

// conRamas is llegaConRamas for any n: from 200 or more, one step leads to
// the sum of its digits or to one of them, all below 200 for n below 10^22.
func (t *tabla) conRamas(n int) (bool, int) {
	if n < 200 {
		return t.con[n], t.pasosCon[n]
	}
	mejor, ok := 0, false
	for _, y := range append(cifras(n, 10), sumaCifras(n, 10)) {
		if okY, p := t.conRamas(y); okY && (!ok || p+1 < mejor) {
			mejor, ok = p+1, true
		}
	}
	return ok, mejor
}

func main() {
	hasta := flag.Int("hasta", 10_000_000, "hasta qué número se juega la regla")
	flag.Parse()
	N := *hasta

	menor := criba(N)
	compuesto := make([]bool, N+1)
	for i := 2; i <= N; i++ {
		compuesto[i] = menor[i] != i
	}

	t7 := nuevaTabla(7)
	var malas1, malas4, primos, primosSin, maxSin, maxCon int
	var escapanSin []int
	pasosSin, pasosCon := map[int]int{}, map[int]int{}
	for n := 1; n <= N; n++ {
		sin, ps := t7.sinRamas(n)
		con, pc := t7.conRamas(n)
		if sin != (n%3 != 0) {
			malas1++
		}
		if !con {
			malas4++
		}
		if sin {
			pasosSin[ps]++
			maxSin = max(maxSin, ps)
		}
		pasosCon[pc]++
		maxCon = max(maxCon, pc)
		if n >= 2 && !compuesto[n] {
			primos++
			if sin {
				primosSin++
			} else {
				escapanSin = append(escapanSin, n)
			}
		}
	}

	fmt.Printf("EL SIETE DEL CAPITÁN (F422) · la regla jugada literal, sumando cifras de verdad, de 1 a %d\n\n", N)
	fmt.Printf("(1) sin ramas, llega al 7 ⟺ 3 no lo divide: %d números que no cumplen\n", malas1)
	fmt.Printf("    pasos hasta el 7 (los que llegan): el que más tarda, %d\n", maxSin)
	fmt.Printf("(2) primos: %d · llegan sin ramas: %d · no llegan: %v\n", primos, primosSin, escapanSin)
	fmt.Printf("(4) con las ramas del capitán, no llegan: %d · el que más tarda, %d pasos\n", malas4, maxCon)
	ver := func(nombre string, h map[int]int) {
		var ks []int
		for k := range h {
			ks = append(ks, k)
		}
		sort.Ints(ks)
		fmt.Printf("    %s:", nombre)
		for _, k := range ks {
			fmt.Printf(" %d pasos %d ·", k, h[k])
		}
		fmt.Println()
	}
	ver("sin ramas", pasosSin)
	ver("con ramas", pasosCon)

	fmt.Printf("\n(3) ¿el 7 es especial? cuántos de 1 a %d llegan a cada cifra, sin ramas:\n", N)
	for meta := 1; meta <= 9; meta++ {
		tm := nuevaTabla(meta)
		llegan, multiplos3 := 0, 0
		for n := 1; n <= N; n++ {
			if ok, _ := tm.sinRamas(n); ok {
				llegan++
				if n%3 == 0 {
					multiplos3++
				}
			}
		}
		fmt.Printf("    al %d: %d (%.4f de todos), de ellos múltiplos de 3: %d\n", meta, llegan, float64(llegan)/float64(N), multiplos3)
	}

	fmt.Printf("\n(5) ¿en qué bases el reloj está completo? (jugado literal, contra el 2 raíz primitiva de b-1)\n")
	coinciden := 0
	var completas []int
	for b := 4; b <= 40; b++ {
		got, quiero := relojCompleto(b), raizPrimitiva(2, b-1)
		if got == quiero {
			coinciden++
		}
		if got {
			completas = append(completas, b)
		}
	}
	fmt.Printf("    coinciden %d de 37 bases · reloj completo en: %v\n", coinciden, completas)

	fmt.Printf("\n(6) el puente con los primos: ¿en qué parte de los primos p el 2 es raíz primitiva (reloj completo en base p+1)?\n")
	const artin = 0.3739558136192
	res := resultado{hasta: N, primos: primos, escapan: escapanSin, maxSin: maxSin, maxCon: maxCon, bases: completas}
	for _, tope := range []int{1000, 100000, N} {
		cuenta, total := 0, 0
		for p := 3; p <= tope; p++ {
			if menor[p] == p {
				total++
				if dosEsRaizPrimitiva(p, menor) {
					cuenta++
				}
			}
		}
		fmt.Printf("    primos hasta %d: %d de %d = %.5f · la constante de Artin %.5f\n", tope, cuenta, total, float64(cuenta)/float64(total), artin)
		res.artin = append(res.artin, [3]float64{float64(tope), float64(cuenta), float64(total)})
	}

	fmt.Printf("\n(6) el encuentro del capitán: la rama que llega espera, y al juntarse se suman (de 1 a %d)\n", N)
	var noTerminan, conJuntadas, maxPasosE, maxJuntadas int
	for n := 1; n <= N; n++ {
		ok, p, j, _ := encuentro(n, 7, 10)
		if !ok {
			noTerminan++
		}
		if j > 0 {
			conJuntadas++
		}
		maxPasosE, maxJuntadas = max(maxPasosE, p), max(maxJuntadas, j)
	}
	_, _, _, del9 := encuentro(9, 7, 10)
	fmt.Printf("    no terminan con todas sus ramas en el 7: %d · se juntaron ramas en %d · a lo sumo %d pasos y %d juntadas\n",
		noTerminan, conJuntadas, maxPasosE, maxJuntadas)
	fmt.Printf("    el 9 del capitán: %s\n", del9)
	res.noTerminan, res.conJuntadas, res.maxPasosE = noTerminan, conJuntadas, maxPasosE

	ruta := filepath.Join("galeria", "laminas", "10-el-telar", "el-siete-del-capitan.svg")
	if err := os.WriteFile(ruta, []byte(lamina(res)), 0o644); err != nil {
		fmt.Println("⚠ no pude escribir la lámina:", err)
		return
	}
	fmt.Println("\nlámina:", filepath.ToSlash(ruta))
}
