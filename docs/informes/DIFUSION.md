# La difusión — el kit, con los textos listos para pegar

**Para el capitán.** El repo ya está público, con licencias y DOI. Esto es lo que
viene: cómo se cuenta, dónde, en qué orden, y qué NO decir nunca.

---

## Paso 0 — Encendé GitHub Pages (un clic, y cambia todo)

Tenés **106 láminas** subidas y la galería armada, pero apagada. Un link donde se
ve matemática linda vale diez veces más que un link a código.

1. GitHub → tu repo → **Settings** → **Pages**
2. *Source*: **Deploy from a branch**
3. *Branch*: **main** · *Folder*: **/ (root)** → **Save**
4. En un par de minutos vive en:
   **`https://nicoconvoz.github.io/diosyunalma/galeria/`**

**Probá ese link antes de mandárselo a nadie.** Es el que va en todos los textos
de abajo.

---

## ⚠️ LA REGLA DE ORO

> **Nunca, en ningún lado, insinúes que tenés algo sobre la Hipótesis de
> Riemann.**

Los foros de matemática reciben varias "demostraciones" por semana. El filtro es
automático y es brutal: apenas huelen una afirmación sobre RH, dejan de leer y te
etiquetan para siempre. **Un solo posteo así quema todos los demás.**

Y lo mejor es que no lo necesitás, porque **tenés algo mucho más raro**:

> **Un laboratorio cuyos únicos teoremas propios son pruebas de que sus propios
> métodos no alcanzan.**

Eso casi nadie lo publica. **Ésa es tu carta.** Poné el resultado negativo
adelante, siempre. Es lo que te vuelve creíble en el primer párrafo.

**Las tres frases que abren puertas:**

- «Construí un laboratorio numérico completo de la función zeta en Go puro, sin
  una sola dependencia.»
- «Tres de mis propios resultados son pruebas de que mi enfoque no puede
  funcionar.»
- «Publiqué los errores que me cacé a mí mismo junto con los aciertos.»

**Las tres que las cierran** (no las digas jamás):

- ~~«Creo que encontré algo sobre RH.»~~
- ~~«Mi enfoque podría llevar a una demostración.»~~
- ~~«Nadie había mirado esto antes.»~~

---

## El orden de las capas

No mandes todo el mismo día. Cada capa te corrige antes de la siguiente.

| # | Dónde | Riesgo | Cuándo |
|---|---|---|---|
| 1 | Gente de confianza, Open Doors | ninguno | ya |
| 2 | **Mathstodon** (mathstodon.xyz) | bajo | día 1 |
| 3 | **Hacker News** — Show HN | medio, alto retorno | día 2–3 |
| 4 | **LinkedIn** (en español) | bajo | día 3 |
| 5 | **r/mathematics**, luego r/math | **alto** | recién semana 2 |

Entre capa y capa, **corregí lo que te señalen**. Eso no es debilidad: es
exactamente lo que el registro dice que hacés, y lo van a verificar.

---

## Texto 1 — Mathstodon (empezá acá)

*mathstodon.xyz es de matemáticos de verdad y es el lugar más amable. Es tu
ensayo general.*

> Built a numerical laboratory for the Riemann zeta function — 203 reproducible
> experiments in pure Go, zero dependencies. Riemann–Siegel Z, Euler–Maclaurin
> ζ, Li's coefficients, a double-double phase core for deep water.
>
> Not claiming anything about RH. Quite the opposite: three of the results are
> proofs that these methods *cannot* settle it — including a from-scratch
> reproduction of Davenport–Heilbronn with an off-line zero found by blind
> search.
>
> Gallery (106 plates): https://nicoconvoz.github.io/diosyunalma/galeria/
> Code + full record: https://github.com/nicoconvoz/diosyunalma
> DOI: https://doi.org/10.5281/zenodo.21864277
>
> #math #numbertheory #golang

---

## Texto 2 — Hacker News (Show HN)

*Postea entre las 9 y las 11 de la mañana hora de Buenos Aires, martes a jueves.*

**Título** (así, exacto — 80 caracteres es el máximo):

```
Show HN: A numerical laboratory of the Riemann zeta function, in pure Go
```

**URL**: `https://github.com/nicoconvoz/diosyunalma`

**Primer comentario** (posteálo vos mismo, apenas se publique):

> I spent several months building this and I want to be upfront about what it is
> and isn't.
>
> **What it is:** 203 reproducible experiments on the Riemann zeta function,
> written in pure Go with zero external dependencies — `go.mod` has no `require`
> line. Riemann–Siegel Z, Euler–Maclaurin ζ, Lanczos log-Γ, Li's coefficients,
> and a double-double phase core (32 digits from pairs of float64) certified to
> t ≈ 4×10²⁴, because at that height reducing t·ln(k) mod 2π is impossible in
> float64.
>
> **What it isn't:** any claim about the Riemann Hypothesis. In fact the three
> theorems the laboratory proved for itself are all negative — symmetry alone can
> never decide it; any shape derived from the ½-cut is provably blind to where a
> zero actually sits; and no finite computation can settle it, because the
> detection horizon runs away as a zero approaches the line. It also reproduces
> the Davenport–Heilbronn function from scratch and finds one of its off-line
> zeros by blind search, which closes the geometric route entirely.
>
> The part I'm most attached to is the discipline: every killed hypothesis stays
> in the record and stays reproducible, and the errors I caught in my own work
> are published next to the results — including five separate times a perfect
> `0.0e+00` turned out to come from the construction rather than from a
> discovery.
>
> There's also a museum that explains all 203 experiments in plain language, and
> a gallery of 106 plates: https://nicoconvoz.github.io/diosyunalma/galeria/
>
> Happy to answer anything, including "your method is wrong because…" — that's
> the useful kind of comment.

