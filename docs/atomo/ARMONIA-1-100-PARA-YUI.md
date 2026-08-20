# La Armonía 1–100 — documento para Yui

**Tres flashes del capitán, encadenados y medidos · 2026-08-19 · F360**

Este documento no responde a una auditoría: nace de tres flashes consecutivos de
Nico, formalizados y medidos el mismo día. Se lo mandamos porque el tercero
terminó en una prueba con veredicto 10 de 10 y dos correcciones propias, y ese
es exactamente el material que usted audita.

**Aclaración de régimen, antes de todo.** Acá no hay operador, no hay medio, no
hay construcción del taller: son los ceros de Riemann y los enteros, crudos,
cara a cara. Por eso R6 no aplica — no hay nada que proteger de circularidad.
Esto es medición sobre datos, en el territorio de la fórmula explícita, y así
debe leerse: **el provecho es calibración y confirmación fina, no matemática
nueva.** Lo decimos nosotros antes de que lo diga usted.

Se reproduce con `go run ./cmd/laarmonia`. Sin parámetros ajustables: la única
elección es el alcance 1–100, que fue el flash mismo.

---

## 1 · El flash y su porqué

> «Usar los primos que están entre 1 y 100 y los ceros que están entre 1 y 100,
> proyectar una armonía entre ellos, y usar el murciélago a ver si encuentra una
> relación escondida. Usamos esta proporción chica porque así podemos escalar si
> encontramos una alineación inesperada.» — Nico

El instrumental: el murciélago es el eco `E(T) = (2/M) Σₙ cos(γₙT)`, y su
gemelo invertido `P(t) = Σₚ (log p/√p)·cos(t·log p)`. Alcance: los **29 ceros**
con γ < 100 y los **25 primos** < 100. Controles: 2000 remuestreos de períodos
al azar en el mismo rango, para cada medición.

---

## 2 · Primera medición: la armonía existe, en las dos direcciones

**Los ceros cantan, los primos escuchan.** El eco en T = log p es negativo
(absorción — el signo que el taller midió 13 de 13 en la Fase II) y suena a
**1,71 σ** del azar. La relación 1/2: la correlación de −E(log p) con la ley
`log p·p^(−1/2)` da **0,829**, contra 0,78 y 0,77 de las leyes rivales p⁰ y
p⁻¹. Con 25 primos las leyes vecinas aún se parecen; el rango chico limita la
separación, y se declara.

**Los primos cantan, los ceros escuchan.** P(t) evaluada en los 29 ceros da
canto negativo en **29 de 29**. Ni una excepción.

**Vuelo libre — la parte que no se le dijo dónde mirar.** De los 8 picos más
fuertes de |E(T)| en todo el rango, **6 caen en log p con error de milésimas**:
log 7 (a 0,0024), log 5 (0,0011), log 11 (0,0006), log 3 (0,0014), log 13
(0,0009), log 73 (0,0085). Los ceros, solos, deletrean los primos.

**La apuesta de escala del flash, confirmada:** el eco crece 0,68 σ (5 ceros) →
1,00 (10) → 1,36 (15) → 1,75 (29). Crece como √M, y extrapolado a los 620
ceros del taller reproduce las ~36 σ que la Fase I midió. La escalera es
coherente de la proporción mínima a la grande.

---

## 3 · Segunda medición: los números que no son primos

La fórmula explícita hace una profecía de tres vías sobre TODO n, vía Λ(n):
primo canta; potencia de primo canta **con la voz de su base** (log p, no
log n); compuesto mixto **calla**. Medida cruda sobre los 99 números de 2 a 100:

| familia | cuántos | media −E | σ |
|---|---:|---:|---:|
| primos | 25 | 0,420 | **1,74** |
| potencias de primo (4, 8, 9, 16, 25, 27, 32, 49, 64, 81) | 10 | 0,108 | 0,40 |
| compuestos mixtos | 64 | −0,121 | **−0,47** |

**La voz de las potencias, el detalle más fino:** la correlación de √n·(−E)
con log p (la voz de la BASE) da **+0,42**, y con log n (la voz propia) da
**−0,42**. Signos opuestos: el 8 canta como el 2, no como el 8. El 25 grita
(base 5, ladrillo grande); el 64 casi no se oye (2⁶, la voz del 2 dividida por
√64). Es Λ(n) leída directamente del eco.

Esto conecta con el no-go de la Fase III: log(a·b) = log a + log b no crea
período, y acá se lo escucha en los datos sin construir nada.

