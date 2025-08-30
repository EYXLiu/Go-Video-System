package handlers

import (
	"go-video-system/service"
	"net/http"
	"path/filepath"

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

func DownloadVideo(c *gin.Context) {
	videoID := c.Param("video_id")
	res := c.Query("res")
	video, err := service.RetrieveVideo(videoID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "video not found"})
		return
	}

	filePath, ok := video.Resolutions[res]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "resolution not available"})
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+filepath.Base(filePath))
	c.Header("Content-Type", "video/mp4")
	c.File(filePath)
}
