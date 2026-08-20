# El lema formal del cielo — el pez, demostrado

**Para la mesa de los tres · 2026-08-16 · responde el PEDIDO F329
entero: enunciado exacto, hipótesis mínimas, umbral, prueba completa de
la desigualdad y del límite, m = 1 / m ≥ 2 separados, tasa, caso sin
líder estricto, y separación demostrado/observado.** Nombre «cielo»:
provisional, no usado como objeto. Trinidad: intacta.

---

## LEMA (candidato formal — sin bautizar)

**Hipótesis mínimas** (menos que las esperadas — hallazgo de la
formalización): **H0** (rᵢ > 1), **H1** (m finito), **H4** (fondo sobre
la línea, cerrado bajo conjugación, con N_fondo(T) ≤ (T/2π)log T), y
**HL (solo para m ≥ 2)**: líder estricto r_L > r₂ := maxᵢ≠L rᵢ.
**Ni H2 ni H3 se usan** — declarado.

**Enunciado.** Sea δ_L = log r_L. Para todo entero n ≥ 3:

    m = 1:   |λₙ/r_Lⁿ + 2cos(nθ_L)| ≤ [(4/π)n·log n + 4 + 2r_L⁻ⁿ] / r_Lⁿ
    m ≥ 2:   |λₙ/r_Lⁿ + 2cos(nθ_L)| ≤ [(4/π)n·log n + (6m−2) + 2(m−1)r₂ⁿ + 2r_L⁻ⁿ] / r_Lⁿ

y en consecuencia   **λₙ/r_Lⁿ + 2cos(nθ_L) → 0**   (n → ∞). ∎(cand.)

*(r₂ no aparece en el caso m = 1: la convención de la Trinidad,
mantenida. El umbral exacto de la desigualdad es n ≥ 3 — lo exige solo
la cota del coro; el límite no necesita umbral.)*

## Prueba de la desigualdad (Tarea 1 — desde la formulación vigente)

**(P1) La descomposición exacta** (la de F1-F3, sin nada nuevo):

    λₙ = coroₙ + Σᵢ ℓᵢ(n),   ℓᵢ(n) = 4 − 2cos(nθᵢ)(rᵢⁿ + rᵢ⁻ⁿ)

**(P2) Separar el término líder, EXACTO** (sin aproximar):

    ℓ_L(n) = −2cos(nθ_L)·r_Lⁿ + 4 − 2cos(nθ_L)·r_L⁻ⁿ

Por lo tanto, como identidad:

    λₙ/r_Lⁿ + 2cos(nθ_L) = [ coroₙ + Σᵢ≠L ℓᵢ(n) + 4 − 2cos(nθ_L)·r_L⁻ⁿ ] / r_Lⁿ

**(P3) Acotar cada pieza del numerador** (valor absoluto, triángulo —
suma finita por H1):

    (C1) 0 ≤ coroₙ ≤ (4/π)n·log n, ∀n ≥ 3        [L5 bajo H4, sellada F299-301]
    (C2) [m ≥ 2] |ℓᵢ| ≤ 4 + 2(rᵢⁿ + rᵢ⁻ⁿ) ≤ 6 + 2r₂ⁿ, ∀i ≠ L
         [|cos| ≤ 1; rᵢ⁻ⁿ ≤ 1; rᵢ ≤ r₂]  ⟹  Σᵢ≠L |ℓᵢ| ≤ (m−1)(6 + 2r₂ⁿ)
    (C3) |4| = 4      (C4) |2cos(nθ_L)·r_L⁻ⁿ| ≤ 2r_L⁻ⁿ

Sumando: numerador ≤ (4/π)n·log n + (m−1)·6 + 4 + 2(m−1)r₂ⁿ + 2r_L⁻ⁿ,
y (m−1)·6 + 4 = 6m−2. Para m = 1, C2 es vacía y queda la versión sin
r₂. Dividir por r_Lⁿ da la desigualdad. ∎

## Prueba del límite (Tarea 2 — término por término, sin esconder nada)

Con δ_L > 0 (H0) y r_L = e^{δ_L}:

    (L-a) n·log n / r_Lⁿ → 0:  e^{nδ_L} ≥ (nδ_L)³/6  [tres términos de la
          serie de la exponencial] ⟹ n·log n/e^{nδ_L} ≤ 6·log n/(δ_L³·n²) → 0.
          Elemental, sin hipótesis extra.
    (L-b) (6m−2)/r_Lⁿ → 0:  constante sobre exponencial creciente.  [H0]
    (L-c) [m ≥ 2] 2(m−1)·(r₂/r_L)ⁿ = 2(m−1)·e^{−n(δ_L−δ₂)} → 0
          ⟺ δ_L > δ₂ ⟺ HL. **Único lugar donde entra líder estricto.**
    (L-d) 4/r_Lⁿ → 0.   (L-e) 2r_L⁻ⁿ/r_Lⁿ = 2/r_L²ⁿ → 0.   [H0]

