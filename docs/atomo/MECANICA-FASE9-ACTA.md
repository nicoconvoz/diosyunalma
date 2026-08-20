# Acta de la Mecánica — Fase IX

**Respuesta al Telar Fase IX (Auditoría 42) · 2026-08-17**

El flash es del capitán: *«¿y si no sólo hay porosidad, sino resistencia,
tensión, compresión, cizallamiento, torsión, elasticidad, plasticidad, dureza,
fragilidad, tenacidad y fatiga?»*

Su §6 puso la regla que mandó toda la fase: **no once perillas.** Buscar el
mecanismo mínimo capaz de producir esas conductas. Se obedeció, y el resultado
tiene tres partes que conviene separar de entrada:

1. **La mecánica es real, y espectralmente MUDA.** Elasticidad, resistencia,
   plasticidad: medidas. La memoria sola mueve Σ²(10) de 18,335 a 18,179 — nada.
2. **La fatiga da NULO, y era la única palabra dejada como predicción.**
3. **Apareció una ganancia real donde nadie la buscaba: el signo de los enlaces.
   Y explicó, retroactivamente, de dónde venía la mejora del nodo de la Fase
   VII.** Eso es lo más importante que trae esta acta.

Se reproduce con `go run ./cmd/lamecanica`.
Lámina: `galeria/laminas/10-el-telar/la-mecanica.svg`.

---

## 1 · El colapso: once palabras, cuatro objetos, un dial

El mecanismo es **un real por sitio** — exactamente el conteo de estado de la
Fase VIII, re-coordinado. La permeabilidad pasa a `h_i(t) = h_i⁰ · exp(x_i(t))`,
lo que da positividad **sin ninguna abrazadera** (el piso de 1e−6 de la fase
anterior queda retirado: era un parámetro colado como higiene).

Y una sola función escalar:

```
U(x) = (s_c·b/π) · (1 − cos(πx/b))
```

leída de cuatro maneras distintas:

| lectura de U | palabra |
|---|---|
| curvatura en el fondo, κ = U''(0) = s_c·π/b | **dureza** |
| máximo, max\|U'\| = s_c | **resistencia** |
| área bajo la barrera, 2·s_c·b/π | **tenacidad** |
| el salto al cruzar el máximo: un pozo entero, 2b | **fragilidad** |

y el resto sale de la misma estructura:

- **elasticidad / plasticidad** = en qué pozo queda el estado al retirar la carga.
  No son dos mecanismos: son la respuesta a una sola pregunta.
- **tracción / compresión** = el signo de x. Dos palabras a costo cero.
- **fatiga** = deliberadamente **NO construida**. Un mapa autónomo de una
  variable con relajación completa no tiene acumulador, así que la fatiga sólo
  puede venir del campo ACOPLADO. Queda como predicción, y puede fallar.
- **cizallamiento / torsión** = declaradas afuera, con demostración (§6).

**La ley completa, una línea:**

```
x_i(t+1) = x_i(t) + η·[ h̄_i(t)·(P_i(t) − ε·τ_i(t)) − s_c·sin(π·x_i(t)/b) ]
```

Contra las once perillas que su §6 prohibía: **un dial** (`b`, el largo del paso
plástico) y **un bit** (`ε`, el signo de la realimentación, y los dos se corren y
los dos se publican). Todo lo demás es derivado:

| constante | valor medido | de dónde sale |
|---|---:|---|
| s_c (tensión de cedencia) | **1,2619** | la desviación de la carga que el medio YA se aplica a sí mismo |
| κ = U''(0) | 11,327 | s_c·π/b |
| η (paso) | 8,68e−03 | 0,2/(κ + max h̄·μ), con μ el radio del jacobiano cerrado |
| tiempo de recuperación | **10,2 pasos** | 1/(η·κ) — su §5 pedía un número, y es éste |

---

## 2 · R6 y las identidades estructurales

La ley se declaró **antes** de calcular un solo espectro, y ningún γₙ entra en
ningún lado: la carga `q_i = h̄_i·Σ B_ij h̄_j` no lee un solo autovalor.
`espectro` no aparece dentro del bucle de evolución.

