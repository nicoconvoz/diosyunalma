package main

import (
	"fmt"
	"math"
)

// audit20: first audit element, brought forward. If the deep-truncation picket
// is simply RECONSTRUCTING the actual zeros (the counting equation becomes the
// definition of the zeros as N grows), the convergence is a duality statement,
// not an independent model. Measure it: k-th model point vs k-th real zero.
func audit20() {
	fmt.Println("AUDITORIA (elemento 1, adelantado): distancia del modelo a los CEROS REALES")
	g := cerosPaso(4000, 0.02)
	for _, N := range []int{97, 997, 9973, 99991} {
		sp := picketPuro(4000, componentes(N))
		n := len(sp)
		if len(g) < n {
			n = len(g)
		}
		sum, mx := 0.0, 0.0
		cerca := 0
		j := 0
		for i := 0; i < n; i++ {
			// nearest real zero (both series are sorted): index-matching would
			// inherit any point-count drift, so match by proximity instead
			for j+1 < len(g) && math.Abs(g[j+1]-sp[i]) < math.Abs(g[j]-sp[i]) {
				j++
			}
			d := math.Abs(sp[i] - g[j])
			sum += d
			if d > mx {
				mx = d
			}
			if d < 0.1 {
				cerca++
			}
		}
		fmt.Printf("   N=%-6d  puntos %d/%d  |dist| media %.4f  max %.3f  a <0,1 del cero real: %d%%\n",
			N, len(sp), len(g), sum/float64(n), mx, 100*cerca/n)
	}
	fmt.Println("   (espaciado medio real ~1,0: distancia media 0,1 = un decimo de gap)")
}
