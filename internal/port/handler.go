package port

import (
	"net/http"
	"os"
	"sarva/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	fileService *service.FileService
}

func NewHandler(fileService *service.FileService) *Handler {
	return &Handler{fileService: fileService}
}

func (h *Handler) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed"})
		return
	}

	filePath := "./" + file.Filename
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "File save failed"})
		return
	}

	err = h.fileService.UploadFile(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "File processing failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully"})
	os.Remove(filePath)
}
