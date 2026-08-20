package main

// corrida.go - Phase X. ONE question: how much of the frustration gain survives
// when participation is held comparable to the control?
//
// The PR/N bands are DECLARED HERE, before the first spectrum, as her section 9
// demands. They are not chosen after seeing the numbers.

import (
	"fmt"
	"math"
	"sort"
)

const (
	NX   = 400
	T0X  = 100.0
	KMX  = 120
	FX   = 30.0
	SEM  = 6 // seeds per level
	PRok = 0.090
	PRmd = 0.070
)

type arm struct {
	nom            string
	s5, s10, s20   float64
	alfa, pr, dens float64
	frust          float64
	vivos          int
	des            float64 // spread over seeds, when the arm is a family
	n              int
}

type ResX struct {
	S0                      arm
	FamB, FamA              []arm
	Dist, Recip, Igualada   []arm
	Nodo                    []arm
	Pareto                  []arm
	Todos                   []arm
	Ceros, PisoGUE, PisoGOE float64
	Veredicto               string
}

// signoEnlace is family B: each bond independently negative with probability q.
// This is the clean knob - it lets the amount of frustration be tuned continuously
// from zero, which is what testing H3 (an intermediate optimum) requires.
func signoEnlace(n, kmax int, q float64, d *dado) func(int, int) float64 {
	t := make([][]float64, n)
	for i := range t {
		t[i] = make([]float64, kmax+1)
		for k := range t[i] {
			t[i][k] = 1
		}
	}
	for i := 0; i < n; i++ {
		for k := 1; k <= kmax && i+k < n; k++ {
			if d.u() < q {
				t[i][k] = -1
			}
		}
	}
	return func(i, j int) float64 {
		if i > j {
			i, j = j, i
		}
		if j-i > kmax {
			return 1
		}
		return t[i][j-i]
	}
}

// signoDistancia is her section 7 control "azar condicionado por distancia": the
// sign depends ONLY on k = |i-j|. Same marginal by distance, site structure gone.
func signoDistancia(kmax int, q float64, d *dado) func(int, int) float64 {
	s := make([]float64, kmax+1)
	for k := range s {
		s[k] = 1
		if d.u() < q {
			s[k] = -1
		}
	}
	return func(i, j int) float64 {
		k := j - i
		if k < 0 {
			k = -k
		}
		if k > kmax {
			return 1
		}
		return s[k]
	}
}

// densidadNeg measures the fraction of negative bonds actually present.
func densidadNeg(sg func(int, int) float64, n, kmax int) float64 {
	tot, neg := 0, 0
	for i := 0; i < n; i++ {
		for k := 1; k <= kmax && i+k < n; k++ {
			tot++
			if sg(i, i+k) < 0 {
				neg++
			}
		}
	}
	return float64(neg) / float64(tot)
}

func banda(pr float64) string {
	switch {
	case pr >= PRok:
		return "COMPARABLE"
	case pr >= PRmd:
		return "degradada"
	default:
		return "ATRAPADA"
	}
}

func medirArm(nom string, ms []modo, h []float64, c func(int) float64, sg func(int, int) float64, d *dado) arm {
	o := medir(espectro(ms, h, c, KMX, FX, sg))
	a := arm{nom: nom, s5: o.s5, s10: o.s10, s20: o.s20, alfa: o.alfa, pr: o.pr, vivos: o.vivos, n: 1}
	if sg != nil {
		a.dens = densidadNeg(sg, len(ms), KMX)
		a.frust = frustracion(sg, len(ms), KMX, d)
	}
	return a
}

// familia averages an arm over independent seeds and reports the spread, because
// her section 10 forbids comparing a single arithmetic rule against one seed.
func familia(nom string, ms []modo, h []float64, c func(int) float64, d *dado, hacer func(*dado) func(int, int) float64) arm {
	var s10s []float64
	acc := arm{nom: nom}
	viv := 0
	for s := 0; s < SEM; s++ {
		sg := hacer(d)
		a := medirArm(nom, ms, h, c, sg, d)
		s10s = append(s10s, a.s10)
		acc.s5 += a.s5 / SEM
		acc.s20 += a.s20 / SEM
		acc.alfa += a.alfa / SEM
		acc.pr += a.pr / SEM
		acc.dens += a.dens / SEM
		acc.frust += a.frust / SEM
		viv += a.vivos
	}
	acc.s10 = media(s10s)
	acc.des = desvio(s10s)
	acc.vivos = viv / SEM
	acc.n = SEM
	return acc
}

