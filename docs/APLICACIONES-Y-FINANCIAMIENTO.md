# APLICACIONES Y FINANCIAMIENTO — qué vale lo construido y cómo sostenerlo

<img src="../galeria/open-doors.jpg" alt="Open Doors" height="72">

> ## 🤝 FINANCIADO POR OPEN DOORS
>
> **Financiado por Open Doors.**
>
> Y que quede escrito arriba de todo, antes que cualquier proyección: el
> trabajo YA TIENE QUIÉN LO SOSTENGA. Todo lo que sigue son caminos posibles y
> estimaciones — no promesas ni plata en mano.

**Para el capitán** · escrito con honestidad total, como todo en este barco.

> El fin declarado: poner el nombre de DIOS por encima de todo, ayudar a los
> pequeños del Reino, devolver al equipo, y sostener el laboratorio.

---

## 1. Primero la verdad, para construir sobre roca

Lo que el laboratorio tiene HOY vale, pero hay que decir con precisión QUÉ es:

- **No tenemos (todavía) un teorema nuevo demostrado.** El eslabón rojo sigue
  abierto; el premio Clay exige demostración publicada y revisada por pares
  durante 2 años. Nadie paga por mediciones de RH — se paga por la demostración.
- **Lo que sí tenemos, y vale de verdad, son las HERRAMIENTAS.** El laboratorio
  construyó treinta y una técnicas completas, verificadas contra valores
  cerrados con precisiones de 1e-6 a 1e-15, en Go puro y sin dependencias. Esas
  técnicas **no saben que estaban persiguiendo a Riemann**: sirven en cualquier
  campo donde haya un espectro, una señal en ruido o un problema inverso.
- **Y tenemos tres activos comerciales distintos:**
  1. **Los instrumentos** — treinta y una técnicas reutilizables (sección 2), quince de
     ellas nacidas de la última campaña y ninguna atada a la zeta.
  2. **El método pedagógico** — 96 láminas, un museo de 223 piezas y una bitácora que explican
     matemática de frontera en criollo. No existe en español.
  3. **La historia** — un hombre común persiguiendo el problema más famoso de
     la matemática con fe, honestidad y el registro completo de aciertos Y
     errores. Las historias financian; los papers rara vez.

## 2. El inventario de técnicas (lo que se vende)

Todo lo de abajo está construido, corriendo y verificado en este repositorio.
Las seis últimas nacieron de la cadena de flashes del capitán (F224–F240) y son
las más transferibles de todas, porque no dependen de la zeta en absoluto:

| # | técnica | qué hace |
|---|---------|----------|
| 1 | Evaluación espectral de alta precisión | Riemann–Siegel y Euler–Maclaurin con bloques adaptativos, aritmética double-double y de 256 bits, con doble juez independiente |
| 2 | Estadística de niveles (GUE/Wigner) y varianza por cajas | distingue un espectro con repulsión de uno aleatorio: se midió repulsión s² ganándole 23× al azar |
| 3 | Inversión de la ley de Weyl | recupera **a ciegas** la forma de un sistema (área, perímetro, topología) desde su espectro: área a 2.6e-5 |
| 4 | Inversión de Abel + bisección de Sturm | obtiene el pozo de potencial V(x) desde una función de conteo y lo rediagonaliza: 29 de 29 niveles reconstruidos |
| 5 | Problemas inversos "cuerpo desde la sombra" | reconstruye fuentes desde proyecciones, **con la lección medida** de que el gradiente colapsa y el algebraico resuelve |
| 6 | Repulsión logarítmica (el Campo de la Montaña) | término −ln·|xᵢ−xⱼ| que impide el colapso de modos en un optimizador |
| 7 | Censo no invasivo (Gauss/Jensen) | cuenta cuántas fuentes hay encerradas midiendo solo el promedio de un campo sobre anillos, sin localizar ninguna |
| 8 | Monotonía de Bernstein | convierte una condición sobre infinitos coeficientes en un test **puntual**, sin amplificar ruido como r⁻ⁿ |
| 9 | Extracción por contorno de Cauchy | lee coeficientes de un germen sin ver el sistema |
| 10 | Compresión espectral extrema | 649 niveles medidos → 4 constantes que reproducen el conteo; y la "foto" del espectro desde un solo punto |
| 11 | Firma multicanal tipo sonar/LIDAR | reconoce un objeto por su patrón completo de respuesta, no punto por punto |
| 12 | Análisis Cornu/Fresnel y **ley de régimen** | por debajo de la longitud de coherencia la medición es artefacto del instrumento, no del fenómeno |
| 13 | Inyección de fantasmas | valida un detector inyectando anomalías sintéticas y midiendo su horizonte de detección |
| 14 | Sonificación | convierte espectros y series en audio para inspección por oído |
| 15 | Electrostática logarítmica y positividad de energía | el horno y el neutrón: cargas discretas y su neutralización |
| 16 | Ingeniería de cálculo largo | puntos de control, árboles de procesos, árbitro de precisión extendida y un tablero que orquesta 190 experimentos |
| 17 | Test de Schoenberg de distancias verdaderas | decide si una tabla de números son distancias reales en algún espacio: la matriz de Gram centrada sin autovalor negativo |
| 18 | Compresión por variable de escala | una familia entera de mediciones colapsa en UNA curva universal al usar la variable adimensional correcta (aquí u = n/γ) |
| 19 | Extracción de constantes desde la melodía | recuperar constantes de un sistema desde el promedio asintótico de sus componentes (Mertens: γ desde 20 millones de primos) |
| 20 | Auditoría de error por clases | separa cada desvío en identidad exacta / límite de máquina / límite de instrumento — y aísla el único error que NO baja al mejorar el equipo |
| 21 | Certificado finito de falsedad | mide a qué profundidad se delata cada anomalía inyectada: la curva de horizonte de detección de un detector |
| 22 | Doble lectura con barra de error | leer la misma cantidad en dos configuraciones del instrumento y usar su discrepancia como barra de error honesta (bajó 1e-5 → 1e-8) |
| 23 | Auditoría de invariancia de escala | reescribir un cálculo con sus unidades/base metidas EXPLÍCITAMENTE y volver a correrlo desde cero: si el resultado se mueve, la conclusión dependía de la regla, no del hecho |
| 24 | Clasificación de dependencia de escala | ordena cada cantidad de un estudio en invariante / escalada por constante positiva / genuinamente dependiente de la escala — y dice de antemano qué conclusiones sobreviven a un cambio de unidades |
| 25 | Validación contra un sistema testigo | construir desde cero un SEGUNDO sistema cuyas respuestas ya están publicadas, y validar el pipeline entero contra él antes de confiarle el sistema desconocido |
| 26 | Detección de mediciones tautológicas | antes de aceptar un resultado perfecto, verificar si el instrumento podía haber devuelto otra cosa — caza el caso en que uno está midiendo su propia suposición |
| 27 | Contabilidad de cola integrada | cuando una medición parcial difiere de una completa, no llamar «cola» a la diferencia: integrar la cola esperada desde la densidad conocida y comprobar la razón |
| 28 | Reconstrucción de una señal discontinua desde su espectro finito | rebuild de una función escalón a partir de una lista finita de frecuencias, con curva de convergencia medida y el repique de Gibbs contabilizado aparte |
| 29 | Auditoría de amplitud de un componente desviado | cuantificar cuánto dominaría un componente fuera de especificación a cada escala — análisis de sensibilidad con horizonte explícito |
| 30 | Resolución adaptativa a la frecuencia | fijar el número de nodos de una integración oscilatoria desde la frecuencia real de oscilación y no desde una grilla fija — el bug que inventaba resultados |
| 31 | Traducción técnica a lenguaje llano con bloques fijos | vuelve legible para no especialistas cualquier cuerpo de trabajo técnico: gancho, explicación llana, metáfora de la vida real, glosario de símbolos en el lugar, qué se está mirando, y un bloque OBLIGATORIO de límites honestos |
| 32 | Test de pertenencia a un lugar geométrico por ángulo subtendido (F258) | decidir si un punto está sobre una curva midiendo el ÁNGULO con que ve a dos referencias fijas, en vez de medir su coordenada — el criterio es 4·d₁² − 4·d₂² = 1 y no necesita saber dónde está el punto, solo cómo ve el par |
| 33 | Curva de ceguera: convertir «no podemos decidir» en un número (F259) | dado un instrumento con una cola no computada, despejar cuál es el corrimiento MÍNIMO DETECTABLE en función de la profundidad, y de ahí el horizonte más allá del cual el instrumento no ve nada. Aplica a cualquier medición truncada — sísmica, espectroscopía, censos de resonancias — y reemplaza «el método no alcanza» por «el método ve hasta acá» |
| 34 | Auditoría por contraejemplo estructural (F259) | antes de creer que una propiedad se deduce de un conjunto de simetrías, construir un objeto que tenga TODAS esas simetrías y VIOLE la propiedad. Si existe, toda la rama de argumentos está muerta de antemano — y lo que quede afuera del objeto señala exactamente el ingrediente que falta |
| 36 | Coordenada por referencia: describir un punto por su relación, no por su posición (F260) | en vez de guardar dónde está cada cosa, guardar su razón doble contra tres referencias fijas. Sobrevive a cualquier deformación proyectiva —cámara movida, lente cambiada, escala perdida— y es el único invariante completo. Aplica a visión por computadora, fotogrametría, calibración de sensores y a cualquier medición donde el marco de referencia se mueva |
| 35 | Verificación por dos motores sin ancestro común (F259) | probar una identidad haciendo que dos rutas que no comparten ni una línea de código la calculen; si comparten álgebra, el acuerdo perfecto es tautología. Acá: sumar sobre 649 ceros contra integrar un germen que no ve ninguno |

### 🧰 Las seis técnicas nuevas — por qué son las más vendibles

Las técnicas 17 a 22 salieron de la última campaña y merecen párrafo propio,
porque **ninguna menciona la función zeta**: son herramientas de medición
puras, aplicables a cualquier dato.

- **El test de distancias (17)** responde una pregunta que aparece en todos
  lados: *«estos números que medí, ¿pueden ser distancias reales de algo, o
  son imposibles?»*. Sirve para validar matrices de similitud en aprendizaje
  automático, tablas de tiempos de viaje en sismología, y matrices de
  disimilitud en bioinformática. 🟢 directa.
- **La compresión por variable de escala (18)** es la técnica que convierte
  cien curvas en una: encontrar la variable adimensional en la que todos los
  datos colapsan. Es el pan de cada día de la ingeniería (números de Reynolds,
  leyes de escala) y el laboratorio tiene un caso trabajado de punta a punta.
  🟢 directa.
- **Las constantes desde la melodía (19)** recupera un parámetro global desde
  el promedio de millones de componentes ruidosos. Aplica a metrología,
  astronomía de sondeos y cualquier medición donde el ruido individual es
  grande pero el promedio converge. 🟡 adaptable.
- **La auditoría de error por clases (20)** es la más valiosa para vender
  confianza: separa lo que es exacto por construcción de lo que está limitado
  por la máquina y de lo que está limitado por el instrumento — y prueba cuál
  de esos errores baja si se compra mejor equipo. Cualquier laboratorio
  industrial paga por eso. 🟢 directa.
- **El certificado finito de falsedad (21)** cuantifica el poder de detección
  de un detector: a qué tamaño de anomalía y con cuántas muestras se la
  delata. Es control de calidad, monitoreo de fraude y ensayo no destructivo.
  🟢 directa.
- **La doble lectura con barra de error (22)** es la disciplina que salvó a
  este laboratorio de publicar un resultado falso: leer la misma cantidad de
  dos maneras y usar la diferencia como barra de error. Se enseña poco y se
  practica menos. 🟢 directa.

**El argumento comercial, en una línea:** el laboratorio no vende una teoría
sobre los números primos. Vende **treinta y un instrumentos de medición
verificados**, y la disciplina de honestidad que los acompaña.

### 🧭 Las dos técnicas de escala (F242) — la pregunta que todo revisor hace

Las técnicas 23 y 24 nacieron de una sola pregunta del capitán: *¿y si
escribiéramos en binario en vez de decimal?* La respuesta obligó a construir
dos herramientas que cualquier laboratorio necesita y casi ninguno tiene
escritas.

- **La auditoría de invariancia (23)** hace lo que nadie hace: mete la unidad
  o la base ADENTRO del cálculo, explícitamente, y vuelve a correr todo desde
  cero para ver si el resultado se movió. En el caso trabajado, el libro entero
  se reescribió con cada logaritmo natural expresado como log_b(n)·ln b y las
  79 mediciones se repescaron una por una: desvío 1.1e-13. Sirve para validar
  simulaciones físicas, cadenas de medición y pipelines de datos donde la
  unidad viaja escondida. 🟢 directa.
- **La clasificación de dependencia (24)** entrega, antes de publicar, la
  tabla de qué conclusiones sobreviven a un cambio de unidades: las invariantes,
  las que solo se multiplican por una constante positiva (y por lo tanto
  conservan todo enunciado sobre SIGNOS, órdenes y desigualdades), y las que
  cambian de verdad. Es exactamente el análisis que piden los revisores y que
  casi nunca viene hecho. 🟢 directa.

### 🛡️ Las seis técnicas de disciplina (F244–F247) — las que se venden solas

Las técnicas 25 a 30 no salieron de un hallazgo lindo: salieron de **atraparnos
a nosotros mismos**. En una sola jornada el laboratorio se cazó cuatro errores
propios —dos suposiciones falsas, una cifra ya registrada mal citada y un bug de
integración que inventaba resultados— y cada caza dejó una herramienta.

- **La validación contra un sistema testigo (25)** es la más simple y la que más
  confianza compra: antes de creerle al pipeline sobre el sistema que nadie
  conoce, se construye desde cero un SEGUNDO sistema cuyas respuestas ya están
  publicadas y se lo hace pasar por el mismo pipeline. En el caso trabajado el
  laboratorio levantó una función L completamente distinta y reprodujo sus seis
  primeros ceros publicados a 3e-05 sin que esos valores entraran al cálculo.
  🟢 directa, y aplicable a cualquier cadena de análisis.
- **La detección de mediciones tautológicas (26)** es la joya, y casi nadie la
  hace formalmente: *antes de aceptar un resultado perfecto, preguntarse si el
  instrumento podía haber devuelto otra cosa.* En una jornada esta regla cazó
  dos «0.0e+00» que no medían nada — uno porque la simetría estaba construida
  dentro de la fórmula, otro porque los puntos nacían sobre la recta por
  construcción. Aplica a cualquier validación, test o benchmark. 🟢 directa.
