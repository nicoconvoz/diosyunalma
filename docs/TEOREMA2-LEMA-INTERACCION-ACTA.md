# El Lema de Interacción — el acta de los dos pedidos

**Documento para la auditora · 2026-08-15 · respuesta al §14 de su
auditoría «EL ESCUDO CAE»**: (1) la desigualdad del golpe, derivada línea
por línea; (2) la versión exacta de Dirichlet, el lema de agenda y la
N₀ explícita. Verificación numérica: `go run ./cmd/losdospedidos` (F311).
Convención: la de F304 (log natural, techos arriba), más m = número de
cuartetos, r_max = maxᵢ rᵢ > 1, δ = log r_max.

> Regla del sello vigente: nada de este documento es una demostración de
> RH; el alcance es una configuración FINITA de cuartetos sobre fondo en
> la línea.

---

## 1. Pedido 1 — la desigualdad del golpe, línea por línea

**Afirmación.** En una cita simultánea (‖n·θᵢ‖ < ε para todo i, con
0 < ε ≤ 1):

    Σᵢ ℓᵢ,n ≤ 2ε²(m−1) + 4 − 2(1 − ε²/2)·r_maxⁿ

y con la elección ε = 1:  Σᵢ ℓᵢ,n ≤ 2m + 2 − r_maxⁿ.

**Derivación, sin saltos.**

*(L1)* Para todo x real: 1 − cos x = 2·sin²(x/2) ≤ 2·(x/2)² = x²/2,
o sea **cos x ≥ 1 − x²/2**. (Una línea: |sin t| ≤ |t|.)

*(L2)* Si ‖n·θᵢ‖ < ε entonces cos(n·θᵢ) = cos(‖n·θᵢ‖) ≥ 1 − ε²/2 > 0
(el coseno solo depende de la distancia al origen módulo 2π; positivo
porque ε ≤ 1 < π/2... de hecho basta ε ≤ √2).

*(L3)* Rᵢⁿ + Rᵢ⁻ⁿ = rᵢⁿ + rᵢ⁻ⁿ (la expresión es simétrica bajo R → 1/R,
y rᵢ = max(Rᵢ, 1/Rᵢ)); además rᵢⁿ + rᵢ⁻ⁿ ≥ 2 (AM–GM) y rᵢⁿ + rᵢ⁻ⁿ ≥ rᵢⁿ
(el segundo sumando es positivo).

*(L4)* Como cos(n·θᵢ) ≥ 1 − ε²/2 > 0 y (Rᵢⁿ + Rᵢ⁻ⁿ) > 0:

    ℓᵢ,n = 4 − 2·cos(n·θᵢ)·(Rᵢⁿ + Rᵢ⁻ⁿ) ≤ 4 − 2(1 − ε²/2)·(rᵢⁿ + rᵢ⁻ⁿ)

*(L5)* Para cada i ≠ i_max, usar (L3) con la cota "≥ 2":

    ℓᵢ,n ≤ 4 − 4(1 − ε²/2) = 2ε²

*(L6)* Para i = i_max, usar (L3) con la cota "≥ r_maxⁿ":

    ℓ_max,n ≤ 4 − 2(1 − ε²/2)·r_maxⁿ

*(L7)* Sumar (m−1) términos de (L5) y uno de (L6):

    Σᵢ ℓᵢ,n ≤ 2ε²(m−1) + 4 − 2(1 − ε²/2)·r_maxⁿ   ∎

**Verificado:** en las 24 079 citas reales (ε = 1) del mejor escudo del
laboratorio en n ≤ 2×10⁵: **cero violaciones**; ídem con tres cuartetos
en las citas triples.

## 2. Pedido 2 — Dirichlet exacto, el lema de agenda y la N₀

**(2a) El teorema exacto que se usa** (aproximación simultánea de
Dirichlet; prueba: palomar en el toro):

    Para reales α₁,…,α_m y entero Q ≥ 2, existe n con 1 ≤ n ≤ Q^m
    y enteros p₁,…,p_m tales que |n·αᵢ − pᵢ| ≤ 1/Q para todo i.

