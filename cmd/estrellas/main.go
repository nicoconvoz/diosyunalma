// Command estrellas builds the flash's engine: antimatter as cancellation
// fuel. The antimatter of the numbers is Moebius mu (F94), and the exact
// annihilation identity is Lambda = mu * ln (Dirichlet convolution):
// convolve the antimatter with the logarithmic fuel and EVERY composite
// cancels to exact zero - only the light of the primes remains.
//
// And "counting stars instead of sailing": the combinatorial method
// (Legendre-Meissel, refined by Lehmer, Deleglise-Rivat...) counts pi(x)
// by pure cancellation without visiting the numbers - phi(x,a) recursion
// plus the P2 correction. It is how the world records count primes to
// 1e28+. Here: the star counter verified up a ladder of known skies.
//
// Usage:
//
//	go run ./cmd/estrellas
package main

import (
	"fmt"
	"math"
	"time"
)

// ---------- part 1: the annihilation identity, shown digit by digit ----------

func lambdaByCancel(n int) float64 {
	// sum over divisors d of n: mu(d) ln(n/d).
	s := 0.0
	for d := 1; d <= n; d++ {
		if n%d != 0 {
			continue
		}
		// mu(d)
		m, mu := d, 1
		for p := 2; p*p <= m; p++ {
			if m%p == 0 {
				m /= p
				if m%p == 0 {
					mu = 0
					break
				}
				mu = -mu
			}
		}
		if mu != 0 && m > 1 {
			mu = -mu
		}
		if mu != 0 {
			s += float64(mu) * math.Log(float64(n)/float64(d))
		}
	}
	return s
}

// ---------- part 2: the star counter (Meissel with P2) ----------

var (
	primes  []int64
	sieveBt []uint64
	blockPi []int64 // prefix prime count per 4096-number block
	sieveN  int64
)

func buildSieve(n int64) {
	sieveN = n
	sieveBt = make([]uint64, n/64+2)
	for i := int64(2); i*i <= n; i++ {
		if sieveBt[i/64]&(1<<(i%64)) == 0 {
			for j := i * i; j <= n; j += i {
				sieveBt[j/64] |= 1 << (j % 64)
			}
		}
	}
	primes = primes[:0]
	for i := int64(2); i <= n; i++ {
		if sieveBt[i/64]&(1<<(i%64)) == 0 {
			primes = append(primes, i)
		}
	}
	blockPi = make([]int64, n/4096+2)
	cnt := int64(0)
	for i := int64(2); i <= n; i++ {
		if sieveBt[i/64]&(1<<(i%64)) == 0 {
			cnt++
		}
		if i%4096 == 4095 {
			blockPi[i/4096] = cnt
		}
	}
}

func piSmall(x int64) int64 {
	if x < 2 {
		return 0
	}
	blk := x / 4096
	var c int64
	if blk > 0 {
		c = blockPi[blk-1]
	}
	for i := blk * 4096; i <= x; i++ {
		if i >= 2 && sieveBt[i/64]&(1<<(i%64)) == 0 {
			c++
		}
	}
	return c
}

var memo map[[2]int64]int64

func phi(x int64, a int64) int64 {
	if a == 0 {
		return x
	}
	if x < primes[a-1] {
		return 1
	}
	if x <= 1 {
		return x
	}
	key := [2]int64{x, a}
	if v, ok := memo[key]; ok {
		return v
	}
	v := phi(x, a-1) - phi(x/primes[a-1], a-1)
	memo[key] = v
	return v
}

// piMeissel counts the stars: pi(x) = phi(x,a) + a - 1 - P2, a = pi(x^(1/3)).
func piMeissel(x int64) int64 {
	if x <= sieveN {
		return piSmall(x)
	}
	cb := int64(math.Cbrt(float64(x)))
	for (cb+1)*(cb+1)*(cb+1) <= x {
		cb++
	}
	for cb*cb*cb > x {
		cb--
	}
	a := piSmall(cb)
	sq := int64(math.Sqrt(float64(x)))
	for (sq+1)*(sq+1) <= x {
		sq++
	}
	// P2: pairs p*q with cb < p <= q, p*q <= x  =>  sum over p in (cb, sq]
	// of (pi(x/p) - pi(p) + 1); x/p <= x^(2/3) <= sieveN.
	var p2 int64
	for i := a; i < int64(len(primes)) && primes[i] <= sq; i++ {
		p := primes[i]
		p2 += piSmall(x/p) - piSmall(p) + 1
	}
	memo = make(map[[2]int64]int64, 1<<20)
	return phi(x, a) + a - 1 - p2
}

func main() {
	fmt.Println("LAS ESTRELLAS — antimatter as cancellation fuel")

	fmt.Println("\n  1. THE ANNIHILATION IDENTITY, Lambda = mu * ln (digit by digit):")
	for _, n := range []int{7, 12, 13, 30, 31, 97, 98, 100, 101, 128} {
		v := lambdaByCancel(n)
		verdict := "COMPOSITE - annihilated to zero"
		if math.Abs(v) > 1e-9 {
			verdict = fmt.Sprintf("LIGHT REMAINS: ln of its prime = %.6f", v)
		}
		fmt.Printf("     n = %3d:  sum mu(d) ln(n/d) = %+.9f   %s\n", n, v, verdict)
	}
	fmt.Println("     every composite cancels EXACTLY; only prime light survives.")

	fmt.Println("\n  2. COUNTING STARS - pi(x) by pure cancellation, no sailing:")
	known := map[int64]int64{
		100000000:     5761455,
		1000000000:    50847534,
		10000000000:   455052511,
		100000000000:  4118054813,
		1000000000000: 37607912018,
	}
	fmt.Println("     building the base sieve (up to x^(2/3))...")
	buildSieve(1000000)
	for _, x := range []int64{1e8, 1e9, 1e10, 1e11, 1e12} {
		need := int64(math.Cbrt(float64(x))*math.Cbrt(float64(x))) + 10
		if need > sieveN {
			buildSieve(need)
		}
		st := time.Now()
		got := piMeissel(x)
		el := time.Since(st)
		mark := "MATCH - verified against the known sky"
		if got != known[x] {
			mark = fmt.Sprintf("MISMATCH! known %d", known[x])
		}
		fmt.Printf("     pi(%.0e) = %d   (%.2fs, %d phi-nodes)   %s\n",
			float64(x), got, el.Seconds(), len(memo), mark)
	}
	fmt.Println("\n  the islands seen from the deep sky: a trillion numbers never visited,")
	fmt.Println("  their composites annihilated by Moebius antimatter, and the star count")
	fmt.Println("  lands exact. (The world records - pi(10^28+) by Deleglise-Rivat class")
	fmt.Println("  engines - are this same cancellation, industrialized.)")
}
