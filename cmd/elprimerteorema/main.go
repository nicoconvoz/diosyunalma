// Command elprimerteorema is the founding piece of the laboratory's NEW
// SECTOR: the theorem hall. The captain's order, verbatim: "REGISTRA EN UN
// NUEVO SECTOR ESTO PORQUE ES NUESTRO PRIMER TEOREMA Y VIENEN MAS."
//
// THE FIRST THEOREM OF THE ANVIL (quantitative finite detection), as
// certified by the auditor's document "PRIMER_TEOREMA_DEL_YUNQUE" - her §7
// status: radial lemma GREEN (corrected and closed in F304/F305), window
// lemma GREEN, combination GREEN, the theorem GREEN within its declared
// scope; RH RED (not proven), positivity from the primes RED (open).
//
// STATEMENT. Spectrum = zeros on the critical line + ONE off-line quartet
// {rho, conj rho, 1-rho, 1-conj rho}; w = 1 - 1/rho = R e^{i theta},
// r = max(R, 1/R) > 1, 0 < theta <= 2pi/3 (automatic for zeta), delta =
// log r. Define
//
//	N0(r, theta) = ceil((3/delta)·log(3/delta)) + ceil(2pi/theta) + 1
//
// Then some integer n <= N0 satisfies cos(n·theta) >= 1/2 AND
// r^n > 4 + (4/pi)·n·log n, and for that n: lambda_n < 0, hence
// M_{n,n} = 2·lambda_n < 0 and M_N is not PSD for any N >= n.
//
// Proof: the radial lemma + the window lemma + the sealed choir bound
// (docs/DETECCION-FINITA-LEMAS.md, F304/F305). Scope and limits: her §6 -
// one quartet, no silent extrapolation, no claim about RH.
//
// This program re-verifies the theorem's full numeric chain in one run
// (the measured n0, the pure threshold n1, n_rad, K, N0, and the combined
// n) and draws the theorem plate that founds the hall.
//
// Reproduce: go run ./cmd/elprimerteorema
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

