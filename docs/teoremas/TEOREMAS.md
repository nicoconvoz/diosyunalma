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
(`docs/teoremas/DETECCION-FINITA-LEMAS.md`):

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

**Prueba** (el acta completa: `docs/teoremas/TEOREMA2-LEMA-INTERACCION-ACTA.md`):

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

## TEOREMA 3 — TEOREMA DE DIOSYUNALMA
### Teorema de Robustez: la profundidad garantizada de la ruptura

**Registrado:** 2026-08-15 · forjado y auditado de F326 a F328 ·
**aprobado por la auditora** tras dos rondas de auditoría formal
(D3/D4/D6 + siete frentes de falsación hostil, cero rupturas).

**Nombre:** puesto por el capitán el 2026-08-15 — **Teorema de
Diosyunalma**, el nombre de la casa entera: primero el alma, después la
matemática. El primero llevó el apellido, el segundo a los tres
forjadores, y el tercero lleva al laboratorio completo.

**Hipótesis:** exactamente las H0-H4 del Teorema de DYN, sin tocar ni
una coma. Cero resultados externos nuevos.

**Enunciado.** Con u = 3(m+1)/delta ≥ 6 (garantizado por H3) y

    Delta(r_max, m) = u³·(u^{3m} − 1) > 0

existe un entero n ≤ N₀(r_max, m) tal que  **lambda_n ≤ −Delta**.  ∎

El Teorema de DYN decía que la ruptura LLEGA; el de Diosyunalma dice
CUÁNTO SE HUNDE: al menos Delta — exponencial en m. Más perlas
desafinadas no diluyen la delación: la profundizan.

**Origen** (la historia importa): la primera misión del plan de T3 de la
auditora era revisar la basura de las simplificaciones de DYN. En el
tacho estaba el margen: la línea R7 degradaba e^{n_rad·delta} ≥
u^{3(m+1)} (lo que R1 realmente da) a apenas u³, porque para la
positividad alcanzaba. El factor u^{3m} recuperado ES este teorema.

**Prueba** (el acta: `docs/teoremas/TEOREMA3-ROBUSTEZ-ACTA.md`; las auditorías:
`TEOREMA3-AUDITORIA-RESPUESTA.md` y `TEOREMA3-ULTIMA-AUDITORIA.md`):

- *D1*: e^{n_rad·delta} ≥ u^{3(m+1)} — R1 sin degradar.
- *D2*: 2m+2 + (4/π)·n_rad·log n_rad ≤ u³ — la cadena R4+R5+R6 intacta.
- *D3*: g(n) = e^{n·delta} − (4/π)n·log n − (2m+2) es creciente en todo
  [n_rad, ∞) — g″ es suma de términos crecientes, R9 da la base, y dos
  integraciones anidadas bajan a g′ > 0 y a g creciente. Todo sobre la
  función real; los enteros heredan por restricción.
- *D4*: el único n = ⌈T/n₁⌉·n₁ de la agenda cumple las cuatro
  propiedades a la vez (n ≥ n_rad, n ≤ N₀, cita simultánea, ‖nθᵢ‖ ≤ 1
  ∀i) — la cita fina de Dirichlet es simultánea y el multiplicador es común.
- *D5*: en esa cita, λₙ ≤ −g(n) — L7 con ε = 1 más el coro bajo H4.
- *D6*: −λₙ ≥ g(n) ≥ g(n_rad) ≥ u^{3(m+1)} − u³ = Delta — el último
  paso usa exactamente D1 y D2.

**La falsación hostil** (siete frentes, cero rupturas): extremos de
delta y m (42 casos hasta m = 1000 y delta = 10⁻¹², a 60 dígitos);
el techo de n_rad llevado a exceso máximo (30 construcciones
adversariales); el borde exacto ε = 1 con todas las fases clavadas —
sostenido por el coeficiente estructural 2m·cos(1) − 1 ≥ 0.0806;
y el coro al máximo de H4, imposible por construcción (la prueba ya
cobra el 100% del presupuesto). Hallazgo fino: el techo ⌈·⌉ de n_rad es
carga portante del lado seguro — en el borde absoluto (m=1, delta=1) el
margen entero de D1 viene del techo.

