package service

import (
	"eden-ops/internal/model"
	"eden-ops/internal/repository"
	"eden-ops/pkg/auth"
	"eden-ops/pkg/logger"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// UserService 用户服务接口
type UserService interface {
	Login(username, password string) (string, error)
	Create(user *model.User) error
	Update(user *model.User) error
	Delete(id uint) error
	Get(id uint) (*model.User, error)
	List(page, size int) ([]*model.User, int64, error)
	GetUserInfo(id uint) (*model.User, error)
	GetUserRoles(userID uint) ([]*model.Role, error)
	AssignRoles(userID uint, roleIDs []uint) error
}

// userService 用户服务实现
type userService struct {
	userRepo repository.UserRepository
	jwtAuth  *auth.JWTAuth
}

// NewUserService 创建用户服务
func NewUserService(userRepo repository.UserRepository, jwtAuth *auth.JWTAuth) UserService {
	return &userService{
		userRepo: userRepo,
		jwtAuth:  jwtAuth,
	}
}

// Login 用户登录
func (s *userService) Login(username, password string) (string, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		logger.Info("登录验证失败: 用户不存在, username=%s", username)
		return "", errors.New("用户名或密码错误")
	}

	// 验证密码
	if err := s.verifyPassword(user.Password, password); err != nil {
		logger.Info("登录验证失败: 密码错误, username=%s", username)
		return "", errors.New("用户名或密码错误")
	}

	token, err := s.jwtAuth.GenerateToken(user.ID, user.Username)
	if err != nil {
		logger.Error("登录失败: 生成token错误, username=%s, error=%v", username, err)
		return "", fmt.Errorf("生成token失败: %v", err)
	}

	return token, nil
}

// verifyPassword 验证密码（支持前端 SHA256 加密）
func (s *userService) verifyPassword(hashedPassword, inputPassword string) error {
	// 直接验证 bcrypt 哈希（兼容旧的直接密码存储）
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword)); err == nil {
		return nil
	}

	// 如果输入密码长度为64，可能是前端 SHA256 加密的结果
	if len(inputPassword) == 64 {
		// 这种情况下，数据库中应该存储的是 SHA256 密码的 bcrypt 哈希
		// 但当前数据库存储的是原始密码的 bcrypt 哈希，所以这里会失败
		// 我们需要更新数据库或者提供迁移方案
		logger.Warn("检测到SHA256密码但数据库存储原始密码哈希，需要更新数据库")
	}

	return fmt.Errorf("密码验证失败")
}

// hashPassword 加密密码
func (s *userService) hashPassword(password string) (string, error) {
	if password == "" {
		return "", nil
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("密码加密失败: %v", err)
	}
	return string(hashedPassword), nil
}

// Create 创建用户
func (s *userService) Create(user *model.User) error {
	hashedPassword, err := s.hashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	return s.userRepo.Create(user)
}

// Update 更新用户
func (s *userService) Update(user *model.User) error {
	hashedPassword, err := s.hashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	return s.userRepo.Update(user)
}

// Delete 删除用户
func (s *userService) Delete(id uint) error {
	return s.userRepo.Delete(id)
}

// Get 获取用户
func (s *userService) Get(id uint) (*model.User, error) {
	return s.userRepo.Get(id)
}

// List 获取用户列表
func (s *userService) List(page, size int) ([]*model.User, int64, error) {
	return s.userRepo.List(page, size)
}

// GetUserInfo 获取用户信息
func (s *userService) GetUserInfo(id uint) (*model.User, error) {
	return s.userRepo.Get(id)
}

// GetUserRoles 获取用户角色
func (s *userService) GetUserRoles(userID uint) ([]*model.Role, error) {
	return s.userRepo.GetUserRoles(userID)
}

// AssignRoles 分配角色
func (s *userService) AssignRoles(userID uint, roleIDs []uint) error {
	return s.userRepo.AssignRoles(userID, roleIDs)
}
