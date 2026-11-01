# Acta de los Dos Espejos — el primo es el que menos se refleja

**Flash del capitán · 2026-09-15 · F396**

> «Si ponemos entre dos espejos a un número primo, ¿qué observamos adelante y
> atrás?»

Se reproduce con `go run ./cmd/losdosespejos`.

---

## §1 · La imagen, traducida exacto

Parado en n y mirando a distancia k para los dos lados se ve el par (n−k, n+k).
Los dos reflejos son primos justo cuando ese par es de primos — y como
**(n−k) + (n+k) = 2n**:

> **⚡ Contar los reflejos de n ES contar las particiones de Goldbach de 2n.**

El flash preguntó **Goldbach** sin saberlo. (Eso es álgebra, no un hallazgo: se
declara como identidad.)

## §2 · La trampa, declarada antes de contar — DECIMOQUINTA

k tiene que tener **paridad opuesta** a n, o los dos reflejos salen pares y
mueren al nacer. Eso tacha la mitad de las distancias para **todo** n, primo o
compuesto: no dice nada sobre primalidad.

## §3 · El único reflejo que solo tiene un primo

A distancia **k = 0** el número se ve **a sí mismo**: 2n = n + n es partición de
Goldbach solo cuando n es primo. Cierto, exclusivo de los primos — **y es la
definición**, así que se declara y no se vende. El resto cuenta k ≥ 1.

## §4 · Un defecto propio, cazado antes de registrar

La primera corrida tenía un **error nuestro**: al buscar los divisores impares
no se sacaba primero el factor 2, así que en los números pares el primo grande
del final se perdía (el 100002 = 2·3·7·2381 aparecía con 2 divisores en vez de
3) y se inventaba una clase espuria de «0 divisores» con 99 miembros. Corregido
antes de anotar nada; las tablas de abajo son las de después del arreglo. La
conclusión no cambió, pero el número de divisores sí.

## §5 · La respuesta al flash

Medido sobre los 1.000 números de [100.000, 101.000):

| grupo | cuántos | S(n) medio | reflejos medios |
|---|---:|---:|---:|
| **PRIMOS** | 81 | **1,0000** | **1.065,6** |
| compuestos | 919 | 1,5554 | 1.660,9 |

> **El primo ve un 35,8 % MENOS reflejos que un compuesto de su tamaño.**

Ejemplos a ojo: el 100003 (primo) tiene 1.070 reflejos; el 100050 = 2·3·5²·23·29
tiene **3.104**; el 100485 tiene **3.947**.

**Y la razón es la de F395 dada vuelta:** un primo **no tiene divisores chicos**.
Los divisores son los que enriquecen la estructura simétrica; el primo no tiene
ninguno, así que su factor es el mínimo posible y **es el que menos se refleja**.

## §6 · La ley detrás: Hardy–Littlewood, verificada

S(n) = ∏ (q−1)/(q−2) sobre los primos impares q que dividen a n:

| divisores impares de n | cuántos | S(n) medio | reflejos medios | **reflejos / S(n)** |
|---:|---:|---:|---:|---:|
| 1 | 184 | 1,0030 | 1.069,4 | **1.066,2** |
| 2 | 445 | 1,3450 | 1.432,6 | **1.065,2** |
| 3 | 306 | 1,8441 | 1.962,4 | **1.064,2** |
| 4 | 62 | 2,5348 | 2.697,5 | **1.064,2** |

Los reflejos crudos varían en un factor 2,5; **divididos por S(n) quedan
clavados en 1.065 ± 1** — un 0,2 %. Y la dispersión relativa cae de 0,3879 a
0,0131: **dividir por S aplana 29,6 veces**. Los divisores explican casi toda la
diferencia; lo que queda es el ruido de los primos.

Detalle fino que cierra solo: la razón compuestos/primos de los reflejos vale
**1,5586**, y la razón de sus S(n) vale **1,5554**. Coinciden al 0,2 %.

## §7 · Lo que NO se ve

- Que el número de reflejos **nunca llegue a cero** es la **conjetura de
  Goldbach**, abierta desde **1742**. Acá se cuenta cuántos hay; nadie demostró
  que no pueda haber ninguno.
- La fórmula que gobierna todo esto es de **Hardy y Littlewood (1923)** y
  **también es una conjetura**. Lo medido es que sus números dan; eso no la
  demuestra.
- Este tramo llega a 101.000. Nada se afirma más allá.

## Clasificación

- **Matemática (identidad, y de otros):** reflejos de n = particiones de
  Goldbach de 2n (álgebra); el factor singular de Hardy–Littlewood (conjetura
  ajena, 1923).
- **Medición (nuestra):** primos 1.065,6 contra compuestos 1.660,9 (−35,8 %);
  reflejos/S(n) constante a 1.065 ± 1 sobre un rango de 2,5× en S; aplanado
  29,6×; la coincidencia 1,5586 contra 1,5554.
- **Nuevo (si algo):** ningún teorema. La traducción limpia del flash a
  Goldbach, y la verificación cuantitativa de que **la escasez de reflejos del
  primo es exactamente su escasez de divisores** — el mismo mecanismo de F395,
  visto en el espejo.

Estructura cerrada no es hipótesis demostrada.

**Todavía no.**
