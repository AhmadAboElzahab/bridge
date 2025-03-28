package utils

import (
	"bytes"
	"errors"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/buckket/go-blurhash"
	"github.com/chai2010/webp"
	"github.com/google/uuid"
)

// Confirm this is the correct import

// Allowed image extensions
var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
}

// Max file size (5MB)
const maxFileSize = 5 * 1024 * 1024

// ProcessImageUpload handles image validation, conversion to WebP, saving, and BlurHash generation.
func ProcessImageUpload(file *multipart.FileHeader, uploadDir string) (string, string, error) {

	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", "", errors.New("invalid file type, only JPG and PNG allowed")
	}

	// Validate file extension
	ext := filepath.Ext(file.Filename)
	if !allowedExtensions[ext] {
		return "", "", errors.New("invalid file type, only JPG and PNG allowed")
	}

	// Validate file size
	if file.Size > maxFileSize {
		return "", "", errors.New("file too large, max size is 5MB")
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		return "", "", errors.New("failed to open file")
	}
	defer src.Close()

	// Read file into memory
	imgData, err := io.ReadAll(src)
	if err != nil {
		return "", "", errors.New("failed to read file")
	}

	// Decode image
	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return "", "", errors.New("invalid image format")
	}

	// Generate BlurHash
	hash, err := blurhash.Encode(4, 3, img)
	if err != nil {
		return "", "", errors.New("failed to generate BlurHash")
	}

	// Convert image to WebP
	var webpBuffer bytes.Buffer
	err = webp.Encode(&webpBuffer, img, &webp.Options{Quality: 80})
	if err != nil {
		return "", "", errors.New("failed to convert image to WebP")
	}

	// Generate unique filename with .webp extension
	newFileName := uuid.New().String() + ".webp"
	savePath := filepath.Join(uploadDir, newFileName)

	// Save WebP file to disk
	if err := os.WriteFile(savePath, webpBuffer.Bytes(), 0644); err != nil {
		return "", "", errors.New("failed to save file")
	}

	// Return saved file path and BlurHash
	return savePath, hash, nil
}
