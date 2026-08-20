# Acta del Punto Medio — Fase XIII

**Respuesta a Fase XIII (Auditoría 46) · 2026-08-19**

Su orden fue inusual y se cumplió a la letra: **ninguna batería espectral.**
Volver al origen, generalizar la identidad de Nico sin reemplazarla por q−p=2,
introducir g, el centro m y las anclas, decidir el estatus de los ±1/2, y si no
hay nada más, demostrarlo.

Hay algo más. Y también hay una parte que se reduce a q−p=g, y se declara.

Por orden del capitán, el resultado queda bautizado **TEOREMA DEL PUNTO MEDIO**
e inscripto como **Teorema 6** en el libro de los teoremas del taller
(`docs/teoremas/TEOREMAS.md`), con su alcance y su honestidad declarados.

Se reproduce con `go run ./cmd/puntomedio`.
Lámina: `galeria/laminas/10-el-telar/el-punto-medio.svg`.

---

## 1 · La demostración formal (su §16.1)

**Teorema del Punto Medio (Nico).** Para impares p < q:

```
(p+1)/2 = (q−1)/2   ⟺   q = p + 2
```

y el valor compartido c es entero, con p = 2c−1, q = 2c+1.

*Prueba.* Multiplicando por 2: p+1 = q−1 ⟺ q = p+2. Como p y q son impares,
(p+1)/2 y (q−1)/2 son enteros. Si c es el valor común, 2c−1 = p y 2c+1 = q. ∎

**Verificado además por fuerza bruta** sobre ~25 millones de pares impares
≤ 20 000: cero fallas. Y una precisión que importa: **la forma base es una
identidad de los IMPARES** — no usa primalidad. Lo específico de los primos
aparece recién en la ley del centro (§3).

## 2 · La generalización por gap (su §16.2–3)

Sean p < q primos impares, g = q−p (par), r = g/2, m = (p+q)/2, y las anclas
a±(n) = (n±1)/2.

**(I) Geometría pura:** m ∈ ℤ y (p, q) = (m−r, m+r). Toda pareja de impares es
un par simétrico alrededor de un centro entero.

**(II) La identidad de las anclas — la forma general que contiene el caso g=2:**

```
a⁻(q) − a⁺(p) = r − 1
```

| g | r | Δ anclas | lectura |
|---:|---:|---:|---|
| 2 | 1 | **0** | ancla compartida — el teorema del capitán |
| 4 | 2 | 1 | anclas adyacentes: se tocan |
| ≥6 | g/2 | r−1 | hueco de exactamente (g−4)/2 enteros |

Verificada sobre 17 982 pares consecutivos de primos ≤ 200 000: cero fallas. El
diccionario geométrico entero de la Fase XII es esta única identidad leída caso
por caso.

## 3 · La ley del centro — lo que los primos agregan

**(III) Para p, q > 3 primos: 3 | m ⟺ 6 ∤ g.**

*Prueba.* Módulo 3: p, q primos > 3 ⟹ p, q ≢ 0. Si 6 ∤ g entonces r ≢ 0
(mod 3). De los dos residuos posibles de p (1 o 2), la exclusión
q = p + 2r ≢ 0 elimina exactamente uno; en la rama sobreviviente m = p + r ≡ 0.
Si 6 | g, entonces r ≡ 0 y m ≡ p ≢ 0. ∎

Verificada sobre **todos los 304 590 pares de primos en (3, 6000]** — no sólo
consecutivos: cero fallas.

**La tabla de los centros** (pares consecutivos ≤ 200 000):

| g | m mod 6 | pares |
|---:|---|---:|
| 2 | **{0}** | 2159 |
| 4 | {3} | 2134 |
| 6 | {2, 4} | 3455 |
| 8 | {3} | 1391 |
| 10 | {0} | 1705 |
| 12 | {1, 5} | 1821 |
| 14 | {0} | 959 |
| 16 | {3} | 637 |
| 18 | {2, 4} | 1035 |

**La firma es periódica en g con período 6.** Los centros de cada clase de gap
viven en progresiones aritméticas mod 6 **disjuntas por clase**: la clase del
gap se lee en el residuo del centro. El caso gemelo es la clase m ≡ 0 (mod 6),
y su ancla compartida c = m/2 es siempre múltiplo de 3.

## 4 · El estatus de los ±1/2 (su §16.4) — veredicto pedido, veredicto dado

Su §8 pedía decidir entre trivial, estructural o falso. La respuesta honesta es
**doble**:

- **Como aritmética: TRIVIALES.** Son la paridad — el paso de un impar a sus dos
  enteros vecinos. No crean estructura nueva. (Su «caso trivial», confirmado.)
- **Como coordenada: ESTRUCTURALES.** T(n) = (n−1)/2 es la biyección estándar
  impares → enteros, y **divide todo gap por 2** (verificado: T(q)−T(p) = g/2,
  17 982 pares, cero fallas). Bajo esa lente los gemelos son *exactamente* los
  pares de imágenes consecutivas, y el ancla compartida es c = m/2. Es un cambio
  de coordenadas que compacta la pregunta de los gaps — útil, no creador.

Y el «caso falso» de su §8 queda descartado: al generalizar NO queda sólo
q−p=2 — queda la identidad (II) y la ley (III).

## 5 · Prueba de honestidad (su §14), sin anestesia

1. La forma base **se reduce a paridad**: identidad de los impares. Declarado.
2. Los ±1/2 **no aportan aritmética** más allá de la coordenada. Declarado.
3. La ley del centro **es matemática clásica** (criba mod 3): verdadera,
   demostrada, y no nueva para el mundo — nueva como organización para este
   taller. Declarado sin vueltas: acá no se descubrió matemática desconocida;
   se demostró, ordenó y nombró la que el hallazgo manual tenía adentro.
4. Consecuencia demostrable sin simulación (su §16.6): para todo par gemelo
   > (3,5), el punto medio m ≡ 0 (mod 6) y el ancla c ≡ 0 (mod 3). Los centros
   gemelos viven en la red 6ℤ — sin excepción posible, por (III).

## 6 · El puente al operador: propuesta en papel, sin correr nada (su §10 y §16.7)

Su §15 prohibió correr. Se obedece. Lo que queda escrito para cuando usted
decida:

- La Fase XII probó **pares** (la proyección demasiado específica, como usted
  dijo). La ley del centro sugiere probar **centros**: una máscara definida no
  por «estos dos sitios son gemelos» sino por «este sitio está en la red de
  centros m ≡ 0 (mod 6)» — una propiedad de SITIO derivada de relaciones, que
  invierte la lección de la Fase X (sitio localiza, enlace no) en territorio
  nuevo: habría que decidir primero, en papel, si una marca de sitio sobre la
  red 6ℤ es gauge o no lo es (la Fase IX probó que una FASE de sitio es gauge
  pura; una marca de *amplitud* de sitio no lo es).
- Falsador anticipado: la red 6ℤ es periódica, así que su control natural es la
  red 6ℤ + corrimiento. Si el corrimiento rinde igual, la aritmética del centro
  no aporta y el canal se cierra en una corrida.
- Nada de esto se ejecuta hasta que usted lo bendiga con la regla congelada.

---

El teorema es de Jesús Nicolás Astorga, encontrado a mano. Las demostraciones,
la generalización, la ley del centro y los veredictos de honestidad, de este
taller. El nombre lo puso el capitán y quedó inscripto en el libro.

Estructura cerrada no es hipótesis demostrada.

**Todavía no.**
