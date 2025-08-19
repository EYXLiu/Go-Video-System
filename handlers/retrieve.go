package handlers

import (
    "net/http"
    "go-video-system/service"
    "github.com/gin-gonic/gin"
)

func RetrieveVideo(c *gin.Context) {
	videoID := c.Param("video_id")
	video, err := service.RetrieveVideo(videoID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "video not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"video": video})
}