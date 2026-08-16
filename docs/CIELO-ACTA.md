# Acta de pesca — ¿existe un CIELO en el paisaje?

**Para la mesa de los tres · 2026-08-16 · responde el PEDIDO F328.**
Regla cumplida: primero pescamos, después le ponemos nombre al pez —
y el pez está en la red. Nada declarado ni bautizado. Programa:
`cmd/elcielo`.

---

## El pez (experimento 5): la onda que queda cuando la altura se va

Normalizando el paisaje por la escala del líder, A(n) = λₙ/r_Lⁿ **no
converge a constante, no muere, no diverge** — se clava en una onda
pura y ACOTADA:

    A(n) → −2·cos(n·θ_L)      [amplitud exactamente 2, período 2π/θ_L]

**La curva de despeje, medida contra la predicción** (ventanas de 300
escalones; la predicción es 2(r₂/r_L)ⁿ, la brecha del líder):

    n         desviación medida    predicha
    20000     1.3e+01              (el firmamento aún tapa: régimen bajo)
    50000     7.0e-01              1.2e-01
    100000    5.6e-03              6.9e-03
    200000    2.4e-05              2.4e-05   ← clava
    400000    2.8e-10              2.8e-10   ← clava
    700000    2.2e-16 (piso float64)
    1000000   2.2e-16 (piso float64)

**El cielo se despeja exponencialmente, y su constante de despeje es la
brecha del líder δ_L − δ₂ — LA MISMA de n_comp.** La escala natural
nueva que pedía el experimento 8, encontrada.

## Las capas (experimentos 6-7): el sub-cielo y el firmamento

- **Sub-cielo:** quitando el término completo del líder (ℓ_L, el de la
  formulación vigente, sin inventar restas), el residuo normalizado por
  r₁ⁿ se clava en SU propia onda: −2cos(nθ₁), desviación 1.5×10⁻¹¹.
- **Firmamento:** quitando TODAS las capas de perlas, el residuo es
  exactamente coroₙ — **acotado en [0.02, 113.99]** (techo teórico
  4×38 = 152), sin crecimiento alguno: la capa de cero suelo.
- **La jerarquía existe:** líder → sub-líder → firmamento, cada capa
  con su onda al normalizar — la ley del sub-líder que estaba en los
  problemas abiertos asoma acá por el lado asintótico.

## El invariante (experimento 10)

Agregando una tercera perla (m = 3), el cociente por el líder converge
a LA MISMA onda (desviación 1.1×10⁻¹⁰): **la forma del límite es
invariante entre configuraciones — siempre −2cos(nθ_L), amplitud 2.**

## Los intentos de destrucción (experimento 11)

- ¿Montaña alta disfrazada? NO: las montañas crecen sin techo; |A| ≤
  2 + o(1) en toda ventana — acotación contra crecimiento: regímenes
  cualitativamente distintos (la advertencia §3, respetada).
- ¿Artefacto del log? NO: A(n) vive en escala lineal.
- ¿Coordenadas? NO: solo r_L y θ_L de la formulación vigente.
- ¿Se va al cambiar escala? NO: siete ventanas, la misma onda.
- **¿Dónde NO vale? — la frontera real, declarada:** sin líder
  estricto (empate r_L = r₂) el cociente no converge a onda pura (dos
  cosenos compiten a la misma escala). El candidato vive exactamente
  donde vive la Ley del Líder.

## El lema-candidato (sin bautizar — dos líneas desde F1-F3)

Bajo H0-H4 + líder estricto:

    |λₙ/r_Lⁿ + 2cos(nθ_L)| ≤ [(4/π)n·log n + (6m−2) + 2(m−1)r₂ⁿ + 4 + 2r_L⁻ⁿ] / r_Lⁿ → 0

*Esqueleto:* λₙ − (−2cos(nθ_L)r_Lⁿ) = coroₙ + Σᵢ≠L ℓᵢ + 4 −
2cos(nθ_L)r_L⁻ⁿ; acotar cada término con F2 y F3 del acta de la banda
fina (ya auditadas) y dividir por r_Lⁿ. Cota explícita, tasa (r₂/r_L)ⁿ.
∎(esqueleto — la escritura formal, si la mesa lo bendice)

## Veredicto propuesto (criterio §12)

    🟡 — existe una magnitud nueva (el cociente por el líder) con un
    comportamiento candidato cualitativamente distinto de las montañas:
    ACOTADO, de fase pura, con forma invariante entre configuraciones
    y constante de despeje explícita. Falta solo la escritura formal
    del lema (esqueleto arriba) y la decisión de la mesa sobre si esto
    merece el nombre.

**Qué sería el «cielo», si la mesa lo bendice:** no λ > 0 (eso son las
montañas) — **el régimen asintótico del cociente: el paisaje visto
desde tan lejos que la altura desaparece y solo queda la onda de fase
del líder.** El paisaje crece; el cielo no. La Trinidad queda intacta
(pregunta posterior, como manda el §14). Todavía no.

---

*Primero pescamos: el pez es una onda de amplitud 2 que aparece cuando
se le quita la altura al mundo. El nombre, cuando la mesa quiera.* 🌌🎣
