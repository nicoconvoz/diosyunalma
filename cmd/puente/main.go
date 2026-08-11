// Command puente is EL PUENTE DE MANDO: one program that starts the whole
// laboratory. It knows every experiment in cmd/, runs them on demand from a
// dashboard, streams their output live, and shows the plate each one draws -
// so the captain never has to launch anything by hand again.
//
//	go run ./cmd/puente          # opens the bridge at http://localhost:8118
//
// The engine lives here; the catalogue of experiments lives in catalogo.go.
package main

import (
	"archive/zip"
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const puerto = "8118"

// maxParalelo caps how many experiments compile and run at once so the
// machine stays usable; lookouts do not consume a slot.
const maxParalelo = 3

// maxLineas is the size of each experiment's output ring buffer.
const maxLineas = 400

// Exp is one experiment of the laboratory.
type Exp struct {
	Dir     string   // directory under cmd/
	Emoji   string   // one emoji
	Titulo  string   // short Spanish title
	Criollo string   // one plain-Spanish line
	Cat     string   // hall key
	Salida  string   // file it writes (svg/wav/png/log), or ""
	Vel     string   // rapido | medio | lento | vigia
	Nota    string   // flags, port, warning
	Args    []string // extra arguments (e.g. -cazar)
}

// Sala is a hall of the bridge: a named group of experiments.
type Sala struct {
	Clave, Emoji, Nombre, Detalle string
}

var salas = []Sala{
	{"caras", "💎", "Las Siete Caras", "la campaña final: el eslabón rojo visto desde todos sus ángulos"},
	{"acta", "📜", "El Acta y el Molde", "el intento formal, el reloj de sol, los premios, el Campo"},
	{"atomo", "⚛️", "El Átomo", "los planos, el espejo, la esfera, el murciélago, el sótano"},
	{"armonizacion", "🌀", "La Armonización", "la dimensión 0, el cambiaformas, el horno, el neutrón"},
	{"mar", "🌊", "El Mar y el Tren", "la caza, la cartografía, las costas y los curlicues"},
	{"instrumentos", "🔧", "Los Instrumentos", "motores y servidores: el faro, el tren, los cascos"},
	{"archivo", "🗄️", "El Archivo", "los experimentos de la primera era"},
}

// ---- state ----

type estado struct {
	Fase      string    `json:"fase"` // listo | cola | corriendo | ok | falla | detenido
	Inicio    time.Time `json:"-"`
	Fin       time.Time `json:"-"`
	Veredicto string    `json:"veredicto"`
	Lineas    []string  `json:"-"`
	Codigo    int       `json:"codigo"`
	pid       int
	parar     chan struct{}
}

type puente struct {
	mu      sync.Mutex
	est     map[string]*estado
	laminas map[string]string // plate filename -> newest path on disk
	exes    map[string]string // experiment dir -> compiled binary, when present
	sem     chan struct{}
	raiz    string
}

func nuevoPuente(raiz string) *puente {
	p := &puente{
		est:     map[string]*estado{},
		laminas: map[string]string{},
		exes:    map[string]string{},
		sem:     make(chan struct{}, maxParalelo),
		raiz:    raiz,
	}
	for _, e := range catalogo {
		p.est[e.Dir] = &estado{Fase: "listo"}
	}
	p.indexarLaminas()
	go func() {
		for range time.Tick(3 * time.Second) {
			p.indexarLaminas()
		}
	}()
	return p
}

// ---- launching an experiment ----
//
// The shop runs from source with `go run`, but the bridge also has to work on
// a machine where Go was never installed - an investor's laptop, a fresh
// desktop. So a compiled binary in bin/ always wins, and Go is only the
// fallback for the workshop.

func rutaEjecutable(raiz, dir string) string {
	nombre := dir
	if runtime.GOOS == "windows" {
		nombre += ".exe"
	}
	ruta := filepath.Join(raiz, "bin", nombre)
	if fi, err := os.Stat(ruta); err == nil && !fi.IsDir() {
		return ruta
	}
	return ""
}

func hayGo() bool {
	_, err := exec.LookPath("go")
	return err == nil
}

func (p *puente) comandoDe(e Exp) (*exec.Cmd, error) {
	if exe := rutaEjecutable(p.raiz, e.Dir); exe != "" {
		return exec.Command(exe, e.Args...), nil
	}
	if hayGo() {
		return exec.Command("go", append([]string{"run", "./cmd/" + e.Dir}, e.Args...)...), nil
	}
	return nil, fmt.Errorf("falta el ejecutable bin/%s y no hay Go instalado para compilarlo", e.Dir)
}

func (p *puente) buscar(dir string) (Exp, bool) {
	for _, e := range catalogo {
		if e.Dir == dir {
			return e, true
		}
	}
	return Exp{}, false
}

// esRelleno reports whether a line is a decorative rule or a file receipt -
// the kind of line that looks important and says nothing.
func esRelleno(l string) bool {
	if l == "" || strings.HasPrefix(l, "escrita:") || strings.HasPrefix(l, "written:") {
		return true
	}
	adorno, total := 0, 0
	for _, r := range l {
		if r == ' ' {
			continue
		}
		total++
		switch r {
		case '═', '─', '━', '-', '=', '·', '_', '*':
			adorno++
		}
	}
	return total == 0 || float64(adorno)/float64(total) > 0.55
}

func tieneDigito(l string) bool {
	return strings.ContainsAny(l, "0123456789")
}

// veredictoDe picks the most meaningful line of an experiment's output: the
// judge's line when there is one, otherwise the last line that actually says
// something. Decorative rules and "escrita: x.svg" receipts never win.
func veredictoDe(lineas []string) string {
	marcas := []string{"⚖", "🏆", "★", "VEREDICTO", "→"}
	limpias := make([]string, 0, len(lineas))
	for _, l := range lineas {
		if t := strings.TrimSpace(l); !esRelleno(t) {
			limpias = append(limpias, t)
		}
	}
	marcada := func(l string) bool {
		for _, m := range marcas {
			if strings.Contains(l, m) {
				return true
			}
		}
		return false
	}
	for _, exigirDigito := range []bool{true, false} {
		for i := len(limpias) - 1; i >= 0; i-- {
			if marcada(limpias[i]) && (!exigirDigito || tieneDigito(limpias[i])) {
				return limpias[i]
			}
		}
	}
	for i := len(limpias) - 1; i >= 0; i-- {
		if tieneDigito(limpias[i]) {
			return limpias[i]
		}
	}
	if len(limpias) > 0 {
		return limpias[len(limpias)-1]
	}
	return ""
}

func (p *puente) correr(e Exp) {
	p.mu.Lock()
	st := p.est[e.Dir]
	if st.Fase == "corriendo" || st.Fase == "cola" {
		p.mu.Unlock()
		return
	}
	st.Fase = "cola"
	st.Lineas = nil
	st.Veredicto = ""
	st.Codigo = 0
	st.parar = make(chan struct{})
	p.mu.Unlock()

	go func() {
		vigia := e.Vel == "vigia"
		if !vigia {
			p.sem <- struct{}{}
			defer func() { <-p.sem }()
		}

		p.mu.Lock()
		st.Fase = "corriendo"
		st.Inicio = time.Now()
		p.mu.Unlock()

		cmd, err := p.comandoDe(e)
		if err != nil {
			p.terminar(st, "falla", 1, err.Error())
			return
		}
		cmd.Dir = p.raiz
		salida, err2 := cmd.StdoutPipe()
		if err2 != nil {
			p.terminar(st, "falla", 1, "no pude abrir la salida: "+err2.Error())
			return
		}
		cmd.Stderr = cmd.Stdout
		if err := cmd.Start(); err != nil {
			p.terminar(st, "falla", 1, "no pude arrancar: "+err.Error())
			return
		}

		p.mu.Lock()
		st.pid = cmd.Process.Pid
		p.mu.Unlock()

		listo := make(chan struct{})
		go func() {
			select {
			case <-st.parar:
				matarArbol(cmd.Process.Pid)
			case <-listo:
			}
		}()

		sc := bufio.NewScanner(salida)
		sc.Buffer(make([]byte, 0, 256*1024), 4*1024*1024)
		for sc.Scan() {
			linea := sc.Text()
			p.mu.Lock()
			st.Lineas = append(st.Lineas, linea)
			if len(st.Lineas) > maxLineas {
				st.Lineas = st.Lineas[len(st.Lineas)-maxLineas:]
			}
			p.mu.Unlock()
		}
		errEspera := cmd.Wait()
		close(listo)

		p.mu.Lock()
		detenido := st.Fase == "deteniendo"
		lineas := append([]string(nil), st.Lineas...)
		p.mu.Unlock()

		fase, codigo := "ok", 0
		switch {
		case detenido:
			fase = "detenido"
		case errEspera != nil:
			fase, codigo = "falla", 1
			if ee, ok := errEspera.(*exec.ExitError); ok {
				codigo = ee.ExitCode()
			}
		}
		p.terminar(st, fase, codigo, veredictoDe(lineas))
	}()
}

func (p *puente) terminar(st *estado, fase string, codigo int, veredicto string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	st.Fase = fase
	st.Codigo = codigo
	st.Fin = time.Now()
	if veredicto != "" {
		st.Veredicto = veredicto
	}
	if st.Inicio.IsZero() {
		st.Inicio = st.Fin
	}
}

func (p *puente) detener(dir string) {
	p.mu.Lock()
	st, ok := p.est[dir]
	if !ok || (st.Fase != "corriendo" && st.Fase != "cola") {
		p.mu.Unlock()
		return
	}
	st.Fase = "deteniendo"
	parar := st.parar
	p.mu.Unlock()
	if parar != nil {
		select {
		case <-parar:
		default:
			close(parar)
		}
	}
}

// ---- casting off ----
//
// A previous voyage that was closed by shutting the window leaves its
// process holding the port, and the next launch finds the door blocked.
// So the bridge always clears its own berth before opening: it asks the
// system who is holding the port and sends that whole tree away - never
// itself, and never the `go run` parent that carries it.

// duenosDelPuerto asks the system which processes hold the given TCP port.
func duenosDelPuerto(port string) []int {
	var salida []byte
	var err error
	if runtime.GOOS == "windows" {
		salida, err = exec.Command("netstat", "-ano", "-p", "tcp").Output()
	} else {
		salida, err = exec.Command("lsof", "-ti", "tcp:"+port).Output()
	}
	if err != nil {
		return nil
	}
	yo, papa := os.Getpid(), os.Getppid()
	vistos := map[int]bool{yo: true, papa: true}
	var pids []int
	for _, linea := range strings.Split(string(salida), "\n") {
		campos := strings.Fields(linea)
		var pid int
		if runtime.GOOS == "windows" {
			// PROTO  DIRECCION-LOCAL  DIRECCION-REMOTA  ESTADO  PID
			if len(campos) < 4 || !strings.HasSuffix(strings.ToUpper(campos[0]), "TCP") {
				continue
			}
			if !strings.HasSuffix(campos[1], ":"+port) {
				continue
			}
			pid, _ = strconv.Atoi(campos[len(campos)-1])
		} else {
			if len(campos) != 1 {
				continue
			}
			pid, _ = strconv.Atoi(campos[0])
		}
		if pid <= 0 || vistos[pid] {
			continue
		}
		vistos[pid] = true
		pids = append(pids, pid)
	}
	return pids
}

// soltarAmarras frees the bridge's berth and returns the ports it cleared.
func soltarAmarras(port string) []int {
	pids := duenosDelPuerto(port)
	for _, pid := range pids {
		matarArbol(pid)
	}
	return pids
}

// escuchar opens the port, clearing whoever holds it and retrying while the
// system releases the socket.
func escuchar(port string) (net.Listener, error) {
	if previos := soltarAmarras(port); len(previos) > 0 {
		fmt.Printf("   amarras sueltas: cerré %d proceso(s) del viaje anterior %v\n", len(previos), previos)
	}
	var ultimo error
	for intento := 0; intento < 20; intento++ {
		ln, err := net.Listen("tcp", ":"+port)
		if err == nil {
			return ln, nil
		}
		ultimo = err
		soltarAmarras(port)
		time.Sleep(150 * time.Millisecond)
	}
	return nil, ultimo
}

// matarArbol kills the process and every child it spawned. On Windows a
// `go run` parent leaves the compiled binary running, so the whole tree
// must go - the shop learned this the hard way.
func matarArbol(pid int) {
	if runtime.GOOS == "windows" {
		_ = exec.Command("taskkill", "/PID", strconv.Itoa(pid), "/T", "/F").Run()
		return
	}
	if pr, err := os.FindProcess(pid); err == nil {
		_ = pr.Kill()
	}
}

// ---- plate index ----
//
// Plates live in two places: freshly written ones land in the root (the
// workbench) and archived ones sit in the gallery halls (the museum). The
// index is rebuilt in the background so the dashboard never stats every file
// per refresh; the newest copy of a name always wins.

func (p *puente) indexarLaminas() {
	exes := map[string]string{}
	for _, e := range catalogo {
		if ruta := rutaEjecutable(p.raiz, e.Dir); ruta != "" {
			exes[e.Dir] = ruta
		}
	}
	idx := map[string]string{}
	mtime := map[string]time.Time{}
	poner := func(ruta string) {
		fi, err := os.Stat(ruta)
		if err != nil || fi.IsDir() {
			return
		}
		n := filepath.Base(ruta)
		if t, ok := mtime[n]; !ok || fi.ModTime().After(t) {
			idx[n], mtime[n] = ruta, fi.ModTime()
		}
	}
	for _, pat := range []string{
		filepath.Join(p.raiz, "*"),
		filepath.Join(p.raiz, "galeria", "laminas", "*", "*"),
		filepath.Join(p.raiz, "galeria", "sonidos", "*"),
		filepath.Join(p.raiz, "luz", "*"),
	} {
		m, _ := filepath.Glob(pat)
		for _, ruta := range m {
			poner(ruta)
		}
	}
	p.mu.Lock()
	p.laminas = idx
	p.exes = exes
	p.mu.Unlock()
}

// rutaLamina must be called with p.mu held.
func (p *puente) rutaLamina(archivo string) string {
	if archivo == "" {
		return ""
	}
	return p.laminas[archivo]
}

// ---- API ----

type fichaJSON struct {
	Dir       string `json:"dir"`
	Emoji     string `json:"emoji"`
	Titulo    string `json:"titulo"`
	Criollo   string `json:"criollo"`
	Cat       string `json:"cat"`
	Vel       string `json:"vel"`
	Nota      string `json:"nota"`
	Salida    string `json:"salida"`
	Fase      string `json:"fase"`
	Segundos  int    `json:"segundos"`
	Veredicto string `json:"veredicto"`
	Lineas    int    `json:"lineas"`
	Lamina    bool   `json:"lamina"`
	Codigo    int    `json:"codigo"`
	Motor     string `json:"motor"` // compilado | go | falta
}

func (p *puente) apiEstado(w http.ResponseWriter, r *http.Request) {
	conGo := hayGo()
	p.mu.Lock()
	fichas := make([]fichaJSON, 0, len(catalogo))
	corriendo, hechos, compilados := 0, 0, 0
	for _, e := range catalogo {
		motor := "falta"
		if p.exes[e.Dir] != "" {
			motor, compilados = "compilado", compilados+1
		} else if conGo {
			motor = "go"
		}
		st := p.est[e.Dir]
		seg := 0
		switch st.Fase {
		case "corriendo", "deteniendo":
			seg = int(time.Since(st.Inicio).Seconds())
			corriendo++
		case "ok", "falla", "detenido":
			seg = int(st.Fin.Sub(st.Inicio).Seconds())
			if st.Fase == "ok" {
				hechos++
			}
		}
		fichas = append(fichas, fichaJSON{
			Dir: e.Dir, Emoji: e.Emoji, Titulo: e.Titulo, Criollo: e.Criollo,
			Cat: e.Cat, Vel: e.Vel, Nota: e.Nota, Salida: e.Salida,
			Fase: st.Fase, Segundos: seg, Veredicto: st.Veredicto,
			Lineas: len(st.Lineas), Lamina: p.rutaLamina(e.Salida) != "", Codigo: st.Codigo,
			Motor: motor,
		})
	}
	p.mu.Unlock()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"fichas": fichas, "corriendo": corriendo, "hechos": hechos, "total": len(catalogo),
		"compilados": compilados, "hayGo": conGo,
	})
}

