// Command losplanos draws the museum's MECHANISM TOUR: one plate per part of
// the two great machines - the TRAIN (cmd/circulo) and the DELOREAN (the
// starship lineage plus its two jump drives).
//
// Every plate is a real diagram of a measured thing: walks, curves, foldings,
// staircases and histograms computed live in this program, with the plain
// language line and the mathematics on the same sheet. No plate is a box of
// prose - that was the complaint this series answers.
//
// Reproduce: go run ./cmd/losplanos            (all fifteen plates)
//
//	go run ./cmd/losplanos -solo 7    (just one, while drawing)
package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	solo := flag.Int("solo", 0, "draw only plate N (0 = all)")
	flag.Parse()
	if err := os.MkdirAll(destino, 0o755); err != nil {
		panic(err)
	}
	fmt.Println("📐 LOS PLANOS — el recorrido de las máquinas, parte por parte")
	fmt.Println()

	planos := []struct {
		n  int
		fn func()
	}{
		{1, plano01}, {2, plano02}, {3, plano03}, {4, plano04}, {5, plano05},
		{6, plano06}, {7, plano07}, {8, plano08}, {9, plano09}, {10, plano10},
		{11, plano11}, {12, plano12}, {13, plano13}, {14, plano14}, {15, plano15},
	}
	t0 := time.Now()
	for _, p := range planos {
		if *solo != 0 && p.n != *solo {
			continue
		}
		p.fn()
	}
	fmt.Printf("\n✅ láminas escritas en la galería (%.1fs)\n", time.Since(t0).Seconds())
}