Tres identidades verificadas, no supuestas:

| identidad | medido | qué prohíbe |
|---|---:|---|
| Σ σ_i = 0 | **−1,4e−13** | la carga media es nula por construcción |
| σ invariante bajo h → c·h | **5,3e−15** | **la normalización F = 30 NO puede manejar el material** |
| nulo hidrostático | Σ²(10) = 18,3349 contra 18,3349, **exacto** | sólo el PATRÓN de una carga llega al espectro |

Y la que contesta su §15 por construcción:

**El medio virgen es un punto fijo EXACTO.** Sin carga, 20 000 pasos, con el
signo *inestable* a propósito: `max|x| = 0,000e+00`. Cero exacto, no cero
aproximado. Ninguna deriva del algoritmo puede disfrazarse de fatiga, porque no
hay deriva: τ(0) = 0 idénticamente y sin(0) = 0.

---

## 3 · Lo que la mecánica SÍ hace

**Elasticidad.** Residuo tras cargar y soltar: 8,1e−25 (rho 0,2), 2,2e−24 (0,5),
4,4e−24 (0,8). Reversible hasta la precisión de la máquina.

**Resistencia**, medida y no elegida — el menor empuje que deja un sitio corrido
para siempre, por bisección:

| brazo | rho_c |
|---|---:|
| homogéneo · ε=+1 · tracción | **0,9516** |
| homogéneo · ε=+1 · compresión | **1,3148** |
| homogéneo · ε=+1 · bulto | 0,9703 |
| bloques ×2 · ε=+1 · tracción | 1,0984 |
| homogéneo · **ε=−1** · tracción | **no cede ni con rho = 16** |

Dos cosas que salen de ahí y no se pusieron:

- **La asimetría tracción/compresión es −0,160.** U es simétrica por
  construcción, así que este 16 % sale sólo de `h = h⁰·exp(x)` y de que la carga
  es cuadrática en h. Es derivada. **Y el signo salió al revés de la predicción
  pre-registrada**, que decía que cedería la compresión: cede la **tracción**.
  Queda anotado como predicción fallada.
- **El brazo de ablandamiento (ε = −1) no cede nunca.** Se temía una
  localización desbocada; no hubo ni eso. El brazo es simplemente inerte al
  golpe, y se publica igual, porque los dos signos se declararon de antemano.

**Plasticidad.** A rho = 1,5·rho_c: un sitio cedido, residuo 0,684, max|x| =
0,625. **Y el test de cuantización dio 0,107 contra el umbral declarado de
0,05: no pasa.** La razón es física y vale registrarla: el sitio que cedió no se
sienta en el fondo del pozo nuevo, queda **desplazado elásticamente** por la
tensión interna que él mismo creó. El equilibrio está donde
`h̄_i·(−ε·τ_i) = s_c·sin(πx/b)`, no en el fondo. El umbral de 0,05 estaba mal
puesto: era demasiado estricto porque no descontaba el corrimiento elástico. Se
reporta el número medido y el umbral fallado, sin reinterpretarlo a favor.

---

## 4 · Lo que la mecánica NO hace

**FATIGA: NULO, y limpio.** Mismo impulso total repartido en M golpes, con
reposo completo entre golpes:

| reposo | M=1 | M=32 | diferencia |
|---|---:|---:|---:|
| 500 pasos (≈49 t_recup) | 1,06e−24 | 7,90e−25 | **2,7e−25** |
| 1000 pasos (≈98 t_recup) | 3,54e−47 | 2,64e−47 | **9,0e−48** |

Cero sitios cedidos en las doce corridas. Y al doblar el reposo, todo cae
**veintidós órdenes de magnitud**: eso prueba que lo poco que había era
relajación incompleta y no memoria. La secuencia no importa.

