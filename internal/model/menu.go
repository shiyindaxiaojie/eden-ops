package model

import (
	"time"

	"gorm.io/gorm"
)

// Menu 菜单模型
type Menu struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	ParentID  uint           `gorm:"default:0" json:"parent_id"`
	Name      string         `gorm:"size:32;not null" json:"name"`
	Path      string         `gorm:"size:128" json:"path"`
	Component string         `gorm:"size:128" json:"component"`
	Sort      int            `gorm:"default:0" json:"sort"`
	Icon      string         `gorm:"size:32" json:"icon"`
	Status    int            `gorm:"default:1" json:"status"` // 1: 正常, 0: 禁用
	Hidden    bool           `gorm:"default:false" json:"hidden"`
	Cache     bool           `gorm:"default:false" json:"cache"`
	Type      int            `gorm:"default:1" json:"type"` // 1: 菜单, 2: 按钮
	Roles     []*Role        `gorm:"many2many:sys_role_menu;" json:"roles,omitempty"`
	Children  []*Menu        `gorm:"-" json:"children,omitempty"`
}

// TableName 表名
func (Menu) TableName() string {
	return "sys_menu"
}

// AfterFind 查询后钩子
func (m *Menu) AfterFind(tx *gorm.DB) error {
	// 如果有子菜单，递归调用 AfterFind
	if len(m.Children) > 0 {
		for _, child := range m.Children {
			if err := child.AfterFind(tx); err != nil {
				return err
			}
		}
	}
	return nil
}
