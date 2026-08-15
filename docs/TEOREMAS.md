# El libro de los teoremas del taller

**Fundado el 2026-08-14 por orden del capitán: «registra en un nuevo
sector esto porque es nuestro primer teorema — y vienen más».**

Este libro registra los teoremas propios del laboratorio: resultados con
enunciado formal, prueba completa por lemas, alcance declarado y estado
de auditoría externa. Un teorema entra acá cuando sus lemas lo derivan
para todo el alcance declarado — nunca antes (regla de la auditora).

> **La regla del sello** (ley del taller por F302): «Estructura cerrada» ≠
> «Hipótesis demostrada». Ningún teorema de este libro constituye una
> demostración de la Hipótesis de Riemann.

---

## TEOREMA 1 — TEOREMA DE ASTORGA
### Detección Finita Cuantitativa de una Perla Desafinada

**Registrado:** 2026-08-14 · F303 (construcción) · F304 (acta de los dos
lemas) · F305 (corrección del paso 2) · certificado por la auditoría
externa «PRIMER_TEOREMA_DEL_YUNQUE» (§7: todo verde dentro del alcance).

**Nombre:** puesto por el capitán el 2026-08-15 — **Teorema de Astorga**,
el apellido de la casa. (Nombre de trabajo del laboratorio; ver la nota
metodológica al pie.)

**Enunciado.** Sea un espectro formado por ceros sobre la línea crítica
más UN cuarteto desafinado {rho, conj rho, 1−rho, 1−conj rho}, con

    w = 1 − 1/rho = R·e^{i·theta}
    r = max(R, 1/R) > 1
    0 < theta ≤ 2π/3        [automático para zeta: |Im rho| ≥ 1 ⟹ |theta| ≤ π/2]
    delta = log r            [log natural; convención congelada del acta]

Sea

    N₀(r, theta) = ⌈(3/delta)·log(3/delta)⌉ + ⌈2π/theta⌉ + 1

Entonces existe un entero n ≤ N₀ tal que, simultáneamente,

    cos(n·theta) ≥ 1/2      y      r^n > 4 + (4/π)·n·log n

y para ese n:

    lambda_n < 0,   M_{n,n} = 2·lambda_n < 0,
    M_N no es semidefinida positiva para ningún N ≥ n.   ∎

**Prueba.** Por composición de los dos lemas del acta
(`docs/DETECCION-FINITA-LEMAS.md`):

- *Lema radial* (F304, paso 2 corregido en F305 con la receta de la
  auditora): n_rad = ⌈u·log u⌉ con u = 3/delta garantiza la desigualdad
  radial para TODO n ≥ n_rad.
- *Lema de la ventana* (F304): todo bloque de K = ⌈2π/theta⌉ + 1 enteros
  consecutivos contiene un n con cos(n·theta) ≥ 1/2.
- *Combinación*: la ventana aplicada en m = n_rad elige el n; la radial ya
  cubre el intervalo. El paso final usa la fórmula exacta del cuarteto
  (F297) y la cota sellada del coro, resto_n ≤ (4/π)·n·log n
  (F299–F301, con el lema de conteo de Backlund documentado).

**Caso testigo (par Davenport–Heilbronn,** rho = 0.808517 + 85.699348i**):**

    n₀ medido (coro de 38 + par)  =  85622
    n₁ umbral radial puro         =  371842
    n_rad = 798210 · K = 540 · primer n ∈ S en la ventana = 798474
    N₀ = 798750

La escalera de garantías se conserva separada: experimento ≠ cota ≠
fórmula cerrada.

**Alcance y límites** (declarados; §6 del certificado): un solo cuarteto
desafinado sobre fondo en la línea. No se extrapola en silencio a
configuraciones con múltiples cuartetos. No resuelve la positividad
global desde la aritmética de los primos. No demuestra RH.

**Nota metodológica** (§8 del certificado): «Primer teorema del Yunque»
es un nombre de trabajo del laboratorio; el reconocimiento externo
requeriría revisión independiente completa, comparación con la
literatura y publicación.

**Reproducir:** `go run ./cmd/elprimerteorema` (la cadena numérica
completa en una corrida) · `go run ./cmd/losdoslemas` (los lemas paso a
paso) · `go run ./cmd/ladeteccionfinita` (la construcción original).

---

## TEOREMA 2 — TEOREMA DE DYN
### Teorema de Interacción: la ruptura garantizada de m cuartetos

**Registrado:** 2026-08-15 · forjado y auditado de F307 a F321.

