package main

// EL CERO QUE NO SE ALCANZA - the captain asks whether there is a relation
// between time, space, matter and the absolute zero of temperature that the
// universe can reach.
//
// There is, it is exact, and the answer has a twist: the universe does NOT
// approach absolute zero. It approaches a FLOOR set by dark energy, and that
// floor is about 2.6e-30 K - thirty orders of magnitude below today's cold,
// but not zero.
//
// And on the way there is a familiar face: the same zeta function this house
// has been chasing for four hundred findings is literally inside the formulas
// for the energy and the photon count of the cosmic background.
//
// Verdicts fixed BEFORE computing:
//   V1 (the relation asked for) in an expanding universe T is inversely
//      proportional to the scale factor: T*a = constant. So temperature IS a
//      clock - it names the age. Computed back to the CMB's release.
//   V2 (the third law, Nernst 1912) absolute zero cannot be reached in a
//      finite number of steps. Stated as the theorem it is, not measured.
//   V3 (THE FLOOR, and the real answer) with a cosmological constant the
//      universe settles into de Sitter space, which has its OWN temperature
//      T = hbar*H/(2*pi*k). Computed from the measured Hubble constant.
//   V4 (zeta in the sky) the blackbody energy density carries zeta(4) and the
//      photon number density carries zeta(3). Both computed from their series
//      and checked against the closed forms, so the connection is verified and
//      not merely asserted.
//   V5 (honest limits) all of this assumes LambdaCDM and a dark energy that is
//      a true constant. If it evolves, the floor moves. Said, not hidden.
// Nernst 1912, Gibbons-Hawking 1977, Planck 1900. Nothing here is ours.

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

const (
	hbar = 1.054571817e-34 // J s
	kB   = 1.380649e-23    // J/K
	c    = 2.99792458e8    // m/s
	G    = 6.67430e-11     // m3 kg-1 s-2
	// Hubble constant, Planck 2018: 67.4 km/s/Mpc
	H0Mpc = 67.4
	Mpc   = 3.0856775814913673e22 // m
	Tcmb  = 2.72548               // K, measured by COBE/FIRAS
	edadU = 13.797e9              // years
)

