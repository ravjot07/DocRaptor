package handlers

import (
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ravjot07/docraptor-backend/config"
	"github.com/ravjot07/docraptor-backend/models"
	"github.com/ravjot07/docraptor-backend/utils"

	"github.com/gin-gonic/gin"
)

// UploadDocHandler handles uploading a new documentation page via a Markdown file
func UploadDocHandler(c *gin.Context) {
	// Retrieve the file from the request
	file, err := c.FormFile("file")
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "File is required")
		fmt.Println("Error retrieving file:", err)
		return
	}

	fmt.Println("File retrieved:", file.Filename)

	// Validate the file extension
	if !isMarkdownFile(file.Filename) {
		utils.RespondWithError(c, http.StatusBadRequest, "Only Markdown (.md) files are allowed")
		fmt.Println("Invalid file extension:", file.Filename)
		return
	}

	fmt.Println("File extension validated.")

	// Open the uploaded file
	openedFile, err := file.Open()
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to open the uploaded file")
		fmt.Println("Error opening file:", err)
		return
	}
	defer openedFile.Close()

	fmt.Println("File opened successfully.")

	// Read the file content
	fileBytes, err := ioutil.ReadAll(openedFile)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to read the uploaded file")
		fmt.Println("Error reading file:", err)
		return
	}
	content := string(fileBytes)

	fmt.Println("File content read successfully.")

	// Extract title from the first Markdown heading
	title := extractTitle(content)
	if title == "" {
		title = deriveIDFromFilename(file.Filename) // Fallback to ID if no title found
		fmt.Println("No title found. Using ID as title:", title)
	} else {
		fmt.Println("Title extracted:", title)
	}

	// Derive ID from the filename (excluding the extension)
	id := deriveIDFromFilename(file.Filename)
	fmt.Println("Document ID derived:", id)

	// Check if a document with the same ID already exists
	var existingDoc models.Doc
	if err := config.DB.Where("id = ?", id).First(&existingDoc).Error; err == nil {
		utils.RespondWithError(c, http.StatusConflict, "Document with this ID already exists")
		fmt.Println("Duplicate document ID:", id)
		return
	}

	fmt.Println("No duplicate document found. Proceeding to save.")

	// Save the file to disk
	savedFilePath, err := saveFileToDisk(c, file)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to save the file")
		fmt.Println("Error saving file to disk:", err)
		return
	}

	fmt.Println("File saved to disk at:", savedFilePath)

	// Create a new document record in the database
	newDoc := models.Doc{
		ID:       id,
		Title:    title,
		Content:  content,
		FilePath: savedFilePath,
	}

	if err := config.DB.Create(&newDoc).Error; err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to create document in the database")
		fmt.Println("Error creating document in DB:", err)
		return
	}

	fmt.Println("Document record created successfully in DB.")

	c.JSON(http.StatusCreated, gin.H{
		"message": "Document uploaded successfully",
		"id":      newDoc.ID,
	})
}

// Helper function to check if the uploaded file is a Markdown file
func isMarkdownFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".md"
}

// Helper function to derive the document ID from the filename
func deriveIDFromFilename(filename string) string {
	return strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
}

// Helper function to extract title from Markdown content
func extractTitle(content string) string {
	// Regex to find first H1 or H2 heading
	re := regexp.MustCompile(`(?m)^(#|##)\s+(.*)$`)
	matches := re.FindStringSubmatch(content)
	if len(matches) >= 3 {
		return matches[2]
	}
	return ""
}

// Function to save the uploaded file to disk
func saveFileToDisk(c *gin.Context, fileHeader *multipart.FileHeader) (string, error) {
	// Define the upload directory
	uploadDir := "./uploads"

	// Create the directory if it doesn't exist
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err := os.Mkdir(uploadDir, 0755)
		if err != nil {
			return "", fmt.Errorf("failed to create upload directory: %v", err)
		}
		fmt.Println("Created uploads directory.")
	}

	// Sanitize the filename
	safeFilename := sanitizeFilename(fileHeader.Filename)
	filePath := filepath.Join(uploadDir, safeFilename)

	// Save the uploaded file to the specified path
	if err := c.SaveUploadedFile(fileHeader, filePath); err != nil {
		return "", fmt.Errorf("failed to save uploaded file: %v", err)
	}

	fmt.Println("File saved at:", filePath)

	return filePath, nil
}

// Function to sanitize filenames
func sanitizeFilename(filename string) string {
	// Remove any path traversal characters
	sanitized := filepath.Base(filename)
	// Optionally, remove or replace unwanted characters
	re := regexp.MustCompile(`[^a-zA-Z0-9_\-\.]`)
	sanitized = re.ReplaceAllString(sanitized, "_")
	return sanitized
}
