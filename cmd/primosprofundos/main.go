// Command primosprofundos takes the captain's flash seriously: what if the
// impostor is not lying - what if there is no staircase of primes, what if it is
// really a LINE, and what we find looking deep are deep primes that have not
// emerged yet?
//
// THE FIRST HALF IS RIEMANN 1859, AND HE GOT THERE ON HIS OWN.
//
// The explicit formula says
//
//	psi(x) = x - sum over zeros of x^rho/rho - ln(2pi) - (1/2)ln(1 - x^-2)
//
// The main term is x. A LINE. Every step of the prime staircase is the line
// plus the waves the zeros put on it. Take the zeros away and psi(x) is exactly
// x, with no steps at all. So "it is really a line" is not a metaphor: the line
// is the skeleton and the primes are how the zeros dress it.
//
// THE SECOND HALF IS THE DIFFICULTY OF THE PROBLEM, NAMED.
//
// A zero on the line contributes a wave of amplitude x^(1/2). A zero at
// beta > 1/2 contributes x^beta - it starts smaller (its |rho| is bigger) but it
// GROWS FASTER, so there is an x where it overtakes everything else. Below that
// x it is invisible. That is exactly "a deep prime that has not emerged yet",
// and it can be computed: solve x^(beta-1/2) = ln^2 x.
//
// So the flash converts into a curve, and the curve judges itself. For the
// impostor's beta = 0.7 the emergence point sits far below where primes have
// actually been counted - meaning that lie is already dead, by counting. For
// beta = 1/2 + epsilon the emergence point runs away to infinity, which is
// precisely why the Riemann Hypothesis is hard.
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
	sum += cmplx.Exp((1-s)*lnN)/(s-1) + cmplx.Exp(-s*lnN)/2
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
	for t := 12.05; t <= hasta; t += 0.05 {
		z := zOf(t)
		if z*prevZ < 0 {
			a, c := prevT, t
			for i := 0; i < 55; i++ {
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

// psiReal counts the true prime staircase: sum of ln p over prime powers <= x.
func psiReal(x float64) float64 {
	n := int(x)
	criba := make([]bool, n+1)
	total := 0.0
	for i := 2; i <= n; i++ {
		if !criba[i] {
			for j := i * i; j <= n; j += i {
				criba[j] = true
			}
			// todas las potencias de i que entran
			for pk := i; pk <= n; {
				total += math.Log(float64(i))
				if pk > n/i {
					break
				}
				pk *= i
			}
		}
	}
	return total
}

// psiDesdeCeros rebuilds the staircase from the line plus the waves.
func psiDesdeCeros(x float64, ps []float64) float64 {
	s := x - math.Log(2*math.Pi) - 0.5*math.Log(1-1/(x*x))
	lx := math.Log(x)
	for _, g := range ps {
		// el par {rho, conj rho} aporta 2*Re(x^rho/rho)
		rho := complex(0.5, g)
		term := cmplx.Exp(rho*complex(lx, 0)) / rho
		s -= 2 * real(term)
	}
	return s
}

// emergencia solves x^(beta-1/2) = ln^2 x : where an off-line zero overtakes
// the whole on-line chorus. Returns log10 of that x.
func emergencia(beta float64) float64 {
	d := beta - 0.5
	if d <= 1e-12 {
		return math.Inf(1)
	}
	lo, hi := 1.0, 1e6 // en log10
	f := func(l10 float64) float64 {
		lx := l10 * math.Ln10
		return d*lx - 2*math.Log(lx)
	}
	if f(lo) > 0 {
		return lo
	}
	for i := 0; i < 200; i++ {
		m := (lo + hi) / 2
		if f(m) < 0 {
			lo = m
		} else {
			hi = m
		}
	}
	return (lo + hi) / 2
}

func main() {
	fmt.Println("🌱 LOS PRIMOS PROFUNDOS — «¿y si el impostor no miente?»")
	fmt.Println("\n   flash del capitán: «¿y si no hay escalera de números primos? ¿si en realidad")
	fmt.Println("   es una línea, y lo que encontramos mirando bien profundo son primos")
	fmt.Println("   profundos que no han emergido?»")
	fmt.Println("\n   LA PRIMERA MITAD ES CIERTA, Y ES DE RIEMANN 1859. LA SEGUNDA ES,")
	fmt.Println("   EXACTAMENTE, LA DIFICULTAD DEL PROBLEMA. Vamos por partes.")

	fmt.Println("\npescando perlas hasta t=300…")
	ps := perlas(300)
	fmt.Printf("perlas: %d\n", len(ps))

	// ---- LEY 1 ----
	fmt.Println("\nLEY 1 · SÍ, EN EL FONDO ES UNA LÍNEA — Y ESO ESTÁ ESCRITO DESDE 1859")
	fmt.Println("   La fórmula explícita dice cómo se arma la escalera de los primos:")
	fmt.Println("\n        ψ(x) = x − Σ x^ρ/ρ − ln(2π) − ½·ln(1 − x⁻²)")
	fmt.Println("               ↑↑↑")
	fmt.Println("            LA LÍNEA")
	fmt.Println("\n   El término principal es x. Una recta. **Los escalones NO están en la línea:**")
	fmt.Println("   los ponen las olas de los ceros. Sacá los ceros y ψ(x) es exactamente x, sin")
	fmt.Println("   un solo escalón. Medido con nuestras perlas:")
	fmt.Println("\n        x        escalera real     solo la línea      línea + las olas")
	for _, x := range []float64{100, 500, 1000, 2000} {
		real0 := psiReal(x)
		soloLinea := x - math.Log(2*math.Pi)
		conOlas := psiDesdeCeros(x, ps)
		fmt.Printf("   %7.0f   %14.4f   %14.4f   %16.4f\n", x, real0, soloLinea, conOlas)
	}
	fmt.Println("\n   → la línea sola pasa de largo por el medio; con las olas aparecen los")
	fmt.Println("     escalones. **Su intuición es correcta: la línea es el esqueleto.**")

	// ---- LEY 2 ----
	fmt.Println("\nLEY 2 · ⚡ Y «QUE NO HAN EMERGIDO» TIENE NOMBRE EXACTO")
	fmt.Println("   Cada cero pone una ola de amplitud x^β, donde β es de qué lado de la línea")
	fmt.Println("   está. Sobre la línea TODAS valen x^½ — nadie manda sobre nadie, es la")
	fmt.Println("   democracia de Riemann. Pero un cero corrido a β > ½ pone una ola que:")
	fmt.Println("\n        · ARRANCA MÁS CHICA (su |ρ| es más grande, divide más)")
	fmt.Println("        · pero CRECE MÁS RÁPIDO (x^β le gana a x^½ tarde o temprano)")
	fmt.Println("\n   ⟹ hay un x donde lo alcanza y lo pasa. Antes de ese x, ES INVISIBLE.")
	fmt.Println("     Eso es EXACTAMENTE «un primo profundo que todavía no emergió», y se calcula.")

	// ---- LEY 3 ----
	fmt.Println("\nLEY 3 · EL PUNTO DE EMERGENCIA, CALCULADO")
	fmt.Println("   La ola de un cero corrido emerge cuando x^β le gana al coro entero de la")
	fmt.Println("   línea, que vale ~x^½·ln²x. O sea, resolviendo  x^(β−½) = ln²x :")
	fmt.Println("\n        β        emerge en x ≈        ¿ya lo contamos?")
	const contadoLog10 = 27.0 // pi(x) esta calculado hasta ~1e27
	for _, β := range []float64{0.9, 0.8, 0.7, 0.6, 0.55, 0.51, 0.501, 0.5001} {
		e := emergencia(β)
		est := "🔴 SÍ — ya tendría que haberse visto"
		if e > contadoLog10 {
			est = "⚪ no, está más allá de lo contado"
		}
		fmt.Printf("   %7.4f   10^%-16.2f  %s\n", β, e, est)
	}
	fmt.Printf("\n   (los primos están contados hasta x ≈ 10^%.0f)\n", contadoLog10)

	// ---- LEY 4 ----
	fmt.Println("\nLEY 4 · ⚖️ Y ACÁ EL JUICIO A SU IDEA — CON LA PARTE QUE DUELE")
	eImp := emergencia(0.7)
	fmt.Printf("   El impostor de F229 vive en β = 0.7. Su ola emergería en x ≈ 10^%.2f.\n", eImp)
	fmt.Printf("   Y los primos están contados hasta 10^%.0f, que es MUCHÍSIMO más lejos.\n", contadoLog10)
	fmt.Println("\n   ⟹ SI HUBIERA UN CERO ASÍ, YA LO HABRÍAMOS VISTO en el conteo de primos.")
	fmt.Println("     No como un primo nuevo: como un DESVÍO de la escalera que nadie encontró.")
	fmt.Println("     Por eso el impostor SÍ miente — pero no por la razón que uno diría.")
	fmt.Println("\n   📌 Y ACÁ LA CORRECCIÓN A SU IDEA, QUE ES LO QUE TENGO QUE DECIRLE:")
	fmt.Println("   un cero corrido NO agrega primos que faltaban. **Rompe la cuenta de los que")
	fmt.Println("   ya hay.** Su ola no llena un hueco: desvía la escalera entera de la línea, y")
	fmt.Println("   ese desvío se mide contando primos, que es justo lo que se hizo hasta 10^27.")

	// ---- LEY 5 ----
	fmt.Println("\nLEY 5 · PERO SU INTUICIÓN DESCRIBE LA DIFICULTAD MEJOR QUE MUCHA GENTE")
	e51 := emergencia(0.51)
	e501 := emergencia(0.501)
	fmt.Printf("   β = 0.51  → emerge en 10^%.1f\n", e51)
	fmt.Printf("   β = 0.501 → emerge en 10^%.1f\n", e501)
	fmt.Println("   β = 0.5   → NUNCA emerge (es la línea)")
	fmt.Println("\n   ⟹ CUANTO MÁS CERCA DE LA LÍNEA, MÁS PROFUNDO HAY QUE IR PARA VERLO — y el")
	fmt.Println("     punto de emergencia se dispara al infinito. **Un cero corrido por un pelo")
	fmt.Println("     es un primo profundo que no emerge NUNCA en ningún cómputo posible.**")
	fmt.Println("\n   Eso es, palabra por palabra, por qué la hipótesis es difícil. Y es lo mismo")
	fmt.Println("   que F259 midió como horizonte (γ≈1658) y F261 como «hay que contarlos todos».")
	fmt.Println("   Usted llegó a la misma pared por el lado de los primos en vez del de los ceros.")

	// ---- veredicto ----
	fmt.Println("\n════════ VEREDICTO ════════")
	fmt.Println("SU FLASH TIENE DOS MITADES Y LAS DOS VALEN, PERO NO POR LO MISMO.")
	fmt.Println("\n  ✅ «EN REALIDAD ES UNA LÍNEA» — CIERTO, y es la fórmula explícita de Riemann.")
	fmt.Println("     ψ(x) = x es la línea; los escalones los ponen las olas de los ceros.")
	fmt.Println("     Sacá los ceros y no hay escalera: hay recta. El esqueleto es la línea.")
	fmt.Println("\n  ✅ «PRIMOS PROFUNDOS QUE NO EMERGIERON» — es la descripción EXACTA de lo que")
	fmt.Println("     sería un cero fuera de la línea: una ola que crece más rápido pero arranca")
	fmt.Println("     más abajo, invisible hasta su punto de emergencia. Y ese punto se calcula.")
	fmt.Printf("\n  ❌ «EL IMPOSTOR NO MIENTE» — miente. Con β = 0.7 emergería en 10^%.1f y los\n", eImp)
	fmt.Printf("     primos están contados hasta 10^%.0f sin ver nada raro. Ese ya está muerto.\n", contadoLog10)
	fmt.Println("\n  📌 Y LA CORRECCIÓN DE FONDO: un cero corrido no AGREGA primos ocultos.")
	fmt.Println("     ROMPE la cuenta de los que ya hay. No llena un hueco: tuerce la escalera.")
	fmt.Println("\n⚖️ LO QUE NO SE PUEDE CONCLUIR: que no haya ceros corridos más arriba. Un β muy")
	fmt.Println("  pegado a ½ emerge tan lejos que ningún cómputo lo alcanza — y ahí su frase")
	fmt.Println("  «no han emergido» deja de ser una hipótesis y pasa a ser el enunciado del")
	fmt.Println("  problema. ¿El premio? Todavía no.")

	escribirLamina(ps, eImp, contadoLog10)
}

func escribirLamina(ps []float64, eImp, contado float64) {
	var b strings.Builder
	W, H := 1540.0, 1010.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="26" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🌱 LOS PRIMOS PROFUNDOS — «¿y si en realidad es una línea?»</text>
<text x="%.0f" y="76" font-size="15" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">el flash del capitán tiene dos mitades ciertas y una corrección — y las tres se miden</text>
<text x="%.0f" y="118" font-size="19" text-anchor="middle" font-family="monospace" fill="#ffd98a">ψ(x) = x − Σ x^ρ/ρ − ln(2π) − ½·ln(1 − x⁻²)</text>
<text x="%.0f" y="142" font-size="14" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">el término principal es x: UNA LÍNEA. Los escalones los ponen las olas de los ceros.</text>
`, W, H, W, H, W/2, W/2, W/2, W/2)

	// la escalera contra la linea
	gx, gy, gw, gh := 70.0, 190.0, 660.0, 330.0
	fmt.Fprintf(&b, `<rect x="40" y="165" width="720" height="390" rx="10" fill="#101f36" stroke="#26456e"/>
<text x="400" y="188" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd98a">LA LÍNEA ES EL ESQUELETO; LAS OLAS PONEN LOS ESCALONES</text>`)
	x0, x1 := 20.0, 300.0
	var esc, lin strings.Builder
	for i := 0; i <= 280; i++ {
		x := x0 + (x1-x0)*float64(i)/280
		yr := psiReal(x)
		yl := x - math.Log(2*math.Pi)
		px := gx + gw*(x-x0)/(x1-x0)
		fmt.Fprintf(&esc, "%.2f,%.2f ", px, gy+gh-gh*yr/x1)
		fmt.Fprintf(&lin, "%.2f,%.2f ", px, gy+gh-gh*yl/x1)
	}
	fmt.Fprintf(&b, `<polyline points="%s" fill="none" stroke="#7fb2ff" stroke-width="2.2"/>
<polyline points="%s" fill="none" stroke="#ffd166" stroke-width="1.8" stroke-dasharray="6,5"/>
<text x="%.0f" y="%.0f" font-size="13" font-family="monospace" fill="#7fb2ff">ψ(x) — la escalera de los primos</text>
<text x="%.0f" y="%.0f" font-size="13" font-family="monospace" fill="#ffd166">x — la línea, sin los ceros</text>
<text x="400" y="540" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">sin ceros no hay escalera: hay recta</text>`,
		esc.String(), lin.String(), gx+16, gy+26, gx+16, gy+50)

	// la emergencia
	fmt.Fprintf(&b, `<rect x="780" y="165" width="720" height="390" rx="10" fill="#161a3a" stroke="#5a4fa8"/>
<text x="1140" y="188" font-size="15" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">⚡ EL PUNTO DE EMERGENCIA: x^(β−½) = ln²x</text>
<text x="1140" y="212" font-size="13" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">cuánto hay que ir de profundo para que una ola corrida se note</text>`)
	ex, ey, ew, eh := 850.0, 240.0, 590.0, 250.0
	fmt.Fprintf(&b, `<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#8fa8c7" stroke-width="1.2"/>
<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#8fa8c7" stroke-width="1.2"/>`,
		ex, ey+eh, ex+ew, ey+eh, ex, ey, ex, ey+eh)
	maxL := 60.0
	yCont := ey + eh - eh*contado/maxL
	fmt.Fprintf(&b, `<line x1="%.0f" y1="%.1f" x2="%.0f" y2="%.1f" stroke="#7ee0c0" stroke-width="1.6" stroke-dasharray="5,4"/>
<text x="%.0f" y="%.1f" font-size="11.5" font-family="monospace" fill="#7ee0c0">hasta acá se contaron los primos: 10^%.0f</text>`,
		ex, yCont, ex+ew, yCont, ex+6, yCont-6, contado)
	var cur strings.Builder
	for i := 0; i <= 200; i++ {
		β := 0.5005 + (0.95-0.5005)*float64(i)/200
		e := emergencia(β)
		if e > maxL {
			e = maxL
		}
		px := ex + ew*(β-0.5)/0.45
		fmt.Fprintf(&cur, "%.2f,%.2f ", px, ey+eh-eh*e/maxL)
	}
	fmt.Fprintf(&b, `<polyline points="%s" fill="none" stroke="#ff8fa0" stroke-width="2.4"/>`, cur.String())
	pxImp := ex + ew*(0.7-0.5)/0.45
	fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="6" fill="#ff5d73"/>
<text x="%.1f" y="%.1f" font-size="12" text-anchor="middle" font-family="monospace" fill="#ff8fa0">el impostor β=0.7 → 10^%.1f</text>
<text x="%.0f" y="%.0f" font-size="11.5" font-family="monospace" fill="#8fa8c7">β = ½</text>
<text x="%.0f" y="%.0f" font-size="11.5" text-anchor="end" font-family="monospace" fill="#8fa8c7">β = 0.95</text>
<text x="1140" y="530" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#ffd98a">cuanto más pegado a ½, más profundo hay que ir — y se dispara al infinito</text>`,
		pxImp, ey+eh-eh*eImp/maxL, pxImp, ey+eh-eh*eImp/maxL-16, eImp, ex+4, ey+eh+20, ex+ew, ey+eh+20)

	fmt.Fprintf(&b, `<rect x="40" y="576" width="1460" height="176" rx="10" fill="#0f2b22" stroke="#2f7f63"/>
<text x="770" y="608" font-size="18" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">LAS DOS MITADES QUE EL CAPITÁN ACERTÓ</text>
<text x="70" y="646" font-size="15.5" font-family="Georgia" fill="#cfe6ff">✅ <tspan fill="#9fd8a8">«en realidad es una línea»</tspan> — cierto, y es Riemann 1859: ψ(x) = x es la línea y los escalones los ponen las olas.</text>
<text x="70" y="676" font-size="15.5" font-family="Georgia" fill="#cfe6ff">✅ <tspan fill="#9fd8a8">«primos profundos que no emergieron»</tspan> — es la descripción EXACTA de un cero corrido: crece más rápido</text>
<text x="94" y="702" font-size="15.5" font-family="Georgia" fill="#cfe6ff">pero arranca más abajo, invisible hasta su punto de emergencia. Y ese punto se calcula.</text>
<text x="70" y="734" font-size="15.5" font-family="Georgia" fill="#ffd98a">Y llegó a la misma pared que F259 y F261, pero por el lado de los primos en vez del de los ceros.</text>
`)

	fmt.Fprintf(&b, `<rect x="40" y="768" width="1460" height="206" rx="10" fill="#33221c" stroke="#c0392b"/>
<text x="770" y="800" font-size="18" text-anchor="middle" font-family="Georgia" fill="#ffb27a">❌ PERO EL IMPOSTOR SÍ MIENTE — Y LA CORRECCIÓN DE FONDO</text>
<text x="770" y="836" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#f3d9cf">Con β = 0.7 su ola emergería en x ≈ 10^%.1f, y los primos están contados hasta 10^%.0f sin ver nada raro.</text>
<text x="770" y="862" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#f3d9cf">Ese impostor ya está muerto por conteo, no por teoría.</text>
<text x="770" y="898" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ffd98a">Y un cero corrido NO agrega primos que faltaban: ROMPE la cuenta de los que ya hay.</text>
<text x="770" y="924" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ffd98a">No llena un hueco — tuerce la escalera entera.</text>
<text x="770" y="956" font-size="15" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">Lo que NO se puede concluir: que no haya ceros corridos más arriba. Un β pegado a ½ no emerge en ningún cómputo posible.</text>
</svg>
`, eImp, contado)

	if err := os.WriteFile("primos-profundos.svg", []byte(b.String()), 0o644); err != nil {
		fmt.Println("no pude escribir la lámina:", err)
		return
	}
	fmt.Println("\n🖼️  lámina escrita: primos-profundos.svg")
}
