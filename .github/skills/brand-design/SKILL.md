---
name: brand-design
description: Guía visual para diseñar interfaces web coherentes con la identidad de AulaViva.
---

# Skill de diseño de marca de AulaViva

Aplica esta guía en cada pantalla, componente y estado visual de AulaViva.
Prioriza una experiencia educativa profesional, cálida y clara; no produzcas
una plantilla SaaS genérica.

## Identidad y personalidad

AulaViva comunica confianza, cercanía, organización, modernidad y calidez
humana. La estética es editorial contemporánea: espacios amplios, jerarquía
tipográfica marcada y detalles relacionados con educación usados con moderación.
La interfaz debe sentirse profesional, no infantil ni administrativa.

## Paleta cromática

- **Verde bosque** `#244A3D`: color principal, navegación y acciones primarias.
- **Verde suave** `#DCE9DF`: superficies de apoyo y estados positivos sutiles.
- **Marfil** `#F8F6F0`: fondo predominante.
- **Coral** `#E98268`: acento para acciones o estados que requieran atención.
- **Gris oscuro** `#343A36`: texto principal.
- **Gris cálido** `#8B918B`: texto secundario y elementos auxiliares.

Usa el coral como acento, nunca como color dominante. Comprueba el contraste
de texto y no dependas únicamente del color para comunicar un estado.

## Tipografía y jerarquía

Combina una serif editorial expresiva para títulos seleccionados y destacados
con una sans-serif moderna para navegación, formularios, tablas y texto
operativo. Se pueden usar DM Sans, Manrope, Plus Jakarta Sans o alternativas
equivalentes disponibles sin introducir complejidad innecesaria.

Establece una jerarquía clara entre título de página, contexto, sección,
etiqueta y valor. Mantén una medida de lectura cómoda y evita usar muchos
tamaños o pesos en la misma vista.

## Componentes y superficies

- El dashboard administrativo debe usar sidebar, encabezado contextual,
  resumen visual, eventos recientes, pendientes y tabla de autorizaciones.
- Varía las jerarquías: no llenes la pantalla con tarjetas idénticas.
- Usa superficies marfil o blancas cálidas, bordes sutiles y separación clara.
- Prefiere bordes discretos y sombras mínimas; el contenido debe definir la
  jerarquía.
- Los botones primarios usan verde bosque; las acciones de atención pueden
  usar coral con texto de contraste.
- Los estados `pending`, `authorized` y `not_authorized` deben tener texto
  explícito además de cualquier color o icono.
- El portal de acudientes debe ser más sencillo que el dashboard y priorizar
  identificación del alumno, datos del evento, respuesta y confirmación.

## Espaciado, iconografía y decoración

Usa una escala de espaciado consistente y generosa. Alinea formularios,
tablas y acciones con una retícula estable. Utiliza iconos sobrios y
consistentes, nunca emojis como sustitutos. Las ilustraciones o motivos
educativos deben ser pequeños, funcionales y no competir con la tarea
principal.

## Responsive y accesibilidad

Diseña primero el portal para móviles; los controles de respuesta deben ser
grandes, claros y fáciles de alcanzar. El dashboard debe adaptarse sin perder
información ni acciones esenciales en pantallas estrechas.

Usa HTML semántico, etiquetas asociadas a campos, foco visible, navegación
completa por teclado, objetivos táctiles adecuados y contraste WCAG AA como
referencia. Incluye estados de carga, vacío, error y éxito comprensibles.
Proporciona texto alternativo cuando una imagen tenga significado.

## Evitar

- Dashboards genéricos o aspecto de plantilla SaaS.
- Gradientes excesivos, sombras fuertes o colores saturados.
- Tarjetas repetitivas sin jerarquía.
- Interfaces infantilizadas.
- Emojis como iconos.
- Decoración sin propósito.
- Texto pequeño, bajo contraste o estados comunicados solo por color.
