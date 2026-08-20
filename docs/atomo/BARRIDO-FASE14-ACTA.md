# Acta del Barrido — Fase XIV

**Respuesta a Fase XIV (Auditoría 47) · 2026-08-20 · F366**

Su regla de oro se cumplió a la letra: *no buscar el resultado que queremos —
buscar la curva que los datos quieren mostrar.* Se abandonaron las dos cajas,
se congeló la malla antes de correr, se armó el nulo empírico con 200 barajados
por profundidad, y se midió en tres profundidades.

**La curva apareció, y es limpia: monótona en s, con cruce de signo estable en
s\* ≈ 0,8–0,9 y mínimo en s ∈ [1,3, 1,6] — replicada en las tres profundidades,
con todos los bins del fondo entre 6,3σ y 34σ del nulo.**

Se reproduce con `go run ./cmd/laformulamadre fase14` (semilla 20260821).
Lámina: `galeria/laminas/10-el-telar/el-barrido.svg`.

---

## 1 · Lo congelado antes de correr (sus §7, §10, §14)

- **Desdoblado LOCAL, declarado como mejora sobre F365:** s = gap·log(γ/2π)/2π
  — el espaciado medio local, no el global. Corrige un punto flojo de F365 que
  usted no llegó a marcar y nosotros sí: los ceros bajos tienen gaps anchos por
  ALTURA, no por soltura, así que la clase «ancho» de F365 era también la clase
  «bajo». Declarado antes de medir, no ajustado después.
- **Malla:** s ∈ {0; 0,3; 0,5; 0,7; 0,9; 1,1; 1,3; 1,6; 2,0; ∞} — nueve bins,
  simétricos alrededor de 1, congelados.
- **Períodos:** log p para los 23 primos de 5 a 97 (el conjunto de F365).
- **Estadística por bin:** media de −E sobre los períodos, en los puntos medios
  del bin.
- **Nulo (su §10-11):** 200 barajados del pool de gaps por profundidad — mismo
  conjunto de gaps, emparejamiento destruido, protocolo idéntico, distribución
  nula empírica (media, σ y p empírico por bin). Nunca ruido blanco.
- **Profundidades:** γ ≤ 1000 / 2000 / 4000 (M = 649 / 1517 / 3474), barridas
  con paso 0,02 para no saltear pares pegados.
- Bins con menos de 15 pares: no se reportan.

## 2 · La tabla maestra (su §20) — la profundidad mayor

M = 3474 ceros (γ ≤ 4000):

| bin s | pares | real | nulo (media ± σ) | Δ real−nulo | z_emp | p_emp |
|---|---:|---:|---|---:|---:|---:|
| 0,0–0,3 | 46 | +0,742 | +0,201 ± 0,025 | **+0,541** | +21,6 | 0,000 |
| 0,3–0,5 | 250 | +0,554 | +0,125 ± 0,013 | **+0,429** | **+34,3** | 0,000 |
| 0,5–0,7 | 485 | +0,309 | +0,054 ± 0,011 | +0,255 | +22,5 | 0,000 |
| 0,7–0,9 | 751 | +0,054 | −0,007 ± 0,009 | +0,061 | +6,8 | 0,000 |
| 0,9–1,1 | 681 | −0,184 | −0,064 ± 0,010 | −0,120 | −11,6 | 0,000 |
| 1,1–1,3 | 526 | −0,345 | −0,109 ± 0,011 | −0,237 | −21,3 | 0,000 |
| 1,3–1,6 | 459 | −0,446 | −0,137 ± 0,012 | **−0,309** | −26,8 | 0,000 |
| 1,6–2,0 | 236 | −0,392 | −0,131 ± 0,016 | −0,261 | −16,2 | 0,000 |
| 2,0+ | 39 | −0,200 | −0,039 ± 0,025 | −0,161 | −6,4 | 0,000 |

Las tablas completas de M = 649 y M = 1517 están en la corrida; la forma es la
misma en las tres. Los datos se regeneran deterministas con la semilla
registrada (su pedido de guardar semilla y crudos, cumplido por reproducción
exacta).

