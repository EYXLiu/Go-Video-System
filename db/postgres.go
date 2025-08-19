package db

import (
	"log"
	"os"
	"time"

	"encoding/json"

	"go-video-system/model"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var Db *pgxpool.Pool

func PostgresInit() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	pool, err := pgxpool.New(Ctx, os.Getenv("POSTGRES_DSN"))
	if err != nil {
		log.Fatal(err)
	}
	_, err = Db.Exec(Ctx, `
        CREATE TABLE IF NOT EXISTS videos (
            video_id TEXT PRIMARY KEY,
            user_id TEXT NOT NULL,
            filename TEXT,
            status TEXT,
            duration DOUBLE PRECISION,
            resolutions JSONB,
            created_at BIGINT,
            updated_at BIGINT
        );`)
	if err != nil {
		log.Fatal(err)
	}
	Db = pool
}

func PostgresUpdate(videoID string, status string, resolutions map[string]string) error {
	resolutionsJSON, err := json.Marshal(resolutions)
	if err != nil {
		return err
	}

	_, err = Db.Exec(Ctx,
		`UPDATE videos SET status=$1, resolutions=$2, updated_at=$3 WHERE video_id=$4`,
		status,
		resolutionsJSON,
		time.Now().Unix(),
		videoID,
	)
	return err
}

func PostgresSet(video model.Video) error {
	resolutions, err := json.Marshal(video.Resolutions)
	if err != nil {
		return err
	}

	_, err = Db.Exec(Ctx,
		`INSERT INTO videos (
			video_id, user_id, filename, status, duration, resolutions, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		video.VideoID,
		video.UserID,
		video.FileName,
		video.Status,
		video.Duration,
		resolutions,
		video.CreatedAt,
		video.UpdatedAt,
	)
	return err
}

func PostgresGet(videoID string) (model.Video, error) {
	var video model.Video
	var resolutions []byte

	row := Db.QueryRow(Ctx,
		`SELECT * FROM videos WHERE video_id=$1`, videoID,
	)
	err := row.Scan(
		&video.VideoID,
		&video.UserID,
		&video.FileName,
		&video.Status,
		&video.Duration,
		&resolutions,
		&video.CreatedAt,
		&video.UpdatedAt,
	)
	if err != nil {
		return video, err
	}

	err = json.Unmarshal(resolutions, &video.Resolutions)
	return video, err
}
