# PA-001: Gestión de personas autorizadas para recoger estudiantes

## Estado

SPEC aprobada e implementada. Los endpoints descritos están disponibles en el
backend actual y la feature Angular correspondiente está integrada. La
funcionalidad continúa destinada a datos ficticios del MVP académico.

## Contexto del problema

La Escuela Primaria Sofía no cuenta con un registro centralizado de las
personas autorizadas para recoger a cada estudiante. Cuando llega una persona
diferente a la habitual, el personal docente o administrativo no dispone de
una fuente confiable y rápida para consultar quién está autorizado por la
familia.

Esto genera incertidumbre, pérdida de tiempo y dificultades operativas.

## Objetivo

Permitir que el personal administrativo consulte estudiantes y gestione las
personas autorizadas asociadas a cada uno mediante una interfaz web clara,
rápida y adecuada para un MVP académico.

Registrar una persona en AulaViva representa únicamente que el personal la
incluyó en el registro operativo del estudiante. No representa una verificación
de identidad, una autorización legal ni una confirmación independiente de la
familia.

## Usuarios y alcance

### Usuarios

- Personal administrativo o docente que consulta y mantiene la información.
- La funcionalidad no incorpora un portal separado para las personas
  autorizadas.

### Incluido en el MVP

1. Listar estudiantes.
2. Buscar estudiantes por nombre.
3. Consultar la ficha de un estudiante.
4. Listar las personas autorizadas de un estudiante.
5. Registrar una persona autorizada.
6. Editar sus datos.
7. Eliminar un registro.
8. Identificar estudiantes sin personas autorizadas registradas.
9. Mostrar estados de carga, error y ausencia de información.

### Fuera del alcance

- Verificación de identidad o validación contra documentos oficiales.
- Firma, consentimiento digital o aprobación por parte de la familia.
- Fotografías, biometría o archivos adjuntos.
- Control de acceso, autenticación, roles avanzados o auditoría completa.
- Notificaciones, integración real con WhatsApp o lectura de documentos.
- Registro histórico de cambios o recuperación de registros eliminados.
- Sustituir las solicitudes y respuestas del módulo de autorizaciones
  existente.

## Datos de una persona autorizada

Cada registro debe contener:

- Nombre completo.
- Tipo de documento.
- Número de documento.
- Parentesco o relación con el estudiante.
- Teléfono de contacto.

Los nombres técnicos de propiedades, endpoints y estados estarán en inglés.
La interfaz de usuario y los mensajes visibles estarán en español.

## Reglas de negocio

1. Una persona autorizada debe pertenecer a un estudiante existente.
2. Los campos obligatorios deben validarse antes de persistir.
3. No se permiten registros duplicados para el mismo estudiante con el mismo
   tipo y número de documento.
4. El mismo adulto puede estar autorizado para recoger a varios estudiantes;
   la unicidad se aplica por estudiante, tipo de documento y número de
   documento, no globalmente.
5. El estudiante asociado no puede cambiarse mediante la edición del registro.
6. Los datos creados o modificados deben persistir después de reiniciar el
   backend.
7. La eliminación solo debe afectar al registro seleccionado.
8. Consultar, editar o eliminar un registro inexistente debe producir un error
   explícito y no modificar otros registros.
9. Los errores deben conservar la estructura y los códigos HTTP definidos por
   los contratos existentes.
10. El sistema no verifica que la persona registrada sea realmente quien
    afirma ser, ni debe presentar el registro como una verificación de
    identidad.
11. Las operaciones válidas deben ser atómicas desde la perspectiva del
    cliente: una operación rechazada no debe dejar un registro parcialmente
    guardado.

## Estados de interacción

Estos estados describen la experiencia de la pantalla; no son estados de
negocio persistidos de una persona autorizada:

