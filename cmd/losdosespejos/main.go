package main

// LOS DOS ESPEJOS - the captain's flash: put a prime between two mirrors, what
// do we see in front and behind?
//
// The image has an exact translation. Standing at n and looking at distance k
// in both directions gives the pair (n-k, n+k). Both reflections are prime
// exactly when n-k and n+k are a prime pair - and since (n-k) + (n+k) = 2n,
// counting the reflections of n IS counting the Goldbach partitions of 2n.
//
// So the flash asks, without knowing it: how many Goldbach partitions does a
// prime have, compared with its neighbours of the same size?
//
// Verdicts fixed BEFORE running:
//   V1 (trap #15, declared) k must have parity opposite to n or the two
//      reflections are even and dead on arrival. That halves the candidates
//      for EVERY n, prime or composite. No information about primality.
//   V2 (identity, declared as such) reflections of n = Goldbach partitions of
//      2n. That is algebra, not a discovery.
//   V3 (the measurement) Hardy-Littlewood predicts the count tracks the
//      singular series S(n) = prod over odd primes q|n of (q-1)/(q-2).
//      Measured: R(n)/S(n) must be far flatter than R(n) itself.
//   V4 (the punchline, predicted before looking) a prime has NO small odd
//      divisors, so S(p) = (p-1)/(p-2) is essentially 1 - the MINIMUM. So a
//      prime between two mirrors should see FEWER reflections than almost any
//      composite of its size. If the flash has an answer, that is the answer.
//   V5 (the one thing only a prime does) at k = 0 the prime sees ITSELF:
//      2p = p + p is a partition only when p is prime. Exact, and trivial -
//      declared so it is not sold as a finding.
// Goldbach is open since 1742. Nothing here proves anything about it, or RH.

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func criba(n int) []bool {
	es := make([]bool, n+1)
	for i := 2; i <= n; i++ {
		es[i] = true
	}
	for p := 2; p*p <= n; p++ {
		if !es[p] {
			continue
		}
		for m := p * p; m <= n; m += p {
			es[m] = false
		}
	}
	return es
}

// serie returns the Hardy-Littlewood singular factor of n: prod (q-1)/(q-2)
// over the ODD prime divisors of n.
func serie(n int, es []bool) float64 {
	s := 1.0
	m := n
	for m%2 == 0 {
		m /= 2
	}
	for q := 3; q*q <= m; q += 2 {
		if m%q != 0 {
			continue
		}
		if es[q] {
			s *= float64(q-1) / float64(q-2)
		}
		for m%q == 0 {
			m /= q
		}
	}
	if m > 2 && es[m] {
		s *= float64(m-1) / float64(m-2)
	}
	return s
}

// reflejos counts k >= 1 with n-k and n+k both prime (the mirror images of n).
func reflejos(n int, es []bool) int {
	c := 0
	k0 := 1
	if n%2 == 1 {
		k0 = 2
	}
	for k := k0; k <= n-2; k += 2 {
		if es[n-k] && es[n+k] {
			c++
		}
	}
	return c
}

type fila struct {
	n     int
	primo bool
	r     int
	s     float64
	divs  int
}

