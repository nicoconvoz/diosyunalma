# Acta del Descenso al Mínimo — el átomo del ladrillo es el 2

**Flash del capitán · 2026-09-15 · F400** (corrección y continuación de F399)

> «Esto sirvió pero no me entendiste: a ese 5 que sale ahí quiero hacerle el
> equivalente a ± √2, hasta que quede un mínimo número con √2.»
>
> *(y a mitad de camino)* «Quiero el átomo de los ladrillos jejej»

Se reproduce con `go run ./cmd/eldescensoalminimo`.

**Me había equivocado en F399.** El capitán no quería subir el b: quería la
operación aplicada **al 5 mismo**, recursivamente. Y esa pregunta tiene una
respuesta exacta y un fondo que se reproduce a sí mismo.

**El átomo es el 2.**

---

## §1 · El 5 sí se parte igual

La operación que hizo el 23 a partir del 5 es **n ↦ n² − 2**. Para partir el 5
del mismo modo hace falta un a con a² − 2 = 5, o sea **a = √7**:

> **5 = (√7 + √2)(√7 − √2) = 4,059964873438 × 1,231537748691 = 5,000000000000**

Funciona porque 7 − 2 = 5. Ninguno de los dos factores es entero. Y al √7 se le
puede hacer lo mismo, y al que salga también, para siempre.

## §2 · El descenso

| paso | el número | t (con n = 2·cosh t) | t anterior / 2 |
|---:|---:|---:|---:|
| 0 | 5,000000000000 | 1,566799236972 | — |
| 1 | 2,645751311065 | 0,783399618486 | **0,783399618486** |
| 2 | 2,155400498994 | 0,391699809243 | **0,391699809243** |
| 4 | 2,009596961618 | 0,097924952311 | **0,097924952311** |
| 8 | 2,000037458306 | 0,006120309519 | **0,006120309519** |
| 10 | 2,000002341137 | 0,001530077380 | **0,001530077380** |

**El descenso converge al 2.** Y en cada escalón el parámetro t se **parte al
medio, exacto** — la coincidencia de las dos últimas columnas no es
aproximación, es identidad.

## §3 · Por qué se detiene ahí: el 2 se reproduce

> **(2 + √2)(2 − √2) = 3,414213562373 × 0,585786437627 = 2,000000000000**

**Da 2. El mismo 2 con el que empezó.** Es el único número que la operación
devuelve intacto, y por eso es el fondo. La cuenta, en un renglón:

    x² − 2 = x  ⟹  x² − x − 2 = 0  ⟹  (x − 2)(x + 1) = 0  ⟹  x = 2

| vecino | n² − 2 | ¿vuelve al mismo? |
|---:|---:|---|
| 1,9 | 1,610000 | no, se va |
| **2,0** | **2,000000** | **SÍ ← el punto fijo** |
| 2,1 | 2,410000 | no, se va |

**Abajo del 2 la operación se desbarranca** (1,9 → 1,61 → 0,59 → −1,65);
**arriba se dispara** (2,1 → 2,41 → 3,81 → 12,5). El 2 es el filo del cuchillo.

## §4 · El nombre de lo que pasa: la operación duplica un ángulo

Si n = 2·cosh(t), entonces **n² − 2 = 2·cosh(2t)**. La operación del capitán no
hace más que **duplicar t**. Por eso bajar es partirlo al medio, y por eso el
fondo es t = 0, donde 2·cosh(0) = **2**.

Con ángulos reales es la fórmula del ángulo doble: 4cos²θ − 2 = 2cos(2θ) — la
recurrencia de **Chebyshev**.

## §5 · Para arriba, con enteros: los números de Lucas

| paso | número | ¿primo? | |
|---:|---:|---|---|
| 1 | 3 | SÍ | L(2) |
| 2 | 7 | SÍ | L(4) |
| 3 | 47 | SÍ | L(8) |
| 4 | 2.207 | SÍ | L(16) |
| 5 | 4.870.847 | **no** | L(32) |
| 6 | 23.725.150.497.407 | **no** | L(64) |

Son los **números de Lucas en las potencias de dos**. La cadena, subiendo,
camina esa sucesión. **La primalidad se midió, no se supuso: se cae en el
quinto escalón.**

## §6 · Trampa decimonovena, declarada

Escribir el descenso entero de una vez da la **raíz anidada**:

> 2 = √(2 + √(2 + √(2 + …)))  — verificado con 60 raíces: **2,000000000000000**

Es una **identidad clásica** de cualquier libro de análisis. Se reporta con el
nombre puesto y no como hallazgo. Lo mismo Lucas y Chebyshev.

## §7 · La respuesta, y el átomo del ladrillo

En F397 la metáfora fue que partir un primo es como partir un ladrillo: con las
manos no, pero bajando de nivel sí. El capitán pidió entonces **el átomo del
ladrillo**. Acá está, y tiene número:

> **El átomo es el 2.**

- el 5 se parte en (√7 + √2)(√7 − √2), y cada resultado se vuelve a partir igual
- el descenso baja siempre y **no se pasa**: 5 → 2,6458 → 2,1554 → 2,0388 → …
- **el 2 se reproduce**: (2 + √2)(2 − √2) = 2. Ahí la escalera se termina sola
- abajo se desbarranca, arriba se dispara: es el filo exacto
- y el motivo de fondo: la operación duplica un ángulo, y bajar lo parte al
  medio hasta llegar a cero

**Un átomo de verdad: lo partís y te devuelve lo mismo.**

## Clasificación

- **Matemática (demostrado, y de otros):** el punto fijo de x² − 2 = x; la
  identidad n² − 2 = 2cosh(2t) con n = 2cosh(t) (Chebyshev); los números de
  Lucas y L(2n) = L(n)² − 2; la raíz anidada clásica.
- **Medición (nuestra):** el descenso desde 5 con t partiéndose al medio escalón
  por escalón; el comportamiento de los vecinos del 2; la primalidad medida de
  la cadena de Lucas, que se cae en L(32).
- **Nuevo (si algo):** ningún teorema. Lo propio es **la pregunta bien
  perseguida** — que «seguir partiendo hasta el mínimo» tenga una respuesta
  exacta, un número que se reproduce, y un porqué con nombre.

Estructura cerrada no es hipótesis demostrada.

**Todavía no.**
