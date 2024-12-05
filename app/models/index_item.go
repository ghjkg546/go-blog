package models

type IndexItem struct {
	ID    int32  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Title string `gorm:"column:title;comment:标题" json:"name"` // 标题

}

//// TableName Category's table name
//func (*Category) TableName() string {
//	return TableNameCategory
//}
