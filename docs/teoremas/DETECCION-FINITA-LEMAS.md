# La detección finita — el acta de los dos lemas

**Documento formal para la auditora · 2026-08-14 · respuesta a los siete
puntos del §12 de la auditoría «LA DETECCIÓN FINITA»** (F304,
`cmd/losdoslemas` verifica cada paso). Las fórmulas van en bloques
monoespaciados. La regla de la auditora manda: no se sella una fórmula
porque funciona en el experimento; se sella cuando sus lemas la derivan
para todo el alcance declarado.

> **La regla del sello** (ley del taller por F302): «Estructura cerrada» ≠
> «Hipótesis demostrada». Nada de este documento es una demostración de RH.

---

## 0. Convención congelada (§12.6 y §9 de la auditoría)

Todas las cantidades de este documento usan UNA sola convención:

    log  = logaritmo NATURAL (ln), en todo el documento
    rho  = cero desafinado, con w = 1 − 1/rho
    R    = |w|
    r    = max(R, 1/R) > 1          [el radio de crecimiento relevante]
    delta = log r                    [LA DEFINICIÓN DE δ — punto §12.1-2]
    u    = 3/delta                   [requiere delta ≤ 1, o sea u ≥ 3]
    techo = ⌈·⌉ redondeo hacia arriba

**Valores oficiales para el par Davenport–Heilbronn** (entrada: el cero
medido a ciegas por este laboratorio, rho = 0.808517 + 85.699348i, con
R = 0.999957995624542 a precisión float64 completa):

    r     = 1/R           = 1.000042006139900
    delta = log r         = 0.0000420052568...
    u     = 3/delta       = 71419.65...
    theta = |arg(w)|      = 0.011668416599736

**Umbral radial oficial** (§12.6, zanjando los 66 escalones): el primer
entero n con r^n > 4 + (4/π)·n·log n, buscado exhaustivamente bajo esta
convención, es

    n₁ = 371842

y da lo mismo con el r completo que con el r redondeado a 1.0000420061:
el valor 371908 que circuló no se reproduce bajo ninguna de las dos — se
propone congelar 371842 como constante oficial.

---

## 1. Lema radial (§12.3), prueba completa

**LEMA R.** Sea delta = log r ∈ (0, 1] y u = 3/delta ≥ 3. Entonces

    n_rad := ⌈ u·log u ⌉

cumple  r^n > 4 + (4/π)·n·log n  para TODO entero n ≥ n_rad.

**Demostración.** Sea n* = u·log u (valor real, sin techo) y
g(n) = e^{n·delta} − 4 − (4/π)·n·log n.

*Paso 1 — g(n_rad) > 0.* Como n_rad ≥ n*, se tiene
n_rad·delta ≥ n*·delta = 3·log u, luego

    e^{n_rad·delta} ≥ e^{3·log u} = u³

Para el otro lado: n_rad ≤ n* + 1 ≤ (1 + 1/3.29)·n* ≤ 1.31·n*
(porque n* = u·log u ≥ 3·log 3 = 3.29), y

    log n_rad ≤ log n* + 1/n* ≤ (log u + log log u) + 0.31 ≤ 2.29·log u

(usando log log u ≤ log u y 0.31 ≤ 0.29·log u pues log u ≥ log 3 = 1.09).
Entonces

    (4/π)·n_rad·log n_rad ≤ (4/π)·1.31·u·log u·2.29·log u ≤ 3.82·u·(log u)²

y basta ver  u³ > 4 + 3.82·u·(log u)², que dividiendo por u se reduce al

    LEMITA AUXILIAR:  u² ≥ 4·(log u)² + 2   para todo u ≥ 3

(4/u ≤ 4/3 < 2 y 3.82 ≤ 4). El lemita: en u = 3 es 9 ≥ 6.83; y
d/du[u² − 4(log u)²] = 2u − 8·log u/u > 0 ⟺ u² > 4·log u, cierto en
u = 3 (9 > 4.39) y creciente. ∎ (lemita)

*Paso 2 — g creciente desde n_rad.* **(Corregido 2026-08-14, auditoría
«la pizza» §6: la versión anterior encadenaba «≤ 4·log u», que FALLA en
u = 3 — 1.28·(2.29·log 3 + 1) = 4.50 > 4·log 3 = 4.39. Crédito a la
auditora; la ruta corregida es la que ella propuso en su §15: comparar
directo contra 3u², sin el paso intermedio innecesariamente fuerte.)*

g'(n) = delta·e^{n·delta} − (4/π)(log n + 1). En n = n_rad:

    delta·e^{n_rad·delta} ≥ delta·u³ = 3u²
    (4/π)(log n_rad + 1) ≤ 1.28·(2.29·log u + 1) = 2.94·log u + 1.28

y la comparación directa:  3u² > 2.94·log u + 1.28  para todo u ≥ 3
(en u = 3: 27 > 4.50, y la diferencia crece: d/du[3u² − 2.94·log u] =
6u − 2.94/u > 0). Luego g'(n_rad) > 0. Además
g''(n) = delta²·e^{n·delta} − (4/π)/n > 0 en n ≥ n_rad (el primer término
es ≥ delta²·u³ = 9u ≥ 27 y el segundo ≤ 1.28/n_rad, minúsculo), así que
g' crece y se mantiene positivo. Con g(n_rad) > 0 y g' > 0 en adelante,
g(n) > 0 para todo n ≥ n_rad. ∎

