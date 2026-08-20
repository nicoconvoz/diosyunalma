# Mapa de familias del Telar — Berry–Keating y su vecindario

**Entrega de las secciones 7, 10 y 13 del Pliego del Átomo — Fase II (Yui).** Base propia:
`docs/atomo/PLIEGO-ATOMO-ACTA.md` (Fase I, hallazgo F349, `cmd/elpliego`). Fecha: 2026-08-17.

Esto **no** demuestra RH, no afirma que el operador exista, y no acerca ni aleja la hipótesis. Es un mapa:
qué familia cae, bajo qué hipótesis exacta cae, y qué familia sigue en pie. Una estructura cerrada no es una
hipótesis probada.

## §0 — Reglas de lectura

Su §7 pone la condición que gobierna todo: *«No llamar 'descartada' a una familia más allá de las hipótesis
exactas que el no-go citado cubre.»* La obedecemos al pie, y donde nuestra propia acta se pasó de la raya lo
decimos en §12, en exhibición.

**Veredictos.** **DESCARTADA**: hay un enunciado —teorema ajeno, violación de R6, o medición con controles—
que la cierra *dentro de sus hipótesis explícitas*; fuera de ellas sigue viva. **MODIFICABLE**: falla por un
defecto localizado con dirección de reparación identificada. **ABIERTA**: la evidencia no la cierra.

**Estados de R6.** **Limpio**: no lee {γₙ} ni nada computado de ellos. **Condicional**: depende de si ζ/ξ
entran en D — Yui exige decidirlo antes de probar candidatos; acá no lo decidimos, lo marcamos. **Violado**:
lee los γₙ, o lee la fase de ζ sobre la recta crítica, que es lo mismo con otro nombre.

# PARTE I — Las siete preguntas de su §7

## 7(a) — La parte del conteo que xp SÍ reproduce

Reproduce el término suave **entero, constante incluida**. No «casi». El área de la región clásica, sin
aproximar, es T·ln(T/(l_x·l_p)) − T + l_x·l_p; con l_x·l_p = 2πħ y ħ = 1 queda Área/h = (T/2π)·ln(T/2πe) +
1, y la regla semiclásica ⟨N⟩ = Área/h − α/(4π) con α = π/2 resta exactamente 1/8 y deja el **7/8** de
Riemann (Berry & Keating, *SIAM Review* **41** (1999) 236, sec. 6). No reproduce la parte fluctuante, y dos
caveats van puestos: **(i)** la derivación es heurística — las órbitas truncadas son acotadas pero no
periódicas, así que no hay toro EBK genuino ni índice de Maslov de un lazo cerrado (Sierra,
arXiv:1601.01797); **(ii)** la constante no es robusta — los modelos xp general-covariantes H = U(x)p +
V(x)/p dan −1/2 (arXiv:1110.3203, *J. Phys. A* **45** (2012) 055209) y el compacto de Berry–Keating de 2011
no da ningún 7/8.

**El conteo suave no distingue familias: lo saca cualquiera de la familia xp. Es piso, no evidencia.**

## 7(b) — El problema del espectro continuo

Conviene matar una creencia difundida: **no es que xp «no sea autoadjunto».** Lo es, esencialmente, en
L²(0,∞), con índices de deficiencia ambos cero. El problema es peor: el espectro de H_BK sobre L²(ℝ₊,dx) es
**puramente continuo**, igual a toda la recta real, absolutamente continuo, de multiplicidad uno, y **no
tiene un solo autovalor**; las autofunciones generalizadas ψ_E(x) = x^(−1/2+iE) no son de cuadrado
integrable (Endres & Steiner, *J. Phys. A* **43** (2010) 095204, arXiv:0912.3183). Sin autovalores no hay
niveles, y sin niveles no hay γₙ. No es un problema técnico de dominio: es la razón por la cual toda la
familia necesita un paso adicional —confinar, perturbar o cambiar de objeto— y ese paso es donde se juega
todo.

## 7(c) — Qué cambia al confinarlo

**En el espacio de fases: no se puede.** La región {x ≥ l_x, p ≥ l_p, xp ≤ T} es una figura clásica; x y p
no conmutan, así que no existe proyección sobre ella. Los propios Berry y Keating llaman a cerrar el espacio
de fases un problema central sin resolver.