*Prueba (bosquejo del palomar):* los Q^m + 1 puntos
({j·α₁}, …, {j·α_m}), j = 0..Q^m, caen en las Q^m cajas de lado 1/Q del
cubo unidad; dos comparten caja; su diferencia n = j − j' cumple lo
pedido. ∎

**(2b) Traducción angular:** con αᵢ = θᵢ/2π queda: para todo Q ≥ 2
existe n ≤ Q^m con **‖n·θᵢ‖ ≤ 2π/Q para todo i**. Infinitud: si todas
las θᵢ/2π son racionales con denominador común D, los múltiplos de D dan
retornos exactos; si alguna es irracional, tomar Q → ∞ produce
soluciones nuevas sin fin (el error baja de todo umbral). ∎

**(2c) El lema de agenda** (una cita más allá de cualquier objetivo T,
con ε = 1). *(Precisado 2026-08-15 a pedido de la auditoría de F311,
Pedido 1: las hipótesis mínimas y cada desigualdad con su renglón.)*

**Hipótesis mínima: T entero, T ≥ 1.** Nada más.

    Aplicar (2b) con Q = ⌈2πT⌉ ≥ 2: existe n₁ entero, 1 ≤ n₁ ≤ Q^m
    ≤ (2πT + 1)^m, con ‖n₁·θᵢ‖ ≤ 2π/Q para todo i. Sea J = ⌈T/n₁⌉
    y n = J·n₁. Las cuatro desigualdades, renglón por renglón:

    (i)   J ≤ T:        n₁ ≥ 1 ⟹ T/n₁ ≤ T ⟹ ⌈T/n₁⌉ ≤ ⌈T⌉ = T
                        (T es entero, y ⌈·⌉ es monótona).
    (ii)  n ≥ T:        J ≥ T/n₁ ⟹ n = J·n₁ ≥ T.
    (iii) n ≤ T + n₁:   J < T/n₁ + 1 ⟹ n = J·n₁ < T + n₁.
    (iv)  ‖n·θᵢ‖ ≤ 1:   la norma circular ‖x‖ = dist(x, 2πℤ) es
                        subaditiva SIN restricción (es la métrica del
                        cociente ℝ/2πℤ), luego ‖J·n₁·θᵢ‖ ≤ J·‖n₁·θᵢ‖
                        ≤ J·(2π/Q) ≤ T·2π/⌈2πT⌉ ≤ T·2π/(2πT) = 1.  ∎

Verificado en batería: 2000 pares (T, n₁) al azar — cero violaciones de
las cuatro desigualdades.

**(2d) El umbral radial con m perlas — derivación completa, sin
abreviaturas.** *(Reescrito 2026-08-15 a pedido de la nota final de la
auditoría: cada constante con su renglón. La hipótesis queda AJUSTADA a
δ ≤ 1 — más simple que el (m+1)/2 anterior, suficiente, y automática
para ζ; el viejo colchón «(m+1)³ ≥ 2m+2» queda superado por la cota más
limpia 2m+2 ≤ u_m — cambio declarado.)*

**Hipótesis: δ = log r_max ≤ 1** (o sea r_max ≤ e). Automática para los
ceros relevantes por el LEMA δ-ζ (§2d-bis): δ ≤ log √2 = 0.3466.

**LEMA RADIAL-m.** Con δ ∈ (0, 1], m ≥ 1, u = u_m = 3(m+1)/δ y
n_rad,m = ⌈u·log u⌉, se cumple

    r_maxⁿ > 2m + 2 + (4/π)·n·log n     para todo n ≥ n_rad,m.

