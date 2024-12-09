package main

import (
	"os"

	"github.com/ravjot07/docraptor-backend/config"
	"github.com/ravjot07/docraptor-backend/handlers"
	"github.com/ravjot07/docraptor-backend/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	config.InitDB()

	router.Use(middleware.LoggerMiddleware())

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to DocRaptor Backend",
		})
	})

	// api routes
	api := router.Group("/api")
	{
		api.GET("/doc", handlers.GetDocsHandler)
		api.GET("/docs/:id", handlers.GetDocByIDHandler)
		api.POST("/upload", handlers.UploadDocHandler)
		api.PUT("/docs/:id", handlers.UpdateDocHandler)
		api.DELETE("/docs/:id", handlers.DeleteDocHandler)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)

}
