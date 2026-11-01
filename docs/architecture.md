# Arquitectura de AulaViva

## 1. Objetivo y límites

Esta arquitectura guía la implementación del MVP de autorizaciones escolares
para una institución de aproximadamente 300 alumnos. Busca mantener separadas
las reglas de negocio, los detalles técnicos y la presentación sin introducir
microservicios, contenedores de inyección ni una plataforma de estado global.

El sistema será una aplicación Angular servida junto a una API Go que usa
SQLite. El despliegue y la ejecución local deben seguir siendo simples. Las
funcionalidades fuera de la especificación, como autenticación, notificaciones
automáticas o integración real con WhatsApp, no forman parte de esta
arquitectura.

## 2. Arquitectura general

La solución tiene dos módulos ejecutables comunicados por HTTP:

```text
Angular
  core / layout / features / shared
             |
             | HTTP + JSON
             v
Go API
  presentation -> application -> domain
       ^              ^
       |              |
  infrastructure ----+
       |
     SQLite
```

La dirección conceptual de dependencias es:

- `domain`: no depende de ninguna capa externa.
- `application`: depende de contratos y tipos del dominio.
- `infrastructure`: implementa contratos del dominio y depende de SQLite.
- `presentation`: depende de la aplicación y traduce HTTP a sus DTO.
- `cmd/api`: compone implementaciones concretas y configura el servidor.

No se permite que el dominio importe HTTP, SQLite o paquetes de Angular; que la
aplicación ejecute SQL; ni que los handlers contengan reglas de negocio.

## 3. Estructura propuesta

### Backend Go

```text
backend/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── domain/
│   │   ├── entities/
│   │   │   ├── authorization.go
│   │   │   ├── course.go
│   │   │   ├── event.go
│   │   │   └── student.go
│   │   ├── errors/
│   │   │   └── errors.go
│   │   └── repositories/
│   │       ├── authorization_repository.go
│   │       ├── course_repository.go
│   │       ├── event_repository.go
│   │       └── student_repository.go
│   ├── application/
│   │   ├── dto/
│   │   │   ├── authorization.go
│   │   │   └── event.go
│   │   └── services/
│   │       ├── authorization_service.go
│   │       └── event_service.go
│   ├── infrastructure/
│   │   └── persistence/
│   │       └── sqlite/
│   │           ├── migrations/
│   │           ├── connection.go
│   │           ├── schema.go
│   │           └── repositories/
│   └── presentation/
│       └── http/
│           ├── handlers/
│           ├── middleware/
│           ├── responses/
│           └── routes/
├── go.mod
└── README.md
```

Las carpetas solo se crean cuando contienen una necesidad real. La estructura
anterior es el destino de implementación; no implica generar archivos vacíos
durante la fase documental.

### Frontend Angular

```text
frontend/
└── src/
    ├── app/
    │   ├── core/
    │   │   ├── config/
    │   │   ├── http/
    │   │   └── services/
    │   ├── features/
    │   │   ├── authorizations/
    │   │   │   ├── components/
    │   │   │   ├── models/
    │   │   │   ├── pages/
    │   │   │   └── services/
    │   │   ├── dashboard/
    │   │   │   ├── components/
    │   │   │   └── pages/
    │   │   └── events/
    │   │       ├── components/
    │   │       ├── models/
    │   │       ├── pages/
    │   │       └── services/
    │   ├── layout/
    │   │   ├── admin-layout/
    │   │   ├── guardian-layout/
    │   │   ├── header/
    │   │   ├── navigation/
    │   │   └── sidebar/
    │   ├── shared/
    │   │   ├── components/
    │   │   ├── directives/
    │   │   ├── pipes/
    │   │   └── utils/
    │   ├── app.config.ts
    │   └── app.routes.ts
    └── styles.css
```

No se crearán carpetas `guards`, `interceptors` o `utils` hasta que exista una
necesidad concreta. `core` aloja singletons y configuración; `shared` no
conoce una funcionalidad específica; `features` contiene páginas, componentes,
modelos y servicios de su propio flujo; `layout` compone la estructura visual.

## 4. Responsabilidades por capa

### Domain

Define `Event`, `Course`, `Student` y `Authorization`, sus tipos de estado,
errores de negocio y contratos de repositorio. Sus reglas incluyen estados
válidos, asociación evento-alumno y transición única desde `pending`. No debe
conocer JSON, SQL, HTTP ni detalles de Angular.

Las interfaces se definen donde se consumen como contratos útiles, no como
reflejo automático de cada struct. Los repositorios pueden separarse cuando
los casos de uso requieran capacidades distintas.