func main() {
	fmt.Println("🪞🪞 LOS DOS ESPEJOS — un primo puesto entre dos espejos: ¿qué se ve adelante y atrás?")
	fmt.Println()

	fmt.Println("§1 · la imagen, traducida exacto")
	fmt.Println("   Parado en n y mirando a distancia k para los dos lados se ve el par (n−k, n+k).")
	fmt.Println("   Los dos reflejos son primos justo cuando ese par es de primos — y como (n−k)+(n+k) = 2n,")
	fmt.Println("   ⚡ CONTAR LOS REFLEJOS DE n ES CONTAR LAS PARTICIONES DE GOLDBACH DE 2n.")
	fmt.Println("   (V2: eso es álgebra, no un hallazgo. El flash preguntó Goldbach sin saberlo.)")
	fmt.Println()

	fmt.Println("§2 · la trampa, declarada antes de contar (decimoquinta)")
	fmt.Println("   k tiene que tener paridad OPUESTA a n, o los dos reflejos salen pares y mueren al nacer.")
	fmt.Println("   Eso tacha la mitad de las distancias para TODO n, primo o compuesto: no dice nada de primalidad.")
	fmt.Println()

	fmt.Println("§3 · el único reflejo que solo tiene un primo (V5, exacto y trivial)")
	fmt.Println("   A distancia k = 0 el número se ve A SÍ MISMO: 2n = n + n es partición de Goldbach")
	fmt.Println("   solo cuando n es primo. Es cierto, es exclusivo de los primos — y es la definición,")
	fmt.Println("   así que se declara y no se vende. El resto del experimento cuenta k ≥ 1.")
	fmt.Println()

	const LO, HI = 100000, 101000
	es := criba(2*HI + 10)
	var fs []fila
	for n := LO; n < HI; n++ {
		d := 0
		m := n
		for m%2 == 0 {
			m /= 2
		}
		for q := 3; q*q <= m; q += 2 {
			if m%q == 0 && es[q] {
				d++
				for m%q == 0 {
					m /= q
				}
			}
		}
		if m > 2 && es[m] {
			d++
		}
		fs = append(fs, fila{n, es[n], reflejos(n, es), serie(n, es), d})
	}

	fmt.Printf("§4 · medido sobre los %d números del tramo [%d, %d)\n\n", len(fs), LO, HI)
	fmt.Println("   algunos ejemplos, para ver el efecto a ojo:")
	fmt.Println()
	fmt.Println("           n        ¿primo?    divisores impares    S(n)      reflejos")
	ejemplos := []int{100003, 100019, 100057, 100002, 100050, 100170, 100485, 100947}
	for _, e := range ejemplos {
		for _, f := range fs {
			if f.n != e {
				continue
			}
			p := "no"
			if f.primo {
				p = "SÍ"
			}
			fmt.Printf("      %8d        %-3s            %d            %.4f      %5d\n", f.n, p, f.divs, f.s, f.r)
		}
	}
	fmt.Println()

	// V4: primes vs composites
	var rp, rc, sp, sc float64
	var np, nc int
	for _, f := range fs {
		if f.primo {
			rp += float64(f.r)
			sp += f.s
			np++
		} else {
			rc += float64(f.r)
			sc += f.s
			nc++
		}
	}
	fmt.Println("§5 · LA RESPUESTA AL FLASH: ¿quién ve más reflejos?")
	fmt.Println()
	fmt.Println("        grupo                cuántos     S(n) medio     reflejos medios")
	fmt.Printf("        PRIMOS               %5d       %.4f         %8.1f\n", np, sp/float64(np), rp/float64(np))
	fmt.Printf("        compuestos           %5d       %.4f         %8.1f\n", nc, sc/float64(nc), rc/float64(nc))
	fmt.Printf("\n   razón compuestos/primos: %.4f   → el primo ve %.1f %% MENOS reflejos que un compuesto de su tamaño\n",
		(rc/float64(nc))/(rp/float64(np)), 100*(1-(rp/float64(np))/(rc/float64(nc))))
	fmt.Println()
	fmt.Println("   ⚡ Y la razón es exactamente la de F395 dada vuelta: un primo NO TIENE divisores chicos.")
	fmt.Println("      Los divisores de n son los que enriquecen la estructura simétrica; el primo no tiene ninguno,")
	fmt.Println("      así que su factor S(n) es el MÍNIMO posible y ve la menor cantidad de reflejos de su tamaño.")
	fmt.Println()

	// V3: does R/S flatten?
	fmt.Println("§6 · la ley detrás: los reflejos siguen a los divisores (Hardy–Littlewood)")
	fmt.Println()
	fmt.Println("        divisores impares de n     cuántos     S(n) medio     reflejos medios     reflejos / S(n)")
	porD := map[int][]fila{}
	for _, f := range fs {
		porD[f.divs] = append(porD[f.divs], f)
	}
	var ds []int
	for d := range porD {
		ds = append(ds, d)
	}
	sort.Ints(ds)
	for _, d := range ds {
		g := porD[d]
		if len(g) < 15 {
			continue
		}
		var mr, ms float64
		for _, f := range g {
			mr += float64(f.r)
			ms += f.s
		}
		mr /= float64(len(g))
		ms /= float64(len(g))
		fmt.Printf("                 %d                  %5d       %.4f         %8.1f            %8.1f\n", d, len(g), ms, mr, mr/ms)
	}
	// spread comparison
	var vr, vrs []float64
	for _, f := range fs {
		vr = append(vr, float64(f.r))
		vrs = append(vrs, float64(f.r)/f.s)
	}
	fmt.Printf("\n   dispersión relativa de los reflejos crudos:      %.4f\n", cv(vr))
	fmt.Printf("   dispersión relativa de reflejos / S(n):         %.4f   → dividir por S aplana %.1f×\n", cv(vrs), cv(vr)/cv(vrs))
	fmt.Println("   Los divisores explican casi toda la diferencia. Lo que queda es el ruido de los primos.")
	fmt.Println()

	fmt.Println("§7 · lo que se ve y lo que NO se ve")
	fmt.Println()
	fmt.Println("   Lo que se ve, parado en un primo entre dos espejos:")
	fmt.Println("     · a sí mismo, a distancia cero — y eso solo le pasa a un primo")
	fmt.Printf("     · unos %.0f pares de reflejos primos en promedio a este tamaño...\n", rp/float64(np))
	fmt.Println("     · ...que son MENOS que los de casi cualquier compuesto vecino, por no tener divisores chicos")
	fmt.Println()
	fmt.Println("   Lo que NO se ve, y hay que decirlo:")
	fmt.Println("     · que el número de reflejos nunca llegue a cero es la CONJETURA DE GOLDBACH, abierta desde 1742.")
	fmt.Println("       Acá se cuenta cuántos hay; nadie demostró que no pueda haber ninguno.")
	fmt.Println("     · la fórmula que gobierna todo esto es de Hardy y Littlewood (1923) y es una CONJETURA también.")
	fmt.Println("       Lo medido acá es que sus números dan; eso no la demuestra.")
	fmt.Println()
	fmt.Println("   Estructura cerrada no es hipótesis demostrada. Todavía no.")

	dibujar(fs, ds, porD, rp/float64(np), rc/float64(nc))
	fmt.Println("\n🖼️  lámina escrita: galeria/laminas/10-el-telar/los-dos-espejos.svg")
}

