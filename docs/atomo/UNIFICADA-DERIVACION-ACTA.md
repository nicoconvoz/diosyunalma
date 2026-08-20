# La Derivación Formal de la Curva Madre

**Formalización ordenada tras F387 · 2026-08-20 · sin experimentos nuevos:
toda medición citada proviene de F385–F387, ya registradas.**

---

## 0 · Definiciones (nada más entra después)

- **Trayectoria**: S_n = Σ_{m=1}^{n} m^(−σ)·e^(−i·t·ln m), puntos del plano
  complejo. σ ∈ (0,1) fija (la pintura), t > 0 fija (el dial de altura),
  n ∈ ℕ (el paso).
- **Centro**: C = valor regularizado de la serie completa (para σ ≤ 1, la
  continuación analítica, computada por Euler–Maclaurin). Criterio fijado en
  F385 §2.
- **Radio**: r(n) = |S_n − C| = |T_n|, donde T_n = Σ_{m>n} m^(−σ)e^(−it·ln m)
  es la cola (regularizada).
- **Escala del reloj**: τ = t/2π. **Ángulo local**: θ = t/n.

## 1 · La derivación de la curva madre, término por término

La cola, sacando factor común su primer término:

T_n = n^(−σ)·e^(−it·ln n) · Σ_{j≥1} (1+j/n)^(−σ) · e^(−it·ln(1+j/n))

**Las dos aproximaciones — las condiciones de la orden, con su rango:**

- **(A1) Congelar la amplitud**: (1+j/n)^(−σ) ≈ 1. Error relativo σ·j_ef/n.
- **(A2) Linealizar la fase**: ln(1+j/n) ≈ j/n, así que cada término gira
  θ = t/n respecto del anterior. El término cuadrático despreciado aporta
  fase t·j²/(2n²).
- **(A3) Alcance efectivo**: la serie casi-geométrica está dominada por sus
  primeros j_ef ~ 1/(2|sin(θ/2)|) términos; con eso, los errores de A1 y A2
  son O(σ/(n·sin)) y O(t/(n²·sin²)) — pequeños para n > τ.

Bajo A1+A2 la cola es una serie geométrica exacta:

Σ_j e^(−iθj) = e^(−iθ)/(1−e^(−iθ)) = **1/(e^(iθ)−1)**

y por lo tanto:

**T_n ≈ −n^(−s)/(e^(iθ)−1) ⟹ R(n) = n^(−σ)/(2·|sin(t/2n)|)**

**De dónde sale cada símbolo — la lista de la orden:**

| símbolo | origen |
|---|---|
| **n^(−σ)** | la amplitud del PRIMER término de la cola: la pintura en el paso actual. A1 mide toda la cola en esa unidad. |
| **t** | el único lugar donde entra la altura: la fase dibujada es t·ln m, y localmente el reloj marca θ = t/n radianes por paso. τ = t/2π es donde marca exactamente una vuelta por paso. |
| **sin(t/2n)** | la interferencia de MEDIO ÁNGULO entre pasos consecutivos: \|e^(iθ)−1\| = 2·sin(θ/2) es la CUERDA del círculo unitario entre dos direcciones separadas θ. La cola de un polígono casi-geométrico mide (primer lado)/(cuerda del giro). |
| **el factor 2** | es el 2 de esa cuerda — la distancia entre dos puntos vecinos del círculo unitario. Nada más. |
| **fase exacta (yapa)** | 1/(e^(iθ)−1) = e^(−iθ/2)/(2i·sin(θ/2)): la cola apunta en ángulo recto más medio paso — se usa en la ley de vueltas. |

**Validez**: n > τ (debajo de τ la cola contiene puntos estacionarios — las
resonancias t·ln(1+j/n) = 2πj, que son exactamente la escalera armónica
n_j = τ/j de F377: la misma condición, visto desde el otro lado). Medido
(F387 §1): la curva sigue la trayectoria con mediana 0,993–1,004 en seis
casos y tres bandas — **error global 0,7%**.

## 2 · Observable β

- **Definición**: β = dlnr/dlnn ajustado en [2t, 6t].
- **Derivación**: n ≫ t ⟹ θ → 0 ⟹ sin(θ/2) ≈ θ/2 = t/2n ⟹
  R → n^(−σ)·(n/t) = **n^(1−σ)/t ⟹ β = 1−σ**, obligado. La corrección
  siguiente: ln R = (1−σ)ln n − ln t + θ²/24 + O(θ⁴), y como θ ∝ 1/n el
  ajuste en ventana finita corre la pendiente en −θ²/12 promediado.
