// Command elimpostor takes Finding 229's impostor back out of the drawer and
// harmonizes it at dimension 0, which is the captain's order.
//
// THE IMPOSTOR, as it was registered:
//
//	P(s) = (s-a)(s-conj a)(s-(1-a))(s-(1-conj a)),  a = 0.7 + 3i
//
// It carries EVERY symmetry of the book - functional equation, Schwarz
// reflection, the shapeshifter sigma - to 1.6e-16, and its four roots sit at
// Re = 0.70 and Re = 0.30. Off the line. That is why F229 killed symmetry-alone
// as a route, and why F259's Davenport-Heilbronn buried it.
//
// SO THE GEOMETRY CANNOT SEE IT. THE QUESTION IS WHETHER DIMENSION 0 CAN.
//
// Carry the roots through the shapeshifter w = 1 - 1/rho and build Li's price
// over the pairs, exactly as F232 does for the real pearls:
//
//	lambda_n = sum over conjugate pairs of [ 2 - 2 Re(w^n) ]
//
// On the skin |w| = 1 and every term sits in [0, 4]. Off the skin one pair has
// |w| < 1 and the other |w| > 1 - their product is exactly 1, which is F225's
// north times south - and the second one grows like r^n. The moment its phase
// comes round, the price goes NEGATIVE and stays catchable.
//
// THE HONEST LIMIT, AND IT IS THE WHOLE POINT. This works because the impostor
// has FOUR roots and we know all four. Zeta has infinitely many and we hold 649,
// which is why F259 measured the laboratory's horizon at gamma ~ 1658. Dimension
// 0 convicts an impostor you can enumerate. Zeta is not enumerable, and neither
// is Davenport-Heilbronn - of which we know exactly one off-line zero. Same wall,
// same reason.
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

func w(ρ complex128) complex128 { return 1 - 1/ρ }

// aporte is one conjugate pair's share of the price, F232's unconditional form.
func aporte(n int, ρ complex128) float64 {
	return 2 - 2*real(cmplx.Pow(w(ρ), complex(float64(n), 0)))
}

