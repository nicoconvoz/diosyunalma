# Telar — Contrato R6 y conjunto D admisible

**Para la auditora · responde sus secciones 2, 11 y 13 del Pliego del Átomo Fase II.** Pidió tres
cosas: la definición formal de R6 y de D, la tabla de dependencias que todo candidato entrega antes
de ser probado, y las reglas de parada convertidas en algo que un programa pueda correr. Están las
tres, más cuatro candidatos auditados de verdad. Todo número es del acta `PLIEGO-ATOMO-ACTA.md`
(`cmd/elpliego`, hallazgo F349) o lleva cita con nombre. Encuadre, para que no se lea de más: esto
es un contrato de auditoría, no construye operador alguno y no mueve la hipótesis un milímetro.

## §1 — R6, con cuantificadores

**1.1 El objeto que se audita.** Sea `Z = {γ_1, γ_2, …}` el conjunto de ordenadas de los ceros no
triviales y `D` el conjunto de datos admisibles de §2. Un **candidato** no es un operador: es una
cuaterna declarada `C = (F, φ, θ, H)` — `F` la familia declarada (el espacio de objetos donde vive
H), `φ` el procedimiento de construcción (un algoritmo finito, escrito y ejecutable), `θ` la lista
**completa** de parámetros libres con su valor *y la historia de cómo se los fijó*, y
`H = φ(D_C ; θ)` con `D_C ⊆ D`. Un operador entregado sin `φ` y sin `θ` no es un candidato: es un
número. No se audita.

**1.2 R6 en el modelo de oráculos.** El constructor trabaja con dos oráculos. `O_D` contesta sobre
datos admisibles («¿cuál es el primo mil?», «¿cuánto vale Γ(1/4 + it/2)?»). `O_Z` contesta cualquier
pregunta cuya respuesta dependa de dónde están los ceros: «¿cuánto vale γ_100?», «¿cuánto da Σ²(10)
sobre los ceros?», y también **«¿este candidato se parece más a los ceros que aquel otro?»** — esa
última es la que más se cuela.

> **R6 — NO CIRCULARIDAD.** `C` satisface R6 si y sólo si existe una transcripción completa de `φ`
> —incluida la historia de cada `θ_i`— en la que
>
>     c(C) := #{consultas a O_Z durante la construcción de H} = 0
>
> Con cuantificadores:
>
>     (i)   ∀ hoja s del árbol de derivación de H:  s ∈ D
>     (ii)  ∀ θ_i ∈ θ:  ∃ derivación finita π_i con todas sus hojas en D tal que π_i ⊢ (θ_i = v_i)
>     (iii) ∀ bifurcación b en la historia de φ —elección de familia, de miembro, de peso, de
>           corte, de semilla, de norma—: la condición de b no depende de una respuesta de O_Z
>
> `c(C)` se reporta siempre, aunque valga cero. Un solo bit ya viola R6.

**1.3 Por qué se define sobre la derivación y no sobre la información.** Acá está el punto fino, y
conviene decirlo antes de que alguien lo use en contra: **informacionalmente nada es no circular.**
Los primos determinan a ζ, ζ determina sus ceros; todo objeto construido desde los primos es, en ese
sentido, una función de `Z`. Si «circular» significara «la salida es función de la entrada», el
contrato prohibiría también la respuesta correcta. Por eso R6 no habla de información sino de
**procedimiento**: la pregunta auditable no es «¿la salida depende de los ceros?» sino **«¿hubo que
saber dónde están los ceros para escribir la definición?»**. Eso sí es decidible, dado que el
candidato entregue transcripción honesta — y por eso el contrato entero se apoya en §3.

