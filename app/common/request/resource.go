package request

type BatchSave struct {
	Content    string `form:"content" json:"content"`
	CategoryId uint   `form:"category_id" json:"category_id"`
	DiskTypeId int    `form:"disk_type_id" json:"disk_type_id"`
}

type Crawl struct {
	DetailUrl string `form:"detail_url" json:"detail_url"`
	NameRule  string `form:"name_rule" json:"name_rule"`
	LinkRule  string `form:"link_rule" json:"link_rule"`
}

type TransSave struct {
	Ids string `form:"ids" json:"ids"`
	Fid string `form:"fid" json:"fid"`
}

type BatchShare struct {
	PageSize   int    `form:"page_size" json:"page_size"`
	Fid        string `form:"fid" json:"fid"`
	CategoryId uint   `form:"category_id" json:"category_id"`
}
