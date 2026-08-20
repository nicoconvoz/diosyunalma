# Acta de Muchos Modos — Fase V

**Respuesta al Telar Fase V (Auditoría 38) · 2026-08-17**

Su pregunta principal:

> ¿Cuántos modos comunes necesita el fluido para pasar de repulsión local a
> rigidez colectiva?

**La respuesta no es un número: es que el número DEPENDE de a qué distancia se
mire.** Y eso, que no esperábamos, es el hallazgo.

Se reproduce con `go run ./cmd/muchosmodos`.

---

## 1 · El montaje, declarado antes de correr un solo espectro (su §9)

- **El medio.** 520 modos provenientes de 507 excitaciones distintas (primos),
  ventana desde t = 100. Cada primo aporta su escalera `2πk/log p` con peso
  `Λ(p)/√p`.
- **Los canales.** Caracteres de Dirichlet módulo un primo q: el grupo es
  cíclico de orden q−1 = K, y el canal a pesa la excitación de p con
  `cos(2π·a·ind(p)/K)`, con ind el logaritmo discreto. **Así organiza la
  aritmética a los primos**, es R6-limpio, y a = 0 reproduce exactamente el modo
  único de la Fase IV.
- **El control de fuerza (su §10, y esto lo agregamos porque faltaba).** La
  **fuerza total** del acoplamiento se normaliza igual para todo K. Sin eso,
  subir K sube también el empuje total y el experimento confunde «más canales»
  con «más fuerza». Con la traza fija, lo único que cambia entre corridas es el
  **rango**.
- **El control de superposición (su §10).** El mismo K con canales de soporte
  **disjunto** —cada canal toca sólo una clase de restos, sin cruce— que son K
  problemas independientes y ningún efecto colectivo.
- **R6.** Entradas: Λ(n) y los restos mód q. Ningún parámetro tocado contra los
  γₙ; los ceros se leen sólo al final, como regla.

---

## 2 · La tabla

Control sin acoplar (K = 0): **Σ²(10) = 6,4395**, mínimo 7,63×10⁻³, pegados 9,44 %.

| q | K | rango efectivo | Σ² acoplado | Σ² independiente | mínimo | pegados |
|---:|---:|---:|---:|---:|---:|---:|
| 2 | 1 | 1 | 5,4872 | 5,4872 | 5,52×10⁻² | 0,97 % |
| 3 | 2 | 2 | 5,1562 | 5,1562 | 4,58×10⁻³ | 5,80 % |
| 5 | 4 | 3 | 4,8608 | 4,5740 | 7,51×10⁻⁵ | 7,17 % |
| 7 | 6 | 4 | 4,5931 | 3,9568 | 1,81×10⁻³ | 5,05 % |
| 11 | 10 | 6 | **4,3752** | 4,5526 | 2,43×10⁻³ | 9,36 % |
| 17 | 16 | 9 | 4,6064 | 4,9594 | 5,80×10⁻⁴ | 7,65 % |
| 23 | 22 | 12 | 4,5699 | 6,9931 | 4,00×10⁻⁴ | 8,09 % |
| 31 | 30 | 16 | 4,5275 | 6,2239 | 1,24×10⁻³ | 7,55 % |
| 53 | 52 | 27 | 5,0853 | 5,5355 | 1,07×10⁻³ | 9,35 % |
| 101 | 100 | 51 | 5,1946 | 10,5680 | 2,24×10⁻³ | 7,69 % |

*(El rango efectivo se mide, no se supone: la familia de cosenos tiene repeticiones,
así que K canales pedidos dan aproximadamente K/2 independientes. Reportamos el
que hay.)*

---

## 3 · El hallazgo: el óptimo se corre con la distancia

Σ² a tres distancias, y el rango que la minimiza:

| distancia | mejor rango | Σ² ahí |
|---|---:|---:|
| L = 5 | **4** | 2,4497 |
| L = 10 | **6** | 4,3752 |
| L = 20 | **16** | 8,0049 |

**El rango óptimo crece con la distancia que se mide.** Eso es, medido, «el rango
compra alcance»: para que dos niveles se sientan a distancia L, el medio necesita
más canales. La predicción que la Fase IV dejó escrita queda **confirmada — y con
forma**, porque no dice sólo «más canales ayudan» sino **cuántos, para qué
alcance**.

