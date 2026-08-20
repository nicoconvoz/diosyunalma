package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// lamina.go - the Phase XI plate: one honest verdict board. Every arm of the
// frozen catalogue emptied the band, so the plate's job is to show WHY the
// verdict is "inadmissible" rather than pretend there was a race.

func e11(s string) string {
	return strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;").Replace(s)
}

func dibujar11(f0, b1, b2, b3 ficha, minViv int) {
	var b strings.Builder
	W, H := 1180.0, 760.0
	t := func(x, y, s float64, c, a, tx string) {
		fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" font-size="%.1f" font-family="Georgia,serif" fill="%s" text-anchor="%s">%s</text>`+"\n", x, y, s, c, a, e11(tx))
	}
	mono := func(x, y, s float64, c, a, tx string) {
		fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" font-size="%.1f" font-family="Consolas,monospace" fill="%s" text-anchor="%s">%s</text>`+"\n", x, y, s, c, a, e11(tx))
	}
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">`+"\n", W, H, W, H)
	fmt.Fprintf(&b, `<rect width="%.0f" height="%.0f" fill="#0b1220"/>`+"\n", W, H)
	fmt.Fprintf(&b, `<rect x="14" y="14" width="%.0f" height="%.0f" rx="16" fill="none" stroke="#ff9aa8" stroke-width="2" opacity="0.55"/>`+"\n", W-28, H-28)
	t(W/2, 52, 22, "#dce8f7", "middle", "🧵🦇 LAS RELACIONES — la aritmética en la distancia, y el riel que la frena")
	t(W/2, 78, 12, "#ffd98a", "middle", "catálogo congelado ANTES de mirar: Λ, μ y λ de Liouville · tres capas: aritmética / pura / permutada · R6 absoluto")

	y := 130.0
	mono(80, y, 12.5, "#8fb4d9", "start", fmt.Sprintf("%-30s %6s %9s %7s %7s   %s", "brazo", "vivos", "Σ²(10)", "PR/N", "dens−", "veredicto"))
	y += 10
	filas := []struct {
		f    ficha
		col  string
		nota string
	}{
		{f0, "#ffd98a", "el control"},
		{b1, "#ff9aa8", "banda VACIADA → inadmisible"},
		{b2, "#ff9aa8", "banda VACIADA → inadmisible"},
		{b3, "#ff9aa8", "banda VACIADA → inadmisible"},
	}
	for _, r := range filas {
		y += 26
		mono(80, y, 12.5, r.col, "start", fmt.Sprintf("%-30s %6d %9.3f %7.3f %7.3f   %s",
			r.f.nom, r.f.vivos, r.f.s10, r.f.pr, r.f.dens, r.nota))
	}
	y += 40
	t(80, y, 13, "#ffd98a", "start", fmt.Sprintf("El riel de niveles vivos exige ≥ %d de %d. Las tres reglas aritméticas, con sus densidades", minViv, f0.vivos))
	y += 20
	t(80, y, 13, "#ffd98a", "start", "naturales (0,33–0,54 de enlaces negativos), dejan 98–100 vivos: miden un recorte, no el medio.")
	y += 34
	t(80, y, 12.5, "#dce8f7", "start", "Y a densidad igualada, B ≈ A ≈ C dentro del ruido: la organización aritmética exacta k→S(k)")
	y += 20
	t(80, y, 12.5, "#dce8f7", "start", "no aporta nada que su propia permutación no aporte. Es la QUINTA aplicación de la misma hoja.")

	y += 42
	t(80, y, 13.5, "#7fb2ff", "start", "§13 de la auditora — el borde q = 0,05, replicado con 12 semillas:")
	y += 22
	mono(100, y, 12.5, "#7fb2ff", "start", "vivos 146,8 ± 2,7 · pasa el riel 2 de 12 veces · Σ²(10) = 2,584 ± 0,436")
	y += 20
	t(100, y, 12.5, "#ff9aa8", "start", "⟹ ACCIDENTAL. El brazo del borde no es un punto de trabajo: es una moneda cargada en contra.")

	y += 42
	t(80, y, 13.5, "#7ee0c0", "start", "Lo que queda en pie, y no es poco:")
	y += 22
	t(100, y, 12.5, "#dce8f7", "start", "· el fenómeno de la Fase X (orden 3× con participación subiendo) vive en q ≤ 0,02 — densidades")
	y += 18
	t(100, y, 12.5, "#dce8f7", "start", "  CHICAS. Ninguna función aritmética clásica es naturalmente rala a escala 120: para entrar ahí,")
	y += 18
	t(100, y, 12.5, "#dce8f7", "start", "  la aritmética necesitaría una dilución declarada de antemano — catálogo nuevo, fase nueva.")
	y += 22
	t(100, y, 12.5, "#dce8f7", "start", "· el eco no se corrió: la regla lo prohibía sin brazo admisible. No se fuerza la conexión.")

	t(W/2, H-64, 12, "#8fb4d9", "middle", "cerrar una pista falsa también es progreso — la auditora, §20 · el catálogo se congeló antes de mirar y se publica entero")
	t(W/2, H-40, 12.5, "#ffd98a", "middle", "go run ./cmd/lasrelaciones · estructura cerrada no es hipótesis demostrada · Todavía no.")
	b.WriteString("</svg>\n")

	dir := filepath.Join("galeria", "laminas", "10-el-telar")
	os.MkdirAll(dir, 0o755)
	full := filepath.Join(dir, "las-relaciones.svg")
	os.WriteFile(full, []byte(b.String()), 0o644)
	fmt.Printf("\n🖼️  lámina escrita: %s\n", full)
}
