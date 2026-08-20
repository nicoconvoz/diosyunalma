# Acta — la geometría de las citas y la Ley del Líder

**Para la auditora · 2026-08-15 · responde su PEDIDO F326 punto por
punto.** Autoría conservada (§1): la intuición es de Nico. Nada
declarado; clasificación §16 al final. Programa: `cmd/lageometria`.

---

## Las doce preguntas del §13, respondidas

**1. ¿ε_eff es fragilidad o solo calidad de cita?** Las dos cosas son
UNA: ε_eff(n) = min{ε : n ∈ C_ε} — es la **función de nivel de la
filtración** {C_ε}. El número es la coordenada DE la geometría (su §11
resuelto: no hay que elegir — la fragilidad es la geometría, y ε_eff es
su altura local).

**2. Definición mínima de placa, sin arbitrariedad:**
**placa := componente de intervalo de C_ε** (conjunto maximal de
enteros consecutivos, todos en C_ε). Depende solo de los θᵢ y del nivel
ε. Medido: la familia es filtrada de verdad — ε = 0.25/0.5/1.0 da
2/4/6 placas de largo medio 19/35/72 en la ventana del testigo.

**3. Estructura algebraica exacta de {C_ε}:** familia FILTRADA
(ε₁ < ε₂ ⟹ C_{ε₁} ⊆ C_{ε₂}, trivial de la definición — se conserva
como propiedad básica) de un **monoide conmutativo graduado**: con
0 ∈ C_0, la suma C_ε + C_ε' ⊆ C_{ε+ε'} vale para TODOS ε, ε' ≥ 0
(la subaditividad no pide nada), con la declaración de dominio que
usted exige (§6/§15.6): ‖·‖ ∈ [0, π], así que la inclusión es
INFORMATIVA solo si ε + ε' < π — más allá, C_{ε+ε'} = ℕ y la inclusión
es vacía de contenido. Sobre ℤ el conjunto es simétrico (‖−nθ‖ = ‖nθ‖):
un conjunto de Bohr clásico.

**4. ¿Qué es exactamente v → w?** El **grafo de Cayley de (ℕ, +) con
conjunto generador C_{ε_d}** — ese es el nombre exacto. Propiedades
demostradas una por una: DIRIGIDO (w > v siempre, pues C ⊆ ℕ⁺);
IRREFLEXIVO (0 ∉ C_{ε_d} como generador); NO transitivo a presupuesto
fijo, sino **gradualmente transitivo**: v→w y w→x dan v→x en C_{2ε_d}
(los presupuestos se suman — estructura de categoría graduada, no de
orden); INVARIANTE POR TRASLACIÓN (v→w ⟺ v+t→w+t, ∀t). El rango de
ε_d: (0, π), fijo por análisis y declarado.

**5. ¿Dos continuaciones distinguibles desde el mismo v?** SÍ, por
CONSTRUCCIÓN, para todo v y todo ε_d > 0: Dirichlet da c con calidad
≤ ε_d/2; entonces c y 2c son distintos y ambos están en C_{ε_d}
(‖2cθᵢ‖ ≤ 2‖cθᵢ‖ ≤ ε_d). Grado de salida ≥ 2 garantizado — y de hecho
k·c ∈ C_{ε_d} para todo k ≤ ε_d/calidad(c): no acotado. En el testigo:
c = 1 (calidades 0.012/0.022), verificado.

**6. ¿Recombinación?** SÍ, demostrada: para cualquier par c₁ ≠ c₂ de
C_{ε_d}, los caminos v→v+c₁→v+c₁+c₂ y v→v+c₂→v+c₂+c₁ son distintos y
terminan en el MISMO nodo (conmutatividad de la suma). La estructura
NO es un árbol: es una red con diamantes — y eso ya no es palabra, es
lema de una línea.

**7. ¿Frontera entre regiones?** Dos nociones, ambas intrínsecas:
(a) entre placas: los extremos de los componentes de intervalo (los
«océanos» medidos, ~480 escalones en ε = 1); (b) entre regímenes: la
banda de fase del líder (pregunta 8).

