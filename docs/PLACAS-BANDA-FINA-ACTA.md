# El punto amarillo, cerrado — banda fina → N* → λₙ < 0

**Para la auditora · 2026-08-15 · el único punto que usted pidió cerrar
antes de sentarnos a escribir el Teorema de las Placas.** Cadena F1-F6,
todos los cuantificadores a la vista, N* explícito.

---

## Enunciado (candidato a lema de la banda fina)

**∀** configuración bajo H0-H4 (las de DYN, intactas) **con líder
estricto** — ∃ único L con r_L > rᵢ ∀i ≠ L; sea r₂ := maxᵢ≠L rᵢ,
δ_L = log r_L, δ₂ = log r₂ (para m = 1 no hay competidores y δ₂ no se
usa) — y **∀** n entero con

    n ≥ N* := max( n_rad , n_comp )
    n_rad  = ⌈u·log u⌉,  u = 3(m+1)/δ_L          [el de DYN, sin cambios:
                                                  δ_L ES el δ de DYN, pues
                                                  el líder tiene r máximo]
    n_comp = ⌈ log(2(m−1)/cos 1) / (δ_L − δ₂) ⌉  [si m ≥ 2; := 0 si m = 1]

vale la implicación:  **‖n·θ_L‖ ≤ 1  ⟹  λₙ < 0**, con el mayorante
explícito λₙ ≤ (4/π)n·log n + (6m−2) + 2(m−1)r₂ⁿ − 2cos(1)·r_Lⁿ. ∎(cand.)

Nótese el cuantificador fuerte: vale para TODO n de la banda fina más
allá de N* — no solo para citas conjuntas ni para n elegidos por
agenda alguna. Las demás perlas quedan SIN restricción de fase.

## La cadena F1-F6

    (F1) [∀n en banda fina] ‖nθ_L‖ ≤ 1 ⟹ cos(nθ_L) = cos(‖nθ_L‖) ≥ cos(1)
         [paridad + período + cos decreciente en [0, π]]
         ⟹ ℓ_L = 4 − 2cos(nθ_L)(r_Lⁿ + r_L⁻ⁿ) ≤ 4 − 2cos(1)·r_Lⁿ
         [r_Lⁿ + r_L⁻ⁿ ≥ r_Lⁿ > 0; orientación ✓]
    (F2) [∀n, ∀i ≠ L] ℓᵢ ≤ 4 + 2(rᵢⁿ + rᵢ⁻ⁿ) ≤ 6 + 2r₂ⁿ
         [cos ≥ −1; rᵢ⁻ⁿ ≤ 1; rᵢ ≤ r₂ por definición de r₂]
    (F3) [∀n ≥ 3] coroₙ ≤ (4/π)·n·log n    [L5 bajo H4, sellada F299-F301;
         n ≥ N* ≥ n_rad ≥ 11 ≥ 3 lo garantiza]
    (F4) sumando F1+F2+F3 (m−1 competidores):
         λₙ ≤ (4/π)n·log n + (6m−2) + 2(m−1)r₂ⁿ − 2cos(1)·r_Lⁿ
         [(m−1)·6 + 4 = 6m−2; el mayorante del enunciado]
    (F5) [∀n ≥ n_comp, m ≥ 2] 2(m−1)·r₂ⁿ ≤ cos(1)·r_Lⁿ
         [⟺ n(δ_L − δ₂) ≥ log(2(m−1)/cos 1); exige δ_L > δ₂: exactamente
         la hipótesis líder-estricto — ahí y solo ahí se usa. m = 1: vacua]
    (F6) [∀n ≥ n_rad] cos(1)·r_Lⁿ > (4/π)n·log n + (6m−2)
         [el radial de DYN, RECICLADO con corchete duplicado:
          cos(1) > 1/2 ⟹ basta e^{nδ_L} > corchete₂(n) := (8/π)n·log n + (12m−4).
          En n_rad:  corchete₂ = 2·corchete_DYN + (8m−8), y
          corchete_DYN ≤ u³ (D2) ∧ 8m−8 < 4(2m+2) ≤ 4u ≤ u³ (u ≥ 6)
          ⟹ corchete₂(n_rad) ≤ 3u³ ≤ u^{3(m+1)} ≤ e^{n_rad·δ_L} (D1;
          3 ≤ u^{3m} pues u^{3m} ≥ 6³). Extensión a todo n ≥ n_rad por el
          patrón D3 ya auditado: G := e^{nδ_L} − corchete₂ tiene
          G″ = δ²e^{nδ} − (8/π)/n suma de crecientes, G″(n_rad) > 0 y
          G′(n_rad) > 0 con márgenes enormes (verificados en batería),
          y dos integraciones anidadas dan G > 0 en [n_rad, ∞).]

