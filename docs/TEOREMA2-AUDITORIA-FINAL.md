# Auditoría final completa del Lema de Interacción — el informe A-H

**Para la auditora · 2026-08-15 · respuesta al encargo «la cuenta del
pancho» (F317): auditoría de extremo a extremo, reconstrucción desde
cero, caza de cuantificadores e INTENTOS ACTIVOS DE ROMPER EL TEOREMA.**
Regla seguida: no dar por cerrado ningún paso por haber sobrevivido antes;
no defender el teorema — descubrir qué matemática tenemos.

**RESULTADO PRINCIPAL DE LA AUDITORÍA: se encontró UNA HIPÓTESIS OCULTA
(tarea 5/8/10), se declaró, y el acta quedó parcheada. Con el parche, el
veredicto interno es 🟡 CASI SELLO — el sello es de la auditora, nunca
del autor.**

---

## A. Enunciado exacto que queda probado (post-parche)

**TEOREMA DE INTERACCIÓN (Parte A).** Sea una configuración formada por:
un fondo numerable de ceros sobre la línea crítica, más m cuartetos
{ρᵢ, ρ̄ᵢ, 1−ρᵢ, 1−ρ̄ᵢ}. Bajo las hipótesis H0–H4 de la sección B, existe
un entero n con

    n ≤ N₀(r_max, m) = n_rad,m + ⌈2π·n_rad,m⌉^m      [n_rad,m = ⌈u_m·log u_m⌉, u_m = 3(m+1)/δ]

tal que λₙ < 0 y por tanto M_{n,n} = 2λₙ < 0 y M_N no es semidefinida
positiva para ningún N ≥ n. **Cuantificador: EXISTE un n ≤ N₀ (no «para
todo n»); la negatividad de M vale para TODO N ≥ n.**

## B. Lista de hipótesis (completa, sin interpretación)

    H0. rᵢ = max(|wᵢ|, 1/|wᵢ|) > 1 para todo i — los m cuartetos están
        ESTRICTAMENTE fuera de la línea. [Antes implícita en «fuera de
        la piel»; ahora explícita. Sin ella δ = 0 y u_m no está definido.]
    H1. m entero, m ≥ 1, finito.
    H2. |Im ρᵢ| ≥ 1 para todo cuarteto.
    H3. δ = log r_max ≤ 1  (r_max ≤ e). [Automática si H2 y los cuartetos
        vienen de la franja: lema δ-ζ da r ≤ √2.]
    H4. ⚠️ LA HIPÓTESIS ANTES OCULTA — densidad del fondo: el conteo de
        pares del fondo en la línea cumple N_fondo(T) ≤ (T/2π)·log T
        para T ≥ 2 (y en particular el fondo es localmente finito, con
        la convergencia apareada estándar de λₙ).

## C. Prueba reconstruida como cadena de lemas (con cuantificadores)

    L1 (GOLPE; usa H0, H1). Para todo n con ‖n·θᵢ‖ < 1 ∀i («cita»):
        Σᵢ ℓᵢ,n ≤ 2m + 2 − r_maxⁿ.   [L1-L7 del acta §1; por cita, no ∀n]
    L2 (DIRICHLET; sin hipótesis). ∀Q ≥ 2 entero ∃n ∈ [1, Q^m] con
        ‖n·θᵢ‖ ≤ 2π/Q ∀i.   [palomar; existencial]
    L3 (AGENDA; usa L2). ∀T ∈ ℤ, T ≥ 1, ∃n ∈ [T, T + ⌈2πT⌉^m] que es
        cita con ε = 1.   [existencial, posterior a T]
    L4 (RADIAL-m; usa H0, H1, H3). ∀n ≥ n_rad,m:
        r_maxⁿ > 2m + 2 + (4/π)·n·log n.   [universal desde n_rad,m;
        R0–R10 del acta, con R6 y R8 autocontenidos]
    L5 (CORO; usa H4). ∀n ≥ 3: coro_n ≤ (4/π)·n·log n.   [universal;
        la cadena F299–F301: por-par ≤ min(4, n²/γ²) con C = 1
        (φ = 2·arctan(1/2γ), exacto en la línea), conteo H4, cola por
        sumación parcial con borde −N(x)/x² ≤ 0]
    L6 (ENSAMBLE). Tomar T = n_rad,m en L3: existe cita n con
        n_rad,m ≤ n ≤ N₀. En ese n: λₙ = coro_n + Σℓᵢ,n
        ≤ (4/π)n·log n + 2m + 2 − r_maxⁿ   [L5 y L1; n ≥ n_rad,m ≥ 11 ≥ 3]
        < 0                                 [L4 en ese mismo n]  ∎

    Cada n construido satisface SIMULTÁNEAMENTE cita (L3) y umbral
    radial (L4 es universal desde n_rad,m, y la cita es ≥ n_rad,m):
    no hay conflicto de cuantificadores.

