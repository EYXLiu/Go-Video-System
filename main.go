package main

import (
	"go-video-system/db"
	"go-video-system/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	db.Init()
	defer db.Close()

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/upload/init", handlers.UploadInitHandler)
	r.POST("/upload/chunk", handlers.UploadChunkHandler)
	r.POST("/upload/complete", handlers.UploadCompleteHandler)

	r.GET("/video/:video_id", handlers.RetrieveVideo)
	r.GET("/videos/:video_id/download", handlers.DownloadVideo)

	r.Run(":8080")
}
