package models

import (
	"strconv"
)

type User struct {
	ID        int32  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	UserName  string `json:"username" gorm:"column:username;size:50;not null;comment:用户登录名"`
	Name      string `json:"name" gorm:"size:30;not null;comment:用户名称"`
	Mobile    string `json:"mobile" gorm:"size:24;not null;index;comment:用户手机号"`
	Avatar    string `json:"avatar" gorm:"comment:头像"`
	Email     string `json:"email" gorm:"size:200;comment:邮箱"`
	Password  string `json:"password" gorm:"not null;default:'';comment:用户密码"`
	Score     int64  `gorm:"column:score;" json:"score"`
	CreatedAt int64  `gorm:"autoCreateTime"` // 使用时间戳秒数填充创建时间
	UpdatedAt int64  `gorm:"autoUpdateTime"` // 使用时间戳毫秒数填充更新时间
}

// TableName User's table name
func (*User) TableName() string {
	return "user"
}

func (user User) GetUid() string {
	return strconv.Itoa(int(user.ID))
}

type File struct {
	Category int    `json:"category"`
	FID      string `json:"fid"`
	FileName string `json:"file_name"`
	FileType int    `json:"file_type"`
	PDirFID  string `json:"pdir_fid"`
}
