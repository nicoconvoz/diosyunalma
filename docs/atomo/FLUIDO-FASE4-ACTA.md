# Acta del Fluido — Fase IV

**Respuesta a la Hipótesis del Fluido Resonante (Auditoría 37) · 2026-08-17**

> «Es como una onda en un fluido: algunas grandes, otras pequeñas, otras
> medianas, todas resonando y creando una melodía única.»
> — **Jesús Nicolás Astorga**

**La idea es del capitán.** La auditora la formalizó como ruta de investigación y
nos pidió atacarla, formalizarla o descartarla con experimentos. Esta acta hace
las tres cosas a la vez: **formaliza un modelo mínimo, lo mide, y descarta una
versión — dejando viva la otra.** Se reproduce con `go run ./cmd/elfluido`.

---

## 1 · Por qué el cambio de lenguaje no es cosmético

La Fase III cerró toda familia cuyas órbitas cerradas se puedan **concatenar**.
El §6 de su documento lo dice mejor de lo que lo habíamos dicho nosotros:

> La pregunta no es «¿cómo prohibimos la concatenación?» sino «¿podemos construir
> un sistema en el que la concatenación ni siquiera sea la operación
> fundamental?»

Eso es exactamente lo que hace un medio. **En un fluido no hay órbitas que pegar:
hay modos que se acoplan.** El no-go de la Fase III no se esquiva con un truco —
deja de aplicar, porque su hipótesis (ii) ni siquiera se puede enunciar. Ese es el
mérito real de la intuición del capitán, y es estructural, no poético.

---

## 2 · El medio mínimo, definido antes de mirar nada

Respondiendo a su §5 («no aceptar todavía la palabra fluido como si fuera una
respuesta»), acá está el objeto matemático concreto:

- **Una excitación por primo.** El primo p es una excitación de escala
  característica **log p**, que aporta al medio la escalera de frecuencias
  `ω = 2πk/log p`. (Ésa es la respuesta a su §7: la variable natural que produce
  la escala log p es la frecuencia de la escalera; la potencia p^m aparece como
  **armónico** de la misma excitación, no como órbita nueva — y por eso no hay
  nada que concatenar.)
- **Un campo común.** Todas las excitaciones se acoplan a **un mismo modo del
  medio** con la fuerza que la propia fórmula explícita le asigna a esa
  excitación: `Λ(p)/√p = log p · p^{−1/2}`. **Nada se ajusta.**
- **El operador.** `H = D + g·|v⟩⟨v|`, con D las frecuencias libres y v las
  amplitudes de arriba.

Ese acoplamiento es de **rango uno**, así que el espectro acoplado es **exacto**:
los niveles son las raíces de la secular

    1 = g · Σᵢ vᵢ² / (E − ωᵢ)

una raíz entre cada par de frecuencias libres consecutivas. Sin diagonalizador,
sin truncamiento, y con toda la familia en g disponible de una.

**R6: limpio.** La única entrada es Λ(n). Ningún parámetro se tocó mirando los
γₙ; el espectro se midió antes de compararlo con ellos.

---

## 3 · Lo que el medio SÍ hace

Con 78 excitaciones (primos ≤ 400) y 53 099 modos en [100, 1000]:

| g | Σ²(10) | espaciado mínimo | pegados (< 0,1) | eco k·log p |
|---:|---:|---:|---:|---:|
| 0 (sin acoplar) | 7,5842 | 8,50×10⁻⁶ | 9,16 % | 203,4 |
| 0,01 | 7,3633 | 2,27×10⁻³ | 1,73 % | 66,1 |
| 1 | 6,8716 | 2,27×10⁻³ | 1,66 % | 103,1 |
| 10 | 6,8716 | 2,27×10⁻³ | 1,66 % | 103,4 |
| 100 | 6,8716 | 2,27×10⁻³ | 1,66 % | 103,4 |
| *ceros verdaderos* | *0,3364* | *0,2365* | *0,00 %* | *182,1* |

**Las ondas se sienten entre sí, y eso está medido.** El espaciado mínimo crece
**270 veces** apenas se enciende el acoplamiento, y los niveles pegados caen de
9,16 % a 1,66 % — un factor 5,5. **Hay repulsión genuina, generada por el medio,
que las escaleras desacopladas de la Fase III no tenían.** Ese es un hallazgo
positivo y es del capitán.

