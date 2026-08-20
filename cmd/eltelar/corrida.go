package main

import (
	"fmt"
	"math"
)

// corrida.go - the four sections of the Phase II answer.

func main() {
	fmt.Println("🧵 EL TELAR — Fase II: ¿existe una estructura NO CIRCULAR cuya geometría")
	fmt.Println("   crezca como ln(E/2π) y cuyo espectro tenga eco en k·log p?")

	const t0, t1 = 100.0, 1000.0
	const tope = 3.2
	d := &dado{s: 20260817}

	// the true sea, used ONLY as a ruler
	verdad := cerosVerdaderos(t0, t1, 0.02)
	per := periodosAritmeticos(tope)
	azar := periodosAlAzar(d, len(per), per, tope)
	fmt.Printf("\n§0 · la regla · %d ceros verdaderos en [%.0f,%.0f] · %d períodos aritméticos hasta T=%.1f\n",
		len(verdad), t0, t1, len(per), tope)

	// -----------------------------------------------------------------------
	// §1 · THE R6 AUDIT: every candidate declares its inputs BEFORE being tested
	// -----------------------------------------------------------------------
	fmt.Println("\n§1 · AUDITORÍA R6 — cada candidato declara sus entradas ANTES de ser probado")
	type cand struct{ nombre, entradas, veredicto string }
	for _, c := range []cand{
		{"H = diag(γₙ)", "los γₙ", "CIRCULAR — regla de parada 1"},
		{"matriz GUE al azar", "una semilla y nada más", "UNIVERSAL SIN ARITMÉTICA — regla 2"},
		{"la caja que respira", "la ecuación funcional (θ), NINGÚN cero", "ADMISIBLE — se prueba abajo"},
		{"los primos solos", "la ecuación funcional (θ) + Λ(n) hasta un tope", "ADMISIBLE — se prueba abajo"},
	} {
		fmt.Printf("     %-22s entradas: %-48s → %s\n", c.nombre, c.entradas, c.veredicto)
	}
	fmt.Println("     nota: θ y la ley suave descienden de la ECUACIÓN FUNCIONAL, no de dónde están los ceros.")
	fmt.Println("     Si eso no se admite en D, los dos últimos candidatos caen y no queda ninguno.")

	// -----------------------------------------------------------------------
	// §2 · THE BREATHING BOX ALONE
	// -----------------------------------------------------------------------
	fmt.Println("\n§2 · LA CAJA QUE RESPIRA, SOLA — ¿qué da la geometría sin aritmética?")
	nivCaja := nivelesDe(suave, t0, t1)
	spCaja := desplegar(nivCaja)
	dC, fC, rC := murcielago(nivCaja, per, tope)
	_, _, rCa := murcielago(nivCaja, azar, tope)
	fmt.Printf("     largo efectivo: ln(T/2π) va de %.4f a %.4f en el tramo (por construcción)\n",
		largoCaja(t0), largoCaja(t1))
	fmt.Printf("     CONTEO      : %d niveles contra %d ceros verdaderos — diferencia %d ✓ (por construcción)\n",
		len(nivCaja), len(verdad), len(nivCaja)-len(verdad))
	fmt.Printf("     CORRELACIÓN : Σ²(10) = %.6f · Σ²(20) = %.6f   (una reja perfecta da ~0)\n",
		sigma2(spCaja, 10, 600), sigma2(spCaja, 20, 600))
	fmt.Printf("     ARITMÉTICA  : eco |E| dentro %.6f · fuera %.6f · razón %.3f · en períodos al azar %.3f\n",
		dC, fC, rC, rCa)
	var dif []float64
	for i := 0; i < len(nivCaja) && i < len(verdad); i++ {
		dif = append(dif, math.Abs(nivCaja[i]-verdad[i]))
	}
	fmt.Printf("     IDENTIDAD   : |nivel − γ| medio %.4f · peor %.4f (en unidades de altura)\n", media(dif), maxi(dif))
	fmt.Println("     ⟹ la caja que respira acierta el CONTEO y no sabe nada más: sin estructura fina,")
	fmt.Println("       sin correlaciones, sin eco. La geometría correcta es necesaria y VACÍA.")
	fmt.Println("     ⚠ y es ESTIPULADA, no derivada: le pusimos ln(T/2π) a mano. Bajo el contrato R6 eso")
	fmt.Println("       la deja en TAUTOLÓGICA hasta que la ley se derive de algo que no sean los ceros.")

	// -----------------------------------------------------------------------
	// §3 · THE PRIMES ALONE
	// -----------------------------------------------------------------------
	fmt.Println("\n§3 · LOS PRIMOS SOLOS — la mitad aritmética, sin mirar un solo cero")
	fmt.Println("     N_p(T) = θ(T)/π + 1 − (1/π)·Σ_{n=p^k ≤ P} Λ(n)·sin(T·log n)/(√n·log n)")
	fmt.Printf("     %8s %7s %11s %11s %11s %9s %9s\n", "P", "términos", "|niv−γ| medio", "peor", "Σ²(10)", "eco k·logp", "eco azar")
	var filas []filaP
	var mejor []float64
	for _, P := range []int{10, 100, 1000, 10000, 100000} {
		ns, lam := mangoldt(P)
		cuenta := func(T float64) float64 { return suave(T) + sPrimos(T, ns, lam) }
		niv := nivelesDe(cuenta, t0, t1)
		sp := desplegar(niv)
		var dd []float64
		for i := 0; i < len(niv) && i < len(verdad); i++ {
			dd = append(dd, math.Abs(niv[i]-verdad[i]))
		}
		_, _, r := murcielago(niv, per, tope)
		_, _, ra := murcielago(niv, azar, tope)
		s2 := sigma2(sp, 10, 600)
		fmt.Printf("     %8d %7d %11.5f %11.5f %11.4f %9.3f %9.4f\n",
			P, len(ns), media(dd), maxi(dd), s2, r, ra)
		filas = append(filas, filaP{P, len(ns), media(dd), maxi(dd), s2, r, ra})
		mejor = niv
	}
	fmt.Println("     ⚠ TAUTOLOGÍA DECLARADA: este candidato lleva log n en su definición, así que su ECO está")
	fmt.Println("       garantizado por construcción y NO informa (regla §14.6 del contrato). Las columnas que")
	fmt.Println("       sí informan son Σ² y la distancia a los γₙ — y ésas no estaban puestas a mano.")
	fmt.Printf("     para comparar, los ceros VERDADEROS con el mismo protocolo: Σ²(10) = %.4f · eco %.3f · azar %.4f\n",
		sigma2(desplegar(verdad), 10, 600), razonDe(verdad, per, tope), razonDe(verdad, azar, tope))

	// how many levels land within a fraction of a mean gap
	if len(mejor) > 0 {
		var d1, d2, d3 int
		for i := 0; i < len(mejor) && i < len(verdad); i++ {
			g := 2 * math.Pi / largoCaja(verdad[i]) // local mean gap
			e := math.Abs(mejor[i]-verdad[i]) / g
			if e < 0.5 {
				d1++
			}
			if e < 0.1 {
				d2++
			}
			if e < 0.01 {
				d3++
			}
		}
		n := len(verdad)
		fmt.Printf("     con el mejor tope: %d de %d niveles dentro de medio espaciado (%.1f%%), %d dentro de 0,1 (%.1f%%), %d dentro de 0,01 (%.1f%%)\n",
			d1, n, 100*float64(d1)/float64(n), d2, 100*float64(d2)/float64(n), d3, 100*float64(d3)/float64(n))
	}

	// -----------------------------------------------------------------------
	// §4 · THE VERDICT
	// -----------------------------------------------------------------------
	fmt.Println("\n§4 · LAS CUATRO CAPAS, SEPARADAS COMO PIDIÓ LA AUDITORA")
	fmt.Printf("     %-24s %-10s %-14s %-12s %s\n", "candidato", "conteo", "correlaciones", "identidad", "aritmética")
	fmt.Printf("     %-24s %-10s %-14s %-12s %s\n", "H = diag(γₙ)", "sí", "sí", "sí", "sí — CIRCULAR")
	fmt.Printf("     %-24s %-10s %-14s %-12s %s\n", "matriz GUE", "reescalable", "sí", "no", "NO")
	fmt.Printf("     %-24s %-10s %-14s %-12s %s\n", "caja que respira", "sí", "NO (rígida)", "no", "NO")
	fmt.Printf("     %-24s %-10s %-14s %-12s %s\n", "los primos solos", "sí", "medida arriba", "medida arriba", "medida arriba")

	// -----------------------------------------------------------------------
	// §5 · THE SIGN, MEASURED - not argued
	// -----------------------------------------------------------------------
	fmt.Println("\n§5 · EL SIGNO, MEDIDO SOBRE NUESTROS PROPIOS CEROS")
	fmt.Println("     la fórmula de Selberg entra con + en las órbitas; la explícita con − en los primos.")
	fmt.Println("     eso NO es convención: el coeficiente de Fourier de la densidad en τ = k·log p tiene")
	fmt.Println("     signo propio. Medido con ventana sobre los 620 ceros:")
	fmt.Printf("     %10s %8s %14s %14s %s\n", "período", "n=p^k", "D(τ) medido", "−Λ(n)/(π√n)", "signo")
	neg, tot := 0, 0
	var signos [][3]float64
	for _, tau := range per {
		n := math.Exp(tau)
		nr := math.Round(n)
		D := coefDensidad(verdad, tau, t0, t1)
		pred := -mangoldtDe(int(nr)) / (math.Pi * math.Sqrt(nr))
		sg := "NEGATIVO ✓"
		if D > 0 {
			sg = "positivo ✗"
		} else {
			neg++
		}
		tot++
		fmt.Printf("     %10.5f %8.0f %14.6f %14.6f %s\n", tau, nr, D, pred, sg)
		signos = append(signos, [3]float64{tau, D, pred})
	}
	fmt.Printf("     ⟹ %d de %d períodos aritméticos dan coeficiente NEGATIVO — el espectro de absorción,\n", neg, tot)
	fmt.Println("       leído en nuestros propios datos. Es una restricción de diseño: el candidato debe")
	fmt.Println("       DERIVAR ese signo, no elegirlo. (Hejhal 1976, Berry 1986; la lectura es de Connes.)")

	dibujar(Resultado{
		T0: t0, T1: t1, Verdad: verdad, Periodos: per, Azar: azar,
		NivCaja: nivCaja, S2caja10: sigma2(spCaja, 10, 600), S2caja20: sigma2(spCaja, 20, 600),
		EcoCaja: rC, EcoCajaAzar: rCa, DifCaja: media(dif),
		Filas: filasPlanas(filas), MejorNiv: mejor,
		S2verdad:  sigma2(desplegar(verdad), 10, 600),
		EcoVerdad: razonDe(verdad, per, tope), EcoVerdadAzar: razonDe(verdad, azar, tope),
		LargoT0: largoCaja(t0), LargoT1: largoCaja(t1), Tope: tope,
		Signos: signos, SignosNeg: neg, SignosTot: tot,
	})
}

