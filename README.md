# Go-Video-System
Tech Stack: Go, Gin, Postgres, MinIO (AWS S3 compatible), Redis, Docker, FFmpeg, ReactJS (for frontend API calls)

# About
Implemented a Go API backend using Gin to upload videos using chunking and ReactJS API calls <br/>
Used Redis for chunked uploads for tracking upload sessions <br/>
Metadata for files (location, resolutions, length) is saved in Postgres <br/>
Uploaded files are saved to MinIO (S3-compatible endpoints, just switch endpoint from localhost to s3 host) <br/>
FFmpeg is used for thumbnail and resolution generation, following best practices <br/>
* Should switch to HLS for auto resolution switching in the future 

# Endpoints
Follow the `frontend_api_call.js` file for upload video frontend call samples <br/>
Designed for scalable video handling - following chunking method <br/> 
**POST** `/upload/init` initialize a new upload session <br/>
**POST** `/upload/chunk` upload a single video chunk <br/>
**POST** `/upload/complete` merge chunks, generate thumbnail and resolutions, upload metadata to Postgres <br/>
**GET** `/video/:video_id` get metadata for videos -> can access through the resolutions links <br/>