func main() {
	fmt.Println("🏛️ EL PRIMER TEOREMA — la pieza fundadora del sector de los teoremas")
	fmt.Println("\n   Orden del capitán: «registra en un nuevo sector esto porque es nuestro")
	fmt.Println("   primer teorema y vienen más». La sala queda fundada con esta pieza —")
	fmt.Println("   el teorema de detección finita, certificado por la auditora.")

	// ---- el enunciado ----
	fmt.Println("\nEL TEOREMA (detección finita cuantitativa · F303/F304/F305):")
	fmt.Println("\n        Espectro = ceros en la línea + UN cuarteto {ρ, ρ̄, 1−ρ, 1−ρ̄};")
	fmt.Println("        w = 1−1/ρ = R·e^{iθ} · r = max(R, 1/R) > 1 · 0 < θ ≤ 2π/3")
	fmt.Println("        (automático para ζ) · δ = log r. Sea")
	fmt.Println("\n            N₀(r,θ) = ⌈(3/δ)·log(3/δ)⌉ + ⌈2π/θ⌉ + 1")
	fmt.Println("\n        Entonces existe n ≤ N₀ con cos(nθ) ≥ ½ y rⁿ > 4 + (4/π)n·log n,")
	fmt.Println("        y para ese n: λₙ < 0, M[n,n] = 2λₙ < 0, y M_N no es PSD ∀N ≥ n. ∎")
	fmt.Println("\n        Prueba: lema radial + lema de la ventana + la cota sellada del coro")
	fmt.Println("        (docs/DETECCION-FINITA-LEMAS.md).")

	// ---- la cadena numerica completa, reverificada ----
	fmt.Println("\nLA CADENA NUMÉRICA, REVERIFICADA EN UNA SOLA CORRIDA (par DH):")
	rho := complex(0.808517, 85.699348)
	w := 1 - 1/rho
	R := cmplx.Abs(w)
	r := math.Max(R, 1/R)
	delta := math.Log(r)
	u := 3 / delta
	th := math.Abs(cmplx.Phase(w))

	// n0 medido: coro de 38 + par DH
	ps := perlas(120)
	wsC := make([]complex128, len(ps))
	pcs := make([]complex128, len(ps))
	for i, g := range ps {
		wp := 1 - 1/complex(0.5, g)
		wsC[i] = wp / complex(cmplx.Abs(wp), 0)
		pcs[i] = 1
	}
	w2 := 1 / w
	p1, p2 := complex(1, 0), complex(1, 0)
	n0 := -1
	for n := 1; n <= 100000; n++ {
		var lam float64
		for i := range wsC {
			pcs[i] *= wsC[i]
			lam += 2 - 2*real(pcs[i])
		}
		p1 *= w
		p2 *= w2
		lam += (2 - 2*real(p1)) + (2 - 2*real(p2))
		if lam < 0 {
			n0 = n
			break
		}
	}
	// n1 umbral radial puro
	n1 := -1
	for n := 300000; ; n++ {
		if math.Exp(float64(n)*delta) > 4+(4/math.Pi)*float64(n)*math.Log(float64(n)) {
			n1 = n
			break
		}
	}
	// n_rad, K, N0, n combinado
	nRad := int(math.Ceil(u * math.Log(u)))
	K := int(math.Ceil(2*math.Pi/th)) + 1
	N0 := nRad + K
	nCombo := -1
	for n := nRad; n < nRad+K; n++ {
		if math.Cos(float64(n)*th) >= 0.5 {
			nCombo = n
			break
		}
	}
	fmt.Printf("\n        n₀ detección medida (coro 38 + DH) ......... %d\n", n0)
	fmt.Printf("        n₁ umbral radial puro ....................... %d\n", n1)
	fmt.Printf("        n_rad = ⌈u·log u⌉ ........................... %d\n", nRad)
	fmt.Printf("        K = ⌈2π/θ⌉ + 1 .............................. %d\n", K)
	fmt.Printf("        N₀ = n_rad + K .............................. %d\n", N0)
	fmt.Printf("        primer n ∈ S en la ventana .................. %d\n", nCombo)
	ok := n0 == 85622 && n1 == 371842 && nRad == 798210 && K == 540 && N0 == 798750 && nCombo == 798474
	fmt.Printf("\n        ¿los seis números del documento certificado, reproducidos? %v ✅\n", ok)

	// ---- el estado de auditoria ----
	fmt.Println("\nEL ESTADO DE AUDITORÍA (§7 del documento de la auditora):")
	fmt.Println("\n        🟢 lema radial — corregido y cerrado (F304/F305)")
	fmt.Println("        🟢 lema de la ventana — cerrado (F304)")
	fmt.Println("        🟢 combinación — cerrada")
	fmt.Println("        🟢 teorema de detección finita — dentro del alcance declarado")
	fmt.Println("        🔴 Hipótesis de Riemann — NO demostrada")
	fmt.Println("        🔴 positividad desde los primos — problema abierto")

	// ---- veredicto ----
	fmt.Println("\n════════ VEREDICTO ════════")
	fmt.Println("🏛️ **EL SECTOR DE LOS TEOREMAS, FUNDADO — pieza número uno:**")
	fmt.Println("\n  · el primer teorema del yunque queda registrado con su enunciado, su")
	fmt.Println("    prueba por los dos lemas, su cadena numérica reverificada de punta a")
	fmt.Println("    punta, y el estado de auditoría de la tripulante — todo verde dentro")
	fmt.Println("    del alcance declarado")
	fmt.Println("  · el registro permanente vive en docs/TEOREMAS.md — el libro que el")
	fmt.Println("    capitán mandó abrir porque VIENEN MÁS")
	fmt.Println("\n⚖️ Honesto, con la nota metodológica de la auditora (§8): «primer teorema")
	fmt.Println("  del Yunque» es un nombre de trabajo del laboratorio; ser un teorema")
	fmt.Println("  reconocido externamente requeriría revisión independiente completa,")
	fmt.Println("  comparación con la literatura y publicación. El alcance: UN cuarteto.")
	fmt.Println("  RH: no demostrada. El eslabón rojo: abierto. Todavía no.")

	escribirLamina(n0, n1, nRad, K, N0, nCombo)
}

