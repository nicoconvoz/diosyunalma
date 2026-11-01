# ACTA — F417 · LAS FRACCIONES DE PRIMOS

**Fecha:** 2026-09-19
**Código:** `cmd/lasfracciones`
**Lámina:** `galeria/laminas/10-el-telar/las-fracciones.svg`

## Pedido del capitán

> «¿qué relación hay entre fracciones de números primos, por ejemplo 3/7, 7/3,
> 2/3, 3/2, 2/7, 7/2 — y entre sus ramas? ¿da alguna pista?»

Escribió los seis pares ordenados que se pueden armar con {2, 3, 7}. Esos seis
objetos llevan **tres estructuras distintas adentro**, y las tres son reales.

## Trampa declarada ANTES de calcular (trampa nº 27)

La mediante de 2/3 y 3/7 es (2+3)/(3+7) = 5/10 = **½**, y ½ es EL número de la
casa. Huele a señal. **Control obligatorio**: contar cuántas veces las mediantes
de fracciones de primos caen en ½ y comparar contra fracciones de enteros
cualesquiera del mismo rango.

## (1) Como NÚMEROS — la fracción ES la resolución del sonar

El logaritmo convierte la fracción en una **resta**: ln(p/q) = ln p − ln q. Y esa
resta es exactamente la cantidad del número clave de F415, γ_max = 2π/|ln p − ln p′|.

| fracción | valor | ln | γ_max |
|---|---|---|---|
| 3/7 | 0,428571 | −0,847298 | 7,42 |
| 7/3 | 2,333333 | +0,847298 | 7,42 |
| 2/3 | 0,666667 | −0,405465 | 15,50 |
| 3/2 | 1,500000 | +0,405465 | 15,50 |
| 2/7 | 0,285714 | −1,252763 | 5,02 |
| 7/2 | 3,500000 | +1,252763 | 5,02 |

**La fracción y su vuelta dan el mismo |ln| con signo opuesto**: son el mismo
intervalo musical, subiendo o bajando.

**Y ln(3/2) = 0,4055 es EXACTAMENTE el número que ya estaba en la tabla de F415
para separar el 2 del 3.** El capitán llegó a esa cantidad por su cuenta,
escribiendo fracciones.

Además se encadenan: ln(7/2) = ln(7/3) + ln(3/2), verificado al sexto decimal.
O sea 7/2 = (7/3)·(3/2). **Las fracciones de primos forman un grupo**, y ese
grupo es ℚ⁺ entero, generado libremente por los primos.

## (2) Como SÍMBOLOS — reciprocidad cuadrática

En teoría de números, (3/7) escrito así **es el símbolo de Legendre**: pregunta
si 3 es un cuadrado módulo 7. Y la relación entre (p/q) y (q/p) es la **ley de
reciprocidad cuadrática** de Gauss.

| par | (p/q) | (q/p) | producto | predicho |
|---|---|---|---|---|
| 3, 7 | −1 | +1 | −1 | −1 |
| 3, 11 | +1 | −1 | −1 | −1 |
| 5, 13 | −1 | −1 | +1 | +1 |
| 7, 11 | −1 | +1 | −1 | −1 |
| 11, 19 | +1 | −1 | −1 | −1 |

**Verificado en 4.371 pares de primos impares hasta 500: CERO violaciones.**
V1 confirmado.

El sentido en criollo: p y q se preguntan «¿sos un cuadrado en mi mundo?», y
**siempre se contestan lo mismo… salvo que los dos sean 3 módulo 4**, y ahí se
contestan al revés. Gauss 1796, que la demostró ocho veces distintas.

## (3) Las RAMAS — lo que el capitán preguntó, y es contable

El decimal de a/p se repite con período ord₁₀(p), y las fracciones a/p se parten
en exactamente **(p−1)/ord₁₀(p) collares** distintos. **Esas son las ramas.**

| p | período | ramas | collar de 1/p |
|---|---|---|---|
| 3 | 1 | **2** | 3 |
| 7 | 6 | **1** | 142857 |
| 11 | 2 | **5** | 09 |
| 13 | 6 | **2** | 076923 |
| 17 | 16 | **1** | 0588235294117647 |
| 19 | 18 | **1** | 052631578947368421 |
| 31 | 15 | **2** | 032258064516129 |

**Collares construidos dígito por dígito para todo primo < 200: 0 desajustes.**
V2 confirmado.

**El 7 tiene UNA sola rama**: 1/7 = 142857, 3/7 = 428571, 2/7 = 285714 — son la
**misma rueda girada**. El 11 tiene cinco ruedas distintas. El 13, dos.

¿Cada cuánto un primo tiene una sola rama? Medido hasta 100.000: **3.617 de
9.590 = 0,377164**, contra la **constante de Artin (1927) = 0,3739558**.
Diferencia 0,0032. V3 confirmado.

**⛔ Y eso NO está demostrado.** Es conjetura abierta desde 1927.

## (4) LA PISTA — que es lo que se preguntó

Ninguna fracción de primos vale 1 salvo la trivial: ningún producto 2^a · 3^b ·
7^c da 1 salvo a = b = c = 0. **Eso es factorización única**, y en logaritmo
significa que ln2, ln3 y ln7 son **linealmente independientes sobre ℚ**.

Medido: la combinación más cercana a cero con coeficientes hasta 40 en valor
absoluto es

−13·ln2 + 33·ln3 − 14·ln7 = **0,00055009** — nunca cero.

**Y ESA es la pista.** Las frecuencias del sonar **nunca entran en fase**. Si dos
logaritmos de primos fueran conmensurables, la música se repetiría y los primos
tendrían un patrón periódico. No lo son. Por eso la música no se repite jamás, y
por eso el sonar necesita **infinitos** ceros para ser exacto — que es
exactamente el precio que se midió en F415.

V4 confirmado.

## (5) La trampa, y una corrección a mi propio control

| | aciertos | de | tasa |
|---|---|---|---|
| mediantes = ½ entre fracciones de PRIMOS < 50 | 478 | 43.890 | 1,0891 % |
| control con enteros cualesquiera del mismo rango | 552 | 43.890 | 1,2577 % |

**Señal / silencio = 0,8659.** Por debajo de uno: la mediante ½ es aritmética de
números chicos, no estructura. Trampa 27 declarada y enterrada.

**⚠ CORRECCIÓN A MI PROPIO CONTROL.** El primer control que escribí daba vuelta
numerador y denominador. **Eso no era un control**: la lista de fracciones de
primos ya contiene p/q **y** q/p, así que el intercambio manda el conjunto sobre
sí mismo y da el número idéntico por construcción — y en efecto dio 478 de
43.890, exactamente igual. Un control que no puede diferir del caso medido no
mide nada. Lo reemplacé por el de enteros cualesquiera, que sí puede fallar.

## Qué es de quién

- Reciprocidad cuadrática: **Gauss, 1796**.
- Los períodos decimales y sus collares: **Gauss, Disquisitiones, arts. 315–317**.
- La conjetura sobre la densidad: **Artin, 1927**. Abierta.
- La independencia lineal de los ln p: consecuencia directa de la factorización
  única, **Euclides**.

Lo propio y medido: haber visto que el **|ln| de la fracción ES el número clave
de nuestro F415**, haber contado las ramas una por una construyendo los collares,
y haber medido la mejor combinación que se acerca a cero.

## El cierre

Nada de esto predice un primo. **Ordena lo que ya sabíamos** — pero lo ordena
alrededor de una sola idea: la fracción de dos primos, mirada en logaritmo, es
la distancia entre dos notas del sonar, y ninguna de esas distancias es un
múltiplo racional de otra.

Todavía no.