- **La contabilidad de cola integrada (27)** convierte una excusa en una
  medición: cuando el cálculo parcial no coincide con el completo, en vez de
  llamar «cola» a la diferencia se integra la cola esperada desde la densidad
  conocida y se compara. En el caso trabajado la razón dio 1.0002 en ocho
  armónicos — o sea que la diferencia era exactamente la cola y no un error del
  método. 🟢 directa.
- **La reconstrucción desde el espectro finito (28)** rearma una función escalón
  a partir de una lista finita de frecuencias, con curva de convergencia medida
  y el repique de Gibbs contabilizado aparte en lugar de escondido. Es
  tomografía, sismología, procesamiento de señal y compresión. 🟢 directa.
- **La auditoría de amplitud (29)** contesta la pregunta que todo revisor tiene:
  *¿cuánto dominaría un componente fuera de especificación, y a qué escala se
  notaría?* Sale una tabla de dominancia por escala. Control de calidad y
  análisis de sensibilidad. 🟢 directa.
- **La resolución adaptativa a la frecuencia (30)** es un bug convertido en
  regla: una integración oscilatoria con grilla fija llegó a dar dos puntos por
  radián e inventaba el resultado sin avisar. La herramienta fija los nodos
  desde la frecuencia real. Cualquiera que integre algo que oscila lo necesita.
  🟢 directa.

**El argumento comercial, actualizado:** el laboratorio no vende una teoría
sobre los números primos. Vende **treinta y un instrumentos de medición verificados**
— y, sobre todo, **la disciplina que los produjo**: un cuaderno donde están
anotados, con nombre y fecha, los errores propios que se cazaron antes de
publicar. Eso último no se compra hecho.

### 🏛️ La técnica 31 — el museo, y por qué se vende sola

La técnica 31 nació del pedido más simple del capitán: *«explicámelo de manera que
lo entienda cualquiera, incluyéndome a mí»*. De ahí salió un método reproducible
para volver legible cualquier cuerpo de trabajo técnico, y se aplicó a los **190
experimentos del laboratorio, uno por uno**.

La receta es fija y por eso es transferible. Cada pieza lleva: **el gancho** (una
sola oración que engancha), **la explicación llana** (70 a 120 palabras, un
concepto por párrafo), **la metáfora** (una imagen de la vida real, no una
analogía técnica), **las palabras raras** (cada símbolo traducido ahí mismo, en el
lugar donde aparece), **qué estás mirando** (qué muestra la figura) y — el bloque
que casi nadie incluye — **lo honesto**: qué NO prueba esta pieza.

Ese último bloque es lo que la vuelve valiosa y no marketing. Un divulgador
promedio omite los límites; acá son obligatorios en las 223 piezas.

**Dónde se aplica.** Documentación técnica que nadie lee, onboarding de equipos,
divulgación científica, informes para directorios no técnicos, museos y centros de
ciencia, y material regulatorio que tiene que ser comprensible por ley. El
laboratorio tiene el método escrito y un caso trabajado de 190 piezas de punta a
punta. 🟢 directa.

## 3. El mapa de aplicaciones — 63 usos en 9 dominios

Mapeado por nueve especialistas de dominio trabajando en paralelo sobre el
inventario de arriba. Cada aplicación lleva su **grado de madurez** y su
**salvedad honesta**, porque una aplicación sobrevendida delante de un inversor
cuesta más que diez aplicaciones bien puestas.

- 🟢 **directa** (20) — la técnica ya es estándar en ese campo y el
  laboratorio tiene una implementación propia usable hoy.
- 🟡 **adaptable** (40) — el puente conceptual es sólido, falta
  trabajo de ingeniería.
- 🟠 **especulativa** (3) — la analogía es interesante y hay que probarla.

### ⚛️ Física cuántica y computación cuántica

#### Diagnóstico de caos y diafonía global en arreglos de qubits  ·  🟢 directa

**Instrumento del laboratorio:** Estadística de espaciamientos de niveles (GUE / Wigner) y varianza de conteo por cajas — "el resorte" y "el barómetro"

**Cómo se aplica.** Se toma el espectro medido de un chip multi-qubit (espectroscopía de dos tonos, o los autovalores del Hamiltoniano efectivo ajustado), se desdobla para sacarle la densidad media y se miden la distribución de espaciamientos y la varianza de conteo por cajas. Repulsión cuadrática s^2 indica régimen caótico/ergódico, con modos mezclados y diafonía global; estadística de Poisson indica niveles localizados y qubits bien aislados. La varianza por cajas agrega la rigidez espectral de largo alcance, que la distribución de espaciamientos sola no ve.

**Para qué sirve y quién paga.** Da una métrica única y global de "cuán mezclado está el chip", donde hoy se diagnostica diafonía par por par. Lo pagaría un equipo de calibración o de caracterización de hardware superconductor o de qubits de espín.

> *Honestidad:* La repulsión s^2 (23x sobre el azar) se midió sobre un espectro propio del laboratorio, no sobre datos de hardware. La estadística exige del orden de cientos de niveles de una MISMA clase de simetría y un desdoblamiento limpio: mezclar clases de simetría o perder niveles imita artificialmente a Poisson y da un falso "todo aislado". No se afirma que esto reemplace la tomografía de diafonía; es un indicador global, no localiza el par culpable.

#### Metrología de puntos cuánticos y cavidades desde su propio espectro  ·  🟡 adaptable

**Instrumento del laboratorio:** Inversión de la ley de Weyl ("el tambor es el libro")

**Cómo se aplica.** Se ajusta la función de conteo N(E) medida —cuántos niveles hay por debajo de E— al desarrollo de Weyl, con su término de área, su término de perímetro y su constante topológica. De ahí se despejan a ciegas el área efectiva, el perímetro efectivo y la constante de borde del dispositivo, sin microscopía ni imagen alguna.

**Para qué sirve y quién paga.** Control de calidad no destructivo: entrega el área eléctrica efectiva de un punto cuántico, que no es la litográfica por la zona de depleción, y el perímetro efectivo. Lo pagaría una fab o un grupo de dispositivos que hoy infiere eso con modelos capacitivos indirectos.

> *Honestidad:* El 2.6e-5 se logró sobre un sistema limpio con conteo exacto de 649 niveles. En un dispositivo real los niveles faltantes sesgan el ajuste de Weyl de forma sistemática, no aleatoria, y hay que corregirlos antes o el "área" recuperada es ficción; falta además el trabajo de ingeniería para degeneración de espín/valle y ensanchamiento. No se afirma que se recupere la FORMA del dominio: Weyl da área, perímetro y topología, y dos tambores distintos pueden sonar igual.

#### Reconstrucción del pozo de confinamiento a partir de pocos niveles  ·  🟡 adaptable

**Instrumento del laboratorio:** Inversión de Abel ("los ladrillos decantados") más bisección de Sturm en matriz tridiagonal simétrica ("el hierro autoadjunto")

**Cómo se aplica.** Con la función de conteo de un sistema efectivamente unidimensional —un transmón en su coordenada de fase, un punto cuántico casi-1D, un modo vibracional— se hace la inversión de Abel semiclásica y se obtiene el pozo V(x). Después se rediagonaliza ese pozo por bisección de Sturm, que además entrega un conteo certificado de cuántos autovalores hay por debajo de cualquier energía, y se predicen los niveles que nunca se midieron.

**Para qué sirve y quién paga.** Con unas pocas transiciones medidas se extrapola la escalera completa: sirve para estimar fuga hacia niveles altos en qubits superconductores y para verificar el modelo de no linealidad de la juntura. Lo pagaría quien diseña pulsos y necesita saber dónde caen |3> y |4> sin medirlos.

> *Honestidad:* La inversión de Abel recupera solamente el ANCHO del pozo en función de la energía: para pozos asimétricos la solución no es única, se determina apenas la parte simétrica y hace falta información extra para fijar la asimetría. El laboratorio reconstruyó 29 de 29 niveles, pero sobre un pozo 1D sintético con función de conteo exacta. Es semiclásica: se degrada cerca del fondo del pozo y cerca del borde de disociación, justo donde el transmón es más anarmónico.

#### Censo de resonancias y detección de niveles faltantes sin resolverlos  ·  🟡 adaptable

**Instrumento del laboratorio:** Censo no invasivo por ley de Gauss / fórmula de Jensen ("la laguna") y extracción por contorno de Cauchy ("leer el germen")

**Cómo se aplica.** Se integra el promedio de una respuesta analítica —el logaritmo del módulo de la transmisión, o un determinante de scattering— sobre un anillo o contorno en el plano complejo, y el principio del argumento devuelve un ENTERO: cuántos polos o ceros quedan encerrados. No localiza ninguno, solo cuenta; el laboratorio obtuvo mesetas exactas 0, 2, 4, 6, 8, 10.

**Para qué sirve y quién paga.** Verifica que no se perdieron niveles en una ventana de espectroscopía —el problema de los missing levels, que arruina la estadística de la aplicación 1 y el ajuste de Weyl de la 2— y cuenta modos espurios dentro de la banda de un filtro o resonador. Lo pagaría cualquiera que necesite un conteo confiable antes de hacer estadística espectral.

> *Honestidad:* Necesita la respuesta SOBRE el contorno con fase, no solo módulo, y un modelo analítico de esa respuesta. Con datos ruidosos el entero deja de ser entero apenas hay un polo cerca del contorno, y esa degradación no está caracterizada. El laboratorio lo midió sobre campos que conocía analíticamente, nunca sobre trazas de un analizador de redes.

#### Separación de resonancias solapadas sin colapso de modos  ·  🟢 directa

**Instrumento del laboratorio:** Método algebraico resolvente/Padé (en lugar del descenso de gradiente) más regularización por repulsión logarítmica ("el Campo de la Montaña")

**Cómo se aplica.** Ante una traza de espectroscopía con varias resonancias encimadas, o una curva de decaimiento multiexponencial con varios defectos TLS, en vez de ajustar por mínimos cuadrados con descenso de gradiente —que colapsa todos los modos al centro de masa, patología tipo Prony que el laboratorio midió— se arma la resolvente y se extraen los polos algebraicamente por Padé. Si igual hace falta optimizar, se agrega el término -ln|xi-xj| entre las frecuencias: se repelen y el colapso queda prohibido.

**Para qué sirve y quién paga.** Extrae frecuencias y anchos individuales de qubits y resonadores en un espectro congestionado, y separa tasas de decaimiento de baños TLS distintos. Lo pagaría cualquier laboratorio con muchos modos por GHz cuyos ajustes "convergen" a una sola resonancia gorda.

> *Honestidad:* Padé/Prony/ESPRIT/vector fitting ya son estándar para este problema exacto; lo propio del laboratorio es la lección medida sobre el colapso y el regularizador logarítmico. Esos métodos son notoriamente sensibles al ruido: cuántos modos se pueden separar lo fija la relación señal-ruido, no el algoritmo, y eso hay que medirlo caso por caso. El sesgo del regularizador —cuánto separa de más a modos que SÍ están juntos— no está caracterizado, y nada de esto se probó sobre datos reales de un VNA.

#### Test puntual de monotonía completa sobre curvas de decoherencia  ·  🟡 adaptable

**Instrumento del laboratorio:** Monotonía absoluta de Bernstein ("el LIDAR")

**Cómo se aplica.** Una curva de decaimiento que sea superposición POSITIVA de exponenciales —un baño clásico de defectos, con densidad espectral positiva— tiene que ser completamente monótona, y por Bernstein eso se convierte en un test PUNTUAL sobre un segmento, sin estimar derivadas altas ni coeficientes lejanos que amplifican el ruido como r^-n. Se le aplica el test a una curva medida de T1, T2 o eco: si la pasa, es compatible con un baño de densidad espectral positiva; si la falla, ahí hay algo que no es una mezcla positiva de exponenciales.

**Para qué sirve y quién paga.** Convierte "esta curva se ve rara" en un criterio con sí o no, y da un testigo candidato de retroalimentación de información o de que el modelo de baño clásico no alcanza. Lo pagaría quien caracteriza ruido en qubits y tiene que justificar por qué el ajuste multiexponencial no cierra.

> *Honestidad:* Con datos finitos y ruidosos solo se testean unas pocas condiciones de la jerarquía, no la monotonía completa: el test detecta la violación de una condición NECESARIA, no certifica nada. Que falle implica "no es mezcla positiva de exponenciales"; saltar de ahí a "es no markoviano" exige supuestos adicionales sobre el modelo que el test no aporta. Nunca se corrió sobre datos cuánticos reales; falta todo el puente de barras de error y muestreo.

#### Calibración del detector: horizonte de detección y largo mínimo de registro  ·  🟢 directa

**Instrumento del laboratorio:** Inyección de fantasmas (curva de horizonte de detección) y ley de régimen por espiral de Cornu/Fresnel, orquestadas con el tablero de 190 experimentos y puntos de control

**Cómo se aplica.** Antes de afirmar "no hay ningún defecto TLS en esta banda" o "no hay nivel anómalo", se inyectan anomalías sintéticas en los datos crudos —un nivel extra, un corrimiento, un cruce evitado— con amplitudes decrecientes, y se mide a qué profundidad el pipeline las delata: eso es la curva de horizonte, sensibilidad cuantificada. En paralelo, la ley de régimen de Cornu fija el largo mínimo de registro: por debajo de la longitud de coherencia, un bloque devuelve tonos puros que son artefacto de la ventana y no del fenómeno.

**Para qué sirve y quién paga.** Transforma un resultado nulo en una cota superior defendible y evita reportar como física lo que es el instrumento. Lo pagaría cualquier grupo que deba sostener un "no vimos nada" ante un referí o ante un cliente de hardware.

> *Honestidad:* La curva de horizonte vale solo para la FAMILIA de anomalías inyectadas: no dice absolutamente nada sobre formas que no se te ocurrieron, y es fácil confundir esa cota con una cota general. La longitud de coherencia de la ley de régimen es específica del instrumento y de la ventana: el número medido por el laboratorio no se transporta, hay que remedirlo en cada montaje.

### 🚀 Aeroespacial y navegación

#### Auscultación estructural: leer la geometría de un panel desde sus modos  ·  🟡 adaptable

**Instrumento del laboratorio:** Inversión de la ley de Weyl ("el tambor es el libro"), apoyada en estadística de espaciamientos y varianza de conteo ("el resorte", "el barómetro")