func main() {
	fmt.Println("🎭 EL IMPOSTOR, ARMONIZADO EN LA DIMENSIÓN 0")
	fmt.Println("\n   orden del capitán: «traeme de nuevo la fórmula del impostor… y armonicémosla")
	fmt.Println("   en la dimensión 0».")

	a := complex(0.7, 3.0)
	raices := []complex128{a, cmplx.Conj(a), 1 - a, 1 - cmplx.Conj(a)}
	P := func(s complex128) complex128 {
		p := complex(1, 0)
		for _, r := range raices {
			p *= s - r
		}
		return p
	}
	sigma := func(s complex128) complex128 { return 1 - cmplx.Conj(s) }

	// ---- LEY 1: la fórmula, de vuelta sobre la mesa ----
	fmt.Println("\nLEY 1 · LA FÓRMULA, TAL CUAL QUEDÓ EN F229")
	fmt.Println("\n        P(s) = (s−a)(s−ā)(s−(1−a))(s−(1−ā))        con  a = 0.7 + 3i")
	fmt.Println("\n   Sus cuatro raíces:")
	for _, r := range raices {
		fmt.Printf("        %+.4f %+.4fi     Re = %.2f  ← %s\n", real(r), imag(r), real(r),
			map[bool]string{true: "SOBRE la línea", false: "FUERA de la línea"}[math.Abs(real(r)-0.5) < 1e-12])
	}
	fmt.Println("\n   Y tiene TODAS las simetrías del libro. Re-verificado ahora, no de memoria:")
	pruebas := []complex128{complex(0.3, 7.2), complex(0.8, 13.6), complex(2.0, 21.3), complex(-0.7, 9.1)}
	peorFE, peorSch, peorSig := 0.0, 0.0, 0.0
	for _, s := range pruebas {
		den := math.Max(cmplx.Abs(P(s)), 1e-300)
		if d := cmplx.Abs(P(s)-P(1-s)) / den; d > peorFE {
			peorFE = d
		}
		if d := cmplx.Abs(P(cmplx.Conj(s))-cmplx.Conj(P(s))) / den; d > peorSch {
			peorSch = d
		}
		if d := cmplx.Abs(P(sigma(s))-cmplx.Conj(P(s))) / den; d > peorSig {
			peorSig = d
		}
	}
	fmt.Printf("        ecuación funcional  P(s) = P(1−s) .......... %.1e\n", peorFE)
	fmt.Printf("        espejo de Schwarz   P(s̄) = conj P(s) ...... %.1e\n", peorSch)
	fmt.Printf("        el cambiaformas     P(σ(s)) = conj P(s) .... %.1e\n", peorSig)
	fmt.Println("   → perfecto en las tres. La geometría NO lo puede distinguir de ξ.")

	// ---- LEY 2: al disco ----
	fmt.Println("\nLEY 2 · AL DISCO — Y ACÁ YA SE LE VE EL DISFRAZ")
	fmt.Println("   Pasamos las cuatro raíces por el cambiaformas w = 1 − 1/ρ. Sobre la línea")
	fmt.Println("   toda perla cae EN la piel (|w| = 1). Estas no:")
	fmt.Println("\n        raíz                    w = 1 − 1/ρ              |w|          ¿en la piel?")
	var wA, wB complex128
	for i, r := range raices {
		ww := w(r)
		if i == 0 {
			wA = ww
		}
		if i == 2 {
			wB = ww
		}
		fmt.Printf("   %-16s %-24s %-12.9f %s\n",
			fmt.Sprintf("%+.2f%+.2fi", real(r), imag(r)),
			fmt.Sprintf("%+.6f%+.6fi", real(ww), imag(ww)),
			cmplx.Abs(ww),
			map[bool]string{true: "sí", false: "NO — está afuera"}[math.Abs(cmplx.Abs(ww)-1) < 1e-9])
	}
	prod := cmplx.Abs(wA) * cmplx.Abs(wB)
	fmt.Printf("\n   ⚡ Y fijate el detalle: un par tiene |w| < 1 y el otro |w| > 1, y su producto es\n")
	fmt.Printf("     %.15f — o sea EXACTAMENTE 1 (desvío %.1e). Es el norte × sur = 1 de F225.\n",
		prod, math.Abs(prod-1))
	fmt.Println("     El impostor no puede escapar de esa ley: si uno se hunde, el otro se levanta.")
	fmt.Println("     Y ahí está su perdición, porque el que se levanta CRECE COMO rⁿ.")

	// ---- LEY 3: el precio ----
	fmt.Println("\nLEY 3 · ⚡ EL PRECIO EN LA DIMENSIÓN 0 — Y ACÁ LO AGARRA")
	fmt.Println("   Armamos el precio de Li con la fórmula incondicional de F232, la misma que")
	fmt.Println("   usamos con las perlas de verdad:")
	fmt.Println("\n        λₙ = Σ sobre pares {ρ, ρ̄} de [ 2 − 2·Re(wⁿ) ]")
	fmt.Println("\n   Sobre la piel cada término vive entre 0 y 4 y NUNCA puede ser negativo.")
	fmt.Println("   Fuera de la piel no hay techo. Medido sobre el impostor:")
	fmt.Println("\n        n         λₙ del impostor         ¿negativo?")
	lam := func(n int) float64 { return aporte(n, a) + aporte(n, 1-a) }
	primerNeg := -1
	for _, n := range []int{1, 2, 5, 10, 25, 50, 100, 200} {
		v := lam(n)
		fmt.Printf("   %7d   %22.6f   %s\n", n, v,
			map[bool]string{true: "◀── SÍ, CAYÓ", false: "no"}[v < 0])
	}
	for n := 1; n <= 100000; n++ {
		if lam(n) < 0 {
			primerNeg = n
			break
		}
	}
	peor, peorN := 0.0, 0
	for n := 1; n <= 2000; n++ {
		if v := lam(n); v < peor {
			peor, peorN = v, n
		}
	}
	fmt.Printf("\n   → PRIMER NEGATIVO EN n = %d. Y sigue empeorando: en n = %d vale %.4g.\n", primerNeg, peorN, peor)
	fmt.Println("   ⟹ LA DIMENSIÓN 0 LO CONDENA. Donde la simetría no vio nada, el precio lo canta.")

	// ---- LEY 4: contra las perlas de verdad ----
	fmt.Println("\nLEY 4 · Y CONTRA LAS PERLAS DE VERDAD, PARA QUE NO SEA UN TRUCO")
	fmt.Println("\npescando perlas hasta t=500…")
	ps := perlas(500)
	lamReal := func(n int) float64 {
		s := 0.0
		for _, g := range ps {
			s += aporte(n, complex(0.5, g))
		}
		return s
	}
	fmt.Printf("perlas: %d\n", len(ps))
	fmt.Println("\n        n        λₙ del impostor        λₙ de ζ (medido)")
	for _, n := range []int{1, 10, 50, 100, 200} {
		fmt.Printf("   %7d   %20.6f   %20.6f\n", n, lam(n), lamReal(n))
	}
	negReal := 0
	for n := 1; n <= 200; n++ {
		if lamReal(n) < 0 {
			negReal++
		}
	}
	fmt.Printf("   → el impostor cae en n = %d; ζ NO cae ni una vez en n = 1..200 (%d negativos).\n", primerNeg, negReal)
	fmt.Println("     El instrumento distingue. No es que diga que todo está bien: dice la verdad")
	fmt.Println("     sobre el impostor y la verdad sobre ζ, y son distintas.")

	// ---- LEY 5: el límite ----
	fmt.Println("\nLEY 5 · ⚖️ Y AHORA EL LÍMITE, QUE ES TODO EL ASUNTO")
	fmt.Println("   ¿Por qué funcionó acá y no cierra el problema? Por UNA razón, y es aritmética")
	fmt.Println("   pura, no falta de esfuerzo:")
	fmt.Println("\n        EL IMPOSTOR TIENE 4 RAÍCES Y LAS CONOCEMOS A LAS CUATRO.")
	fmt.Println("\n   El precio de Li es una suma sobre TODAS las raíces. Con cuatro, la suma es")
	fmt.Println("   exacta y el veredicto es exacto. Con ζ son INFINITAS y tenemos 649:")
	fmt.Printf("\n        raíces del impostor ....... 4      conocidas: 4    (el 100%%)\n")
	fmt.Printf("        ceros de ζ ................ ∞      conocidos: %d   (el 0%%)\n", len(ps))
	fmt.Println("\n   Por eso F259 midió el horizonte del laboratorio en γ ≈ 1658: más arriba, la")
	fmt.Println("   cola de lo que no vemos pesa más que todo lo que medimos.")
	fmt.Println("\n   📌 Y LO MISMO PASA CON DAVENPORT–HEILBRONN, que es el impostor grande de F259:")
	fmt.Println("   de él conocemos UN cero fuera de la línea, no todos. Así que NO podemos armarle")
	fmt.Println("   el precio tampoco. **Es inmune a esta prueba por la misma razón que ζ.**")
	fmt.Println("   Sería deshonesto decir que la dimensión 0 «caza impostores»: caza los que se")
	fmt.Println("   pueden CONTAR ENTEROS. Y ni ζ ni Davenport–Heilbronn se pueden contar enteros.")

	// ---- veredicto ----
	fmt.Println("\n════════ VEREDICTO ════════")
	fmt.Println("EL IMPOSTOR DE F229, ARMONIZADO EN LA DIMENSIÓN 0, SE CAE.")
	fmt.Printf("  · sus tres simetrías, perfectas ................. %.0e / %.0e / %.0e\n", peorFE, peorSch, peorSig)
	fmt.Printf("  · pero al disco, un par sale de la piel ......... |w| = %.9f y %.9f\n", cmplx.Abs(wA), cmplx.Abs(wB))
	fmt.Printf("  · con producto exactamente 1 (F225) ............. desvío %.1e\n", math.Abs(prod-1))
	fmt.Printf("  · y el precio de Li se hunde en ................. n = %d\n", primerNeg)
	fmt.Printf("  · mientras ζ aguanta n = 1..200 sin caer ........ %d negativos\n", negReal)
	fmt.Println("\nLO QUE ESTO DICE, Y ES BUENO:")
	fmt.Println("  ✅ La dimensión 0 VE lo que la geometría no ve. El cambiaformas no es un adorno:")
	fmt.Println("     convierte «está fuera de la línea» en «el precio explota», que es medible.")
	fmt.Println("  ✅ Y confirma que el camino de Li tiene DIENTES — no es un termómetro trabado.")
	fmt.Println("\n⚖️ LO QUE NO DICE, Y HAY QUE DECIRLO IGUAL DE FUERTE:")
	fmt.Println("  El impostor cayó porque tiene cuatro raíces y las conocemos todas. ζ tiene")
	fmt.Println("  infinitas. Davenport–Heilbronn también. La dimensión 0 condena a los que se")
	fmt.Println("  pueden contar enteros, y el problema del millón es exactamente el que NO se puede.")
	fmt.Println("\n¿El premio? Todavía no. Pero el disfraz de éste ya no le sirve.")

	escribirLamina(a, raices, wA, wB, prod, primerNeg, peor, peorN, peorFE, len(ps), negReal, lam)
}

