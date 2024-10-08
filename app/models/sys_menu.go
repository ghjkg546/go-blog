// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package models

import (
	"time"
)

const TableNameSysMenu = "sys_menu"

// SysMenu 菜单管理
type SysMenu struct {
	ID         uint     `gorm:"column:id;primaryKey;autoIncrement:true;comment:ID" json:"id"`                                          // ID
	ParentID   uint     `gorm:"column:parent_id;not null;comment:父菜单ID" json:"parent_id"`                                           // 父菜单ID
	TreePath   string    `gorm:"column:tree_path;comment:父节点ID路径" json:"tree_path"`                                               // 父节点ID路径
	Name       string    `gorm:"column:name;not null;comment:菜单名称" json:"name"`                                                    // 菜单名称
	Type       int32     `gorm:"column:type;not null;comment:菜单类型（1-菜单 2-目录 3-外链 4-按钮）" json:"type"`                       // 菜单类型（1-菜单 2-目录 3-外链 4-按钮）
	RouteName  string    `gorm:"column:route_name;comment:路由名称（Vue Router 中用于命名路由）" json:"route_name"`                      // 路由名称（Vue Router 中用于命名路由）
	RoutePath  string    `gorm:"column:route_path;comment:路由路径（Vue Router 中定义的 URL 路径）" json:"route_path"`                   // 路由路径（Vue Router 中定义的 URL 路径）
	Component  string    `gorm:"column:component;comment:组件路径（组件页面完整路径，相对于 src/views/，缺省后缀 .vue）" json:"component"` // 组件路径（组件页面完整路径，相对于 src/views/，缺省后缀 .vue）
	Perm       string     `gorm:"column:perm;comment:【按钮】权限标识" json:"perm"`                                                      // 【按钮】权限标识
	AlwaysShow int32      `gorm:"column:always_show;comment:【目录】只有一个子路由是否始终显示（1-是 0-否）" json:"always_show"`           // 【目录】只有一个子路由是否始终显示（1-是 0-否）
	KeepAlive  int32      `gorm:"column:keep_alive;comment:【菜单】是否开启页面缓存（1-是 0-否）" json:"keep_alive"`                       // 【菜单】是否开启页面缓存（1-是 0-否）
	Visible    int32      `gorm:"column:visible;not null;default:1;comment:显示状态（1-显示 0-隐藏）" json:"visible"`                    // 显示状态（1-显示 0-隐藏）
	Sort       int32      `gorm:"column:sort;comment:排序" json:"sort"`                                                                // 排序
	Icon       string     `gorm:"column:icon;comment:菜单图标" json:"icon"`                                                            // 菜单图标
	Redirect   string     `gorm:"column:redirect;comment:跳转路径" json:"redirect"`                                                    // 跳转路径
	CreatedAt time.Time   `gorm:"<-:create" column:created_at;comment:创建时间 " json:"create_time"`                                   // 创建时间
	UpdatedAt time.Time   `gorm:"column:updated_at;comment:更新时间" json:"update_time"`                                               // 更新时间
	Params     string     `gorm:"column:params;comment:路由参数" json:"params"`                                                        // 路由参数
	Children    []SysMenu `gorm:"-" json:"children"`                                                                                   //
	Level    uint         `gorm:"-" json:"level"`                                                                                      //
}



type Tree struct {
	SysMenu
	Children []Tree `json:"children"`
}

// TableName SysMenu's table name
func (*SysMenu) TableName() string {
	return TableNameSysMenu
}
