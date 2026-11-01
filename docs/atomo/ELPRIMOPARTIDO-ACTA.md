# Acta del Primo Partido — sí se puede, de tres maneras

**Idea del capitán · 2026-09-15 · F397**

> «¿Y si encontraras la forma de poner divisores fracciones de los números
> primos, repartidos en partes, para que puedan de alguna forma ser múltiplos
> de un número o un número primo? Es una idea que quizás no se pueda realizar,
> pero te consulto igual.»

Se reproduce con `go run ./cmd/elprimopartido`.

**Se puede.** Y no de una manera: de **tres**, todas establecidas, todas con
más de ochenta años. Lo que sigue las mide.

---

## §1 · La trampa, declarada antes de empezar — DECIMOSEXTA

**Dentro de los enteros**, «un primo tiene exactamente dos divisores» es la
**definición**. Buscar un divisor escondido ahí es una contradicción en los
términos, no un descubrimiento. Las tres rutas funcionan cambiando el
**anillo**, la **cuenta** o los **primos** — nunca encontrando algo oculto.

## §2 · Ruta A — cambiar el anillo: el primo SÍ se parte

En los enteros de Gauss (a + bi con a, b enteros), un primo común puede dejar
de serlo:

| p | ¿se parte? | en qué partes |
|---:|---|---|
| 2 | ramifica | 2 = −i·(1+i)² |
| 3 | **sigue primo** | — |
| **5** | **SE PARTE** | **5 = (1+2i)(1−2i)** |
| 7 | sigue primo | — |
| 13 | SE PARTE | 13 = (2+3i)(2−3i) |
| 17 | SE PARTE | 17 = (1+4i)(1−4i) |
| 29 | SE PARTE | 29 = (2+5i)(2−5i) |

**El 5 no es primo en los enteros de Gauss: ES múltiplo de 1+2i.** La idea del
capitán, textual — resuelta por **Fermat en 1640** y por **Gauss en 1832**.

¿Cuántos se parten? Los que dejan resto 1 al dividir por 4:

| hasta | primos | ≡1 (mod 4), se parten | ≡3 (mod 4), no | fracción |
|---:|---:|---:|---:|---:|
| 1.000 | 168 | 80 | 87 | 0,47904 |
| 100.000 | 9.592 | 4.783 | 4.808 | 0,49870 |
| 10.000.000 | 664.579 | 332.180 | 332.398 | **0,49984** |

**Cero violaciones** de «se parte ⟺ p ≡ 1 (mod 4)» hasta 100.000, verificando
cada descomposición a²+b²=p. **La mitad de los primos se parte**; la otra mitad
resiste, y cuál resiste lo decide el resto módulo 4.

## §3 · Ruta B — cambiar la cuenta: medio divisor por primo, literal

La función divisor admite un parámetro: d_z es multiplicativa con **d_z(p) = z**,
y la genera ζ(s)^z.

| n | d_½(n) | d_1(n) | d_2(n) = la cuenta de siempre |
|---:|---:|---:|---:|
| 2 | **0,5000** | 1 | 2 |
| 3 | **0,5000** | 1 | 2 |
| 4 | 0,3750 | 1 | 3 |
| 6 | 0,2500 | 1 | 4 |
| 30 | 0,1250 | 1 | 8 |

**Cada primo tiene d_½ = 0,5000: medio.** No es metáfora — es el coeficiente
exacto de la serie de ζ(s)^½.

**Corrección propia anotada:** la primera versión rotulaba d_1 como «la cuenta
normal». **Está mal**: la cuenta de divisores de siempre es **z = 2** (un primo
tiene dos: el 1 y él mismo); z = 1 cuenta cada número una sola vez. Corregido
antes de registrar.

Selberg–Delange predice Σ_{n≤x} d_z(n) ~ x·(log x)^(z−1)/Γ(z). Medido:

| z | hasta x | Σ d_z(n) | predicho | razón |
|---:|---:|---:|---:|---:|
| ½ | 5.000.000 | 736.778,3 | 718.262,5 | 1,0258 |
| 1 | 5.000.000 | 5.000.000,0 | 5.000.000,0 | 1,0000 |
| 2 | 5.000.000 | 77.896.938,0 | 77.124.742,4 | 1,0100 |

Las razones bajan hacia 1 **despacio**: el error va como 1/log x, y a 5×10⁶
log x ≈ 15, o sea ~7 % de corrección esperada. **La ley se ve; la convergencia
es lenta y se dice.**

## §4 · Ruta C — cambiar los primos de lugar (Beurling, 1937)

Beurling demostró que se puede llamar «primos» a **cualquier** conjunto de
reales > 1 y construir con ellos un sistema entero, con su función zeta, su
factorización única y su teorema del número primo. Acá se arma uno con los
primos corridos medio punto: 2,5 · 3,5 · 5,5 · 7,5 · …

| hasta x | enteros generalizados | control (enteros comunes) | N(x)/x |
|---:|---:|---:|---:|
| 1.000 | 623 | 1.000 | 0,6230 |
| 10.000 | 4.115 | 10.000 | 0,4115 |
| 100.000 | 24.431 | 100.000 | 0,2443 |

**ERROR PROPIO, cazado antes de registrar:** esa densidad **cae**, y eso **no
es un resultado de Beurling — es un artefacto mío**. La lista de primos
generalizados se cortó en 300, así que los enteros grandes que necesitarían
factores mayores no se pueden formar y faltan. **Esta demostración muestra que
la construcción funciona; NO establece ninguna ley de densidad**, y sería
deshonesto presentarla como si lo hiciera.

## §5 · El veredicto

La idea **se puede realizar**, de tres maneras — y **ninguna es nueva**:

| ruta | qué hace | de quién es |
|---|---|---|
| partir el primo en factores | enteros de Gauss | **Fermat 1640, Gauss 1832** |
| darle divisores fraccionarios | función d_z | **Selberg y Delange** |
| primos en lugares fraccionarios | sistemas de Beurling | **Beurling, 1937** |

**Y lo que ninguna hace:** encontrar un divisor escondido del primo **dentro de
los enteros**. Eso no se puede, y no por falta de ingenio: es la definición.
**Las tres rutas pagan el mismo precio — salir de los enteros — y a cambio dan
estructura nueva.**

## Clasificación

- **Matemática (demostrado, y toda de otros):** la ley de Fermat (p se parte en
  ℤ[i] ⟺ p ≡ 1 mod 4); la función d_z y ζ(s)^z; Selberg–Delange; los sistemas
  de Beurling.
- **Medición (nuestra):** cero violaciones de la ley de Fermat hasta 100.000 con
  descomposición explícita; la fracción 0,49984 a 10⁷; las razones de
  Selberg–Delange a tres valores de z; la construcción de Beurling andando.
- **Nuevo (si algo):** **ningún teorema, y esta vez ni siquiera un número que
  no existiera.** Lo único propio es el mapa: las tres rutas puestas una al lado
  de otra, con el precio de cada una dicho — y las dos correcciones propias
  anotadas.

Estructura cerrada no es hipótesis demostrada.

**Todavía no.**
