# Acta del Primo Construido — fabricado con dos números que no son enteros

**Orden del capitán · 2026-09-15 · F398**

> «Construime un número primo a partir de múltiplos que no sean enteros. Hacé
> el intento, aunque digas que no se pueda. Necesitamos experimentar.»

Se reproduce con `go run ./cmd/elprimoconstruido`.

**Se pudo.** Y en las dos direcciones: fabricando el primo desde los factores, y
partiendo primos ya dados.

---

## §1 · Hacia adelante — el primo no se busca, se FABRICA

Se eligen **primero** dos números irracionales, se multiplican, y recién
después se pregunta qué salió:

| factor 1 | factor 2 | producto | ¿primo? |
|---:|---:|---:|---|
| 3,414213562373… | 0,585786437627… | **2** | SÍ |
| 4,414213562373… | 1,585786437627… | **7** | SÍ |
| 6,414213562373… | 3,585786437627… | **23** | SÍ |
| 8,414213562373… | 5,585786437627… | **47** | SÍ |
| 10,414213562373… | 7,585786437627… | **79** | SÍ |
| 7,828427124746… | 2,171572875254… | **17** | SÍ |
| 9,828427124746… | 4,171572875254… | **41** | SÍ |

**Ninguno de los dos factores es entero.** Son irracionales, con infinitos
decimales que no se repiten nunca. Su producto da un entero exacto, y ese
entero es primo. Textual:

> **7 = 3,414213562373… × 1,585786437627…**
> **23 = (5+√2)(5−√2)**

## §2 · Hacia atrás — partir un primo dado

| p | anillo | los dos factores | producto |
|---:|---|---|---:|
| 7 | ℤ[√2] | 4,414213562 × 1,585786438 | 7,000000000 |
| 11 | ℤ[√3] | 4,464101615 × 2,464101615 | 11,000000000 |
| 17 | ℤ[√2] | 7,828427125 × 2,171572875 | 17,000000000 |
| 19 | ℤ[√5] | 5,472135955 × 3,472135955 | 19,000000000 |
| 29 | ℤ[√5] | 11,472135955 × 2,527864045 | 29,000000000 |
| 41 | ℤ[√2] | 9,828427125 × 4,171572875 | 41,000000000 |

## §3 · No es capricho: una congruencia decide qué primo alcanza cada anillo

| anillo | regla | alcanzados hasta 10.000 | fracción | violaciones |
|---|---|---:|---:|---:|
| ℤ[√2] | p ≡ ±1 (mod 8) | 603 / 1.228 | 0,4910 | **0** |
| ℤ[√3] | p ≡ ±1 (mod 12) | 608 / 1.228 | 0,4951 | **0** |
| ℤ[√5] | p ≡ ±1 (mod 5) | 610 / 1.228 | 0,4967 | **0** |
| ℤ[i] | p ≡ 1 (mod 4) | 609 / 1.228 | 0,4959 | **0** |
| ℤ[√−2] | p ≡ 1,3 (mod 8) | 606 / 1.228 | 0,4935 | **0** |

**Cada anillo alcanza casi exactamente la mitad**, y cuál alcanza lo dice el
resto de una división. Reciprocidad cuadrática (Gauss, 1832), verificada sobre
1.228 primos.

**Corrección propia, anotada:** la primera versión de la regla daba **una
violación** en ℤ[√3] (el 3) y otra en ℤ[√5] (el 5). No era un fallo de la
matemática sino de mi regla: esos primos **ramifican** — dividen al
discriminante — y por eso son alcanzables sin cumplir la congruencia de
partición. Agregado el caso, quedaron cero. También se corrigió el signo: cuando
la norma sale negativa, la factorización correcta es (b√d+a)(b√d−a), no
(a+b√d)(a−b√d); antes se mostraban productos negativos.

## §4 · ¿Cuántos anillos hacen falta para alcanzarlos a todos?

| anillos usados | primos NO alcanzados hasta 10.000 |
|---|---:|
| ℤ[√2] | 625 |
| + ℤ[√3] | 312 |
| + ℤ[√5] | 160 |
| + ℤ[i] | 81 |
| + ℤ[√−2] | **0** |

Ningún anillo solo alcanza a todos. **Con cinco, no queda ninguno afuera** en
ese tramo.

## §5 · La ruta que alcanza a todos, y por qué no vale — TRAMPA DECIMOSÉPTIMA

p = √p × √p funciona para **todo** primo (la raíz de un primo es irracional,
demostrado desde los griegos). Pero **por eso mismo no dice nada**: también
factoriza al 4, al 6 y al 9 exactamente igual. **Factoriza a todos por igual,
primos y compuestos.** Se reporta porque el capitán pidió el intento completo,
con el cartel puesto.

## §6 · Lo que se ganó y lo que no

**Se ganó, y es literal lo que se pidió:**
- un primo fabricado eligiendo **primero** dos irracionales: (5+√2)(5−√2) = 23
- primos partidos en dos factores no enteros, verificados a doce decimales
- la ley exacta de qué primo se deja partir en qué anillo, con **cero
  violaciones** en 1.228 primos

**No se ganó, y hay que decirlo con todas las letras:**
- **el primo sigue siendo primo entre los enteros.** Nada de esto lo vuelve
  compuesto ahí.
- dentro de cada anillo, el número que se parte deja de ser primo — **pero los
  factores nuevos son primos ELLOS**. El problema no desaparece: **se muda de
  casa.**
- ningún anillo alcanza a todos los primos; siempre hay que elegir cuál.
- **nada de esto es nuevo:** Gauss, 1832, y la reciprocidad cuadrática.

## Clasificación

- **Matemática (demostrado, y de otros):** la irracionalidad de √p; la ley de
  partición en cuerpos cuadráticos y la reciprocidad cuadrática (Gauss 1832);
  la ramificación de los primos que dividen al discriminante.
- **Medición (nuestra):** las construcciones verificadas a doce decimales; las
  fracciones 0,4910 a 0,4967 por anillo con cero violaciones sobre 1.228
  primos; la cobertura acumulada 625 → 312 → 160 → 81 → 0.
- **Nuevo (si algo):** ningún teorema. Lo propio es la **fabricación mostrada
  hacia adelante** — elegir los irracionales primero y que el primo salga
  después — y las dos correcciones propias anotadas.

Estructura cerrada no es hipótesis demostrada.

**Todavía no.**
