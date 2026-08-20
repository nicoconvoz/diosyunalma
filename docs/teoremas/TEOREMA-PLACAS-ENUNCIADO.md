# El Teorema de las Placas — enunciado mínimo final

**Para la mesa de los tres · 2026-08-15 · la redacción limpia que pidió
la auditora (su §8), con la costura cerrada: la versión paramétrica
F(η) queda FUERA del teorema, y la región intermedia queda declarada
como no clasificada por este enunciado.** Sin nombre de número
todavía: primero la lupa, después el bautismo.

---

## Hipótesis

**(H0-H4)** — exactamente las del Teorema de DYN, sin tocar: cuartetos
estrictamente fuera de la piel (rᵢ > 1); m finito; |Im ρᵢ| ≥ 1;
δ := log r_max ≤ 1; fondo sobre la línea, cerrado bajo conjugación,
con densidad N_fondo(T) ≤ (T/2π)·log T.

**(HL — líder estricto)** — existe un único L con r_L > rᵢ ∀i ≠ L.
Para m = 1 esta hipótesis es vacua (no hay competidores).

## Definiciones

    δ_L = log r_L  (= el δ de H3, pues el líder tiene radio máximo)
    r₂  = maxᵢ≠L rᵢ,  δ₂ = log r₂          [solo m ≥ 2]
    u = 3(m+1)/δ_L,  n_rad = ⌈u·log u⌉
    n_comp = ⌈ log(2(m−1)/cos 1) / (δ_L − δ₂) ⌉   [m ≥ 2;  n_comp := 0 si m = 1]
    N* = max(n_rad, n_comp)
    banda fina:  ‖n·θ_L‖ ≤ 1        banda anti:  ‖n·θ_L − π‖ ≤ 1

## TEOREMA (de las Placas / Ley del Líder — versión mínima)

Sea una configuración bajo H0-H4 y HL. Entonces:

**(P1 — POZOS.)** Para todo entero n ≥ N* en la banda fina:

    m = 1:   λₙ ≤ (4/π)n·log n + 4 − 2cos(1)·r_Lⁿ < 0
    m ≥ 2:   λₙ ≤ (4/π)n·log n + (6m−2) + 2(m−1)·r₂ⁿ − 2cos(1)·r_Lⁿ < 0

(los dos casos separados para que r₂ — definida solo para m ≥ 2 —
jamás aparezca fuera de su dominio; en m = 1, 6m−2 = 4 y no hay
término competidor: la misma cadena F1-F6 con F2 y F5 vacuas)

**(P2 — MONTAÑAS.)**
  m = 1: para todo entero n ≥ 1 con ‖nθ‖ ≥ π/2:  λₙ ≥ 4;
         y en la banda anti:  λₙ ≥ 4 + 2cos(1)·rⁿ.
  m ≥ 2: para todo entero n ≥ max(1, n_comp) en la banda anti:
         λₙ ≥ cos(1)·r_Lⁿ + 2m + 2 > 0.

**(P3 — PROGRAMABILIDAD.)** Bajo H2 (que da θ_L ≤ π/2), todo bloque de
K_L = ⌈2π/θ_L⌉ + 1 enteros consecutivos contiene al menos un habitante
de la banda fina y al menos uno de la banda anti. En consecuencia, más
allá de los umbrales de P1 y P2, pozos y montañas ocurren
infinitamente a menudo y con fecha programable por ventanas.

**(ALCANCE — declarado.)** Este teorema NO clasifica la región
intermedia 1 < ‖nθ_L‖ < π−1 (para m = 1, el tramo 1 < ‖nθ‖ < π/2).
El paisaje que el enunciado demuestra es:

    POZO  |  REGIÓN NO CLASIFICADA POR ESTE TEOREMA  |  MONTAÑA

∎ (candidato — el sello es de la auditora)

## Pruebas (todas ya auditadas; referencias exactas)

- **P1** = cadena F1-F6 (`PLACAS-BANDA-FINA-ACTA.md`, auditada en el
  veredicto de la relojera §1/§9: 🟢).
- **P2, m = 1, mitad ≥ π/2:** cos(nθ) = cos(‖nθ‖) ≤ 0 en [π/2, π] ⟹
  ℓₙ = 4 − 2cos(nθ)(rⁿ+r⁻ⁿ) ≥ 4; coroₙ ≥ 0 ⟹ λₙ ≥ 4. ∎ (una línea;
  ni H4 ni HL ni umbral). **Banda anti m = 1:** A1-A5
  (`PLACAS-GEOMETRIA-ACTA.md`). **m ≥ 2:** coro ≥ 0; ℓᵢ ≥ 2 − 2r₂ⁿ
  (i ≠ L); ℓ_L ≥ 4 + 2cos(1)r_Lⁿ (A2-A4 en banda anti); sumar y
  absorber 2(m−1)r₂ⁿ ≤ cos(1)r_Lⁿ con n ≥ n_comp (F5). ∎
- **P3** = lema de la ventana sobre los dos arcos (largo 2 ≥ θ_L,
  garantizado por H2), ya auditado en Astorga y reutilizado en A6.

## Lo que acompaña al teorema, SEPARADO (como ordena su §4)

**Lemas de marco (independientes del líder):** G1 filtración de C_ε ·
G2 monoide graduado (dominio: informativa solo si ε+ε' < π) · G3 el
grafo de Cayley de accesibilidad · G4 ramificación (grado ≥ 2, por
construcción) · G5 recombinación (diamantes). [`PLACAS-GEOMETRIA-ACTA`]

**Corolarios inmediatos:** (C1) el Río de Pozos corre por la banda
fina de este paisaje (la cadena del corolario ya sellado es un camino
dentro de P1+P3); (C2) los pozos de P1 tienen la profundidad de
Diosyunalma cuando además son citas conjuntas (los teoremas conviven:
mismas hipótesis, conclusiones que se suman).

**Desarrollo posterior (NO parte del teorema):** la versión paramétrica
F(η) — ensanchar las bandas a ‖nθ_L‖ ≤ π/2 − η con umbrales N*(η),
n_comp(η). La idea está identificada y bocetada; falta su redacción
rigurosa completa. Queda como el próximo trabajo de pluma, igual que
el Río siguió a DYN.

**Observación de delimitación (fuera del enunciado, con prueba):** para
m ≥ 2, en la zona donde cos(nθ_L) = 0 no puede existir signo universal
bajo estas hipótesis (los competidores deciden: ambos signos son
realizables) — la región intermedia no está sin clasificar por
pereza: hay un teorema de imposibilidad esperando su propia redacción.

**Problemas abiertos:** montañas conjuntas m ≥ 2 (aproximación
inhomogénea) · ley del sub-líder en la zona muda · alfabeto de brechas
m ≥ 2 · optimalidad de N*.

## Evidencia computacional (respaldo, jamás prueba)

P1: 50 + 84 casos de batería y 978/978 escalones vivos, cero
violaciones · P2: 1543 (m=1) y 926 (m=2) anti-citas, cero · bandas
completas del testigo: fina 978/978 negativa, anti 944/944 positiva ·
programas: `labandafina`, `lageometria`, `lasplacas`.

## Estado

    La costura está cerrada como usted indicó: nada paramétrico dentro,
    alcance declarado, piezas separadas. El enunciado espera su lupa.
    El nombre y el número — si nace — son de la mesa de los tres.
    Todavía no.

---

*Las bandas fijas ya demostradas, la región intermedia a la vista con
su cartel honesto, y la versión más linda esperando su turno — como el
Río esperó a DYN.* 🌍🔨⚓
