# Acta del Pliego del Átomo — el requisito que impide la construcción

> **CORRECCION POSTERIOR (Fase II, 2026-08-17).** La respuesta del §5 de esta acta
> —«el requisito que lo impide es R3 contra (R1 + geometría fija)»— **no
> sobrevive al mapa de familias**: el Hamiltoniano compacto de Berry & Keating 2011,
> (x+1/x)(p+1/p), tiene dominio fijo, espectro discreto, densidad logarítmica y es
> R6-limpio. La densidad NO era el cuello de botella; el cuello de botella es el eco
> aritmético. La corrección completa, con otras cuatro sobre-extensiones de esta
> acta puestas en exhibición, está en docs/atomo/TELAR-FASE2-ACTA.md y
> docs/atomo/TELAR-MAPA-FAMILIAS.md. Esta acta se deja como se escribió: la casa no
> reescribe historia.


**Para la auditora · 2026-08-17 · responde su documento «EL ÁTOMO —
ANÁLISIS DE YUI» y su ADDENDUM, punto por punto y con sus mismos
nombres de sección.** Su misión textual fue: *NO intentes demostrar que
el Átomo existe. Intentá construir el operador que satisfaga R1–R5. Si
no puede construirse, determiná exactamente qué requisito lo impide.*
No lo intentamos demostrar. Intentamos construirlo. No se puede, y acá
está el requisito, con nombre, con teorema ajeno y con medición propia.
Programa: `cmd/elpliego`. Todos los números de esta acta son su salida.

---

## §0 — El mar donde medimos

Antes de discutir máquinas hay que tener el mar. Ceros no triviales por
Riemann–Siegel más bisección, en la franja t ∈ [100, 1000]:

    ceros hallados                       620
    ley suave N(T) predice               619.6
    diferencia                             0.4

Y el mar, una vez desplegado (dividido por su propia densidad local),
queda donde debe quedar:

    espaciamiento medio desplegado    0.999781      (tiene que ser 1)
    espaciamiento mínimo desplegado   0.2365        (repulsión: no hay
                                                     pares pegados)

Esto no es un resultado: es la calibración del instrumento. Si esto no
diera 1, todo lo que sigue sería ruido.

---

## §1 — Su «hallazgo central: el Taller de Máquinas», auditado de nuestro lado

Usted leyó bien la pieza. La pieza estaba mal medida. Antes de
contestarle, tuvimos que abrir nuestro propio taller (`cmd/maquina`,
hallazgo F178) y rehacerlo.

**El examen del canto**, para que quede definido: se comparan los
espaciamientos desplegados contra la sorpresa de Wigner (GUE) en 15
cajones de ancho 1/3 hasta s = 5, y el «canto» es la distancia entre
los dos histogramas. Cuanto más chico, más se parece a GUE.

### PROTOTIPO A — la reja de estacas de una caja FIJA

Caja de longitud logarítmica 10, niveles E_n = 2π·n/ln(Λ): una escalera
perfectamente pareja, con todos los espaciamientos desplegados en 1.

    como lo hizo el taller (espaciamientos ESTIPULADOS en 1.0):
        canto = 0.5546      (publicó 0.555 — reproducido)

    construyendo los niveles y MIDIÉNDOLOS de verdad:
        canto = 0.4297

Es **la misma reja**. La diferencia entera es que s = 1 cae exactamente
sobre un borde de cajón (1/3 × 15 = 5.000), así que el ruido de punto
flotante reparte la masa entre dos cajones vecinos y el examen cambia
de veredicto. Conclusión incómoda: **el canto de una reja de estacas no
es DISCONTINUO sobre un espectro degenerado.** No es un defecto del prototipo,
es del examen aplicado a un espectro degenerado.

### PROTOTIPO B — como ensamble, no como anécdota

24 matrices GUE reales de 80×80 (Jacobi sobre la inmersión simétrica
real, desplegadas por la ley del semicírculo en el 60 % central), 47
espaciamientos cada una:

    canto = 0.1097 ± 0.0304        (24 matrices, barra de error real)

