# Go-Video-System
Tech Stack: Go, Gin, Postgres, MinIO (AWS S3 compatible), Redis, Docker, FFmpeg, ReactJS (for frontend API calls)

# About
Implemented a Go API backend using Gin to upload videos using chunking and ReactJS API calls <br/>
Used Redis for chunked uploads for tracking upload sessions <br/>
Metadata for files (location, resolutions, length) is saved in Postgres <br/>
Uploaded files are saved to MinIO (S3-compatible endpoints, just switch endpoint from localhost to s3 host) <br/>
Files are automatically streamed through the resolution link using MinIO/S3 206 partial content api calls <br/>
FFmpeg is used for thumbnail and resolution generation 
- Should switch to HLS for auto resolution switching in the future 

# Endpoints
Follow the `frontend_api_call.js` file for upload video frontend call samples <br/>
Designed for scalable video handling - following chunking method <br/> 
**POST** `/upload/init` initialize a new upload session <br/>
**POST** `/upload/chunk` upload a single video chunk <br/>
**POST** `/upload/complete` merge chunks, generate thumbnail and resolutions, upload metadata to Postgres <br/>
**GET** `/video/:video_id` get metadata for videos -> can access video and thumbnail through the resolutions links <br/>

# Frontend
Follow the below for a sample <video> div <br/>
```
<video id="video_player" width="640" height="360" poster="" controls preload="none">
  Your browser does not support the video tag.
</video>
<button id="video_play_button">Play Video</button>
```
