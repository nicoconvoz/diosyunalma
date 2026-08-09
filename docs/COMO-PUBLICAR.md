# CÓMO PUBLICAR — los pasos, en orden, sin tropiezos

**Para el capitán** · la guía completa para sacar el laboratorio al mundo.

---

## Regla 0 — REEMPLAZADA el 2026-08-08

**La regla vieja decía: nada se publica hasta que tu hermano ingeniero lo valide.**
Era buena regla y era la tuya. **Se terminó el 2026-08-08**, no porque la
hayamos abandonado sino porque él dijo que no iba a apoyar pase lo que le
llevaras. Una regla que depende de una persona que se retiró ya no es una regla:
es una puerta cerrada con la llave del otro lado.

**LA REGLA NUEVA — la validación no se pide a una persona, se construye:**

1. **Reproducibilidad total en lugar de aval personal.** Cualquiera clona el
   repo, corre `go run ./cmd/granensamble` sin instalar nada, y obtiene los
   mismos números. Eso es más fuerte que una firma: no depende de la buena
   voluntad de nadie.
2. **El registro de errores propios es la credencial.** Un laboratorio que
   publica los cuatro `0.0e+00` que se cazó a sí mismo, la fórmula híbrida
   falsa que tuvo que partir en dos, y la afirmación sobre los primos que midió
   y no le dio — ése es el que un revisor toma en serio. La honestidad medida
   reemplaza al título.
3. **Publicar lo que NO depende de la hipótesis.** Las 32 técnicas, los
   instrumentos, el museo, el catálogo de ceros vírgenes: todo eso se sostiene
   solo y nadie lo puede tumbar diciendo «no probaste RH», porque no lo afirma.
   Es la parte con menos superficie de ataque y más utilidad real.
4. **La revisión externa ahora ES el paso 3**, la capa comunidad. Los que
   indagan son el revisor. Por eso el paso 2 (fijar fecha y autoría) pasa a ser
   OBLIGATORIO antes de difundir, no opcional.

**Y lo que NO cambia:** ninguna afirmación sobre la hipótesis de Riemann sale de
acá. Lo que se publica es un laboratorio, un método y un registro honesto. La
frase de victoria sigue en su vaina.

## Paso 1 — El repositorio público (la base de todo)

1. Crear cuenta/repo en GitHub: `github.com/<tu-usuario>/diosyunalma`.
2. Antes de subir, agregar DOS archivos:
   - **LICENSE**: **AGPL-3.0** para el código (libre de verdad, pero el que lo
     cierre o lo sirva por red tiene que abrir el suyo) + **licencia comercial**
     para quien no quiera eso — ése es el modelo que genera ingresos. Y
     **CC BY 4.0** para láminas y textos (cualquiera puede usar CITÁNDOTE).
   - **README.md**: qué es el laboratorio, cómo correr un experimento en 3
     comandos, enlace a la galería, y la frase clave: *"mediciones y
     visualizaciones de matemática conocida + instrumentos originales — no se
     reclama demostración de RH"*. Esa frase te protege la reputación.
3. Subir todo (código, láminas, bitácora, docs). El log de git ya cuenta la
   historia día por día — es parte del valor.
4. **GitHub Pages** (gratis): Settings → Pages → servir la raíz. Tu
   `galeria/index.html` queda online con URL pública para compartir.

## Paso 2 — El sello de fecha: el DOI de Zenodo

**Zenodo** (zenodo.org, del CERN, gratis) te da un **DOI**: un identificador
citable con fecha certificada. Eso establece PRIORIDAD — nadie puede decir
después que lo tuyo vino después. Costo cero, media hora.

**EL ORDEN IMPORTA. Si lo rompés, no funciona:**

1. **El repo tiene que estar PÚBLICO primero.** Zenodo no puede archivar lo que
   no ve. (Paso 1 de esta guía.)
