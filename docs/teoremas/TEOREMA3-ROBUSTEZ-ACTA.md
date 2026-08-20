# Acta del candidato a Teorema 3 — la cota de robustez

**Para la auditora · 2026-08-15 · respuesta a la Primera Misión para Doc
(§7 del PLAN TEOREMA 3): conservar los márgenes cuantitativos descartados
en las simplificaciones y derivar una cota explícita para −λₙ.**
Reglas cumplidas: H0-H4 sin tocar · ningún resultado externo nuevo ·
ninguna simulación como demostración universal · y el punto de pérdida
del margen, localizado exactamente (§1).

**Estado: 🔵 CANDIDATO — T3 no declarado. Nace solo si pasa el criterio
§10 del plan, y los cuantificadores los audita la relojera.**

---

## 1. Dónde se perdió el margen (localizado, como pedía la misión)

En el acta de DYN, la línea **(R7)** dice:

    (R7) g(n_rad) > 0:  por (R1), e^{n_rad·δ} ≥ u³; y por (R4)+(R5)+(R6):
         2m+2 + (4/π)·n_rad·log n_rad ≤ u·[1 + 3(log u)²] ≤ u³

**Ahí está el descarte.** (R1) no da u³: da mucho más. La cadena exacta:

    n_rad ≥ u·log u   [definición, antes del techo]
    n_rad·δ ≥ u·log u·δ = 3(m+1)·log u        [porque u·δ = 3(m+1)]
    e^{n_rad·δ} ≥ e^{3(m+1)·log u} = u^{3(m+1)}

R7 solo necesitaba positividad, así que degradó u^{3(m+1)} a u³ (válido
porque m ≥ 1 ⟹ 3(m+1) ≥ 6 > 3). El factor descartado es **u^{3m}
entero** — para el testigo m = 2, un factor de 7.6×10²⁹. La misión era
recuperarlo; recuperado queda.

## 2. La derivación, con cuantificadores exactos (D1-D6)

Sea una configuración bajo H0-H4 (idénticas al acta de DYN, sin tocar),
δ = log r_max, u = 3(m+1)/δ, n_rad = ⌈u·log u⌉, y
N₀ = n_rad + (2π·n_rad + 1)^m.

    (D1) [∀ configuración] e^{n_rad·δ} ≥ u^{3(m+1)}         [§1, desde R1 exacta]
    (D2) [∀ configuración] 2m+2 + (4/π)·n_rad·log n_rad ≤ u³ [R4+R5+R6, sin cambios]
    (D3) [∀n ≥ n_rad] g(n) := e^{nδ} − (4/π)·n·log n − (2m+2) es CRECIENTE
         [R8+R9+R10, sin cambios: g′ > 0 en n_rad y g″ > 0 después]
    (D4) [∃n cita] el lema de agenda (L2+L3, sin cambios) da una cita con
         ε = 1 en [n_rad, n_rad + n₁], y ese n cumple n ≤ N₀
    (D5) [en esa cita] L7 con ε = 1:  λₙ ≤ coroₙ + 2m+2 − r_maxⁿ
         y con la cota del coro (L5, bajo H4):  λₙ ≤ −g(n)
    (D6) encadenando:  −λₙ ≥ g(n) ≥[D3] g(n_rad) ≥[D1,D2] u^{3(m+1)} − u³

**ENUNCIADO CANDIDATO (Robustez DYN).** Bajo H0-H4, con

    Δ(r_max, m) = u³·(u^{3m} − 1) > 0,   u = 3(m+1)/log r_max ≥ 6

existe n ≤ N₀(r_max, m) tal que  **λₙ ≤ −Δ(r_max, m)**.  ∎(candidato)

Positividad: u ≥ 6 > 1 ⟹ u^{3m} > 1 ⟹ Δ > 0 (D-positividad trivial).
Cada eslabón D1-D5 es un lema YA AUDITADO del acta de DYN (R1, R4-R6,
R8-R10, L2-L3, L5, L7) usado sin modificación — la novedad es únicamente
NO tirar el margen en D1 y encadenar en D6.

## 3. Los tres estantes (evidencia, jamás demostración)

