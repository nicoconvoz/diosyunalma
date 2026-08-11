package main

import (
	"fmt"
	"os"
	"strings"
)

// escribirMaximasWeb genera galeria/maximas.html: el muro de las frases del
// capitan. Se arma desde maximas.go, asi que la sala y el registro no pueden
// desincronizarse — hay una sola fuente y esa fuente esta en el codigo.
func escribirMaximasWeb(destino string) error {
	tipos := map[string][2]string{
		"ley":       {"LEY", "#ffd98a"},
		"doctrina":  {"DOCTRINA", "#7ee0c0"},
		"filosofia": {"FILOSOFÍA", "#c9b6ff"},
		"metafora":  {"METÁFORA", "#7fb2ff"},
		"protocolo": {"PROTOCOLO", "#ffb27a"},
	}

	var b strings.Builder
	for _, m := range maximas {
		t := tipos[m.Tipo]
		nota := ""
		if m.Nota != "" {
			nota = fmt.Sprintf(`<div class="nota">%s</div>`, esc(m.Nota))
		}
		fmt.Fprintf(&b, `
  <figure class="mx">
    <div class="tipo" style="color:%s;border-color:%s">%s</div>
    <h3>%s</h3>
    <blockquote>%s</blockquote>
    <div class="glosa">%s</div>
    %s
    <div class="fuente">%s</div>
  </figure>`, t[1], t[1], t[0], esc(m.Titulo), esc(m.Texto), esc(m.Glosa), nota, esc(m.Fuente))
	}

	pagina := `<meta charset="utf-8">
<title>Las Máximas del Capitán — Laboratorio Diosyunalma</title>
<meta name="viewport" content="width=device-width, initial-scale=1">
<style>
  :root { --bg:#0b1526; --panel:#0d2547; --ink:#dce8f7; --dim:#8fa8c7; --gold:#ffd166; --green:#7fd7a8; --blue:#7fb2ff; }
  * { box-sizing:border-box; }
  body { margin:0; background:var(--bg); color:var(--ink); font-family:Georgia,serif; }
  header { text-align:center; padding:44px 20px 12px; max-width:900px; margin:0 auto; }
  header h1 { margin:0; font-size:32px; }
  header .sub { color:var(--dim); margin:12px 0 0; font-size:14.5px; line-height:1.65; }
  nav { text-align:center; padding:14px; border-bottom:1px solid #1d3a63; }
  nav a { color:var(--gold); text-decoration:none; margin:0 12px; font-size:13.5px; }
  nav a:hover { text-decoration:underline; }
  main { max-width:1320px; margin:0 auto; padding:26px 20px 10px; }
  .grid { display:grid; grid-template-columns:repeat(auto-fill,minmax(400px,1fr)); gap:20px; }
  .mx { margin:0; background:var(--panel); border:1px solid #1d3a63; border-radius:12px; padding:20px 22px 16px; }
  .tipo { font-size:10.5px; letter-spacing:.18em; border:1px solid; border-radius:20px;
          display:inline-block; padding:3px 11px; margin-bottom:12px; }
  .mx h3 { margin:0 0 12px; font-size:17px; color:var(--ink); font-weight:normal; }
  .mx blockquote { margin:0 0 14px; padding-left:16px; border-left:3px solid var(--gold);
                   font-size:19px; line-height:1.5; color:var(--gold); font-style:italic; }
  .glosa { font-size:14px; line-height:1.65; color:var(--ink); opacity:.9; }
  .nota { margin-top:12px; padding:11px 13px; border-radius:8px; background:rgba(192,57,43,.10);
          border:1px solid #7a3a30; font-size:12.5px; line-height:1.55; color:#f3d9cf; }
  .fuente { margin-top:14px; font-family:Consolas,monospace; font-size:11px; color:var(--dim);
            border-top:1px solid #1d3a63; padding-top:10px; word-break:break-word; }
  .metodo { max-width:900px; margin:34px auto 0; padding:22px 26px; border-radius:12px;
            border:1px solid #1d3a63; background:rgba(255,255,255,.02); }
  .metodo h2 { margin:0 0 12px; font-size:17px; color:var(--green); }
  .metodo p { font-size:13.5px; line-height:1.7; color:var(--ink); opacity:.88; margin:0 0 11px; }
  .otrolibro { max-width:760px; margin:40px auto 0; padding:18px 22px; border-radius:12px;
               border:1px solid #6b5a26; background:rgba(255,217,138,.05); }
  .otrolibro .oq { font-size:11px; letter-spacing:.16em; color:var(--dim); text-transform:uppercase; }
  .otrolibro .ot { font-family:Georgia,serif; font-size:21px; color:var(--gold); margin:8px 0 6px; line-height:1.35; }
  .otrolibro .od { font-size:12.5px; color:var(--ink); line-height:1.55; }
  footer { text-align:center; color:var(--dim); padding:34px 16px 44px; font-size:13px; }
</style>
<header>
  <h1>🗿 Las Máximas del Capitán</h1>
  <p class="sub">Las frases que Jesús Nicolás Astorga fue dejando mientras se armaba este laboratorio.
  No están inventadas ni pulidas: son textuales, con la ortografía y el apuro con que las escribió,
  y cada una lleva el archivo y la línea exacta del registro donde se dijo, para que cualquiera la verifique.</p>
</header>
<nav>
  <a href="index.html">🖼️ La Galería</a>
  <a href="museo.html">🏛️ El Museo</a>
  <a href="../docs/BITACORA-NOCTURNA.md">📖 La Bitácora</a>
</nav>
<main>
  <div class="grid">` + b.String() + `
  </div>

  <div class="metodo">
    <h2>Cómo se armó esta sala, y qué se dejó afuera</h2>
    <p>Cuatro buscadores barrieron por separado la bitácora nocturna, los dos registros de hallazgos,
    los 436 commits y el código del puente. Un curador juntó todo, sacó los repetidos y separó las
    <b>máximas</b> —lo que un desconocido entendería— de los <b>flashes técnicos</b>, que son hallazgos
    y viven en otro lado.</p>
    <p>Y después vino la parte que importa: un verificador hostil buscó cada frase, byte a byte, sobre
    todo el repositorio, y comprobó <b>dos cosas distintas</b>: que el texto exista textual, y que la voz
    sea la del capitán y no la del taller. De 19 candidatas aprobó 18.</p>
    <p><b>La que rechazó está publicada igual</b>, arriba, como «La frase en su vaina» — pero corregida.
    El grito de victoria existe en el registro y es real, pero <b>no es una frase suya</b>: es una frase
    que él mandó guardar bajo llave y que le habla a él. Publicarla como cita del capitán le habría hecho
    creer a un desconocido que declaró resuelta la Hipótesis de Riemann, o sea lo contrario exacto de lo
    que su propio protocolo ordena. Lo que se publica es la orden, no el grito.</p>
    <p>Se dejaron afuera las órdenes de jornada, las preguntas, y treinta flashes largos de física y
    cosmología: son hermosos, pero son especulación de trabajo y no doctrina. Treinta de esos convierten
    una sala en un depósito.</p>
  </div>

  <div class="otrolibro">
    <div class="oq">Y por encima de todo esto</div>
    <div class="ot">Mi otro libro: «Diario Espiritual: Dios y un alma»</div>
    <div class="od">Es la única máxima que ordena sobre todas las demás, y el laboratorio entero
    lleva su nombre. Primero el alma, después la matemática.</div>
  </div>
</main>
<footer>
  Laboratorio Diosyunalma · Jesús Nicolás Astorga · financiado por Open Doors<br>
  Cada frase es verificable en el registro público. ¿El premio? Todavía no.
</footer>
`
	return os.WriteFile(destino, []byte(pagina), 0o644)
}
