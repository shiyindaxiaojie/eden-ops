package service

import (
	"eden-ops/internal/model"
	"eden-ops/internal/repository"
	"eden-ops/pkg/auth"
	"errors"
	"fmt"
	"log"

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
	log.Printf("尝试登录: 用户名=%s", username)

	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		log.Printf("用户不存在: %v", err)
		return "", fmt.Errorf("用户不存在: %v", err)
	}

	log.Printf("找到用户: ID=%d, 用户名=%s", user.ID, user.Username)

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Printf("密码验证失败: %v", err)
		return "", errors.New("用户名或密码错误")
	}

	log.Printf("密码验证成功，生成token")

	token, err := s.jwtAuth.GenerateToken(user.ID, user.Username)
	if err != nil {
		log.Printf("生成token失败: %v", err)
		return "", fmt.Errorf("生成token失败: %v", err)
	}

	log.Printf("登录成功: 用户ID=%d, 用户名=%s", user.ID, user.Username)

	return token, nil
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
