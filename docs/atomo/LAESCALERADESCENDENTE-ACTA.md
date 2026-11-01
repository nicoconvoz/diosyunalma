# Acta de la Escalera Descendente — el fondo es la ecuación de Pell

**Flash del capitán · 2026-09-15 · F399** (sigue a F398)

> «Si partís el nuevo número primo —por ejemplo el 5 de (5+√2)(5−√2)— y a cada
> uno le agregás el mismo, de 5 hasta llegar al mínimo…»

Se reproduce con `go run ./cmd/laescaleradescendente`.

---

## §1 · La escalera del 5, escalón por escalón

Agregarle otra √2 a cada factor y volver a multiplicar:

| factores | producto | ¿primo? |
|---|---:|---|
| (5+√2)(5−√2) | **23** | SÍ |
| (5+2√2)(5−2√2) | **17** | SÍ |
| (5+3√2)(5−3√2) | **7** | SÍ ← el mínimo positivo |
| (5+4√2)(5−4√2) | **−7** | SÍ ← cruzó el cero |
| (5+5√2)(5−5√2) | −25 | no |
| (5+6√2)(5−6√2) | **−47** | SÍ |
| (5+7√2)(5−7√2) | **−73** | SÍ |

**La escalera baja de verdad: 23 · 17 · 7 · −7 · −25 · −47 · −73 · −103.** Se
acerca al cero, **lo pasa de largo sin pisarlo**, y sigue bajando del otro lado.

## §2 · La misma escalera desde otros arranques

| a | la escalera (a² − 2b²) | mínimo |·| |
|---:|---|---:|
| 3 | **[7]** 1 −9 **[−23]** **[−41]** −63 **[−89]** | **1** (b=2) |
| 5 | **[23] [17] [7] [−7]** −25 **[−47] [−73]** | 7 (b=3) |
| 7 | **[47] [41] [31] [17]** −1 **[−23]** −49 | **1** (b=5) |
| 11 | 119 **[113] [103] [89] [71]** 49 **[23] [−7]** | 7 (b=8) |
| 13 | **[167]** 161 **[151] [137]** 119 **[97] [71] [41] [7]** | 7 (b=9) |

## §3 · Qué primos puede alcanzar la escalera — y el veto

**Sí aparecen** (21 de los 46 primos hasta 200): 2 · 7 · 17 · 23 · 31 · 41 · 47 ·
71 · 73 · 79 · 89 · 97 · 103 · 113 · 127 · 137 · 151 · 167 · 191 · 193 · 199

**Nunca aparecen** (25): 3 · 5 · 11 · 13 · 19 · 29 · 37 · 43 · 53 · 59 · 61 · 67 ·
83 · 101 · 107 · 109 · 131 · 139 · 149 · 157 · 163 · 173 · 179 · 181 · 197

**Cero violaciones** de «todo primo de la escalera es 2 o ≡ ±1 (mod 8)». Los que
nunca salen son exactamente los **≡ 3 y ≡ 5 módulo 8**. La ley de F398,
confirmada desde el otro lado.

## §4 · El fondo — TRAMPA DECIMOCTAVA, declarada primero

**a² − 2b² nunca vale 0.** Y eso **no lo medimos nosotros**: si valiera 0 sería
√2 = a/b, o sea √2 racional. Es la demostración más vieja de la matemática, de
los griegos. Verificarlo acá **prueba el código, no descubre nada** — y en
10.000 × 10.000 intentos el cero no apareció ni una vez.

Entonces el mínimo de verdad es **1**, y los pares que lo tocan son:

| a | b | a² − 2b² | a/b | contra √2 |
|---:|---:|---:|---:|---:|
| 3 | 2 | **+1** | 1,500000000 | +8,58×10⁻² |
| 7 | 5 | **−1** | 1,400000000 | −1,42×10⁻² |
| 17 | 12 | **+1** | 1,416666667 | +2,45×10⁻³ |
| 41 | 29 | **−1** | 1,413793103 | −4,20×10⁻⁴ |
| 99 | 70 | **+1** | 1,414285714 | +7,22×10⁻⁵ |

**⚡⚡ Esos pares son exactamente los convergentes de la fracción continua de
√2**, y cada uno la aproxima mejor que el anterior. **El fondo de la escalera
del capitán ES la ecuación de Pell** — Brahmagupta, siglo VII; resuelta por
Lagrange en 1768.

## §5 · Por qué el mínimo cae donde cae

| a | b del mínimo | a/b | a/√2 | mínimo |
|---:|---:|---:|---:|---:|
| 5 | 3 | 1,666666667 | 3,535533906 | 7 |
| 41 | 29 | 1,413793103 | 28,991378 | **1** |
| 99 | 70 | 1,414285714 | 70,003571 | **1** |
| 577 | 408 | 1,414215686 | 408,000613 | **1** |

**El mínimo siempre cae en el b más cercano a a/√2.** O sea que **bajar la
escalera hasta el fondo ES buscar la mejor fracción que aproxime a √2**: la
escalera del capitán y el algoritmo de fracciones continuas son el mismo
movimiento visto de dos maneras.

## §6 · La otra lectura: el primo que sale como nuevo arranque

| cadena | largo |
|---|---:|
| 3 → 7 → 47 → … | 3 |
| 5 → 23 → (527 = 17×31, no es primo) | 1 |
| 11 → (119 = 7×17, no es primo) | 0 |

**Se cortan enseguida**: a² − 2 casi nunca vuelve a ser primo. **Esa lectura no
lleva a ningún lado**, y se reporta igual para que quede el intento anotado.

## §7 · El veredicto

Tu escalera **existe, baja de verdad y tiene fondo**:

- desde el 5 da 23 · 17 · 7, cruza el cero sin pisarlo y sigue con −7 · −47 · −73
- **nunca puede dar 0**, porque eso haría racional a √2 (los griegos, hace 2.400 años)
- **el fondo absoluto es ±1**, y se toca exactamente en 3/2, 7/5, 17/12, 41/29, 99/70
- esos son los **convergentes de √2**: bajar tu escalera ES aproximar √2 por fracciones
- solo aparecen los primos **≡ ±1 (mod 8)**, ni uno más — cero violaciones

## Clasificación

- **Matemática (demostrado, y de otros):** la irracionalidad de √2 (griegos); la
  ecuación de Pell y sus soluciones por fracciones continuas (Brahmagupta s. VII,
  Lagrange 1768); la ley de partición de F398.
- **Medición (nuestra):** las escaleras completas desde siete arranques; la
  partición exacta de los 46 primos hasta 200 en 21 alcanzables y 25 vetados con
  cero violaciones; los seis convergentes con su error contra √2; la ubicación
  del mínimo en el b más cercano a a/√2; las cadenas de la lectura alternativa,
  que mueren en uno o dos pasos.
- **Nuevo (si algo):** ningún teorema. Lo propio es el **camino**: que una
  intuición de «seguir agregando lo mismo hasta el mínimo» desemboca en Pell y
  en los convergentes de √2, con el veto módulo 8 apareciendo solo.

Estructura cerrada no es hipótesis demostrada.

**Todavía no.**