**En el espacio de posiciones: sí, y da la densidad equivocada.** Con u = ln x, xp se vuelve −i·d/du sobre
la caja: la cuerda de guitarra, k_n ≈ π·n/L, N(k) ≈ (L/π)·k, densidad = L/π **constante**, mientras los
ceros piden dN/dT = (1/2π)·log(T/2π), que crece sin techo. Por grande que se elija L, la recta y la curva se
cruzan una vez y se separan para siempre.

**El no-go, con su alcance exacto:** Endres & Steiner clasifican *todas* las realizaciones autoadjuntas de
H_BK sobre grafos cuánticos **compactos**, derivan fórmula de traza exacta y asintótica de Weyl, y concluyen
que **ni H_BK ni H²_BK pueden dar como autovalores los ceros no triviales**. Verificado por contenido; **no
pudimos verificar el número de teorema** que suele citarse. Cubre H_BK y H²_BK, sobre grafos compactos, con
cualquier realización autoadjunta — y nada más: ni otros hamiltonianos, ni dominios no compactos, ni
familias de dominios, ni límites. Confirmación independiente con matiz: Bolte, Egger & Keppeler (*J. Phys.
A* **50** (2017) 105201, arXiv:1610.06472) cuantizan Berry–Keating como operador de diferencias sobre una
red finita y periódica y obtienen densidad media **logarítmica**, pero sólo combinando el límite continuo,
el **límite de volumen infinito** y el semiclásico. La red finita fija no da el logaritmo; el logaritmo
aparece cuando la caja se va al infinito.

## 7(d) — Qué cambia si el confinamiento depende de E

Nuestra Fase I midió, por ventanas de altura y con L = 2π/brecha:

    altura T   brecha media   L medida   ln(T/2π)
    ─────────────────────────────────────────────
      177.8      1.900731      3.3057     3.3427
      351.8      1.572490      3.9957     4.0252
      550.3      1.408931      4.4595     4.4726
      734.7      1.320367      4.7587     4.7616
      910.1      1.261054      4.9825     4.9757

La mejor caja fija del tramo, L = 4.300, se equivoca hasta 20.2 % (minimax). La caja respira:
3.306 → 4.982, ×1.507, siguiendo ln(T/2π) sin ajustar nada. **Pero esa lectura es CONDICIONAL, y
acá lo corregimos:** L = 2π/brecha sólo es «una longitud» si ya se supuso que la máquina es −i·d/du sobre un
intervalo de largo L. Lo incondicional es la densidad, y la densidad es de Riemann; lo nuestro es la
verificación de que el tramo la sigue.

**Y hay contraejemplo publicado a la lectura amplia.** Berry & Keating 2011, el hamiltoniano compacto H = (x
+ 1/x)(p + 1/p): todas las órbitas clásicas **acotadas**; tras simetrizar, una familia de realizaciones
autoadjuntas sobre **todo el eje x positivo** etiquetadas por un ángulo α, con **espectro discreto de
energías reales**, y los **dos primeros términos** de la densidad asintótica iguales a los de la densidad de
alturas de los ceros (*J. Phys. A* **44** (2011) 285203). Hamiltoniano fijo, dominio fijo, espectro
discreto, densidad logarítmica, sin leer un solo γₙ. Lo que lo hace posible es que el símbolo **no es lineal
en p**: el área de {H ≤ E} crece como E·ln E, no como L·E. **La caja no tiene que respirar en el espacio;
puede respirar en el espacio de fases.** El confinamiento dependiente de E es *una* salida, no la única, y
ambas evaden Endres–Steiner porque el no-go cubre H_BK, no cualquier hamiltoniano.

## 7(e) — ¿Puede derivarse esa dependencia sin los γₙ?

Sí, y hay tres derivaciones publicadas que no leen los ceros: **(1)** por área de espacio de fases
(Berry–Keating 2011); **(2)** por límite de volumen (Bolte–Egger–Keppeler 2017); **(3)** por órbitas
cerradas añadidas a xp — Sierra & Rodríguez-Laguna, *PRL* **106** (2011) 200201, arXiv:1102.5356: H =
x(p+l_p²/p) tiene órbitas periódicas cerradas, que xp no tenía, y su espectro «coincide con los ceros de
Riemann **promedio**». **La respuesta honesta: la ley de conteo NO es el cuello de botella.** Está resuelta
de tres maneras R6-limpias desde 2011. Lo que ninguna tiene es el eco.

