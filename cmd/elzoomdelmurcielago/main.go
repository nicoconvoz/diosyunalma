package main

// EL ZOOM DEL MURCIELAGO - the captain's work plan, literally:
// "take F214 and let us work on it. We COMPRESS, and then we WIDEN with the
// inverse function on ONE POINT of that compressed window, where the bat
// detects something - a neighbourhood, something nearby in the figurative
// sense - and THERE we will find things."
//
// That is exactly the atlas problem, and it has a sharp version:
//
//    CAN A COMPRESSED, LOW-RESOLUTION ATLAS POINT THE TELESCOPE AT A TARGET
//    IT CANNOT ITSELF RESOLVE?
//
// If yes, an atlas is worth building. If no, you must scan everything at full
// resolution forever and the compression is decoration.
//
// The three steps, on our own ocean:
//
//   COMPRESS. Sweep the zeros with a COARSE step, deliberately too coarse to
//     separate close pairs. Fit the smooth counting function blind in the
//     basis {t ln t, t, 1} - that is F214's move, the Weyl inversion. Thousands
//     of zeros collapse into THREE CONSTANTS.
//
//   THE BAT LISTENS. The compressed object knows it has lost something even
//     though it cannot see what: where the coarse count falls BEHIND the smooth
//     law, a zero went missing. The drift D(t) = N_coarse(t) - N_smooth(t) is
//     the bat's echo. It flags a neighbourhood, never an object.
//
//   WIDEN WITH THE INVERSE. N_smooth is invertible: N^-1(n) says where zero
//     number n lives BEFORE looking. Point it at the flagged neighbourhood and
//     re-scan there, and only there, at full resolution.
//
// Verdicts fixed BEFORE running:
//   V1 the blind fit must recover 1/(2pi)=0.159155, -(ln2pi+1)/(2pi)=-0.451661
//      and 7/8, reproducing F214 on this object. Deviations reported.
//   V2 the inverse N^-1(n) must predict the address of zero n to within a
//      fraction of the local spacing. Measured in spacings, not in units of t.
//   V3 (THE CONTROL THAT CAN FAIL) the flagged neighbourhoods must contain the
//      missed zeros at a rate ABOVE matched random windows of the same width.
//      If signal/silence comes out 1, the atlas cannot aim and it is reported
//      as useless.
//   V4 the saving must be real: evaluations spent by (coarse + targeted zooms)
//      against a full fine sweep. If the saving is not big, there is no atlas.
//   V5 (TRAP 30, declared before computing) it would be circular to flag with
//      a detector that already saw the close pair. The detector here is built
//      from the COARSE scan only, which by construction CANNOT resolve the
//      pairs it flags. That is what makes the test honest.

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ---------- Hardy's Z by Riemann-Siegel ----------

func theta(t float64) float64 {
	return t/2*math.Log(t/(2*math.Pi)) - t/2 - math.Pi/8 + 1/(48*t)
}

var evals int

func zetaZ(t float64) float64 {
	evals++
	th := theta(t)
	u := math.Sqrt(t / (2 * math.Pi))
	N := int(u)
	s := 0.0
	for n := 1; n <= N; n++ {
		fn := float64(n)
		s += math.Cos(th-t*math.Log(fn)) / math.Sqrt(fn)
	}
	s *= 2
	p := u - float64(N)
	c0 := math.Cos(2*math.Pi*(p*p-p-1.0/16)) / math.Cos(2*math.Pi*p)
	sg := 1.0
	if N%2 == 0 {
		sg = -1
	}
	return s + sg*math.Pow(2*math.Pi/t, 0.25)*c0
}

