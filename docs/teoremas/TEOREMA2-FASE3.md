# Teorema 2 — respuestas del taller al §10 del «Pancho Completo»

**Documento para la auditora · 2026-08-15 · las cuatro preguntas de su
§10, respondidas** — dos por derivación, una por eliminación de parámetro,
y el porqué del k = 2. Verificación: `go run ./cmd/laanchura` (F309).
Convención: la de F304/F308 (δ = log r, A(n) = e^{nδ}+e^{−nδ},
B(n) = |cos(n·Δθ/2)|, base DH, n₀_sola = 85622).

---

## Q2 + Q3 — la anchura efectiva, DERIVADA

Cerca de un cero de la envolvente, B sube lineal: B(n) ≈ |n − n_k|·Δθ/2.
La pausa protege mientras B(n)·A(n) < C*, de donde la anchura total:

    P(k, Δθ) = 4·C* / (Δθ · A(n_zona))

No hay ajuste: es el cruce lineal de la envolvente con el umbral del
predictor.

## Q1 — por qué k = 0 protege mucho más que k = 2

Con las pausas ALINEADAS (Δθ*_k = (2k+1)π/n₀), todas caen en el MISMO
lugar n₀ — pero Δθ*_k crece como (2k+1), así que

    w_k = w₀ / (2k+1)        w₀ = 60584 · w₁ = 20195 · w₂ = 12117 escalones

Mismo lugar, pausa cinco veces más angosta: la de k = 2 no puede sostener
la zona contra las resonancias vecinas. Medido:

    k = 0: margen +15718 · k = 1: +7103 · k = 2: −12437 · k = 3: −12431

La escalera 1 : 1/3 : 1/5 predice el ORDEN de los escudos y la caída de
k = 2 — el fallo honesto de la FASE 2, ahora explicado.

## Q4 — la eliminación de C* (candidato)

Las gemelas duplican la amplitud efectiva de las resonancias, de modo que
necesitan aproximadamente la MITAD de la amplitud de la perla sola:

    C*_der = A(n₀_sola)/2 = 18.250      [calibrado: 20.284 · desvío −10.0%]

El test de las doce τ, repetido con el C* derivado — **sin calibrar
nada**:

    error mediano con C* calibrado: 1.5%
    error mediano con C* DERIVADO:  4.5%

El predictor sobrevive sin su calibración, con el costo declarado. La ley
predictiva queda **sin parámetros libres** — candidata para el Lema de
Interacción de la FASE 3. (El −10% del C* derivado viene de ignorar la
fluctuación del coro: derivación a primer orden, declarado.)

## El sello, registrado

La frase de su §9 queda en el registro del laboratorio: *«sí sello que
existe una ley predictiva experimental reproducible dentro del alcance
declarado»* — la primera ley predictiva sellada del Teorema 2. El teorema
universal sigue en rojo, como ella marca.

## Límites

Una base (DH), una ventana (n ≤ 3×10⁵), derivaciones a primer orden, DOS
cuartetos. Herencia para la FASE 3: la fórmula de anchura P(k,Δθ), el
predictor sin parámetros, y el candidato a Lema de Interacción («dos
desafinadas pueden retrasarse, jamás salvarse»).

---

*«La FASE 1 descubrió el batido. La FASE 2 consiguió que el batido dijera
dónde mirar.» — y esta hoja agrega: la FASE 3 ya sabe cuánto mide la
pausa.*
