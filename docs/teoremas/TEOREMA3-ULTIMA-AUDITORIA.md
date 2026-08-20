# Última auditoría del candidato T3 — los hielos del cóctel

**Para la auditora · 2026-08-15 · responde su PEDIDO F322 (última
auditoría T3).** Complementa la respuesta F327
(`docs/teoremas/TEOREMA3-AUDITORIA-RESPUESTA.md`), que cerró D3/D4/D6 y el primer
ataque; acá van los puntos NUEVOS de su pedido: los tres refinamientos
de §4, los tres ataques quirúrgicos de §7, y la evaluación de futuro
del auxiliar (§9). H0-H4 intactas; cero externos nuevos.

---

## §4 — los tres refinamientos de D3

**(4a) Dominio continuo/discreto, sin laguna.** g se define sobre la
variable REAL n ∈ [n_rad, ∞), donde es C² (n_rad ≥ 11 > 0; log n y 1/n
sin singularidades). TODA la maquinaria de derivadas (R8-R10, D3a-D3c de
F327) opera sobre esa función real. Los lemas la evalúan después en
enteros n — que son un subconjunto del intervalo real, así que la
monotonía demostrada sobre ℝ se hereda por restricción. **En ningún paso
se deriva una función discreta**: se acota la real y se restringe. Sin
laguna.

**(4b) La inferencia completa, sin «por convexidad»:**

    (i)   g″(n) = δ²e^{nδ} − (4/π)/n. El término δ²e^{nδ} es
          estrictamente creciente (δ > 0 por H0); el término −(4/π)/n
          también es creciente (sube desde negativo hacia 0). Suma de
          crecientes: g″ estrictamente creciente en [n_rad, ∞).
    (ii)  R9 da la base: g″(n_rad) > 0. Con (i): g″(t) > 0 ∀t ≥ n_rad.
    (iii) TFC: g′(n) = g′(n_rad) + ∫_{n_rad}^n g″(t)dt. El integrando es
          > 0 por (ii) y g′(n_rad) > 0 por R8 ⟹ g′(n) > 0 ∀n ≥ n_rad.
    (iv)  TFC otra vez: g(n) = g(n_rad) + ∫_{n_rad}^n g′(t)dt ≥ g(n_rad),
          con igualdad solo en n = n_rad. g creciente. ∎

La palabra «convexidad» no se usa: la cadena es monotonía de g″ + dos
integraciones, cada una escrita.

**(4c) Dónde entra n_rad ≥ 11, exactamente — en TRES lugares:**

    1. L5 (la cota del coro) exige n ≥ 3: lo da n ≥ n_rad ≥ 11 ≥ 3.
    2. La agenda exige T = n_rad entero ≥ 1 (el paso J = ⌈T/n₁⌉ ≤ T
       de D4-(iii) usa T entero positivo): lo da n_rad ≥ 11 ≥ 1.
    3. R0/R2: n* = u·log u ≥ 6·log 6 = 10.75 (de u ≥ 6) es lo que
       produce n_rad ≥ 11 Y alimenta la constante 1.094 = 1 + 1/10.75
       de R2, que a su vez sostiene D2.
    Y el origen de u ≥ 6: H3 (δ ≤ 1) ⟹ u = 3(m+1)/δ ≥ 3·2/1 = 6, ∀m ≥ 1.

## §5 y §6 — cerrados en F327, con el punto final de su lista

La traza del único n (cuatro propiedades, cada desigualdad con su
porqué) y la cadena D6 orientada están completas en la respuesta F327.
El punto final de su §5 — **«el mismo n puede introducirse en L7, L5 y
D5»** — queda explícito: L7 recibe de (iii) exactamente su hipótesis
(‖nθᵢ‖ ≤ 1 ∀i, con ε = 1 ∈ (0, √2)); L5 recibe de (iv) su umbral
(n ≥ 11 ≥ 3) y H4 es global; D5 es literalmente la suma de esas dos
conclusiones evaluadas EN ese n. Un solo n, tres lemas alimentados. ∎

## §7 — los tres ataques quirúrgicos (los hielos)

**(a) Techo de n_rad MÁXIMO.** Construcción adversarial: para cada m
elegí u tal que u·log u queda 10⁻³⁰ por encima de un entero — el exceso
del techo es máximo (≈ 1), el peor caso para D2 (n_rad más grande que
lo necesario). 30 casos (m = 1..10, δ ≤ 1), a 60 dígitos: **D1 y D2
aguantan los 30, cero fallas.** (El caso general está cubierto por R2:
n_rad ≤ n* + 1 ≤ 1.094·n*, válido siempre que n* ≥ 10.75 — o sea
siempre, por 4c-3.)

