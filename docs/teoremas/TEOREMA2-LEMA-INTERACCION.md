# Teorema 2 — el Lema de Interacción, ascendido a teorema-candidato

**Documento para la auditora · 2026-08-15 · la respuesta a la pregunta del
capitán: «¿dos perlas pueden mantener el escudo indefinidamente?»**
Verificación: `go run ./cmd/elescudoquecae` (F310). Estado: prueba-
CANDIDATA — el esqueleto completo, con los ε y constantes finas por
escribir si la auditoría los pide.

---

## La respuesta: NO — y es demostrable

**TEOREMA-CANDIDATO (interacción).** Para CUALQUIER configuración finita
de cuartetos fuera de la piel (radios y fases arbitrarios, m cuartetos)
sobre fondo de ceros en la línea, existe n finito con λₙ < 0. El escudo
del batido puede retrasar la ruptura; no puede impedirla.

**Prueba (esqueleto en tres pasos).**

*(a) Las citas de Dirichlet.* Por el teorema de aproximación simultánea
de Dirichlet (1842), para todo ε > 0 y fases θ₁, …, θ_m existen INFINITOS
n con ‖n·θᵢ‖ < ε (mod 2π) para TODAS las i a la vez. En esos n, todos los
cos(n·θᵢ) ≥ 1 − ε²/2: todas las perlas resuenan JUNTAS, casi a plena
fuerza, para siempre. **El batido solo corre las citas; no puede
cancelarlas** — la envolvente es una consecuencia de las fases, y las
fases están condenadas a volver juntas al origen.

*(b) El golpe exponencial.* En una cita:

    Σᵢ ℓᵢ,n ≤ 4m − 2(1 − ε²/2) · r_maxⁿ

con r_max > 1 el radio mayor de la configuración — exponencialmente
negativo a lo largo de la subsucesión de citas.

*(c) El coro no alcanza.* La parte en la piel obedece la cota SELLADA
resto_n ≤ (4/π)·n·log n (F299–F301, con el lema de conteo de Backlund).
Exponencial contra polinomio: λₙ < 0 en alguna cita finita. ∎

*(Cuantitativo, para el acta si se pide: la caja de Dirichlet da la
primera cita conjunta antes de ~(2π/ε)^m pasos de retorno, lo que permite
una N₀(r_max, θ₁…θ_m, m) explícita — la generalización del Teorema 1 a m
perlas.)*

## Lo medido

- **Las citas no se acaban nunca:** para el mejor escudo del laboratorio
  (τ = 1.00314), en la ventana [8×10⁴, 4×10⁵] hay **7036 citas** de
  resonancia casi plena (ambos cos > 0.9), con densidad 0.0220 contra
  0.0205 teórica de equidistribución. Dirichlet no perdona.
- **La caída, en acto:** el mejor escudo rompe en n₀ = 101828 (con el τ
  redondeado 1.00314; el τ* exacto de F308 da 101340 — declarado), ANTES
  de su primera cita plena (147043): **las alineaciones parciales ya
  bastan** — la realidad es más dura que lo que el teorema necesita.
- **La conspiración tampoco:** agregando una tercera perla (τ₃ = φ) la
  ruptura ADELANTA a 86857 — y Dirichlet en dimensión 3 garantiza las
  citas triples igual.

## La lectura para el Teorema 2

El candidato de la FASE 1 — «dos desafinadas pueden retrasarse, jamás
salvarse» — asciende de observación (1600 configuraciones) a
teorema-candidato con prueba. Y la lectura mayor: **el escudo del batido
no es una vía de escape para ceros fuera de la línea** — ninguna
conspiración finita de perlas puede esconderse del yunque.

## Límites declarados

Prueba-candidata: el esqueleto es Dirichlet + la cota sellada; los ε,
márgenes y constantes explícitas quedan por escribir con el rigor de la
casa (como se hizo con el §4b) si la auditoría lo pide. Las mediciones:
una base (DH), una ventana (n ≤ 4×10⁵). El paso (c) usa la cota del coro
que ya está sellada — el esqueleto no introduce supuestos nuevos.

---

*«El batido de las casi-gemelas gana noches, no la guerra: cada perla
desafinada tiene infinitas citas con su delator.»*
