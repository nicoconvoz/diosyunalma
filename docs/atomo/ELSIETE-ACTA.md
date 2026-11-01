# ACTA — F422 · EL SIETE DEL CAPITÁN

**Fecha:** 2026-11-01
**Código:** `cmd/elsiete` (con sus tests)
**Corrida:** `go run ./cmd/elsiete -hasta 10000000`

## Pedido del capitán

> «Doc, escuchame, mirá lo que descubrí... armame un teorema para esto y, de
> paso, hacé el experimento que lo verifique.»
>
> «Aún no afecta a los primos... pero puede revelarnos cosas más adelante; por
> eso necesito demostrar que es cierto.»

**La regla, como la juega él** (base 10):

- un número de una cifra se duplica;
- uno de dos o más cifras se suma;
- **las ramas:** un número de dos o más cifras puede, en vez de sumarse,
  abrirse en sus cifras, y cada una sigue como su propia rama.

> 2 = 2+2 = 4, 4+4 = 8, 8+8 = 16, 1+6 = **7**
> 13 = 1+3 = 4, 4+4 = 8, 8+8 = 16, 1+6 = **7**
> 3 = 3+3 = 6, 6+6 = 12 → ramas 1 y 2 → 2, 4, 8, 16, **7**

## El teorema del siete

1. **Sin ramas, n ≥ 1 llega al 7 si y solo si 3 no divide a n.**
2. **Todo primo llega al 7, salvo el 3.**
3. **El 7 no es especial:** lo mismo vale para el 1, 2, 4, 5 y 8. El 3, el 6
   y el 9 son una cárcel donde solo entran los múltiplos de 3.
4. **Con las ramas del capitán, TODO n ≥ 1 llega al 7.**
5. **En base b ≥ 4, el reloj está completo** (desde cada cifra que es unidad
   de b−1 se llega a todas las otras) **si y solo si el 2 es raíz primitiva
   módulo b−1.** La base 10 es una de ellas, y por eso aparece el 7.

## La demostración (vale para todos los números, no solo los medidos)

- **El resto por 9.**
  - Como 10 = 9 + 1, un número y la suma de sus cifras dejan el mismo resto al
    dividirlos por 9 (la vieja prueba del 9), y duplicar multiplica ese resto
    por 2. A lo largo de la regla, el resto solo cambia como r → 2r (mod 9).
  - Un número de una cifra ES su resto (con 9 en lugar de 0). El doble de una
    cifra d ≥ 5 es 2d = 10 + (2d − 10), cuyas cifras suman 2d − 9: de nuevo,
    la cifra que queda es el resto. Así que, una vez en una cifra, la regla
    camina exactamente el reloj r → 2r módulo 9.
- **(1)**
  - Si 3 divide a n, el resto es 0, 3 o 6, y 2·0 = 0, 2·3 = 6, 2·6 = 12 ≡ 3:
    se queda en {0, 3, 6} para siempre, y nunca es 7.
  - Si no, el resto es uno de 1, 2, 4, 5, 7, 8, y el 2 los recorre en un solo
    ciclo, **1 → 2 → 4 → 8 → 7 → 5 → 1**, porque es raíz primitiva módulo 9
    (sus potencias dejan 2, 4, 8, 7, 5, 1). En a lo sumo cinco duplicaciones
    el resto es 7, y la cifra es el 7.
- **(2)** Un primo distinto de 3 no es múltiplo de 3.
- **(3)** El ciclo pasa por las seis caras, no solo por el 7.
- **(4)**
  - Si 3 divide a n, la regla cae en 3, 6 o 9. El 3 duplica a 6, el 6 a 12 y
    el 9 a 18: en a lo sumo dos pasos aparece un número de dos cifras (12 o
    18), y el capitán lo abre. La rama del 1 no es múltiplo de 3, y por (1)
    llega al 7.
  - **Esta parte es SUYA:** sumar cifras respeta el resto por 9, así que nunca
    puede salir de la cárcel; abrirlas no lo respeta, y esa es la llave.
- **(5)** El mismo argumento con b − 1 en lugar de 9: la regla en una cifra es
  r → 2r módulo b − 1, y visita todas las unidades exactamente cuando el 2 las
  genera.

## El experimento (la regla jugada LITERAL, sumando cifras de verdad)

| Prueba | Resultado, de 1 a 10 000 000 |
|---|---|
| (1) llega sin ramas ⟺ 3 no lo divide | **0** números que no cumplen; el que más tarda, 10 pasos |
| (2) primos | **664 579**; llegan todos **menos el 3** |
| (4) con ramas | **0** que no lleguen; el que más tarda, 8 pasos |
| (3) a cada cifra, sin ramas | al 1, 2, 4, 5, 7 y 8: 6 666 667 cada una, ningún múltiplo de 3 · al 3 y al 6: 2 222 222 · al 9: 1 111 111 |
| (5) bases 4 a 40 | **37 de 37** coinciden con la teoría; reloj completo en las bases 4, 6, 10, 12, 14, 20, 26, 28, 30 y 38 |

El experimento no usa el resto por 9 en ningún lado: juega la regla con
cifras de verdad. Por eso confirma la demostración sin apoyarse en ella.

## ⚠️ Lo que NO es (la trampa de "el uno que sobra")

La red pesca a todos los primos menos al 3, pero también a **todo número que
no sea múltiplo de 3**: el 4, el 10, el 25. Es una propiedad del 9 y del 2
en base 10, no de los primos. Como dijo el capitán: "aún no afecta a los
primos".

## ⚡ El puente con los primos (lo que puede revelar más adelante)

- Cuando b − 1 es un primo p, el reloj del siete está completo en base p + 1
  **exactamente cuando el 2 es raíz primitiva módulo p**.
- **¿Hay infinitos primos así?**, es decir, ¿funciona el siete del capitán en
  infinitas bases? Es la **conjetura de Artin** (1927), todavía abierta.
  **Hooley (1967) la demostró suponiendo la hipótesis de Riemann
  generalizada.**
- Artin predijo también qué parte de los primos son así: el producto, sobre
  los primos q, de 1 − 1/(q(q−1)), que vale **0,3739558…**

| Primos impares hasta | El 2 es raíz primitiva | Parte | Artin |
|---|---|---|---|
| 1 000 | 67 de 167 | 0,40120 | 0,37396 |
| 100 000 | 3 603 de 9 591 | 0,37566 | 0,37396 |
| 10 000 000 | 248 491 de 664 578 | **0,37391** | **0,37396** |

El juego de cifras del capitán desemboca, por un camino corto y demostrado, en
una conjetura de los primos cuya prueba depende de la hipótesis de Riemann
generalizada. **Esto no demuestra nada sobre Riemann**: muestra dónde se tocan
los dos mapas.