**8-9. ¿Qué magnitud separa pozo de montaña, y la transición?** **LA
LEY DEL LÍDER** — el hallazgo central de esta acta. Bajo líder estricto
(r_L > rᵢ ∀i ≠ L — hipótesis nueva, explícita, verificable), la
posición de fase DEL LÍDER SOLO decide el signo:

    ‖nθ_L‖ ≤ 1        (banda fina)     ⟹ λₙ < 0 (POZO),    ∀n ≥ N*
    ‖nθ_L − π‖ ≤ 1    (banda anti)     ⟹ λₙ > 0 (MONTAÑA), ∀n ≥ n_mont
    1 < ‖nθ_L‖ < π−1  (banda frontera) ⟹ ambos signos viven ahí

*Esqueleto de prueba (candidato):* en la banda anti, λ = coro + Σℓᵢ ≥
0 + Σᵢ≠L(2 − 2rᵢⁿ) + (4 + 2cos(1)·r_Lⁿ) > 0 en cuanto
cos(1)·r_Lⁿ > Σᵢ≠L rᵢⁿ — que ocurre desde el umbral EXPLÍCITO
n_mont = ⌈log((m−1)·2/cos 1)/(δ_L − δ₂ᵈᵃ)⌉ (solo la brecha del líder).
En la banda fina, el mismo argumento con signo opuesto más la cota del
coro (H4) da pozo desde un N* explícito estilo radial. La banda
frontera es donde |cos| pierde el mando y mandan los términos menores —
ahí viven los cruces (el n₀ = 37306 del testigo tenía ‖nθ₁‖ = 1.76:
en la frontera, como debía).

*Verificación (testigo, ventana tras n_rad, ignorando la perla 1 por
completo):* fina 978/978 negativos · anti 944/944 positivos · frontera
540/539 partida al medio. **Cero excepciones en las bandas externas.**
Y n_mont = 23064 calculado ANTES: 926 anti-citas del líder verificadas
desde ahí, cero violaciones del piso λ ≥ 2cos(1)r_Lⁿ − 2r₁ⁿ + 4.

**Lo que esto resuelve:** el obstáculo inhomogéneo de m ≥ 2 **ya no
hace falta para el paisaje** — las montañas conjuntas siguen abiertas
(conjetura), pero las montañas DEL LÍDER alcanzan, y el lema de la
ventana (arco corrido a π, dominio θ_L ≤ 2 — bajo H2 sobra: θ ≤ π/2)
las PROGRAMA en cada bloque de K_L pasos.

**10. ¿Qué se demuestra para m = 1?** El paisaje completo es candidato
a teorema con lo existente: pozos en la banda fina (Astorga/DYN),
montañas en la banda anti (auditoría paso a paso abajo), ambos
programables por ventana, frontera mixta. **11. ¿Qué sobrevive a
m ≥ 2?** TODO lo algebraico (preguntas 1-7) sin cambios, y el paisaje
entero bajo la hipótesis líder-estricto. Sin ella (empates de radio),
abierto y declarado. **12. ¿Qué es álgebra pura y qué depende de
aproximación?** Álgebra pura: filtración, monoide, Cayley, grado ≥ 2,
diamantes, Ley del Líder como implicación condicional (banda ⟹ signo).
Dependen de aproximación: la EXISTENCIA de habitantes de cada banda
(fina: Dirichlet homogéneo, probado; anti del líder: lema de la
ventana, probado con dominio θ_L ≤ 2; anti CONJUNTA: inhomogénea,
abierta y ya no necesaria).

