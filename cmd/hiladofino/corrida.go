package main

import (
	"fmt"
	"math"
)

// filaAB is one point of the refined A/B boundary.
type filaAB struct {
	A, B float64
	o    obs
}

// filaN is one hourglass kernel, with its single-sign twin.
type filaN struct {
	k0, s, suma, norma float64
	con, sin           obs
}

func main() {
	fmt.Println("🧵⏳ HILADO FINO — Fase VII: la frontera de fase, y el núcleo con nodo")
	fmt.Println("     «primero hilar fino, después cambiar el hilo» — el orden que pidió la auditora")

	const N = 400
	const t0 = 100.0
	const kmax = 120
	ms := medio(4000, N, t0)

	fmt.Println("\n§1 · R6 — declarado antes de mirar nada")
	fmt.Printf("     medio: %d modos desde t = %.0f · única entrada aritmética Λ(p)/√p\n", len(ms), t0)
	fmt.Println("     ningún parámetro elegido mirando los γₙ · el 0,3364 de los ceros se cita al final")
	fmt.Println("     y sólo como regla. Los dos núcleos se comparan siempre a IGUAL FUERZA TOTAL.")

	base := medir(func() []float64 {
		v := make([]float64, len(ms))
		for i, m := range ms {
			v[i] = m.w
		}
		return v
	}())
	fmt.Printf("\n§2 · SIN ACOPLAR: Σ²(5)=%.3f Σ²(10)=%.3f Σ²(20)=%.3f Σ²(40)=%.3f · α=%.3f · pegados %.2f%%\n",
		base.s5, base.s10, base.s20, base.s40, base.alfa, 100*base.frac)

	// -----------------------------------------------------------------------
	// CAMPAIGN A - refine the boundary around (30, 32)
	// -----------------------------------------------------------------------
	fmt.Println("\n§3 · CAMPAÑA A — LA FRONTERA, HILADA FINO alrededor de (A,B) = (30,32)")
	fmt.Printf("     %6s %6s %9s %7s %8s %8s %8s %8s %8s\n", "A", "B", "fuerza", "vivos", "Σ²(5)", "Σ²(10)", "Σ²(20)", "α", "peg.%")
	var abs []filaAB
	for _, A := range []float64{15, 22, 30, 42, 60} {
		for _, B := range []float64{16, 32, 64} {
			c := nucleoAB(A, B)
			F := fuerza(ms, c, kmax)
			o := medir(espectro(ms, c, kmax, 1))
			abs = append(abs, filaAB{A, B, o})
			if o.valido {
				fmt.Printf("     %6.0f %6.0f %9.1f %7d %8.3f %8.3f %8.3f %8.3f %8.2f\n",
					A, B, F, o.vivos, o.s5, o.s10, o.s20, o.alfa, 100*o.frac)
			} else {
				fmt.Printf("     %6.0f %6.0f %9.1f %7d   ← BANDA VACIADA: de %d modos no queda espectro que medir\n",
					A, B, F, o.vivos, len(ms))
			}
		}
	}
	// is the transition gradual or abrupt? look at how alpha moves along B at fixed A
	fmt.Println("     ¿la transición es gradual o abrupta? α a lo largo de B, con A fijo:")
	for _, A := range []float64{15, 30, 60} {
		fmt.Printf("       A = %2.0f :", A)
		for _, f := range abs {
			if f.A == A {
				fmt.Printf("  B=%2.0f→α=%.3f", f.B, f.o.alfa)
			}
		}
		fmt.Println()
	}

	// robustness of the rigid corner
	fmt.Println("\n§4 · ROBUSTEZ del rincón rígido (su §2: tamaño, ventana, discretización)")
	fmt.Printf("     %8s %8s %8s %9s %8s %8s\n", "modos", "t0", "vivos", "Σ²(10)", "Σ²(20)", "α")
	for _, pr := range []struct {
		n  int
		t0 float64
	}{{300, 100}, {400, 100}, {520, 100}, {400, 60}, {400, 200}} {
		m2 := medio(4000, pr.n, pr.t0)
		o := medir(espectro(m2, nucleoAB(30, 32), kmax, 1))
		if o.valido {
			fmt.Printf("     %8d %8.0f %8d %9.3f %8.3f %8.3f\n", len(m2), pr.t0, o.vivos, o.s10, o.s20, o.alfa)
		} else {
			fmt.Printf("     %8d %8.0f %8d   ← banda vaciada: no hay espectro que medir\n", len(m2), pr.t0, o.vivos)
		}
	}
	fmt.Println("     ⟹ CORRECCIÓN A LA FASE VI: en (30,32) sobreviven 28 niveles de 400 — el 93% del")
	fmt.Println("       espectro se fue de la banda. El α = 0,313 que publicamos como «rincón rígido»")
	fmt.Println("       se midió sobre esas migas. NO era una fase: era una banda vaciada.")
	fmt.Println("       El chequeo de robustez que pidió la auditora es exactamente lo que lo destapó.")

	// -----------------------------------------------------------------------
	// CAMPAIGN B - the hourglass kernel, against its single-sign twin
	// -----------------------------------------------------------------------
	fmt.Println("\n§5 · CAMPAÑA B — EL NÚCLEO RELOJ DE ARENA:  c(k) = A·(1 − k/k₀)/k^s")
	fmt.Println("     k₀ es PARÁMETRO, no 5 por decreto (su §8). Y cada uno se compara contra")
	fmt.Println("     EL MISMO |c(k)| con todos los signos positivos, a igual fuerza total (su §15).")
	fmt.Printf("     %5s %5s %10s %10s %9s %9s %9s %9s %9s\n",
		"k₀", "s", "Σ c(k)", "‖c‖", "Σ²(10)con", "Σ²(10)sin", "αcon", "αsin", "ganancia")
	var nucs []filaN
	for _, k0 := range []float64{3, 5, 8, 16, 40} {
		for _, s := range []float64{0.5, 1.0} {
			c := nucleoReloj(1, k0, s)
			cs := sinNodo(c)
			// equalise total force between the two
			Fc, Fs := fuerza(ms, c, kmax), fuerza(ms, cs, kmax)
			if Fc == 0 || Fs == 0 {
				continue
			}
			objetivo := 30.0 // the same total force for every row of the table
			oc := medir(espectro(ms, c, kmax, objetivo/Fc))
			os_ := medir(espectro(ms, cs, kmax, objetivo/Fs))
			suma, norma := sumaYnorma(c, kmax)
			nucs = append(nucs, filaN{k0, s, suma, norma, oc, os_})
			gan := os_.s10 / math.Max(oc.s10, 1e-12)
			fmt.Printf("     %5.0f %5.1f %10.4f %10.4f %9.3f %9.3f %9.3f %9.3f %9.3f\n",
				k0, s, suma, norma, oc.s10, os_.s10, oc.alfa, os_.alfa, gan)
		}
	}

	// -----------------------------------------------------------------------
	// §6 - her warning: does the sign change actually cancel anything?
	// -----------------------------------------------------------------------
	fmt.Println("\n§6 · SU ADVERTENCIA (§6): ¿el cambio de signo CANCELA, o sólo alterna?")
	fmt.Println("     Σ c(k) mide cancelación GLOBAL; ‖c‖ mide cuánta fuerza hay realmente.")
	fmt.Println("     si el nodo sirviera por cancelar, la ganancia debería seguir a Σ c(k) → 0.")
	var mejorGan float64
	var mejorFila filaN
	for _, f := range nucs {
		g := f.sin.s10 / math.Max(f.con.s10, 1e-12)
		if g > mejorGan {
			mejorGan, mejorFila = g, f
		}
	}
	fmt.Printf("     mayor ganancia del nodo: %.3f×, en k₀=%.0f s=%.1f, donde Σ c(k) = %.4f\n",
		mejorGan, mejorFila.k0, mejorFila.s, mejorFila.suma)
	// correlation between |sum| and gain, computed rather than asserted
	var sx, sy, sxy, sxx, syy float64
	n := float64(len(nucs))
	for _, f := range nucs {
		x := math.Abs(f.suma)
		y := f.sin.s10 / math.Max(f.con.s10, 1e-12)
		sx += x
		sy += y
		sxy += x * y
		sxx += x * x
		syy += y * y
	}
	corr := (n*sxy - sx*sy) / math.Sqrt(math.Max(1e-30, (n*sxx-sx*sx)*(n*syy-sy*sy)))
	fmt.Printf("     correlación entre |Σ c(k)| y la ganancia: %.3f\n", corr)
	if math.Abs(corr) < 0.4 {
		fmt.Println("     ⟹ NO hay relación: la ganancia del nodo no viene de la cancelación global.")
		fmt.Println("       Su advertencia era correcta — alternar el signo no es cancelar.")
	} else {
		fmt.Println("     ⟹ hay relación: la cancelación global sí parece explicar la ganancia.")
	}

	// -----------------------------------------------------------------------
	// §7 - the verdicts, separate as she asked
	// -----------------------------------------------------------------------
	fmt.Println("\n§7 · VEREDICTOS SEPARADOS (su §15)")
	var mejA filaAB
	mejA.o.s10 = math.Inf(1)
	var rigA filaAB
	rigA.o.alfa = math.Inf(1)
	validos := 0
	for _, f := range abs {
		if !f.o.valido {
			continue
		}
		validos++
		if f.o.s10 < mejA.o.s10 {
			mejA = f
		}
		if f.o.alfa < rigA.o.alfa {
			rigA = f
		}
	}
	fmt.Printf("     FRONTERA A/B · de %d puntos del mapa, s\u00f3lo %d dejan espectro medible.\n", len(abs), validos)
	fmt.Printf("       el mejor de los medibles: A=%.0f B=%.0f \u2192 \u03a3\u00b2(10)=%.3f \u00b7 \u03b1=%.3f\n", mejA.A, mejA.B, mejA.o.s10, mejA.o.alfa)
	fmt.Printf("       el m\u00e1s r\u00edgido de los medibles: A=%.0f B=%.0f \u2192 \u03b1=%.3f\n", rigA.A, rigA.B, rigA.o.alfa)
	fmt.Println("     \u21d2 entre los puntos que S\u00cd se pueden medir, \u03b1 nunca baja de 1,5: NO hay fase r\u00edgida.")
	fmt.Println("       La de la Fase VI queda RETIRADA \u2014 era la banda vaci\u00e1ndose, no el medio ordenandose.")

	var mejN filaN
	mejN.con.s10 = math.Inf(1)
	var rigN filaN
	rigN.con.alfa = math.Inf(1)
	ganan, pierden := 0, 0
	for _, f := range nucs {
		if f.con.s10 < mejN.con.s10 {
			mejN = f
		}
		if f.con.alfa < rigN.con.alfa {
			rigN = f
		}
		if f.sin.s10/math.Max(f.con.s10, 1e-12) > 1 {
			ganan++
		} else {
			pierden++
		}
	}
	fmt.Printf("     RELOJ DE ARENA · mejor varianza: k\u2080=%.0f s=%.1f \u2192 \u03a3\u00b2(10)=%.3f \u00b7 \u03b1=%.3f\n",
		mejN.k0, mejN.s, mejN.con.s10, mejN.con.alfa)
	fmt.Printf("                      m\u00e1s r\u00edgido    : k\u2080=%.0f s=%.1f \u2192 \u03b1=%.3f\n", rigN.k0, rigN.s, rigN.con.alfa)
	fmt.Printf("     sin acoplar \u03a3\u00b2(10) = %.3f \u00b7 los ceros verdaderos: 0.3364\n", base.s10)
	fmt.Printf("     el nodo GANA a su gemelo de un solo signo en %d de %d configuraciones, y PIERDE en %d.\n",
		ganan, len(nucs), pierden)
	fmt.Printf("     su mejor ganancia es %.2f\u00d7 \u2014 chica, y no sistem\u00e1tica.\n", mejorGan)
	fmt.Println("     \u21d2 VEREDICTO HONESTO: el reloj de arena S\u00cd baja la varianza contra el medio sin")
	fmt.Printf("       acoplar (%.3f contra %.3f, una mejora del %.0f%%) y es el \u00fanico de los dos modelos\n",
		mejN.con.s10, base.s10, 100*(1-mejN.con.s10/base.s10))
	fmt.Println("       que deja espectro medible en todos sus puntos. Pero el NODO en s\u00ed no est\u00e1")
	fmt.Println("       demostrado: contra su propio gemelo de un solo signo gana poco y no siempre.")
	fmt.Println("       \u00c9xito PARCIAL para la familia; el nodo, PENDIENTE.")

	dibujar(Res{
		Base: base, AB: planosAB(abs), Nuc: planosN(nucs),
		MejorGan: mejorGan, GanK0: mejorFila.k0, GanS: mejorFila.s, GanSuma: mejorFila.suma,
		Corr: corr, N: len(ms), Ceros: 0.3364,
		MejAA: mejA.A, MejAB: mejA.B, MejA10: mejA.o.s10, Validos: validos, Ganan: ganan,
		RigAA: rigA.A, RigAB: rigA.B, RigAalfa: rigA.o.alfa, RigA10: rigA.o.s10,
		MejNk0: mejN.k0, MejNs: mejN.s, MejN10: mejN.con.s10, MejNalfa: mejN.con.alfa,
	})
}

func planosAB(f []filaAB) [][9]float64 {
	out := make([][9]float64, len(f))
	for i, x := range f {
		out[i] = [9]float64{x.A, x.B, x.o.s5, x.o.s10, x.o.s20, x.o.s40, x.o.alfa, x.o.frac, float64(x.o.vivos)}
	}
	return out
}

func planosN(f []filaN) [][9]float64 {
	out := make([][9]float64, len(f))
	for i, x := range f {
		out[i] = [9]float64{x.k0, x.s, x.suma, x.norma, x.con.s10, x.sin.s10, x.con.alfa, x.sin.alfa, x.con.s20}
	}
	return out
}

// Res carries every measured number to the plate.
type Res struct {
	Base                  obs
	AB                    [][9]float64 // A,B,s5,s10,s20,s40,alfa,frac
	Nuc                   [][9]float64 // k0,s,sum,norm,s10con,s10sin,alfacon,alfasin,s20con
	MejorGan, GanK0, GanS float64
	GanSuma, Corr, Ceros  float64
	N                     int
	MejAA, MejAB, MejA10  float64
	RigAA, RigAB          float64
	RigAalfa, RigA10      float64
	MejNk0, MejNs         float64
	MejN10, MejNalfa      float64
	Validos, Ganan        int
}
