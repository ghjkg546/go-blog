package models

type Feedback struct {
	ID           int32  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Content      string `gorm:"column:content" json:"content"`
	Message      string `gorm:"column:message" json:"message"`
	UserID       int32  `gorm:"column:user_id" json:"user_id"`
	Status       int32  `gorm:"column:status" json:"status"`
	CreatedAt    int64  `gorm:"autoCreateTime"` // 使用时间戳秒数填充创建时间
	CreatedAtStr string `gorm:"-" json:"create_time_str"`
	UpdatedAt    int64  `gorm:"autoUpdateTime"` // 使用时间戳毫秒数填充更新时间
	User         User   `gorm:"foreignKey:user_id" json:"user"`
}

func (*Feedback) TableName() string {
	return "feedback"
}
