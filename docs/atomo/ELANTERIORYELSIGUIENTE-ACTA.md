# Acta del Anterior y el Siguiente — cada divisor veta un hueco

**Flash del capitán · 2026-09-15 · F395**

> «¿Hay alguna relación entre la composición de los divisores de un número
> anterior a un número primo y el número primo siguiente?»

Se reproduce con `go run ./cmd/elanterioryelsiguiente`.

**Sí, y es exacta.** Cada divisor primo impar de p−1 **prohíbe** un hueco
determinado hasta el primo siguiente. Se deriva en un renglón, se verifica en
664.578 huecos sin una sola excepción, y después se mide cuánto pesa de verdad
— que es la parte donde hay que ser honesto.

---

## §1 · La trampa, declarada antes de contar — DECIMOCUARTA

«p−1 es par» vale para todo primo impar **y para todo compuesto impar**: es la
definición de impar, ya cazada en F272 y F393. El 2 queda **excluido** de todas
las cuentas; solo se miran los divisores **impares** de p−1, que son los que
llevan información.

## §2 · La regla, en un renglón

    si q | (p−1)  ⟹  p ≡ 1 (mod q)  ⟹  p + (q−1) ≡ 0 (mod q)
    y como p + q − 1 > q, ese número es COMPUESTO
    ⟹ el hueco hasta el primo siguiente NO PUEDE valer q − 1

**El 3 veta el hueco 2 · el 5 veta el 4 · el 7 veta el 6 · el 11 veta el 10.**

| p | p−1 | divisores impares | huecos vetados | siguiente | hueco |
|---:|---:|---|---|---:|---:|
| 7 | 6 | 3 | 2 | 11 | 4 |
| 11 | 10 | 5 | 4 | 13 | 2 |
| 17 | 16 | — | ninguno | 19 | 2 |
| 29 | 28 | 7 | 6 | 31 | 2 |
| 31 | 30 | 3, 5 | 2, 4 | 37 | 6 |

## §3 · Verificación (un TEOREMA verificándose, no un hallazgo)

Sobre los 664.578 huecos hasta 10⁷:

| q | hueco vetado | veces con q \| p−1 | veces con q ∤ p−1 |
|---:|---:|---:|---:|
| 3 | 2 | **0** | 58.979 |
| 5 | 4 | **0** | 58.621 |
| 7 | 6 | **0** | 99.987 |
| 11 | 10 | **0** | 54.431 |
| 13 | 12 | **0** | 65.513 |

**Cero violaciones**, como debe ser: esto está demostrado, no medido. Si
imprimiera otra cosa, el error estaría en el código.

## §4 · Lo que sí se mide: cuánto corre el veto la distribución

| q | hueco medio con q \| p−1 | con q ∤ p−1 | corrimiento |
|---:|---:|---:|---:|
| 3 | 15,5205 | 14,5741 | **+0,9464 (+6,49 %)** |
| 5 | 14,7636 | 15,1417 | −0,3780 (−2,50 %) |
| 7 | 15,9445 | 14,8680 | **+1,0766 (+7,24 %)** |
| 11 | 14,9832 | 15,0543 | −0,0711 (−0,47 %) |
| 13 | 15,4695 | 15,0089 | +0,4607 (+3,07 %) |

Y la escalera, contando cuántos primos chicos dividen a p−1:

| divisores impares de p−1 | cuántos primos | hueco medio |
|---:|---:|---:|
| 0 | 135.148 | 14,9715 |
| 1 | 267.272 | 14,7793 |
| 2 | 184.507 | 15,0708 |
| 3 | 63.513 | 15,7942 |
| 4 | 12.536 | **17,0108** |
| 5 | 1.462 | **18,2079** |
| 6 | 128 | **22,0469** |

**Cuantos más primos chicos dividen a p−1, más lejos queda el siguiente** —
porque cada divisor tacha un candidato cercano.

**Irregularidad anotada, no maquillada:** la escalera **no es monótona al
principio** — con 1 divisor el hueco medio (14,78) es *menor* que con 0 (14,97).
De 2 en adelante sube limpio. La clase de 0 divisores exige que p−1 = 2^a por
primos todos mayores a 31, que es una condición fuerte y selecciona p atípicos;
no se ofrece explicación medida para esa inversión, se deja señalada.

## §5 · El control del espejo, y su límite

Por el mismo renglón, q | (p+1) veta el hueco q+1:

| q | veta hueco | hueco medio con q \| p+1 | con q ∤ p+1 | corrimiento |
|---:|---:|---:|---:|---:|
| 3 | 4 | 14,5741 | 15,5205 | −0,9464 |
| 5 | 6 | 15,4649 | 14,9081 | +0,5569 |
| 7 | 8 | 14,3623 | 15,1842 | −0,8219 |

**Advertencia de lectura que hay que hacer:** para **q = 3** este control **no
es independiente**. Todo primo p > 3 cumple p ≡ 1 o p ≡ 2 (mod 3), así que
3 | p−1 y 3 | p+1 son **complementarios por necesidad** — por eso los números
14,5741 y 15,5205 aparecen intercambiados en las dos tablas: es la misma
medición dada vuelta, no una segunda. Para q ≥ 5 sí son clases distintas (p
puede no caer en ninguna), y ahí el control vale.

Lo que el espejo sí muestra: **el veto no es propiedad del «anterior»** — es de
la aritmética modular, y el anterior y el siguiente vetan huecos **distintos**.

## §6 · ¿Y cuánto predice, realmente?

| | |
|---|---:|
| hueco medio global | 15,0473 |
| desvío típico del hueco | 12,5697 |
| **varianza explicada por los divisores impares de p−1** | **0,0012 = 0,12 %** |

**Ésta es la parte que hay que decir fuerte.** El efecto es **real y
demostrable**, pero **chico**: explica el 0,12 % de la variación. Saber los
divisores de p−1 te deja **tachar candidatos**, no te dice cuál es el siguiente.
**Vetar no es predecir.**

## Clasificación

- **Matemática (demostrado, y elemental):** q | (p−1) ⟹ el hueco ≠ q−1;
  q | (p+1) ⟹ el hueco ≠ q+1. Es el mismo argumento con el que Hardy y
  Littlewood decidían qué constelaciones de primos son **admisibles**. **No es
  matemática nueva**: la casa lo mide, no lo descubre.
- **Medición (nuestra):** cero violaciones en 664.578 huecos; los corrimientos
  de la media por divisor; la escalera de 14,97 a 22,05 con su irregularidad
  declarada; el 0,12 % de varianza explicada.
- **Nuevo (si algo):** ningún teorema. La **cuantificación** de un efecto que
  se suele enunciar y no medir, con la honestidad de decir que pesa 0,12 %.

Estructura cerrada no es hipótesis demostrada.

**Todavía no.**
