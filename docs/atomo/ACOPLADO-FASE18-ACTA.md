# Acta del Acoplado — Fase XVIII

**Respuesta a Fase XVIII (Auditoría 51) · 2026-08-20 · F370**

Su hipótesis: repulsión y campo deben resolverse **juntos**. Se construyó sin
perillas, se pre-registró la predicción, y el veredicto por sus criterios del
§16–17 es **APOYO PARCIAL CUANTITATIVO**: la auto-consistencia mejora a los dos
modelos de referencia a la vez, la comparación decisiva (simultáneo contra a
posteriori) da **2,5×**, el control asesino funciona — y el déficit restante
quedó con nombre nuevo, más fino que el anterior.

Se reproduce con `go run ./cmd/laformulamadre fase18` (semilla 20260825).
Lámina: `galeria/laminas/10-el-telar/el-acoplado.svg`.

---

## 1 · La construcción, sin perillas (sus §5–7)

Los ceros reales están clavados por la ecuación de conteo completa. El modelo
hace **exactamente eso** con material GUE:

```
u_k = u_{k−1} + w_k          w_k ~ Wigner (GUE, media 1)   [la repulsión]
γ_k resuelve:  N_liso(γ_k) + S(γ_k) = u_k                  [el clavado]
```

con S la misma fluctuación truncada (n ≤ 97) de las Fases XVI–XVII. El campo
entra **adentro** del clavado: donde S sube, las posiciones se corren Y los
gaps locales se comprimen a la vez, consistentemente — tamaños y posiciones
emergen juntos, que es su hipótesis hecha ecuación. Cero parámetros: la fórmula
fija todo. *(La otra familia de su §5 — el gas de Coulomb — necesita temperatura
y cronograma de relajación: perillas. Por su §6 se corrió la implementación sin
perillas y el gas queda diferido, declarado.)*

## 2 · La tabla (su §9)

| modelo | amplitud | cruce s\* | pendiente | s desv (Wigner: 0,42) |
|---|---:|---|---|---:|
| **REAL** | **0,349** | 0,865 | −0,95 | **0,42** |
| R) repulsión sola (S=0) | 0,024 | dispersos | ≈0 | 0,43 |
| F) campo solo (sin repulsión) | 0,285 | 0,94 | −0,80 | — |
| **AC) ACOPLADO coherente** | **0,138** | **0,88–0,96** | −0,10…−0,37 | 0,58 |
| AF) acoplado, fase destruida | **0,037** | dispersos | ≈0 | 0,58 |
| AP) a posteriori (XVII) | 0,055 | 0,62–0,75 | −0,2…−0,4 | 0,60 |
| AD) densidad sola (XVI) | 0,115 | ~0,70 | −0,3…−0,5 | **3,67** ⚠ |

## 3 · El veredicto contra sus criterios del §16, punto por punto

- **Sin circularidad, sin parámetros ajustados** ✓ — la regla es la ecuación de
  conteo, la amplitud la fija la fórmula.
- **Mejora cuantitativa sobre AMBAS referencias** ✓ — 0,138 contra 0,115
  (densidad) y 0,055 (a posteriori). **La comparación que usted marcó como
  decisiva (su §10): simultáneo = 2,5 × a posteriori.** La auto-consistencia no
  es irrelevante: es el factor dominante entre los dos.
- **El cruce se mueve a la zona real** ✓ — de 0,70 (XVI–XVII) a 0,88–0,96,
  encerrando al 0,865 real (con leve sobre-tiro).
- **T4 con el patrón correcto y magnitud que crece** ✓ — 0,116 en el primer
  bin: el doble del a posteriori (0,055), 43 % del real (0,267). El factor
  faltante bajó de ~3 a ~2,3 **sin tocar una perilla**.
- **La coherencia es necesaria** ✓ — fase destruida: 0,037, el piso GUE. Muere
  como debe.
- **Residuo** — máx 0,250 contra 0,312 (densidad) y 0,300 (a posteriori):
  mejora real pero modesta; la amplitud global sigue un factor ~2,5 abajo.
- **Estable entre semillas** ✓ — cruces 0,88/0,90/0,96, amplitudes consistentes.

**Ninguno de sus criterios de fracaso del §17 se cumplió.** En particular el
central: la diferencia simultáneo/a-posteriori NO resultó irrelevante — resultó
el efecto más grande de la fase.

## 4 · Lo que su §12 mandó registrar en vez de ocultar — y es la pista nueva

Dos hechos de estadística de espaciados, reportados sin maquillaje:

1. **El acoplado paga un precio que los ceros reales no pagan:** su desvío de
   espaciados es 0,58 contra 0,42 de Wigner — el campo ensancha la distribución.
   **Los ceros reales llevan la señal completa (0,349) manteniendo el 0,42
   exacto.** El mecanismo real es más eficiente: más señal, cero distorsión.
2. De paso quedó expuesto que el viejo modelo de densidad (XVI) **siempre tuvo
   la estadística rota** (desvío 3,67 por la zona baja del espectro) — su §12
   aplicado retroactivamente. El acoplado la arregla casi del todo.

**Y de ahí sale el déficit con nombre nuevo.** El brazo «campo solo» da 0,285 —
más que el acoplado. O sea: **la repulsión iid (paseo de gaps Wigner
independientes) DILUYE la coherencia del campo**. Pero la repulsión real de los
ceros no es un paseo iid: el GUE verdadero tiene **rigidez de largo alcance**
(la varianza del conteo crece como log, no linealmente). Un paseo iid difunde y
emborrona el clavado; una secuencia rígida lo preservaría. **El ingrediente que
falta ya no es el acoplamiento (probado): es la RIGIDEZ de largo alcance del
material GUE** — nuestros u_k son ruido browniano donde deberían ser casi una
regla. Hipótesis para la fase siguiente, sin perillas: reemplazar el paseo iid
por una secuencia con rigidez GUE genuina (por ejemplo, autovalores CUE/GUE
reales en vez de gaps Wigner independientes) y repetir el juicio idéntico.

## 5 · La escalera de las tres fases, en una línea

XVI: mover tamaños da ⅓ · XVII: mover posiciones después, menos aún · XVIII:
**moverlos juntos da el doble que separados y encierra el cruce** — y lo que
queda apunta a la rigidez, no a la fuerza. Su cierre era la dirección correcta:
no aumentar la fuerza — cambiar la dinámica. Se cambió, y respondió.

HL y el Telar siguen congelados, como manda.

---

Estructura cerrada no es hipótesis demostrada.

**Todavía no.**