**Derivación, renglón por renglón.** Sea n* = u·log u y
g(n) = e^{n·δ} − (2m+2) − (4/π)·n·log n.

    (R0) u ≥ 3(m+1) ≥ 6, y n* = u·log u ≥ 6·log 6 = 10.75.

    (R1) n_rad·δ ≥ n*·δ = 3(m+1)·log u ≥ 3·log u
         ⟹ e^{n_rad·δ} ≥ u³.

    (R2) n_rad ≤ n* + 1 ≤ (1 + 1/10.75)·n* ≤ 1.094·n*     [por (R0)]

    (R3) log n_rad ≤ log(1.094·u·log u)
                   = log 1.094 + log u + log log u
                   ≤ 0.09 + log u + log u                  [log log u ≤ log u, u ≥ e]
                   ≤ 2.06·log u                            [0.09 ≤ 0.06·log u pues log u ≥ log 6 = 1.79]

    (R4) (4/π)·n_rad·log n_rad ≤ 1.28·(1.094·u·log u)·(2.06·log u)
                                ≤ 2.89·u·(log u)² ≤ 3·u·(log u)².

    (R5) 2m + 2 ≤ u:  equivale a δ ≤ 3(m+1)/(2m+2) = 3/2,
         y la hipótesis da δ ≤ 1 < 3/2.  ⟹ (2m+2)/u ≤ 1.

    (R6) LEMITA: u² ≥ 3·(log u)² + 1 para u ≥ 6.
         *(Cerrado autocontenido 2026-08-15, con la ruta de presentación
         sugerida por la auditoría — adoptada con crédito.)*
         Sea h(u) = u² − 3·(log u)²  y  q(u) = u² − 3·log u.

         (R6a) q'(u) = 2u − 3/u > 0 para u ≥ 6:
               2u − 3/u ≥ 12 − 1/2 = 23/2 > 0.
         (R6b) q(6) = 36 − 3·log 6 = 36 − 5.375 = 30.62 > 0.
         (R6c) De (R6a) + (R6b): q(u) > 0 para todo u ≥ 6,
               o sea u² > 3·log u en todo el rango.
         (R6d) h'(u) = 2u − 6·log(u)/u = (2/u)·(u² − 3·log u)
               = (2/u)·q(u) > 0 para u ≥ 6:  h es creciente.
         (R6e) h(6) = 36 − 3·(log 6)² = 36 − 9.635 = 26.36 > 1.
         (R6f) De (R6d) + (R6e): h(u) ≥ h(6) > 1 para todo u ≥ 6,
               es decir u² > 3·(log u)² + 1. ∎

    (R7) g(n_rad) > 0:  por (R1), e^{n_rad·δ} ≥ u³; y por (R4)+(R5)+(R6):
         2m+2 + (4/π)n_rad·log n_rad ≤ u·[1 + 3(log u)²] ≤ u·u² = u³
         con desigualdad estricta en (R6) para u ≥ 6.  ✓

    (R8) g'(n) = δ·e^{n·δ} − (4/π)(log n + 1) > 0 en n = n_rad.
         *(Cerrado autocontenido 2026-08-15, con la ruta sugerida por la
         auditoría de la servilleta — adoptada con crédito: la frase
         anterior «la brecha solo se agranda» era plausible, no
         demostración.)*
         Por un lado δ·e^{n_rad·δ} ≥ δ·u³ = 3(m+1)·u² ≥ 6·u²  [m ≥ 1];
         por el otro (4/π)(log n_rad + 1) ≤ 2.64·log u + 1.28  [por R3].
         Basta entonces:  6·u² > 2.64·log u + 1.28  para u ≥ 6.
         Sea H(u) = 6u² − 2.64·log u − 1.28.

         (R8a) H'(u) = 12u − 2.64/u ≥ 72 − 0.44 = 71.56 > 0 para u ≥ 6.
         (R8b) H(6) = 216 − 2.64·log 6 − 1.28 = 216 − 4.73 − 1.28
               = 209.99 > 0.
         (R8c) De (R8a) + (R8b): H(u) > 0 para todo u ≥ 6.
         (R8d) Luego δ·e^{n_rad·δ} ≥ 6u² > 2.64·log u + 1.28
               ≥ (4/π)(log n_rad + 1), o sea g'(n_rad) > 0. ∎

    (R9) g''(n) = δ²·e^{n·δ} − (4/π)/n > 0 en n ≥ n_rad:
         δ²·e^{n·δ} ≥ δ²·u³ = 9(m+1)²·u ≥ 9·4·6 = 216 ≫ (4/π)/n.  ✓

    (R10) g'' > 0 ⟹ g' creciente ⟹ g' > 0 para todo n ≥ n_rad
          ⟹ g creciente ⟹ g(n) ≥ g(n_rad) > 0 para todo n ≥ n_rad. ∎