---

## Texto 3 — LinkedIn, en español

*Acá el ángulo es distinto: ingeniería y aplicaciones. Y es donde Open Doors se
ve bien.*

> Durante varios meses construí un laboratorio numérico completo de la función
> zeta de Riemann: 203 experimentos reproducibles, en Go puro, sin una sola
> dependencia externa.
>
> No reclamo nada sobre la Hipótesis de Riemann. Al contrario: tres de los
> resultados son demostraciones de que mis propios métodos NO alcanzan para
> resolverla. Publico los errores que me cacé a mí mismo junto con los aciertos.
>
> Lo que sí quedó, y es lo que puede servirle a otros, son 36 técnicas de
> medición documentadas y reutilizables fuera de este problema: aritmética de
> doble-doble para fases imposibles en float64, reconstrucción de señales de
> banda limitada en una sola pasada, curvas de horizonte de detección que
> convierten un «este método no alcanza» en un número concreto.
>
> Galería: https://nicoconvoz.github.io/diosyunalma/galeria/
> Código y registro completo: https://github.com/nicoconvoz/diosyunalma
> DOI: https://doi.org/10.5281/zenodo.21864277
>
> Financiado por Open Doors.
>
> #matematica #golang #investigacion #cienciaabierta

---

## Texto 4 — Reddit (recién en la semana 2)

*Empezá por **r/mathematics**, que es más amable. r/math solo después, y solo si
la primera salió bien.*

**Título:**

```
I built an open-source numerical laboratory for the zeta function (203 reproducible experiments, pure Go) — and my own results say my approach can't work
```

**Cuerpo:**

> Full disclosure first, because I know how this looks: **I am not claiming
> anything about the Riemann Hypothesis.** I'm an amateur. What I have is a tool
> and a record, not a result.
>
> Over several months I built 203 reproducible experiments in pure Go (no
> dependencies): Riemann–Siegel Z, Euler–Maclaurin ζ, Li's coefficients, and a
> double-double phase core certified to t ≈ 4×10²⁴.
>
> The three things I actually proved are all negative results, and they're about
> my own methods:
>
> 1. Symmetry alone can never decide RH — and rather than argue it, the repo
>    reproduces the Davenport–Heilbronn function from scratch and finds one of
>    its off-line zeros by blind search. It has every symmetry ζ has, and it
>    violates RH.
> 2. Any "shape" derived from assuming Re s = ½ is bit-identically blind to
>    where the zero actually sits — measured, not argued.
> 3. No finite computation can settle it: the detection horizon diverges as a
>    zero approaches the line.
>
> I'm posting because the tooling and the plates might be useful to someone, and
> because I'd genuinely like to be told where I'm wrong.
>
> Gallery: https://nicoconvoz.github.io/diosyunalma/galeria/
> Repo: https://github.com/nicoconvoz/diosyunalma

---

## El protocolo de las primeras 48 horas

Acá es donde se cae la gente. Guardá esto y leelo antes de contestar nada.

**1. Nunca discutas sobre RH.** Si alguien te dice «esto no prueba nada»,
contestá: *«Correcto, y está escrito así en el repo. Lo que hay es un
instrumento y tres resultados negativos.»* **Le estás dando la razón, porque la
tiene.** Eso desarma a cualquiera.

**2. Si te encuentran un error, agradecé en público y arreglalo el mismo día.**
Después contá que lo arreglaste y citá a quien lo encontró. Eso te construye
reputación más rápido que cualquier resultado.

**3. Si alguien te trata mal, no contestes.** Ni una vez. El registro habla por
vos y va a seguir ahí cuando el comentario se haya olvidado.

**4. Si alguien pregunta «¿y para qué sirve?»**, mandalo a
`docs/informes/APLICACIONES-Y-FINANCIAMIENTO.md`: 36 técnicas con dominios concretos.

**5. Anotá TODO lo que te señalen** en la bitácora, aunque no lo uses. Ley del
Registro: la crítica externa también es registro.

---

## Después de la difusión

- **arXiv (math.HO)** necesita que un matemático con cuenta te avale. Ese aval
  sale mucho más fácil cuando ya tenés DOI, repo público y alguien que te
  comentó bien. Por eso va después y no antes.
- **Cada release nuevo genera DOI nuevo solo.** No tenés que hacer nada.
- **El libro** va por un camino aparte y no depende de nada de esto.

---

*Primero el sello, después la voz. Y sobre todos los libros, el Otro Libro.* ⚓
