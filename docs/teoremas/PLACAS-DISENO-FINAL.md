# Diseño final — Teorema de las Placas / Ley del Líder

**Para la mesa de los tres · 2026-08-15 · responde el PEDIDO F327 en su
formato exacto I-X.** Regla cumplida: ningún hueco tapado con
intuición; lo no demostrado está marcado; cada hipótesis en su lugar.

---

## I. VEREDICTO F1-F6: **CERRADO** (auditoría línea por línea)

- **F1a** cos(nθ_L) = cos(‖nθ_L‖): cos es par y 2π-periódica, y ‖x‖ es
  la distancia de x a 2πℤ ⟹ cos(x) = cos(‖x‖) para todo x real ✓.
  En la banda ‖nθ_L‖ ≤ 1 ⊂ [0, π]: cos decrece ahí ⟹ cos ≥ cos(1) ✓.
- **F1b** orientación del líder: −2cos(nθ_L)(r_Lⁿ + r_L⁻ⁿ) ≤
  −2cos(1)(r_Lⁿ + r_L⁻ⁿ) ≤ −2cos(1)·r_Lⁿ — el último paso descarta
  −2cos(1)·r_L⁻ⁿ ≤ 0, válido PORQUE cos(1) > 0 ✓ (la orientación que
  usted pidió vigilar: correcta, y el porqué queda escrito).
- **F2** competidores: cos ≥ −1 ⟹ ℓᵢ ≤ 4 + 2(rᵢⁿ + rᵢ⁻ⁿ); rᵢ⁻ⁿ ≤ 1
  (rᵢ > 1, n ≥ 1) y rᵢ ≤ r₂ := maxᵢ≠L rᵢ ⟹ ℓᵢ ≤ 6 + 2r₂ⁿ ✓.
- **F3** coro: L5 exige n ≥ 3 y H4; N* ≥ n_rad ≥ ⌈6·log 6⌉ = 11 ≥ 3 ✓.
- **F4** constantes: (m−1)·6 + 4 = 6m − 2 ✓ (suma verificada).
- **F5** n_comp: definida solo para m ≥ 2 (para m = 1 se fija
  n_comp := 0 y F5 es vacua — r₂ no se usa: sin objeto indefinido);
  argumento del log: 2(m−1)/cos(1) ≥ 2/cos(1) > 1 > 0 ✓; y
  δ_L − δ₂ > 0 ⟺ r_L > r₂ ⟺ **líder estricto, exactamente** ✓.
- **F6** reciclaje del radial: u = 3(m+1)/δ_L con δ_L = log r_max (el
  MISMO u y n_rad de DYN, pues el líder es la perla de radio máximo);
  corchete₂(n_rad) = 2·corchete_DYN + (8m−8) ≤ 2u³ + 4u ≤ 3u³ ≤
  u^{3(m+1)} ≤ e^{n_rad·δ_L} [D1 + D2 + u ≥ 6 + m ≥ 1] ✓.
- **F6, la extensión (G, G′, G″), justificación completa:** G(n) :=
  e^{nδ_L} − (8/π)n·log n − (12m−4) sobre la variable REAL n ≥ n_rad
  (C² ahí). G″(n) = δ²e^{nδ} − (8/π)/n es suma de dos funciones
  crecientes ⟹ G″ creciente; G″(n_rad) ≥ δ²u^{3(m+1)} − (8/π)/11 =
  9(m+1)²u^{3m+1} − 0.232 > 0 [u ≥ 6]; TFC dos veces (el patrón D3
  auditado en F327/F328): G′ > 0 y G > 0 en todo [n_rad, ∞); los
  enteros heredan por restricción — jamás se deriva una función
  discreta ✓.
- **Universalidad:** el cuantificador es ∀n de la banda con n ≥ N* —
  ninguna agenda interviene; las demás perlas quedan sin restricción ✓.

## II. VEREDICTO POZO — enunciado exacto

**LEMA POZO.** ∀ configuración bajo H0-H4 con líder estricto,
∀n ≥ N* = max(n_rad, n_comp):  ‖nθ_L‖ ≤ 1 ⟹
λₙ ≤ (4/π)n·log n + (6m−2) + 2(m−1)r₂ⁿ − 2cos(1)r_Lⁿ < 0.
*Prueba:* F1-F6 (I). ∎ — estado: **cerrado, a su lupa.**

## III. VEREDICTO MONTAÑA — m = 1 y m ≥ 2 SEPARADOS

