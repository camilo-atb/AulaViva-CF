# AUT-001: Gestión de autorizaciones escolares

## Contexto del problema

La Escuela Primaria Sofía gestiona por papel las autorizaciones para salidas,
jornadas de vacunación y eventos. Los documentos pueden perderse o regresar
sin firma, y el personal no tiene una vista confiable del estado de cada
alumno. Las familias ya utilizan WhatsApp y no deben instalar una aplicación
nueva.

## Objetivo

Digitalizar el ciclo de una autorización escolar: crear un evento, generar una
solicitud individual por alumno, compartir un enlace, recibir una respuesta y
consultar el resultado consolidado desde el panel administrativo.

## Usuarios involucrados

- **Personal administrativo o docente:** crea eventos, genera solicitudes y
  consulta resultados.
- **Acudiente:** abre el enlace individual, revisa la información del evento y
  autoriza o no autoriza la participación del alumno.

El MVP no implementa autenticación ni gestión avanzada de roles; el acceso al
panel se considera parte del entorno controlado de la institución.

## Alcance funcional

1. Mostrar un dashboard administrativo con resumen y eventos recientes.
2. Listar y crear eventos con nombre, descripción, fecha y curso.
3. Disponer de cursos y alumnos de demostración.
4. Generar una solicitud individual para cada alumno del curso del evento.
5. Producir un enlace individual por solicitud para compartirlo manualmente.
6. Mostrar al acudiente el alumno, evento, fecha y descripción.
7. Registrar una respuesta de autorizar o no autorizar.
8. Consultar solicitudes pendientes, autorizadas y no autorizadas.

Quedan fuera WhatsApp Business API, calificaciones, tareas, aplicación móvil,
personas autorizadas para recoger alumnos, notificaciones automáticas y
conexiones con sistemas externos.

## Flujo principal

1. El personal crea un evento y selecciona un curso.
2. El sistema valida los datos y genera una solicitud por cada alumno del
   curso, cada una con un identificador o enlace individual.
3. El personal copia y comparte los enlaces con los acudientes.
4. El acudiente abre su enlace y consulta los datos del evento y del alumno.
5. El acudiente elige **Autorizar** o **No autorizar**.
6. El sistema registra la respuesta y muestra una confirmación.
7. El personal consulta el consolidado por estado desde el dashboard o el
   listado del evento.

## Reglas de negocio

- Un evento requiere nombre, fecha y curso; la descripción puede ser opcional.
- La fecha debe ser válida y el curso debe existir.
- Solo se generan solicitudes para los alumnos del curso seleccionado.
- Una solicitud pertenece a un único evento y alumno.
- Una solicitud nueva inicia en estado `pending`.
- Una solicitud respondida no puede modificarse mediante una segunda respuesta.
- Una respuesta válida cambia el estado a `authorized` o `not_authorized` y
  conserva la fecha de respuesta.
- El enlace de una solicitud no permite consultar otras solicitudes.
- Crear o regenerar solicitudes no debe duplicar la combinación evento-alumno.

## Estados de una solicitud

| Estado | Significado |
| --- | --- |
| `pending` | Aún no se ha registrado una respuesta. |
| `authorized` | El acudiente autorizó la participación. |
| `not_authorized` | El acudiente no autorizó la participación. |

## Diseño de la solución

La aplicación tendrá un frontend Angular y un backend Go con API REST JSON. El
frontend separará las vistas administrativas del portal de acudientes y
encapsulará el acceso HTTP en servicios. El backend separará handlers, reglas
de negocio y repositorios. SQLite almacenará eventos, cursos, alumnos y
solicitudes, con una inicialización controlada del esquema.

El enlace compartido contendrá un identificador opaco de solicitud, no un
listado ni información personal innecesaria. WhatsApp solo se simulará
mediante la acción de copiar o compartir el enlace.

La arquitectura detallada, la estructura de carpetas, las dependencias entre
capas, el modelo relacional y la estrategia de pruebas están documentados en
[docs/architecture.md](../docs/architecture.md). Este documento define el
comportamiento del producto; la arquitectura no debe ampliar su alcance.

## Criterios de aceptación