| Estado | Comportamiento esperado |
| --- | --- |
| `loading` | Mostrar una indicación clara mientras se consulta o guarda información. |
| `ready` | Mostrar estudiantes o personas autorizadas disponibles. |
| `empty` | Explicar que no hay resultados y ofrecer la siguiente acción. |
| `error` | Mostrar un mensaje comprensible y permitir reintentar cuando sea posible. |
| `success` | Confirmar que la operación terminó correctamente. |

Un estudiante sin personas autorizadas debe identificarse mediante texto
explícito, por ejemplo “Sin personas autorizadas registradas”, y no solo por
color o un icono.

## Criterios de aceptación

### CA-01: Registro de una persona autorizada

**Dado** un estudiante existente y un formulario con nombre completo, tipo de
documento, número de documento, relación y teléfono válidos,

**Cuando** el personal confirma el registro,

**Entonces** el sistema crea una persona autorizada asociada únicamente a ese
estudiante, responde con `201 Created` y la persona aparece en su listado con
los datos guardados.

### CA-02: Consulta de estudiantes y personas autorizadas

**Dado** un conjunto de estudiantes, algunos con personas autorizadas y otros
sin registros,

**Cuando** el personal lista o busca estudiantes por nombre y abre la ficha de
uno,

**Entonces** el sistema devuelve únicamente las coincidencias, muestra la
información básica del estudiante, lista sus personas autorizadas y señala
explícitamente cuando la lista está vacía.

### CA-03: Edición y eliminación aislada

**Dado** una persona autorizada existente asociada a un estudiante,

**Cuando** el personal actualiza sus datos válidos o confirma su eliminación,

**Entonces** la edición conserva el estudiante asociado y refleja los nuevos
datos, o la eliminación responde con `204 No Content`; en ambos casos ningún
otro registro resulta modificado.

### CA-04: Validaciones y duplicados

**Dado** un estudiante existente,

**Cuando** el personal envía campos obligatorios ausentes o intenta registrar
otra persona con el mismo tipo y número de documento para ese estudiante,

**Entonces** el sistema rechaza la operación con `400 Bad Request` para datos
inválidos o `409 Conflict` para el duplicado, devuelve el error JSON
documentado y no crea ni modifica datos parcialmente.

### CA-05: Estado vacío y persistencia

**Dado** un estudiante sin personas autorizadas y una persona registrada
válidamente,

**Cuando** el personal consulta el estudiante, reinicia la conexión o el
backend y vuelve a consultar, y después elimina el registro,

**Entonces** el estado vacío se muestra de forma explícita antes y después del
registro, la persona permanece disponible tras el reinicio y la eliminación
solo deja vacío el listado de ese estudiante.

## Métricas de éxito propuestas

Las siguientes son metas del MVP, no resultados medidos. Deben comprobarse con
datos ficticios, pruebas automatizadas y una validación manual controlada.

| Métrica | Forma de medición | Objetivo esperado | Momento en que se mide |
| --- | --- | --- | --- |
| Tiempo de registro | Medir desde la apertura del formulario hasta la confirmación de un registro válido | Menos de 60 segundos | Validación manual del flujo MVP |
| Estudiantes consultables | Comparar estudiantes demo cargados en SQLite con los devueltos por la consulta | 100% | Prueba funcional de consulta |
| Operaciones válidas persistidas | Crear, editar y eliminar registros; reiniciar la conexión y comparar el estado esperado | 100% | Pruebas de repositorio e integración |
| Registros incompletos aceptados | Contar operaciones aceptadas con algún campo obligatorio ausente o inválido | 0 | Pruebas de validación |
| Estudiantes sin registros identificables | Comparar estudiantes sin personas con los que muestran el estado vacío | 100% identificables | Prueba funcional de listado y ficha |

No se afirmará que estos objetivos fueron alcanzados hasta ejecutar y registrar
las mediciones correspondientes.

## Modelo de datos propuesto

La funcionalidad añade una relación de uno a muchos:

```text
students 1 ─── N authorized_people
```

