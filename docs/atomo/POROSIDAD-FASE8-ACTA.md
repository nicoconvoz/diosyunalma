# Acta de la Porosidad — Fase VIII

**Respuesta al Telar Fase VIII (Auditoría 41) · 2026-08-17**

El flash es del capitán: *«¿y si en este sistema hay cierta porosidad: mayor o
menor porosidad, materiales más blandos y otros más duros?»*

La respuesta corta, con todas las letras y antes de cualquier detalle: **el flash
acertó en que el medio no debe suponerse uniforme, y la heterogeneidad mejora de
verdad. Pero la mejora NO es aritmética.** Un campo de dureza ordenado y sin un
solo primo adentro hace exactamente lo mismo. Ése es el resultado negativo de
esta fase y es el importante.

Se reproduce con `go run ./cmd/porosidad`.
Lámina: `galeria/laminas/10-el-telar/porosidad.svg`.

---

## 1 · La traducción matemática (su §7)

Su §7 prohibía usar «duro», «blando» y «poroso» como variables sin definirlas, y
preguntaba si porosidad y dureza son dos magnitudes o una. **Son una.**

Cada sitio lleva una **permeabilidad local** `h_i > 0`, y el núcleo pasa de
depender sólo de la distancia a depender también del lugar:

```
homogéneo    H_ij = f(|i−j|)
heterogéneo  C(i,j) = √(h_i · h_j) · f(|i−j|)
```

La raíz del producto no es adorno: es la única forma factorizada que mantiene la
matriz **simétrica** (y por lo tanto el operador autoadjunto, R1) sin introducir
una dirección privilegiada. `C(i,j) = h_i · f(|i−j|)` no es simétrica y rompe R1
de entrada.

Con eso, las cinco palabras de su §7 colapsan en una sola magnitud leída en dos
sentidos: **poroso / blando = h grande** (la influencia pasa fácil, el sitio
responde mucho); **duro = h chico**. No hacen falta dos variables.

`f` es la cola larga sobreviviente de la Fase VII: `f(k) = A/|k|^s` con `s = 0.5`,
y su variante con nodo `(1 − |k|/k₀)/|k|^s`, `k₀ = 5`.

**La fuerza total está fija en 30 en todos los brazos** (normalización de traza —
el control de su §10 de la fase anterior). Sin eso, un medio heterogéneo tiene
más acoplamiento total que el homogéneo y la comparación no dice nada.

---

## 2 · R6, la parte más importante (su §13)

La regla aritmética se declaró **antes** de calcular un solo espectro:

```
h_i = 1 / log(p_i)
```

El razonamiento es físico, no ajustado: el modo de un primo chico es una onda
larga y floja, y deja pasar; el de un primo grande es corta y rígida. Ningún γₙ
se mira, ni para elegir la regla, ni para elegir las regiones, ni para ajustar
nada. Se corrió además una **segunda** regla aritmética independiente, `h_i =
d(k_i)` (número de divisores), que no está correlacionada con la amplitud del
modo, precisamente para que la conclusión no dependa de una fórmula sola.

Auditoría R6 completa: **limpia**. Los γₙ aparecen una única vez en toda la
corrida, como regla de medir al final: Σ²(10) = 0,3364.

---

## 3 · La matriz de comparación (su §11)

400 modos, misma fuerza total, misma cola larga. Sólo cambia el material.

| medio | vivos | Σ²(5) | **Σ²(10)** | Σ²(20) | α | PR/N |
|---|---:|---:|---:|---:|---:|---:|
| homogéneo, sin nodo | 187 | 5,765 | **18,335** | 62,336 | 1,717 | 0,105 |
| homogéneo, CON nodo | 291 | 2,794 | **5,426** | 9,647 | 0,894 | 0,042 |
| **mezclado**, sin nodo (5 semillas) | 194 | — | **17,501 ± 1,435** | — | — | 0,100 |
| **mezclado**, CON nodo (5 semillas) | 302 | — | **4,899 ± 0,596** | — | — | 0,038 |
| aritmético 1/log p, sin nodo | 219 | 4,740 | **14,236** | 47,647 | 1,665 | 0,080 |
| aritmético 1/log p, CON nodo | 311 | 2,916 | **4,700** | 8,996 | 0,813 | 0,029 |
| aritmético divisores, sin nodo | 197 | 4,713 | **14,374** | 45,822 | 1,641 | 0,065 |
| aritmético divisores, CON nodo | 300 | 2,990 | **5,787** | 9,182 | 0,809 | 0,036 |

