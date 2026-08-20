# Acta de la Truncación — Fase XX

**Respuesta a Fase XX (Auditoría 52) · 2026-08-20 · F372**

Su orden fue estricta y se cumplió estricta: **no se agregó nada.** El campo
puro, la única variable el tope N de la suma. Resultado: **las tres
convergencias de su §7 ocurrieron — y por su §16 se frenó, se auditó el primer
elemento, y ese elemento reencuadra la campaña entera.** No se celebra: se
reporta.

Se reproduce con `go run ./cmd/laformulamadre fase20` (y `audit20`).
Lámina: `galeria/laminas/10-el-telar/la-truncacion.svg`.

---

## 1 · La tabla principal (su §6)

| N | amplitud | cruce s\* | pendiente | s desv | T4 bin1 | residuo máx |
|---:|---:|---:|---:|---:|---:|---:|
| 97 | 0,267 | 0,847 | −0,615 | 0,42 | 0,215 | 0,0960 |
| 997 | 0,315 | 0,921 | −0,995 | 0,40 | 0,205 | 0,0784 |
| 9973 | 0,353 | 0,869 | −1,011 | 0,39 | 0,256 | 0,0467 |
| 99991 | 0,374 | **0,861** | **−0,867** | 0,39 | 0,292 | **0,0265** |
| **REAL** | **0,348** | **0,866** | **−0,924** | **0,42** | **0,267** | — |

*(El modelo es determinista; la estabilidad se verificó con dos semillas de
nulos: 0,353/0,869/−1,011 contra 0,346/0,867/−1,015 — idénticos.)*

## 2 · Las convergencias, contra su lista

- **Amplitud**: 0,267 → 0,315 → 0,353 → 0,374 — **atraviesa** el 0,348 real en
  el tercer escalón (leve sobre-tiro en el cuarto; no monótona-al-blanco, y se
  reporta tal cual).
- **Cruce**: 0,847 → 0,921 → 0,869 → **0,861** contra 0,866 real — en los dos
  escalones profundos, a **cinco milésimas**.
- **Espaciados**: 0,42 → 0,40 → 0,39 — deriva LEVE hacia abajo (7 %); no se
  corrige, se reporta (su §9).
- **Pendiente**: −0,867 contra −0,924. **T4**: 0,215 → 0,292, encerrando el
  0,267 real.
- **El residuo total colapsa monótono: 0,096 → 0,027.** A N = 99991 la curva
  entera del modelo acompaña a la real bin por bin dentro de 0,027 — la
  amplitud de la señal es 0,35.

Por su tabla del §18: fila cuarta — **resultado fuerte → reproducir y auditar
independientemente.** Se frenó ahí.

## 3 · El primer elemento de la auditoría, adelantado — y lo que revela

Antes de interpretar, la pregunta que su §16 manda hacerse: ¿hay una dependencia
no-accidental entre el modelo y aquello con lo que se compara? La medición
directa — distancia de cada punto del modelo al cero real **más cercano**:

| N | distancia media | máx | a < 0,1 del cero real |
|---:|---:|---:|---:|
| 97 | 0,060 | 0,231 | 81 % |
| 997 | 0,035 | 0,468 | 99 % |
| 9973 | 0,024 | 0,084 | **100 %** |
| 99991 | **0,018** | 0,100 | **100 %** |

**El modelo no es un modelo que "se parece" al espejo: a truncación profunda,
sus puntos SON los ceros de Riemann, reconstruidos desde los primos** — a menos
de dos centésimas de espaciado, cada uno. La ecuación de conteo con S completa
es, en el límite, la definición misma de los ceros: la convergencia de la
Fase XX es la dualidad ceros↔primos operando, no un mecanismo independiente que
casualmente acierta.

## 4 · Lo que esto cierra, y con qué palabras

Con las reglas de su §12 («converger ≠ explicar») y esta auditoría adelantada,
la conclusión honesta de la campaña del Espejo (Fases XV–XX) es:

1. **La curva del Espejo del Punto Medio es aritmética pura.** Cruce en 0,866,
   pendiente −0,92, amplitud 0,35, T4 — todo está codificado en la fórmula
   explícita, sin una gota de aleatoriedad independiente. La cadena de fases lo
   demostró por eliminación (densidad ✗, desplazamiento ✗, acoplado parcial,
   rigidez externa ✗) y por reconstrucción (campo puro ✓ convergente).
2. **La imagen de Berry, medida en nuestro observable**: la estadística local
   "tipo GUE" de los ceros no necesita ruido externo — los términos profundos
   de la misma fórmula la fabrican (a N = 97 los puntos ya siguen a los ceros a
   0,06; el resto de la suma sólo afina).
3. **No es física nueva ni una explicación mecánica nueva**: es la dualidad de
   siempre, ahora demostrada de punta a punta sobre el observable que el
   Teorema del Punto Medio del capitán abrió. El valor está en el CIRCUITO
   COMPLETO, medible y reproducible: teorema a mano → espejo → curva → identidad
   de cuatro términos → eliminación de mecanismos → reconstrucción convergente.

## 5 · La auditoría que queda (su §16, para su orden)

Adelantamos el elemento 1 (dependencia modelo↔referencia: revelada y central).
Quedan, si usted los ordena: la segunda implementación independiente del
cálculo, el efecto del rango de γ, y la sensibilidad numérica fina. Con el
elemento 1 a la vista, la pregunta para usted es si esos pasos agregan algo —
la "coincidencia" ya tiene explicación exacta: el modelo converge a los ceros
mismos.

HL y el Telar siguen congelados. Y la deriva de espaciados (0,42 → 0,39) queda
como el único cabo suelto numérico de la fase.

---

XVI: tamaños. XVII: posiciones. XVIII: juntos. XIX: la rigidez muere.
**XX: se quitó todo el ruido, la fórmula habló — y lo que dijo es que ella ya
era los ceros.**

Estructura cerrada no es hipótesis demostrada.

**Todavía no.**