**Verificado en grilla:** m = 1..10 × δ ∈ {0.01, 0.1, 0.3466, 0.7, 1.0}
— las diez líneas (R1)–(R10), cero violaciones en los 50 casos.

**(2d-bis) El fundamento de |Im ρ| ≥ 1** *(Precisado 2026-08-15 tras la
nota «los aderezos» de la auditoría: la versión anterior deducía de más —
el argumento de η solo excluye ceros en el segmento REAL (0,1), no ceros
complejos con 0 < |Im ρ| < 1, y «cero real» era un desliz terminológico.
Crédito a la auditora. La estructura corregida sigue exactamente sus dos
opciones del §2.)*

**PRIMERO — la hipótesis, explícita y sin deducción pretendida:**

    HIPÓTESIS DEL TEOREMA: todo cuarteto de la configuración cumple
    |Im ρᵢ| ≥ 1.

Es una restricción declarada del conjunto considerado — parte del
alcance del teorema, exigida siempre, también para configuraciones
sintéticas. El teorema NO afirma que esta hipótesis se deduzca de nada.

**SEGUNDO — por qué los ceros de ζ la satisfacen (justificación
independiente, con su naturaleza declarada):**

    (i) Sobre el segmento real (0,1) no hay ceros — elemental:
        η(σ) = Σ (−1)^{n−1} n^{−σ} > 0 (alternada decreciente) y
        1 − 2^{1−σ} < 0 en (0,1) ⟹ ζ(σ) < 0 ≠ 0. ∎
        [Esto cubre SOLO Im ρ = 0; se declara como tal.]

    (ii) Para 0 < |Im ρ| < 1, la ausencia de ceros es un TEOREMA
        ESTABLECIDO POR CONTEO RIGUROSO: el principio del argumento
        aplicado a ξ sobre el borde de la región (el método que Backlund
        1914 volvió riguroso) da N(T) exacto para alturas bajas, y el
        resultado clásico es N(14) = 0 — no hay ceros no triviales con
        0 < |Im ρ| < 14.13. Cómputos certificados: Gram (1903, primeros
        ceros), Backlund (1914, conteo riguroso), van de Lune–te
        Riele–Winter (1986, 1.5×10⁹ ceros), Platt–Trudgian (2021,
        verificación hasta 3×10¹²). Se cita como resultado computacional
        riguroso de la literatura — NO como observación del primer cero
        conocido.

    (iii) Nota de honestidad instrumental: el motor propio del
        laboratorio (perlas(), bisección de Riemann–Siegel) re-encuentra
        γ₁ = 14.134725 en cada corrida, pero eso es CORROBORACIÓN, no la
        fuente rigurosa: un barrido de cambios de signo de Z(t) no
        excluye por sí solo ceros fuera de la línea ni pares perdidos.
        La fuente rigurosa es (ii).

Luego **todo cero NO TRIVIAL de ζ** cumple |Im ρ| ≥ 14 > 1 [(i) + (ii)],
y con ello r ≤ √2 (LEMA δ-ζ) y δ ≤ 0.3466 ≤ 1. ∎

    LEMA δ-ζ (F312): |Im ρ| ≥ 1 ⟹ r ≤ √2.
    Prueba: |w|² = ((β−1)² + γ²)/(β² + γ²) con 0 < β < 1, |γ| ≥ 1;
    numerador y denominador están ambos entre γ² y 1 + γ², de modo que
    |w|² ∈ [1/2, 2] y r = max(|w|, 1/|w|) ≤ √2. ∎

**(2e) LA N₀ EXPLÍCITA.** Tomando T = n_rad,m en el lema de agenda:

    N₀(r_max, m) = n_rad,m + (2π·n_rad,m + 1)^m

