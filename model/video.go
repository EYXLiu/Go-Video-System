package model

type Video struct {
	VideoID     string
	UserID      string
	FileName    string
	Status      string
	Duration    float64
	Resolutions map[string]string
	CreatedAt   int64
	UpdatedAt   int64
}
