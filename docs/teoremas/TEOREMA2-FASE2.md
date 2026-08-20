# Teorema 2 — FASE 2: el batido se puede predecir

**Documento para la auditora · 2026-08-14 · respuesta a la hoja «El Batido
Protector» (FASE 2)** — su §14 ejecutado al pie de la letra: predicción
ANTES de la medición. Verificación: `go run ./cmd/elbatidoprotector`
(F308). Estado: la pregunta del §15 tiene respuesta afirmativa dentro de
los límites declarados. Ningún teorema se afirma todavía.

---

## 0. Convención

La de F304 (δ = log natural de r, techos hacia arriba), más:

    A(n) = e^{n·δ} + e^{−n·δ}          [la amplitud del cuarteto]
    B(n) = |cos(n·Δθ/2)|               [la envolvente del batido]
    perla base: par DH · n₀_sola = 85622 (la zona de ruptura)

## 1. Lema 2.1 de la hoja — VERIFICADO

Las pausas viven exactamente donde el §6 dice: B(n_k) = 0 en
n_k = (2k+1)·π/|Δθ| (máximo medido en los ceros: 1e−13, ruido float).

## 2. La fórmula a priori — y el tesoro

Para colocar la pausa k SOBRE la zona de ruptura:

    Δθ*_k = (2k+1)·π / n₀_sola     ⟹     τ*_k = 1 + Δθ*_k/θ₁

Predicho ANTES de medir (par DH): τ* = 1.00314 (k=0) · 1.00943 (k=1) ·
1.01572 (k=2). Medido DESPUÉS:

    k = 0:  n₀ = 101340   (+18.4%)   ← EL ESCUDO MAYOR
    k = 1:  n₀ =  92725   (+8.3%)
    k = 2:  n₀ =  73185   (−14.5%)   ← NO protege (fallo honesto, abajo)

**El tesoro:** la FASE 1, con su barrido ciego (paso 0.005), había
encontrado como máximo 89454 (+4.5%) en τ = 1.010 ≈ k=1. La fórmula
señaló un punto ENTRE los nodos de la grilla — τ = 1.00314 — y ahí estaba
un escudo 4.1 veces mayor. **La teoría le dijo al experimento dónde
mirar.**

## 3. El predictor de un parámetro — candidato a Lema 2.2

    C* = A(n₀_gemelas)                 [calibrado SOLO en las gemelas]
    n₀_pred(τ) = primer n con B(n)·A(n) ≥ C*

Sin coro, sin ajustes por caso. Conjunto de prueba de doce τ (de 0.4 a
2.5, incluyendo las tres pausas predichas y los intervalos musicales):

    error MEDIANO: 1.5% · diez de doce bajo el 5%
    el peor: 24.1% en τ = 1.01572 (la pausa k = 2)

**El fallo honesto:** la pausa k = 2, más angosta (la anchura de la pausa
decrece con k), medio-yerra la zona: ni protege ni el predictor la
acierta. La anchura decreciente de las pausas es la pregunta natural de
la FASE 3.

## 4. El protocolo del §14, paso a paso

    1 · perla base fijada ........ DH, convención F304
    2 · zona de ruptura sola ..... n₀ = 85622
    3 · pausas predichas ......... n_k = (2k+1)π/Δθ (lema 2.1)
    4 · alineación con la zona ... Δθ*_k = (2k+1)π/n₀
    5 · n₀ medido ................ DESPUÉS de predecir
    6 · comparación .............. mediana 1.5%; el tesoro en k = 0

## 5. Respuesta a la pregunta del laboratorio (§15)

**SÍ**: la diferencia angular que coloca una pausa del batido exactamente
sobre la zona de ruptura es Δθ* = (2k+1)·π/n₀_sola, y la pausa ancha
(k = 0) es el escudo mayor. Con las palabras de la hoja: la FASE 1
descubrió el batido; la FASE 2 debía descubrir si podemos predecirlo —
**podemos**. La teoría cuantitativa de la armonía tiene su primera
fórmula predictiva.

## 6. Límites declarados

Una perla base (DH; la generalización a otras bases es FASE 3), una
ventana (n ≤ 3×10⁵), un parámetro calibrado (C*, en las gemelas), DOS
cuartetos. El outlier k = 2 queda como el borde conocido del predictor.

---

*Una simulación descubre; una identidad explica; un lema reduce; un
teorema cierra todos los pasos. La FASE 2 agrega el escalón que faltaba
entre los dos primeros: una fórmula que PREDICE — y que encontró lo que
la búsqueda ciega no vio.*