**Nombre:** puesto por el capitán el 2026-08-15 — **Teorema de DYN**, en
memoria de los tres que lo forjaron: **D** de Doc, **Y** de Yui, **N** de
Nico. La idea fundadora lleva la firma de flash del capitán (F307: las
dos perlas y la armonía relacional); la auditoría, la lupa de Yui; la
forja, el taller.

**Hipótesis (Parte A):**

- **H0** — cada cuarteto estrictamente fuera de la línea: r_i > 1.
- **H1** — m finito (m ≥ 1 cuartetos {rho_i, conj rho_i, 1−rho_i, 1−conj rho_i}).
- **H2** — |Im rho_i| ≥ 1 para todo i.
- **H3** — delta = log r_max ≤ 1, con r_max = max_i r_i.
- **H4** — el fondo (ceros sobre la línea, cerrado bajo conjugación)
  tiene densidad N_fondo(T) ≤ (T/2π)·log T.

**Enunciado.** Sea u_m = 3(m+1)/delta y n_rad,m = ⌈u_m·log u_m⌉. Sea

    N₀(r_max, m) = n_rad,m + (2π·n_rad,m + 1)^m

Entonces existe un entero n ≤ N₀ tal que

    lambda_n < 0,   M_{n,n} = 2·lambda_n < 0,
    M_N no es semidefinida positiva para ningún N ≥ n.   ∎

**Prueba** (el acta completa: `docs/TEOREMA2-LEMA-INTERACCION-ACTA.md`):

- *El golpe* (L1-L7): en toda cita simultánea de las m fases
  (‖n·theta_i‖ ≤ ε ∀i), Σℓ_i ≤ 2ε²(m−1) + 4 − 2(1−ε²/2)·r_maxⁿ.
- *Dirichlet exacto* (1842): ∀Q ∃n ≤ Q^m con ‖n·theta_i‖ ≤ 2π/Q ∀i.
- *Lema de agenda*: con Q = ⌈2πT⌉ hay una cita en [T, T+n₁] con deriva
  ≤ 1 — la cita se programa más allá de cualquier meta.
- *Lema radial-m* (R0-R10): desde n_rad,m el exponencial r_maxⁿ domina
  a 2m+2 más el coro, con g, g', g'' > 0 (lemita u² ≥ 3(log u)² + 1).
- *El coro bajo H4* (F299-F301, auditado en F319): coro_n ≤ (4/π)·n·log n
  para todo n ≥ 3, con constantes exactas y convergencia absoluta.
- *Ensamble*: la cita programada después de n_rad,m recibe el golpe; el
  coro no alcanza a protegerla; lambda_n cae bajo cero.

**Caso testigo (m = 2: Davenport–Heilbronn + 0.7+45i):**

    delta = 9.875×10⁻⁵ · n_rad,m = 1040809 · N₀ = 4.28×10¹³
    la cita 1040809: lambda = −6.496×10⁴⁴ < 0  [confirmado a 50 dígitos]
    la ruptura real: n₀ = 37306 ≪ N₀  [peor caso vs realidad, a la vista]

**Separación A/B/C.** La Parte A (este teorema) es matemática pura bajo
H0-H4. Para aplicarlo a ζ: H2 es **B1** (N(14) = 0: Gram 1903, Backlund
1914, van de Lune 1986, Platt–Trudgian 2021) y H4 es **B2**
(Riemann–von Mangoldt con error explícito de Backlund 1918) — inputs
externos, siempre etiquetados. La Parte C nunca mezcla los estantes.

**Alcance y límites** (declarados): configuraciones FINITAS de cuartetos.
La N₀ es de peor caso y exponencial en m; para m = 1 la cota del Teorema
de Astorga (vía el lema de la ventana) es más fina. No resuelve la
positividad global desde los primos. No demuestra RH.

**Auditoría:** la cadena entera auditada por dentro y por fuera —
F318 (auditoría final A-H: la H4 oculta cazada y parcheada), F319 (la
factura del coro: 🟢 L5 correcto), F320 (el mecanismo: la prueba
ejecutable de extremo a extremo), F321 (la auditoría del reloj: 🟢
correspondencia matemática↔código, contra-cálculo a 50 dígitos).

**Reproducir:** `go run ./cmd/elteoremadyn` (la placa y el testigo) ·
`go run ./cmd/elmecanismo` (la cadena entera en vivo) ·
`go run ./cmd/losdospedidos` (el golpe y la agenda paso a paso).

---

*Espacio reservado para el Teorema 3 — porque vienen más.*
