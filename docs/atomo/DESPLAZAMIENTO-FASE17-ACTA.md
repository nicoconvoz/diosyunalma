# Acta del Desplazamiento — Fase XVII

**Respuesta a Fase XVII (Auditoría 50) · 2026-08-20 · F369**

Su orden fue: *no persigas la curva — construí una causa que pueda fallar.* Se
construyó, sin una sola perilla, se pre-registró la predicción, y se la sometió
a juicio. **Falló donde sus criterios del §18 mandan declararlo, y dejó dos
huellas verdaderas que reforman la hipótesis en vez de enterrarla.**

Se reproduce con `go run ./cmd/laformulamadre fase17` (semilla 20260824).
Lámina: `galeria/laminas/10-el-telar/el-desplazamiento.svg`.

---

## 1 · La regla, sin circularidad (sus §4 y §8)

Derivada de la misma fórmula explícita, **cero parámetros libres**:

```
S(γ)  = −(1/π)·Σ_{n≤97, Λ>0} Λ(n)·sin(γ log n)/(√n·log n)
δγ(γ) = −S(γ)/ρ̄(γ)          ρ̄ = log(γ/2π)/2π
```

Un cero clavado por N_liso(γ) + S(γ) = n − ½ queda corrido de su posición lisa
exactamente en −S/ρ̄: la **amplitud la fija la fórmula**, la truncación n ≤ 97
es la de las fases anteriores. Nada se leyó de la curva real — ni cruce, ni
amplitud, ni pendiente, ni residuo. La predicción quedó pre-registrada en el
encabezado del código antes de correr.

## 2 · La tabla que importa (su §7)

| modelo | amplitud | cruce s\* | pendiente |
|---|---:|---|---|
| **REAL** (3474 ceros) | **0,351** | 0,863 | −0,93 |
| a) GUE puro | 0,026 | 0,58–0,65 | ≈ 0 |
| b) GUE + densidad | 0,124 | 0,67–0,90 | −0,06…−0,50 |
| **c) GUE + DESPLAZAMIENTO** | **0,082** | 0,61–0,70 | −0,11…−0,32 |
| d) densidad + desplazamiento | 0,082 | 0,66–0,72 | −0,20…−0,45 |
| e) fase DESTRUIDA | **0,029** | dispersos | ≈ 0 |
| f) desplazamiento ×0,5 | 0,115 | 0,76 | −0,23…−0,33 |
| f) desplazamiento ×2 | 0,051 | 0,65–0,86 | −0,11…−0,29 |

Residuos máximos contra el real: densidad 0,312 · desplazamiento 0,300 ·
combinado 0,341.

## 3 · El veredicto contra sus criterios, punto por punto

**Del §18 (fracaso):**
- *«El desplazamiento no mejora al modelo de densidad»* — **SE CUMPLE**: 0,082
  contra 0,124 en amplitud, cruce que no avanza hacia 0,863, pendiente que no
  se empina, residuo prácticamente igual. **La hipótesis, en su forma lineal,
  FALLA el examen cuantitativo y se declara.**
- Nota de robustez (su §11): la respuesta a la amplitud no es monótona (×0,5 da
  0,115; ×1 da 0,082; ×2 da 0,051) — el corrimiento grande EMBORRONA el eco en
  vez de amplificarlo. Reportado como robustez, no usado para elegir.

**Del §17 (lo que sí pasó):**
- *«Pierde la señal cuando se destruye la coherencia de fase»* — **SE CUMPLE**:
  con las fases φₙ aleatorizadas (misma distribución de amplitudes) el brazo cae
  a 0,029 — el nivel del GUE puro. Lo poco que el desplazamiento produce, lo
  produce **por coherencia**, no por el tamaño del corrimiento.
- **La prueba clave de su §9 — T4 — dio la segunda huella:** el desplazamiento
  genera la estructura seno·seno con el **patrón de signos correcto** (positivo
  en apretados, negativo en anchos), cosa que el modelo de densidad NO hace (su
  T4 es ruido incoherente). Pero la magnitud es ⅓ de la real (0,085 contra
  0,267 en el primer bin) — el mismo factor 3 que persigue a toda esta campaña.

## 4 · Conclusión (su §19): REFORMULAR, con dirección precisa

Ni mecanismo candidato ni hipótesis muerta: **parcialmente explicada, a
reformular** — y la reforma tiene dirección señalada por los propios datos:

1. El desplazamiento lineal **sí** es el único brazo que produce T4 con la
   forma correcta, y **sí** depende de la coherencia de fase (los dos sellos
   cualitativos de la hipótesis).
2. Lo que no alcanza es la **magnitud**, y el defecto del modelo es visible en
   su construcción: aplicamos el corrimiento SOBRE gaps Wigner ya sorteados —
   desplazamiento y repulsión actúan en secuencia, no juntos. En los ceros
   reales la repulsión opera sobre posiciones ya corridas: el sistema es
   **auto-consistente**, y esa retroalimentación es exactamente lo que un
   corrimiento aplicado a posteriori no puede fabricar.
3. **Hipótesis reformulada, para su evaluación:** un modelo auto-consistente —
   repulsión GUE y campo de la fórmula explícita resolviéndose JUNTOS (por
   ejemplo, gas de Coulomb 1D en el potencial de la fórmula explícita, o la
   dinámica secuencial con el corrimiento aplicado dentro del sorteo del gap,
   no después). Cero parámetros nuevos: los mismos objetos, acoplados en vez de
   encadenados.

## 5 · Lo que sigue congelado

Hardy–Littlewood sigue afuera (el déficit sigue nombrándose con la propia
fórmula explícita). El Telar sigue congelado. Los tres, hasta su orden.

---

La causa se construyó para poder fallar, y falló donde debía fallar — dejando
mejor delimitado lo que la próxima causa tiene que ser. Eso también es avanzar:
ahora sabemos que el espejo no se explica ni moviendo tamaños (XVI) ni moviendo
posiciones a posteriori (XVII); pide las dos cosas **a la vez**.

Estructura cerrada no es hipótesis demostrada.

**Todavía no.**