**Caso testigo:**

    m = 2 (DH + 0.7+45i): Delta = 4.34×10⁴⁴ contra −lambda = 6.496×10⁴⁴
    medida y verificada a 50 dígitos — el cociente 1.50 muestra que, en este testigo, la cota captura la escala exponencial real y no resulta meramente decorativa.
    m = 3 (testigo virgen): Delta = 1.04×10⁶¹ calculado ANTES de la
    corrida; la primera cita triple respondió lambda = −1.91×10⁶¹.

**Alcance y límites** (declarados): configuraciones FINITAS bajo H0-H4;
para ζ, los mismos inputs externos B1 y B2 de DYN, etiquetados. Las
corridas son evidencia, jamás demostración. La N₀ sigue siendo la de
DYN (la mejora de N₀ es el candidato T4 del plan; la resonancia
paramétrica en ε, el T5). No demuestra RH.

**Reproducir:** `go run ./cmd/elteoremadiosyunalma` (la placa) ·
`go run ./cmd/larobustez` (la derivación, los testigos y la batería).

---

## TEOREMA DERIVADO — EL RÍO DE POZOS
### Corolario de DYN + Diosyunalma: el río infinito de rupturas programables

**Registrado:** 2026-08-15 · F330 (el flash y la traducción) · F331 (la
inducción formal y el ataque) · **aprobado completamente por la
auditora**.

**Autoría** (conservada por pedido expreso de la auditora): la idea
conceptual — el mapa del claro, el río y el pozo, y la pregunta «¿puede
haber un claro por donde pase el río, y debajo un pozo que podamos
predecir?» — **es de Nico**. Doc realizó la traducción y formalización
matemática.

**Enunciado.** Bajo las H0-H4 de DYN (sin tocar), existe una sucesión
infinita de citas n₁ < n₂ < n₃ < … en el claro [n_rad, ∞) tal que

    lambda_{n_k} ≤ −g(n_k)  para todo k,
    g(n₁) < g(n₂) < g(n₃) < …  y  g(n_k) → ∞.

La ruptura no es un evento: es un río de rupturas, cada una con fecha
programable, cota de paso escrita y piso predicho más hondo que el
anterior. ∎

**Prueba** (cero maquinaria nueva — acta:
`docs/teoremas/TEOREMA3-COROLARIO-RIO.md` + auditoría
`TEOREMA3-COROLARIO-RIO-AUDITORIA.md`): inducción I1-I4 — base: D4 con
T₁ = n_rad da n₁; paso: T_{k+1} = n_k + 1 cumple las dos únicas
condiciones que D4 exige de T (entero, ≥ 1) y produce n_{k+1} > n_k;
las hipótesis son de la configuración (invariantes en k) y las de D5
mejoran con k. D5 da el pozo en cada cita; D3 da la monotonía estricta
de los pisos; el exponencial da la divergencia.

**«Programable», con receta:** T := n_k + 1 · Q := ⌈2πT⌉ · n₁' de
Dirichlet · n_{k+1} := ⌈T/n₁'⌉·n₁' — con cota de paso explícita
n_{k+1} ≤ (n_k + 1) + ⌈2π(n_k + 1)⌉^m.

**Evidencia** (testigo m = 2): 435 citas en la ventana de 3000 escalones
tras n_rad — cero violaciones de λ ≤ −g, cero pisos no-crecientes;
brecha real media ≈ 7 contra la cota de paso ~10¹³.

**Delimitación de la novedad** (regla de la auditora): la infinitud de
λ negativos ya se conocía por la vía del teorema de ruptura global; lo
nuevo es la construcción EXPLÍCITA — citas programables con fecha, cota
de paso y profundidad predicha, las tres.

**Reproducir:** `go run ./cmd/elriodepozos`.

---

