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
- Consulta de estudiantes y gestión de personas autorizadas para recogerlos.
- Persistencia con SQLite.

No incluye integración real con WhatsApp, calificaciones, tareas, aplicación
móvil, notificaciones automáticas, roles avanzados ni autenticación. El
registro de personas autorizadas utiliza datos ficticios y no representa una
verificación legal de identidad.

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
│   ├── gestion-autorizaciones.md
│   └── gestion-personas-autorizadas.md
├── .github/
│   └── skills/
│       └── brand-design/
│           └── SKILL.md
├── frontend/
│   ├── src/app/core/
│   ├── src/app/features/
│   │   └── authorized-people/
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
[docs/architecture.md](docs/architecture.md). Los contratos de la nueva
funcionalidad están en
[specs/gestion-personas-autorizadas.md](specs/gestion-personas-autorizadas.md).

## Estado actual

Fase 4 de Spec Driven Development: backend Go y frontend Angular implementados
y validados. El frontend consume datos reales del backend; no contiene datos
inventados para las vistas principales. PA-001 incluye listado y búsqueda de
estudiantes, ficha, consulta y CRUD de personas autorizadas.

El proyecto está preparado para ejecutarse localmente y contiene configuración
base para un despliegue posterior en Vercel y Render. Esta entrega no crea
servicios ni publica la aplicación.

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

## Preparación de despliegue

### Backend en Render

El archivo [render.yaml](render.yaml) define un Web Service Go con:

- `rootDir: backend`.
- Build: `go build -o bin/aulaviva ./cmd/api`.
- Start: `./bin/aulaviva`.
- Plan `free`.
- SQLite en `aulaviva.db`.

Variables requeridas en Render:

```text
AULAVIVA_DB_PATH=aulaviva.db
AULAVIVA_ALLOWED_ORIGIN=https://aula-viva-cf.vercel.app
```

Render inyecta `PORT`; el backend lo prioriza sobre `AULAVIVA_PORT`.

El plan gratuito no ofrece almacenamiento persistente para SQLite. La base
puede reiniciarse o perderse durante un redeploy, por lo que esta configuración
es apropiada únicamente para una demostración académica con datos ficticios.
No se deben guardar datos reales ni presentar este despliegue como una
instancia de producción con persistencia garantizada.

### Frontend en Vercel

Configura el proyecto con:

```text
Root Directory: frontend
Install Command: npm ci
Build Command: npm run build
Output Directory: dist/aulaviva/browser
```

La configuración de producción apunta a
`https://aulaviva-cf.onrender.com/api`. Angular 18 genera el archivo
`index.html` dentro de `dist/aulaviva/browser`; el archivo
[vercel.json](frontend/vercel.json) publica ese directorio y redirige las rutas
de Angular a `index.html` para permitir recargas en rutas internas.

La aplicación no implementa autenticación. Antes de utilizar datos reales,
debe añadirse una protección administrativa; CORS no reemplaza autenticación.
La demo pública debe utilizar únicamente datos ficticios.

## Documentación SDD

- [AGENTS.md](AGENTS.md): contexto persistente y convenciones.
- [Especificación funcional](specs/gestion-autorizaciones.md).
- [SPEC PA-001](specs/gestion-personas-autorizadas.md).
- [Arquitectura](docs/architecture.md).
- [Skill de diseño de marca](.github/skills/brand-design/SKILL.md).
- [Documentación del backend](backend/README.md).
