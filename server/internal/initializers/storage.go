package initializers

import (
	"log"
	"os"
	"strings"

	"github.com/AhmadAboElzahab/bridge/internal/storage"
)

// Storage is the active storage driver. Set by ConnectStorage().
var Storage storage.Driver

// ConnectStorage initialises the storage backend based on STORAGE_DRIVER env var.
//
// STORAGE_DRIVER=local  (default) — saves files to ./storage/ on disk
// STORAGE_DRIVER=s3               — AWS S3
// STORAGE_DRIVER=r2               — Cloudflare R2 (S3-compatible)
//
// Required env vars for s3/r2:
//
//	STORAGE_ACCESS_KEY_ID
//	STORAGE_SECRET_ACCESS_KEY
//	STORAGE_BUCKET
//	STORAGE_PUBLIC_URL      — public base URL for uploaded files
//	STORAGE_REGION          — defaults to "auto" (correct for R2; use region name for S3)
//	STORAGE_ENDPOINT        — custom endpoint; required for R2, omit for S3
func ConnectStorage() {
	driver := strings.ToLower(os.Getenv("STORAGE_DRIVER"))

	switch driver {
	case "s3", "r2":
		s3Driver, err := storage.NewS3Driver(
			os.Getenv("STORAGE_ACCESS_KEY_ID"),
			os.Getenv("STORAGE_SECRET_ACCESS_KEY"),
			envOrDefault("STORAGE_REGION", "auto"),
			os.Getenv("STORAGE_ENDPOINT"),
			os.Getenv("STORAGE_BUCKET"),
			os.Getenv("STORAGE_PUBLIC_URL"),
		)
		if err != nil {
			log.Fatal("Failed to initialise S3/R2 storage:", err)
		}
		Storage = s3Driver
		log.Printf("Storage: %s driver (bucket: %s)\n", driver, os.Getenv("STORAGE_BUCKET"))

	default:
		Storage = &storage.LocalDriver{BasePath: "./storage"}
		log.Println("Storage: local driver (./storage)")
	}
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
