package models

// SysDict mapped from table <sys_dict_item>
type SysDict struct {
	ID         int32         `gorm:"column:id;primaryKey" json:"id"`
	Name       string        `gorm:"column:name" json:"name"`
	Code       string        `gorm:"column:code" json:"code"`
	Status     int32         `gorm:"column:status" json:"status"`
	Remark     string        `gorm:"column:remark" json:"remark"`
	CreateTime int32         `gorm:"column:create_time" json:"create_time"`
	UpdateTime int32         `gorm:"column:update_time" json:"update_time"`
	DictItems  []SysDictItem `gorm:"foreignKey:dict_id" json:"dictItems"`
}

// TableName SysDictItem's table name
func (*SysDict) TableName() string {
	return "sys_dict"
}
