# Respuesta a la auditoría del Corolario del Río de Pozos

**Para la auditora · 2026-08-15 · responde su PEDIDO F323, punto por
punto.**

## §0 — Autoría, conservada en el registro

Tal cual quedó en F330 y acá se ratifica: **la idea conceptual del Río
de Pozos es de Nico** (el flash y la historia-mapa del claro, el río y
el pozo); Doc hizo después la traducción y formalización matemática.
La distinción queda en las tres copias del registro.

## §4 — La iteración, por inducción formal (I1-I4)

Fijada una configuración bajo H0-H4 (invariante en toda la inducción:
las hipótesis son de la CONFIGURACIÓN, no del paso k — mismo m, mismos
θᵢ, mismo δ, mismo n_rad).

**(I1) Base.** D4 con T₁ = n_rad (entero ≥ 11 ≥ 1 ✓) produce n₁ con:
n₁ ≥ T₁ = n_rad, n₁ ≤ T₁ + ⌈2πT₁⌉^m, y ‖n₁θᵢ‖ ≤ 1 ∀i. Existe n₁. ✓

**(I2) Paso inductivo.** Hipótesis: existe n_k, cita con n_k ≥ n_rad.
Sea T_{k+1} = n_k + 1 — entero (n_k ∈ ℤ) y ≥ 1 ✓. Las condiciones que
D4 exige de su parámetro T son exactamente dos: T entero y T ≥ 1
(la traza de F327 usa T entero en J = ⌈T/n₁⌉ ≤ T, y T ≥ 1 en
2π/Q ≤ 1/T). Ambas se cumplen. D4 aplica y produce n_{k+1} con:

    n_{k+1} ≥ T_{k+1} = n_k + 1 > n_k          [estrictamente mayor]
    n_{k+1} ≤ T_{k+1} + ⌈2π·T_{k+1}⌉^m         [cota explícita, §7]
    ‖n_{k+1}·θᵢ‖ ≤ 1  ∀i                       [cita, mismo sentido]

**(I3) Hipótesis preservadas.** H0-H4: globales, sin dependencia de k ✓.
Las de D4: verificadas en I2 para cada k ✓. Las de D5 en n_{k+1}:
n_{k+1} > n_k ≥ n_rad ≥ 11 ≥ 3 (umbral de L5) ✓ y la condición angular
ε = 1 es la línea tercera de I2 ✓. Nada se pierde al pasar de k a k+1.

**(I4) Conclusión.** Por I1 + I2 + inducción sobre ℕ: existe una
sucesión n₁ < n₂ < n₃ < … estrictamente creciente e infinita (cada
paso produce un término nuevo; ninguna hipótesis se agota). ∎

## §5 — D5 en cada cita

Por I3, cada n_k cumple TODAS las hipótesis de D5 (cita ε = 1, n ≥ 3,
H4 global): λ_{n_k} ≤ −g(n_k), para todo k. ∎

## §6 — Crecimiento y divergencia de los pozos

D3 (auditada en F327/F328: g estrictamente creciente en [n_rad, ∞), con
igualdad solo en el punto inicial): n_k < n_{k+1}, ambos ≥ n_rad ⟹
g(n_k) < g(n_{k+1}). Estricto para todo k. Y como los n_k son enteros
estrictamente crecientes, n_k ≥ n_rad + (k−1) → ∞; con la cota ya
auditada g(n) = e^{nδ} − (4/π)n·log n − (2m+2) → ∞ (el exponencial
manda, R7-R10), sigue g(n_k) → ∞. ∎

## §7 — «Programable», con procedimiento y cota escrita

No es existencia abstracta: cada paso es un PROCEDIMIENTO con receta y
cota. Dado n_k:

    1. T := n_k + 1;  Q := ⌈2π·T⌉                    [dos operaciones]
    2. n₁' := la cita fina de Dirichlet para Q        [búsqueda ≤ Q^m]
    3. n_{k+1} := ⌈T/n₁'⌉·n₁'                          [una división y un producto]

**Cota explícita del paso:**  n_{k+1} ≤ (n_k + 1) + ⌈2π·(n_k + 1)⌉^m.
Es la misma N₀ de DYN con n_rad reemplazado por n_k + 1 — enorme y de
peor caso, como toda cota de palomar, pero ESCRITA. (En el testigo la
realidad es incomparablemente más fina: 435 citas en 3000 escalones,
brecha media ≈ 7 contra una cota de paso ~10¹³ — evidencia, no parte
de la prueba.)

## §9 — El ataque hostil, punto por punto

- **¿Puede detenerse la iteración?** No: D4 solo exige T entero ≥ 1, y
  T_{k+1} = n_k + 1 lo cumple siempre. No hay recurso que se agote —
  Dirichlet vale para todo Q, y Q crece sin techo disponible que romper.
- **¿Dependencia circular?** No: T_{k+1} depende solo de n_k (ya
  construido), y la cota del paso depende solo de T_{k+1}. El grafo de
  dependencias es una cadena hacia adelante, sin ciclos.
- **¿Se pierde una hipótesis en k → k+1?** No — I3: H0-H4 son de la
  configuración (invariantes), las de D4 se re-verifican en cada paso, y
  las de D5 mejoran con k (n_{k+1} más adentro del claro que n_k).
- **¿n_k → ∞ sin g estrictamente creciente?** Imposible en este
  corolario: todos los n_k viven en [n_rad, ∞) donde D3 da monotonía
  ESTRICTA — el único empate posible de D3 (el punto inicial) requeriría
  n_k = n_{k+1} = n_rad, excluido por n_{k+1} > n_k.
- **¿«Cita» cambia de sentido?** No: en cada iteración la conclusión
  angular es la misma desigualdad ‖nθᵢ‖ ≤ 1 ∀i (I2, línea 3) — el Q
  interno cambia con k, pero la DEFINICIÓN de cita (ε = 1) es idéntica;
  ningún paso usa una noción más débil.
- **¿Contraejemplo bajo H0-H4?** No lo encontramos: la cadena es
  composición de lemas auditados (D4 por inducción + D5 + D3), sin
  eslabón nuevo que atacar por fuera de los ya atacados en F327/F328.

## §10 — Acordado

Las 435 citas del testigo son apoyo; los cuantificadores los cierra la
inducción I1-I4, no la corrida.

## §12 — La pregunta final, respondida

**Sí: DYN + Diosyunalma convierten la ruptura única en un río infinito
de rupturas programables, cada una con su pozo más profundo que el
anterior** — con la novedad delimitada como usted exige: la infinitud
de λ negativos por otra vía no es nueva; lo nuevo es la construcción
explícita (receta §7) de citas con fecha, cota de paso escrita y pisos
crecientes predichos.

## §11 — Veredicto solicitado (autoevaluación; el sello es suyo)

    🟢 COROLARIO CERRADO — inducción I1-I4 completa, D5 y D3 aplicados
    con hipótesis verificadas en cada k, programabilidad con receta y
    cota escrita, ataque hostil sin ruptura. La declaración oficial y
    el sello, como siempre, son de la relojera. Todavía no.

---

*La idea nació del capitán; la inducción dice que el río realmente
corre — y cada vuelta de la agenda es una pala más honda.* 🌊🕳️📐
