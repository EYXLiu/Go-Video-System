package db

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var S3Client *minio.Client
var bucketName string

func S3Init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	bucketName = os.Getenv("S3_VIDEO_BUCKET")

	endpoint := os.Getenv("S3_ENDPOINT")
	accessKey := os.Getenv("S3_ACCESS_KEY")
	secretKey := os.Getenv("S3_SECRET_KEY")
	useSSL := os.Getenv("S3_USE_SSL") == "true"
	region := os.Getenv("S3_REGION")

	S3Client, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
		Region: region,
	})
	if err != nil {
		log.Fatal(err)
	}
	err = S3Client.MakeBucket(Ctx, bucketName, minio.MakeBucketOptions{Region: region})
	if err != nil {
		log.Fatal(err)
	}
}

func S3Upload(uploadID string, chunkNum int, file io.Reader, size int64) error {
	objectName := fmt.Sprintf("uploads/%s/%d", uploadID, chunkNum)
	_, err := S3Client.PutObject(Ctx, bucketName, objectName, file, size, minio.PutObjectOptions{})
	return err
}

func S3Merge(uploadID string, totalChunks int) (string, error) {
	filename := uploadID + ".mp4"
	localPath := filepath.Join(os.TempDir(), filename)

	outFile, err := os.Create(localPath)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	for i := range totalChunks {
		chunk := fmt.Sprintf("uploads/%s/%d", uploadID, i)
		object, err := S3Client.GetObject(context.TODO(), bucketName, chunk, minio.GetObjectOptions{})
		if err != nil {
			return "", nil
		}

		_, err = io.Copy(outFile, object)
		if err != nil {
			return "", nil
		}

		_ = S3Client.RemoveObject(context.TODO(), bucketName, chunk, minio.RemoveObjectOptions{})
	}

	return localPath, nil
}

func S3PresignedURL(key string, expiry time.Duration) (string, error) {
	reqParams := make(map[string][]string)
	presignedURL, err := S3Client.PresignedGetObject(
		Ctx,
		bucketName,
		key,
		expiry,
		reqParams,
	)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}