## D. Auditoría independiente por lema

    L1  🟢  cadena L1-L7 releída desde ℓᵢ,n literal; signos y sentidos
            de desigualdad correctos; cos > 0 garantizado por ε = 1;
            verificada además en 24079 citas dobles y 7675 triples.
    L2  🟢  palomar estándar; Q^m + 1 puntos, Q^m cajas; n = j − j' ≥ 1.
    L3  🟢  las cuatro desigualdades con T entero (F312/F313); la
            subaditividad de ‖·‖ es la del cociente ℝ/2πℤ, sin
            restricción; Q = ⌈2πT⌉ ≥ ⌈2π⌉ = 7 ≥ 2 para T ≥ 1.
    L4  🟢  R0–R10 releído: techos correctos (R2 usa n* ≥ 10.75),
            R6/R8 autocontenidos (F316/F317), R9 uniforme en n, R10
            propaga. δ > 0 garantizado por H0; log definidos (u ≥ 6).
    L5  🟠→🟢  AQUÍ ESTABA EL HALLAZGO: la cota invoca el conteo tipo
            Backlund, que NO es interno a la configuración abstracta —
            es una propiedad del fondo. Declarada como H4, el lema es
            correcto y su cadena (F299–F301) se releyó entera: C = 1,
            sumación parcial con el signo del borde, primitiva exacta.
    L6  🟢  el ensamble usa L4 universal + L3 existencial ≥ T: legal.

## E. Dependencias externas (para la aplicación a ζ — Parte C)

    B1. |Im ρ| > 1 para los ceros no triviales de ζ: conteo riguroso
        N(14) = 0 (Backlund 1914; van de Lune 1986; Platt–Trudgian
        2021). EXTERNO.
    B2. ⚠️ NUEVA, antes camuflada dentro de «la cota sellada»: la
        densidad H4 para el fondo de ζ — N(T) ≤ (T/2π)·log T — viene de
        Riemann–von Mangoldt con el error explícito de Backlund (1918).
        EXTERNO, del mismo rango que B1. [El conteo N(T) mayora en
        particular al subconjunto de ceros en la línea, así que H4 se
        satisface en el escenario ¬RH que la Parte C contempla.]
    Las verificaciones numéricas propias (motor de perlas, grillas) son
    CORROBORACIÓN instrumental, no fuentes de B1/B2.

## F. Intentos de falsación y resultados

    F-1. Bordes: δ = 1, m = 1 (u = 6 exacto): todas las desigualdades
         sobreviven con holgura estricta (R6: 36 > 10.63; R5: 4 ≤ 6). ✗
    F-2. m grande (δ ≤ 1): R5 pide δ ≤ 3/2 — cubierto; n_rad,m crece
         pero es finito. ✗
    F-3. θᵢ = 0 (cuarteto real): excluido por H2 (w real ⟹ ρ real ⟹
         Im ρ = 0 < 1). Además la agenda no lo necesitaría. ✗
    F-4. Cuarteto CON R = 1 (en la línea): rompe H0 (δ = 0, u_m
         indefinido) — por eso H0 es ahora explícita. Con H0, ✗.
    F-5. ⚡ EL ATAQUE QUE FUNCIONÓ (contra la versión SIN H4): un fondo
         en la línea con densidad superexponencial (N_fondo(T) ~ e^{cT},
         c > δ) hace crecer coro_n más rápido que r_maxⁿ y la prueba se
         rompe: la cita llega y el coro la tapa igual, para siempre. El
         esqueleto de contraejemplo muestra que ALGUNA hipótesis de
         densidad es NECESARIA — no es un adorno. Con H4 declarada, el
         ataque muere: densidad polinomial ⟹ coro polinomial ⟹ el
         exponencial gana. ✗ (post-parche)
    F-6. Cuantificadores (tarea 7): revisados uno a uno en C — sin
         confusión existe/para-todo; «infinitas citas» no se usa (basta
         UNA cita ≥ T, que L3 da); «configuración finita» vs «todos los
         ceros»: el teorema es sobre m cuartetos finitos, y la Parte C
         lo aplica a cualquier subconjunto finito de ceros de ζ fuera
         de la línea — no afirma nada sobre infinitos cuartetos. ✗

## G. Problemas encontrados

    G-1. (mayor, corregido) H4 oculta: la cota del coro dependía de una
         propiedad del fondo no declarada en las hipótesis de la Parte A
         — exactamente la clase de fuga que la separación A/B/C existe
         para impedir. Declarada como H4; B gana B2. Acta parcheada.
    G-2. (menor, corregido) H0 (r_max > 1) implícita — ahora explícita;
         sin ella δ = 0 y el lema radial no arranca.
    G-3. (nota) la cota del coro exige n ≥ 3; la cita del ensamble
         cumple n ≥ n_rad,m ≥ 11 — satisfecha siempre, anotada.

## H. Veredicto final

    Estado ANTES del parche:  🟠 HABÍA UN SALTO (H4 sin declarar).
    Estado DESPUÉS del parche: 🟡 CASI SELLO — la cadena L1-L6 es, a
    juicio de esta auditoría interna, correcta bajo H0–H4, con los
    cuantificadores en regla y los ataques F-1..F-6 fracasados
    (post-parche). Las correcciones fueron declarativas (hipótesis), no
    estructurales.

    El sello 🟢 no lo emite el autor: la regla del laboratorio manda que
    la auditoría final independiente — la de Yui — decida si el
    resultado merece el sello de teorema dentro del alcance declarado
    (Parte A bajo H0–H4; Parte C para ζ con los inputs externos B1 y B2
    etiquetados).

---

*La cuenta del pancho, pagada: un salto encontrado y parcheado, dos
hipótesis con nombre, seis ataques fracasados, y el reloj abierto sobre
la mesa — engranaje por engranaje.*