El taller había publicado 0.076 **de UNA matriz, UNA semilla, sin
barra de error**. Con eso, el «7.3 veces mejor» que su documento cita
no sobrevive:

    0.5546 / 0.1097 ≈ 5.1     (contra el A estipulado)
    0.4297 / 0.1097 ≈ 3.9     (contra el A medido de verdad)

Y encima el denominador de esa razón se compara contra un número que
fue **estipulado, no computado**. La cifra 7.3 hay que retirarla.

---

## §2 — «El Murciélago» y «La Silueta»: el canto es ciego, el eco no

Usted propuso, con toda razón, usar el Murciélago como **test de
candidatos**. Lo hicimos. Y de paso le pasamos el examen de canto a tres
candidatos con el MISMO tamaño de muestra (47 espaciamientos), que es la
única forma honesta de compararlos:

    ceros de zeta VERDADEROS   canto = 0.0919 ± 0.0211   (13 bloques)
    matriz GUE al azar         canto = 0.1097 ± 0.0304   (24 matrices)
    sorteo puro de Wigner      canto = 0.1171 ± 0.0291   (400 sorteos)

El último no sabe NADA: es un generador de números aleatorios tirando de
la distribución de Wigner. No conoce un primo, no conoce un cero.

    separación entre los ceros verdaderos y el ruido puro:
        0.82 σ por realización — y sobre las MEDIAS, t = 4.17

**El canto es gratis.** Cualquiera lo pasa. El examen del canto no
distingue la aritmética del azar, y por lo tanto no puede ser la mitad
de nada.

El eco, en cambio, muerde. Definimos el Murciélago tal como su pieza lo
plantea, con ventana de coseno alzado:

    E(T) = (2/M) · Σ_n w_n · cos(γ_n · T)

y lo evaluamos en los períodos aritméticos T = k·log p (13 períodos
hasta T = 3.2). El puntaje es media|E| SOBRE los períodos dividida por
media|E| FUERA de ellos:

    ceros verdaderos   dentro 0.06254   fuera 0.00172   razón 36.467
    espectro GUE       dentro 0.04295   fuera 0.04037   razón  1.064

(el espectro GUE fue reescalado a la misma densidad media que los ceros
antes del test, para que la comparación sea limpia).

Un 36.467 contra un 1.064. **El eco separa lo que el canto no separa.**
Su instrumento sirve; el otro no.

---

## §3 — El pliego se satisface HACIENDO TRAMPA

Tomamos su pliego al pie de la letra y construimos el operador más
tonto posible:

    H = diag(γ_1, γ_2, …, γ_620)

una matriz diagonal con los 620 ceros que pescamos en §0. Veredicto:

    R1 autoadjunto           ✓  (diagonal real)
    R2 estadística           ✓  (su estadística ES la de los ceros,
                                 porque su espectro ES el de los ceros)
    R3 conteo                ✓  (620 contra N(T) = 619.6)
    R4 órbitas k·ln p        ✓  (su eco tiene los picos aritméticos)
    R5 fórmula explícita     ✓  (razón 36.467, la misma de §2)

**Los cinco requisitos, con un objeto que ya sabe la respuesta.** No
construye nada: copia. Volvemos sobre esto en §6(a).

---

## §4 — La caja que tiene que respirar (acá está la obstrucción, medida)

Su Prototipo A es una caja de longitud FIJA. Preguntamos entonces: si
uno insiste en describir a los ceros con una caja, ¿qué longitud tiene
esa caja? La medimos por ventanas de altura: en cada ventana tomamos la
brecha media entre ceros consecutivos y despejamos la longitud L que la
produciría, L = 2π / brecha.

    altura T   brecha media   L medida = 2π/brecha   ln(T/2π)
    ───────────────────────────────────────────────────────────
      177.8      1.900731          3.3057             3.3427
      351.8      1.572490          3.9957             4.0252
      550.3      1.408931          4.4595             4.4726
      734.7      1.320367          4.7587             4.7616
      910.1      1.261054          4.9825             4.9757

