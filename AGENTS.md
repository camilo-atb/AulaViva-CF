# AulaViva

## Contexto

AulaViva es un MVP web para la Escuela Primaria Sofía, con aproximadamente 300
alumnos. Digitaliza las autorizaciones escolares que hoy se envían en papel,
sin exigir a las familias instalar una aplicación. El enlace para cada
acudiente se comparte por canales existentes, como WhatsApp.

El objetivo del MVP es permitir que personal autorizado cree eventos, genere
solicitudes individuales, reciba respuestas y consulte el estado consolidado.
Las calificaciones, tareas, notificaciones automáticas y la gestión de
personas autorizadas para recoger alumnos no forman parte del producto.

## Stack y arquitectura

- Frontend: Angular, TypeScript, componentes standalone, Angular Router y CSS.
- Backend: Go con API REST y JSON.
- Persistencia: SQLite con inicialización o migraciones controladas.
- Herramientas: Git, GitHub, VS Code y GitHub Copilot.

El frontend consume el backend mediante servicios HTTP. El backend sigue una
arquitectura por capas inspirada en Clean Architecture:
`domain -> application -> infrastructure/presentation`. El dominio no conoce
frameworks ni almacenamiento; la aplicación coordina casos de uso mediante
interfaces; infraestructura implementa SQLite; presentación adapta HTTP. El
servidor es el único responsable de acceder a SQLite. La arquitectura completa,
sus carpetas y contratos están en
[docs/architecture.md](docs/architecture.md).

Las decisiones arquitectónicas fueron aprobadas: la API usa `/api/events` y
`/api/requests`, responde `409 Conflict` ante respuestas duplicadas, usa
fechas `YYYY-MM-DD` sin hora, mantiene el MVP sin autenticación, SQLite y la
separación Domain/Application/Infrastructure/Presentation, y conserva
componentes standalone con organización feature-based en Angular.

## Alcance y restricciones

El MVP incluye dashboard administrativo, eventos, cursos y alumnos de
demostración, generación de solicitudes por alumno, enlaces individuales,
portal de acudientes, respuestas de autorización o rechazo y consulta de
resultados pendientes, aprobados y rechazados.

No implementar integración real con WhatsApp, calificaciones, tareas,
aplicación móvil nativa, notificaciones automáticas, roles avanzados ni
integraciones académicas externas. No inventar funcionalidades fuera de la
[especificación](specs/gestion-autorizaciones.md).

## Convenciones de desarrollo

- La documentación y la interfaz dirigida al usuario estarán en español.
- El código, símbolos, nombres de archivos de código, endpoints y contratos
  técnicos estarán en inglés.
- Mantener cambios pequeños, explícitos y fáciles de ejecutar localmente.
- Preferir tipos y validaciones claras antes que conversiones inseguras.
- Los errores deben comunicarse explícitamente; no ocultar fallos con valores
  predeterminados que parezcan exitosos.
- No añadir dependencias o infraestructura sin justificar su necesidad.
- Usar constructores `New...` e inyección manual de dependencias en Go.
- Usar `context.Context` en casos de uso y operaciones de persistencia.
- Crear interfaces pequeñas únicamente cuando expresen un contrato útil o
  permitan sustituir una dependencia para pruebas; no crear interfaces por
  cada struct.
- Evitar variables globales, dependencias circulares, punteros sin necesidad y
  funciones con responsabilidades mezcladas.
- Mantener los detalles SQL dentro de infraestructura y la lógica de negocio
  fuera de handlers y repositorios.

## Reglas para frontend

- Usar componentes standalone y una navegación definida con Angular Router.
- Encapsular las llamadas HTTP en servicios tipados.
- Organizar el código por `core`, `shared`, `features` y `layout`; cada
  funcionalidad puede contener únicamente las carpetas que necesite.
- Mantener las páginas como composición de componentes y servicios, sin lógica
  de negocio ni peticiones HTTP directas en componentes de presentación.
