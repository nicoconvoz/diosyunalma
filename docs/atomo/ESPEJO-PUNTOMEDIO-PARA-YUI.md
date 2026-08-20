# El Espejo del Punto Medio — documento para Yui

**F365 · 2026-08-20 · dos flashes del capitán, encadenados sobre la Fase XIII**

Este documento no responde a una auditoría: nace de dos flashes de Nico dados el
mismo día, y termina en una señal que **sobrevivió a su control decisivo y
creció al replicarla con el doble de ceros**. Se lo mandamos porque de acá puede
partir una fase nueva, y queremos su diseño antes de dar un paso más.

**Régimen:** datos crudos — ceros reales de Riemann (calculados con
Riemann–Siegel en el propio taller) y enteros. No hay operador ni construcción,
así que R6 no aplica; es medición en el territorio de la fórmula explícita.

Se reproduce con `go run ./cmd/laformulamadre` (corre las dos partes seguidas).

---

## 1 · El primer flash: el Punto Medio adentro de la fórmula madre

> «Incluí el teorema del punto medio en la fórmula madre y usá el murciélago —
> buscala.» — Nico

Traducción: ¿los ceros escuchan los **centros** de los primos gemelos (12, 18,
30, … — la red 6ℤ del Teorema 6)? Tres sondas sobre 649 ceros:

| capa | primos (control positivo) | centros gemelos | pareja log(p·q) |
|---|---:|---:|---:|
| eco de un cero, −E(T) | **+24,9 σ** | −1,7 σ | −1,3 σ |
| factor de forma K(T) (de a dos) | **+23,3 σ** | −0,7 σ | — |

**Resultado: silencio en las dos capas.** Coherente con la ley Λ que el taller
midió en F360 — los centros son todos compuestos, y la orquesta no conoce
palabras, sólo letras. Hasta se vigiló el contagio (el 12 vive a 0,08 de
log 11 y log 13) y aun así calla. **Nulo limpio con deslinde preciso**: el
Teorema del Punto Medio vive en una capa que los ceros no escuchan de frente;
si entra en la fórmula madre, no entra por el período del centro.

## 2 · El segundo flash: dar vuelta el teorema, sobre los ceros

> «En vez de usar los pares normales, usá la fórmula del teorema del punto
> medio CON CEROS — quiero ver qué aparece en esos casos.» — Nico

Tres trasplantes: los **puntos medios** de ceros consecutivos
mₙ = (γₙ + γₙ₊₁)/2, las **anclas** aₙ = (γₙ−1)/2 (la coordenada del teorema,
que divide todo gap por 2), y los **ceros gemelos** (pares con gap < 0,7 del
espaciado medio) contra los **anchos** (gap > 1,3).

### 2a · Lo esperado, que confirma la coordenada

- Los puntos medios **invierten** la voz de los primos (−68% del eco crudo):
  es el factor de fase cos(gap·T/2) del corrimiento — trigonometría, declarada
  como tal.
- Las anclas hacen **la mudanza exacta que el álgebra manda**: la voz
  desaparece de T = log p (−0,3 σ) y reaparece en T = 2·log p = **log p²**
  (−18,9 σ). La coordenada que divide gaps por 2 corre el eco una octava —
  sobre los ceros igual que sobre los impares.

### 2b · ⚡ Lo que apareció: los ceros gemelos cantan al revés de los anchos

| centros de pares | 649 ceros | **1517 ceros (réplica)** | control: gaps barajados |
|---|---:|---:|---:|
| **apretados** (gap < 0,7·medio) | +12,4 σ | **+23,5 σ** | +1,9 σ / +6,2 σ |
| **anchos** (gap > 1,3·medio) | −8,0 σ | **−20,6 σ** | −2,9 σ / −8,3 σ |

Lectura de la tabla, con todas las letras:

- Los centros de pares **apretados** cantan los primos con el signo de la
  **absorción** — casi a plena voz (−E = +0,25 contra +0,28 de los ceros
  crudos), a pesar del corrimiento.
- Los centros de pares **anchos** cantan **invertido** (emisión).
- **El control decisivo**: reconstruir cada centro con un gap BARAJADO del
  mismo conjunto (idéntico factor de fase, emparejamiento destruido) derrumba
  la señal a un tercio, en las dos profundidades. La parte que sobrevive al
  barajado es trigonometría; **el excedente ×3 es emparejamiento real**: el
  gap de un par sabe en qué fase de las ondas de los primos está parado.
- **Réplica**: al doblar los ceros (649 → 1517) la señal **crece** en los dos
  brazos y el control queda proporcionalmente igual de abajo.

**En criollo: dónde el coro de los ceros se aprieta, lo deciden los primos.**

## 3 · La honestidad, antes de que usted la pida

1. **La física de fondo es conocida en teoría.** Que los gaps de los ceros
   conversan con los primos es el territorio de la correlación de pares de
   Montgomery (1973) y de las correcciones aritméticas de Bogomolny–Keating,
   donde entran las constantes de Hardy–Littlewood de los gemelos. No
   afirmamos haber descubierto esa física.
2. **Lo propio del taller es la FORMA**: el trasplante del Punto Medio la
   convierte en un partidor limpio de dos signos (+/−) con control interno
   (el barajado de gaps), medible desde 649 ceros con estadística enorme. No
   conocemos esta presentación en dos brazos con ese control — pero no hemos
   hecho búsqueda bibliográfica seria, y lo decimos.
3. Los umbrales 0,7 y 1,3 se eligieron ANTES de mirar (simetría alrededor del
   espaciado medio), pero no se barrieron: la dependencia con el umbral está
   sin medir.
4. El §2a es álgebra confirmada, no hallazgo, y así está rotulado.
5. Nada de esto afirma nada sobre RH.

## 4 · De acá puede partir algo nuevo — lo que le proponemos diseñar

1. **El barrido del umbral**: medir la señal como función continua del corte de
   gap (no dos clases sino el espectro entero de clases). Si el cruce de signo
   ocurre en un gap característico, ese gap es un número nuevo del taller.
2. **La pregunta Hardy–Littlewood**: la teoría dice que la conversación
   gap↔primos lleva las constantes de los gemelos adentro. ¿Se puede LEER una
   constante de Hardy–Littlewood desde nuestro partidor de dos brazos? Sería
   la primera vez que el taller extrae un número aritmético profundo de los
   ceros por esta vía.
3. **El puente con el Telar**: la Fase X terminó con una pregunta viva («¿por
   qué el signo de enlace no localiza y el de sitio sí?») y este espejo sugiere
   dónde mirar: en nuestro operador, ¿los pares apretados de NUESTROS niveles
   también saben la fase de algo? Un partidor gemelos/anchos sobre el espectro
   del brazo ganador de la Fase X costaría una corrida.
4. **Más profundidad**: γ = 2000 fue la réplica; el taller navega mucho más
   hondo. La curva señal(M) contra el control(M) decidiría si el excedente ×3
   es constante o crece.

Elija usted el orden, o descarte lo que no sirva. Los dos flashes son de Jesús
Nicolás Astorga; las mediciones, los controles y las declaraciones de límite,
de este taller.

Estructura cerrada no es hipótesis demostrada.

**Todavía no.**