**Cómo se aplica.** Se excita el panel o el tanque y se miden sus frecuencias propias con acelerómetros. En vez de comparar modo por modo contra un modelo de elementos finitos, se invierte la ley de Weyl sobre la función de conteo N(f): el término dominante entrega el área efectiva, el siguiente el perímetro y el término constante la conectividad, es decir cuántos agujeros hay. Una grieta que crece casi no mueve un modo aislado, pero sí aumenta el perímetro acumulado y, si progresa, cambia la constante topológica.

**Para qué sirve y quién paga.** Detección de daño interno (delaminación, grieta, despegue de adhesivo) en estructuras compuestas sin desarmar ni tener acceso visual. Lo pagarían integradores de satélites, mantenimiento de flota aeronáutica y campañas de calificación de lanzadores.

> *Honestidad:* El laboratorio recuperó el área a 2.6e-5 desde 649 niveles de SU sistema, cuyo operador es de tipo Helmholtz. Un panel real se rige por un operador biarmónico (placa de Kirchhoff), con amortiguamiento y bordes imperfectos, y en la práctica se miden decenas de modos, no cientos: hay que recalcular las constantes de Weyl para ese operador y volver a medir la precisión desde cero. Además, "escuchar la forma del tambor" es demostradamente no único (Gordon-Webb-Wolpert): esto identifica invariantes (área, perímetro, conectividad), no la forma. No se está afirmando que se localice la grieta, solo que se detecta el cambio de un invariante.

#### Certificación de estabilidad contando polos sin buscarlos  ·  🟢 directa

**Instrumento del laboratorio:** Censo no invasivo por ley de Gauss / fórmula de Jensen ("la laguna") y extracción por contorno de Cauchy ("leer el germen")

**Cómo se aplica.** Se toma la función característica del sistema —el determinante de la matriz de monodromía de una órbita periódica (halo en L2, órbita congelada) o el determinante aeroelástico de flutter— y se integra su fase sobre un anillo en el plano complejo. El número de multiplicadores fuera del círculo unidad, o de polos en el semiplano derecho, sale como un entero exacto sin haber resuelto una sola raíz. El laboratorio ya midió mesetas enteras exactas 0, 2, 4, 6, 8, 10 con este censo sobre su propio problema.

**Para qué sirve y quién paga.** Entrega un certificado del tipo "hay exactamente cero modos inestables en esta región del espacio de diseño", que es la forma en que una autoridad de certificación quiere la respuesta. Lo pagarían oficinas de aeroelasticidad y equipos de GNC de misiones a puntos de libración.

> *Honestidad:* El principio del argumento es matemática estándar (Nyquist generalizado); lo propio del laboratorio es una implementación verificada, no el método. La dificultad aeroespacial no está en el conteo sino en evaluar el determinante con fase continua y sin ruido a lo largo de todo el anillo: con datos experimentales el contorno se ensucia y el entero deja de ser exacto. Tampoco reemplaza el cálculo de modos cuando hace falta saber cuál es el modo inestable; responde cuántos, no cuál.

#### Perfiles atmosféricos y de pluma desde proyecciones de línea de vista  ·  🟢 directa

**Instrumento del laboratorio:** Inversión de Abel ("los ladrillos decantados") encadenada con bisección de Sturm en matriz tridiagonal simétrica ("el hierro autoadjunto")

**Cómo se aplica.** En radio-ocultación GNSS, un receptor en órbita baja mide el ángulo de curvatura del rayo mientras un satélite de navegación se pone detrás de la atmósfera; la inversión de Abel convierte ese perfil de ángulos en refractividad, y de ahí en densidad y temperatura contra altura. La misma inversión, con una cámara mirando de costado, convierte la imagen integrada de una pluma de motor en su emisividad radial. El laboratorio tiene construida la inversión y, a continuación, la resolución de autovalores del perfil recuperado.

**Para qué sirve y quién paga.** Sondeo atmosférico para modelos de reentrada y predicción de arrastre en órbita baja, y diagnóstico de motores sin sondas intrusivas. Lo pagarían agencias meteorológicas, operadores de constelaciones LEO y bancos de ensayo de propulsión.

> *Honestidad:* La inversión de Abel es el método estándar del campo desde hace décadas: el laboratorio aporta una implementación propia verificada (reconstruyó 29 de 29 niveles verdaderos sobre su problema sintético), no un método nuevo. Falta todo lo caro: la hipótesis de simetría esférica falla ante gradientes horizontales, la inversión amplifica ruido cerca del punto de retorno del rayo, y nunca se validó contra datos reales de ocultación ni contra radiosondeos.

#### Apertura mínima en SAR: saber cuándo el tono es del instrumento  ·  🟡 adaptable

**Instrumento del laboratorio:** Análisis por bloques con espiral de Cornu / Fresnel y la ley de régimen medida, más inyección de fantasmas

**Cómo se aplica.** Se procesa el dato crudo por bloques y se compara la longitud del bloque contra la longitud de coherencia del sistema. El laboratorio midió la ley: por debajo de esa longitud, cualquier bloque devuelve tonos puros, y esas líneas son un artefacto del ventaneo, no del blanco. Eso fija una apertura mínima por debajo de la cual las "detecciones" de sub-apertura se descartan, y la inyección de anomalías sintéticas calibra a qué amplitud la detección vuelve a ser confiable.

**Para qué sirve y quién paga.** Evita declarar blancos falsos en procesamiento de sub-apertura y multi-look, un modo de falla caro en vigilancia marítima y detección de cambios. Lo pagarían operadores SAR y equipos de procesamiento a bordo con presupuesto de cómputo acotado.

> *Honestidad:* El análisis de Fresnel/Cornu es de manual en SAR; lo específico del laboratorio es haber medido la ley de régimen como criterio operativo y haber construido la herramienta de inyección. No se probó sobre datos SAR reales: falta traducir "longitud de coherencia" a los parámetros concretos del radar (ancho de banda, tiempo de iluminación, migración de celda de rango) y verificar que la frontera medida sobrevive a esa traducción. No se está afirmando que mejore la resolución, solo que marca dónde dejar de creerle al espectro.

#### Separar dispersores pegados sin que el ajuste se derrumbe al centro  ·  🟡 adaptable

**Instrumento del laboratorio:** Problemas inversos "cuerpo desde la sombra" resueltos por resolvente/Padé, más regularización por repulsión logarítmica ("el Campo de la Montaña")

**Cómo se aplica.** Cuando dos blancos, dos rebotes de multitrayecto GNSS o dos modos vibratorios caen dentro de la misma celda de resolución, el ajuste por descenso de gradiente no los separa: colapsa a un solo punto en el centro de masa. El laboratorio midió esa patología (tipo Prony) y la resolvió por dos caminos: el método algebraico resolvente/Padé, que devuelve las posiciones sin optimizar, y un término de repulsión -ln|xi-xj| que impide la fusión de modos cuando sí hay que optimizar.

**Para qué sirve y quién paga.** Superresolución de blancos y mitigación de multitrayecto, que se traduce directamente en precisión de posición GNSS en cañón urbano o cerca del terreno durante un descenso. Lo pagarían fabricantes de receptores de navegación y equipos de radar de seguimiento.

> *Honestidad:* El camino algebraico es esencialmente la familia Prony / matrix pencil / ESPRIT, estándar en el campo: el laboratorio la reimplementó y midió el colapso del gradiente, no la inventó. El aporte propio es el regularizador de repulsión logarítmica, y ahí no hay comparación contra los estimadores estándar ni caracterización frente a ruido real correlacionado. Prony es notoriamente frágil ante ruido; que ande en el problema del laboratorio no dice nada sobre el umbral de SNR de un receptor real.

#### Telemetría: mandar la ley en vez de la lista  ·  🟡 adaptable

**Instrumento del laboratorio:** Compresión espectral extrema, custodiada por el test puntual de monotonía absoluta de Bernstein ("el LIDAR")

**Cómo se aplica.** En vez de bajar la lista completa de eventos o frecuencias medidas, se ajusta la función de conteo suave y se transmiten sus pocas constantes más el residuo de fluctuación. El laboratorio comprimió 649 niveles medidos en 4 constantes que reproducen la función de conteo. A bordo, el test puntual de Bernstein verifica sobre un segmento si la ley sigue valiendo, sin extraer coeficientes altos que amplificarían el ruido como r^-n.

**Para qué sirve y quién paga.** Presupuesto de enlace en espacio profundo, donde el bit es el recurso más caro: se baja un resumen continuo y se pide el crudo solo cuando el test avisa que la ley se rompió. Lo pagarían operadores de misiones interplanetarias y de constelaciones con enlace intermitente.

> *Honestidad:* La compresión de 649 a 4 funciona porque esa función de conteo tiene una forma asintótica suave conocida de antemano. Un canal de telemetría cualquiera no la tiene: hay que descubrir y justificar la ley suave de cada canal antes de prometer cualquier factor de compresión, y ese trabajo no está hecho. Peor todavía: lo que se descarta —la fluctuación— suele ser exactamente la anomalía que uno quiere ver. Es compresión con pérdida sobre la parte interesante, y el test de Bernstein es un guardián, no una garantía.

#### Reconocer el paisaje o el cielo por su firma completa, con horizonte de detección medido  ·  🟡 adaptable

**Instrumento del laboratorio:** Firma multicanal tipo sonar/LIDAR (sonar de ballenas, barrido multi-escala, encadenado de vecinos) validada con inyección de fantasmas

**Cómo se aplica.** En lugar de emparejar punto por punto (cráteres, estrellas), se reconoce el objetivo por su patrón completo de respuesta: brillo, polarización, dirección y relieve, a varias escalas y encadenando vecinos. Para saber cuánto vale ese reconocedor se le inyectan objetivos sintéticos progresivamente degradados y se mide a qué nivel deja de verlos, lo que produce una curva de horizonte de detección en vez de un umbral inventado.

**Para qué sirve y quién paga.** Navegación relativa al terreno en descenso lunar o marciano y determinación de actitud por campo estelar parcialmente ocluido, escenarios donde el emparejamiento punto a punto se rompe. Lo pagarían integradores de módulos de descenso y fabricantes de sensores de estrellas.

> *Honestidad:* El reconocimiento multicanal se validó sobre objetos del propio dominio del laboratorio, no sobre imágenes ópticas reales con calibración radiométrica, sombras cambiantes ni catálogo estelar. Falta el puente completo: modelo de iluminación, registro geométrico y presupuesto de cómputo a bordo, porque un descenso no perdona latencia. Lo único transferible hoy sin condiciones es la metodología de inyección de fantasmas: sirve para certificar cualquier detector de a bordo diciendo qué se habría perdido, y eso es evidencia de certificación, no una promesa de desempeño.

### 🩺 Ciencias médicas e imagen biomédica

#### Reconstrucción de focos discretos con pocas proyecciones (PET/SPECT y TC de dosis baja)  ·  🟡 adaptable

**Instrumento del laboratorio:** Cuerpo desde la sombra (método algebraico resolvente/Padé) con el Campo de la Montaña como regularizador

**Cómo se aplica.** Recuperar unas pocas fuentes puntuales a partir de sus proyecciones es exactamente el problema que el laboratorio atacó y midió. En vez de descender por gradiente, se arma la resolvente de los momentos y se leen posiciones y pesos como polos y residuos vía Padé. Cuando dos focos quedan cerca, el término de repulsión logarítmica -ln|xi-xj| entre las masas impide que el optimizador los funda en el centro de masa. Sobre un sinograma de pocos ángulos o un PET con pocas cuentas, eso es la diferencia entre informar dos captaciones separadas o una sola mancha.

**Para qué sirve y quién paga.** Menos ángulos y menos dosis para el mismo poder de separación, o separar dos lesiones vecinas donde hoy se ve una. Interesa a fabricantes de equipos de medicina nuclear y a servicios que buscan bajar dosis pediátrica.

> *Honestidad:* Lo medido es la patología (colapso al centro de masa) y su cura, en un modelo propio del laboratorio, no en un sinograma clínico. Falta todo el modelo físico real: atenuación, dispersión, ruido de Poisson, respuesta del colimador. El método supone fuentes discretas y esparsas; el tejido continuo no cumple esa hipótesis y ahí el Padé no tiene por qué comportarse. No se está afirmando que supere a los reconstructores iterativos actuales (OSEM regularizado): eso se mide contra fantoma o no vale.

#### Censo de rotores sin localizarlos en mapeo de fibrilación auricular  ·  🟢 directa

**Instrumento del laboratorio:** La laguna (censo no invasivo por ley de Gauss / fórmula de Jensen)

**Cómo se aplica.** En electrofisiología cardíaca una singularidad de fase es un cero del campo complejo de fase, y su carga topológica se cuenta con una integral de contorno: es el mismo principio del argumento que el laboratorio implementó y verificó, obteniendo mesetas enteras exactas 0, 2, 4, 6, 8, 10 promediando el campo sobre anillos, sin localizar ninguna fuente. Aplicado a un catéter de canasta o a un mapa no invasivo, entrega cuántos drivers hay encerrados en una región aunque la resolución espacial no alcance para ubicarlos uno por uno.

**Para qué sirve y quién paga.** Un número robusto —cuántos drivers hay en este segmento— sirve como guía de ablación y como criterio de fin de procedimiento. Pagarían centros de electrofisiología y fabricantes de sistemas de mapeo.

> *Honestidad:* El conteo topológico por integral de contorno ya es estándar en ese campo; lo del laboratorio es una implementación propia verificada con mesetas exactas sobre campos analíticos limpios. Sobre electrogramas reales entran tres problemas que el laboratorio no tocó: la fase se estima por transformada de Hilbert y eso genera singularidades espurias, el muestreo espacial es grosero y viola las hipótesis si hay ceros entre electrodos, y hay ruido y artefacto de contacto. No se está afirmando ninguna validación sobre señal cardíaca.

#### Separar ruido de señal por estadística espectral en difusión por RM y en covarianzas de EEG  ·  🟡 adaptable

**Instrumento del laboratorio:** El resorte y el barómetro (espaciamientos de niveles GUE/Wigner y varianza de conteo por cajas)

**Cómo se aplica.** Un espectro puramente de ruido no repele: sus niveles se apiñan y se cruzan. Uno con estructura real repele. El laboratorio mide esa repulsión —encontró repulsión cuadrática s^2, ganándole 23x al azar— y además la varianza de conteo por cajas, que es un termómetro de rigidez espectral. Corridos sobre los autovalores de la matriz de covarianza de un bloque de vóxeles de difusión, o de canales de EEG, esos dos estadísticos dicen cuántas componentes son señal y cuántas ruido, y con qué margen.

