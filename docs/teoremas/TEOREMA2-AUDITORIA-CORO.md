# Auditoría del engranaje F299–F301 — la factura del coro

**Para la auditora · 2026-08-15 · respuesta al pedido «H4 y el coro»
(F318): auditar L5 desde sus primeras líneas, bajo H4, con intento
activo de ruptura.** Solo el coro — R0-R10, Dirichlet, agenda y golpe no
se tocan (no apareció dependencia que lo exigiera).

---

## 1. El coro, reconstruido desde cero

**Definición.** El fondo es un conjunto de ceros SOBRE la línea
(β = 1/2), **cerrado bajo conjugación** (los ceros vienen en pares
{½+iγ, ½−iγ}, γ > 0) — ver §2 sobre esta cláusula. Su contribución:

    coro_n = Σ_pares a_n(γ),    a_n(γ) = 2 − 2·Re(wⁿ) = 4·sin²(n·φ/2)

con w = e^{iφ} (|w| = 1 exacto en la línea, F214) y el ángulo EXACTO

    φ(γ) = arg((ρ−1)/ρ) = 2·arctan(1/(2γ))       [derivado en F299]

**Cota por par:** a_n = 4·sin²(nφ/2) ≤ min(4, (nφ)²) — de sin² ≤ 1 y
|sin x| ≤ |x|. **Constante C = 1:** arctan(x) < x ⟹ φ < 1/γ ⟹
(nφ)² < n²/γ². Ambas verificadas antes (F298/F299: 0 violaciones;
φ·γ → 1 por abajo, o sea C = 1 es AJUSTADA, no mejorable).

**Dónde entra H4 (dos veces):**

    (a) conteo:  Σ_{γ≤n} 4 ≤ 4·N_f(n) ≤ 4·(n/2π)·log n = (2/π)·n·log n
    (b) cola:    Σ_{γ>n} n²/γ² = n²·T(n), y por sumación parcial
        T(x) = 2∫ₓ^∞ N_f(t)/t³ dt − N_f(x)/x²        [borde ≤ 0: se descarta]
             ≤ (1/π)∫ₓ^∞ log t/t² dt = (log x + 1)/(π·x)
        [la integración por partes exige N_f(t)/t² → 0: lo da H4, pues
        N_f/t² ≤ log t/(2πt) → 0; y la primitiva −(log t+1)/t es exacta]

**El 4/π, exacto y no «del orden de»:**

    coro_n ≤ (2/π)·n·log n + n²·(log n + 1)/(π·n)
           = (3/π)·n·log n + n/π
           ≤ (4/π)·n·log n      ⟺  n/π ≤ (1/π)·n·log n  ⟺  log n ≥ 1

**y de ahí el umbral n ≥ 3** (el menor entero con log n ≥ 1). Cada
constante tiene partida de nacimiento: 2/π del conteo, 1/π de la cola,
1/π del término n/π absorbido, total 4/π.

## 2. Cuantificadores, dominio y convergencia (las cinco preguntas)

    ① ¿∀n ≥ 3 o solo los n de la agenda? — PARA TODO n ≥ 3: cada paso
      de §1 es universal en n; la agenda no interviene en L5.
    ② ¿H4 basta por sí sola? — Sí para la COTA. Pero hay una cláusula
      DEFINICIONAL que debe declararse (no es hipótesis matemática
      extra): el fondo cerrado bajo conjugación — sin ella λₙ ni
      siquiera es real. Queda escrita en la definición de configuración.
      La finitud local está implicada por H4 (N_f(T) < ∞ ∀T).
    ③ ¿Hipótesis adicional? — Solo la cláusula de ②, definicional.
    ④ ¿Tipo de convergencia? — ABSOLUTA: a_n(γ) ≥ 0 (¡suma de términos
      no negativos!). No hay convergencia condicional, ni agrupación,
      ni simetrización en L5. [La convergencia apareada solo aparecía
      en documentos anteriores para λ con espectros infinitos generales;
      aquí los cuartetos son finitos y el coro es de términos ≥ 0.]
    ⑤ ¿Dependencia del orden de sumación? — Ninguna: términos no
      negativos ⟹ el valor es el supremo de las sumas parciales,
      independiente del orden (Tonelli para series).

## 3. Intento activo de romper L5

**Ataque: el fondo adversarial de densidad máxima.** Se construyó el
peor fondo que H4 permite — pares colocados exactamente donde
N_f(T) = (T/2π)·log T incrementa (densidad al tope) y cada par cobrando
su cota min(4, n²/γ²) entera (fase peor imaginable):

    n = 10:     coro_max = 14.6    contra cota 29.3     ✗ no rompe
    n = 100:    coro_max = 309.9   contra cota 586.3    ✗ no rompe
    n = 1000:   coro_max = 4513.3  contra cota 8795.2   ✗ no rompe
    n = 10000:  coro_max = 52819   contra cota 117270   ✗ no rompe

**Por qué no existe contraejemplo:** la cota es un TEOREMA derivado de
H4 por desigualdades puntuales válidas par a par — el adversario solo
puede acercarse a (3/π)·n·log n + n/π, que queda estrictamente debajo de
(4/π)·n·log n para n ≥ 3. El margen (~45%) es el colchón del paso
log n ≥ 1. Un fondo que rompa la cota tendría que violar H4.

## 4. Interno vs externo

    INTERNO (F299–F301, matemática pura): la implicación
    «H4 ⟹ coro_n ≤ (4/π)·n·log n para n ≥ 3».
    EXTERNO (solo para la Parte C, aplicación a ζ): que el fondo de ζ
    satisface H4 — eso es B2 (Riemann–von Mangoldt + error explícito de
    Backlund 1918), ya etiquetado en F318.

## 5. Verificación de constantes (pedido §5, una por una)

    4/π ......... exacta: 2/π (conteo) + 1/π (cola) + 1/π (absorción
                  de n/π con log n ≥ 1) — §1
    C = 1 ....... exacta y ajustada: φ = 2·arctan(1/2γ) < 1/γ, con
                  φ·γ → 1 por abajo (no mejorable)
    n ≥ 3 ....... el menor entero con log n ≥ 1 — origen exacto en §1
    la cola ..... (log x + 1)/(π·x), primitiva exacta −(log t+1)/t
    el borde .... −N_f(x)/x² ≤ 0, descartable en cota superior (F300)
    transición .. en L6 el coro se evalúa en la cita n ≥ n_rad,m ≥ 11 ≥ 3:
                  el umbral se cumple siempre; λₙ ≤ (4/π)n·log n + 2m+2
                  − r_maxⁿ < 0 por L4

## 6. Veredicto sobre L5

    🟢 L5 CORRECTO — F299–F301 demuestra exactamente
    «coro_n ≤ (4/π)·n·log n para todo n ≥ 3» bajo H4, con una
    explicitación DEFINICIONAL (no hipótesis nueva): el fondo cerrado
    bajo conjugación, ahora escrita en la definición de configuración.
    Convergencia absoluta, sin dependencia de orden. Sin contraejemplo
    posible bajo H4 (§3). Constantes exactas, no «del orden de».

---

*La factura del coro, abierta renglón por renglón: cada peso tiene su
recibo, y el adversario con la billetera llena de H4 no llegó ni al 60%
de la cota.* 🌭📜🔍