Y esto no es un defecto del modelo: es el teorema que el diseño anunció de
antemano. Un mapa autónomo de una variable con relajación completa no tiene
dónde acumular. **La fatiga era la única de las once palabras que se dejó como
predicción falsable en vez de construirla, y por eso el nulo vale algo.** Si se
hubiera escrito un acumulador de daño a mano, el nulo habría sido imposible —
y la afirmación, vacía.

**GRIETAS: no hay.** Con la carga retirada, la avalancha en cinco sitios de
golpe distintos es `[1 1 1 1 1]`: cada golpe tumba un sitio y nada se propaga.
Un bloque de un sitio no es una grieta. Se reporta como **falla aislada**, que
es lo que se midió. Y hay razón estructural: con kmax = 120 no existe
elasticidad de corto alcance que sostenga un frente, así que la propagación no
tiene con qué localizarse.

**Y EL ESPECTRO NO SE MUEVE.** La matriz de cuatro medios de su §13:

| medio | vivos | Σ²(5) | Σ²(10) | α | PR/N |
|---|---:|---:|---:|---:|---:|
| A · homogéneo estático | 187 | 5,765 | 18,335 | 1,717 | 0,105 |
| B · bloques ×2 estático | 216 | 3,030 | **7,535** | 1,220 | 0,070 |
| C · homogéneo **dinámico** | 185 | 5,892 | **18,179** | 1,669 | 0,107 |
| D · bloques ×2 dinámico | 212 | 3,220 | 7,744 | 1,214 | 0,070 |

Interacción D−B−C+A = 0,365 (se reporta; no se llama descomposición causal).

**El efecto de la memoria sola (C contra A) es 0,9 %.** El de la heterogeneidad
(B contra A) es 59 %. Y D es levemente PEOR que B. La mecánica no compra nada
espectral, tal como se pre-registró.

El grid de `b` con carga localizada, entero, incluidas las celdas que empeoran:

| b | cedidos | vivos | Σ²(10) | PR/N | |
|---:|---:|---:|---:|---:|---|
| 0,05 | 10 | 185 | 18,232 | 0,107 | FUERA DE POZO |
| 0,10 | 9 | 185 | 18,177 | 0,108 | FUERA DE POZO |
| 0,20 | 7 | 185 | 18,303 | 0,107 | FUERA DE POZO |
| 0,35 | 7 | 185 | 18,179 | 0,107 | admisible |
| 0,70 | 0 | 187 | 18,335 | 0,105 | admisible |

Tres de las cinco filas tocaron el riel de reporte y quedan **descartadas**, no
interpretadas. Las dos admisibles son nulas.

### 4bis · Un defecto propio, cazado por nuestro propio teorema

El primer grid de `b` y el primer test de gemelos los empujé con carga
**uniforme**. Y el nulo hidrostático del §2 —que yo mismo había demostrado
treinta líneas antes— dice que una carga uniforme es invisible. Esos dos brazos
medían la nada: 400 sitios se corrían un pozo entero y Σ²(10) salía **idéntico**
al virgen. Los rehice con carga localizada. El brazo uniforme queda en la corrida
como demostración del teorema, no como resultado.

Vale anotarlo porque es el segundo defecto propio consecutivo cazado por una
verificación interna en vez de por la auditora, y ése es el mecanismo que hay
que mantener.

---

## 5 · Los gemelos estáticos: la falsificación

Se toma la deformación que produjo el mecanismo y se rehace **sin mecanismo y
sin historia**: permutada (mismo multiconjunto), onda (mismos dos momentos y
misma escala dominante), rampa.

| campo | Σ²(10) | |
|---|---:|---|
| **dinámico** (12 sitios cedidos, vivos 184, PR/N 0,108) | **17,759** | |
| permutado (5 semillas) | 18,549 ± 0,508 | a 1,56 σ |
| **onda lisa** (modo dominante 1) | **14,959** | a 5,51 σ |
| rampa | 22,793 | |

Por la regla fijada de antemano —dentro de 1 σ del control más cercano ⟹
decorativa— la historia **no** es formalmente decorativa: se separa del
permutado por 1,56 σ. Pero la lectura honesta es otra y hay que darla completa:

