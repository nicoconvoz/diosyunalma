package main

import (
	"fmt"
	"sort"
	"strings"
)

// THE MEETING (the captain, 2026-11-01): "9 = 9+9 = 18, acá se divide, 1+1 =
// 2, 2+2 = 4, 4+4 = 8; el otro ya estaba en 8 cuando lo alcanza; se
// encuentran y siguen juntas, 8 de una más 8 de otra, 16, y finalmente 7".
//
// The rule with the meeting, in base b:
//
//	a one-digit number doubles; a number of two digits or more adds its
//	digits, unless the sum falls in the prison (it shares a factor with b-1:
//	3, 6 or 9 in base 10), and then it opens into branches, one per nonzero
//	digit;
//	the smallest branch moves, the others wait;
//	two branches on the same number join: they add up (8 + 8 = 16) and go on
//	as one;
//	a branch that reaches the goal has arrived; it ends when all have.
//
// Why it always ends in base 10: the free branches all live on the same turn
// of the clock of 9 (1 → 2 → 4 → 8 → 7 → 5), each at a fixed number of hops
// from the 7. A move brings the branch that moves one hop closer and leaves
// the others where they are, and joining two equal branches is the same as
// one more doubling. The total of hops left only goes down, so it ends, and
// it ends in the 7. A branch in the prison reaches 12 or 18 in at most two
// steps and opens.
//
// ⚠️ WHAT FAILED (my own conjecture, refuted by the experiment): I expected
// the branches to meet in another base exactly when 2 is a primitive root
// of b-1. Base 5 says no: there they meet anyway, because opening a number
// mixes the turns again; and in some bases the branches multiply without
// ending. The meeting outside base 10 stays an open question.

// paso1 moves one branch one step: it returns the branches it becomes (one,
// or several when it opens).
func paso1(v, b int) []int {
	if v < b {
		return []int{2 * v}
	}
	if s := sumaCifras(v, b); mcd(s, b-1) == 1 {
		return []int{s}
	}
	var r []int
	for _, d := range cifras(v, b) {
		if d != 0 {
			r = append([]int{d}, r...) // most significant first
		}
	}
	return r
}

// juntar plays the rule with the meeting on a set of branches. With meta >= 0
// it ends when every branch is on meta; with meta < 0, when a single branch
// is left (they all met). It says whether it ended, in how many steps, how
// many joins, and the road, written as the captain writes it.
func juntar(ramas []int, meta, b int) (bool, int, int, string) {
	var camino strings.Builder
	pasos, juntadas := 0, 0
	visto := map[string]bool{}
	for {
		sort.Ints(ramas)
		fin := len(ramas) == 1 && meta < 0
		if meta >= 0 {
			fin = true
			for _, r := range ramas {
				fin = fin && r == meta
			}
		}
		if fin {
			return true, pasos, juntadas, camino.String()
		}
		clave := fmt.Sprint(ramas)
		if visto[clave] {
			return false, pasos, juntadas, camino.String()
		}
		visto[clave] = true

		// the smallest branch that has not arrived moves
		k := 0
		for meta >= 0 && ramas[k] == meta {
			k++
		}
		nuevas := paso1(ramas[k], b)
		resto := append(append([]int{}, ramas[:k]...), ramas[k+1:]...)
		pasos++
		if len(nuevas) > 1 {
			partes := make([]string, len(nuevas))
			for i, d := range nuevas {
				partes[i] = fmt.Sprint(d)
			}
			fmt.Fprintf(&camino, " → ramas %s y %s", strings.Join(partes[:len(partes)-1], ", "), partes[len(partes)-1])
		} else {
			fmt.Fprintf(&camino, " → %d", nuevas[0])
		}
		// each new branch joins the one waiting on its number, if any
		for _, v := range nuevas {
			for {
				i := -1
				for j, r := range resto {
					if r == v {
						i = j
						break
					}
				}
				if i < 0 || v == meta {
					if i >= 0 { // two arrived at the goal: they are together there
						resto = append(resto[:i], resto[i+1:]...)
						juntadas++
					}
					break
				}
				resto = append(resto[:i], resto[i+1:]...)
				juntadas++
				fmt.Fprintf(&camino, " se junta con %d: %d", v, 2*v)
				v = 2 * v
			}
			resto = append(resto, v)
		}
		ramas = resto
	}
}

// encuentro plays the captain's rule with the meeting from n, toward meta.
func encuentro(n, meta, b int) (bool, int, int, string) {
	if n <= 0 {
		return false, 0, 0, ""
	}
	ok, pasos, juntadas, camino := juntar([]int{n}, meta, b)
	return ok, pasos, juntadas, fmt.Sprint(n) + camino
}
