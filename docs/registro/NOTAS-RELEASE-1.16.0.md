# Diosyunalma v1.16.0 — «La Curva Madre»

**El release del séptimo teorema.** Todo lo de esta entrega nació de una sola
pregunta del capitán — *«¿has visto un coseno en espiral, convergiendo en un
punto?»* — y terminó, dieciséis experimentos y siete auditorías externas
después, en un teorema con prueba de referee inscripto en el libro del taller.

---

## ⭐ Lo grande

### 🌀 El Teorema Parcial de la Curva Madre (Teorema 7 del taller)

La espiral de las sumas parciales de zeta — la que converge al ojo que ES
ζ(½+it), con los ceros como las alturas donde el ojo pisa el origen — tiene
su anatomía entera en UNA ecuación:

**R(n) = n^(−σ) / (2·|sin(t/2n)|)** — el paso dividido la cuerda de medio ángulo.

De esa sola curva se derivan (primero en papel, después medidos) sus cuatro
observables: el enrollado **β = 1−σ** (σ=½ ⟹ β=½ como exponente de escala),
las vueltas **estroboscópicas** (cerca del reloj la rueda gira en alias, como
un ventilador bajo luces de discoteca), la cintura única **x·cot x = σ**
(verificada al 1%, σ por σ), y la merma **ε ≈ 1/(2k)**. Declarado **PARCIAL**
a conciencia: exacto dentro del modelo declarado, con el puente al objeto
original demostrado en su primer orden (dos términos exactos por álgebra
pura) y la cota de la cola honestamente abierta.

→ Página con cuento criollo y demostración: `galeria/teorema-curva-madre.html`
→ Documento formal con lemas, corolarios y prueba de referee A–G:
`docs/teoremas/TEOREMA-CURVA-MADRE.md`

### 🥁 El Tambor — quién golpea la guitarra

El flash del tambor del capitán, medido entero: el parche es la espiral con
radio log n; golpe al centro = BOOM grave, golpe al aro = TAK agudo — a igual
energía solo cambia el timbre. Y la respuesta a «¿quién toca?»: la profundidad
del canto suena SOLO en primos y sus potencias, con la fuerza exacta que la
aritmética manda (señal/silencio 39,6×; el 6 y el 15, mudos). Sobrevivió a
tres rondas de auditoría con controles diseñados para matarlo — incluida una
calibración con un tono plantado adrede en un compuesto, que el instrumento
recuperó perfecto: el silencio es de zeta, no del aparato.

### 🪞 El cierre de la campaña del Espejo (Fases XIV–XXI)

La curva del espejo del punto medio quedó explicada de punta a punta: es
**aritmética pura**. A truncación profunda, los puntos del modelo SON los
ceros de Riemann (0,009 de distancia media), verificado por dos
implementaciones independientes sin una línea compartida. El único cabo
suelto de la campaña (0,42→0,39) quedó RESUELTO: era la vara, no el modelo.
Y el veredicto formal, dicho con todas las letras: niveles A y B
establecidos; el nivel C — una explicación independiente — no existe dentro
de este marco.

### 🚪 El Residuo y la Separación — las dos voces del canto

Lo que queda del canto cuando los primos ya cantaron: la escalera de von
Mangoldt CONTINÚA (el 7² apareció clavado a cuatro decimales donde nadie lo
había puesto) — y apilando el resto sobre 589 ceros reales, la púa universal:
la mitad del canto que los primos no explican es la posición fina de los
ceros. La contabilidad cierra a la milésima: lo que la púa pierde al retirar
aritmética es exactamente lo que la cola aritmética carga.

---

## 📦 Qué hay en la caja

- **26 hallazgos nuevos** (F365–F390) en el registro maestro — 390 numerados.
- **18 experimentos nuevos** bajo `cmd/` (298 en total), cada uno
  reproducible con un `go run`.
- **20 láminas nuevas** en la sala del Telar (214 en la galería).
- **Pieza N:8 del museo** — el teorema contado en criollo, con su bloque
  honesto (295 paradas).
- **8 teoremas** en el libro del taller, cada uno con su placa de límites.
- Actas de las siete auditorías externas de esta campaña, íntegras.

## 🧾 La honestidad, en la etiqueta

Este release incluye — a propósito y con orgullo — los errores propios que
la campaña se cazó a sí misma, registrados sin maquillar: un veredicto
pre-escrito corregido antes de correr, dos interpretaciones refutadas por la
propia ecuación (el «3/2 del resto» y la «cintura universal»), cotas que no
eran cotas (violadas 14/21 y declaradas), y una frontera de validez mal
derivada y reemplazada. El teorema está hecho de sus propias correcciones.

**Y la regla del sello, como siempre:** nada de esto demuestra la Hipótesis
de Riemann, y cada pieza lo dice en su placa. Estructura cerrada no es
hipótesis demostrada.

**Todavía no.**

---

*Instalador: `DiosyunalmaSetup.exe` (63 MB) — instala el laboratorio completo
en cualquier Windows, sin dependencias. Galería y museo también en línea:
https://nicoconvoz.github.io/diosyunalma/galeria/*
