# Environment Variables

The app loads variables from `.env` using `godotenv`. In Docker, `.env.docker` is used instead.

Copy the appropriate example file to get started:

```bash
# Local development
cp .env.exapmle .env

# Docker
cp .env.docker.example .env.docker
```

---

## Server

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `PORT` | No | — | Declared but not used — server hardcodes `:8080` |
| `GIN_MODE` | No | `debug` | Set to `release` in production to disable debug output |

---

## Database

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `DB_URL` | **Yes** | — | PostgreSQL DSN |

**Format:**
```
host=localhost user=myuser password=mypassword dbname=bridge port=5432 sslmode=disable
```

In Docker, `host=` must match the Postgres service name (e.g. `host=db`).

**Connection pool** (configured in code, not env):
- Max open connections: 25
- Max idle connections: 5
- Connection max lifetime: 5 minutes

**Docker-only Postgres container variables** (used by the `db` service in docker-compose):

| Variable | Required | Description |
|----------|----------|-------------|
| `POSTGRES_USER` | Yes | Must match `user=` in `DB_URL` |
| `POSTGRES_PASSWORD` | Yes | Must match `password=` in `DB_URL` |
| `POSTGRES_DB` | Yes | Must match `dbname=` in `DB_URL` |

---

## Authentication

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `JWT_SECRET` | **Yes** | — | Secret key used to sign and verify JWT tokens |

Generate a strong secret:
```bash
openssl rand -base64 48
```

**Token behaviour:**
- Algorithm: HS256
- Expiry: 24 hours from issue time
- Header format: `Authorization: Bearer <token>`

---

## CORS

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `CORS_ALLOWED_ORIGINS` | No | `http://localhost:3000` | Comma-separated list of allowed origins |

**Examples:**
```ini
# Single origin
CORS_ALLOWED_ORIGINS=http://localhost:3000

# Multiple origins
CORS_ALLOWED_ORIGINS=https://app.yourdomain.com,https://admin.yourdomain.com
```

---

## Storage

Controls where uploaded images are stored. See [Storage](storage.md) for full setup guides.

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `STORAGE_DRIVER` | No | `local` | Storage backend: `local`, `s3`, or `r2` |
| `STORAGE_ACCESS_KEY_ID` | s3/r2 only | — | Access key for S3 or R2 |
| `STORAGE_SECRET_ACCESS_KEY` | s3/r2 only | — | Secret key for S3 or R2 |
| `STORAGE_ENDPOINT` | r2 only | — | Custom endpoint URL for R2 or other S3-compatible services |
| `STORAGE_BUCKET` | s3/r2 only | — | Bucket name |
| `STORAGE_REGION` | s3/r2 only | `auto` | Region — use `auto` for R2, region name for S3 |
| `STORAGE_PUBLIC_URL` | s3/r2 only | — | Public base URL for uploaded files (no trailing slash) |

### Driver: `local`

No extra variables needed. Files are saved to `./storage/` and served at `/storage/*`.

```ini
STORAGE_DRIVER=local
```

### Driver: `r2` (Cloudflare R2)

```ini
STORAGE_DRIVER=r2
STORAGE_ACCESS_KEY_ID=your_r2_access_key_id
STORAGE_SECRET_ACCESS_KEY=your_r2_secret_access_key
STORAGE_ENDPOINT=https://<account-id>.r2.cloudflarestorage.com
STORAGE_BUCKET=bridge-media
STORAGE_REGION=auto
STORAGE_PUBLIC_URL=https://media.yourdomain.com
```

### Driver: `s3` (AWS S3)

```ini
STORAGE_DRIVER=s3
STORAGE_ACCESS_KEY_ID=your_aws_access_key_id
STORAGE_SECRET_ACCESS_KEY=your_aws_secret_access_key
STORAGE_ENDPOINT=
STORAGE_BUCKET=bridge-media
STORAGE_REGION=us-east-1
STORAGE_PUBLIC_URL=https://bridge-media.s3.us-east-1.amazonaws.com
```

---

## Full Example

```ini
# Server
PORT=8080
GIN_MODE=release

# Database
DB_URL=host=localhost user=bridge password=secret dbname=bridge port=5432 sslmode=disable

# Auth
JWT_SECRET=your_very_long_random_secret_here

# CORS
CORS_ALLOWED_ORIGINS=https://app.yourdomain.com

# Storage
STORAGE_DRIVER=r2
STORAGE_ACCESS_KEY_ID=abc123
STORAGE_SECRET_ACCESS_KEY=xyz789
STORAGE_ENDPOINT=https://abc123.r2.cloudflarestorage.com
STORAGE_BUCKET=bridge-media
STORAGE_REGION=auto
STORAGE_PUBLIC_URL=https://media.yourdomain.com
```
