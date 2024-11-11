package models

// FileInfo represents each item in the "list" array in "data"
type FileInfo struct {
	Fid           string `json:"fid"`
	FileName      string `json:"file_name"`
	PdirFid       string `json:"pdir_fid"`
	Category      int    `json:"category"`
	FileType      int    `json:"file_type"`
	CreatedAt     int64  `json:"created_at"`
	CreateTimeStr string `json:"create_time_str"`
}

type ShareItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
