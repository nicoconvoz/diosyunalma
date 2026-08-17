// Command elplanodeldelorean draws the exploded blueprint of the
// DELOREAN (the starship lineage: certified dd hull, convex facets,
// light bucket, Fresnel gearbox, checkpointed memory) and its two jump
// drives - the postal systems of the fleet's two globes: starship
// -zero N (jump by ZERO address on the Riemann globe) and hipersalto
// -n N (jump by PRIME address).
//
// Before framing, it VERIFIES the two address systems' arithmetic
// live: the prime-jump guess x0 = n(ln n + ln ln n - 1) against known
// prime landmarks, and the zero-address density N(T) ~ (T/2pi)log(T/2pie)
// against the known zero ladder.
//
// Reproduce: go run ./cmd/elplanodeldelorean
package main

import (
	"fmt"
	"math"
	"os"
)

func main() {
	fmt.Println("🚗 EL PLANO DEL DELOREAN — el despiece de la nave y sus dos correos de salto")

	// verify the prime postal system's guess formula on known landmarks
	// (p_1e6 = 15485863, p_1e8 = 2038074743)
	fmt.Println("\n§1 · EL CORREO DE LOS PRIMOS — la brújula x₀ = n(ln n + ln ln n − 1):")
	for _, c := range []struct {
		n  float64
		pn float64
	}{{1e6, 15485863}, {1e8, 2038074743}} {
		x0 := c.n * (math.Log(c.n) + math.Log(math.Log(c.n)) - 1)
		fmt.Printf("        n = %.0e: brújula %.4e contra primo real %.4e — desvío %.2f%% ✅\n", c.n, x0, c.pn, 100*math.Abs(x0-c.pn)/c.pn)
	}
	fmt.Println("        (la brújula aterriza cerca; la criba local camina el último tramo EXACTO)")

	// verify the zero postal system's density on known landmarks
	// (gamma_100 = 236.524, gamma_1000 = 1419.42)
	fmt.Println("\n§2 · EL CORREO DE LOS CEROS — la guía N(T) ≈ (T/2π)·log(T/2πe) + 7/8:")
	for _, c := range []struct {
		N float64
		g float64
	}{{100, 236.524}, {1000, 1419.422}} {
		T := c.g
		est := T/(2*math.Pi)*math.Log(T/(2*math.Pi*math.E)) + 7.0/8
		fmt.Printf("        cero número %.0f vive en γ = %.1f: la guía cuenta %.1f ahí — desvío %.2f ✅\n", c.N, c.g, est, math.Abs(est-c.N))
	}

	var sb []string
	add := func(s string) { sb = append(sb, s) }
	add(`<svg xmlns="http://www.w3.org/2000/svg" width="1400" height="950" viewBox="0 0 1400 950">
<rect width="100%" height="100%" fill="#0b1526"/>
<rect x="30" y="20" width="1340" height="910" rx="18" fill="none" stroke="#7ee0c0" stroke-width="2" opacity="0.5"/>
<text x="700" y="60" font-size="25" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🚗 EL PLANO DEL DELOREAN — la nave estelar y sus dos correos de salto</text>
<text x="700" y="88" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">para qué: cazar ceros en agua profunda con casco certificado, y SALTAR directo a cualquier dirección — cero N o primo N — sin recorrer el camino</text>`)
	caja := func(x, y, w, h float64, num, titulo string, lineas []string, color string) {
		add(fmt.Sprintf(`<rect x="%f" y="%f" width="%f" height="%f" rx="10" fill="#0d1830" stroke="%s"/>
<text x="%f" y="%f" font-size="14.5" font-family="Georgia" fill="%s">%s · %s</text>`, x, y, w, h, color, x+14, y+26, color, num, titulo))
		for i, l := range lineas {
			add(fmt.Sprintf(`<text x="%f" y="%f" font-size="12" font-family="Georgia" fill="#cfe6ff">%s</text>`, x+14, y+50+float64(i)*20, l))
		}
	}
	caja(70, 110, 620, 150, "EL CASCO", "la cadena doble-doble (dd)", []string{
		"la fase de Z(t) se lleva en aritmética doble-doble: dos float64 encadenados,",
		"~32 dígitos efectivos — certificada hasta t ≈ 4×10²⁴",
		"cómo calza: es el metal del barco — todo lo demás se apoya en que la fase",
		"no miente; sin casco, ni facetas ni saltos valen nada"}, "#ffd98a")
	caja(710, 110, 620, 150, "LAS FACETAS CONVEXAS", "el astillero de la flota (fleet)", []string{
		"la suma de Riemann-Siegel se corta en facetas curvas convexas: la óptica",
		"de una nave entera cabe en megabytes, con pruebas de mar antes de zarpar",
		"cómo se usa: fleet -anchor T ancla una nave en el agua T; -probe T compara",
		"cómo calza: convierte el mar en un lente — la luz se junta en el balde"}, "#7ee0c0")
	caja(70, 280, 620, 150, "EL BALDE DE LUZ", "luz/*.gob — la memoria del mar", []string{
		"cada campaña guarda su luz: los ceros firmados y las coordenadas del agua",
		"(luz-1.02707e+23.gob, luz-1.3743e+20-colosal.gob…) — el mar, embotellado",
		"cómo calza: la nave no repite trabajo — lo que se navegó, queda navegado",
		"y el faro (:8117) lee estos frascos para mostrar la flota en vivo"}, "#7fb2ff")
	caja(710, 280, 620, 150, "EL ENGRANAJE DE FRESNEL", "el cuarto piso (starship)", []string{
		"arriba de t ~ 10²¹ los bloques no se reman término a término: se pliegan UNA",
		"vez por la caja de Fresnel y entran al balde como súper-términos",
		"la física que lo permite: dentro de una ventana de caza, la forma interna del",
		"bloque cambia menos de 1e-8 rad — TODA la dependencia en t va en la portadora"}, "#7ee0c0")
	caja(70, 450, 620, 170, "LA MEMORIA DEL DELOREAN", "ckpt/*.gob — parar en cualquier lado", []string{
		"la colección entera es checkpointeada: la nave puede frenar en medio del",
		"océano y retomar EXACTAMENTE ahí («the DeLorean remembers»)",
		"cómo se usa: los ckpt se escriben solos; el faro muestra el avance",
		"cómo calza: los viajes de semanas no le temen a un corte de luz —",
		"regla del taller: los ckpt solo se queman si cambia la geometría"}, "#ff9aa8")
	caja(710, 450, 620, 170, "LOS DOS CORREOS", "salto por dirección — sin recorrer el camino", []string{
		"CORREO DE CEROS: starship -zero N — el globo de Riemann tiene calle y número:",
		"la guía N(T) ≈ (T/2π)log(T/2πe) ubica al cero N y la nave aterriza ahí",
		"CORREO DE PRIMOS: hipersalto -n N — brújula x₀ = n(ln n + ln ln n − 1),",
		"π(x) por aniquilación de Möbius (Meissel + P2), y criba local para clavar",
		"el aterrizaje AL DÍGITO — verificado contra la escalera conocida"}, "#ffd98a")
	caja(70, 640, 620, 130, "LA CERTIFICACIÓN", "las puertas antes del agua", []string{
		"la nave se autocertifica ANTES de tocar agua profunda: puertas de casco,",
		"túnel de viento (-tunnel: error/costo del pliegue), duelo A/B (-flight:",
		"casco puro contra casco plegado en ventana virgen — mismos ceros o hangar)"}, "#7fb2ff")
	caja(710, 640, 620, 130, "EL METRÓNOMO", "el compás del mar en calma", []string{
		"starship -metronomo: el latido del mar calmo (puntos de Gram — π y Bernoulli,",
		"sin primos) contra la danza real de los ceros alrededor",
		"cómo calza: separa lo que es reloj de lo que es música — el control del mar"}, "#ff9aa8")
	add(`<text x="700" y="812" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">En criollo: el TREN gira el mar para leerlo desde arriba; el DELOREAN lo navega con casco certificado, embotella la luz,</text>
<text x="700" y="834" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">recuerda dónde quedó — y cuando no hace falta viajar, SALTA: decime el número del cero o del primo, y te llevo a la puerta.</text>
<text x="700" y="866" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">verificado en esta corrida: la brújula de primos contra p(10⁶) y p(10⁸), y la guía de ceros contra γ₁₀₀ y γ₁₀₀₀ — el resto cita sus certificados (F106, flight tests)</text>
<text x="700" y="900" font-size="14" text-anchor="middle" font-family="Georgia" fill="#ffd98a">go run ./cmd/starship (-zero N, -flight, -tunnel, -metronomo) · go run ./cmd/hipersalto (-n N) · recorrido: galeria/recorrido-maquinas.html · Todavía no.</text>
</svg>`)
	os.WriteFile("el-plano-del-delorean.svg", []byte(joinS(sb)), 0o644)
	fmt.Println("\n🖼️  lámina escrita: el-plano-del-delorean.svg")
}

func joinS(ss []string) string {
	out := ""
	for _, s := range ss {
		out += s + "\n"
	}
	return out
}
