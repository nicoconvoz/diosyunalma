# Acta de los Ceros y la Raíz — el espejo del exponente

**Flash del capitán · 2026-09-15 · F394** (sigue a F393)

> «¿Hay alguna relación entre los ceros y la división por la raíz?»

Se reproduce con `go run ./cmd/loscerosylaraiz`.

**Sí, y es un espejo exacto.** F393 dejó establecido que σ = ½ significa que
cada primo entra al producto de Euler **dividido por su propia raíz**, p^(−½).
La pregunta del capitán es qué hace el otro lado. Respuesta: cada cero entra a
la fórmula explícita **multiplicado por la raíz de x**, x^(+½). Mismo exponente,
signo opuesto, a los dos lados de la misma ecuación.

---

## R6 declarado de entrada (no circularidad)

Las alturas γ se calculan de Z(t), que **es** la línea crítica: vienen con el ½
puesto. Por lo tanto **las alturas no son evidencia de nada** acá. Lo que se
mide son los **dos exponentes**, que quedan libres y se barren.

Precisión de las alturas: la aproximación de Riemann–Siegel usada lleva solo el
término de corrección principal, así que a t bajo tiene ~0,0025 de error (la
primera altura sale 14,137197 contra 14,134725 tabulado). Para reconstruir a
ln x ≈ 8 eso son 0,02 radianes de fase: no cambia ninguna conclusión, pero se
anota.

## §1 · La trampa, declarada antes de correr — DECIMOTERCERA

|x^ρ| = x^β es **álgebra pura**. Decir «si β = ½ entonces cada cero pesa √x» es
repetir la definición del módulo, no medir nada. Lo que **no** es álgebra es que
la suma sobre las alturas reales reproduzca ψ(x) − x. Eso es lo que se mide.

## §2 · La fórmula explícita: los ceros reconstruyendo a los primos

Con 649 alturas hasta γ = 1000, evaluando en 1.600 puntos de x entre 200 y 5.000:

| | |
|---|---:|
| correlación ψ(x) − x contra −Σ x^ρ/ρ | **0,9804** |
| residuo RMS | 2,366 |
| amplitud típica de ψ(x) − x | 10,900 |
| **control: alturas al azar en el mismo rango** | **−0,0323** |

Las alturas verdaderas reconstruyen a los primos; unas cualquiera no reconstruyen
nada. La fórmula explícita es un **teorema** (Riemann 1859, von Mangoldt 1895):
acá solo se verifica que la cuenta cierra y que el control no sobrevive.

## §3 · El exponente, dejado libre

β se barrió de 0,30 a 0,70 **sin imponer nada**, midiendo qué valor reconstruye
mejor a los primos:

| β | residuo RMS | correlación |
|---:|---:|---:|
| 0,30 | 8,814 | 0,9752 |
| 0,40 | 6,305 | 0,9791 |
| 0,48 | 2,823 | 0,9803 |
| **0,50** | **2,366** | **0,9804 ← MÍNIMO** |
| 0,52 | 3,030 | 0,9803 |
| 0,60 | 13,369 | 0,9795 |
| 0,70 | 42,562 | 0,9768 |

Afinado en pasos de 0,002: **β = 0,5000**. Y el pozo es agudo — en 0,70 el
residuo es **18 veces peor**. El ½ no es decorativo: es lo que sostiene el
balance.

**Lectura honesta:** esto no demuestra que los ceros estén en ½ (R6: las alturas
ya vienen de la línea). Lo que muestra es que el **exponente que casa ceros con
primos es ése y no otro**, y que el ajuste lo elige solo.

## §4 · El otro lado, medido aparte

Proyectando log|Z(t)| sobre los tonos ln p en la ventana [10000,5 · 11000,5]:

| primo | amplitud medida | 1/(2√p) predicho | razón |
|---:|---:|---:|---:|
| 2 | 0,352874 | 0,353553 | 0,9981 |
| 3 | 0,288133 | 0,288675 | 0,9981 |
| 7 | 0,189190 | 0,188982 | 1,0011 |
| 23 | 0,104190 | 0,104257 | 0,9994 |
| 37 | 0,082259 | 0,082199 | 1,0007 |

**Exponente ajustado: amplitud ∝ p^(−0,5069)**, o sea α = 0,5069 contra ½. El
tambor de F380–F382 ya sabía que la amplitud vale Λ(m)/(2√m·ln m); lo nuevo acá
es que se ajusta el **exponente solo**, sin suponerlo.

## §5 · La respuesta

| lado | qué hace con la raíz | exponente medido |
|---|---|---:|
| cada **primo** | entra **dividido** por su propia raíz: p^(−½) | **α = 0,5069** |
| cada **cero** | entra **multiplicado** por la raíz de x: x^(+½) | **β = 0,5000** |

**Mismo exponente, signo opuesto.** Los primos se achican por la raíz; los ceros
se agrandan por la raíz. Y es el mismo ½ porque la ecuación funcional ata σ con
1−σ: **lo que un lado divide, el otro lo multiplica**, así que ½ es el único
punto de equilibrio. Por eso mover un cero fuera de la línea rompería el balance
en los dos lados a la vez.

## Clasificación

- **Matemática (demostrado, y de otros):** |x^ρ| = x^β (álgebra); la fórmula
  explícita (Riemann 1859, von Mangoldt 1895); la amplitud del tambor
  Λ(m)/(2√m·ln m); la ecuación funcional.
- **Medición (nuestra):** correlación 0,9804 con control en −0,0323; el barrido
  libre de β con mínimo en 0,5000 y pozo 18× a 0,70; el ajuste independiente
  del exponente del lado primo, α = 0,5069.
- **Nuevo (si algo):** ningún teorema. Los **dos exponentes medidos por
  separado**, uno de cada lado de la ecuación, cayendo en el mismo ½ — y el
  barrido que muestra que ese ½ es estructural y no decorativo.

Estructura cerrada no es hipótesis demostrada.

**Todavía no.**