func main() {
	fmt.Println("🌡️  EL CERO QUE NO SE ALCANZA — tiempo, espacio, materia y la temperatura mínima")
	fmt.Println()
	fmt.Println("   La pregunta del capitán tiene respuesta exacta, y trae una sorpresa: el universo NO se")
	fmt.Println("   dirige al cero absoluto. Tiene un PISO, y ese piso no es cero.")
	fmt.Println()

	H0 := H0Mpc * 1000 / Mpc // s^-1
	fmt.Printf("   constantes usadas · H₀ = %.4e s⁻¹ (67,4 km/s/Mpc, Planck 2018) · T_fondo = %.5f K (COBE)\n\n", H0, Tcmb)

	// ---------------------------------------------------------------- V1
	fmt.Println("§1 · LA RELACIÓN QUE PEDISTE: la temperatura ES un reloj")
	fmt.Println()
	fmt.Println("   En un universo que se expande, la temperatura del fondo cae con el estiramiento del espacio:")
	fmt.Println()
	fmt.Println("        T × a  =  constante          (a es el factor de escala: cuánto se estiró el espacio)")
	fmt.Println()
	fmt.Println("   O sea que espacio y temperatura son la misma cosa vista de dos maneras, y como el espacio")
	fmt.Println("   se estira con el tiempo, la temperatura te DICE la edad. Es un termómetro-calendario:")
	fmt.Println()
	fmt.Println("        cuando el universo era…     el espacio medía     la temperatura era     y eso pasaba a los…")
	for _, z := range []float64{1100, 10, 1, 0.5, 0} {
		a := 1 / (1 + z)
		T := Tcmb / a
		// rough matter/Lambda age, enough for the ladder
		edad := edadU * math.Pow(a, 1.5)
		if z == 1100 {
			edad = 380000
			fmt.Printf("        el fondo se liberó           %.5f            %8.1f K            %.0f años\n", a, T, edad)
			continue
		}
		fmt.Printf("        %.1f veces más chico          %.5f            %8.3f K            %.2f mil millones de años\n",
			1+z, a, T, edad/1e9)
	}
	fmt.Println()
	fmt.Println("   ⚡ Ahí está la relación entre el tiempo, el espacio y la temperatura: son UN SOLO número,")
	fmt.Println("      y la materia lo siente como calor. Cuando el espacio duplica su tamaño, todo se enfría")
	fmt.Println("      a la mitad — sin que nadie saque calor de ningún lado.")
	fmt.Println()

	// ---------------------------------------------------------------- V2
	fmt.Println("§2 · el cero absoluto, y por qué no se toca")
	fmt.Println()
	fmt.Println("   La tercera ley de la termodinámica (Nernst, 1912) dice que NO se puede llegar al cero")
	fmt.Println("   absoluto en un número finito de pasos. No es una dificultad técnica: es un teorema.")
	fmt.Println("   Cada paso de enfriamiento se lleva una fracción de lo que queda, nunca el resto entero.")
	fmt.Println()
	fmt.Println("   ⚡ Y fijate que ésa es EXACTAMENTE la forma de la criba de F409: cada primo se lleva")
	fmt.Println("      una fracción de lo que queda y nunca termina. La misma estructura, en otro mundo.")
	fmt.Println()

	// ---------------------------------------------------------------- V3
	fmt.Println("§3 · EL PISO — y acá está la sorpresa")
	fmt.Println()
	fmt.Println("   Si la energía oscura es una constante cosmológica, el universo no se enfría para siempre:")
	fmt.Println("   se acomoda en un espacio de De Sitter, que tiene su PROPIA temperatura, la de Gibbons y")
	fmt.Println("   Hawking (1977). El horizonte irradia, y esa radiación pone un piso:")
	fmt.Println()
	fmt.Println("        T_piso  =  ħ · H / (2π · k_B)")
	fmt.Println()
	Tds := hbar * H0 / (2 * math.Pi * kB)
	fmt.Printf("        con H₀ = %.4e s⁻¹:\n\n", H0)
	fmt.Printf("        T_piso  =  %.6e K\n\n", Tds)
	fmt.Printf("   ⚡⚡ El universo hoy está a %.5f K y va a seguir enfriándose… hasta %.2e K, y AHÍ SE DETIENE.\n", Tcmb, Tds)
	fmt.Printf("      Eso es un factor de %.1e — treinta órdenes de magnitud — y después nada más.\n", Tcmb/Tds)
	fmt.Println()
	fmt.Println("   📌 O sea: la temperatura mínima que el universo puede alcanzar NO es el cero absoluto.")
	fmt.Println("      Es un número chiquito pero distinto de cero, y quien lo fija es la energía oscura —")
	fmt.Println("      o sea, la geometría del espacio mismo. Tiempo, espacio y temperatura, otra vez atados.")
	fmt.Println()

	// how long
	fmt.Println("        ¿cuánto falta? El enfriamiento pasa a ser exponencial en la fase de De Sitter:")
	tH := 1 / H0 / (3.156e7 * 1e9) // Hubble time in Gyr
	n := math.Log(Tcmb / Tds)
	fmt.Printf("        el tiempo de Hubble vale %.1f mil millones de años, y hacen falta ~%.0f de ésos\n", tH, n)
	fmt.Printf("        → del orden de %.0f mil millones de años para tocar el piso. El universo tiene %.1f.\n", n*tH, edadU/1e9)
	fmt.Println()

	// ---------------------------------------------------------------- V4
	fmt.Println("§4 · Y ACÁ APARECE UNA VIEJA CONOCIDA: ζ está en el cielo")
	fmt.Println()
	fmt.Println("   La misma función zeta que este laboratorio persigue hace 410 hallazgos está DENTRO de las")
	fmt.Println("   fórmulas del fondo cósmico. No como metáfora: como factor numérico.")
	fmt.Println()
	// zeta(3) and zeta(4) from their series
	z3, z4 := 0.0, 0.0
	for n := 1; n <= 20000000; n++ {
		f := float64(n)
		z3 += 1 / (f * f * f)
		z4 += 1 / (f * f * f * f)
	}
	fmt.Println("        magnitud del fondo cósmico          lleva adentro     valor (sumado acá)     forma cerrada")
	fmt.Printf("        densidad de ENERGÍA (Stefan)          ζ(4)              %.9f            π⁴/90 = %.9f\n", z4, math.Pow(math.Pi, 4)/90)
	fmt.Printf("        cuántos FOTONES hay por metro³        ζ(3)              %.9f            (Apéry, irracional)\n", z3)
	fmt.Println()
	// photon density
	nGamma := 2 * z3 / (math.Pi * math.Pi) * math.Pow(kB*Tcmb/(hbar*c), 3)
	u := 8 * math.Pow(math.Pi, 5) * math.Pow(kB*Tcmb, 4) / (15 * math.Pow(2*math.Pi*hbar, 3) * math.Pow(c, 3))
	fmt.Printf("   ⚡⚡ Con esos números y T = %.5f K:\n", Tcmb)
	fmt.Printf("        fotones del fondo por metro cúbico:  %.4e   (el valor publicado es ~4,1e8)\n", nGamma)
	fmt.Printf("        densidad de energía del fondo:       %.4e J/m³\n", u)
	fmt.Println()
	fmt.Println("   📌 Cada metro cúbico de espacio vacío tiene ahí, ahora mismo, cuatrocientos millones de")
	fmt.Println("      fotones del Big Bang — y para contarlos hace falta ζ(3), el número de Apéry.")
	fmt.Println()

	// ---------------------------------------------------------------- V5
	fmt.Println("§5 · el vínculo más hondo: la temperatura ES tiempo imaginario")
	fmt.Println()
	fmt.Println("   Y todavía hay una atadura más directa entre el tiempo y la temperatura, y es exacta:")
	fmt.Println()
	fmt.Println("        β  =  ħ / (k_B · T)")
	fmt.Println()
	fmt.Println("   En la física cuántica a temperatura T, el sistema se comporta como si el tiempo fuera")
	fmt.Println("   CIRCULAR, con una vuelta que dura exactamente β en tiempo imaginario. Más caliente,")
	fmt.Println("   vuelta más corta; más frío, vuelta más larga.")
	fmt.Println()
	fmt.Println("        temperatura         la vuelta del tiempo imaginario dura…")
	for _, T := range []float64{1e9, 300, Tcmb, 1e-9, Tds} {
		b := hbar / (kB * T)
		et := ""
		if T == Tcmb {
			et = "   ← el fondo cósmico"
		}
		if T == Tds {
			et = "   ← el piso de De Sitter"
		}
		fmt.Printf("        %12.3e K        %.4e segundos%s\n", T, b, et)
	}
	fmt.Println()
	fmt.Printf("   ⚡⚡ En el piso de De Sitter la vuelta dura %.2e segundos — comparable a la edad del universo\n", hbar/(kB*Tds))
	fmt.Println("      multiplicada por miles de millones. El frío absoluto sería una vuelta de duración INFINITA:")
	fmt.Println("      el tiempo dejaría de cerrarse sobre sí mismo. Por eso el cero no se alcanza — habría que")
	fmt.Println("      estirar un círculo hasta hacerlo recto, y eso no se hace en un número finito de pasos.")
	fmt.Println()

	fmt.Println("§6 · LA RESPUESTA, JUNTA")
	fmt.Println()
	fmt.Println("   ✅ Sí, hay relación, y es exacta en los cuatro:")
	fmt.Println("      · ESPACIO y TEMPERATURA son el mismo número: T × a = constante")
	fmt.Println("      · TIEMPO y TEMPERATURA también, porque el espacio se estira con el tiempo:")
	fmt.Println("        la temperatura del fondo te dice la edad del universo")
	fmt.Println("      · la MATERIA lo siente como calor, y el fondo la baña con 4×10⁸ fotones por metro cúbico")
	fmt.Println("      · y la temperatura ES una vuelta de tiempo imaginario: β = ħ/(k·T)")
	fmt.Println()
	fmt.Printf("   ⚡⚡⚡ Y el cero absoluto NO se alcanza — pero no solo por Nernst: es que hay un PISO en %.2e K,\n", Tds)
	fmt.Println("        puesto por la energía oscura. El universo se va a enfriar treinta órdenes de magnitud")
	fmt.Println("        y ahí se queda, tibio en su propio horizonte, para siempre.")
	fmt.Println()
	fmt.Println("   ⚖️ LO QUE HAY QUE DECLARAR: nada de esto es nuestro (Nernst 1912, Planck 1900, Gibbons y")
	fmt.Println("      Hawking 1977) y todo depende de que el modelo ΛCDM sea correcto y de que la energía")
	fmt.Println("      oscura sea una CONSTANTE de verdad. Si evoluciona, el piso se mueve. Y esto es física,")
	fmt.Println("      no matemática: se mide con telescopios y se puede refutar con una medición nueva.")
	fmt.Println()
	fmt.Println("   Estructura cerrada no es hipótesis demostrada. Todavía no.")

	dibujar(Tds, z3, z4, nGamma)
	fmt.Println("\n🖼️  lámina escrita: galeria/laminas/10-el-telar/el-cero-que-no-se-alcanza.svg")
}

