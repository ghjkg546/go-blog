package request

type BatchSave struct {
	Content    string `form:"content" json:"content"`
	CategoryId uint   `form:"category_id" json:"category_id"`
}