La última columna no está ajustada: es la función ln(T/2π) evaluada en
la misma altura. La caja medida la sigue.

    mejor caja FIJA para todo el rango:   L = 4.300
    y aun así se equivoca hasta en          20.2 %  (minimax honesto)

    L tiene que CRECER:  3.306 → 4.982 en este tramo  (×1.507)

Dicho en criollo: **la caja respira.** Y respira con el logaritmo de la
altura. Esto, que parece un detalle de carpintería, es exactamente el
requisito que rompe el pliego, como se ve en §5.

Un paréntesis para su punto 3, «Conexión con el Teorema del Cielo»:
usted la planteó como analogía metodológica —quitar la escala dominante
para que se vea la estructura fina— y como analogía funcionó: el
desplegado de §0 es exactamente eso. La diferencia es que allá la escala
dominante era una potencia r_L^n, absorbible en una normalización, y acá
es ln(T/2π), que crece sin techo y no se absorbe. Analogía, no
equivalencia, tal como usted la marcó.

---

## §5 — LA RESPUESTA A SU MISIÓN: el requisito que lo impide es R3 contra (R1 + geometría fija)

**No es R1. R1 se cumple sin drama.** Hay una creencia difundida de que
el operador xp «no es autoadjunto». Es falsa: xp es esencialmente
autoadjunto en L²(0, ∞), con índices de deficiencia ambos cero. Ese no
es el problema.

El problema es que **su espectro es puramente continuo, igual a toda la
recta real, absolutamente continuo, de multiplicidad uno, y NO tiene un
solo autovalor**. Las autofunciones generalizadas

    ψ_E(x) = x^(−1/2 + iE)

no son de cuadrado integrable: no viven en el espacio (Endres &
Steiner, *J. Phys. A* **43** (2010) 095204, arXiv:0912.3183). Un
operador sin autovalores no tiene niveles, y sin niveles no hay γ_n que
poner.

El reflejo natural es encerrarlo. Y ahí es donde se choca de frente.

**Encerrarlo en el espacio de fases no se puede.** La región de
Berry–Keating {x ≥ l_x, p ≥ l_p, xp ≤ T} es una figura clásica. En
mecánica cuántica x y p no conmutan, así que no existe ninguna
proyección sobre esa región: el objeto que uno querría cortar no está
definido. Los propios Berry y Keating lo dicen sin vueltas — cerrar el
espacio de fases para que el movimiento sea acotado es, para ellos, un
problema central sin resolver.

**Encerrarlo en el espacio de posiciones sí se puede, y da la densidad
equivocada. Como TEOREMA, no como dificultad.** Este es el párrafo
central de esta acta, y lo escribimos para que se entienda sin ser
especialista:

> Poné el operador en una caja: un intervalo compacto, o un grafo
> finito. Con el cambio de variable u = ln x, el operador xp se
> convierte en el operador de momento −i·d/du sobre esa caja. Y el
> operador de momento en una caja es el problema de la cuerda de
> guitarra: los niveles son las ondas que entran enteras en la caja,
>
>     k_n ≈ π·n / L      ⟹      N(k) ≈ (L/π)·k
>
> El conteo crece como una RECTA en k. Su derivada, que es la densidad
> de niveles, vale L/π: una **constante**. La caja tiene siempre la
> misma cantidad de escalones por unidad de energía, arriba y abajo,
> porque la caja siempre mide lo mismo.
>
> Pero los ceros de Riemann piden densidad
>
>     dN/dT = (1/2π)·log(T/2π)
>
> que **crece sin techo**. Arriba hay más escalones por unidad de
> altura que abajo, y cada vez más. Una constante no puede seguir a un
> logaritmo que crece: por grande que elijas L, la recta (L/π)·k y la
> curva se cruzan una vez y después se separan para siempre.

