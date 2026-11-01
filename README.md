# AulaViva

AulaViva es un MVP web para digitalizar las autorizaciones escolares de la
Escuela Primaria Sofía. Permite que el personal cree eventos, genere
solicitudes individuales y consulte las respuestas de los acudientes mediante
enlaces compartibles, sin exigir una aplicación adicional.

## Problema y objetivo

Las autorizaciones en papel se pierden, pueden regresar sin firma y dificultan
conocer el estado de cada alumno. El objetivo es ofrecer un flujo digital
sencillo para crear un evento, solicitar una respuesta y consultar el
consolidado de autorizaciones.

## Funcionalidades del MVP

- Dashboard administrativo.
- Listado y creación de eventos con curso, fecha y descripción.
- Cursos y alumnos de demostración.
- Solicitudes individuales y enlaces compartibles.
- Portal web para acudientes.
- Respuestas de autorizar o no autorizar.
- Consulta de estados pendientes, autorizados y no autorizados.
- Persistencia con SQLite.

No incluye integración real con WhatsApp, calificaciones, tareas, aplicación
móvil, notificaciones automáticas, roles avanzados ni registro de personas
autorizadas para recoger alumnos.

## Stack tecnológico

- Angular y TypeScript, con componentes standalone, Angular Router y CSS.
- Go para la API REST y JSON.
- SQLite para persistencia.

## Estructura del repositorio

```text
.
├── AGENTS.md
├── README.md
├── docs/
│   └── architecture.md
├── specs/
│   └── gestion-autorizaciones.md
├── .github/
│   └── skills/
│       └── brand-design/
│           └── SKILL.md
├── frontend/
│   ├── src/app/core/
│   ├── src/app/features/
│   ├── src/app/layout/
│   └── src/app/shared/
└── backend/
```

## Arquitectura

El backend es una aplicación Go única, organizada por capas: `domain` contiene
reglas y contratos, `application` coordina casos de uso, `infrastructure`
implementa SQLite y `presentation` expone HTTP. `cmd/api` ensambla las
dependencias manualmente.

El frontend Angular usa componentes standalone y organización por
funcionalidades. `core` contiene servicios transversales, `shared` elementos
reutilizables, `features` las áreas de negocio (`dashboard`, `events` y
`authorizations`) y `layout` las estructuras administrativas y del portal.
Las páginas componen la interfaz; los servicios encapsulan HTTP y los modelos
tipan los contratos.

La creación de eventos y solicitudes es transaccional. SQLite refuerza la
unicidad evento-alumno, y las respuestas se registran mediante una actualización
atómica de solicitudes pendientes para conservar la primera respuesta válida.
El frontend obtiene los cursos demo mediante `GET /api/courses`.
La descripción completa está en
[docs/architecture.md](docs/architecture.md).

## Estado actual

Fase 4 de Spec Driven Development: backend Go y frontend Angular implementados
y validados. El frontend consume datos reales del backend; no contiene datos
inventados para las vistas principales.

El proyecto está preparado para ejecutarse localmente. Esta entrega contiene el
código fuente, las pruebas y la documentación del MVP; no incluye despliegue
en producción ni configuración de hosting.

## Requisitos previos

- Go 1.22 o compatible.
- Node.js y npm.
- Un navegador moderno.

## Ejecución local

### Backend

Desde `backend/`:

```powershell
go test ./...
go run ./cmd/api
```

La API estará disponible en `http://localhost:8080`.

### Frontend

Desde `frontend/`:

```powershell
npm install
npm run build
npm start
```

La aplicación estará disponible en `http://localhost:4200`. El backend debe
estar ejecutándose para que dashboard, eventos y portal consulten datos reales.

## Documentación SDD

- [AGENTS.md](AGENTS.md): contexto persistente y convenciones.
- [Especificación funcional](specs/gestion-autorizaciones.md).
- [Arquitectura](docs/architecture.md).
- [Skill de diseño de marca](.github/skills/brand-design/SKILL.md).
- [Documentación del backend](backend/README.md).
