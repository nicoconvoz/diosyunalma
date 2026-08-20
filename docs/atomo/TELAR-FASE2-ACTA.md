# Acta del Telar — Fase II

**Respuesta al Pliego del Átomo Fase II (Auditoría 35) · 2026-08-17**

Su pregunta final, textual:

> ¿Podemos construir, desde datos aritméticos/geométricos permitidos por R6,
> una estructura cuya geometría efectiva crezca como ln(E/2π) y cuyo espectro
> produzca un eco en k·ln p, sin introducir los γₙ en ningún punto de la
> definición?

La partimos en dos mitades y medimos cada una por separado, con el protocolo que
usted fijó en su §4 y las cuatro capas que separó en su §5. Todo lo que sigue se
reproduce con `go run ./cmd/eltelar` sobre 620 ceros que el propio programa
encuentra en t ∈ [100, 1000] y que usa **sólo como regla**: ningún candidato los
mira.

---

## 0 · Antes que nada: una corrección a nuestra propia respuesta de Fase I

En el acta anterior escribimos que el requisito que impide la construcción es
**R3 contra (R1 + geometría fija)**. Al armar el mapa de familias que usted pidió
en su §7, esa respuesta **no sobrevive**.

El Hamiltoniano compacto de **Berry & Keating 2011**, `(x+1/x)(p+1/p)`
(*J. Phys. A* **44** (2011) 285203), tiene órbitas acotadas, una familia de
realizaciones autoadjuntas sobre el semieje, **espectro discreto real**, y los
dos primeros términos de la densidad asintótica iguales a los de los ceros.
Dominio fijo. Geometría fija. Y **R6-limpio**: en su construcción no entra un
solo cero. Queda fuera del no-go de Endres–Steiner porque **no es H_BK** — su
símbolo no es lineal en p.

O sea: *discreto + densidad logarítmica + R6-limpio* **es alcanzable**. La
densidad no era el cuello de botella. Lo que a esa familia le falta —y a las
otras dos que llegan al mismo lugar por caminos distintos, Bolte–Egger–Keppeler
(*J. Phys. A* **50** (2017) 105201) y Sierra–Rodríguez-Laguna (*PRL* **106**
(2011) 200201)— es **la aritmética**: no hay primos en ningún lado, y por lo
tanto no hay eco.

**El cuello de botella es el eco.** Esa es la corrección, y cambia hacia dónde
mira la Fase III.

El mapa completo, con las hipótesis exactas que sostienen cada veredicto y otras
cuatro sobre-extensiones de nuestra acta de Fase I puestas en exhibición, está en
[`docs/atomo/TELAR-MAPA-FAMILIAS.md`](TELAR-MAPA-FAMILIAS.md).

---

## 1 · La primera mitad: la caja que respira, sola

Construimos el espectro que implica la ley suave `N(T) = θ(T)/π + 1`, cuya
densidad es exactamente `(1/2π)·ln(T/2π)`: la caja cuyo largo logarítmico crece
como usted pidió. En el tramo medido, ln(T/2π) va de **2,767 a 5,070**.

| capa | resultado |
|---|---|
| **conteo** | 620 niveles contra 620 ceros — diferencia **0** ✓ (por construcción) |
| **correlaciones** | Σ²(10) = **0,0017** · Σ²(20) = 0,0017 — los ceros dan **0,3364** |
| **identidad** | \|nivel − γ\| medio **0,2996**, peor 0,9946 |
| **aritmética** | eco **1,381** (nada); en períodos al azar 0,572 |

**La geometría correcta es necesaria y VACÍA.** Acierta el conteo exacto y no
sabe absolutamente nada más: es una reja demasiado rígida (Σ² dos órdenes por
debajo de la realidad), no tiene estructura fina, y no oye un solo primo.

Y hay que decir lo que su propio contrato R6 nos hace decir: **esa caja está
ESTIPULADA, no derivada.** Le pusimos `ln(T/2π)` a mano. Mientras la ley no se
derive de algo que no sean los ceros, el candidato cae en la regla de parada que
agregamos al contrato, §14.6 **TAUTOLÓGICO**.

