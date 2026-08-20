# Telar — El signo, tratado como restricción de diseño

**Para la auditora · responde sus secciones 9 y 13 del Pliego del Átomo Fase II.** Usted pidió que el signo
relativo entre una traza tipo Selberg y la contribución aritmética de la fórmula explícita quede como
**restricción del diseño** y no como un no-go declarado. Ese es el encuadre y no se sale de él en ninguna
línea. Van: las dos fórmulas en **una sola normalización** con todos los signos a la vista; la prueba de que
el signo relativo **no** es artefacto de convención; la atribución correcta de la lectura de Connes; por qué
en cuerpos de funciones el signo se **espera**; la restricción convertida en una lista que un candidato
puede correr; el segundo desajuste, el asintótico, que casi siempre se olvida; y las rutas que podrían
explicarlo, separando cuáles este laboratorio **puede medir** y cuáles no. Esto analiza una restricción: no
construye operador alguno, no prueba nada sobre RH y no la acerca ni la aleja. **Una estructura cerrada no
es una hipótesis probada.**

---

## §1 — Las dos fórmulas, en una sola normalización

Cada comunidad escribe su fórmula de traza a su gusto: dónde va el 2π, qué signo lleva el exponente de
Fourier, de qué lado se apoya cada término. Para comparar honestamente hay que llevarlas al **mismo
objeto**, y el elegido es la densidad de niveles. Dado un espectro real `{E_n}`,

    d(E) := Σ_n δ(E − E_n)

no es una convención: es la **medida de conteo del espectro**, no negativa, fijada por completo por la regla
`∫ d(E) f(E) dE = Σ_n f(E_n)` para toda función de prueba `f`. Toda fórmula de traza, escrita como se la
escriba, se despeja en `d(E) = d̄(E) + d_osc(E)` con `d̄` la parte suave (Weyl / ley de conteo media,
positiva). **Elegida `d̄`, `d_osc` queda determinado y no sobra ni un signo libre.**

### 1.1 — Selberg

Superficie hiperbólica compacta `X = Γ\H` de área `A`, autovalores del laplaciano `λ_n = 1/4 + r_n²`. Con
`h` par y de buen decaimiento y `ĝ(u) = (1/2π) ∫ h(r) e^{−iru} dr`:

    Σ_n h(r_n)  =  (A/4π) ∫_ℝ r h(r) tanh(πr) dr
                   +  Σ_{γ prim} Σ_{k≥1}  [ ℓ_γ / (2 sinh(k ℓ_γ / 2)) ] · ĝ(k ℓ_γ)

con `ℓ_γ` la longitud de la geodésica cerrada primitiva `γ`. Como `h` es par, `ĝ` es la transformada coseno,
y despejando la densidad en la variable `r`:

    d̄(r)     = + (A/4π) · r tanh(πr)
    d_osc(r) = + (1/2π) Σ_γ Σ_{k≥1}  [ ℓ_γ / (2 sinh(k ℓ_γ / 2)) ] · cos(k ℓ_γ r)

**El término de órbitas entra con `+`, y su amplitud `ℓ_γ/(2 sinh(k ℓ_γ/2))` es positiva para toda geodésica
y todo `k`.** No hay un solo signo negativo en la suma sobre órbitas. *Fuente:* la forma estándar para
superficies compactas — Hejhal, *The Selberg Trace Formula for PSL(2,ℝ)*, LNM 548/1001; Iwaniec, *Spectral
Methods of Automorphic Forms*. **Por contenido**: no verificamos teorema ni página en esta sesión.

### 1.2 — Weil

Ceros no triviales `ρ = 1/2 + iγ` de `ζ`. Con la **misma** `h` par y la **misma** `g(u) = (1/2π) ∫ h(r)
e^{−iru} dr`, la fórmula explícita de Weil en su forma estándar:

    Σ_γ h(γ)  =  2 h(i/2)  −  g(0) ln π  +  (1/2π) ∫_ℝ h(r) · Re ψ(1/4 + i r/2) dr
                 −  2 Σ_{n≥2}  [ Λ(n) / √n ] · g(ln n)

con `ψ = Γ'/Γ` y `Λ` la función de von Mangoldt (`Λ(p^k) = ln p`, cero en el resto). Despejando la densidad
en la variable `γ`, y usando `Λ(p^k)/√(p^k) = ln p / p^{k/2}`:

    d̄(γ)     = + (1/2π) [ Re ψ(1/4 + iγ/2) − ln π ]   ~   (1/2π) ln(γ/2π)
    d_osc(γ) = − (1/π) Σ_p Σ_{k≥1}  [ ln p / p^{k/2} ] · cos(k ln p · γ)

