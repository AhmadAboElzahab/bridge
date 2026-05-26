package storage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

// LocalDriver saves files to a directory on disk and returns a relative URL path
// that the Gin static file server can serve.
type LocalDriver struct {
	// BasePath is the root directory where files are stored, e.g. "./storage"
	BasePath string
}

func (d *LocalDriver) Upload(_ context.Context, folder, ext, _ string, data []byte) (string, error) {
	dir := filepath.Join(d.BasePath, folder)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	filename := uuid.New().String() + ext
	fullPath := filepath.Join(dir, filename)

	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return "storage/" + folder + "/" + filename, nil
}