- Usar interfaces o tipos TypeScript para los DTO, formularios reactivos cuando
  el formulario lo justifique y tipado estricto.
- Centralizar la traducción visual de errores HTTP cuando sea transversal;
  mantener los errores específicos cerca de la funcionalidad que los presenta.
- Usar Signals u otro mecanismo reactivo de Angular para estado local derivado,
  sin introducir una librería global de estado para este MVP.
- Aplicar lazy loading solo a límites de navegación que aporten valor, como
  el área administrativa y el portal de acudientes.
- Mantener estados de carga, vacío y error visibles y comprensibles.
- Aplicar la identidad visual de
  [.github/skills/brand-design/SKILL.md](.github/skills/brand-design/SKILL.md)
  en todas las pantallas.
- Diseñar primero la experiencia móvil del portal de acudientes y mantener
  responsive el panel administrativo.
- Usar HTML semántico, navegación por teclado, foco visible, contraste
  suficiente y etiquetas asociadas a los controles.

## Reglas para backend

- Mantener separados dominio, aplicación, infraestructura y presentación.
- El dominio contiene entidades, tipos, errores y contratos de repositorio; no
  importa paquetes de infraestructura o HTTP.
- La aplicación contiene casos de uso y DTO; depende de contratos del dominio,
  nunca de SQLite.
- Infraestructura contiene conexión, esquema, migraciones y repositorios
  SQLite. Presentación contiene handlers, rutas, middleware y JSON.
- Permitir dependencias hacia capas externas solo desde las capas externas:
  `presentation -> application -> domain` y
  `infrastructure -> domain`; `main` ensambla todo.
- Validar entradas en los límites de la API y devolver JSON consistente.
- Usar códigos HTTP apropiados sin exponer detalles internos de SQLite.
- Pasar `context.Context` desde HTTP hasta persistencia y manejar errores de
  forma explícita con `errors.Is` o `errors.As` cuando corresponda.
- Construir dependencias manualmente desde `cmd/api/main.go`; no usar un
  contenedor de inyección.
- Generar solicitudes de forma determinista y evitar duplicados para un mismo
  evento y alumno.
- Hacer idempotente la respuesta de una solicitud: una solicitud respondida no
  puede cambiarse mediante una segunda respuesta, incluso bajo concurrencia.

## Reglas para persistencia

- Persistir eventos, cursos, alumnos y solicitudes en SQLite.
- Inicializar el esquema de forma controlada y reproducible.
- Usar consultas parametrizadas y restricciones de integridad apropiadas.
- Crear el evento y sus autorizaciones en una única transacción.
- Reforzar la regla evento-alumno con una restricción `UNIQUE` y aceptar como
  fuente de verdad la condición atómica de actualización de una solicitud
  pendiente.
- No almacenar secretos ni datos innecesarios.
- Mantener la persistencia detrás de una interfaz para poder probar la lógica
  sin acoplarla a SQLite.

## Reglas de seguridad

- Tratar todos los datos recibidos del cliente como no confiables.
- Validar identificadores, fechas, textos y estados permitidos.
- Los enlaces individuales deben identificar únicamente la solicitud necesaria
  y no exponer listados ni datos de otros alumnos.
- No registrar respuestas ni datos personales innecesarios en logs.
- Configurar CORS de forma explícita para el entorno local y revisar cualquier
  cambio antes de ampliarlo.

## Spec Driven Development

El orden de trabajo es obligatorio:

1. Preparación: contexto, estructura, especificaciones y skills.
2. Revisión: comprobar que estos documentos no se contradigan.
3. Planificación: proponer la arquitectura mínima que satisfaga la spec.
4. Implementación: comenzar solo tras una solicitud expresa.
5. Validación: comprobar cada criterio de aceptación con pruebas funcionales.

Antes de implementar, resolver cualquier ambigüedad de la spec. Una
funcionalidad no se considera terminada hasta verificar objetivamente sus
criterios de aceptación. No desarrollar frontend ni backend durante las fases
de preparación y revisión.
