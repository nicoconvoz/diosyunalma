// Command climb reframes the plateau: not a closed club but a mountain
// everyone is climbing at a private speed.
//
// Finding 65 found the positive plateau {12, 18, 24, 42} growing while the
// valleys fade. But the four-decade trajectories tell a richer story: EVERY
// gap's deviation is rising toward (or past) zero — the membership question
// "who is in the plateau?" becomes a crossing-time question: "when does each
// gap's deviation cross zero?" The crossing order so far: 12 (before 10⁸),
// 18 (10⁸→10⁹), 24 (10⁹→10¹⁰), 42 (10¹⁰→10¹¹).
//
// This command measures each gap's climb rate per decade, notes whether the
// climb accelerates or decelerates, and issues bracketed forecasts for the
// 10¹² run — recorded BEFORE that walker returns.
//
// PRE-REGISTERED for 10¹²: gaps 12, 18, 24, 42 stay positive with 18 and 24
// still climbing; gap 6 stays negative; gap 12 stays the leader inside
// [+1.9, +2.2]; gap 36 stays the deepest of 6..60 inside [−8.6, −7.6];
// gaps 30 and 54 land inside their printed brackets, crossing or nearly.
//
// Usage:
//
//	go run ./cmd/climb
package main

import "fmt"

// deviations (%) at 10^8, 10^9, 10^10, 10^11 (Findings 54, 63, 65).
var table = map[int][]float64{
	6:  {-0.97, -0.69, -0.53, -0.41},
	12: {1.92, 2.20, 2.21, 2.13},
	18: {-0.70, 0.14, 0.47, 0.69},
	24: {-0.47, -0.10, 0.11, 0.57},
	30: {-6.08, -2.10, -1.11, -0.50},
	36: {-12.13, -12.51, -10.72, -9.39},
	42: {1.92, -0.03, -0.37, 0.55},
	48: {-7.12, -9.07, -6.03, -4.35},
	54: {-18.92, -5.77, -2.34, -1.03},
	60: {-9.24, -7.51, -4.86, -3.80},
}

var order = []int{6, 12, 18, 24, 30, 36, 42, 48, 54, 60}

func main() {
	fmt.Println("THE CLIMB — every gap scales the mountain at its own speed")
	fmt.Println("\n   d    now      climb/decade (last three)   next-step forecast [linear, decelerating]")
	for _, d := range order {
		v := table[d]
		r1, r2, r3 := v[1]-v[0], v[2]-v[1], v[3]-v[2]
		lin := v[3] + r3
		dec := lin
		if r2 != 0 && r3/r2 > 0 && r3/r2 < 1 {
			dec = v[3] + r3*(r3/r2)
		}
		lo, hi := lin, dec
		if lo > hi {
			lo, hi = hi, lo
		}
		fmt.Printf("  %2d  %+6.2f     %+5.2f  %+5.2f  %+5.2f          [%+5.2f, %+5.2f]\n",
			d, v[3], r1, r2, r3, lo, hi)
	}

	fmt.Println("\ncrossing order so far: 12 (before 10^8), 18, 24, 42 — one per decade.")
	fmt.Println("\nPRE-REGISTERED for 10^12 (recorded before that walker returns):")
	fmt.Println("  - 12, 18, 24, 42 stay positive; 18 and 24 still climbing")
	fmt.Println("  - 6 stays negative; 12 stays the leader inside [+1.9, +2.2]")
	fmt.Println("  - 36 stays the deepest of 6..60 inside [-8.6, -7.6]")
	fmt.Println("  - 30 and 54 land inside their printed brackets")
	fmt.Println("\nthe question is no longer WHO is in the plateau - it is WHEN each")
	fmt.Println("gap arrives, and what sets the private speed of every climb.")
}
