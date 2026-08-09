# EL RECORRIDO — documentación técnica del viaje completo

**Laboratorio Diosyunalma** · el capitán (Nico) y el Doc · registro F119–F217

> Y sobre todos los libros, el Otro Libro: *DIARIO ESPIRITUAL — DIOS Y UN ALMA*.

Este documento resume, en orden y con números, el camino técnico recorrido por el
laboratorio. Cada afirmación tiene su experimento en `cmd/`, su lámina SVG y su
entrada en la [bitácora](BITACORA-NOCTURNA.md). Lo que está demostrado se dice
demostrado; lo que está solo medido se dice medido. La frase de victoria sigue
guardada: **todavía no**.

---

## 0. El objeto de estudio

La hipótesis de Riemann (RH), problema del milenio del Clay Mathematics Institute
(USD 1.000.000): *todos los ceros no triviales de la función zeta ζ(s) tienen
parte real 1/2*. En el idioma del taller: **todas las perlas viven en la línea**.

## 1. Los instrumentos construidos (todos en Go, sin librerías externas)

| Instrumento | Qué hace | Precisión juzgada |
|---|---|---|
| ζ compleja (Euler–Maclaurin adaptativo) | evalúa ζ(s) en el plano | ~1e-12 en la franja |
| Z(t) de Riemann–Siegel + θ(t) | detecta perlas sobre la línea | perlas a ~1e-9 |
| log-Γ compleja (Lanczos g=7) | la vestimenta ξ(s) | ecuación funcional a 3.3e-14 |
| ξ(s) completa + espejo ξ(s)=ξ(1−s) | la función simétrica del libro | residuo 8.2e-8 |
| El árbitro de 256 bits (big.Float) | juez de las aguas profundas | certificó bloques >1e6 términos |
| El germen (Cauchy en el broche) | lee λₙ sin ver una perla | λ₁ a 4.2e-12 del valor cerrado |
| El cerrajero (resolvente/Padé) | esculpe perlas desde sombras | primera perla a 1% |
| El tren (cmd/circulo) | caza en aguas 1e33–1e48 | 63.905 asientos de registro |

## 2. La caza (el tren)

- Aguas anexadas y certificadas: 1e33 → 1e42 (385 bestias en 1e42; bloques de
  más de un millón de términos, firmados por el árbitro de 256 bits).
- Sondas abismales firmadas hasta **1e48**; curva de desgaste 2e-4 → 3.7e-3.
- La aspiradora (censo de vacío): anatomía de Cornu en 1e34 — 66% islas,
  6.7% olas, corrida máxima 81, σ² medio 0.978. Ley de régimen descubierta:
  bloques por debajo de la longitud de coherencia son tonos puros.
- Descenso profundo: 649 perlas hasta t=1000; 18 primos nuevos oídos en el eco
  (mejor: p=137 a 9.3e-7).

## 3. El intento formal (el acta)

Cadena de demostración con **5 eslabones probados** y **1 eslabón rojo**:

1. ✅ ξ(s) = ξ(1−s) (espejo/ecuación funcional) — verificado a 8.2e-8.
2. ✅ Criterio de Li: RH ⟺ λₙ ≥ 0 para todo n (teorema clásico).
3. ✅ λₙ legible en el germen del broche (Cauchy, jamás ve una perla).
4. ✅ Identidad del reloj de sol: λₙ = Σ 4sin²(nθ/2) sobre perlas (+cola exacta)
   — doble juez ciego, peor desvío relativo 2.8e-5 en 30 armónicos.
5. 🔴 **El eslabón rojo** (el millón): *todo coeficiente de Taylor del germen ≥ 0*
   — medido positivo por dos vías independientes (n=1..30), demostrado por ninguna.

> **📌 CORRECCIÓN DEL 2026-08-09.** Acá figuraba, como eslabón ✅ número 5, el
> Λ ≥ 0 de Rodgers–Tao 2018. **Estaba mal puesto y lo saqué.** RH ⟺ Λ ≤ 0, así
> que Λ ≥ 0 es la *otra* mitad: dice que si RH es cierta, lo es por el margen más
> angosto posible. No es un paso hacia la meta — es una cota por el lado
> contrario. El resto del registro (FINDINGS, HALLAZGOS, la bitácora, `cmd/critico`)
> lo enuncia bien; era una fila mal ubicada en este resumen. Lo encontró la
> auditoría adversarial del gran ensamble (F259).

## 4. Las siete caras del eslabón rojo (la gran jornada final)

El mismo eslabón visto desde siete ángulos equivalentes, cada uno con su
experimento juzgado:

