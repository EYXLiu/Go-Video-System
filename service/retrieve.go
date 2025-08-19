package service

import (
	"go-video-system/db"
	"go-video-system/model"
)

func RetrieveVideo(videoID string) (model.Video, error) {
	return db.PostgresGet(videoID)
}