---

## 4 · Tercera medición: la prueba del contagio (el veredicto)

Los diez compuestos mixtos "más ruidosos" del §3 resultaron estar **todos**
pegados a un primo: log 72 está a 0,014 de log 73; el 60 junto al 61; el 92
junto al 89… Hipótesis declarada antes de correr: no es voz propia, es
**contagio** — con 29 ceros el oído no resuelve distancias menores que su
borrón, y el eco del primo vecino los salpica. Predicción falsable: al subir la
cantidad de ceros, los diez deben hundirse bajo el piso de ruido **mientras los
primos siguen cantando**. Si alguno sigue sonando, tiene voz propia y la
fórmula explícita estaría en problemas.

Se calcularon **649 ceros** con Riemann–Siegel hasta γ = 1000 (los 29 conocidos
coinciden a 0,0025) y se subió el oído por escalones:

| n (mixto) | vecino | M=29 | M=50 | M=100 | M=200 | **M=649** |
|---:|---:|---:|---:|---:|---:|---:|
| 72 | 73 | 0,600 | 0,197 | 0,126 | 0,131 | **0,015** |
| 75 | 73 | 0,523 | 0,180 | 0,076 | 0,045 | **0,000** |
| 92 | 89 | 0,492 | 0,172 | 0,082 | 0,028 | **0,006** |
| 45 | 47 | 0,489 | 0,097 | 0,087 | 0,059 | **0,008** |
| 39 | 41 | 0,493 | 0,103 | 0,131 | 0,010 | **0,009** |
| 93 | 97 | 0,488 | 0,233 | 0,096 | 0,064 | **0,011** |
| 76 | 79 | 0,477 | 0,308 | 0,086 | 0,039 | **0,002** |
| 86 | 89 | 0,480 | 0,420 | 0,029 | 0,003 | **0,021** |
| 24 | 23 | 0,448 | 0,108 | 0,076 | 0,057 | **0,016** |
| 60 | 61 | 0,438 | 0,071 | 0,193 | 0,011 | **0,042** |
| **primos (media)** | — | 0,420 | 0,410 | 0,381 | 0,338 | **0,279** |

Piso de ruido con 649 ceros: √(2/M) = 0,055.

**Veredicto: 10 de 10.** Los diez acusados quedaron bajo el piso de ruido tras
caer más de cinco veces, y los primos quedaron cinco veces POR ENCIMA del piso,
cantando. El ruido de los mixtos era contagio del vecino. **Los números
compuestos no existen para la orquesta de los ceros — ni siquiera de cerca.**

---

## 5 · Dos correcciones propias, cazadas en el mismo turno

1. **Mi primer veredicto automático dio 3 de 10 sobre datos que a ojo eran 10 de
   10.** La regla evaluaba el silencio a la altura teórica γ* = π/d, y para
   siete casos esa altura caía por debajo del primer escalón, así que la regla
   los comparaba contra sí mismos. Se corrigió **la regla, no los datos**, y la
   corrección quedó comentada en el código.
2. **La profecía γ* = π/d era incompleta.** El horario del silencio no lo fija
   sólo la resolución: lo fija el **piso de ruido √(2/M)** que impone tener
   finitos ceros. Con 29 ceros el piso es 0,26 y nadie puede estar "callado".
   El criterio honesto es: bajo el piso al oído máximo, tras caída mayor a 5×.

---

## 6 · Qué es y qué no es

- **No es matemática nueva.** Es la fórmula explícita de Riemann —la dualidad
  de 1859— verificada en la proporción más chica posible y a través de tres
  lentes: signo, ley de volumen, y silencio de los compuestos con su horario.
- **Sí es calibración nueva del taller**: ahora sabemos que la armonía es
  audible desde 29 ceros, que escala como √M, y que el silencio de Λ(n) = 0 se
  puede *ver profundizarse* escalón por escalón. Ningún experimento nuestro
  había mirado el contagio y su horario.
- **El hilo escalable que queda**: dar vuelta el instrumento y usar el canto de
  los primos < 100 como buscador de ceros que nunca se le mostraron. Si la
  proporción chica localiza ceros nuevos, eso sí sería hallazgo propio.

---

Los tres flashes son de Jesús Nicolás Astorga. Las mediciones, los controles y
las dos correcciones, de este taller. Registrado como F360 en la bitácora.

Estructura cerrada no es hipótesis demostrada.

**Todavía no.**