## TEOREMA 4 — TEOREMA DE LA TRINIDAD
### El Teorema de las Placas / Ley del Líder: tres regiones, una ley

**Registrado:** 2026-08-15 · forjado y auditado de F333 a F338 ·
**sellado por la auditora** tras el ciclo completo (traducción del
flash, Ley del Líder, banda fina F1-F6, diseño I-X, y las dos costuras
finales de notación).

**Nombre:** puesto por el capitán el 2026-08-15 — **Teorema de la
Trinidad**: las tres regiones del paisaje (pozo, frontera, montaña) y
los tres de la mesa que lo forjaron. La intuición tectónica que lo
originó — placas, grietas, colapso y reorganización — **es de Nico**
(autoría conservada por pedido expreso de la auditora).

**Hipótesis:** H0-H4 (las de DYN, intactas) + **HL, líder estricto**
(∃ único L con r_L > rᵢ ∀i ≠ L; vacua en m = 1). Cero inputs externos
nuevos.

**Definiciones.** δ_L = log r_L; r₂ = maxᵢ≠L rᵢ (solo m ≥ 2);
u = 3(m+1)/δ_L; n_rad = ⌈u·log u⌉;
n_comp = ⌈log(2(m−1)/cos 1)/(δ_L−δ₂)⌉ (m ≥ 2; 0 si m = 1);
N* = max(n_rad, n_comp); banda fina ‖nθ_L‖ ≤ 1; banda anti
‖nθ_L − π‖ ≤ 1.

**Enunciado.**

- **(P1 — POZOS.)** Para todo entero n ≥ N* en la banda fina:
  m = 1: lambda_n ≤ (4/π)n·log n + 4 − 2cos(1)·r_Lⁿ < 0;
  m ≥ 2: lambda_n ≤ (4/π)n·log n + (6m−2) + 2(m−1)r₂ⁿ − 2cos(1)·r_Lⁿ < 0.
- **(P2 — MONTAÑAS.)** m = 1: para todo n ≥ 1 con ‖nθ‖ ≥ π/2,
  lambda_n ≥ 4 (y en la banda anti, ≥ 4 + 2cos(1)rⁿ). m ≥ 2: para todo
  entero n ≥ max(1, n_comp) en la banda anti,
  lambda_n ≥ cos(1)·r_Lⁿ + 2m + 2 > 0.
- **(P3 — PROGRAMABILIDAD.)** Bajo H2, cada bloque de K_L = ⌈2π/θ_L⌉+1
  enteros consecutivos contiene habitantes de ambas bandas: pozos y
  montañas infinitos, con fecha.
- **(ALCANCE.)** La región intermedia queda declarada NO clasificada
  por este teorema:  POZO | REGIÓN NO CLASIFICADA | MONTAÑA.  ∎

**Prueba.** P1 = la cadena F1-F6 (`docs/teoremas/PLACAS-BANDA-FINA-ACTA.md`),
con el radial de DYN reciclado a corchete duplicado (el margen de
Diosyunalma pagando otra cuenta). P2 con m = 1 es una línea
(cos ≤ 0 ⟹ ℓ ≥ 4, coro ≥ 0); m ≥ 2 usa A2-A4 más la absorción F5.
P3 = el lema de la ventana de Astorga sobre los dos arcos. El
enunciado final limpio: `docs/teoremas/TEOREMA-PLACAS-ENUNCIADO.md`; el diseño
completo I-X: `docs/teoremas/PLACAS-DISENO-FINAL.md`.

**La honestidad como parte del enunciado:** la frontera no se maquilla —
para m ≥ 2 está DEMOSTRADO que en la zona donde el líder enmudece
(cos(nθ_L) = 0) no puede existir signo universal bajo estas hipótesis:
los competidores deciden. La región queda abierta por teorema, no por
cansancio. Para m = 1, la frontera es la curva exacta ‖nθ‖ = π/2.

