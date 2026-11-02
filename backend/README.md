# AulaViva Backend

API REST del MVP de AulaViva, implementada en Go con SQLite.

## Ejecución local

Desde esta carpeta:

```powershell
go run ./cmd/api
```

La API escucha en `http://localhost:8080`. Se pueden personalizar la ruta de
SQLite y el puerto:

```powershell
$env:AULAVIVA_DB_PATH = "data\aula.db"
$env:AULAVIVA_PORT = "8081"
go run ./cmd/api
```

La primera ejecución crea el esquema y carga cursos y alumnos ficticios.

Variables de entorno disponibles:

```text
AULAVIVA_DB_PATH       Ruta de SQLite (por defecto: aulaviva.db)
AULAVIVA_PORT          Puerto local alternativo (por defecto: 8080)
PORT                   Puerto priorizado por proveedores como Render
AULAVIVA_ALLOWED_ORIGIN Origen permitido por CORS
```

En Render, `AULAVIVA_ALLOWED_ORIGIN` debe ser exactamente
`https://aula-viva-cf.vercel.app`, sin una barra final.

En Render gratuito puede usarse `AULAVIVA_DB_PATH=aulaviva.db`, pero el
filesystem no es persistente: la base puede perderse al reiniciar o desplegar
el servicio. Esta configuración solo es adecuada para una demostración
académica con datos ficticios. No se debe utilizar para datos reales sin
contratar almacenamiento persistente o migrar posteriormente a una base de
datos administrada.

El endpoint `GET /api/courses` expone los cursos demo para el formulario del
frontend. En desarrollo se permite CORS desde `http://localhost:4200`.

Endpoints principales:

- `GET /api/courses`
- `GET /api/events`
- `POST /api/events`
- `GET /api/events/{id}`
- `GET /api/requests/{token}`
- `POST /api/requests/{token}/response`
- `GET /api/students?search={name}`
- `GET /api/students/{id}`
- `GET /api/students/{id}/authorized-people`
- `POST /api/students/{id}/authorized-people`
- `PUT /api/authorized-people/{id}`
- `DELETE /api/authorized-people/{id}`

Las personas autorizadas utilizan los tipos de documento `CC`, `TI`, `CE`,
`PASSPORT` y `OTHER`. La funcionalidad usa datos ficticios y no verifica
identidad legal.

## Pruebas

```powershell
go test ./...
```

La implementación está separada en `domain`, `application`,
`infrastructure/persistence/sqlite` y `presentation/http`. Las dependencias se
ensamblan manualmente desde `cmd/api/main.go`.

El despliegue reproducible está descrito en
[../render.yaml](../render.yaml). El backend continúa sin autenticación y solo
debe exponerse públicamente con datos ficticios hasta incorporar protección
administrativa.