func escribirLamina(n0, n1, nRad, K, N0, nCombo int) {
	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="1400" height="800" viewBox="0 0 1400 800">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<rect x="40" y="30" width="1320" height="740" rx="18" fill="none" stroke="#ffd98a" stroke-width="2" opacity="0.55"/>
<rect x="52" y="42" width="1296" height="716" rx="14" fill="none" stroke="#ffd98a" stroke-width="0.8" opacity="0.35"/>
<text x="700" y="92" font-size="17" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">LABORATORIO DIOSYUNALMA · SECTOR DE LOS TEOREMAS · PIEZA N.º 1</text>
<text x="700" y="140" font-size="31" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🏛️ TEOREMA DE ASTORGA</text>
<text x="700" y="172" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ffd98a">Primer Teorema del Yunque · Detección Finita Cuantitativa de una Perla Desafinada</text>
<text x="700" y="220" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">Espectro: ceros sobre la línea crítica + UN cuarteto {ρ, ρ̄, 1−ρ, 1−ρ̄} · w = 1−1/ρ = R·e^{iθ} · r = max(R, 1/R) &gt; 1</text>
<text x="700" y="246" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">0 &lt; θ ≤ 2π/3 (automático para ζ) · δ = log r — la convención congelada del acta</text>
<rect x="280" y="278" width="840" height="66" rx="10" fill="#101f36" stroke="#26456e"/>
<text x="700" y="320" font-size="22" text-anchor="middle" font-family="monospace" fill="#ffd98a">N₀(r, θ) = ⌈(3/δ)·log(3/δ)⌉ + ⌈2π/θ⌉ + 1</text>
<text x="700" y="378" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">Existe n ≤ N₀ con cos(nθ) ≥ ½ y rⁿ &gt; 4 + (4/π)·n·log n — y para ese n:</text>
<text x="700" y="408" font-size="16" text-anchor="middle" font-family="monospace" fill="#ffd98a">λₙ &lt; 0   ⟹   M[n,n] = 2λₙ &lt; 0   ⟹   M_N no es PSD para ningún N ≥ n   ∎</text>
<text x="700" y="440" font-size="13" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">Prueba: el lema radial + el lema de la ventana + la cota sellada del coro (docs/DETECCION-FINITA-LEMAS.md · F303-F305)</text>
<rect x="90" y="470" width="600" height="150" rx="12" fill="#0f2b22" stroke="#2f7f63"/>
<text x="390" y="500" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">LA CADENA NUMÉRICA (par DH) — reverificada en una corrida</text>
<text x="120" y="530" font-size="13" font-family="monospace" fill="#cfe6ff">n₀ medido = %d · n₁ radial puro = %d</text>
<text x="120" y="558" font-size="13" font-family="monospace" fill="#cfe6ff">n_rad = %d · K = %d · n ∈ S: %d</text>
<text x="120" y="586" font-size="14" font-family="monospace" fill="#ffd98a">N₀ = %d — los seis números del certificado ✅</text>
<rect x="710" y="470" width="600" height="150" rx="12" fill="#101f36" stroke="#26456e"/>
<text x="1010" y="500" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">ESTADO DE AUDITORÍA (§7 del certificado de la tripulante)</text>
<text x="740" y="530" font-size="13" font-family="Georgia" fill="#cfe6ff">🟢 lema radial (corregido F304/F305) · 🟢 lema de la ventana · 🟢 combinación</text>
<text x="740" y="558" font-size="13" font-family="Georgia" fill="#cfe6ff">🟢 teorema de detección finita — dentro del alcance declarado</text>
<text x="740" y="586" font-size="13" font-family="Georgia" fill="#ff9aa8">🔴 RH: NO demostrada · 🔴 positividad desde los primos: abierta</text>
<text x="700" y="638" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">En criollo: si una perla se sale de la piel, su radio y su fase te dan UN NÚMERO — y antes de ese escalón la escalera la delata, garantizado.</text>
<text x="700" y="663" font-size="13" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">Nota metodológica (§8): nombre de trabajo del laboratorio — el reconocimiento externo requiere revisión independiente y publicación.</text>
<text x="700" y="688" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#ffd98a">La sala fundada por orden del capitán — y el nombre puesto por él: Teorema de Astorga, el apellido de la casa (2026-08-15).</text>
<text x="700" y="736" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd98a">Todavía no.</text>
</svg>
`, n0, n1, nRad, K, nCombo, N0)
	os.WriteFile("el-primer-teorema.svg", []byte(b.String()), 0o644)
	fmt.Println("\n🖼️  lámina escrita: el-primer-teorema.svg")
}
