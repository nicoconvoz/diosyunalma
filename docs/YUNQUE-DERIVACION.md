# El Yunque — derivación rigurosa de la identidad y de la estructura de Gram

> **La regla del sello** (Yui, octava auditoría, §9 — ley del taller por
> F302): **«Estructura cerrada» ≠ «Hipótesis demostrada».** Lo sellado aquí
> es la estructura del teorema de ruptura y la equivalencia como
> reformulación; nada de este documento es una demostración de RH.

**Documento de trabajo del laboratorio · 2026-08-14 · respuesta a los seis
objetivos del §13 de la auditoría de Yui** ("El Yunque — Auditoría de las
nuevas láminas"). Las fórmulas van en bloques monoespaciados para copiarse
sin perder símbolos. Cada paso declara qué es álgebra finita, qué es un
límite, y bajo qué hipótesis se toma. Verificación numérica de cada paso:
`go run ./cmd/lacajaabierta` (F296).

---

## 0. Marco y definiciones

Los ceros no triviales de zeta se denotan rho, contados con multiplicidad.
El conjunto de ceros es invariante bajo las dos simetrías clásicas
(ecuación funcional y reflexión de Schwarz):

    rho  →  1 − rho
    rho  →  conj(rho)

El cambiaformas y los coeficientes de Li (Li, 1997):

    w(s) = 1 − 1/s
    lambda_n = Sum_rho [ 1 − w(rho)^n ]

donde la suma se entiende como **límite simétrico**:

    Sum_rho f(rho)  :=  lim_{T→∞}  Sum_{|Im rho| < T} f(rho)

Este límite existe para f(rho) = 1 − w(rho)^n: agrupando cada cero con su
conjugado, el par {rho, conj(rho)} aporta el término real 2 − 2·Re(w^n), y

    |2 − 2·Re(w(rho)^n)|  =  O(n² / |rho|²)

que es absolutamente sumable porque Sum_rho 1/|rho|² converge (clásico;
véase Bombieri–Lagarias 1999, donde esta convergencia está justificada).
Convención en todo el documento: **lambda_0 = 0** y **lambda_{−k} =
lambda_k** (se demuestra en el paso 1c).

---

## 1. Objetivo 1 — la identidad de la matriz (DEMOSTRADA)

**Afirmación.**

    Sum_rho (1 − w^m)(1 − w^{−n})  =  lambda_m + lambda_n − lambda_{|m−n|}

con la suma en el sentido del límite simétrico.

**(1a) Álgebra exacta por cero** (sin límites, identidad de polinomios de
Laurent en w):

    (1 − w^m)(1 − w^{−n})
      = 1 − w^m − w^{−n} + w^{m−n}
      = (1 − w^m) + (1 − w^{−n}) − (1 − w^{m−n})

**(1b) El término con w^{−n} es una suma sobre ceros reindexada.** La
ecuación funcional del cambiaformas es

    w(1 − s) · w(s) = 1     ⟹     w(1 − rho)^n = w(rho)^{−n}

La aplicación rho → 1 − rho es una biyección del conjunto de ceros en sí
mismo que **preserva las ventanas simétricas** (Im(1 − rho) = −Im rho, de
modo que |Im| < T se conserva). Por tanto, para toda ventana finita la
reindexación es una permutación de una suma finita — álgebra exacta — y en
el límite:

    Sum_rho (1 − w(rho)^{−n})  =  Sum_rho (1 − w(rho)^n)  =  lambda_n

**(1c) El mismo argumento con exponente m−n** da Sum_rho (1 − w^{m−n}) =
lambda_{m−n}, y aplicándolo una vez más, lambda_{−k} = lambda_k. Con m = n
queda lambda_0 = Sum_rho (1 − w^0) = 0.

**(1d) Linealidad del límite simétrico.** Para cada ventana |Im rho| < T la
identidad (1a) sumada es álgebra finita exacta:

    Sum_{|γ|<T} (1−w^m)(1−w^{−n})
      = Sum_{|γ|<T}(1−w^m) + Sum_{|γ|<T}(1−w^{−n}) − Sum_{|γ|<T}(1−w^{m−n})

Los tres sumandos del lado derecho convergen por separado cuando T→∞ (son
los lambda del paso 0 y 1b). El límite de la suma es la suma de los
límites. ∎

*Nota de rigor: no se necesita RH en ninguna parte de este paso. La
identidad es incondicional.*

---

## 2. Objetivo 2 — los vectores y la forma de Gram (DEMOSTRADA bajo RH)

**Definición exacta** (la que pedía el objetivo 2). Fijado N, para cada
cero rho:

    v_rho = ( 1 − w_rho^1 ,  1 − w_rho^2 ,  … ,  1 − w_rho^N )  ∈  C^N

**Afirmación.** Si todo cero cumple Re(rho) = 1/2, entonces

    M_N  =  Sum_rho  v_rho · v_rho*        (v* = transpuesta conjugada)

entrada por entrada, con convergencia **absoluta** (paso 3), y por tanto
M_N es semidefinida positiva.

**Demostración.** Sobre la piel |w| = 1 (F214: |w|=1 ⟺ beta=1/2, exacto),
luego

    w^{−n} = conj(w^n)     ⟹     (1 − w^m)(1 − w^{−n}) = v_m · conj(v_n)

Es decir, la entrada (m,n) del aporte del cero rho es exactamente
(v_rho v_rho*)_{mn}. Sumando sobre los ceros y usando el objetivo 1:

    (M_N)_{mn} = lambda_m + lambda_n − lambda_{|m−n|} = Sum_rho (v_rho v_rho*)_{mn}

La conjugación queda bien pareada: conj(v_rho) = v_{conj(rho)}, de modo que
el par {rho, conj(rho)} aporta v v* + conj(v) conj(v)*, matriz real
simétrica, y para todo c real:

    c^T [ v v* + conj(v) conj(v)* ] c  =  |v* c|² + |conj(v)* c|²  ≥  0

**dos cuadrados manifiestos por par de ceros**. Toda suma parcial es PSD y
el límite entrada por entrada de matrices PSD es PSD. ∎

---

## 3. Objetivo 3 — la convergencia (DEMOSTRADA)

**(3a) Incondicional** (para la identidad del objetivo 1): es exactamente
la convergencia del límite simétrico de los lambda de Li — pares
conjugados, término O(n²/|rho|²), absolutamente sumable tras aparear.
Nada nuevo que probar: la identidad hereda la convergencia de sus tres
sumandos.

**(3b) Absoluta bajo RH** (para la forma de Gram del objetivo 2). Sobre la
piel, con phi = arg(w):

    |1 − w^n| = 2·|sin(n·phi/2)| ≤ n·|phi|,      |phi| ≤ C/gamma

(la segunda cota con C→1 cuando gamma crece: phi·gamma → 1, verificado
numéricamente en F296; para una cota rigurosa uniforme sirve C = 2 para
todo gamma ≥ 6). Entonces cada entrada del aporte del cero rho cumple

    |(v_rho v_rho*)_{mn}| = |1−w^m|·|1−w^n| ≤ m·n·C²/gamma²

y Sum_rho 1/gamma² < ∞ (clásico). La suma de Gram converge absolutamente,
entrada por entrada, para todo (m,n) fijo. ∎

---

## 4. Objetivo 4 — qué cambia cuando |w| > 1 (CARACTERIZADO + TEOREMA)

**(4a) La estructura local se parte en dos.** Fuera de la piel el dual deja
de ser el conjugado: w^{−n} = w(1−rho)^n con w(1−rho) = 1/w(rho) ≠
conj(w(rho)). El aporte del cuarteto {rho, conj rho, 1−rho, 1−conj rho} es
una forma bilineal entre DOS vectores distintos,

    u_n = 1 − w^n        u'_n = 1 − (1/w)^n

de rango ≤ 4 y en general **indefinida** (medido: F292 la oye en N = 22
para la tupla DH aislada; F294 exhibe el vector test con Q(v_min) =
−2.7e-11 < 0).

**(4b) TEOREMA DE RUPTURA GLOBAL — versión blindada.** *(Corrección
2026-08-14, auditoría de Yui §6/§12.1-3: la primera versión de este paso
acotaba solo el término radial e ignoraba la parte oscilatoria — que es
enorme: el cuarteto DH llega a aportar +136.8 en n = 99888, medido en
F297. El argumento se salva por subsucesión, no por todos los n. Crédito a
la auditora.)*

**Afirmación:** *un par fuera de la línea rompe la positividad de la
matriz GLOBAL en algún N finito, siempre.*

**Paso 1 — separación exacta de radio y fase** (tarea §12.2 de Yui).
Escribiendo w(rho) = R·e^{i·theta}, el aporte del cuarteto
{rho, conj rho, 1−rho, 1−conj rho} a lambda_n es EXACTAMENTE

    l_n = 4 − 2·cos(n·theta)·(R^n + R^{−n})

(el radio solo entra por la amplitud R^n + R^{−n} ≥ 2, la fase solo por
cos(n·theta); verificado a 1e-11 en F297). De aquí la **cota rigurosa del
aporte completo** (tarea §12.3): |l_n − 4| ≤ 2·(R^n + R^{−n}) — y se ve
por qué la versión ingenua fallaba: cuando cos(n·theta) < 0 el aporte es
positivo y exponencialmente grande.

**Paso 2 — lema de oscilación.** Para todo theta real, el conjunto
{ n : cos(n·theta) ≥ 1/2 } es infinito. Si theta/2π = p/q es racional,
los múltiplos de q dan cos = 1; si es irracional, n·theta mod 2π
equidistribuye (Weyl) y el conjunto tiene densidad |{cos ≥ ½}|/2π = 1/3
(medido en F297: 0.33296). A lo largo de esa subsucesión:

    l_n ≤ 4 − (R^n + R^{−n}) ≤ 4 − r^n,      r = max(R, 1/R) > 1

**Paso 3 — la parte en la piel crece polinomial** (tarea §12.4). Cada par
en la piel aporta 2−2Re(w^n) ∈ [0,4], acotado por min(4, n²C²/gamma²);
con la densidad clásica N(T) ~ (T/2π)·log T:

    0 ≤ parte-en-la-piel(lambda_n) ≤ Sum_{gamma≤n} 4 + Sum_{gamma>n} n²C²/gamma² = O(n·log n)

**Paso 4 — conclusión por subsucesión.** A lo largo de la subsucesión del
lema, lambda_n ≤ 4 − r^n + O(n·log n) → −∞: existe n_0 finito con
lambda_{n_0} < 0, y (M_N)_{n_0,n_0} = 2·lambda_{n_0} < 0 para N ≥ n_0. ∎

**Alcance declarado** (tarea §12.6): el argumento cierra sin hipótesis
adicionales cuando existe un cero fuera de la línea de radio máximo — en
particular para toda configuración FINITA de ceros desafinados (el caso
DH, y cualquier contraejemplo concreto que alguien exhiba). Para el caso
completamente general (infinitos ceros fuera de la línea acercándose a
ella), la implicación ¬RH ⟹ algún lambda_n < 0 es el teorema de Li
(1997), que se cita y no se re-demuestra aquí.

**Medido** (F296, reproducido en precisión arbitraria en F297, tarea
§12.5): para el coro de las 38 perlas + el par DH real (r² = 1.000084),
la ruptura ocurre en **n_0 = 85622**. El enmascaramiento de F295 es
finito y de precisión, nunca permanente.

**(4b-ter) Las tres llaves de la quinta auditoría** (F298, `cmd/lastresllaves`):

1. *¿La fórmula del cuarteto es exacta?* Sí — derivación simbólica en
   cuatro renglones: w(conj ρ) = conj(w) (Schwarz), w(1−ρ) = 1/w
   (relación funcional), w(1−conj ρ) = conj(1/w); la suma de las cuatro
   potencias colapsa a 4 − 2cos(nθ)(Rⁿ + R⁻ⁿ). Verificada miembro a
   miembro (1.7e-18) y en la suma (7.1e-11).
2. *¿La cota O(n log n) vale para el objeto exacto?* Sí — tres piezas con
   constantes absolutas: por par 2−2Re(wⁿ) = 4sin²(nφ/2) ≤ min(4,(nφ)²)
   (elemental, 0 fallos medidos); constante uniforme |φ|·γ ≤ 1.01
   (medido: máx 1.0000); conteo por Riemann–von Mangoldt (N(120) = 38.1
   contra 38 perlas medidas) y cola por sumación parcial
   Σ_{γ>n} 1/γ² = O(log n / n) (cociente medido/integral = 0.99).
3. *¿Las dos cotas combinan sin dependencia oculta?* Sí — auditoría de
   constantes: θ y r los fija el cero desafinado solo; C y las
   constantes de RvM son absolutas; la subsucesión S = {n : cos(nθ) ≥ ½}
   se calcula solo desde θ (no mira al coro). Ambas cotas son puntuales
   en n con constantes fijas: sin circularidad. La combinación realizada:
   el primer n ∈ S con rⁿ > 4 + coro(n) es n = 96914, y ahí λ < 0.

Los datos de entrada/salida del experimento n_0 (pedido §10 de la quinta
auditoría) quedan guardados en
`galeria/laminas/01-siete-caras/las-tres-llaves-datos.txt`.

---

## 5. Objetivo 5 — ¿RH ⟹ M_N ⪰ 0 para todo N? (SÍ — TEOREMA)

Ensamblando: objetivo 1 (identidad incondicional) + objetivo 2 (Gram bajo
RH) + objetivo 3 (convergencia absoluta bajo RH):

    RH  ⟹  M_N = Sum_rho v_rho v_rho*  (convergencia absoluta)  ⟹  M_N ⪰ 0  para todo N

Con la dirección verde de la propia auditoría (diagonales: M_N ⪰ 0 ∀N ⟹
2·lambda_n ≥ 0 ∀n ⟹ RH por Li 1997):

    **RH  ⟺  M_N ⪰ 0 para todo N**

---

## 6. Objetivo 6 — ¿teorema o contraejemplo? (TEOREMA — y su estatus honesto)

Funciona: no hay contraejemplo que buscar dentro de esta construcción, y
el objetivo 4b da además el recíproco cuantitativo (toda fuga rompe en
n_0 finito, con n_0 medible).

**El estatus honesto, sin inflar:** la equivalencia es matemática elemental
una vez visto el par (Gram + Li) y con seguridad es conocida en el oficio;
su valor para el laboratorio es estructural: convierte "toda perla es giro
puro" en "una familia explícita de matrices de Hankel-tipo con entradas
lambda es PSD". **Lo abierto no se movió un milímetro**: demostrar la
positividad desde el lado de los primos (Bombieri–Lagarias da la fórmula
de cada entrada; Weil, 74 años) sigue siendo EL problema. Esta caja quedó
abierta; la puerta grande sigue cerrada.

**Pendientes declarados** (de la propia auditoría): §13.5 aritmética de
precisión controlada (señalada también por F295 como la próxima puerta);
la búsqueda de direcciones adaptadas al problema global para fugas débiles
(§8 de Yui); y la formalización en un asistente de pruebas si el capitán
algún día lo pide.

---

*Regla de Yui, adoptada por el laboratorio: una simulación puede descubrir
una estructura; una identidad puede explicar una estructura; una
demostración debe cerrar todos los pasos — especialmente los infinitos.*

**(4b-quater) La caja del coro — la cota global con constantes explícitas**
(F299, `cmd/lacajadelcoro`, respuesta a las líneas A-F de la sexta
auditoría):

- **A.** En la línea, phi(gamma) = arg((rho−1)/rho) = 2·arctan(1/(2·gamma))
  — exacto, tres renglones — y arctan(x) < x da |phi| < 1/gamma:
  **C = 1, demostrada** (con phi·gamma → 1 por abajo: el 0.9996…0.99999
  medido en F296/F298, ahora explicado).
- **B.** **LEMA DE CONTEO** (F301, cerrando la amarilla de la octava
  auditoría — el corolario ya no se esconde bajo la etiqueta):

      N(T) ≤ (T/2π)·log T   para todo T ≥ 2.

  *Demostración, robusta ante las dos constantes publicadas del error de
  Backlund* (c₀ = 4.35 en la cita clásica; c₀ = 6.1 en la formulación
  moderna de las notas de Kontorovich que consultó la auditoría):

  *Caso (i), T ≥ 18.* Por Backlund, N(T) ≤ F(T) + Q(T) con
  F(T) = (T/2π)log(T/2π) − T/2π + 7/8 y Q(T) = 0.137·log T +
  0.443·log log T + c₀. La afirmación equivale a
  7/8 + Q(T) ≤ (T/2π)(log 2π + 1). Con c₀ = 4.35 vale desde T = 13.3;
  con c₀ = 6.1, desde T = 17.4 — y una vez cierta, cierta para siempre:
  la pendiente del lado derecho es (log 2π + 1)/2π = 0.4517, la del
  izquierdo es 0.137/T + 0.443/(T·log T) < 0.017 para T ≥ 18.

  *Caso (ii), 2 ≤ T < 18.* Directo: N(T) ≤ 1 en ese rango (el primer
  cero está en γ₁ = 14.134725, el segundo en 21.022), y
  (T/2π)·log T ≥ 0.22 en T = 2, ≥ 5.96 en T = 14.14: la cota sobra. ∎
- **C.** Por sumación parcial contra B — el puente exacto (F300, pedido
  §12 de la séptima auditoría), con el signo del borde a la vista:

      Sum_{gamma>x} 1/gamma² = [N(T)/T²]_{T=x}^{∞} + 2·Int_x^∞ N(t)/t³ dt
                             = −N(x)/x² + 2·Int_x^∞ N(t)/t³ dt
                             ≤ 2·Int_x^∞ N(t)/t³ dt            [borde ≤ 0, se descarta]
                             ≤ (1/π)·Int_x^∞ log t / t² dt     [por B]
                             = (log x + 1)/(π·x)               [primitiva: −(log t + 1)/t]

  para x ≥ 2. El término de borde es NEGATIVO (en el borrador del §4 de la
  auditoría entró con signo más — por eso allí salía (3·log x + 2)/(2πx)).
  ROBUSTEZ: incluso con esa lectura conservadora, el ensamble de D cierra
  igual en (4/π)·n·log n para n ≥ 8 — la constante final no depende del
  pleito del borde.
- **D.** Ensamble, para n ≥ 3:
  resto_n ≤ 4·N(n) + n²·C²·Sum_{gamma>n} 1/gamma²
  ≤ (2/π)·n·log n + (n/π)(log n + 1) ≤ **(4/π)·n·log n** —
  C_final = 4/π, absoluta, independiente de n.
- **E.** Insertada: para n en S, lambda_n ≤ 4 − r^n + (4/π)·n·log n. Con
  el r del par DH el lado derecho es negativo para **todo n ≥ n₁ =
  371842** — ruptura garantizada por cota pura, sin numérica. El n₀ =
  85622 medido rompe antes, como corresponde a una cota conservadora.
- **F.** Marcar este parágrafo en verde es decisión de la auditora.