## 7(f) — El eco k·log p

**Lo que el canto no hace.** El examen de un solo espaciado contra Wigner no distingue aritmética de azar:

    ceros verdaderos    canto = 0.0919 ± 0.0211   (13 bloques)
    matriz GUE al azar  canto = 0.1097 ± 0.0304   (24 matrices)
    sorteo puro Wigner  canto = 0.1171 ± 0.0291   (400 sorteos)

    separación ceros vs. ruido puro: 0.82 σ por realización; sobre las medias
    t = 4.17 (sorteo puro), t = 2.09 (GUE); sin cajones (KS) t(ceros,GUE) = 1.43

**Lo que la varianza de número sí hace**, con los mismos 47 espaciados: Σ²(10) da ceros 0.327 ± 0.066, GUE
0.549 ± 0.112, sorteo SIN memoria 1.441 ± 0.797, con t(ceros vs sorteo) = 18.8. **Lo que el eco hace**, con
E(T) = (2/M)·Σₙ wₙ·cos(γₙT) evaluado en T = k·log p:

    ceros verdaderos, períodos k·log p    razón 36.467
    espectro GUE emparejado en densidad   razón  1.064            [acta §2]
    ceros verdaderos, períodos AL AZAR    razón 0.0069 ± 0.0018
                                          (máx 0.0115, 120 sorteos)

Cuatro órdenes de magnitud, y ninguno de los 120 sorteos se acercó. **Y el encuadre que nos corrigió el
refutador:** el 36× lo carga el denominador. Contra el control emparejado, los ceros suenan 2.87 veces más
fuerte EN los períodos aritméticos y 12.7 veces más callados AFUERA. **Lo raro no es el grito en log p: es
el silencio entre los primos.** Consecuencia: ninguna de las tres familias R6-limpias de 7(e) tiene eco
mostrado, así que por su regla —«un candidato sin eco no pasa a la siguiente fase»— quedan **abiertas pero
no aprobadas**.

**Falsa tensión que conviene despejar.** No hay conflicto entre «GUE» y «aritmética» para un espectro de
longitudes de primos: los períodos k·log p = log(p^k) son todos distintos por factorización única
(multiplicidad exactamente 1), {log p} es ℚ-linealmente independiente, y los pesos Λ(n)/√n cumplen la regla
de suma genérica de Hannay–Ozorio, dando K_diag(τ) = τ, el valor GUE (Bogomolny, arXiv:0708.4223). El
mecanismo Poisson dispara **sólo** para superficies hiperbólicas aritméticas, cuyo espectro de longitudes
tiene multiplicidades exponencialmente grandes. Caveat: la curva GUE completa, más allá de la aproximación
diagonal, descansa sobre la conjetura de pares de primos de Hardy–Littlewood. «GUE para zeta» es él mismo
conjetural.

## 7(g) — El signo

En la fórmula de traza de Selberg la suma sobre órbitas periódicas entra con **+1**; en la fórmula explícita
de Weil la suma sobre primos entra con **−2**. El signo **relativo** es libre de convención: no se arregla
normalizando. **No es un no-go demostrado** — Yui tiene razón en tratarlo como restricción de diseño. **La
observación es vieja**: Hejhal 1976, Berry 1986, a quien Connes cita; lo de Connes es la **interpretación**
— los ceros como espectro de **absorción**, el espacio de Pólya–Hilbert apareciendo en negativo. **Tiene
explicación a nivel de teorema en cuerpos de funciones**: los ceros viven en H¹, que entra en
Grothendieck–Lefschetz con signo (−1)¹ — o sea, lo que Yui pide («orientación, fase, índice o estructura
tipo absorción») ya está contestado en el caso análogo, y la respuesta es un **índice cohomológico**.
**Connes probó una fórmula de traza cuya validez es equivalente a RH** para funciones L con Grossencharakter
(*Selecta Math.* **5** (1999) 29–106): **reducción**, no demostración. Y hay un segundo desajuste, el
«problema asintótico» de Sierra: 2·sinh(m·λ_p/2) iguala p^(m/2) sólo cuando m → ∞.

