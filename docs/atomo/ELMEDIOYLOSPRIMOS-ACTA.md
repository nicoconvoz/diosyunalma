# Acta del Medio y los Primos — la trampa y el puente, separados

**Flash del capitán · 2026-09-15 · F393**

> «¿Qué relación hay entre el ½ de la parábola y la división por dos del número
> primo que da una fracción? Y la relación de la división por números primos
> del mismo número primo… ¿qué relación hay en todo eso?»

Se reproduce con `go run ./cmd/elmedioylosprimos`.

La pregunta mezcla **dos cosas que hay que separar**: una trampa que la casa ya
había cazado, y el puente más real que existe entre el ½ y los primos.

---

## §1 · La trampa, medida contra su control — DUODÉCIMA aparición

«Dividir el primo por dos da una fracción», y esa fracción es siempre ½.
Verdadero, y sin una sola excepción. Medido hasta 10⁶:

| conjunto | cuántos | con frac(n/2) = ½ exacto | % |
|---|---:|---:|---:|
| primos impares | 78.497 | 78.497 | **100,0000** |
| **compuestos impares** | 421.502 | 421.502 | **100,0000** |

**Señal/silencio = 1,0000.** No es una propiedad de los primos: es la
**definición de impar**. Todo primo salvo el 2 es impar porque si fuera par lo
dividiría el 2.

**Territorio ya tenido:** F272 «El uno que sobra» cazó exactamente esto como la
**séptima** aparición de la trampa del 0.0e+00, con el mismo control. Esta es la
duodécima de la casa y la segunda vez que cae el mismo flash — se registra
igual, porque la regla dice que lo que se mide se anota.

## §2 · «Dividir el primo por primos»

De todas las divisiones p/q con q primo ≤ p, da entero **exactamente una**:
p/p = 1. Medido en los primeros ocho primos, y es la **definición de primo**
(no lo divide nadie salvo él mismo y el 1), no un hallazgo.

Lo que sí construye algo es la fracción que **sobrevive** a cada primo:
(p−1)/p, que es la criba. El primer primo da (2−1)/2 = **½**, y ese ½ es el
mínimo sobre todos los primos que existen — pero eso ya es **F276**, registrado.

## §3 · LA RESPUESTA: qué significa σ = ½ en términos de dividir por primos

ζ(s) = ∏ₚ (1 − p^(−s))^(−1). Cada primo entra en el producto **pesando p^(−σ)**:

| σ | lo que pesa cada primo | el 2 | el 3 | el 97 |
|---|---|---:|---:|---:|
| 1,0 | 1/p — el borde | 0,500000 | 0,333333 | 0,010309 |
| 0,7 | 1/p^0,7 | 0,615572 | 0,463463 | 0,040669 |
| **0,5** | **1/√p ← LA LÍNEA** | **0,707107** | **0,577350** | **0,101535** |
| 0,3 | 1/p^0,3 | 0,812252 | 0,719223 | 0,253494 |

**σ = ½ es exactamente «dividir cada primo por su PROPIA raíz cuadrada».** El 2
pesa 1/√2, el 3 pesa 1/√3, el 97 pesa 1/√97. Ésa es la relación directa entre
el ½ y dividir por primos — y no tiene nada que ver con dividir por 2.

## §4 · De dónde sale el 1, de dónde el 0, y por qué el medio cae en ½

Σₚ p^(−σ) hasta 10⁷: diverge en σ ≤ 1 (3,04 en σ=1 y sigue subiendo), converge
en σ > 1 (2,02 en 1,1). Cada eslabón con su dueño:

- **σ > 1**: el producto de Euler converge y ningún factor es cero ⟹ **ζ no
  tiene ceros ahí**. Ese borde en 1 lo ponen **los primos**: Σ 1/p diverge
  (Euler, 1737) — son bastantes.
- **σ < 0**: la ecuación funcional refleja lo anterior ⟹ solo quedan los ceros
  triviales. Ese 0 lo pone **la simetría** s ↦ 1−s, no los primos.
- **0 ≤ σ ≤ 1**: la franja que sobra, y **½ es el único punto que es su propio
  espejo**.

**½ = (0 + 1)/2: el 1 lo ponen los primos, el 0 lo pone la simetría, el ½ es el
punto fijo.**

## §5 · El puente hondo: leer el ½ en los primos, sin mirar un solo cero

Von Koch, 1901: **RH ⟺ π(x) = Li(x) + O(√x·log x)**. El ½ del exponente de √x
**es** el ½ de la línea crítica: el supremo de las partes reales de los ceros
vale exactamente el exponente del error al contar primos.

Medido con la criba hasta 10⁷, ajustando el pico de |ψ(x) − x| por bloques
diádicos:

| medición | θ |
|---|---:|
| exponente crudo, de los primos solos | **0,4809** |
| control de Cramér (moneda 1/ln n) | **0,4888** |
| corrigiendo por log²x | 0,3091 — **sobre-corrige** |

**Calibración del instrumento** (sin esto el número no vale nada): con un sesgo
x^0,70 plantado adrede mide **0,6995**; con uno de x^0,90 mide **0,8996**.

### Hasta dónde llega esto, y ni un paso más

- El θ crudo (0,4809) y el de la moneda (0,4888) son **el mismo número**: los
  primos se desvían como si tiraran una moneda. Ésa es la lectura fuerte.
- La corrección por log²x sale por debajo de ½ y **hay que decirlo**: el log²x
  es la cota del peor caso de RH, no el tamaño típico, y a 10⁷ se lleva puesta
  parte de la señal.
- El instrumento distingue 0,70 de 0,90 sin dudar, pero a 10⁷ **no puede
  distinguir ½ de 0,52**. Esta medición no puede confirmar RH ni de casualidad;
  solo podría **refutarla**, y no lo hizo.
- **θ = ½ no es una conclusión de esta tabla: ES la Hipótesis de Riemann.**

## Clasificación

- **Matemática (demostrado, y de otros):** p impar ⟺ frac(p/2) = ½; p/q entero
  solo en p/p; el producto de Euler y la ausencia de ceros en σ > 1 (Euler 1737);
  la ecuación funcional y el punto fijo ½; la equivalencia de Von Koch (1901).
- **Medición (nuestra):** el control que mata la trampa (100 % contra 100 %);
  las sumas Σp^(−σ); θ = 0,4809 contra la moneda 0,4888, con calibración.
- **Nuevo (si algo):** ningún teorema. La **separación limpia** entre el ½ que
  es tautología y el ½ que es aritmética, con el control que lo decide. Y la
  frase corta que contesta el flash: **dividir por dos no dice nada; dividir por
  la raíz lo dice todo.**

Estructura cerrada no es hipótesis demostrada.

**Todavía no.**