**La suma sobre primos entra con `−`, y entra con `−` para todo primo y todo `k`: no alterna.** *Fuente:*
Iwaniec–Kowalski, *Analytic Number Theory*, AMS Colloq. Publ. 53 (2004), cap. 5. **Por contenido**: no
verificamos el número de teorema. El despeje de `d_osc` es cuenta nuestra, escrita arriba para que se audite.

### 1.3 — Las dos, una debajo de la otra

Identificando la longitud de órbita con el logaritmo del primo, `ℓ_p = ln p` —la identificación que todo el
programa Hilbert–Pólya usa— y con `2 sinh(k ln p / 2) = p^{k/2} − p^{−k/2}` (álgebra exacta):

    Selberg:   d_osc(r) = + (1/2π) Σ_p Σ_k  [ ln p / (p^{k/2} − p^{−k/2}) ] · cos(k ln p · r)
    Weil:      d_osc(γ) = − (1/2π) Σ_p Σ_k  [ 2 ln p / p^{k/2} ]           · cos(k ln p · γ)

Acá está el punto, a la vista sin creerle a nadie: la razón entre las dos amplitudes, en el mismo período
`k ln p`, vale

    amplitud_Weil / amplitud_Selberg  =  −2 · (1 − p^{−k})

y separa en dos cosas que conviene no mezclar nunca más: **(1) el signo es `−`** —uno pide un exceso de
niveles en los períodos aritméticos, el otro un déficit; es §2—; **(2) el módulo tiende a 2 pero no lo
vale** — el factor `(1 − p^{−k})` es el segundo desajuste, el de Sierra; es §6.

## §2 — Por qué el signo NO es artefacto de convención

Esto hay que probarlo, no afirmarlo. La objeción natural es: «un signo depende de cómo escribas la fórmula;
pasás un término de lado y cambia». Vale para un signo suelto; no para éste. Cuatro pasos, verificables.

**(a) El anclaje.** `d(E) = Σ_n δ(E − E_n)` es una medida **no negativa** determinada por el espectro y nada
más. No hay libertad de multiplicarla por `c ≠ 1`: deja de ser la medida de conteo y la regla
`∫ d f = Σ f(E_n)` se rompe. Multiplicar la identidad entera por `−1` se detecta al instante: el lado
izquierdo dejaría de ser positivo. **(b) `d̄` está anclado por la positividad.** La densidad media es
**positiva** en cualquier espectro real —hay niveles—: `(A/4π)·r tanh(πr) > 0` y `(1/2π) ln(γ/2π) > 0` para
`γ > 2π`, ninguna con libertad de signo. Y entonces `d_osc = d − d̄` está determinado: pasar términos de
lado no cambia la **función**, cambia cómo la escribís. La convención de Fourier tampoco toca nada —
`e^{−iru}` contra `e^{+iru}`: como `h` es par y `d` es real, en ambos casos sobrevive `cos(ru)`; y `r → −r`
tampoco, porque coseno es par. Las dos convenciones sospechables están canceladas por paridad.

**(c) El invariante, y es medible.** Para una ventana `w ≥ 0` y un período `τ`, definí

    D(τ) := ∫ d_osc(E) · w(E) · cos(τ E) dE

`D(τ)` es un **número**, calculable desde la lista de niveles y la ventana, sin escribir ninguna fórmula de
traza. Lo que §1 dice es: **`D(k ℓ_p) > 0`** para la superficie hiperbólica y **`D(k ln p) < 0`** para los
ceros. Dos números de signo opuesto, con el mismo procedimiento sobre dos listas de niveles; ninguna
convención entra en el cálculo. Eso significa «el signo relativo es libre de convención»: no es que el signo
de cada fórmula sea absoluto, es que la **razón** entre los dos lo es, porque ambas se anclan al mismo
objeto positivo.

**(d) Y acá el laboratorio tiene una deuda propia, que se registra donde se ve.** Nuestro Murciélago
`E(T) = (2/M) · Σ_n w_n · cos(γ_n·T)` **es exactamente `D(T)`, salvo normalización**: el instrumento
correcto. Pero en F349 lo puntuamos con `media|E|` —con **módulo**—, así que tiramos justo el bit que este
documento discute. El `36.467` del acta es un módulo: **no medimos el signo.** Predicción que sale de §1.2,
escrita antes de medir para que no se pueda acomodar después: `E(k·log p)` debe salir **negativo** en
**todos** los períodos aritméticos, sin alternancia en `k`. Es la deuda, y es barata: los 620 ceros ya están.