La familia que lo tiene resuelto estructuralmente es la de Connes: si los ceros aparecen como **co-núcleo**
—lo que falta, no lo que hay— el signo sale solo. Sierra & Townsend (*PRL* **101** (2008) 110201,
arXiv:0805.4079) muestran que ese modelo de absorción emerge como el **límite del nivel de Landau más bajo**
de una partícula cargada en campos eléctrico y magnético cruzados: la traducción física más limpia del signo
que encontramos publicada.

# PARTE II — Las familias

## Las seis de su §10

**F1 · Diagonalizaciones que incorporan los γₙ — DESCARTADA, R6-violado.** H = diag(γ₁,…,γ₆₂₀) cumple R1
(autoadjunto), R2 (su estadística ES la de los ceros), R3 (620 contra N(T) = 619.6), R4 y R5 (razón 36.467):
los cinco requisitos con un objeto que ya sabe la respuesta. La decide **R6 y nada más**.

**F2 · Matrices finitas ajustadas a los ceros — DESCARTADA por R6, NO por densidad.** Nuestra acta §8
escribió «lejos de… cualquier matriz finita fija». Como razón de densidad está mal: la §3 de la misma acta
prueba que una matriz finita fija sí puede tener el conteo correcto. Cae porque el ajuste lee los ceros — y
con la razón correcta el alcance es más chico: **una matriz finita construida sin mirar los γₙ no está
cubierta por nada de lo que sabemos.**

**F3 · Grafos compactos de geometría fija — DESCARTADA dentro de Endres–Steiner, y sólo ahí.** Hipótesis
exactas: operador H_BK o H²_BK; dominio grafo cuántico compacto; realización cualquiera autoadjunta. Fuera
de esas tres no está descartado nada — ni (x+1/x)(p+1/p) sobre el semieje, ni x(p+l_p²/p), ni los límites de
familias de grafos de largo creciente, ni las realizaciones **no locales** sobre intervalos finitos (F9,
F12b). R6: **limpio**. Cae por el lado del conteo, no de la circularidad, que es la única forma decente de
caer.

**F4 · Laplacianos de superficies aritméticas — DESCARTADA CONDICIONALMENTE.** El espectro de longitudes de
una superficie hiperbólica aritmética tiene multiplicidades exponencialmente grandes, lo que produce
estadística **Poisson**, no GUE (Bogomolny–Georgeot–Giannoni–Schmit; Luo–Sarnak). **La hipótesis que falta
declarar:** el argumento supone que los γₙ son GUE, o sea Montgomery–Odlyzko, «creído, no probado» según
nuestra propia acta dos párrafos antes de usarlo como si estuviera probado. Enunciado correcto: *condicional
a Montgomery–Odlyzko*.

**F5 · Modelos que sólo produzcan GUE sin eco — DESCARTADOS con el Murciélago.** Los decide nuestra medición
con controles emparejados —misma cantidad de niveles, misma densidad media—:
1.064 contra 36.467, y 120 sorteos de períodos al azar que nunca pasaron de 0.0115. Alcance
honesto: **Nivel 2** en su escala, restricción que descarta familias, no teorema. Corolario que duele: una
matriz GUE no aporta información; la división útil no es «estadística + afinación» sino «espectro + eco», y
el eco es el único de los dos que cuesta.

**F6 · Parámetros que codifican indirectamente los γₙ — DESCARTADOS por R6, y son los más difíciles de
detectar.** Ejemplo vivo y publicado: Sierra arXiv:1601.01797 — fermión de Dirac sin masa en Rindler con
potenciales delta sobre los **enteros libres de cuadrados**. La geometría viene de la aritmética (bien),
pero la extensión autoadjunta está *«ajustada a la fase de la función zeta sobre la recta crítica»* (cita
verificada del resumen). **La violación de R6 no estaba en el espacio, ni en el potencial, ni en los primos:
estaba en la elección de la extensión autoadjunta.** Un parámetro de contorno puede codificar la respuesta
entera; eso debería entrar al protocolo de auditoría de su §11.

## Las que pidió verificar

