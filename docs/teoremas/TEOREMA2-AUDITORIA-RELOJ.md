# Auditoría del reloj — matemática vs. implementación

**Para la auditora · 2026-08-15 · respuesta al pedido «F319 — auditoría
reloj: matemática ↔ código».** Objeto: `cmd/elmecanismo` contra las actas
(`TEOREMA2-LEMA-INTERACCION-ACTA.md`, F299-F301, F318, F319). Regla
aceptada en su totalidad: **una corrida con cero fallos NUNCA sustituye
una demostración universal** — cada salida del programa queda
clasificada abajo en una de las cuatro categorías.

**Resultado adelantado:** un hallazgo real, cazado y corregido DURANTE
esta auditoría (§7-A: el chequeo de H4 verificaba 6 instancias y se
presentaba como verificación de la ventana — ahora chequea los 38 puntos
de salto, que por monotonía SÍ es la verificación completa). Después de
esa corrección: ninguna pieza de código hace algo distinto de lo que
afirma su lema. Todo lo que el programa no cubre está marcado en §6.

---

## 1. El diccionario — cada flecha del mapa

Categorías: **[U]** prueba matemática universal (vive en el acta) ·
**[F]** verificación computacional finita · **[E]** input externo ·
**[I]** ejemplo numérico concreto.

| Flecha | Lema que la justifica | Código que la implementa | Hipótesis | ¿Instancia o desigualdad general? | Queda sin verificar automáticamente |
|---|---|---|---|---|---|
| H0-H4 → (configuración admisible) | definiciones §0 y §3 del acta | bloque `§H`: r₁,r₂,δ por fórmula; comparaciones directas; H4 en los 38 saltos γₖ | — | H0-H3: **[I]** sobre esta configuración. H4: verificación COMPLETA de la ventana (§7-A) | que la configuración elegida sea representativa (no lo necesita: el teorema es ∀ configuración) |
| → L1 (golpe) | L1-L7 del acta (§1) | bucle L1: `slack = 2m+2 − e^{nδ} − (ℓ₁+ℓ₂)` en cada cita | citas con ε=1 | **[F]**: 20239 citas reales, la desigualdad L7 instanciada en ε=1 | L7 universal en ε y en m: **[U]** en el acta |
| → L2 (Dirichlet) | palomar en el toro (§2a) | bucle L2: 50 pares al azar × Q ∈ {5,10,20}, busca n ≤ Q² | m=2 | **[F]**: el experimento; la PRUEBA del palomar es **[U]** | el teorema general de Dirichlet — solo acta |
| → L3 (agenda) | lema de agenda (§2b), desigualdades (i)-(iv) | bucle L3: 2000 pares (T,n₁) al azar, chequea J≤T, n≥T, n≤T+n₁, J·2π/Q≤1 | T≥1, n₁≥1 | **[F]** del esqueleto aritmético | la subaditividad ‖J·n₁θ‖ ≤ J·‖n₁θ‖ — solo acta (§6) |
| → L4 (radial) | R0-R10 (§2d) | bucle L4: R2 (1.094), R3 (2.06), R5, R6 (u²≥3(log u)²+1), R7, R8, R9 evaluadas en la grilla | m=1..10, δ≤1 | **[F]**: 50 instancias de afirmaciones cuyo carácter universal lo prueban q(u) y H(u) en el acta | R6/R8 universales vía q,H; R10 (propagación por monotonía) es puramente lógica: sin código |
| → L5 (coro) | F299-F301 (cota (4/π)n·log n) | bucle L5: coro recursivo + comparación directa contra la cota, n ∈ [3, 1.1×10⁶] | H4, n≥3 | **[F]** de la CONCLUSIÓN (1.1M instancias); la cadena conteo+cola+absorción es **[U]** (F299-301, auditada en F319) | la implicación universal H4 ⟹ cota — solo acta |
| → L6 (λₙ < 0) | ensamble §3 del acta | bloque L6: n_rad,m, búsqueda de cita ≥ n_rad, λ = coro + ℓ₁ + ℓ₂ | todas | **[I]**: UN testigo (n = 1040809) de la existencia que el acta garantiza ∀n ≤ N₀ | el cuantificador ∃n ≤ N₀ en el peor caso — solo acta |
| → M_{n,n} < 0 | diagonal del yunque: M_{n,n} = 2λₙ (λ₀ = 0) | **NO está en este programa** (implementada en `cmd/elyunque`, F292) | — | interpretación | marcada: paso solo-escrito respecto de este programa |
| aplicación a ζ | Parte C | **NO está en el programa** — B1, B2 son **[E]** | — | — | correctamente fuera: el programa es Parte A |

## 2. Correspondencia fina, eslabón por eslabón

**Norma circular.** `mod2pi(x)` = reducción mod 2π + plegado a [0, π] =
exactamente ‖x‖ (distancia al múltiplo de 2π más cercano). Es la norma
de las actas. Detalle declarado: el código usa `< 1` estricto donde el
acta escribe ≤ ε — dirección CONSERVADORA (toda cita que el código
acepta es cita del acta).

