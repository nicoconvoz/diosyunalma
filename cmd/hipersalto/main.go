// Command hipersalto is the DeLorean's twin jump drive: the fleet already
// jumps by ZERO address (starship -zero N, the Riemann globe); this
// engine jumps by PRIME address — "take me to prime number N" — landing
// EXACTLY on the Nth prime without enumerating the N-1 before it.
//
// The physics is Finding 106's annihilation: pi(x) counted by Moebius
// cancellation (Meissel + P2) up to the neighborhood of the guess
// x0 = n(ln n + ln ln n - 1), then a local segmented sieve walks the last
// stretch and pins the landing to the digit. Verified against the known
// prime ladder.
//
// Usage:
//
//	go run ./cmd/hipersalto           # the verification ladder
//	go run ./cmd/hipersalto -n 1e9    # jump to prime number one billion
package main

import (
	"flag"
	"fmt"
	"math"
	"time"
)

var (
	primes  []int64
	sieveBt []uint64
	blockPi []int64
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
	key := [2]int64{x, a}
	if v, ok := memo[key]; ok {
		return v
	}
	v := phi(x, a-1) - phi(x/primes[a-1], a-1)
	memo[key] = v
	return v
}

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
	var p2 int64
	for i := a; i < int64(len(primes)) && primes[i] <= sq; i++ {
		p := primes[i]
		p2 += piSmall(x/p) - piSmall(p) + 1
	}
	memo = make(map[[2]int64]int64, 1<<20)
	return phi(x, a) + a - 1 - p2
}

// segment sieves [lo, hi] and returns a composite-bit slice.
func segment(lo, hi int64) []uint64 {
	bits := make([]uint64, (hi-lo)/64+2)
	for _, p := range primes {
		if p*p > hi {
			break
		}
		start := (lo + p - 1) / p * p
		if start < p*p {
			start = p * p
		}
		for j := start; j <= hi; j += p {
			bits[(j-lo)/64] |= 1 << ((j - lo) % 64)
		}
	}
	return bits
}

// hyperjump lands exactly on prime number n.
func hyperjump(n int64) int64 {
	nf := float64(n)
	x0 := int64(nf * (math.Log(nf) + math.Log(math.Log(nf)) - 1))
	need := int64(math.Cbrt(float64(x0))*math.Cbrt(float64(x0))) + 10
	sq := int64(math.Sqrt(float64(x0))) + 2
	if sq > need {
		need = sq
	}
	if need > sieveN {
		buildSieve(need)
	}
	c := piMeissel(x0) // primes up to the guess, by pure cancellation
	// walk the last stretch with a segmented sieve.
	const W = int64(1 << 22)
	x := x0
	for {
		if c >= n {
			// walk backward inside a segment.
			lo := x - W
			bits := segment(lo, x)
			for y := x; y > lo; y-- {
				if y >= 2 && bits[(y-lo)/64]&(1<<((y-lo)%64)) == 0 {
					if c == n {
						return y
					}
					c--
				}
			}
			x = lo
		} else {
			lo := x + 1
			hi := x + W
			bits := segment(lo, hi)
			for y := lo; y <= hi; y++ {
				if bits[(y-lo)/64]&(1<<((y-lo)%64)) == 0 {
					c++
					if c == n {
						return y
					}
				}
			}
			x = hi
		}
	}
}

func main() {
	nFlag := flag.Float64("n", 0, "jump to prime number N")
	flag.Parse()

	fmt.Println("EL HIPERSALTO — jump drive of the prime world")

	if *nFlag > 0 {
		n := int64(*nFlag)
		st := time.Now()
		p := hyperjump(n)
		fmt.Printf("\n  prime #%d = %d   (%.1fs, no enumeration of the %d before it)\n",
			n, p, time.Since(st).Seconds(), n-1)
		return
	}

	known := map[int64]int64{
		1000000:     15485863,
		100000000:   2038074743,
		1000000000:  22801763489,
		10000000000: 252097800623,
	}
	for _, n := range []int64{1e6, 1e8, 1e9, 1e10} {
		st := time.Now()
		p := hyperjump(n)
		mark := "MATCH - the jump landed on the digit"
		if p != known[n] {
			mark = fmt.Sprintf("MISMATCH! known %d", known[n])
		}
		fmt.Printf("  prime #%.0e = %d   (%.1fs)   %s\n", float64(n), p, time.Since(st).Seconds(), mark)
	}
	fmt.Println("\n  the fleet's two globes now both have postal systems: the DeLorean")
	fmt.Println("  jumps by zero address (-zero N), the hyperjump by prime address (-n N).")
	fmt.Println("  no sailing, no enumeration: cancellation flies you there directly.")
}