**F7 · Sierra & Rodríguez-Laguna, H = x(p + l_p²/p)** — *PRL* **106** (2011) 200201, arXiv:1102.5356,
resumen verificado. **Da:** repara el defecto estructural de xp —«las trayectorias clásicas no son cerradas,
lo que deja el modelo incompleto»— con órbitas periódicas cerradas, y su espectro coincide con los ceros
**promedio**; se generaliza a funciones L de Dirichlet con distintas extensiones autoadjuntas. **Falla:**
«promedio» es la palabra clave — ley suave, no ceros individuales, y sin eco. **R6:** limpio para la parte
suave; **condicional** para las funciones L, porque no pudimos verificar si la elección de extensión lee ζ —
exactamente el punto de F6. **MODIFICABLE:** tiene el ingrediente que xp no tenía; le falta que los períodos
sean k·log p.

**F8 · Sierra & Townsend, niveles de Landau** — *PRL* **101** (2008) 110201, arXiv:0805.4079, resumen
verificado. **Da:** el modelo de absorción de Connes emerge como límite del nivel de Landau más bajo de una
partícula cargada en un plano con potencial eléctrico y campo magnético uniforme — realización física del
**signo**. **Falla:** exploratorio por confesión propia; «sugieren un papel» para los niveles superiores en
la parte fluctuante, sin establecerlo; no da ceros individuales. **R6:** hereda el de Connes, condicional a
D. **ABIERTA** como mecanismo del signo, **MODIFICABLE** como candidato a operador.

**F9 · Berry & Keating 2011, (x+1/x)(p+1/p)** — *J. Phys. A* **44** (2011) 285203, resumen verificado.
**Da:** órbitas acotadas, familia de realizaciones autoadjuntas sobre el semieje etiquetadas por α,
**espectro discreto real**, y los dos primeros términos de la densidad asintótica iguales a los de los
ceros. **Falla:** el tercer término, sin 7/8; y no hay primos en ningún lado — sin eco. **R6: LIMPIO.**
**ABIERTA — y es la familia que más nos corrige:** demuestra que «discreto + densidad logarítmica +
R6-limpio» es alcanzable, y por lo tanto que ése no era el cuello de botella.

**F10 · Bender, Brody & Müller** — *PRL* **118** (2017) 130201, arXiv:1608.03679. **Da:** un operador Ĥ de
límite clásico 2xp, no hermítico en el sentido usual pero con iĤ PT-simétrico de simetría rota; «si las
autofunciones obedecen una condición de contorno adecuada, entonces los autovalores asociados corresponden a
los ceros no triviales» (resumen); si se demostrara la autoadjunción vía operador métrico, eso demostraría
RH. **Falla — y está publicado como crítica formal:** Bellissard, arXiv:1704.02644, «da argumentos que
muestran que la estrategia propuesta por los autores para demostrar la Hipótesis de Riemann no funciona en
realidad» (resumen). Sus objeciones, según las fuentes que pudimos leer: el operador de momento sobre el
semieje no admite extensión autoadjunta; la zeta de Hurwitz en Re(z) = 1/2 no es de cuadrado integrable
sobre el semieje; y Ĥ se obtiene del operador de dilatación por un cambio de base **no unitario**. La
réplica (arXiv:1705.06767) sostiene que eso «ya fue discutido y no afecta las conclusiones», sin conceder
nada; por separado Müller (arXiv:1704.04705) muestra cómo *una versión* del operador se construye
rigurosamente sobre un espacio lineal bien definido — lo que confirma, leído al derecho, que la original
necesitaba esa aclaración. **R6:** violado o condicional según D — la condición de contorno *es* una
ecuación en ζ. **MODIFICABLE, en disputa abierta:** ofrece un **diccionario**, no un operador; no ofrece lo
único que se le pedía, la realidad del espectro.

**F11 · El espacio de clases de adeles de Connes.** **Da:** la única familia del mapa que ataca los tres
problemas estructurales a la vez — **el tamaño** (no un intervalo sino un espacio con infinitos lugares; el
corte no es un borde: es su Ruta B), **el signo** (los ceros como absorción, como co-núcleo: sale solo) y
**la aritmética** (los primos son lugares, no parámetros; las órbitas periódicas del sector de Riemann dan
la realización geométrica). Estado reciente verificado: Connes, Consani & Moscovici, *Zeta zeros and prolate
wave operators*, **Ann. Funct. Anal. 15 (2024) 87**, arXiv:2310.18423 — análogo semilocal del operador de
onda prolato, cuya parte positiva del espectro realiza los ceros bajos y cuyo espacio de Sonin (parte
negativa) realiza el comportamiento ultravioleta. **Falla:** no demuestra RH; la fórmula de traza es una
**reducción** —su validez equivale a RH— y la realización espectral todavía no es un operador autoadjunto
tipo Hilbert–Pólya con todos los ceros de forma incondicional. **R6:** la **construcción** usa adeles,
lugares y primos, no los γₙ; lo que lee los ceros es la **evaluación**. Ésa es la separación
construcción/evaluación que Yui pide en su §2, y es la familia donde esa separación salva al candidato en
vez de hundirlo. **ABIERTA — la puerta principal.**