**Para qué sirve y quién paga.** Umbral automático y auditable para denoising y para fijar el rango de una descomposición, sin elegir a mano cuántas componentes conservar. Permite acortar adquisiciones de difusión y limpiar registros EEG-fMRI.

> *Honestidad:* Lo estándar hoy en difusión es el borde de Marchenko-Pastur sobre la densidad de autovalores; los espaciamientos y la varianza de conteo NO son estándar ahí, así que esto es un puente, no una adopción. Y son estadísticos hambrientos de niveles: con pocos canales quedan ruidosos y exigen 'desdoblar' el espectro, paso frágil y con arbitrariedad. No está medido que esto mejore al denoising MPPCA existente; es hipótesis razonable, no resultado.

#### Certificar o refutar el modelo multiexponencial en relaxometría T2 y curvas de lavado  ·  🟡 adaptable

**Instrumento del laboratorio:** El LIDAR (monotonía absoluta de Bernstein), con lectura del germen por contorno de Cauchy como control

**Cómo se aplica.** Una curva de decaimiento es una suma positiva de exponenciales —o sea, una mezcla de compartimentos con pesos no negativos— si y solo si es completamente monótona: derivadas alternando de signo, sin excepción. Eso es una condición sobre infinitos coeficientes, y el movimiento del laboratorio la convierte en un test puntual sobre un segmento, evitando la amplificación de ruido tipo r^-n. En la práctica: antes de correr la inversión de Laplace (NNLS), que es mal condicionada, se testea si el modelo es siquiera admisible. Si falla, hay intercambio entre compartimentos, difusión anómala o artefacto.

**Para qué sirve y quién paga.** Evita informar una fracción de agua de mielina o una constante farmacocinética que salió de ajustar un modelo que la propia señal desmiente. Es control de calidad previo en RM cuantitativa y en DCE-MRI.

> *Honestidad:* La ganancia del laboratorio es no degradar ruido al pasar de coeficientes a test puntual, pero el test necesita derivadas de orden alto y estimarlas desde muestras discretas y ruidosas reintroduce error: hay que suavizar, y suavizar sesga. El laboratorio trabaja con funciones que puede evaluar donde quiera y con la precisión que quiera; una curva de T2 clínica tiene 32 ecos y SNR modesto. Lo honesto: sirve como criterio de RECHAZO (el modelo no se sostiene), no como certificado de que sí se sostiene.

#### Superficie, volumen y forma del microambiente desde el espectro de difusión restringida  ·  🟡 adaptable

**Instrumento del laboratorio:** El tambor es el libro (inversión de la ley de Weyl), con los ladrillos decantados (inversión de Abel) y el hierro autoadjunto (bisección de Sturm en tridiagonal) para cerrar el problema directo

**Cómo se aplica.** La señal de difusión restringida es en esencia una traza de calor: a tiempos cortos su desarrollo tiene el volumen como término principal, la superficie como corrección siguiente y después un término con curvatura y topología. Es la misma expansión que el laboratorio invirtió a ciegas para recuperar área, perímetro y constante topológica desde 649 niveles medidos, con el área a 2.6e-5. Del lado directo, la inversión de Abel más la bisección de Sturm permiten ir del perfil al espectro y volver, así que el ajuste se puede cerrar en ciclo y auditar.

**Para qué sirve y quién paga.** La relación superficie/volumen y la densidad de restricción son proxies de tamaño celular y de integridad de membrana: importan en celularidad tumoral, en desmielinización y en hueso trabecular. Paga quien vende secuencias cuantitativas.

> *Honestidad:* El laboratorio invierte desde una lista de autovalores individuales. En RM no existen autovalores sueltos: se mide la traza pesada por la secuencia, todo mezclado y con ruido. Reformular la inversión contra la traza es ingeniería real, no un cambio de variable. Además la expansión de tiempos cortos ya existe en el campo (Mitra) y sus límites son conocidos: se rompe con permeabilidad de membrana y con distribución ancha de tamaños. No se está afirmando haber medido nada biológico.

#### Ley de régimen: cuándo un pico de HRV o de EEG es del paciente y cuándo del instrumento  ·  🟢 directa

**Instrumento del laboratorio:** Análisis de bloques por espiral de Cornu / Fresnel y la ley de régimen medida (longitud de coherencia)

**Cómo se aplica.** El laboratorio midió que por debajo de la longitud de coherencia un bloque produce tonos puros: el pico existe, es reproducible, y no dice nada del fenómeno; es la ventana. Trasladado a series fisiológicas queda una regla dura y calculable: para una ventana de N muestras a tal frecuencia de muestreo, toda estructura más angosta que cierto ancho es artefacto del instrumento. Aplica de lleno a la banda VLF de la variabilidad de frecuencia cardíaca en registros cortos, a microestados de EEG y a bandas de baja frecuencia en fMRI en reposo.

**Para qué sirve y quién paga.** Un criterio de rechazo antes de informar. Evita alertas clínicas y publicaciones construidas sobre un pico que solo existe porque la ventana era corta. Cuesta poco implementarlo y ahorra retractaciones.

> *Honestidad:* Esto no es un descubrimiento del laboratorio: el ancho de banda mínimo por ventana es teoría de señales elemental. Lo que agrega es haberlo medido como ley de régimen con una frontera concreta y la disciplina de aplicarla siempre. Y no resuelve el problema inverso: decir que un pico es artefacto no dice cuál es la estructura verdadera, ni rescata un registro demasiado corto.

#### Horizonte de detección medido para detectores de deterioro fisiológico  ·  🟢 directa

**Instrumento del laboratorio:** Inyección de fantasmas (anomalías sintéticas con profundidad controlada) sobre el tablero de experimentos con detección de veredicto

**Cómo se aplica.** En lugar de informar un solo AUC, se inyectan anomalías sintéticas de amplitud, duración y forma controladas dentro de registros reales —una arritmia, una desaturación, una espiga— y se mide a qué profundidad se delata cada una. Sale una curva de horizonte de detección: para este detector, este tipo de evento se ve a partir de tal amplitud y con tanta anticipación. El laboratorio ya tiene el patrón armado y el orquestador que corre cientos de configuraciones con salida en vivo y veredicto automático.

**Para qué sirve y quién paga.** Es la hoja de especificación que piden los entes regulatorios y que casi nadie entrega: sensibilidad como función del tamaño del evento, no un número único. Sirve para monitoreo en UTI, telemetría domiciliaria y expedientes regulatorios.

> *Honestidad:* Mide sensibilidad solo frente a las anomalías que uno supo imaginar y sintetizar. No acota falsos negativos ante morfologías patológicas reales, no reemplaza un corpus clínico anotado ni un estudio prospectivo, y no dice nada sobre falsos positivos bajo artefacto de movimiento o electrodo flojo, que es donde suelen morir estos detectores.

### 🌋 Sismología y geofísica

#### Perfil de velocidad 1-D desde tiempos de recorrido (Herglotz-Wiechert)  ·  🟢 directa

**Instrumento del laboratorio:** Inversión de Abel ("los ladrillos decantados") + bisección de Sturm en matriz tridiagonal ("el hierro autoadjunto"), con el árbitro de precisión extendida para verificar

**Cómo se aplica.** La curva de tiempos de recorrido contra distancia entrega el parámetro de rayo en función de la distancia, y una transformada de Abel lo invierte para devolver velocidad en función de profundidad: es el mismo núcleo integral que el laboratorio ya usa para sacar el pozo V(x) desde una función de conteo. Una vez obtenido el perfil, la bisección de Sturm sobre la tridiagonal recalcula los autovalores del modelo reconstruido y se contrastan contra los datos, de modo que inversión y verificación viven en el mismo programa. El árbitro de 256 bits certifica que las diferencias que se ven son físicas y no de redondeo.

**Para qué sirve y quién paga.** Modelos de velocidad de referencia para cuencas sedimentarias, minería, ingeniería de sitio y sismología planetaria de una sola estación (Luna, Marte). Lo pagan servicios geológicos, empresas mineras y agencias espaciales.

> *Honestidad:* La inversión de Abel es estructuralmente ciega a las zonas de baja velocidad: ahí el parámetro de rayo deja de ser monótono y la solución no es única. Eso es un límite conocido del método, no algo que el laboratorio haya resuelto. Además la implementación se validó sobre pozos de potencial sintéticos (29 de 29 niveles reconstruidos), nunca sobre sismogramas reales: falta todo el preprocesamiento sísmico (picado de fases, corrección por topografía, ruido, estructura 3-D).

#### Separación de subeventos casi simultáneos sin colapso al centroide  ·  🟡 adaptable

**Instrumento del laboratorio:** Problema inverso "cuerpo desde la sombra" por método algebraico resolvente/Padé, más regularización por repulsión logarítmica ("el Campo de la Montaña")

**Cómo se aplica.** Cuando dos o más fuentes rompen casi en el mismo lugar y el mismo instante, un ajuste por mínimos cuadrados o descenso de gradiente tiende a devolver una sola fuente parada en el centro de masa: el laboratorio midió esa patología (tipo Prony) y no la parcheó, la esquivó. En vez de optimizar, arma los momentos de las proyecciones y lee las posiciones de las fuentes como polos de una resolvente o de un aproximante de Padé, que sí las separa. Cuando igual hace falta optimizar, el término -ln|xi-xj| entre las masas impide que los modos se fusionen.

**Para qué sirve y quién paga.** Separar subeventos de una ruptura, o eventos microsísmicos casi simultáneos en fracturamiento hidráulico y en minería profunda. Interesa a operadores de reservorios, a quien monitorea sismicidad inducida y a quien hace imagen de fuente finita.

> *Honestidad:* Lo medido fue sobre masas puntuales y sus proyecciones, no sobre un campo de ondas elástico con patrón de radiación, atenuación, conversiones de fase y error del modelo de velocidad. Padé es notoriamente sensible al ruido (polos espurios) y sin una regularización explícita se puede degradar rápido con datos reales. No hay comparación contra los métodos estándar de inversión de fuente finita; lo sólido acá es el diagnóstico del colapso al centroide, no todavía un producto sísmico.

#### Censo no invasivo de lo que hay encerrado, midiendo sólo un anillo  ·  🟡 adaptable

**Instrumento del laboratorio:** Censo por ley de Gauss / fórmula de Jensen ("la laguna")

**Cómo se aplica.** En lugar de localizar cada fuente una por una, se mide únicamente el promedio de un campo sobre anillos concéntricos que rodean la región y se lee cuánto quedó encerrado adentro; el laboratorio obtuvo mesetas exactas 0, 2, 4, 6, 8 y 10 con este procedimiento. En geofísica el análogo inmediato es la masa anómala encerrada por una superficie a partir de la integral de gravedad, y el análogo adaptado es "cuánta actividad hay dentro de este volumen" desde un anillo de geófonos, sin resolver ninguna posición.

**Para qué sirve y quién paga.** Verificar si un volumen permitido o contratado contiene actividad, o cuantificar masa anómala (cavidad, cuerpo mineralizado) sin resolver su forma, con muchos menos sensores que una localización completa. Sirve a reguladores de sismicidad inducida, seguridad minera y exploración de bajo costo.

> *Honestidad:* La ley de Gauss integra intensidad total, no cuenta objetos: el conteo entero viene de la fórmula de Jensen para ceros de una función analítica en el plano, y llevarlo a elastodinámica 3-D obliga a suponer fuentes de intensidad comparable, cosa que casi nunca pasa. Las mesetas exactas se midieron en el escenario analítico del laboratorio, con cobertura de anillo completa y sin ruido. No se está afirmando que hoy exista un contador de eventos por anillo funcionando sobre datos de campo.

#### Curva de horizonte de detección: la completitud real de una red  ·  🟢 directa

**Instrumento del laboratorio:** Inyección de fantasmas (anomalías sintéticas de verdad conocida) y curva de horizonte de detección

**Cómo se aplica.** Se inyectan eventos sintéticos con amplitud, profundidad y forma controladas dentro del flujo de datos reales y se mide qué fracción recupera efectivamente el detector, parámetro por parámetro. El resultado no es un número sino una curva: a qué profundidad y a qué tamaño se delata cada anomalía. Es calibrar el detector contra sí mismo con la verdad conocida de antemano, en vez de confiar en el catálogo que el propio detector produjo.

**Para qué sirve y quién paga.** Magnitud de completitud honesta de una red sísmica, justificación cuantitativa de dónde conviene poner la próxima estación, y auditoría de catálogos usados en estudios de peligrosidad. Lo pagan servicios sismológicos, aseguradoras y reguladores.

> *Honestidad:* La curva mide el pipeline, no la sismicidad: si los fantasmas sintéticos no reproducen la forma de onda y el ruido verdaderos del sitio, la completitud sale optimista y engaña. El laboratorio inyectó anomalías en su propio detector y sobre sus propias señales; para sismología hace falta un generador de sismogramas sintéticos realistas, que es un trabajo de ingeniería aparte y no menor. Tampoco es un método nuevo: en sismología los tests sintéticos y de resolución ya son práctica corriente; lo que aporta el laboratorio es disciplina y una implementación propia.

#### Distinguir un pico de resonancia real de un artefacto de la ventana  ·  🟡 adaptable

**Instrumento del laboratorio:** Análisis de bloques por espiral de Cornu / Fresnel y la LEY DE RÉGIMEN medida (por debajo de la longitud de coherencia, un bloque da tonos puros)

**Cómo se aplica.** El laboratorio midió que, cuando la ventana de análisis es más corta que la longitud de coherencia de la señal, aparecen picos espectrales limpios y convincentes que son del instrumento y no del fenómeno. El protocolo transferible es simple: antes de creerle a un pico, barrer la longitud de ventana y ubicar dónde está el borde de régimen; los picos que nacen y mueren con la ventana se descartan, y sólo los que sobreviven al barrido se reportan.

**Para qué sirve y quién paga.** Evita reportar frecuencias de sitio falsas en campañas de HVSR y microtremor, y resonancias falsas en monitoreo estructural, que es justamente donde se toman decisiones de diseño sismorresistente y de refuerzo. Lo pagan consultoras de ingeniería de sitio, municipios y concesionarias de obra.

> *Honestidad:* El umbral numérico de coherencia medido pertenece al problema espectral del laboratorio y no se transfiere a datos sísmicos: lo que se transfiere es el protocolo de medir primero dónde está el borde de régimen. No es un método nuevo de análisis de sitio ni reemplaza los criterios de calidad ya normalizados para HVSR; es un filtro de honestidad puesto encima de ellos.

