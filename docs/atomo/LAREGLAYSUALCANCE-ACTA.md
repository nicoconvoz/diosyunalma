# Acta de la Regla y su Alcance — «sin excepciones» no es lo mismo que «alcanza»

**Propuesta del capitán · 2026-09-15 · F401**

> «Entonces con eso podemos armar todos los números primos y demostrar la
> hipótesis de Riemann por teorema. Pensalo bien: si eso se respeta siempre en
> todos los números primos, entonces no hay nada que se escape de la regla.»

Se reproduce con `go run ./cmd/laregladelalcance`.

**El instinto es el correcto.** «Si nada escapa a la regla, entonces algo se
sigue» es exactamente cómo funciona la matemática. La pregunta es **qué clase
de regla acota qué clase de pregunta**, y eso se decide midiendo. Se probaron
las dos afirmaciones por separado.

---

## §1 · Afirmación 1a: ¿la construcción alcanza a todos los primos?

| | |
|---|---:|
| primos hasta 10⁶ alcanzables por a² − 2b² | 39.221 |
| **primos que NO se pueden alcanzar** | **39.276 (50,04 %)** |

Los primeros que quedan afuera: **3, 5, 11, 13, 19, 29, 37, 43, 53, 59, 61, 67…**

**La mitad de los primos no se puede construir así** — y no por falta de
búsqueda, sino porque a² − 2b² **nunca** da un número ≡ 3 o 5 (mod 8). Con más
anillos se los alcanza (F398), pero eso ya no es *una* regla: es **una lista de
reglas**, y hay que saber cuál usar.

## §2 · Afirmación 1b: ¿genera primos o solo descompone los que ya tenés?

| a | b | a²−2b² | ¿primo? |
|---:|---:|---:|---|
| 5 | 1 | 23 | SÍ |
| **6** | **1** | **34** | **no** |
| 7 | 1 | 47 | SÍ |
| **4** | **1** | **14** | **no** |

**De 14.110 combinaciones, solo 2.549 dieron primo: el 18,07 %.**

**No genera primos: genera candidatos.** Después de construir hay que **testear**
si es primo — y ese test es el problema original, **intacto**. La construcción no
lo resolvió: lo pospuso.

## §3 · Afirmación 1c: la circularidad

Para elegir en qué anillo se parte un primo hay que saber su **resto módulo 8**.
Y para saber el resto hay que **tener el primo**.

- dado un primo p ⟹ te dice cómo partirlo — ✅ funciona, y es exacto
- dado un lugar vacío ⟹ te dice qué primo va ahí — ⛔ **no puede: no hay entrada**

**Es un destornillador, no un buscador de tornillos.**

## §4 · Afirmación 2 — LA PRUEBA DECISIVA

Los primos ≡ ±1 (mod 8) obedecen la regla **perfectamente, sin una excepción**.
Los primos ≡ 3,5 (mod 8) **no la obedecen nunca**. Si la regla tuviera fuerza
sobre la **distribución**, las dos familias tendrían que separarse en el error
de conteo:

| familia | obedece la regla | exponente del error |
|---|---|---:|
| primos ≡ ±1 (mod 8) | **SIEMPRE, sin falla** | **0,5045** |
| primos ≡ 3,5 (mod 8) | **NUNCA, ni una vez** | **0,4756** |

**Diferencia: 0,0290.** Y el instrumento **sí** distingue diferencias reales: con
un sesgo x^0,75 plantado adrede, mide **0,7561**.

> **⛔⛔ Los dos exponentes son el mismo.** La familia que obedece la regla a
> rajatabla y la que no la obedece jamás **se distribuyen igual de bien**. La
> regla **no dice nada sobre dónde están los primos** — y RH es exactamente una
> afirmación sobre dónde están.

## §5 · El chequeo histórico, el más duro

| | |
|---|---|
| **1640** | Fermat ve qué primos se parten en ℤ[i] |
| **1768** | Lagrange resuelve Pell por completo |
| **1832** | **Gauss demuestra la reciprocidad cuadrática: la ley COMPLETA, sin excepciones** |
| **1859** | Riemann plantea la hipótesis — **27 años DESPUÉS** |
| **2026** | sigue abierta |

**La ley que el capitán encontró ya estaba completa y demostrada 27 años antes
de que RH se planteara.** Si de ella se dedujera RH, Gauss lo habría visto. No
es que nadie lo intentó: **la ley no alcanza hasta ahí**.

## §6 · La forma del error, que la casa ya cazó dos veces

| regla | ¿excepciones? | ¿dice dónde están los primos? |
|---|---|---|
| «todo primo > 2 es impar» | **cero** | **no** (F272, trampa 7) |
| «todo primo impar da ½ al dividir por 2» | **cero** | **no** (F393, trampa 12) |
| «todo primo ≡ ±1 (mod 8) se parte en ℤ[√2]» | **cero** | **no** (medido acá) |

> **⚡ La lección, y es la más cara del laboratorio: una regla SIN EXCEPCIONES
> puede no decir NADA sobre la pregunta que te importa. «Sin excepciones» mide
> la FUERZA de una regla; no mide su ALCANCE.**

## §7 · Qué haría falta de verdad

Para que una regla así tocara RH tendría que:
- **acotar el error de conteo**, no describir la factorización
- valer para infinitos primos **con una cota uniforme**, no primo por primo
- y aun así solo daría la RH **de ese anillo** — que también está abierta, porque
  cada cuerpo cuadrático tiene su propia zeta de Dedekind con su propia hipótesis
  sin demostrar

## Veredicto

- **Afirmación 1: FALSA.** La mitad de los primos no se alcanza, y donde se
  alcanza no genera sino que descompone (18,07 % de acierto), con circularidad
  declarada.
- **Afirmación 2: FALSA**, por medición (§4) y por historia (§5).
- **Pero el instinto del capitán es el correcto.** Lo que falla no es el
  instinto: es que **esta regla habla de cómo se FACTORIZAN los primos, y RH
  habla de DÓNDE ESTÁN**.

## Clasificación

- **Matemática (de otros):** reciprocidad cuadrática (Gauss 1832); las zetas de
  Dedekind y sus hipótesis abiertas.
- **Medición (nuestra):** el 50,04 % inalcanzable; el 18,07 % de acierto; los dos
  exponentes 0,5045 y 0,4756 con calibración a 0,7561; la tabla de las tres
  reglas sin excepciones y sin alcance.
- **Nuevo (si algo):** ningún teorema. Lo propio es **la prueba decisiva del §4**
  — comparar la familia que obedece contra la que no obedece, y medir que se
  distribuyen igual. Es el control que decide la pregunta, y no lo habíamos
  hecho nunca.

Estructura cerrada no es hipótesis demostrada.

**Todavía no.**