**(b) Resonancia EXACTA en el límite ε = 1.** El peor caso legal:
TODAS las m fases clavadas en ‖nθᵢ‖ = 1 y TODAS las perlas con
r = r_max. L7 exige m·[4 − 2cos(1)(rⁿ+r⁻ⁿ)] ≤ 2m+2 − rⁿ. El margen es

    (2m·cos(1) − 1)·rⁿ + [términos ≥ −2m−2 acotados]

y el coeficiente dominante **2m·cos(1) − 1 ≥ 2cos(1) − 1 = 0.0806 > 0
para todo m ≥ 1** — positivo por estructura (cos 1 > ½ es la holgura de
la parábola de L1 en el borde), no por suerte. Batería: m hasta 100,
rⁿ desde 1 hasta e⁷⁰⁰, 30 casos: **cero fallas**; el caso rⁿ = 1 exacto
da margen 0.161m + 1 > 0. El límite cerrado ε = 1 está DENTRO del
dominio válido (ε < √2) con desigualdad estricta.

**(c) Coro al MÁXIMO permitido por H4.** No hay ataque posible por
construcción: la cadena D5 descuenta el **100% del presupuesto H4**
((4/π)·n·log n entero — en el testigo: 1.8×10⁷ contra r_maxⁿ = 4.3×10⁴⁴,
38 órdenes de diferencia). Cualquier coro legal — incluido el
adversarial de densidad máxima de F319, que solo alcanza ~60% del
presupuesto — es MENOR o igual que lo ya descontado: agrega margen,
jamás lo quita. Un coro que rompa la cota tendría que violar H4, y eso
está fuera de las hipótesis.

**(+ los extremos de δ y m** — δ → 1, δ → 0, m = 1, m grande — quedaron
verificados en F327 §7-A/D: 42 casos hasta m = 1000 y δ = 10⁻¹², cero
fallas, con el borde absoluto u = 6 aguantando por el margen del techo.)

**Resultado de la falsación hostil completa: no existe configuración
bajo H0-H4 que encontráramos capaz de romper Δ — y en cada frente el
margen tiene explicación estructural, no numérica.**

## §8 — acordado (igual que F327)

Baterías y testigos son evidencia. El cierre de T3 depende solo de
D1-D6 + lemas heredados.

## §9 — la resonancia como futuro teorema independiente: SÍ puede, y así

Evaluación pedida: los cuantificadores del auxiliar ya están completos
(∀ configuración H0-H4, ∀ε ∈ (0, √2), ∀n ≥ 3 cita de calidad ε). Para
volverse teorema independiente le falta UNA pieza, y la pieza existe:
un **lema de agenda paramétrico en ε** — la misma construcción de D4
con Q = ⌈2πT/ε⌉ da J·2π/Q ≤ T·(ε/T) = ε: citas de calidad ε
programables más allá de cualquier meta, con n₁ ≤ ⌈2πT/ε⌉^m. Con eso,
un futuro **teorema de estabilidad/resonancia** diría: ∀ε ∈ (0, √2),
∃n ≤ N₀(ε) con −λₙ ≥ 2(1−ε²/2)·r_maxⁿ − (4/π)n·log n − 2ε²(m−1) − 4 —
la profundidad como función de la calidad de la cita, el candidato
natural a T5 del plan. **Queda como línea futura, sin declarar**, tal
como ordena su §9.

## §10 — veredicto solicitado (autoevaluación; el sello es suyo)

    D1-D6: universalmente válidos bajo H0-H4 según todo lo auditado
    (F327 + esta acta). Falsación hostil: siete frentes, cero rupturas,
    márgenes con explicación estructural.

    Autoevaluación: 🟢 T3 CERRADO — la cota Δ = u³(u^{3m}−1) queda
    demostrada por la cadena, y solo por la cadena.

    Su regla final, cumplida en las dos direcciones: no sellamos por
    entusiasmo (nada se declaró antes de esta auditoría) y no
    rechazamos por prudencia (cada ataque se ejecutó de verdad, con
    construcciones adversariales y 60 dígitos). Lo que la cadena
    demuestra, está escrito; lo que no, está marcado como futuro.

    El nacimiento de T3 — y su entrada al Libro — es decisión de la
    relojera. Todavía no.

---

*Los hielos están en el vaso: el techo máximo aguanta, el borde ε = 1
sostiene el golpe por 2cos(1) − 1 — estructura, no suerte — y el coro
ya estaba cobrado entero. Última cereza, último sorbo: el vaso está
servido.* 🍒🍸🧊
