# ACTA — F418 · LA APERTURA

**Fecha:** 2026-09-19
**Código:** `cmd/laapertura`
**Lámina:** `galeria/laminas/10-el-telar/la-apertura.svg`

## Pedido del capitán

> «imaginá a los números como bolitas horizontales con el centro en una línea,
> entonces tenemos que
>
> ```
>         ooo---------->3
>          oo---------->2
>           o ---------->1
> ```
>
> Medí la apertura del espacio de los números primos teniendo en cuenta este
> concepto, a ver si nos da una pista.»

## Trampa declarada ANTES de calcular (trampa nº 28)

En el dibujo aparecen dos **medios**: el ángulo mínimo sale arctan(½) y las
filas pares tienen sus bolitas corridas medio lugar. Huele a la línea crítica.
**Control obligatorio**: medir si ese ½ distingue primos de no-primos.

## (1) La bolita del centro

Leído literal: n impar tiene una bolita **justo en la línea**; n par tiene un
**hueco**. De ahí sale algo exacto y único:

**El 2 es el único primo con un hueco en el medio.** 148.932 primos de 148.933
hasta 2.000.000 tienen bolita central.

**⚠ CONTROL, porque acá está la tentación.** Si la bolita del centro explicara
algo, todos los impares serían primos. Medido: **14,89 % de los impares son
primos** hasta 2 millones, contra 2/ln x = 13,78 % predicho.

→ **La bolita del centro es PARIDAD. No dice nada de primalidad.** V1 confirmado.

## (2) El ángulo de apertura — acá sí se mide, y es EXACTO

De una fila prima a la siguiente la fila crece g bolitas, o sea **g/2 de cada
lado**, en un solo escalón hacia arriba. Entonces:

**ángulo = arctan(g / 2)**

| hueco g | g/2 | ángulo | qué es |
|---|---|---|---|
| **1** | 0,5 | **26,565°** | **SOLO entre el 2 y el 3 — una vez en toda la historia** |
| **2** | 1 | **45,000°** | **PRIMOS GEMELOS — cuarenta y cinco EXACTOS** |
| 4 | 2 | 63,435° | primos primos |
| 6 | 3 | 71,565° | sexy primes |
| 34 | 17 | 86,634° | |

Verificado sobre todos los primos hasta 2.000.000: **hueco mínimo g = 1, en
p = 2, una sola vez**. Después de ahí todos los primos son impares y el hueco
más chico posible pasa a ser 2, para siempre. V2 confirmado.

## (3) Los 45 grados

Pares que abren a 45° exactos hasta 2.000.000: **14.871**.

| forma de Hardy–Littlewood | predice | razón medido/predicho |
|---|---|---|
| cruda, 1,3203236·x/ln²x | 12.545 | 1,1855 |
| **honesta, 1,3203236·∫dt/ln²t** | **14.796** | **1,0051** |

**⚠ ERROR PROPIO CORREGIDO.** Escribí primero la forma cruda **y además la
dupliqué**, porque 1,3203236 ya ES 2·C₂ y yo volví a multiplicar por dos. Dio
25.089 contra 14.871 medidos, razón 0,59. Corregido: sin el doble y con la
integral, la razón es 1,0051. Es exactamente la lección de F414 — hay que
**sumar** la densidad, no evaluarla al final.

**⛔ Que haya infinitos 45° es conjetura abierta.** Nadie lo demostró.
V3 confirmado tras la corrección.

## (4) Cómo se abre el cono

Apilando las filas primas, una por escalón:

| escalón k | medio ancho | semiángulo medido | arctan(ln k/2) | con Cipolla |
|---|---|---|---|---|
| 10 | 14,5 | 55,408° | 49,023° | 46,892° |
| 100 | 270,5 | 69,711° | 66,525° | 68,710° |
| 1.000 | 3.959,5 | 75,826° | 73,853° | 75,690° |
| 10.000 | 52.364,5 | 79,188° | 77,749° | 79,146° |
| 100.000 | 649.854,5 | **81,252°** | 80,145° | **81,225°** |