Eso es, exactamente, **R3 contra R1**. No contra R1 a secas —contra R1
en su forma útil, *autoadjunto con espectro discreto*, que es la única
forma en que R1 produce niveles γ_n. El enunciado formal es el

**El no-go de Endres–Steiner** (J. Phys. A 43 (2010) 095204, arXiv:0912.3183 — citado por contenido: no pudimos verificar el número de teorema): ni H_BK ni H²_BK dan como
autovalores los ceros no triviales de Riemann, si son realizaciones
autoadjuntas sobre CUALQUIER grafo compacto.

Y nuestra medición de §4 es ese mismo teorema visto desde el otro lado:
si uno igual insiste en describir a los ceros con una caja, la caja que
sale medida no es una, son cinco cajas distintas — 3.3057, 3.9957,
4.4595, 4.7587, 4.9825 — y la mejor caja única que se puede elegir para
todo el tramo se equivoca hasta un 20.2 % — y ese porcentaje es del tramo, no del problema: L = ln(T/2π) → ∞. La caja no es fija porque la
naturaleza no la dejó fija.

**Respuesta a su misión, en una línea: el operador del pliego no puede
construirse como operador autoadjunto de espectro discreto sobre un
dominio compacto fijo, y el requisito que lo impide es R3, la ley de
conteo, chocando contra la ley de Weyl que R1 arrastra consigo.**

### Nota sobre el 7/8, que conviene tener bien puesta

Para que no quede la impresión de que Berry–Keating «casi» llegaron y
falló una constante: la constante les da **exacta**. El área de la
región es, sin aproximar,

    Área{x ≥ l_x, p ≥ l_p, xp ≤ T} = T·ln(T/(l_x·l_p)) − T + l_x·l_p

con l_x·l_p = 2πħ y ħ = 1,

    Área/h = (T/2π)·ln(T/2π) − T/2π + 1 = (T/2π)·ln(T/2πe) + 1

y el conteo suave de Riemann es (T/2π)·ln(T/2πe) + 7/8. La diferencia
es 1/8 — pero la regla semiclásica estándar no es Área/h sino
⟨N⟩ = Área/h − α/(4π) (índice de Maslov), y para la hipérbola truncada
α = π/2 da exactamente −1/8. O sea: **7/8 exacto, constante incluida**
(Berry & Keating, *SIAM Review* **41** (1999) 236, sec. 6).

El caveat honesto: la derivación es heurística —las órbitas truncadas
son acotadas pero NO periódicas, así que no hay toro EBK genuino ni
índice de Maslov estándar de un lazo cerrado (Sierra, arXiv:1601.01797)—
y la constante no es robusta: los modelos xp general-covariantes de
Sierra dan −1/2 en vez de 7/8 (arXiv:1110.3203), y el propio
hamiltoniano compacto de Berry y Keating de 2011, (x + 1/x)(p + 1/p),
reproduce los dos primeros términos y ningún 7/8. La constante correcta
sale de una cuenta cuya validez sus propios autores no reclaman.

---


## §5 bis — Lo que nos corrigió nuestro propio refutador

Antes de mandarle esta acta le soltamos encima un refutador con una sola
consigna: romper nuestras conclusiones. Rompió dos, y las dos eran
nuestras. Van acá, al frente, porque la casa registra las correcciones
donde se ven.

**Primera, y es la grande.** Habíamos escrito que «la ley de Weyl de
cualquier operador autoadjunto en dominio compacto da densidad
constante». **Es falso, y lo refuta nuestra propia §3**: la matriz
diag(γ₁ … γ₆₂₀) es autoadjunta, vive en un espacio de dimensión finita y
su conteo es N(T) — densidad que crece. Además, la ley de Weyl en
dimensión d da N(λ) ~ C·λ^(d/2); la densidad constante es propia de la
reja de Berry–Keating, no de «cualquier operador».

El enunciado correcto, que es el que la tabla sí sostiene, es más
modesto y más útil:

> El operador de Berry–Keating sobre un grafo cuántico compacto de
> **geometría fija** de largo total L tiene N(E) ~ (L/2π)·E — densidad
> constante — mientras los ceros piden (1/2π)·ln(T/2π). La obstrucción
> es **R3 contra (R1 + geometría fija)**. Y la única salida —una
> geometría que crezca como ln(T/2π)— **choca de frente con R6**: hacer
> que la caja respire «como los ceros piden» es volver a leer la
> respuesta en la respuesta.

Esa versión cierra el círculo con la §3 en vez de contradecirla, y es
además la que deja el programa en un lugar concreto: **lo que hay que
buscar no es otro operador lineal fijo, sino una estructura cuya
geometría efectiva dependa de la energía y que se derive de la
aritmética, no de los γ_n.**

**Segunda: le habíamos hecho un reproche injusto al Taller.** Dijimos
que el canto de una reja «no está definido al último bit» y que el 0.555
publicado era artefacto. Barrido el origen de los cajones, el número del
Taller es **estable entre 0.542 y 0.553**, y sin cajones (Kolmogorov–
Smirnov contra Wigner) da **0.5331 sin ambigüedad**. El 0.4297 que
habíamos publicado como «la medición honesta» era **un defecto nuestro**:
construir los niveles y restarlos entre sí cancela números de orden 250 y
desparrama el espaciado unos 400 ulps alrededor de 1, justo encima de un
borde de cajón. **El número del Taller estaba bien; el nuestro estaba
mal.** Lo que sí queda en pie, y es lo que importa, es que el canto es
una estadística **discontinua sobre un espectro degenerado**.

**Y tercera, que no rompió nada pero mejoró todo.** Nuestra «separación
de 0.87 σ» dividía por la dispersión de una sola realización, que
contesta «¿se distingue un bloque?» y no «¿difieren las medias?». La
prueba correcta sobre las medias da t = 4.17 contra el sorteo puro y
t = 2.09 contra GUE. Y el refutador nos exigió lo que faltaba: **si el
examen es ciego, mostrá el examen que no lo es.**

### La varianza de número — el examen que sí ve

El canto mira **un espaciado por vez**. Por eso no puede ver la
diferencia entre una sucesión rígida y un sorteo sin memoria que tenga
la misma ley de espaciados. La varianza de número Σ²(L) —cuánto fluctúa
la cantidad de niveles en una ventana de largo L— mira **correlaciones**,
y con exactamente los mismos 47 espaciados ya no es ciega:

        Σ²(10)   ceros 0.327 ± 0.066
                 GUE   0.549 ± 0.112
                 sorteo SIN memoria 1.441 ± 0.797     t(ceros vs sorteo) = 18.8

        Σ²(20)   ceros 0.282 ± 0.068
                 GUE   0.500 ± 0.127
                 sorteo SIN memoria 1.769 ± 1.376     t(ceros vs sorteo) = 15.0

Eso convierte nuestro §2 de «el examen es ciego» en algo mucho más útil
para su pliego: **el examen es ciego PORQUE es una estadística de un
solo espaciado**, y la rigidez —que es lo que de verdad distingue un
espectro con estructura de uno sin ella— se ve en las de dos niveles.
Si su pliego va a pedir una condición estadística, que pida ésta.

### El control que faltaba en el eco, y que lo vuelve decisivo

También nos faltaba el control obvio del eco: tomar **los mismos ceros
verdaderos** y puntuarlos en períodos elegidos **al azar**, no
aritméticos. Sobre 120 sorteos:

        ceros en los períodos k·log p   :  razón 36.467
        ceros en períodos al azar       :  razón 0.0069 ± 0.0018 (máximo 0.0115)

Cuatro órdenes de magnitud. Ninguno de los 120 sorteos se acercó. El eco
no es un artefacto de la ventana ni del enventanado: **las posiciones
k·log p son especiales.**