**1.4 Las seis formas de violarlo.**

    V1  EXPLÍCITA   la lista {γ_n} aparece literalmente en φ.
    V2  DERIVADA    aparece f({γ_n}): la brecha media medida, las cajas L = 2π/brecha de §4 del
                    acta (3.3057 … 4.9825), el Σ²(10) = 0.327 ± 0.066, la razón 36.467. Son
                    mediciones NUESTRAS sobre los ceros: ninguna entra en una definición.
    V3  AJUSTE      θ_i elegido minimizando discrepancia contra datos de ceros. Cuenta aunque sea
                    a ojo, de un solo parámetro, y aunque el resultado sea redondo.
    V4  LAVADA      θ_i no es literalmente función de Z pero está en biyección con una: «el valor
                    que hace que el conteo dé 620» es V4, porque 620 es nuestro conteo medido.
    V5  SELECCIÓN   la familia F es admisible pero el MIEMBRO se eligió mirando los ceros. Elegir
                    es definir. Es la que más se escapa: el candidato final no muestra cicatriz.
    V6  FUGA        se ajustó tras ver el resultado y se reportó lo ajustado como si fuera la
                    definición original. Se previene con una sola cosa: congelar antes (§3).

**1.5 Construcción contra evaluación.** Su §2 pide separarlos. **De construcción**: entran en la
definición de `H`, cumplen R6 entero, `c` cuenta sobre ellos. **De evaluación**: entran sólo en la
medición —ancho de ventana del Murciélago, cajones del canto, la L de Σ²(L), tramos de altura,
cantidad de realizaciones— y **sí** pueden mirar los ceros, porque los ceros son el objeto de
comparación, no la máquina. Con dos condiciones que no se negocian: se congelan **antes** de ver la
salida y son **idénticos** para candidato y controles. El acta enseñó por qué: mover el origen de
los cajones movía el canto de la misma reja entre 0.542 y 0.553. Un parámetro de evaluación elegido
después es un resultado inventado.

**1.6 Los casos sutiles que usted marcó.** *Una constante que «casualmente» coincide:* la regla es
**derivabilidad, no intención**. Si alguien pone una pared en 4.300 y no exhibe la derivación desde
`D`, no importa de dónde jure haberla sacado — 4.300 es exactamente la mejor caja fija de nuestra
§4, un número computado DESDE los ceros. Coincidir con una cantidad Z-derivada **invierte la carga
de la prueba**: queda `PENDIENTE-DERIVACIÓN` hasta que aparezca `π_i`. No es «circular probado»; es
«no admitido». *Un parámetro ajustado a los ceros* es V3; *ζ(s)* y *N(T)* se deciden en §2, ítems 3
y 4, por canal y por rol.

**1.7 El segundo eje: DERIVADO contra ESTIPULADO.** R6 sola no alcanza: un candidato puede tener
`c(C) = 0` y no explicar nada, porque metió la respuesta a mano desde una fórmula que resulta ser
Z-libre. Por eso todo veredicto lleva dos ejes — Eje 1, circularidad: `LIMPIO / CIRCULAR /
PENDIENTE-DERIVACIÓN`; Eje 2, contenido: `DERIVADO / ESTIPULADO`. `ESTIPULADO` quiere decir que la
propiedad que se celebra fue puesta como hipótesis, no obtenida como consecuencia. No es trampa; es
vacío. Su §6 pide exactamente esto al preguntar si `L(E) ≈ ln(E/2π)` puede derivarse «y no
simplemente ajustarse a los ceros».

## §2 — El conjunto D, ítem por ítem, con la razón

**1. Primos y funciones aritméticas — ADMITIDO.** `p_n`, `Λ(n)`, `μ`, `φ`, `d(n)`, `ψ(x)` y sus
sumas finitas: se computan por criba desde la divisibilidad, `c = 0` exacto. Es donde vive su Ruta
C: si `L(E)` va a emerger sin leer la respuesta, sale de acá.

**2. π, e, γ_Euler, Γ(s), log, funciones especiales — ADMITIDO**, porque cada una tiene definición propia y evaluación a precisión arbitraria sin tocar `Z`.

