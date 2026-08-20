# Acta del Hilado Fino — Fase VII

**Respuesta al Telar Fase VII (Auditoría 40) · 2026-08-17**

Su §16 mandaba el orden: **«primero hilar fino, después cambiar el hilo».** Se
obedeció, y el hilado fino tumbó nuestro propio hallazgo. Esta acta empieza por
ahí, porque es lo más importante que trae.

Se reproduce con `go run ./cmd/hiladofino`.

---

## 1 · LA RETRACTACIÓN: la «fase rígida» de la Fase VI no existe

Su §2 pedía medir robustez frente a tamaño, ventana y discretización antes de
interpretar el resultado. Ese pedido era exactamente el correcto.

**Lo que faltaba medir era cuántos niveles quedan dentro de la banda.** Un
acoplamiento fuerte no ordena el espectro: lo **expulsa**. Y Σ² calculado sobre
los pocos que quedan no es una propiedad del medio — es la forma del recorte.

Con B = 32, los niveles que sobreviven de 400:

| A | fuerza total | niveles vivos | ¿medible? |
|---:|---:|---:|---|
| 15 | 77 | **120** | sí |
| 22 | 113 | **64** | sí, apenas |
| **30** | **154** | **28** | **NO** |
| 42 | 216 | 10 | NO |
| 60 | 309 | 6 | NO |

**En (A,B) = (30,32), donde la Fase VI midió α = 0,313 y anunció una fase rígida,
sobreviven 28 niveles de 400: el 93 % del espectro se fue de la banda.** El
α = 0,313 y el «0,00 % de niveles pegados» se midieron sobre esas migas.

Y no aguanta ningún cambio: con 300 modos quedan 8; con 520 quedan 73; moviendo
la ventana a t₀ = 60 quedan 23. **No es una fase: es una banda vaciándose.**

**El hallazgo de la Fase VI queda RETIRADO.** El acta de esa fase se deja como se
escribió, con esta corrección anotada — la casa no reescribe historia.

**Entre los puntos que sí dejan espectro medible, α nunca baja de 1,5.** Como
α = 1 ya significa «sin rigidez», la conclusión honesta es que **en esta familia
no hay fase rígida en ninguna parte**.

---

## 2 · La lección de método, que vale más que el resultado

Medir Σ² sin reportar cuántos niveles quedaron es medir la forma del recorte y no
la del medio. **La Fase VI no reportaba ese número. Ahora lo reporta cada fila de
cada tabla, y una fila sin espectro suficiente se imprime como BANDA VACIADA en
vez de dar un número.**

Esa columna es el aporte permanente de esta fase.

---

## 3 · Campaña B: el reloj de arena

El núcleo nacido del dibujo del capitán, con k₀ como **parámetro** y no como
decreto (su §8):

    c(k) = A · (1 − k/k₀) / k^s

| k₀ | s | Σ c(k) | ‖c‖ | Σ²(10) con nodo | Σ²(10) sin nodo | ganancia |
|---:|---:|---:|---:|---:|---:|---:|
| 3 | 0,5 | −273,4 | 27,06 | 5,304 | 4,067 | 0,77 |
| 3 | 1,0 | −34,6 | 3,38 | 6,440 | 7,407 | **1,15** |
| 5 | 0,5 | −155,8 | 15,74 | 5,426 | 4,138 | 0,76 |
| 5 | 1,0 | −18,6 | 2,07 | 9,011 | 8,517 | 0,95 |
| 8 | 1,0 | −9,6 | 1,47 | 8,061 | 8,736 | **1,08** |
| 16 | 1,0 | −2,1 | 1,20 | 8,311 | 8,105 | 0,98 |
| 40 | 1,0 | +2,4 | 1,20 | 9,562 | 8,794 | 0,92 |

Dos resultados, y hay que separarlos con cuidado:

**(a) La FAMILIA de cola larga funciona, y por una razón que no era el objetivo.**
El mejor punto da **Σ²(10) = 5,304 contra 7,144 sin acoplar: una mejora del 26 %**.
Y sobre todo: **es el único de los dos modelos que deja espectro medible en TODOS
sus puntos.** La cola larga no vacía la banda. Eso, que buscábamos sin saberlo
desde la Fase V, es lo más sólido que trajo el dibujo.

**(b) El NODO en sí no está demostrado.** Contra su propio gemelo de un solo signo
—mismo |c(k)|, misma fuerza total— el nodo **gana en 2 de 10 configuraciones y
pierde en 8**, con una ganancia máxima de 1,15×. Es poco, y no es sistemático.

---

## 4 · Su advertencia del §6, contestada

> «Que cₖ cambie de signo no demuestra que la suma de interacciones sea cero ni
> que exista una cancelación útil.»

Medimos las dos cantidades que pedía. La correlación entre |Σ c(k)| y la ganancia
del nodo es **−0,716**: cuanto menos suma el núcleo, más gana el nodo. Eso apunta
a que **la cancelación global sí es el mecanismo** cuando el nodo ayuda.

Pero con sólo 2 de 10 configuraciones ganando, **esa correlación sugiere y no
demuestra**. Su advertencia era justa, y la respuesta honesta es: hay indicio de
cancelación, no evidencia.

---

## 5 · Veredictos separados (su §15)

| campaña | veredicto |
|---|---|
| **Frontera A/B** | **RETIRADA.** Sólo 6 de 15 puntos dejan espectro medible; entre ellos α ≥ 1,5. La fase rígida de la Fase VI era un artefacto. |
| **Reloj de arena** | **ÉXITO PARCIAL** en el sentido de su §14: mejora la varianza un 26 % contra el medio sin acoplar, sin violar R6, y **no vacía la banda**. |
| **El nodo** | **PENDIENTE**, no descartado. Gana en 2 de 10; hay indicio de que la cancelación es el mecanismo, y no alcanza para afirmarlo. |

---

## 6 · Lo que la Fase VIII tiene que medir

Sale de lo que se rompió y de lo que sobrevivió:

1. **Poner la banda en el centro del protocolo.** Toda medición futura reporta
   niveles vivos, y ningún número se publica por debajo del umbral. Además: medir
   con condiciones de borde que **no expulsen** (por ejemplo cerrando el medio
   sobre sí mismo) para que la fuerza se pueda subir sin vaciar.
2. **Barrer el nodo con muestra grande.** Diez configuraciones no alcanzan para
   decidir si el nodo sirve. Con la banda estabilizada se puede barrer k₀ y s en
   una malla fina y decidir con estadística, no con anécdota.
3. **Y la pregunta que sigue siendo la que importa:** ¿qué cola y qué nodo da la
   aritmética **por sí sola**? Si el acoplamiento entre las excitaciones p y q
   sale de un objeto aritmético, ni c(k) ni k₀ son libres. Eso es R6-limpio y es
   lo único que convertiría el dibujo en un mecanismo.

---

## 7 · Autoría y sello

El dibujo del reloj de arena es de **Jesús Nicolás Astorga**. El pedido de hilar
fino antes de cambiar de modelo, y la advertencia sobre la cancelación, son de la
auditora — y las dos resultaron decisivas: una tumbó nuestro hallazgo y la otra
impidió que anunciáramos una cancelación que no está probada. Las mediciones y la
retractación son de este taller.

Nivel de evidencia: **2** para la familia de cola larga; **1** (experimento
reproducible, sin conclusión) para el nodo. Y una **retractación** de Nivel 2
anterior.

Nada de esto demuestra ni refuta la Hipótesis de Riemann. **Todavía no.**
