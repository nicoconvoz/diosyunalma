package main

// espejo.go - the captain's next flash: apply the Midpoint Theorem's formula TO
// THE ZEROS instead of to prime pairs. Three transplants, measured raw:
//   A) the MIDPOINTS m_n = (g_n + g_{n+1})/2 of consecutive zeros: do they
//      still carry the primes' voice at T = log p?
//   B) the ANCHORS a_n = (g_n - 1)/2 (the theorem's coordinate, which halves
//      every gap): the algebra predicts their echo lives at DOUBLED periods -
//      the anchor map should move the primes' voice to T = 2 log p, the slots
//      of the prime SQUARES. Measured, not assumed.
//   C) TWIN ZEROS: pairs closer than 0.7 of the mean spacing against pairs
//      wider than 1.3 - do the centers of CLOSE pairs carry more prime voice
//      than the centers of far pairs?
// Exploration mode: nothing registered unless something significant appears.

import (
	"fmt"
	"math"
)

func ecoDe(xs []float64, T float64) float64 {
	s := 0.0
	for _, x := range xs {
		s += math.Cos(x * T)
	}
	return 2 * s / float64(len(xs))
}

// zMedia: z-score of the mean of -E over the given periods against 2000 random
// periods drawn from the same window, for an arbitrary level set.
func zMedia(xs []float64, Ts []float64, tMin, tMax float64, d *dado) (float64, float64) {
	var vals []float64
	for _, T := range Ts {
		vals = append(vals, -ecoDe(xs, T))
	}
	var res []float64
	for r := 0; r < 2000; r++ {
		res = append(res, -ecoDe(xs, tMin+(tMax-tMin)*d.u()))
	}
	se := desvio(res) / math.Sqrt(float64(len(Ts)))
	return (media(vals) - media(res)) / math.Max(se, 1e-12), media(vals)
}

func espejo(g []float64, d *dado) {
	fmt.Println("\n🪞📐 EL ESPEJO — el Teorema del Punto Medio aplicado A LOS CEROS")

	primos := []int{5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97}
	var Tp []float64
	for _, p := range primos {
		Tp = append(Tp, math.Log(float64(p)))
	}
	tMin, tMax := math.Log(5.0), math.Log(100.0)

	// the transplanted objects
	var mids, anclas []float64
	for i := 0; i+1 < len(g); i++ {
		mids = append(mids, (g[i]+g[i+1])/2)
	}
	for _, x := range g {
		anclas = append(anclas, (x-1)/2)
	}

	zZ, mZ := zMedia(g, Tp, tMin, tMax, d)
	zM, mM := zMedia(mids, Tp, tMin, tMax, d)
	fmt.Println("\n§A · LOS PUNTOS MEDIOS de ceros consecutivos — ¿conservan la voz de los primos?")
	fmt.Printf("   ceros crudos    : -E medio %+.4f → %+6.1fσ   (el control positivo)\n", mZ, zZ)
	fmt.Printf("   PUNTOS MEDIOS   : -E medio %+.4f → %+6.1fσ\n", mM, zM)

	fmt.Println("\n§B · LAS ANCLAS a = (γ−1)/2 — la coordenada que divide gaps por 2")
	fmt.Println("   la álgebra manda: cos(a·T) resuena donde cos(γ·T/2) — la voz debe MUDARSE")
	fmt.Println("   de T = log p a T = 2·log p = log p². Se mide en los dos lugares:")
	var T2 []float64
	for _, p := range primos {
		T2 = append(T2, 2*math.Log(float64(p)))
	}
	zA1, _ := zMedia(anclas, Tp, tMin, tMax, d)
	zA2, _ := zMedia(anclas, T2, 2*tMin, 2*tMax, d)
	fmt.Printf("   anclas en T = log p   : %+6.1fσ   (¿quedó algo en el lugar viejo?)\n", zA1)
	fmt.Printf("   anclas en T = 2·log p : %+6.1fσ   (el lugar que la coordenada predice)\n", zA2)

	fmt.Println("\n§C · CEROS GEMELOS — centros de pares APRETADOS contra pares ANCHOS")
	esp := media(diffs(g))
	var cerca, lejos []float64
	for i := 0; i+1 < len(g); i++ {
		gp := g[i+1] - g[i]
		m := (g[i] + g[i+1]) / 2
		if gp < 0.7*esp {
			cerca = append(cerca, m)
		} else if gp > 1.3*esp {
			lejos = append(lejos, m)
		}
	}
	zC, mC := zMedia(cerca, Tp, tMin, tMax, d)
	zL, mL := zMedia(lejos, Tp, tMin, tMax, d)
	fmt.Printf("   centros de gemelos (gap<0,7·medio, %d pares): -E %+.4f → %+6.1fσ\n", len(cerca), mC, zC)
	fmt.Printf("   centros de anchos  (gap>1,3·medio, %d pares): -E %+.4f → %+6.1fσ\n", len(lejos), mL, zL)

	// the decisive control: REBUILD each midpoint with a SHUFFLED gap from the
	// pool. If the +12/-8 split survives foreign gaps, it is the deterministic
	// phase factor cos(gap*T/2) - trigonometry of the coordinate, not new
	// arithmetic. If it dies, the pairing gap<->position carries information.
	fmt.Println("\n§E · EL CONTROL DECISIVO — los mismos centros con GAPS BARAJADOS")
	sp := diffs(g)
	bar := append([]float64(nil), sp...)
	for i := len(bar) - 1; i > 0; i-- {
		j := int(d.u() * float64(i+1))
		bar[i], bar[j] = bar[j], bar[i]
	}
	var cerB, lejB []float64
	for i := 0; i+1 < len(g); i++ {
		m := g[i] + bar[i]/2 // midpoint rebuilt with a foreign gap
		if bar[i] < 0.7*esp {
			cerB = append(cerB, m)
		} else if bar[i] > 1.3*esp {
			lejB = append(lejB, m)
		}
	}
	zCB, mCB := zMedia(cerB, Tp, tMin, tMax, d)
	zLB, mLB := zMedia(lejB, Tp, tMin, tMax, d)
	fmt.Printf("   gemelos BARAJADOS (%d): -E %+.4f → %+6.1fσ\n", len(cerB), mCB, zCB)
	fmt.Printf("   anchos  BARAJADOS (%d): -E %+.4f → %+6.1fσ\n", len(lejB), mLB, zLB)

	fmt.Println("\n§ LECTURA DEL ESPEJO")
	fmt.Printf("   voz conservada por los puntos medios: %.0f%% de la de los ceros crudos\n",
		100*zM/math.Max(zZ, 1e-9))
	if zA2 > 3 && zA1 < zA2/2 {
		fmt.Println("   ✓ LA MUDANZA SE CUMPLE: las anclas cantan en log p², como la coordenada")
		fmt.Println("     del teorema predice — la mitad de los gaps corre la voz una octava.")
	}
	if math.Abs(zC-zL) > 3 {
		fmt.Println("   ⚡ los centros de gemelos y de anchos NO cantan igual: los pares apretados")
		fmt.Println("     llevan otra cantidad de voz de los primos. ESO merece más ceros.")
	} else {
		fmt.Println("   · gemelos y anchos cantan parecido: el gap no cambia la voz del centro.")
	}
}

