# Licencias del laboratorio

Este repositorio lleva **dos licencias**, una para el código y otra para el
contenido. No es un capricho: son cosas distintas y se usan distinto.

| Qué | Licencia | Archivo | En criollo |
|---|---|---|---|
| **El código** — todo `cmd/` y cualquier fuente Go | **AGPL-3.0** | [LICENSE](LICENSE) | libre de usar, estudiar y modificar — pero si lo metés en tu producto o lo servís por red, tenés que abrir tu código |
| **El código, para uso cerrado** | **Licencia comercial** | [LICENCIA-COMERCIAL.md](LICENCIA-COMERCIAL.md) | te exime de la AGPL a cambio de un acuerdo |
| **El contenido** — `galeria/` (láminas y sonidos), `docs/` (bitácora, hallazgos, informe, museo) y los textos | **CC BY 4.0** | [LICENSE-CONTENIDO.txt](LICENSE-CONTENIDO.txt) | usalo, copialo, modificalo, hasta comercialmente — **pero citá de dónde salió** |

---

## Cómo citarme

Si usás una lámina, un texto o un hallazgo, alcanza con esto:

```
Jesús Nicolás Astorga y RESOURCES OPEN DOORS S.A.S, "Laboratorio Diosyunalma", 2026.
https://github.com/nicoconvoz/diosyunalma
Licencia CC BY 4.0 (https://creativecommons.org/licenses/by/4.0/)
```

Si usás el código, tenés que cumplir la AGPL: dejar el aviso de copyright y la
licencia, y **publicar el fuente de lo que construyas encima** — también si lo
servís por red (artículo 13). Si eso no te sirve, hay
[licencia comercial](LICENCIA-COMERCIAL.md).

---

## Por qué dos y no una

- **AGPL en el código** para que el trabajo circule sin que nadie lo cierre: el
  que lo use y mejore, devuelve. Y el artículo 13 cubre el caso moderno — no
  alcanza con no distribuir el binario: si lo servís por red, también abrís.
  Sobre esa base se apoya la **licencia comercial**: quien no quiera abrir su
  código, paga. Ese es el modelo con el que viven MySQL, Qt y Ghostscript.
- **CC BY 4.0 en el contenido** porque las láminas y los textos no son
  software: son obra. La MIT no está pensada para eso. Y CC BY pide lo único
  que importa acá: **que digan de dónde salió**.

---

## Lo que estas licencias NO dicen

Ninguna de las dos afirma nada sobre matemática. Vale la frase que va también
en el README:

> Este laboratorio publica **mediciones, instrumentos y visualizaciones de
> matemática conocida**, más algunos instrumentos originales. **No se reclama
> aquí ninguna demostración de la Hipótesis de Riemann.** Cada resultado viene
> con su límite declarado, y los errores del propio laboratorio están
> publicados junto con los aciertos.

---

## Dependencias de terceros

**Ninguna.** El `go.mod` no tiene una sola línea `require`: todo está
construido sobre la biblioteca estándar de Go. Por eso no hace falta ningún
archivo de atribuciones de terceros ni auditoría de licencias aguas arriba.

La única excepción es el propio texto legal de CC BY 4.0 incluido en
`LICENSE-CONTENIDO.txt`, que se transcribe tal cual lo publica Creative
Commons — como corresponde, y como ellos mismos piden.

---

## Titularidad

La obra tiene **dos titulares de copyright**, por decisión expresa del autor:

- **Jesús Nicolás Astorga** — autor
- **RESOURCES OPEN DOORS S.A.S** — organización de ingeniería que financia el trabajo
  (en pantalla y en los créditos figura con su marca, *Open Doors*)

Qué significa co-titularidad, dicho sin vueltas: **los dos son dueños de la
obra en igualdad**. Cualquiera de los dos puede licenciarla, relicenciarla o
disponer de ella. No es un crédito honorífico: es propiedad compartida, y
figura así en las dos licencias.

Se hizo por decisión del capitán, sabiendo el alcance. Vale dejar escrito que
**financiar no otorga copyright por sí solo**: sin un contrato que lo diga, la
obra habría quedado 100% del autor. Esto es un acto voluntario suyo.

## Datos ya confirmados por el capitán

- **Autor:** Jesús Nicolás Astorga
- **Co-titular:** RESOURCES OPEN DOORS S.A.S
- **Repositorio:** https://github.com/nicoconvoz/diosyunalma
- **Año:** 2026 — si el repo se publica más adelante, actualizar en `LICENSE`
  (línea 3) y en `LICENSE-CONTENIDO.txt` (línea 3)

El módulo de Go también dice `github.com/nicoconvoz/diosyunalma`, que es lo
que hace falta para que `go get` funcione: si el módulo y la URL del repo no
coinciden, Go se niega a bajarlo.

GitHub va a detectar la AGPL-3.0 y te va a poner la etiqueta en la portada del
repositorio. La CC BY y la comercial quedan como archivos, que es lo normal.

## El artículo 13, y por qué nos toca de lleno

El puente de mando **es un servidor HTTP**. Bajo AGPL eso no es un detalle: el
artículo 13 obliga a que cualquiera que use el programa POR RED pueda obtener
su fuente. Por eso el puente sirve, y tiene que seguir sirviendo:

- **`/fuente`** — la oferta, en pantalla
- **`/fuente.zip`** — el fuente correspondiente de verdad: todo `cmd/`, el
  `go.mod` y las licencias, armado en vivo

Si alguien alguna vez saca ese enlace del pie de la página, el laboratorio
queda incumpliendo su propia licencia. Está anotado acá para que no pase.
