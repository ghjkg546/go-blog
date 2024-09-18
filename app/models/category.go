package models

const TableNameCategory = "category"

// Category mapped from table <category>
type Category struct {
	ID            int32  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name          string `gorm:"column:name;comment:标题" json:"name"`               // 标题
	Slug          string `gorm:"column:slug;comment:slug" json:"slug"`             // 别名
	ParentID      int32  `gorm:"column:parent_id;comment:父级id" json:"parent_id"`   // 父级id
	Status        int32  `gorm:"column:status;default:1;comment:状态" json:"status"` // 状态
	CreatedAt     int64  `gorm:"autoCreateTime"`                                   // 使用时间戳秒数填充创建时间
	UpdatedAt     int64  `gorm:"autoUpdateTime"`                                   // 使用时间戳毫秒数填充更新时间
	CreateTimeStr string `gorm:"-" json:"create_time_str"`
	UpdateTimeStr string `gorm:"-" json:"update_time_str"`
}

// TableName Category's table name
func (*Category) TableName() string {
	return TableNameCategory
}
