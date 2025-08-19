package handlers

import (
	"net/http"
	"strconv"

	"go-video-system/service"

	"github.com/gin-gonic/gin"
)

func UploadInitHandler(c *gin.Context) {
	userID := c.PostForm("user_id")
	fileName := c.PostForm("file_name")
	totalChunksStr := c.PostForm("total_chunks")

	totalChunks, err := strconv.Atoi(totalChunksStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid total_chunks"})
		return
	}

	upload, err := service.UploadInit(userID, fileName, totalChunks)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"upload_id": upload.UploadID})
}

func UploadChunkHandler(c *gin.Context) {
	uploadID := c.PostForm("upload_id")
	chunkNumStr := c.PostForm("chunk_num")

	chunkNum, err := strconv.Atoi(chunkNumStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chunk_num"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file missing"})
		return
	}

	err = service.UploadChunk(uploadID, chunkNum, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "chunk uploaded"})
}

func UploadCompleteHandler(c *gin.Context) {
	uploadID := c.PostForm("uploadID")

	video, err := service.UploadComplete(uploadID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"video_id": video.VideoID, "status": video.Status})
}