func diffs(g []float64) []float64 {
	var v []float64
	for i := 0; i+1 < len(g); i++ {
		v = append(v, g[i+1]-g[i])
	}
	return v
}

// espejoC replicates only the twin/wide split plus its shuffle control.
func espejoC(g []float64, d *dado) {
	primos := []int{5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97}
	var Tp []float64
	for _, p := range primos {
		Tp = append(Tp, math.Log(float64(p)))
	}
	tMin, tMax := math.Log(5.0), math.Log(100.0)
	esp := media(diffs(g))
	var cerca, lejos []float64
	for i := 0; i+1 < len(g); i++ {
		gp := g[i+1] - g[i]
		m := (g[i] + g[i+1]) / 2
		if gp < 0.7*esp {
			cerca = append(cerca, m)
		} else if gp > 1.3*esp {
			lejos = append(lejos, m)
		}
	}
	zC, mC := zMedia(cerca, Tp, tMin, tMax, d)
	zL, mL := zMedia(lejos, Tp, tMin, tMax, d)
	fmt.Printf("   gemelos (%d): -E %+.4f -> %+6.1f sigmas | anchos (%d): -E %+.4f -> %+6.1f sigmas\n",
		len(cerca), mC, zC, len(lejos), mL, zL)
	sp := diffs(g)
	for i := len(sp) - 1; i > 0; i-- {
		j := int(d.u() * float64(i+1))
		sp[i], sp[j] = sp[j], sp[i]
	}
	var cerB, lejB []float64
	for i := 0; i+1 < len(g); i++ {
		m := g[i] + sp[i]/2
		if sp[i] < 0.7*esp {
			cerB = append(cerB, m)
		} else if sp[i] > 1.3*esp {
			lejB = append(lejB, m)
		}
	}
	zCB, _ := zMedia(cerB, Tp, tMin, tMax, d)
	zLB, _ := zMedia(lejB, Tp, tMin, tMax, d)
	fmt.Printf("   control barajado: gemelos %+6.1f sigmas | anchos %+6.1f sigmas\n", zCB, zLB)
}
