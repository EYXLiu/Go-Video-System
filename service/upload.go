package service

import (
	"mime/multipart"
	"time"

	"go-video-system/db"
	"go-video-system/model"

	"github.com/google/uuid"
)

func UploadInit(userID string, filename string, totalChunks int) (model.UploadSession, error) {
	session := model.UploadSession{
		UploadID:  uuid.New().String(),
		UserID:    userID,
		FileName:  filename,
		Chunks:    totalChunks,
		Uploaded:  0,
		Status:    "STARTED",
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	err := db.RedisSet(session)
	return session, err
}

func UploadChunk(uploadID string, chunkNum int, file *multipart.FileHeader) error {
	f, err := file.Open()
	if err != nil {
		return err
	}
	defer f.Close()

	return db.S3Upload(uploadID, chunkNum, f, file.Size)
}

func UploadComplete(uploadID string) (model.Video, error) {
	upload, err := db.RedisGet(uploadID)
	if err != nil {
		return model.Video{}, err
	}

	path, err := db.S3Merge(uploadID, upload.Chunks)
	if err != nil {
		return model.Video{}, err
	}

	video := model.Video{
		VideoID:     uuid.New().String(),
		UserID:      upload.UserID,
		FileName:    upload.FileName,
		Status:      "PROCESSING",
		Duration:    0,
		Resolutions: map[string]string{},
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}

	err = processVideo(path, &video)
	if err != nil {
		return model.Video{}, err
	}

	video.Status = "READY"
	err = db.PostgresSet(video)
	if err != nil {
		return model.Video{}, err
	}

	return video, nil
}