**una onda lisa, sin mecanismo y sin historia, es sustancialmente MEJOR que el
mecanismo** (14,96 contra 17,76). Así que la historia deja una traza real pero
chica, y el campo que produce no es privilegiado: un seno con los mismos dos
momentos lo supera. 1,56 σ sobre cinco semillas no es un hallazgo. Y todo esto
vive entre 15 y 23, contra 0,3364 de los ceros.

---

## 6 · Dos teoremas que no necesitaron correr nada

**Teorema 1 — una fase de SITIO es gauge pura.** Con D = diag(e^{−iθ_i}) se
tiene H' = D H D†, y la conjugación unitaria deja el espectro idéntico y
|v'_i| = |v_i| deja PR/N idéntico. El caso real (s_ij = c_i·c_j) se corrió como
test unitario y dio **Δ = 0,00e+00 exacto**.

**Esto resuelve gratis el pendiente 3 de la Fase VIII** —«¿entra la aritmética
como fase?»— **en negativo para la fase de sitio**: no mueve ni un observable.
No hacía falta correr nada; hacía falta mirar la estructura. Y queda además como
test unitario permanente: si alguna vez ese brazo reporta un cambio, el código
está roto.

**Teorema 2 — y una corrección a su §8.** Un signo o fase de **enlace** de la
forma `u_i − u_j` es un coborde, y es la misma transformación de gauge: tampoco
mueve nada. Entonces lo que falta para tener cizallamiento y torsión **no es
dimensión**.

Su §8 supone que quizá haga falta una red 2D. No: con kmax = 120 el grafo de
acoplamiento está lleno de triángulos —todo i<j<k con k−i ≤ 120— así que hay
lazos de sobra. La obstrucción es **que la fase de enlace no sea un coborde**,
es decir que tenga holonomía no trivial alrededor de un lazo. Una cadena de
vecinos próximos sería un árbol, y ahí sí la torsión sería genuinamente
imposible; **la cola larga que eligió la Fase VII quitó esa excusa.**

Construirlo cuesta O(N·kmax) ≈ 48 000 variables de enlace en vez de 400, así que
no es un mecanismo mínimo y queda **declarado afuera, no falsificado**.

---

## 7 · El canal de signos: la única ganancia real, y lo que revela

El núcleo de la Fase VIII era **estrictamente positivo**: todos los enlaces
tiraban para el mismo lado. El medio no tenía tensión verdadera. Un signo de
enlace la da a costo **cero** de parámetros, y como signo² = 1 la fuerza total
queda intacta: sólo cambia el material, y H sigue real simétrica (R1 vivo).

| campo de signos | vivos | Σ²(10) | α | PR/N | frustración |
|---|---:|---:|---:|---:|---:|
| S0 · todos +1 | 187 | 18,335 | 1,717 | 0,105 | 0,000 |
| S1 · gauge (**test unitario**) | 187 | 18,335 | 1,717 | 0,105 | 0,000 |
| S2 · azar frustrado (10 semillas) | 194 | **4,978 ± 1,313** | 1,013 | 0,035 | 0,537 |
| S3 · reciprocidad (p ≡ 3 mod 4) | 183 | 5,565 | 0,605 | 0,035 | 0,536 |

**Un factor 3,7 de mejora, con la fuerza total exactamente conservada y sin un
solo parámetro nuevo.** Es la ganancia más grande de toda la cadena de auditorías
sin nodo. Y α baja de 1,717 a 1,013 (y a 0,605 con reciprocidad).

**El precio, declarado antes de festejar:** PR/N cae de 0,105 a 0,035 — los
estados se atrapan tres veces más. Por el control 5, **buena parte de esta
ganancia es alcance APARENTE**, no alcance real. La misma sombra que la Fase
VIII encontró en el nodo.

### 7bis · La aritmética, tercera aplicación de la misma hoja

