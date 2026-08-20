# Acta de la Frontera — Fase X

**Respuesta al Telar Fase X (Auditoría 43) · 2026-08-18**

Su pregunta era una sola: **¿cuánta de la mejora de la frustración sobrevive
cuando se exige participación comparable?**

Respuesta corta: **sobrevive, y con creces — pero no por lo que creíamos.** Hay
rigidez real (gana su H2), y en el camino apareció por qué la Fase IX vio
localización: **no venía de la frustración, venía de la codificación.**

Se reproduce con `go run ./cmd/lafrontera`.
Lámina: `galeria/laminas/10-el-telar/la-frontera.svg`.

---

## 1 · Los dos rieles

**Riel 1, declarado antes de la primera medición** (su §9). El control S0 tiene
PR/N = 0,105:

- **COMPARABLE**: PR/N ≥ 0,090 · **degradada**: 0,070–0,090 · **ATRAPADA**: < 0,070
- **Regla de veredicto, fijada también antes:** H2 gana sólo si algún brazo
  COMPARABLE baja Σ²(10) más de 2 sigmas por debajo de S0.

**Riel 2, que su §4 pedía y yo omití en la primera corrida.** Un brazo es
admisible sólo si conserva **≥ 80 % de los niveles vivos del control** (≥ 149 de
187). Sin eso, Σ² se mide sobre otro recorte y no es comparable.

Lo anoto porque la primera corrida no lo tenía, y el resultado cambia: sin riel,
el «mejor» brazo daba **Σ²(10) = 0,988 — por debajo del piso GUE — pero con 87
niveles de 187.** La banda vaciándose. Es *exactamente* el error que obligó a
retractar la Fase VI. Con el riel puesto, **siete brazos quedan descartados** y
el ganador es otro. El hallazgo sobrevive; ese número no.

---

## 2 · El veredicto: gana H2

Familia B — el dial limpio, cada enlace negativo con probabilidad q:

| brazo | vivos | Σ²(10) | α | PR/N | dens− | frustr | |
|---|---:|---:|---:|---:|---:|---:|---|
| S0 control | 187 | 18,335 | 1,717 | 0,105 | 0 | 0 | |
| B · q=0,005 | 184 | 14,551 ± 1,745 | 1,704 | 0,129 | 0,005 | 0,015 | admisible |
| B · q=0,010 | 178 | 11,334 ± 0,965 | 1,697 | 0,156 | 0,010 | 0,028 | admisible |
| **B · q=0,020** | **170** | **6,324 ± 0,781** | 1,590 | **0,189** | 0,020 | 0,057 | **admisible** |
| B · q=0,050 | 148 | 2,270 ± 0,464 | 1,055 | 0,226 | 0,051 | 0,140 | inadmisible |
| B · q=0,100 | 122 | 1,389 ± 0,323 | 0,725 | 0,247 | 0,100 | 0,242 | inadmisible |
| B · q=0,200 | 97 | 1,080 ± 0,256 | 0,414 | 0,199 | 0,199 | 0,392 | inadmisible |
| B · q=0,350 | 87 | 0,988 ± 0,275 | 0,349 | 0,266 | 0,352 | 0,491 | inadmisible |

**El mejor brazo admisible: Σ²(10) = 6,324 contra 18,335 del control — 15,4
sigmas — con PR/N = 0,189 contra 0,105, es decir la participación SUBIÓ, y
conservando 170 de 187 niveles.**

Es la primera vez en toda la cadena de auditorías que Σ² baja **sin** que la
participación caiga. No es esconder los estados: es ordenarlos. Su §16 pedía
demostrar esa flecha antes de cualquier otra cosa, y está demostrada.

**Su H3, con honestidad.** El dial tiene un mínimo interior en q ≈ 0,35, o sea
«más frustración = mejor» es falso. Pero ese mínimo vive **entero en la zona
inadmisible**: entre los brazos que pasan el riel, Σ² cae monótona. Así que
**H3 no queda establecida sobre filas admisibles**, y no se afirma.

---

## 3 · El hallazgo: la frustración NO determina la respuesta

Su §6 preguntaba si igualar un solo número alcanza. **No alcanza.** Tres familias
a la MISMA frustración de triángulos (≈0,22):

| familia | Σ²(10) | PR/N | vivos |
|---|---:|---:|---:|
| signo por **ENLACE** (B, q=0,10) | 1,389 | **0,247** | 122 |
| signo por **SITIO** (A, p=0,30) | 7,883 | **0,046** | 181 |
| signo por **DISTANCIA** (q=0,10) | 2,442 | **0,281** | 137 |