El brazo **mezclado** es el control aleatorio de su §6: los mismos valores de `h`,
permutados al azar. Mismo histograma, misma media, misma dispersión — sólo cambia
el orden. Es el único control que puede separar «heterogeneidad» de «aritmética».

**Niveles vivos entre 187 y 311 de 400 en todos los brazos: la banda no se vacía
en ninguno.** Es la disciplina que dejó la Fase VII, aplicada de entrada y no al
final: sin esto, ningún Σ² de esta tabla valdría nada.

---

## 4 · La prueba que decide, y el control que le faltaba

Contra el mezclado, el campo aritmético gana:

| | aritmético | mezclado | distancia |
|---|---:|---:|---:|
| sin nodo | 14,236 | 17,501 ± 1,435 | **2,28 σ** |
| con nodo | 4,700 | 4,899 ± 0,596 | 0,33 σ |

Y aquí es donde esta hoja podría haberse quedado cantando victoria. **No.** Las
dos reglas aritméticas, que no tienen nada que ver entre sí, dieron casi el mismo
número (14,236 y 14,374). Eso huele a que el medio no está sintiendo los primos:
está sintiendo que el campo `h` **tiene orden**, cualquiera sea su origen.

Así que se agregó un control que su documento no pedía: campos igual de ordenados
y **sin nada de aritmética**.

| campo de dureza | Σ²(10) | distancia a la aritmética |
|---|---:|---:|
| aritmético 1/log p | 14,236 | — |
| **rampa lisa** (monótona) | **30,731** | 11,49 σ |
| **onda lisa** (senoidal) | **15,398** | **0,81 σ** |
| mezclado | 17,501 ± 1,435 | 2,28 σ |

**La onda lisa empata al campo aritmético a 0,81 σ sin usar un solo primo.**

Por lo tanto, en el lenguaje de su §16:

- **Éxito experimental inicial: SÍ.** La heterogeneidad mejora Σ²(10) un 22,4 %
  de forma reproducible y sin vaciar la banda.
- **Éxito estructural: NO.** La mejora tampoco aparece en el control aleatorio,
  es cierto — pero sí aparece en un control **ordenado y no aritmético**. La
  condición de su §16 se cumple contra el azar y se cae contra el orden.
- **Éxito aritmético: NO.** Su §6 avisaba de no atribuir a la aritmética una
  mejora que la heterogeneidad sola explica. El aviso pegó.
- **Fracaso útil: exactamente esto.**

Lo que el medio premia es que la dureza tenga **estructura a escala intermedia**.
La aritmética está en esa clase, pero no se distingue de una onda cualquiera.

### 4bis · Y no es «cualquier orden»

Este es el hallazgo lateral que vale la pena guardar: **la rampa monótona es
peor que el azar** (30,731 contra 17,501). Ordenar el medio de blando a duro en
una sola pendiente lo **arruina** — casi el doble de varianza que revolverlo.

Entonces no es «orden bueno, azar malo». Es más fino: el medio necesita variación
a escala intermedia, ni plana, ni revuelta, ni monótona. Un gradiente global
concentra la respuesta en un extremo y la desordena.

---

## 5 · Localización (su §8 y su §9)

Medio con tres zonas —blanda, media, dura— variando sólo el contraste:

| contraste | vivos | Σ²(10) | Σ²(20) | α | **PR/N** |
|---:|---:|---:|---:|---:|---:|
| ×1 (homogéneo) | 187 | 18,335 | 62,336 | 1,717 | 0,105 |
| **×2** | 216 | **7,535** | 16,448 | **1,220** | 0,070 |
| ×5 | 263 | 15,306 | 50,653 | 1,673 | 0,044 |
| ×20 | 281 | 30,456 | 107,071 | 1,778 | 0,024 |
| ×100 | 284 | 39,794 | 140,146 | 1,761 | 0,017 |

