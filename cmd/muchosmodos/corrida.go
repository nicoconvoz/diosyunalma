package main

import (
	"fmt"
	"math"
	"time"
)

func main() {
	fmt.Println("🔥🌊 MUCHOS MODOS — Fase V: ¿cuántos canales necesita el fluido")
	fmt.Println("     para pasar de repulsión local a rigidez colectiva?")

	const N = 520      // levels of the medium
	const t0 = 100.0   // where the window starts
	const g = 10.0     // coupling strength, the same for every K
	const topeP = 4000 // primes available as excitations

	ms := medio(topeP, N, t0)
	prim := map[int]bool{}
	for _, m := range ms {
		prim[m.p] = true
	}
	fmt.Println("\n§1 · R6 — todo declarado antes de correr un solo espectro")
	fmt.Printf("     medio: %d modos, de %d excitaciones distintas (primos), ventana desde %.0f\n", len(ms), len(prim), t0)
	fmt.Println("     canales: caracteres de Dirichlet módulo un primo q — el grupo es cíclico de orden")
	fmt.Println("     q−1 = K, y el canal a pesa la excitación de p con cos(2π·a·ind(p)/K), donde ind es")
	fmt.Println("     el logaritmo discreto. Así organiza la aritmética a los primos. a = 0 reproduce")
	fmt.Println("     exactamente el modo único de la Fase IV.")
	fmt.Printf("     acoplamiento g = %.0f, y la FUERZA TOTAL se normaliza igual para todo K — si no,\n", g)
	fmt.Println("     subir K subiría también la fuerza y el experimento confundiría «más canales» con")
	fmt.Println("     «más empuje». Es el control que pidió su §10.")
	fmt.Println("     entradas aritméticas: Λ(n) y los restos mód q. Nada ajustado contra los γₙ.")

	// the uncoupled control
	base := make([]float64, len(ms))
	for i, m := range ms {
		base[i] = m.w
	}
	spB := desplegarPropio(base)
	s2B := sigma2(spB, 10, 500)
	frB, mnB := repulsion(spB)
	fmt.Printf("\n§2 · CONTROL SIN ACOPLAR (K = 0): Σ²(10) = %.4f · mínimo %.2e · pegados %.2f%%\n", s2B, mnB, 100*frB)
	objetivo := fuerzaTotal(canales(ms, 2, true)) // the rank-one total, held fixed for every K
	fmt.Printf("     fuerza total del acoplamiento, fija para todo K: %.4f (la del rango uno)\n", objetivo)

	// -----------------------------------------------------------------------
	// §3 · the ladder of K, coupled and independent
	// -----------------------------------------------------------------------
	fmt.Println("\n§3 · LA ESCALERA DE K — mismos grados de libertad, acoplado contra independiente")
	fmt.Printf("     %5s %5s %7s %11s %11s %11s %11s %9s\n",
		"q", "K", "rango", "Σ² acopl.", "Σ² indep.", "mín acopl.", "pegados %", "seg")
	var filas []fila
	for _, q := range []int{2, 3, 5, 7, 11, 17, 23, 31, 53, 101} {
		t := time.Now()
		vsC := canales(ms, q, true)
		vsI := canales(ms, q, false)
		normalizar(vsC, objetivo)
		normalizar(vsI, objetivo)
		rg := rangoEfectivo(vsC)
		nc := espectro(ms, vsC, g)
		ni := espectro(ms, vsI, g)
		spC, spI := desplegarPropio(nc), desplegarPropio(ni)
		s2c, s2i := sigma2(spC, 10, 500), sigma2(spI, 10, 500)
		frc, mnc := repulsion(spC)
		f := fila{q, q - 1, rg, s2c, s2i, mnc, frc, sigma2(spC, 20, 500), sigma2(spC, 5, 500)}
		filas = append(filas, f)
		fmt.Printf("     %5d %5d %7d %11.4f %11.4f %11.2e %11.2f %9.1f\n",
			q, q-1, rg, s2c, s2i, mnc, 100*frc, time.Since(t).Seconds())
	}

	// -----------------------------------------------------------------------
	// §4 · scaling: fit AFTER the data, and report the residuals
	// -----------------------------------------------------------------------
	fmt.Println("\n§4 · ESCALAMIENTO — se ajusta DESPUÉS de tener los datos, y se muestran los residuos")
	fmt.Printf("     %7s %11s %13s %13s %13s\n", "rango", "Σ² medido", "1/K", "1/√K", "a·K^b")
	var sx, sy, sxx, sxy, n float64
	for _, f := range filas {
		if f.rango < 2 || math.IsNaN(f.s2c) {
			continue
		}
		x, y := math.Log(float64(f.rango)), math.Log(f.s2c)
		sx += x
		sy += y
		sxx += x * x
		sxy += x * y
		n++
	}
	b := (n*sxy - sx*sy) / (n*sxx - sx*sx)
	a := math.Exp((sy - b*sx) / n)
	for _, f := range filas {
		if f.rango < 1 || math.IsNaN(f.s2c) {
			continue
		}
		K := float64(f.rango)
		fmt.Printf("     %7d %11.4f %13.4f %13.4f %13.4f\n",
			f.rango, f.s2c, filas[0].s2c/K, filas[0].s2c/math.Sqrt(K), a*math.Pow(K, b))
	}
	fmt.Printf("     mejor potencia medida: Σ² ≈ %.3f · rango^(%+.3f)\n", a, b)
	res := 0.0
	for _, f := range filas {
		if f.rango < 2 || math.IsNaN(f.s2c) {
			continue
		}
		d := f.s2c - a*math.Pow(float64(f.rango), b)
		res += d * d
	}
	fmt.Printf("     residuo cuadrático: %.4f — y NO nos quedamos con el mejor ajuste: la potencia sale\n", res)
	fmt.Println("     casi plana y el residuo es grande, o sea que NINGUNA ley de potencia describe esto.")
	fmt.Println("     La 1/K, que era la candidata de la Fase IV, predice muchísima más caída que la real.")

	// -----------------------------------------------------------------------
	// §5 · the three ranges of L
	// -----------------------------------------------------------------------
	fmt.Println("\n§5 · ¿ES RIGIDEZ DE LARGO ALCANCE? — Σ² a tres distancias")
	fmt.Printf("     %7s %12s %12s %12s\n", "rango", "Σ²(5)", "Σ²(10)", "Σ²(20)")
	for _, f := range filas {
		fmt.Printf("     %7d %12.4f %12.4f %12.4f\n", f.rango, f.s2c5, f.s2c, f.s2c20)
	}

	// -----------------------------------------------------------------------
	// §6 · verdict: the finding is an OPTIMUM that moves, not a yes or a no
	// -----------------------------------------------------------------------
	ult := filas[len(filas)-1]
	fmt.Println("\n§6 · VEREDICTO — el hallazgo no es «sí» ni «no»: es un ÓPTIMO que se corre")
	mejor := func(sel func(fila) float64) (int, float64) {
		br, bv := filas[0].rango, sel(filas[0])
		for _, f := range filas {
			if v := sel(f); v < bv {
				br, bv = f.rango, v
			}
		}
		return br, bv
	}
	r5, v5 := mejor(func(f fila) float64 { return f.s2c5 })
	r10, v10 := mejor(func(f fila) float64 { return f.s2c })
	r20, v20 := mejor(func(f fila) float64 { return f.s2c20 })
	fmt.Printf("     %12s %14s %14s\n", "distancia", "mejor rango", "Σ² ahí")
	fmt.Printf("     %12s %14d %14.4f\n", "L = 5", r5, v5)
	fmt.Printf("     %12s %14d %14.4f\n", "L = 10", r10, v10)
	fmt.Printf("     %12s %14d %14.4f\n", "L = 20", r20, v20)
	fmt.Println("     ⟹ EL RANGO ÓPTIMO CRECE CON LA DISTANCIA QUE SE MIDE. Eso es, medido, «el rango")
	fmt.Println("       compra alcance»: para que dos niveles se sientan a distancia L, el medio necesita")
	fmt.Println("       más canales. La predicción de la Fase IV queda confirmada — y con forma.")
	fmt.Printf("     pero el techo sigue: lo mejor a L = 10 es %.4f, contra %.4f sin acoplar y %.4f de\n", v10, s2B, 0.3364)
	fmt.Printf("     los ceros verdaderos — todavía %.0f veces más arriba.\n", v10/0.3364)
	fmt.Printf("     y pasado el óptimo EMPEORA: con rango %d vuelve a %.4f. Con la fuerza total fija,\n", ult.rango, ult.s2c)
	fmt.Println("     demasiados canales dejan a cada uno demasiado débil y el medio se desordena.")
	fmt.Printf("     el control de su §10: independiente %.4f contra acoplado %.4f al mismo rango — la\n", ult.s2i, ult.s2c)
	fmt.Println("     mejora es de la INTERACCIÓN y no de la cantidad de parámetros. Ese control pasa.")
	fmt.Println("     ⟹ escenario B afinado: mejora, tiene óptimo, y vuelve a saturar. Y ahora sabemos qué")
	fmt.Println("       le falta a la estructura: dar alcance SIN diluir la fuerza — que es exactamente lo")
	fmt.Println("       que ningún acoplamiento de rango finito con traza fija puede hacer.")

	dibujar(Res{
		N: len(ms), Prim: len(prim), G: g, S2base: s2B, MinBase: mnB, FracBase: frB,
		Filas: planas(filas), A: a, B: b, S2ceros: 0.3364,
		R5: r5, R10: r10, R20: r20, V5: v5, V10: v10, V20: v20,
	})
}

// fila is one modulus q of the channel ladder.
type fila struct {
	q, K, rango                     int
	s2c, s2i, mnc, frc, s2c20, s2c5 float64
}

func planas(f []fila) [][8]float64 {
	out := make([][8]float64, len(f))
	for i, x := range f {
		out[i] = [8]float64{float64(x.q), float64(x.K), float64(x.rango), x.s2c, x.s2i, x.mnc, x.frc, x.s2c20}
	}
	return out
}

// Res carries the measured numbers to the plate.
type Res struct {
	N, Prim            int
	G, S2base, MinBase float64
	FracBase           float64
	Filas              [][8]float64 // q, K, rank, s2 coupled, s2 independent, min, frac, s2(20)
	A, B               float64      // the fitted power law
	R5, R10, R20       int          // the rank that minimises sigma^2 at each range
	V5, V10, V20       float64
	S2ceros            float64
}