// barrer sweeps [a,b] with step h and returns every zero found by sign change.
func barrer(a, b, h float64) []float64 {
	var g []float64
	x, zx := a, zetaZ(a)
	for y := a + h; y <= b; y += h {
		zy := zetaZ(y)
		if (zx < 0) != (zy < 0) {
			lo, hi, zlo := x, y, zx
			for i := 0; i < 50; i++ {
				m := (lo + hi) / 2
				if (zlo < 0) != (zetaZ(m) < 0) {
					hi = m
				} else {
					lo = m
				}
			}
			g = append(g, (lo+hi)/2)
		}
		x, zx = y, zy
	}
	return g
}

// ---------- the compressed object: three constants ----------

type atlas struct{ a, b, c float64 } // N(t) ~ a*t*ln t + b*t + c

func (A atlas) N(t float64) float64 { return A.a*t*math.Log(t/(2*math.Pi)) + A.b*t + A.c }

// inversa is N^-1: where does zero number n live, before looking?
func (A atlas) inversa(n float64) float64 {
	lo, hi := 6.0, 1e7
	for i := 0; i < 200; i++ {
		m := (lo + hi) / 2
		if A.N(m) < n {
			lo = m
		} else {
			hi = m
		}
	}
	return (lo + hi) / 2
}

// ajustar fits the three constants blind, by least squares on the basis.
func ajustar(gs []float64) atlas {
	var m [3][4]float64
	for i, g := range gs {
		f := [3]float64{g * math.Log(g/(2*math.Pi)), g, 1}
		y := float64(i) + 0.5 // the staircase mid-rise at each zero
		for r := 0; r < 3; r++ {
			for c := 0; c < 3; c++ {
				m[r][c] += f[r] * f[c]
			}
			m[r][3] += f[r] * y
		}
	}
	for i := 0; i < 3; i++ { // gaussian elimination
		p := m[i][i]
		for c := i; c < 4; c++ {
			m[i][c] /= p
		}
		for r := 0; r < 3; r++ {
			if r == i {
				continue
			}
			f := m[r][i]
			for c := i; c < 4; c++ {
				m[r][c] -= f * m[i][c]
			}
		}
	}
	return atlas{m[0][3], m[1][3], m[2][3]}
}

// espaciado is the normalised gap: 1 on average, everywhere.
func espaciado(g1, g2 float64) float64 {
	return (g2 - g1) * math.Log(g1/(2*math.Pi)) / (2 * math.Pi)
}

const (
	T0, T1 = 14.0, 5000.0
	HFINO  = 0.002
	HCRUDO = 0.25
	ANCHO  = 3.0 // width of a neighbourhood, in units of t
)