**Conclusión.** Para n ≥ N* en la banda fina: por F4,
λₙ ≤ P(n) + 2(m−1)r₂ⁿ − 2cos(1)r_Lⁿ con P(n) = (4/π)n log n + 6m−2;
por F5, 2(m−1)r₂ⁿ ≤ cos(1)r_Lⁿ ⟹ λₙ ≤ P(n) − cos(1)r_Lⁿ; por F6,
P(n) < cos(1)r_Lⁿ ⟹ **λₙ < 0**. Cada desigualdad con su orientación. ∎

**Corolario de programabilidad** (con lo ya auditado): el lema de la
ventana sobre el arco fino (largo 2 ≥ θ_L, garantizado bajo H2 por
θ_L ≤ π/2) pone habitantes de la banda fina en cada bloque de
K_L = ⌈2π/θ_L⌉ + 1 enteros consecutivos: los pozos de la banda fina
son PROGRAMABLES por ventana más allá de N*.

## Dónde entra cada hipótesis (el mapa de dependencias)

    H0 (r > 1): δ_L > 0 — F6 y la definición de u.
    H1 (m finito): la suma F4 es finita.
    H2 (|Im ρ| ≥ 1): NO se usa en el signo — solo en el corolario de
       programabilidad (θ_L ≤ π/2 ≤ 2). Declarado.
    H3 (δ ≤ 1): u ≥ 6 — F6 (vía D1/D2 de Diosyunalma).
    H4 (densidad): F3, y solo ahí.
    LÍDER ESTRICTO (δ_L > δ₂): F5, y solo ahí. Para m = 1 no hace falta.

## Verificación (evidencia, jamás prueba)

- **F6 en la grilla** (m = 1..10 × δ ≤ 1, 50 casos): corchete₂ ≤ 3u³,
  la desigualdad F6 en n_rad, G′ > 0 y G″ > 0 — cero violaciones.
- **F5 en el borde exacto** n_comp (84 casos sintéticos, m = 2..8,
  brechas de líder hasta el 1%): cero violaciones.
- **La cadena entera en vivo** (testigo m = 2: n_rad = 1 040 809,
  n_comp = 23 064, N* = 1 040 809): los **978 escalones de banda fina**
  de la ventana [N*, N*+3000] — mayorante F4: cero violaciones; signo
  λ < 0: **cero violaciones**. (Los 978/978 de F334, ahora con su
  teorema-candidato debajo.)

`go run ./cmd/labandafina`

## Estado

    El punto amarillo queda cerrado: banda fina → N* → λₙ < 0, con
    F1-F6 escritas, N* explícito, cuantificadores fuertes (∀n de la
    banda, sin agenda), y el mapa de qué hipótesis entra dónde.
    Las dos bandas externas del paisaje están ahora al mismo nivel:
      · anti  ⟹ montaña, ∀n ≥ n_mont   [F333/F334, A1-A6]
      · fina  ⟹ pozo,    ∀n ≥ N*       [esta acta, F1-F6]
    Si su lupa no encuentra grieta, la razón seria existe: sentarnos a
    escribir el TEOREMA DE LAS PLACAS / LEY DEL LÍDER. El sello, el
    ensamble y el nombre del próximo gran teorema son de la mesa de los
    tres. Todavía no.

---

*El amarillo era la mitad fina del paisaje. Cerrado con el radial de
DYN reciclado — corchete duplicado, tres u³ de sobra — y el líder
pagando la cuenta de sus competidores desde n_comp.* 🟡→🟢⛏️
