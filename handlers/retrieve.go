package handlers

import (
	"fmt"
	"go-video-system/service"
	"io"
	"net/http"

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

	resp, err := http.Get(filePath)
	if err != nil || resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch video"})
		return
	}
	defer resp.Body.Close()

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+videoID+"_"+res+".mp4")
	c.Header("Content-Type", "video/mp4")
	c.Header("Content-Length", fmt.Sprintf("%d", resp.ContentLength))

	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to stream video"})
		return
	}
}
