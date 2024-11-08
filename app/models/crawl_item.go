package models

// Category mapped from table <category>
type CrawlItem struct {
	ID            int32  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name          string `gorm:"column:name;comment:标题" json:"name"` // 标题
	Url           string `gorm:"column:url;comment:url" json:"url"`  // 别名
	CreatedAt     int64  `gorm:"autoCreateTime"`                     // 使用时间戳秒数填充创建时间
	CreateTimeStr string `gorm:"-" json:"create_time_str"`
}

// TableName Category's table name
func (*CrawlItem) TableName() string {
	return "crawl_item"
}
