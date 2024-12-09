package handlers

import (
	"net/http"
	"os"

	"github.com/ravjot07/docraptor-backend/config"
	"github.com/ravjot07/docraptor-backend/models"
	"github.com/ravjot07/docraptor-backend/utils"

	"github.com/gin-gonic/gin"
)

// GetDocsHandler fetches all the documentation pages
func GetDocsHandler(c *gin.Context) {
	var docs []models.Doc
	if err := config.DB.Find(&docs).Error; err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve documents")
		return
	}

	// Create a list of document summaries
	var docList []map[string]string
	for _, doc := range docs {
		docList = append(docList, map[string]string{
			"id":    doc.ID,
			"title": doc.Title,
		})
	}

	c.JSON(http.StatusOK, docList)
}

// GetDocByIDHandler retrieves a specific documentation page by ID
func GetDocByIDHandler(c *gin.Context) {
	id := c.Param("id")
	var doc models.Doc
	if err := config.DB.Where("id = ?", id).First(&doc).Error; err != nil {
		utils.RespondWithError(c, http.StatusNotFound, "Document not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":        doc.ID,
		"title":     doc.Title,
		"content":   doc.Content,  // Raw Markdown
		"file_path": doc.FilePath, // Optional
	})
}

// UpdateDocHandler updates a documentation page by ID
func UpdateDocHandler(c *gin.Context) {
	id := c.Param("id")
	content := c.PostForm("content")
	title := c.PostForm("title")

	if content == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Content is required")
		return
	}

	var doc models.Doc
	if err := config.DB.Where("id = ?", id).First(&doc).Error; err != nil {
		utils.RespondWithError(c, http.StatusNotFound, "Document not found")
		return
	}

	// Update fields
	doc.Content = content
	if title != "" {
		doc.Title = title
	}

	if err := config.DB.Save(&doc).Error; err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to update document")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Document updated successfully",
	})
}

// DeleteDocHandler deletes a documentation page by ID
func DeleteDocHandler(c *gin.Context) {
	id := c.Param("id")
	var doc models.Doc
	if err := config.DB.Where("id = ?", id).First(&doc).Error; err != nil {
		utils.RespondWithError(c, http.StatusNotFound, "Document not found")
		return
	}

	if err := config.DB.Delete(&doc).Error; err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to delete document")
		return
	}

	// Optional: Remove the file from disk
	if doc.FilePath != "" {
		err := os.Remove(doc.FilePath)
		if err != nil {
			utils.RespondWithError(c, http.StatusInternalServerError, "Failed to delete file from disk")
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Document deleted successfully",
	})
}