### Application

Coordina los casos de uso de crear evento, generar solicitudes, consultar
resultados y registrar respuestas. Recibe `context.Context`, valida
invariantes de negocio, transforma DTO de entrada y salida y depende de
interfaces del dominio. No construye respuestas HTTP ni contiene consultas
SQL.

La creación del evento y sus solicitudes debe recibir una capacidad
transaccional del repositorio o un servicio de persistencia equivalente, de
forma que todo confirme o todo se revierta.

### Infrastructure

Abre y configura SQLite, aplica el esquema o migraciones y contiene las
consultas parametrizadas. Sus repositorios convierten filas a entidades del
dominio y traducen errores de SQLite a errores que la aplicación pueda
interpretar. La transacción de creación se inicia y confirma aquí, bajo la
coordinación del caso de uso.

### Presentation

Define rutas, handlers y serializadores JSON. Valida forma y tipos de entrada,
extrae parámetros, invoca servicios de aplicación y traduce errores conocidos a
códigos HTTP. No decide reglas de autorización ni accede a SQLite. Middleware
solo se añade para necesidades reales, como CORS explícito o logging seguro.

### Frontend

Las páginas representan rutas y coordinan la composición. Los componentes
reutilizables presentan datos y emiten eventos de UI. Los servicios de cada
feature encapsulan `HttpClient`, DTO y estados de carga/error; `core` aloja
configuración y preocupaciones transversales; `layout` contiene navegación y
estructura. Formularios complejos usan Reactive Forms y el estado local
derivado puede usar Signals.

## 5. Reglas de implementación Go

- Usar structs para entidades, comandos, resultados y respuestas, con tipos
  propios cuando hagan explícitos estados o valores del dominio.
- Usar funciones `New...` para construir servicios, repositorios y handlers con
  sus dependencias requeridas.
- Inyectar dependencias manualmente en `cmd/api/main.go`; no usar un contenedor.
- Pasar `context.Context` desde el handler al servicio y al repositorio.
- Comprobar y devolver cada error; usar `errors.Is`/`errors.As` para decisiones
  sobre errores conocidos.
- Evitar globales mutables, punteros innecesarios y métodos que hagan demasiadas
  cosas.
- Mantener interfaces pequeñas y justificadas por desacoplamiento o pruebas.
- No crear interfaces artificiales para structs sin consumidores alternativos.
- Mantener SQL, escaneo de filas y transacciones dentro de infraestructura.
- Facilitar pruebas sustituyendo repositorios mediante implementaciones pequeñas
  en memoria o dobles de prueba.

## 6. Reglas de implementación Angular

- Usar componentes standalone y tipado estricto.
- Separar páginas, componentes de presentación, modelos y servicios HTTP.
- No realizar peticiones HTTP ni lógica de negocio en componentes visuales.
- Centralizar la configuración de `HttpClient` y los errores verdaderamente
  transversales en `core`.
- Añadir interceptors o guards solo cuando exista una necesidad funcional.
- Usar lazy loading para separar el área administrativa y el portal cuando
  reduzca el acoplamiento; no usarlo como ritual.
- Preferir Signals para estado local simple y derivado; no introducir store
  global para este MVP.
- Usar Reactive Forms para la creación de eventos y validación declarativa.
- Mantener estados de carga, vacío, error y éxito explícitos y accesibles.

## 7. Contrato Angular-Go

### JSON y nombres

Los DTO HTTP usan `camelCase`, UTF-8 y objetos JSON. Los nombres internos de
Go y TypeScript permanecen en inglés; las etiquetas y mensajes visibles son
españoles. No se exponen filas SQLite directamente.

Las fechas de evento son fechas de calendario y usan `YYYY-MM-DD`. Los
timestamps técnicos, si se necesitan para auditoría mínima, usan ISO 8601 en
UTC. El frontend no debe interpretar una fecha de calendario como un instante
local que pueda cambiar de día.

### Respuestas y errores

Respuestas exitosas:

- `200 OK`: consultas y actualizaciones con cuerpo.
- `201 Created`: creación de eventos.
- `204 No Content`: operación exitosa sin representación.

Errores:

```json
{
  "error": {
    "code": "authorization_already_answered",
    "message": "La solicitud ya tiene una respuesta."
  }
}
```

Se usan `400` para entrada inválida, `404` para recursos inexistentes, `409`
para conflictos de estado o unicidad y `500` para fallos inesperados sin
detalles internos. El frontend transforma `code` en mensajes en español y no
depende de comparar textos técnicos.

