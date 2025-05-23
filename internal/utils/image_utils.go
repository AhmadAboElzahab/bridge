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

// allowedExtensions defines the permitted image file extensions for upload.
var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
}

// maxFileSize defines the maximum file size allowed for upload (5MB).
const maxFileSize = 5 * 1024 * 1024

// ProcessImageUpload handles image upload, validation, and conversion to WebP.
// It also generates a BlurHash for the uploaded image.
//
// This function performs the following:
//   - Validates the file type and size
//   - Decodes the image and encodes it into WebP format
//   - Generates a BlurHash string from the image
//   - Saves the converted WebP file to the specified directory
//
// Parameters:
//   - file: the multipart uploaded file header
//   - uploadDir: the destination directory to save the processed image
//
// Returns:
//   - savePath: full path where the WebP image is saved
//   - hash: generated BlurHash string representing the image preview
//   - error: any error encountered during processing (nil if successful)
func ProcessImageUpload(file *multipart.FileHeader, uploadDir string) (string, string, error) {
	// Ensure the directory exists
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", "", errors.New("invalid file type, only JPG and PNG allowed")
	}

	// Validate extension
	ext := filepath.Ext(file.Filename)
	if !allowedExtensions[ext] {
		return "", "", errors.New("invalid file type, only JPG and PNG allowed")
	}

	// Validate file size
	if file.Size > maxFileSize {
		return "", "", errors.New("file too large, max size is 5MB")
	}

	// Read file into memory
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

	// Save file with new UUID name
	newFileName := uuid.New().String() + ".webp"
	savePath := filepath.Join(uploadDir, newFileName)

	if err := os.WriteFile(savePath, webpBuffer.Bytes(), 0644); err != nil {
		return "", "", errors.New("failed to save file")
	}

	return savePath, hash, nil
}
