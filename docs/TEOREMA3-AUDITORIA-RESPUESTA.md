# Respuesta a la auditoría del candidato T3 — D3, D4, D6 y el ataque hostil

**Para la auditora · 2026-08-15 · responde su NOTA F321 (auditoría T3
Robustez), punto por punto.** Su regla, adoptada: no se sella porque la
fórmula sea hermosa — se sella porque no encontramos dónde romperla.

---

## §4 de su nota — D3 sin saltos: g creciente en TODO [n_rad, ∞)

g(n) = e^{nδ} − (4/π)·n·log n − (2m+2), como función de variable REAL
n ∈ [n_rad, ∞). Es C² ahí (n_rad ≥ 11 > 0, log n definido). Tres pasos:

**(D3a) g″ > 0 en todo el intervalo, no solo en n_rad.**

    g″(n) = δ²·e^{nδ} − (4/π)·(1/n)

El primer término es ESTRICTAMENTE CRECIENTE en n (δ > 0 por H0) y el
segundo es estrictamente decreciente en valor absoluto (1/n decrece).
Luego g″ es estrictamente creciente, y por tanto

    ∀n ≥ n_rad:  g″(n) ≥ g″(n_rad) > 0        [R9 da la base]

**(D3b) g′ > 0 en todo el intervalo.** Por el teorema fundamental del
cálculo, para n ≥ n_rad:

    g′(n) = g′(n_rad) + ∫_{n_rad}^{n} g″(t) dt  ≥  g′(n_rad)  >  0
                        [integrando > 0 por D3a]   [R8 da la base]

**(D3c) g creciente.** De nuevo por el TFC:

    g(n) = g(n_rad) + ∫_{n_rad}^{n} g′(t) dt  ≥  g(n_rad)
                        [integrando > 0 por D3b]

con igualdad solo en n = n_rad. Los n que usamos son enteros dentro del
intervalo real: la monotonía sobre ℝ los cubre. **D3 cerrado: R9 aporta
la base de g″, la monotonía de g″ la extiende a todo el intervalo, y dos
integraciones anidadas bajan de g″ a g′ y de g′ a g. Sin saltos.** ∎

## §5 de su nota — D4: el MISMO n cumple las cuatro condiciones

Construcción explícita (la de la agenda, con cada cuantificador a la
vista). Sea T = n_rad = ⌈u·log u⌉ — entero, y ≥ 11 ≥ 1 (u ≥ 6). Sea
Q = ⌈2πT⌉. Por **L2 (Dirichlet)** existe n₁ con

    1 ≤ n₁ ≤ Q^m   y   ‖n₁·θᵢ‖ ≤ 2π/Q  ∀i ∈ {1,…,m}     [∃n₁; ∀i]

Sea J = ⌈T/n₁⌉ y **n = J·n₁**. Ese ÚNICO n cumple:

    (i)   n ≥ n_rad:    J ≥ T/n₁  ⟹  n = J·n₁ ≥ T = n_rad.
    (ii)  n ≤ N₀:       J − 1 < T/n₁  ⟹  n < T + n₁ ≤ n_rad + Q^m
                        y Q = ⌈2π·n_rad⌉ ≤ 2π·n_rad + 1
                        ⟹  n ≤ n_rad + (2π·n_rad + 1)^m = N₀.
    (iii) condición angular de L7 con ε = 1, PARA TODO i:
                        ‖n·θᵢ‖ = ‖J·(n₁θᵢ)‖ ≤ J·‖n₁θᵢ‖     [subaditividad
                        de la norma circular, iterada J−1 veces]
                        ≤ J·(2π/Q).
                        Ahora J = ⌈T/n₁⌉ ≤ T   [n₁ ≥ 1, T entero]
                        y 2π/Q ≤ 2π/(2πT) = 1/T   [Q ≥ 2πT]
                        ⟹  ‖n·θᵢ‖ ≤ T·(1/T) = 1 = ε.  ∀i. ✓
    (iv)  hipótesis de D5 preservadas en ese n:
                        L7 pide exactamente (iii) con ε = 1 ✓;
                        L5 pide n ≥ 3: n ≥ n_rad ≥ 11 ✓;
                        H4 es global, no depende de n ✓.

**Lo que la agenda posterior a T NO puede hacer:** empujar n fuera de
[n_rad, N₀] — (i) y (ii) lo encierran por construcción — ni perder la
condición angular — (iii) la deriva del múltiplo está acotada por 1
ANTES de elegir nada, para todos los i simultáneamente, porque n₁ es la
MISMA cita fina para las m fases (Dirichlet es simultáneo) y el factor J
es común. **D4 cerrado: un solo n, las cuatro propiedades, cada
desigualdad con su porqué.** ∎

## §6 de su nota — D6: la cadena, orientación por orientación

    −λₙ  ≥  g(n)                    [D5: λₙ ≤ coroₙ + (2m+2) − r_maxⁿ por L7
                                     con ε=1 (usa (iii)-(iv) de D4), y
                                     coroₙ ≤ (4/π)n·log n por L5+H4 (usa n≥3);
                                     sumar cotas superiores de λ da cota
                                     inferior de −λ: orientación ✓]
    g(n) ≥  g(n_rad)                [D3c, porque n ≥ n_rad por (i): ✓]
    g(n_rad) ≥ u^{3(m+1)} − u³      [EXACTAMENTE D1 y D2 y nada más:
                                     a ≥ A (D1: e^{n_rad·δ} ≥ u^{3(m+1)})
                                     y b ≤ B (D2: corchete ≤ u³)
                                     ⟹ a − b ≥ A − B: orientación ✓]
    u^{3(m+1)} − u³ = u³(u^{3m}−1) = Δ   [álgebra]