### Recursos

| Método | Ruta | Uso |
| --- | --- | --- |
| `GET` | `/api/events` | Lista administrativa de eventos. |
| `POST` | `/api/events` | Crea un evento y sus solicitudes. |
| `GET` | `/api/events/{eventId}` | Detalle y resultados del evento. |
| `GET` | `/api/courses` | Lista los cursos demo disponibles para crear eventos. |
| `GET` | `/api/requests/{requestToken}` | Vista individual del acudiente. |
| `POST` | `/api/requests/{requestToken}/response` | Registra una respuesta. |

Los estados de autorización compartidos por ambos módulos son exactamente
`pending`, `authorized` y `not_authorized`.

## 8. Modelo relacional SQLite

```text
courses 1 --- N students
courses 1 --- N events
events  1 --- N authorizations
students 1 --- N authorizations
```

Tablas y restricciones:

- `courses`: `id` como clave primaria, `name` obligatorio y único.
- `students`: `id` como clave primaria, `course_id` obligatorio con clave
  foránea a `courses`, y `full_name` obligatorio.
- `events`: `id` como clave primaria, `course_id` obligatorio con clave foránea,
  `name`, `description` y `event_date` obligatorios según el contrato.
- `authorizations`: `id` como clave primaria, `event_id` y `student_id`
  obligatorios con claves foráneas, `request_token` opaco y único, `status`
  restringido a los tres estados, y `responded_at` nulo solo para `pending`.

La tabla `authorizations` tiene `UNIQUE(event_id, student_id)` y un índice por
`event_id`, otro por `student_id` y otro por `status` cuando las consultas lo
justifiquen. Las claves foráneas deben estar activas en cada conexión SQLite.

Al crear un evento, una transacción inserta el evento y sus autorizaciones a
partir de los alumnos del curso. Si falla cualquier inserción, se revierte
todo. Para responder, se ejecuta una actualización atómica equivalente a:

```sql
UPDATE authorizations
SET status = ?, responded_at = ?
WHERE request_token = ? AND status = 'pending';
```

Si las filas afectadas son cero, se distingue entre solicitud inexistente y
solicitud ya respondida, sin sobrescribir la primera respuesta. Esto garantiza
que dos respuestas concurrentes no puedan modificar el mismo registro.

## 9. Estrategia de pruebas

### Backend

- Unitarias de servicios de aplicación con repositorios sustituibles:
  validación, generación por alumno, estados y errores.
- Repositorios contra SQLite temporal: esquema, claves foráneas, unicidad,
  transacciones y persistencia tras reconexión.
- Handlers y endpoints con servidor HTTP de prueba: JSON, códigos, contratos,
  errores de entrada y recursos inexistentes.
- Pruebas específicas para regeneración sin duplicados y respuestas repetidas,
  incluyendo una prueba de concurrencia sobre la actualización atómica.

### Frontend

- Componentes importantes: formulario de evento, resumen de resultados y
  portal de respuesta, incluyendo estados de carga/error/confirmación.
- Servicios HTTP con `HttpTestingController` para verificar rutas, DTO y
  manejo de errores.
- Validación de formulario para campos obligatorios y fecha.
- Flujo del acudiente desde la carga de solicitud hasta una respuesta y su
  confirmación visual.

Las pruebas deben cubrir primero los cuatro criterios de aceptación de la
especificación. No se requiere una plataforma de pruebas distribuida ni
pruebas de carga para este MVP.

## 10. Adecuación al MVP

La propuesta mantiene un único backend, una base SQLite, inyección manual,
interfaces limitadas a contratos de repositorio y estado frontend local. Las
cuatro capas del backend son límites de responsabilidad, no servicios
desplegables separados. `core`, `shared`, `features` y `layout` reflejan
necesidades reales de Angular y no obligan a crear carpetas vacías.

Para aproximadamente 300 alumnos, esta separación facilita pruebas y futuras
extensiones sin introducir microservicios, eventos, CQRS, un ORM, un store
global o un sistema de autenticación prematuro.

## Decisiones aprobadas

Las siguientes decisiones fueron aprobadas antes de iniciar la implementación:

- Las rutas principales serán `/api/events` y `/api/requests`.
- Una solicitud ya respondida produce `409 Conflict`.
- Las fechas de evento son `YYYY-MM-DD`, sin hora.
- El MVP no implementa autenticación; se ejecuta localmente con datos ficticios.
- Se mantiene SQLite y la separación Domain, Application, Infrastructure y
  Presentation en Go.
- Angular usará componentes standalone y organización por funcionalidades.
