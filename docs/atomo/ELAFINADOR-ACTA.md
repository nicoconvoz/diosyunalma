# ACTA — F416 · EL AFINADOR

**Fecha:** 2026-09-19
**Código:** `cmd/elafinador`
**Lámina:** `galeria/laminas/10-el-telar/el-afinador.svg`

## Pedido del capitán

> «entonces igualá a raíz de x la fórmula, y en el cambiaformas poné el
> equivalente a poner todos los números en una variable para ver si alguno
> desentona»

Dos órdenes, y las dos son operaciones **exactas**, no aproximaciones.

## Trampa declarada ANTES de calcular (trampa nº 26 del catálogo)

Que el cambiaformas mande la línea crítica al círculo unidad es **geometría de
Möbius elemental**. Es una reescritura exacta de la Hipótesis, no una
demostración de nada. La trampa sería presentar «todos los ceros dan
exactamente 1» como si fuera evidencia a favor: no lo es, porque solo se pueden
medir los ceros que ya encontramos, y de los que faltan el instrumento no dice
una palabra. **Control obligatorio**: inyectar ceros falsos fuera de la línea y
medir si el instrumento los detecta, y publicar **dónde se queda sordo**.

## (A) Igualar la fórmula a √x

Cada par de ceros conjugados aporta

2·√x · (½·cos(γ·ln x) + γ·sen(γ·ln x)) / (¼ + γ²)

Al dividir por √x, la amplitud queda **2/|ρ|, que ya no depende de x en
absoluto**. La x sobrevive en un solo lugar: adentro del coseno, como ln x.

**Eso ES «poner todos los números en una variable»**: la variable es u = ln x, y
todo número del mundo es un punto de esa única recta.

Medido con la escalera real cribada hasta un millón:

E(x) = (ψ(x) − x + ln2π + ½·ln(1−1/x²)) / √x

| x | E(x) | cota RH (ln²x/8π) |
|---|---|---|
| 100 | −0,779252 | 0,8438 |
| 1.000 | −0,163078 | 1,8986 |
| 10.000 | +0,115588 | 3,3753 |
| 100.000 | +0,157248 | 5,2739 |
| 1.000.000 | −0,415240 | 7,5944 |

**Peor excursión en todo el rango: E(222) = −0,828217** (cota allí 1,1614).

Acotado, y siempre por debajo de la cota de Schoenfeld (1976) que vale bajo la
Hipótesis. V1 confirmado.

## (B) El cambiaformas como afinador

w(s) = 1 − 1/s = (s−1)/s. Sobre la línea crítica s = ½ + iγ:

|w| = |−½ + iγ| / |½ + iγ| = 1 **EXACTO**

El cambiaformas de la casa manda la línea crítica al **círculo unidad**. Una
sola variable por cero — el radio — y «afinado» significa exactamente 1.

**Medido: 79 ceros pescados a ciegas hasta altura 200. Peor desvío de |w|
respecto de 1: 2,220 × 10⁻¹⁶ — precisión de máquina.** V2 confirmado.

Y el cambiaformas resultó tener **orden 3**: s → 1−1/s → 1/(1−s) → s. Verificado
numéricamente con error 0 y 10⁻¹⁴. El número de la casa sigue en pie: w(2) = ½.

## (C) El control que PUEDE fallar — y falla arriba

Ceros falsos inyectados fuera de la línea:

| σ | γ | \|w\| | 1−\|w\| | predicho δ/(¼+γ²) |
|---|---|---|---|---|
| 0,600 | 14,13 | 0,9995002510 | 4,997 × 10⁻⁴ | 4,999 × 10⁻⁴ |
| 0,600 | 100 | 0,9999900003 | 1,000 × 10⁻⁵ | 1,000 × 10⁻⁵ |
| 0,600 | 500 | 0,9999996000 | 4,000 × 10⁻⁷ | 4,000 × 10⁻⁷ |
| 0,510 | 14,13 | 0,9999500114 | 4,999 × 10⁻⁵ | 4,999 × 10⁻⁵ |
| 0,501 | 500 | 0,9999999960 | 4,000 × 10⁻⁹ | 4,000 × 10⁻⁹ |

Sí desentona: cualquier cero fuera de la línea cae fuera del círculo, y la
fórmula predicha coincide con lo medido en todas las cifras.

**⚠ PERO el desvío muere como 1/γ².** A altura 500 un cero corrido 0,1 mueve el
radio apenas 4 partes en diez millones. **El afinador se queda sordo en altura**,
y eso hay que decirlo con el número puesto. V3 y V4 confirmados, incluida la
parte incómoda.

## (D) El eco desentonado — por qué el √x es el límite

Un cero en σ da un eco de tamaño x^σ; dividido por √x queda x^(σ−½):

| x | σ = 0,51 | σ = 0,6 |
|---|---|---|
| 10⁶ | 1,15 | 3,98 |
| 10¹² | 1,32 | 15,8 |
| 10³⁰ | 2,00 | 1.000 |
| 10¹⁰⁰ | 10,0 | 10¹⁰ |

Un solo cero desafinado termina tapando a toda la orquesta. Por eso el √x no es
un detalle técnico: **es la afinación entera**.

## Qué es de quién

La transformación de Möbius y el hecho de que mande una recta a un círculo son
matemática del siglo XIX. La cota |ψ(x) − x| < √x·ln²x/(8π) bajo la Hipótesis es
de **Schoenfeld (1976)**.

Lo propio y medido: la excursión acotada de E(x) sobre el millón entero, el
desvío de máquina en los 79 ceros de pesca propia, la fórmula de sensibilidad
δ/(¼+γ²) verificada contra lo medido, y **el número que dice dónde el
instrumento se queda sordo**.

## El cierre

El capitán pidió una variable donde se vea si alguno desentona, y la variable
existe y es una sola: **el radio**.

Estar afinado es estar en el círculo. **La Hipótesis de Riemann es que NINGUNO
desentona.**

Todavía no.