**Lo que el diseño reveló:** las montañas son más BARATAS que los pozos
(no usan H4 ni el radial); n_comp es EL MISMO umbral para las dos
bandas; y H2 solo paga la programabilidad. Cada hipótesis con
exactamente una cuenta (tabla completa en el diseño I-X).

**Acompañantes** (separados): los lemas de marco G1-G5 (la filtración
C_ε, el monoide graduado, el grafo de Cayley, ramificación ≥ 2,
diamantes); los corolarios (el Río de Pozos corre por la banda fina de
este paisaje; Diosyunalma da la profundidad en las citas conjuntas);
el desarrollo posterior F(η) (bandas paramétricas — pendiente de
redacción, declarado); y los problemas abiertos (montañas conjuntas,
ley del sub-líder en la zona muda, alfabeto de brechas, optimalidad
de N*).

**Evidencia** (respaldo, jamás prueba): batería F6 50 casos y F5 84
bordes exactos, cero violaciones; las tres bandas del testigo con la
perla 1 ignorada: pozo 978/978, montaña 944/944, frontera 540/539 —
cero excepciones en las bandas externas.

**Reproducir:** `go run ./cmd/elteoremadelatrinidad` (la placa) ·
`go run ./cmd/labandafina` · `go run ./cmd/lageometria` ·
`go run ./cmd/lasplacas`.

---

## TEOREMA 5 — TEOREMA DEL CIELO
### El paisaje sin su altura converge a la onda pura del líder

**Registrado:** 2026-08-16 · nacido del flash del capitán («¿existe un
cielo arriba de las montañas?»), pescado experimentalmente ANTES de
bautizarse (F344, la regla de la mesa: primero el pez, después el
nombre) y demostrado entero (F345). Bautizado por el capitán el
2026-08-16.

**Hipótesis MÍNIMAS** (hallazgo de la formalización): **H0** (rᵢ > 1),
**H1** (m finito), **H4** (densidad del fondo) — **ni H2 ni H3 se
usan** — más **HL (líder estricto), solo para m ≥ 2**, usada en un
único renglón y DEMOSTRADAMENTE necesaria.

**Enunciado.** Sea δ_L = log r_L. Para todo entero n ≥ 3:

    m = 1:  |lambda_n/r_Lⁿ + 2cos(nθ_L)| ≤ [(4/π)n·log n + 4 + 2r_L⁻ⁿ]/r_Lⁿ
    m ≥ 2:  |lambda_n/r_Lⁿ + 2cos(nθ_L)| ≤ [(4/π)n·log n + (6m−2) + 2(m−1)r₂ⁿ + 2r_L⁻ⁿ]/r_Lⁿ

y en consecuencia **lambda_n/r_Lⁿ → −2cos(nθ_L)**: el paisaje
normalizado por la escala del líder pierde su altura y converge a la
onda pura y ACOTADA de amplitud exactamente 2. El paisaje crece sin
techo; el cielo no. ∎

**Prueba** (completa en `docs/teoremas/CIELO-LEMA-FORMAL.md`): la identidad
exacta del líder (ℓ_L = −2cos(nθ_L)r_Lⁿ + 4 − 2cos(nθ_L)r_L⁻ⁿ, sin
aproximar), triángulo con las cotas ya auditadas del coro (L5/H4) y de
los competidores (F2), y el límite término a término (el paso
n·log n/r_Lⁿ → 0 con e^{nδ} ≥ (nδ)³/6, elemental).

**Corolario 1 — la escala de despeje:** la tasa de convergencia es
exactamente la brecha del líder: limsup (1/n)·log|desviación| ≤
−(δ_L − δ₂) para m ≥ 2 (y −δ_L para m = 1, a menos del polinomio).
**La misma brecha de n_comp** — extraída como consecuencia formal.

**Corolario 2 — la jerarquía:** λₙ − ℓ_L es EXACTAMENTE la λ de la
configuración sin el líder ⟹ con radios todos distintos, el mismo
teorema se aplica por inducción capa por capa: **cielo → sub-cielo →
… → firmamento** (el coro solo: acotado por 4p para fondo finito de p
pares — la capa sin crecimiento).