**F12 · de Branges y la positividad de Weil.** Hay que separar dos cosas que suelen ir juntas. **(a) La ruta
de de Branges, en la forma exacta que Conrey & Li examinan: DESCARTADA.** arXiv:math/9812166, **IMRN 2000,
no. 18, 929–940**: dan ejemplos mostrando que las condiciones de positividad de de Branges —que implicarían
la GRH— **no se satisfacen** por las funciones que definen los espacios de Hilbert con núcleo reproductor
asociados a ζ. Hipótesis exacta: **esas** condiciones, **esos** espacios; no es un no-go sobre espacios de
funciones enteras en general. **(b) La positividad de Weil como criterio: ABIERTA, con movimiento
reciente.** Connes–van Suijlekom y Connes–Consani–Moscovici convierten la forma de Weil en objetos finitos
—matrices de Galerkin para un corte de primos y una banda—. Y en junio de 2026, Masatoshi Suzuki, *Weil's
quadratic form via the screw function*, arXiv:2606.09096 (resumen verificado directamente), estudia esa
forma mediante funciones continuas en vez de distribuciones y **formula una conjetura**: que un operador
autoadjunto cuyos autovalores sean las partes imaginarias de los ceros puede obtenerse como **límite**,
cuando a → ∞, de operadores autoadjuntos provenientes de realizaciones **no locales** del operador
diferencial de primer orden sobre [−a, a]. Todo obtenido **sin suponer RH**; es conjetura, no teorema, y el
propio resumen lo dice. **Por qué nos importa más que ninguna otra cosa del mapa:** ese objeto es
literalmente **la caja que respira** — la respiración es el parámetro a, y lo que reemplaza a la condición
de contorno fija es la **no localidad**. Queda fuera de Endres–Steiner por dos vías: no es un compacto fijo
y la realización no es local.

**F13 · Sierra: primos como espejos en Rindler** — arXiv:1404.4252 y su continuación arXiv:1601.01797.
**Da:** lo más cercano publicado a su Ruta C — fermión de Dirac sin masa en Rindler con potenciales delta
interpretados como **espejos parcialmente reflectantes en movimiento**, colocados según los primos, con
períodos **log p** en el tiempo propio del observador acelerado; H_BK emerge de las componentes quirales;
los ceros aparecen como autovalores discretos sumergidos en el continuo; clase AIII / GUE quiral. **La
geometría efectiva sale de los primos.** **Falla:** el resultado se obtiene para valores del corrimiento de
fase relacionados con la fase de ζ, y la continuación de 2016 lo dice sin ambigüedad. **R6: VIOLADO tal como
está publicado, pero la violación está LOCALIZADA** en la fase de la extensión; el resto —espejos en log p,
Rindler, emergencia de H_BK— es limpio. **MODIFICABLE, con el defecto mejor localizado del mapa.** La
pregunta de Fase III es de una línea: ¿existe una elección de fase derivable de los primos que dé lo que la
fase de ζ da a mano?

**F14 · LeClair, flujo espectral** — arXiv:2406.01828, *Adv. Theor. Math. Phys.* (2025). Los niveles E_n(σ)
tienden a los ceros cuando σ → 1/2, a partir de una matriz S construida con el **producto de Euler**, cuya
unitariedad es la condición que hace hermítico al hamiltoniano. **No pudimos verificar** si la ecuación
trascendente que define E_n lee la fase de ζ. **R6 indeterminado ⇒ SIN CLASIFICAR.** La anotamos porque el
producto de Euler como origen del operador es el tipo de input que la Ruta C necesita.

# PARTE III — Tabla final

