# Storage

The API stores uploaded images (maid and user avatars) via a pluggable storage driver selected at startup through the `STORAGE_DRIVER` environment variable.

All images are:
- Validated (JPG, PNG, WebP only — max 5 MB)
- Converted to **WebP** at 80% quality
- Given a **BlurHash** string for client-side blur-up placeholders
- Uploaded to the active driver under a UUID filename

---

## Drivers

### `local` (default)

Files are saved to `./storage/{folder}/` on disk. The Gin server exposes them at `/storage/*`.

```ini
STORAGE_DRIVER=local
```

**Pros:** zero configuration, works offline  
**Cons:** files are lost if the container restarts (unless a Docker volume is mounted), not suitable for multi-instance deployments

**Docker volume mount** (already configured in `docker-compose.prod.yml`):
```yaml
volumes:
  - storage_data:/app/storage
```

---

### `r2` — Cloudflare R2

S3-compatible object storage with **no egress fees**. Recommended for production.

```ini
STORAGE_DRIVER=r2
STORAGE_ACCESS_KEY_ID=your_r2_access_key_id
STORAGE_SECRET_ACCESS_KEY=your_r2_secret_access_key
STORAGE_ENDPOINT=https://<account-id>.r2.cloudflarestorage.com
STORAGE_BUCKET=bridge-media
STORAGE_REGION=auto
STORAGE_PUBLIC_URL=https://media.yourdomain.com
```

**Setup steps:**
1. Go to Cloudflare dashboard → R2 → Create bucket
2. Create an API token with Object Read & Write permissions
3. Set your bucket to **Public** access (or use a custom domain)
4. Copy the endpoint URL from the bucket settings (format: `https://<account-id>.r2.cloudflarestorage.com`)
5. For `STORAGE_PUBLIC_URL`, use either the `r2.dev` public URL or your custom domain

---

### `s3` — AWS S3

```ini
STORAGE_DRIVER=s3
STORAGE_ACCESS_KEY_ID=your_aws_access_key_id
STORAGE_SECRET_ACCESS_KEY=your_aws_secret_access_key
STORAGE_ENDPOINT=
STORAGE_BUCKET=bridge-media
STORAGE_REGION=us-east-1
STORAGE_PUBLIC_URL=https://bridge-media.s3.us-east-1.amazonaws.com
```

**Setup steps:**
1. Create an S3 bucket in the AWS console
2. Create an IAM user with `s3:PutObject` and `s3:GetObject` permissions on the bucket
3. Generate an access key for the IAM user
4. Configure bucket policy to allow public reads (or use CloudFront)
5. Leave `STORAGE_ENDPOINT` empty — the SDK uses the default AWS endpoint

---

## How Uploads Work

```
POST /api/auth/signup (multipart avatar)
         │
         ▼
  ProcessImageUpload()          internal/utils/image_utils.go
  ├── validate extension + size
  ├── decode image
  ├── generate BlurHash (4×3 components)
  ├── encode to WebP (quality 80)
  └── initializers.Storage.Upload(ctx, "users", webpBytes)
            │
            ├── LocalDriver  → ./storage/users/uuid.webp
            │                  returns "storage/users/uuid.webp"
            │
            └── S3Driver     → bucket/users/uuid.webp
                               returns "https://media.domain.com/users/uuid.webp"
```

The returned public URL (or relative path for local) is stored in `users.avatar`. The BlurHash is stored in `users.blurhash`.

---

## Switching Drivers

No code changes are required. Change `STORAGE_DRIVER` (and the related credentials) in `.env` and restart the server.

Files already uploaded to the previous driver are **not migrated automatically**. If you switch from local to cloud, existing local paths stored in the database will become broken links. Plan a migration script if needed.

---

## Adding a New Driver

1. Create `internal/storage/yourdriver.go` implementing the `storage.Driver` interface:

```go
type Driver interface {
    Upload(ctx context.Context, folder string, data []byte) (publicURL string, err error)
}
```

2. Add a case to the `switch` in `internal/initializers/storage.go`
3. Document the required env vars in `docs/environment-variables.md`