- **Predicción**: β = 1−σ − ⟨θ²/12⟩_[2t,6t] ≈ 1−σ − 0,008.
- **Medición** (F385): 0,4927 / 0,6926 / 0,2927 para σ = 0,5 / 0,3 / 0,7.
- **Error**: |β−(1−σ)| = 0,0073 — **igual al sesgo de ventana predicho**. La
  desviación del valor límite no es misterio: es la curvatura θ²/24 de la
  propia curva madre.

## 3 · Observable Δn (ley de vueltas alias/directo)

- **Definición**: vuelta = avance de 2π del ángulo desenrollado de S_n − C
  muestreado en n entero; Δn_k = largo del ciclo k.
- **Derivación**: φ(n) = arg(S_n−C) = −t·ln n − θ/2 + const (fase exacta del
  §1). Incremento por paso: −θ + O(θ/n). El MUESTREO entero con
  desenrollado de arco corto (|salto| < π) lee:
  - θ < π (n > 2τ): lee −θ → **Δn = 2π/θ = n/τ** (directo);
  - π < θ < 2π (τ < n < 2τ): lee 2π−θ = 2π(n−τ)/n → **Δn = n/(n−τ)** (alias).
- **Predicción**: la ley a trozos anterior, con error de cuantización ±1 paso.
- **Medición** (F387 §3): medianas Δn medido/predicho **0,996–1,003**, seis
  casos; tabla ciclo a ciclo 8/7,74 · 6/6,20 · 5/5,41 …
- **Error**: ≤ 0,4% en mediana; ±1 paso por ciclo individual.

## 4 · Observable ε_k — y la relación con k

- **Definición**: ε_k = 1 − r̄_k/r̄_{k−1} sobre ciclos angulares.
- **Derivación**: ε_k ≈ −Δln R por ciclo = [σ − (θ/2)cot(θ/2)]·Δn_k/n_k
  (regla de la cadena sobre la curva madre: dlnR/dlnn = −σ + (θ/2)cot(θ/2)).
  **La relación con k**: k cuenta vueltas estroboscópicas; en la zona alias
  dn/dk = n/(n−τ) ⟹ (n−τ)·dn ≈ τ·dk ⟹ **(n−τ)² ≈ 2τk**. Con θ → 2π⁻,
  (θ/2)cot(θ/2) → −n/(n−τ), y entonces
  ε_k ≈ n/(n−τ)² ≈ τ/(2τk) = **1/(2k)**.
- **Predicción**: potencia α = −1; forma exacta [σ−(θ/2)cot(θ/2)]·Δn/n.
- **Medición** (F386, F387 §4): α ajustado −0,90…−1,25 con la potencia
  ganando a exponencial e hipérbola en 6/6 casos; razón medido/derivado por
  ciclo mediana 1,07–1,53.
- **Error**: el grueso viene de la discretización (Δn = 3–8 pasos ⟹ ±1 paso
  ≈ ±15–30% por ciclo); la tendencia y el exponente quedan; el coeficiente
  fino de 1/(2k) queda como asintótica derivada, verificada en orden y no en
  dígito. Se declara.

## 5 · Observable cintura

- **Definición**: mínimo del radio suavizado (ventana logarítmica ±5%).
- **Derivación**: dlnR/dlnn = 0 ⟺ **(θ/2)·cot(θ/2) = σ**. Con x = θ/2:
  x·cot x decrece de 1 a −∞ en (0, π) ⟹ raíz única x*(σ) ⟹
  **n*/τ = π/x***. (Existencia y unicidad: monotonía de x·cot x.)
- **Predicción** (bisección, sin números a mano): 2,695 (σ=0,5) ·
  2,323 (σ=0,3) · 3,412 (σ=0,7).
- **Medición** (F387 §5): 2,679–2,697 (σ=0,5, cuatro casos) · 2,303 · 3,420.
- **Error**: ≤ 1%. El «3,7τ» de F386 medía otra cosa (el detector de cruce
  sostenido disparando tarde sobre la cola ruidosa) — corregido en F387.

## 6 · Lo que la curva NO es

- No fue ajustada: cero parámetros libres; se derivó de la definición de la
  trayectoria antes de compararla (F387, encabezado del código).
- No es de los ceros: es la anatomía del instrumento — idéntica en cero,
  casi-cero y control. La única firma del cero sigue siendo la posición del
  ojo C respecto del origen.
- No toca a RH: describe la herramienta, no el misterio. Su valor es que
  ahora todo apartamiento futuro de esta curva es, por construcción, señal.

Estructura cerrada no es hipótesis demostrada.

**Todavía no.**
