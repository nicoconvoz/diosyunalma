# Acta del Porqué — Fase XVI

**Respuesta a Fase XVI (Auditoría 49) · 2026-08-20 · F368**

Su regla metodológica se cumplió al pie: **derivar primero, comparar después —
y si una explicación conocida mata la novedad, decirlo.** El resultado es el
contrario y se dice con números: **la explicación conocida explica un tercio de
la curva, y la identidad algebraica nombra con precisión el término que falta.**

Se reproduce con `go run ./cmd/laformulamadre fase16` (semilla 20260823).
Lámina: `galeria/laminas/10-el-telar/el-porque.svg`.

---

## 1 · La derivación exacta (sus §2–4): la identidad de cuatro términos

Con m = γₙ + gₙ/2, la trigonometría da, exacta:

```
cos(mT) = cos(γₙT)·cos(gₙT/2) − sin(γₙT)·sin(gₙT/2)
```

y contra el nulo analítico (fase poblacional × factor de gap del bin), el
excedente se descompone **como identidad, no como metáfora**:

```
E(s) = T1 + T2 + T3 + T4
T1 = −2·⟨Cov_bin(cos γT, cos gT/2)⟩     covarianza interna coseno
T2 = −2·⟨(C_bin − C_pob)·c̄_bin⟩         SELECCIÓN × factor, coseno
T3 = +2·⟨Cov_bin(sin γT, sin gT/2)⟩     covarianza interna seno
T4 = +2·⟨(S_bin − S_pob)·s̄_bin⟩         SELECCIÓN × factor, seno
```

**Verificación de la identidad, bin por bin:** la suma T1+T2+T3+T4 reproduce el
E(s) medido (con nulo barajado) en todos los bins — p. ej. 0,327 contra 0,347;
−0,283 contra −0,287 en las puntas. De paso queda verificado que el barajado
estima bien el producto de marginales (el nulo analítico ≈ el nulo empírico).

**El factor cos(gT/2) de su §4, derivado y medido:** F(s) va de +0,455 (pares
apretados) a −0,643 (anchos), cruzando en s ≈ 0,78. Es la parte esperable del
corrimiento — y NO alcanza para explicar E(s), porque multiplica a la selección
pero no la crea.

**La respuesta a su §3 (¿qué parte es de qué?):** la masa de los términos de
SELECCIÓN es **2,33** contra **0,38** de las covarianzas internas — el 86 % del
espejo es selección: *qué ceros terminan en cada clase de gap*. Y dentro de la
selección **domina T4, el término seno·seno** — el dato fino de la fase: la
correlación fuerte no es con la cresta de densidad (coseno) sino con el
**corrimiento** (seno): los ceros de pares apretados no sólo viven donde la onda
de los primos aprieta — viven **corridos con ella**.

## 2 · GUE puro (sus §6 y §10): descartado, con números

Tres espectros GUE puros (gaps Wigner sobre la densidad lisa real, sin primos),
por el pipeline idéntico: **|E| máximo 0,050 contra 0,35 del real — plano.**
Sin primos no hay eco que repartir: GUE solo no produce la curva. La curva no es
estadística de matriz aleatoria a secas.

## 3 · El modelo mecanismo: GUE + fórmula explícita

Construido como manda su §13 — derivado antes de comparar: densidad modulada
por la fórmula explícita truncada en los MISMOS n ≤ 97 que mide el eco,

```
ε(γ) = −(2/log(γ/2π)) · Σ_{n≤97} Λ(n)·cos(γ log n)/√n
```

y gaps Wigner con espaciado local 1/(ρ̄·(1+ε)). Tres semillas, mismo pipeline.

| | real | modelo GUE+fórmula |
|---|---:|---:|
| signo y caída monótona | ✓ | ✓ |
| amplitud | ±0,35 | **±0,10 (⅓)** |
| cruce s\* | 0,863 | **0,70 ± 0,01** |
| pendiente | −0,92 | **−0,6** |

**El modelo produce la CLASE de curva pero se queda corto en todo lo
cuantitativo.** El residuo (real − modelo, su §15) conserva la mayor parte de la
estructura: llega a 0,34 donde la señal real llega a 0,35.

## 4 · El veredicto (su §10, rama «ningún modelo reproduce»)

Su §10 preveía las dos salidas. Cayó la segunda: *«si ningún modelo reproduce la
estructura, identificar qué término falta»*. **Identificado, con nombre de
pila:**

- El modelo de densidad-modulada mueve los **tamaños** de los gaps donde la
  onda encresta. Eso produce selección tipo T2 (coseno) — y por eso da la curva
  chica.
- Pero la identidad dice que el espejo real está dominado por **T4** — el
  acople seno·seno, que es **corrimiento coherente de las POSICIONES** de los
  ceros con la onda de los primos, no sólo compresión local de gaps.
- **El término que falta es el campo de desplazamiento**: el próximo modelo
  debe mover las posiciones con la onda (el corrimiento que la fórmula
  explícita también dicta), no sólo modular la densidad. Ésa es la predicción
  concreta y falsable que esta fase deja lista.

## 5 · La pendiente −1 (su §8)

Recorridas sus seis opciones: no es la normalización (sobrevive al desdoblado
global, Fase XV), no es el procedimiento de control (el segundo control aplana
todo), no es GUE puro (plano), y el modelo densidad-sola da −0,6 — **la
pendiente completa −0,92 pertenece al término faltante T4** igual que el resto
del déficit. Queda como propiedad del mecanismo completo, pendiente de la
versión con desplazamiento.

## 6 · Hardy–Littlewood (su §7): todavía NO pertinente

Su criterio: HL entra sólo si existe una predicción que explique lo que GUE no
explica. Lo que GUE+fórmula no explica quedó identificado como el término de
desplazamiento de la MISMA fórmula explícita — no hace falta (todavía) una
constante de pares de primos para nombrar el déficit. HL queda cerrado para
esta fase y re-evaluable cuando el modelo con desplazamiento diga qué residuo
deja. El Telar sigue congelado, como usted mandó.

## 7 · Los cuatro objetos de su §15, entregados

E observado(s) ✓ · E GUE-puro(s) ≈ 0 ✓ · E GUE+fórmula(s) ✓ · residuo bin por
bin ✓ — impresos en la corrida y dibujados en la lámina.

---

La pregunta de la fase era «¿por qué cruza?». La respuesta que alcanzamos:
**cruza porque la clase del gap selecciona la fase del cero en la onda de los
primos — un 86 % selección, 14 % covarianza interna — y de esa selección, la
parte mayor es corrimiento (seno), no densidad (coseno). La mitad conocida del
mecanismo da un tercio de la curva; el término que falta tiene nombre, T4, y
modelo propuesto para la próxima fase.**

Estructura cerrada no es hipótesis demostrada.

**Todavía no.**
