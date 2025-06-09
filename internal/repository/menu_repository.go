package repository

import (
	"eden-ops/internal/model"

	"gorm.io/gorm"
)

// MenuRepository 菜单仓储接口
type MenuRepository interface {
	Create(menu *model.Menu) error
	Update(menu *model.Menu) error
	Delete(id uint) error
	FindByID(id uint) (*model.Menu, error)
	List() ([]*model.Menu, error)
	ListByRoleID(roleID uint) ([]*model.Menu, error)
}

// MenuRepositoryImpl 菜单仓储实现
type MenuRepositoryImpl struct {
	db *gorm.DB
}

// NewMenuRepository 创建菜单仓储实例
func NewMenuRepository(db *gorm.DB) MenuRepository {
	return &MenuRepositoryImpl{db: db}
}

// Create 创建菜单
func (r *MenuRepositoryImpl) Create(menu *model.Menu) error {
	return r.db.Create(menu).Error
}

// Update 更新菜单
func (r *MenuRepositoryImpl) Update(menu *model.Menu) error {
	return r.db.Save(menu).Error
}

// Delete 删除菜单
func (r *MenuRepositoryImpl) Delete(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 删除角色菜单关联
		if err := tx.Where("menu_id = ?", id).Delete(&model.RoleMenu{}).Error; err != nil {
			return err
		}
		// 删除菜单
		return tx.Delete(&model.Menu{}, id).Error
	})
}

// FindByID 根据ID查找菜单
func (r *MenuRepositoryImpl) FindByID(id uint) (*model.Menu, error) {
	var menu model.Menu
	err := r.db.First(&menu, id).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

// List 获取菜单列表
func (r *MenuRepositoryImpl) List() ([]*model.Menu, error) {
	var menus []*model.Menu
	err := r.db.Order("sort").Find(&menus).Error
	return menus, err
}

// ListByRoleID 根据角色ID获取菜单列表
func (r *MenuRepositoryImpl) ListByRoleID(roleID uint) ([]*model.Menu, error) {
	var menus []*model.Menu
	err := r.db.Joins("JOIN sys_role_menu ON sys_role_menu.menu_id = sys_menu.id").
		Where("sys_role_menu.role_id = ?", roleID).
		Order("sort").
		Find(&menus).Error
	return menus, err
}