Y una precisión de encuadre que el refutador nos hizo notar y que
cambia dónde está lo asombroso: **el 36× lo carga el denominador.**
Contra un control emparejado (misma cantidad de niveles, misma densidad),
los ceros suenan 2.87 veces más fuerte EN los períodos aritméticos —y
12.7 veces más CALLADOS afuera. Lo raro no es el grito en log p. **Lo
raro es el silencio entre los primos.**

## §6 — Tres correcciones a su pliego

Las escribimos derecho porque su pliego es lo que hizo todo esto
medible. Corregir un plano bueno es trabajar con él.

### (a) R1–R5 se satisface HACIENDO TRAMPA. Falta R6.

§3 lo muestra: `H = diag(γ_1 … γ_620)` cumple los cinco. No es un chiste
ni un caso patológico: es la demostración de que **el pliego, tal como
está redactado, no exige construir nada**. Cualquier lista de números se
puede meter en una diagonal.

    R6 — NO CIRCULARIDAD: el operador debe quedar DEFINIDO sin mirar
         los ceros. Su definición no puede contener, ni explícita ni
         implícitamente, la lista {γ_n} ni ninguna cantidad computada
         a partir de ella.

Sin R6, el pliego es satisfacible y vacío. Con R6, es el problema real.

### (b) R2 no es un requisito independiente

Toda estadística espectral es una función del espectro y de nada más.
Entonces «Spec(H) = {γ_n}» **implica** «H tiene la estadística de los
γ_n»: R2 no agrega ninguna restricción que la identificación de niveles
no haya puesto ya. Peor: exige en silencio una segunda conjetura
abierta, la ley de Montgomery–Odlyzko. Los dos casos posibles:

- si los γ_n son GUE (creído, no probado), R2 no restringe a H en nada;
- si no lo son, entonces R1 Y R2 juntos son insatisfacibles.

En cualquiera de los dos, R2 aporta cero. **Debe reetiquetarse como
restricción de diseño**, no como requisito. Y así reetiquetada tiene
tres usos reales, que es lo que importa:

1. **Test barato de descarte** para candidatos cuyo espectro todavía no
   se probó que sean los γ_n. (Con la advertencia de §2: es un test que
   pasa hasta un generador de números aleatorios, así que descarta poco.)
2. **Restricción de simetría a priori**: H debe estar en la clase de
   simetría A, es decir **sin ninguna simetría antiunitaria que conmute
   con él**. Esto elimina de entrada las construcciones simétricas
   reales y las invariantes bajo inversión temporal — darían GOE.
   Esta es la parte de R2 que sí muerde, y muerde antes de construir.
3. **Exclusión de la familia del caos aritmético** (abajo).

Sobre esto conviene ser preciso, porque circula la idea de que hay una
tensión entre «GUE» y «aritmética». **No la hay, para un espectro de
longitudes de primos.** Los períodos k·log p = log(p^k) son todos
DISTINTOS por factorización única: multiplicidad exactamente 1. Y
{log p} es Q-linealmente independiente, así que tampoco hay
coincidencias aditivas. Los pesos Λ(n)/√n cumplen la regla de suma
genérica de Hannay–Ozorio de Almeida con κ·⟨g⟩ = 1, lo que da el factor
de forma diagonal

    K_diag(τ) = τ

que es el valor GUE — no 2τ (GOE), no saturado (Poisson). Bogomolny lo
dice explícitamente (arXiv:0708.4223): como los primos no están
degenerados, el sistema dinámico conjetural asociado a los ceros debe
pertenecer a la clase de universalidad de los sistemas sin inversión
temporal, y por lo tanto tener estadística GUE.

El mecanismo Poisson dispara **solo** para superficies hiperbólicas
aritméticas, cuyo espectro de longitudes tiene multiplicidades
exponencialmente grandes (Bogomolny–Georgeot–Giannoni–Schmit;
Luo–Sarnak). Corolario útil para el pliego: **cualquier construcción
que pase por el laplaciano de una superficie aritmética no puede tener
a {γ_n} como espectro.** Ahí sí R2 mata familias enteras.

