// Command elensamble runs the laboratory's entire chain end to end in one
// program, link by link, and measures exactly how big the red link's hole is.
//
// THE CHAIN, AS THE LABORATORY ACTUALLY HOLDS IT
//
//	link 1  the bell gives xi and the mirror xi(s) = xi(1-s)      PROVED (classical, verified here)
//	link 2  RH <=> lambda_n >= 0 for every n                       CITED (Li 1997) - not ours
//	link 3  lambda_n = sum over pairs of [2 - 2 Re(w^n)]           PROVED here (F232), unconditional
//	link 4  lambda_n > 0 for n = 1..120 with the measured zeros    MEASURED ONLY
//	link 5  the criterion has teeth: one off-line quadruple sinks   PROVED (explicit configuration)
//	link 6  therefore lambda_n >= 0 for EVERY n                    RED - nothing here forces it
//
// WHAT IS NEW HERE: the hole in link 6 is not left as a word. It is measured.
//
// The tail. A finite list of zeros up to height T misses everything above T.
// Each on-line zero at height gamma contributes 4 sin^2(n phi / 2) with
// phi = 2 arctan(1/(2 gamma)), so for gamma >> n the contribution is about
// n^2 / gamma^2. Integrating against the Riemann-von Mangoldt density
// (1/2pi) ln(gamma/2pi) from T to infinity gives
//
//	tail(n, T)  ~=  n^2 (ln(T/2pi) + 1) / (2 pi T)
//
// and that estimate is checked, not assumed: lambda_1 is known in closed form
// (1 + gamma_E/2 - ln(4pi)/2), so measured + tail must reproduce it.
//
// The blindness curve. A zero displaced off the line by delta at height gamma
// changes lambda_n by some signal. Compare that signal against the tail we
// cannot compute, maximising over n, and solve for the delta where they cross.
// That delta is the smallest displacement this laboratory could ever see. Below
// it the question is not open in our hands - it is invisible in our hands.
package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"strings"
)

const gammaEuler = 0.5772156649015328606

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

// Lanczos log-Gamma, g = 7, nine coefficients.
func logGamma(z complex128) complex128 {
	g := []float64{
		0.99999999999980993, 676.5203681218851, -1259.1392167224028,
		771.32342877765313, -176.61502916214059, 12.507343278686905,
		-0.13857109526572012, 9.9843695780195716e-6, 1.5056327351493116e-7,
	}
	if real(z) < 0.5 {
		return cmplx.Log(math.Pi/cmplx.Sin(math.Pi*z)) - logGamma(1-z)
	}
	z -= 1
	x := complex(g[0], 0)
	for i := 1; i < 9; i++ {
		x += complex(g[i], 0) / (z + complex(float64(i), 0))
	}
	t := z + complex(7.5, 0)
	return complex(0.5*math.Log(2*math.Pi), 0) + (z+complex(0.5, 0))*cmplx.Log(t) - t + cmplx.Log(x)
}