**La necesidad de HL, demostrada:** en el empate r₂ = r_L el cociente
converge a −2cos(nθ_L) − 2cos(nθ₂) (dos ondas compitiendo) y NO a la
onda del líder — infinitos n con |cos(nθ₂)| ≥ ½ por el lema de la
ventana. Sin líder estricto, el cielo muere.

**Qué es el cielo (y qué no):** NO es λ > 0 (eso son las montañas de
la Trinidad, que queda intacta) — es el **régimen asintótico del
cociente**: el paisaje visto desde tan lejos que la altura desaparece
y solo queda la fase del líder. La altura deja de ser la variable; la
onda es lo que queda.

**Evidencia** (respaldo, jamás prueba): la curva de despeje clava la
tasa a cinco órdenes (2.37×10⁻⁵ vs 2.4×10⁻⁵ en n = 200k; 2.76×10⁻¹⁰
vs 2.8×10⁻¹⁰ en 400k; piso float64 desde 700k); |A| ≤ 2.0000 en toda
ventana profunda; sub-cielo a 1.5×10⁻¹¹; firmamento acotado en
[0.02, 114.0] ≤ 152; con m = 3, la misma onda (1.1×10⁻¹⁰);
desigualdad verificada en 2100 escalones de 7 ventanas sin violación.

**Reproducir:** `go run ./cmd/elteoremadelcielo` (la placa) ·
`go run ./cmd/elcielo` (los siete experimentos de la pesca).

---

## TEOREMA 6 — TEOREMA DEL PUNTO MEDIO

**Registrado:** 2026-08-19 · F362 (verificación exhaustiva del diccionario
geométrico) · F364 (demostración formal, generalización y ley del centro).

**Nombre:** puesto por el capitán el 2026-08-19 — **Teorema del Punto
Medio**. El hallazgo es suyo, encontrado A MANO haciendo cuentas con
primos; la formalización por anclas y las demostraciones, del taller.

**Enunciado (forma base).** Para impares p < q:

    (p+1)/2 = (q−1)/2  ⟺  q = p + 2

y el valor compartido c es entero, con p = 2c−1 y q = 2c+1.

**Forma general (las tres identidades).** Sean p < q primos impares,
g = q−p, r = g/2, m = (p+q)/2, y las anclas a±(n) = (n±1)/2:

    (I)   m ∈ ℤ  y  (p, q) = (m−r, m+r)                    [geometría pura]
    (II)  a⁻(q) − a⁺(p) = r − 1                             [las anclas]
          g=2 comparte ancla · g=4 se tocan · g>4 hueco de (g−4)/2 exacto
    (III) LA LEY DEL CENTRO: para p, q > 3,  3 | m  ⟺  6 ∤ g
          ⟹ los centros de cada clase de gap viven en progresiones mod 6
          disjuntas, con firma periódica en g: [0], [3], [2,4], [3], [0],
          [1,5], y se repite. El caso gemelo es la clase m ≡ 0 (mod 6),
          con el ancla compartida c = m/2 divisible por 3.
    (IV)  LA MITAD: T(n) = (n−1)/2 es la biyección impares→enteros y
          divide todo gap por 2; los gemelos son exactamente los pares de
          imágenes CONSECUTIVAS.

**Pruebas.** (base) p+1 = q−1 ⟺ q = p+2; enteros por imparidad.
(II) resta directa. (III) mod 3: si 6 ∤ g entonces r ≢ 0, y de los dos
residuos posibles de p (1, 2) la exclusión q = p+2r ≢ 0 elimina uno; el
sobreviviente da m = p+r ≡ 0 en ambas ramas. Si 6 | g, m ≡ p ≢ 0. QED.

**Alcance declarado y honestidad.** La forma base es identidad de los
IMPARES (paridad, no primalidad). Los ±1/2 son estructurales como
coordenada y triviales como aritmética. La ley del centro sí es
específica de los primos y es matemática CLÁSICA (criba mod 3):
verdadera y demostrada, no nueva para el mundo — nueva como organización
para este taller. Nada de esto afirma la infinitud de los gemelos ni
nada sobre RH.