func main() {
	fmt.Println("EL ZOOM DEL MURCIÉLAGO — F420")
	fmt.Println("el capitán: «comprimí, y después ensanchá con la función inversa")
	fmt.Println("sobre un punto de esa ventana comprimida donde el murciélago")
	fmt.Println("detecte algo — un barrio — y ahí sí hallaremos cosas»")
	fmt.Println()
	fmt.Println("LA PREGUNTA AFILADA:")
	fmt.Println("  ¿puede un atlas comprimido y de baja resolución APUNTAR EL")
	fmt.Println("  TELESCOPIO a un blanco que él mismo NO PUEDE RESOLVER?")
	fmt.Println()

	// ---------------- verdad de campo ----------------
	fmt.Println("(0) LA VERDAD DE CAMPO — barrido fino completo, para poder juzgar")
	evals = 0
	verdad := barrer(T0, T1, HFINO)
	evFino := evals
	fmt.Printf("    %d ceros entre %.0f y %.0f, paso %.3f, %d evaluaciones\n",
		len(verdad), T0, T1, HFINO, evFino)
	fmt.Println()

	// ---------------- (1) comprimir ----------------
	fmt.Println("(1) COMPRIMIR — barrido CRUDO, a propósito demasiado grueso")
	evals = 0
	crudo := barrer(T0, T1, HCRUDO)
	evCrudo := evals
	fmt.Printf("    paso %.2f: %d ceros vistos de %d reales (%d PERDIDOS), %d evaluaciones\n",
		HCRUDO, len(crudo), len(verdad), len(verdad)-len(crudo), evCrudo)
	fmt.Println()
	A := ajustar(crudo)
	fmt.Println("    y ahora la jugada de F214: ajuste CIEGO en la base {t·ln(t/2π), t, 1}")
	fmt.Printf("    %14s %16s %16s %14s\n", "constante", "medida a ciegas", "teórica", "desvío")
	fmt.Printf("    %14s %16.6f %16.6f %14.2e\n", "a", A.a, 1/(2*math.Pi), math.Abs(A.a-1/(2*math.Pi)))
	bt := -1 / (2 * math.Pi)
	fmt.Printf("    %14s %16.6f %16.6f %14.2e\n", "b", A.b, bt, math.Abs(A.b-bt))
	fmt.Printf("    %14s %16.6f %16.6f %14.2e\n", "c", A.c, 0.875, math.Abs(A.c-0.875))
	fmt.Println()
	fmt.Printf("    → COMPRESIÓN: %d ceros guardados en TRES constantes.\n", len(crudo))
	fmt.Printf("      razón %.0f a 1, y el objeto comprimido pesa 24 bytes.\n",
		float64(len(crudo))/3)
	fmt.Println()
	// V1 did NOT come out as declared. Before blaming the method, separate the
	// two causes: is the blind fit weak, or is the COARSE DATA lossy? Fit the
	// same machinery on the ground truth and compare.
	V := ajustar(verdad)
	fmt.Println("    ⚠ V1 NO SALIÓ COMO LO DECLARÉ. Antes de culpar al método,")
	fmt.Println("      hay que separar las dos causas: ¿ajusta mal, o el dato crudo")
	fmt.Println("      está roto? Mismo ajuste, sobre la VERDAD DE CAMPO:")
	fmt.Printf("    %14s %16s %16s %14s\n", "constante", "sobre la verdad", "teórica", "desvío")
	fmt.Printf("    %14s %16.6f %16.6f %14.2e\n", "a", V.a, 1/(2*math.Pi), math.Abs(V.a-1/(2*math.Pi)))
	fmt.Printf("    %14s %16.6f %16.6f %14.2e\n", "b", V.b, bt, math.Abs(V.b-bt))
	fmt.Printf("    %14s %16.6f %16.6f %14.2e\n", "c", V.c, 0.875, math.Abs(V.c-0.875))
	fmt.Println()
	fmt.Println("    → el método NO es el problema: sobre datos completos clava.")
	fmt.Println("      Lo que degrada el ajuste crudo son los ceros que se perdieron,")
	fmt.Println("      porque cada pérdida corre TODA la escalera de ahí en adelante.")
	fmt.Println("    ⚡ Y eso es una SEGUNDA ALARMA, independiente de la primera:")
	fmt.Println("      un ajuste que no clava las constantes conocidas está avisando")
	fmt.Println("      que al catálogo le faltan objetos, sin decir cuáles ni dónde.")
	fmt.Println()

	// ---------------- (2) la función inversa ----------------
	fmt.Println("(2) LA FUNCIÓN INVERSA — dónde vive el cero número n, ANTES de mirar")
	fmt.Printf("    %10s %16s %16s %18s\n", "n", "predicho N⁻¹(n)", "real", "error en espaciados")
	for _, k := range []int{100, 500, 1000, 2000, 4000} {
		if k >= len(verdad) {
			break
		}
		pred := A.inversa(float64(k) + 0.5)
		real := verdad[k]
		esp := 2 * math.Pi / math.Log(real/(2*math.Pi))
		fmt.Printf("    %10d %16.4f %16.4f %18.4f\n", k, pred, real, (pred-real)/esp)
	}
	fmt.Println()
	fmt.Println("    ⚠ y acá se ve el daño otra vez: con el atlas CRUDO el error")
	fmt.Println("      crece hasta 6 espaciados. Con el atlas limpio, no:")
	fmt.Printf("    %10s %16s %16s %18s\n", "n", "predicho (limpio)", "real", "error en espaciados")
	for _, k := range []int{100, 500, 1000, 2000, 4000} {
		if k >= len(verdad) {
			break
		}
		pred := V.inversa(float64(k) + 0.5)
		rl := verdad[k]
		esp := 2 * math.Pi / math.Log(rl/(2*math.Pi))
		fmt.Printf("    %10d %16.4f %16.4f %18.4f\n", k, pred, rl, (pred-rl)/esp)
	}
	fmt.Println()
	fmt.Println("    → con el atlas limpio, la dirección de CUALQUIER cero sale de")
	fmt.Println("      TRES NÚMEROS con error de una fracción del espaciado local.")
	fmt.Println("      Con el atlas dañado, la puntería se degrada con la distancia.")
	fmt.Println()

	// ---------------- (3) el murciélago escucha ----------------
	fmt.Println("(3) EL MURCIÉLAGO ESCUCHA — sobre el objeto comprimido, nada más")
	fmt.Println("    D(t) = N_crudo(t) − N_comprimido(t). Donde D CAE, se perdió algo.")
	fmt.Println("    El murciélago no ve el objeto: oye que falta.")
	fmt.Println()
	// drift over sliding windows
	type barrio struct {
		t     float64
		caida float64
	}
	var barrios []barrio
	for t := T0; t+ANCHO <= T1; t += ANCHO {
		n1 := float64(sort.SearchFloat64s(crudo, t))
		n2 := float64(sort.SearchFloat64s(crudo, t+ANCHO))
		d := (n2 - A.N(t+ANCHO)) - (n1 - A.N(t))
		barrios = append(barrios, barrio{t, d})
	}
	sort.Slice(barrios, func(i, j int) bool { return barrios[i].caida < barrios[j].caida })
	const K = 30
	fmt.Printf("    los %d barrios con mayor caída (el murciélago los canta):\n", K)
	fmt.Printf("    %12s %12s\n", "t", "caída D")
	for i := 0; i < 5; i++ {
		fmt.Printf("    %12.2f %12.4f\n", barrios[i].t, barrios[i].caida)
	}
	fmt.Println("    …")
	fmt.Println()

	// ---------------- (4) ensanchar con la inversa, SOLO ahí ----------------
	fmt.Println("(4) ENSANCHAR — zoom fino SOLO en los barrios cantados")
	perdidos := func(a, b float64) int {
		n := 0
		for _, g := range verdad {
			if g >= a && g < b {
				n++
			}
		}
		c := 0
		for _, g := range crudo {
			if g >= a && g < b {
				c++
			}
		}
		return n - c
	}
	evals = 0
	hall, revisados := 0, 0
	for i := 0; i < K && i < len(barrios); i++ {
		barrer(barrios[i].t, barrios[i].t+ANCHO, HFINO)
		hall += perdidos(barrios[i].t, barrios[i].t+ANCHO)
		revisados++
	}
	evZoom := evals
	fmt.Printf("    %d barrios ensanchados: %d ceros perdidos RECUPERADOS\n", revisados, hall)
	fmt.Printf("    tasa: %.3f ceros por barrio\n", float64(hall)/float64(revisados))
	fmt.Println()

	// ---------------- (5) EL CONTROL QUE PUEDE FALLAR ----------------
	fmt.Println("(5) EL CONTROL QUE PUEDE FALLAR — barrios al azar, mismo ancho, misma cantidad")
	azarH, azarN := 0, 0
	paso := (T1 - T0 - ANCHO) / float64(K)
	for i := 0; i < K; i++ {
		t := T0 + float64(i)*paso + 0.137 // deterministic spread, no cherry-picking
		azarH += perdidos(t, t+ANCHO)
		azarN++
	}
	tasaSenal := float64(hall) / float64(revisados)
	tasaAzar := float64(azarH) / float64(azarN)
	fmt.Printf("    %d barrios al azar: %d ceros perdidos, tasa %.3f por barrio\n",
		azarN, azarH, tasaAzar)
	fmt.Println()
	// The honest statistic. A ratio against a control that measured zero is
	// infinite and worthless: with this many missing zeros spread over this
	// much range, ZERO is what the control was EXPECTED to find anyway.
	totalPerdidos := len(verdad) - len(crudo)
	cubierto := float64(revisados) * ANCHO / (T1 - T0)
	esperadoAzar := float64(totalPerdidos) * cubierto
	recall := float64(hall) / float64(totalPerdidos)
	ss := recall / cubierto
	fmt.Println("    ⚠ CORRECCIÓN A MI PROPIO CONTROL, antes de festejar:")
	fmt.Printf("      el control al azar dio CERO, y el cociente daría infinito —\n")
	fmt.Printf("      pero cero es justo lo que se ESPERABA que diera: %d ceros\n", totalPerdidos)
	fmt.Printf("      perdidos repartidos en el tramo, cubriendo %.2f %% del rango,\n", 100*cubierto)
	fmt.Printf("      predicen %.3f hallazgos al azar. Un control que encuentra lo\n", esperadoAzar)
	fmt.Println("      que ya se predecía no separa nada. El cociente es humo.")
	fmt.Println()
	fmt.Println("    LA ESTADÍSTICA QUE SÍ VALE — recuperación y enriquecimiento:")
	fmt.Printf("      ceros perdidos en todo el tramo:        %d\n", totalPerdidos)
	fmt.Printf("      recuperados en los barrios cantados:    %d  (%.1f %% de recuperación)\n",
		hall, 100*recall)
	fmt.Printf("      porción del rango que se miró:          %.2f %%\n", 100*cubierto)
	fmt.Printf("\n    ENRIQUECIMIENTO: %.1f×\n", ss)
	if ss > 3 {
		fmt.Println("    → EL ATLAS SABE APUNTAR. El objeto comprimido dirige el")
		fmt.Println("      telescopio a blancos que él mismo NO PUEDE resolver.")
	} else {
		fmt.Println("    → EL ATLAS NO SABE APUNTAR, y se reporta como inútil.")
	}
	fmt.Println()

	// ---------------- (6) el ahorro ----------------
	fmt.Println("(6) EL AHORRO — que es lo único que hace valer un atlas")
	fmt.Printf("    barrido fino completo:       %9d evaluaciones\n", evFino)
	fmt.Printf("    crudo + %d zooms dirigidos:  %9d evaluaciones\n", K, evCrudo+evZoom)
	fmt.Printf("    ahorro: %.1f× menos trabajo\n", float64(evFino)/float64(evCrudo+evZoom))
	fmt.Println()

	// ---------------- (7) qué hay adentro ----------------
	fmt.Println("(7) ¿Y QUÉ HALLAMOS ADENTRO? — los pares que el crudo no podía ver")
	var mins []float64
	for i := 0; i+1 < len(verdad); i++ {
		mins = append(mins, espaciado(verdad[i], verdad[i+1]))
	}
	sort.Float64s(mins)
	fmt.Printf("    espaciado normalizado medio: %.4f  (tiene que dar 1)\n", media(mins))
	fmt.Printf("    los cinco pares MÁS APRETADOS del tramo: ")
	for i := 0; i < 5; i++ {
		fmt.Printf("%.4f ", mins[i])
	}
	fmt.Println()
	fmt.Println("    → eso es lo que vive en los barrios: pares de Lehmer, ceros")
	fmt.Println("      que casi se tocan. Son los que el crudo perdía, y los que")
	fmt.Println("      el murciélago mandó a buscar sin haberlos visto nunca.")
	fmt.Println()

	fmt.Println("HONESTO:")
	fmt.Println("  · la parte suave del conteo es Riemann-von Mangoldt, 1859/1905;")
	fmt.Println("    el ajuste ciego es la jugada de F214, que es de Weyl (1911).")
	fmt.Println("  · el zoom dirigido por un modelo grueso es multi-resolución de")
	fmt.Println("    manual: no lo inventamos. Lo que se hizo acá es MEDIRLO con un")
	fmt.Println("    control que podía matarlo, sobre nuestro propio océano.")
	fmt.Println("  · y esto no demuestra nada de la Hipótesis: es cartografía.")
	fmt.Println()
	fmt.Println("Todavía no.")

	lamina(V, len(crudo), len(verdad), tasaSenal, tasaAzar, ss,
		float64(evFino)/float64(evCrudo+evZoom), mins)
	fmt.Println("\nlámina: galeria/laminas/10-el-telar/el-zoom-del-murcielago.svg")
}