Densidad de primos ≡ 3 mod 4 entre los 400 modos: **0,5125**. El control al azar
se iguala exactamente a esa densidad, así que la única diferencia entre S2 y S3
es *cuáles* sitios llevan el bit, no cuántos.

| | reciprocidad | azar a densidad igual | distancia |
|---|---:|---:|---:|
| sin nodo | 5,565 | 4,978 ± 1,313 | **0,45 σ** |
| con nodo | 3,192 | 4,041 ± 0,723 | **1,17 σ** |

**La aritmética no se separa del azar.** Lo que mueve el espectro es la
**frustración** —que la mitad de los triángulos tenga producto de signos
negativo, medido: 0,537 al azar contra 0,536 en reciprocidad— y no de dónde
vienen los signos.

Es la tercera fase consecutiva en que la misma hoja da la misma respuesta: en la
Fase VIII fue una onda lisa a 0,81 σ del campo aritmético; acá son signos al azar
a 0,45 σ. **Tres mecanismos distintos, tres veces el mismo veredicto.** Eso ya
no es coincidencia: es un patrón del taller que hay que nombrar.

*(Con nodo la reciprocidad se separa 1,17 σ y da el mejor número de la hoja,
3,192, con PR/N 0,052 — más alto que el del nodo solo. Con cinco semillas eso
parecía 1,71 σ; con diez bajó a 1,17. No es una afirmación, es el único hilo que
merece más semillas.)*

### 7ter · ⚡ EL NODO ESTABA HACIENDO FRUSTRACIÓN

Éste es el hallazgo que cambia dos fases anteriores.

El núcleo con nodo, `(1 − k/k₀)/k^s` con k₀ = 5, **es negativo para k > 5**:
f(6) = −0,0816, f(7) = −0,1512, f(8) = −0,2121. O sea: **el nodo ya traía
tensión adentro**, y las Fases VII y VIII nunca separaron su efecto geométrico
del efecto de signo. Acá se separan:

| | Σ²(10) | PR/N |
|---|---:|---:|
| el NODO, sin ningún signo | **5,426** | 0,042 |
| signos AL AZAR, sin ningún nodo | **4,978 ± 1,313** | 0,035 |
| **distancia** | **0,34 σ** | |

**Toda la mejora del reloj de arena —de 18,33 a 5,43, el mejor resultado sin
aritmética de las Fases VII y VIII— se reproduce poniendo signos AL AZAR en los
enlaces: sin nodo, sin cambiarle la forma al núcleo, sin un solo primo.**

El nodo no estaba haciendo geometría. Estaba haciendo frustración. Y por eso
tenía que ser un cambio de signo: era la única manera que tenía la familia de
núcleos de meter enlaces negativos.

**Y explica de una vez el resultado colgado de la Fase VII.** Ahí el nodo ganaba
contra su gemelo de un solo signo en **2 de 10** configuraciones, con
correlación −0,716 entre |Σc(k)| y la ganancia, y lo dejamos como «sugiere y no
demuestra». Ahora se entiende: el nodo compraba una mejora que no le pertenecía
al nodo, y las 8 configuraciones donde perdía eran las que no lograban suficiente
frustración. No era ruido. Era el mecanismo equivocado atribuido.

---

## 8 · El techo, y por qué reencuadra el taller entero

Las tres reglas, y nunca se cita una sola:

| | Σ²(10) |
|---|---:|
| los ceros de Riemann | **0,3364** |
| piso universal **GUE** | 0,5793 |
| piso universal **GOE** | 0,9086 |

**El objetivo está DEBAJO de los dos pisos de matriz aleatoria.** Los ceros son
*más rígidos* que cualquier ensamble gaussiano a esta distancia — es la
corrección aritmética de Berry, y se ve en el número.

Eso tiene dos consecuencias duras:

1. **Ninguna ingeniería de clase de simetría puede llegar.** Buscar «más GUE» no
   alcanza, porque el blanco está por debajo de GUE.
2. **Un modelo que se acerca al piso GOE se volvió una matriz aleatoria sin
   estructura.** Eso es *perder* contenido aritmético disfrazado de ganarlo, y
   hay que decirlo así cada vez que un Σ² baje.