La tabla existente `students` conserva su estructura y relación con `courses`.
Cada persona autorizada referencia exactamente un estudiante mediante una clave
foránea.

### Tabla `authorized_people`

| Campo | Tipo SQLite | Restricciones |
| --- | --- | --- |
| `id` | `INTEGER` | Clave primaria autoincremental |
| `student_id` | `INTEGER` | Obligatorio, clave foránea a `students(id)` |
| `full_name` | `TEXT` | Obligatorio, no vacío después de normalizar espacios |
| `document_type` | `TEXT` | Obligatorio, valor permitido por el contrato |
| `document_number` | `TEXT` | Obligatorio, no vacío después de normalizar espacios |
| `relationship` | `TEXT` | Obligatorio, no vacío después de normalizar espacios |
| `phone` | `TEXT` | Obligatorio, formato básico validable |

Restricciones adicionales:

- `FOREIGN KEY (student_id) REFERENCES students(id)`.
- `UNIQUE (student_id, document_type, document_number)`.
- La comparación de unicidad debe ser consistente con la normalización
  definida por la implementación, evitando aceptar duplicados por diferencias
  accidentales de espacios.
- La eliminación de un estudiante no forma parte del MVP; si en el futuro se
  permite, deberá definirse explícitamente la política para sus registros
  dependientes.

Índices propuestos:

- Índice sobre `authorized_people(student_id)` para listar rápidamente las
  personas de un estudiante.
- Índice sobre `students(full_name)` para soportar la búsqueda por nombre.
- La restricción `UNIQUE` crea el índice necesario para detectar duplicados.

La migración deberá ser posterior a las migraciones existentes, reproducible y
no destructiva. Esta SPEC no crea ni modifica migraciones.

## Contratos HTTP implementados

Los endpoints siguientes están implementados y cubiertos por pruebas de rutas
o de servicios. La búsqueda y las respuestas de detalle incluyen siempre
`authorizedPeople` como arreglo.

### Recursos

| Método | Ruta | Propósito |
| --- | --- | --- |
| `GET` | `/api/students?search={name}` | Listar estudiantes; `search` filtra por nombre cuando se proporciona. |
| `GET` | `/api/students/{studentId}` | Consultar la ficha de un estudiante y sus personas autorizadas. |
| `GET` | `/api/students/{studentId}/authorized-people` | Listar las personas autorizadas de un estudiante. |
| `POST` | `/api/students/{studentId}/authorized-people` | Registrar una persona autorizada. |
| `PUT` | `/api/authorized-people/{authorizedPersonId}` | Editar los datos de una persona autorizada. |
| `DELETE` | `/api/authorized-people/{authorizedPersonId}` | Eliminar una persona autorizada. |

La búsqueda debe ser acotada y devolver coincidencias de estudiantes; no debe
exponer personas autorizadas de otros estudiantes. El detalle del estudiante
puede incluir la lista anidada para evitar una segunda consulta en la ficha,
mientras que el endpoint de colección permite refrescarla de forma
independiente.

### Request de creación

`POST /api/students/12/authorized-people`

```json
{
  "fullName": "Laura Gómez",
  "documentType": "CC",
  "documentNumber": "123456789",
  "relationship": "Tía",
  "phone": "+57 300 555 0101"
}
```

### Response de creación

```json
{
  "id": 41,
  "studentId": 12,
  "fullName": "Laura Gómez",
  "documentType": "CC",
  "documentNumber": "123456789",
  "relationship": "Tía",
  "phone": "+57 300 555 0101"
}
```

### Response de estudiante

```json
{
  "id": 12,
  "courseId": 3,
  "fullName": "Sofía Martínez",
  "authorizedPeople": [
    {
      "id": 41,
      "studentId": 12,
      "fullName": "Laura Gómez",
      "documentType": "CC",
      "documentNumber": "123456789",
      "relationship": "Tía",
      "phone": "+57 300 555 0101"
    }
  ]
}
```

