# Acta de la Unificada — respuesta a la Orden Unificada F385–F386 (Auditoría 60)

**Respuesta al protocolo maestro · 2026-08-20 · F387**

Su pregunta unificadora tiene respuesta, y es la que usted exigía: **una sola
ley de escala genera los cuatro observables.** Se derivó primero, se verificó
después, y dos interpretaciones mías de F386 murieron en el camino — se
registran como fallos, como manda su punto 9.

Se reproduce con `go run ./cmd/launificada`. Lámina:
`galeria/laminas/10-el-telar/la-unificada.svg`.

---

## La frase de su §25, que ahora se puede escribir

> **Partiendo de** la trayectoria S_n = Σ_{m≤n} m^(−σ)·e^(−it·ln m) y su ojo
> C = ζ, **se obtiene** por resumación geométrica de la cola (los pasos giran
> casi-igual cerca de n: Σ_j e^(−itj/n) = 1/(e^(it/n)−1)):
>
> **S_n − ζ ≈ −n^(−s)/(e^(iθ)−1), θ = t/n ⟹ R(n) = n^(−σ)/(2|sin(θ/2)|)**
>
> **De ella se siguen simultáneamente** β = 1−σ (límite θ→0), la merma
> ε = [σ−(θ/2)cot(θ/2)]·Δn/n con asintótica 1/(2k), la ley de vueltas a
> trozos Δn = n/(n−τ) (alias, n<2τ) y n/τ (directo, n>2τ), y la cintura
> (θ/2)cot(θ/2) = σ, **bajo las condiciones** n > τ (sin resonancias en la
> cola) y t ≫ 1. **Los experimentos independientes son compatibles dentro
> de** 0,7% (curva), 0,4% (vueltas), 1% (cintura) y ~15–50% (merma por
> ciclo, dominada por la discretización de ciclos).

## A · OBSERVADO

- La curva madre contra el radio medido punto a punto, seis casos, tres
  bandas (1,2τ..8τ]: medianas **0,993–1,004**. [Corregido por la Auditoría
  63: la curva madre es la trayectoria radial del MODELO linealizado, y
  reproduce la observada dentro del error medido en el dominio estudiado.]
- Vueltas: Δn medido contra la ley a trozos: medianas **0,996–1,003**.
- Cintura (mínimo real del radio suavizado): 2,679 / 2,697 / 2,693 (σ=0,5),
  2,303 (σ=0,3), 3,420 (σ=0,7), en unidades de τ.

## B · AJUSTADO

- β sobre la propia curva madre en [2t, 6t]: 0,6927 / 0,4927 / 0,2927 —
  **idéntico a lo medido en F385 hasta el sesgo de ventana** (−0,007): la
  ecuación reproduce hasta el defecto del ajuste.
- ε medido/derivado por ciclo: medianas 1,07–1,53 (la discretización de los
  ciclos domina el error; el orden y la tendencia 1/(2k) están).

## C · DERIVADO (sin usar el resultado para construir la fórmula)

1. **β(σ) = 1−σ**: θ→0 ⟹ R → n^(1−σ)/t. Su pregunta del §22 queda
   respondida: la construcción OBLIGA a β = 1−σ, y por lo tanto σ=1/2 obliga
   β=1/2. Derivado, no supuesto.
2. **Cintura**: dlnR/dlnn = −σ + (θ/2)cot(θ/2) = 0 ⟹ (θ/2)cot(θ/2) = σ.
   Resuelta por bisección: n*/τ = 2,695 (σ=0,5), 2,323 (0,3), 3,412 (0,7).
   Contra lo medido: 2,679 / 2,303 / 3,420. **A ~1%.**
3. **Vueltas estroboscópicas**: el paso entero muestrea la fase; para
   θ∈(π,2π) el giro visible es el alias 2π−θ = 2π(n−τ)/n ⟹ Δn = n/(n−τ);
   para θ<π, Δn = 2π/θ = n/τ.
4. **Merma**: ε = −dlnR/dlnn·Δlnn; cerca de τ, (n−τ)² ≈ 2τk ⟹ **ε_k ≈ 1/(2k)**
   — la potencia α = −1 que F386 midió como −0,90…−1,25, derivada.

## D · INTERPRETADO — y los dos fallos de F386, registrados

- **FALLO 1**: el «p ≈ 3/2 del término de resto» era un espejismo de rango —
  p_k no es constante (decrece como n/(n−τ): 6,4 → 1,9 medido) y su mediana
  1,5–1,7 era una foto del régimen intermedio. Caso C/D de su §4: la
  interpretación del resto era incorrecta.
- **FALLO 2**: «σ se cancela y la cintura es universal» era FALSO. La cintura
  depende de σ — y la ecuación la clava caso por caso. El «3,7τ universal»
  de F386 era el detector de cruce sostenido disparando tarde sobre la cola
  ruidosa (caso C de su §11: dependía del criterio, no de la espiral).
- Lo que sí queda: toda esta arquitectura es del INSTRUMENTO (idéntica en
  cero, casi-cero y control — su regla §15 respetada). Los ceros siguen
  viviendo solo en la posición del ojo.

## Tabla maestra actualizada (su §19)

| Observable | Resultado | Estado |
|---|---|---|
| curva R(n) | n^(−σ)/(2·sin(θ/2)·) | **derivada + verificada 0,7%** |
| razón por vuelta q_k | → 1 | derivado (límite de la curva) |
| merma ε_k | ≈ 1/(2k) | derivada; verificación por ciclo gruesa (~15–50%) |
| exponente α | −1 | derivado |
| exponente local p | NO constante; n/(n−τ) | corregido (fallo 1) |
| enrollado β | = 1−σ | **derivado** |
| cintura n*/τ | π/x*, con x*cot(x*)=σ | **derivada + verificada 1%** |
| coeficiente «3,7» | artefacto del detector | resuelto (fallo 2) |
| posición del ojo | única firma del cero | sin cambio |

Su §23 intacto: nada de esto toca a RH. La ecuación describe el instrumento
con el que miramos — y ahora que el instrumento está entendido hasta el
hueso, lo que NO explique esta curva será, por fin, señal.

Estructura cerrada no es hipótesis demostrada.

**Todavía no.**
