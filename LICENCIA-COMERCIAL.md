# Licencia comercial

El código de este laboratorio se publica bajo **AGPL-3.0** ([LICENSE](LICENSE)).
Para muchos usos eso alcanza y no hay que pagar nada.

Para otros, no alcanza. Para esos existe esta licencia comercial.

---

## ¿Necesitás una licencia comercial?

**NO la necesitás si:**

- Estudiás el código, lo corrés, lo modificás para vos.
- Lo usás en investigación académica y publicás lo que hagas.
- Lo usás en un proyecto propio que también es AGPL.
- Enseñás con él.

**SÍ la necesitás si:**

- Querés incorporar este código —o partes— a un **producto cerrado** que
  distribuís sin publicar tu fuente.
- Querés **servirlo por red** (SaaS, API, servicio interno accesible a
  terceros) sin publicar el código completo de tu servicio. *Esto es lo
  específico de la AGPL: la sección 13 obliga a entregar el fuente incluso a
  quien solo interactúa por red.*
- Tu organización tiene una política que prohíbe copyleft fuerte.
- Necesitás garantías, soporte o indemnidad que la AGPL explícitamente no da.

Regla práctica: **si tu abogado se pone nervioso con la AGPL, necesitás esta
licencia.**

---

## Qué otorga

Una licencia comercial reemplaza las obligaciones de la AGPL por términos
negociados. Típicamente:

| Punto | AGPL-3.0 | Licencia comercial |
|---|---|---|
| Publicar tu código derivado | obligatorio | **no** |
| Entregar fuente a usuarios de red (§13) | obligatorio | **no** |
| Uso en producto propietario | no permitido | **sí** |
| Precio | gratis | a convenir |
| Soporte y garantía | ninguna | a convenir |

La forma del precio se acuerda en cada caso: pago único, suscripción anual,
por instalación, o regalía sobre ingresos del producto que lo incorpore.

---

## Qué NO cubre ninguna licencia

Hay que decirlo claro, porque hace a la honestidad de este laboratorio:

- **La matemática no se licencia.** Los teoremas, las identidades y las
  mediciones publicadas acá son de dominio público el minuto en que se
  publican. Nadie necesita permiso para usarlos, ni podríamos dárselo.
- Lo que sí está licenciado es el **código** (AGPL o comercial) y el
  **contenido** —láminas, textos, bitácora— bajo
  [CC BY 4.0](LICENSE-CONTENIDO.txt).
- Esta licencia **no vende resultados ni conclusiones**: vende software y el
  derecho a integrarlo sin abrir el tuyo.

---

## Qué hay acá que valga la pena licenciar

El catálogo completo, con las 36 técnicas y sus dominios de aplicación, está
en [docs/APLICACIONES-Y-FINANCIAMIENTO.md](docs/APLICACIONES-Y-FINANCIAMIENTO.md).
Lo más maduro, en una línea cada uno:

- **Motor de fase doble-doble** — 32 dígitos con pares de float64; reducción
  mod 2π exacta donde el float64 se rinde. Techo certificado t ≈ 4×10²⁴.
- **Balde de luz de banda limitada** — deposita todos los términos UNA vez
  sobre una grilla de Nyquist; después se lee cualquier punto gratis.
- **Plegado de Fresnel / Landsberg–Schaar** — resuelve 10⁸ términos con 7,
  error 1.3×10⁻¹², aceleración del riel ~8×10⁹.
- **Transmisión adaptativa (`gearFor`)** — elige el algoritmo por medición del
  régimen, no por suposición.
- **Curva de ceguera** — convierte «este método no alcanza» en un número: hasta
  qué profundidad el instrumento ve, y desde dónde no.
- **Detección de mediciones tautológicas** — verifica si el instrumento podía
  haber devuelto otra cosa antes de aceptar un resultado perfecto.

---

## Cómo pedirla

Escribí a **mza.nicolas.astorga@gmail.com** con:

1. Qué parte del laboratorio querés usar.
2. En qué producto o servicio.
3. Escala aproximada (usuarios, instalaciones, ingresos previstos).

Se responde con una propuesta concreta.

---

## Titularidad

La obra tiene dos co-titulares de copyright:

- **Jesús Nicolás Astorga** — autor
- **RESOURCES OPEN DOORS S.A.S** — organización de ingeniería que financia el
  trabajo

**Toda licencia comercial requiere el acuerdo de ambos.** Cualquier propuesta
se evalúa y se firma entre las dos partes.

---

*Este documento describe la disponibilidad de una licencia comercial; no es en
sí mismo un contrato ni una oferta vinculante. Los términos definitivos se
acuerdan por escrito en cada caso.*