*Nota de registro:* la lámina de F303 citaba el lemita en su versión para
n* sin techo (u² ≥ 3(log u)² + 2); la versión robusta al techo que usa
esta prueba es la de constante 4. Ambas son ciertas; la oficial es la 4.

---

## 2. Lema de la ventana (§12.4), prueba completa

**LEMA V.** Sea 0 < theta ≤ 2π/3 y K = ⌈2π/theta⌉ + 1. Entonces todo
bloque de K enteros consecutivos {m, m+1, …, m+K−1} contiene un n con
cos(n·theta) ≥ 1/2.

**Demostración.** El conjunto A = {x mod 2π : cos x ≥ 1/2} es el arco
[−π/3, π/3] módulo 2π, de longitud |A| = 2π/3. Sea x = m·theta (posición
sin reducir) y sea L el extremo izquierdo de la primera copia del arco A
en la recta real con L ≥ x (las copias son [2πk − π/3, 2πk + π/3], k ∈ ℤ;
la distancia de x a la próxima copia es < 2π). Sea t ≥ 0 mínimo con
x + t·theta ≥ L. Entonces:

    (i)  si t = 0, x ya está en una copia de A (tomar n = m);
    (ii) si t ≥ 1, por minimalidad x + (t−1)·theta < L, luego
         x + t·theta < L + theta ≤ L + 2π/3,
         o sea x + t·theta ∈ [L, L + 2π/3) ⊆ copia de A.

Un paso de tamaño theta ≤ |A| no puede saltar el arco: ésa es la línea
(ii). Y t está acotado: t·theta ≤ (L − x) + theta < 2π + theta, luego
t ≤ 2π/theta + 1 ≤ K − 1 + 1 = K, y de hecho t ≤ ⌈2π/theta⌉ ≤ K − 1,
así que n = m + t ∈ {m, …, m+K−1}. ∎

**LEMA V-ζ (la hipótesis es automática para zeta).** Todo cero rho con
|Im rho| ≥ 1 cumple |1/rho| ≤ 1, luego w = 1 − 1/rho pertenece al disco
cerrado de centro 1 y radio 1, que está contenido en el semiplano
Re w ≥ 0; por tanto |theta| = |arg w| ≤ π/2 < 2π/3. ∎
(Los ceros no triviales de zeta tienen |Im rho| ≥ 14.13.)

---

## 3. La combinación en el mismo intervalo (§12.5)

Aplicar el LEMA V con m = n_rad: existe n ∈ [n_rad, n_rad + K − 1] con
cos(n·theta) ≥ 1/2. El LEMA R da la desigualdad radial para TODO
n ≥ n_rad — en particular para ESE mismo n. Las dos garantías conviven en
el mismo entero sin pedirse nada: la ventana elige el n, la radial ya
cubre el intervalo entero. ∎

**Realizada para el par DH:** el primer n ∈ S dentro de
[n_rad, n_rad + K − 1] = [798210, 798749] es n = 798474, y allí
4 − r^n + (4/π)·n·log n ≈ −3.7×10¹⁴ < 0.

---

## 4. El teorema formal (§12.7)

**TEOREMA (detección finita cuantitativa).** Sea un espectro formado por
ceros sobre la línea crítica más UN cuarteto desafinado
{rho, conj rho, 1−rho, 1−conj rho} de radio r > 1 y fase theta con
0 < theta ≤ 2π/3 (automática si |Im rho| ≥ 1, LEMA V-ζ). Con la
convención del §0, sea

    N₀(r, theta) = ⌈(3/delta)·log(3/delta)⌉ + ⌈2π/theta⌉ + 1

Entonces existe un entero n ≤ N₀ con cos(n·theta) ≥ 1/2 y
r^n > 4 + (4/π)·n·log n (LEMAS R y V, §3), y para ese n:

    lambda_n ≤ [4 − 2·cos(n·theta)·(R^n + R^{−n})] + resto_n
             ≤ 4 − r^n + (4/π)·n·log n
             < 0

(la primera desigualdad es la fórmula exacta del cuarteto, derivación
§4b paso 1, más la cota sellada del coro, §4b-quater; la segunda usa
cos ≥ 1/2 y R^n + R^{−n} ≥ r^n). Por tanto M_{n,n} = 2·lambda_n < 0 y
M_N no es semidefinida positiva para ningún N ≥ n. ∎

**Para el par DH:** N₀ = 798210 + 539 + 1 = 798750.

**Alcance declarado:** un cuarteto desafinado sobre fondo en la línea —
exactamente el alcance de la cadena sellada en F293–F302. La escalera de
garantías se mantiene separada: n₀ = 85622 (medido) < n₁ = 371842 (cota
pura) < N₀ = 798750 (fórmula cerrada). Reducir las constantes de N₀ es
trabajo declarado futuro (§13.8 de la hoja de ruta).

---

*Verificación numérica de cada paso: `go run ./cmd/losdoslemas` — el
lemita auxiliar en grilla, g, g' y g'' positivos, el lema de la ventana
contra siete fases adversariales (cero violaciones), la hipótesis V-ζ, y
la combinación realizada. Una simulación descubre; una identidad explica;
una demostración cierra todos los pasos — especialmente los infinitos.*