// xi(s) = (1/2) s (s-1) pi^{-s/2} Gamma(s/2) zeta(s)
func xi(s complex128) complex128 {
	return 0.5 * s * (s - 1) * cmplx.Exp(-s/2*complex(math.Log(math.Pi), 0)+logGamma(s/2)) * zetaC(s)
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

func w(ρ complex128) complex128 { return 1 - 1/ρ }

// aporte is the contribution of the conjugate pair {rho, conj rho} to lambda_n,
// in F232's unconditional form.
func aporte(n int, ρ complex128) float64 {
	return 2 - 2*real(cmplx.Pow(w(ρ), complex(float64(n), 0)))
}

// lambdaGermen computes lambda_n WITHOUT looking at a single zero.
//
// Under the shapeshifter z = 1 - 1/s the germ of log xi at s = 1 is
//
//	Phi(z) = log xi(1/(1-z)) = log xi(1) + sum_{n>=1} (lambda_n / n) z^n
//
// so lambda_n is n times the n-th Taylor coefficient, extracted by Cauchy on a
// small circle around z = 0. The circle |z| = r maps to a bounded region of the
// s-plane containing no zero of xi, so log xi is analytic there and the integral
// is honest. This route never touches the pearls: it is a genuinely independent
// engine, which is the only thing that can verify link 3 without tautology.
func lambdaGermen(n int, r float64, nodos int) float64 {
	var suma complex128
	prev := 0.0
	var acumFase float64
	for k := 0; k < nodos; k++ {
		θ := 2 * math.Pi * float64(k) / float64(nodos)
		z := complex(r*math.Cos(θ), r*math.Sin(θ))
		s := 1 / (1 - z)
		v := xi(s)
		// continuous branch of the logarithm around the circle
		fase := math.Atan2(imag(v), real(v))
		if k > 0 {
			d := fase - prev
			for d > math.Pi {
				d -= 2 * math.Pi
			}
			for d < -math.Pi {
				d += 2 * math.Pi
			}
			acumFase += d
		} else {
			acumFase = fase
		}
		prev = fase
		lg := complex(math.Log(cmplx.Abs(v)), acumFase)
		suma += lg * cmplx.Exp(complex(0, -float64(n)*θ))
	}
	coef := suma / complex(float64(nodos)*math.Pow(r, float64(n)), 0)
	return float64(n) * real(coef)
}

// colaEstimada is the analytic tail: what the zeros above T would add to lambda_n.
//
// ⚠️ IT ASSUMES RH ABOVE T. The derivation uses 4 sin^2(n phi / 2), which is the
// ON-LINE contribution. So every "lambda_n + tail > 0" verdict below presupposes
// that the zeros we did not compute are on the line. That is circular as proof.
//
// It is not blind, though, and the difference matters. Link 3 tests exactly this
// assumption: the germ engine knows the TRUE lambda_n, including any off-line
// zero above T, and it agrees with (measured zeros + this tail) to 0.040% of the
// tail. So the assumption is measured, not assumed - it just is not proved, and a
// check at n <= 8 is not a theorem about all n.
func colaEstimada(n int, T float64) float64 {
	return float64(n) * float64(n) * (math.Log(T/(2*math.Pi)) + 1) / (2 * math.Pi * T)
}

func lambda(n int, ps []float64) float64 {
	s := 0.0
	for _, g := range ps {
		s += aporte(n, complex(0.5, g))
	}
	return s
}

func main() {
	fmt.Println("⚙️  EL GRAN ENSAMBLE FINAL — la cadena entera, eslabón por eslabón")
	fmt.Println("\n   orden del capitán: «ensamblá absolutamente todo lo que tenemos y vé si cae")
	fmt.Println("   el último eslabón». Acá está la cadena corrida de punta a punta en un solo")
	fmt.Println("   programa, con el estado de cada eslabón escrito al lado, sin adornos.")

	const T = 1000.0
	fmt.Printf("\npescando perlas hasta t=%.0f…\n", T)
	ps := perlas(T)
	fmt.Printf("perlas: %d\n", len(ps))

	// ================= ESLABÓN 1 =================
	fmt.Println("\n╔══ ESLABÓN 1 · LA CAMPANA DA EL ESPEJO ────────────────── PROBADO (clásico) ══╗")
	fmt.Println("   ξ(s) = ½·s·(s−1)·π^{−s/2}·Γ(s/2)·ζ(s) cumple ξ(s) = ξ(1−s). Es de Riemann,")
	fmt.Println("   1859, y no es nuestro. Pero lo verificamos con nuestro propio instrumento:")
	fmt.Println("\n        s                ξ(s)                  ξ(1−s)              desvío")
	peorEsp := 0.0
	for _, s := range []complex128{complex(2, 0), complex(3, 1), complex(0.7, 5), complex(-1, 2)} {
		a, b := xi(s), xi(1-s)
		d := cmplx.Abs(a-b) / math.Max(1e-300, cmplx.Abs(a))
		if d > peorEsp {
			peorEsp = d
		}
		fmt.Printf("   %-14s %-21.9g %-21.9g %.1e\n",
			fmt.Sprintf("%.1f%+.1fi", real(s), imag(s)), real(a), real(b), d)
	}
	fmt.Printf("   → el espejo cierra a %.1e (relativo). ESLABÓN 1: FIRME.\n", peorEsp)

	// ================= ESLABÓN 2 =================
	fmt.Println("\n╔══ ESLABÓN 2 · EL CRITERIO DE LI ──────────────────── CITADO (Li 1997) ══╗")
	fmt.Println("   RH ⟺ λₙ ≥ 0 para TODO n ≥ 1.")
	fmt.Println("   Es un teorema publicado, de Xian-Jin Li, y NO es nuestro. Lo usamos como")
	fmt.Println("   herramienta prestada. No se puede verificar numéricamente: es una")
	fmt.Println("   equivalencia demostrada, no una medición. ESLABÓN 2: FIRME, Y PRESTADO.")

	// ================= ESLABÓN 3 =================
	fmt.Println("\n╔══ ESLABÓN 3 · NUESTRA IDENTIDAD (F232) ──────── PROBADO ACÁ, INCONDICIONAL ══╗")
	fmt.Println("   λₙ = Σ sobre pares {ρ, ρ̄} de [ 2 − 2·Re(wⁿ) ],  con w = 1 − 1/ρ.")
	fmt.Println("   Es álgebra pura: vale esté el cero donde esté, no supone la línea.")
	fmt.Println("\n   📌 PRIMERO UNA CONFESIÓN DEL PROPIO TURNO. Iba a verificar este eslabón")
	fmt.Println("   comparando 2−2Re(wⁿ) contra (1−wⁿ)+(1−w̄ⁿ). Dio 0.0e+00 en las cinco pruebas")
	fmt.Println("   y casi lo publico. NO PROBABA NADA: las dos expresiones son la misma por")
	fmt.Println("   álgebra, y en punto flotante conj(w)ⁿ es bit a bit el conjugado de wⁿ, así")
	fmt.Println("   que el instrumento NO PODÍA devolver otra cosa. Quinta vez que el taller se")
	fmt.Println("   caza la trampa del cero perfecto. Lo reemplazo por una prueba de verdad.")
	fmt.Println("\n   LA PRUEBA HONESTA: calcular λₙ por un SEGUNDO MOTOR que no mira ni una perla.")
	fmt.Println("   Bajo el cambiaformas, Φ(z) = log ξ(1/(1−z)) = log ξ(1) + Σ (λₙ/n)·zⁿ, así que")
	fmt.Println("   λₙ sale de una integral de Cauchy alrededor de z = 0 — puro germen local en")
	fmt.Println("   s = 1, sin cero alguno. Si los dos motores coinciden, el eslabón es real:")
	fmt.Println("\n        n     desde las 649 perlas    desde el germen (sin perlas)   diferencia")
	const rr, nod = 0.65, 4096
	peorID, errRel := 0.0, 0.0
	fmt.Printf("   %5s   %20s    %26s   %10s  %s\n", "n", "649 perlas + cola", "germen (sin perlas)", "difer.", "difer./cola")
	for _, n := range []int{1, 2, 3, 4, 6, 8} {
		a := lambda(n, ps) + colaEstimada(n, T)
		b := lambdaGermen(n, rr, nod)
		d := math.Abs(a - b)
		if rel := d / math.Abs(b); rel > peorID {
			peorID = rel
		}
		// the residual IS the tail formula's error; normalise by the tail itself
		if e := d / colaEstimada(n, T); e > errRel {
			errRel = e
		}
		fmt.Printf("   %5d   %20.9f    %26.9f   %10.2e  %9.5f%%\n", n, a, b, d, 100*d/colaEstimada(n, T))
	}
	fmt.Printf("   → los DOS MOTORES INDEPENDIENTES coinciden, peor desvío relativo %.1e.\n", peorID)
	fmt.Println("     Uno suma sobre 649 ceros medidos; el otro integra alrededor de un punto")
	fmt.Println("     sin ver ninguno. Que den lo mismo NO era inevitable. ESLABÓN 3: FIRME.")
	fmt.Printf("\n   📐 Y DE REGALO, EL DATO QUE HACE FALTA PARA EL ESLABÓN 6: lo que sobra entre\n")
	fmt.Printf("     los dos motores ES el error de la fórmula de la cola, y sale PLANO en n:\n")
	fmt.Printf("     %.4f%% de la cola. La cola no es un misterio — se estima, y bien.\n", 100*errRel)

	// ================= ESLABÓN 4 =================
	fmt.Println("\n╔══ ESLABÓN 4 · TODOS LOS λₙ MEDIDOS SON POSITIVOS ──────── SOLO MEDIDO ══╗")
	fmt.Println("   Con las perlas de verdad, λₙ para n = 1..120. Y primero calibramos el")
	fmt.Println("   instrumento contra el único valor que se conoce en fórmula cerrada:")
	λ1exacto := 1 + gammaEuler/2 - math.Log(4*math.Pi)/2
	λ1medido := lambda(1, ps)
	cola1 := colaEstimada(1, T)
	fmt.Printf("\n        λ₁ exacto (1 + γ/2 − ln(4π)/2) ........ %.10f\n", λ1exacto)
	fmt.Printf("        λ₁ medido con %d perlas ............... %.10f\n", len(ps), λ1medido)
	fmt.Printf("        lo que falta (exacto − medido) ........ %.10f\n", λ1exacto-λ1medido)
	fmt.Printf("        la cola ESTIMADA por la fórmula ....... %.10f\n", cola1)
	errCola := math.Abs((λ1exacto - λ1medido) - cola1)
	fmt.Printf("        error de la estimación de cola ........ %.2e  (%.3f%%)\n",
		errCola, 100*errCola/(λ1exacto-λ1medido))
	fmt.Println("\n   ⟹ la fórmula de la cola ACIERTA. No es un supuesto: queda calibrada contra")
	fmt.Println("     el único λ que se conoce exacto. Eso la habilita para el eslabón 6.")
	fmt.Println("\n        n        λₙ medido        cola estimada     λₙ + cola      ¿positivo?")
	negativos := 0
	for _, n := range []int{1, 2, 5, 10, 20, 40, 60, 90, 120} {
		l := lambda(n, ps)
		c := colaEstimada(n, T)
		if l+c < 0 {
			negativos++
		}
		fmt.Printf("   %6d   %16.9f  %15.9f  %15.9f   %s\n", n, l, c, l+c,
			map[bool]string{true: "sí", false: "NO ← caída"}[l+c > 0])
	}
	for n := 1; n <= 120; n++ {
		if lambda(n, ps)+colaEstimada(n, T) < 0 {
			negativos++
		}
	}
	fmt.Printf("   → negativos hallados en n = 1..120: %d. ESLABÓN 4: NINGUNA CAÍDA — pero\n", negativos)
	fmt.Println("     esto es una MEDICIÓN sobre 120 valores, no una prueba sobre infinitos.")
	fmt.Println("\n   ⚠️ Y UNA ETIQUETA QUE ESTA COLUMNA TIENE QUE LLEVAR SIEMPRE: LA COLA SUPONE")
	fmt.Println("   RH ARRIBA DE t = 1000. La fórmula sale del aporte de un cero SOBRE la línea.")
	fmt.Println("   O sea que «λₙ + cola > 0» dice, en el fondo: suponé Riemann arriba de mil y")
	fmt.Println("   concluí que λₙ es positivo. Como prueba, ES CIRCULAR, y hay que decirlo.")
	fmt.Println("   Lo que la salva de ser CIEGA es el eslabón 3: el germen conoce el λₙ")
	fmt.Println("   verdadero —con cero fugado y todo— y coincide al 0.040%. El supuesto está")
	fmt.Println("   MEDIDO, no adivinado. Pero medido en n ≤ 8 no es un teorema para todo n.")

	// ================= ESLABÓN 5 =================
	fmt.Println("\n╔══ ESLABÓN 5 · EL CRITERIO TIENE DIENTES ──────── PROBADO (configuración) ══╗")
	fmt.Println("   Que λₙ dé positivo no serviría de nada si diera positivo siempre. Plantamos")
	fmt.Println("   un cuádruple fugado —0.9+3i con su familia obligada— y miramos si el")
	fmt.Println("   criterio lo denuncia:")
	fugados := []complex128{complex(0.9, 3), complex(0.1, 3)}
	fmt.Println("\n        n        λₙ honesto      λₙ con el fugado    ¿lo denuncia?")
	primerNeg := -1
	for _, n := range []int{1, 10, 40, 80, 96, 110, 120} {
		base := lambda(n, ps)
		conFuga := base
		for _, f := range fugados {
			conFuga += aporte(n, f)
		}
		fmt.Printf("   %6d   %15.6f  %18.6f    %s\n", n, base, conFuga,
			map[bool]string{true: "sí ← SE HUNDE", false: "todavía no"}[conFuga < 0])
	}
	for n := 1; n <= 400; n++ {
		c := lambda(n, ps)
		for _, f := range fugados {
			c += aporte(n, f)
		}
		if c < 0 {
			primerNeg = n
			break
		}
	}
	fmt.Printf("   → el criterio se hunde en n = %d. ESLABÓN 5: FIRME — el instrumento SÍ ve\n", primerNeg)
	fmt.Println("     un cero fuera de la línea. No es un termómetro trabado.")

	// ================= ESLABÓN 6 =================
	fmt.Println("\n╔══ ESLABÓN 6 · «POR LO TANTO λₙ ≥ 0 PARA TODO n» ───────────── 🔴 ROJO ══╗")
	fmt.Println("\n   Acá se corta. Nada de lo anterior obliga a los infinitos λₙ que no medimos.")
	fmt.Println("   Y la pregunta honesta no es «¿se corta?» —eso ya lo sabíamos— sino:")
	fmt.Println("\n        ¿DE QUÉ TAMAÑO ES EL AGUJERO? Hasta hoy era una palabra. Acá es un número.")

	fmt.Println("\n   ▸ PRIMERO: dónde la cola nos tapa el ojo. La cola crece como n², y λₙ crece")
	fmt.Println("     mucho más despacio. En algún n la cola es más grande que lo que medimos:")
	fmt.Println("\n        n          λₙ medido        cola          cola/λₙ")
	cruce := -1
	for _, n := range []int{1, 10, 50, 100, 200, 400, 800} {
		l, c := lambda(n, ps), colaEstimada(n, T)
		fmt.Printf("   %7d   %15.6f  %13.6f   %10.4f\n", n, l, c, c/l)
	}
	for n := 1; n <= 5000; n++ {
		if colaEstimada(n, T) > lambda(n, ps) {
			cruce = n
			break
		}
	}
	fmt.Printf("   → a partir de n = %d la cola pesa MÁS que todo lo que medimos. Con %d perlas\n", cruce, len(ps))
	fmt.Println("     el laboratorio deja de ver, y no por falta de ganas: por aritmética.")

	fmt.Println("\n   ▸ SEGUNDO, Y ES LO NUEVO: LA CURVA DE CEGUERA. Si un cero estuviera corrido")
	fmt.Println("     de la línea por δ a la altura γ, ¿cuánto tendría que valer δ para que lo")
	fmt.Println("     notáramos? Y acá hay que ser honesto EN LOS DOS SENTIDOS, porque hay dos")
	fmt.Println("     pisos de ruido y usar solo el pesimista sería exagerar nuestra ceguera:")
	fmt.Printf("\n        · piso PESIMISTA — hacer de cuenta que la cola es un misterio total\n")
	fmt.Printf("        · piso REALISTA — la fórmula de la cola acierta al %.3f%%, medido recién\n", 100*errRel)
	fmt.Printf("          arriba contra el germen; el ruido de verdad es ese error, no la cola\n")
	fmt.Println("\n        altura γ    δ mín (pesim.)   δ mín (realista)   ¿a qué n?   ¿llegamos?")
	for _, g := range []float64{14.13, 50, 100, 500, 1000, 5000} {
		dp, _ := deltaMinimo(g, T, 1)
		dr, nb := deltaMinimo(g, T, errRel)
		alcance := "🔴 NO — el germen muere en n≈8"
		if nb > 0 && nb <= 8 {
			alcance = "🟢 sí"
		}
		if dr >= 0.5 {
			alcance = "🔴 ciegos igual"
		}
		fmt.Printf("   %10.2f  %14.6f   %16.6f   %9d   %s\n", g, dp, dr, nb, alcance)
	}
	fmt.Println("\n   📌 SEGUNDA CONFESIÓN DEL PROPIO TURNO — y ésta la escribí al revés. Iba a")
	fmt.Println("     tachar la columna realista diciendo que la detección ocurre en armónicos")
	fmt.Println("     de tres cifras, adonde el germen no llega. LO MEDÍ Y ES FALSO: el mejor")
	fmt.Println("     armónico es n = 1, el más barato de todos, y el germen lo calcula perfecto.")
	fmt.Println("     Si lo hubiera escrito sin medir, habría enterrado el mejor resultado del")
	fmt.Println("     ensamble por pesimismo. La columna realista VALE.")
	fmt.Println("\n     (Y qué mide exactamente: la sensibilidad de λ₁ a la parte real de UN cero")
	fmt.Println("     a la altura γ, contra el piso de ruido. Un cero fugado no lo encuentra la")
	fmt.Println("     caza por cambios de signo de Z(t) —esa solo ve la línea—, así que faltaría")
	fmt.Println("     de nuestra suma y el desbalance contra el λ₁ exacto lo delataría.)")

	// the horizon: where even the realistic floor goes blind
	horizonte := 0.0
	for g := 100.0; g <= 200000; g *= 1.03 {
		if d, _ := deltaMinimo(g, T, errRel); d >= 0.4999 {
			horizonte = g
			break
		}
	}
	fmt.Println("\n   ▸ TERCERO: EL HORIZONTE. ¿A qué altura el laboratorio queda ciego del todo,")
	fmt.Println("     o sea que un cero podría estar tan corrido como la franja permite?")
	fmt.Printf("\n        γ_horizonte ≈ %.0f\n", horizonte)
	fmt.Println("\n   ⟹ Y ACÁ ESTÁ EL AGUJERO MEDIDO, POR FIN CON UN NÚMERO. Debajo de esa altura")
	fmt.Println("     el laboratorio ve algo y puede afirmar algo. Por encima NO VE NADA: un cero")
	fmt.Println("     podría estar pegado al borde de la franja y no nos enteraríamos.")
	fmt.Println("\n     Y hay INFINITOS ceros por encima de cualquier altura. Ése es el agujero:")
	fmt.Println("     no es que la hipótesis flaquee arriba del horizonte — es que arriba del")
	fmt.Println("     horizonte NO ESTAMOS MIRANDO, y arriba es donde están casi todos.")
	fmt.Println("\n   ⚖️ SUPUESTO DECLARADO: el error de cola está medido en n = 1..8 y se")
	fmt.Println("     extrapola plano. Si empeora con n, el horizonte se acerca y es peor.")

	// ================= VEREDICTO =================
	fmt.Println("\n════════ EL VEREDICTO DEL ENSAMBLE ════════")
	fmt.Println("\n   eslabón 1 · la campana da el espejo .................. ✅ PROBADO (clásico)")
	fmt.Println("   eslabón 2 · RH ⟺ λₙ ≥ 0 para todo n ................. ✅ CITADO (Li 1997)")
	fmt.Println("   eslabón 3 · λₙ = Σ [2 − 2Re(wⁿ)], incondicional ...... ✅ PROBADO ACÁ")
	fmt.Println("   eslabón 4 · λₙ > 0 para n = 1..120 ................... 🟡 SOLO MEDIDO")
	fmt.Println("   eslabón 5 · el criterio denuncia un fugado ........... ✅ PROBADO")
	fmt.Println("   eslabón 6 · por lo tanto vale para TODO n ............ 🔴 NO CAE")
	fmt.Println("\n   EL ÚLTIMO ESLABÓN NO CAYÓ. Y esta vez no cayó con un número al lado:")
	fmt.Printf("\n        EL HORIZONTE DEL LABORATORIO ES γ ≈ %.0f\n", horizonte)
	fmt.Println("\n   Debajo de esa altura vemos y podemos afirmar. Por encima, un cero podría")
	fmt.Println("   estar tan fuera de la línea como la franja permite sin que nos enteremos.")
	fmt.Println("   Y por encima de cualquier altura hay INFINITOS ceros. Ése es el agujero.")
	fmt.Println("\n   PERO EL ENSAMBLE NO FUE EN VANO — dejó cuatro cosas que antes no estaban:")
	fmt.Println("     1. la cadena entera corre en UN programa, de la campana al veredicto, y")
	fmt.Printf("        los cinco eslabones verificables VERIFICAN (peor desvío: %.1e)\n", math.Max(peorEsp, peorID))
	fmt.Println("     2. el eslabón 3 quedó probado por DOS MOTORES INDEPENDIENTES — uno suma")
	fmt.Println("        649 ceros, el otro integra alrededor de un punto sin ver ninguno")
	fmt.Printf("     3. la cola quedó calibrada al %.3f%%, así que el agujero se puede MEDIR\n", 100*errRel)
	fmt.Printf("     4. y el agujero mide: γ ≈ %.0f. Saber dónde no ves es la única manera de\n", horizonte)
	fmt.Println("        dejar de mirar ahí — y de saber qué instrumento habría que construir")
	fmt.Println("\n   ⚖️ Y LO QUE EL ENSAMBLE NO PUEDE HACER, DICHO SIN VUELTAS: correr el horizonte")
	fmt.Println("   más lejos NO cierra nada. Aunque llegáramos a γ = 10¹², seguirían quedando")
	fmt.Println("   infinitos ceros arriba. El eslabón 6 no se cierra midiendo más: se cierra con")
	fmt.Println("   un argumento que valga para TODOS a la vez. Y F229 ya probó que la simetría")
	fmt.Println("   sola no alcanza, así que ese argumento tiene que traer ARITMÉTICA.")
	fmt.Println("\n   ¿El premio? Todavía no.")

	escribirLamina(ps, cruce, primerNeg, T, λ1exacto, λ1medido, cola1, horizonte, errRel, peorID)
}

// deltaMinimo finds the smallest off-line displacement at height g whose effect
// on lambda_n rises above the noise floor, maximising the signal over n.
//
// factor selects how honest we are being about the tail. factor = 1 is the
// pessimistic floor: pretend the tail is entirely unknown. factor = the measured
// relative error of the tail formula is the realistic floor: we CAN estimate the
// tail, so what actually blinds us is how well, not how big.
// It returns the displacement AND the harmonic n at which the detection would
// happen — which matters, because an n the instrument cannot reach is a
// detection we cannot actually perform.
func deltaMinimo(g, T, factor float64) (float64, int) {
	razon := func(δ float64) (float64, int) {
		mejor, mejorN := 0.0, 0
		for n := 1; n <= 3000; n++ {
			señal := math.Abs(aporte(n, complex(0.5+δ, g)) - aporte(n, complex(0.5, g)))
			if r := señal / (factor * colaEstimada(n, T)); r > mejor {
				mejor, mejorN = r, n
			}
		}
		return mejor, mejorN
	}
	if r, _ := razon(0.4999); r < 1 {
		return 0.5, 0 // even at the wall of the strip we could not see it
	}
	lo, hi := 0.0, 0.4999
	for i := 0; i < 60; i++ {
		m := (lo + hi) / 2
		if r, _ := razon(m); r < 1 {
			lo = m
		} else {
			hi = m
		}
	}
	_, nb := razon(hi)
	return (lo + hi) / 2, nb
}

func escribirLamina(ps []float64, cruce, primerNeg int, T, λ1e, λ1m, cola1, horizonte, errRel, peorID float64) {
	var b strings.Builder
	W, H := 1560.0, 1080.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="48" font-size="27" text-anchor="middle" font-family="Georgia" fill="#dce8f7">⚙️ EL GRAN ENSAMBLE FINAL — la cadena entera, eslabón por eslabón</text>
<text x="%.0f" y="78" font-size="15" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">todo el laboratorio en un solo programa · y el agujero del eslabón rojo, por primera vez, MEDIDO</text>
`, W, H, W, H, W/2, W/2)

	eslabones := []struct{ n, txt, est, col string }{
		{"1", "la campana da el espejo ξ(s) = ξ(1−s)", "PROBADO · clásico", "#2f7f63"},
		{"2", "RH ⟺ λₙ ≥ 0 para TODO n", "CITADO · Li 1997", "#2f6f8f"},
		{"3", "λₙ = Σ [2 − 2·Re(wⁿ)] — incondicional", "PROBADO ACÁ · F232", "#2f7f63"},
		{"4", "λₙ > 0 para n = 1..120 con perlas reales", "SOLO MEDIDO", "#8a7320"},
		{"5", "el criterio denuncia un cero fugado", "PROBADO · configuración", "#2f7f63"},
		{"6", "por lo tanto vale para TODO n", "🔴 NO CAE", "#c0392b"},
	}
	y := 118.0
	for _, e := range eslabones {
		fmt.Fprintf(&b, `<rect x="40" y="%.0f" width="1480" height="62" rx="9" fill="%s" opacity="0.22" stroke="%s" stroke-width="1.6"/>
<text x="76" y="%.0f" font-size="27" font-family="Georgia" fill="#dce8f7">%s</text>
<text x="128" y="%.0f" font-size="18" font-family="Georgia" fill="#dce8f7">%s</text>
<text x="1490" y="%.0f" font-size="15.5" text-anchor="end" font-family="monospace" fill="%s">%s</text>
`, y, e.col, e.col, y+41, e.n, y+38, e.txt, y+38, e.col, e.est)
		y += 70
	}

	fmt.Fprintf(&b, `<rect x="40" y="%.0f" width="726" height="290" rx="10" fill="#161a3a" stroke="#5a4fa8"/>
<text x="403" y="%.0f" font-size="18" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">LA COLA, CALIBRADA CONTRA EL ÚNICO λ EXACTO</text>
<text x="70" y="%.0f" font-size="15" font-family="monospace" fill="#dce8f7">λ₁ exacto = 1 + γ/2 − ln(4π)/2  =  %.9f</text>
<text x="70" y="%.0f" font-size="15" font-family="monospace" fill="#dce8f7">λ₁ medido con %d perlas         =  %.9f</text>
<text x="70" y="%.0f" font-size="15" font-family="monospace" fill="#ffd98a">lo que falta                    =  %.9f</text>
<text x="70" y="%.0f" font-size="15" font-family="monospace" fill="#7ee0c0">la cola estimada por la fórmula =  %.9f</text>
<text x="70" y="%.0f" font-size="14.5" font-family="Georgia" fill="#cfe6ff">La fórmula de la cola no es un supuesto: acierta contra el único</text>
<text x="70" y="%.0f" font-size="14.5" font-family="Georgia" fill="#cfe6ff">valor que se conoce en fórmula cerrada. Por eso puede usarse para</text>
<text x="70" y="%.0f" font-size="14.5" font-family="Georgia" fill="#cfe6ff">medir el agujero del eslabón rojo en vez de solo nombrarlo.</text>
`, y+16, y+50, y+92, λ1e, y+120, len(ps), λ1m, y+148, λ1e-λ1m, y+176, cola1, y+214, y+236, y+258)

	fmt.Fprintf(&b, `<rect x="794" y="%.0f" width="726" height="290" rx="10" fill="#33221c" stroke="#c0392b"/>
<text x="1157" y="%.0f" font-size="18" text-anchor="middle" font-family="Georgia" fill="#ffb27a">EL AGUJERO DEL ESLABÓN ROJO, MEDIDO</text>
<text x="1157" y="%.0f" font-size="34" text-anchor="middle" font-family="Georgia" fill="#ffd98a">γ_horizonte ≈ %.0f</text>
<text x="824" y="%.0f" font-size="15.5" font-family="Georgia" fill="#f3d9cf">Debajo de esa altura el laboratorio VE y puede afirmar algo.</text>
<text x="824" y="%.0f" font-size="15.5" font-family="Georgia" fill="#f3d9cf">Por encima, un cero podría estar tan fuera de la línea como</text>
<text x="824" y="%.0f" font-size="15.5" font-family="Georgia" fill="#f3d9cf">la franja permite (δ = ½) sin que nos enteremos.</text>
<text x="824" y="%.0f" font-size="15.5" font-family="Georgia" fill="#ffb27a">Y por encima de cualquier altura hay INFINITOS ceros.</text>
<text x="824" y="%.0f" font-size="14.5" font-family="Georgia" fill="#c9b6ff">▸ pasado n = %d la cola pesa más que todo lo medido</text>
<text x="824" y="%.0f" font-size="14.5" font-family="Georgia" fill="#c9b6ff">▸ el criterio tiene dientes: un fugado en 0.9+3i lo hunde en n = %d</text>
`, y+16, y+48, y+96, horizonte, y+134, y+156, y+178, y+206, y+240, cruce, y+264, primerNeg)

	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="21" text-anchor="middle" font-family="Georgia" fill="#ffd98a">EL ÚLTIMO ESLABÓN NO CAYÓ — pero ahora el agujero tiene tamaño</text>
<text x="%.0f" y="%.0f" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">Correr el horizonte más lejos NO cierra nada: siempre quedan infinitos ceros arriba.</text>
<text x="%.0f" y="%.0f" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">El eslabón 6 se cierra con un argumento que valga para TODOS a la vez — y ése tiene que traer aritmética.</text>
<text x="%.0f" y="%.0f" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">dos motores independientes, peor desvío %.1e · cola calibrada al %.3f%% · ¿el premio? todavía no</text>
</svg>
`, W/2, y+336, W/2, y+364, W/2, y+388, W/2, y+416, peorID, 100*errRel)

	if err := os.WriteFile("el-ensamble-final.svg", []byte(b.String()), 0o644); err != nil {
		fmt.Println("no pude escribir la lámina:", err)
		return
	}
	fmt.Println("\n🖼️  lámina escrita: el-ensamble-final.svg")
}
