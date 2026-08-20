# Acta de la Separación — respuesta a la Auditoría 56 (batería completa)

**Respuesta a la guía «El residuo y la puerta hacia los ceros» · 2026-08-20 · F384**

Sus cinco experimentos (A–E) y su prueba de oro (§15), corridos enteros a
máxima potencia. Decisión metodológica declarada: F_Λ es **sintética**
(amplitudes Λ(m)/(√m·ln m), fase cero — la truncación explícita y controlada
de su §7), para que el residuo limpio no herede ruido de medición. Ventana
base [10000, 10500], dt = 0,005, **588 ceros reales** apilados.

Se reproduce con `go run ./cmd/laseparacion`. Lámina:
`galeria/laminas/10-el-telar/la-separacion.svg`.

---

## A · Ampliar la escalera — la púa mientras la aritmética se retira ENTERA

| N retirado | tonos | varianza R_N | púa s=0 |
|---:|---:|---:|---:|
| 40 | 19 | 0,851 | −4,895 |
| 100 | 35 | 0,745 | −4,777 |
| 500 | 114 | 0,600 | −4,569 |
| 2000 | 333 | 0,489 | −4,361 |
| 10000 | 1280 | 0,388 | **−4,133** |

**La púa sobrevive a la retirada de TODA la aritmética controlable** (1280
tonos, hasta m = 10000): baja lenta, no colapsa. Su pregunta del §8:
permanece, disminuyendo despacio — y el porqué de esa disminución lo da B.

## B · Control aritmético — y la contabilidad que cierra a la milésima

La cola sintética Λ(41..10000) sola — sin zeta, sin ceros — apilada en los
ceros reales da **−0,761**. Y la púa real perdió −4,895 → −4,133 = **−0,762**
al retirar exactamente esa cola. **El libro cierra a la milésima**: lo que la
púa pierde al ampliar N es exactamente lo que la cola aritmética carga. La
aritmética truncada por sí sola NO fabrica la púa (−0,76 contra −4,9); pero
cada tramo de escalera lleva su cuota exacta de púa — la dualidad, contada
moneda por moneda.

## C · Centros falsos (su §10)

R_2000 apilado en: ceros reales **−4,361** · corridos +0,2 → **+0,455** ·
azar con igual densidad → **+0,029** · huecos permutados → **−0,010** ·
medios → **+0,226**. La púa está atada a la POSICIÓN de los ceros; ni la
densidad ni los intervalos permutados la ven.

## D · Cambio de ventana (su §11)

t ≈ 10000 / 30000 / 100000: púa **−4,361 / −4,400 / −4,300**, controles
+0,23 / +0,36 / +0,44. Tres alturas, misma púa.

## E · Cambio de resolución (su §12)

dt = 0,01 / 0,005 / 0,0025: hombros ±0,04 y ±0,20 **idénticos a la milésima**
(−0,64, +0,45); el fondo s=0 se PROFUNDIZA con la resolución (−3,7 → −4,4 →
−5,1). Eso no es artefacto: es la firma de una singularidad logarítmica —
cuanto más cerca muestreás de log|t−γ|, más hondo lo ves. El perfil es
estable; el pozo es infinito, como manda log.

## ORO · Su §15 — la descomposición F ≈ F_Λ + α·F_ceros

F_ceros se construyó SOLO desde los ceros: núcleo mollificado log(|t−γ|/W),
W = 2, destendido. Dos confesiones de método antes de los números:

1. El primer intento usó corte duro |t−γ| < 10 y dio α = 0,14: cada cero
   entrando o saliendo de la franja saltaba log(10) — un serrucho que aplasta
   la regresión. Instructivo, y queda declarado.
2. Con el núcleo mollificado, la regresión punto a punto da **α = 0,30,
   correlación 0,55, 29,7%** del residuo fino explicado por los ceros solos.
   NO clava el 1 — y se reporta tal cual: mezcla la zona singular
   (coeficiente 1) con la zona media, donde el producto truncado no tiene por
   qué empatar.

El estimador que la teoría realmente predice — la **forma de la púa apilada**
contra log|s| (el promedio sobre 588 ceros aísla la singularidad):

| perfil | pendiente medida | Hadamard |
|---|---:|---:|
| canto completo tras Λ≤40 | **0,877** | 1 |
| residuo limpio R_2000 | 0,693 | (menos: la cola se llevó su cuota) |
| control en centros al azar | −0,033 | 0 |

La cascada de varianza: **1,711 → 0,485 (tras Λ) → 0,341 (tras Λ y ceros)**.

## Veredicto, con el semáforo de su §16

- La descomposición F = F_Λ + F_ceros + E **dejó de ser dibujo**: cada pieza
  es medible, la contabilidad A↔B cierra a la milésima, y la púa tiene la
  forma log|s| con pendiente 0,88 donde la teoría pone 1.
- Lo que NO cerró, dicho con todas las letras: el α punto a punto (0,30) está
  lejos del 1 — la zona media entre ceros no está capturada por el núcleo
  local truncado. Ahí vive el E(t) de su §13, con varianza 0,34 aún sin dueño.
- Rojo intacto: ningún teorema, nada sobre RH. Correlación y contabilidad,
  no causa.

Su frase de trabajo, respondida: la escalera y la púa **SE PUEDEN separar** —
cuantitativa, reproducible y establemente — y la parte de cada una quedó
pesada. Lo que queda sin separar (la voz media, var 0,34) es el siguiente
objeto, y tiene nombre en su §13: E(t).

Estructura cerrada no es hipótesis demostrada.

**Todavía no.**
