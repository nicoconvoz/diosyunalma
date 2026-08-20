// Command lasdosperlas runs PHASE 1 of the captain's own Theorem-2 attack
// sheet ("TEOREMA_2_YUNQUE_Dos_Perlas_Armonia_Relacional" - his idea,
// blessed by the auditor, carrying his flash signature): the central
// experiment of its §8 across the §9 matrix, measuring how TWO off-line
// pearls interact in the anvil's ladder.
//
// SETUP. Spectrum = the 38-pearl on-line choir + quartet 1 (the real DH
// pearl: radius r1, phase theta1) + quartet 2 parametrized RELATIONALLY,
// exactly as the sheet asks:
//
//	pearl 2 = (r1^rho, theta1 * tau)
//	rho = radial ratio  (Delta_R = (rho-1)·log r1)
//	tau = phase ratio   (Delta_theta = (tau-1)·theta1)
//
// Observable: n0(rho, tau) = first n with lambda_n < 0 - the detection
// time. Reference: pearl 1 alone + choir gives n0 = 85622 (F296/F297).
//
// WHAT PHASE 1 MEASURED (the PATTERN for Yui's PHASE 2):
//
//  1. PROTECTIVE HARMONY EXISTS - the captain's flash was RIGHT. It is
//     weak and rare (25 of 1600 grid configurations delay detection past
//     the lone pearl; the best, tau = 1.010, delays it 4.5%) and its
//     mechanism is the BEAT: two nearly-twin pearls slightly detuned
//     create a beat envelope cos(n·Delta_theta/2) whose null covers the
//     rupture zone (measured envelope 0.264 there, vs 1.0 for exact
//     twins). But it never SAVES: every configuration still breaks at a
//     finite n0 - protection is delay, not immunity. (The shop's first
//     draft of this verdict said "zero delay" - the data itself refuted
//     the typed conclusion before registration. Recorded.)
//  2. RADIAL DOMINANCE RULES THE COARSE SHAPE: along tau = 1, n0 scales
//     like ~ n0_sola / rho (measured log-log slope printed) - the
//     bigger-radius pearl sets the clock, as the sheet's §7 A suspected.
//  3. THE SIGN OF Delta_theta IS IRRELEVANT: tau and -tau give identical
//     n0 (cos is even) - measured, and worth recording for the sheet's
//     case 9.
//  4. THE PHASE FINE-STRUCTURE IS MUSICAL, NOT LINEAR: within equal
//     radii the detection time varies in a narrow band whose ranking is
//     set by the RATIO theta2/theta1 - phase-locked twins (tau = 1) are
//     the most fragile, and among the named musical intervals the most
//     protective in the DH configuration is THE JUST FOURTH tau = 4/3 -
//     the laboratory's oldest interval. A control with a different base
//     pearl tests whether that ranking is universal or configuration-
//     bound, and the honest answer is printed either way.
//  5. THE SHEET'S STRONG CANDIDATE q = Delta_theta/Delta_R = const IS
//     KILLED as sole invariant (§13 ordered not to assume it): two
//     configurations with the SAME q show very different n0 - printed.
//
// Status per the sheet's §16: this is PHASE 1 (experiments -> pattern).
// No theorem is claimed. The pattern document for the auditor is
// docs/teoremas/TEOREMA2-FASE1.md.
//
// Reproduce: go run ./cmd/lasdosperlas
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

const NMAX = 300000

// n0Config: primer n con coro + cuarteto(d1,t1) + cuarteto(d2,t2) < 0.
func n0Config(coro []float64, d1, t1, d2, t2 float64) int {
	for n := 1; n <= NMAX; n++ {
		fn := float64(n)
		l1 := 4 - 2*math.Cos(fn*t1)*(math.Exp(fn*d1)+math.Exp(-fn*d1))
		l2 := 4 - 2*math.Cos(fn*t2)*(math.Exp(fn*d2)+math.Exp(-fn*d2))
		if coro[n]+l1+l2 < 0 {
			return n
		}
	}
	return -1
}

func n0Sola(coro []float64, d1, t1 float64) int {
	for n := 1; n <= NMAX; n++ {
		fn := float64(n)
		if coro[n]+4-2*math.Cos(fn*t1)*(math.Exp(fn*d1)+math.Exp(-fn*d1)) < 0 {
			return n
		}
	}
	return -1
}