**⚠ CORRECCIÓN A MI PROPIO VEREDICTO V4.** Declaré que el semiángulo seguiría
arctan(ln k/2). Es demasiado crudo y se queda corto en todo el rango. La que
acierta es la de Cipolla, la misma que medimos en F413:
p_k ≈ k(ln k + ln ln k − 1). Con esa, 81,225° contra 81,252° medidos.

→ **El cono NO abre recto. El exceso de apertura es el logaritmo**, y el
semiángulo trepa hacia los 90° sin llegar nunca.

### Y algo que NO es un error, aunque lo parezca

| cerca de | hueco medio | MEDIA de arctan | arctan de la MEDIA |
|---|---|---|---|
| 1.000 | 6,533 | 67,947° | 72,979° |
| 10.000 | 9,183 | 71,265° | 77,714° |
| 1.000.000 | 13,851 | 75,152° | 81,784° |

Las dos últimas columnas **no son lo mismo**, y la diferencia no es un fallo:
**arctan es cóncava**, así que el promedio de los ángulos es menor que el ángulo
del promedio. Es la desigualdad de Jensen. Significa que **los huecos chicos
pesan más en la apertura** de lo que sugiere el hueco medio: el cono abre más
despacio de lo que el logaritmo solo haría creer.

## (5) La otra apertura — en cuántas CAJAS se abre la fila

Una fila de n bolitas se dobla en un rectángulo a×b si n = a·b, así que la
cantidad de cajas es **d(n)**, el número de divisores.

| hasta | cajas promedio | ln n + 2γ − 1 |
|---|---|---|
| 10.000 | 9,366800 | 9,364772 |
| 100.000 | 11,667500 | 11,667357 |
| 1.000.000 | **13,970034** | **13,969942** |

Dirichlet, 1849. **Y otra vez aparecen el logaritmo y γ.**

→ **Un primo se abre en DOS. Siempre.** Es el piso de la apertura, y la única
caja posible es el palo original y su vuelta. V5 confirmado.

## (6) La trampa, enterrada

Filas pares con corrimiento exacto de ½: **50.000 de 50.000 = 100 %**. O sea
**todas**, primas o no. Es paridad, no es la línea crítica. Y el ½ de arctan(½)
sale simplemente de que el hueco 1 se reparte entre los dos lados.
Trampa 28 declarada y enterrada. V6 confirmado.

## LA PISTA

**La apertura de los primos tiene piso y techo, y los dos son exactos.**

- **Piso absoluto: 26,565° = arctan(½)**, una sola vez en toda la historia
  infinita, entre el 2 y el 3.
- **Piso repetible: 45° exactos**, y son los primos gemelos. Que sean infinitos
  está abierto.
- **Techo: 90°**, al que el cono se acerca sin llegar, y la velocidad con que se
  acerca **ES el logaritmo**.
- Y la otra apertura, la de las cajas, tiene también **piso 2 en los primos** y
  promedio **ln n + 2γ − 1** en todos.

Mirando bolitas, sin una sola fórmula, el dibujo del capitán devuelve **el
logaritmo dos veces por caminos distintos**.

## Qué es de quién

Nada de esto es nuevo. arctan(g/2) es trigonometría de escuela aplicada a los
huecos entre primos, que es el objeto más estudiado del tema. El promedio de
divisores es de **Dirichlet, 1849**. La conjetura de los gemelos con su
constante es de **Hardy y Littlewood, 1923**, y sigue abierta.

Lo propio: haber traducido el dibujo a un ángulo medible, haber verificado que
el piso absoluto ocurre **una sola vez**, y haber encontrado la brecha de Jensen
entre el ángulo medio y el ángulo del hueco medio.

**Y esto no predice un primo.** Reordena el hueco en forma de ángulo.

Todavía no.
