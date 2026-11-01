# ACTA — F415 · EL SONAR EXACTO

**Fecha:** 2026-09-19
**Código:** `cmd/elsonarexacto`
**Lámina:** `galeria/laminas/10-el-telar/el-sonar-exacto.svg`

## Pedido del capitán

> «alinea la profundidad y el logaritmo con porcentajes con el sonar de nuestro
> murciélago (buscar en hallazgos antiguos) a ver que encontramos quiero un
> número clave para calcular no una aproximación»

Dos exigencias: alinear tres piezas propias (F414 la profundidad, F413 el
porcentaje, F359 el sonar del murciélago) y entregar un número **exacto**, no
una estimación.

## Trampa declarada ANTES de calcular (trampa nº 26 del catálogo)

La trampa disponible era responder «la constante clave es γ = 0,5772…» o
«es 1,3203…» y vender una constante suelta como si calculara primos. Cualquier
constante fija que no dependa de x no puede dar exactitud en una cuenta que
depende de x. **Control**: cualquier número que se proponga como «clave» tiene
que mejorar una cuenta medible frente a la aproximación cruda, y la mejora se
reporta con el número medido, no con adjetivos.

## Lo que se hizo

1. Criba hasta x = 100 y construcción de ψ(x) = Σ_{p^k ≤ x} ln p (la escalera
   real, contada una por una).
2. Pesca a ciegas de ceros de ζ sobre la línea crítica por cambio de signo de
   la función Z de Hardy, hasta altura 200: **79 ceros**.
3. Evaluación de la fórmula explícita
   ψ(x) = x − Σ_ρ x^ρ/ρ − ln(2π) − ½·ln(1 − 1/x²)
   con k = 0, 1, 5, 10, 20, 40, 79 ecos, en x = 20, 24, 30, 50, 70, 90.
4. Cálculo de la resolución γ_max = 2π/|ln p − ln p′| para pares de primos
   vecinos.

## Resultados medidos

| x | error sin ecos | error con 79 ecos | mejora |
|---|---|---|---|
| 20 | 0,7340 | 0,0147 | **50,0×** |
| 30 | 1,5243 | 0,0455 | **33,5×** |
| 50 | — | — | **5,2×** |
| 70 | — | — | **12,2×** |
| 90 | — | — | **2,3×** |

ψ(30) verdadero = 28,476498562.

**Número clave calculable — la ley de resolución:**

γ_max ≈ 2π / |ln p − ln p′|

| par | γ_max |
|---|---|
| 2 y 3 | 15,5 |
| 11 y 13 | 37,6 |
| 97 y 101 | 155,5 |

Cuanto más juntos están dos primos, más alto hay que escuchar para separarlos.

## El alineamiento con el murciélago (F359)

El término de cada cero oscila como cos(γ · ln x). Los ceros **cantan en el
logaritmo de x**, que es exactamente la profundidad medida en F414. Por eso sus
picos caen en los logaritmos de los primos, que es lo que el murciélago de F359
escuchaba sin saber por qué: los períodos T = k · ln p de aquella pieza son los
mismos que aparecen acá dentro del coseno. Las tres piezas no se parecen: son
la misma, vista desde tres lados.

## Correcciones propias

1. **Error escrito y corregido**: dije que sin ningún eco la fórmula daba 30.
   Es falso: da 28,16, porque los términos constantes −ln(2π) y −½ln(1−1/x²) ya
   están adentro antes del primer eco.
2. **Predicción propia que FALLÓ**: predije por escrito que en x = 24 (lejos de
   un escalón de la escalera) la convergencia se vería limpia, sin oscilación.
   No se ve. También oscila. El efecto dominante no es el fenómeno de Gibbs
   junto a los saltos, sino que la suma sobre ceros converge **condicionalmente**:
   el orden de los ecos importa y las sumas parciales oscilan en todos lados.

## Qué es de quién

La fórmula explícita es de **Riemann (1859)**, demostrada rigurosamente por
**von Mangoldt (1895)**. No es de la casa y es de lo más hermoso que existe.

Lo que aportó la casa: medir cuánto mejora de verdad con ecos contados de una
pesca propia (de 2× a 50× según el punto), calcular la ley de resolución que el
capitán pedía como número para calcular, y mostrar que el eje del murciélago de
F359 y la profundidad de F414 son literalmente el mismo eje que el de la
fórmula.

## El precio, dicho completo

El número exacto existe, pero cuesta: la igualdad requiere **infinitos ceros**.
Se cambió una aproximación corta por una igualdad con una suma infinita adentro.
Eso no es un truco: es la estructura real del problema.

## El cierre

El √x que aparece en la fórmula es el tamaño **máximo** de cada eco, si y solo
si todos los ceros están sobre la línea. Dicho de otro modo:

**La Hipótesis de Riemann es, literalmente, que ningún eco del murciélago suena
más fuerte que √x.**

Todavía no.
