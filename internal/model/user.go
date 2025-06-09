package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Username  string         `gorm:"size:32;uniqueIndex;not null" json:"username"`
	Password  string         `gorm:"size:128;not null" json:"-"`
	Nickname  string         `gorm:"size:32" json:"nickname"`
	Email     string         `gorm:"size:128" json:"email"`
	Phone     string         `gorm:"size:32" json:"phone"`
	Avatar    string         `gorm:"size:255" json:"avatar"`
	Status    int            `gorm:"default:1" json:"status"` // 1: 正常, 0: 禁用
	RoleIDs   []int64        `gorm:"-" json:"role_ids,omitempty"`
	Roles     []*Role        `gorm:"many2many:sys_user_role;" json:"roles,omitempty"`
}

// TableName 表名
func (User) TableName() string {
	return "sys_user"
}

// AfterFind 查询后钩子
func (u *User) AfterFind(tx *gorm.DB) error {
	u.RoleIDs = make([]int64, 0, len(u.Roles))
	for _, role := range u.Roles {
		u.RoleIDs = append(u.RoleIDs, int64(role.ID))
	}
	return nil
}
