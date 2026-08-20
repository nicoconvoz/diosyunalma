# Teorema Parcial de la Curva Madre de la Espiral

**Construido por orden de la Auditoría 62 · revisado por las Auditorías
63–64 · 2026-08-20 · F389–F390**
**Rango: TEOREMA PARCIAL, por decreto del capitán (2026-08-20) — exacto
dentro del modelo A1–A3, con el puente al objeto original en CASO B (primer
orden demostrado, cola abierta).**
**Base:** F385–F388. **Naturaleza:** teorema DENTRO DEL MODELO (Nivel I),
con la aproximación a la trayectoria original declarada aparte (Nivel II) y
neutralidad total sobre los ceros (Nivel III). Nada experimental es premisa;
nada de Riemann es premisa; ninguna aproximación se llama exacta.

---

## 1 · Definiciones

Sean σ ∈ (0,1), t > 0, s = σ + it. Para n ∈ ℕ:

- S_n = Σ_{m≤n} m^(−s), con m^(−s) = m^(−σ)·e^(−it·ln m). (La trayectoria.)
- **C := ζ(s) = ζ(σ+it)** — el valor de la única continuación meromorfa de
  Σ_{m≥1} m^(−s) (definida por la serie en Re s > 1) a σ ∈ (0,1), t > 0
  (lejos del polo s = 1). **Representación usada** (Euler–Maclaurin con un
  término de corrección): para todo corte X,
  ζ(s) = S_X + X^(1−s)/(s−1) − X^(−s)/2 + R(X), donde R(X) tiene la forma
  del término B₂ (s/12)·X^(−s−1) más su resto clásico; en magnitud
  |R(X)| ≤ c·|s|·X^(−σ−1) con c de orden 1/12 (clásico, no de este taller).
  **Conexión con los experimentos**: el «ojo» numérico es exactamente esa
  fórmula evaluada en X = N; su error E_C = R(N) es la fila «centro» de la
  tabla del §1b. La definición no se cambia por conveniencia numérica.
- T_n = C − S_n (la cola). r(n) = |T_n|.
- θ = t/n (ángulo local), τ = t/(2π), x = θ/2.

**El modelo local.** Se define la cola modelada

  T_n^M := n^(−s) · Σ_{j≥1} e^(−iθj)   (suma de Abel)

es decir, la cola verdadera con dos sustituciones declaradas:

- **(A1)** (1+j/n)^(−σ) ↦ 1 — congelar la amplitud al primer término.
  Parámetro pequeño: j/n sobre el alcance efectivo.
- **(A2)** ln(1+j/n) ↦ j/n — linealizar la fase. Término omitido dominante:
  t·j²/(2n²) por término.
- **(A3)** la identidad de resumación Σ_{j≥1} z^j = z/(1−z) con z = e^(−iθ),
  θ ∉ 2πℤ, en el sentido de Abel. **A3 es exacta: error propio cero.**

A1 y A2 son HIPÓTESIS DE MODELADO, no aproximaciones acotadas: los
términos dominantes obtenidos formalmente en la expansión de primer orden
son σG/n y θG²/n (con G = 1/(2|sin x|)) — expansión asintótica formal, no
igualdad ni cota — y por F388 esos términos NO son cotas superiores del
error total (violadas como cotas en 14/21 regímenes, hasta ×4,4). La cota
rigurosa de orden completo queda ABIERTA.

## 1b · El puente modelo ↔ objeto original

Se define E_n por la identidad

  **T_n = T_n^M + E_n**

y la descomposición conceptual |E_n| ≤ E_A1(n) + E_A2(n) + E_rest(n),
más, en la práctica numérica, E_C (estimación del centro). Estado real de
cada componente — ninguna celda rellenada artificialmente:

