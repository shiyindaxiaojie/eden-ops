package service

import (
	"eden-ops/internal/model"
	"eden-ops/internal/repository"
)

// MenuService 菜单服务接口
type MenuService interface {
	Create(menu *model.Menu) error
	Update(menu *model.Menu) error
	Delete(id uint) error
	Get(id uint) (*model.Menu, error)
	List() ([]*model.Menu, error)
	ListByRoleID(roleID uint) ([]*model.Menu, error)
}

// MenuServiceImpl 菜单服务实现
type MenuServiceImpl struct {
	menuRepo repository.MenuRepository
}

// NewMenuService 创建菜单服务实例
func NewMenuService(menuRepo repository.MenuRepository) MenuService {
	return &MenuServiceImpl{
		menuRepo: menuRepo,
	}
}

// Create 创建菜单
func (s *MenuServiceImpl) Create(menu *model.Menu) error {
	return s.menuRepo.Create(menu)
}

// Update 更新菜单
func (s *MenuServiceImpl) Update(menu *model.Menu) error {
	return s.menuRepo.Update(menu)
}

// Delete 删除菜单
func (s *MenuServiceImpl) Delete(id uint) error {
	return s.menuRepo.Delete(id)
}

// Get 获取菜单
func (s *MenuServiceImpl) Get(id uint) (*model.Menu, error) {
	return s.menuRepo.FindByID(id)
}

// List 获取菜单列表
func (s *MenuServiceImpl) List() ([]*model.Menu, error) {
	return s.menuRepo.List()
}

// ListByRoleID 根据角色ID获取菜单列表
func (s *MenuServiceImpl) ListByRoleID(roleID uint) ([]*model.Menu, error) {
	return s.menuRepo.ListByRoleID(roleID)
}