func razonDe(niv, per []float64, tope float64) float64 {
	_, _, r := murcielago(niv, per, tope)
	return r
}

func maxi(v []float64) float64 {
	if len(v) == 0 {
		return math.NaN()
	}
	m := v[0]
	for _, x := range v {
		if x > m {
			m = x
		}
	}
	return m
}

// filaP is one truncation of the prime-only candidate.
type filaP struct {
	P                       int
	terms                   int
	medio, peor, s2, r, raz float64
}

func filasPlanas(f []filaP) [][7]float64 {
	out := make([][7]float64, len(f))
	for i, x := range f {
		out[i] = [7]float64{float64(x.P), float64(x.terms), x.medio, x.peor, x.s2, x.r, x.raz}
	}
	return out
}

// Resultado carries every measured number to the plate.
type Resultado struct {
	T0, T1                        float64
	Verdad, Periodos, Azar        []float64
	NivCaja                       []float64
	S2caja10, S2caja20            float64
	EcoCaja, EcoCajaAzar, DifCaja float64
	Filas                         [][7]float64 // P, terms, mean, worst, sigma2, echo, echo-random
	MejorNiv                      []float64
	S2verdad                      float64
	EcoVerdad, EcoVerdadAzar      float64
	LargoT0, LargoT1, Tope        float64
	Signos                        [][3]float64 // period, measured D, predicted -Lambda/(pi sqrt n)
	SignosNeg, SignosTot          int
}