| Familia | Qué da | Qué falla | R6 | Veredicto | Hipótesis exacta que la decide |
|---|---|---|---|---|---|
| **F1** diag(γₙ) | R1–R5 completos | no construye: copia | violado | **DESCARTADA** | R6, sin más |
| **F2** matriz finita ajustada | conteo correcto posible | el ajuste lee los ceros | violado | **DESCARTADA** | R6 — **no** densidad (nuestra §3 lo refuta) |
| **F3** grafo compacto geom. fija | R1 y espectro discreto | densidad constante vs. log | limpio | **DESCARTADA** | Endres–Steiner: **H_BK o H²_BK**, **grafo compacto**, **cualquier** realización autoadjunta |
| **F4** laplaciano sup. aritmética | espectro discreto, aritmética | multiplicidad exponencial → Poisson | limpio | **DESCARTADA (condicional)** | condicional a **Montgomery–Odlyzko** (γₙ = GUE) |
| **F5** GUE sin eco | estadística de un nivel | 1.064 contra 36.467 | limpio | **DESCARTADA** | Murciélago con control emparejado (Nivel 2, no teorema) |
| **F6** parámetros que codifican γₙ | lo que se le pida | circular por el contorno | violado | **DESCARTADA** | R6 aplicado a la **extensión autoadjunta**, no sólo al espacio |
| **F7** x(p+l_p²/p) | órbitas cerradas + ceros promedio | ni ceros individuales ni eco | limpio / cond. | **MODIFICABLE** | «promedio» ≠ ceros; faltan períodos k·log p |
| **F8** niveles de Landau | absorción de Connes como LLL | exploratorio; sin ceros | condicional | **ABIERTA** (signo) | los autores «sugieren» la parte fluctuante |
| **F9** (x+1/x)(p+1/p) | discreto + densidad log, dominio fijo | sin 7/8, sin eco | **limpio** | **ABIERTA** | fuera del no-go: **no es H_BK**; símbolo no lineal en p |
| **F10** Bender–Brody–Müller | diccionario ceros ↔ contorno | no autoadjunto; crítica de Bellissard | violado / cond. | **MODIFICABLE (en disputa)** | contorno = ecuación en ζ; ψ no de cuadrado integrable en Re = 1/2 |
| **F11** adeles de Connes | tamaño + signo + primos a la vez | reducción, no demostración | limpio en construcción | **ABIERTA — puerta principal** | la validez de la fórmula de traza **equivale** a RH |
| **F12a** de Branges (forma Conrey–Li) | positividad ⇒ GRH | no se satisface para ζ | limpio | **DESCARTADA** | Conrey–Li, IMRN 2000: **esas** condiciones, **esos** espacios |
| **F12b** positividad de Weil | criterio equivalente a RH | operador conjeturado, no probado | limpio | **ABIERTA** | Suzuki 2026: **conjetura**; límite a→∞; realizaciones **no locales** |
| **F13** primos como espejos (Rindler) | L(E) desde los primos; H_BK emergente | fase ajustada a ζ crítica | violado, **localizado** | **MODIFICABLE** | la violación está sólo en la fase de la extensión |
| **F14** flujo espectral (Euler) | E_n(σ) → ceros | no verificado si lee la fase de ζ | **indeterminado** | **SIN CLASIFICAR** | pendiente de auditoría de inputs |

# PARTE IV — Las puertas abiertas

Lo que la evidencia **no** cierra, en orden de cuánto nos corrige:

1. **La densidad ya no es el problema.** F9 lo demuestra: hamiltoniano fijo, dominio fijo, espectro
   discreto, densidad logarítmica, R6-limpio, publicado en 2011. Nuestra Fase I puso el peso ahí; la
   literatura dice que está en la aritmética. **Es la corrección más grande de este documento.**
2. **Símbolo no lineal sobre dominio fijo.** El ln E no exige una caja que respire en el espacio: exige un
   recinto de fases cuya área crezca como E·ln E. Hay al menos dos maneras (F9, y el límite de la red de
   F3), y ninguna toca R6.
3. **La no localidad como sustituto del borde.** Suzuki (F12b): límite de realizaciones no locales de
   −i·d/dx sobre [−a,a] con a → ∞. Evade Endres–Steiner por dos vías simultáneas; es la formalización más
   cercana a nuestra «caja que respira», y es conjetura abierta, no teorema.
4. **La fase derivada de los primos.** F13 tiene la geometría correcta y una única violación de R6,
   localizada. ¿Se puede derivar esa fase de los primos en lugar de leerla de ζ?