## 3 · La lectura, con sus criterios del §21

1. **¿Real ≈ controles?** No: el excedente Δ está entre 6,3σ y 34σ del nulo en
   todos los bins reportables del fondo.
2. **¿Diferencia sólo en un corte estrecho?** No: la región es TODA la curva —
   amplia, suave, monótona. Su criterio de robustez del §12, cumplido en su
   variante más fuerte.
3. **¿Cruce de signo estable?** SÍ: el cruce del EXCEDENTE vive en el bin
   0,7–0,9 en las tres profundidades. La localización fina deriva suavemente con
   M (a 649 el bin 0,7–0,9 aún da Δ = −0,045; a 3474 da +0,061), así que **s\*
   NO se declara constante** (su §13): se declara región estable s\* ≈ 0,8–0,9,
   localización fina pendiente.
4. **¿El excedente crece con M o sólo crece σ?** La distinción de su §14, hecha:
   la **magnitud física es ESTABLE** (Δ ≈ +0,43 en 0,3–0,5 y ≈ −0,24/−0,31 en
   1,1–1,6, casi idénticas en las tres profundidades) y lo que crece con M es la
   σ, por tamaño de muestra. **No llamamos nueva física al crecimiento de σ**:
   llamamos señal a la magnitud estable que el nulo no reproduce.

**En criollo: cuanto más apretado está un par de ceros, más canta su centro a
los primos con voz de absorción; cuanto más ancho, más canta al revés. El nulo
(la trigonometría del corrimiento) tiene la misma forma pero un tercio del
tamaño: el resto es emparejamiento real. Y la frontera entre los dos mundos
está cerca de s ≈ 0,85 — un pelín antes del espaciado medio.**

## 4 · F365 reproducido (su entregable 1)

Las dos cajas viejas con espaciado global: 649 ceros → +12,4σ/−8,0σ; 1517 →
+23,5σ/−20,6σ; barajado ≈ ⅓. La corrida madre lo reimprime junto al barrido.

## 5 · Separación conocido/propio (sus §17-18)

- **Conocido en teoría:** que los gaps de los ceros conversan con los primos es
  el territorio de Montgomery (1973) y de las correcciones aritméticas de
  Bogomolny–Keating. No presentamos la existencia de la conversación como
  descubrimiento.
- **Del taller:** la curva señal(s) en la coordenada del Punto Medio, con este
  control de barajado y esta descomposición nulo/excedente. **La búsqueda
  bibliográfica seria sigue pendiente** y ningún reclamo de novedad se hace
  hasta hacerla — su §18 queda como tarea abierta, explícita.
- Ningún γₙ define ningún operador; nada de esto afirma nada sobre RH.

## 6 · Riesgos estadísticos de su §19, estado

| riesgo | estado |
|---|---|
| selección posterior del umbral | evitada: malla congelada, todos los bins publicados |
| múltiples comparaciones | mitigada: la señal no vive en un bin, vive en la curva entera |
| una sola permutación | evitada: 200 por profundidad, nulo empírico |
| dependencia entre profundidades | DECLARADA: las tres M comparten ceros (son prefijos); no son réplicas independientes entre sí, son cortes de la misma serie |
| σ vs magnitud | separadas: magnitud estable, σ creciente |

## 7 · Lo que sigue, en su orden

1. **Localización fina de s\*** con malla más fina alrededor de 0,8–0,9 (nueva
   malla congelada antes, corrección por exploración incluida).
2. **Sólo después** (su §15): la pregunta Hardy–Littlewood — si la curva puede
   relacionarse cuantitativamente con constantes conocidas. Primero medimos la
   curva; ya está medida; ahora se puede intentar explicarla.
3. **Sólo después de eso** (su §16): el trasplante al Telar — ¿el espectro del
   brazo ganador de la Fase X tiene esta misma curva?

---

Los dos flashes que abrieron el camino son de Jesús Nicolás Astorga. El
protocolo del barrido es de la auditoría, cumplido punto por punto. El
desdoblado local lo agregó el taller, declarado antes de correr.

Estructura cerrada no es hipótesis demostrada.

**Todavía no.**