Misma frustración, y PR/N va de 0,046 a 0,281 — un factor 6. Σ²(10) va de 7,9 a
1,4.

**Lo que decide no es cuánta frustración hay, sino dónde vive el signo.** El
signo derivado de los sitios **atrapa** los estados; el signo puesto en los
enlaces, no.

**Y eso reinterpreta la Fase IX.** Allá toda la familia de signos era por sitio,
`s_ij = (−1)^(u_i·u_j)`. La localización que reportamos —PR/N de 0,105 a 0,035—
**no venía de la frustración: venía de la codificación.** La sombra que le
pusimos al hallazgo de la Fase IX era real, pero mal atribuida.

---

## 4 · La aritmética: canal cerrado

Ahora el control iguala densidad **y** frustración de triángulos (10 controles
aceptados de 112 sorteos):

| | Σ²(10) | PR/N | dens− | frustr |
|---|---:|---:|---:|---:|
| reciprocidad p ≡ 3 mod 4 | 5,565 | 0,035 | 0,267 | 0,536 |
| azar IGUALADO | 5,310 ± 1,431 | 0,038 | 0,270 | 0,530 |
| **distancia** | **0,18 σ** | | | |

**No se separa.** Su §14 dice qué hacer con eso y se hace: **cerrar este canal**
y no insistir con esta codificación. Es la cuarta fase consecutiva en que la
misma hoja da la misma respuesta.

Con nodo: reciprocidad 3,193 contra azar+nodo 4,375 ± 1,100 → **1,07 σ**. No
alcanza, y además los tres brazos con nodo están ATRAPADOS (PR/N 0,042–0,052).
El hilo que la Fase IX marcó para más semillas queda **cerrado, no pendiente**.

### 4bis · Un residuo estructural que sí queda

Enlaces **independientes** a densidad 0,267 dan frustración ≈ 0,449. La
reciprocidad, a la misma densidad, da **0,536**. El campo aritmético vive en una
zona del plano (densidad, frustración) a la que el azar por enlace **no llega**:
por eso el sorteo tuvo que hacerse desde la familia por sitio.

Es un hecho estructural real sobre el campo de la reciprocidad — **aunque el
espectro no lo distinga**. Se registra como lo que es: una propiedad del campo,
no una señal espectral.

---

## 5 · Veredictos separados

| hipótesis | veredicto |
|---|---|
| **H1** (todo es localización) | **NO** — hay mejora con PR/N mayor que el control |
| **H2** (rigidez real) | **SÍ** — 6,324 contra 18,335 a PR/N 0,189, 15,4 σ |
| **H3** (óptimo intermedio) | existe, pero sólo entre brazos inadmisibles: **no se afirma** |
| frustración como causa | **REATRIBUIDA**: lo que decide es la codificación del signo |
| atribución aritmética | **CERRADA** — 0,18 σ con densidad y frustración igualadas |
| nodo + reciprocidad | **CERRADO** — 1,07 σ, y atrapado |

Y las reglas, que nunca se citan solas: ceros 0,3364 · piso GUE 0,5793 · piso
GOE 0,9086. Su §12 tiene razón y se obedece — acá no se persiguió 0,3364. El
mejor brazo admisible está 18,8 veces arriba de los ceros. La pregunta de esta
fase era otra, y se contestó.

---

## 6 · Lo que queda abierto

1. **¿Por qué el signo por enlace no localiza y el de sitio sí?** Es la pregunta
   que reemplaza a todas las anteriores. Un signo de sitio es casi un gauge
   (c_i·c_j lo es exactamente); el de enlace no. La sospecha: lo que atrapa es
   la *estructura de coborde parcial*, no el signo.
2. **Empujar el riel.** Con q = 0,05 el brazo queda a un nivel de ser admisible
   (148 contra 149) y da Σ² = 2,270 a PR/N 0,226. Vale medir con más modos para
   que la banda aguante — ahí puede estar el punto de verdad.
3. **La distancia.** El control por distancia dio PR/N 0,281, el más alto de la
   hoja. Que un signo que sólo depende de k funcione tan bien no estaba previsto.

---

El pedido de separar rigidez de localización es de la auditoría. El riel de
niveles vivos que faltaba lo cazó este taller, otra vez adentro y a tiempo.

Estructura cerrada no es hipótesis demostrada.

**Todavía no.**