func cv(v []float64) float64 {
	var m float64
	for _, x := range v {
		m += x
	}
	m /= float64(len(v))
	s := 0.0
	for _, x := range v {
		s += (x - m) * (x - m)
	}
	return math.Sqrt(s/float64(len(v))) / m
}

func dibujar(fs []fila, ds []int, porD map[int][]fila, mp, mc float64) {
	var b strings.Builder
	W, H := 1060, 660
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`, W, H, W, H)
	b.WriteString(`<rect width="100%" height="100%" fill="#0d1017"/>`)
	txt := func(x, y float64, sz int, fill, anc, s string) {
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="%d" fill="%s" text-anchor="%s" font-family="Georgia,serif">%s</text>`, x, y, sz, fill, anc, s)
	}
	txt(530, 32, 20, "#e8e2d0", "middle", "LOS DOS ESPEJOS — el primo es el que MENOS se ve reflejado")
	txt(530, 54, 12, "#8a8676", "middle", "mirar a distancia k para los dos lados es contar las particiones de Goldbach de 2n · F396")

	// scatter: reflections vs n, primes highlighted
	x0, y0, w, h := 20.0, 76.0, 1020.0, 280.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" fill="#131826" stroke="#2a3040"/>`, x0, y0, w, h)
	txt(x0+w/2, y0+22, 13, "#e8e2d0", "middle", "cuántos reflejos tiene cada número del tramo (los primos en amarillo)")
	rmax := 0
	for _, f := range fs {
		if f.r > rmax {
			rmax = f.r
		}
	}
	for _, f := range fs {
		px := x0 + 40 + float64(f.n-fs[0].n)/float64(len(fs))*(w-70)
		py := y0 + h - 40 - float64(f.r)/float64(rmax)*(h-80)
		if f.primo {
			fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="2.6" fill="#f2c14e"/>`, px, py)
		} else {
			fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="1.3" fill="#5ec8ff" opacity="0.45"/>`, px, py)
		}
	}
	fmt.Fprintf(&b, `<line x1="%.0f" y1="%.1f" x2="%.0f" y2="%.1f" stroke="#f2c14e" stroke-dasharray="5,4"/>`,
		x0+40, y0+h-40-mp/float64(rmax)*(h-80), x0+w-30, y0+h-40-mp/float64(rmax)*(h-80))
	txt(x0+w-36, y0+h-46-mp/float64(rmax)*(h-80), 11, "#f2c14e", "end", fmt.Sprintf("media de los PRIMOS: %.0f", mp))
	fmt.Fprintf(&b, `<line x1="%.0f" y1="%.1f" x2="%.0f" y2="%.1f" stroke="#5ec8ff" stroke-dasharray="5,4"/>`,
		x0+40, y0+h-40-mc/float64(rmax)*(h-80), x0+w-30, y0+h-40-mc/float64(rmax)*(h-80))
	txt(x0+w-36, y0+h-46-mc/float64(rmax)*(h-80)+16, 11, "#5ec8ff", "end", fmt.Sprintf("media de los compuestos: %.0f", mc))
	txt(x0+w/2, y0+h-12, 10, "#8a8676", "middle", "n de 100.000 a 101.000")

	// bars by divisor count
	bx, by, bw, bh := 20.0, 372.0, 640.0, 250.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" fill="#131826" stroke="#2a3040"/>`, bx, by, bw, bh)
	txt(bx+bw/2, by+22, 13, "#e8e2d0", "middle", "más divisores chicos, más reflejos")
	n := 0
	var vals []float64
	for _, d := range ds {
		if len(porD[d]) < 15 {
			continue
		}
		var mr float64
		for _, f := range porD[d] {
			mr += float64(f.r)
		}
		vals = append(vals, mr/float64(len(porD[d])))
	}
	if len(vals) > 0 {
		vmax := 0.0
		for _, v := range vals {
			vmax = math.Max(vmax, v)
		}
		for _, d := range ds {
			if len(porD[d]) < 15 {
				continue
			}
			v := vals[n]
			px := bx + 70 + float64(n)*110
			hh := v / vmax * (bh - 110)
			col := "#6ee7a0"
			if d == 1 {
				col = "#f2c14e"
			}
			fmt.Fprintf(&b, `<rect x="%.1f" y="%.1f" width="70" height="%.1f" fill="%s" opacity="0.85"/>`, px, by+bh-50-hh, hh, col)
			txt(px+35, by+bh-56-hh, 12, "#e8e2d0", "middle", fmt.Sprintf("%.0f", v))
			txt(px+35, by+bh-32, 11, "#c9c3b0", "middle", fmt.Sprintf("%d div.", d))
			n++
		}
	}
	txt(bx+bw/2, by+bh-14, 10, "#8a8676", "middle", "en amarillo la columna donde viven los primos (1 divisor: ellos mismos)")

	// the honest box
	hx := 680.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" fill="#1a1216" stroke="#5a2a30"/>`, hx, by, 360.0, bh)
	txt(hx+180, by+22, 13, "#ff6b6b", "middle", "lo que NO se ve")
	txt(hx+20, by+52, 11, "#c9c3b0", "start", "Que los reflejos nunca lleguen a CERO")
	txt(hx+20, by+70, 11, "#c9c3b0", "start", "es la conjetura de GOLDBACH, abierta")
	txt(hx+20, by+88, 11, "#c9c3b0", "start", "desde 1742. Acá se cuentan cuántos hay;")
	txt(hx+20, by+106, 11, "#c9c3b0", "start", "nadie probó que no pueda haber ninguno.")
	fmt.Fprintf(&b, `<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#5a2a30"/>`, hx+20, by+124, hx+340, by+124)
	txt(hx+20, by+148, 11, "#c9c3b0", "start", "La ley que gobierna el efecto es de")
	txt(hx+20, by+166, 11, "#c9c3b0", "start", "Hardy y Littlewood (1923) y también")
	txt(hx+20, by+184, 11, "#c9c3b0", "start", "es una CONJETURA. Medir que sus")
	txt(hx+20, by+202, 11, "#c9c3b0", "start", "números dan no la demuestra.")
	txt(hx+20, by+230, 11, "#8a8676", "start", "⚠️ trampa 15: la paridad de k tacha la")
	txt(hx+20, by+246, 11, "#8a8676", "start", "mitad para todo n. Todavía no.")
	b.WriteString(`</svg>`)
	ruta := filepath.Join("galeria", "laminas", "10-el-telar", "los-dos-espejos.svg")
	_ = os.MkdirAll(filepath.Dir(ruta), 0o755)
	_ = os.WriteFile(ruta, []byte(b.String()), 0o644)
}