**ε = 1 cableado.** El programa instancia ε = 1 en L1 y L6. L7 es
universal en ε en el acta; el código no lo verifica para otros ε.
Declarado como instancia.

**El fondo sobre la línea.** Las perlas salen de ceros de Z(t) con t
real (sobre la línea por construcción) y el código NORMALIZA w := w/|w|
— eso fuerza |w| = 1 exacto, que es la hipótesis del coro. La distancia
entre la fase medida y φ = 2·arctan(1/2γ) se mide: 7×10⁻¹⁸.

**L7 en ε = 1.** El acta: Σℓ ≤ 2ε²(m−1) + 4 − 2(1−ε²/2)·r_maxⁿ. Con
ε = 1: 2(m−1) + 4 − r_maxⁿ = 2m+2 − r_maxⁿ. El código compara contra
`2m+2 − e^{nδ}` con δ = log r_max: la MISMA expresión. ✓

**L4: las constantes del código contra las del acta**, una por una:
1.094 = R2 ✓ · 2.06 = R3 ✓ · 2m+2 ≤ u = R5 ✓ · u² ≥ 3(log u)²+1 = R6
(el lemita, con su prueba vía h y q en el acta) ✓ · g > 0 = R7 ✓ ·
g' > 0 = R8 ✓ · g'' > 0 = R9 ✓. R9 se evalúa en n = n_rad; su validez
∀n ≥ n_rad la da la monotonía (acta). R10 no genera cálculo: es lógica.

## 3. Auditoría especial de L5 (las cinco preguntas)

① **¿El mismo término por-par?** Sí: la recursión calcula
2 − 2·Re(wⁿ) = 4·sin²(nφ/2); verificado contra la fórmula directa
4·sin²(nφ/2) en 4 puntos de control: desviación máxima **2.8×10⁻¹⁰**
(medida, no asumida).
② **¿Usa realmente H4?** Con precisión: el programa VERIFICA que este
fondo satisface H4 (los 38 saltos — corregido en §7-A) y verifica la
CONCLUSIÓN de L5 directamente. La implicación H4 ⟹ cota no se
"ejecuta": es la prueba universal de F299-F301. El paso intermedio
(conteo + cola) sí está implementado aparte, en el test adversarial de
F319 (fondo de densidad máxima H4: no llega al 60% de la cota).
③ **¿El 4/π por la misma cadena?** En el programa el 4/π aparece solo
como cota final. La cadena 2/π + 1/π + 1/π es solo-acta: **marcada**.
④ **¿Comprobación finita o cota universal?** Comprobación FINITA de
1.1×10⁶ instancias de una cota cuya universalidad (∀n ≥ 3) es teorema
escrito. Nunca al revés.
⑤ **¿Precisión que oculte una violación?** No puede: el margen mínimo
medido es **87.7% de la cota** (en n = 56; el peor caso relativo de
toda la corrida) y la deriva medida es 2.8×10⁻¹⁰ — una violación
tendría que esconderse en una rendija 10¹² veces más chica que el
margen. Además los términos son no negativos: sin cancelación posible
dentro del coro.

## 4. Auditoría especial de L6

**Separación exigida — y es aún más limpia de lo que parecía:**

- El **ejemplo concreto** n₀ = 37306 (primera ruptura real) **NO es una
  cita**: a 50 dígitos, ‖n₀·θ₁‖ = 1.762 > 1 (y ℓ₁ = +5.90 ahí, positiva).
  La ruptura empírica va por FUERA de la ruta del teorema — es dato, no
  conclusión. El cruce está clavado: λ(37305) = +0.3193,
  λ(37306) = −0.0321.
- El **testigo del teorema** es la cita n = 1040809 (primera con ε = 1
  después de n_rad,m = 1040809): ‖nθ₁‖ = 0.8042, ‖nθ₂‖ = 0.7247 — ambas
  < 1 con margen ~0.2 (contra una tolerancia de fase de ~3×10⁻¹⁰: la
  condición de cita no depende del borde). Ahí λ = −6.49607464×10⁴⁴ < 0:
  la conclusión que el acta garantiza para ALGÚN n ≤ N₀, exhibida.
- **¿n₀ satisface lo necesario para L5?** Sí: n₀ = 37306 ≥ 3 (umbral del
  coro) y la cota del coro rige en él; L6 (la ruta por citas) no se le
  aplica ni se le atribuye.
- **¿El signo depende de redondeo?** No. Contra-cálculo independiente a
  50 dígitos con los ceros EXACTOS de mpmath (no los del scanner Go):
  λ(1040809) = −6.49607464285×10⁴⁴ — coincide con el float64 del Go en
  todos los dígitos mostrados; λ(37306) = −0.03211271, margen 0.032
  contra un presupuesto de error float64 de ~10⁻⁶ en esa escala.
- **M_{n,n} = 2λₙ < 0**: interpretación del acta (λ₀ = 0 ⟹ diagonal
  = 2λₙ); implementada en `cmd/elyunque` (F292), no en este programa.
  Marcada como tal.