**3. ζ(s) y ξ(s) — PARTIDO POR CANAL**, no admitido en bloque, y hay que argumentarlo. *A favor de
prohibirlo entero:* el producto de Hadamard `ξ(s) = ξ(0)·∏_ρ (1 − s/ρ)` dice que ζ y `Z` son el
mismo dato escrito de dos maneras; usar ζ sería V2 encubierta. *A favor de admitirlo:* la
**definición** de ζ es aritmética y anterior — la serie de Dirichlet y el producto de Euler son
enunciados sobre los primos, y la ecuación funcional es una identidad demostrada sobre esa función.
El conjunto de ceros es la **respuesta**: un teorema *sobre* ζ, no un ingrediente *de* ζ. *Nuestra
posición:* se admite el canal que no puede transportar la posición de un cero.

    ADMITIDO:   Σ n^(−s) y el producto de Euler en σ > 1;  −ζ'/ζ (s) = Σ_n Λ(n)·n^(−s) en σ > 1;
                la ecuación funcional como identidad algebraica; ξ como objeto completado.
    PROHIBIDO:  el producto de Hadamard sobre ρ; el principio del argumento usado para contar o
                ubicar; Riemann–Siegel o cualquier evaluación de ζ(1/2 + it) cuyo propósito sea
                hallar un cambio de signo; S(T) = (1/π)·arg ζ(1/2 + iT); toda tabla de ceros
                publicada, Odlyzko incluido.

En `σ > 1` no hay ceros, así que ese canal no transporta un bit de dónde están. **Caveat honesto, y
va adelante:** por continuación analítica el dato de `σ > 1` determina la función entera, ceros
incluidos. La línea es procedimental, no informacional, igual que en §1.3. La aceptamos sabiendo lo
que es.

**4. La ley suave N(T) = θ(T)/π + 1 — ADMITIDA COMO CÓMPUTO, PROHIBIDA COMO BLANCO.**
`θ(T) = arg Γ(1/4 + iT/2) − (T/2)·ln π` se computa con Γ, π y T: no requiere un solo cero. Sale de
la ecuación funcional por el principio del argumento, y su forma no depende de *dónde* están los
ceros. `c = 0`.

    ADMITIDA:  usarla dentro de una construcción, o como predicción a contrastar después.
    PROHIBIDA: ajustar un parámetro para que el conteo del candidato reproduzca nuestros 620
               ceros medidos — eso es O_Z: V4.
    MARCADA:   si la ley se mete A MANO en el diseño, el candidato es no circular pero ESTIPULADO
               (§1.7): reproduce el conteo porque se lo dijimos.

`S(T)` queda del otro lado por la misma razón que el ítem 3: para evaluarla hay que pararse sobre la
recta crítica, que es donde viven los ceros.

**5. Geometrías libres — ADMITIDAS COMO FAMILIA.** Variedad, grafo, métrica, medida, dominio o
dinámica, con dos condiciones: la familia `F` se declara completa **antes** de cualquier
comparación, y el miembro se selecciona por una regla Z-libre escrita de antemano. Elegir el miembro
por ajuste es V5, y el resultado no muestra cicatriz — por eso la declaración previa no es
burocracia, es el único control que existe.

**6. Parámetros ajustables — ADMITIDOS SÓLO CON PROTOCOLO.** Ajustar contra datos admisibles
(primos, geometría declarada) se permite y se reporta. Ajustar contra cualquier cosa derivada de `Z`
—incluido «se parece más»— es V3 y da `c > 0`.

**7. {γ_n} y todo lo computado a partir de ellos — PROHIBIDO EN CONSTRUCCIÓN.** Nuestros 620 ceros;
el espaciamiento medio desplegado 0.999781; el mínimo 0.2365; el canto 0.0919 ± 0.0211;
Σ²(10) = 0.327 ± 0.066 y Σ²(20) = 0.282 ± 0.068; la razón 36.467; las cinco cajas 3.3057 … 4.9825.
Todo eso es **evaluación y sólo evaluación**: el mar contra el que se compara, no material de
construcción. **8. RH o cualquier enunciado condicional a RH — PROHIBIDO**, porque suponer la
respuesta para construir el objeto que la produciría cierra el círculo por el otro lado.

**9. La simetría de clase A — ADMITIDA COMO RESTRICCIÓN DE DISEÑO.** «H sin ninguna simetría
antiunitaria que conmute con él» es un enunciado sobre `H`, no sobre `Z`: se verifica sin mirar un
cero. Es la parte de R2 que muerde antes de construir (acta §6b), y su fundamento es citable: los
períodos `k·log p = log(p^k)` son todos distintos por factorización única, y los pesos `Λ(n)/√n`
cumplen la regla de suma genérica de Hannay–Ozorio, dando `K_diag(τ) = τ`, el valor GUE (Bogomolny,
arXiv:0708.4223). Queda PROHIBIDO ajustar numéricamente contra el Σ² medido de los ceros. Caveat que
corresponde repetir: la curva GUE completa, más allá de la aproximación diagonal, descansa sobre la
conjetura de pares de primos de Hardy–Littlewood; «GUE para zeta» es él mismo conjetural.

