# Acta de los Tres Regímenes — Fase VI

> **RETRACTACION POSTERIOR (Fase VII, 2026-08-17).** El hallazgo central de esta
> acta —la «esquina rígida» en (A,B) = (30,32) con α = 0,313— **NO SOBREVIVE** al
> chequeo de robustez que pidió la auditora. En ese punto sobreviven 28 niveles de
> 400 dentro de la banda: el 93 % del espectro fue expulsado, y el α se midió
> sobre esas migas. No era una fase del medio: era una banda vaciándose. Entre los
> puntos que SÍ dejan espectro medible, α nunca baja de 1,5, o sea que no hay fase
> rígida en ninguna parte de esta familia. La corrección completa, con la tabla de
> niveles vivos, está en docs/atomo/HILADOFINO-FASE7-ACTA.md. Esta acta se deja como se
> escribió: la casa no reescribe historia.


**Respuesta a la Hipótesis de los Tres Regímenes (Auditoría 39) · 2026-08-17**

> «Podría tener tres modos: voltaje alto, amperaje alto o una mezcla equilibrada
> de ambos.»
> — **Jesús Nicolás Astorga**

**La intuición es del capitán y resultó ser el diagnóstico exacto de lo que la
Fase V no podía resolver.** Se reproduce con `go run ./cmd/tresregimenes`.

---

## 1 · Las dos variables reales (su §5: encontrarlas, no suponerlas)

Su documento prohíbe tomar la metáfora literalmente. Buscamos las variables
matemáticas, y estaban ahí, atadas:

**La Fase V controlaba una sola cantidad — la fuerza total — y esa única
restricción ataba dos cosas distintas:**

- **A** — cuánto empuja **un** par de niveles sobre otro. *(su «voltaje»)*
- **B** — a **cuántos** pares llega ese empuje, el alcance. *(su «amperaje»)*

Con la traza fija, comprar alcance se pagaba en empuje. **Ése era el nudo**, y es
exactamente lo que su intuición señaló: una perilla donde había que tener dos.

El modelo, fuera de la diagonal:

    H_ij = A · amp_i · amp_j · exp(−|i−j|/B)         amp = Λ(p)/√p

La **fuerza total pasa a ser una cantidad derivada**, no el control. Eso es lo que
permite la prueba decisiva.

R6: única entrada aritmética Λ(n); 400 modos desde t = 100; ningún parámetro
elegido mirando los γₙ.

---

## 2 · La prueba decisiva: misma fuerza, distinto reparto

Si el resultado dependiera sólo del total, su hipótesis perdía. Buscamos pares
con fuerza casi igual y repartos distintos:

| fuerza total | A | B | Σ²(10) |
|---:|---:|---:|---:|
| 1,54 | 0,3 | 32 | **7,097** |
| 1,82 | 10 | 0,5 | **4,103** |

**Razón 1,73×.** El resultado **no depende sólo del total: depende del reparto.**

**La hipótesis queda apoyada: son dos grados de libertad independientes, no uno.**

---

## 3 · El mapa, y el observable que casi se nos escapa

Su §9 pedía un mapa y no un óptimo suelto. Fue la decisión correcta, porque el
óptimo de Σ² es la pregunta equivocada.

**Lo que mide rigidez no es el valor de Σ², es cómo CRECE con la distancia.**
Definimos α por Σ²(L) ∝ L^α:

- **α = 1** → la varianza crece como la distancia: **sin rigidez** (Poisson).
- **α → 0** → la varianza casi no crece: **rígido**, que es lo que tienen los ceros.

Y ahí aparece la estructura:

| A \ B | 0,5 | 2 | 8 | 32 |
|---:|---:|---:|---:|---:|
| 0,3 | 1,096 | 1,124 | 1,076 | 1,009 |
| 1 | 1,254 | 0,951 | 0,968 | 1,067 |
| 3 | 1,139 | 1,074 | 1,352 | 1,508 |
| 10 | 1,027 | 1,560 | 1,681 | 1,791 |
| 30 | 1,360 | 1,619 | 1,112 | **0,313** |

