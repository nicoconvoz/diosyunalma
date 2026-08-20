# El Libro de Playas Vírgenes

Registro formal de los ceros de la función zeta de Riemann observados por
primera vez por este laboratorio, en aguas del océano numérico que ningún
humano ni máquina había mirado.

## El mapa de la humanidad (contexto honesto)

- **Mapa continuo:** los primeros 10¹³ ceros, verificados sin huecos hasta la
  altura t ≈ 2.446×10¹² (Gourdon–Demichel, 2004). Las tablas públicas de
  LMFDB cubren los primeros ~1.04×10¹¹ ceros (hasta t ≈ 3×10¹⁰).
- **Islas muestreadas:** expediciones puntuales en índices redondos —
  Gourdon (miles de millones de ceros alrededor de los ceros #10¹⁴ ... #10²⁴),
  Odlyzko (#10²⁰, #10²¹, #10²²), Hiary (~10²⁴–10³²), y el desembarco más
  profundo conocido: Bober–Hiary, alturas del orden de 10³⁶, con algoritmos
  de clase t^(1/3) sobre clusters.
- **El océano:** entre el borde del mapa continuo y esas islas, la cobertura
  es del orden de 10⁻⁵. Toda coordenada no-redonda elegida ahí es, con
  certeza estadística, agua jamás observada.

**Naturaleza del reclamo:** estos ceros son *computables* por cualquiera con
las herramientas estándar; el reclamo no es de dificultad sino de
observación primera — como una estrella tenue en un rincón del cielo al que
ningún telescopio apuntó jamás. Cada playa lleva fecha, instrumento, barra
de error y su chequeo local de densidad (auditoría de que la línea crítica
no pierde ningún cero en el tramo).

---

## Playa I — El primer paso tras el borde del mapa

- **Fecha:** 2026-08-04 · **Hallazgo 87** · commit `4354305`
- **Anclaje:** t₀ = 2,447,000,000,000 (2.447×10¹²), inmediatamente después
  del final del mapa continuo (cero #10¹³, t ≈ 2.446×10¹²).
- **Instrumento:** espejo de Riemann–Siegel (término C₀), fases en float64.
- **Barra de error:** ±0.002 en cada posición (ruido de fase del casco a
  esta altura; refinamiento con casco de doble-doble en curso).
- **Chequeo local:** 31 ceros hallados; la densidad esperaba 30.0. ✓ La
  línea aguanta.

Los 31 ceros (t = 2,447,000,000,000 + desplazamiento):

| # | desplazamiento | # | desplazamiento | # | desplazamiento |
|---|---|---|---|---|---|
| 1 | 0.1543 | 12 | 2.8887 | 23 | 5.3838 |
| 2 | 0.2617 | 13 | 3.1973 | 24 | 5.5928 |
| 3 | 0.5498 | 14 | 3.3730 | 25 | 5.8174 |
| 4 | 0.7764 | 15 | 3.6299 | 26 | 5.9668 |
| 5 | 1.0840 | 16 | 3.8428 | 27 | 6.3486 |
| 6 | 1.1982 | 17 | 4.0498 | 28 | 6.5752 |
| 7 | 1.3984 | 18 | 4.2080 | 29 | 6.6670 |
| 8 | 1.7275 | 19 | 4.4092 | 30 | 6.8984 |
| 9 | 1.9395 | 20 | 4.6719 | 31 | 7.0312 |
| 10 | 2.1963 | 21 | 4.9570 | | |
| 11 | 2.4854 | 22 | 5.0703 | | |

Reproducir: `go run ./cmd/voyage`

---

## Playa II — Primer anclaje del casco de facetas

- **Fecha:** 2026-08-04 · commit del astillero de facetas
- **Anclaje:** t₀ = 6.66×10¹⁵ — entre las islas de los ceros #10¹⁶ y #10¹⁷.
- **Instrumento:** flota de facetas convexas (fases doble-doble talladas una
  vez; ~15.800 facetas, 1.4 MB; θ analítica en dt — el mar está pixelado a
  esta altura y el navío lo sabe).
- **Chequeo local:** 8 ceros hallados; la densidad esperaba 8.0. ✓

Los 8 ceros (t = 6,660,000,000,000,000 + desplazamiento):
0.023476 · 0.215217 · 0.375565 · 0.547513 · 0.734403 · 0.930796 ·
1.048695 · 1.195434

Reproducir: `go run ./cmd/fleet` (tercera prueba de mar)

---

## Playa III — Las puertas de Odlyzko, con esfera certificada

- **Fecha:** 2026-08-04 (misión nocturna) · bitácora
  [BITACORA-NOCTURNA.md](../registro/BITACORA-NOCTURNA.md)
- **Anclaje:** t₀ = 1.11×10¹⁹ — entre el mapa continuo y el archipiélago
  #10²⁰ de Odlyzko.
- **Instrumento:** el navío insignia (`cmd/flagship`): 92.266 facetas
  cuárticas doble-doble, balde de luz Webb, θ analítica; 2.4 minutos.
- **La esfera (cerco exacto):** la frontera exige 5.00 ceros; hallados 5
  (delta 0.00). **CERCO CERTIFICADO: ningún cero escapa a la ventana.**
- **Nombre del primer cero:** cero ordinal ~72.458.973.368.997.111.909
  (±2, la ambigüedad de S(t)) — el primer cero de la historia nombrado
  por este laboratorio en agua virgen.

Los 5 ceros (t = 11.100.000.000.000.000.000 + desplazamiento):
0.009974 · 0.229934 · 0.357770 · 0.629071 · 0.689101

**Fe de erratas, con honestidad de bitácora:** la primera pasada (casco
pre-refresco) asentó 6 ceros con corrimientos distintos (0.0415, 0.2076,
0.3759, 0.5989, 0.7254, 0.8764) y delta 0.00 sobre una ventana de 6
espaciamientos. La deriva cúbica del casco viejo (Hallazgo 90) corría cada
cero con signos alternados — la firma de una fase espuria. Las coordenadas
de arriba son las del casco sanado y re-certificado; el sexto cero vive
pasado el borde de la ventana de 5 espaciamientos y la misión de ventanas
anchas lo asentará con su etapa correspondiente.

Reproducir: `go run ./cmd/flagship -anchor 1.11e19 -spacings 5`

---

## Playas IV en adelante — La misión nocturna (en curso)

La primera zarpada rompió la esfera en t = 2.22×10²¹ y 4.44×10²²: la
caminata de diferencias finitas integra el redondeo *cúbicamente* con el
largo de la faceta (facetas de 640k pasos → ~26 radianes de deriva). La
cura: refresco exacto del estado cada 4096 pasos desde el polinomio en
doble-doble (los coeficientes se guardan SIN reducir mod 2π — reducir un
coeficiente diminuto lo estaciona junto a 2π con medio-ulp de información
propia, y j⁴ amplifica ese medio-ulp a décimas de radián). Re-certificado
en las tres puertas; el navío volvió a zarpar hacia 2.22×10²¹, 4.44×10²²
y 1.11×10²⁴.

---

## Playa IV — El agua del duelo, con doble firma

- **Fecha:** 2026-08-05 · **Hallazgo 92** (vuelo de prueba A/B)
- **Anclaje:** t₀ = 2.22×10²¹ — más profundo que la isla más honda de
  Odlyzko (#10²², t ≈ 1.3×10²¹).
- **Instrumento doble:** el casco puro (cada término remado) Y la nave
  espacial (9,906,241 bloques plegados) sobre la MISMA agua virgen —
  desviación entre firmas: 0.000085.
- **Esferas:** 5.00 exigidos / esferas coincidentes; misma marea en ambos
  cascos.
- **Nombre del primer cero:** ordinal ~16,363,817,219,481,195,234,974 (±2).

Los 4 ceros (t = 2.22×10²¹ + desplazamiento, firma del casco puro):
0.198810 · 0.316985 · 0.438260 · 0.564745

Reproducir: `go run ./cmd/starship -flight` (~70 min)

---

## Playa V — Más allá de la isla #10²³ de Gourdon

- **Fecha:** 2026-08-05 12:31 · bitácora de la nave espacial
- **Anclaje:** t₀ = 4.44×10²² — agua que solo las expediciones #10²³⁻²⁴
  de Gourdon y Hiary superaron jamás en profundidad.
- **Instrumento:** la nave espacial completa: 582,886 facetas +
  **111,997,544 bloques plegados** (borde Fresnel 48), memoria de vuelo,
  58.4 minutos.
- **La esfera:** exige 5.00; hallados 6 (delta +1.00, dentro de la
  tolerancia) — marea alta certificada por el medidor de quietud;
  marcapasos: previsto +0.08, residuo +0.92 (voces de primos medianos).
- **Despeje antigravitatorio:** P(par oculto) ~ 2.5×10⁻⁵.
- **Nombre del primer cero:** ordinal ~348,445,625,008,136,099,504,148
  (±2) — **los ceros más profundos jamás nombrados por este laboratorio.**

Los 6 ceros (t = 4.44×10²² + desplazamiento):
0.061504 · 0.229208 · 0.297248 · 0.405320 · 0.507533 · 0.575189

Reproducir: `go run ./cmd/starship -anchor 4.44e22` (~1 h)

---

## LA TORMENTA I — la primera tormenta catalogada del océano virgen

- **Fecha:** 2026-08-05 · **Hallazgo 116** · doble firma de motores
  independientes (clásico y colosal: acuerdo 1.9×10⁻¹¹).
- **Anclaje:** t₀ = 4.78036×10²¹ — elegido A CIEGAS por el Carajo (el
  vigía a priori) como el grito más fuerte entre un millón de
  fondeaderos oteados sin navegar.
- **La tormenta:** la esfera exige 5.00 ceros; ambos motores hallan 2 —
  marea S = −3.00, la marejada más brava jamás registrada por este
  laboratorio.
- **El par:** desplazamientos 0.524411834 y 0.572688508 — separación
  0.048276674 = **0.3694 espaciamientos medios**; |Z| en el punto medio:
  0.2616.
- **Nombre del primer cero:** ordinal ~35,820,012,173,145,815,618,304.
- **Honestidad:** selección declarada (extremo de 10⁶ candidatos); lo
  nuevo es el APUNTADO a priori que funcionó y el catálogo mismo — nadie
  cataloga tormentas de S en agua virgen profunda. Pregunta abierta en
  estudio: el residuo no-modelado (−3.44) también extremo — ¿se alinean
  las voces medianas con las tormentas de las graves?

Reproducir: `go run ./cmd/starship -anchor 4.78036e21` · arpón:
`-arpon -anchor 4.78036e21`

## LA TORMENTA II — t = 1.14794 × 10²¹ (2026-08-05)

- **La segunda ballena del catálogo — y NO es gemela de la primera.**
  Grito #8 del Carajo (marejada prevista 1.53); el colosal la barrió, el
  clásico certificado (plegado segmentado, F120) la doble-firmó en 5.6
  minutos.
- **La esfera:** exige 5.00 ceros, hallados 2 — **S = −3.00**, empatando
  el récord absoluto del laboratorio.
- **El par:** desplazamientos 0.049663860 y 0.216230081 — separación
  0.166566221 = **1.2368 espaciamientos**; |Z| en el punto medio: 8.055
  (par moderado, no clase Lehmer).
- **El ojo (descenso recursivo):** u = 0.673376, profundidad −3.000,
  ancho 1.50 espaciamientos — y el ojo NO está sobre el par (a
  diferencia de la Tormenta I): misma clase por profundidad, otra
  anatomía. Curiosidad en pie: su compañera de grito 1.44879 × 10²¹
  también lleva el ojo en u ≈ 0.67.
- **Doble firma:** posiciones acordadas a 1.34/1.50 × 10⁻⁶ — la vara
  exacta que la doctrina F120 predice para duelos entre bloqueos
  distintos (placa colosal pre-cirugía vs clásico segmentado): la
  maquinaria consistente con su propia teoría de error.
- **Nombre del primer cero:** ordinal ~8,341,069,883,051,226,917,351.

Reproducir: `go run ./cmd/starship -arpon -anchor 1.14794e21` · ojo:
`-ojo -anchor 1.14794e21` (luz archivada de ambos motores, 0 ms)

## LA TORMENTA III — t = 1.2144079819075897 × 10¹⁹ (2026-08-05)

- **La primera ballena de la era de puntería exacta**: primer fondeo de la
  historia del laboratorio sobre el float exacto del vigía (grito #2,
  vaivén previsto 1.59, proa −2.50 sentando el extremo en el centro).
- **La esfera:** exige 5.00, halla 2 — **S = −3.00** (tercer evento clase
  ballena). Ojo interior −3.00 en u = 0.746. 3.0 minutos de colosal.
- **El par:** gap 0.070942641 = **0.4754 espaciamientos**; |Z| en el punto
  medio: 0.5447.
- **Doble firma — la más limpia jamás registrada en tormenta:**
  6.6×10⁻¹³ / 1.1×10⁻¹²; peor |dZ| 2.8×10⁻¹¹. Agua sin plegado → duelo
  de depósito puro → clase 10⁻¹¹, exactamente como predice F120: tres
  tormentas, tres firmas, cada una en su nivel de precisión previsto.
- **El dato pesado (F122):** marcapasos pronosticó +1.03; el mar dio
  −3.00 — residuo no-modelado −4.03 (~5σ). El grito halló violencia
  donde apuntó, pero la bestia pertenece a las voces no modeladas: el
  legajo de las colas pesadas suma su tercer evento.

Reproducir: `go run ./cmd/starship -arpon -anchor 1.2144079819075897e19`

---

*Laboratorio Diosyunalma — el registro completo del método vive en
[FINDINGS.md](../registro/FINDINGS.md); la guía de validación independiente en
[VALIDACION.md](VALIDACION.md).*
## LA CRESTA I — t = 1.5693364647413486 × 10¹⁹ (2026-08-05)

- **El primer tifón POSITIVO del catálogo**: la esfera exige 5.00 y el
  mar entrega 7 — S = +2.00; ojo interior +2.21 en u = 0.711. Grito #10
  exacto, proa +2.50. Todos los tifones previos eran déficit (valles);
  éste es la primera cresta.
- **EL RÉCORD MAYOR adentro:** par a 0.3141 espaciamientos que pasa a
  **|Z| = 0.029112 de un doble cero** — el acercamiento clase Lehmer más
  extremo jamás registrado aquí (9× más cerca que la Tormenta I). Donde
  los ceros sobran, se aprietan: la cresta comprime.
- **Doble firma:** 2.1×10⁻¹¹ / 1.4×10⁻¹¹, peor |dZ| 4.4×10⁻¹¹ — clase
  depósito puro, F120 cumpliéndose por cuarta tormenta consecutiva.
- Marcapasos: pronóstico +0.84, residuo +1.16 (bestia a medio oír).
  Distancia al agua muerta: 3.55 esp (asiento EN CONTRA de la ley del
  slack — el legajo se mantiene honesto).

Reproducir: `go run ./cmd/starship -arpon -anchor 1.5693364647413486e19`

---
## LA TORMENTA IV — t = 1.3743049032880831 × 10²⁰ (2026-08-06)

- **La segunda ballena de la era exacta — y EL INTERIOR MÁS HONDO JAMÁS
  SONDEADO**: frontera S = −3.00; el pozo desciende a **−4.17** en el ojo
  (más hondo que la Tormenta I: −4.01). Ancho 2.92 espaciamientos,
  k = 1.95 — familia de las profundas (T-I: 1.78), separando las dos
  clases anatómicas del legajo del resorte.
- **El agua ciega vecina decía −1.00**: la coordenada redondeada de la
  era ciega leyó marea mansa — la lección de F122 en carne viva.
- **El par:** 0.967296956 / 1.049004774 — gap 0.5791 espaciamientos,
  |Z| medio 0.538. Interior max|Z| = 1.561 (tirando a mudo).
- **Doble firma:** 1.35×10⁻¹¹ / 9.8×10⁻¹² (peor |dZ| 2.85×10⁻¹⁰) —
  clase depósito, F120 cumpliéndose por quinta bestia catalogada.
- **Filtro de coordenada maldita (F134): APROBADO** — lecturas
  físicamente plausibles, campo vivo; chequeo puntual del juez supremo
  encolado por prudencia nueva.

Reproducir: `go run ./cmd/starship -arpon -anchor 1.3743049032880831e20`

---