Ninguna afirmación necesita condición no declarada. ∎

## Corolario 1 — la escala de despeje (Tarea/§14, extraída formalmente)

De la desigualdad: para m ≥ 2, el término 2(m−1)e^{−n(δ_L−δ₂)} domina
asintóticamente a los demás (todos O(poly·e^{−nδ_L}) y δ_L − δ₂ < δ_L),
así que

    limsup (1/n)·log |λₙ/r_Lⁿ + 2cos(nθ_L)| ≤ −(δ_L − δ₂)     [m ≥ 2]
    (y ≤ −δ_L a menos del factor polinomial, si m = 1)

**La constante de despeje ES la brecha del líder — la misma de n_comp:
consecuencia formal del lema, no coincidencia numérica.** ∎

## Corolario 2 — la jerarquía (sub-cielo y firmamento, formalizados)

Observación clave: **R(n) = λₙ − ℓ_L(n) = coroₙ + Σᵢ≠L ℓᵢ(n) es
EXACTAMENTE la λ de la configuración reducida** (mismo fondo, los m−1
cuartetos restantes). Entonces:

- **Sub-cielo:** si la configuración reducida tiene su propio líder
  estricto (radios todos distintos — hipótesis declarada), el MISMO
  lema aplicado a ella da R(n)/r₂ⁿ + 2cos(nθ₂) → 0. Por inducción, la
  jerarquía entera: cada capa pelada revela la onda de su líder.
- **Firmamento:** peladas las m capas, el residuo es coroₙ exacto, y
  para fondo FINITO de p pares, 0 ≤ coroₙ ≤ 4p (cada par aporta
  4sin²(nφ/2) ≤ 4): acotado, sin crecimiento. Para fondo infinito bajo
  H4 rige la cota (4/π)n·log n (L5). ∎
  *(La invarianza de la forma entre configuraciones — Tarea 7 — queda
  así demostrada: la forma −2cos(nθ_L) sale del lema para TODA
  configuración bajo sus hipótesis, no de la configuración concreta.)*

## Nota obligatoria — el caso SIN líder estricto (Tarea 8)

HL no es decorativa: **en el empate es FALSA la conclusión.** Si
r₂ = r_L con θ₂ ≠ θ_L (dos líderes empatados), la misma identidad P2 da

    λₙ/r_Lⁿ + 2cos(nθ_L) = −2cos(nθ₂) + o(1)

y −2cos(nθ₂) NO tiende a 0: por el lema de la ventana hay infinitos n
con |cos(nθ₂)| ≥ ½. El cociente converge entonces al par de ondas
−2cos(nθ_L) − 2cos(nθ₂) — exactamente los «dos cosenos a la misma
escala» del acta de pesca, ahora con prueba de que el candidato NO
sobrevive sin HL. La hipótesis es necesaria, demostrado.

## Ataque hostil restante (Tarea 8, completo)

- ¿Agenda escondida? NO: el lema cuantifica sobre TODO n ≥ 3; ninguna
  cita interviene.
- ¿Normalización lineal? SÍ: A(n) es un cociente, sin logaritmos.
- ¿Dominios? r_L > 1 (H0) ⟹ δ_L > 0 y r_L⁻ⁿ ≤ 1; para m ≥ 2,
  1 < r₂ < r_L ⟹ 0 < δ₂ < δ_L. Todos verificados.
- ¿Ventanas? La estabilidad medida es EVIDENCIA (abajo); la prueba es
  P1-P3 + L-a..L-e, y solo ella.

## Evidencia (separada de la prueba, jamás la sustituye)

La curva de despeje del acta de pesca clava la tasa del Corolario 1 a
cinco órdenes (2.4×10⁻⁵ vs 2.4×10⁻⁵ en n = 200k; 2.8×10⁻¹⁰ vs
2.8×10⁻¹⁰ en 400k; piso float64 desde 700k); |A| ≤ 2.0000 en toda
ventana profunda; sub-cielo a 1.5×10⁻¹¹; firmamento en [0.02, 114.0]
≤ 4·38 = 152 ✓; m = 3 con la misma onda a 1.1×10⁻¹⁰.
(`go run ./cmd/elcielo`.)

## Estado propuesto (criterio §15)

    🟢 SELLO (autoevaluación): la desigualdad está demostrada ∀n ≥ 3
    con hipótesis mínimas declaradas (H0+H1+H4, más HL solo si m ≥ 2 y
    solo en L-c), el límite cerrado término a término, la tasa escrita
    como corolario, los casos m = 1 / m ≥ 2 separados sin objetos
    fantasma, la jerarquía formalizada con su hipótesis (radios
    distintos), y la necesidad de HL demostrada por contraejemplo.
    El sello real y el nombre — si el pez lo merece — son de la mesa.
    Todavía no.

---

*Primero vimos el cielo; ahora está demostrado que no era una ilusión:
quitarle la altura al mundo deja una onda de amplitud 2, que se despeja
exactamente a la velocidad de la brecha del líder.* 🌌📐