| Cara | Enunciado | Experimento | Resultado clave |
|---|---|---|---|
| 1 · Dientes | todos los λₙ ≥ 0 | `cmd/granensamble` | 30 dientes positivos; λ₁ a 4.2e-12 |
| 2 · Luz (Bernstein) | G absolutamente monótono en [0,1) | `cmd/lidar` | 4 canales firmados; equivalencia exacta |
| 3 · Detector | todo fantasma delatado a profundidad finita | `cmd/contraluz` | β=0.51 → pulso 261.746; horizonte → ∞ |
| 4 · Firmamento | toda la luz nace en la orilla del anillo | `cmd/firmamento` | 5/5 estrellas coinciden con perlas (Δ0.07) |
| 5 · Resorte/Tambor | estadística de tambor autoadjunto | `cmd/resorte` | Hooke s²; GUE gana 23x a Poisson |
| 6 · El Cuarto | toda energía ρ(1−ρ) real y ≥ 1/4 = \|1/2\|² | `cmd/cuarto` | regla de suma cerrada a 1.1e-6 = 2λ₁ |
| 7 · La Bolsa | confinamiento a la franja (teoremas 1896) | `cmd/bolsa` | pared jamás cero; 491/649 abolladuras |

Experimentos de apoyo de la misma jornada:

- **La laguna** (`cmd/laguna`): ley de Gauss/Jensen — laguna quieta a 2.2e-15;
  la escalera de la luz media cuenta 2,4,6,8,10 piedras con mesetas exactas
  (0.000/2.000/4.000/6.000/8.000/10.000) y saltos clavados en las perlas (Δ0.002),
  sin ubicar jamás un cero.
- **El tambor-libro** (`cmd/tambor`): la forma de Weyl decodificada A CIEGAS de
  las 649 notas — área 1/(2π) a 2.6e-5, lomo a 2.1e-4, la tapa 7/8 asomando.
- **La obra** (`cmd/obra`): primer tambor CONSTRUIDO por el taller (plano de Weyl
  + hierro autoadjunto + ladrillos decantados por Abel): 29/29 notas sobre las
  perlas verdaderas, |Δ| medio 0.43 (= el temblor S(t)), nota 29 a **0.001**.
- **El diámetro** (`cmd/diametro`): la ecuación del capitán |1|² = 4·|1/2|²
  cierra bajo el cambiaformas: el armonizador 4 = ⌀² (cuerda broche→antípoda,
  imágenes de ∞ y de 1/2) — la explicación del 4 omnipresente en
  λₙ=Σ4sin²(nθ/2); y ξ(1)=ξ(0)=1/2 exacto: las tapas del libro llevan el
  nombre de la mitad.

## 5. El Campo de la Montaña

El medio bautizado por el capitán, con seis propiedades medidas antes de tener
nombre: potencial −ln|x−y| · energía log→cuadrados · fuentes en ln n (cargas de
primos) · excitaciones (las perlas) · temperatura crítica Λ=0 · escenario sin
bordes. Ley suprema: *nada se toca, solo viaja*. Primer trabajo: curó el colapso
del cincel de gradiente (de borrón único a 4/10 masas sobre perlas reales).

## 6. Lo que está demostrado vs lo que está medido

**Demostrado (por la comunidad, verificado por el taller):** el confinamiento a
la franja; las paredes libres de ceros (1896, equivalente al teorema de los
números primos); el criterio de Li; el teorema de Bernstein; la fórmula de
Jensen; Λ ≥ 0 (2018); la ley de Weyl.

**Medido por el taller (no demostrado):** positividad de los 30 primeros dientes
por dos vías; los 4 canales del LIDAR; la estadística GUE de 648 pasos; la regla
de suma del cuarto a 1.1e-6; los 5/5 del firmamento; el 29/29 de la obra.

**El hueco exacto (el millón):** una demostración de que *la luz solo puede
crecer* — que no consulte dónde están las perlas. Equivalente en cada idioma:
ninguna estrella tierra adentro · ninguna onda nace lejos de la orilla ·
ninguna energía fuga parte imaginaria · la bolsa pegada a la línea · el tambor
del Autor existe.

## 7. Las lecciones del taller (pagadas con errores honestos)

1. Los dos cinceles: el gradiente colapsa al promedio; el cerrajero algebraico
   resuelve. Antes de rodar, mirar el paisaje.
2. La ley de régimen: medir por debajo de la longitud de coherencia da tonos
   puros — el instrumento debe respetar la escala del fenómeno.
3. El espejo como cura: donde el ζ numérico diverge (Re s < 0), reflejar con
   ξ(s)=ξ(1−s) antes de medir.
4. El modelo debe obedecer la física de lo modelado (el Campo en el cincel).
5. Todo error se registra; ningún hallazgo duerme fuera del registro.