**Hay un punto justo.** La porosidad suave (×2) es el mejor resultado sin nodo de
toda la hoja: 7,535 contra 18,335 del homogéneo — mejor incluso que el campo
aritmético. Y α baja de 1,72 a 1,22, que es el movimiento en la dirección buena.

Pasado ese punto, el medio se **tapona**: PR/N cae a 0,017 (estados atrapados en
una región) y Σ² se va a 39,79. Demasiada porosidad no ayuda, aísla.

### La advertencia de su §9, cumplida

Su §9 pedía distinguir alcance real de alcance aparente. La respuesta directa es
la razón de participación PR/N = (1/Σv⁴)/N: vale 1 para estados extendidos por
todo el medio y tiende a 0 para estados atrapados.

**El nodo baja Σ²(10) de 18,335 a 5,426 — pero PR/N cae de 0,105 a 0,042.** Los
estados se atrapan 2,5 veces más. Buena parte de esa mejora del nodo puede ser
**alcance aparente por localización, no alcance real**, y hay que declararlo
antes de festejar el número. Queda como pendiente medible: separar cuánto de la
ganancia del nodo sobrevive a PR/N constante.

---

## 6 · Veredictos separados (su §17)

| efecto | medida | veredicto |
|---|---|---|
| **heterogeneidad** | 18,335 → 14,236 (22,4 %); mejor punto 7,535 con porosidad suave | **real y reproducible**, sin vaciar la banda |
| **nodo** | 18,335 → 5,426 (homogéneo); 14,236 → 4,700 (aritmético) | **el efecto más grande**, pero con PR/N a la mitad: parte puede ser aparente |
| **aritmético** | 2,28 σ contra el mezclado, **0,81 σ contra una onda lisa** | **NO se sostiene** |

Los tres se suman: aritmético + nodo da 4,700, el mejor de la hoja. Pero la
descomposición dice que el mérito es del nodo y de la heterogeneidad, **no de los
primos**.

Y la escala de todo esto: los ceros verdaderos están en Σ²(10) = 0,3364. El mejor
brazo de esta fase está catorce veces más arriba. La distancia sigue siendo
enorme.

---

## 7 · Sobre su §14: «¿y si los primos son parte del material?»

Ésta era la extensión conceptual fuerte, y esta fase la probó en su forma más
directa: los primos construyendo la permeabilidad del medio en lugar de ser las
excitaciones.

**Medido, y no alcanza en esta forma.** Una permeabilidad derivada de los primos
no se distingue de una onda lisa. La cadena `aritmética → medio → propagación →
espectro` no queda refutada como idea, pero sí queda refutada la versión ingenua:
poner `log p` en la dureza local no le transmite al espectro nada que una
modulación suave cualquiera no transmita igual.

Si esa cadena existe, la aritmética tiene que entrar por algún lado que **no** sea
un campo escalar suave de sitio. Ésa es la información nueva.

---

## 8 · Lo que NO se afirma (su §15)

- No se afirma que los primos sean literalmente un material.
- No se afirma que la porosidad explique nada de RH.
- No se afirma que la heterogeneidad produzca GUE: Σ²(10) sigue en 7,5–14, y los
  ceros están en 0,34.
- No se afirma que el punto justo (×2) sea un mecanismo. Es un punto medido en un
  juguete, con una semilla y un tamaño.
- No se afirma que el nodo mejore el alcance real: PR/N dice que hay que dudarlo.

---

## 9 · Lo que queda abierto

1. Separar la ganancia del nodo a **PR/N constante** — el pendiente que deja su §9.
2. Barrer el punto justo con tamaño y semilla: ¿el óptimo en ×2 se mueve con N?
3. Si la aritmética no entra como campo escalar de sitio, ¿entra como **fase**?
   Un `h` complejo, o una modulación en el signo, no fue probado.
4. La onda lisa que empató: ¿su período importa? Si hay un período óptimo, eso es
   una escala del medio y merece medirse.

---

El flash de la porosidad es de Jesús Nicolás Astorga. Los tres controles
obligatorios y la advertencia de la localización, de la auditora. El control
ordenado-no-aritmético que tumbó la atribución, de este taller.

Estructura cerrada no es hipótesis demostrada.

**Todavía no.**
