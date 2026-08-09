// Command arbol builds the secret family tree of the primes and answers
// the flash: do all primes have ancestors?
//
// In the DIVISIBILITY tree the primes are the fatherless roots: nothing
// divides them; the tree of factorization converges DOWN onto them and
// expands UP into the composites (F110's inheritance).
//
// But there is a second tree, and in it every prime has ancestors: the
// PRATT tree (the very tree used to certify primality). Each prime p has
// its p-1, composite of SMALLER primes - its parents. And since p-1 is
// even for every odd prime, THE NUMBER 2 IS A DIRECT ANCESTOR OF EVERY
// PRIME: 2 is the Adam of the tree. Reading down, every lineage
// converges to 2; reading up, Dirichlet's theorem guarantees every prime
// infinitely many prime descendants (p = 1 mod q has infinite prime
// solutions). Converges at the root; expands without end at the crown.
//
// Here the patriarchal line (largest parent of p-1, iterated) is built
// for every prime up to 2,000,000 and the generations are measured.
//
// Usage:
//
//	go run ./cmd/arbol
package main

import (
	"fmt"
	"math"
)

func main() {
	fmt.Println("EL ÁRBOL — the secret genealogy of the primes")

	const N = 2000000
	// largest prime factor sieve.
	lpf := make([]int32, N+1)
	for p := 2; p <= N; p++ {
		if lpf[p] == 0 {
			for m := p; m <= N; m += p {
				lpf[m] = int32(p)
			}
		}
	}

	// patriarchal depth: depth(2) = 0; depth(p) = 1 + depth(lpf(p-1)).
	depth := make([]int16, N+1)
	depth[2] = 0
	var maxP, maxD = 2, 0
	var chainEnd int
	for p := 3; p <= N; p++ {
		if int(lpf[p]) != p {
			continue // composite
		}
		d := 1 + int(depth[lpf[p-1]])
		depth[p] = int16(d)
		if d > maxD {
			maxD, maxP = d, p
		}
		chainEnd = p
	}
	_ = chainEnd

	// verify Adam: 2 divides every p-1 (trivially true; stated as law).
	fmt.Println("\n  THE LAW OF ADAM: p-1 is even for every odd prime, so 2 is a DIRECT")
	fmt.Println("  ancestor of every prime in existence. All lineages converge at 2;")
	fmt.Println("  Dirichlet guarantees every prime infinite descendants above.")

	// sample lineages.
	fmt.Println("\n  sample patriarchal lineages (largest parent of p-1, iterated):")
	for _, p := range []int{13, 97, 15013, 999983} {
		fmt.Printf("    %d", p)
		q := p
		for q != 2 {
			q = int(lpf[q-1])
			fmt.Printf(" -> %d", q)
		}
		fmt.Printf("   (generation %d)\n", depth[p])
	}

	// the deepest lineage in range.
	fmt.Printf("\n  the deepest lineage up to %d: generation %d, patriarch chain of %d:\n    %d", N, maxD, maxP, maxP)
	q := maxP
	for q != 2 {
		q = int(lpf[q-1])
		fmt.Printf(" -> %d", q)
	}
	fmt.Println()

	// generation census by decade.
	fmt.Println("\n  mean generation by decade (does the tree deepen slowly?):")
	fmt.Println("    primes near   mean gen   ln ln p")
	for _, lo := range []int{1000, 10000, 100000, 1000000} {
		cnt, sum := 0, 0
		for p := lo; p < lo+lo/10 && p <= N; p++ {
			if int(lpf[p]) == p {
				cnt++
				sum += int(depth[p])
			}
		}
		fmt.Printf("    %9d   %8.2f   %7.2f\n", lo, float64(sum)/float64(cnt),
			math.Log(math.Log(float64(lo))))
	}

	// census of the early generations.
	var gens [32]int
	tot := 0
	for p := 2; p <= N; p++ {
		if int(lpf[p]) == p {
			gens[depth[p]]++
			tot++
		}
	}
	fmt.Printf("\n  the generations of the %d primes up to %d:\n    ", tot, N)
	for g := 0; g <= maxD; g++ {
		fmt.Printf("gen%d:%d  ", g, gens[g])
	}
	fmt.Println()
	fmt.Println("\n  the answer, whole: in divisibility the primes are the fatherless")
	fmt.Println("  roots; in the tree of p-1 every prime has ancestors and ALL descend")
	fmt.Println("  from 2 - the tree converges at its Adam and expands infinitely at")
	fmt.Println("  its crown. Read bottom-up it is a genealogy; read top-down it is a")
	fmt.Println("  certificate - the same tree Pratt uses to PROVE primality.")
}
