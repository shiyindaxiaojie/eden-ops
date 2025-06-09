package model

import (
	"time"

	"gorm.io/gorm"
)

// Role 角色模型
type Role struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Name      string         `gorm:"size:32;uniqueIndex;not null" json:"name"`
	Code      string         `gorm:"size:32;uniqueIndex;not null" json:"code"`
	Status    int            `gorm:"default:1" json:"status"` // 1: 正常, 0: 禁用
	Remark    string         `gorm:"size:255" json:"remark"`
	MenuIDs   []uint         `gorm:"-" json:"menu_ids,omitempty"`
	Menus     []*Menu        `gorm:"many2many:sys_role_menu;" json:"menus,omitempty"`
	Users     []*User        `gorm:"many2many:sys_user_role;" json:"users,omitempty"`
}

// TableName 表名
func (Role) TableName() string {
	return "sys_role"
}

// AfterFind 查询后钩子
func (r *Role) AfterFind(tx *gorm.DB) error {
	r.MenuIDs = make([]uint, 0, len(r.Menus))
	for _, menu := range r.Menus {
		r.MenuIDs = append(r.MenuIDs, menu.ID)
	}
	return nil
}

// RoleMenu 角色菜单关联模型
type RoleMenu struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	RoleID    uint      `gorm:"not null" json:"role_id"`
	MenuID    uint      `gorm:"not null" json:"menu_id"`
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

// TableName 表名
func (RoleMenu) TableName() string {
	return "sys_role_menu"
}

// UserRole 用户角色关联模型
type UserRole struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	RoleID    uint      `gorm:"not null" json:"role_id"`
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

// TableName 表名
func (UserRole) TableName() string {
	return "sys_user_role"
}