Cuando no existan personas autorizadas, `authorizedPeople` debe ser una lista
vacía (`[]`), no `null`.

### Request de edición

`PUT /api/authorized-people/41`

```json
{
  "fullName": "Laura Gómez Pérez",
  "documentType": "CC",
  "documentNumber": "123456789",
  "relationship": "Tía",
  "phone": "+57 300 555 0102"
}
```

La edición no recibe `studentId`; la persona permanece asociada al estudiante
original.

### Response de eliminación

`DELETE /api/authorized-people/41`

```text
204 No Content
```

## Códigos HTTP y errores

Las respuestas de error deben usar la estructura existente:

```json
{
  "error": {
    "code": "authorized_person_duplicate",
    "message": "Ya existe una persona con ese documento para este estudiante."
  }
}
```

| Situación | Código HTTP esperado | Código de error propuesto |
| --- | --- | --- |
| Request inválido, campo ausente o formato inválido | `400 Bad Request` | `invalid_request` |
| Estudiante inexistente | `404 Not Found` | `student_not_found` |
| Persona autorizada inexistente | `404 Not Found` | `authorized_person_not_found` |
| Documento duplicado para el mismo estudiante | `409 Conflict` | `authorized_person_duplicate` |
| Error inesperado de persistencia | `500 Internal Server Error` | `internal_error` |

No se deben exponer consultas SQL, rutas internas, trazas ni datos de otros
estudiantes en los mensajes HTTP.

## Diseño UX

### Listado de estudiantes

La pantalla administrativa mostrará:

- Título contextual y una explicación breve de la tarea.
- Campo de búsqueda por nombre con etiqueta visible.
- Lista o tabla de estudiantes con nombre, curso e indicador de personas
  autorizadas.
- Un indicador textual para estudiantes sin personas registradas.
- Estados de carga, error con reintento y ausencia de resultados.

La búsqueda debe ser directa y usable en pantallas estrechas. No se debe
convertir la pantalla en un dashboard de tarjetas repetitivas.

### Ficha del estudiante

La ficha mostrará primero la identidad operativa del estudiante y su curso,
seguida del listado de personas autorizadas. El estado vacío debe explicar qué
falta y presentar una acción clara para registrar la primera persona.

### Formulario de creación y edición

El formulario usará controles etiquetados para:

- Nombre completo.
- Tipo de documento.
- Número de documento.
- Parentesco o relación.
- Teléfono.

Debe mostrar validaciones junto al campo, impedir el envío de datos
incompletos, indicar carga durante el guardado y mostrar confirmación o error
después de la operación. La edición y la eliminación deben requerir acciones
explícitas; la eliminación debe pedir confirmación antes de ejecutarse.

### Identidad visual y accesibilidad

La experiencia debe seguir la skill de diseño de marca:

- Verde bosque `#244A3D` para navegación y acciones principales.
- Verde suave `#DCE9DF` para superficies de apoyo y estados positivos.
- Marfil `#F8F6F0` como fondo predominante.
- Coral `#E98268` como acento de atención, sin usarlo como color dominante.
- Gris oscuro `#343A36` para texto principal.
- Gris cálido `#8B918B` para información secundaria.

Los títulos seleccionados pueden usar una serif editorial y los controles,
tablas y textos operativos una sans-serif moderna. Se deben mantener
jerarquías claras, bordes sutiles y sombras mínimas, evitando gradientes,
emojis, aspecto infantil o una plantilla SaaS genérica.

Los estados no dependerán únicamente del color: deben incluir texto explícito,
contraste suficiente, foco visible, navegación por teclado y controles táctiles
adecuados. La ficha y el formulario deben funcionar en móvil y escritorio.

## Alineación arquitectónica

La implementación futura respetará las capas existentes:

- **Domain:** entidad `AuthorizedPerson`, tipos y reglas de unicidad o
  validación que no dependan de HTTP ni SQLite; contratos de repositorio
  únicamente cuando representen una necesidad real.