**Todo el mapa está entre 0,95 y 1,79 — es decir, sin rigidez. Salvo una esquina.**

---

## 4 · La esquina: su tercer régimen existe

En **A = 30 y B = 32 — los dos altos a la vez**:

- **α = 0,313** contra ~1 en todo el resto del mapa
- **niveles pegados = 0,00 %** — cero colisiones
- **Σ²(20) = 16,81 es MENOR que Σ²(10) = 17,81**: la varianza **deja de crecer**

Eso es rigidez en el sentido estricto: los niveles se acomodan entre sí a lo largo
de la escalera, no sólo con el vecino.

**Y no lo consigue ninguna de las dos perillas sola.** Con A = 30 y B chico, α se
queda en 1,36. Con B = 32 y A chico, en 1,01. **La transición pide las dos arriba
al mismo tiempo** — que es, palabra por palabra, su «mezcla equilibrada de ambos».

Sus tres regímenes, medidos:

| régimen | punto | Σ²(10) | α |
|---|---|---:|---:|
| **I · tensión dominante** (A alto, B bajo) | A=10, B=0,5 | 4,103 | 1,027 |
| **II · flujo dominante** (A bajo, B alto) | A=1, B=8 | 5,028 | 0,968 |
| **III · los dos altos** | A=30, B=32 | 17,810 | **0,313** |

Los regímenes I y II dan la **menor varianza** y ninguna rigidez. El III da la
**rigidez** y paga en varianza. Son fases distintas, no puntos mejores o peores de
una misma.

---

## 5 · Honestidad: lo que cuesta

En la esquina rígida, Σ²(10) vale 17,81 en valor absoluto — **peor** que el mejor
punto del mapa (A = 3, B = 0,5: Σ²(10) = 4,013) y muy lejos de los ceros (0,3364).

**Se gana rigidez y se paga en varianza.** No es todavía el espectro de los ceros.
Es una **fase distinta**, que en ninguno de los mapas anteriores había aparecido.

Contra sus predicciones falsadoras (§10): la separación de A y B **sí** produce
algo que una sola variable no daba; la mejora **no** depende sólo del total; y el
efecto **no** es sólo repulsión local — es justamente el crecimiento a distancia
lo que cambia. La hipótesis **gana evidencia**.

---

## 6 · Lo que la Fase VII tiene que preguntar

El mapa ya no es de un óptimo: **es de una frontera de fase, y hay que
recorrerla.**

1. **¿Existe un camino DENTRO de la esquina rígida que baje la varianza sin
   perder la rigidez?** Refinar la malla alrededor de (A alto, B alto) y seguir la
   curva de α constante, mirando qué pasa con Σ².
2. **¿Dónde está exactamente la frontera?** Medir α en una malla fina para ver si
   la transición es abrupta (frontera de fase) o suave. Si es abrupta, hay un
   mecanismo detrás y hay que nombrarlo.
3. **¿Qué (A,B) da la aritmética por sí sola?** Si el acoplamiento entre las
   excitaciones p y q sale de un objeto aritmético, A y B **no son libres**: están
   determinados. Medir dónde cae ese punto en el mapa es R6-limpio y decide si la
   aritmética elige, sola, el régimen rígido.

La tercera es la que importa: **el mapa dice dónde hay que estar; la aritmética
dirá si la naturaleza está ahí.**

---

## 7 · Autoría y sello

La intuición de los tres regímenes es de **Jesús Nicolás Astorga**. La ruta por
fases y las predicciones falsadoras son de la auditora. La identificación de las
dos variables reales, el mapa y la medición del exponente de rigidez son de este
taller.

Nivel de evidencia: **2** — un mapa reproducible que identifica una fase nueva y
deja tres mediciones concretas.

Nada de esto demuestra ni refuta la Hipótesis de Riemann. **Todavía no.**
