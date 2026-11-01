# Estrategia de producto e ICP

## Propósito

Este documento orienta las decisiones funcionales, de experiencia y de
comunicación de AulaViva para la tercera actividad de Spec Driven Development.
No reemplaza las especificaciones funcionales; define para quién se construye
el producto y qué problemas deben priorizarse.

## Cliente ideal (ICP)

AulaViva está orientada a instituciones educativas de educación primaria
privadas o de gestión independiente, principalmente en Colombia y
Latinoamérica, con aproximadamente 150 a 600 estudiantes.

Estas instituciones suelen:

- Utilizar WhatsApp como canal principal con las familias.
- Gestionar autorizaciones escolares mediante documentos físicos.
- Tener poco o ningún equipo interno de tecnología.
- Necesitar herramientas sencillas, accesibles desde cualquier navegador y
  fáciles de adoptar.
- Buscar menos trabajo administrativo y mejor control de la información.

El decisor de compra o de adopción es normalmente el rector, director o
coordinador administrativo. El personal docente y administrativo participa en
la operación diaria, mientras que padres y acudientes interactúan con los
enlaces o portales de cada proceso.

## Problemas principales del ICP

1. Los documentos físicos se pierden, llegan incompletos o no regresan a
   tiempo.
2. El personal no tiene una vista centralizada y actualizada de las respuestas
   de las familias.
3. Consultar información operativa de un estudiante requiere revisar papeles,
   chats o preguntar a otras personas.
4. Las instituciones necesitan mejorar el control sin comprar una plataforma
   académica completa ni exigir otra aplicación a las familias.
5. Las herramientas tecnológicas complejas generan dependencia de soporte y
   resistencia de adopción.

## Necesidades funcionales

- Crear y consultar eventos escolares y sus autorizaciones.
- Consultar estudiantes y su curso desde un único lugar.
- Gestionar información operativa asociada a cada estudiante con permisos
  claros y validaciones.
- Identificar rápidamente información faltante o estados pendientes.
- Acceder desde un navegador de escritorio, tableta o teléfono.
- Compartir enlaces mediante los canales que la institución ya utiliza.
- Persistir los cambios y mostrar errores comprensibles.

Estas necesidades deben resolverse con flujos breves y una interfaz clara, no
con módulos administrativos extensos.

## Motivaciones de compra y adopción

- Reducir el tiempo dedicado a perseguir y organizar documentos.
- Disminuir la incertidumbre durante actividades y procesos escolares.
- Tener trazabilidad básica sin instalar infraestructura compleja.
- Mantener la comunicación familiar en canales conocidos como WhatsApp.
- Mostrar una mejora visible con una inversión y capacitación pequeñas.
- Empezar con un problema concreto y ampliar solo si el uso demuestra valor.

## Objeciones de adopción

- “Las familias no instalarán otra aplicación.”
- Temor a que el personal no tenga tiempo para aprender una herramienta nueva.
- Dudas sobre privacidad y manejo de información de menores.
- Dependencia de conectividad o de un sistema que parezca difícil de mantener.
- Preferencia por continuar usando hojas de cálculo, papel o chats.
- Preocupación por costos recurrentes y soporte técnico.

La respuesta de producto debe ser un acceso web sencillo, lenguaje no técnico,
datos ficticios en la demostración, mínimo esfuerzo de configuración y
transparencia sobre las limitaciones de seguridad y privacidad.

## User Persona principal: coordinadora administrativa

**Nombre representativo:** Marcela  
**Rol:** coordinadora administrativa de una primaria independiente.  
**Contexto:** coordina comunicaciones operativas para 250 estudiantes y
trabaja con un equipo pequeño. No es especialista en tecnología.  
**Objetivos:** saber qué está pendiente, encontrar información de un estudiante
en segundos y reducir el trabajo manual antes de cada actividad.  
**Frustraciones:** documentos extraviados, información repartida entre chats y
la necesidad de confirmar datos varias veces.  
**Necesita:** una interfaz directa, filtros simples, estados visibles y
confirmación clara de que cada cambio quedó guardado.

## User Persona secundaria: docente de primaria

**Nombre representativo:** Andrés  
**Rol:** docente que acompaña diariamente a un grupo de estudiantes.  
**Contexto:** usa principalmente su teléfono y necesita consultar información
sin interrumpir una actividad o buscar carpetas físicas.  
**Objetivos:** identificar rápidamente el estado de un estudiante y contar con
información operativa confiable cuando recibe a una familia o acompaña una
actividad.  
**Frustraciones:** listas desactualizadas, respuestas ambiguas y sistemas que
requieren demasiados pasos.  
**Necesita:** búsqueda rápida, lectura clara en móvil y acciones explícitas,
sin configuraciones técnicas.

## Propuesta de valor diferenciadora

AulaViva convierte procesos escolares dispersos en flujos digitales pequeños y
comprensibles, accesibles desde un navegador y compatibles con la comunicación
que la institución ya usa. En lugar de intentar reemplazar todos los sistemas
académicos, resuelve tareas operativas concretas con una experiencia cálida,
profesional y adecuada para equipos escolares pequeños.

## Principios de producto derivados

1. **Sencillez antes que amplitud:** cada módulo debe resolver una tarea
   frecuente con pocos pasos.
2. **Información accionable:** mostrar estado, faltantes y siguiente acción,
   no solo registros.
3. **Adopción sin fricción:** navegador, enlaces compartibles y lenguaje claro.
4. **Confianza responsable:** usar datos ficticios en la demostración y no
   presentar el MVP como verificación de identidad o sistema de seguridad
   integral.
5. **Calidez profesional:** conservar la identidad editorial definida en la
   skill de diseño de marca, evitando una interfaz administrativa genérica.
6. **Privacidad por defecto:** minimizar datos, evitar biometría y documentar
   cualquier limitación antes de usar información real.

## Relación con la siguiente funcionalidad

La gestión de personas autorizadas para recoger estudiantes responde
directamente a los problemas de consulta operativa y control de información del
ICP. Debe implementarse como un módulo acotado relacionado con `students`,
reutilizando las capas Go y la organización feature-based de Angular existentes,
sin alterar los contratos actuales de eventos y autorizaciones.

La especificación funcional de ese módulo se elaborará y revisará en la fase
posterior, antes de modificar el código.

## Métricas de éxito propuestas

Estas son metas para el MVP; todavía no son resultados medidos:

| Métrica | Objetivo propuesto |
| --- | --- |
| Tiempo promedio para registrar una persona autorizada | Menos de 60 segundos en un flujo válido. |
| Estudiantes demo consultables | 100% de los estudiantes cargados por el backend. |
| Operaciones válidas reflejadas en SQLite | 100% en pruebas funcionales de crear, editar y eliminar. |
| Registros aceptados con campos obligatorios incompletos | 0%. |
| Estudiantes sin información visible | 100% identificables mediante un estado vacío o indicador explícito. |
| Operaciones con error explicado al usuario | 100% de los casos de validación y recurso inexistente cubiertos por pruebas. |

## Estado de medición

Durante esta fase no se han medido tiempos de uso ni resultados con usuarios
reales. Los objetivos anteriores deberán comprobarse con datos ficticios,
pruebas automatizadas y una validación manual controlada después de aprobar la
SPEC funcional.