// coefDensidad is the Fourier coefficient of the level density at period tau,
// with the smooth part removed: the quantity whose SIGN is the invariant that
// separates a Selberg system (+) from the Riemann zeros (-).
func coefDensidad(niv []float64, tau, t0, t1 float64) float64 {
	M := len(niv)
	s := 0.0
	for i, g := range niv {
		w := 0.5 - 0.5*math.Cos(2*math.Pi*(float64(i)+0.5)/float64(M))
		s += w * math.Cos(tau*g)
	}
	// the same window applied to the smooth density, subtracted
	const pasos = 20000
	h := (t1 - t0) / pasos
	sm := 0.0
	for k := 0; k < pasos; k++ {
		E := t0 + (float64(k)+0.5)*h
		x := (E - t0) / (t1 - t0)
		w := 0.5 - 0.5*math.Cos(2*math.Pi*x)
		sm += w * math.Log(E/(2*math.Pi)) / (2 * math.Pi) * math.Cos(tau*E) * h
	}
	return 2 * (s - sm) / float64(M)
}

// mangoldtDe returns Lambda(n) for a single integer.
func mangoldtDe(n int) float64 {
	for p := 2; p*p <= n; p++ {
		if n%p == 0 {
			m := n
			for m%p == 0 {
				m /= p
			}
			if m == 1 {
				return math.Log(float64(p))
			}
			return 0
		}
	}
	if n > 1 {
		return math.Log(float64(n))
	}
	return 0
}
