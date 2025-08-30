package service

import (
	"fmt"
	"go-video-system/db"
	"go-video-system/model"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func getDuration(filePath string) float64 {
	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries",
		"format=duration", "-of", "default=noprint_wrappers=1:nokey=1", filePath)

	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return 0
	}

	duration, err := strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
	if err != nil {
		fmt.Println(err)
		return 0
	}

	return float64(duration)
}

func getResolutions(filePath string) (map[string]string, error) {
	assets := make(map[string]string)
	outputDir := filepath.Dir(filePath)
	base := filepath.Base(filePath)

	thumbPath := filepath.Join(outputDir, "thumbnail.jpg")
	cmd := exec.Command("ffmpeg", "-i", filePath, "-frames:v", "1", "-q:v", "2", thumbPath)
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("thumbnail generation failed: %v", err)
	}
	assets["thumbnail"] = thumbPath

	resolutions := []string{"720", "480", "360"}
	for _, res := range resolutions {
		outPath := filepath.Join(outputDir, fmt.Sprintf("%s_%sp.mp4", strings.TrimSuffix(base, ".mp4"), res))
		cmd := exec.Command("ffmpeg", "-i", filePath, "-vf", fmt.Sprintf("scale=%s:-2", res), outPath)
		if err := cmd.Run(); err != nil {
			fmt.Println(res, ":", err)
			continue
		}
		assets[res] = outPath
	}

	return assets, nil
}

func processVideo(filePath string, video *model.Video) error {
	video.Duration = getDuration(filePath)
	assets, err := getResolutions(filePath)
	if err != nil {
		return err
	}
	resolutions := make(map[string]string)

	if thumb, ok := assets["thumbnail"]; ok {
		url, err := db.S3UploadFile(thumb, fmt.Sprintf("videos/%s/thumbnail.jpg", video.VideoID))
		if err != nil {
			return err
		}
		resolutions["thumbnail"] = url
	}

	for res, localFile := range assets {
		if res == "thumbnail" {
			continue
		}

		url, err := db.S3UploadFile(localFile, fmt.Sprintf("videos/%s/%sp.mp4", video.VideoID, res))
		if err != nil {
			fmt.Println("upload error:", err)
			continue
		}

		resolutions[res] = url
	}

	video.Resolutions = resolutions

	return nil
}
