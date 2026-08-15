# Teorema 2 — FASE 1: el patrón medido

**Documento para la auditora · 2026-08-14 · resultados del experimento
central (§8) sobre la matriz (§9) de la hoja «Dos Perlas y Armonía
Relacional»** — la hoja del capitán, con su firma de flash. Verificación:
`go run ./cmd/lasdosperlas` (F307). Estado según el §16 de la hoja: FASE 1
(experimentos → patrón). Ningún teorema se afirma.

---

## 0. Protocolo

    espectro = coro (38 perlas medidas en la piel)
             + cuarteto 1 (par DH real: r₁, θ₁)
             + cuarteto 2 = (r₁^ρ, θ₁·τ)      [parametrización relacional]
    observable: n₀(ρ, τ) = primer n con λₙ < 0
    referencia: perla 1 sola + coro → n₀ = 85622  (F296/F297)
    ventana: n ≤ 3×10⁵ · grilla: ρ ∈ [½, 2] × τ ∈ [0.3, 3], 1600 configs
             + barrido fino de τ (ρ = 1, 541 puntos)

## 1. EL RESULTADO MAYOR: la armonía protectora EXISTE

    configuraciones que RETRASAN la ruptura más allá de la sola: 25 de 1600
    la mejor del barrido fino: τ = 1.010 → n₀ = 89454  (+4.5%)

**El flash del capitán acertó.** Y el mecanismo, medido, es el **BATIDO**:
con τ ≈ 1 las dos fases casi iguales producen la envolvente

    cos(n·θ₁) + cos(n·θ₂) = 2·cos(n·θ̄)·cos(n·Δθ/2)

y cuando la PAUSA de la envolvente cae sobre la zona de ruptura, el par
queda escudado:

    envolvente media |cos(n·Δθ/2)| en la zona (n ≈ 83000–88000):
      gemelas exactas (Δθ = 0) ...... 1.000  (sin escudo — lo MÁS frágil)
      τ = 1.010 (la más protectora) . 0.264  (la pausa cubre la zona)
      cuarta justa (τ = 4/3) ........ 0.633  (escudo parcial)

**Pero nunca salva:** las 1600 configuraciones rompen en n finito — la
protección es retraso, no inmunidad.

    CANDIDATO A LEMA DE INTERACCIÓN (FASE 3):
    dos desafinadas pueden retrasarse; jamás salvarse.

## 2. La dominancia radial (§7 A de la hoja)

    ρ:    1      1.25    1.5     2      3      4
    n₀: 71596  61891  52210  42514  27970  20967
    pendiente log-log de n₀ contra ρ: −0.89 ≈ −1  ⟹  n₀ ~ sola/ρ

El radio mayor pone el reloj; la fase solo ajusta el momento exacto.

## 3. Los nueve casos de la matriz (§9)

    gemelas (ρ=τ=1) ......... 71596 (0.836x)   ← lo más frágil
    espejo (τ = −1) ......... 71596 (idéntico: el coseno es par — caso 9)
    octava ½ / octava 2 ..... 73185 / 71602
    quinta 3/2 .............. 73195
    cuarta justa 4/3 ........ 75902 (0.886x)   ← primera con nombre (DH)
    áurea φ ................. 73190
    dominancia ρ=2 .......... 42514 (0.497x)
    sumisa ρ=½ .............. 80210 (0.937x)
    q constante (ρ=τ=3/2) ... 52773

Control con otra perla base (β = 0.7, γ = 45): la cuarta justa sigue
protegiendo más que las gemelas, pero el ranking fino DEPENDE de la
configuración — es estructura diofántica del cruce, no constante
universal. Declarado.

## 4. El candidato fuerte q = Δθ/ΔR = constante — MUERTO como invariante único

    mismo q = 1:  (ρ=τ=1.5) → n₀ = 52773  ·  (ρ=τ=1.1) → n₀ = 69472
    diferencia: 24% — q solo NO organiza el paisaje

(La propia hoja lo ordenaba en su §13: «NO asumir q = constante».)

## 5. Respuesta a la pregunta fundamental (§10)

Sí existe una relación que separa armonía de fragilidad, con DOS
regímenes — no una curva F(ΔR, Δθ) = 0:

    régimen radial:  el RADIO decide quién gana (n₀ ~ sola/ρ)
    régimen de fase: la armonía vive en Δθ CHICO PERO NO CERO —
                     el batido de las casi-gemelas como escudo

Y la respuesta medida a la pregunta del §15 (la intuición pedida al
capitán): la proporción más armónica no es lineal, ni inversa, ni
logarítmica — es **«casi iguales, apenas desafinadas»**: la pausa del
batido como escudo. Los datos quedan sobre la mesa para su intuición.

## 6. Confesión de la casa

El primer borrador del veredicto decía «cero configuraciones retrasan» —
**los propios datos refutaron la conclusión tecleada antes de registrar**
(25 retrasan; la mejor +4.5%). El pecado recurrente del veredicto
tecleado, cazado esta vez por el experimento mismo. Registrado.

## 7. Límites declarados

Una ventana (n ≤ 3×10⁵), una grilla finita, dos perlas base, DOS cuartetos
solamente. El paso a la FASE 2 (patrón → conjetura) es de la auditora y
del capitán.

---

*Una simulación descubre; una identidad explica; un lema reduce; un
teorema cierra todos los pasos. Y esta vez la simulación descubrió que el
flash del capitán tenía razón.*