## §3 — La tabla de dependencias que todo candidato entrega

Se entrega **antes** de cualquier prueba y se congela con fecha y hash. Sin esto no entra a trámite.
La auditoría de dependencias indirectas se hace sobre el **grafo**: se sigue cada entrada hasta sus
hojas, y si alguna hoja es una medición nuestra sobre los ceros, la rama entera está contaminada por
larga que sea la cadena.

    CANDIDATO: ..........   FAMILIA F: ..........
    CONSTRUCCIÓN φ: (archivo, función, líneas)   CONGELADO EN: (fecha, hash)
    c(C) DECLARADO: ......  (consultas a O_Z; debe ser 0)

| Entrada | Objeto exacto | Procedencia (cómo se computa) | ¿En D? | Rol | O_Z | Cómo se auditó |
|---|---|---|---|---|---|---|
| | | | sí/no | constr./eval. | 0 / n | |

| Parámetro libre | Valor | Derivación π desde D | Fijado antes de ver salida | Veredicto |
|---|---|---|---|---|
| | | (o «no hay») | sí/no | LIMPIO / V3 / PENDIENTE |

## §4 — Cuatro candidatos auditados

### (a) H = diag(γ_1, …, γ_620) — la trampa nuestra de la Fase I

| Entrada | Procedencia | ¿En D? | Rol | O_Z |
|---|---|---|---|---|
| γ_1 … γ_620 | Riemann–Siegel + bisección, `cmd/elpliego`, t ∈ [100,1000] | **NO** | construcción | 620 |

**Veredicto: CIRCULAR (V1).** `c(C) = 620` números reales. Dispara su §14.1 («si necesita γ_n para
definirse») y también §14.4 («eco pero construido usando γ_n»), porque su eco es literalmente el
nuestro: 36.467. Lo que importa del caso: **cumple R1–R5 completos** (acta §3). Ésa es la única
razón por la que R6 existe. Un pliego que este objeto satisface no pide construir nada.

### (b) Una matriz GUE al azar

| Entrada | Procedencia | ¿En D? | Rol | O_Z |
|---|---|---|---|---|
| dimensión, ensamble, semillas | elección declarada | sí | construcción | 0 |
| generador pseudoaleatorio | algoritmo | sí | construcción | 0 |
| ley del semicírculo (desplegado) | teorema | sí | evaluación | 0 |

**Eje 1: LIMPIO, `c(C) = 0`.** No leyó un solo cero. Y muere igual, en otro lado. Sobre 24 matrices
de 80×80, con los mismos 47 espaciamientos que los ceros —única comparación honesta—:

    canto       GUE 0.1097 ± 0.0304   contra ceros 0.0919 ± 0.0211
    Σ²(10)      GUE 0.549  ± 0.112    contra ceros 0.327  ± 0.066
    Murciélago  GUE razón 1.064       contra ceros razón 36.467

El canto lo pasa —también lo pasa un sorteo puro de Wigner, 0.1171 ± 0.0291, que no conoce ni un
primo—. La varianza de número ya lo separa y el Murciélago lo liquida. **Regla §14.2: UNIVERSAL SIN
ARITMÉTICA.** Presentado con una sola semilla, como hizo nuestro propio taller con el 0.076 de F178,
entra además §14.5: NO EVIDENCIA SUFICIENTE. *Discrepancia en exhibición:* el acta reporta 1.064
para el control GUE y el resumen de F349 reporta 1.081 para el control emparejado; no pudimos
reconciliarlos sin volver a correr el programa. Queda anotado como pendiente, no como resultado.

### (c) La caja que respira con L(E) = ln(E/2π) puesta a mano