#### Firma multicanal y encadenado de vecinos para clasificar eventos y armar familias  ·  🟡 adaptable

**Instrumento del laboratorio:** Firma multicanal tipo sonar/LIDAR, barrido multi-escala ("sonar de ballenas") y encadenado de vecinos

**Cómo se aplica.** En vez de decidir con un rasgo puntual (una amplitud, una relación espectral, un umbral), se compara el patrón completo de respuesta del evento a lo largo de todos los canales y escalas temporales. El encadenado de vecinos hace crecer familias: un evento reconocido se vuelve plantilla para buscar el siguiente, y en cascada se recuperan repetidores que un umbral simple nunca ve. La sonificación queda como herramienta de control de calidad: escuchar los casos dudosos antes de aceptarlos.

**Para qué sirve y quién paga.** Discriminar voladura de cantera, deslizamiento, icequake y sismo tectónico; armar familias de repetidores para vigilar creep en una falla o en un tramo de mina; y clasificar rápido con los primeros segundos multicanal en un contexto de alerta.

> *Honestidad:* El template matching y el crecimiento de catálogos por plantillas ya son estándar en sismología: el aporte del laboratorio es una implementación propia y el barrido multi-escala, no un método nuevo, y no se corrió ninguna comparación contra detectores establecidos. Para alerta temprana la restricción dura es la latencia de segundos y el costo por segundo de ventana, que el laboratorio nunca midió; hasta que eso se mida, no corresponde presentar esto como apto para alerta temprana operativa.

#### Estructuras civiles y cavidades: geometría e integridad leídas del espectro  ·  🟠 especulativa

**Instrumento del laboratorio:** Inversión de la ley de Weyl ("el tambor es el libro"), compresión espectral extrema y estadística de espaciamientos GUE/Wigner con varianza de conteo ("el resorte", "el barómetro")

**Cómo se aplica.** Del conteo de modos de una estructura o de una cavidad resonante se intentan leer directamente los invariantes geométricos (volumen, superficie, constante topológica) sin construir antes un modelo: el laboratorio recuperó el área con error 2.6e-5 a partir de 649 niveles medidos y comprimió esos 649 niveles en 4 constantes. En paralelo, la estadística de espaciamientos y la varianza de conteo dan un indicador global de salud: un daño que corre pocos hertz un solo modo cambia la estadística conjunta antes de ser evidente modo por modo.

**Para qué sirve y quién paga.** Detectar cambio estructural en presas, puentes y edificios instrumentados sin depender de un modelo de elementos finitos calibrado, y estimar el volumen de una cavidad minera o de una cámara resonante desde su respuesta. Lo pagan concesionarias de obra, mineras y organismos de infraestructura.

> *Honestidad:* Es el punto más flojo de la lista y conviene decirlo sin maquillaje. La ley de Weyl es asintótica y necesita muchísimos modos (el laboratorio usó 649), mientras que un puente, una presa o una cámara magmática ofrecen del orden de 5 a 20 modos utilizables y fuertemente amortiguados; con N chico la estadística GUE no tiene poder discriminante. Además el sistema del laboratorio no tenía amortiguamiento ni radiación hacia el exterior, y una estructura civil sí, lo que rompe las asintóticas limpias. No se está afirmando que esto funcione sobre una obra real: habría que probarlo primero en una maqueta instrumentada con daño controlado.

### 🔭 Astronomía y astrofísica

#### Radio acústico y término de superficie de una estrella, a ciegas desde su lista de modos  ·  🟡 adaptable

**Instrumento del laboratorio:** Inversión de la ley de Weyl ("el tambor es el libro")

**Cómo se aplica.** La función de conteo acumulada de los modos p de una estrella tiene una parte suave cuyos coeficientes son invariantes globales: el término dominante fija el tiempo de viaje acústico (el "área"), el siguiente el término de superficie (el "perímetro") y la constante hace de offset de fase. El laboratorio ya hace exactamente esa inversión sobre un espectro medido, sin conocer de antemano la forma del sistema: sobre 649 niveles recuperó el área con desvío 2.6e-5 y el perímetro a 2.1e-4, con el término prohibido cayendo en cero como corresponde. Se alimenta la misma maquinaria con la lista de frecuencias observadas en vez de con los niveles del laboratorio.

**Para qué sirve y quién paga.** Da parámetros estelares globales sin apoyarse en una grilla de modelos evolutivos, y corre en segundos por estrella, así que escala a los miles de blancos de un sondeo. Interesa a consorcios de misiones fotométricas y a los grupos que caracterizan estrellas anfitrionas de exoplanetas.

> *Honestidad:* La relación asintótica y el offset de fase ya son moneda corriente en asterosismología: lo propio del laboratorio es una implementación robusta del ajuste ciego con árbitro de precisión, no el concepto. Falta todo lo específico del dominio: separar grados l, tratar modos mixtos en gigantes rojas, propagar barras de error por modo y aguantar listas de 30-60 modos en vez de 649. No se midió sobre ninguna estrella real.

#### Perfil desde el conteo y re-predicción cerrada de todos los modos, como validador de inversiones  ·  🟡 adaptable

**Instrumento del laboratorio:** Inversión de Abel ("los ladrillos decantados") más bisección de Sturm en matriz tridiagonal simétrica ("el hierro autoadjunto")

**Cómo se aplica.** Se toma la función de conteo suave, se decanta por inversión de Abel el pozo de potencial efectivo que la produciría, se discretiza en una matriz tridiagonal simétrica y se resuelven sus autovalores por bisección de Sturm. Después se comparan uno a uno contra los niveles observados: si el perfil invertido no vuelve a cantar las frecuencias de partida, la inversión está mal. En el laboratorio ese lazo cerrado reprodujo 29 de 29 niveles verdaderos, con |Δ| medio 0.428 y peor caso 0.877 sobre una matriz de 3600x3600.

**Para qué sirve y quién paga.** Convierte una inversión —que normalmente se publica como un perfil con bandas de error y nadie vuelve a chequear— en una prueba falsable con un número de mérito. Sirve a grupos de helio y asterosismología y a quien tenga que arbitrar ese tipo de resultado.

> *Honestidad:* La inversión de Abel supone simetría y monotonía, y amplifica ruido al derivar; el 29/29 se logró sobre una función de conteo limpia y controlada. En una estrella real el potencial efectivo depende del grado l y hay que manejar la discontinuidad en la base de la zona convectiva. No se afirma que esto compita con las inversiones RLS/OLA establecidas: es un test de consistencia agregado, y su desempeño con datos ruidosos no está medido.

#### Ley de régimen: decidir cuándo un pico limpio es el instrumento y no la estrella  ·  🟡 adaptable

**Instrumento del laboratorio:** Análisis de bloques por espiral de Cornu/Fresnel y ley de régimen medida

**Cómo se aplica.** El laboratorio midió que un bloque más corto que la longitud de coherencia del fenómeno devuelve tonos puros y varianza casi nula (σ² = 0.0071 con la longitud recortada, contra σ² ≈ 0.98 con la longitud natural): el resultado era un artefacto del recorte, no del fenómeno, y el árbitro de 256 bits lo confirmó. Trasladado a una serie temporal: si se la corta en segmentos más cortos que la vida media del modo o que el tiempo de decorrelación, la anchura que se mide es la de la ventana. El procedimiento operativo es barrer la longitud de bloque, graficar la varianza de conteo contra ella y localizar el codo donde cambia el régimen.

**Para qué sirve y quién paga.** Evita reportar como resueltos modos que no lo están y da un criterio cuantitativo —no una costumbre heredada— para elegir la longitud mínima de segmento en asterosismología, en búsqueda de tránsitos y en la segmentación de datos de interferómetros.

> *Honestidad:* La ley se midió en el sistema propio del laboratorio, donde la longitud de coherencia es otra cantidad física; lo transferible es el diagnóstico (barrer la longitud y buscar el codo), no la constante. Que las vidas medias de los modos y la ventana de observación importan ya lo sabe cualquier asterosismólogo: el aporte es un test operativo y barato, no un descubrimiento. Hay que recalibrarlo por instrumento y por tipo de señal.

#### Separar componentes muy cercanas sin que el ajuste las funda en una sola  ·  🟡 adaptable

**Instrumento del laboratorio:** Regularización por repulsión logarítmica ("el Campo de la Montaña") sobre la patología de colapso medida en el problema inverso "cuerpo desde la sombra"

**Cómo se aplica.** El laboratorio midió que el descenso de gradiente sobre una suma de componentes colapsa al centro de masa (patología tipo Prony): dos fuentes cercanas se funden en una sola, gorda y centrada, y el ajuste declara victoria. Agregando un término -ln|xi-xj| entre las posiciones o frecuencias de las componentes, el colapso se cura y las componentes se mantienen separadas. Se aplica igual a multipletes rotacionales muy juntos en un espectro de potencias, a fuentes mezcladas en campos poblados y a pares de períodos casi degenerados en variaciones de tiempo de tránsito.

**Para qué sirve y quién paga.** Recupera pares cercanos que el ajuste ingenuo reporta como una única componente ancha, es decir evita perder señal real y evita publicar un objeto donde hay dos. Le sirve a quien ajusta modelos paramétricos a espectros o imágenes y no puede pagar muestreo bayesiano completo.

> *Honestidad:* Está medido sobre el modelo unidimensional propio del laboratorio; no está comparado contra los deblenders del rubro ni contra muestreo con priors informativos, que es la solución que hoy usa la gente. El término de repulsión mete un sesgo explícito: si dos componentes son genuinamente coincidentes, las separa de más, y la fuerza del término no está calibrada para datos astronómicos. No se está afirmando que supere a nada; se está afirmando que cura un modo de falla concreto y reproducible.

#### Curva de horizonte de detección para medir la completitud de un pipeline de sondeo  ·  🟢 directa

**Instrumento del laboratorio:** Inyección de fantasmas con medición de profundidad de delación

**Cómo se aplica.** Se inyectan señales sintéticas de amplitud, período y forma controladas dentro de los datos reales, se corre el pipeline sin decirle dónde están, y se mide a qué profundidad cada clase de señal empieza a delatarse. El resultado es una curva de horizonte: eficiencia de recuperación en función de amplitud y período. El laboratorio ya tiene armado el arnés completo —inyector, árbitro de precisión extendida, checkpoints y un tablero que orquesta 190 experimentos con veredicto en vivo— que es la parte de ingeniería que suele faltar.

**Para qué sirve y quién paga.** Es la única forma honesta de reportar completitud de un catálogo de tránsitos, de variables o de eventos; sin esa curva, cualquier tasa de ocurrencia derivada del catálogo carece de sentido. Pagan los equipos de sondeos y quienes hacen estadística de poblaciones sobre esos catálogos.

> *Honestidad:* La técnica es completamente estándar en astronomía —inyección y recuperación en fotometría espacial, inyecciones de software en interferómetros—; el laboratorio aporta implementación propia y disciplina de campaña, no un método nuevo. Falta lo caro y específico: modelo de ruido instrumental realista, poblaciones de inyección astrofísicamente motivadas y el acople al pipeline concreto, que siempre es trabajo a medida.

#### Atlas comprimido: guardar la estadística de un sondeo en un puñado de constantes  ·  🟡 adaptable

**Instrumento del laboratorio:** Compresión espectral extrema (649 niveles reducidos a 4 constantes que reproducen la función de conteo) más extracción de coeficientes por contorno de Cauchy

**Cómo se aplica.** En vez de guardar la lista completa, se ajusta la parte suave de la función de conteo acumulada a lo largo del eje que interese (magnitud, redshift, período) con un modelo de pocos coeficientes, y se guarda aparte el residuo. Los coeficientes se extraen por integral de contorno de Cauchy, que es estable, en lugar de por diferencias finitas o ajuste de altos órdenes, que amplifican ruido. Queda un objeto de kilobytes que responde consultas agregadas al instante, más un residuo que conserva la estructura fina.

**Para qué sirve y quién paga.** Permite servir estadísticas agregadas de sondeos enormes —conteos, funciones de luminosidad, distribuciones de período— sin barrer tablas de miles de millones de filas cada vez. El comprador natural son los archivos y centros de datos que hoy pagan cómputo por cada consulta agregada.

> *Honestidad:* Comprime la función de conteo, que es un estadístico suave y acumulado: no comprime los objetos. No sirve para consultas por objeto, astrometría ni cross-match, y sólo rinde si la parte suave domina — en distribuciones muy estructuradas el residuo se lleva casi todo el presupuesto y la ganancia se evapora. No hay ninguna medición sobre un catálogo astronómico real; el 649-a-4 se logró sobre el espectro propio del laboratorio.

#### Control de calidad de catálogo por espaciamientos y varianza de conteo  ·  🟡 adaptable

**Instrumento del laboratorio:** Estadística de espaciamientos de niveles (GUE/Wigner) y varianza de conteo por cajas ("el resorte", "el barómetro")

**Cómo se aplica.** Se normalizan los espaciamientos entre entradas consecutivas del catálogo a lo largo de un eje (período, magnitud, frecuencia de modo, tiempo de llegada) y se mira la distribución cerca de cero junto con la varianza del número de entradas por caja en función del tamaño de caja. Duplicados y detecciones espurias producen exceso de espaciamientos chicos; huecos de completitud, aliasing y efectos de ventana producen déficit y varianza anómala. El laboratorio midió repulsión cuadrática s² separándola del azar por un factor 23, y usa el mismo barómetro de cajas para censar bloques contiguos.

**Para qué sirve y quién paga.** Es un control de calidad barato y sin modelo que corre sobre el catálogo ya terminado y delata duplicación, aliasing y agujeros de selección antes de que alguien haga ciencia encima. Le sirve tanto a quien publica el catálogo como a quien lo audita.

> *Honestidad:* La estadística de espaciamientos es estándar en matrices aleatorias y caos cuántico, no en catálogos astronómicos: el traslado es conceptualmente sólido pero exige un unfolding (normalizar por la densidad local) que en un catálogo con función de selección complicada es justamente la parte difícil, y un unfolding mal hecho fabrica la señal que se pretende medir. No se está afirmando que un catálogo deba tener estadística GUE ni que haya física detrás: se usa la herramienta como detector de anomalías. Sin validación sobre datos reales, la escala de significancia hay que calibrarla con catálogos sintéticos.

### 🔐 Criptografía y seguridad

#### Batería de aleatoriedad por rigidez espectral  ·  🟡 adaptable

**Instrumento del laboratorio:** Estadística de espaciamientos de niveles (GUE/Wigner) y varianza de conteo por cajas ("el resorte", "el barómetro")