| componente | aproximación | término de la expansión formal de 1er orden | ¿cota rigurosa? |
|---|---|---|---|
| A1 | amplitud congelada | σG/n | **NO — pendiente** |
| A2 | fase linealizada | θG²/n | **NO — pendiente** |
| A3 | resumación de Abel | — | **exacta dentro del modelo** |
| E_mix + E_rest | cruzados/superiores | O(σθG³/n²) formal | **NO — pendiente** |
| discretización | ciclos enteros | ±1 paso por ciclo | acotada trivialmente |
| centro C | resto EM en el corte X | c·|s|·X^(−σ−1), c ~ 1/12 | clásica (§1) |

(Los términos de A1/A2 son los que produce la expansión formal de primer
orden — expansión asintótica formal, no igualdad ni cota. F388 midió que
el error total los excede hasta ×4,4.)

**Regla de lectura para todo el documento**: R_M es la trayectoria radial
DEL MODELO. Toda coincidencia con r(n) significa r(n) = R_M(n)·(1+η(n)),
con η medido ≤ 0,7% en el dominio experimental estudiado. La formulación
correcta es: *la curva madre es la trayectoria radial del modelo
linealizado, y reproduce la trayectoria observada dentro del error medido
en el dominio experimental estudiado* — nunca «R_M = trayectoria
original».

## 1c · El puente reforzado — CASO B (cota parcial demostrada)

La forma cerrada del modelo resuma exactamente las partes (it) de la serie
de Euler–Maclaurin de T_n: con 1/(e^(iθ)−1) = −1/2 − (i/2)cot(θ/2) y la
expansión de Bernoulli del cot, cada término del modelo es el término EM
correspondiente con s reemplazado por it. Por lo tanto E_n = T_n − T_n^M
tiene la representación en serie

  E_n = Σ_k (B_{2k}/(2k)!)·[(s)_{2k−1} − (it)^{2k−1}]·n^(−s−2k+1) + ΔR

(diferencia de restos ΔR incluida), cuyos DOS PRIMEROS TÉRMINOS son
exactos por álgebra pura:

- **D₁ = n^(1−s)·[1/(s−1) − 1/(it)]** — el desajuste de deriva;
  |D₁| = (1−σ)·n^(1−σ)/(t·|s−1|).
- **D₂ = σ·n^(−s−1)/12** — la parte σ del término B₂.

**Verificación numérica del puente** (`go run ./cmd/elpuente`, t = 1000,
tres σ): |E − (D₁+D₂)|/|E| = 0,002–0,026 para θ ≤ 1 (n ≳ t); crece a
0,24–0,66 en θ = 2; y en θ = 3 la cola de la serie domina. **Veredicto del
criterio de cierre (su §14): CASO B — cota parcial.** El primer orden del
puente está demostrado (álgebra) y verificado (≤ 2,6% en su zona θ ≤ 1);
queda abierta la cota uniforme de la cola de la serie (las correcciones σ
a los términos de Bernoulli superiores), que es la pieza dominante para
θ ≳ 2 y degenera hacia θ → 2π. Lo que falta está nombrado con precisión:
acotar Σ_{k≥3} |(s)_{2k−1} − (it)^{2k−1}|·|B_{2k}|/(2k)!·n^(−σ−2k+1)
uniformemente en θ < 2π−δ, más ΔR.

## 2 · Lema de la Curva Madre (exacto dentro del modelo)

**Lema.** Para todo n con θ = t/n ∉ 2πℤ:

  T_n^M = n^(−s)/(e^(iθ) − 1),  y por lo tanto
  **R_M(n) := |T_n^M| = n^(−σ) / (2·|sin(t/2n)|)**,
  con fase arg(T_n^M) = −t·ln n − θ/2 − π/2 (mod 2π).

**Demostración (la cadena completa, cada flecha con su hipótesis).**

  T_n = Σ_{j≥1} (n+j)^(−s)                        [cola original]
   → n^(−s)·Σ_j (1+j/n)^(−σ)·e^(−it·ln(1+j/n))   [factorización: IDENTIDAD]
   → n^(−s)·Σ_j e^(−it·ln(1+j/n))                [flecha A1 — introduce E_A1]
   → n^(−s)·Σ_j e^(−iθj) = T_n^M                 [flecha A2 — introduce E_A2]
   → n^(−s)/(e^(iθ)−1)                           [flecha A3: Abel — IDENTIDAD]
   → R_M = n^(−σ)/(2|sin(θ/2)|)                  [cuerda: |e^(iθ)−1| = 2|sin(θ/2)| — IDENTIDAD]

