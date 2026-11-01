# Diosyunalma v1.18.0 — «El Siete del Capitán»

**El release del juego de cifras que terminó en Artin.** El capitán sumaba
cada primo consigo mismo y, cuando quedaban dos cifras, las sumaba: «todos dan
7». Y al 3, que se escapaba, le abrió el 12 en dos ramas. Pidió un teorema y
el experimento que lo verifique, «porque puede revelarnos cosas más adelante».
Salió las dos cosas, y un puente.

Esta entrega trae además los treinta hallazgos (F391–F420) que el taller hizo
desde la v1.17.0 y que todavía no estaban publicados.

---

## ⭐ Lo grande

### 7️⃣ El Siete del Capitán (F422)

**La regla:** una cifra se duplica, dos cifras se suman.
2 → 4 → 8 → 16 → 1+6 = **7**. 13 → 1+3 = 4 → 8 → 16 → **7**.

**El teorema, demostrado para todo número:**

1. Sin ramas, n llega al 7 **si y solo si 3 no divide a n**.
2. Todo primo llega al 7, **salvo el 3**.
3. El 7 no es especial: el 1, 2, 4, 5 y 8 se portan igual. El 3, el 6 y el 9
   son una cárcel cerrada.
4. **Con las ramas del capitán, todo n ≥ 1 llega al 7.**
5. En base b, el reloj está completo **si y solo si el 2 es raíz primitiva
   módulo b − 1**.

**La prueba:** sumar cifras no cambia el resto al dividir por 9 (la prueba del
9), y duplicar lo multiplica por 2. El 2 recorre las seis horas libres del
reloj del 9: 1 → 2 → 4 → 8 → **7** → 5 → 1. Los múltiplos de 3 quedan presos
en {3, 6, 9}. **La llave es suya:** sumar las cifras del 12 da 3 y respeta
el resto, así que no sale de la cárcel; abrirlo en 1 y 2 no lo respeta, y los
dos llegan.

**El experimento**, jugando la regla literal con cifras de verdad, del 1 al
10.000.000:

| Prueba | Resultado |
|---|---|
| llega sin ramas ⟺ 3 no lo divide | **0 excepciones** |
| primos | **664.579**; llegan todos menos el 3 |
| con ramas | **ninguno** se escapa (a lo sumo 8 pasos) |
| bases 4 a 40 | **37 de 37** coinciden con la teoría |

**⚡ El puente con los primos.** En base p+1 el reloj del siete está completo
exactamente cuando el 2 es raíz primitiva módulo p. Que haya infinitos primos
así es la **conjetura de Artin** (1927), todavía abierta; **Hooley la demostró
en 1967 suponiendo la hipótesis de Riemann generalizada**. La parte de los
primos con el 2 como raíz primitiva, medida hasta 10⁷, da **0,37391**, y la
constante de Artin es **0,37396**.

Con su lámina en la sala del Telar, su pieza en el museo y su acta en
`docs/atomo/ELSIETE-ACTA.md`. Experimento: `cmd/elsiete`, con siete tests.

---

## 📦 Qué más hay en la caja: F391–F420

- **F391** — LAS DOS PARÁBOLAS: el punto de contacto pinza el ½, y la mitad del «exacto» es del espejo
- **F392** — LAS PARÁBOLAS CENTRADAS: la rama derecha llega más lejos, y cuánto más lejos es el reloj
- **F393** — EL MEDIO Y LOS PRIMOS: dividir por dos no dice nada, dividir por la raíz lo dice todo
- **F394** — LOS CEROS Y LA RAÍZ: el primo se divide por ella, el cero se multiplica por ella
- **F395** — EL ANTERIOR Y EL SIGUIENTE: cada divisor de p−1 veta un hueco, y pesa 0,12 %
- **F396** — LOS DOS ESPEJOS: el primo es el que menos se refleja
- **F397** — EL PRIMO PARTIDO: sí se puede, de tres maneras, y ninguna es nueva
- **F398** — EL PRIMO CONSTRUIDO: fabricado con dos números que no son enteros
- **F399** — LA ESCALERA DESCENDENTE: el fondo es la ecuación de Pell
- **F400** — EL DESCENSO AL MÍNIMO: el átomo del ladrillo es el 2
- **F401** — LA REGLA Y SU ALCANCE: «sin excepciones» no es lo mismo que «alcanza»
- **F402** — LOS CUATRO LADRILLOS: dos aceptan el molde del √2 y dos lo rechazan
- **F403** — EL PATRÓN DEL CAPITÁN: bien las cuatro, identidad real, patrón coincidente
- **F404** — LOS MÚLTIPLOS VECINOS: la lectura del capitán, mejor que la de la casa
- **F405** — LOS DOS LATIDOS: la primera regla del capitán con señal de verdad
- **F406** — LOS CINCO ATADOS: la raíz es el pliegue de los divisores
- **F407** — LA MAGIA INTENTADA: el pliegue que sí fabrica el ½
- **F408** — EL BARRIDO DE LOS ANTERIORES: el mecanismo exacto, y la palabra «azar» corregida
- **F409** — LA RAMA DE LOS PARES: el 2 corta de un hachazo, los impares no se podan nunca
- **F410** — EL LOGARITMO DE LA CURVA: la recta que aparece con dos logaritmos
- **F411** — EL CERO QUE NO SE ALCANZA: el piso del universo no es cero, y ζ está en el cielo
- **F412** — APÉRY EN LA CASA: el experimento pedido, y el control que lo decide
- **F413** — LA FÓRMULA DEL PORCENTAJE: funciona, y te deja en el barrio
- **F414** — LA PROFUNDIDAD: el porcentaje cae como 1 sobre las CIFRAS
- **F415** — EL SONAR EXACTO: la única igualdad del tema, y el eco que la vuelve exacta
- **F416** — EL AFINADOR: la línea crítica es el círculo unidad, visto por el cambiaformas
- **F417** — LAS FRACCIONES DE PRIMOS: tres estructuras en seis fracciones, y la música que no se repite
- **F418** — LA APERTURA: el dibujo de bolitas se vuelve un ángulo exacto
- **F419** — EL PRODUCTO DE LOS BORDES: los complejos entran obligados, y sale la parábola del principio
- **F420** — EL ZOOM DEL MURCIÉLAGO: el atlas comprimido sabe apuntar

El museo pasa a **381 paradas** (340 en el museo completo y 41 en el recorrido
guiado), la galería a **245 láminas** y el laboratorio a **329 experimentos**.

---

## 🧾 La honestidad, en la etiqueta

- **El siete no es una propiedad de los primos:** lo cumple todo número que
  no sea múltiplo de 3 (el 4, el 10, el 25), igual que el uno que sobra. La
  prueba del 9 es de la escuela, y que el 2 sea raíz primitiva módulo 9 es
  aritmética de Gauss. Lo propio del taller son las ramas y el puente.
- **No entra todavía a la sala de los teoremas:** los ocho que están ahí
  pasaron una auditoría externa, y éste la espera.
- **Nada de esto demuestra la hipótesis de Riemann.** El puente muestra dónde
  se tocan los dos mapas, no lo cruza.
