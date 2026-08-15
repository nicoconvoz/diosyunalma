# El corolario del río de pozos — el flash del capitán, traducido

**Para la auditora · 2026-08-15 · la pregunta del capitán: «¿puede haber
un claro por donde pase el río, y debajo del río un pozo que podamos
predecir?» — respuesta: SÍ, y es un corolario de lo ya auditado, con
cero maquinaria nueva.**

## La historia que lo parió (el mapa de los tres teoremas)

El Claro (Astorga): el niño entra al bosque y el aparente caos tenía un
orden escondido. El Río (DYN): el agua cambia sin parar pero sigue
siendo el mismo río. El Pozo (Diosyunalma): observa, calcula y predice
la profundidad antes de llegar — y el pozo está. El flash une los tres:
un claro, con el río pasando, y pozos predichos debajo.

## La traducción rigurosa (corolario-candidato)

**Enunciado.** Bajo H0-H4 (las de DYN, sin tocar), existe una sucesión
INFINITA de citas n₁ < n₂ < n₃ < … en el claro [n_rad, ∞) tal que

    λ_{n_k} ≤ −g(n_k)   para todo k,   con g(n₁) < g(n₂) < g(n₃) < …

y g(n_k) → ∞: cada pozo se predice más hondo que el anterior, antes de
llegar, y la escalera se hunde sin fondo sobre la sucesión de citas.

**Prueba (esquema — todo lemas ya auditados):**
- *El río es infinito:* iterar el lema de agenda (D4) con T₁ = n_rad y
  T_{k+1} = n_k + 1 — cada iteración produce una cita nueva estrictamente
  mayor, con su propia cota de programación ≤ T_k + ⌈2πT_k⌉^m.
- *El pozo bajo cada cita:* D5 en cada n_k da λ ≤ −g(n_k).
- *Los pisos solo se hunden:* D3 (g creciente en el claro).
- *Sin fondo:* g(n) ≥ e^{nδ} − (4/π)n·log n − (2m+2) → ∞.

**Evidencia** (testigo m = 2, ventana de 3000 escalones tras n_rad):
435 citas, cero violaciones de λ ≤ −g, cero pisos no-crecientes; las
primeras doce a la vista en `go run ./cmd/elriodepozos`. Evidencia de un
testigo — la prueba es el esquema de arriba, y los cuantificadores
exactos del enunciado son material para su auditoría.

**Qué agrega sobre T2+T3:** DYN daba UNA ruptura ≤ N₀; Diosyunalma le
dio profundidad; el corolario dice que la ruptura no es un evento — es
un RÍO de rupturas, cada una con pozo predicho y más hondo. (λₙ < 0
para infinitos n ya se sabía por la vía del teorema de ruptura global;
lo nuevo acá es que las infinitas rupturas caen en citas PROGRAMABLES
con pisos explícitos crecientes — fecha, lugar y profundidad, las tres.)

**Estado: 🔵 corolario-candidato — sin declarar.** El sello, el criterio
y los cuantificadores son suyos. La regla de siempre preside. Todavía no.

---

*El niño ya sabe las tres cosas: dónde está el claro, que el río nunca
deja de pasar, y cuánto mide cada pozo antes de asomarse.* 🌊🕳️📐