func escribirLamina(a complex128, raices []complex128, wA, wB complex128, prod float64,
	primerNeg int, peor float64, peorN int, peorFE float64, nPerlas, negReal int, lam func(int) float64) {
	var b strings.Builder
	W, H := 1540.0, 1020.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="26" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🎭 EL IMPOSTOR, ARMONIZADO EN LA DIMENSIÓN 0</text>
<text x="%.0f" y="76" font-size="15" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">tiene todas las simetrías del libro y la geometría no lo ve — pero el precio de Li lo canta</text>
<text x="%.0f" y="112" font-size="19" text-anchor="middle" font-family="monospace" fill="#ffd98a">P(s) = (s−a)(s−ā)(s−(1−a))(s−(1−ā))   ·   a = 0.7 + 3i</text>
`, W, H, W, H, W/2, W/2, W/2)

	// the disk
	cx, cy, R := 350.0, 380.0, 175.0
	fmt.Fprintf(&b, `<rect x="40" y="140" width="620" height="500" rx="10" fill="#101f36" stroke="#26456e"/>
<text x="350" y="172" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ffd98a">AL DISCO: LAS CUATRO SE SALEN DE LA PIEL</text>
<circle cx="%.0f" cy="%.0f" r="%.0f" fill="none" stroke="#3d6fa8" stroke-width="2"/>
<text x="%.0f" y="%.0f" font-size="12" text-anchor="middle" font-family="monospace" fill="#4fa3d1">la piel |w| = 1</text>
`, cx, cy, R, cx, cy-R-10)
	for _, r := range raices {
		ww := 1 - 1/r
		m := cmplx.Abs(ww)
		px, py := cx+R*m*real(ww)/m, cy-R*m*imag(ww)/m
		col := "#ff5d73"
		fmt.Fprintf(&b, `<circle cx="%.2f" cy="%.2f" r="7" fill="%s"/>`, px, py, col)
		fmt.Fprintf(&b, `<text x="%.2f" y="%.2f" font-size="11" font-family="monospace" fill="#ff8fa0">|w|=%.4f</text>`,
			px+11, py-8, m)
	}
	fmt.Fprintf(&b, `<text x="350" y="592" font-size="14" text-anchor="middle" font-family="monospace" fill="#7ee0c0">|w₁|·|w₂| = %.12f  —  el norte × sur = 1 de F225</text>
<text x="350" y="616" font-size="13" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">si uno se hunde el otro se levanta, y el que sube crece como rⁿ</text>
`, prod)

	// the price
	fmt.Fprintf(&b, `<rect x="690" y="140" width="810" height="500" rx="10" fill="#161a3a" stroke="#5a4fa8"/>
<text x="1095" y="172" font-size="16" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">⚡ EL PRECIO DE LI, Y DÓNDE SE HUNDE</text>
<text x="1095" y="206" font-size="17" text-anchor="middle" font-family="monospace" fill="#dce8f7">λₙ = Σ [ 2 − 2·Re(wⁿ) ]</text>
`)
	// mini plot of lambda_n
	gx, gy, gw, gh := 730.0, 240.0, 730.0, 330.0
	maxN := 260
	lo, hi := 0.0, 0.0
	for n := 1; n <= maxN; n++ {
		v := lam(n)
		if v < lo {
			lo = v
		}
		if v > hi {
			hi = v
		}
	}
	if hi == lo {
		hi = lo + 1
	}
	y0 := gy + gh*(hi-0)/(hi-lo)
	fmt.Fprintf(&b, `<line x1="%.0f" y1="%.1f" x2="%.0f" y2="%.1f" stroke="#8fa8c7" stroke-width="1.4" stroke-dasharray="5,4"/>
<text x="%.0f" y="%.1f" font-size="11" font-family="monospace" fill="#8fa8c7">0</text>`,
		gx, y0, gx+gw, y0, gx-16, y0+4)
	var pts strings.Builder
	for n := 1; n <= maxN; n++ {
		x := gx + gw*float64(n-1)/float64(maxN-1)
		y := gy + gh*(hi-lam(n))/(hi-lo)
		fmt.Fprintf(&pts, "%.2f,%.2f ", x, y)
	}
	fmt.Fprintf(&b, `<polyline points="%s" fill="none" stroke="#ff5d73" stroke-width="2"/>`, pts.String())
	if primerNeg > 0 && primerNeg <= maxN {
		x := gx + gw*float64(primerNeg-1)/float64(maxN-1)
		fmt.Fprintf(&b, `<line x1="%.1f" y1="%.0f" x2="%.1f" y2="%.0f" stroke="#ffd166" stroke-width="1.6" stroke-dasharray="4,4"/>
<text x="%.1f" y="%.0f" font-size="12" text-anchor="middle" font-family="monospace" fill="#ffd166">n = %d</text>`,
			x, gy, x, gy+gh, x, gy-8, primerNeg)
	}
	fmt.Fprintf(&b, `<text x="1095" y="%.0f" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#ff8fa0">cae por primera vez en n = %d, y sigue de largo: %.4g en n = %d</text>`,
		gy+gh+34, primerNeg, peor, peorN)

	fmt.Fprintf(&b, `<rect x="40" y="660" width="1460" height="150" rx="10" fill="#0f2b22" stroke="#2f7f63"/>
<text x="770" y="692" font-size="18" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">LO QUE ESTO DICE, Y ES BUENO</text>
<text x="770" y="726" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">La dimensión 0 VE lo que la geometría no ve: convierte «está fuera de la línea» en «el precio explota», que es medible.</text>
<text x="770" y="752" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">Sus tres simetrías son perfectas a %.0e — y aun así el precio lo condena. El camino de Li tiene DIENTES.</text>
<text x="770" y="784" font-size="15" text-anchor="middle" font-family="monospace" fill="#7ee0c0">impostor: cae en n = %d    ·    ζ con %d perlas: %d negativos en n = 1..200</text>
`, peorFE, primerNeg, nPerlas, negReal)

	fmt.Fprintf(&b, `<rect x="40" y="828" width="1460" height="152" rx="10" fill="#33221c" stroke="#c0392b"/>
<text x="770" y="860" font-size="18" text-anchor="middle" font-family="Georgia" fill="#ffb27a">⚖️ Y LO QUE NO DICE, IGUAL DE FUERTE</text>
<text x="770" y="894" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#f3d9cf">El impostor cayó porque tiene CUATRO raíces y las conocemos a las cuatro. El precio de Li es una suma sobre TODAS.</text>
<text x="770" y="920" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#ffd98a">ζ tiene infinitas y tenemos %d. Davenport–Heilbronn también es inmune: de él conocemos UN cero fuera de la línea, no todos.</text>
<text x="770" y="950" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#f3d9cf">La dimensión 0 condena a los que se pueden CONTAR ENTEROS. El problema del millón es exactamente el que no se puede.</text>
</svg>
`, nPerlas)

	if err := os.WriteFile("el-impostor.svg", []byte(b.String()), 0o644); err != nil {
		fmt.Println("no pude escribir la lámina:", err)
		return
	}
	fmt.Println("\n🖼️  lámina escrita: el-impostor.svg")
}