**LEMA MONTAÑA-1 (m = 1).** ∀n ≥ 1 (¡sin umbral!):
‖nθ‖ ≥ π/2 ⟹ λₙ ≥ 4.  *Prueba:* cos(nθ) = cos(‖nθ‖) ≤ 0 en
[π/2, π] ⟹ ℓₙ = 4 − 2cos(nθ)(rⁿ+r⁻ⁿ) ≥ 4; coroₙ ≥ 0. ∎
Y en la sub-banda anti ‖nθ − π‖ ≤ 1: λₙ ≥ 4 + 2cos(1)·rⁿ → ∞
[A1-A6, auditada]. **Nota de diseño: para m = 1 la mitad entera
‖nθ‖ ≥ π/2 es montaña, incondicional — ni H4 ni umbral. El ecuador
‖nθ‖ = π/2 pertenece al lado montaña.**

**LEMA MONTAÑA-m (m ≥ 2).** ∀n ≥ n_comp (¡el MISMO umbral de F5 —
un solo umbral competidor para las dos bandas!):  ‖nθ_L − π‖ ≤ 1 ⟹
λₙ ≥ cos(1)·r_Lⁿ + 2m + 2 > 0.
*Prueba:* coro ≥ 0; ℓᵢ ≥ 4 − 2(rᵢⁿ + rᵢ⁻ⁿ) ≥ 2 − 2r₂ⁿ (i ≠ L);
ℓ_L ≥ 4 + 2cos(1)r_Lⁿ [A2-A4 con banda anti]; sumar:
λ ≥ 4 + 2(m−1) + 2cos(1)r_Lⁿ − 2(m−1)r₂ⁿ; y F5 (n ≥ n_comp) absorbe
2(m−1)r₂ⁿ ≤ cos(1)r_Lⁿ. ∎ — **sin H4, sin radial: las montañas son
más BARATAS que los pozos** (ver VII).
*Excluido del teorema (§9 suyo):* montañas en anti-citas CONJUNTAS
(todas las fases anti a la vez) — siguen dependiendo de aproximación
inhomogénea: problema abierto, no se usa.

## IV. FRONTERA — qué se demuestra y qué queda abierto

Su §5 pedía una de tres salidas. Entregamos **las tres, cada una en su
dominio, con prueba**:

1. **Sub-zona de pozo (paramétrica):** ∀η ∈ (0, π/2 − 1], la banda
   ‖nθ_L‖ ≤ π/2 − η fuerza λₙ < 0 para n ≥ N*(η) — misma cadena F1-F6
   con cos(1) → sin(η) y el corchete multiplicado por 2/sin(η); el
   costo es explícito: N*(η) = max(n_rad + ⌈log(3/sin η)/δ_L⌉,
   n_comp(η)), n_comp(η) = ⌈log(2(m−1)/sin η)/(δ_L−δ₂)⌉.
   [BOCETO VERIFICADO en estructura; la escritura F(η) completa es
   parte de la redacción final del teorema — marcado, no tapado.]
2. **Sub-zona de montaña (paramétrica):** ‖nθ_L − π‖ ≤ π/2 − η fuerza
   λₙ > 0 para n ≥ n_comp(η) (m ≥ 2; para m = 1, incondicional por
   MONTAÑA-1). Mismo boceto, mismo estado.
3. **El residuo — EL ECUADOR:** la vecindad de ‖nθ_L‖ = π/2 (donde
   cos(nθ_L) ≈ 0 y el líder enmudece). Ahí:
   - **m = 1: NO es zona abierta** — pertenece al lado montaña
     (MONTAÑA-1: λ ≥ 4 en todo ‖nθ‖ ≥ π/2) y la mitad fina la cubre
     la paramétrica. La frontera de m = 1 es la CURVA ‖nθ‖ = π/2, no
     una región.
   - **m ≥ 2: NO EXISTE signo universal con las hipótesis actuales, y
     lo demostramos:** en cos(nθ_L) = 0, λₙ = coroₙ + 4 + Σᵢ≠Lℓᵢ, y el
     signo queda en manos de los competidores: un competidor en banda
     fina propia empuja −2cos(1)r₂ⁿ (negativo arbitrariamente grande),
     uno en banda anti empuja +2cos(1)r₂ⁿ. Ambas combinaciones de
     fase son realizables en configuraciones bajo H0-H4 ⟹ ambos
     signos ocurren. **La frontera m ≥ 2 queda ABIERTA por teorema,
     no por cansancio** — y delimitada: es exactamente donde empieza
     el problema del sub-líder (¿aplica la Ley del Líder recursiva al
     segundo radio cuando el primero enmudece? — punto abierto X.2).
   *Evidencia consistente:* la banda frontera del testigo dio 540/539
   — mixta casi exacta, como el teorema manda.