Tres flechas son identidades; dos (A1, A2) son las hipótesis del modelo y
cargan los errores de la tabla del §1b. La identidad de la cuerda:
e^(iθ)−1 = 2i·sin(θ/2)·e^(iθ/2), de donde módulo y fase. ∎

**Origen de cada factor** (obligatorio por la orden): n^(−σ) = amplitud del
primer término de la cola (A1 mide todo en esa unidad); t/n = fase local
del reloj t·ln m entre pasos consecutivos; sin(t/2n) = módulo de e^(iθ)−1,
la CUERDA del círculo unitario entre dos direcciones separadas θ; el 2 = el
2 de esa cuerda. **La fórmula es exacta dentro del modelo A1–A3, no
necesariamente para la cola original.**

## 3 · Lema de la derivada logarítmica

**Lema.** dln R_M/dln n = −σ + x·cot x, con x = θ/2.

**Demostración.** ln R_M = −σ·ln n − ln 2 − ln|sin x|; como x = t/2n,
dx/dln n = −x, y dln|sin x|/dln n = cot x·(−x). ∎

(De este lema salen los Corolarios I, III y IV; el II usa la fase del Lema 2.)

## 4 · Corolario I — β (exponente de escala)

**Enunciado.** Para θ → 0: dln R_M/dln n = (1−σ) − x²/3 + O(x⁴); en
particular β := lim dln R_M/dln n = **1−σ**, y σ = 1/2 ⟹ β = 1/2.

**Demostración.** x·cot x = 1 − x²/3 − x⁴/45 − ⋯ (serie clásica de cot). ∎

Este 1/2 es un **exponente de escala**, no una razón de contracción por
vuelta (F385 descartó aquella). En ventana finita [2t, 6t] el corrimiento
predicho es −⟨x²/3⟩ ≈ −0,008 (con θ=t/n: x²/3 = θ²/12).

## 5 · Corolario II — ley de vueltas directa/alias

**Hipótesis adicionales**: muestreo en n entero; desenrollado de arco corto
(cada salto de fase se reduce a (−π, π]).

**Enunciado.** El incremento de fase muestreado es −t·ln(1+1/n) =
−θ + O(θ/n). Reducido a (−π, π]:

- **directo** (0 < θ < π, o sea n > 2τ): giro visible θ ⟹ Δn = 2π/θ = n/τ;
- **alias** (π < θ < 2π, o sea τ < n < 2τ): giro visible 2π−θ =
  2π(n−τ)/n ⟹ **Δn = n/(n−τ)**.

**El papel de τ**: en n = τ el reloj da exactamente una vuelta por paso (el
estroboscopio congela: Δn → ∞); en n = 2τ está la frontera de Nyquist del
muestreo (θ = π). Errores separados: de modelo O(θ/n) por paso; de
discretización ±1 paso por ciclo. ∎

## 6 · Corolario III — cintura

**Enunciado.** R_M tiene un único mínimo en n > τ, dado por
**x·cot x = σ**, x = θ/2 ∈ (0, π), y n*/τ = π/x*(σ). La solución existe,
es única, y n*(σ) es estrictamente creciente en σ.

**Demostración.** Por el Lema 3, el punto crítico exige x·cot x = σ.
h(x) = x·cot x cumple h(0⁺) = 1, h(π⁻) → −∞ y h'(x) =
(sin 2x − 2x)/(2 sin²x) < 0 en (0,π) (pues sin u < u): estrictamente
decreciente ⟹ raíz única x*(σ) para todo σ < 1. Monotonía: dx*/dσ =
1/h'(x*) < 0 ⟹ n*/τ = π/x* creciente en σ. ∎

## 7 · Corolario IV — ε_k (condicional, nivel de certeza MENOR)

