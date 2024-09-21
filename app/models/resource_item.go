package models

import (
	"strconv"
)

type NetDiskItem struct {
	Url  string `json:"url"`
	Type int    `json:"type"`
}

type ResourceItem struct {
	ID
	Title          string        `json:"name" gorm:"size:200;not null;comment:用户名称"`
	CategoryId     uint          `json:"category_id" gorm:"comment:用户名称"`
	Description    string        `json:"description" gorm:"size:255;not null;index;comment:用户手机号"`
	CoverImg       string        `json:"cover_img" gorm:"size:255;comment:封面图"`
	DiskItems      string        `json:"disk_items" gorm:"not null;default:'';comment:网盘信息"`
	DiskItemsArray []NetDiskItem ` json:"disk_items_array" gorm:"-"`
	TagIds         string        `json:"tag_ids" gorm:"not null;default:'';comment:tag"`
	SearchId       string        `json:"search_id" gorm:"comment:搜索id"`
	Status         uint          `json:"status" gorm:"comment:状态"`
	Views          uint          `json:"views" gorm:"comment:阅读次数"`
	CreatedAt      int64         `gorm:"autoCreateTime"` // 使用时间戳秒数填充创建时间
	UpdatedAt      int64         `gorm:"autoUpdateTime"` // 使用时间戳毫秒数填充更新时间
	CreateTimeStr  string        `gorm:"-" json:"create_time_str"`
	UpdateTimeStr  string        `gorm:"-" json:"update_time_str"`
	Url            string        `gorm:"-" json:"url"`
}

func (res ResourceItem) GetUid() string {
	return strconv.Itoa(int(res.ID.ID))
}

// TableName 指定表名
func (res ResourceItem) TableName() string {
	return "resource_item"
}