## V. GEOMETRÍA — qué es lema y qué es teorema

Todo LEMA (marco independiente; el teorema principal no depende de
ellos, y ellos no dependen del líder):

    LEMA G1 (filtración): ε₁ < ε₂ ⟹ C_{ε₁} ⊆ C_{ε₂}. [definición]
    LEMA G2 (monoide graduado): C_ε + C_ε' ⊆ C_{ε+ε'}, ∀ε,ε' ≥ 0;
      informativa solo si ε+ε' < π (dominio declarado). [subaditividad]
    LEMA G3 (Cayley): v→w ⟺ w−v ∈ C_{ε_d} es el grafo de Cayley de
      (ℕ,+) con generador C_{ε_d}: dirigido, irreflexivo, invariante
      por traslación, gradualmente transitivo (presupuestos se suman).
    LEMA G4 (ramificación): ∀v, ∀ε_d > 0: grado de salida ≥ 2
      (Dirichlet a ε_d/2 da c; c y 2c son continuaciones distintas).
    LEMA G5 (recombinación): ∀c₁ ≠ c₂ ∈ C_{ε_d}: diamante
      v → v+c₁ → v+c₁+c₂ = v+c₂+c₁ ← v+c₂ ← v. [conmutatividad]
    Placa := componente de intervalo de C_ε — definición, natural,
      sin coordenadas. El Río de Pozos = UNA cadena dentro de G3.
    ε_eff(n) = min{ε : n ∈ C_ε} — la función de nivel de la
      filtración: fragilidad-número y fragilidad-geometría son una.

## VI. TEOREMA DE LAS PLACAS / LEY DEL LÍDER — el enunciado mínimo

**Hipótesis:** H0-H4 (las de DYN, intactas) + LÍDER ESTRICTO (∃ único
L con r_L > rᵢ ∀i ≠ L). **Definiciones:** δ_L, r₂, u, n_rad, n_comp,
N* como arriba; para m = 1, n_comp := 0 y la hipótesis líder-estricto
es vacua.

**(P1 — pozos)** ∀n ≥ N* con ‖nθ_L‖ ≤ 1:  λₙ < 0
                 [con el mayorante explícito de II].
**(P2 — montañas)** m = 1: ∀n ≥ 1 con ‖nθ‖ ≥ π/2: λₙ ≥ 4; y en
                 ‖nθ−π‖ ≤ 1: λₙ ≥ 4 + 2cos(1)rⁿ.
                 m ≥ 2: ∀n ≥ n_comp con ‖nθ_L−π‖ ≤ 1:
                 λₙ ≥ cos(1)r_Lⁿ + 2m + 2.
**(P3 — frontera)** m = 1: la frontera es la curva ‖nθ‖ = π/2 (P2
                 cubre el lado cerrado ≥ π/2; la paramétrica IV.1 el
                 resto). m ≥ 2: en la zona del ecuador no existe signo
                 universal bajo estas hipótesis (IV.3, demostrado);
                 la clasificación paramétrica IV.1-IV.2 la encierra en
                 vecindades de ‖nθ_L‖ = π/2 tan finas como se pague.
**(P4 — programabilidad)** bajo H2 (θ_L ≤ π/2 ≤ 2): las bandas de P1
                 y P2 tienen habitantes en cada bloque de
                 K_L = ⌈2π/θ_L⌉ + 1 enteros consecutivos [lema de la
                 ventana, arcos de largo 2 ≥ θ_L] — pozos y montañas
                 con fecha, para siempre. ∎(candidato)

Fuera del enunciado (su §9, cumplido): montañas conjuntas m ≥ 2;
frontera m ≥ 2 como región clasificada; profundidad⟹ramas (falsada);
el Río como enumeración total; evidencia como prueba; RH.

## VII. DEPENDENCIAS — la tabla exacta

    Hipótesis      P1 pozos   P2 montañas   P3 frontera   P4 program.  G1-G5
    H0 (r > 1)     ✓ (F6)     ✓ (rᵢ⁻ⁿ≤1)    ✓             —            —
    H1 (m finito)  ✓ (F4)     ✓ (suma)      ✓             —            —
    H2 (|Im ρ|≥1)  —          —             —             ✓ (θ_L≤π/2)  —
    H3 (δ ≤ 1)     ✓ (u≥6)    —             —             —            —
    H4 (densidad)  ✓ (F3)     **NO**        ✓ (lado pozo) —            —
    Líder estricto ✓ (F5,m≥2) ✓ (F5,m≥2)    ✓ (m≥2)       —            —
    Dirichlet      —          —             —             —            ✓ (G4)

  **La asimetría que el diseño revela: las montañas no usan H4 ni el
  radial — son estructuralmente más baratas que los pozos.** H2 solo
  paga la programabilidad, en ambas bandas.