5. **Órbitas cerradas con períodos aritméticos.** F7 resolvió «órbitas cerradas» sin resolver «períodos
   k·log p». Nadie que hayamos verificado tiene las dos cosas en el mismo objeto R6-limpio.
6. **El signo como índice, no como accidente.** En cuerpos de funciones sale de que los ceros viven en H¹;
   F8 lo traduce a niveles de Landau. ¿Qué índice cohomológico juega ese papel sobre ℚ?
7. **El silencio entre los primos.** Nuestro hallazgo más raro y el menos explotado: los ceros son 12.7
   veces **más callados** fuera de los períodos aritméticos que un control GUE emparejado. Ningún modelo del
   mapa explica el silencio; todos apuntan al grito.

# §12 — Correcciones a nuestra propia acta, en exhibición

1. **«Ningún dominio compacto fijo puede servir, y eso está probado»** (acta §8) — **over-reach**. Lo
   probado es exactamente Endres–Steiner; Berry–Keating 2011 es contraejemplo publicado a la lectura amplia.
2. **«Lejos de… cualquier matriz finita fija»** (acta §8) — **razón equivocada**. Cae por R6, no por
   densidad; la §3 de la misma acta lo refuta.
3. **«Endres–Steiner 15.6»** en el sello §9, mientras la §5 dice que el número no pudo verificarse —
   **inconsistencia interna**. Se cita por contenido, nunca por número.
4. **«L(E) = ln(E/2π)» como requisito del candidato** (acta §8) — **condicional, no incondicional**. L =
   2π/brecha sólo es una longitud si ya se supuso −i·d/du sobre un intervalo. Lo incondicional es la
   densidad, y es de Riemann.
5. **«Cualquier construcción que pase por el laplaciano de una superficie aritmética no puede tener a {γₙ}
   como espectro»** (acta §6(b)) — **falta la condición**: supone Montgomery–Odlyzko, que la misma acta
   llama «creído, no probado» dos párrafos antes.

# §13 — Lo que no pudimos verificar

- **El número de teorema de Endres–Steiner** («15.6»): contenido verificado, número no. Citado siempre por
  contenido.
- **Los detalles internos de la crítica de Bellissard** (arXiv:1704.02644): verificamos su resumen textual,
  pero el PDF no se pudo extraer, así que las tres objeciones específicas que citamos provienen de **fuentes
  secundarias**.
- **La constante −1/2 de los modelos general-covariantes de Sierra** (arXiv:1110.3203): paper y referencia de
  revista verificados, pero el resumen público no menciona el término constante. Lo arrastramos de Fase I.
- **La ausencia de 7/8 en Berry–Keating 2011:** el resumen dice «los **dos primeros** términos»; que el
  tercero no sea 7/8 es **inferencia nuestra**, no afirmación del paper.
- **F14 (LeClair) entera:** sin verificar si su ecuación trascendente lee la fase de ζ no hay veredicto de
  R6, y sin R6 no hay clasificación.
- **Discrepancia numérica interna, sin resolver:** el eco del control GUE emparejado figura como **1.064** en
  el acta §2 y como **1.081** en el briefing de Fase II. Usamos el del acta, que es el registro escrito. Hay
  que rehacer esa corrida y fijar el número.
- **Preprints de 2026 sobre realizaciones numéricas de la forma de Weil** aparecieron en la búsqueda pero
  **no los auditamos**. Sólo verificamos directamente el de Suzuki (arXiv:2606.09096).

# §14 — Sello

    Se cierra:   el mapa. Quince entradas, cada veredicto con la hipótesis exacta que lo
                 sostiene y ninguna descartada más allá de ella.
    Se corrige:  cinco sobre-alcances propios de la Fase I, en exhibición.
    Se mueve:    el peso del programa. La densidad logarítmica está resuelta desde 2011
                 de tres maneras R6-limpias. El cuello de botella nunca fue la caja: es
                 el eco.
    Se abre:     siete puertas, y la principal —adeles— tiene el tamaño, el signo y los
                 primos a la vez, y no demuestra nada.
    NO se hizo:  nada de esto prueba RH ni la refuta, ni la acerca ni la aleja un
                 milímetro. Una estructura cerrada no es una hipótesis probada.

    Todavía no.

⚛️🦇📏
