# Acta de la Rigidez — Fase XIX

**El escalón con dientes · 2026-08-20 · F371**

La brújula era suya: *no aumentar la fuerza — cambiar la dinámica.* Se cambió el
material por GUE genuino de rigidez perfecta, con la regla pre-registrada de que
si no superaba al paseo, la hipótesis moría. **No lo superó. Murió, y se
declara.** Y en el velorio apareció lo importante: el mejor modelo de toda la
campaña estaba escondido en un brazo de control de la Fase XVIII.

Se reproduce con `go run ./cmd/laformulamadre fase19` (semilla 20260826).
Lámina: `galeria/laminas/10-el-telar/la-rigidez.svg`.

---

## 1 · El material, verificado antes del campo

Ensamble tridiagonal de Hermite (β = 2): autovalores **GUE exactos**,
desplegados por la semicircular, sólo el interior del espectro. La rigidez se
midió antes de usarla — varianza de conteo en ventanas de L = 10:

| material | varianza de conteo |
|---|---:|
| paseo iid (Fase XVIII) | 1,82 |
| **GUE rígido (nuevo)** | **0,59** |
| GUE teórico | ~0,59 |
| Poisson | 10 |

El material es el correcto, clavado por la MISMA ecuación de conteo. Cero
perillas.

## 2 · El juicio

| modelo | amplitud | cruce | s desv |
|---|---:|---|---:|
| REAL | 0,346 | 0,862 | 0,42 |
| rígido SIN campo | 0,029 | dispersos | 0,43 |
| **RÍGIDO + campo** | **0,120** | 0,85–0,91 | 0,56 |
| rígido + campo, fase rota | 0,026 | muerto | — |
| paseo + campo (XVIII) | 0,133 | 0,81–0,89 | 0,57 |

**Veredicto por la regla pre-registrada: la hipótesis de la rigidez MUERE.**
0,120 contra 0,133 del paseo — no mejora ni la amplitud, ni el T4 (≈ iguales),
ni el ensanchamiento de espaciados (0,56 contra 0,57: idéntico). El residuo
mejora apenas (0,227 contra 0,250). La coherencia sigue siendo necesaria (fase
rota → 0,026), pero la rigidez de largo alcance del material **no era el
ingrediente**.

## 3 · Lo que la muerte enseña — el doble conteo de la fluctuación

Los cuatro modelos acoplados de la campaña (densidad, a posteriori, paseo,
rígido) comparten dos síntomas que ahora se leen juntos:

1. Todos se quedan en ~⅓ de la amplitud real.
2. Todos ensanchan los espaciados a 0,56–0,58 mientras los ceros reales llevan
   la señal completa manteniendo 0,42 **exacto**.

El diagnóstico: **todos suman una fluctuación GUE independiente ENCIMA del
campo.** Dos fuentes de varianza independientes → espaciados anchos, y el ruido
ajeno diluye la coherencia → amplitud chica. Pero en los ceros reales no hay dos
fuentes: **la fluctuación "GUE" de los ceros ES aritmética** — son los términos
de la misma fórmula explícita que nuestra truncación en n ≤ 97 tira a la basura.
Los primos grandes, sumados de a miles, fabrican la estadística local tipo GUE
(la vieja imagen de Berry), y por eso son coherentes con los primos chicos que
el eco mide: es un solo objeto, no campo más ruido.

## 4 · Y el ganador estaba en un control

Releído con estos ojos, el brazo **F) campo solo** de la Fase XVIII — la
estacada rígida SIN ruido agregado, clavada apenas por N_liso + S — era el mejor
modelo de toda la campaña y nadie lo coronó porque era un control:

| | amplitud | cruce | pendiente | s desv |
|---|---:|---:|---:|---:|
| campo solo (XVIII, control F) | **0,285** | 0,94 | −0,80 | **0,42** |
| el mejor acoplado (XVIII/XIX) | 0,138 | ~0,90 | ~−0,3 | 0,57 |
| REAL | 0,346 | 0,862 | −0,93 | 0,42 |

**Sin ruido agregado: 82 % de la amplitud real Y el desvío de espaciados
clavado en 0,42.** Todo lo que le agregamos encima —paseo o GUE rígido— sólo
lo empeoró.

## 5 · La hipótesis siguiente, sin perillas, para su bendición

**No hay material: el campo ES el material.** El modelo puro
`N_liso(γ) + S(γ) = k − ½` sin ninguna fluctuación agregada, con la truncación
de S **extendida** (n ≤ 97 → n ≤ 997 → n ≤ 9973, barrido declarado como
robustez): la predicción de Berry es que los términos nuevos fabrican solos la
estadística local que hoy ponemos a mano con ruido, y siendo aritméticos, son
coherentes con el eco por construcción. Si al profundizar la truncación la
amplitud sube de 0,285 hacia 0,346 y el cruce baja de 0,94 hacia 0,862 con los
espaciados quietos en 0,42 — el modelo converge al espejo real por el único
camino que queda: **más aritmética, no más ruido.** Si no converge, también lo
sabremos con la misma regla de siempre.

HL y el Telar siguen congelados.

---

La escalera completa, a esta altura: XVI tamaños (⅓) → XVII posiciones después
(peor) → XVIII juntos (mejor, 2,5×) → XIX material rígido (**muerto**) → y el
saldo: el espejo no quiere ruido de ninguna clase — quiere más fórmula.

Estructura cerrada no es hipótesis demostrada.

**Todavía no.**
