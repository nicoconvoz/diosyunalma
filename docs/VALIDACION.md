# Guía de validación

> **La regla del sello** (auditoría externa de Yui, 2026-08-14, adoptada
> como ley del laboratorio en F302): **«Estructura cerrada» ≠ «Hipótesis
> demostrada».** Las auditorías externas sellaron la estructura del teorema
> de ruptura del yunque (F293-F301); ninguna afirmación de este laboratorio
> constituye una demostración de la Hipótesis de Riemann, y la positividad
> desde el lado de los primos sigue abierta.


Esta guía es para revisar el proyecto desde cero, **sin confiar en nada de lo que
afirma el registro**. Todo hallazgo se regenera desde el código fuente; ningún
resultado depende de datos externos descargados.

Está escrita para alguien que viene a romper el trabajo, no a admirarlo. Los
puntos débiles están señalados a propósito, y los lugares donde el laboratorio se
equivocó están marcados con el número del hallazgo para que se puedan ir a leer.

> **Última verificación de esta guía: 2026-08-11.** Cada número de acá abajo se
> midió con un comando el mismo día, no se copió de una versión anterior. Si una
> cifra no coincide con lo que ve, la guía está vieja: créale al comando.

---

## 1. Requisitos

- **Go 1.22 o superior** — https://go.dev/dl/ (instalador estándar, sin
  configuración extra). No hace falta nada más: ni bibliotecas, ni Python, ni
  cuentas, ni descargas.
- Una terminal cualquiera. En Windows: PowerShell.
- Para los experimentos pesados: ~2 GB de RAM libres. Los livianos corren en
  cualquier máquina.

```bash
go version
```

---

## 2. La batería de tests

Desde la raíz del repositorio:

```bash
go test ./...
```

**Resultado esperado, medido el 2026-08-11:**

| qué | cuánto |
|-----|--------|
| paquetes con tests, todos en `ok` | **6** (`primes`, `pattern`, `control`, `information`, `spectral`, `riemann`) |
| funciones de test | **102** |
| pruebas contando subtests | **176** |
| fallos | **0** |

Para contarlo usted mismo:

```bash
go test ./... -v | grep -c "^--- PASS:"
```

> **Corrección que va acá porque corresponde.** Las versiones anteriores de esta
> guía y del README decían «250 casos». Ese número contaba entradas de tablas
> adentro de los subtests y hoy no es verificable con un comando. Se reemplaza
> por lo que sí se puede contar. Si alguien quiere el número de asserts, va a
> tener que instrumentarlo — y si lo hace, que lo mande.

Los comandos de `cmd/` **no tienen tests**: son capas finas sobre las bibliotecas
y se validan corriéndolos. Eso es una decisión, no un olvido, y es un punto
legítimo para atacar.

Todo el proyecto se desarrolló con TDD estricto: cada función de biblioteca tiene
primero su test en rojo y después su implementación. El historial de git lo
documenta commit por commit.

---

## 3. El tamaño real de lo que va a revisar

Medido el 2026-08-11:

| qué | cuánto |
|-----|--------|
| hallazgos numerados en `FINDINGS.md` | **302** |
| secciones, contando las sub-numeradas con letra (140b, 220f…) | **314** |
| experimentos ejecutables en `cmd/` | **241** |
| láminas en `galeria/laminas/` | **144** |
| sonidos en `galeria/sonidos/` | **8** |
| paradas del museo | **272** |
| máximas del capitán, verificadas una por una | **19** |

Para contarlos:

```bash
rg -o "^## Finding [0-9]+" docs/FINDINGS.md | rg -o "[0-9]+" | sort -n | uniq | wc -l
ls galeria/laminas/*/*.svg | wc -l
```

**El registro no tiene huecos ni duplicados.** Se auditó el 2026-08-10 y tenía
las dos cosas: `F258` estaba duplicado en los dos idiomas, faltaban `F262` y
`F263` en el registro en inglés y faltaba el `135` en el castellano. Está
reparado, y la reparación está registrada en `F271`. Verifíquelo:

```bash
rg -o "^## Finding [0-9]+" docs/FINDINGS.md | sort | uniq -c | grep -v " 1 "
```