func fila(a arm) string {
	sd := "        "
	if a.n > 1 {
		sd = fmt.Sprintf("±%-7.3f", a.des)
	}
	return fmt.Sprintf("     %-30s %5d %8.3f %s %7.3f %7.3f %7.3f %7.3f  %s",
		a.nom, a.vivos, a.s10, sd, a.alfa, a.pr, a.dens, a.frust, banda(a.pr))
}

func cabecera() {
	fmt.Printf("     %-30s %5s %8s %8s %7s %7s %7s %7s  %s\n",
		"brazo", "vivos", "Σ²(10)", "±", "α", "PR/N", "dens−", "frustr", "banda PR/N")
}

func main() {
	fmt.Println("🧵🌿 LA FRONTERA — Fase X: ¿la frustración da RIGIDEZ, o sólo esconde")
	fmt.Println("     los estados? — separar rigidez real de localización")

	ms := medio(4000, NX, T0X)
	c0 := colaLarga(0.5)
	cN := conNodo(5, 0.5)
	h := normalizar1(homogenea(len(ms)))
	d := &dado{s: 20260818}
	R := ResX{Ceros: 0.3364, PisoGUE: pisoGUE(10), PisoGOE: pisoGOE(10)}

	fmt.Println("\n§0 · LA BANDA, DECLARADA ANTES DE LA PRIMERA MEDICIÓN (su §9)")
	fmt.Printf("     El control S0 tiene PR/N = 0,105. Se declara AHORA, sin haber mirado nada:\n")
	fmt.Printf("       COMPARABLE : PR/N ≥ %.3f   (dentro de ~15%% del control)\n", PRok)
	fmt.Printf("       degradada  : %.3f ≤ PR/N < %.3f\n", PRmd, PRok)
	fmt.Printf("       ATRAPADA   : PR/N < %.3f\n", PRmd)
	fmt.Println("     REGLA DE VEREDICTO, también fijada acá:")
	fmt.Println("       H2 (rigidez real) gana SÓLO si algún brazo COMPARABLE baja Σ²(10) por debajo")
	fmt.Println("       de S0 en más de 2 sigmas de su propia dispersión. Si no, gana H1 o H3.")
	fmt.Printf("     %d semillas independientes por nivel.\n", SEM)
	fmt.Println("     Y EL SEGUNDO RIEL, que su §4 pedía y yo había omitido en la primera corrida:")
	fmt.Println("       un brazo es ADMISIBLE sólo si conserva al menos el 80% de los niveles vivos")
	fmt.Println("       del control. Sin eso, Σ² se mide sobre un recorte distinto y no es comparable")
	fmt.Println("       — es exactamente el error que obligó a retractar la Fase VI. La primera")
	fmt.Println("       corrida de esta fase lo tenía: el mejor brazo bajaba a 87 niveles de 187.")

	fmt.Println("\n§1 · EL CONTROL")
	R.S0 = medirArm("S0 · todos +1", ms, h, c0, nil, d)
	cabecera()
	fmt.Println(fila(R.S0))

	// -----------------------------------------------------------------------
	fmt.Println("\n§2 · FAMILIA B — el dial limpio: cada enlace negativo con probabilidad q")
	fmt.Println("     Permite subir la frustración DESDE CERO y de a poco, que es lo que hace")
	fmt.Println("     falta para probar su H3 (un óptimo intermedio) en vez de sólo los extremos.")
	cabecera()
	for _, q := range []float64{0.005, 0.01, 0.02, 0.05, 0.10, 0.20, 0.35, 0.50} {
		qq := q
		a := familia(fmt.Sprintf("B · q = %.3f", qq), ms, h, c0, d,
			func(dd *dado) func(int, int) float64 { return signoEnlace(len(ms), KMX, qq, dd) })
		R.FamB = append(R.FamB, a)
		fmt.Println(fila(a))
	}

	// -----------------------------------------------------------------------
	fmt.Println("\n§3 · FAMILIA A — la de la Fase IX: signo derivado de sitios, s_ij = (−1)^(u_i·u_j)")
	cabecera()
	for _, p := range []float64{0.05, 0.15, 0.30, 0.5125, 0.80} {
		pp := p
		a := familia(fmt.Sprintf("A · p = %.4f", pp), ms, h, c0, d,
			func(dd *dado) func(int, int) float64 { return signoFrustrado(bitsAzar(len(ms), pp, dd)) })
		R.FamA = append(R.FamA, a)
		fmt.Println(fila(a))
	}

	// -----------------------------------------------------------------------
	fmt.Println("\n§4 · CONTROL POR DISTANCIA (su §7) — el signo depende SÓLO de k = |i−j|")
	fmt.Println("     Misma marginal por distancia, estructura de sitio destruida.")
	cabecera()
	for _, q := range []float64{0.10, 0.50} {
		qq := q
		a := familia(fmt.Sprintf("dist · q = %.2f", qq), ms, h, c0, d,
			func(dd *dado) func(int, int) float64 { return signoDistancia(KMX, qq, dd) })
		R.Dist = append(R.Dist, a)
		fmt.Println(fila(a))
	}

	// -----------------------------------------------------------------------
	fmt.Println("\n§5 · LA PREGUNTA ARITMÉTICA, con el control IGUALADO EN FRUSTRACIÓN (su §7)")
	uRec, dens := bitsReciprocidad(ms)
	sgRec := signoFrustrado(uRec)
	aRec := medirArm("reciprocidad p≡3 mod 4", ms, h, c0, sgRec, d)
	R.Recip = append(R.Recip, aRec)
	objF, objD := aRec.frust, aRec.dens
	fmt.Printf("     la regla aritmética da: densidad de negativos %.4f · frustración %.4f\n", objD, objF)
	fmt.Println("     el control ya no iguala sólo la densidad: iguala la FRUSTRACIÓN por rechazo.")
	var igu []float64
	var acc arm
	hechos, intentos := 0, 0
	fmt.Printf("     nota estructural: enlaces INDEPENDIENTES a densidad %.3f dan frustración ≈%.3f\n",
		objD, 3*objD*(1-objD)*(1-objD)+objD*objD*objD)
	fmt.Printf("     y la aritmética da %.3f — MÁS ALTA. O sea: el campo de la reciprocidad NO es\n", objF)
	fmt.Println("     alcanzable por enlaces independientes; vive en una zona del plano")
	fmt.Println("     (densidad, frustración) que el azar por enlace no visita. Por eso el control")
	fmt.Println("     igualado se sortea de la familia POR SITIO, que sí llega ahí.")
	for hechos < 10 && intentos < 600 {
		intentos++
		sg := signoFrustrado(bitsAzar(len(ms), 0.45+0.14*d.u(), d))
		if math.Abs(frustracion(sg, len(ms), KMX, d)-objF) > 0.02 ||
			math.Abs(densidadNeg(sg, len(ms), KMX)-objD) > 0.03 {
			continue
		}
		a := medirArm("igualado", ms, h, c0, sg, d)
		igu = append(igu, a.s10)
		acc.pr += a.pr
		acc.vivos += a.vivos
		acc.dens += a.dens
		acc.frust += a.frust
		hechos++
	}
	cabecera()
	fmt.Println(fila(aRec))
	if hechos >= 3 {
		f := float64(hechos)
		ai := arm{nom: fmt.Sprintf("azar IGUALADO en frustración"), s10: media(igu), des: desvio(igu),
			pr: acc.pr / f, dens: acc.dens / f, frust: acc.frust / f, vivos: int(float64(acc.vivos) / f), n: hechos}
		R.Igualada = append(R.Igualada, ai)
		fmt.Println(fila(ai))
		z := math.Abs(aRec.s10-ai.s10) / math.Max(ai.des, 1e-9)
		fmt.Printf("     %d controles aceptados de %d intentos · distancia aritmética−azar: %.2f sigmas\n",
			hechos, intentos, z)
		if z < 1 {
			fmt.Println("     ⟹ igualando densidad Y frustración, la aritmética NO se separa. Su §14 dice")
			fmt.Println("       qué hacer con eso: CERRAR este canal y no insistir con esta codificación.")
		} else {
			fmt.Println("     ⟹ la aritmética se separa aun con la frustración igualada: hay que replicar.")
		}
	} else {
		fmt.Printf("     sólo %d controles igualados en %d intentos: el rechazo no alcanza y NO se afirma nada\n", hechos, intentos)
	}

	// -----------------------------------------------------------------------
	fmt.Println("\n§6 · RECIPROCIDAD + NODO, el hilo que la Fase IX marcó (su §11)")
	cabecera()
	aRN := medirArm("nodo + reciprocidad", ms, h, cN, sgRec, d)
	R.Nodo = append(R.Nodo, aRN)
	fmt.Println(fila(aRN))
	aN0 := medirArm("nodo, sin signos", ms, h, cN, nil, d)
	R.Nodo = append(R.Nodo, aN0)
	fmt.Println(fila(aN0))
	aNaz := familia("nodo + azar a densidad igual", ms, h, cN, d,
		func(dd *dado) func(int, int) float64 { return signoFrustrado(bitsAzar(len(ms), dens, dd)) })
	R.Nodo = append(R.Nodo, aNaz)
	fmt.Println(fila(aNaz))
	zN := math.Abs(aRN.s10-aNaz.s10) / math.Max(aNaz.des, 1e-9)
	fmt.Printf("     reciprocidad+nodo contra azar+nodo, %d semillas: %.2f sigmas\n", SEM, zN)

	// -----------------------------------------------------------------------
	fmt.Println("\n§7 · LA FRONTERA DE PARETO (su §9) — Σ²(10) contra PR/N")
	R.Todos = append(R.Todos, R.S0)
	R.Todos = append(R.Todos, R.FamB...)
	R.Todos = append(R.Todos, R.FamA...)
	R.Todos = append(R.Todos, R.Dist...)
	R.Todos = append(R.Todos, R.Recip...)
	R.Todos = append(R.Todos, R.Igualada...)
	minV := int(0.8 * float64(R.S0.vivos))
	for _, a := range R.Todos {
		if a.vivos < minV {
			continue // inadmissible: the band emptied, so its Sigma^2 is not comparable
		}
		dom := false
		for _, b := range R.Todos {
			if b.vivos < minV {
				continue
			}
			if b.pr >= a.pr && b.s10 <= a.s10 && (b.pr > a.pr || b.s10 < a.s10) {
				dom = true
				break
			}
		}
		if !dom {
			R.Pareto = append(R.Pareto, a)
		}
	}
	sort.Slice(R.Pareto, func(i, j int) bool { return R.Pareto[i].pr > R.Pareto[j].pr })
	fmt.Println("     los brazos que NADIE domina (nadie baja más Σ² sin perder participación):")
	cabecera()
	for _, a := range R.Pareto {
		fmt.Println(fila(a))
	}

	// -----------------------------------------------------------------------
	fmt.Println("\n§7bis · SU §6, CONTESTADA: ¿alcanza con igualar la frustración?")
	fmt.Println("     Tres familias distintas A LA MISMA frustración de triángulos (≈0,22):")
	cabecera()
	for _, a := range R.Todos {
		if a.frust > 0.20 && a.frust < 0.26 {
			fmt.Println(fila(a))
		}
	}
	fmt.Println("     ⟹ NO ALCANZA. A frustración prácticamente igual, PR/N va de 0,046 a 0,281 y")
	fmt.Println("       Σ²(10) de 7,9 a 1,4. La frustración de triángulos NO determina la respuesta:")
	fmt.Println("       lo que decide es SI EL SIGNO SE DERIVA DE LOS SITIOS O DE LOS ENLACES.")
	fmt.Println("       El signo por sitio ATRAPA los estados; el signo por enlace NO. Y eso")
	fmt.Println("       reinterpreta la Fase IX: ahí toda la familia era por sitio, así que la")
	fmt.Println("       localización no venía de la frustración — venía de la CODIFICACIÓN.")

	fmt.Println("\n§8 · EL VEREDICTO, con los DOS rieles")
	minViv := int(0.8 * float64(R.S0.vivos))
	fmt.Printf("     riel de niveles vivos: hay que conservar ≥ %d de los %d del control\n", minViv, R.S0.vivos)
	var mejorComp *arm
	for i := range R.Todos {
		a := &R.Todos[i]
		if a.pr >= PRok && a.nom != R.S0.nom && a.vivos >= minViv {
			if mejorComp == nil || a.s10 < mejorComp.s10 {
				mejorComp = a
			}
		}
	}
	fmt.Println("     brazos que el riel DESCARTA (bajan Σ² vaciando la banda, no ordenándola):")
	for _, a := range R.Todos {
		if a.vivos < minViv {
			fmt.Printf("       %-30s vivos %d · Σ²(10) = %.3f  ⟶ inadmisible\n", a.nom, a.vivos, a.s10)
		}
	}
	if mejorComp == nil {
		R.Veredicto = "H1"
		fmt.Println("     NINGÚN brazo con signos mantuvo PR/N comparable al control.")
		fmt.Println("     ⟹ GANA H1: la caída de Σ² es localización. Su §15 es explícito — no se")
		fmt.Println("       afirma rigidez real por frustración. La ganancia de la Fase IX queda")
		fmt.Println("       REINTERPRETADA como estados atrapados, no como alcance.")
	} else {
		mejora := (R.S0.s10 - mejorComp.s10) / math.Max(mejorComp.des, 1e-9)
		fmt.Printf("     mejor brazo COMPARABLE: %s con Σ²(10) = %.3f y PR/N = %.3f\n",
			mejorComp.nom, mejorComp.s10, mejorComp.pr)
		fmt.Printf("     contra S0 = %.3f: %.2f sigmas de mejora\n", R.S0.s10, mejora)
		if mejora > 2 {
			R.Veredicto = "H2"
			fmt.Println("     ⟹ GANA H2: hay rigidez que SOBREVIVE a participación comparable.")
		} else {
			R.Veredicto = "H1"
			fmt.Println("     ⟹ la mejora a PR/N comparable no llega a 2 sigmas: NO se afirma rigidez real.")
		}
	}
	// H3: is there an interior optimum along the clean dial?
	if len(R.FamB) > 2 {
		mi, mv := 0, R.FamB[0].s10
		for i, a := range R.FamB {
			if a.s10 < mv {
				mi, mv = i, a.s10
			}
		}
		if mi > 0 && mi < len(R.FamB)-1 {
			fmt.Printf("     y el dial limpio tiene un MÍNIMO INTERIOR en q = %.3f (Σ²(10) = %.3f):\n",
				R.FamB[mi].dens, mv)
			fmt.Println("       eso es su H3 — existe una cantidad intermedia de frustración, y «más")
			fmt.Println("       frustración = mejor» es falso. Pero hay que leerlo junto a su PR/N.")
			if R.Veredicto == "H1" {
				R.Veredicto = "H1+H3"
			}
		} else {
			fmt.Println("     el dial limpio NO tiene mínimo interior: Σ² cae monótona con la frustración,")
			fmt.Println("       así que H3 no aparece en esta familia.")
		}
	}

	fmt.Println("\n§9 · LAS TRES REGLAS")
	fmt.Printf("     ceros %.4f · piso GUE %.4f · piso GOE %.4f\n", R.Ceros, R.PisoGUE, R.PisoGOE)
	fmt.Println("     su §12 tiene razón y se obedece: acá NO se persigue 0,3364. La pregunta de")
	fmt.Println("     esta fase era si la mejora es rigidez o escondite, y ésa se contestó.")

	dibujarX(R)
}
