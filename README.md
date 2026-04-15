# owner-api-proxy

Go-based API proxy service with Cobra CLI, Gin HTTP server, and modular internal structure.

## Project Structure

```text
owner-api-proxy/
├── README.md
└── src/
    ├── app.example.yaml              # Example nested configuration
    ├── app.yaml                      # Local runtime configuration
    ├── go.mod
    ├── go.sum
    ├── cmd/
    │   ├── main.go                   # CLI entrypoint
    │   └── command/
    │       ├── init.go               # Root Cobra command registration
    │       ├── db.go                 # db command group
    │       └── server.go             # server command wrapper
    ├── pkg/
    │   └── cli/
    │       ├── db/
    │       │   ├── model.go          # db model command
    │       │   ├── pull.go           # db pull command
    │       │   ├── migrate.go        # db migrate command
    │       │   └── test/
    │       └── server/
    │           └── server.go         # Starts Gin server via CLI
    └── internal/
        ├── adapter/
        │   └── http/
        │       ├── routes.go         # Top-level HTTP route registration
        │       └── v1/
        │           ├── controller/
        │           │   └── hello.go  # Example handler
        │           └── routes/
        │               └── hello.go  # v1 route bindings
        ├── config/
        │   ├── config.go             # App config loader (Viper)
        │   ├── gin.go                # Gin engine + address config
        │   ├── logrus.go             # Logger config (reserved)
        │   └── postgres.go           # Postgres config (reserved)
        └── database/
            ├── migration/
            │   └── migrate.go
            └── model/
                └── user.go
```

## Configuration Format

Use nested config in `src/app.yaml`:

```yaml
server:
  host: ""
  port: "8080"

postgres:
  host: "localhost"
  user: "postgres"
  password: "postgres"
  port: "5432"
  database: "owner_api_proxy"
```

You can copy `src/app.example.yaml` as a starting point.

## Run

From `src/` directory:

1. Install Package
```bash
go mod tidy
```

2. Run Server
```bash
go run cmd/main.go server
```

Useful commands:

```bash
go run cmd/main.go help
go run cmd/main.go db pull
go run cmd/main.go db migrate
go run cmd/main.go db model --name User
```
