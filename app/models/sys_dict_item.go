package models

const TableNameSysDictItem = "sys_dict_item"

// SysDictItem mapped from table <sys_dict_item>
type SysDictItem struct {
	ID         int32  `gorm:"column:id;primaryKey" json:"id"`
	DictID     int32  `gorm:"column:dict_id" json:"dict_id"`
	Name       string `gorm:"column:name" json:"name"`
	Value      string `gorm:"column:value" json:"value"`
	Status     int32  `gorm:"column:status" json:"status"`
	Sort       int32  `gorm:"column:sort" json:"sort"`
	Remark     string `gorm:"column:remark" json:"remark"`
	CreateTime int32  `gorm:"column:create_time" json:"create_time"`
	UpdateTime int32  `gorm:"column:update_time" json:"update_time"`
}

// TableName SysDictItem's table name
func (*SysDictItem) TableName() string {
	return TableNameSysDictItem
}
