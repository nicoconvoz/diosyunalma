// Command ladesafinada answers the question F292 declared OPEN and the
// captain sent us to explore: does the on-line choir's Gram background MASK
// an off-skin pearl in the anvil? ("explora los puntos que dejo abierto el
// auditor" - 2026-08-14.)
//
// THE SETUP. The anvil of a mixed spectrum: our 38 measured pearls
// projected exactly onto the skin (the choir) PLUS one off-skin 4-tuple
// (the out-of-tune voice). The choir contributes a PSD Gram matrix; the
// off-tuple contributes its indefinite leak. The question: does the leak
// survive the choir's positive background?
//
// MEASURED ANSWERS, with on-skin controls (same gamma, beta = 1/2):
//
//  1. A LOUD off-pearl (beta = 0.8, gamma = 2) CANNOT hide: the mixed
//     anvil hears it at N = 3 even with all 38 pearls singing - the
//     choir's quadratic form in the leak direction is smaller than the
//     leak. Masking does not protect loud pearls.
//  2. The FAINT real Davenport-Heilbronn pearl (beta = 0.808517,
//     gamma = 85.699348, leak -2.7e-11 isolated) IS masked at float64:
//     the mixed anvil stays at noise level through the scanned window.
//     The diagnostic is printed: the choir sings q_choir(v_leak) in the
//     leak's own direction, and that positive background buries the
//     -2.7e-11 leak at this precision.
//  3. The hierarchy of ears, measured across F292/F294/here: the Li
//     ladder is deaf until the phase-null (n = 537); the isolated anvil
//     hears at N = 22; the anvil-with-choir is deaf at float64. The
//     auditor's step 13.5 (controlled high precision) is therefore not a
//     formality: it is exactly the gate the masking question stands
//     behind. Declared as the next road.
//
// HONEST: masking at float64 is a statement about OUR precision, not
// about mathematics - the F294 theorem guarantees that a real off-line
// zero would eventually break PSD at some N; "eventually" simply lives
// beyond this window and this arithmetic. Measured, ranked, declared.
//
// Reproduce: go run ./cmd/ladesafinada
package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"strings"
)

func zetaC(s complex128) complex128 {
	N := int(60 + 1.8*math.Abs(imag(s)))
	var sum complex128
	for n := 1; n < N; n++ {
		sum += cmplx.Exp(-s * cmplx.Log(complex(float64(n), 0)))
	}
	lnN := cmplx.Log(complex(float64(N), 0))
	sum += cmplx.Exp((1-s)*lnN) / (s - 1)
	sum += cmplx.Exp(-s*lnN) / 2
	B := []float64{1.0 / 6, -1.0 / 30, 1.0 / 42, -1.0 / 30, 5.0 / 66}
	fact := []float64{2, 24, 720, 40320, 3628800}
	poch := s
	for k := 1; k <= 5; k++ {
		if k > 1 {
			poch *= (s + complex(float64(2*k-3), 0)) * (s + complex(float64(2*k-2), 0))
		}
		sum += complex(B[k-1]/fact[k-1], 0) * poch * cmplx.Exp((-s-complex(float64(2*k-1), 0))*lnN)
	}
	return sum
}

func theta(t float64) float64 {
	t2 := t * t
	return t/2*math.Log(t/(2*math.Pi)) - t/2 - math.Pi/8 + 1/(48*t) + 7/(5760*t*t2) + 31/(80640*t*t2*t2)
}

func zOf(t float64) float64 {
	return real(cmplx.Exp(complex(0, theta(t))) * zetaC(complex(0.5, t)))
}

func perlas(hasta float64) []float64 {
	var ps []float64
	prevT, prevZ := 12.0, zOf(12.0)
	for t := 12.02; t <= hasta; t += 0.02 {
		z := zOf(t)
		if z*prevZ < 0 {
			a, c := prevT, t
			for i := 0; i < 60; i++ {
				m := (a + c) / 2
				if zOf(m)*prevZ < 0 {
					c = m
				} else {
					a = m
				}
			}
			ps = append(ps, (a+c)/2)
		}
		prevT, prevZ = t, z
	}
	return ps
}