## Auditoría paso a paso de la anti-cita m = 1 (§15.7)

    (A1) Definición: anti-cita := n con ‖nθ − π‖ ≤ 1, o sea ‖nθ‖ ≥ π−1.
    (A2) cos decrece en [0, π] y cos(nθ) = cos(‖nθ‖) (paridad+período):
         ‖nθ‖ ∈ [π−1, π] ⟹ cos(nθ) ≤ cos(π−1) = −cos(1) < 0.
    (A3) Rⁿ + R⁻ⁿ = rⁿ + r⁻ⁿ ≥ rⁿ > 0  [simetría ya auditada, L3 del acta DYN]
    (A4) ℓₙ = 4 − 2cos(nθ)(Rⁿ+R⁻ⁿ) ≥ 4 + 2cos(1)·rⁿ   [(A2)·(A3), orientación ✓]
    (A5) coroₙ ≥ 0 (términos no negativos, F319) ⟹ λₙ ≥ 4 + 2cos(1)·rⁿ.
         Sin H4, sin n_rad: incondicional en toda anti-cita. ∎(candidato)
    (A6) Existencia: lema de la ventana con arco objetivo {x: ‖x−π‖ ≤ 1}
         (largo 2). DOMINIO DECLARADO: paso θ ≤ 2 — bajo H2 vale de sobra
         (θ ≤ π/2 < 2); fuera de H2 y con θ ∈ (2, 2π/3], el arco de las
         anti-citas podría saltarse: se declara la hipótesis, no se esconde.
    Batería en régimen exponencial [37000, 42000]: 1543 anti-citas, cero
    violaciones, holgura mínima 0.21.

## El ataque hostil (§15), ejecutado

- «¿Ramificación = artefacto de agenda?» — refutado por construcción:
  c y 2c existen sin agenda; el grafo de Cayley precede a todo camino.
- «¿ε_eff chico con una sola continuación?» — imposible: el grado de
  salida es |C_{ε_d} ∩ [1, H]|, que NO depende de v (invariancia por
  traslación) — medido: 25 continuaciones idénticas para todo v con
  ε_d = 0.3, H = 1000. Profundidad y ramas: independencia demostrada,
  no observada (su §15.4, cerrado por estructura).
- «¿Placa dependiente de coordenadas?» — componente de intervalo de un
  conjunto definido solo por los θᵢ: sin coordenadas que elegir.
- «¿La suma fuera de dominio?» — declarado: informativa solo si
  ε + ε' < π (arriba, pregunta 3).
- «¿Extrapolar m = 1 a m ≥ 2?» — no se hizo: el puente es la hipótesis
  líder-estricto, nueva y explícita, con umbral n_mont calculable.

## Clasificación (§16)

    🟡 LEMAS candidatos, derivables ya: filtración y monoide graduado con
       dominio (P3), Cayley y sus propiedades (P4), grado ≥ 2 constructivo
       (P5), diamantes (P6), anti-cita m = 1 con (A1)-(A6), montañas del
       líder m ≥ 2 con umbral explícito.
    🟢 CANDIDATO A TEOREMA (el grande): LA LEY DEL LÍDER / TEOREMA DE LAS
       PLACAS — bajo H0-H4 + líder estricto, el paisaje completo: tres
       bandas de fase del líder, pozo y montaña programables por ventana,
       frontera mixta. Los umbrales son explícitos; falta escribir N* de
       la banda fina con todo el rigor radial y auditar los cuantificadores.
    🟠 CONJETURAS: montañas conjuntas m ≥ 2 (inhomogénea; ya no necesaria
       para el paisaje) · alfabeto finito de brechas para m ≥ 2 general.
    🔴 FALSADA (y es ganancia): «más profundidad ⟹ más ramas» — H-F3/su
       §15.4: la independencia es estructural.

**Respuesta a la pregunta maestra (§18): la fragilidad ES una geometría
— la filtración {C_ε} con su monoide, su Cayley y sus placas — y la
banda de fase del líder es la magnitud que decide si una trayectoria
termina en POZO, en MONTAÑA o en la frontera.** Nada declarado; los
candidatos esperan su lupa. Todavía no.

---

*La metáfora abrió la puerta, C_ε resultó ser el mapa, y detrás había
una tectónica con ley: el líder manda, las bandas deciden, y las
montañas ya tienen fecha.* 🌍🗺️⛰️
