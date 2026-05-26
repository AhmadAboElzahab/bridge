package utils

import (
	"bytes"
	"context"
	"errors"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"path/filepath"

	"github.com/AhmadAboElzahab/bridge/internal/initializers"
	"github.com/buckket/go-blurhash"
	"github.com/chai2010/webp"
)

var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
}

const maxFileSize = 5 * 1024 * 1024 // 5 MB

// ProcessImageUpload validates, converts to WebP, generates a BlurHash,
// and uploads the image via the active storage driver.
//
// Parameters:
//   - ctx:    request context (used for upload cancellation)
//   - file:   multipart file header from the request
//   - folder: logical folder name within the storage backend (e.g. "users")
//
// Returns:
//   - publicURL: URL or path where the file can be accessed
//   - blurhash:  BlurHash string for client-side placeholder rendering
//   - error
func ProcessImageUpload(ctx context.Context, file *multipart.FileHeader, folder string) (string, string, error) {
	ext := filepath.Ext(file.Filename)
	if !allowedExtensions[ext] {
		return "", "", errors.New("invalid file type, only JPG, PNG and WebP allowed")
	}

	if file.Size > maxFileSize {
		return "", "", errors.New("file too large, max size is 5 MB")
	}

	src, err := file.Open()
	if err != nil {
		return "", "", errors.New("failed to open file")
	}
	defer src.Close()

	imgData, err := io.ReadAll(src)
	if err != nil {
		return "", "", errors.New("failed to read file")
	}

	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return "", "", errors.New("invalid image format")
	}

	hash, err := blurhash.Encode(4, 3, img)
	if err != nil {
		return "", "", errors.New("failed to generate BlurHash")
	}

	var webpBuf bytes.Buffer
	if err := webp.Encode(&webpBuf, img, &webp.Options{Quality: 80}); err != nil {
		return "", "", errors.New("failed to convert image to WebP")
	}

	publicURL, err := initializers.Storage.Upload(ctx, folder, ".webp", "image/webp", webpBuf.Bytes())
	if err != nil {
		return "", "", err
	}

	return publicURL, hash, nil
}
