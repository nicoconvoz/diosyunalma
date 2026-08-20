# Acta de la Cadena — respuesta a la Guía del Tambor (Auditoría 55)

**Respuesta a la Guía de orientación · 2026-08-20 · F382**

Sus 8 pasos, caminados en su orden, por una **ruta B independiente**: θ con
término asintótico extra, suma de Riemann–Siegel acumulada en orden
descendente, paso de grilla 0,005 (no 0,01), integración trapezoidal (no suma
simple), ventana corrida a [10000,5, 11000,5], clamp 1e-5 (no 1e-4). Nada
compartido con la ruta A salvo las definiciones — que es el experimento.

Se reproduce con `go run ./cmd/lacadena`.

---

## Pasos 1–2 · Las definiciones

F(t) = log|Z(t)| centrada en su media de ventana; c_m = (1/T)∫F(t)e^(−it·ln m)dt;
barra = |c_m|. Exactamente sus PIEZAS A, B, C.

## Paso 3 · La derivación, revisada explícitamente

Las dos escrituras de la amplitud — Λ(m)/(2√m·ln m) y p^(−k/2)/(2k) — se
evaluaron por separado para m = 2, 3, 4, 8, 9, 16, 19, 25, 27 (y 6, 10, 12
como ceros): **diferencia máxima 6,9×10⁻¹⁸** — precisión de máquina. Su §5
queda revisado y firmado.

## Paso 4 · La calibración que su pregunta 5 necesitaba

Antes de medir zeta, el instrumento se calibró con una señal **sintética** de
amplitudes conocidas — incluyendo, adrede, **un tono plantado en el compuesto
10** y nada en el 6:

| m | recuperado | esperado (amplitud/2) |
|---|---:|---:|
| 2 | 0,1503 | 0,1500 |
| 3 | 0,0751 | 0,0750 |
| 5 | 0,0349 | 0,0350 |
| 6 | 0,0002 | 0 |
| **10** | **0,0999** | **0,1000** |

**El aparato oye compuestos cuando suenan.** El factor ½ de la proyección
queda calibrado empíricamente (no solo derivado). Y la conclusión que sella
su control: el silencio de los compuestos en zeta es **de zeta**, no del
método. Ninguna explicación trivial por procedimiento sobrevive a esto.

## Pasos 5–6 · Los tres valores de su §15 y la fuga

| m | ruta B | ruta A (acta) | Λ predice |
|---|---:|---:|---:|
| 2 | 0,3529 | 0,3524 | 0,3536 |
| 19 | 0,1138 | 0,1146 | 0,1147 |
| 6 | 0,0035 | 0,0034 | 0 |

Fuga del m=6 con la ventana, ruta B: T = 500 → 0,0096; T = 1000 → 0,0035;
T = 2000 → 0,0007. La escalera decreciente de su §7, reproducida.

## Paso 7 · Los controles negativos, reproducidos

| conjunto | ruta B | ruta A | Λ |
|---|---:|---:|---:|
| primos reales | 0,1619 | 0,1620 | 0,1624 |
| compuestos puros | 0,0026 | 0,0028 | 0 |
| corridos p+1 | 0,0429 | 0,0429 | 0,0409 |

Dos rutas, mismos números a la milésima. Su §9 puede considerarse cerrado.

## Paso 8 · La frase de trabajo, respondida con el reparto exacto

- **MATEMÁTICA (demostrado):** la identidad algebraica del Paso 3; el factor
  ½ de la proyección (ahora calibrado); el producto de Euler y la expansión
  de von Mangoldt para Re(s) > 1; Λ(m) = 0 en compuestos ⟺ factorización
  única.
- **MEDICIÓN (numérico, dos rutas):** que la proyección de log|Z| sobre la
  línea crítica devuelve esas amplitudes a ~10⁻³; la fuga ~1/T; los controles
  negativos con correlación 1,000; la indiferencia total a la ventana.
- **NUEVO (si algo — y dicho sin inflar):** ningún teorema. Queda en pie:
  (a) el hecho empírico de que la estructura de von Mangoldt es **medible
  sobre la línea crítica** con esta precisión mediante esta proyección;
  (b) la prueba de calibración de que el silencio pertenece a zeta y no al
  aparato; (c) el instrumento mismo — el parche radial log n donde el timbre
  ES la aritmética, que conecta el tambor con la espiral (F375–F379) y su
  escalera armónica en un solo cuerpo.

Sus cinco preguntas del §16: sí, sí, sí (a precisión de máquina), sí (dos
rutas), sí (calibración + falsos primos). El semáforo rojo queda rojo: nada
de esto demuestra von Mangoldt sobre la línea ni toca a RH como teorema.

## Su §18 — la pregunta abierta, anotada como próxima

La información de los **ceros** no vive en los tonos de primos: los tonos
2..40 explican R² ≈ 0,5 del canto. La otra mitad — el residuo — son las púas
de log|Z| en cada cero: ahí vive la posición fina. Su pregunta «¿qué
información adicional contiene esta representación sobre la distribución de
los ceros?» tiene entonces una forma operativa: **estudiar el residuo de la
reconstrucción por primos** — qué queda del canto cuando los primos ya
cantaron lo suyo. Esa es la puerta que su guía dejó señalada, y queda a la
espera de su orden.

Estructura cerrada no es hipótesis demostrada.

**Todavía no.**