**Hipótesis adicionales**: (H1) ciclos del Corolario II en zona alias;
(H2) aproximación continua dn/dk = Δn = n/(n−τ); (H3) régimen 1 ≪ k ≪ τ/2.

**Enunciado.** Bajo H1–H3: (n−τ)² = 2τk·(1+o(1)) y

  ε_k = [σ − x·cot x]·Δn/n = **1/(2k)·(1 + O(√(k/τ)) + O(1/k))**.

**Demostración.** De H2: (n−τ)·dn = n·dk ≈ τ·dk (por H3), integrando
(n−τ)² ≈ 2τk. Cerca de θ → 2π⁻: x·cot x = −n/(n−τ)·(1+O((n−τ)/n)),
entonces ε_k ≈ [σ + n/(n−τ)]·(1/(n−τ)) ≈ n/(n−τ)² ≈ 1/(2k). ∎

**Separación obligada**: el EXPONENTE −1 es robusto bajo H1–H3; el
COEFICIENTE 1/2 hereda los errores o(1) y la discretización (±1 paso en
Δn = 3–8 ⟹ ±15–30% por ciclo): **verificado en orden, no en dígito** (F388).
No se eleva al nivel de β ni de la cintura.

## 8 · TEOREMA (condicional, dentro del modelo)

> **Teorema de la Curva Madre de la Espiral (condicional).**
> Sean σ ∈ (0,1), t > 0, y el modelo local definido por A1–A3 sobre la cola
> de S_n. Entonces, para n > τ con θ = t/n ∉ 2πℤ:
>
> **R_M(n) = n^(−σ)/(2·|sin(t/2n)|)**  (exacto en el modelo),
>
> y en sus respectivos regímenes:
> 1. β = 1−σ (θ→0; Corolario I) — en particular σ=1/2 ⟹ β=1/2 como
>    exponente de escala;
> 2. la ley de vueltas Δn = n/τ (directo, n>2τ) y Δn = n/(n−τ) (alias,
>    τ<n<2τ), bajo muestreo entero y desenrollado de arco corto;
> 3. la cintura única en x·cot x = σ, con n*/τ = π/x* creciente en σ;
> 4. ε_k = 1/(2k)·(1+o(1)) bajo las hipótesis adicionales H1–H3, con el
>    coeficiente pendiente de verificación fina.

**Proposición B (aproximación a la trayectoria original — NO teorema).**
Los términos dominantes obtenidos formalmente en la expansión de primer
orden del error T_n vs T_n^M son σG/n (A1) y θG²/n (A2); además, los dos
primeros términos de la representación en serie de E_n son exactos por
álgebra (D₁ y D₂, §1c — CASO B). Empíricamente una envoltura ×5 cubre los
21 regímenes medidos y la frontera del 1% escala como n−τ ≈ √(τ/(2π·0,01))
(predicho 50, medido 43, a t=1000). **Queda abierta la obtención de una
cota rigurosa de orden completo (la cola de la serie del §1c y ΔR).**

## 9 · No-circularidad (las seis preguntas de la orden)

1. ¿La curva se construyó sin usar β? SÍ — resumación de la cola, F387.
2. ¿Sin usar ε? SÍ. 3. ¿Sin usar la cintura? SÍ. 4. ¿Sin datos de ceros?
SÍ — el modelo no contiene ceros; los σ-controles y las alturas de control
dan la misma curva, y el esqueleto ya había sido verificado sin zeta (F379).
5. ¿Las consecuencias se derivan después? SÍ — predicciones pre-registradas
en el encabezado de `cmd/launificada` antes de comparar. 6. ¿Los datos solo
validan? SÍ — tabla de validación abajo; ningún parámetro fue ajustado.

## 10 · Validación experimental (los datos NO son premisa)

| pieza | predicción (sin datos) | medición (F385–F388) | error | fuente del error |
|---|---|---|---|---|
| curva | R_M | medianas 0,993–1,004 | ≤0,7% | A1/A2 + centro EM |
| β | 1−σ−0,008 | 0,6927/0,4927/0,2927 | 0,0073 | ventana (predicha) |
| vueltas | a trozos | 0,996–1,003 | ≤0,4% | discretización ±1 |
| cintura | 2,323/2,695/3,412 | 2,303/2,679–2,697/3,420 | ≤1% | suavizado |
| ε exponente | −1 | −0,90…−1,25 | — | ruido + ciclos |
| ε coeficiente | 1/2 | orden correcto | 15–50% | discretización |

