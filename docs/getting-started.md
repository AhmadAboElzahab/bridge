# Getting Started

## Prerequisites

| Tool | Version | Purpose |
|------|---------|---------|
| Go | 1.24+ | Runtime and build |
| PostgreSQL | 15+ | Database |
| Air | latest | Hot reload for development |
| swag | latest | Swagger doc generation |

Install Air and swag:

```bash
go install github.com/air-verse/air@latest
go install github.com/swaggo/swag/cmd/swag@latest
```

---

## Local Setup (without Docker)

### 1. Clone and install dependencies

```bash
git clone https://github.com/AhmadAboElzahab/bridge.git
cd bridge
go mod tidy
```

### 2. Configure environment

```bash
cp .env.exapmle .env
```

Edit `.env` with your local database credentials. See [Environment Variables](environment-variables.md) for every option.

### 3. Run migrations

```bash
go run ./cmd/migrations/main.go
```

### 4. (Optional) Seed the database

```bash
go run ./cmd/seed/main.go
```

### 5. Start the server

With hot reload:
```bash
air
```

Without hot reload:
```bash
go run ./cmd/main.go
```

The server starts on **http://localhost:8080**.

---

## Local Setup (with Docker — development)

Uses `docker-compose.yml` with Air hot reload inside the container.

```bash
cp .env.docker.example .env.docker
# Edit .env.docker with your values

docker compose up --build
```

The app mounts your local directory into the container so code changes trigger a reload automatically.

---

## Regenerate Swagger Docs

Run this from the project root whenever you add or change a handler annotation:

```bash
swag init
```

This updates `docs/swagger.json`, `docs/swagger.yaml`, and `docs/docs.go`.

---

## Project Structure

```
cmd/
  main.go              — server entrypoint (hardcodes :8080)
  migrations/main.go   — GORM AutoMigrate runner
  seed/main.go         — database seeder
internal/
  controllers/         — HTTP handlers (thin layer, no business logic)
  middlewares/         — JWT auth middleware
  models/              — GORM structs
  services/            — business logic and query builders
  utils/               — stateless helpers
  initializers/        — DB and storage client setup
  constants/           — shared constants
  routes/routes.go     — all route registrations
docs/                  — Swagger files + this documentation
```

See [Architecture](architecture.md) for a deeper breakdown.