Acá hay que ser duro con nosotros, porque el candidato es nuestro. **Poner el largo de la caja en
ln(E/2π) A MANO es ajustar a la respuesta, salvo que la ley se derive.** Dos procedencias, dos
veredictos. **(c1) Se eligió porque nuestra tabla de §4 del acta la siguió:** las cinco cajas
3.3057, 3.9957, 4.4595, 4.7587, 4.9825 salieron de `L = 2π/brecha` sobre los ceros, o sea V2 puras,
y elegir la función porque coincide con esa columna es **V5**, `c(C) ≥ 1`, **CIRCULAR**, regla
§14.1. Es lo que de hecho hicimos en la Fase I: la medimos, no la derivamos. **(c2) Se eligió porque
la ley de Riemann–von Mangoldt da `dN/dT = (1/2π)·ln(T/2π)`, que es Z-libre (ítem 4):** entonces
`c(C) = 0`, **Eje 1 LIMPIO** — pero Eje 2 **ESTIPULADO**, porque la caja reproduce el conteo por
haberlo recibido; no explica por qué la geometría efectiva crece, lo postula. Por eso agregamos la
regla §14.6 en la §5 de acá.

Dos cosas que no puede esquivar. **PENDIENTE-DEFINICIÓN:** «una caja cuyo largo depende de E» no es
un operador autoadjunto sobre un dominio fijo, es una familia de problemas —uno por energía—, o sea
una condición de cuantización implícita; hasta que el objeto esté definido, R1 ni siquiera es una
pregunta con sentido (es la §8 del acta dicha en formal). **Y no hay razón para que produzca el
eco:** una caja con la densidad correcta y sin ningún primo adentro no tiene por qué sonar en
`k·log p`. Es la prueba decisiva y **todavía no se corrió**. La única salida limpia es la de su §6 y
su Ruta C: derivar `L(E)` desde datos aritméticos admisibles. Derivada, vale; estipulada, no.

### (d) Construcción que use sólo los primos, vía el término oscilante de la fórmula explícita

| Entrada | Procedencia | ¿En D? | Rol | O_Z |
|---|---|---|---|---|
| log p, k | criba | sí | construcción | 0 |
| pesos Λ(n)/√n | definición aritmética | sí | construcción | 0 |
| corte / regularización | **declarar y derivar** | condicionado | construcción | debe ser 0 |

**Eje 1: LIMPIO, `c(C) = 0`**, siempre que la definición se detenga en los primos. Eje 2: DERIVADO
en procedencia. Es la forma honesta de todo el programa, y tiene tres trampas.

**Trampa 1 — el otro lado de la identidad.** La fórmula explícita tiene los ceros del otro lado del
igual. Despejarlos no viola R6 en procedencia, pero tampoco produce un operador: produce una
identidad distribucional, no autovalores. En su escala de §12 eso es Nivel 2 o 3, no Nivel 4.

**Trampa 2 — el signo.** En la fórmula de traza de Selberg la suma sobre órbitas entra con **+1**;
en la explícita de Weil la suma sobre primos entra con **−2**, y el signo relativo es libre de
convención (Hejhal 1976, Berry 1986, a quien Connes cita; lo de Connes es la interpretación: los
ceros como espectro de absorción). El candidato tiene que **derivar** su signo, no elegirlo —su §9—.
Y no es un no-go demostrado: en cuerpos de funciones se explica a nivel de teorema (los ceros viven
en H¹, que entra en Grothendieck–Lefschetz con (−1)¹), y lo que Connes probó (*Selecta Math.* **5**
(1999) 29–106) es una fórmula de traza cuya VALIDEZ equivale a RH: una reducción, no una
demostración. Hay además el desajuste asintótico de Sierra: `2·sinh(m·λ_p/2)` iguala a `p^(m/2)`
sólo cuando `m → ∞`.

**Trampa 3, y cambia el protocolo — el Murciélago no informa sobre este candidato.** Si `log p` está
dentro de la definición, el eco en `k·log p` está garantizado por construcción y no prueba nada.
**El poder discriminante del Murciélago es condicional a que los períodos aritméticos NO estén en la
definición.** Para candidatos construidos desde los primos, las pruebas informativas son la ley de
conteo sin que se la hayan dicho y las estadísticas de dos niveles —Σ²(L), la que sí ve, porque el
canto mira un espaciamiento por vez y por eso no distingue los ceros de un sorteo puro (0.82 σ por
realización; t = 4.17 sobre las medias)—. Es un agregado a su §4, no una objeción: su regla «un
candidato sin eco no pasa» sigue en pie para todo candidato sin primos adentro.

