# Diosyunalma v1.18.1 — «El Encuentro»

**El release de la regla completa.** El v1.18.0 publicó el Siete del Capitán
con una parte de su regla salteada. El capitán la marcó el mismo día: «te
salteaste una parte de la regla». Esta entrega la agrega, la demuestra, la
verifica, y mete al siete adentro de la sala de los teoremas.

---

## ⭐ Lo grande

### 🤝 El encuentro: la segunda mitad de la regla del capitán

> «9 = 9+9 = 18, acá se divide, 1+1 = 2, 2+2 = 4, 4+4 = 8; el otro ya
> estaba en 8 cuando lo alcanza; se encuentran y siguen juntas, 8 de una más
> 8 de otra, 16, y finalmente 7.»

**La regla completa:**

1. Una cifra se duplica. Dos cifras se suman, salvo que la suma caiga en la
   cárcel (3, 6 o 9): ahí se abren en ramas.
2. Avanza la rama más chica; las demás esperan.
3. Dos ramas en el mismo número se juntan: se suman y siguen como una.
4. Termina cuando todas llegaron al 7.

**La ley del encuentro, demostrada:** con la regla completa, todas las ramas
de todo número se encuentran y terminan juntas en el 7.
- Las ramas libres viven todas en la misma vuelta del reloj del 9
  (1 → 2 → 4 → 8 → 7 → 5), cada una a un número fijo de saltos del 7.
- Moverse acerca a la que se mueve y no aleja a ninguna, y juntar dos iguales
  es lo mismo que duplicar.
- Lo que falta por recorrer solo baja: termina, y termina en el 7.

**Verificada**, jugando la regla completa del 1 al 10.000.000:

| Prueba | Resultado |
|---|---|
| números que no terminan con todas sus ramas juntas en el 7 | **0** |
| números donde las ramas se juntaron | **3.333.333**: exactamente los múltiplos de 3 |
| el camino más largo | 34 pasos y 10 juntadas |

El 9 del capitán sale tal cual lo escribió:
`9 → 18 → ramas 1 y 8 → 2 → 4 → 8 se junta con 8: 16 → 7`.

### 🏛️ El siete entra a la sala de los teoremas

Como **Teorema 8 del taller**, con su página propia, su tarjeta en la
galería, su pieza en la sala del museo (que pasa a nueve paradas) y su
sección en el libro de teoremas. Entró sin la auditoría externa de sus
hermanos, por decisión del capitán: «es para eso, para que los que lo vean lo
auditen». En cada lugar se dice de frente.

---

## 🧾 La honestidad, en la etiqueta

- **Una conjetura propia, refutada por el experimento:** yo había anunciado
  que, en otra base, las ramas se encuentran solo cuando el 2 es raíz
  primitiva. En base 5 se juntan igual, porque abrir los números vuelve a
  mezclar las vueltas, y en algunas bases las ramas se multiplican sin
  terminar. El encuentro fuera de la base 10 queda como pregunta abierta.
- **Nada de esto es una propiedad de los primos** ni demuestra la hipótesis
  de Riemann: el puente con Artin muestra dónde se tocan los dos mapas.
