package repository

import (
	"eden-ops/internal/model"

	"gorm.io/gorm"
)

// RoleRepository 角色仓储接口
type RoleRepository interface {
	Create(role *model.Role) error
	Update(role *model.Role) error
	Delete(id uint) error
	FindByID(id uint) (*model.Role, error)
	List(page, pageSize int) ([]*model.Role, int64, error)
	AssignMenus(roleID uint, menuIDs []uint) error
	AssignUserRoles(userID uint, roleIDs []uint) error
}

// RoleRepositoryImpl 角色仓储实现
type RoleRepositoryImpl struct {
	db *gorm.DB
}

// NewRoleRepository 创建角色仓储实例
func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &RoleRepositoryImpl{db: db}
}

// Create 创建角色
func (r *RoleRepositoryImpl) Create(role *model.Role) error {
	return r.db.Create(role).Error
}

// Update 更新角色
func (r *RoleRepositoryImpl) Update(role *model.Role) error {
	return r.db.Save(role).Error
}

// Delete 删除角色
func (r *RoleRepositoryImpl) Delete(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 删除角色菜单关联
		if err := tx.Where("role_id = ?", id).Delete(&model.RoleMenu{}).Error; err != nil {
			return err
		}
		// 删除角色用户关联
		if err := tx.Where("role_id = ?", id).Delete(&model.UserRole{}).Error; err != nil {
			return err
		}
		// 删除角色
		return tx.Delete(&model.Role{}, id).Error
	})
}

// FindByID 根据ID查找角色
func (r *RoleRepositoryImpl) FindByID(id uint) (*model.Role, error) {
	var role model.Role
	err := r.db.Preload("Menus").First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// List 获取角色列表
func (r *RoleRepositoryImpl) List(page, pageSize int) ([]*model.Role, int64, error) {
	var roles []*model.Role
	var total int64

	err := r.db.Model(&model.Role{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Preload("Menus").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&roles).Error

	return roles, total, err
}

// AssignMenus 分配菜单
func (r *RoleRepositoryImpl) AssignMenus(roleID uint, menuIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 删除原有的角色菜单关联
		if err := tx.Where("role_id = ?", roleID).Delete(&model.RoleMenu{}).Error; err != nil {
			return err
		}

		// 创建新的角色菜单关联
		roleMenus := make([]model.RoleMenu, len(menuIDs))
		for i, menuID := range menuIDs {
			roleMenus[i] = model.RoleMenu{
				RoleID: roleID,
				MenuID: menuID,
			}
		}

		return tx.Create(&roleMenus).Error
	})
}

// AssignUserRoles 分配用户角色
func (r *RoleRepositoryImpl) AssignUserRoles(userID uint, roleIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 删除原有的用户角色关联
		if err := tx.Where("user_id = ?", userID).Delete(&model.UserRole{}).Error; err != nil {
			return err
		}

		// 创建新的用户角色关联
		userRoles := make([]model.UserRole, len(roleIDs))
		for i, roleID := range roleIDs {
			userRoles[i] = model.UserRole{
				UserID: userID,
				RoleID: roleID,
			}
		}

		return tx.Create(&userRoles).Error
	})
}
