package model

type UploadSession struct {
	UploadID  string
	UserID    string
	FileName  string
	Chunks    int
	Uploaded  int
	Status    string
	CreatedAt int64
	UpdatedAt int64
}
