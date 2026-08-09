# Licencias del laboratorio

Este repositorio lleva **dos licencias**, una para el código y otra para el
contenido. No es un capricho: son cosas distintas y se usan distinto.

| Qué | Licencia | Archivo | En criollo |
|---|---|---|---|
| **El código** — todo `cmd/` y cualquier fuente Go | **MIT** | [LICENSE](LICENSE) | hacé lo que quieras con él, incluso venderlo; solo dejá el aviso de autoría |
| **El contenido** — `galeria/` (láminas y sonidos), `docs/` (bitácora, hallazgos, informe, museo) y los textos | **CC BY 4.0** | [LICENSE-CONTENIDO.txt](LICENSE-CONTENIDO.txt) | usalo, copialo, modificalo, hasta comercialmente — **pero citá de dónde salió** |

---

## Cómo citarme

Si usás una lámina, un texto o un hallazgo, alcanza con esto:

```
Jesús Nicolás Astorga y RESOURCES OPEN DOORS S.A.S, "Laboratorio Diosyunalma", 2026.
https://github.com/nicoconvoz/numerosprimos
Licencia CC BY 4.0 (https://creativecommons.org/licenses/by/4.0/)
```

Si usás el código, no hace falta citar en la documentación: alcanza con dejar
el archivo `LICENSE` donde esté, como pide MIT.

---

## Por qué dos y no una

- **MIT en el código** es la licencia más permisiva que existe entre las
  serias. Que alguien agarre el motor doble-doble, o el balde de luz, o el
  tren, y lo use en su propio trabajo es **exactamente lo que queremos** —
  ése es medio el punto de publicar.
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
- **Repositorio:** https://github.com/nicoconvoz/numerosprimos
- **Año:** 2026 — si el repo se publica más adelante, actualizar en `LICENSE`
  (línea 3) y en `LICENSE-CONTENIDO.txt` (línea 3)

El módulo de Go también dice `github.com/nicoconvoz/numerosprimos`, que es lo
que hace falta para que `go get` funcione: si el módulo y la URL del repo no
coinciden, Go se niega a bajarlo.

GitHub va a detectar la MIT sola y te va a poner la etiqueta en la portada del
repositorio. La CC BY va a quedar como archivo, que es lo normal para
licencias de contenido.
