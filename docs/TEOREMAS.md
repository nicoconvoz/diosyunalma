# El libro de los teoremas del taller

**Fundado el 2026-08-14 por orden del capitán: «registra en un nuevo
sector esto porque es nuestro primer teorema — y vienen más».**

Este libro registra los teoremas propios del laboratorio: resultados con
enunciado formal, prueba completa por lemas, alcance declarado y estado
de auditoría externa. Un teorema entra acá cuando sus lemas lo derivan
para todo el alcance declarado — nunca antes (regla de la auditora).

> **La regla del sello** (ley del taller por F302): «Estructura cerrada» ≠
> «Hipótesis demostrada». Ningún teorema de este libro constituye una
> demostración de la Hipótesis de Riemann.

---

## TEOREMA 1 — Detección Finita Cuantitativa de una Perla Desafinada

**Registrado:** 2026-08-14 · F303 (construcción) · F304 (acta de los dos
lemas) · F305 (corrección del paso 2) · certificado por la auditoría
externa «PRIMER_TEOREMA_DEL_YUNQUE» (§7: todo verde dentro del alcance).

**Enunciado.** Sea un espectro formado por ceros sobre la línea crítica
más UN cuarteto desafinado {rho, conj rho, 1−rho, 1−conj rho}, con

    w = 1 − 1/rho = R·e^{i·theta}
    r = max(R, 1/R) > 1
    0 < theta ≤ 2π/3        [automático para zeta: |Im rho| ≥ 1 ⟹ |theta| ≤ π/2]
    delta = log r            [log natural; convención congelada del acta]

Sea

    N₀(r, theta) = ⌈(3/delta)·log(3/delta)⌉ + ⌈2π/theta⌉ + 1

Entonces existe un entero n ≤ N₀ tal que, simultáneamente,

    cos(n·theta) ≥ 1/2      y      r^n > 4 + (4/π)·n·log n

y para ese n:

    lambda_n < 0,   M_{n,n} = 2·lambda_n < 0,
    M_N no es semidefinida positiva para ningún N ≥ n.   ∎

**Prueba.** Por composición de los dos lemas del acta
(`docs/DETECCION-FINITA-LEMAS.md`):

- *Lema radial* (F304, paso 2 corregido en F305 con la receta de la
  auditora): n_rad = ⌈u·log u⌉ con u = 3/delta garantiza la desigualdad
  radial para TODO n ≥ n_rad.
- *Lema de la ventana* (F304): todo bloque de K = ⌈2π/theta⌉ + 1 enteros
  consecutivos contiene un n con cos(n·theta) ≥ 1/2.
- *Combinación*: la ventana aplicada en m = n_rad elige el n; la radial ya
  cubre el intervalo. El paso final usa la fórmula exacta del cuarteto
  (F297) y la cota sellada del coro, resto_n ≤ (4/π)·n·log n
  (F299–F301, con el lema de conteo de Backlund documentado).

**Caso testigo (par Davenport–Heilbronn,** rho = 0.808517 + 85.699348i**):**

    n₀ medido (coro de 38 + par)  =  85622
    n₁ umbral radial puro         =  371842
    n_rad = 798210 · K = 540 · primer n ∈ S en la ventana = 798474
    N₀ = 798750

La escalera de garantías se conserva separada: experimento ≠ cota ≠
fórmula cerrada.

**Alcance y límites** (declarados; §6 del certificado): un solo cuarteto
desafinado sobre fondo en la línea. No se extrapola en silencio a
configuraciones con múltiples cuartetos. No resuelve la positividad
global desde la aritmética de los primos. No demuestra RH.

**Nota metodológica** (§8 del certificado): «Primer teorema del Yunque»
es un nombre de trabajo del laboratorio; el reconocimiento externo
requeriría revisión independiente completa, comparación con la
literatura y publicación.

**Reproducir:** `go run ./cmd/elprimerteorema` (la cadena numérica
completa en una corrida) · `go run ./cmd/losdoslemas` (los lemas paso a
paso) · `go run ./cmd/ladeteccionfinita` (la construcción original).

---

*Espacio reservado para el Teorema 2 — porque vienen más.*
