# Acta de las Parábolas Centradas — la anatomía del contacto en el (0,0)

**Flash del capitán · 2026-09-10 · F392** (segundo flash sobre F391)

> «Ubicá las dos parábolas justo en el centro, apoyadas en el 0 absoluto, y
> decime la relación de distancia positiva y negativa, qué relación tiene el
> ½ con cada altura, qué ángulo separa la distancia en x y en y, y cuánta
> separación hay del 0 a cada punto y de una mitad de la parábola a la otra.»

Se reproduce con `go run ./cmd/lasparabolascentradas`.

---

## El marco

x = σ − β con β medido por Newton en el plano (nunca supuesto); y = |ζ(β+x+iγ)|²
para la de arriba, y = −|ζ(β−x+iγ)|² para la de abajo (su espejo por el origen).
Quince ceros hallados a ciegas en 0,2 ≤ σ ≤ 0,8, en dos ventanas (10–50 y
998–1004). Todo lo que sigue se predijo con la fórmula de F391 antes de mirar
la columna medida: cerca del cero,

    y = a·x² − a·L·x³ + …      a = |ζ'(ρ)|²      L = log(γ/2π) = 2θ'(γ)

## Las respuestas, una por una

**1 · La relación de distancia positiva contra negativa.** A la misma altura y,
la rama derecha llega más lejos que la izquierda:

    x₊ + x₋ = L·y/a        x₊/|x₋| ≈ 1 + L·√(y/a)

Medido contra predicho en 17 filas: γ=14,13 a y=2,5×10⁻⁴ → x₊+x₋ = +3,243×10⁻⁴
contra +3,243×10⁻⁴; γ=1003 → +2,039×10⁻³ contra +2,029×10⁻³ (0,5 %, el
siguiente orden). La razón derecha/izquierda vale 1,016 en γ≈14 y **1,107**
en γ≈1000, a la misma distancia horizontal w = 0,02: es 1 + L·w.

**2 · Qué relación tiene el ½ con cada altura.** Ninguna: **β = 0,500000000000
en los 15 ceros** (peor desvío 2×10⁻¹²). El punto de contacto no se mueve con
la altura. Lo que sí se mueve es la forma: la inclinación L sube como
log(γ/2π) (de 0,81 en γ=14 a 5,07 en γ=1000, determinista) y la apertura
a = |ζ'(ρ)|² vaga cero a cero (de 0,63 a 21,4, sin ley visible en 15 muestras).

**3 · El ángulo.** En el punto de contacto las dos tangentes son horizontales:
el ángulo entre las dos parábolas ahí es **0**. El único ángulo intrínseco del
par es el del **eje**: los puntos medios de cada cuerda horizontal se corren a
la derecha linealmente con la altura, con pendiente dx/dy = **L/(2a)**. Medido
0,6443 contra predicho 0,6444 en γ=14; 0,1188 contra 0,1182 en γ=999. El ángulo
ψ = atan(L/2a) (32,8° en γ=14, 6,8° en γ=999, en unidades crudas σ contra
|ζ|²) **depende de la escala elegida para y**; la pendiente L/(2a) no. Se
reporta el ángulo porque se pidió, y se dice de qué depende.

**4 · Las separaciones (ventana x = ±0,02).**
- Del 0 al punto derecho (w, f₊(w)) y al izquierdo (−w, f₊(−w)): 0,020002 y
  0,020002 en γ=14; 0,021451 y 0,022140 en γ=999 — dominadas por w porque y ≪ x.
- Ancho de la base a la altura y = a·w²: x₊ − x₋ = 0,040003 (γ=14) a 0,040137
  (γ=1000), o sea 2w más la corrección de la inclinación.
- **De una mitad a la otra de la MISMA parábola**: f₊(w) − f₊(−w) = **−2aLw³**,
  medido −8,159×10⁻⁶ contra −8,161×10⁻⁶ (γ=14) y −1,743×10⁻³ contra
  −1,738×10⁻³ (γ=999). Ésta es la separación que mide algo: es el reloj.
- **Entre las dos curvas (arriba menos abajo)**: g(x) = f₊(x) + f₊(−x). Su
  simetría g(x) − g(−x) = **0,0e+00** en los 15 — **por construcción**
  (undécima aparición de la trampa, declarada antes de correr). No tiene mérito.

## Veredictos pre-registrados

- **V1** (trampa #11) confirmada: 0,0e+00 en todos.
- **V2**: el reloj leído del ajuste cúbico, L medido / L = 0,9997 (γ<50) y
  0,9976 (γ≈1000); peor desvío **2,4×10⁻³** < 1 %. **La asimetría entre las dos
  mitades ES el reloj.**
- **V3**: x₊ + x₋ = L·y/a al 0,01 % (bajo) y 0,5 % (alto); pendiente del eje
  = L/(2a) a la milésima.
- **V4**: β quieto en ½; L determinista; a errante.

## Clasificación

- **Matemática (demostrado):** la paridad de g por construcción; el desarrollo
  y = a x² − a L x³ con L = −Re(ζ''/ζ')(ρ) = 2θ'(γ) para un cero sobre la
  línea (Cauchy–Riemann + ζ = e^(−iθ)Z, F391); de ahí x₊ + x₋ = L y/a y la
  pendiente L/(2a) por perturbación.
- **Medición:** las cuatro tablas, 15 ceros, al 0,01–0,5 %.
- **Nuevo (si algo):** ningún teorema. Un dibujo exacto de qué es del espejo
  (el contacto en ½, la simetría vertical) y qué es de zeta (la inclinación,
  que lee log(γ/2π) desde la forma; la apertura a, que es de cada cero).

Estructura cerrada no es hipótesis demostrada.

**Todavía no.**