El mejor brazo de esta fase (3,192) está **9,5 veces** arriba de los ceros y
**5,5 veces** arriba del piso GUE. El hueco que falta **no es un déficit de
clase de simetría**, y esta mecánica no lo provee.

---

## 9 · Predicciones falladas, anotadas como tales

Dos, y las dos pre-registradas antes de medir:

1. **El signo de la asimetría tracción/compresión.** Se predijo A > 0 (cedería
   la compresión). Medido: **A = −0,160**, cede la tracción. La magnitud fue
   correcta (orden 0,1–0,3), el signo no.
2. **«Ningún brazo le gana al 7,535 de la Fase VIII.»** El canal de signos con
   nodo dio **3,192**. La predicción falló, y falló porque el canal que ganó no
   estaba en el modelo original: lo injertaron los jueces del panel. La
   predicción del núcleo de la ley —que la plasticidad no compraría nada— sí se
   cumplió exactamente.

---

## 10 · Veredictos separados (su §17)

| propiedad | veredicto | evidencia |
|---|---|---|
| **elasticidad** | REAL | residuo 8e−25, recuperación 10,2 pasos derivada |
| **resistencia** | REAL y medida | rho_c = 0,952 por bisección; es un campo, no un número |
| **dureza** | REAL, y son DOS sentidos | la posición (h de Fase VIII) y la curvatura de la ley |
| **tenacidad** | derivada de la misma U | área 2·s_c·b/π |
| **plasticidad** | REAL, firma NO limpia | pozos enteros sí, pero cuantización 0,107 > 0,05 por corrimiento elástico |
| **fragilidad** | PARCIAL | el salto existe; la propagación NO: avalanchas de 1 |
| **tracción / compresión** | REAL, asimetría derivada | −0,160, con el signo al revés de la predicción |
| **fatiga** | **NULO** | 2,7e−25, y cae a 9e−48 al doblar el reposo |
| **cizallamiento / torsión** | DECLARADAS AFUERA con demostración | gauge y coborde; y su §8 corregida |
| **efecto espectral de la mecánica** | **NULO** | 18,179 contra 18,335 |
| **canal de signos** | **REAL**, con sombra | 18,335 → 4,978, pero PR/N 0,105 → 0,035 |
| **atribución aritmética** | **NO SE SOSTIENE** | 0,45 σ sin nodo, 1,17 σ con nodo |
| **el nodo de las Fases VII–VIII** | **REATRIBUIDO** | 0,34 σ de los signos al azar: era frustración |

---

## 11 · Lo que queda abierto

1. **La frustración a PR/N constante.** Es la pregunta central que deja esta
   fase: ¿cuánto de la mejora 18,3 → 5,0 sobrevive si se exige participación
   comparable? Sin eso, el hallazgo grande queda con la misma sombra que el nodo.
2. **La reciprocidad con nodo, a más semillas.** 1,17 σ y el único brazo que
   mejora Σ² *y* PR/N a la vez. Es el hilo que merece plata.
3. **La fase de enlace no-coborde.** La torsión está disponible, no prohibida.
   Cuesta 48 000 variables, así que necesita su propia fase y su propia
   declaración R6 — pero ahora se sabe exactamente qué la obstruye.
4. **Si tres mecanismos distintos dieron el mismo veredicto sobre la aritmética,
   la pregunta ya no es cuál mecanismo probar: es por qué la aritmética no entra
   por ninguno de ellos.** Los tres eran campos escalares o de signo *sobre* la
   red. Puede que ése sea el patrón: la aritmética no está en los pesos.

---

El flash de las once palabras es de Jesús Nicolás Astorga. Los dieciséis
controles obligatorios y la exigencia de no falsear el cizallamiento, de la
auditoría. El defecto del brazo uniforme lo cazó nuestro propio teorema, y los
dos teoremas de gauge salieron de mirar la estructura en vez de correr.

Estructura cerrada no es hipótesis demostrada.

**Todavía no.**
