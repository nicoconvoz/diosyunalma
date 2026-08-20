# Acta de la Máquina — Fase III

**Respuesta al Telar Fase III (Auditoría 36) · 2026-08-17**

Su regla de oro, textual:

> NO BUSCAR UNA LISTA QUE IMITE LOS CEROS. BUSCAR UNA MÁQUINA QUE NO NECESITE
> CONOCERLOS.

Su §5 fija el primer paso: **de primo a órbita** — decir qué objeto es una órbita
p^k, por qué su período es log(p^k), **derivar** el peso Λ(n)/(√n·log n) y
**derivar** el signo en vez de imponerlo.

Fuimos por ahí. El resultado es en su mayor parte negativo, que es la clase útil
de resultado: **encontramos un obstáculo nuevo que cierra de golpe la familia más
obvia**, y una derivación que conecta las dos mitades que la Fase II había dejado
separadas.

Todo se reproduce con `go run ./cmd/lamaquina`. Única entrada aritmética: Λ(n).
Los γₙ aparecen sólo como regla de medir.

---

## 1 · El obstáculo de la concatenación (no-go nuevo)

**Hipótesis exactas** — y sólo dentro de ellas vale:

1. Las órbitas primitivas del sistema tienen longitudes `{log p}`, una por primo.
2. Dos órbitas cerradas que comparten un punto se pueden **concatenar** en otra
   órbita cerrada.

La hipótesis (2) es cierta en **todo grafo cuántico conexo** y en **todo flujo con
un punto recurrente común** — es decir, en la enorme mayoría de los sistemas
dinámicos que uno construiría.

**Consecuencia inmediata.** Si (1) y (2), el espectro de longitudes es cerrado
bajo la suma. Y como

    log a + log b = log(ab)

el sistema está **obligado** a tener una órbita de longitud `log n` para **todo**
entero n que sea producto de primos, o sea para todo n.

**Pero la fórmula explícita le da peso exactamente CERO a todo n que no sea
potencia de primo:** Λ(6) = Λ(10) = Λ(12) = … = 0.

**Medido:** hasta n = 4096, un sistema que concatena está obligado a llevar
**4 095 longitudes**; la fórmula explícita permite **604** (las potencias de
primo) y **prohíbe 3 491 — el 85,3 %**. Las primeras prohibidas son
log 6, log 10, log 12, log 14, log 15, log 18, log 20, log 21…

**⟹ Ningún sistema cuyas órbitas se puedan concatenar puede tener el espectro de
longitudes de la fórmula explícita.** Sólo quedan dos salidas:

- **las órbitas no se tocan** (y entonces no hay concatenación posible), o
- **las concatenaciones se cancelan exactamente** — y cancelar exactamente para
  los 3 491 compuestos, y para todos los que siguen, no es un detalle técnico:
  es una condición extraordinariamente fuerte que el candidato debería explicar.

Esto responde su §5 de la manera más dura posible: **una órbita p^k no puede ser
una órbita cerrada de un sistema dinámico corriente.** Es exactamente por eso que
la máquina que buscamos tiene que ser rara.

---

## 2 · La única salida, construida y medida

Si las órbitas no se tocan, el objeto es una **unión disjunta de lazos**, uno por
primo, de longitud log p. Su espectro es la unión de las escaleras `2πm/log p`.
Lo construimos con los 78 primos hasta 400:

| capa | resultado |
|---|---|
| **conteo** | densidad **constante** 59,00 por unidad — los ceros piden (1/2π)·ln(T/2π), que crece ✗ |
| **correlaciones** | Σ²(10) = **7,58** contra **0,336** de los ceros; espaciado mínimo 8,5×10⁻⁶ ✗ |
| **aritmética** | eco en k·log p = **198,4** ✓ (los ceros dan 182,1) |

La aritmética es **genuina** — cada escalera *es* un primo, no hay nada puesto a
mano — pero los niveles se amontonan: es una superposición de relojes
independientes, sin repulsión ninguna. **Muere en conteo y en correlaciones.**

Es la contracara exacta de Berry–Keating 2011, que tiene el conteo y las
correlaciones y no tiene aritmética.

---

## 3 · El peso, derivado

Su §12 pide derivar el peso Λ(n)/(√n·log n) en vez de dejarlo como suma de
Fourier. Se deriva, y sale algo lindo.

- Una órbita **estable** (un círculo) aporta peso ℓ y **no se amortigua** con las
  repeticiones. La fórmula explícita sí se amortigua, como p^{−m/2}. **Un círculo
  no puede ser un primo.**
- Una órbita **hiperbólica** de longitud ℓ y exponente de inestabilidad λ aporta
  `ℓ / (2·senh(mλ/2))` en la fórmula de traza, que decae en m con **tasa λ**.