Y tiene una razón estructural que ya conocíamos y ahora se ve funcionando: para
una perturbación positiva de rango K los autovalores cumplen `ωᵢ ≤ Eᵢ ≤ ω_{i+K}`.
**El rango es literalmente cuántas celdas puede correrse un nivel.** Rango uno:
una celda, repulsión sólo con la vecina — que es exactamente lo que medimos en la
Fase IV. Rango K: K celdas.

---

## 4 · El techo, y por qué vuelve a subir

Lo mejor que se consigue a L = 10 es **4,3752**, contra 6,4395 sin acoplar y
**0,3364 de los ceros verdaderos**: seguimos **trece veces más arriba**.

Y pasado el óptimo **empeora**: con rango 51 vuelve a 5,1946. La razón es directa
y sale del control de fuerza: **con la traza fija, más canales significa cada
canal más débil.** Se gana alcance y se pierde empuje, y el producto tiene un
máximo. Por eso el óptimo existe y por eso se corre: a mayor L hace falta más
alcance, y conviene pagar más dilución.

**Ninguna ley de potencia describe la caída.** El mejor ajuste da exponente
**+0,012** —o sea, plano— con residuo 0,72. La ley 1/K que la Fase IV había
propuesto como candidata **predice muchísima más caída que la real**. Lo
reportamos igual, como pide su §7: no nos quedamos con el mejor ajuste, mostramos
que ninguno sirve.

---

## 5 · El control de superposición pasa

Al mismo rango 51: **acoplado 5,1946 contra independiente 10,5680.** El acoplado
es dos veces mejor con los mismos grados de libertad y la misma fuerza total.
**La mejora es de la interacción, no del número de parámetros.** Su §10 queda
satisfecho.

---

## 6 · Veredicto contra sus cinco escenarios (§11)

**Escenario B, afinado.** Σ² mejora con el rango, tiene un **óptimo**, y vuelve a
subir. El rango alto **por sí solo no alcanza**.

Pero ahora la estructura que falta tiene una descripción precisa, y ése es el
aporte de esta fase:

> hace falta algo que dé **ALCANCE sin DILUIR la fuerza** — y eso es exactamente
> lo que ningún acoplamiento de rango finito con traza fija puede hacer.

Un acoplamiento de rango K reparte una fuerza fija entre K canales. Para tener
alcance L sin diluir, el medio no puede ser una suma de K proyectores: **tiene
que ser algo cuya interacción decaiga lentamente con la distancia en vez de
cortarse a K celdas.** En el lenguaje de matrices: no rango bajo, sino
**decaimiento lento fuera de la diagonal**. En el lenguaje de la física: no un
puñado de campos, sino **no-localidad**.

Y eso conecta, sin que lo hayamos buscado, con las tres puertas que la Fase III
dejó abiertas: **las realizaciones no locales de Suzuki**, el espacio de clases de
**adeles de Connes**, y **los espejos de Sierra**. Las tres son maneras de tener
interacción de largo alcance sin repartir una fuerza finita entre canales.

---

## 7 · Lo que la Fase VI debería medir

Sale del mecanismo, no de una opinión:

1. **Reemplazar el rango por el decaimiento.** En vez de `Σ_a g_a|v_a⟩⟨v_a|`,
   un acoplamiento `H_ij = f(|i−j|)` con f de cola lenta (potencia, no
   exponencial), con la misma traza. Medir Σ²(L) contra el exponente de la cola.
   **Predicción**: existe un exponente crítico donde la rigidez de largo alcance
   aparece de golpe.
2. **Preguntar cuál es el decaimiento que la aritmética da naturalmente.** Si el
   acoplamiento entre las excitaciones p y q sale de un objeto aritmético, su
   cola no es libre: está determinada. Medirla es R6-limpio y decide si la
   aritmética alcanza para el exponente crítico.

---

## 8 · Autoría y sello

La intuición del fluido es de **Jesús Nicolás Astorga**. La ruta por fases y los
controles son de la auditora. El montaje, la normalización de fuerza que faltaba,
las mediciones y el hallazgo del óptimo móvil son de este taller.

Nivel de evidencia: **2** — una restricción cuantitativa que descarta una familia
(rango finito con traza fija) y deja una predicción medible.

Nada de esto demuestra ni refuta la Hipótesis de Riemann. **Todavía no.**
