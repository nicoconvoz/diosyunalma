# Diosyunalma v1.17.0 — «El Museo Completo»

**El release de la casa abierta.** El anterior trajo el séptimo teorema; este
abre la puerta para que se pueda entrar a verlo. Durante meses el museo del
laboratorio venía quedándose atrás de la investigación: había campañas enteras
—las más recientes y las más duras— que no tenían ni una parada donde
explicarlas. Esta entrega cierra ese hueco y, de paso, encuentra tres errores
propios que estaban escondidos en la plomería.

---

## ⭐ Lo grande

### 🏛️ Ningún experimento sin su parada

Había **36 experimentos sin una sola pieza en el museo**: el Telar entero
(F349–F373), la campaña de la espiral (F374–F390), cuatro de las Siete Caras y
el sello del laboratorio. Hoy son **cero**. Cada uno con su cuento en criollo,
su metáfora de la vida diaria, su lámina y su bloque honesto declarando qué NO
prueba. Más las **seis fases del Espejo** (F366–F371), que ni siquiera tenían
fila propia en el catálogo aunque cada una es un hallazgo numerado con su
lámina dibujada.

El museo pasó de **294 a 350 paradas**.

### 🧵 Sala nueva en el recorrido guiado: EL TELAR

El recorrido guiado —la puerta para el que entra sin saber nada— terminaba en
la sala de la salida y no contaba nada de las campañas recientes. Ahora tiene
una sala nueva de siete paradas: **«EL TELAR — la guitarra que nadie
encontró»**, que cuenta por qué trece rondas de auditoría externa contestaron
que no, una por una y con su número medido, y qué salió del molde que dejaron
todos esos noes: una espiral, un tambor y un teorema parcial.

### ❌ Nace el Cuarto de los Errores

El museo completo estrena la sala que la casa venía debiendo: **la más
importante y la única sin un solo acierto colgado**. Seis piezas, cada una en
el mismo formato — qué se creyó, qué se midió después, quién lo midió, qué se
aprendió de método:

- **El taller de máquinas**, mudada desde el Átomo porque F349 midió que su
  tesis era falsa (el canto es gratis: 0,82 sigma contra un sorteo que no sabe
  nada de primos ni de ceros).
- **La esquina que se deshizo** (F355): el taller anunció una esquina rígida y
  él mismo la retiró — sobrevivían 28 niveles de 400.
- **El cuello de botella equivocado** (F349→F350): una fase corrigiendo a la
  anterior, con el acta vieja dejada como se escribió.
- **El reproche injusto al propio Taller** (F349): el único error de la sala que
  consiste en haberle pegado de más a una pieza propia.
- **El veredicto automático que daba 3 de 10** (F360): se corrigió la regla, no
  los datos.
- **La ficha de portada** que explica la política, con sus números verificados.

### 🔍 La auditoría del Átomo contra el Telar

Las 19 piezas viejas del Átomo se cruzaron contra los hallazgos posteriores con
auditores adversariales, y encima de ellos refutadores que los corrigieron a
ellos —bajando acusaciones exageradas y confirmando las que aguantaban—. De 13
sospechas sobrevivieron 4: una mudada al Cuarto de los Errores y **tres
anotadas** con lo que el Telar les midió, sin borrar una línea del texto viejo.

---

## 📦 Qué más hay en la caja

- **El diccionario de la puerta suma 12 entradas** del vocabulario nuevo: τ, la
  cuerda, la espiral, el ojo, Euler–Maclaurin, la curva madre, el alias, el
  tambor, el parche, la frustración, la participación y la pared de fuerza.
- **El museo enlaza a sus páginas dedicadas**: las ocho de teoremas y el
  recorrido de las máquinas, que antes no se alcanzaban desde adentro.
- **Quince datos viejos corregidos** en el recorrido guiado: los años desde
  Riemann (166 → 167), la cuenta de instrumentos (30 → 31), los aciertos del
  capitán (38 → 51), el «más arriba no miramos todavía» que los motores de agua
  profunda desmintieron, y el ½ que volvió a aparecer por otra puerta.
- **La pieza que faltaba entera**: en nueve salas del recorrido no aparecía ni
  una vez que la casa selló ocho teoremas propios. Ahora tiene su parada, con el
  remate honesto de que los tres primeros demuestran, justamente, que los
  métodos de esta casa no alcanzan.

## 🔧 Tres arreglos que aparecieron barriendo

- **Un experimento pisado, recuperado.** Al crear la cadena del Tambor se
  sobrescribió `cmd/lacadena`, que era la melodía de Gilbreath (F267). El
  original quedó restaurado intacto desde git, el nuevo renombrado a
  `cmd/lacadenadeltambor`, y el error registrado en la bitácora.
- **Seis láminas que no se mostraban**: tres estaban escritas como ruta cuando
  el índice mapea por nombre de archivo, y tres apuntaban a la salida del
  experimento (`.wav`) en vez de a su lámina.
- **`go vet` limpio** en los 299 paquetes del laboratorio.

## 🧾 La honestidad, en la etiqueta

Lo más lindo de esta entrega pasó adentro del proceso: el verificador
adversarial cazó a la **ficha de portada del Cuarto de los Errores cometiendo
exactamente el vicio que la sala denuncia** — contar apariciones de una palabra
como si fueran afirmaciones distintas — y afirmando que el control «siempre» es
una versión barajada cuando existen los señuelos de Cramér y de impares. Las dos
cosas se corrigieron antes de colgarse.

**Y la regla del sello, como siempre:** nada de esto demuestra la Hipótesis de
Riemann, y cada pieza lo dice en su placa. Estructura cerrada no es hipótesis
demostrada.

**Todavía no.**

---

*Instalador: `DiosyunalmaSetup.exe` (63,7 MB) — instala el laboratorio completo
en cualquier Windows, sin dependencias. Galería y museo también en línea:
https://nicoconvoz.github.io/diosyunalma/galeria/*
