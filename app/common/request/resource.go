package request

type BatchSave struct {
	Content    string `form:"content" json:"content"`
	CategoryId uint   `form:"category_id" json:"category_id"`
}

type BatchShare struct {
	PageSize   int    `form:"page_size" json:"page_size"`
	Fid        string `form:"fid" json:"fid"`
	CategoryId uint   `form:"category_id" json:"category_id"`
}