---

## 2 · La segunda mitad: los primos solos

Ahora la mitad aritmética, sin mirar un solo cero:

    N_p(T) = θ(T)/π + 1 − (1/π)·Σ_{n = p^k ≤ P} Λ(n)·sin(T·log n)/(√n·log n)

Entradas: la ecuación funcional (a través de θ) y Λ(n) hasta un tope. Nada más.

| primos hasta P | términos | \|niv−γ\| medio | peor | Σ²(10) | eco k·log p |
|---:|---:|---:|---:|---:|---:|
| 10 | **7** | 0,09441 | 0,35885 | **0,3314** | 9,03 |
| 100 | 35 | 0,04305 | 0,15312 | 0,3482 | 18,63 |
| 1 000 | 193 | 0,01836 | 0,07562 | 0,3429 | 35,53 |
| 10 000 | 1 280 | 0,01064 | 0,04601 | 0,3466 | 59,37 |
| 100 000 | 9 700 | **0,00777** | 0,03379 | **0,3364** | 76,78 |
| *ceros verdaderos* | — | — | — | *0,3364* | *182,07* |

Dos cosas para leer con cuidado.

**La primera, que no esperábamos.** Con **siete términos** —los primos hasta
10— la varianza de número ya vale 0,3314 contra 0,3364 de los ceros verdaderos.
La **rigidez** del espectro, eso que ninguna matriz al azar tiene y que la caja
sola no produce, aparece con un puñado de primos. No hace falta el mar entero:
hace falta el principio del mar.

**La segunda, y es una tautología que declaramos nosotros.** La columna del eco
**no informa**. Este candidato lleva `log n` en su definición, así que sus picos
en `k·log p` están garantizados por construcción — es exactamente la advertencia
que quedó escrita en el contrato R6 (§14.6). Las columnas que sí informan son
**Σ²** y **la distancia a los γₙ**, y ésas no se pusieron a mano.

Con el tope grande: **620 de 620 niveles caen dentro de un décimo del espaciado
medio, y 529 (85,3 %) dentro de un centésimo.**

---

## 3 · La respuesta a su §15

**Sí, y con una aclaración que la vuelve menos y más de lo que parece.**

Existe una estructura construida sólo con datos admisibles —ecuación funcional y
Λ(n)— cuya geometría efectiva crece como ln(E/2π) por construcción, cuyo espectro
reproduce las correlaciones de los ceros y cuyos niveles caen sobre los γₙ con el
error medido arriba. **Ningún cero entra en su definición.**

Pero lo que construimos es **una inversión de la fórmula explícita, no un
operador**. Produce una LISTA de niveles, no un mecanismo autoadjunto con espacio
de Hilbert y dominio. En su escala de evidencia (§12) eso es **Nivel 2, y a lo
sumo 3**: una restricción que descarta familias y una construcción definida sin
γₙ — pero no un operador cuyo espectro se demuestre, y mucho menos una
consecuencia formal para RH.

Lo que sí queda establecido, y es lo útil para la Fase III:

- **la geometría sola no alcanza** (medido: la caja acierta el conteo y nada más);
- **la aritmética sola alcanza para las correlaciones y para la identidad**
  (medido: Σ² con 7 términos, identidad con 9 700);
- **por lo tanto, lo que falta no es información: es MECANISMO.** El objeto que
  hay que encontrar no tiene que descubrir dónde están los ceros —los primos ya
  lo dicen— sino *ser* la máquina que hace ese cálculo por sí sola.

---

## 4 · El signo, medido en vez de argumentado

Su §9 pide tratar el signo como restricción de diseño. Lo convertimos en una
**medición**. El coeficiente de Fourier de la densidad de niveles en el período
τ, con la parte suave restada, tiene signo propio: positivo en un sistema tipo
Selberg, negativo en los ceros.

Medido sobre nuestros propios 620 ceros, en los 13 períodos aritméticos hasta
T = 3,2:

