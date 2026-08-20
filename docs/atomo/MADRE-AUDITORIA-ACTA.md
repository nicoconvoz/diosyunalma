# Acta de la Auditoría de la Curva Madre — respuesta a la Auditoría 61

**F388 · 2026-08-20 · formato A–L exigido por su §16 · se atacó la curva y
la curva sangró donde tenía que sangrar: dos fallos registrados, ninguno
reparado por redacción.**

Se reproduce con `go run ./cmd/laauditoriamadre`.

---

## A · Definiciones (notación unificada — su §14 corregido)

s = σ + it (una sola vez, acá). S_n = Σ_{m≤n} m^(−s)·(con m^(−s) =
m^(−σ)e^(−it·ln m)). C = continuación analítica (EM). T_n = C − S_n.
r(n) = |T_n|. τ = t/2π. θ = t/n. G = 1/(2·|sin(θ/2)|).

## B · Lema A1 (congelar amplitud)

- **Hipótesis**: (1+j/n)^(−σ) ≈ 1 sobre el alcance efectivo.
- **Derivación**: el primer orden de la corrección es −(σ/n)·Σ j·e^(−iθj);
  con |Σ j·e^(−iθj)| = 1/(4sin²(θ/2)) exacto, el tamaño relativo es
  **σG/n** — «~» significa: término dominante EXACTO de la serie de
  correcciones, no cota superior.
- **Cota**: **FALLO REGISTRADO** — como cota superior estricta, σG/n es
  violada en 14 de 21 regímenes (hasta ×4,4): los órdenes siguientes y el
  error EM del propio centro C se le suman. Envoltura práctica que cubre
  los 21 regímenes medidos: **5×(σG/n + θG²/n)**. La cota rigurosa de orden
  completo queda **PENDIENTE** y así se declara.

## C · Lema A2 (linealizar fase)

- **Hipótesis**: ln(1+j/n) ≈ j/n; resto cuadrático t·j²/(2n²) por término.
- **Derivación**: primer orden −(it/2n²)·Σ j²e^(−iθj), tamaño relativo
  **θG²/n** (con |Σ j²x^j| = |x(1+x)|/|1−x|³ ≤ 2·(2sin(θ/2))^(−3)).
- **Cota**: mismo estado que B — término dominante correcto, cota estricta
  fallida, envoltura ×5 vigente.

## D · Lema A3 (alcance efectivo)

- «Dominada» tiene definición rigurosa: la identidad geométrica es la suma
  de Abel del modelo linealizado y G es su peso: j_ef = G no es heurística,
  es el módulo de la resumación. **A3 no es una aproximación: error propio
  cero.** CERRADO.

## E · Proposición (curva madre)

Bajo A1+A2, T_n = −n^(−s)/(e^(iθ)−1) EXACTO en el modelo linealizado ⟹
R(n) = n^(−σ)/(2|sin(θ/2)|). Contra la trayectoria: error por régimen
0,20–0,78% en la zona lisa (θ ≤ 3), 0,55–2,7% en la zona alias (tabla §1
del programa, 21 regímenes × 3 σ). **Dominio de validez medido y corregido:
el 1% se alcanza en n−τ = 43 pasos (σ=0,5, t=1000), y la fórmula correcta
de la frontera es la de A2: n−τ ≈ √(τ/(2π·tol)) — predice 50. Mi frontera
anterior (σ/(2π·tol) ≈ 8 pasos, basada en A1) estaba MAL: cerca de τ manda
A2, y la frontera escala como √τ, no como constante. Fallo registrado.**

## F · Corolario β — CERRADO

β = 1−σ obligado en el límite; el sesgo de ventana −⟨θ²/12⟩ ≈ −0,008
predicho contra −0,0073 medido. Asintótico y ventana finita, separados.

## G · Corolario Δn — CERRADO