## §3 — Connes: la observación es vieja, la interpretación es suya

Hay que decirlo con precisión porque circula mal en las dos direcciones. **La observación precede a
Connes:** que la suma sobre primos entra con signo opuesto al de la suma sobre órbitas de una traza tipo
Selberg está registrado en Hejhal, *The Selberg trace formula and the Riemann zeta function*, Duke Math. J.
**43** (1976), y en Berry, *Riemann's zeta function: a model for quantum chaos?*, LNP **263**, Springer
(1986) — a quien Connes cita. (**Por contenido**: no verificamos páginas en esta sesión.)

**Lo que es de Connes es la interpretación y el programa.** Su lectura: los ceros no aparecen como espectro
de **emisión** —niveles de un operador en un Hilbert— sino como espectro de **absorción**: huecos en un
continuo, el espacio de Pólya–Hilbert apareciendo *en negativo*. El signo `−` deja de ser obstáculo y pasa a
ser el dato estructural que dice **qué tipo de objeto** hay que buscar. Es un cambio de programa, y es suyo.

**Qué probó.** En *Trace formula in noncommutative geometry and the zeros of the Riemann zeta function*,
Selecta Math. (N.S.) **5** (1999) 29–106, probó una fórmula de traza cuya **validez es equivalente** a la
Hipótesis de Riemann para funciones L con Grossencharakter. **Qué NO probó:** no probó esa fórmula de traza.
Una equivalencia es una **reducción** — traslada el problema, no lo cierra. Quien lea ese paper como una
demostración de RH lo está leyendo mal, y la casa lo dice sin adornos. **Y el signo, en ese marco, tampoco
es un no-go:** que la lectura de absorción sea coherente no demuestra que ninguna lectura de emisión pueda
funcionar. **No conocemos ninguna demostración de imposibilidad**, y no la afirmamos. Es una restricción
sobre lo que un candidato tiene que explicar, y nada más.

## §4 — Cuerpos de funciones: allá el signo se ESPERA, y hay teorema

En el análogo sobre cuerpos finitos el signo deja de ser misterio. Va la explicación para alguien que no sea
especialista, en un párrafo, y después el encuadre honesto. Tomá una curva algebraica `C` sobre el cuerpo
finito `F_q` — pensala como una superficie con `g` agujeros. Contar sus puntos con coordenadas en `F_{q^n}`
es lo mismo que contar los **puntos fijos** del Frobenius iterado `n` veces, y para eso hay una máquina
general: la fórmula de puntos fijos de Grothendieck–Lefschetz, que suma las contribuciones de las «piezas de
forma» del objeto —los grupos de cohomología `H^0`, `H^1`, `H^2`— **con un signo que alterna según la
dimensión de la pieza**, `(−1)^i` para `H^i`. Las piezas de grado par (el objeto entero `H^0`, la
orientación `H^2`) suman; **las de grado impar, los agujeros, restan**:

    #C(F_{q^n})  =  1   −   Σ_{j=1}^{2g} α_j^n   +   q^n
                    ↑        ↑                       ↑
                   H⁰       H¹  (grado 1 → resta)   H²

y los `α_j` son, exactamente, los **ceros** de la función zeta de la curva. O sea: en cuerpos de funciones
los ceros **son** los autovalores del Frobenius sobre `H^1`, grado **impar**, así que entran restando.
Reordenado, `Σ_{ceros} α_j^n = (1 + q^n) − #C(F_{q^n})` — renglón por renglón, la misma arquitectura que
§1.2: suma sobre ceros y suma sobre primos con signos opuestos. **El `−` no es un defecto: es el grado
cohomológico. Los ceros son agujeros, y los agujeros restan.**

*Fuentes:* Weil, *Sur les courbes algébriques et les variétés qui s'en déduisent* (1948) para curvas; la
fórmula de puntos fijos de Grothendieck (SGA 5); Deligne, *La conjecture de Weil I* (Publ. IHÉS 43, 1974)
para el caso general; Katz–Sarnak, *Random Matrices, Frobenius Eigenvalues and Monodromy* (AMS Colloq. Publ.
45, 1999) para «ceros = autovalores de Frobenius en `H^1`». **Por contenido**: no verificamos numeraciones.