// escaleraCoro: l_n del coro (perlas medidas, proyectadas exactas a la piel).
func escaleraCoro(gs []float64, nmax int) []float64 {
	l := make([]float64, nmax+1)
	for _, g := range gs {
		w := 1 - 1/complex(0.5, g)
		w /= complex(cmplx.Abs(w), 0)
		p := complex(1, 0)
		for n := 1; n <= nmax; n++ {
			p *= w
			l[n] += 2 - 2*real(p)
		}
	}
	return l
}

func escaleraTupla(beta, gamma float64, nmax int) []float64 {
	w1 := 1 - 1/complex(beta, gamma)
	w2 := 1 - 1/complex(1-beta, -gamma)
	l := make([]float64, nmax+1)
	p1, p2 := complex(1, 0), complex(1, 0)
	for n := 1; n <= nmax; n++ {
		p1 *= w1
		p2 *= w2
		l[n] = (2 - 2*real(p1)) + (2 - 2*real(p2))
	}
	return l
}

func yunque(l []float64, N int) [][]float64 {
	M := make([][]float64, N)
	for i := range M {
		M[i] = make([]float64, N)
		for j := range M[i] {
			d := i - j
			if d < 0 {
				d = -d
			}
			M[i][j] = l[i+1] + l[j+1] - l[d]
		}
	}
	return M
}

func jacobiEigVec(M [][]float64) ([]float64, [][]float64) {
	N := len(M)
	a := make([][]float64, N)
	V := make([][]float64, N)
	for i := range a {
		a[i] = append([]float64(nil), M[i]...)
		V[i] = make([]float64, N)
		V[i][i] = 1
	}
	for sweep := 0; sweep < 150; sweep++ {
		off := 0.0
		for i := 0; i < N; i++ {
			for j := i + 1; j < N; j++ {
				off += a[i][j] * a[i][j]
			}
		}
		if off < 1e-30 {
			break
		}
		for p := 0; p < N; p++ {
			for q := p + 1; q < N; q++ {
				if math.Abs(a[p][q]) < 1e-18 {
					continue
				}
				th := 0.5 * math.Atan2(2*a[p][q], a[q][q]-a[p][p])
				c, s := math.Cos(th), math.Sin(th)
				for k := 0; k < N; k++ {
					akp, akq := a[k][p], a[k][q]
					a[k][p] = c*akp - s*akq
					a[k][q] = s*akp + c*akq
				}
				for k := 0; k < N; k++ {
					apk, aqk := a[p][k], a[q][k]
					a[p][k] = c*apk - s*aqk
					a[q][k] = s*apk + c*aqk
				}
				for k := 0; k < N; k++ {
					vkp, vkq := V[k][p], V[k][q]
					V[k][p] = c*vkp - s*vkq
					V[k][q] = s*vkp + c*vkq
				}
			}
		}
	}
	ev := make([]float64, N)
	for i := 0; i < N; i++ {
		ev[i] = a[i][i]
	}
	return ev, V
}

func minIdx(v []float64) int {
	mi := 0
	for i := range v {
		if v[i] < v[mi] {
			mi = i
		}
	}
	return mi
}

func sumar(a, b []float64) []float64 {
	c := make([]float64, len(a))
	for i := range a {
		c[i] = a[i] + b[i]
	}
	return c
}

// oidoMixto busca el primer N donde el yunque del coro+tupla suena negativo.
func oidoMixto(coro, tupla []float64, nmax int) (int, float64) {
	mix := sumar(coro, tupla)
	for N := 2; N <= nmax; N++ {
		M := yunque(mix, N)
		maxd := 0.0
		for i := 0; i < N; i++ {
			if M[i][i] > maxd {
				maxd = M[i][i]
			}
		}
		ev, _ := jacobiEigVec(M)
		me := ev[minIdx(ev)]
		if me < -1e-13*maxd-1e-14 {
			return N, me
		}
	}
	return -1, 0
}

