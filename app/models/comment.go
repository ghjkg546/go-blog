package models

const TableNameComment = "comment"

// Comment mapped from table <comment>
type Comment struct {
	ID             int32  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Content        string `gorm:"column:content" json:"content"`
	ParentID       int32  `gorm:"column:parent_id" json:"parent_id"`
	UserID         int32  `gorm:"column:user_id" json:"user_id"`
	ResourceItemId int32  `gorm:"column:resource_item_id" json:"resource_item_id"`
	CreatedAt      int64  `gorm:"autoCreateTime"` // 使用时间戳秒数填充创建时间
	CreatedAtStr   string `gorm:"_" json:"create_time_str"`
	UpdatedAt      int64  `gorm:"autoUpdateTime"` // 使用时间戳毫秒数填充更新时间
	User           User   `gorm:"foreignKey:user_id" json:"user"`
}

// TableName Comment's table name
func (*Comment) TableName() string {
	return TableNameComment
}
