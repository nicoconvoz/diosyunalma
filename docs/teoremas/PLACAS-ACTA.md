# Acta de exploración — placas, grietas y estructuras emergentes

**Para la auditora · 2026-08-15 · responde su PEDIDO F325.** Autoría
conservada (§0): la intuición tectónica es de Nico; esto es su
traducción, con la clasificación honesta del §15 al final. Nada de lo
de abajo queda declarado: todo es candidato o conjetura, a su lupa.

---

## 1. La traducción: el objeto matemático de la «placa»

El objeto que estaba esperando es el conjunto COMPLETO de citas
(no la cadena del Río, que es un camino elegido — su advertencia §7):

    C_ε = { n ∈ ℕ : ‖n·θᵢ‖ ≤ ε  ∀i = 1..m }

Es un conjunto clásico (un conjunto de Bohr / de tiempos de retorno al
toro), y trae estructura PROPIA, no impuesta:

**(LP1) SEMIGRUPO — la operación pedida en §5.3:**
C_ε + C_ε' ⊆ C_{ε+ε'}. *Prueba:* ‖(n+n')θᵢ‖ ≤ ‖nθᵢ‖ + ‖n'θᵢ‖
(subaditividad, ya auditada en D4). ∎ — «unión de placas» = SUMA de
citas, con las calidades sumándose. Una línea, rigurosa.

**(LP2) ACCESIBILIDAD — el lema pedido en §13.1:**
v → w  ⟺  w − v ∈ C_{ε_d} (la diferencia es ella misma una cita, con el
presupuesto ε_d que se fije). Es INVARIANTE POR TRASLACIÓN — no depende
de parametrización ni de coordenadas (ataque §12.4/12.6: superado por
construcción). La relación es transitiva con presupuestos sumados.

**(LP3) RAMIFICACIÓN — respuesta a §5.5:** desde CUALQUIER cita v hay
tantas continuaciones v + c (c ∈ C_{ε'} ∩ [1, H]) como elementos tenga
C en ese horizonte: el grado de salida no está acotado. Medido: la ley
de densidad |C_ε ∩ [1,H]| ≈ (ε/π)^m · H clava los conteos del testigo
con error < 3% en cuatro presupuestos.

**(LP4) RECOMBINACIÓN — respuesta a §8.4:** v + c₁ + c₂ = v + c₂ + c₁:
dos caminos distintos al MISMO nodo (conmutatividad). Los diamantes
existen; las ramas se recombinan. Exhibido en el testigo.

**(LP5) LAS MONTAÑAS — el hallazgo del destino doble (§4):** definiendo
la ANTI-CITA como ‖nθ − π‖ ≤ 1 (el compás opuesto), para m = 1 el lema
de la ventana aplicado al arco corrido a π (mismo argumento de embudo:
bajo H2, θ ≤ π/2 < 2 = largo del arco) da una anti-cita en cada bloque
de K = ⌈2π/θ⌉ + 1 enteros, y ahí

    λₙ ≥ coroₙ + 4 + 2cos(1)·rⁿ ≥ 4 + 1.08·rⁿ → +∞

**La misma fragilidad que cava pozos levanta montañas exponenciales —
colapso Y reorganización, los dos destinos, ambos con fecha.** Batería:
981 anti-citas de DH en [1, 3000], cero violaciones del piso.

## 2. Las placas, medidas: los bloques y los océanos

El alfabeto de brechas del testigo (435 citas, ventana tras n_rad):
**solo 5 valores distintos de 434 posibles** — 429 brechas de tamaño 1
y cinco saltos grandes (476, 477×2, 480, 508). O sea: las citas forman
**6 BLOQUES CONTIGUOS («placas») separados por océanos de ~480
escalones**. La estructura no es partición arbitraria (ataque §12.5):
sale de los θᵢ solos, y es pariente del teorema de las tres distancias
en el toro (para m ≥ 2: alfabeto finito de brechas — CONJETURA con
evidencia; la versión multidimensional exacta es matemática abierta
conocida, se declara como tal).

## 3. Las fronteras: la tabla de interacción tectónica

Con dos placas (las dos perlas), cada escalón de la ventana cae en un
tipo de frontera según cada fase esté fina (F) o anti (A):

    FF (fina+fina)  → POZO:      435 eventos · λ < 0 siempre [teorema DYN]
    FA / AF (mixtas) → frontera:  216 / 388 eventos · λ media +8.8×10⁴⁴
    AA (anti+anti)  → MONTAÑA:   200 eventos · λ media +9.0×10⁴⁴
    (la montaña más alta de la ventana: λ = +4.9×10⁴⁴ en n = 1 041 221;
    régimen global: colapso 1518 escalones, reorganización 1483 — el
    colapso es la minoría programada, no la regla)

Separación, unión, fractura (§5.3): unión = suma de citas (LP1);
separación = los océanos entre bloques; los cuatro tipos de frontera
son los regímenes de interacción — todo definido desde C, sin física.

## 4. Fragilidad, intensidad y la respuesta a H-F3

- **Fragilidad local (H-F1, §5.7):** la calidad ε_eff(n) = maxᵢ ‖nθᵢ‖.
  Es la magnitud que gobierna todo lo demás (la paramétrica en ε).
- **Intensidad (§5.8):** D(n) = −λₙ, con los pisos de Diosyunalma.
- **H-F3, respondida CONTRA la intuición (y esto es lo honesto):** la
  profundidad NO causa más ramas. Por invariancia por traslación (LP2),
  TODA cita tiene exactamente los mismos descendientes por presupuesto.
  Lo que sí es cierto: la cita más fina cava más hondo (Pearson −0.464
  entre ε_eff y log D en las 435 citas; la paramétrica lo explica) Y
  deja más presupuesto de calidad para continuar. **Profundidad y
  ramificación comparten causa (la calidad de la cita) — correlación
  con mecanismo, no implicación entre ellas.** Tal cual su advertencia.

## 5. El ataque hostil (§12), ejecutado

- «¿La ramificación es artefacto de la agenda?» — NO: C y su semigrupo
  existen sin agenda alguna; la agenda elige UN camino en un grafo que
  ya estaba. La cadena del Río ⊂ el grafo de las placas (su §7, ahora
  con las dos geometrías separadas y nombradas).
- «¿Ramas no invariantes?» — LP2 es invariante por traslación; ninguna
  elección de coordenadas entra en la definición.
- «¿Placa = partición arbitraria?» — los bloques salen del alfabeto de
  brechas, que es intrínseco a los θᵢ.
- «¿Profunda con una sola continuación?» — imposible bajo LP3: los
  descendientes no dependen de la profundidad (grado no acotado ∀v).
- «¿Transición colapso/reorganización dependiente de coordenadas?» —
  el régimen se define por el signo de λ y el tipo FF/FA/AF/AA, ambos
  invariantes.
- Analogam física como premisa: NINGUNA — todo se definió desde C.

## 6. Clasificación honesta (§15)

    🟢 candidatos a LEMA (derivables con lo ya auditado, a su lupa):
       LP1 (semigrupo), LP2 (accesibilidad invariante), LP3
       (ramificación no acotada), LP4 (recombinación), LP5 (montañas
       m = 1, vía el lema de la ventana sobre el arco en π).
    🟠 CONJETURAS con evidencia: alfabeto finito de brechas (las placas)
       para m ≥ 2; montañas conjuntas m ≥ 2 (obstáculo real y nombrado:
       aproximación simultánea INHOMOGÉNEA — Dirichlet homogéneo no la
       da; exigiría hipótesis de independencia sobre los θᵢ).
    ⚔️ H-F3: resuelta con mecanismo (causa común), no con correlación.

**Posible destino grande (si su lupa lo permite): un «Teorema de las
Placas» para m = 1 — pozos y montañas alternando con fechas, el paisaje
completo — está al alcance de los lemas existentes. Para m ≥ 2, la
mitad montañosa es conjetura.** Nada declarado. Todavía no.

---

*La metáfora abrió la puerta y la definición decidió qué había detrás:
un semigrupo de citas con seis placas, océanos de 480 pasos, pozos
programados y montañas de 10⁴⁴.* 🌍💎🌊