y en la cita n ≤ N₀ así garantizada: λₙ ≤ [2m + 2 − r_maxⁿ] +
(4/π)·n·log n < 0 por (2d). ∎

**Los números, honestos:** para el par DH (m = 2, δ = 4.2×10⁻⁵):
n_rad,m ≈ 2.6×10⁶ y N₀ ≈ 2.7×10¹⁴ — **enorme y finita**, como
corresponde a una cota de peor caso por palomar («aunque sea inicialmente
muy grande», la vara de la casa). La realidad rompe en 1.0×10⁵: cuatro
mil millones de veces antes. Para m = 1 esta N₀ es mucho más gruesa que
la del Teorema 1 (que usa el lema de la ventana, más fino) — declarado:
ésta es de propósito general.

## 3. El teorema, con los palillos puestos — la separación A / B / C

*(Estructurado 2026-08-15 según la recomendación de la nota «los palillos
chinos»: dónde termina la demostración interna y dónde empieza el input
externo, marcado con precisión.)*

**PARTE A — TEOREMA DE INTERACCIÓN (demostración interna, autocontenida).**
*(Hipótesis completadas 2026-08-15 por la auditoría final A-H de F318:
H0 explicitada y H4 — antes oculta dentro de «la cota sellada» — cazada
por el intento de falsación F-5 y declarada.)*
Hipótesis: **H0** rᵢ > 1 para todo i (cuartetos estrictamente fuera de
la línea); **H1** m finito ≥ 1; **H2** |Im ρᵢ| ≥ 1; **H3** δ = log r_max
≤ 1; **H4** el fondo en la línea cumple la densidad N_fondo(T) ≤
(T/2π)·log T para T ≥ 2 (sin H4 la prueba se rompe: un fondo
superexponencial taparía las citas para siempre — F-5 del informe).
Conclusión: existe n ≤ N₀(r_max, m) con λₙ < 0.
Ingredientes: Dirichlet por palomar (§2a-b), lema de agenda (§2c), golpe
L1-L7 (§1), lema radial-m R0-R10 (§2d), lema δ-ζ (que deduce δ ≤ log √2
DE la hipótesis |Im ρᵢ| ≥ 1 — interno), y la cota sellada del coro
(F299-F301, interna al registro del laboratorio). **Nada de la Parte A
usa hechos sobre dónde están los ceros de ζ.**

**PARTE B — INPUT EXTERNO (resultado de la teoría de ζ, citado, no
demostrado aquí).** Los ceros no triviales de ζ cumplen |Im ρ| > 1: el
segmento real por el argumento de η (elemental, §2d-bis i), y
0 < |Im ρ| < 1 por el conteo riguroso N(14) = 0 (Backlund 1914; van de
Lune et al. 1986; Platt–Trudgian 2021 — §2d-bis ii). **Este bloque es
EXTERNO: no es consecuencia del Lema de Interacción ni de ninguna pieza
de la Parte A.**

**PARTE B (continuación) — segundo input externo, B2:** la densidad H4
para el fondo de ζ viene de Riemann–von Mangoldt con el error explícito
de Backlund (1918) — N(T) ≤ (T/2π)·log T mayora en particular al
subconjunto de ceros en la línea. EXTERNO, del mismo rango que B1.

**PARTE C — APLICACIÓN A ζ (A + B1 + B2).** Por B1 los cuartetos de
ceros de ζ cumplen H2 (y H3 vía δ-ζ); por B2 el fondo cumple H4; H0 y
H1 son la definición del escenario (una configuración finita de ceros
fuera de la línea). Luego la conclusión de A aplica: ninguna
conspiración finita de ceros de ζ fuera de la línea puede esconderse
del yunque. ∎

La subida de «candidato» a «teorema del libro» es decisión de la
auditora — ahora con cada pieza sosteniendo exactamente lo que dice
sostener.

## 4. Límites

Configuraciones FINITAS; la N₀ es de peor caso (exponencial en m);
mediciones en una base y una ventana. Nada sobre RH.

---

*Los dos pedidos, entregados: el golpe sin saltos y la cita con fecha,
hora y cota. La agenda de Dirichlet no tiene página en blanco.*