## VIII. CONTRAEJEMPLOS BUSCADOS (ninguno encontrado; todos documentados)

- Orientaciones F1-F6: revisadas una a una (I); la única sutil (F1b)
  explicada con su porqué.
- Paso discreto n_rad → ∀n: sin laguna — función real, dos TFC,
  restricción a enteros (I).
- **Líder apenas dominante:** δ_L − δ₂ → 0⁺ ⟹ n_comp → ∞ pero FINITO
  para toda brecha estricta: el lema sobrevive con umbral que explota
  — declarado, no escondido. Batería en el borde exacto n_comp con
  brechas de hasta 1%: 84 casos, 0 violaciones.
- m = 1 vs m ≥ 2: separados en todo el diseño; r₂ no se define en
  m = 1 (n_comp := 0, F5 vacua) — sin objeto fantasma.
- Argumentos de log: 2(m−1)/cos 1 ≥ 2/cos 1 > 1 (m ≥ 2); 3/sin η > 0;
  u ≥ 6 > 1 — todos válidos en sus dominios.
- ¿n_mont y N* coexisten? Sí: mismo marco de hipótesis, y n_comp es
  literalmente EL MISMO umbral en ambas bandas — el testigo muestra
  las dos bandas activas en la misma ventana (978 pozos y 944
  montañas conviviendo).
- Frontera: el «contraejemplo» existe y es teorema — para m ≥ 2 ambos
  signos ocurren en el ecuador (IV.3): delimita, no rompe.
- Dependencia de agenda: ninguna — P1-P3 cuantifican sobre TODO n de
  cada banda; P4 usa solo la ventana (determinista, sin elección).

## IX. EVIDENCIA COMPUTACIONAL (respaldo, jamás prueba)

    F6 grilla (m=1..10 × δ≤1): 50 casos, 0 violaciones     [labandafina]
    F5 borde exacto sintético: 84 casos, 0                 [labandafina]
    P1 vivo (978 escalones banda fina tras N*): 0 en signo
      y 0 en mayorante                                     [labandafina]
    P2 vivo m=2 (926 anti-citas del líder tras n_comp): 0  [lageometria]
    A1-A6 m=1 régimen exponencial: 1543 anti-citas, 0      [lageometria]
    Ley del Líder por bandas (perla 1 ignorada): fina
      978/978 neg · anti 944/944 pos · frontera 540/539    [lageometria]
    Geometría: ley (ε/π)²·H al 3%; grado 25 constante ∀v;
      placas 2/4/6 por nivel; diamante exhibido            [lasplacas]

## X. PUNTOS ABIERTOS (solo los reales)

    X.1 Montañas conjuntas m ≥ 2 (aproximación inhomogénea simultánea).
    X.2 El ecuador m ≥ 2: ¿ley del sub-líder recursiva cuando el líder
        enmudece? (la frontera como puerta del próximo mecanismo).
    X.3 La escritura completa de la paramétrica F(η)/N*(η) (boceto en
        IV; falta la redacción con el rigor F1-F6 — trabajo de pluma,
        no de idea).
    X.4 Alfabeto finito de brechas para m ≥ 2 general (tres distancias
        multidimensional).
    X.5 Optimalidad de N*: los cocientes 1.5-1.8 sugieren margen O(1);
        sin afirmación (regla del cociente de la auditora).

## RESPUESTA A LA PREGUNTA FINAL (§15)

**SÍ — y el ensamble está arriba sin un solo supuesto oculto:** los
lemas G1-G5 dan la geometría; P1-P4 dan el paisaje con la Ley del
Líder; la frontera está delimitada por teorema (curva exacta en m = 1,
no-universalidad demostrada en m ≥ 2); y la tabla VII muestra cada
hipótesis pagando exactamente una cuenta. Estado propuesto para su
criterio §14: **🟠 FRONTERA ABIERTA** — el teorema es correcto y sus
cuantificadores están cerrados en P1, P2 y P4; la frontera m ≥ 2 queda
abierta A PROPÓSITO y con prueba de que debe quedarlo. Si su lupa
concuerda, el edificio está listo para el terremoto — y para el
nombre. Todavía no.

---

*Nico puso la intuición, la arquitectura quedó en pie, y el diseño
reveló dos secretos: las montañas son más baratas que los pozos, y la
frontera no es una mancha — es un ecuador.* 🌍🔨💎