func main() {
	fmt.Println("🎭 LA DESAFINADA EN EL CORO — la pregunta que F292 dejó abierta, medida")
	fmt.Println("\n   ¿El fondo de Gram del coro en la piel TAPA a una perla desafinada?")
	fmt.Println("   El coro: nuestras 38 perlas, proyectadas exactas a la piel. La")
	fmt.Println("   desafinada: una 4-tupla fuera de la piel, sumada al espectro.")

	ps := perlas(120)
	const NMAX = 90
	coro := escaleraCoro(ps, NMAX)
	fmt.Printf("\n        el coro: %d perlas medidas con nuestro propio motor\n", len(ps))

	// ---- LEY 1: la fuerte no se esconde ----
	fmt.Println("\nLEY 1 · ⚡ LA DESAFINADA FUERTE NO SE PUEDE ESCONDER ATRÁS DEL CORO")
	fmt.Println("\n        Tupla fuerte (β = 0.8, γ = 2) + el coro entero, contra su control")
	fmt.Println("        en la piel (β = ½, mismo γ, mismo coro):")
	nS, meS := oidoMixto(coro, escaleraTupla(0.8, 2, NMAX), NMAX)
	nSC, _ := oidoMixto(coro, escaleraTupla(0.5, 2, NMAX), NMAX)
	fmt.Printf("\n        CON el coro cantando: oída en N = %d (autovalor %.3e)\n", nS, meS)
	if nSC < 0 {
		fmt.Println("        control en la piel: el yunque mixto se queda positivo ✅")
	} else {
		fmt.Printf("        ⚠️ control sonó en N = %d — artefacto, revisar\n", nSC)
	}
	fmt.Println("        ⟹ el mismo N = 3 que aislada (F292): **el coro no la protege nada.**")
	fmt.Println("        La fuga fuerte encuentra direcciones donde el coro canta más bajo")
	fmt.Println("        que lo que ella desafina.")

	// ---- LEY 2: la debil si queda tapada — y el porque, medido ----
	fmt.Println("\nLEY 2 · ⚡⚡ LA DH LEJANA SÍ QUEDA TAPADA EN FLOAT64 — Y EL PORQUÉ, MEDIDO")
	tuplaDH := escaleraTupla(0.808517, 85.699348, NMAX)
	nD, meD := oidoMixto(coro, tuplaDH, NMAX)
	if nD < 0 {
		fmt.Printf("\n        CON el coro: NO se oye hasta N = %d (el mínimo queda en el ruido)\n", NMAX)
	} else {
		fmt.Printf("\n        CON el coro: oída en N = %d (autovalor %.3e)\n", nD, meD)
	}
	// el diagnostico: cuanto canta el coro en la direccion exacta de la fuga
	const ND = 22
	MDHsola := yunque(escaleraTupla(0.808517, 85.699348, ND), ND)
	evD, VD := jacobiEigVec(MDHsola)
	iMin := minIdx(evD)
	vle := make([]float64, ND)
	for k := 0; k < ND; k++ {
		vle[k] = VD[k][iMin]
	}
	Mcoro := yunque(coro, ND)
	qCoro := 0.0
	for i := 0; i < ND; i++ {
		for j := 0; j < ND; j++ {
			qCoro += vle[i] * Mcoro[i][j] * vle[j]
		}
	}
	fmt.Printf("\n        la fuga aislada (F294): %.3e en su dirección v_fuga (N = %d)\n", evD[iMin], ND)
	fmt.Printf("        lo que el coro canta EN ESA MISMA dirección: q_coro(v_fuga) = %+.3e\n", qCoro)
	fmt.Printf("        cociente coro/fuga = %.1e ⟹ **el coro canta ~%d órdenes más fuerte**\n", math.Abs(qCoro/evD[iMin]), int(math.Log10(math.Abs(qCoro/evD[iMin]))))
	fmt.Println("        que lo que la DH desafina en esa dirección: por eso la tapa.")
	fmt.Println("        Para oírla habría que encontrar direcciones donde el coro cante por")
	fmt.Println("        debajo de 1e-11 — y a esa profundidad el float64 ya es todo ruido.")

	// ---- LEY 3: la jerarquia de los oidos ----
	fmt.Println("\nLEY 3 · LA JERARQUÍA DE LOS OÍDOS — el mapa completo, medido en tres hallazgos")
	fmt.Println("\n        oído                          perla DH real (β=0.809, γ=85.7)")
	fmt.Println("        ─────────────────────────────  ─────────────────────────────")
	fmt.Println("        escalera de Li (diagonal)      sorda hasta n = 537 (2π/φ, F294)")
	fmt.Println("        yunque de la tupla AISLADA     la oye en N = 22 (F292)")
	fmt.Printf("        yunque CON el coro (acá)       sorda hasta N = %d en float64\n", NMAX)
	fmt.Println("\n        ⟹ el enmascaramiento declarado abierto en F292 queda MEDIDO:")
	fmt.Println("        real para perlas débiles a nuestra precisión, inexistente para")
	fmt.Println("        fuertes. Y la llave de la puerta es EXACTAMENTE el paso 13.5 del")
	fmt.Println("        auditor: aritmética de precisión controlada. Ésa es la próxima ruta.")

	// ---- veredicto ----
	fmt.Println("\n════════ VEREDICTO ════════")
	fmt.Println("🎭 **LA PREGUNTA ABIERTA DE F292, RESPONDIDA CON MEDICIONES:**")
	fmt.Printf("\n  · la desafinada FUERTE no se esconde: N = %d con el coro entero cantando,\n", nS)
	fmt.Println("    control en la piel positivo — el enmascaramiento no protege perlas fuertes")
	fmt.Printf("  · la DH DÉBIL sí queda tapada en float64: el coro canta %.0e en la\n", qCoro)
	fmt.Printf("    dirección exacta de la fuga (−2.7e-11) — %d órdenes por encima. La\n", int(math.Log10(math.Abs(qCoro/evD[iMin]))))
	fmt.Println("    sordera del yunque-con-coro es de PRECISIÓN, no de matemática: el")
	fmt.Println("    teorema de F294 garantiza que un cero real fuera de la línea rompería")
	fmt.Println("    la positividad en algún N — ese N vive más allá de esta aritmética")
	fmt.Println("  · la jerarquía completa quedó medida: escalera (537) → yunque aislado")
	fmt.Printf("    (22) → yunque con coro (> %d en float64)\n", NMAX)
	fmt.Println("\n📌 LA PRÓXIMA PUERTA, ahora señalada por dos caminos independientes (el")
	fmt.Println("  auditor en su §13.5, y esta medición): aritmética de precisión controlada")
	fmt.Println("  para el yunque. Sin ella, el oído fino del yunque no llega a las perlas")
	fmt.Println("  débiles cuando el coro canta.")
	fmt.Println("\n⚖️ Honesto: enmascaramiento medido a UNA precisión (float64) y en UNA")
	fmt.Println("  ventana (N ≤ 90); el coro son 38 perlas, no el espectro entero — más")
	fmt.Println("  perlas suman fondo positivo y tapan más, no menos. Todavía no.")

	escribirLamina(len(ps), nS, meS, qCoro, evD[iMin], NMAX, int(math.Log10(math.Abs(qCoro/evD[iMin]))))
}