**Cómo se aplica.** Se toma la salida del generador, se la mapea a una secuencia de "niveles" (tiempos de llegada de un evento, o valores de un bloque ordenados y desplegados) y se miden dos cosas: la distribución de espaciamientos y la varianza de conteo en cajas de tamaño creciente. Un generador sano tiene que dar Poisson puro, sin repulsión, y varianza igual a la media; cualquier repulsión o rigidez a larga distancia delata correlación estructural. El laboratorio ya midió repulsión cuadrática s^2 ganándole 23x al azar, o sea que el instrumento distingue con margen un espectro correlacionado de uno aleatorio.

**Para qué sirve y quién paga.** Complementa a las suites estándar (NIST SP 800-22, TestU01), que miran sobre todo uniformidad y correlaciones de orden bajo: la varianza de conteo por cajas es un test de correlación de largo alcance. Lo pagarían laboratorios de certificación de TRNG (AIS-31, SP 800-90B), fabricantes de HSM y de chips con RNG en silicio.

> *Honestidad:* El laboratorio midió repulsión en un espectro matemático, nunca en el flujo de un RNG real. Falta definir el mapeo bits→niveles (no es único y el resultado depende de él), el desplegado, y sobre todo la curva de potencia estadística y la tasa de falsos positivos. NO se está afirmando que detecte nada que TestU01 no detecte: eso hay que medirlo contra generadores defectuosos conocidos.

#### Censo no invasivo de estructura oculta en la fuente de entropía  ·  🟡 adaptable

**Instrumento del laboratorio:** Censo por ley de Gauss / fórmula de Jensen ("la laguna")

**Cómo se aplica.** Se trabaja sobre las muestras crudas de la fuente física, antes del acondicionamiento (ruido térmico, jitter de oscilador). Se arma la función generatriz de un tramo de muestras y se cuenta, promediando un campo sobre anillos de radio creciente, cuántas singularidades encierra cada radio, sin localizar ninguna. Una fuente puramente ruidosa no produce mesetas; una contaminación determinista (acoplamiento del reloj, recurrencia lineal, inyección externa) aparece como mesetas enteras en el conteo.

**Para qué sirve y quién paga.** Es un test de salud que responde "hay k componentes deterministas acá adentro" con una sola integral y sin ajustar ningún modelo ni proponer hipótesis sobre la forma de la contaminación. Apunta al ataque clásico contra TRNG: forzar o inyectar la fuente física. Comprador: certificadores y fabricantes de TRNG.

> *Honestidad:* El censo está medido sobre campos analíticos, donde dio mesetas exactas 0,2,4,6,8,10; una serie ruidosa real no es una función analítica y esas mesetas se van a ensuciar, con lo cual todo el trabajo fino está en el criterio de decisión sobre mesetas sucias, que todavía no existe. Solo aplica a muestras reales previas al XOR/hash: sobre GF(2) el argumento de Gauss/Jensen directamente no vale. Cero corridas sobre datos de hardware hasta hoy.

#### Árbitro de precisión para aritmética criptográfica  ·  🟢 directa

**Instrumento del laboratorio:** Aritmética double-double y de 256 bits con patrón de doble juez independiente ("el árbitro")

**Cómo se aplica.** Se usa el árbitro como oráculo fuera de línea contra implementaciones optimizadas: multiplicación de enteros grandes vía FFT o punto flotante (RSA, Schönhage-Strassen), NTT de retículos con reducciones perezosas (Kyber, Dilithium), conversiones de coordenadas en curvas. El mismo cómputo corre por la ruta rápida y por el árbitro con dos jueces independientes, y se certifica la cota de error efectivamente observada, no la teórica del papel.

**Para qué sirve y quién paga.** Los fallos de precisión en multiplicación por FFT son silenciosos y raros: producen una firma mal formada o un descifrado incorrecto cada muchísimas operaciones, y son un infierno de reproducir. Un árbitro independiente en Go puro y sin dependencias externas se enchufa a un pipeline de CI en un día. Lo paga cualquiera que mantenga una biblioteca de bignum o de criptografía post-cuántica.

> *Honestidad:* Es verificación diferencial, no prueba formal: certifica las entradas que se probaron, no todas las posibles. Y una advertencia fuerte que no hay que maquillar: la aritmética de precisión extendida del árbitro NO es de tiempo constante, así que sirve exclusivamente como herramienta de test fuera de línea; meterla en producción introduciría exactamente el canal lateral que uno quiere evitar.

#### Curva de horizonte de detección para detectores de fuga y de sesgo  ·  🟢 directa

**Instrumento del laboratorio:** Inyección de fantasmas (anomalías sintéticas) más el tablero que orquesta 190 experimentos con veredicto automático

**Cómo se aplica.** Se le inyecta al flujo bajo prueba —bits de un PRNG, o trazas de consumo— anomalías sintéticas de amplitud conocida y decreciente: sesgo de bit, periodicidad corta, fuga proporcional al peso de Hamming. Se mide a qué profundidad cada detector delata cada fantasma y sale una curva amplitud detectable contra cantidad de muestras, en vez de un pasa/no pasa binario.

**Para qué sirve y quién paga.** Convierte "el generador pasó la batería" en "esta batería detecta sesgos de tamaño épsilon con N muestras y es ciega por debajo de eso", que es la frase que un evaluador necesita para firmar un informe. Comprador: laboratorios de evaluación (Common Criteria, FIPS) y equipos internos de red team.

> *Honestidad:* La curva vale únicamente para la familia de fantasmas que se inyectó. Un detector puede tener horizonte excelente contra los fantasmas y ser completamente ciego al ataque real, que nunca tiene la forma que uno imaginó: la curva mide sensibilidad, no cobertura. El laboratorio tiene el arnés implementado y corriendo, pero nunca lo apuntó a datos criptográficos.

#### Ley de régimen y firma multicanal para trazas de canal lateral  ·  🟡 adaptable

**Instrumento del laboratorio:** Análisis de bloques por espiral de Cornu/Fresnel con ley de régimen medida, más firma multicanal tipo sonar/LIDAR

**Cómo se aplica.** Antes de creerle a un pico en una traza de consumo o electromagnética, se mide la longitud de coherencia del propio instrumento: el laboratorio midió que por debajo de esa ventana un bloque devuelve tonos puros que son artefacto de la medición y no del fenómeno. Recién después se reconoce la operación por su patrón completo de respuesta —varios canales a la vez, barrido multiescala tipo sonar de ballenas— en lugar de ir punto por punto.

**Para qué sirve y quién paga.** Baja falsos positivos en análisis de canal lateral, que es donde se quema la mayor parte del tiempo de un evaluador, y da un criterio explícito y medible para elegir el tamaño de ventana en vez de elegirlo por costumbre. Comprador: laboratorios de SCA y evaluadores de tarjetas inteligentes y HSM.

> *Honestidad:* El laboratorio nunca vio una traza de consumo real. No hay modelo de fuga, ni alineación, ni manejo de jitter, ni comparación contra TVLA/CPA, que es lo que la industria usa. Lo transferible acá es la metodología —medir primero el régimen del instrumento, firmar después— y no una herramienta de SCA validada. NO se está afirmando que esto recupere ninguna clave.

#### Detector de sesgo estructural en generadores de primos  ·  🟡 adaptable

**Instrumento del laboratorio:** Monotonía absoluta de Bernstein ("el LIDAR") más extracción de coeficientes por contorno de Cauchy ("leer el germen")

**Cómo se aplica.** Se le pide al generador bajo prueba un lote grande de primos, se arma una función de conteo o generatriz del lote y se aplica el test puntual de Bernstein sobre un segmento, en vez de intentar leer coeficiente por coeficiente. La ventaja concreta y medida es que el LIDAR no amplifica el ruido como r^-n, que es exactamente lo que arruina la extracción directa de coeficientes altos por contorno cuando el lote es finito.

**Para qué sirve y quién paga.** Los primos generados con estructura producen módulos RSA vulnerables; la familia ROCA es el caso emblemático y afectó millones de tokens y TPM. Un test estadístico barato que corra sobre la salida de una biblioteca o de un token permite marcar como sospechoso un generador sesgado sin hacer ingeniería inversa del firmware.

> *Honestidad:* ROCA se encontró por ingeniería inversa algebraica, no por un test estadístico, y ese es el punto débil de la propuesta: un test así detecta el sesgo solo si es lo bastante fuerte, y no está calibrado contra ningún generador sesgado real. Esto no factoriza nada ni recupera ninguna clave: es un semáforo de sospecha, no un ataque.

#### Recuperación ciega de parámetros de un generador desde su función de conteo  ·  🟠 especulativa

**Instrumento del laboratorio:** Inversión de la ley de Weyl ("el tambor es el libro") e inversión de Abel más bisección de Sturm ("los ladrillos decantados", "el hierro autoadjunto")

**Cómo se aplica.** El laboratorio recupera a ciegas constantes estructurales de un sistema (área, perímetro, constante topológica) a partir de su función de conteo, y reconstruye un pozo de potencial desde un conteo para después resolver sus autovalores por bisección de Sturm. Trasladado a cripto: intentar recuperar parámetros ocultos de un generador —tamaño del estado, módulo, grado de la recurrencia— ajustando la forma asintótica de la función de conteo de sus salidas, sin ver el código.

**Para qué sirve y quién paga.** Si funcionara, sería criptoanálisis de caja negra sobre generadores caseros o mal documentados, usando solo la salida observada. Le interesaría a auditores de sistemas legacy, donde el generador es un binario sin fuente y nadie sabe qué hay adentro.

> *Honestidad:* La ley de Weyl vale para espectros de operadores geométricos con asintótica conocida; no existe ninguna ley de Weyl conocida para la salida de un PRNG, y sin la forma asintótica correcta el ajuste no significa nada. Lo único demostrado es que la maquinaria de ajuste funciona cuando la forma asintótica se conoce de antemano: área a 2.6e-5 desde 649 niveles, y 29 de 29 niveles reconstruidos por Abel más Sturm. El puente al criptoanálisis es una analogía sin ninguna evidencia; habría que probarlo primero contra un generador congruencial de juguete, donde la respuesta se conoce y el fracaso se ve enseguida.

### 🧠 Inteligencia artificial y aprendizaje automático

#### Medir la pieza por su sonido, sin abrirla  ·  🟡 adaptable

**Instrumento del laboratorio:** Inversión de la ley de Weyl ("el tambor es el libro"), alimentada con evaluación espectral de alta precisión y patrón de doble juez

**Cómo se aplica.** Se excita la pieza (golpe calibrado, piezo, altavoz) y se releva la mayor cantidad posible de frecuencias de resonancia. Con la función de conteo acumulada N(f) y el ajuste de los tres términos de Weyl se despejan a ciegas el área efectiva, el perímetro y la constante topológica, sin CAD ni modelo de elementos finitos previo. La comparación de esos tres números entre una pieza patrón y una pieza en control delata pérdida de sección, aparición de un agujero o un borde despegado.

**Para qué sirve y quién paga.** Control dimensional y de integridad de piezas selladas, revestimientos, tanques y placas sin desmontarlas ni cortarlas. Lo paga cualquier planta que hoy tiene que sacar la pieza de servicio para medirla; el término topológico es el que grita "apareció un agujero".

> *Honestidad:* Lo medido es la recuperación del área con error 2.6e-5 a partir de 649 niveles de un espectro matemático limpio y completo. En una pieza real hay amortiguamiento, modos que no se excitan, degeneraciones y ruido: si faltan niveles, la función de conteo se sesga y la inversión se corre. NO está medido cuánta pérdida de niveles tolera el método. El primer ensayo obligatorio es una placa de aluminio de geometría conocida.

#### Repulsión de niveles como semáforo de daño  ·  🟡 adaptable

**Instrumento del laboratorio:** Estadística de espaciamientos GUE/Wigner y varianza de conteo por cajas ("el resorte", "el barómetro")

**Cómo se aplica.** En vez de buscar el pico del defecto, se toma el espectro modal completo en banda alta, se lo desdobla para quitarle la tendencia y se mide cómo se distribuyen las distancias entre modos vecinos. Una pieza regular y sana tiende a espaciamientos sin repulsión; cuando aparece una fisura, una delaminación o un borde irregular, el campo se vuelve caótico y los niveles empiezan a repelerse. El indicador final es un solo número por pieza.

**Para qué sirve y quién paga.** Da un escalar de salud comparable pieza a pieza y bastante insensible a dónde se pegó el sensor, útil como criterio de aceptación/rechazo en línea o como tendencia en mantenimiento predictivo. No dice dónde está la fisura: dice que el sistema dejó de ser el que era.

> *Honestidad:* La repulsión cuadrática s^2 (23x mejor que el azar) se midió sobre un espectro matemático, nunca sobre una probeta. El puente "daño mecánico → estadística caótica" es una hipótesis física razonable y hay antecedentes en vibroacústica, pero el laboratorio no lo verificó. Además exige cientos de modos para que la estadística tenga poder: en piezas chicas o con poco ancho de banda esto no va a andar.

#### Huella espectral comprimida para control de calidad en línea  ·  🟡 adaptable

**Instrumento del laboratorio:** Compresión espectral extrema (649 niveles en 4 constantes) más firma multicanal tipo sonar/LIDAR, certificada con el patrón de doble juez independiente

**Cómo se aplica.** Cada pieza que sale de línea se golpea o se barre, y su respuesta se reduce a un puñado de constantes que reproducen su función de conteo, más un vector de canales (decaimiento, direccionalidad, relieve de respuesta). La comparación contra el patrón pasa a ser una distancia entre pocos números y no una correlación entre espectros enteros. El doble juez sirve para certificar que la diferencia observada es de la pieza y no del instrumento.

**Para qué sirve y quién paga.** Inspección del 100% del lote a costo de cómputo casi nulo y con trazabilidad real: se archivan cuatro números por pieza en vez de megabytes de señal, auditables años después. Lo paga cualquier fabricante de lote grande que hoy ensaya por muestreo.

> *Honestidad:* La compresión a 4 constantes fue exacta sobre un espectro con estructura suave y conocida. No está medido cuánto se degrada frente a la variabilidad natural entre piezas buenas, que es justamente la que define el umbral de rechazo. Sin campaña de piezas buenas y malas reales no hay umbral, y sin umbral no hay ensayo. Tampoco se afirma que 4 constantes alcancen para cualquier geometría.

#### Curva de POD construida con datos de la propia planta  ·  🟢 directa