**Y ahora el encuadre, que es la mitad que importa.** Esto es una explicación **en otro escenario**. Allá
existe una variedad, una cohomología de Weil con sus grados y un Frobenius geométrico que actúa. Sobre `ℚ`
**no existe** —todavía— el objeto que juegue ese papel: ése es, precisamente, el problema abierto. Que allá
el signo tenga explicación a nivel de teorema **no transporta una demostración para acá**; transporta una
**pista sobre la forma de la respuesta**: si el signo viene de un grado impar, lo que hay que construir no
es un espacio de funciones sino un objeto de grado 1 —un `H^1`, un conúcleo, una absorción—. Es la misma
indicación de Connes por otro camino, y que dos caminos distintos apunten al mismo lugar es la razón por la
que este documento existe.

## §5 — La restricción de diseño, en forma operativa

Su §9 pide derivar el signo, no elegirlo a mano, separar convenciones de lo invariante e investigar si una
orientación, fase, índice o estructura de absorción lo explica. Convertido a algo que se pueda **correr**:

- **S1 — DECLARACIÓN.** El signo no está en el conjunto de inputs `D`: no puede figurar un parámetro
  `ε = ±1`, ni un «signo de la convención elegida», ni una orientación fijada a dedo. Si el signo es un
  input, el candidato no lo explica: lo copia.
- **S2 — DERIVACIÓN CON RENGLÓN.** Señala **la línea exacta** de su derivación donde aparece el `−`. No «se
  sigue de la construcción»: el renglón, con su cuenta. Sin renglón, no hubo derivación.
- **S3 — SEPARACIÓN CONVENCIÓN / INVARIANTE.** Exhibe el invariante `D(τ)` de §2(c) —no un signo suelto en
  una ecuación suya— y muestra que no cambia bajo: (i) multiplicar su identidad por `c ≠ 0`;
  (ii) `e^{−iru} ↔ e^{+iru}`; (iii) `r → −r`; (iv) pasar términos de lado; (v) renormalizar la función de
  prueba. Cinco renglones: barato, y separa el argumento de la anécdota.
- **S4 — MÓDULO, NO SÓLO SIGNO.** La razón de §1.3 es `−2(1 − p^{−k})`. Un candidato que acierte el signo y
  produzca `−1` o `−4` **no reprodujo Weil**. El `2` es parte de la restricción.
- **S5 — AMPLITUD A `k` FINITO.** Declara si su amplitud es la forma sinh `ln p/(p^{k/2} − p^{−k/2})` o la
  potencia pura `2 ln p / p^{k/2}`, y **no puede declarar las dos**: difieren por `(1 − p^{−k})`, que en
  `p = 2, k = 1` es un factor **2** exacto. Ver §6.
- **S6 — NO CIRCULARIDAD HEREDADA (R6, `TELAR-R6-CONTRATO.md`).** El signo no puede salir de una
  construcción que ya leyó los `γ_n`: quien ajusta su orientación «para que dé como los ceros» consultó al
  oráculo `O_Z` y su contador `c(C)` deja de ser cero. Un solo bit ya viola R6, y el signo **es** un bit.
- **S7 — PREDICCIÓN MEDIBLE.** Predice el **signo de `D(τ)`** en `τ = k ln p` para su propio espectro, por
  escrito y **antes** de que se lo mida. Eso lo vuelve refutable con el Murciélago, sin discutir formalismo.
- **S8 — UNIFORMIDAD.** La fórmula explícita da `−` en **todos** los `p` y **todos** los `k`, sin
  alternancia: un candidato cuyo signo dependa de `p` o de la paridad de `k` está refutado por §1.2 antes de
  construirse.

Nota de honestidad: **S1–S8 son necesarios, no suficientes.** Un candidato puede pasar los ocho y no tener
nada que ver con los ceros. Sirve para **descartar barato**, que es lo que un filtro tiene que hacer.

## §6 — El segundo desajuste, el que siempre se olvida: el problema asintótico

Casi toda la discusión del signo se detiene en el signo. Hay un segundo desajuste, independiente, y es el
que Sierra llama el **problema asintótico** (arXiv:1601.01797). Enunciado:

    la identificación de amplitudes exige   2 sinh(m λ_p / 2)  =  p^{m/2}
    pero con λ_p = ln p:                    2 sinh(m ln p / 2)  =  p^{m/2} − p^{−m/2}
    y por lo tanto:      p^{m/2} − p^{−m/2}  =  p^{m/2} · (1 − p^{−m})