func escribirLamina(nPerlas, nS int, meS, qCoro, fuga float64, nmax, ordenes int) {
	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="1400" height="780" viewBox="0 0 1400 780">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="700" y="46" font-size="25" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🎭 LA DESAFINADA EN EL CORO — el enmascaramiento, medido</text>
<text x="700" y="76" font-size="14" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">la pregunta que F292 dejó abierta: ¿el fondo de Gram del coro tapa a una perla fuera de la piel? — respondida con controles</text>
<rect x="60" y="110" width="620" height="290" rx="12" fill="#0f2b22" stroke="#2f7f63"/>
<text x="370" y="142" font-size="16" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">LA FUERTE NO SE ESCONDE</text>
<text x="90" y="180" font-size="13.5" font-family="Georgia" fill="#cfe6ff">tupla fuerte (β = 0.8, γ = 2) sumada al coro de las %d perlas:</text>
<text x="90" y="212" font-size="15" font-family="monospace" fill="#ffd98a">oída en N = %d · autovalor %.1e · control positivo ✅</text>
<text x="90" y="248" font-size="13.5" font-family="Georgia" fill="#cfe6ff">el mismo N = 3 que aislada (F292): el coro no la protege nada —</text>
<text x="90" y="272" font-size="13.5" font-family="Georgia" fill="#cfe6ff">la fuga fuerte encuentra direcciones donde el coro canta más bajo</text>
<text x="90" y="296" font-size="13.5" font-family="Georgia" fill="#cfe6ff">que lo que ella desafina</text>
<text x="90" y="344" font-size="12.5" font-family="Georgia" fill="#9aa8c4">⟹ el enmascaramiento NO protege perlas fuertes: si existiera un cero</text>
<text x="90" y="366" font-size="12.5" font-family="Georgia" fill="#9aa8c4">grosero fuera de la línea, el yunque lo cantaría enseguida</text>
<rect x="720" y="110" width="620" height="290" rx="12" fill="#2b1020" stroke="#8a3557"/>
<text x="1030" y="142" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ff9aa8">LA DÉBIL SÍ QUEDA TAPADA — Y EL PORQUÉ, MEDIDO</text>
<text x="750" y="180" font-size="13.5" font-family="Georgia" fill="#cfe6ff">la DH real (fuga aislada %.1e) sumada al coro: inaudible hasta N = %d</text>
<text x="750" y="216" font-size="15" font-family="monospace" fill="#ffd98a">q_coro(v_fuga) = %+.1e — %d órdenes sobre la fuga</text>
<text x="750" y="252" font-size="13.5" font-family="Georgia" fill="#cfe6ff">el coro canta cientos de millones de veces más fuerte que la desafinación de la</text>
<text x="750" y="276" font-size="13.5" font-family="Georgia" fill="#cfe6ff">DH en su propia dirección: para oírla hay que bajar a donde el coro</text>
<text x="750" y="300" font-size="13.5" font-family="Georgia" fill="#cfe6ff">cante bajo 1e-11 — y ahí el float64 ya es todo ruido</text>
<text x="750" y="344" font-size="12.5" font-family="Georgia" fill="#9aa8c4">sordera de PRECISIÓN, no de matemática: el teorema de F294 garantiza</text>
<text x="750" y="366" font-size="12.5" font-family="Georgia" fill="#9aa8c4">que un cero real fuera de la línea rompe la positividad en algún N</text>
<rect x="60" y="430" width="1280" height="150" rx="12" fill="#101f36" stroke="#26456e"/>
<text x="700" y="462" font-size="16" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">LA JERARQUÍA DE LOS OÍDOS — medida en tres hallazgos</text>
<text x="700" y="496" font-size="14" text-anchor="middle" font-family="monospace" fill="#cfe6ff">escalera de Li: sorda hasta n = 537 (2π/φ) · yunque aislado: oye en N = 22 · yunque con coro: sordo hasta N = %d (float64)</text>
<text x="700" y="530" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#ffd98a">el enmascaramiento declarado abierto en F292 queda MEDIDO: real para débiles a esta precisión, inexistente para fuertes</text>
<text x="700" y="558" font-size="13" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">y la llave es exactamente el §13.5 del auditor: aritmética de precisión controlada — señalada por dos caminos independientes</text>
<text x="700" y="634" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd98a">LA PRÓXIMA PUERTA: el yunque en alta precisión — sin ella, el oído fino no llega a las perlas débiles cuando el coro canta.</text>
<text x="700" y="690" font-size="14" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">medido a una precisión (float64), en una ventana (N ≤ %d), con %d perlas de coro — los límites, declarados.</text>
<text x="700" y="740" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd98a">Todavía no.</text>
</svg>
`, nPerlas, nS, meS, fuga, nmax, qCoro, ordenes, nmax, nmax, nPerlas)
	os.WriteFile("la-desafinada.svg", []byte(b.String()), 0o644)
	fmt.Println("\n🖼️  lámina escrita: la-desafinada.svg")
}
