# Acta de la Máscara — Fase XII

**Respuesta a Fase XII (Auditoría 45) · 2026-08-19**

**Atribución, primero y exacta (su §1 y §17):** la relación la descubrió **Nico,
a mano, haciendo cuentas con primos**: `(p+1)/2 = (q−1)/2 ⟺ q = p+2`. Doc hizo
la lámina. Esta acta es la validación computacional que faltaba, y se registra
como trabajo posterior al hallazgo.

Su cierre ordenaba: *no intentes demostrar que la intuición tenía razón —
encontrá exactamente dónde falla.* Se obedeció, y el resultado tiene dos mitades
de valor desigual: **la geometría es exacta sin excepción** (la intuición no
falla en la matemática), y **la máscara rala da nulo limpio a este tamaño** (la
intuición no opera todavía en el espectro — y el porqué quedó medido).

Se reproduce con `go run ./cmd/lamascara`.
Lámina: `galeria/laminas/10-el-telar/la-mascara.svg`.

---

## 1 · La validación algebraica (su §14.1–2): EXACTA

Formalización del centro compartido: cada primo impar p tiene dos **anclas**
enteras, `a±(p) = (p±1)/2` — los dos enteros que flanquean su mitad. Con eso, el
diccionario geométrico entero de la lámina se vuelve tres identidades:

| gap g | regla geométrica | identidad |
|---|---|---|
| 2 | mismo centro | a⁺(p) = a⁻(q) — la igualdad de Nico |
| 4 | se tocan sin superponerse | a⁻(q) − a⁺(p) = 1 (anclas adyacentes) |
| > 4 | hueco de (g−4)/2 | exactamente (g−4)/2 enteros entre las anclas |

Verificado sobre **los 9590 pares de primos impares consecutivos hasta 100 000:
1224/1224 gemelos comparten ancla · 1215 pares g=4 adyacentes · 7151 huecos
exactos · CERO fallas.**

**La igualdad del capitán es un teorema chico y verdadero.** La geometría de la
lámina es identidad aritmética, no coincidencia visual. Eso queda demostrado en
el sentido de su §16, primer eslabón.

## 2 · La máscara, congelada antes de todo espectro (su §7 y §12)

- **Objeto "relación entre primos" en el operador (su §14.3):** el enlace (i,j)
  entre dos modos; cada modo pertenece a un primo.
- **Regla congelada:** enlace marcado ⟺ |i−j| ≤ 120 y |p_i − p_j| = g. Marcado =
  signo −1, amplitud intacta (signo² = 1 conserva F = 30 exacta).
- **Clases declaradas, todas corridas y publicadas:** g = 2, 4, 6.
- **Controles:** azar a igual cantidad de enlaces (6 semillas), y permutada a
  igual cantidad **e igual distribución por distancia** (6 semillas).
- Rieles de Fase X vigentes: vivos ≥ 149/187, PR/N ≥ 0,090.
- Predicción pre-registrada: máscara ≈ controles, tras cinco hojas iguales.

## 3 · Lo medido

| brazo | enlaces | vivos | Σ²(10) | PR/N | riel |
|---|---:|---:|---:|---:|---|
| S0 · sin máscara | 0 | 187 | 18,335 | 0,105 | ok |
| **M gemelos g=2** | **62** | 181 | **15,497** | 0,091 | ok |
| A · azar emparejado | 62 | 186 | 16,444 ± 0,271 | 0,105 | ok |
| C · permutada por k | 62 | 182 | 15,871 ± 1,548 | 0,085 | atrapada |
| **M g=4** | 66 | 185 | 15,491 | 0,111 | ok |
| A / C | 66 | 185–186 | 15,515 ± 0,936 / 16,410 ± 1,156 | 0,10–0,11 | ok |
| **M g=6** | 114 | 185 | 15,738 | 0,131 | ok |
| A / C | 114 | 184–186 | 15,767 ± 0,496 / 15,973 ± 1,572 | 0,12 | ok |

Densidades: **0,15 %–0,28 %** de los 40 740 enlaces. Ralas de verdad.

## 4 · El veredicto (sus §12–13)

