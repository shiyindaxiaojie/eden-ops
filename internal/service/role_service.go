package service

import (
	"eden-ops/internal/model"
	"eden-ops/internal/repository"
)

// RoleService 角色服务接口
type RoleService interface {
	Create(role *model.Role) error
	Update(role *model.Role) error
	Delete(id uint) error
	Get(id uint) (*model.Role, error)
	List(page, pageSize int) ([]*model.Role, int64, error)
	AssignMenus(roleID uint, menuIDs []uint) error
	AssignUserRoles(userID uint, roleIDs []uint) error
}

// RoleServiceImpl 角色服务实现
type RoleServiceImpl struct {
	roleRepo repository.RoleRepository
}

// NewRoleService 创建角色服务实例
func NewRoleService(roleRepo repository.RoleRepository) RoleService {
	return &RoleServiceImpl{
		roleRepo: roleRepo,
	}
}

// Create 创建角色
func (s *RoleServiceImpl) Create(role *model.Role) error {
	return s.roleRepo.Create(role)
}

// Update 更新角色
func (s *RoleServiceImpl) Update(role *model.Role) error {
	return s.roleRepo.Update(role)
}

// Delete 删除角色
func (s *RoleServiceImpl) Delete(id uint) error {
	return s.roleRepo.Delete(id)
}

// Get 获取角色
func (s *RoleServiceImpl) Get(id uint) (*model.Role, error) {
	return s.roleRepo.FindByID(id)
}

// List 获取角色列表
func (s *RoleServiceImpl) List(page, pageSize int) ([]*model.Role, int64, error) {
	return s.roleRepo.List(page, pageSize)
}

// AssignMenus 分配菜单
func (s *RoleServiceImpl) AssignMenus(roleID uint, menuIDs []uint) error {
	return s.roleRepo.AssignMenus(roleID, menuIDs)
}

// AssignUserRoles 分配用户角色
func (s *RoleServiceImpl) AssignUserRoles(userID uint, roleIDs []uint) error {
	return s.roleRepo.AssignUserRoles(userID, roleIDs)
}