### CA-01: Creación y generación

Dado un curso con alumnos de demostración, cuando el personal crea un evento
válido y confirma la generación, entonces el sistema persiste el evento y crea
exactamente una solicitud `pending` por cada alumno de ese curso, cada una con
un enlace individual.

### CA-02: Respuesta del acudiente

Dado un enlace válido de una solicitud `pending`, cuando el acudiente consulta
el portal y selecciona **Autorizar** o **No autorizar**, entonces el sistema
persiste el estado correspondiente, registra la respuesta y muestra una
confirmación visible.

### CA-03: Consulta de resultados

Dado un evento con solicitudes en distintos estados, cuando el personal abre su
detalle, entonces puede ver cada alumno y su estado, y los totales de
`pending`, `authorized` y `not_authorized` coinciden con las solicitudes
persistidas.

### CA-04: Persistencia e idempotencia

Dado un evento y una solicitud respondida, cuando se reinicia el backend y se
consulta el evento, los datos y el estado siguen disponibles; si se intenta
responder de nuevo la misma solicitud, el sistema rechaza la operación sin
alterar la respuesta original.

## Restricciones técnicas

- Angular con TypeScript, componentes standalone, Angular Router y CSS.
- Go para la API REST y JSON para los contratos.
- SQLite como única base de datos del MVP.
- Ejecución local sencilla, sin microservicios ni dependencias innecesarias.
- El código y los símbolos técnicos estarán en inglés; esta especificación está
  en español.
- Las solicitudes deben estar relacionadas con un curso, un alumno y un evento,
  y la combinación evento-alumno debe ser única.
- La creación de un evento y sus solicitudes debe ser atómica.
- El registro de una respuesta debe conservar la primera respuesta válida
  mediante una operación atómica que solo actualice solicitudes `pending`.
- Los DTO HTTP usan JSON con `camelCase`; las fechas se intercambian como
  `YYYY-MM-DD` cuando representan una fecha de calendario.
- Los errores HTTP usan `{ "error": { "code": "...", "message": "..." } }`.
- Los estados públicos de una solicitud son `pending`, `authorized` y
  `not_authorized`.
- Decisiones aprobadas: las rutas principales son `/api/events` y
  `/api/requests`; una solicitud ya respondida devuelve `409 Conflict`; las
  fechas no incluyen hora; el MVP no implementa autenticación; y se mantienen
  SQLite, las cuatro capas de Go y componentes standalone feature-based en
  Angular.

## Contrato HTTP resumido

La API REST expone, como mínimo, los siguientes recursos:

| Método | Ruta | Propósito |
| --- | --- | --- |
| `GET` | `/api/events` | Listar eventos administrativos. |
| `POST` | `/api/events` | Crear un evento y sus solicitudes en una transacción. |
| `GET` | `/api/events/{eventId}` | Consultar un evento y sus resultados. |
| `GET` | `/api/courses` | Listar cursos demo para el formulario de eventos. |
| `GET` | `/api/requests/{requestToken}` | Consultar una solicitud individual. |
| `POST` | `/api/requests/{requestToken}/response` | Registrar una respuesta única. |

Las respuestas exitosas usan JSON y los códigos `200` para consultas o
actualizaciones, `201` para creación y `204` cuando no haya cuerpo. Se usa
`400` para entrada inválida, `404` para recursos inexistentes, `409` para una
solicitud ya respondida o una violación de unicidad y `500` para fallos
inesperados. No se fija aquí un esquema de autenticación porque está fuera del
MVP.

## Casos de error relevantes

- Datos obligatorios ausentes o con formato inválido al crear un evento:
  responder con error de validación y no persistir datos parciales.
- Curso inexistente: responder con recurso no encontrado.
- Enlace o solicitud inexistente: responder con recurso no encontrado, sin
  revelar información de otras solicitudes.
- Solicitud ya respondida: rechazar la segunda respuesta y conservar la
  primera.
- Estado de respuesta distinto de `authorized` o `not_authorized`: rechazarlo
  como entrada inválida.
- Fallo de SQLite: informar un error controlado sin presentarlo como una
  operación exitosa.