**Batería de la derivación** (experimento sobre la fórmula): g(n_rad) ≥
Δ > 0 en 50 casos (m = 1..10, δ ∈ {0.01, …, 1}): 0 violaciones.

**Testigo m = 2** (el de DYN, F320/F321): Δ = 4.34×10⁴⁴ contra la
−λ = 6.496×10⁴⁴ medida (y verificada a 50 dígitos) en la cita 1040809.
**El cociente 1.50 muestra que, en este testigo, la cota captura la escala exponencial real y no resulta meramente decorativa.** (El exceso 1.50 se reparte entre el
techo del ceiling, el coseno real > ½, y el coro real ≪ su cota.)

**Testigo NUEVO m = 3** (construido para esta misión): tres cuartetos
reales (DH, 0.7+45i, 0.75+62i) + el coro de las 38 perlas; δ = 9.875×10⁻⁵,
n_rad,3 = 1 422 703, **Δ(m=3) = 1.04×10⁶¹ calculado ANTES de la corrida**;
la primera cita triple después del umbral cayó en n = 1 423 112 y ahí
λ = −1.91×10⁶¹ ≤ −Δ ✅ (cociente 1.84). La fórmula predijo el piso de una
configuración que el laboratorio nunca había armado, y la realidad obedeció.

Reproducir: `go run ./cmd/larobustez`.

## 4. Lo que la cota dice — y regalo para la misión de Nico

DYN decía: la ruptura LLEGA (∃n ≤ N₀: λₙ < 0). La robustez dice CUÁNTO:
llega con profundidad al menos Δ, y Δ = u³(u^{3m}−1) **crece
exponencialmente en m** — más cuartetos desafinados no diluyen la
delación: la profundizan a ritmo exponencial.

Para la pregunta guía de Nico (¿relación estructural entre resonancia y
profundidad?): de L7 sale la relación exacta **paramétrica en ε** —
en una cita de calidad ε,

    −λₙ ≥ 2(1 − ε²/2)·r_maxⁿ − coroₙ − 2ε²(m−1) − 4

la profundidad crece cuando la cita es más fina (ε → 0 duplica el
coeficiente del golpe de 1 a 2) y el término 2ε²(m−1) — el «ruido de los
acompañantes» — se apaga como ε². La resonancia perfecta duplica la
profundidad del golpe del líder y silencia a los demás: relación
estructural, no estadística. Queda a disposición como punto de partida
de esa misión.

## 5. Declaraciones (reglas §7 del plan)

- **H0-H4:** intactas, ni una coma.
- **Resultados externos nuevos:** NINGUNO (la cadena usa solo lemas
  internos ya auditados; para ζ siguen siendo B1 y B2, sin cambios).
- **Simulaciones:** los tres estantes de §3 son evidencia; la única
  prueba es D1-D6, escrita arriba con sus cuantificadores.
- **¿Se pierde margen en algún paso?** El único descarte que queda es
  D2/R4 (constantes de techo 1.094·2.06·1.28 → 3, y la absorción
  u[1+3(log u)²] ≤ u³) y la degradación cos(‖nθ‖) ≥ ½ del golpe en ε = 1
  — ambos acotados por el cociente medido ~1.5-1.8: el margen restante
  es un factor O(1), no exponencial. La pérdida exponencial que había
  (u^{3m}) era exactamente la que esta misión recuperó.

## 6. Para el criterio §10

Enunciado preciso ✓ · hipótesis explícitas (las de DYN, sin tocar) ✓ ·
lemas demostrados (todos heredados y auditados; un solo encadenamiento
nuevo, D6) ✓ · cuantificadores escritos en D1-D6 — **auditoría de la
relojera pendiente** · falsación activa: la batería y el testigo m = 3
construido a ciegas no la rompieron — **el intento hostil formal es de
la auditora** · separación prueba/evidencia: §2 vs §3 ✓.

**T3 todavía no declarado. El nacimiento lo decide el criterio, no el
entusiasmo.** Todavía no.

---

*La misión era recuperar lo que el taller tiró a la basura al
simplificar. Adentro del tacho había un factor u^{3m} — y con él, la
respuesta a «¿cuánto se hunde?»: al menos u³·(u^{3m}−1), exponencial en
el número de perlas desafinadas.* 🛡️📐