## 5. Auditoría de precisión numérica

- **Tipo:** float64 (IEEE 754 doble, 53 bits, ~15.9 dígitos). El
  contra-cálculo de control: mpmath a 50 dígitos.
- **Fases:** el argumento n·θ ~ 3×10⁶ tiene error absoluto ~3×10⁻¹⁰;
  `math.Mod` es exacto sobre su entrada. Las decisiones de cita tienen
  margen ≥ 0.2 en los casos que importan: 9 órdenes por encima del ruido.
- **Exponenciales:** exponente máximo n·δ = 108.6 ≪ 709 (sin overflow);
  e^{−nδ} → 10⁻⁴⁸ sin consecuencia (suma un 0 donde iba un 10⁻⁴⁸).
- **Cancelación catastrófica:** no hay. En n = 1040809 los dos términos
  gigantes (ℓ₁ = −1.3×10¹⁹, ℓ₂ = −6.5×10⁴⁴) tienen el MISMO signo. El
  único λ cercano a cero (n₀) quedó re-derivado a 50 dígitos.
- **Márgenes MEDIDOS, no asumidos:** golpe: slack mínimo 1.071 (en
  n = 4762) contra ruido ~10⁻¹⁵-10⁻⁸ según escala · coro: margen mínimo
  87.7% · deriva de recursión: 2.8×10⁻¹⁰. Ninguna desigualdad de la
  corrida está cerca del límite numérico.
- **Tolerancias declaradas:** 10⁻⁹ en L1 y 10⁻¹² en L3 — ambas quedaron
  irrelevantes: los márgenes reales son ≥ 1.07 y ≥ estrictos por enteros.

## 6. Cobertura — qué cubre el programa y qué solo las pruebas escritas

**Cubierto por el programa:** hipótesis H0-H4 sobre la configuración
viva (H4 con verificación completa de ventana); instancias masivas de
L1, L3, L4, L5; el testigo de L6 con su λ; los márgenes numéricos.

**Cubierto SOLO por las pruebas escritas (huecos declarados, no
defectos):** ① la prueba del palomar (L2); ② la subaditividad de la
norma circular en L3; ③ la universalidad de TODAS las baterías (en ε,
en n, en configuraciones); ④ R10; ⑤ el cuantificador ∃n ≤ N₀ del peor
caso (el programa exhibe un testigo temprano; N₀ = 4.3×10¹³ no es
corrible); ⑥ el paso M_{n,n} = 2λₙ (otro programa); ⑦ toda la Parte C
(ζ, B1, B2) — correctamente ausente: el programa es Parte A.

## 7. Intento activo de romper el reloj (los cinco ataques)

- **A. Caso que cumple las hipótesis pero el programa trata mal —
  ENCONTRADO Y CORREGIDO:** el chequeo de H4 original evaluaba
  N_f(T) ≤ (T/2π)log T en 6 valores de T y se reportaba como «H4 en la
  ventana». Eso era una instancia vestida de verificación — exactamente
  la confusión que esta auditoría persigue. Corregido: como N_f solo
  salta en los γₖ y la cota es creciente, chequear los 38 saltos ES la
  verificación completa. (El resultado no cambió — true antes y después
  — pero el ESTATUS del chequeo sí: de [I] a verificación completa.)
- **B. Falso positivo/negativo:** análisis de dos regímenes en L1 — en
  n chicos el slack puede ser fino pero el ruido es 10⁻¹⁵; en n grandes
  el ruido crece a ~10⁻⁸ pero el slack crece como 0.08·r_maxⁿ. Medido:
  slack mínimo global 1.071. Sin zona ciega.
- **C. Dependencia del orden de suma:** imposible en el coro (términos
  ≥ 0, orden fijo, deriva medida 2.8×10⁻¹⁰); λ suma 3 términos.
- **D. Cuantificadores código vs actas:** dos diferencias, ambas
  declaradas y conservadoras — `<` estricto vs ≤, y ε = 1 instanciado.
  Ninguna en la dirección peligrosa.
- **E. Overflow/underflow/cancelación:** cuantificados en §5 — ninguno.

## 8. Veredicto solicitado (autoevaluación; el grado final es de la auditora)

    Tras la corrección §7-A:
    🟢 CORRESPONDENCIA COMPLETA — cada cálculo del programa implementa
    exactamente la afirmación matemática que dice implementar, con las
    limitaciones DECLARADAS (§6): las pruebas universales viven en las
    actas y el programa nunca las sustituye; exhibe hipótesis
    chequeadas, instancias masivas, márgenes medidos y un testigo.
    Un hallazgo real (H4 instancia-vs-ventana) cazado y corregido por
    esta misma auditoría — el reloj del plano y el reloj que gira son
    el mismo reloj, y ahora también lo es su hoja de inspección.

---

*El plano dice cómo debería funcionar el reloj. El código lo hace
girar. Esta acta comprobó, tornillo por tornillo, que son el mismo
reloj — y donde la hoja de inspección exageraba (H4), se corrigió la
hoja.* ⚙️🕰️🔍