Lo que salga con cuenta mayor a 1 son las sub-numeradas con letra (140 y 140b
comparten el número), no duplicados.

---

## 4. Recorrido de revisión, en orden

### Nivel 1 — la ley de paridad (minutos)

```bash
go run ./cmd/lab
go run ./cmd/residue
go run ./cmd/decompose
```

Qué verificar: las ventanas palindrómicas de gaps contienen una cantidad impar
(o cero) de no-múltiplos de 3, con significancia mayor a 5σ contra señuelos que
preservan el multiconjunto exacto de gaps. La constante 0.8198… **no está
ajustada**: es el producto de Euler Π(q−3)(q−1)/(q−2)² sobre primos q ≥ 5, y
`cmd/decompose` la deriva.

### Nivel 2 — los ceros medidos (minutos a horas según límite)

```bash
go run ./cmd/zeta        # diez ceros de zeta desde los primos
go run ./cmd/radio3      # las estaciones de las tribus mod 3, 4, 5
go run ./cmd/symphony    # mod 7, mod 8, y el dial complejo
go run ./cmd/encore -mods 11,13
go run ./cmd/baton       # cosecha profunda mod 5 + el director
```

Los picos del periodograma tienen que coincidir con los ceros publicados de cada
función L. Fuentes externas de contraste:

| dial | fuente externa |
|------|----------------|
| ζ | https://www.lmfdb.org/zeros/zeta/ (γ₁ = 14.134725…) |
| χ mod 3 | https://www.lmfdb.org/L/1/3/3.2/r1/0/0 |
| χ mod 5 | https://www.lmfdb.org/L/1/5/5.4/r0/0/0 |
| χ mod 7 | https://www.lmfdb.org/L/1/7/7.6/r1/0/0 |
| χ mod 11 | https://www.lmfdb.org/L/1/11/11.10/r1/0/0 |
| χ mod 13 | https://www.lmfdb.org/L/1/13/13.12/r0/0/0 |

Tolerancia esperada: 0.005–0.02 en los límites por defecto. Los dials mod 8 y el
complejo mod 5 quedaron como **predicciones del laboratorio** (no se encontró la
etiqueta LMFDB correspondiente): son los que más valor tiene contrastar por otra
vía.

### Nivel 3 — la reconstrucción inversa (minutos)

```bash
go run ./cmd/sundial
```

Usando **solo** los ceros medidos en el nivel 2, la fórmula explícita reconstruye
la posición de los primos. Es la prueba de que los ceros medidos no son
artefactos del instrumento: cierran el círculo en las dos direcciones.

### Nivel 4 — la campaña de Riemann, y su resultado NEGATIVO

Este es el nivel que más conviene revisar, porque es donde el laboratorio dice
que **fracasó**.

```bash
go run ./cmd/elensamble    # la cadena entera, juzgada eslabón por eslabón
go run ./cmd/davenport     # la función que mata la rama geométrica
```

`cmd/davenport` construye desde cero la función de **Davenport–Heilbronn**
(1936) y encuentra por barrido ciego uno de sus ceros fuera de la línea crítica.
Esa función **cumple todas las simetrías que este laboratorio demostró** y viola
igual la hipótesis. Es la evidencia más clara de que la ruta geométrica está
cerrada, y está en el registro con esas palabras.

Los tres teoremas propios que cierran esa ruta —la simetría sola nunca decide,
toda forma derivada del corte de ½ es ciega a dónde está el cero, y ningún
cómputo finito alcanza porque el horizonte de detección se escapa— están en
`FINDINGS.md` con sus números y sus programas.

### Nivel 5 — la campaña de la tabla del capitán (F264–F280)

La más reciente, y la más fácil de auditar porque son números chicos:

```bash
go run ./cmd/losgemelos      # primos gemelos: el centro es múltiplo de 6
go run ./cmd/lasuma          # Goldbach, con el centro suelto
go run ./cmd/elanterior      # la regla que un múltiplo de 3 mata en tres renglones
go run ./cmd/lacadena        # la melodía que no puede desafinar
go run ./cmd/elrelieve       # el triángulo de Gilbreath
go run ./cmd/losdoses        # el sesgo real, contra su control barajado
go run ./cmd/losciclos       # la hipótesis del compás, muerta por tres pruebas
go run ./cmd/elunoquesobra   # la rueda, y por qué servir para todos no alcanza
go run ./cmd/losdoscentros   # la tautología que era la puerta: Fermat, 1643
go run ./cmd/elcentrocompartido # compartir centro ES ser gemelos
go run ./cmd/elmedioqueune   # el 2 cae en el ½ exacto de la dimensión 0
go run ./cmd/lamelodiadelacriba # los primos hasta 100, armonizados
go run ./cmd/laformadelproblema # la forma del problema, dibujada
go run ./cmd/laley           # x−y=0 en la dimensión 0, y cuál es la puerta
go run ./cmd/eleuler         # RH ⟺ toda perla es un giro puro
```

**Seis de esos quince son kills**, y tres contienen una corrección al propio
laboratorio escrita adentro del programa. Si va a revisar una sola cosa de todo
el proyecto, que sea este nivel: es donde el criterio se ve funcionando en vivo.

Cada uno imprime sus propias leyes, su veredicto **derivado de la medición** y su
bloque de límites. Varios contienen correcciones al propio laboratorio escritas
adentro del programa.

---

## 5. La metodología que hay que atacar

Para invalidar un hallazgo, estos son los puntos débiles que el propio proyecto
declara y que conviene intentar romper:

1. **Los controles.** Cada detección se compara contra señuelos de gaps barajados
   que preservan el multiconjunto exacto (`control/shuffle.go`). Pregunta a
   hacerse: ¿el señuelo es realmente equivalente en todo salvo el orden? Si se
   encuentra un señuelo más honesto que mate un hallazgo, el hallazgo muere.

2. **El umbral.** Nada entra al registro con |z| < 5. Verificar que los cómputos
   de z usen la desviación estándar **de los señuelos**, no la teórica.

3. **El testigo correcto.** Y este es el que más caro salió. En `F269` el
   laboratorio comparó una caminata contra √n —una moneda— cuando el testigo
   válido era el barajado, porque la anti-persistencia ya venía metida en la
   construcción. La afirmación se retractó **en la misma lámina donde se
   descubrió**. Busque más casos: es el error más fácil de cometer acá.

4. **La pre-registración.** Los resultados esperados están escritos en los
   comentarios del código fuente *antes* de mirar los datos. Revisar el historial
   de git para confirmar el orden temporal.

5. **La trampa del `0.0e+00`.** **NUEVE veces** un resultado perfecto resultó venir
   de la construcción y no de los números — una identidad algebraica disfrazada
   de medición. Las nueve están en el registro. Cada vez que vea un cero exacto o
   un acuerdo perfecto, la pregunta correcta es: *¿esto podría haber fallado
   alguna vez?* Si la respuesta es no, no es un hallazgo.

6. **Los veredictos tipeados a mano.** `F265` cazó un defecto de método: la
   palabra «CERO fallos» estaba escrita a mano en la salida mientras la variable
   medida decía otra cosa. Es un defecto **cazable a máquina en los 280
   hallazgos**: comparar cada veredicto afirmado contra la variable que describe.
   Si encuentra otro, es un hallazgo suyo.

7. **Las correcciones.** Están escritas **adentro del hallazgo que revisan**,
   nunca editadas por encima. Verificar que cada corrección sea consistente con
   lo que reemplaza.

8. **Las hipótesis muertas.** Siguen en el registro y siguen siendo
   reproducibles. Un registro que solo puede reproducir a sus sobrevivientes es
   un folleto de ventas.

   *(Acá hubo un conteo de hipótesis muertas que quedó viejo y se sacó. No se
   reemplaza por otro número hasta auditarlas una por una — inventar la cifra
   sería exactamente lo que este registro no hace.)*

---

## 6. Dónde el laboratorio se mató solo

Estos son los lugares donde el propio proyecto tumbó algo suyo. Son los mejores
puntos de entrada para un revisor, porque muestran el criterio en acción:

- **`F259`** — el gran ensamble: dieciocho jueces refutaron por unanimidad el
  cierre de la cadena, y ahí apareció Davenport–Heilbronn.
