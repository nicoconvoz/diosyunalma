# Guía de validación

Esta guía es para revisar el proyecto desde cero, sin confiar en nada de lo que
afirma el registro. Todo hallazgo se regenera desde el código fuente; ningún
resultado depende de datos externos descargados.

## 1. Requisitos

- **Go 1.22 o superior** — https://go.dev/dl/ (instalador estándar, sin
  configuración extra).
- Una terminal cualquiera. En Windows: PowerShell.
- Para los experimentos pesados: ~2 GB de RAM libres. Los livianos corren en
  cualquier máquina.

Verificación de la instalación:

```bash
go version
```

## 2. Correr la batería de tests

Desde la raíz del repositorio:

```bash
go test ./...
```

Resultado esperado: **todos los paquetes en `ok`, 250 casos, cero fallos.**
Los tests cubren las bibliotecas (`primes/`, `pattern/`, `control/`,
`information/`, `spectral/`, `riemann/`). Los comandos de `cmd/` son capas
finas sobre esas bibliotecas y se validan corriéndolos.

Todo el proyecto se desarrolló con TDD estricto: cada función de biblioteca
tiene primero su test en rojo y después su implementación. El historial de git
lo documenta commit por commit.

## 3. Qué revisar y en qué orden

El registro completo está en [FINDINGS.md](FINDINGS.md): 47 hallazgos, con las
hipótesis muertas incluidas y las correcciones a la vista. Sugerencia de
recorrido para una revisión independiente:

### Nivel 1 — la ley de paridad (minutos)

```bash
go run ./cmd/lab
go run ./cmd/residue
```

Qué verificar: las ventanas palindrómicas de gaps contienen una cantidad impar
(o cero) de no-múltiplos de 3, con significancia mayor a 5σ contra señuelos que
preservan el multiconjunto exacto de gaps. La constante 0.8198... no está
ajustada: es el producto de Euler Π(q−3)(q−1)/(q−2)² sobre primos q≥5
(`cmd/decompose` la deriva).

### Nivel 2 — los ceros medidos (minutos a horas según límite)

```bash
go run ./cmd/zeta        # diez ceros de zeta desde los primos
go run ./cmd/radio3      # las estaciones de las tribus mod 3, 4, 5
go run ./cmd/symphony    # mod 7, mod 8, y el dial complejo
go run ./cmd/encore -mods 11,13
go run ./cmd/baton       # cosecha profunda mod 5 + el director
```

Qué verificar: los picos del periodograma deben coincidir con los ceros
publicados de cada función L. Fuentes externas de contraste:

| dial | fuente externa |
|------|----------------|
| ζ | https://www.lmfdb.org/zeros/zeta/ (γ₁ = 14.134725...) |
| χ mod 3 | https://www.lmfdb.org/L/1/3/3.2/r1/0/0 |
| χ mod 5 | https://www.lmfdb.org/L/1/5/5.4/r0/0/0 |
| χ mod 7 | https://www.lmfdb.org/L/1/7/7.6/r1/0/0 |
| χ mod 11 | https://www.lmfdb.org/L/1/11/11.10/r1/0/0 |
| χ mod 13 | https://www.lmfdb.org/L/1/13/13.12/r0/0/0 |

Tolerancia esperada: 0.005–0.02 en los límites por defecto. Los dials mod 8 y
el complejo mod 5 quedaron como predicciones del laboratorio (no se encontró
la etiqueta LMFDB correspondiente); son los que más valor tiene contrastar por
otra vía.

### Nivel 3 — la reconstrucción inversa (minutos)

```bash
go run ./cmd/sundial
```

Qué verificar: usando solo los ceros medidos en el nivel 2, la fórmula
explícita reconstruye la posición de los primos (escalones de ψ). Es la prueba
de que los ceros medidos no son artefactos del instrumento: cierran el círculo
en ambas direcciones.

## 4. La metodología que hay que atacar

Para invalidar un hallazgo, estos son los puntos débiles que el propio
proyecto declara y que conviene intentar romper:

1. **Los controles.** Cada detección se compara contra señuelos de gaps
   barajados que preservan el multiconjunto exacto (`control/shuffle.go`).
   Pregunta a hacerse: ¿el señuelo es realmente equivalente en todo salvo el
   orden? Si se encuentra un señuelo más honesto que mate un hallazgo, el
   hallazgo muere.
2. **El umbral.** Nada entra al registro con |z| < 5. Verificar que los
   cómputos de z usen la desviación estándar de los señuelos, no la teórica.
3. **La pre-registración.** Los resultados esperados están escritos en los
   comentarios del código fuente *antes* de mirar los datos. Revisar el
   historial de git para confirmar el orden temporal.
4. **Las correcciones.** Los hallazgos 2, 9 y 17 contienen errores corregidos
   que quedaron documentados dentro del propio registro. Verificar que las
   correcciones sean consistentes con lo que las reemplaza.
5. **Las hipótesis muertas.** Las hipótesis falsificadas siguen en el
   registro y siguen siendo reproducibles — las de la primera era y las de
   la campaña de la hipótesis de Riemann, incluidos los cinco casos en que
   un `0.0e+00` perfecto resultó venir de la construcción y no de un
   descubrimiento. Un registro que solo puede reproducir a sus
   sobrevivientes es un folleto de ventas.

   *(Acá decía «siete»: era el conteo de la primera era y quedó viejo. No se
   reemplaza por otro número hasta auditarlas una por una — inventar la
   cifra sería exactamente lo que este registro no hace.)*

## 5. Qué es conocido y qué es propio

Importante para calibrar expectativas: la existencia de los ceros de las
funciones L, la fórmula explícita y la ley de Chebotarev/Dirichlet son
matemática establecida (siglos XIX–XX). Lo que este laboratorio aporta no es
teoría nueva sino:

- la **ley de paridad de ventanas palindrómicas** con sus constantes derivadas
  (no encontrada en la literatura tras búsqueda; pendiente de revisión
  experta),
- la **medición casera** de ~60 ceros de 7–9 funciones L desde una sola criba,
  con verificación externa cero a cero,
- el **método**: controles de multiconjunto, pre-registración en el código,
  correcciones a la vista.

La sección "Answered along the way" de FINDINGS.md separa explícitamente qué
resultado reproduce literatura conocida y qué no tiene referencia encontrada.

## 6. Contacto con el registro

- [FINDINGS.md](FINDINGS.md) — los 47 hallazgos, en orden, con sus números.
- [README.md](../README.md) — el resumen y la tabla de comandos.
- `git log` — la cronología completa, un experimento por commit.
