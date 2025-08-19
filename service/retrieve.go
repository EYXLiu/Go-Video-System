package service

import (
	"time"

	"go-video-system/db"
	"go-video-system/model"
)

func RetrieveVideo(videoID string) (model.Video, error) {
	video, err := db.PostgresGet(videoID)
	if err != nil {
		return model.Video{}, err
	}

	signedRes := make(map[string]string)
	for res, objectKey := range video.Resolutions {
		url, err := db.S3PresignedURL(objectKey, 15*time.Minute)
		if err != nil {
			return model.Video{}, err
		}
		signedRes[res] = url
	}
	video.Resolutions = signedRes

	return video, nil
}