- La fórmula explícita pide `2·log p · p^{−m/2}`, que decae con **tasa log p**.

Igualando **tasas** (medido: el cociente da 1,000000 exacto para m ≥ 8):

> **λ = log p.** El exponente de inestabilidad de la órbita es su propia longitud.

En unidades donde el tiempo es la longitud, eso es **exponente de Lyapunov
exactamente 1** — que es, precisamente, el flujo de `xp` (ẋ = x, ṗ = −p).

**Ese es el puente entre las dos mitades que la Fase II dejó separadas:**
Berry–Keating aporta la inestabilidad correcta, los primos aportan las
longitudes, y la máquina necesita las dos a la vez. No es una analogía: la tasa
lo fuerza.

> **Corrección nuestra, en el mismo turno.** Primero despejamos λ de la amplitud
> **completa** y anunciamos que tendía a log p. Es falso: eso da
> λ = log p − (2·log 2)/m, y el cociente λ/log p ni siquiera es monótono (p=2 da
> 1,000; p=11 da 0,630; p=97 da 0,714). Lo que la fórmula fuerza es la **tasa**,
> no la constante. La constante que sobra es justamente el factor del §4.

---

## 4 · El signo: lo que sobra cuando las amplitudes ya se emparejaron

Con λ = ℓ = log p exacto, la razón entre lo que pide Weil y lo que da Selberg,
medida:

| p | m | Selberg (+) | Weil (−) | razón |
|---:|---:|---:|---:|---:|
| 2 | 1 | +0,98025814 | −0,98025814 | −1,000000 |
| 3 | 1 | +0,95142615 | −1,26856820 | −1,333333 |
| 7 | 1 | +0,85806572 | −1,47096981 | −1,714286 |
| 97 | 2 | +0,04716698 | −0,09432394 | −1,999787 |

La razón es exactamente **−2·(1 − p^{−m})**.

Dos lecturas, las dos importantes:

1. **El signo negativo sobrevive a todo el ajuste.** No sale de emparejar
   amplitudes ni de elegir λ. Hay que **derivarlo de otra cosa** — orientación,
   índice, grado cohomológico, o la estructura de absorción de Connes. Su §6 pide
   convertirlo en propiedad derivada del candidato: acá queda demostrado que
   ningún ajuste de amplitud lo va a producir.
2. **El factor (1 − p^{−m}) es el «problema asintótico» de Sierra, medido:** la
   amortiguación de Selberg y la de Weil coinciden sólo cuando m → ∞. A m finito
   no hay flujo hiperbólico estándar que dé exactamente el peso de Weil.

En la Fase II ya habíamos medido que el signo está **en los datos**: 13 de 13
períodos aritméticos con coeficiente de Fourier negativo. Ahora sabemos además
que **no se puede fabricar desde la dinámica sola**.

---

## 5 · La matriz de evaluación (su §9), llenada

| candidato | R6 | ¿operador? | conteo | correlaciones | eco | muere en |
|---|---|---|---|---|---|---|
| inversión aritmética (F350) | limpio | **NO** | sí | sí | tautológico | no es un operador |
| grafo o flujo que concatena | limpio | sí | — | — | **IMPOSIBLE** | §1: longitudes prohibidas |
| escaleras disjuntas | limpio | sí | **NO** | **NO** | sí | conteo y correlaciones |
| Berry–Keating 2011 | limpio | sí | sí | ? | **NO** | no tiene aritmética |

**El hueco es exactamente uno, y ahora tiene forma:**

> un operador con la **inestabilidad de xp** (Lyapunov 1) cuyas **órbitas cerradas
> sean los primos** y que **no permita concatenarlas**, con el signo derivado de
> su estructura y no elegido.

Las tres condiciones juntas son incompatibles con cualquier grafo, con cualquier
flujo con punto recurrente común, y con cualquier unión disjunta de lazos. Por
eso las tres puertas que quedan abiertas —adeles de Connes, realizaciones no
locales de Suzuki, y los primos como espejos de Sierra— son todas **no locales o
no geométricas en el sentido usual**. Eso ya no es una casualidad: es lo que este
no-go predice.

---

## 6 · Nivel de evidencia y sello

En su escala (§11): esta acta aporta **Nivel 2** — una restricción que descarta
familias, con hipótesis exactas y medición. No aporta Nivel 3 nuevo: no
entregamos un operador. Lo que sí hace es **estrechar el lugar donde puede estar**,
y decir por qué las familias que quedan tienen la forma rara que tienen.

Nada de esto demuestra ni refuta la Hipótesis de Riemann. **Todavía no.**