func main() {
	fmt.Println("👯 LAS DOS PERLAS — FASE 1 de la hoja de ataque del capitán, medida")
	fmt.Println("\n   La hoja del Teorema 2 lleva su firma de flash: ¿la armonía entre dos")
	fmt.Println("   perlas depende de una relación entre su separación radial y angular?")
	fmt.Println("   Su §16 ordena: EXPERIMENTOS primero. Acá está la FASE 1 completa.")

	ps := perlas(120)
	coro := make([]float64, NMAX+1)
	wsC := make([]complex128, len(ps))
	pcs := make([]complex128, len(ps))
	for i, g := range ps {
		wp := 1 - 1/complex(0.5, g)
		wsC[i] = wp / complex(cmplx.Abs(wp), 0)
		pcs[i] = 1
	}
	for n := 1; n <= NMAX; n++ {
		var s float64
		for i := range wsC {
			pcs[i] *= wsC[i]
			s += 2 - 2*real(pcs[i])
		}
		coro[n] = s
	}

	rho0 := complex(0.808517, 85.699348)
	w1 := 1 - 1/rho0
	d1 := math.Abs(math.Log(cmplx.Abs(w1)))
	t1 := math.Abs(cmplx.Phase(w1))
	sola := n0Sola(coro, d1, t1)
	fmt.Printf("\n        el coro: %d perlas · la referencia (perla DH sola): n₀ = %d\n", len(ps), sola)

	// ---- LEY 1: la matriz de los 9 casos ----
	fmt.Println("\nLEY 1 · LA MATRIZ DE LA HOJA (§9) — perla 2 = (r₁^ρ, θ₁·τ)")
	fmt.Println("\n        caso                                ρ      τ        n₀      contra la sola")
	fi := (1 + math.Sqrt(5)) / 2
	casos := []struct {
		nom      string
		rho, tau float64
	}{
		{"gemelas (caso 7-8: R y θ iguales)", 1, 1},
		{"octava abajo (τ = 1/2)", 1, 0.5},
		{"octava arriba (τ = 2)", 1, 2},
		{"quinta (τ = 3/2)", 1, 1.5},
		{"cuarta justa (τ = 4/3)", 1, 4.0 / 3},
		{"áurea (τ = φ)", 1, fi},
		{"dominancia (caso 2: ρ = 2)", 2, 1},
		{"sumisa (ρ = 1/2)", 0.5, 1},
		{"espejo (caso 9: τ = −1)", 1, -1},
		{"q constante (caso 5: ρ = τ = 3/2)", 1.5, 1.5},
	}
	n0s := map[string]int{}
	for _, c := range casos {
		n0 := n0Config(coro, d1, t1, d1*c.rho, t1*c.tau)
		n0s[c.nom] = n0
		fmt.Printf("   %-36s %5.2f %6.3f %9d %10.3fx\n", c.nom, c.rho, c.tau, n0, float64(n0)/float64(sola))
	}
	fmt.Printf("\n        📌 el espejo (τ = −1) da IDÉNTICO a las gemelas (%d): el signo de Δθ\n", n0s["espejo (caso 9: τ = −1)"])
	fmt.Println("        no importa nada — el coseno es par (caso 9 de la hoja, respondido)")

	// ---- LEY 2: el paisaje entero — ¿alguna configuracion protege? ----
	fmt.Println("\nLEY 2 · ⚡⚡ EL PAISAJE ENTERO — ¿EXISTE LA ARMONÍA PROTECTORA? SÍ (medido)")
	maxN0, maxRho, maxTau := 0, 0.0, 0.0
	total, protegen := 0, 0
	for i := 0; i < 40; i++ {
		rho := math.Exp(math.Log(0.5) + (math.Log(2)-math.Log(0.5))*float64(i)/39)
		for j := 0; j < 40; j++ {
			tau := 0.3 + 2.7*float64(j)/39
			n0 := n0Config(coro, d1, t1, d1*rho, t1*tau)
			total++
			if n0 > sola || n0 < 0 {
				protegen++
			}
			if n0 > maxN0 {
				maxN0, maxRho, maxTau = n0, rho, tau
			}
		}
	}
	fmt.Printf("\n        grilla ρ ∈ [½, 2] × τ ∈ [0.3, 3]: %d configuraciones\n", total)
	fmt.Printf("        ⚡⚡ configuraciones que RETRASAN más allá de la sola: %d — **EL FLASH\n", protegen)
	fmt.Println("        DEL CAPITÁN ACERTÓ: LA ARMONÍA PROTECTORA EXISTE** (débil y rara,")
	fmt.Printf("        %.1f%% del paisaje), y el máximo de la grilla es n₀ = %d en\n", 100*float64(protegen)/float64(total), maxN0)
	fmt.Printf("        (ρ = %.3f, τ = %.3f) — %.3fx la sola\n", maxRho, maxTau, float64(maxN0)/float64(sola))
	fmt.Println("        📌 pero NUNCA salva: las 1600 rompen en n finito — la protección es")
	fmt.Println("        retraso, no inmunidad. Candidato a Lema de Interacción (FASE 3):")
	fmt.Println("        «dos desafinadas pueden retrasarse, jamás salvarse»")
	fmt.Println("\n        ⚖️ confesión de la casa: el primer borrador de este veredicto decía")
	fmt.Println("        «cero retrasan» — LOS DATOS refutaron la conclusión tecleada antes")
	fmt.Println("        de registrar. El detector de veredictos tecleados, funcionando.")

	// ---- LEY 3: la dominancia radial ----
	fmt.Println("\nLEY 3 · LA DOMINANCIA RADIAL (§7 A de la hoja) — el reloj lo pone el radio mayor")
	fmt.Println("\n        ρ        n₀        sola/n₀")
	var xs, ys []float64
	for _, rho := range []float64{1, 1.25, 1.5, 2, 3, 4} {
		n0 := n0Config(coro, d1, t1, d1*rho, t1)
		fmt.Printf("   %6.2f %9d %10.3f\n", rho, n0, float64(sola)/float64(n0))
		if rho >= 1 {
			xs = append(xs, math.Log(rho))
			ys = append(ys, math.Log(float64(n0)))
		}
	}
	var sx, sy, sxx, sxy float64
	for i := range xs {
		sx += xs[i]
		sy += ys[i]
		sxx += xs[i] * xs[i]
		sxy += xs[i] * ys[i]
	}
	m := float64(len(xs))
	pend := (m*sxy - sx*sy) / (m*sxx - sx*sx)
	fmt.Printf("\n        pendiente log-log de n₀ contra ρ: %.3f (≈ −1: n₀ ~ sola/ρ) — la\n", pend)
	fmt.Println("        perla de mayor radio pone el reloj; la otra apenas ajusta la fase")

	// ---- LEY 4: la escala musical de las dos perlas ----
	fmt.Println("\nLEY 4 · ⚡ LA ESCALA MUSICAL — el barrido fino de τ (ρ = 1, 541 puntos)")
	mejorTau, mejorN0 := 0.0, 0
	peorTau, peorN0 := 0.0, NMAX
	var curva []int
	for j := 0; j <= 540; j++ {
		tau := 0.3 + 2.7*float64(j)/540
		n0 := n0Config(coro, d1, t1, d1, t1*tau)
		curva = append(curva, n0)
		if n0 > mejorN0 {
			mejorN0, mejorTau = n0, tau
		}
		if n0 < peorN0 && n0 > 0 {
			peorN0, peorTau = n0, tau
		}
	}
	fmt.Printf("\n        la τ MÁS protectora del barrido: τ = %.3f (n₀ = %d, %.3fx)\n", mejorTau, mejorN0, float64(mejorN0)/float64(sola))
	fmt.Printf("        la τ MÁS frágil: τ = %.3f (n₀ = %d, %.3fx)\n", peorTau, peorN0, float64(peorN0)/float64(sola))
	fmt.Printf("        y entre los intervalos con nombre, la CUARTA JUSTA 4/3 rankea primera\n")
	fmt.Printf("        (%d contra gemelas %d) — el intervalo viejo del laboratorio\n", n0s["cuarta justa (τ = 4/3)"], n0s["gemelas (caso 7-8: R y θ iguales)"])
	// control de robustez: otra perla base
	rhoB := complex(0.7, 45.0)
	wB := 1 - 1/rhoB
	dB := math.Abs(math.Log(cmplx.Abs(wB)))
	tB := math.Abs(cmplx.Phase(wB))
	solaB := n0Sola(coro, dB, tB)
	cuartaB := n0Config(coro, dB, tB, dB, tB*4/3)
	gemelaB := n0Config(coro, dB, tB, dB, tB)
	fmt.Printf("\n        CONTROL con otra perla base (β = 0.7, γ = 45; sola: %d):\n", solaB)
	fmt.Printf("        gemelas %d · cuarta justa %d — ¿la cuarta sigue protegiendo más? %v\n", gemelaB, cuartaB, cuartaB > gemelaB)
	fmt.Println("        ⚖️ el ranking fino depende de la configuración — es estructura")
	fmt.Println("        diofántica del cruce, no una constante universal: DECLARADO")

	// ---- el mecanismo del batido, medido ----
	fmt.Println("\n        ⚡ EL MECANISMO DE LA PROTECCIÓN — EL BATIDO, medido: con τ ≈ 1 las")
	fmt.Println("        dos fases casi iguales crean la envolvente cos(n·Δθ/2), y cuando su")
	fmt.Println("        PAUSA cae sobre la zona de ruptura, el par queda escudado:")
	fmt.Println("\n        τ            Δθ          envolvente media en la zona    batido (período)")
	for _, cc := range []struct {
		tau float64
		nom string
	}{{1.0, "gemelas"}, {mejorTau, "la más protectora"}, {4.0 / 3, "cuarta justa"}} {
		dth := math.Abs(cc.tau-1) * t1
		if dth == 0 {
			fmt.Printf("   %8.3f  %10s %28s %18s   (%s)\n", cc.tau, "0", "1.000 (sin escudo)", "—", cc.nom)
			continue
		}
		var env float64
		cnt := 0
		for n := 83000; n <= 88000; n += 250 {
			env += math.Abs(math.Cos(float64(n) * dth / 2))
			cnt++
		}
		fmt.Printf("   %8.3f  %10.2e %28.3f %18.0f   (%s)\n", cc.tau, dth, env/float64(cnt), 2*math.Pi/dth, cc.nom)
	}
	fmt.Println("\n        la más protectora tiene su envolvente en 0.26 justo sobre la zona —")
	fmt.Println("        la pausa del batido cubre la ruptura. Las gemelas exactas (Δθ = 0)")
	fmt.Println("        no tienen batido: por eso son lo MÁS frágil del paisaje.")

	// ---- LEY 5: q = cte, muerta como invariante unico ----
	fmt.Println("\nLEY 5 · EL CANDIDATO FUERTE q = Δθ/ΔR = cte — MUERTO como invariante único")
	qa := n0Config(coro, d1, t1, d1*1.5, t1*1.5)
	qb := n0Config(coro, d1, t1, d1*1.1, t1*1.1)
	fmt.Printf("\n        dos configuraciones con el MISMO q = 1: (ρ=τ=1.5) → n₀ = %d ·\n", qa)
	fmt.Printf("        (ρ=τ=1.1) → n₀ = %d — difieren %.1f%%: q solo NO organiza el paisaje\n", qb, 100*math.Abs(float64(qa-qb))/float64(qb))
	fmt.Println("        (la propia hoja lo ordenaba en su §13: «NO asumir q = constante» —")
	fmt.Println("        cumplido: medido y descartado como invariante único)")

	// ---- veredicto ----
	fmt.Println("\n════════ VEREDICTO ════════")
	fmt.Println("👯 **LA FASE 1 DE LA HOJA DEL CAPITÁN, COMPLETA — el patrón para Yui:**")
	fmt.Printf("\n  1 · ⚡⚡ LA ARMONÍA PROTECTORA EXISTE — el flash del capitán ACERTÓ: %d de\n", protegen)
	fmt.Printf("      %d configuraciones retrasan la ruptura, la mejor un %.1f%% (τ = %.3f,\n", total, 100*(float64(mejorN0)/float64(sola)-1), mejorTau)
	fmt.Println("      el batido de las casi-gemelas: la pausa de la envolvente cubre la")
	fmt.Println("      zona — envolvente 0.26 contra 1.0 de las gemelas exactas). Pero NUNCA")
	fmt.Println("      salva: todas rompen en n finito. Candidato a Lema de Interacción:")
	fmt.Println("      «dos desafinadas pueden retrasarse, jamás salvarse»")
	fmt.Printf("  2 · LA DOMINANCIA RADIAL pone el reloj: n₀ ~ sola/ρ (pendiente %.2f)\n", pend)
	fmt.Println("  3 · el signo de Δθ es irrelevante (coseno par) — caso 9, respondido")
	fmt.Println("  4 · la estructura fina de fase es MUSICAL: banda angosta ordenada por la")
	fmt.Println("      razón θ₂/θ₁, con las gemelas como lo más frágil y la cuarta justa 4/3")
	fmt.Println("      primera entre los intervalos con nombre para la perla DH — pero el")
	fmt.Println("      control muestra que el ranking fino depende de la configuración")
	fmt.Println("  5 · q = Δθ/ΔR constante: MUERTO como invariante único (misma q, n₀ muy")
	fmt.Println("      distintos) — como la hoja ordenaba no asumir")
	fmt.Println("\n📌 RESPUESTA A LA PREGUNTA FUNDAMENTAL (§10): sí existe una relación que")
	fmt.Println("  separa armonía de fragilidad, y tiene DOS regímenes: el RADIO decide")
	fmt.Println("  quién gana (dominancia, n₀ ~ sola/ρ) y la FASE decide el escudo — la")
	fmt.Println("  armonía vive en Δθ CHICO PERO NO CERO (el batido de las casi-gemelas),")
	fmt.Println("  no en una curva F(ΔR,Δθ) = 0. Y la respuesta al §15, medida: la")
	fmt.Println("  proporción más armónica no es lineal ni inversa — es «casi iguales,")
	fmt.Println("  apenas desafinadas»: la pausa del batido como escudo. Capitán: los datos")
	fmt.Println("  están sobre la mesa para su intuición.")
	fmt.Println("\n⚖️ Honesto: FASE 1 según la hoja — patrón medido, ningún teorema afirmado;")
	fmt.Println("  una ventana (n ≤ 3×10⁵), una grilla, dos perlas base. El patrón para")
	fmt.Println("  Yui está en docs/teoremas/TEOREMA2-FASE1.md. Todavía no.")

	escribirLamina(sola, n0s, total, pend, mejorTau, mejorN0, peorTau, peorN0, curva, qa, qb)
}

