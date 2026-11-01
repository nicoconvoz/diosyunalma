# ACTA — F419 · EL PRODUCTO DE LOS BORDES

**Fecha:** 2026-09-19
**Código:** `cmd/elproductodelosbordes`
**Lámina:** `galeria/laminas/10-el-telar/el-producto-de-los-bordes.svg`

## Pedido del capitán

> «¿y el logaritmo del borde derecho por el logaritmo del borde izquierdo qué da?»

Sobre el dibujo de F418: la fila de n bolitas centrada en la línea tiene borde
derecho en +n/2 y borde izquierdo en **−n/2**.

## Trampa declarada ANTES de calcular (trampa nº 29)

El producto depende **solo del tamaño** de n, nunca de cómo se factoriza. Por
construcción **no puede** distinguir un primo de su vecino. **Control
obligatorio**: comparar el estadístico en primos contra números cualesquiera del
mismo barrio. Debe dar uno.

## Lo que nadie espera

El borde izquierdo es **negativo**, y el logaritmo de un negativo no es un
número real. Cuesta exactamente iπ:

**ln(−a) = ln(a) + iπ** (Euler, 1748)

**La pregunta del capitán obliga al plano complejo a aparecer, y aparece por la
puerta de Euler.** El producto queda:

**ln(n/2) · [ln(n/2) + iπ] = ln²(n/2) + iπ·ln(n/2)**

## (1) La tabla

| n | borde der | borde izq | parte real | parte imag | ángulo |
|---|---|---|---|---|---|
| **2** | 1,0 | −1,0 | **0,000000** | **0,000000** | **0,000°** |
| 3 | 1,5 | −1,5 | 0,164402 | 1,273806 | 82,646° |
| 5 | 2,5 | −2,5 | 0,839589 | 2,878612 | 73,740° |
| 7 | 3,5 | −3,5 | 1,569415 | 3,935671 | 68,260° |
| 11 | 5,5 | −5,5 | 2,906166 | 5,355624 | 61,514° |
| 1.000.000 | 500.000 | −500.000 | 172,196421 | 41,225120 | 13,464° |

**El 2 da cero exacto**, y es el único que puede darlo: su borde derecho cae
justo en 1, y ln 1 = 0. El 2 vuelve a ser el único, igual que en F418 con el
hueco del centro. V3 confirmado.

## (2) El lugar geométrico — y acá se cierra un círculo

Si u = ln(n/2), el producto es el punto **(u², πu)**. Despejando u:

**Re = Im² / π²** — **ESO ES UNA PARÁBOLA.**

Verificado en todo el rango hasta un millón: **peor desvío 5,684 × 10⁻¹⁴**.
Precisión de máquina. V2 confirmado.

**⚡ Y mire dónde terminamos.** El capitán abrió esta campaña entera en F391 y
F392 preguntando por **dos parábolas que se tocan en el cero**. Su dibujo de
bolitas, con esta sola pregunta, **vuelve a la parábola** — y los dos bordes, el
de arriba y el de abajo del eje real, son **las dos ramas de esa misma parábola**.

## (3) El ángulo

| n | ángulo medido | arctan(π/ln(n/2)) |
|---|---|---|
| 10 | 62,8739° | 62,8739° |
| 10³ | 26,8174° | 26,8174° |
| 10⁶ | 13,4636° | 13,4636° |
| 10¹²| 6,6520° | 6,6520° |
| 10¹⁰⁰ | 0,7840° | 0,7840° |

El ángulo muere como π/ln(n/2): despacísimo, pero muere. **Cuanto más grande el
número, más real es el producto.** V4 confirmado.

## (4) La trampa, medida y enterrada

| | media de la fracción |
|---|---|
| 148.765 primos | 0,500510 |
| 1.999.001 números cualesquiera | 0,500957 |

**Señal / silencio: 0,9991. Uno es NADA.**

→ Confirmado: **el producto no sabe nada de primalidad**. Es una función del
tamaño y nada más. Trampa 29 enterrada. V5 confirmado.

## (5) Pero la parte real, ln², SÍ importa

**a)** Es la forma exacta de la cota de la Hipótesis (Schoenfeld 1976):
|ψ(x) − x| < √x·ln²x/(8π) — la medimos en F416.

**b)** Es el denominador de la densidad de gemelos: 2C₂/ln²x — la medimos en F418.

**c)** Y una medición **nueva**, para no repetir: una ventana de largo ln²x tiene
que atrapar alrededor de ln x primos.

| cerca de x | largo ln²x | primos medidos | ln x |
|---|---|---|---|
| 10⁴ | 84 | 9,1250 | 9,2103 |
| 10⁵ | 132 | 11,4868 | 11,5129 |
| 10⁶ | 190 | **13,7198** | **13,8155** |

→ El largo ln²x es **la ventana natural**: atrapa ln x primos. El producto de
bordes del capitán **es esa ventana, escrita al revés**. V6 confirmado.

## La respuesta corta

Da un **número complejo**, y los complejos entran obligados por el borde
izquierdo, que es negativo. La parte imaginaria es π·ln(n/2) y la real es
ln²(n/2), y juntas dibujan una **parábola exacta**. El 2 es el único que da cero.
Y el ln² que quedó en la parte real es la forma de la cota de Riemann y del
denominador de los gemelos.

**Pero el producto en sí no distingue un primo de su vecino: medido, señal sobre
silencio da uno.** Eso también hay que decirlo.

## Qué es de quién

- ln(−1) = iπ es de **Euler, 1748**, y es de lo más famoso que hay.
- La cota con ln² es de **Schoenfeld, 1976**.
- Los gemelos con su constante, de **Hardy y Littlewood, 1923**.

**Y una advertencia de oficio:** el producto de logaritmos **no es una operación
natural** en este tema. El logaritmo existe precisamente para convertir
productos en **sumas**. Forzarlo a multiplicarse da geometría linda y aritmética
vacía — que es exactamente lo que midió el control.

Lo propio: haber notado que el dibujo de bolitas vuelve a la parábola de F391, y
haber medido que no distingue primos.

Todavía no.