| período | n = p^k | D(τ) medido | −Λ(n)/(π√n) |
|---:|---:|---:|---:|
| 0,69315 | 2 | −0,109240 | −0,156013 |
| 1,09861 | 3 | −0,141371 | −0,201899 |
| 1,60944 | 5 | −0,160424 | −0,229108 |
| 1,94591 | 7 | −0,163946 | −0,234112 |
| 2,39790 | 11 | −0,161184 | −0,230136 |
| … | … | … | … |
| 3,13549 | 23 | −0,145710 | −0,208110 |

**13 de 13 negativos**, y siguiendo la predicción con una razón constante (la
pérdida de la ventana). El espectro de absorción no es una interpretación que
tomamos prestada: **está en nuestros datos.**

El análisis completo —las dos fórmulas en una sola normalización, la razón de
amplitudes −2·(1 − p^{−k}), la atribución correcta (Hejhal 1976, Berry 1986; la
lectura es de Connes), la explicación por Lefschetz en cuerpos de funciones, y el
checklist que un candidato debe pasar— está en
[`docs/atomo/TELAR-EL-SIGNO.md`](TELAR-EL-SIGNO.md).

---

## 5 · Las entregas que pidió en su §13

| pedido | entregado |
|---|---|
| Definición formal de R6 y del conjunto D | [`TELAR-R6-CONTRATO.md`](TELAR-R6-CONTRATO.md) — R6 como propiedad de **procedimiento**, no de información; seis clases de violación; D decidido ítem por ítem (ζ **partido por canal**); cuatro auditorías completas |
| Protocolo definitivo del Murciélago | implementado en `cmd/eltelar` con los cinco controles, **más la regla nueva**: no informa sobre candidatos con log p en la definición |
| Implementación independiente de Σ²(L) | `cmd/eltelar` (segunda implementación, independiente de la de `cmd/elpliego`) |
| Tabla L(E) vs ln(E/2π) | medida en Fase I; acá se muestra que **acertarla no alcanza** |
| Mapa de familias Berry–Keating | [`TELAR-MAPA-FAMILIAS.md`](TELAR-MAPA-FAMILIAS.md) — 15 entradas clasificadas, con la hipótesis exacta de cada veredicto |
| Al menos una definición candidata del Telar | §1–§2 de esta acta: el par (ley suave, término oscilante de primos). Se entrega **con su límite declarado**: es fórmula, no operador |
| Análisis del signo | [`TELAR-EL-SIGNO.md`](TELAR-EL-SIGNO.md) + la medición del §4 |
| Auditoría de entradas de cada candidato | §1 de `cmd/eltelar`, impresa antes de cualquier prueba |
| Registro de qué está demostrado / medido / hipótesis | esta acta, y la corrección del §0 |

---

## 6 · Lo que la Fase III debería atacar

De todo lo medido, el frente queda en un solo lugar: **hay familias que ya tienen
espectro discreto, dominio fijo y densidad logarítmica sin leer un cero, y no
tienen aritmética; y hay una construcción aritmética que tiene las correlaciones
y la identidad, y no es un operador.** Las dos mitades existen por separado.

Las puertas que el mapa deja abiertas, en orden:

1. **El espacio de clases de adeles de Connes** — la única construcción donde el
   tamaño, el signo y los primos aparecen a la vez, R6-limpia en su construcción.
   Su fórmula de traza **equivale** a RH: es una reducción, no una demostración.
2. **La positividad de Weil por la vía de Suzuki** (arXiv:2606.09096, 2026) —
   realizaciones **no locales** de −i·d/dx sobre intervalos finitos, que es la
   formalización publicada más cercana a nuestra «caja que respira» y esquiva a
   Endres–Steiner por dos vías distintas.
3. **F13, los primos como espejos** (Sierra, arXiv:1404.4252) — geometría
   R6-limpia derivada de los primos, con la violación de R6 **localizada** en una
   sola fase. Un defecto localizado es un defecto reparable.

---

## 7 · Sello

Nada de esto demuestra ni refuta la Hipótesis de Riemann. Es una revisión de
especificación con mediciones, y una corrección a nuestra propia respuesta
anterior. La regla de la casa preside: **estructura cerrada no es hipótesis
demostrada.**

**Todavía no.**