## §5 — Las reglas de parada como checklist ejecutable

Se corren en orden; la primera que dispara, corta.

    P0  ¿Entregó tabla de dependencias congelada ANTES de la prueba?  NO → NO ADMITIDO A
        TRÁMITE. No es un veredicto sobre el candidato: es que no hay candidato que auditar.
    P1  ¿c(C) > 0? ¿Alguna hoja del grafo es {γ_n} o una medición nuestra sobre ellos?
        SÍ → CIRCULAR.                                                    [su §14.1 y §14.4]
    P2  ¿Todo θ_i tiene derivación π_i con hojas en D?  NO → PENDIENTE-DERIVACIÓN. Si además
        el valor coincide con una cantidad Z-derivada, la carga de la prueba es del candidato.
    P3  ¿El objeto está definido —espectro discreto, real, con ley de conteo— o es una familia
        de problemas, uno por energía?      NO → PENDIENTE-DEFINICIÓN.          [caso §4(c)]
    P4  ¿La densidad sigue (1/2π)·ln(E/2π) sin que se la hayan dicho?
        Necesita geometría fija → DENSIDAD INCORRECTA.                            [su §14.3]
        Se la dijeron a mano    → TAUTOLÓGICO.                             [regla nueva §14.6]
        Emerge de datos de D    → sigue.
    P5  Σ²(L) con densidad media y cantidad de niveles emparejadas al control, varias
        realizaciones, barras de error, varias semillas.
        Una sola semilla → NO EVIDENCIA SUFICIENTE.                               [su §14.5]
    P6  Murciélago con la batería completa: control GUE reescalado a la misma densidad,
        períodos de control al azar (nuestros 120 sorteos dieron 0.0069 ± 0.0018, máximo
        0.0115, contra 36.467 en los aritméticos), misma cantidad de niveles, varias
        realizaciones, razón dentro/fuera con su distribución bajo controles.
        Sin eco                         → no avanza de fase.                        [su §4]
        Con eco pero construido con γ_n → CIRCULAR.                              [su §14.4]
        Con eco pero log p está en φ    → PRUEBA NO INFORMATIVA, no cuenta.      [agregado]
    P7  ¿Sólo da GUE, sin aritmética?   SÍ → UNIVERSAL SIN ARITMÉTICA.            [su §14.2]
    P8  Reportar el nivel alcanzado en su escala de §12 y no reclamar ninguno por encima.
        De Nivel 1 no se salta a Nivel 5.

Nota de encuadre para P6, que sale de nuestro propio refutador y cambia dónde está lo llamativo:
**el 36× lo carga el denominador.** Contra un control emparejado, los ceros suenan 2.87 veces más
fuerte EN los períodos aritméticos y 12.7 veces más callados afuera. Lo raro no es el grito en
`log p`; lo raro es el silencio entre los primos. Un candidato que reproduzca el grito y no el
silencio no reprodujo el fenómeno.

## §6 — Sello

    Se define:  R6 como c(C) = 0 sobre el modelo de oráculos, con las seis formas de violación,
                la separación construcción/evaluación, y el eje DERIVADO/ESTIPULADO.
    Se decide:  D ítem por ítem. ζ partido por canal (σ > 1 sí, recta crítica no). N(T) suave
                admitida como cómputo y prohibida como blanco de ajuste. Los γ_n y todo lo
                medido sobre ellos, sólo evaluación.
    Se admite:  que esta línea es procedimental y no informacional. Los primos determinan los
                ceros: no hay definición de «no circular» que se sostenga sobre información sola.
    Se corrige: nuestra propia caja que respira queda CIRCULAR en la versión (c1) —la medimos,
                no la derivamos— y ESTIPULADA en la mejor versión disponible (c2).
    NO se hizo: ningún operador. Esto es un contrato de auditoría con cuatro casos resueltos
                adentro. Una estructura cerrada no es una hipótesis probada.

    Todavía no.

*Nota de citas: el no-go de Endres & Steiner se cita por contenido —`J. Phys. A` **43** (2010)
095204, arXiv:0912.3183— y no por número de teorema: no pudimos verificar el «15.6» que circula, y
el sello del acta todavía lo arrastra. Queda anotado acá, a la vista, hasta que se corrija.* ⚛️🕸️