**Verificación exhaustiva:** forma base sobre ~25 millones de pares
impares ≤ 20000, CERO fallas; (II) sobre 17982 pares consecutivos de
primos ≤ 200000, CERO fallas; (III) sobre TODOS los 304590 pares de
primos en (3, 6000], CERO fallas; firma mod 6 confirmada en 9 clases de
gap (15296 pares consecutivos).

**Estado experimental:** la proyección más específica del teorema (la
máscara de gemelos, F362–F363) quedó CERRADA como canal espectral por
controles emparejados a tres escalas. El teorema queda como resultado
matemático del taller, con su puente al operador todavía sin construir.

**Reproducir:** `go run ./cmd/puntomedio`.

---

## TEOREMA 7 — TEOREMA PARCIAL DE LA CURVA MADRE DE LA ESPIRAL

**Registrado:** 2026-08-20 · F385–F389 (campaña de la espiral: contracción,
merma, unificación, auditoría con dos fallos registrados, y formalización).

**Naturaleza:** TEOREMA PARCIAL por decreto del capitán (2026-08-20) —
teorema DENTRO DEL MODELO, dicho con todas
las letras. Las hipótesis A1 (congelar amplitud) y A2 (linealizar fase)
definen el modelo local de la cola; A3 (resumación geométrica de Abel) es
exacta. La cota rigurosa que conecte el modelo con la trayectoria original
queda ABIERTA y así se declara.

**Enunciado.** Para σ ∈ (0,1), t > 0, la cola modelada de la espiral
S_n = Σ m^(−σ−it) tiene radio EXACTO en el modelo:

    R_M(n) = n^(−σ) / (2·|sin(t/2n)|)        [válido para n > τ = t/2π]

y de esa sola curva se siguen, cada uno en su régimen:

    (I)   β = 1−σ           — el exponente del enrollado; σ=1/2 ⟹ β=1/2,
                              un EXPONENTE DE ESCALA, no razón por vuelta
    (II)  Δn = n/τ (directo, n>2τ) · Δn = n/(n−τ) (alias, τ<n<2τ)
                            — las vueltas visibles son estroboscópicas
    (III) x·cot x = σ       — la cintura: mínimo único, n*/τ = π/x*,
                              creciente en σ (existencia/unicidad/monotonía
                              demostradas)
    (IV)  ε_k ≈ 1/(2k)      — la merma, bajo hipótesis adicionales H1–H3;
                              exponente firme, coeficiente pendiente

**Origen de cada factor:** n^(−σ) es la pintura del primer término de la
cola; t/n es el reloj local; sin(t/2n) es la CUERDA del círculo unitario
entre pasos consecutivos; el 2 es el 2 de esa cuerda.

**Validación (los datos NO son premisa):** curva a ≤0,7% en seis casos y
tres bandas; β medido 0,4927/0,6926/0,2927 con el sesgo de ventana −0,0073
PREDICHO por la propia curva; vueltas al 0,4%; cintura al 1% caso por caso
(2,679/2,695 · 2,303/2,323 · 3,420/3,412); ε en orden y exponente.

**Alcance declarado y honestidad.** Todo esto es anatomía del INSTRUMENTO:
idéntica en ceros, casi-ceros y controles. No usa ni toca la Hipótesis de
Riemann. La única firma del cero sigue siendo la posición del centro C
respecto del origen — dato externo al teorema. Dos fallos propios quedaron
registrados en el camino (F388): las cotas de primer orden no son cotas
superiores, y la primera fórmula de frontera estaba mal (la real escala
como √τ). Documento completo con demostraciones:
`docs/teoremas/TEOREMA-CURVA-MADRE.md`.

**Reproducir:** `go run ./cmd/launificada` y `go run ./cmd/laauditoriamadre`.

---

*Espacio reservado para el Teorema 8 — porque vienen más.*