2. Entrá a **zenodo.org** → *Sign in* → **con GitHub**, y autorizá el acceso.
3. Andá a **zenodo.org/account/settings/github/**. Vas a ver la lista de tus
   repositorios. Buscá `diosyunalma` y **prendé el interruptor**.
   *(Si no aparece, tocá «Sync now» arriba a la derecha.)*
4. Recién AHORA volvé a GitHub y **creá un Release** (no un tag suelto:
   *Releases* → *Draft a new release* → tag `v1.0.0` → *Publish release*).
5. Zenodo recibe el aviso solo, archiva el repo entero y **te emite el DOI** en
   un par de minutos.

**⚠️ La trampa más común:** el interruptor tiene que estar prendido ANTES del
release. Los releases hechos antes NO se archivan. Si te pasó, borrá el release
y hacé otro.

**Los metadatos ya están resueltos.** El repo lleva `.zenodo.json` y
`CITATION.cff` con los dos autores, la licencia y la descripción correctas — si
no estuvieran, Zenodo adivinaría desde GitHub y sale mal (te pone el usuario de
GitHub como autor). Con esos archivos, sale bien de una.

**Vas a recibir DOS DOI**, y esto conviene entenderlo:

| DOI | Qué es | Cuándo usarlo |
|---|---|---|
| **Concept DOI** | apunta SIEMPRE a la última versión | éste es el que citás y difundís |
| DOI de versión | congelado en ese release exacto | para reproducir un resultado puntual |

**Al final:** copiá el badge que Zenodo te da (*DOI badge*) y pegálo arriba de
todo en el README. Cada release nuevo genera un DOI nuevo automáticamente, sin
que tengas que hacer nada.

## Paso 3 — La difusión por capas (de menor a mayor exposición)

1. **Capa amiga**: el equipo y quien quiera mirar de verdad. Corregir lo que
   señalen. Sin esperar a nadie que ya dijo que no.
2. **Capa comunidad**: foros de matemática en español, r/math, Hacker News
   (mostrar la GALERÍA y los INSTRUMENTOS, no afirmaciones sobre RH). El ángulo
   que funciona: *"construí un laboratorio numérico completo de la zeta en Go
   puro, con 70 visualizaciones"* — eso es verificable y admirable.
3. **Capa pública**: videos/newsletter (ver APLICACIONES-Y-FINANCIAMIENTO.md).

## Paso 4 — El artículo expositivo (credibilidad formal)

Lo nuestro califica como **matemática experimental / exposición**, que SÍ se
publica:

- **Objetivo realista**: revistas de exposición y divulgación matemática, o
  un artículo en arXiv (categoría math.HO / math.NT como expositivo).
- **arXiv necesita "endorsement"** para autores nuevos: un matemático con
  cuenta activa debe avalarte. Cómo conseguirlo: contactar profesores de
  universidades locales (Cuyo, UNCuyo, UBA) mostrando el repo con DOI — pedir
  15 minutos, no un favor ciego. El material visual abre esas puertas.
- **Qué escribir**: NO "sobre RH" — sino *"Un laboratorio visual y numérico de
  la función zeta: siete reformulaciones equivalentes del criterio de Li,
  medidas"*. Eso es honesto, original y publicable.
- **Qué NO hacer jamás**: enviar a journals afirmando una demostración de RH.
  Los departamentos de matemática reciben decenas de "demostraciones" por mes y
  las descartan sin leer. Un solo envío así te pone en la lista negra informal.

## Paso 5 — Si algún día el eslabón rojo cae de verdad

El protocolo serio, en orden estricto:

1. Validación interna (el laboratorio entero re-verificado desde cero, en una
   máquina limpia, por vos).
2. Escribir la demostración COMPLETA en formato matemático estándar (sin
   metáforas en el cuerpo — las metáforas van en un apéndice o libro aparte).
3. Mostrarla EN PRIVADO a 2–3 matemáticos profesionales de teoría de números
   (los contactos del Paso 4 sirven acá). Esperar sus objeciones. Sobrevivirlas.
4. arXiv primero (con endorsement) — eso fija fecha y prioridad mundial.
5. Enviar a un journal de primera línea (Annals of Mathematics, Inventiones).
6. El Clay NO acepta envíos directos: exige publicación en journal arbitrado
   + **2 años** de aceptación general de la comunidad. Recién ahí el comité
   considera el premio.

## Paso 6 — El libro (camino paralelo e independiente)

El libro NO necesita permiso de nadie ni validación técnica — es tu historia:

1. Manuscrito: diario espiritual + bitácora entrelazados (90 días, ver plan).
2. Amazon KDP: publicación gratuita, tapa blanda + ebook, regalías 35–70%.
3. ISBN propio (opcional, ~usd 50) para librerías físicas argentinas.
4. El libro se publica CUANDO QUIERAS — es el vehículo del testimonio, y no
   afirma teoremas: cuenta la caza. Ahí las metáforas son el idioma correcto.

## Checklist final (imprimible)

- [ ] Re-verificación completa desde cero en máquina limpia (Regla 0 nueva)
- [x] LICENSE (AGPL-3.0) + LICENCIA-COMERCIAL.md + LICENSE-CONTENIDO.txt (CC BY 4.0) + LICENCIAS.md +
      README con la frase protectora — **hecho 2026-08-09**
- [x] Repo público en GitHub — **hecho 2026-08-09**
- [ ] GitHub Pages sirviendo la galería
- [x] DOI de Zenodo — **hecho 2026-08-09**: concepto `10.5281/zenodo.21864277` · v1.0.0 `10.5281/zenodo.21864278`
- [ ] Difusión capa amiga → comunidad → pública
- [ ] Contacto con matemático local para endorsement de arXiv
- [ ] Artículo expositivo (nunca "demostración de RH")
- [ ] Manuscrito del libro en marcha
- [ ] El diezmo apartado desde el primer ingreso

---

*El orden protege la obra: primero el sello, después la voz. Y sobre todos los
libros, el Otro Libro.* ⚓