Ninguna hipótesis nueva entra en el último paso: D1 es R1 sin degradar,
D2 es la cadena R4+R5+R6 intacta. ∎

## §7 de su nota — el ataque hostil, ejecutado

**A. Extremos analíticos.**
- **δ → 1 (borde H3):** u ≥ 6 ⟺ δ ≤ 3(m+1)/6 = (m+1)/2, y H3 (δ ≤ 1)
  lo garantiza para todo m ≥ 1. En el borde absoluto (m=1, δ=1, u=6):
  n_rad = 11, y D1 da margen 11 − 6·log 6 = 0.249 > 0 en el exponente;
  D2 da 37.58 ≤ 216. Aguanta, con margen finito pero positivo.
- **δ → 0:** H0 estricta (rᵢ > 1) ⟹ δ > 0 ⟹ u, n_rad, N₀, Δ finitos.
  No hay división por cero posible bajo las hipótesis.
- **m = 1:** Δ = u³(u³−1) ≥ 216·215 = 46440 > 0. **m grande:** Δ crece;
  positiva y finita para todo m finito (H1).
- **θᵢ degenerados (0 o π):** la cadena no usa θᵢ ≠ 0 en ningún paso —
  Dirichlet, la agenda y L7 valen para fases arbitrarias (‖n·0‖ = 0 ≤ 1
  trivialmente). Sin ruptura. (Bajo H2 además θ = 0 es imposible, pero
  la prueba ni lo necesita.)

**B. Auditoría de techos (⌈·⌉): ¿alguno invierte una desigualdad?**
Hay tres techos en la cadena y los tres empujan en dirección SEGURA:
- n_rad = ⌈u·log u⌉ en D1: el techo AGRANDA n_rad·δ — es carga
  portante de la dirección correcta (en el borde u=6, el margen 0.249
  viene íntegramente del techo; con piso en vez de techo D1 fallaría
  por igualdad justa — el techo no es decorativo, y está del lado bueno).
- El mismo techo en D2: agranda el corchete — dirección PELIGROSA, pero
  es exactamente lo que R2 (n_rad ≤ 1.094·n*) acota, con R0 (n* ≥ 10.75)
  válido siempre que u ≥ 6. Verificado arriba en el borde.
- Q = ⌈2πT⌉ en D4: agranda Q, lo que ACHICA 2π/Q — dirección segura
  para (iii); y en (ii) se usa Q ≤ 2πT + 1, la cota superior correcta.

**C. ¿Δ demasiado grande para lo que D5 garantiza?** No hay hueco
posible: Δ no se postula — es el extremo INFERIOR derivado de la misma
cadena que D5 alimenta. La pregunta equivale a «¿puede g(n_rad) <
u^{3(m+1)} − u³?», y eso negaría D1 o D2, que son R1 y R4-R6 ya
auditadas. Las pérdidas restantes (constantes de techo de D2, la
degradación cos ≥ ½ del golpe en ε = 1) operan TODAS a favor del
enunciado: hacen la −λ real MÁS grande que Δ, jamás menor — cocientes
medidos 1.50 (m=2) y 1.84 (m=3), factores O(1).

**D. Batería hostil de extremos (evidencia, no prueba):** 42 casos con
m ∈ {1, 2, 5, 10, 50, 100, 1000} × δ ∈ {10⁻¹², 10⁻⁸, 10⁻⁴, 0.1, 0.99, 1},
verificados a 60 dígitos (mpmath, en espacio logarítmico): D1, D2 y
Δ > 0 sin una sola falla, incluido el borde u = 6.

**E. Nota de representación (declarada, no matemática):** para
m ≳ 20 (con δ ~ 10⁻⁴), log Δ > 709 y Δ desborda el float64 — Δ sigue
siendo matemáticamente finita; cualquier verificación numérica en ese
rango debe hacerse en espacio logarítmico, como la batería D.

**Resultado del ataque: no encontramos dónde romperla.**

## §8 de su nota — acordado

Los tres estantes quedan como evidencia. La justificación universal de
T3 es exclusivamente D1-D6 + lemas heredados. Ni la batería ni los
testigos entran en la prueba.

## §9 de su nota — la paramétrica en ε: dominio especificado, no declarada

Cuantificadores exactos del auxiliar: para toda configuración bajo
H0-H4, para todo ε con **0 < ε < √2**, y para todo n ≥ 3 que sea cita
de calidad ε (‖nθᵢ‖ ≤ ε ∀i):

    −λₙ ≥ 2(1 − ε²/2)·r_maxⁿ − coroₙ − 2ε²(m−1) − 4

El dominio importa: en ε = √2 el coeficiente 2(1−ε²/2) se anula y la
desigualdad queda vacía de contenido (sigue siendo verdadera pero no
muerde); para ε > √2 el paso L2 (cos ≥ 1−ε²/2 > 0) pierde su signo.
Queda como resultado auxiliar, sin declarar, a la espera de su auditoría.

## §10 — veredicto solicitado (autoevaluación; el sello es suyo)

    D3: cerrado sin saltos (§ arriba, dos integraciones anidadas).
    D4: cerrado — un solo n con las cuatro propiedades, traza completa.
    D6: cerrado — orientaciones verificadas, último paso usa solo D1+D2.
    Ataque hostil: ejecutado en A-E; sin ruptura encontrada.

    Autoevaluación: 🟢 T3 CERRADO bajo H0-H4 — sujeta a que su lupa
    no encuentre lo que la nuestra no vio. El sello, como siempre, es
    de la relojera. Todavía no.

---

*El cóctel sigue en la mesa, y el taller trae el próximo sorbo: D3 con
sus dos integraciones, D4 con su único n de cuatro propiedades, y un
ataque que terminó confirmando que hasta los techos empujan del lado
bueno.* 🍸🛡️
