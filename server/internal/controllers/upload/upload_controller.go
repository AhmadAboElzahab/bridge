package upload

import (
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/AhmadAboElzahab/bridge/internal/initializers"
	"github.com/AhmadAboElzahab/bridge/internal/utils"
	"github.com/gin-gonic/gin"
)

const maxUploadSize = 100 * 1024 * 1024 // 100 MB

var imageExts = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".webp": true,
}

type UploadController struct{}

func NewUploadController() *UploadController {
	return &UploadController{}
}

// Upload godoc
// @Summary Upload a file
// @Description Upload an image, video, or document. Images are converted to WebP and a blurhash is returned.
// @Tags upload
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "File to upload"
// @Param folder query string false "Storage folder (default: uploads)"
// @Success 200 {object} object{url=string,blurhash=string}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/upload [post]
func (uc *UploadController) Upload(ctx *gin.Context) {
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, maxUploadSize)

	file, err := ctx.FormFile("file")
	if err != nil {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "No file provided", err.Error())
		return
	}

	folder := ctx.DefaultQuery("folder", "uploads")
	folder = strings.Trim(folder, "/")

	ext := strings.ToLower(filepath.Ext(file.Filename))

	// Images: convert to WebP and generate blurhash
	if imageExts[ext] {
		url, hash, err := utils.ProcessImageUpload(ctx.Request.Context(), file, folder)
		if err != nil {
			utils.ErrorJSON(ctx, http.StatusBadRequest, "Image upload failed", err.Error())
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"url": url, "blurhash": hash})
		return
	}

	// All other files stored as-is
	src, err := file.Open()
	if err != nil {
		utils.ErrorJSON(ctx, http.StatusInternalServerError, "Failed to open file", err.Error())
		return
	}
	defer src.Close()

	data, err := io.ReadAll(src)
	if err != nil {
		utils.ErrorJSON(ctx, http.StatusInternalServerError, "Failed to read file", err.Error())
		return
	}

	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = mime.TypeByExtension(ext)
	}
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	url, err := initializers.Storage.Upload(ctx.Request.Context(), folder, ext, contentType, data)
	if err != nil {
		utils.ErrorJSON(ctx, http.StatusInternalServerError, "Upload failed", err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"url": url})
}
