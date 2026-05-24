# Deployment

---

## Development (Docker with hot reload)

Uses `docker-compose.yml` + `Dockerfile` + [Air](https://github.com/air-verse/air).

Your local directory is mounted into the container so code changes trigger an automatic rebuild.

```bash
# First time
cp .env.docker.example .env.docker
# Edit .env.docker

docker compose up --build

# Subsequent starts
docker compose up
```

**Services:**
- `app` — Go API on `:8080` (hot-reload via Air)
- `db` — PostgreSQL 15 on `:5432`

---

## Production (Docker with Nginx + SSL)

Uses `docker-compose.prod.yml` + `Dockerfile.prod` + Nginx + Certbot (Let's Encrypt).

### First deploy

```bash
# 1. Copy and fill in env file
cp .env.docker.example .env.docker
# Set GIN_MODE=release, strong JWT_SECRET, real CORS origins, DB creds, storage config

# 2. Build and start
docker compose -f docker-compose.prod.yml up --build -d

# 3. Run migrations
docker exec bridge_app ./migrate

# 4. (Optional) Seed
docker exec bridge_app ./seed
```

### Services in production

| Service | Container | Description |
|---------|-----------|-------------|
| `app` | `bridge_app` | Go API binary (no hot reload) |
| `db` | `bridge_db` | PostgreSQL 15 with health check |
| `nginx` | `bridge_nginx` | Reverse proxy, ports 80 + 443 |
| `certbot` | `bridge_certbot` | Auto-renews Let's Encrypt certificates |

### Volumes in production

| Volume | Mounted at | Purpose |
|--------|-----------|---------|
| `db_data` | `/var/lib/postgresql/data` | Postgres data (persistent) |
| `storage_data` | `/app/storage` | Uploaded files when using `local` driver |
| `certbot_www` | `/var/www/certbot` | ACME challenge files |
| `certbot_conf` | `/etc/letsencrypt` | SSL certificates |

### Entrypoint

`entrypoint.sh` runs on container start:
```sh
./migrate   # runs GORM AutoMigrate
./bridge    # starts the API server
```

---

## Updating the production deployment

```bash
# Pull latest code
git pull

# Rebuild and restart app only (zero-downtime if behind Nginx)
docker compose -f docker-compose.prod.yml up --build -d app
```

---

## Nginx Configuration

Nginx sits in front of the Go app and handles:
- SSL termination (443 → app:8080)
- HTTP → HTTPS redirect
- Static file serving (optional)

Config file: `nginx/nginx.conf`

---

## SSL Certificate Renewal

The `certbot` container polls every 12 hours and renews certificates automatically when they are within 30 days of expiry. No manual action needed.

To force a renewal:
```bash
docker exec bridge_certbot certbot renew --force-renewal
docker exec bridge_nginx nginx -s reload
```

---

## Useful Commands

```bash
# View logs
docker compose -f docker-compose.prod.yml logs -f app

# Open a shell in the app container
docker exec -it bridge_app sh

# Run migrations manually
docker exec bridge_app ./migrate

# Check DB connection
docker exec -it bridge_db psql -U myuser -d bridge

# Restart a single service
docker compose -f docker-compose.prod.yml restart app
```

---

## Checklist Before Going Live

- [ ] `GIN_MODE=release` in `.env.docker`
- [ ] Strong `JWT_SECRET` (min 32 random characters)
- [ ] `CORS_ALLOWED_ORIGINS` set to your actual frontend domain
- [ ] `DB_URL` pointing to production DB with strong password
- [ ] Storage driver configured (`r2` or `s3` recommended for production)
- [ ] SSL certificate obtained via Certbot
- [ ] Postgres data volume is backed up regularly
- [ ] `storage_data` volume is backed up (or using cloud storage — no backup needed)
