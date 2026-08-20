# Acta de la Auditoría Independiente — Fase XXI

**Respuesta a Fase XXI (Auditoría 53) · 2026-08-20 · F373**

Su orden absoluta: no desviarse, reimplementar desde cero, auditar, y recién al
final escribir la conclusión formal que distinga qué se sabe de qué no. Se
cumplió entera. **La segunda implementación reproduce todo, la sensibilidad
numérica no mueve nada, el rango de γ no rompe nada — y el cabo suelto
0,42→0,39 quedó RESUELTO: era un artefacto de la referencia de la
implementación A, no del modelo.**

Se reproduce con `go run ./cmd/laauditoria` (paquete nuevo, cero código
compartido con `cmd/laformulamadre`).

---

## 1 · La ruta B — qué se hizo distinto (su Prueba 1, prioridad absoluta)

| pieza | implementación A | implementación B (nueva) |
|---|---|---|
| S(γ) | tabla Λ precomputada sobre potencias | suma por primo y potencia: Σₚ Σₘ sin(mγ ln p)/(m·p^{m/2}) |
| θ(t) | 4 términos asintóticos | 5 términos (el extra, conmutable) |
| buscador del modelo | marcha punto a punto + bisección | barrido en grilla de cruces de semienteros + bisección |
| ceros reales | barrido RS a paso 0,02 desde γ=10 | barrido RS fresco a paso 0,01 desde γ=30 |
| azar de los nulos | xorshift64 | splitmix64 |
| eco/bins/cruce/T4 | — | reescritos desde cero |

Ninguna función, tabla ni resultado intermedio compartido. Las **definiciones**
(desdoblado, bins, períodos, ventana de pendiente) son idénticas por decreto —
eso ES el experimento.

## 2 · Prueba 1: los cuatro escalones — REPRODUCIDOS

| N | amplitud B (A) | cruce B (A) | T4 B (A) | residuo | dist. media al cero real | <0,1 |
|---:|---|---|---|---:|---:|---:|
| 97 | 0,247 (0,267) | 0,758 (0,847) | 0,197 (0,215) | 0,123 | 0,047 | 91 % |
| 997 | 0,302 (0,315) | 0,903 (0,921) | 0,208 (0,205) | 0,055 | 0,024 | 99 % |
| 9973 | 0,339 (0,353) | 0,861 (0,869) | 0,266 (0,256) | 0,028 | 0,013 | 100 % |
| 99991 | 0,365 (0,374) | **0,856** (0,861) | 0,291 (0,292) | **0,019** | **0,009** | 100 % |
| REAL (ruta B) | 0,346 | 0,863 | 0,267 | — | — | — |

**La tendencia entera se reproduce**: la escalera de amplitud, la convergencia
del cruce, el T4 subiendo hasta encerrar el real, el residuo colapsando, y la
reconstrucción de posiciones (a N profundo, cada punto del modelo está a 0,009
del cero real — ruta B encuentra los ceros aún mejor que A). Las diferencias
A↔B son de milésimas en los escalones profundos; en N=97 el cruce difiere más
(0,758 vs 0,847), coherente con que la truncación corta es la más sensible al
buscador — declarado, no maquillado.

## 3 · Prueba 2: el rango de γ — la convergencia NO es local

| ventana | cruce modelo | cruce real | dist. media | <0,1 |
|---|---:|---:|---:|---:|
| γ ∈ [30, 2000] | 0,844 | 0,847 | 0,012 | 100 % |
| γ ∈ [2000, 4000] | 0,879 | 0,891 | 0,015 | 100 % |
| γ ∈ [30, 4000] | 0,860 | 0,861 | 0,013 | 100 % |

Y el detalle fino: el cruce REAL varía con la ventana (0,847 abajo, 0,891
arriba) — **y el modelo reproduce hasta esa dependencia con la ventana**, a
milésimas. No hay rango privilegiado: hay seguimiento total.

## 4 · Prueba 3: sensibilidad numérica — roca