func (p *puente) apiSalida(w http.ResponseWriter, r *http.Request) {
	dir := r.URL.Query().Get("dir")
	p.mu.Lock()
	st, ok := p.est[dir]
	var texto, fase string
	if ok {
		texto = strings.Join(st.Lineas, "\n")
		fase = st.Fase
	}
	p.mu.Unlock()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]any{"texto": texto, "fase": fase})
}

func (p *puente) apiCorrer(w http.ResponseWriter, r *http.Request) {
	if dir := r.URL.Query().Get("dir"); dir != "" {
		if e, ok := p.buscar(dir); ok {
			p.correr(e)
		}
	}
	if cat := r.URL.Query().Get("sala"); cat != "" {
		for _, e := range catalogo {
			if e.Cat == cat && e.Vel != "vigia" {
				p.correr(e)
			}
		}
	}
	// "todo" runs only the quick experiments: the slow ones sweep zeros to
	// t=1000 and would hold the queue for the rest of the afternoon.
	if r.URL.Query().Get("todo") == "1" {
		for _, e := range catalogo {
			if e.Vel == "rapido" {
				p.correr(e)
			}
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

func (p *puente) apiDetener(w http.ResponseWriter, r *http.Request) {
	if dir := r.URL.Query().Get("dir"); dir != "" {
		p.detener(dir)
	}
	if r.URL.Query().Get("todo") == "1" {
		p.mu.Lock()
		dirs := make([]string, 0, len(p.est))
		for d, st := range p.est {
			if st.Fase == "corriendo" || st.Fase == "cola" {
				dirs = append(dirs, d)
			}
		}
		p.mu.Unlock()
		for _, d := range dirs {
			p.detener(d)
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

func (p *puente) apiLamina(w http.ResponseWriter, r *http.Request) {
	e, ok := p.buscar(r.URL.Query().Get("dir"))
	if !ok {
		http.NotFound(w, r)
		return
	}
	p.mu.Lock()
	ruta := p.rutaLamina(e.Salida)
	p.mu.Unlock()
	if ruta == "" {
		http.NotFound(w, r)
		return
	}
	switch strings.ToLower(filepath.Ext(ruta)) {
	case ".svg":
		w.Header().Set("Content-Type", "image/svg+xml; charset=utf-8")
	case ".wav":
		w.Header().Set("Content-Type", "audio/wav")
	case ".png":
		w.Header().Set("Content-Type", "image/png")
	default:
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	}
	http.ServeFile(w, r, ruta)
}

// ---- the library ----
//
// The research is as much a part of the laboratory as the experiments, so the
// bridge serves the whole docs/ folder: listed with a plain-Spanish blurb, and
// read inside the dashboard itself.

var fichasDoc = map[string][2]string{
	"DIFUSION.md":                      {"📢 La Difusión", "el kit para sacar el laboratorio al mundo: los textos ya escritos, el orden de las capas y qué no decir nunca"},
	"EL-MUSEO.md":                      {"🏛️ El Museo", "el recorrido de ocho salas explicado a lo criollo, para que lo entienda cualquiera"},
	"RECORRIDO.md":                     {"El Recorrido Técnico", "el viaje completo con números: instrumentos, la caza, el acta y las siete caras"},
	"INFORME-TECNICO.md":               {"El Informe Técnico", "el informe formal, el salón de honor y los treinta y ocho que nacieron del capitán"},
	"BITACORA-NOCTURNA.md":             {"La Bitácora Nocturna", "el registro fechado de cada hallazgo, cada error y cada corrección"},
	"HALLAZGOS-ES.md":                  {"Los Hallazgos", "los hallazgos numerados, en castellano"},
	"FINDINGS.md":                      {"Findings", "el registro de la primera era, en inglés"},
	"VALIDACION.md":                    {"La Validación", "guía paso a paso para que un revisor reproduzca todo por su cuenta"},
	"APLICACIONES-Y-FINANCIAMIENTO.md": {"Aplicaciones y Financiamiento", "en qué ramas aplica lo construido y cómo se sostiene el laboratorio"},
	"COMO-PUBLICAR.md":                 {"Cómo Publicar", "los pasos para sacar la obra al mundo sin tropiezos"},
	"CEROS-VIRGENES.md":                {"Los Ceros Vírgenes", "el catálogo de agua que nadie había navegado"},
	"EL-TREN.md":                       {"El Tren", "el motor de caza y su bitácora de aguas anexadas"},
	"NAVE-VIVA.md":                     {"La Nave Viva", "la arquitectura de la flota y sus motores"},
	"LEEME-PAQUETE.txt":                {"Léeme del paquete", "la hoja de bienvenida para quien recibe el laboratorio"},
	"README.md":                        {"README del proyecto", "la puerta de entrada del repositorio, en inglés"},
}

// ordenDoc puts the documents a visitor should read first at the top.
var ordenDoc = []string{"RECORRIDO.md", "INFORME-TECNICO.md", "BITACORA-NOCTURNA.md", "HALLAZGOS-ES.md",
	"VALIDACION.md", "APLICACIONES-Y-FINANCIAMIENTO.md", "COMO-PUBLICAR.md", "FINDINGS.md"}

func (p *puente) rutaDoc(f string) string {
	f = filepath.Base(f)
	if !strings.HasSuffix(f, ".md") && !strings.HasSuffix(f, ".txt") {
		return ""
	}
	for _, dir := range []string{"docs", "."} {
		ruta := filepath.Join(p.raiz, dir, f)
		if fi, err := os.Stat(ruta); err == nil && !fi.IsDir() {
			return ruta
		}
	}
	return ""
}

type docJSON struct {
	Archivo string `json:"archivo"`
	Titulo  string `json:"titulo"`
	Detalle string `json:"detalle"`
	Lineas  int    `json:"lineas"`
	KB      int    `json:"kb"`
}

func (p *puente) apiDocs(w http.ResponseWriter, r *http.Request) {
	vistos := map[string]bool{}
	var lista []docJSON
	agregar := func(ruta, nombre string) {
		if vistos[nombre] {
			return
		}
		b, err := os.ReadFile(ruta)
		if err != nil {
			return
		}
		vistos[nombre] = true
		titulo, detalle := nombre, ""
		if m, ok := fichasDoc[nombre]; ok {
			titulo, detalle = m[0], m[1]
		}
		lista = append(lista, docJSON{
			Archivo: nombre, Titulo: titulo, Detalle: detalle,
			Lineas: strings.Count(string(b), "\n") + 1, KB: (len(b) + 512) / 1024,
		})
	}
	for _, nombre := range ordenDoc {
		if ruta := p.rutaDoc(nombre); ruta != "" {
			agregar(ruta, nombre)
		}
	}
	entradas, _ := os.ReadDir(filepath.Join(p.raiz, "docs"))
	for _, e := range entradas {
		if !e.IsDir() {
			if ruta := p.rutaDoc(e.Name()); ruta != "" {
				agregar(ruta, e.Name())
			}
		}
	}
	for _, nombre := range []string{"README.md"} {
		if ruta := p.rutaDoc(nombre); ruta != "" {
			agregar(ruta, nombre)
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]any{"docs": lista})
}

func (p *puente) apiDoc(w http.ResponseWriter, r *http.Request) {
	ruta := p.rutaDoc(r.URL.Query().Get("f"))
	if ruta == "" {
		http.NotFound(w, r)
		return
	}
	b, err := os.ReadFile(ruta)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write(b)
}

// ---- AGPL section 13: the source, offered to whoever uses the bridge ----
//
// The bridge is a network service, so the licence is not satisfied by shipping
// a LICENSE file: anyone who merely interacts with it over the network must be
// able to GET the corresponding source. That is what these two routes do, and
// why they are not decoration.

// paginaFuente is the human-facing offer.
const paginaFuente = `<!DOCTYPE html><html lang="es"><head><meta charset="utf-8">
<title>El código fuente — Laboratorio Diosyunalma</title>
<style>body{margin:0;background:#0b1526;color:#dce8f7;font-family:Georgia,serif;
padding:46px 22px;line-height:1.65}main{max-width:720px;margin:0 auto}
h1{color:#ffd166;font-size:23px}code{font-family:Consolas,monospace;color:#7fd7a8;font-size:13px}
a{color:#7fb2ff}.caja{border:1px solid #1d3a63;border-radius:12px;padding:18px 22px;margin:22px 0}
.dim{color:#8fa8c7;font-size:13px}</style></head><body><main>
<h1>El código fuente de este programa</h1>
<p>El puente de mando es <b>software libre</b> bajo la
<a href="https://www.gnu.org/licenses/agpl-3.0.html" target="_blank">GNU Affero General Public License v3</a>.
Se entrega SIN NINGUNA GARANTÍA.</p>
<p class="dim">Copyright © 2026 Jesús Nicolás Astorga y RESOURCES OPEN DOORS S.A.S</p>
<div class="caja">
<p><b>Descargá el fuente completo de esta aplicación:</b></p>
<p><a href="/fuente.zip">⬇ fuente-diosyunalma.zip</a></p>
<p class="dim">Incluye todo <code>cmd/</code>, el <code>go.mod</code> y las licencias:
es exactamente lo que hace falta para reconstruir este programa.</p>
</div>
<div class="caja">
<p><b>¿Necesitás usarlo sin abrir tu propio código?</b></p>
<p>Existe una <a href="/api/doc?f=LICENCIA-COMERCIAL.md" target="_blank">licencia comercial</a>
que te exime de las obligaciones de la AGPL.</p>
</div>
<p class="dim">La sección 13 de la AGPL exige que quien usa este programa por red
pueda obtener su fuente. Esta página es esa oferta, y el enlace de arriba la cumple.</p>
<p><a href="/">← volver al puente</a></p>
</main></body></html>`

func (p *puente) fuente(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(paginaFuente))
}

// fuenteZip streams the corresponding source of the running program.
func (p *puente) fuenteZip(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", `attachment; filename="fuente-diosyunalma.zip"`)
	z := zip.NewWriter(w)
	defer z.Close()

	agregar := func(rel string) {
		abs := filepath.Join(p.raiz, rel)
		b, err := os.ReadFile(abs)
		if err != nil {
			return
		}
		f, err := z.Create(filepath.ToSlash(rel))
		if err != nil {
			return
		}
		_, _ = f.Write(b)
	}

	for _, suelto := range []string{
		"go.mod", "LICENSE", "NOTICE", "LICENSE-CONTENIDO.txt",
		"LICENCIA-COMERCIAL.md", "LICENCIAS.md", "README.md",
	} {
		agregar(suelto)
	}
	// todo el arbol de fuentes Go
	_ = filepath.Walk(filepath.Join(p.raiz, "cmd"), func(ruta string, fi os.FileInfo, err error) error {
		if err != nil || fi.IsDir() || !strings.HasSuffix(ruta, ".go") {
			return nil
		}
		if rel, err := filepath.Rel(p.raiz, ruta); err == nil {
			agregar(rel)
		}
		return nil
	})
}

// ---- opening the helm ----
//
// The bridge must land on screen by itself: nothing is left to the captain's
// hand. A single opener can fail silently, so several are tried in turn and
// the arrival is CONFIRMED by watching for the first request to reach us.

// visitado turns true the moment a browser actually loads the dashboard.
var visitado atomic.Bool

func aperturas(url string) [][]string {
	switch runtime.GOOS {
	case "windows":
		return [][]string{
			{"cmd", "/c", "start", "", url},
			{"rundll32", "url.dll,FileProtocolHandler", url},
			{"explorer", url},
			{"powershell", "-NoProfile", "-Command", "Start-Process", url},
		}
	case "darwin":
		return [][]string{{"open", url}}
	default:
		return [][]string{{"xdg-open", url}, {"gio", "open", url}, {"sensible-browser", url}}
	}
}

// abrirTimon tries every opener until the dashboard is actually loaded, and
// says plainly what happened either way.
func abrirTimon(url string) {
	for i, ap := range aperturas(url) {
		if visitado.Load() {
			break
		}
		if i > 0 {
			fmt.Printf("   (el navegador no respondió; probando otra puerta…)\n")
		}
		_ = exec.Command(ap[0], ap[1:]...).Start()
		for espera := 0; espera < 24 && !visitado.Load(); espera++ {
			time.Sleep(250 * time.Millisecond)
		}
	}
	if visitado.Load() {
		fmt.Println("   ✓ timón en pantalla: el navegador ya está mostrando el laboratorio")
		return
	}
	fmt.Println("   ⚠ no pude abrir el navegador solo. Abrilo a mano y entrá a:")
	fmt.Println("     " + url)
}

// esRaiz recognises the laboratory by what it holds, not by how it was built:
// a distributed copy has no go.mod and no sources, only bin/ and the gallery.
func esRaiz(dir string) bool {
	for _, marca := range []string{"go.mod", "galeria", "bin", "cmd"} {
		if _, err := os.Stat(filepath.Join(dir, marca)); err == nil {
			return true
		}
	}
	return false
}

// hallarRaiz looks in the working directory first, then around the executable,
// so the bridge works both from the shop and from a double-clicked copy.
func hallarRaiz() (string, bool) {
	if cwd, err := os.Getwd(); err == nil && esRaiz(cwd) {
		return cwd, true
	}
	if exe, err := os.Executable(); err == nil {
		aqui := filepath.Dir(exe)
		for _, cand := range []string{aqui, filepath.Dir(aqui)} {
			if esRaiz(cand) {
				return cand, true
			}
		}
	}
	return "", false
}

func main() {
	raiz, ok := hallarRaiz()
	if !ok {
		fmt.Println("⚠ no encuentro el laboratorio: poné el puente en la carpeta que contiene galeria/ y bin/")
		os.Exit(1)
	}

	// modo generador: escribe el museo como pagina estatica y sale. Existe
	// porque GitHub Pages no puede correr un servidor, y el museo tiene que
	// poder verse desde la web igual que desde el puente.
	if len(os.Args) > 2 && os.Args[1] == "-museo" {
		destino := os.Args[2]
		if !filepath.IsAbs(destino) {
			destino = filepath.Join(raiz, destino)
		}
		if err := escribirMuseoWeb(raiz, destino); err != nil {
			fmt.Println("⚠ no pude escribir el museo:", err)
			os.Exit(1)
		}
		fi, _ := os.Stat(destino)
		fmt.Printf("🏛️  museo escrito: %s (%.0f KB)\n", destino, float64(fi.Size())/1024)

		// y de paso el muro de las maximas, que vive al lado del museo y se
		// arma de la misma fuente unica: cmd/puente/maximas.go
		mx := filepath.Join(filepath.Dir(destino), "maximas.html")
		if err := escribirMaximasWeb(mx); err != nil {
			fmt.Println("⚠ no pude escribir las máximas:", err)
			os.Exit(1)
		}
		fmi, _ := os.Stat(mx)
		fmt.Printf("🗿 máximas escritas: %s (%.0f KB · %d frases)\n", mx, float64(fmi.Size())/1024, len(maximas))
		return
	}

	p := nuevoPuente(raiz)

	porSala := map[string]int{}
	for _, e := range catalogo {
		porSala[e.Cat]++
	}
	sort.SliceStable(catalogo, func(i, j int) bool {
		oi, oj := ordenSala(catalogo[i].Cat), ordenSala(catalogo[j].Cat)
		if oi != oj {
			return oi < oj
		}
		return catalogo[i].Titulo < catalogo[j].Titulo
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		visitado.Store(true)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(pagina))
	})
	http.HandleFunc("/api/estado", p.apiEstado)
	http.HandleFunc("/api/salida", p.apiSalida)
	http.HandleFunc("/api/correr", p.apiCorrer)
	http.HandleFunc("/api/detener", p.apiDetener)
	http.HandleFunc("/api/lamina", p.apiLamina)
	http.HandleFunc("/api/laminam", p.apiLaminaMuseo)
	http.HandleFunc("/api/museo", p.apiMuseo)
	http.HandleFunc("/api/doc", p.apiDoc)
	http.HandleFunc("/api/docs", p.apiDocs)
	http.HandleFunc("/fuente", p.fuente)        // AGPL art. 13
	http.HandleFunc("/fuente.zip", p.fuenteZip) // el fuente correspondiente
	http.Handle("/galeria/", http.StripPrefix("/galeria/",
		http.FileServer(http.Dir(filepath.Join(raiz, "galeria")))))

	url := "http://localhost:" + puerto
	fmt.Println("⚓ EL PUENTE DE MANDO — Laboratorio Diosyunalma")
	fmt.Printf("   %d experimentos en %d salas\n", len(catalogo), len(salas))
	for _, s := range salas {
		fmt.Printf("   %s %-22s %d\n", s.Emoji, s.Nombre, porSala[s.Clave])
	}
	listos := 0
	for _, e := range catalogo {
		if rutaEjecutable(raiz, e.Dir) != "" {
			listos++
		}
	}
	switch {
	case listos == len(catalogo):
		fmt.Printf("\n   ⚙ los %d experimentos vienen compilados: esta máquina no necesita nada instalado\n", listos)
	case listos > 0 && hayGo():
		fmt.Printf("\n   ⚙ %d de %d compilados; el resto se compila al vuelo con Go\n", listos, len(catalogo))
	case hayGo():
		fmt.Println("\n   ⚙ modo taller: cada experimento se compila al vuelo con Go")
	default:
		fmt.Printf("\n   ⚠ solo %d de %d compilados y no hay Go instalado: los que falten no van a arrancar\n", listos, len(catalogo))
	}
	fmt.Println()

	ln, err := escuchar(puerto)
	if err != nil {
		fmt.Println("⚠ no pude abrir el puerto", puerto, "—", err)
		fmt.Println("  (alguien lo tiene tomado y no lo suelta; cerralo a mano y volvé a arrancar)")
		os.Exit(1)
	}

	fmt.Printf("   timón: %s   (Ctrl+C para amarrar)\n", url)
	go abrirTimon(url)
	if err := http.Serve(ln, nil); err != nil {
		fmt.Println("⚠ el puente se soltó —", err)
		os.Exit(1)
	}
}

func ordenSala(cat string) int {
	for i, s := range salas {
		if s.Clave == cat {
			return i
		}
	}
	return len(salas)
}