| clase | contra azar | contra permutada | veredicto |
|---|---:|---:|---|
| g=2 | **+3,50 σ** | **+0,24 σ** | la distribución por distancia explica la ventaja |
| g=4 | +0,03 σ | +0,80 σ | indistinguible |
| g=6 | +0,06 σ | +0,15 σ | indistinguible |

**Nulo limpio, con un matiz que usted anticipó textualmente.** La máscara de
gemelos le gana al azar crudo por 3,5 σ — pero **no** a su permutada por
distancia (0,24 σ). Es su caso «B ≈ C y A distinto: la geometría de distancia
explica la mejora». Lo que la máscara gemela aporta sobre el azar es *en qué
distancias k caen los gemelos* (k ∈ [3, 80], porque los gemelos tienen log p
casi iguales y sus armónicos quedan cerca en el espectro ordenado) — y esa
distribución, copiada sin aritmética, rinde lo mismo. **La correspondencia
exacta de qué pares son gemelos no hace falta.** Por su §13: no hay señal
aritmética detectable por este canal **a este tamaño**.

El eco no se corrió: la regla lo permitía sólo con señal sobreviviente.

## 5 · La nota de escala, que es la conclusión operativa

Con N = 400 modos el medio contiene poquísimos pares gemelos: la máscara marca
62 enlaces de 40 740. **La máscara existe pero casi no toca al operador** — los
tres brazos apenas se despegan de S0 (−2,6 a −2,8 en Σ², compatible con
cualquier puñado de signos). La pregunta de su §15 no queda respondida en
negativo: queda respondida **«no decidible a N = 400»**. Decidir si la geometría
de Nico opera pide un medio más grande (más pares gemelos adentro), no otra
regla — y agrandar N es una perilla declarable de antemano, no un ajuste.

## 6 · Veredictos separados

| pregunta | veredicto |
|---|---|
| ¿el diccionario geométrico es exacto? | **SÍ — teorema, 9590/9590** |
| ¿la máscara gemela mueve el espectro? | apenas; compatible con su densidad |
| ¿hay señal aritmética sobre controles emparejados? | **NO a este tamaño** |
| ¿la correspondencia exacta par↔gemelo importa? | NO: la permutada por k rinde igual |
| ¿el canal queda cerrado? | **NO — queda «no decidible», pendiente de escala** |

---

El hallazgo es de Jesús Nicolás Astorga, a mano. La formalización por anclas, la
verificación exhaustiva y los controles emparejados, de este taller. La lupa se
usó como usted mandó: para buscar dónde falla. En la matemática no falla; en el
operador, a este tamaño, todavía no le llega la voz.

## Anexo · La escalada bendecida (agregado el mismo día)

Usted preguntó si M_gemelos − M_permutada se separa de cero al crecer N. Se
corrió con todo congelado (kmax, F, máscara, control por k, 6 semillas) y el
capitán ordenó cortar en el tercer escalón, con la serie ya elocuente:

| N | enlaces g=2 | vivos | Σ² S0 | Σ² gemela | Σ² permutada | Δ |
|---:|---:|---:|---:|---:|---:|---:|
| 400 | 62 | 181/187 | 18,335 | 15,497 | 15,698 ± 1,892 | +0,11 σ |
| 800 | 129 | 591/599 | 9,972 | 9,974 | 9,409 ± 0,888 | −0,64 σ |
| 1600 | 269 | 1437/1438 | 6,488 | 5,719 | 5,770 ± 0,567 | +0,09 σ |

**Δ oscila alrededor de cero: no hay separación sistemática.** La regla
pre-registrada (crecimiento monótono y último punto > 2 σ) quedó rota en el
segundo escalón. El escalón N = 3200 se cortó por orden del capitán tras ~6,5
horas de horno, con el veredicto ya determinado por la no-monotonía.

**El canal de la correspondencia exacta par↔gemelo queda CERRADO con datos**, a
tres escalas y con la banda cada vez más sana (a N = 1600 sobreviven 1437 de
1438 niveles). Lo que la máscara gemela aporta es su distribución por
distancia, en toda escala medida.

De regalo, una calibración que vale para todo lo que sigue: el propio S0 afina
al crecer el mar — 18,3 → 10,0 → 6,5 — y la banda deja de perder niveles. El
tamaño era la perilla más barata que teníamos sin usar.

---

Estructura cerrada no es hipótesis demostrada.

**Todavía no.**