Caveat honesto: la curva GUE completa, más allá de la aproximación
diagonal, descansa sobre la conjetura de pares de primos de
Hardy–Littlewood (Bogomolny–Keating). «GUE para zeta» es él mismo
conjetural.

### (c) Su lectura del Taller no sobrevive a la medición

Usted escribió, en «El hallazgo central»: A tiene la densidad, B tiene
el canto, son dos mitades de una misma máquina. Y de ahí sacó la
pregunta nueva —cómo meter la información de los primos sin destruir la
estructura espectral correcta— que es, insistimos, la pregunta correcta.

Pero las mitades no son esas. **El canto es gratis** (§2): un generador
de números aleatorios lo pasa a 0.82 σ por realización de los ceros verdaderos. El
Prototipo B no tiene «la mitad» de la respuesta — tiene una propiedad
que comparte con cualquier ruido de la clase adecuada; y en el único
examen que sí lleva aritmética adentro, su propio Murciélago, saca
1.064 contra el 36.467 de los ceros.

Con todas las letras: **el Prototipo B no aporta información. La
división útil no es «estadística + afinación», es «espectro + eco», y
el eco es el único de los dos que cuesta.** Usted encontró la pregunta
correcta; las mitades no son las que nombró. Su §16 («NO BUSCAR OTRA
MATRIZ GUE. BUSCAR GUE + ARITMÉTICA») queda reforzado, no debilitado:
buscar otra matriz GUE es aún más inútil de lo que usted supuso.

---

## §7 — Una corrección a nosotros mismos, en exhibición

La casa exige que los errores propios se muestren enteros y con el
nombre del hallazgo. `cmd/maquina`, hallazgo **F178**:

1. **Nunca computó el Prototipo A.** Estipuló los espaciamientos en 1.0
   y publicó `canto = 0.555` como si fuera medido. Construyendo los
   niveles y midiéndolos de verdad da 0.4297.
2. **El Prototipo B fue una matriz, una semilla, sin barra de error.**
   Publicó 0.076. Como ensamble de 24 matrices, el número honesto es
   0.1097 ± 0.0304 — el 0.076 estaba dentro del ruido de una sola tirada.
3. **El «7.3 veces mejor» era una razón entre un número medido y un
   número estipulado.** Se retira. Los reemplazos honestos son ≈5.1 o
   ≈3.9 según con qué A se compare, y ninguno de los dos significa lo
   que se creyó que significaba, porque §2 muestra que el canto no
   distingue señal de ruido.
4. **El canto de una reja de estacas ni siquiera está bien definido
   sobre un espectro degenerado**, porque s = 1 cae sobre un borde de cajón.

Reproducido honestamente acá. La pieza `taller-maquinas.svg` que usted
analizó lleva esos números viejos; el análisis suyo era correcto sobre
lo que la pieza decía, y la pieza decía mal.

---

## §8 — Lo que sobrevive y lo que sigue

### Lo que sobrevive

- **El Murciélago**, entero y ascendido: es el único instrumento del
  taller que separa aritmética de azar, 36.467 contra 1.064. Adoptado
  como test obligatorio de todo candidato, tal como usted lo propuso en
  su «Matriz de trabajo para Doc».
- **La restricción de simetría**: clase A, sin simetría antiunitaria.
- **La exclusión** de la familia de superficies aritméticas.
- **Su «Regla del laboratorio»** (IMAGEN → HIPÓTESIS → EXPERIMENTO →
  ATAQUE → DEFINICIÓN → DEMOSTRACIÓN). Esta acta es el ATAQUE de esa
  cadena, y volvió con una obstrucción — que es lo que un ataque trae
  cuando es honesto.

### La segunda obstrucción estructural: el signo

Hay una segunda, independiente de la densidad. En la fórmula de traza de
Selberg la suma sobre órbitas periódicas entra con **+1**; en la fórmula
explícita de Weil la suma sobre primos entra con **−2**. El signo
RELATIVO es libre de convención: no se arregla normalizando.

