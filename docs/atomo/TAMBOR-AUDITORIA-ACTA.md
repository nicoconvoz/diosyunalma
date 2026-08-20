# Acta del Tambor — respuesta a la Auditoría 54

**Respuesta a la auditoría preliminar de «El Tambor» · 2026-08-20 · F380–F381**

Su auditoría llegó sobre la lámina sola, sin el código ni las definiciones — y
aún así apuntó exactamente a los dos huecos que había que cerrar: la prueba de
las tres señales (su §10) y el control negativo (su §11). Se corrieron los dos,
más la robustez de ventana (su §8). **La afirmación central sobrevivió a los
controles diseñados para matarla, y la fórmula de amplitud no es ajuste: se
deriva.**

Se reproduce con `go run ./cmd/eltambor` (la medición original) y
`go run ./cmd/elparche` (los controles de esta acta).

---

## 1 · Sus 14 preguntas, una por una

1. **¿Qué es τ?** τ = t₀/(2π), la escala de escalera de F376–F379 (el mismo τ
   de la ley Q = x/(eˣ−1)). Con t₀ = 10000: τ = 1591,549.
2. **¿De dónde sale 7,372?** ln τ = ln(1591,549) = 7,3725. El 7,3827 medido en
   los pares espejo difiere en 0,0102 porque los compañeros se redondean a
   enteros (τ/2 = 795,77 → 796) y los tonos se miden por cruces por cero sobre
   ventana finita. Diferencia explicada, no escondida.
3. **¿La función completa?** F(t) = log|ζ(½+it)| = log|Z(t)|, con Z de
   Riemann–Siegel. Las barras: c_m = (1/T)∫ F(t)·e^(−it·ln m) dt sobre la
   ventana; la altura es |c_m|.
4. **¿Qué es t?** La altura continua sobre la línea crítica — el mismo dial de
   toda la campaña. Muestreada a paso 0,01.
5. **¿Cómo se calcula cada barra?** Proyección de F sobre el tono ln m
   (pregunta 3). Sin intervención manual: `for m := 2; m <= 40; m++`.
6. **¿p^(−k/2)/(2k) se deriva o fue ajuste?** SE DERIVA. Del producto de
   Euler: log ζ(s) = Σ_m Λ(m)/(ln m)·m^(−s) (expansión de von Mangoldt). El
   término de m tiene amplitud Λ(m)/(√m·ln m); la proyección sobre coseno
   la parte a la mitad: Λ(m)/(2√m·ln m) = p^(−k/2)/(2k) si m = p^k, 0 si no.
   Cero parámetros libres, cero ajuste.
7. **¿Identidad conocida?** Sí: el producto de Euler + von Mangoldt. Es
   matemática de libro — el valor del experimento es que el tambor la MIDE.
8. **¿Qué significa «suena»?** |c_m| compatible con Λ(m)/(2√m·ln m) dentro del
   error de ventana. Medido: primo 2 → 0,3524 contra 0,3536; primo 19 →
   0,1146 contra 0,1147.
9. **¿Qué significa «callan»?** Coeficiente EXACTAMENTE CERO en la expansión
   (Λ(m) = 0 para m compuesto no potencia — teorema, factorización única). En
   la medición: residuo de fuga de ventana finita que DECRECE al alargar T
   (m=6: 0,0094 con T=500 → 0,0034 con T=1000 → 0,0001 con T=2000).
10. **¿Sobrevive al cambiar la ventana?** Sí: [10k,11k], [30k,31k],
    [100k,101k] dan |c_2| = 0,3524 / 0,3530 / 0,3530 (predicho 0,3536). La
    ventana original fue elegida a priori y nada depende de ella.
11. **¿Sobrevive a otra normalización?** Sí: T = 500, 1000, 2000 — los primos
    quietos, los compuestos encogiéndose (como manda la fuga).
