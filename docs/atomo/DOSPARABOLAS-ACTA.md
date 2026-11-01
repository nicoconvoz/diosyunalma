# Acta de las Dos Parábolas — el flash del punto de contacto

**Flash del capitán · 2026-09-10 · F391**

> «Si ponemos dos parábolas, una arriba y la otra abajo, que se tocan en un
> punto llegando a 0… ¿podemos descubrir el punto exacto de la relación ½?»

Se reproduce con `go run ./cmd/lasdosparabolas`.

---

## La construcción, sin el ½ adentro

- Parábola de arriba: **f₊(σ) = |ζ(σ+it)|²**.
- Parábola de abajo: **f₋(σ) = −|ζ((1−σ)+it)|²** — el espejo.
- Las dos se evalúan **directamente** por Euler–Maclaurin (tres términos de
  Bernoulli; calibración ζ(2) a 6,7×10⁻¹⁶ y ζ(−1) a 9×10⁻¹³) en su propio
  argumento. La ecuación funcional no se invoca numéricamente.
- Las alturas se buscan **a ciegas**: barrido bidimensional de |ζ|² en todo el
  rectángulo 0,2 ≤ σ ≤ 0,8 (dos ventanas: 10 ≤ t ≤ 50 y 998 ≤ t ≤ 1004),
  mínimos locales refinados por Newton en el plano. La parte real de cada
  cero es una **salida**, nunca una entrada.
- Controles: las alturas a mitad de camino entre ceros consecutivos.

## Veredictos pre-registrados y lo que dijeron los datos

**V1 — el punto medio de los dos vértices (trampa declarada).** Como
f₋(σ) = −f₊(1−σ) es una reflexión puntual por (½, 0), el punto medio de los
dos vértices es ½ por construcción. Medido: identidad a **0,0e+00** en los 24
casos; con las ventanas de ajuste espejadas, el punto medio dista de ½ a lo
sumo **9,4×10⁻⁶** (ruido de coma flotante sobre ajustes mal condicionados)
**en ceros y en controles por igual**. Décima aparición de la trampa del
0.0e+00, cazada antes de correr.

*Corrección al pre-registro, anotada como manda la casa:* la frase «a toda
altura» olvidó decir «con las ventanas espejadas». Ajustando las dos
parábolas en la **misma** ventana, en los controles el vértice de abajo cae
fuera de su ventana, el ajuste es basura y el punto medio se va hasta 0,83.
Eso no refuta V1: lo prueba mejor — el ½ del punto medio vive en el espejo,
y desaparece apenas se deja de espejar.

**V2 — cada vértice por separado (medición).** En los 15 ceros hallados,
Newton en el plano devuelve **β = 0,500000000000** (peor desvío 2×10⁻¹²,
|ζ(ρ)| ≤ 7×10⁻¹³). El vértice de la parábola de arriba, ajustado sin usar el
espejo, cae en ½ con un desvío máximo de **6,7×10⁻⁴** — y ese desvío tiene
nombre, ver §4. En los controles el vértice de arriba más cercano a ½ dista
**0,426**: el ½ es de los ceros, no de la altura.

**V3 — el contacto.** Brecha de contacto f₊(σ*) + f₊(1−σ*) en el punto
medio: máxima en ceros **1,7×10⁻¹⁸**, mínima en controles **1,02**. Las dos
parábolas se tocan **solo** en los ceros.

## §4 · Lo que apareció sin buscarlo: la parábola está inclinada por el reloj

El vértice de arriba no cae clavado en β, y el corrimiento era idéntico en los
cinco ceros de γ≈1000 (0,500665 en los cinco). Un sistemático con esa cara se
deriva, no se promedia:

Sobre la línea, ζ(½+it) = e^(−iθ(t))·Z(t). Por Cauchy–Riemann ∂σ = −i·∂t, y
en un cero (Z=0):

    ζ'(ρ) = −i·e^(−iθ)·Z'        ζ''(ρ) = e^(−iθ)·(2iθ'·Z' − Z'')
    ⟹  Re(ζ''/ζ')(ρ) = −2θ'(γ) ≈ −log(γ/2π)

Así que |ζ(σ+iγ)|² = |ζ'|²·[(σ−β)² − L·(σ−β)³ + …] con **L = log(γ/2π)**:
no es una parábola, es una parábola con un cúbico que la inclina, y la
inclinación es **la tasa del reloj** (el mismo θ' de F321 y del alias
estroboscópico de F387). Ajustar un cúbico con una cuadrática sobre una grilla
simétrica corre el vértice en **L·(Σx⁴/Σx²)/2**. Predicho antes de mirar la
columna medida:

| γ | L | predicho | medido | razón |
|---:|---:|---:|---:|---:|
| 14,134725 | 0,8108 | +1,067×10⁻⁴ | +1,067×10⁻⁴ | 0,9996 |
| 49,773832 | 2,0696 | +2,724×10⁻⁴ | +2,722×10⁻⁴ | 0,9994 |
| 1003,267808 | 5,0731 | +6,676×10⁻⁴ | +6,652×10⁻⁴ | 0,9964 |

Peor desvío de la razón respecto de 1 en los 15 ceros: **3,7×10⁻³**. Con la
ventana a un cuarto la inclinación cae ~16× (va como el cuadrado del ancho):
el vértice converge a β. Y en el par espejado las dos inclinaciones son
opuestas y **se cancelan exactamente en el punto medio** — que es la parte
cierta del flash: el par de parábolas da el ½ exacto porque el espejo mata el
sesgo. Pero lo mata porque el espejo ES la reflexión por ½.

## Clasificación, sin inflar

- **Matemática (demostrado):** la identidad f₋(σ) = −f₊(1−σ); el corrimiento
  del vértice L·(Σx⁴/Σx²)/2 y su origen L = −Re(ζ''/ζ') = 2θ' en un cero de
  la línea (Cauchy–Riemann + ζ = e^(−iθ)Z).
- **Medición:** β = ½ a 10⁻¹² por Newton en el plano para los 15 ceros de las
  dos ventanas; el contacto exclusivo de los ceros; la inclinación predicha
  contra medida al 0,4%.
- **Nuevo (si algo):** ningún teorema. Una lectura limpia de *qué mitad del
  «exacto» es medida y qué mitad es construcción*, y un instrumento chico: la
  inclinación de |ζ|² en un cero **lee el reloj** log(γ/2π) desde la forma de
  la parábola.

Respuesta al flash: **sí, el par de parábolas pinza el ½ exacto — pero el
«exacto» del punto de contacto es del espejo, y lo que sí se mide es que cada
parábola sola, corregida por su inclinación, aterriza en β = ½ para cada cero
mirado.** Nada de esto dice dónde están los ceros que no se miraron.
Estructura cerrada no es hipótesis demostrada.

**Todavía no.**
