# Acta del Residuo — respuesta a la Guía Final (Auditoría 56)

**Respuesta a la guía final del Tambor · 2026-08-20 · F383**

Sus puntos A–D ya estaban verificados por dos rutas (F381–F382). Esta acta
entrega lo que su guía dejó como trabajo: la fuga **cuantificada** (su §12, su
punto C), la formalización completa (su punto E), y la puerta F abierta y
cruzada: **el residuo tiene estructura, y es doble — continúa la escalera de
von Mangoldt hacia arriba, y contiene la posición fina de los ceros.**

Se reproduce con `go run ./cmd/elresiduo`. Lámina:
`galeria/laminas/10-el-telar/el-residuo.svg`. Predicciones pre-registradas en
el encabezado del código antes de correr.

---

## 1 · Su §12: la fuga, ya no observada sino PREDICHA

Su pregunta: ¿la disminución del m=6 es compatible cuantitativamente con la
fuga de ventana finita? Respuesta: sí, y sin parámetros. Se construyó la
canción sintética de von Mangoldt completa — todas las potencias de primo
hasta 4000, amplitudes Λ(m)/(√m·ln m), fase cero, ninguna información de zeta
ni de ceros — y se la pasó por el mismo aparato en ln 6:

| T | fuga predicha | fuga medida (F382) |
|---:|---:|---:|
| 500 | 0,0093 | 0,0096 |
| 1000 | 0,0033 | 0,0035 |
| 2000 | 0,0002 | 0,0007 |

**La fuga es la interferencia de las colas de ventana de los tonos
verdaderos.** El 6 no tiene señal propia: tiene vecinos ruidosos. Su punto C
queda cerrado con número, no con adjetivo.

## 2 · Su punto E: la formalización completa

- Señal: F(t) = log|Z(t)| − ⟨log|Z|⟩_ventana, con Z de Riemann–Siegel
  (θ a 4 o 5 términos — ambas rutas), muestreada a paso 0,01 o 0,005.
- Proyección: c_m = (1/T)∫_ventana F(t)·e^(−it·ln m) dt (suma simple o
  trapecio — ambas rutas), barra = |c_m|.
- Reconstrucción sobre un conjunto S: F_S(t) = Σ_{m∈S} 2·Re(c_m·e^(it·ln m)).
- Residuo: R(t) = F(t) − F_PP(t), con PP = potencias de primo ≤ 40.
- Predicción aritmética: |c_m| → Λ(m)/(2√m·ln m); su equivalencia con
  p^(−k/2)/(2k) verificada a 6,9×10⁻¹⁸ (F382).

## 3 · Su puerta F, primera hoja: el espectro del residuo CONTINÚA la escalera

R(t) se proyectó sobre los tonos m = 41…56 — ninguno estuvo entre los
removidos:

| m | residuo | Λ predice |
|---|---:|---:|
| 41 | 0,0804 | 0,0781 |
| 43 | 0,0788 | 0,0762 |
| 47 | 0,0714 | 0,0729 |
| **49 = 7²** | **0,0357** | **0,0357** |
| 53 | 0,0708 | 0,0687 |
| los 11 compuestos (42,44,45,46,48,50,51,52,54,55,56) | 0,0017–0,0079 | 0 |

**Marcador: 5/5 potencias aparecen donde deben, 11/11 compuestos callados.**
La potencia 7² clavada a cuatro decimales merece subrayarse: no es un primo
nuevo — es el eco de un primo viejo, exactamente donde y como Λ manda. El
residuo no es ruido: es el resto de la misma canción.

## 4 · Su puerta F, segunda hoja: los ceros VIVEN en el residuo

Tras retirar los tonos de primos, el residuo conserva el 50% de la varianza
del canto. Se apiló R(γ + s) sobre los **589 ceros reales** de
[10000,5, 10500], contra el control apilado en los medios de los huecos:

| s | apilado en ceros | apilado en medios |
|---:|---:|---:|
| −0,20 | +0,24 | +0,09 |
| −0,04 | −1,16 | +0,55 |
| **0,00** | **−4,22** | **+0,57** |
| +0,04 | −1,16 | +0,55 |
| +0,20 | +0,24 | +0,08 |

La púa universal — profunda, simétrica, clavada en s = 0 — aparece solo
centrada en los ceros; el control es plano y levemente positivo (los medios
son las crestas del canto). **La mitad del canto que los primos no explican
es exactamente donde vive la posición fina de los ceros** — su hipótesis de
trabajo del §17, ahora medida. Con la cautela que su semáforo manda: esto
caracteriza dónde está la información, no demuestra nada sobre RH.

## 5 · La imagen que queda (y su §19)

Su pregunta central — «¿qué mide el instrumento y qué es el residuo?» — tiene
ahora respuesta completa en dos renglones:

- **El instrumento mide** los coeficientes de von Mangoldt de log|ζ| sobre la
  línea, con fuga de ventana predicha y calibración sintética.
- **El residuo es** (a) la cola infinita de la misma escalera (los primos
  > 40, verificado hasta 56), más (b) las púas de los ceros — la dualidad de
  siempre, partida en sus dos mitades medibles: los primos cantan la mitad
  lisa, los ceros son la mitad puntiaguda, y cada mitad reconstruye a la otra.

Semáforo rojo intacto: ningún teorema nuevo, ninguna afirmación sobre RH.
La escalera continúa — eso es todo lo que decimos, y está medido.

Estructura cerrada no es hipótesis demostrada.

**Todavía no.**