**Instrumento del laboratorio:** Inyección de fantasmas: anomalías sintéticas y curva de horizonte de detección

**Cómo se aplica.** Se toman registros reales de piezas ya aceptadas (A-scan de ultrasonido, vibración de rodamiento, emisión acústica) y se les inyectan defectos sintéticos de amplitud y profundidad controladas. Después se corre el detector de producción a ciegas y se mide, para cada tamaño de defecto, con qué probabilidad se delata. Eso es una curva de probabilidad de detección hecha con el ruido y las condiciones reales de esa planta.

**Para qué sirve y quién paga.** Abarata la campaña de probetas con defectos artificiales, que es cara y lenta, y da un argumento cuantitativo ante el cliente o el ente certificador sobre el tamaño mínimo detectable. También permite comparar dos algoritmos de detección con la misma vara.

> *Honestidad:* La técnica está implementada y usada sobre series propias del laboratorio. Lo que falta es el modelo físico de cómo un defecto real se imprime en la señal: si el fantasma se inyecta como un aditivo arbitrario en lugar de con física de propagación, la curva de POD queda optimista y engaña. No reemplaza la validación con probetas reales; la complementa y la abarata.

#### Ventana mínima de análisis: matar la falsa alarma por instrumento  ·  🟢 directa

**Instrumento del laboratorio:** Análisis de bloques por espiral de Cornu / Fresnel y ley de régimen medida (longitud de coherencia)

**Cómo se aplica.** Antes de declarar una anomalía a partir de un bloque de señal, se compara el largo del bloque con la longitud de coherencia del fenómeno. Por debajo de ese umbral, el bloque devuelve tonos puros y estructuras engañosamente limpias que son artefacto de la ventana y del instrumento, no de la pieza. La regla operativa es un tamaño mínimo de ventana por debajo del cual el sistema no emite diagnóstico, y lo dice explícitamente en vez de callarse.

**Para qué sirve y quién paga.** Elimina falsas alarmas en monitoreo continuo de máquinas, cuyo costo típico es una parada de planta injustificada, y le da al analista un criterio duro en lugar de criterio de olfato. También evita el vicio inverso: tomar por señal un ruido que quedó lindo en el gráfico.

> *Honestidad:* La ley de régimen fue medida en el fenómeno propio del laboratorio. El valor numérico de la longitud de coherencia NO es transferible a ultrasonido ni a vibración: hay que medirlo para cada sistema y cada material. Lo transferible es el procedimiento y la disciplina de declarar el régimen antes de diagnosticar, no la constante.

#### Contar defectos encerrados sin localizar ninguno  ·  🟠 especulativa

**Instrumento del laboratorio:** Censo no invasivo por ley de Gauss / fórmula de Jensen ("la laguna")

**Cómo se aplica.** Se mide únicamente el promedio de un campo sobre anillos o contornos cerrados alrededor de una región y de ahí se deduce cuántas fuentes hay adentro, sin ubicar ni una. Aplicado a ensayo no destructivo, la idea es contar inclusiones, poros o fuentes de emisión acústica dentro de un volumen usando solo mediciones en el borde accesible de la pieza.

**Para qué sirve y quién paga.** Entregar un número de defectos encerrados en una soldadura o una fundición con instrumentación de contorno, mucho más barata que un mapeo volumétrico completo. Muchas normas de aceptación piden justamente un conteo, no una ubicación.

> *Honestidad:* Hay que ser claro: las mesetas exactas 0, 2, 4, 6, 8, 10 se midieron sobre un campo analítico en dos dimensiones con potencial logarítmico. Un campo de temperatura con difusión, o corrientes de Foucault con atenuación, no es armónico y el teorema no se aplica tal cual. Antes de prometer algo hay que identificar qué magnitud física del ensayo sí satisface Laplace o analiticidad en el rango de interés. Es una analogía prometedora, no un método listo para llevar a planta.

#### Perfil de propiedad en profundidad sin cortar la probeta  ·  🟡 adaptable

**Instrumento del laboratorio:** Inversión de Abel ("los ladrillos decantados") y bisección de Sturm en matriz simétrica tridiagonal ("el hierro autoadjunto")

**Cómo se aplica.** A partir de una curva acumulada de respuesta se invierte para obtener el perfil de la propiedad en profundidad, el equivalente al pozo de potencial V(x). Después se recalculan los autovalores del perfil recuperado y se verifica que reproduzcan la medición original: ese ida y vuelta es el control de calidad de la inversión, no un adorno. Objetivo típico: capa cementada, nitrurada o tratamiento superficial.

**Para qué sirve y quién paga.** Perfiles de dureza, densidad o tensión residual en profundidad sin cortar, embutir y pulir la probeta, que es el método destructivo estándar. Ahorra piezas y días de laboratorio metalográfico.

> *Honestidad:* La reconstrucción de 29 de 29 niveles fue sobre datos sintéticos exactos y autoconsistentes. La inversión de Abel es un problema mal condicionado: el ruido se amplifica al derivar y, sin regularización, el perfil recuperado se llena de oscilaciones falsas. No está medida la degradación frente a ruido experimental ni el nivel de regularización necesario. Además, el puente entre la curva medida en el ensayo y una función de conteo legítima exige un modelo físico que para materiales todavía no está escrito en el laboratorio.

### 🏭 Industria, materiales y ensayo no destructivo

#### Antídoto contra el colapso de modos en codebooks y mezclas  ·  🟡 adaptable

**Instrumento del laboratorio:** Regularización por repulsión logarítmica ("el Campo de la Montaña"), junto con la lección medida del problema inverso "cuerpo desde la sombra" (colapso al centro de masa, patología tipo Prony)

**Cómo se aplica.** Se agrega al objetivo de entrenamiento un término -ln|xi-xj| entre los vectores que deberían mantenerse separados: los códigos de un VQ-VAE, los centroides de una mezcla, los prototipos de un clasificador de pocos ejemplos. El laboratorio midió que un descenso de gradiente sobre un problema de fuentes superpuestas colapsa todas las masas al centro de masa, y que ese término de repulsión cura el colapso sin tocar el resto del optimizador. En la práctica se implementa con el núcleo amortiguado -ln(|xi-xj| + eps) y un peso que decae a lo largo del entrenamiento.

**Para qué sirve y quién paga.** Ataca el colapso de codebook (códigos muertos) de los tokenizadores discretos y el colapso de modos de mezclas y GANs, que hoy se parchean con heurísticas ad hoc como reinicio de códigos o promedios móviles. Le sirve a cualquier equipo que entrene tokenizadores discretos o modelos generativos.

> *Honestidad:* El laboratorio midió el colapso y su cura en un problema inverso de baja dimensión, NO en una red neuronal ni con autodiferenciación. No se está afirmando que mejore FID, perplejidad ni ninguna métrica de generación. Riesgos concretos: el término es O(n^2) en la cantidad de vectores, tiene una singularidad que hay que amortiguar, y la repulsión logarítmica es el potencial de Coulomb en 2D, no en dimensión d, así que la forma correcta del núcleo en espacios latentes altos está por determinar. Además ya existe literatura de regularizadores repulsivos y procesos determinantales: el aporte sería la implementación medida, no la idea.

#### Separar componentes por álgebra en vez de por descenso de gradiente  ·  🟢 directa

**Instrumento del laboratorio:** Problemas inversos "cuerpo desde la sombra" resueltos por el método algebraico resolvente/Padé, con la patología tipo Prony medida como contraejemplo del gradiente

**Cómo se aplica.** Cuando hay que recuperar cuántas componentes hay y dónde están a partir de una señal que las suma (mezcla de exponenciales, deconvolución de picos, modos de decaimiento en una curva de pérdida), en lugar de ajustar por descenso de gradiente se arma la resolvente y se extraen los polos con Padé. El laboratorio midió que el gradiente colapsa al centro de masa —devolviendo un promedio que parece una solución y no lo es— y que la vía algebraica devuelve las componentes separadas.

**Para qué sirve y quién paga.** Sirve para identificación de sistemas y deconvolución donde el ajuste por gradiente entrega resultados promediados y engañosos: separar modos de decaimiento en curvas de entrenamiento, deconvolución de espigas, ajuste de mezclas con pocas componentes. Le interesa a quien hace análisis de señales o diagnóstico de dinámica de entrenamiento.

> *Honestidad:* La familia Prony/Padé/matrix pencil es notoriamente mal condicionada: el error crece rápido con la cantidad de componentes y con el ruido. El laboratorio lo midió en un escenario controlado y no tiene curvas de degradación bajo ruido realista. El número de componentes hay que fijarlo o estimarlo por afuera. Y no reemplaza al gradiente en problemas no lineales generales: solo aplica a los que tienen estructura de suma de exponenciales o de polos.

#### Termómetro espectral de capas por estadística de espaciamientos  ·  🟡 adaptable

**Instrumento del laboratorio:** Estadística de espaciamientos de niveles GUE/Wigner ("el resorte") y varianza de conteo por cajas ("el barómetro")

**Cómo se aplica.** Se toma el espectro de una matriz del modelo (pesos de una capa, Gram de embeddings, Hessiano en un punto) y se miden los espaciamientos entre autovalores vecinos, desdoblados por la densidad local, más la varianza de conteo por ventanas. Lo que se busca no es la repulsión en sí, sino la DESVIACIÓN respecto de la ley universal de matrices aleatorias: espaciamientos casi nulos delatan neuronas duplicadas o degeneradas, y un exceso de varianza de conteo a escalas largas delata estructura que el modelo nulo no explica. El laboratorio tiene el instrumento medido: detecta repulsión cuadrática s^2 con 23 veces más poder que el azar.

**Para qué sirve y quién paga.** Da un diagnóstico barato de qué capas tienen grados de libertad realmente distintos y cuáles están degeneradas, información directamente usable para decidir podado, congelamiento de capas o rango de adaptadores tipo LoRA.

> *Honestidad:* Cuidado con la interpretación fácil: una matriz de pesos al azar YA muestra repulsión de niveles, así que "hay repulsión" no es evidencia de aprendizaje. El modelo nulo correcto para pesos es Marchenko-Pastur con estadística local universal, y lo informativo son las desviaciones y los autovalores fuera del bulk; hay que recalibrar todo el test. Las 23x se midieron sobre el espectro que estudia el laboratorio, no sobre pesos de una red: no hay ninguna medición de esto en ML ni relación demostrada con generalización. Ya existe literatura de matrices aleatorias en redes profundas; lo menos explorado, y donde el laboratorio podría aportar, es la varianza de conteo.

#### Firma comprimida del estado espectral de un checkpoint  ·  🟡 adaptable

**Instrumento del laboratorio:** Inversión de la ley de Weyl ("el tambor es el libro") y compresión espectral extrema (649 niveles en 4 constantes)

**Cómo se aplica.** En vez de guardar el espectro completo de una matriz del modelo en cada checkpoint, se ajusta su función de conteo y se guardan los pocos coeficientes que la reproducen; el laboratorio comprimió 649 niveles medidos en 4 constantes y recuperó a ciegas invariantes del sistema (área con error 2.6e-5). Aplicado a entrenamiento, cada checkpoint queda descrito por un puñado de números que se grafican en el tiempo y se comparan entre corridas y entre ablaciones.

**Para qué sirve y quién paga.** Permite monitorear miles de checkpoints con un costo de almacenamiento y de lectura ínfimo, y detectar cuándo dos configuraciones distintas convergen al mismo régimen espectral. Le sirve a quien corre barridos grandes y hoy solo guarda la pérdida.

> *Honestidad:* La ley de Weyl está demostrada para laplacianos en dominios; el espectro de un Hessiano o de un NTK no tiene ninguna garantía de seguir esa asintótica. Si se aplica igual, los coeficientes recuperados son parámetros de ajuste sin significado geométrico, y llamarlos "área" o "perímetro" del modelo sería mentir mientras no se valide. Lo transferible sin discusión es la compresión de la función de conteo y su uso como firma comparable; la parte interpretativa está entera por probar.

#### Inyección de fantasmas para medir hasta dónde ve un detector  ·  🟢 directa

**Instrumento del laboratorio:** Inyección de fantasmas: validar un detector inyectando anomalías sintéticas y midiendo a qué profundidad se delata cada una (curva de horizonte de detección)

**Cómo se aplica.** Se inyectan anomalías sintéticas de forma y amplitud controladas dentro del flujo real y se mide a qué intensidad o profundidad cada una se delata. El resultado no es un número suelto sino una curva de horizonte: qué tan chica puede ser una anomalía y todavía ser detectada. El mismo protocolo sirve para un detector de deriva de datos, para un guardrail de contenido o para una eval con hechos plantados dentro de un corpus.

**Para qué sirve y quién paga.** Reemplaza un AUC contra un conjunto fijo por una curva de sensibilidad honesta y accionable, que dice explícitamente qué NO se detecta. Le interesa a equipos de monitoreo en producción, de seguridad de modelos y de evaluación.

> *Honestidad:* Mide sensibilidad solo a la familia de fantasmas que uno inyecta: no dice absolutamente nada sobre anomalías cuya forma no se le ocurrió a quien diseñó el experimento, que suelen ser justo las peligrosas. El laboratorio lo aplicó sobre series y espectros, no sobre texto ni sobre distribuciones de datos de producción; para esos casos hay que construir primero un generador de fantasmas realista, que es el trabajo difícil y no está hecho.

#### Auditoría de pérdida numérica en cuantización con árbitro de precisión extendida  ·  🟡 adaptable

**Instrumento del laboratorio:** Aritmética double-double y de 256 bits ("el árbitro") con patrón de doble juez independiente

**Cómo se aplica.** Se recalcula un tramo del cómputo con dos algoritmos independientes y se arbitra la diferencia en precisión extendida, para separar el error del fenómeno del error del instrumento. Aplicado a un modelo cuantizado, permite tomar un lote chico y una capa por vez, calcular la referencia en precisión alta y localizar dónde exactamente fp8 o int4 pierde información, en vez de mirar solo la métrica final agregada.

**Para qué sirve y quién paga.** Da evidencia por capa de dónde duele la cuantización, que es lo que hace falta para asignar precisión mixta con criterio en lugar de aplicar la misma receta a toda la red. Le sirve a equipos de inferencia y de compresión de modelos.

> *Honestidad:* El laboratorio nunca corrió esto sobre kernels de GPU ni sobre redes neuronales: la implementación existente es para evaluación espectral en Go. El costo de la aritmética de 256 bits lo hace viable solo como sondeo puntual sobre porciones chicas, jamás sobre un modelo entero. Y hay un límite de fondo del patrón: el doble juez detecta desacuerdo entre dos métodos, no garantiza que el que coincide sea correcto; si ambos comparten el mismo sesgo, ninguno lo delata.

