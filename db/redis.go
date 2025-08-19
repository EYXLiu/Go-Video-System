package db

import (
	"log"
	"os"

	"encoding/json"

	"go-video-system/model"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client

func RedisInit() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	Rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
}

func RedisSet(upload model.UploadSession) error {
	data, err := json.Marshal(upload)
	if err != nil {
		return err
	}
	return Rdb.Set(Ctx, upload.UploadID, data, 0).Err()
}

func RedisIncrement(uploadID string, chunkNum int) error {
	upload, err := RedisGet(uploadID)
	if err != nil {
		return err
	}
	upload.Chunks = chunkNum
	return RedisSet(upload)
}

func RedisGet(uploadID string) (model.UploadSession, error) {
	var upload model.UploadSession

	data, err := Rdb.Get(Ctx, uploadID).Result()
	if err != nil {
		return upload, err
	}
	err = json.Unmarshal([]byte(data), &upload)
	return upload, err
}