- **`F265`** — la afirmación «su fórmula ES Goldbach» resultó una
  sobreafirmación: con X entero solo alcanza los múltiplos de 4. Corregido en el
  hallazgo, con el ejemplo estrella que lo desmentía nombrado.
- **`F267`** — la sexta trampa del `0.0e+00`, declarada en el mismo hallazgo que
  la contiene.
- **`F269`** — la retractación del testigo equivocado, ya mencionada.
- **`F270`** — una hipótesis del propio autor, muerta por tres pruebas
  independientes, con la trampa de las cazas de ciclos declarada antes de medir.
- **`F271`** — una fuga de un archivo privado adentro del instalador. `.gitignore`
  protege al repositorio, **no al instalador**: son dos listas distintas.

---

## 7. Qué es conocido y qué es propio

Importante para calibrar expectativas. **La mayor parte de lo que hay acá
reproduce matemática establecida**, y el registro lo dice en cada ficha.

Es matemática conocida, de los siglos XIX y XX: la existencia de los ceros de las
funciones L, la fórmula explícita, Chebotarev/Dirichlet, la conjetura de Goldbach
(1742), la de los primos gemelos (de Polignac, 1849), la de Gilbreath (1958) y el
sesgo de Lemke Oliver–Soundararajan (2016). Varios hallazgos recientes son
**reencuentros** con esos objetos, y están etiquetados como tales.

Lo que este laboratorio aporta no es teoría nueva sino:

- la **ley de paridad de ventanas palindrómicas** con sus constantes derivadas
  (no encontrada en la literatura tras búsqueda; **pendiente de revisión
  experta**),
- la **medición casera** de ~70 ceros de 7–9 funciones L desde una sola criba,
  con verificación externa cero a cero,
- **tres resultados negativos** sobre por qué la ruta geométrica no puede cerrar,
  con su contraejemplo reproducido desde cero,
- el **método**: controles de multiconjunto, pre-registración en el código,
  correcciones adentro del hallazgo que revisan, y paneles adversariales que
  tienen permiso explícito para matar el trabajo del propio laboratorio.

**Lo que NO hay acá: una demostración de la Hipótesis de Riemann.** El registro
lo dice en cada cierre con dos palabras: *todavía no*.

---

## 8. Las salas para el público, y cómo verificarlas

Tres piezas de divulgación, y las tres se regeneran desde el código:

```bash
go run ./cmd/puente                      # el puente de mando, localhost:8118
go run ./cmd/puente -museo galeria/museo.html   # regenera museo Y máximas
```

- **`galeria/index.html`** — las 144 láminas. Cada una la dibuja su propio
  experimento: corra el comando y compare el SVG.
- **`galeria/museo.html`** — 272 paradas en criollo. Cada pieza cierra con su
  bloque de límites honestos. Se genera desde `cmd/puente/museo.go`, así que el
  museo y el registro no se pueden desincronizar.
- **`galeria/maximas.html`** — 19 frases del autor. **Cada una lleva su archivo y
  su línea exacta.** Verifíquelas: si una no está donde dice, es un hallazgo suyo.
  Una de las diecinueve está publicada como rechazada a propósito, con la
  explicación de por qué no era su voz.

---

## 9. Contacto con el registro

- [FINDINGS.md](FINDINGS.md) — los 280 hallazgos numerados, en inglés, con sus
  números, sus comandos y sus hipótesis muertas.
- [HALLAZGOS-ES.md](HALLAZGOS-ES.md) — el mismo registro en castellano.
- [BITACORA-NOCTURNA.md](BITACORA-NOCTURNA.md) — el diario de a bordo, día por
  día, donde están los flashes en crudo antes de ser formalizados.
- [INFORME-TECNICO.md](INFORME-TECNICO.md) — el informe largo.
- [README.md](../README.md) — el resumen y la tabla de comandos.
- `git log` — la cronología completa, un experimento por commit.
- **DOI:** [10.5281/zenodo.21864277](https://doi.org/10.5281/zenodo.21864277)
- **Licencias:** código AGPL-3.0, contenido CC BY 4.0, más una licencia comercial
  disponible. Ver [LICENCIAS.md](../LICENCIAS.md).

---

*Si encuentra un error, es un hallazgo. Mándelo con el comando que lo reproduce y
entra al registro con su nombre.*
