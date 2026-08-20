# Acta del Cruce — Fase XV

**Respuesta a Fase XV (Auditoría 48) · 2026-08-20 · F367**

Su orden: localizar el cruce, medirlo, e intentar matarlo. **No murió.** Y el
intento de matarlo dejó dos cosas más valiosas que el número: la convergencia
con M, y el veredicto del mundo nulo — que es más fino de lo que esperábamos y
se reporta con cuidado.

Se reproduce con `go run ./cmd/laformulamadre fase15` (semilla 20260822).
Lámina: `galeria/laminas/10-el-telar/el-cruce.svg`.

---

## 1 · Lo congelado antes de medir

Zoom s ∈ [0,50, 1,30) en bins de 0,05 (16 bins) · E(s) = real − nulo medio ·
200 barajados por profundidad · 200 remuestreos bootstrap · s\* = interpolación
lineal en el cambio de signo MONÓTONO (su §5, con todos los cruces impresos) ·
pendiente por cuadrados mínimos en ±0,15 de s\*.

## 2 · s\*(M): converge (su §4)

| M | s\* | bootstrap 95 % | dE/ds |
|---:|---:|---|---:|
| 649 | 0,7249 | [0,702, 0,797] | −1,13 |
| 1517 | 0,8476 | [0,764, 0,878] | −1,00 |
| 3474 | **0,8638** | **[0,844, 0,885]** | −0,94 |

El desplazamiento entre profundidades se ACHICA (+0,123 → +0,016) y el
intervalo se estrecha: **el cruce converge hacia s\* ≈ 0,86–0,89 en vez de
derivar.** El único cruce monótono de la curva es ése — no hay cruces múltiples
en el real. **No se declara constante** (su §15): se declara localizado, con
intervalo, y convergente.

## 3 · La pendiente (su §7): cruce SUAVE

dE/ds ≈ −1,0, estable en las tres profundidades (−1,13 / −1,00 / −0,94). Ni
meseta, ni salto, ni ruido: una bajada limpia de pendiente uno por unidad de s.
La forma completa E(s) se superpone entre profundidades con la única evolución
del cruce corriéndose a su límite (su §8: forma estable, sin llamarla universal).

## 4 · El mundo nulo (su §6 — «esto es fundamental»), reportado sin maquillaje

Cada barajado produjo su propia curva y sus propios cruces: **~3,8 cruces por
barajado** (749–780 en 200), dispersos por toda la ventana — mediana ≈ 0,89–0,90
(cerca del centro de la ventana), con ~14 % a menos de 0,05 del s\* real,
compatible con dispersión uniforme (0,10/0,80 = 12,5 %).

**La lectura honesta tiene dos mitades:**
- La *ubicación* del cruce real no es improbable por sí sola dentro del mundo
  nulo — un barajado cualquiera cruza en muchos lados, incluido cerca de 0,86.
- Lo que el mundo nulo NO tiene es **la estructura**: el real cruza UNA vez,
  monótono, con pendiente −1 y amplitud tres veces mayor; los nulos cruzan
  ~cuatro veces cada uno, en zigzag chico. Lo distintivo no es el punto: es la
  curva entera que pasa por ese punto. (Su §17, cumplida: perseguimos la forma,
  no el sigma.)

## 5 · La batería de refutación (su §13): sobrevive entera

| ataque | cruce |
|---|---:|
| malla corrida media celda | 0,863 |
| desdoblado GLOBAL (el de F365) | 0,844 |
| mitad BAJA de los ceros | 0,847 |
| **mitad ALTA (réplica INDEPENDIENTE)** | **0,895** |

Todos en [0,84, 0,90]. Las dos mitades son subconjuntos disjuntos — la única
réplica genuinamente independiente disponible — y encierran al valor del total.

## 6 · El segundo control (su §10): la curva es del emparejamiento

Centros intactos, distribución de gaps intacta, cantidad de pares intacta —
sólo se permutaron las etiquetas s. Resultado: **curva PLANA** (−0,072 ± 0,031,
sin dependencia alguna en s). Toda la estructura E(s) vive en la asociación
centro↔gap; nada viene de las distribuciones marginales.

## 7 · Criterio de avance (su §12): CUMPLIDO

✔ aparece en las tres profundidades · ✔ estable ante binning · ✔ el mundo nulo
no reproduce la estructura (aunque sí toca la ubicación, y se declara) ·
✔ incertidumbre final ±0,02 · ✔ la forma completa se conserva.

**Por su propia regla, recién ahora se habilita la pregunta «¿POR QUÉ existe
s\*?» — y con ella Hardy–Littlewood.** Queda pedida su orden para esa etapa.
Nota al margen para entonces, sin cargarla de sentido: s\* converge hacia la
zona 0,86–0,89, y la mediana de espaciados GUE anda cerca de 0,80 — comparación
que NO hicimos formalmente y que pertenece a la fase siguiente.

## 8 · Lo que no se hizo

Ni Hardy–Littlewood, ni el Telar, ni función ajustada por los puntos, ni
constante declarada. La dependencia entre profundidades sigue declarada (los M
son prefijos; la única réplica independiente es la de las mitades).

---

El protocolo de esta fase es de la auditoría, cumplido punto por punto,
incluido el intento de asesinato. El cruce salió vivo, con domicilio
[0,844, 0,885] y pendiente −1.

Estructura cerrada no es hipótesis demostrada.

**Todavía no.**