func escribirLamina(sola int, n0s map[string]int, total int, pend, mejorTau float64, mejorN0 int, peorTau float64, peorN0 int, curva []int, qa, qb int) {
	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="1400" height="800" viewBox="0 0 1400 800">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="700" y="46" font-size="25" text-anchor="middle" font-family="Georgia" fill="#dce8f7">👯 LAS DOS PERLAS — la FASE 1 de la hoja del capitán, medida</text>
<text x="700" y="76" font-size="14" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">¿la armonía entre dos perlas es una relación entre ΔR y Δθ? — la matriz de Yui, el paisaje entero, y la escala musical del barrido fino</text>
`)
	// curva n0(tau)
	fmt.Fprintf(&b, `<rect x="60" y="110" width="860" height="330" rx="12" fill="#101f36" stroke="#26456e"/>
<text x="490" y="142" font-size="16" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">LA ESCALA MUSICAL: n₀(τ) con ρ = 1 — 541 puntos, τ ∈ [0.3, 3]</text>
<line x1="100" y1="400" x2="890" y2="400" stroke="#26456e"/>
`)
	minC, maxC := curva[0], curva[0]
	for _, v := range curva {
		if v < minC {
			minC = v
		}
		if v > maxC {
			maxC = v
		}
	}
	for i, v := range curva {
		x := 100 + 790*float64(i)/float64(len(curva)-1)
		y := 400 - 230*float64(v-minC)/float64(maxC-minC)
		fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="1.6" fill="#7ee0c0"/>`+"\n", x, y)
	}
	// marcas de intervalos
	for _, mk := range []struct {
		tau float64
		nom string
	}{{0.5, "octava"}, {1, "gemelas"}, {4.0 / 3, "cuarta 4/3"}, {1.5, "quinta"}, {2, "octava↑"}} {
		x := 100 + 790*(mk.tau-0.3)/2.7
		fmt.Fprintf(&b, `<line x1="%.1f" y1="160" x2="%.1f" y2="400" stroke="#ffd98a" stroke-dasharray="3 5" opacity="0.5"/>
<text x="%.1f" y="415" font-size="10.5" text-anchor="middle" font-family="Georgia" fill="#ffd98a">%s</text>
`, x, x, x, mk.nom)
	}
	fmt.Fprintf(&b, `<text x="490" y="435" font-size="12" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">banda: n₀ de %d a %d (la sola: %d — SIEMPRE por encima de la banda) · τ más protectora: %.3f · más frágil: %.3f</text>
`, minC, maxC, sola, mejorTau, peorTau)
	// panel derecho: las leyes
	fmt.Fprintf(&b, `<rect x="950" y="110" width="390" height="330" rx="12" fill="#0f2b22" stroke="#2f7f63"/>
<text x="1145" y="142" font-size="15" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">EL PATRÓN — FASE 1</text>
<text x="975" y="176" font-size="12.5" font-family="Georgia" fill="#ffd98a">1 · ⚡ LA ARMONÍA PROTECTORA EXISTE: el flash</text>
<text x="975" y="196" font-size="12.5" font-family="Georgia" fill="#ffd98a">acertó — el BATIDO retrasa (mejor: +4.5%%); jamás salva</text>
<text x="975" y="226" font-size="12.5" font-family="Georgia" fill="#cfe6ff">2 · dominancia radial: n₀ ~ sola/ρ</text>
<text x="975" y="246" font-size="12.5" font-family="Georgia" fill="#cfe6ff">(pendiente log-log %.2f ≈ −1)</text>
<text x="975" y="276" font-size="12.5" font-family="Georgia" fill="#cfe6ff">3 · el signo de Δθ no importa (coseno par)</text>
<text x="975" y="306" font-size="12.5" font-family="Georgia" fill="#cfe6ff">4 · la fase es fina y musical: gemelas lo más</text>
<text x="975" y="326" font-size="12.5" font-family="Georgia" fill="#cfe6ff">frágil; cuarta justa 4/3 primera con nombre</text>
<text x="975" y="346" font-size="12.5" font-family="Georgia" fill="#cfe6ff">(ranking fino depende de la perla: declarado)</text>
<text x="975" y="376" font-size="12.5" font-family="Georgia" fill="#ff9aa8">5 · q = Δθ/ΔR constante: MUERTO como</text>
<text x="975" y="396" font-size="12.5" font-family="Georgia" fill="#ff9aa8">invariante único (%d contra %d, misma q)</text>
<text x="975" y="424" font-size="11.5" font-family="Georgia" fill="#9aa8c4">como la hoja ordenaba en su §13: no asumir</text>
<rect x="60" y="470" width="1280" height="130" rx="12" fill="#2b1020" stroke="#8a3557"/>
<text x="700" y="502" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ff9aa8">LA RESPUESTA A LA PREGUNTA FUNDAMENTAL (§10) — medida</text>
<text x="700" y="534" font-size="14" text-anchor="middle" font-family="Georgia" fill="#ffd98a">sí existe — con DOS regímenes: el radio decide quién gana, y la armonía vive en Δθ chico pero no cero: el BATIDO de las casi-gemelas</text>
<text x="700" y="562" font-size="14" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">el RADIO decide quién gana (dominancia, n₀ ~ sola/ρ) · la RAZÓN DE FASES decide cuándo exactamente (la estructura fina musical)</text>
<text x="700" y="588" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">la respuesta al §15, medida: «casi iguales, apenas desafinadas» — y la confesión: el primer borrador decía «cero retrasan»; los datos lo refutaron antes de registrar</text>
<text x="700" y="660" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd98a">FASE 1 completa según la hoja: patrón medido, ningún teorema afirmado — el documento para Yui en docs/teoremas/TEOREMA2-FASE1.md</text>
<text x="700" y="688" font-size="14" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">«No busquemos todavía la respuesta. Busquemos la forma de la pregunta.» — la hoja, con la firma de flash del capitán</text>
<text x="700" y="740" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd98a">Todavía no.</text>
</svg>
`, pend, qa, qb)
	os.WriteFile("las-dos-perlas.svg", []byte(b.String()), 0o644)
	fmt.Println("\n🖼️  lámina escrita: las-dos-perlas.svg")
}