- **Application:** casos de uso para listar estudiantes, consultar ficha,
  listar, crear, editar y eliminar personas autorizadas; DTO de entrada y
  salida.
- **Infrastructure:** migración SQLite, consultas parametrizadas,
  restricciones, índices y adaptación de filas a entidades.
- **Presentation:** handlers, rutas, validación de parámetros y cuerpos,
  serialización JSON y traducción de errores a HTTP.
- **Angular:** una feature de gestión de personas autorizadas, con servicios
  tipados para `HttpClient`, páginas para las rutas y componentes pequeños
  para listado, ficha y formulario.

No se crearán carpetas vacías ni interfaces artificiales. La inyección manual
de dependencias, `context.Context`, el tipado estricto y las pruebas
reemplazables seguirán las reglas de `AGENTS.md` y
`docs/architecture.md`.

## Seguridad y privacidad

- La demostración académica debe utilizar únicamente datos ficticios.
- No deben incorporarse documentos reales de menores ni de sus familias.
- No se almacenan fotografías, biometría, archivos ni información adicional
  que no sea necesaria para el alcance aprobado.
- El teléfono y el número de documento son datos sensibles para el contexto
  operativo y no deben aparecer en logs innecesarios ni en mensajes de error.
- El MVP actual no tiene autenticación ni autorización administrativa. Antes de
  publicar datos reales, el panel debe contar con protección de acceso,
  revisión de privacidad y controles de auditoría adecuados.
- La interfaz no debe afirmar que un registro equivale a validar identidad o
  autorización legal.
- Los endpoints deben devolver únicamente el estudiante solicitado y sus
  registros relacionados, evitando listados cruzados o exposición accidental.

## Pruebas previstas

La implementación deberá añadir pruebas proporcionales al MVP:

### Backend

- Validación de campos obligatorios y formatos.
- Creación y consulta de personas autorizadas.
- Relación con un estudiante existente.
- Rechazo de duplicados por estudiante, tipo y número de documento.
- Edición sin cambiar el estudiante asociado.
- Eliminación únicamente del registro seleccionado.
- Persistencia después de cerrar y reabrir la conexión SQLite.
- Respuestas HTTP para `400`, `404`, `409`, `201` y `204`.

### Frontend

- Listado y búsqueda de estudiantes.
- Estados de carga, error, resultados vacíos y estudiantes sin registros.
- Validación del formulario.
- Creación, edición y eliminación mediante el servicio HTTP.
- Presentación de errores de duplicado y recurso inexistente.

Cada criterio de aceptación debe vincularse a una prueba automatizada o a una
validación manual reproducible antes de declarar terminada la funcionalidad.

## Revisión de compatibilidad y decisiones aprobadas

Esta SPEC no cambia el comportamiento de `events`, `authorizations`,
`courses` ni de sus endpoints existentes. Añade un recurso relacionado con
`students` y mantiene las convenciones de JSON, errores, capas, SQLite y
componentes standalone ya aprobadas.

El alcance es proporcional a una institución de aproximadamente 300 alumnos:
incluye una relación sencilla, operaciones CRUD y búsqueda básica, pero no
incorpora identidad digital, roles, auditoría ni integraciones externas.

Decisiones aprobadas para la implementación:

1. Las rutas son `/api/students` y `/api/authorized-people`.
2. `documentType` usa el catálogo `CC`, `TI`, `CE`, `PASSPORT` y `OTHER`.
3. El teléfono admite formatos nacionales e internacionales, incluido `+57`,
   con una validación de caracteres permitidos y entre 7 y 15 dígitos.
4. El detalle del estudiante incluye siempre `authorizedPeople`, incluso
   vacío, y se conserva el endpoint independiente de colección.

La funcionalidad no implementa autenticación ni verificación de identidad. El
panel administrativo debe protegerse antes de utilizar datos reales.