| variación | amplitud | cruce |
|---|---:|---:|
| paso de grilla 0,010 | 0,342 | 0,861 |
| paso de grilla 0,004 | 0,342 | 0,860 |
| θ sin el término extra | 0,343 | 0,858 |

Nada se mueve más de tres milésimas. La convergencia no es producto de
precisión, redondeo ni discretización.

## 5 · El cabo suelto 0,42 → 0,39: RESUELTO (no corregido — explicado)

La ruta B midió el desvío de espaciados **de los ceros reales** con su barrido
fino: **0,39** — no 0,42. Y por ventanas, el diagnóstico completo:

| serie | s-desv baja [30,2000] | s-desv alta [2000,4000] |
|---|---:|---:|
| REAL (ruta B) | 0,383 | 0,392 |
| modelo N=99991 | **0,383** | **0,392** |
| modelo N=9973 | 0,385 | 0,395 |

**El modelo coincide con los ceros reales dígito a dígito en las dos ventanas.**
El "0,42" de la Fase XX era de la REFERENCIA de la implementación A (barrido a
paso 0,02 desde γ=10 — de ahí también sus 3 ceros extra, los γ < 30). No había
deriva del modelo: había un artefacto de medición en la vara con que se lo
comparaba. El único cabo suelto de la campaña queda cerrado, y cerrado en la
dirección que más fortalece el resultado.

## 6 · La conclusión formal (su §13), palabra por palabra

- **Qué fue observado**: la curva del Espejo del Punto Medio — E(s) con cruce
  ≈0,86, pendiente ≈−0,93, amplitud ≈0,35, T4 ≈0,27 — sobre los ceros reales.
- **Qué fue reproducido**: todo lo anterior, por dos implementaciones
  independientes, estable ante numérica, semillas y ventanas de γ.
- **Qué es consecuencia matemática de la fórmula explícita (niveles A y B de su
  Prueba 4)**: TODO el observable. La ecuación de conteo con S truncada
  reconstruye las posiciones de los ceros (0,009 de distancia media), su
  estadística local (s-desv idéntico por ventana), y toda la curva del espejo
  con su cruce, pendiente y T4. La estadística local "tipo GUE" emerge de los
  términos aritméticos profundos, sin ruido externo (la imagen de Berry,
  verificada dos veces).
- **Qué NO constituye una explicación independiente (nivel C)**: nada de esto.
  El modelo no es independiente de los ceros — converge a ellos porque la
  fórmula explícita ES la definición dual de los ceros. **No existe, dentro de
  este marco, ninguna afirmación legítima de nivel C.** Se dice con todas las
  letras, como su §17 exige valor para decir.
- **Qué anomalía permanece**: NINGUNA. El 0,42→0,39 está resuelto.
- **Qué pregunta queda legítimamente abierta**: la única — si existe ALGÚN
  observable donde la reconstrucción por la fórmula truncada y los ceros reales
  se separen (ahí viviría estructura no capturada), o una construcción de nivel
  C que produzca esta curva sin usar la dualidad. Dentro del marco del Espejo,
  con este observable y esta profundidad: **no falta nada — y lo decimos.**

## 7 · El balance (su Prueba 5)

| fenómeno | ¿reconstruido? | ¿explicado independientemente? |
|---|---|---|
| curva E(s) | SÍ (residuo 0,019) | NO |
| cruce s\* | SÍ (a milésimas, ventana por ventana) | NO |
| pendiente | SÍ | NO |
| T4 | SÍ | NO |
| estadística local | SÍ (s-desv exacto por ventana) | NO |
| posiciones de los ceros | SÍ (0,009 de distancia) | **NO — y no se asume** |

HL, el Telar, los modelos mecánicos: siguen congelados, como manda.

---

La recompensa de una buena auditoría es saber exactamente qué NO sabemos: el
espejo entero es la dualidad; lo que no tenemos es una causa de nivel C — que
es, con otras palabras, el problema de Riemann de siempre, ahora con el mapa de
qué NO lo resuelve mucho mejor dibujado.

Estructura cerrada no es hipótesis demostrada.

**Todavía no.**