Rama del argumento auditada: el salto muestreado −t·ln(1+1/n) cae en
(−π,π] como −θ (directo, θ<π) o 2π−θ (alias, π<θ<2π); transición en n=2τ.
Error de modelo O(θ/n) por paso; discretización ±1 paso por ciclo,
separados. Medianas 0,996–1,003.

## H · Corolario ε_k — PARCIALMENTE CERRADO

Cadena verificada: ε = [σ−(θ/2)cot(θ/2)]·Δn/n; dn/dk = n/(n−τ) ⟹
(n−τ)² ≈ 2τk (válido para 1 ≪ k ≪ τ/2) ⟹ ε_k ≈ 1/(2k). Separación exigida:
**exponente** α=−1 derivado y compatible (−0,90…−1,25 medido); **coeficiente**
1/2 solo verificado en orden (discretización ±15–30% por ciclo domina);
**ruido**: pollera de resonancia; **ajuste finito**: primer ciclo excluido.
El coeficiente queda como asintótica derivada sin verificación al dígito.

## I · Corolario cintura — CERRADO

h(x) = x·cot x estrictamente decreciente en (0,π) (verificado; h' < 0 por
sin 2x < 2x), h(0⁺)=1, h(π⁻)→−∞ ⟹ raíz única para todo σ<1. n*/τ = π/x*.
Sensibilidad: d(n*/τ)/dσ = +1,4 / +2,4 / +5,4 en σ = 0,3/0,5/0,7 — la
cintura es más sensible a σ cuanto más alta, y aún así lo medido quedó al
1%. El «3,7τ» sigue muerto y enterrado como artefacto del detector.

## J · Comparación predicción/medición (resumen)

| pieza | predicho | medido | estado |
|---|---|---|---|
| curva (zona lisa) | — | 0,20–0,78% | dentro de la envoltura ×5 |
| frontera 1% | 50 pasos (A2) | 43 pasos | corregida y compatible |
| β (ventana) | 1−σ−0,008 | desvío 0,0073 | cerrado |
| Δn | ley a trozos | 0,4% | cerrado |
| ε exponente | −1 | −0,90…−1,25 | cerrado |
| ε coeficiente | 1/2 asintótico | orden, no dígito | parcial |
| cintura | 2,323/2,695/3,412 | 2,303/2,679–2,697/3,420 | cerrado (1%) |

## K · Limitaciones

1. Sin cota superior rigurosa de orden completo para A1/A2 (solo término
   dominante exacto + envoltura empírica ×5).
2. El coeficiente de ε_k sin verificación fina.
3. El error del centro C (truncación EM en el corte X) contamina la medición
   del error de la curva en la zona lejana — contribución identificada, no
   separada término a término.
4. Todo vale para n > τ; debajo de τ mandan las resonancias (la escalera de
   F377) y la curva NO aplica.

## L · Veredicto: **PARCIALMENTE CERRADO**

Cerrado: A3, la proposición en su dominio medido, β, Δn, cintura, y la
separación anatomía-universal / firma-del-cero (la firma sigue siendo solo
la posición de C). Parcial: las cotas estrictas de A1/A2 (fallo registrado,
envoltura declarada) y el coeficiente de ε_k. Refutado: mi primera fórmula
de frontera. **Su §12**: sin(θ/2) y cot(θ/2) son EL MISMO objeto —
dln|sin(θ/2)|/dθ = (1/2)cot(θ/2), verificado 0,214821 = 0,214821 — mientras
que el 1/2 de ε_k ≈ 1/(2k) es el de ∫x·dx = x²/2 en (n−τ)² ≈ 2τk: **dos
familias distintas, y se mantienen separadas**, como manda.

Sin cierre total no hay congelamiento total: quedan pendientes la cota
rigurosa y el dígito del coeficiente. El resto de la estructura queda
congelada según su §18, lista para auditoría independiente.

Estructura cerrada no es hipótesis demostrada.

**Todavía no.**
