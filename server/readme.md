# Bridge API

UAE domestic service marketplace — back-office management API.

**Stack:** Go 1.24 · Gin · GORM · PostgreSQL · JWT · Swagger · Docker

---

## Quick Start

```bash
# Install dependencies
go mod tidy

# Configure environment
cp .env.exapmle .env   # then edit .env

# Run migrations
go run ./cmd/migrations/main.go

# Start with hot reload
air
```

Server runs on **http://localhost:8080**  
Swagger UI: **http://localhost:8080/swagger/index.html**

---

## Documentation

All documentation is in the [`docs/`](docs/) folder:

| Document | Description |
|----------|-------------|
| [Getting Started](docs/getting-started.md) | Prerequisites, setup, running locally |
| [Environment Variables](docs/environment-variables.md) | Every env var with options and defaults |
| [Architecture](docs/architecture.md) | Project structure, layers, route map |
| [API Reference](docs/api-reference.md) | All endpoints documented |
| [Storage](docs/storage.md) | Local disk, Cloudflare R2, AWS S3 |
| [Tab & Filter System](docs/tab-filter-system.md) | Per-user views, columns, and filters |
| [Deployment](docs/deployment.md) | Docker dev and production deployment |

---

## Module

```
github.com/AhmadAboElzahab/bridge
```