## 11 · Dominio y papel de los ceros

**Dominio (condición del enunciado, no nota):** n > τ. Debajo de τ la cola
contiene los puntos estacionarios t·ln(1+j/n) = 2πj — las resonancias de la
escalera armónica n_j = τ/j (F377) — y la curva madre NO aplica.

**Neutralidad:** el teorema describe la estructura radial del modelo —
común, en lo medido, a ceros, casi-ceros y controles. No usa la Hipótesis
de Riemann, no la implica, y no dice nada de los ceros: la información
específica del cero vive en la posición del centro C respecto del origen,
que es un dato EXTERNO al teorema. **La pregunta «¿cuándo C coincide con
el origen?» queda explícitamente FUERA de este teorema.**

## 12 · Limitaciones

1. Nivel II sin cota rigurosa completa (abierto, declarado).
2. Coeficiente de ε_k sin verificación al dígito.
3. El centro numérico C arrastra la truncación EM de su corte — separado
   como fuente, no eliminado.

## 13 · Prueba de referee (leído como matemático externo hostil)

- **A. ¿Dónde se usa una aproximación?** Solo en las flechas A1/A2 del
  Lema (§2), nombradas en el renglón donde actúan, con su error en §1b.
- **B. ¿Dónde se usa una identidad exacta?** Factorización de la cola,
  resumación de Abel (A3), identidad de la cuerda, derivada logarítmica
  (Lema §3), serie del cot en §1c, y el álgebra de D₁ y D₂.
- **C. ¿Dónde se usa una hipótesis?** A1/A2 (modelado, §1); muestreo
  entero + desenrollado de arco corto (Corolario II); H1–H3 (Corolario IV).
- **D. ¿Dónde se introduce un dato experimental?** Recién en §10
  (validación) y en la verificación numérica del puente (§1c) — nunca
  antes de la predicción correspondiente.
- **E. ¿Dónde se necesita una cota todavía no demostrada?** En toda
  identificación de R_M con la trayectoria original más allá del primer
  orden del puente: la cola de la serie de §1c y ΔR (CASO B, declarado).
- **F. ¿Se puede reconstruir el teorema sin mirar los experimentos?**
  Sí: §§1–8 son autocontenidos; los experimentos aparecen en §10 y §1c
  (verificación, no premisa).
- **G. ¿Queda alguna frase que diga más de lo que la demostración
  permite?** Tras esta pasada, no: se eliminó «la ecuación ES la
  trayectoria» (acta y programa), «exactos como primer orden» pasó a
  «términos de la expansión formal de primer orden», y «anatomía
  universal» bajó a «estructura radial del modelo». En sentido inverso, el
  puente SUBIÓ de rango: de problema totalmente abierto (CASO C) a cota
  parcial demostrada (CASO B), porque el álgebra lo permitió.

## 14 · Clasificación final de cada resultado

| Resultado | Estado |
|---|---|
| Curva madre dentro de A1–A3 | **Lema condicional** (exacto en el modelo) |
| β = 1−σ | Corolario del modelo |
| Ley de vueltas directa/alias | Corolario del modelo (+ hipótesis de muestreo) |
| x·cot x = σ (cintura) | Corolario matemático (existencia/unicidad/monotonía demostradas) |
| ε_k ~ 1/(2k) | Corolario con hipótesis adicionales (coeficiente asintótico/modelizado) |
| T_n − T_n^M (el puente) | Problema de control de error — **CASO B**: D₁+D₂ demostrados y verificados; cola abierta |
| Validación numérica | Evidencia experimental |
| Relación con los ceros | Pregunta separada (fuera del teorema) |

Estructura cerrada no es hipótesis demostrada.

**Todavía no.**