Y el experimento que usted pidió en su §8 —agregar excitaciones de a una y ver si
la mejora es colectiva o si sólo se superponen espectros— da colectiva:

| primos | modos | Σ² sin acoplar | Σ² acoplado | mejora |
|---:|---:|---:|---:|---:|
| 2 | 256 | 0,2193 | 0,1015 | 0,1178 |
| 8 | 2 304 | 1,4200 | 0,9383 | 0,4817 |
| 32 | 16 728 | 6,7014 | 5,8397 | 0,8617 |
| 78 | 53 099 | 7,5842 | 6,8716 | 0,7126 |

**La mejora crece con la cantidad de ondas** (0,12 → 0,86). No es superposición:
es interacción.

---

## 4 · Lo que el medio NO hace, y el techo tiene nombre

Σ²(10) baja de 7,58 a 6,87 — y ahí se planta. **Satura**: g = 1, g = 10 y g = 100
dan exactamente lo mismo. Los ceros están en 0,3364, veinte veces más abajo.

La razón es exacta y se puede dibujar: **un acoplamiento de rango uno
INTERCALA.** Mete exactamente una raíz entre cada par de niveles vecinos —
entrelazado de Cauchy. Empuja cada nivel, pero **no puede empujarlo más allá de
sus vecinos inmediatos**: cada nivel queda confinado a su propia celda. Por eso
la repulsión de corto alcance aparece y la rigidez de largo alcance no, y por eso
subir g no sirve de nada pasado cierto punto.

**Lo que muere acá NO es la hipótesis del fluido: muere la versión de UN SOLO
MODO COMÚN.** Y muere con una razón exacta, no por cansancio.

---

## 5 · Veredicto contra su criterio de éxito (§16)

> «Éxito NO es que la simulación se parezca a los ceros; éxito sería que el medio
> genere una estructura que no estaba codificada en la salida.»

| criterio | resultado |
|---|---|
| **R6** | **limpio** — sólo Λ(n), ningún parámetro tocado contra los γₙ |
| **repulsión de corto alcance** | **SÍ, y no estaba codificada**: es efecto del medio (×270 en el mínimo) |
| **rigidez de largo alcance** | **NO** — Σ² 6,87 contra 0,336; techo por entrelazado |
| **eco aritmético** | 103 — pero **TAUTOLÓGICO**, lo declaramos nosotros: log p está en la definición |
| **conteo** | densidad constante: el medio no respira. Frente abierto |

Es un **resultado negativo y útil**, que es la clase que sirve: descarta una
familia concreta con su razón exacta, y deja la hipótesis del capitán viva en su
versión fuerte.

---

## 6 · Lo que este experimento predice para la Fase V

No es una opinión: sale del mecanismo que acabamos de medir.

**Si un solo modo común confina cada nivel a su celda, el medio necesita MUCHOS
modos comunes** — acoplamiento de rango alto — para que las excitaciones se
sientan más allá de la vecina inmediata. Es una predicción medible y barata: se
repite el mismo experimento con `H = D + Σ_{a=1}^{K} g_a |v_a⟩⟨v_a|` y se mira
cómo cae Σ² con K. Si Σ² baja como 1/K, hace falta K comparable al número de
excitaciones — y entonces el «fluido» ya no es un campo con un modo sino algo
con tantos grados de libertad como primos, que es una afirmación estructural
fuerte y contrastable.

Y queda una pregunta que este experimento hace inevitable: **¿de dónde saldrían
esos modos comunes sin violar R6?** La respuesta natural —que cada modo del medio
sea a su vez aritmético— es exactamente el lugar donde el fluido se toca con las
tres puertas abiertas de la Fase III (adeles, no-localidad, espejos). No es
casualidad: **todas son maneras de tener muchos canales a la vez.**

---

## 7 · Autoría y sello

La intuición es de **Jesús Nicolás Astorga**. La formalización como ruta de
investigación es de la auditora. La traducción a un operador concreto, la
medición y el descarte de la versión de un modo son de este taller.

Nada de esto demuestra ni refuta la Hipótesis de Riemann. Nivel de evidencia: 2
(una restricción que descarta una familia, con hipótesis exactas y medición).

**Todavía no.**