12. **¿Control contra conjuntos aleatorios?** Corrido (su §11, abajo): la
    respuesta es mejor que un sí.
13. **¿Todo expresable como suma/integral?** Sí — pregunta 3 más la expansión
    de la pregunta 6. Estructura exacta que usted pidió en su §9:
    F(t) ≈ Σ_p Σ_k A(p,k)·cos(t·k·ln p + φ), A(p,k) = p^(−k/2)/k.
14. **¿Relación real con zeta o analogía visual?** Real y literal: F ES zeta.
    El tambor no se parece a la guitarra: es una proyección medida de ella.

## 2 · Su Prueba 10 — las tres señales (el resultado más limpio)

Reconstrucción de F(t) por proyección sobre cada conjunto de tonos, ventana
[10k, 11k]:

| señal | tonos | R² |
|---|---:|---:|
| todos los enteros 2..40 | 39 | 0,4990 |
| primos y potencias | 19 | **0,4993** |
| solo primos | 12 | 0,4632 |

**Duplicar los tonos agregando los 20 compuestos aporta −0,0003 de R²** (ruido
de sobreajuste: los compuestos no traen información). Las potencias sí aportan
(0,4632 → 0,4993). El R² total ~0,5 es lo esperable: faltan los primos > 40 y
la estructura local de los ceros; la comparación ENTRE conjuntos es la prueba.

## 3 · Su Prueba 11 — el control negativo

Mismo procedimiento, conjuntos falsos de igual cardinalidad (12):

| conjunto | medido | Λ predice |
|---|---:|---:|
| primos reales | 0,1620 | 0,1624 |
| compuestos puros | 0,0028 | 0,0000 |
| primos corridos p+1 | 0,0429 | 0,0409 |
| 20 conjuntos al azar | 0,0654 ± 0,0267 | 0,0639 |

Correlación medido↔predicho a través de los 20 conjuntos azarosos: **1,000**.
El detalle que sella el control: los primos corridos p+1 NO callan del todo —
suenan 0,0429 — y Λ lo predice (0,0409): el conjunto {3,4,6,8,...} contiene
potencias de primo por accidente, y suena exactamente por ellas. **Ningún
conjunto suena por densidad ni por procedimiento: cada conjunto suena lo que Λ
dice de sus miembros, miembro por miembro.**

## 4 · Los semáforos, actualizados tras los controles

- Su VERDE queda verde.
- Sus AMARILLOS: la amplitud p^(−k/2)/(2k) pasa a **derivada** (§1.6); la
  ventana pasa a **irrelevante** (§1.10-11); el espejo n ↔ τ/n queda como lo
  que F377-F379 midieron: la inversión de la escalera armónica del reloj.
- Sus ROJOS, con la formulación corregida que usted pidió en su §7: no decimos
  «los compuestos no tienen estructura prima» — decimos que **el logaritmo
  convierte el producto de Euler en suma sobre potencias de primo únicamente**:
  el tono de un compuesto se descompone en los tonos de sus factores y jamás
  aparece como frecuencia propia. Eso ES su formulación compatible con la
  factorización única, y los controles la midieron.
- ROJO que queda rojo, con todas las letras: la expansión de von Mangoldt está
  demostrada para Re(s) > 1; **sobre la línea crítica es aquí una afirmación
  empírica** (medida a 10⁻³ y estable), no un teorema. Y nada de esto es una
  afirmación sobre la Hipótesis de Riemann.

## 5 · Lo nuevo que queda en pie

El tambor no descubrió la expansión de von Mangoldt — la HIZO SONAR: el parche
radial log n del capitán es la representación donde esa identidad se vuelve
timbre (BOOM/TAK), el espejo n ↔ τ/n se vuelve el par de tonos que suma ln τ,
y «quién golpea» tiene respuesta medida: los primos, cada uno desde su radio,
con la fuerza exacta que la aritmética manda.

Estructura cerrada no es hipótesis demostrada.

**Todavía no.**