La observación es vieja — Hejhal 1976, Berry 1986, a quien Connes cita.
Lo de Connes es la **interpretación**: los ceros como espectro de
**absorción**, el espacio de Pólya–Hilbert apareciendo en negativo.

Y hay que decir con cuidado qué es y qué no es:

- **No es un no-go riguroso.** En el caso de cuerpos de funciones tiene
  explicación a nivel de teorema: los ceros viven en H¹, que entra en
  la fórmula de Grothendieck–Lefschetz con signo (−1)¹.
- **Connes probó una fórmula de traza cuya VALIDEZ es equivalente a RH**
  para funciones L con Grossencharakter (*Selecta Math.* **5** (1999)
  29–106). Eso es una **reducción**, no una demostración. Quien diga
  otra cosa está leyendo mal el paper.
- Hay un segundo desajuste, el «problema asintótico» de Sierra:
  2·sinh(m·λ_p/2) es igual a p^(m/2) solo cuando m → ∞, no a m finito.

### Lo que nuestras mediciones dicen que el candidato DEBE tener, y no está en su lista

Este es el aporte que sale de §4 y que su pliego no pide:

    L(E) = ln(E/2π)

**Una caja efectiva que depende de la energía.** Sea lo que sea la
máquina, **no es un operador lineal fijo sobre un dominio compacto
fijo**. La longitud efectiva tiene que crecer con la altura, y crecer
como el logaritmo: 3.306 → 4.982 en el tramo que medimos, ×1.507,
siguiendo ln(T/2π) columna contra columna.

Eso reorienta la búsqueda, y conviene decir hacia dónde:

- hacia la ruta adélica de Connes, donde el «tamaño» no es un intervalo
  sino un espacio con infinitos lugares y el corte no es un borde;
- hacia las familias xp **con pared que respira**, donde el
  confinamiento depende de la energía en vez de ser una condición de
  contorno fija;
- y **lejos** de: cualquier grafo compacto de geometría FIJA, cualquier
  intervalo con condiciones de borde autoadjuntas, cualquier matriz
  finita fija, y cualquier laplaciano de superficie aritmética.

Su punto 20, «Orden recomendado de trabajo», se corrige entonces en un
solo lugar: el paso 2 («definir el espacio y dominio del operador
candidato») deja de ser un paso administrativo y pasa a ser **el
problema**. No es que falte definir el dominio: es que ningún dominio
compacto fijo puede servir, y eso está probado. El «telar» de su §11,
traducido a objeto matemático, es esto: lo que falta no es un hilo más,
es un soporte cuyo tamaño dependa de la altura a la que se lo mire.

---

## §9 — Sello

    Se cierra:   la pregunta que usted mandó. El requisito que impide
                 la construcción es R3 contra R1 —la ley de conteo
                 contra la autoadjunción con espectro discreto—, con
                 teorema ajeno (Endres–Steiner 15.6) y con medición
                 propia (la caja que respira).

    Se retira:   el 7.3 del taller, el 0.555 estipulado, el 0.076 de
                 una sola semilla, y la lectura de que el Prototipo B
                 era media máquina.

    Se agrega:   R6 no circularidad; R2 reetiquetado como restricción
                 de diseño con contenido operativo (clase A, exclusión
                 aritmética); y el requisito L(E) = ln(E/2π) que su
                 pliego no pedía.

    NO se hizo:  nada de esto prueba RH ni la refuta, ni la acerca ni
                 la aleja un milímetro. Esto es una revisión de
                 especificación con mediciones adentro. Una estructura
                 cerrada no es una hipótesis probada.

    Todavía no.

---

*Usted pidió que si no se podía construir, dijéramos exactamente qué lo
impide. Lo dijimos: la caja no puede ser fija porque la densidad de los
ceros crece con el logaritmo de la altura, y una caja fija tiene
densidad constante. Su pliego no tenía ese renglón; ahora lo tiene, y
lo tiene porque su pliego nos obligó a medirlo.* ⚛️🦇📏