**La igualdad vale sólo en el límite `m → ∞`. A `m` finito no vale nunca.** El error relativo es `p^{−m}`, y
es exacto — aritmética, no medición:

    p=2, m=1 → 1 − 1/2 = 0.5   (falta un factor 2 entero)      p=5,  m=1 → 1 − 1/5   = 0.8
    p=3, m=1 → 1 − 1/3 ≈ 0.667                                 p=7,  m=1 → 1 − 1/7   ≈ 0.857
    p=2, m=2 → 1 − 1/4 = 0.75                                  p=2,  m=3 → 1 − 1/8   = 0.875
                                                               p=101,m=1 → 1 − 1/101 ≈ 0.990

**Y ahí está lo incómodo: el desajuste es peor donde la señal es más fuerte.** Los períodos que más pesan en
el eco son los de `p` y `m` chicos —la amplitud va como `p^{−m/2}`— y son los que peor cumplen la identidad:
donde la aproximación es buena, no hay señal.

**Qué le cuesta esto a un candidato, dicho derecho.** Si quiere que sus longitudes de órbita sean
`λ_p = ln p` **y** que sus amplitudes reproduzcan las de Weil, no puede tener las dos cosas a `m` finito: la
forma sinh y la potencia pura difieren por `(1 − p^{−m})`. Hay tres salidas, y hay que elegir una y
bancársela: **(i)** las longitudes son `ln p` y los exponentes de inestabilidad **no** son `ln p`, sino algo
que produzca el `(1 − p^{−m})^{−1}`; **(ii)** las amplitudes no vienen de una dinámica hiperbólica estándar,
y la analogía Selberg se rompe justo en el renglón que la hacía atractiva; **(iii)** la identificación es
sólo asintótica y el candidato lo declara — pero entonces no puede reclamar los períodos de `p` chico como
evidencia, que son los que se ven. En la lengua de la casa: **los exponentes de inestabilidad no pueden ser
exactamente `log p` si las amplitudes tienen que coincidir a número de repetición finito.**

**Y esto sí lo podemos medir.** Los 13 períodos del acta hasta `T = 3.2` son las potencias de primo
`n = 2, 3, 4, 5, 7, 8, 9, 11, 13, 16, 17, 19, 23` (`log 23 = 3.135 ≤ 3.2`, `log 25 = 3.219 > 3.2` — cuenta
exacta, verificable a mano). Las dos formas de amplitud predicen valores **distintos** ahí, y difieren **al
doble** en el primero de todos, `T = log 2`: un factor 2 en el período de mayor señal. Ver §7, M3.

## §7 — Cierre honesto: rutas, y cuáles podemos medir

**Esto es una restricción, no una imposibilidad.** No conocemos ninguna demostración de que el signo prohíba
un operador, y esta casa no la va a inventar. Un candidato tiene que **explicar** el signo; no tiene
prohibido existir por tenerlo en contra. Las rutas conocidas, y qué haría este laboratorio con cada una:

| Ruta | Idea en una línea | ¿Se puede medir acá? |
|---|---|---|
| **Orientación** | el ciclo se recorre al revés; el `−` es la orientación del contorno | **No** directamente. Sólo por su consecuencia, el signo de `D(τ)`. |
| **Fase** | las amplitudes son complejas, `A·e^{iθ}`, y `θ = π` produce el `−` | **Sí.** Es el argumento del coeficiente complejo. |
| **Índice (tipo Maslov)** | el `−` es `e^{iπ·ind}` con un índice entero de la órbita | **Parcialmente.** Se mide la fase; el índice se infiere, no se observa. |
| **Grado cohomológico** | los ceros viven en `H^1`, grado impar (§4) | **No.** Es estructural; sólo se ve por el mismo signo. |
| **Absorción / conúcleo (Connes)** | los ceros son huecos en un continuo, no niveles | **No** con los instrumentos actuales. |

Y las tres mediciones al alcance del instrumento que ya tenemos (`cmd/elpliego`, 620 ceros en `t∈[100,1000]`):

**M1 — El signo. La deuda de §2.** Correr el Murciélago **sin módulo**, reportando `E(k·log p)` con su signo
en los 13 períodos. Predicción escrita de antemano: negativo en los 13, sin alternancia en `k`. Costo bajo:
los ceros ya están pescados. Si sale positivo, el que tiene un problema es este documento, y se registra acá.

