# Acta de las Relaciones — Fase XI

**Respuesta al Telar Fase XI (Auditoría 44) · 2026-08-19**

Su pregunta central era: *¿puede una regla aritmética sobre RELACIONES entre
sitios producir rigidez global, manteniendo participación, de una forma que no
se explique por distancia, azar o simetría?*

Respuesta corta, con todas las letras: **con el catálogo congelado, no se pudo
ni preguntar — y ése es el resultado.** Las tres reglas aritméticas clásicas
tienen densidades naturales de signo que **vacían la banda**, así que ninguna
llega admisible al examen. A densidad igualada, B ≈ A ≈ C dentro del ruido. Y su
§13 quedó decidido: el borde q = 0,05 era **accidental**.

Se reproduce con `go run ./cmd/lasrelaciones`.
Lámina: `galeria/laminas/10-el-telar/las-relaciones.svg`.

---

## 1 · La caja negra, abierta y congelada (su §4)

La regla exacta del brazo «por distancia» de la Fase X, sin descripción verbal:

```
s[k] = −1 con probabilidad q, sorteo INDEPENDIENTE por cada k ∈ {1..120}
       (xorshift64, semilla 20260818 — el mismo dado de la Fase X)
signo del enlace (i,j) = s[|i−j|]     ⟹  matriz de signos TOEPLITZ
sin normalización extra: signo² = 1 deja la fuerza F = 30 intacta
dominio: k = 1..120 = kmax; fuera de rango, +1
```

El PR/N = 0,281 citado corresponde a **q = 0,10 con 6 semillas**. La simetría
que introduce (S depende sólo de k) queda declarada como pide su §11, y la capa
C existe precisamente para vigilarla.

## 2 · El catálogo, congelado antes de mirar (su §12)

Tres reglas, sólo enteros y factorización, transformación a signo declarada:

- **S_Λ**: s[k] = −1 si k es potencia de primo, +1 si no → densidad natural 0,35
- **S_μ**: s[k] = −1 si μ(k) = −1, +1 si μ(k) ∈ {0, +1} → densidad 0,33
- **S_λ**: s[k] = λ(k) de Liouville = (−1)^Ω(k) → densidad 0,54

Nada se agregó después de ver resultados. Las tres se corrieron y se publican
enteras, incluidas las que pierden.

## 3 · Las tres capas, y el muro

Rieles vigentes de la Fase X: vivos ≥ 149 de 187 · PR/N ≥ 0,090.

| brazo | vivos | Σ²(10) | PR/N | dens− | riel |
|---|---:|---:|---:|---:|---|
| S0 · control | 187 | 18,335 | 0,105 | 0 | ok |
| B · S_Λ | **98** | 3,734 | 0,291 | 0,346 | **VACIADA** |
| A · pura q=0,333 (6 sem.) | 109 | 2,376 ± 0,483 | 0,291 | 0,338 | VACIADA |
| C · permutada (6 sem.) | 105 | 2,912 ± 1,076 | 0,291 | 0,334 | VACIADA |
| B · S_μ | **100** | 3,226 | 0,286 | 0,326 | **VACIADA** |
| A · pura q=0,325 | 107 | 3,332 ± 1,913 | 0,290 | 0,313 | VACIADA |
| C · permutada | 108 | 2,738 ± 1,463 | 0,290 | 0,325 | VACIADA |
| B · S_λ | **99** | 1,527 | 0,296 | 0,540 | **VACIADA** |
| A · pura q=0,542 | 102 | 1,940 ± 0,458 | 0,290 | 0,542 | VACIADA |
| C · permutada | 104 | 3,616 ± 1,144 | 0,286 | 0,545 | VACIADA |

**El muro:** toda regla de distancia con densidad ≥ 0,3 expulsa la mitad del
espectro. Los Σ² chicos de esta tabla (1,5–3,7) se miden sobre ~100 niveles de
400 — son la forma del recorte, no del medio. Por la disciplina de la Fase VII,
**ninguna fila es interpretable**, y se dice así.

**Y a densidad igualada, dentro de ese territorio inadmisible: B ≈ A ≈ C.**
S_Λ salió incluso *peor* que su distancia pura (−2,81 σ); S_μ empatado (±0,3 σ);
S_λ el único con un asomo contra su permutada (+1,83 σ, bajo el umbral de 2 y
sobre banda vaciada — no cuenta). La organización aritmética exacta k→S(k) no
aporta nada que su propia permutación no aporte. **Es la quinta aplicación de
la misma hoja, quinta vez la misma respuesta.**

## 4 · El eco no se corrió, por regla

Su §10 exigía la secuencia operador → espectro → eco → log p **sólo** sobre un
brazo admisible. No hubo brazo admisible. No se fuerza la conexión con La
Armonía (su §16, último criterio de falsación, aplicado tal cual).

## 5 · Su §13: el borde, decidido

q = 0,05 con 12 semillas: vivos **146,8 ± 2,7**, pasa el riel **2 de 12 veces**,
Σ²(10) = 2,584 ± 0,436, PR/N = 0,225. **Accidental.** El brazo del borde no es
un punto de trabajo: es una moneda cargada en contra. Cerrado.

## 6 · Lo que queda en pie

1. **El fenómeno real de la Fase X vive en densidades CHICAS** (q ≤ 0,02: el
   único brazo que ordena manteniendo banda y participación). Ninguna función
   aritmética clásica es naturalmente rala a escala k ≤ 120 — Λ, μ y λ traen
   0,33–0,54 de signos negativos puestos.
2. **Si la aritmética quiere entrar por las relaciones, necesita una dilución
   declarada de antemano** — por ejemplo, aplicar el signo aritmético sólo
   sobre un subconjunto ralo de distancias, con la regla de dilución congelada
   ANTES de correr. Eso es un catálogo nuevo y, si usted lo aprueba, una fase
   nueva. No se probó acá porque agregarlo después de ver los resultados es
   exactamente lo que su §12 prohíbe.
3. Su cierre lo dijo antes que nosotros: **cerrar una pista falsa también es
   progreso.** Van cinco idiomas eliminados — dureza escalar (VIII), onda no
   aritmética (VIII), signo de sitio (IX–X), reciprocidad (X), y ahora Λ/μ/λ
   como reglas de distancia a densidad natural.

---

El pedido de abrir la caja negra, las tres capas y el congelamiento previo, de
la auditoría — y los tres se cumplieron a la letra. El catálogo se publicó
entero, con sus derrotas.

Estructura cerrada no es hipótesis demostrada.

**Todavía no.**
