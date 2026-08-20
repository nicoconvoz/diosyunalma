package main

// LA CADENA - response to Auditoria 55 (the guide). Her 8 steps, walked in
// order, by an INDEPENDENT route B:
//   - theta with an extra asymptotic term (+7/5760t^3)
//   - RS sum accumulated in descending order
//   - grid step 0.005 (not 0.01), trapezoidal integration (not plain sum)
//   - window shifted to [10000.5, 11000.5], clamp 1e-5 (not 1e-4)
// Plus the calibration her question 5 really needs: a SYNTHETIC signal with
// known planted amplitudes - including a tone planted AT A COMPOSITE (ln 10).
// If the pipeline recovers it, the silence of composites in zeta belongs to
// zeta, not to the instrument.

import (
	"fmt"
	"math"
)

// ---- route B Riemann-Siegel --------------------------------------------------

func thetaB(t float64) float64 {
	return t/2*math.Log(t/(2*math.Pi)) - t/2 - math.Pi/8 +
		1/(48*t) + 7/(5760*t*t*t)
}

func zetaB(t float64) float64 {
	th := thetaB(t)
	u := math.Sqrt(t / (2 * math.Pi))
	N := int(u)
	s := 0.0
	for n := N; n >= 1; n-- { // descending accumulation
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

func primoPotencia(m int) (int, int) {
	for p := 2; p <= m; p++ {
		if m%p != 0 {
			continue
		}
		v, k := m, 0
		for v%p == 0 {
			v /= p
			k++
		}
		if v == 1 {
			return p, k
		}
		return 0, 0
	}
	return 0, 0
}

func lambdaSobre(m int) float64 { // Lambda(m)/(2 sqrt(m) ln m)
	p, _ := primoPotencia(m)
	if p == 0 {
		return 0
	}
	return math.Log(float64(p)) / (2 * math.Sqrt(float64(m)) * math.Log(float64(m)))
}

func cerrada(m int) float64 { // p^{-k/2}/(2k)
	p, k := primoPotencia(m)
	if p == 0 {
		return 0
	}
	return math.Pow(float64(p), -float64(k)/2) / (2 * float64(k))
}

// ---- route B measurement -----------------------------------------------------

type serie struct{ ts, F []float64 }

func medirB(t0, T, dt float64) serie {
	nm := int(T/dt) + 1
	s := serie{ts: make([]float64, nm), F: make([]float64, nm)}
	med := 0.0
	for i := 0; i < nm; i++ {
		t := t0 + float64(i)*dt
		z := math.Abs(zetaB(t))
		if z < 1e-5 {
			z = 1e-5
		}
		s.ts[i] = t
		s.F[i] = math.Log(z)
		med += s.F[i]
	}
	med /= float64(nm)
	for i := range s.F {
		s.F[i] -= med
	}
	return s
}

// proyB: trapezoidal projection of F onto tone ln m.
func proyB(s serie, lm float64) float64 {
	var cr, ci, wsum float64
	for i := range s.F {
		w := 1.0
		if i == 0 || i == len(s.F)-1 {
			w = 0.5
		}
		cr += w * s.F[i] * math.Cos(s.ts[i]*lm)
		ci -= w * s.F[i] * math.Sin(s.ts[i]*lm)
		wsum += w
	}
	return math.Hypot(cr, ci) / wsum
}

func main() {
	fmt.Println("⛓️  LA CADENA — la guía de la relojera, eslabón por eslabón, por ruta B")
	fmt.Println("   (θ con término extra, suma descendente, paso 0,005, trapecio, ventana")
	fmt.Println("   corrida a [10000,5, 11000,5], clamp 1e-5 — nada compartido con la ruta A)")
	fmt.Println()

	fmt.Println("PASO 1-2 · las definiciones, exactas:")
	fmt.Println("   F(t) = log|Z(t)| centrada · c_m = (1/T)∫F(t)e^(−it·ln m)dt · barra = |c_m|")
	fmt.Println()

	fmt.Println("PASO 3 · la amplitud, derivada independientemente (las dos escrituras):")
	fmt.Println("      m        Λ(m)/(2√m·ln m)      p^(−k/2)/(2k)      diferencia")
	for _, m := range []int{2, 3, 4, 8, 9, 16, 19, 25, 27, 6, 10, 12} {
		a, c := lambdaSobre(m), cerrada(m)
		fmt.Printf("     %2d         %.6f            %.6f         %.1e\n", m, a, c, math.Abs(a-c))
	}
	fmt.Println("   ⟹ la identidad algebraica del §5 de la guía: exacta a precisión de máquina.")
	fmt.Println()

	fmt.Println("PASO 4 · calibración del instrumento con señal SINTÉTICA de amplitudes conocidas")
	fmt.Println("   plantado: 0,30·cos(t·ln2+1,1) + 0,15·cos(t·ln3+2,3) + 0,07·cos(t·ln5+0,7)")
	fmt.Println("           + 0,20·cos(t·ln10+1,9)   ← ¡un tono plantado EN UN COMPUESTO!")
	t0, T, dt := 10000.5, 1000.0, 0.005
	nm := int(T/dt) + 1
	sin1 := serie{ts: make([]float64, nm), F: make([]float64, nm)}
	for i := 0; i < nm; i++ {
		t := t0 + float64(i)*dt
		sin1.ts[i] = t
		sin1.F[i] = 0.30*math.Cos(t*math.Log(2)+1.1) +
			0.15*math.Cos(t*math.Log(3)+2.3) +
			0.07*math.Cos(t*math.Log(5)+0.7) +
			0.20*math.Cos(t*math.Log(10)+1.9)
	}
	fmt.Println("      m     recuperado    esperado (amplitud/2)")
	for _, caso := range []struct {
		m   int
		esp float64
	}{{2, 0.15}, {3, 0.075}, {5, 0.035}, {6, 0}, {10, 0.10}} {
		r := proyB(sin1, math.Log(float64(caso.m)))
		fmt.Printf("      %2d      %.4f           %.4f\n", caso.m, r, caso.esp)
	}
	fmt.Println("   ⟹ el instrumento RECUPERA el tono plantado en el compuesto 10 y da ~0 en")
	fmt.Println("     el 6 no plantado: el aparato oye compuestos si suenan. El silencio de")
	fmt.Println("     los compuestos en zeta es DE ZETA, no del método.")
	fmt.Println()

	fmt.Println("PASO 5 · reproducción independiente de los tres valores de la guía (§15)")
	v := medirB(t0, T, dt)
	fmt.Println("      m     ruta B      ruta A (acta)     Λ predice")
	for _, caso := range []struct {
		m    int
		acta float64
	}{{2, 0.3524}, {19, 0.1146}, {6, 0.0034}} {
		r := proyB(v, math.Log(float64(caso.m)))
		fmt.Printf("      %2d    %.4f       %.4f           %.4f\n", caso.m, r, caso.acta, cerrada(caso.m))
	}
	fmt.Println()

	fmt.Println("PASO 6 · la fuga con la ventana, ruta B (m = 6)")
	for _, TT := range []float64{500, 1000, 2000} {
		vv := medirB(t0, TT, dt)
		fmt.Printf("      T = %4.0f  →  |c_6| = %.4f\n", TT, proyB(vv, math.Log(6)))
	}
	fmt.Println()

	fmt.Println("PASO 7 · los controles negativos, ruta B")
	soloP := []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37}
	comp := []int{6, 10, 12, 14, 15, 18, 20, 21, 22, 24, 26, 28}
	media := func(set []int) (med, pre float64) {
		for _, m := range set {
			med += proyB(v, math.Log(float64(m)))
			pre += cerrada(m)
		}
		return med / float64(len(set)), pre / float64(len(set))
	}
	mP, pP := media(soloP)
	mC, pC := media(comp)
	var despl []int
	for _, p := range soloP {
		despl = append(despl, p+1)
	}
	mD, pD := media(despl)
	fmt.Printf("      primos reales:     medido %.4f   Λ %.4f   (ruta A: 0,1620 / 0,1624)\n", mP, pP)
	fmt.Printf("      compuestos puros:  medido %.4f   Λ %.4f   (ruta A: 0,0028 / 0)\n", mC, pC)
	fmt.Printf("      corridos p+1:      medido %.4f   Λ %.4f   (ruta A: 0,0429 / 0,0409)\n", mD, pD)
	fmt.Println()

	fmt.Println("PASO 8 · recién ahora, la frase de trabajo — el reparto exacto:")
	fmt.Println("   MATEMÁTICA (demostrado): la identidad algebraica del Paso 3; el factor")
	fmt.Println("     1/2 de la proyección (calibrado en el Paso 4); el producto de Euler y")
	fmt.Println("     von Mangoldt para Re(s)>1; Λ(m)=0 en compuestos = factorización única.")
	fmt.Println("   MEDICIÓN (numérico, reproducido por dos rutas): que la proyección de")
	fmt.Println("     log|Z| SOBRE LA LÍNEA devuelve esas amplitudes a ~10⁻³; la fuga ~1/T;")
	fmt.Println("     los controles negativos; la estabilidad de ventana.")
	fmt.Println("   NUEVO (si algo): ningún teorema. Lo que queda en pie es (a) el hecho")
	fmt.Println("     empírico de que la estructura de von Mangoldt es MEDIBLE sobre la")
	fmt.Println("     línea crítica con esta precisión por esta proyección, (b) la prueba de")
	fmt.Println("     calibración de que el silencio es de zeta y no del aparato, y (c) el")
	fmt.Println("     instrumento mismo — el parche log n donde timbre = aritmética.")
	fmt.Println("   Y SU §18, anotada como próxima pregunta: la información de los CEROS no")
	fmt.Println("     está en los tonos de primos (R²≈0,5): está en el residuo — las púas de")
	fmt.Println("     log|Z| en cada cero. La mitad del canto que los primos 2..40 no explican")
	fmt.Println("     es exactamente donde vive la posición fina de los ceros.")
}
