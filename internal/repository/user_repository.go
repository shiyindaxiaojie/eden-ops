package repository

import (
	"eden-ops/internal/model"
	"log"

	"gorm.io/gorm"
)

// UserRepository 用户仓库接口
type UserRepository interface {
	Create(user *model.User) error
	Update(user *model.User) error
	Delete(id uint) error
	Get(id uint) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	List(page, size int) ([]*model.User, int64, error)
	GetUserRoles(userID uint) ([]*model.Role, error)
	AssignRoles(userID uint, roleIDs []uint) error
}

// userRepository 用户仓库实现
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓库
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create 创建用户
func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// Update 更新用户
func (r *userRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// Delete 删除用户
func (r *userRepository) Delete(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 删除用户角色关联
		if err := tx.Where("user_id = ?", id).Delete(&model.UserRole{}).Error; err != nil {
			return err
		}
		// 删除用户
		return tx.Delete(&model.User{}, id).Error
	})
}

// Get 获取用户
func (r *userRepository) Get(id uint) (*model.User, error) {
	var user model.User
	err := r.db.Preload("Roles").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户
func (r *userRepository) GetByUsername(username string) (*model.User, error) {
	log.Printf("查询用户: username=%s", username)

	var user model.User
	err := r.db.Preload("Roles").Where("username = ?", username).First(&user).Error
	if err != nil {
		log.Printf("查询用户失败: %v", err)
		return nil, err
	}

	log.Printf("查询用户成功: ID=%d, Username=%s", user.ID, user.Username)
	return &user, nil
}

// List 获取用户列表
func (r *userRepository) List(page, size int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	err := r.db.Model(&model.User{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Preload("Roles").Offset((page - 1) * size).Limit(size).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetUserRoles 获取用户角色
func (r *userRepository) GetUserRoles(userID uint) ([]*model.Role, error) {
	var roles []*model.Role
	err := r.db.Joins("JOIN sys_user_role ON sys_user_role.role_id = sys_role.id").
		Where("sys_user_role.user_id = ?", userID).
		Find(&roles).Error
	return roles, err
}

// AssignRoles 分配角色
func (r *userRepository) AssignRoles(userID uint, roleIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 删除所有用户角色关联
		if err := tx.Where("user_id = ?", userID).Delete(&model.UserRole{}).Error; err != nil {
			return err
		}
		// 创建新的用户角色关联
		for _, roleID := range roleIDs {
			if err := tx.Create(&model.UserRole{UserID: userID, RoleID: roleID}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
