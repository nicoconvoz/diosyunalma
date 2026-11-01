# ACTA — F420 · EL ZOOM DEL MURCIÉLAGO

**Fecha:** 2026-09-19
**Código:** `cmd/elzoomdelmurcielago`
**Lámina:** `galeria/laminas/10-el-telar/el-zoom-del-murcielago.svg`
**Trabaja sobre:** F214 (el tambor es el libro), F359/F415 (el sonar), F104 (el globo)

## Pedido del capitán

> «tomá este hallazgo, vamos a trabajar sobre él. Vamos a comprimir y luego a
> ensanchar con la función inversa sobre un punto de esa ventana comprimida
> donde el murciélago detecte algo, un barrio o algo cerca en el sentido
> figurado, y ahí sí hallaremos cosas.»

Eso es el problema del atlas, y tiene una versión afilada:

> **¿Puede un atlas comprimido y de baja resolución APUNTAR EL TELESCOPIO a un
> blanco que él mismo NO PUEDE RESOLVER?**

Si sí, vale la pena construir atlas. Si no, hay que barrer todo a resolución
completa para siempre y la compresión es decoración.

## Trampa declarada ANTES de calcular (trampa nº 30)

Sería **circular** marcar el barrio con un detector que ya vio el par apretado.
El detector de acá se arma **solo con el barrido crudo**, que por construcción
**no puede resolver** los pares que marca. Eso es lo que hace honesto el test.

## (0) La verdad de campo

Barrido fino completo, paso 0,002, entre t = 14 y t = 5.000:
**4.520 ceros, 2.719.000 evaluaciones.**

## (1) COMPRIMIR

Barrido **crudo** a propósito demasiado grueso, paso 0,25:
**4.508 ceros vistos de 4.520 — 12 PERDIDOS**, 245.345 evaluaciones.

Ajuste ciego en la base {t·ln(t/2π), t, 1} — la jugada de F214:

| constante | sobre el crudo | teórica | desvío |
|---|---|---|---|
| a | 0,157884 | 0,159155 | 1,27 × 10⁻³ |
| b | −0,152149 | −0,159155 | 7,01 × 10⁻³ |
| c | 0,022647 | 0,875000 | 8,52 × 10⁻¹ |

**⚠ V1 NO SALIÓ COMO LO DECLARÉ.** Antes de culpar al método hay que separar
las dos causas. Mismo ajuste, sobre la verdad de campo:

| constante | sobre la verdad | teórica | desvío |
|---|---|---|---|
| a | 0,159155 | 0,159155 | **1,49 × 10⁻⁷** |
| b | −0,159154 | −0,159155 | **9,27 × 10⁻⁷** |
| c | **0,874955** | 0,875000 | **4,49 × 10⁻⁵** |

→ **El método no es el problema: sobre datos completos clava.** Lo que degrada
el ajuste crudo son los doce ceros perdidos, porque cada pérdida corre **toda**
la escalera de ahí en adelante.

**⚡ Y eso es una SEGUNDA ALARMA, independiente de la primera**: un ajuste que no
clava las constantes conocidas está avisando que al catálogo **le faltan
objetos**, sin decir cuáles ni dónde.

**⚡ Y una mejora sobre nuestro propio F214**: allá la tapa 7/8 solo asomaba en
0,907 con Δ0,032, y se registró la honestidad de que estaba colineal. Acá, sobre
la función de conteo de ceros, sale **0,874955 — Δ4,5 × 10⁻⁵**, unas **700 veces
mejor**.

**Compresión: 4.508 ceros guardados en tres constantes. Razón 1.503 a 1, y el
objeto comprimido pesa 24 bytes.**

## (2) ENSANCHAR CON LA INVERSA

N̄ es invertible: **N̄⁻¹(n) dice dónde vive el cero número n antes de mirarlo.**

| n | predicho (atlas limpio) | real | error en espaciados |
|---|---|---|---|
| 100 | 237,7183 | 237,7697 | −0,0297 |
| 500 | 812,7399 | 812,7711 | −0,0242 |
| 1.000 | 1.420,6768 | 1.420,4165 | +0,2246 |
| 2.000 | 2.516,4609 | 2.516,5711 | −0,1051 |
| 4.000 | 4.507,2378 | 4.507,7452 | −0,5310 |

Con el atlas **crudo** el error crece hasta **6,44 espaciados**; con el limpio se
queda en fracciones. El daño se ve por tercera vez.

## (3) EL MURCIÉLAGO ESCUCHA

D(t) = N_crudo(t) − N_comprimido(t). **Donde D cae, se perdió algo.** El
murciélago **no ve el objeto: oye que falta.**

Los barrios más cantados: t = 4.589,00 (caída −2,140) · 3.866,00 (−2,059) ·
2.939,00 (−1,929) · 2.927,00 (−1,928) · 1.811,00 (−1,700).

## (4) ZOOM SOLO AHÍ, Y (5) EL JUICIO

30 barrios ensanchados a paso fino: **10 ceros perdidos RECUPERADOS**.
30 barrios al azar del mismo ancho: **0**.

**⚠ CORRECCIÓN A MI PROPIO CONTROL, antes de festejar.** El cociente contra cero
daría infinito, y sería humo: con 12 ceros perdidos repartidos en el tramo y
cubriendo el 1,81 % del rango, **el azar tenía que encontrar 0,217 — o sea, cero
era el resultado esperado**. Un control que encuentra lo que ya se predecía no
separa nada.

**La estadística que sí vale:**

| | |
|---|---|
| ceros perdidos en todo el tramo | 12 |
| recuperados en los barrios cantados | **10 (83,3 % de recuperación)** |
| porción del rango que se miró | **1,81 %** |
| **ENRIQUECIMIENTO** | **46,2×** |

→ **EL ATLAS SABE APUNTAR.** El objeto comprimido dirige el telescopio a blancos
que él mismo no puede resolver. V3 confirmado.

## (6) El ahorro

| | evaluaciones |
|---|---|
| barrido fino completo | 2.719.000 |
| crudo + 30 zooms dirigidos | 292.469 |
| **ahorro** | **9,3× menos trabajo** |

Y encontrando el 83 % de lo que el barrido completo encuentra. V4 confirmado.

## (7) Qué hay adentro del barrio

Espaciado normalizado medio: **0,9998** (tiene que dar 1).
Los cinco pares más apretados: **0,0890 · 0,0946 · 0,1105 · 0,1112 · 0,1249**.

Son **pares de Lehmer**: ceros que casi se tocan. Son los que el crudo perdía, y
los que el murciélago mandó a buscar **sin haberlos visto nunca**.

## Qué es de quién

- La parte suave del conteo es **Riemann–von Mangoldt** (1859/1905).
- El ajuste ciego de invariantes desde un espectro es **Weyl (1911)**, vía F214.
- El zoom dirigido por un modelo grueso es **multi-resolución de manual**: no lo
  inventamos.

Lo propio: haberlo **medido con un control que podía matarlo**, sobre nuestro
propio océano; haber separado el daño del método; haber corregido el propio
control cuando dio infinito; y haber sacado el 7/8 con 700 veces más precisión
que en F214.

**Y esto no demuestra nada de la Hipótesis: es CARTOGRAFÍA.** Pero es exactamente
la maquinaria que pedía la entrada «Atlas comprimido» del informe de
aplicaciones — y ahora tiene número de mérito.

Todavía no.