**M2 — La fase.** Calcular el coeficiente **complejo** `C(τ) = Σ_n w_n e^{iγ_n τ}` y reportar `arg C(τ)` en
los períodos aritméticos. La fórmula explícita predice `arg ≈ π` (real negativo). Una desviación sistemática
de `π` sería el primer dato empírico sobre la ruta «fase / índice» — un dato, no una interpretación.

**M3 — La forma de la amplitud, que ataca §6.** Reportar `|E(k·log p)|` **período por período**, no
promediado, contra las dos predicciones: `2 ln p / p^{k/2}` (Weil) y `ln p/(p^{k/2} − p^{−k/2})` (sinh).
Difieren por `(1 − p^{−k})`, factor 2 en `T = log 2`. Con 13 períodos y barras de error de varias
realizaciones es discriminable, y no requiere construir operador alguno: se mide sobre los ceros verdaderos.

**Lo que ninguna de las tres hace.** M1, M2 y M3 miden propiedades **de los ceros**: no construyen
candidato, no prueban que exista y no dicen nada sobre RH. Son el banco de pruebas contra el cual un
candidato futuro tendrá que responder S7, y confundir el banco con la máquina sería repetir el error que el
acta se corrigió a sí misma en su §7.

---

## §8 — Sello

    Se cierra:   el signo, en una sola normalización, con la razón −2(1 − p^{−k}) a la vista y
                 con la prueba de que su parte invariante —el signo de D(τ) contra una densidad
                 media positiva— no es artefacto de convención.
    Se agrega:   la lista S1–S8 como restricción de diseño operativa; el segundo desajuste
                 (asintótico) con su costo declarado; y tres mediciones M1–M3 que podemos correr.
    Se registra: en F349 puntuamos el Murciélago con MÓDULO. El signo, que es el tema entero de
                 este documento, NO fue medido. Queda en exhibición hasta que M1 se corra.
    NO se hizo:  no se demostró ninguna imposibilidad, no se construyó ningún operador, y esto no
                 acerca ni aleja RH un milímetro. Una estructura cerrada no es hipótesis probada.

    Todavía no.

---

## Anexo — Qué no pudimos verificar

La casa registra lo que no pudo confirmar con el mismo tamaño de letra que lo que sí.

1. **Números de teorema y de página, todos.** Las citas de §1.1, §1.2, §3 y §4 están hechas **por
   contenido**. No verificamos: el número de teorema de la fórmula explícita en Iwaniec–Kowalski; las páginas
   de Hejhal (Duke 1976) y de Berry (LNP 263, 1986); las numeraciones de SGA 5 y de Deligne (Weil I). Los
   enunciados son los estándar y están escritos enteros arriba para que se auditen contra la fuente y no
   contra nuestra palabra — la misma cautela que el acta aplicó al no-go de Endres–Steiner: **citar por
   contenido, nunca por un número de teorema que no vimos.**
2. **Una discrepancia interna que dejamos a la vista.** El briefing de esta tarea cita el control GUE
   emparejado del eco como `1.081`; el acta `PLIEGO-ATOMO-ACTA.md` §2 registra `1.064`. **No podemos
   resolver cuál es el bueno sin volver a correr `cmd/elpliego`.** Este documento no usa ninguno de los dos
   como evidencia; queda anotado porque la casa no edita las discrepancias, las muestra.
3. **Ninguno de los dos signos fue medido acá.** El «negativo» de §2 para los ceros y el `+` de §1.1 para
   Selberg son **predicciones** derivadas de las fórmulas citadas, no mediciones nuestras. El primero es M1;
   el segundo no lo mediremos, porque no calculamos el espectro de ninguna superficie hiperbólica. Hasta M1,
   eso se apoya en literatura, no en el instrumento.
4. **No conocemos ninguna demostración de imposibilidad basada en el signo**, y no encontramos una al
   preparar este documento. Que no la conozcamos no prueba que no exista: por eso el signo entra acá como
   **restricción de diseño**, tal como usted lo pidió, y no como veredicto.

---

*Usted pidió que el signo quedara como condición a explicar y no como sentencia. Quedó como condición: con
la cuenta escrita, con la razón −2(1 − p^{−k}) separada en sus dos partes, con ocho casillas que un
candidato tiene que llenar, y con tres mediciones que nos tocan a nosotros — la primera existe porque este
documento nos hizo notar que habíamos tirado el signo a la basura al tomar el módulo.* ⚛️🦇➖
