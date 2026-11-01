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

El endpoint `GET /api/courses` expone los cursos demo para el formulario del
frontend. En desarrollo se permite CORS desde `http://localhost:4200`.

Endpoints principales:

- `GET /api/courses`
- `GET /api/events`
- `POST /api/events`
- `GET /api/events/{id}`
- `GET /api/requests/{token}`
- `POST /api/requests/{token}/response`

## Pruebas

```powershell
go test ./...
```

La implementación está separada en `domain`, `application`,
`infrastructure/persistence/sqlite` y `presentation/http`. Las dependencias se
ensamblan manualmente desde `cmd/api/main.go`.