func dibujar(Tds, z3, z4, nGamma float64) {
	var b strings.Builder
	W, H := 1060, 700
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`, W, H, W, H)
	b.WriteString(`<rect width="100%" height="100%" fill="#0d1017"/>`)
	txt := func(x, y float64, sz int, fill, anc, s string) {
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="%d" fill="%s" text-anchor="%s" font-family="Georgia,serif">%s</text>`, x, y, sz, fill, anc, s)
	}
	txt(530, 32, 20, "#e8e2d0", "middle", "EL CERO QUE NO SE ALCANZA — el universo tiene un piso, y no es cero")
	txt(530, 54, 12, "#8a8676", "middle", "tiempo, espacio, materia y temperatura: un solo número · F411")

	// the chain
	fmt.Fprintf(&b, `<rect x="20" y="76" width="1020" height="120" fill="#121a16" stroke="#2a5a40"/>`)
	txt(530, 102, 13, "#6ee7a0", "middle", "LA RELACIÓN, EXACTA")
	txt(280, 148, 22, "#f2c14e", "middle", "T × a = constante")
	txt(280, 174, 11, "#8a8676", "middle", "espacio y temperatura, el mismo número")
	txt(780, 148, 22, "#5ec8ff", "middle", "β = ħ / (k·T)")
	txt(780, 174, 11, "#8a8676", "middle", "la temperatura ES una vuelta de tiempo imaginario")

	// the cooling curve to the floor
	x0, y0, w, h := 20.0, 212.0, 620.0, 280.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" fill="#131826" stroke="#2a3040"/>`, x0, y0, w, h)
	txt(x0+w/2, y0+24, 13, "#e8e2d0", "middle", "el enfriamiento, y dónde se detiene")
	piso := y0 + h - 60
	fmt.Fprintf(&b, `<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#6ee7a0" stroke-dasharray="5,4" stroke-width="2"/>`,
		x0+50, piso, x0+w-30, piso)
	txt(x0+w-34, piso-8, 11, "#6ee7a0", "end", fmt.Sprintf("EL PISO: %.1e K", Tds))
	var pts []string
	for i := 0; i <= 60; i++ {
		f := float64(i) / 60
		px := x0 + 60 + f*(w-100)
		logT := math.Log10(2.725) - f*30.5
		if logT < math.Log10(Tds) {
			logT = math.Log10(Tds)
		}
		py := piso - (logT-math.Log10(Tds))/32*(h-110)
		pts = append(pts, fmt.Sprintf("%.1f,%.1f", px, py))
	}
	fmt.Fprintf(&b, `<polyline points="%s" fill="none" stroke="#f2c14e" stroke-width="2.5"/>`, strings.Join(pts, " "))
	txt(x0+70, y0+70, 12, "#f2c14e", "start", "hoy: 2,725 K")
	txt(x0+w/2, y0+h-24, 11, "#8a8676", "middle", "el tiempo →   (treinta órdenes de magnitud de enfriamiento, y después nada)")

	// zeta in the sky
	rx := 660.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="380" height="%.0f" fill="#131826" stroke="#2a3040"/>`, rx, y0, h)
	txt(rx+190, y0+24, 13, "#e8e2d0", "middle", "y ζ está en el cielo")
	txt(rx+24, y0+62, 11, "#c9c3b0", "start", "la energía del fondo lleva adentro:")
	txt(rx+190, y0+92, 20, "#5ec8ff", "middle", fmt.Sprintf("ζ(4) = %.6f", z4))
	txt(rx+190, y0+114, 10, "#8a8676", "middle", "= π⁴/90")
	txt(rx+24, y0+152, 11, "#c9c3b0", "start", "y contar los fotones necesita:")
	txt(rx+190, y0+182, 20, "#f2c14e", "middle", fmt.Sprintf("ζ(3) = %.6f", z3))
	txt(rx+190, y0+204, 10, "#8a8676", "middle", "la constante de Apéry, irracional")
	fmt.Fprintf(&b, `<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#2a3040"/>`, rx+24, y0+226, rx+356, y0+226)
	txt(rx+190, y0+252, 12, "#6ee7a0", "middle", fmt.Sprintf("%.1e fotones", nGamma))
	txt(rx+190, y0+272, 11, "#c9c3b0", "middle", "por metro cúbico, ahora mismo, acá")

	// the honest box
	fmt.Fprintf(&b, `<rect x="20" y="508" width="1020" height="176" fill="#1a1216" stroke="#5a2a30"/>`)
	txt(530, 534, 13, "#ff6b6b", "middle", "lo que hay que declarar")
	txt(60, 566, 11, "#c9c3b0", "start", "· nada de esto es nuestro: Nernst 1912 (tercera ley) · Planck 1900 · Gibbons y Hawking 1977 (el piso)")
	txt(60, 590, 11, "#c9c3b0", "start", "· y todo depende de que la energía oscura sea una CONSTANTE de verdad. Si evoluciona, el piso se mueve.")
	txt(60, 614, 11, "#c9c3b0", "start", "· esto es FÍSICA, no matemática: se mide con telescopios y una medición nueva puede refutarlo mañana.")
	txt(530, 650, 12, "#f2c14e", "middle", "Lo que sí es del laboratorio: haber encontrado a ζ esperándonos en la temperatura del cielo.")
	txt(530, 672, 10, "#8a8676", "middle", "Todavía no.")
	b.WriteString(`</svg>`)
	ruta := filepath.Join("galeria", "laminas", "10-el-telar", "el-cero-que-no-se-alcanza.svg")
	_ = os.MkdirAll(filepath.Dir(ruta), 0o755)
	_ = os.WriteFile(ruta, []byte(b.String()), 0o644)
}