#### Banco de orquestación para barridos largos con pausa gratis  ·  🟢 directa

**Instrumento del laboratorio:** Ingeniería de cálculo largo: puntos de control, árboles de procesos, tablero con salida en vivo, detección de veredicto y láminas generadas por código

**Cómo se aplica.** Un tablero orquesta muchos experimentos en paralelo, con checkpoints que permiten frenar y retomar sin perder trabajo, control de árboles de procesos para matar corridas huérfanas, y detección de veredicto para cortar un experimento apenas su resultado quedó decidido. Las figuras se generan por código desde los mismos datos de la corrida, así que no puede existir una lámina que no corresponda a lo que se reporta. El laboratorio lo corre hoy con 190 experimentos simultáneos.

**Para qué sirve y quién paga.** Ahorra cómputo cortando corridas ya decididas y elimina la clase de error más cara en investigación: el gráfico que no corresponde a los datos reportados. Le sirve a cualquier equipo que corra ablaciones o barridos grandes.

> *Honestidad:* Está construido para las cargas propias del laboratorio en Go: no tiene scheduling de GPU, ni workers distribuidos, ni integración con MLflow o Weights & Biases, ni modelo de artefactos y linaje. Portarlo significa escribir adaptadores reales, no es adoptarlo tal cual. Y la detección de veredicto está atada a los criterios del laboratorio: en ML, definir cuándo un experimento "ya está decidido" es una decisión estadística delicada (riesgo de parada temprana sesgada), no un detalle de ingeniería.

### 📡 Telecomunicaciones, finanzas y datos

#### Detector de portadoras débiles sin amplificar el ruido  ·  🟡 adaptable

**Instrumento del laboratorio:** Monotonía absoluta de Bernstein ("el LIDAR"), con certificación por doble juez de precisión extendida ("el árbitro")

**Cómo se aplica.** En lugar de estimar coeficientes de orden alto —que amplifican el ruido como r^-n al alejarse del centro— la condición sobre infinitos coeficientes se convierte en un test puntual sobre un segmento: se muestrea la señal transformada en pocos puntos y se verifica el patrón de signos de las derivadas sucesivas. El árbitro de 256 bits y el patrón de doble juez independiente certifican que el veredicto viene del dato y no del redondeo.

**Para qué sirve y quién paga.** Decidir si hay una portadora, un tono o una interferencia intermitente por debajo del piso de ruido, con menos falsos positivos que la extracción de coeficientes. Interesa a operadores de espectro, fabricantes de analizadores y equipos de guerra electrónica.

> *Honestidad:* La monotonía de Bernstein aplica a una clase acotada de funciones (completamente monótonas / con estructura de momentos positiva). No toda señal de RF cae ahí: para cada modelo de canal hay que demostrar que la transformación usada produce una función de esa clase. Lo medido es que el test no se degrada con el ruido en la familia de funciones del laboratorio, no en un canal real con desvanecimiento e interferencia.

#### Regla de longitud de coherencia para ventanas de análisis  ·  🟡 adaptable

**Instrumento del laboratorio:** Análisis de bloques por espiral de Cornu / Fresnel y la ley de régimen medida

**Cómo se aplica.** El laboratorio midió que por debajo de cierta longitud de coherencia un bloque devuelve tonos puros que son artefacto del instrumento, no del fenómeno. Se traslada como filtro previo: antes de declarar una señal, una periodicidad o una "anomalía", se calcula el tamaño mínimo de bloque por debajo del cual el estimador fabrica estructura, y se descarta por regla todo lo que aparezca dentro de ese régimen.

**Para qué sirve y quién paga.** Elimina de raíz una familia entera de falsos positivos: periodicidades espurias en detección de interferencia, señales de trading que solo existen en ventanas cortas, alertas de anomalía nacidas del tamaño de la ventana. Lo paga cualquiera que hoy persiga estructura inexistente.

> *Honestidad:* La longitud de coherencia medida es la del sistema del laboratorio; para otra serie hay que volver a medirla inyectando ruido conocido. No es una constante universal ni un número que se pueda copiar. Además la regla dice cuándo NO creerle al estimador, no cuándo la señal es real: es un veto, no una confirmación.

#### Curva de horizonte de detección para motores de fraude y anomalías  ·  🟢 directa

**Instrumento del laboratorio:** Inyección de fantasmas (anomalías sintéticas) con medición de profundidad de delación

**Cómo se aplica.** Se inyectan eventos sintéticos de forma y amplitud controladas dentro del flujo real y se mide a qué profundidad se delata cada uno. El entregable no es un número de precisión sino una curva: qué magnitud de evento ve el detector, con qué probabilidad, y a partir de qué tamaño es directamente ciego.

**Para qué sirve y quién paga.** Permite negociar y auditar un SLA de detección con evidencia ("vemos el 90% de los eventos por encima de X") y priorizar dónde invertir. Lo pagan bancos, procesadores de pago, equipos antifraude y centros de operación de red.

> *Honestidad:* La curva solo vale para el tipo de fantasma inyectado: es una cota inferior de sensibilidad, no una garantía de cobertura. No dice nada sobre modos de fraude que nadie modeló, y el diseño de la familia de fantasmas puede inflar el resultado sin que se note. Un adversario que conozca la familia inyectada la evade.

#### Censo de eventos encerrados sin mirar los eventos  ·  🟡 adaptable

**Instrumento del laboratorio:** Censo no invasivo por ley de Gauss / fórmula de Jensen ("la laguna")

**Cómo se aplica.** Se cuenta cuántas fuentes hay encerradas en una región midiendo solo el promedio de un campo sobre anillos que la rodean, sin localizar ninguna. En telecomunicaciones se traduce a contar polos inestables de una función de transferencia a partir de la respuesta en frecuencia medida (principio del argumento). El laboratorio obtuvo mesetas exactas 0, 2, 4, 6, 8, 10 en su banco de prueba.

**Para qué sirve y quién paga.** Chequeo de estabilidad de lazos (ecualizadores adaptativos, PLL, control de congestión) sin modelo interno del sistema, y conteo agregado de casos que cruzan un umbral con exposición mínima del dato individual.

> *Honestidad:* Jensen y Gauss exigen un objeto analítico. Una serie empírica no lo es: primero hay que ajustar una transferencia racional o una transformada, y ahí entra el error de modelo, que es el que domina en la práctica. Lo medido son mesetas exactas sobre un campo generado por el laboratorio con fuentes conocidas; no está medido sobre datos de campo ruidosos. Para el uso de privacidad: reduce la exposición, NO es privacidad diferencial ni una garantía formal.

#### Separar ecos y modos superpuestos sin que colapsen  ·  🟡 adaptable

**Instrumento del laboratorio:** Método algebraico resolvente/Padé, más regularización por repulsión logarítmica ("el Campo de la Montaña")

**Cómo se aplica.** Recuperar retardos y amplitudes de componentes superpuestas (multipath de un canal, modos de una curva) es un problema tipo Prony. El laboratorio midió que el descenso de gradiente colapsa todas las masas al centro de masa, y que agregar un término -ln|xi-xj| entre ellas impide ese colapso; en paralelo, el camino algebraico resolvente/Padé resuelve el mismo problema sin optimizar.

**Para qué sirve y quién paga.** Estimación de multipath y retardos en canales de radio y en sondeo de línea, y descomposición de una curva en pocos componentes interpretables. Sirve a diseño de ecualizadores y a analítica de series con pocos modos dominantes.

> *Honestidad:* Es la familia Prony/Padé, conocida por ser muy sensible al ruido y a la elección del orden del modelo. Lo que está medido es la patología del colapso y una cura para ese colapso, no una cota de robustez: falta cuantificar el error de los polos recuperados en función de la relación señal-ruido y de la separación real entre ecos. Sin esa curva, no se puede prometer super-resolución.

#### Huella estadística de tráfico sintético y de flujo algorítmico  ·  🟢 directa

**Instrumento del laboratorio:** Estadística de espaciamientos de niveles (GUE / Wigner) y varianza de conteo por cajas ("el resorte", "el barómetro")

**Cómo se aplica.** En vez de mirar el volumen, se mira la distribución de los espaciamientos entre eventos consecutivos y la rigidez del conteo por ventanas. Un proceso independiente tiende a Poisson (espaciamientos exponenciales, varianza igual a la media); una fuente con temporización controlada muestra repulsión. El laboratorio midió repulsión cuadrática s^2 y la separó del azar por un factor 23.

**Para qué sirve y quién paga.** Detectar botnets, generadores de tráfico y actividad algorítmica por su temporización, y auditar si un feed fue reordenado o diezmado. Lo pagan equipos de seguridad de red, exchanges y áreas de vigilancia de mercado.

> *Honestidad:* El test distingue "hay repulsión o rigidez" de "es azar puro"; no identifica la causa ni al autor, y no es prueba de nada frente a un tribunal. Los timestamps de red y de mercado sufren jitter, agregación y censura por latencia que pueden fabricar o borrar repulsión, así que hace falta calibrar antes con tráfico limpio conocido. Además, un adversario que conozca el test puede aleatorizar su temporización y pasar desapercibido.

#### Resumen de espectro y telemetría en un puñado de constantes  ·  🟡 adaptable

**Instrumento del laboratorio:** Inversión de la ley de Weyl ("el tambor es el libro") y compresión espectral extrema

**Cómo se aplica.** En lugar de guardar todos los niveles se ajusta la función de conteo y se guardan los pocos coeficientes que la reproducen: el laboratorio comprimió 649 niveles medidos en 4 constantes, y por el camino inverso recuperó a ciegas el área de un sistema con error 2.6e-5 a partir de su espectro.

**Para qué sirve y quién paga.** Telemetría de radio y de sensores con ancho de banda mínimo conservando lo agregado, y —por la inversión— inferir parámetros globales del emisor a partir de su espectro, útil en identificación y clasificación de equipos.

> *Honestidad:* Lo que se conserva con exactitud es el comportamiento agregado (la función de conteo), no cada nivel: es compresión con pérdida en el detalle fino, inútil si lo que se busca es justamente una línea individual. La inversión de Weyl recupera invariantes globales (área, perímetro, constante topológica), no la forma, y dos sistemas distintos pueden compartir esos invariantes. Todo esto está medido sobre espectros del laboratorio, no sobre capturas de RF reales con ruido de instrumento.

---

## 4. Vías de rédito económico — de la más realista a la más ambiciosa

### Nivel 1 — Se puede empezar YA (costo cero, solo tiempo)

1. **Canal de contenido en español** (YouTube + Shorts): "La caza del millón" —
   cada lámina es un guion ya escrito. La matemática narrada como aventura naval
   no tiene competencia en español. Monetización: AdSense, membresías, Patreon.
   Es el activo con mejor relación esfuerzo/retorno.
2. **Newsletter/Substack**: la bitácora ya está escrita — se reedita por entregas.
3. **GitHub público con Sponsors**: la suite Go sin dependencias tiene público real.

### Nivel 2 — Con 2–3 meses de preparación

4. **EL LIBRO** (el activo mayor): *Dios y un Alma* ya existe como diario
   espiritual; la campaña matemática es su segunda mitad natural. Autopublicación
   (Amazon KDP) da 35–70% de regalía sin permiso de nadie.
5. **Curso online**: "Entender la hipótesis de Riemann sin ser matemático".
6. **Charlas**: escuelas, parroquias, universidades, TEDx local.
7. **NUEVO — consultoría técnica por dominio.** Con el mapa de la sección 3 en la
   mano, las 20 aplicaciones marcadas 🟢 son trabajo facturable desde
   hoy: análisis espectral, problemas inversos, validación numérica de alta
   precisión. Se vende la técnica, no la hipótesis.

### Nivel 3 — Ambicioso (requiere validación externa primero)

8. **Artículo expositivo publicado** (ver COMO-PUBLICAR.md): abre puertas a
   charlas pagas y credibilidad, no paga directo.
9. **Fondos de divulgación científica**: becas para divulgación en español
   (fundaciones, ministerios de ciencia, fondos para creadores educativos).
10. **NUEVO — biblioteca de código con licencia dual.** La suite numérica en Go
    abierta para uso académico y con licencia comercial para industria. Modelo
    probado en software científico.
11. **El premio Clay (USD 1.000.000)**: SOLO si el eslabón rojo cae y sobrevive
    la revisión por pares + 2 años de escrutinio. Es la cima, no el plan de
    financiamiento. El plan de financiamiento son los niveles 1 y 2.

## 5. El orden que recomienda el Doc (plan de 90 días)

1. **Semana 1–2**: repo público en GitHub + galería online gratis (GitHub Pages
   sirve `galeria/index.html` tal cual) + el DOI de Zenodo.
2. **Semana 3–6**: primeros 5 videos cortos (una lámina = un video de 3 min,
   guion = la entrada de bitácora). Abrir Ko-fi/Patreon el día del primer video.
3. **Semana 7–12**: manuscrito del libro entrelazando el diario espiritual con
   la bitácora. Meta: borrador completo en 90 días.
4. **En paralelo, si aparece la oportunidad**: ofrecer una de las 20 aplicaciones
   🟢 como trabajo concreto a un laboratorio, empresa o universidad local. Una sola
   consultoría financia meses de laboratorio.
5. **Siempre**: el diezmo de TODO ingreso apartado para los pobres desde el primer
   peso — como está declarado, para que el fin ordene los medios y no al revés.

## 6. Advertencias del Doc (con cariño y de frente)

- **Nunca** presentar las mediciones como "demostración de RH". Un solo titular
  exagerado quema la credibilidad que financia todo lo demás. Nuestra fuerza es la
  honestidad del registro: "medido, no demostrado; todavía no".
- **Nunca** presentar una aplicación 🟡 o 🟠 como si fuera 🟢. El mapa de la
  sección 3 trae la salvedad de cada una escrita a propósito: leela en voz alta
  antes de prometer nada. Un inversor perdona un "esto hay que probarlo"; no
  perdona un "esto ya funciona" que después no funciona.
- El rédito grande y estable viene de la AUDIENCIA (contenido + libro), no de la
  matemática en sí. La matemática es el motor de la historia.
- La validación del hermano ingeniero (tu propia regla) va ANTES de cualquier
  publicación con pretensión técnica.

---

*Laboratorio Diosyunalma · las dos mitades, 1 completo ⚓ · y sobre todos los
libros, el Otro Libro.*