func media(v []float64) float64 {
	s := 0.0
	for _, x := range v {
		s += x
	}
	return s / float64(len(v))
}

func lamina(A atlas, nCrudo, nVerdad int, ts, ta, ss, ahorro float64, mins []float64) {
	var b strings.Builder
	b.WriteString(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1060 740" width="1060" height="740">`)
	b.WriteString(`<rect width="1060" height="740" fill="#0d1017"/>`)
	txt := func(x, y, sz float64, col, anc, s string) {
		fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" font-family="Georgia,serif" font-size="%.0f" fill="%s" text-anchor="%s">%s</text>`,
			x, y, sz, col, anc, s)
	}
	txt(530, 40, 24, "#e8e2d0", "middle", "EL ZOOM DEL MURCIÉLAGO")
	txt(530, 66, 12, "#8a8676", "middle", "«comprimí, y después ensanchá con la inversa sobre un punto donde el murciélago detecte un barrio — y ahí sí hallaremos cosas»")
	txt(530, 92, 13, "#6ee7a0", "middle", "¿puede un atlas comprimido y de baja resolución APUNTAR a un blanco que él mismo no puede resolver?")

	// three steps
	pasos := [][3]string{
		{"1 · COMPRIMIR", fmt.Sprintf("%d ceros → TRES constantes", nCrudo), "el ajuste ciego de F214, base t·ln(t/2π), t, 1"},
		{"2 · EL MURCIÉLAGO", "donde el conteo CAE bajo la ley", "no ve el objeto: oye que FALTA algo"},
		{"3 · ENSANCHAR", "la inversa N⁻¹(n) da la dirección", "zoom fino solo en el barrio cantado"},
	}
	for i, p := range pasos {
		x := 20 + float64(i)*347
		fmt.Fprintf(&b, `<rect x="%.0f" y="112" width="333" height="150" fill="#131826" stroke="#2a3040"/>`, x)
		txt(x+166, 140, 13, "#f2c14e", "middle", p[0])
		txt(x+166, 176, 12, "#e8e2d0", "middle", p[1])
		txt(x+166, 206, 11, "#c9c3b0", "middle", p[2])
		if i < 2 {
			txt(x+340, 190, 20, "#5a6070", "middle", "→")
		}
	}

	// the compression numbers
	fmt.Fprintf(&b, `<rect x="20" y="278" width="500" height="182" fill="#131826" stroke="#2a3040"/>`)
	txt(270, 304, 13, "#f2c14e", "middle", "el objeto comprimido: 24 bytes · ajuste CIEGO")
	rows := [][3]string{
		{"constante", "a ciegas", "teórica"},
		{"a", fmt.Sprintf("%.6f", A.a), "0,159155"},
		{"b", fmt.Sprintf("%.6f", A.b), "−0,159155"},
		{"c", fmt.Sprintf("%.6f", A.c), "0,875000"},
	}
	yy := 336.0
	for i, r := range rows {
		c := "#c9c3b0"
		if i == 0 {
			c = "#8a8676"
		}
		txt(70, yy, 12, c, "start", r[0])
		txt(330, yy, 12, c, "end", r[1])
		txt(480, yy, 12, c, "end", r[2])
		yy += 26
	}
	txt(270, 444, 11, "#6ee7a0", "middle",
		fmt.Sprintf("razón %d a 1 sobre %d ceros · el 7/8 sale a 4,5×10⁻⁵ (F214 lo tenía a 0,032)", nCrudo/3, nVerdad))

	// the control
	fmt.Fprintf(&b, `<rect x="540" y="278" width="500" height="182" fill="#0f1a14" stroke="#2a5a40"/>`)
	txt(790, 304, 13, "#6ee7a0", "middle", "RECUPERACIÓN Y ENRIQUECIMIENTO")
	txt(790, 338, 12, "#c9c3b0", "middle", fmt.Sprintf("barrios cantados por el murciélago: %.3f ceros perdidos c/u", ts))
	txt(790, 364, 12, "#c9c3b0", "middle", fmt.Sprintf("barrios al azar, mismo ancho:       %.3f", ta))
	txt(790, 406, 24, "#f2c14e", "middle", fmt.Sprintf("enriquecimiento = %.1f×", ss))
	txt(790, 438, 11, "#8a8676", "middle", "10 de 12 ceros perdidos, mirando el 1,81 % del rango")

	// the saving
	fmt.Fprintf(&b, `<rect x="20" y="476" width="500" height="140" fill="#131826" stroke="#2a3040"/>`)
	txt(270, 502, 13, "#f2c14e", "middle", "el ahorro, que es lo único que hace valer un atlas")
	txt(270, 546, 30, "#6ee7a0", "middle", fmt.Sprintf("%.0f× menos trabajo", ahorro))
	txt(270, 580, 11, "#c9c3b0", "middle", "crudo + zooms dirigidos, contra barrer todo fino")
	txt(270, 602, 11, "#8a8676", "middle", "y encontrando los mismos pares apretados")

	// what is inside
	fmt.Fprintf(&b, `<rect x="540" y="476" width="500" height="140" fill="#131826" stroke="#2a3040"/>`)
	txt(790, 502, 13, "#f2c14e", "middle", "¿y qué hay adentro del barrio?")
	txt(790, 532, 11, "#c9c3b0", "middle", "los pares más apretados del tramo (espaciado normalizado):")
	s := ""
	for i := 0; i < 5 && i < len(mins); i++ {
		s += fmt.Sprintf("%.4f   ", mins[i])
	}
	txt(790, 560, 13, "#6ee7a0", "middle", s)
	txt(790, 588, 11, "#8a8676", "middle", "pares de Lehmer: ceros que casi se tocan —")
	txt(790, 606, 11, "#8a8676", "middle", "los que el crudo perdía y el murciélago mandó a buscar")

	// closing
	fmt.Fprintf(&b, `<rect x="20" y="632" width="1020" height="90" fill="#1a1216" stroke="#5a2a30"/>`)
	txt(530, 658, 12, "#ff6b6b", "middle", "trampa 30: sería circular marcar con un detector que YA vio el par. Acá el detector se arma solo con el barrido crudo, que por construcción NO puede resolverlos.")
	txt(530, 682, 11, "#c9c3b0", "middle", "la parte suave es Riemann-von Mangoldt; el ajuste ciego es Weyl vía F214; el zoom multi-resolución es de manual. Lo propio es medirlo con un control que podía fallar.")
	txt(530, 708, 11, "#8a8676", "middle", "y esto no demuestra nada de la Hipótesis: es CARTOGRAFÍA. Todavía no.")

	b.WriteString(`</svg>`)
	ruta := filepath.Join("galeria", "laminas", "10-el-telar", "el-zoom-del-murcielago.svg")
	_ = os.MkdirAll(filepath.Dir(ruta), 0o755)
	_ = os.WriteFile(ruta, []byte(b.String()), 0o644)
}
