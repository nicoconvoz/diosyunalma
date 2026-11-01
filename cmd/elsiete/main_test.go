package main

import (
	"strings"
	"testing"
)

// The captain's examples, step by step as he wrote them.
func TestLosEjemplosDelCapitan(t *testing.T) {
	// 2 = 2+2 = 4, 4+4 = 8, 8+8 = 16, 1+6 = 7
	if ok, pasos := llegaSinRamas(2, 7, 10); !ok || pasos != 4 {
		t.Errorf("2 llega al 7 en 4 pasos (4, 8, 16, 7): %v en %d", ok, pasos)
	}
	// 13 = 1+3 = 4, 4+4 = 8, 8+8 = 16, 1+6 = 7
	if ok, pasos := llegaSinRamas(13, 7, 10); !ok || pasos != 4 {
		t.Errorf("13 llega al 7 en 4 pasos (4, 8, 16, 7): %v en %d", ok, pasos)
	}
	// 3 = 3+3 = 6, 6+6 = 12, 1+2 = 3... preso, sin ramas no llega
	if ok, _ := llegaSinRamas(3, 7, 10); ok {
		t.Errorf("sin ramas el 3 queda preso en 3, 6, 12, 3")
	}
	// with branches: 12 opens into 1 and 2, and 2, 4, 8, 16, 7
	if ok, pasos := llegaConRamas(3, 7, 10); !ok || pasos != 7 {
		t.Errorf("con ramas el 3 llega: 6, 12, rama 2, 4, 8, 16, 7 (7 pasos): %v en %d", ok, pasos)
	}
}

// (1) Without branches, n reaches the 7 exactly when 3 does not divide n.
func TestSinRamasLlegaSiYSoloSiNoEsMultiploDeTres(t *testing.T) {
	for n := 1; n <= 20000; n++ {
		if ok, _ := llegaSinRamas(n, 7, 10); ok != (n%3 != 0) {
			t.Fatalf("n = %d: llega %v, y 3 lo divide: %v", n, ok, n%3 == 0)
		}
	}
}

// (3) The 7 is not special: every face of the clock of 9 that 2 walks through
// (1, 2, 4, 5, 7, 8) behaves the same; 3, 6 and 9 are the prison's.
func TestElSieteNoEsEspecial(t *testing.T) {
	for _, meta := range []int{1, 2, 4, 5, 8} {
		for n := 1; n <= 5000; n++ {
			if ok, _ := llegaSinRamas(n, meta, 10); ok != (n%3 != 0) {
				t.Fatalf("meta %d, n = %d: llega %v", meta, n, ok)
			}
		}
	}
}

// (4) With the captain's branches, every n >= 1 reaches the 7; the 0 never.
func TestConRamasTodosLlegan(t *testing.T) {
	for n := 1; n <= 20000; n++ {
		if ok, _ := llegaConRamas(n, 7, 10); !ok {
			t.Fatalf("con ramas, %d no llegó al 7", n)
		}
	}
	if ok, _ := llegaConRamas(0, 7, 10); ok {
		t.Errorf("el 0 duplica a 0 para siempre")
	}
}

// (6) THE MEETING (the captain, 2026-11-01): "9 = 9+9 = 18, acá se divide,
// 1+1 = 2, 2+2 = 4, 4+4 = 8; el otro ya estaba en 8 cuando lo alcanza; se
// encuentran y siguen juntas, 8 de una más 8 de otra, 16, y finalmente 7".
func TestElEncuentroDelCapitan(t *testing.T) {
	ok, pasos, juntadas, camino := encuentro(9, 7, 10)
	if !ok || juntadas != 1 || camino != "9 → 18 → ramas 1 y 8 → 2 → 4 → 8 se junta con 8: 16 → 7" {
		t.Errorf("el 9 del capitán: %v en %d pasos, %d juntadas, camino %q", ok, pasos, juntadas, camino)
	}
	ok, _, juntadas, camino = encuentro(3, 7, 10)
	if !ok || juntadas != 1 || camino != "3 → 6 → 12 → ramas 1 y 2 → 2 se junta con 2: 4 → 8 → 16 → 7" {
		t.Errorf("el 3: %v, %d juntadas, camino %q", ok, juntadas, camino)
	}
	if ok, _, juntadas, _ := encuentro(13, 7, 10); !ok || juntadas != 0 {
		t.Errorf("el 13 no se abre: suma 4, 8, 16, 7")
	}
}

// With the meeting, every n >= 1 ends with all its branches together in the 7.
func TestConElEncuentroTodasLasRamasTerminanEnElSiete(t *testing.T) {
	for n := 1; n <= 100000; n++ {
		if ok, _, _, _ := encuentro(n, 7, 10); !ok {
			t.Fatalf("%d no terminó con todas sus ramas en el 7", n)
		}
	}
}

// (5) In base b (b >= 4), without branches, from every one-digit unit of the
// clock of b-1 the rule reaches every other unit exactly when 2 is a
// primitive root modulo b-1. (Base 3 has a clock of one face, where
// "primitive root" says nothing.)
func TestLaBaseDecideSiElRelojEstaCompleto(t *testing.T) {
	for b := 4; b <= 40; b++ {
		if got, quiero := relojCompleto(b), raizPrimitiva(2, b-1); got != quiero {
			t.Errorf("base %d: el reloj completo dio %v, y el 2 raíz primitiva módulo %d es %v", b, got, b-1, quiero)
		}
	}
}

func TestLaRaizPrimitivaDelNueve(t *testing.T) {
	if !raizPrimitiva(2, 9) {
		t.Errorf("el 2 recorre 2, 4, 8, 7, 5, 1: es raíz primitiva módulo 9")
	}
	if raizPrimitiva(2, 7) {
		t.Errorf("módulo 7 el 2 da 2, 4, 1: no recorre el 3, el 5 ni el 6")
	}
}

// (6) The fast test for "2 is a primitive root modulo the prime p" (the
// order through the prime factors of p-1) agrees with the slow one that
// walks the whole order.
func TestArtinRapidoCoincideConElLento(t *testing.T) {
	menor := criba(3000)
	for p := 3; p <= 3000; p++ {
		if menor[p] != p {
			continue
		}
		if got, quiero := dosEsRaizPrimitiva(p, menor), raizPrimitiva(2, p); got != quiero {
			t.Errorf("p = %d: rápido %v, lento %v", p, got, quiero)
		}
	}
}

// The plate tells the theorem: the clock of 9, the prison, the branches, and
// the measured numbers, with the decimal comma of the house.
func TestLaLaminaCuentaElTeorema(t *testing.T) {
	svg := lamina(resultado{hasta: 10000000, primos: 664579, escapan: []int{3}, maxSin: 10, maxCon: 8,
		bases: []int{4, 6, 10}, artin: [][3]float64{{1000, 67, 167}, {10000000, 248491, 664578}}})
	for _, quiero := range []string{"<svg", "EL SIETE DEL CAPITÁN", "1 → 2 → 4 → 8 → 7 → 5", "664.579", "0,37391", "0,37396", "</svg>"} {
		if !strings.Contains(svg, quiero) {
			t.Errorf("la lámina no dice %q", quiero)
		}
	}
}
