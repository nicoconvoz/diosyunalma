// Command selffocus builds the self-focusing mirror walker.
//
// The flash, as an algorithm: look at the echoes, find the LAST VISIBLE
// image, position yourself on it; from the new position the landscape
// refocuses (each echo becomes the midpoint of its neighbours — the view
// from one step behind); find the new last visible image, jump again — and
// keep walking until unexplored territory: below the floor any single
// mirror can reach.
//
// The single mirror's floor (Finding 84): stop at the last visible image,
// error ~3.9·10⁻⁷ on Stirling at x = 2. PRE-REGISTERED: the walker must
// descend at least two orders of magnitude below that floor before its
// images stop improving.
//
// Formal lineage: iterated hyperasymptotics (Berry–Howls); the walker is
// its greedy pedestrian version, and the last-image rule is applied
// self-similarly at every level of the walk.
//
// Usage:
//
//	go run ./cmd/selffocus
package main

import (
	"fmt"
	"math"
)

func main() {
	// Stirling's series at x = 2; truth ln Gamma(2) = 0.
	bern := []float64{1.0 / 6, -1.0 / 30, 1.0 / 42, -1.0 / 30, 5.0 / 66,
		-691.0 / 2730, 7.0 / 6, -3617.0 / 510, 43867.0 / 798,
		-174611.0 / 330, 854513.0 / 138, -236364091.0 / 2730,
		8553103.0 / 6, -23749461029.0 / 870, 8615841276005.0 / 14322}
	x := 2.0
	s := (x-0.5)*math.Log(x) - x + 0.5*math.Log(2*math.Pi)
	echoes := []float64{}
	for n := 1; n <= len(bern); n++ {
		s += bern[n-1] / (float64(2*n) * float64(2*n-1) * math.Pow(x, float64(2*n-1)))
		echoes = append(echoes, s)
	}

	fmt.Println("THE SELF-FOCUSING WALKER — Stirling at x = 2, truth = 0")
	fmt.Printf("\nsingle-mirror floor (Finding 84): 3.90e-07\n")
	fmt.Println("\nthe walk: position on the last visible image, refocus, jump again")
	fmt.Println("\n  level   echoes   last visible image   error at the focus")

	best := math.Inf(1)
	bestLevel := 0
	level := 0
	for len(echoes) >= 2 {
		// the image sizes are the differences between adjacent echoes;
		// the last visible image is the smallest one.
		smallest, at := math.Inf(1), 0
		for i := 0; i+1 < len(echoes); i++ {
			if d := math.Abs(echoes[i+1] - echoes[i]); d < smallest {
				smallest, at = d, i
			}
		}
		err := math.Abs(echoes[at])
		if e2 := math.Abs(echoes[at+1]); e2 < err {
			err = e2
		}
		fmt.Printf("   %2d      %2d       %.2e             %.2e\n",
			level, len(echoes), smallest, err)
		if err < best {
			best, bestLevel = err, level
		}
		// reposition: the view from one step behind - midpoints.
		next := make([]float64, len(echoes)-1)
		for i := range next {
			next[i] = (echoes[i] + echoes[i+1]) / 2
		}
		echoes = next
		level++
	}

	fmt.Printf("\ndeepest focus: %.2e at level %d — %.0fx below the single-mirror floor\n",
		best, bestLevel, 3.90e-7/best)
	if best < 3.9e-9 {
		fmt.Println("PRE-REGISTRATION MET: two orders below the floor - unexplored territory.")
	} else {
		fmt.Println("PRE-REGISTRATION FAILED: the walker stalled above the target depth.")
	}
	fmt.Println("\nthe walker applies the last-image rule at every level of itself:")
	fmt.Println("look, position, refocus, jump - down to where no single mirror sees.")
}
