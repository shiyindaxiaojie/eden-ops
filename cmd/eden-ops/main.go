package main

import (
	"context"
	"eden-ops/internal/handler"
	"eden-ops/internal/initialize"
	"eden-ops/internal/repository"
	"eden-ops/internal/router"
	"eden-ops/internal/service"
	"eden-ops/internal/task"
	"eden-ops/pkg/auth"
	"eden-ops/pkg/config"
	"eden-ops/pkg/logger"
	"fmt"
	"log"

	"github.com/sirupsen/logrus"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfigFromYAML("configs/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化日志
	if err := logger.Init(cfg.Log); err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}

	// 初始化数据库
	db, err := initialize.InitDB(cfg.Database)
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	// 初始化JWT
	jwtAuth := auth.NewJWTAuth(cfg.JWT.Secret, cfg.JWT.Expire)

	// 初始化仓库
	cloudAccountRepo := repository.NewCloudAccountRepository(db)
	databaseConfigRepo := repository.NewDatabaseConfigRepository(db)
	serverConfigRepo := repository.NewServerConfigRepository(db)
	k8sConfigRepo := repository.NewK8sConfigRepository(db)
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	menuRepo := repository.NewMenuRepository(db)

	// 初始化服务
	cloudAccountService := service.NewCloudAccountService(cloudAccountRepo)
	databaseConfigService := service.NewDatabaseConfigService(databaseConfigRepo)
	serverConfigService := service.NewServerConfigService(serverConfigRepo)
	k8sConfigService := service.NewK8sConfigService(k8sConfigRepo)
	userService := service.NewUserService(userRepo, jwtAuth)
	roleService := service.NewRoleService(roleRepo)
	menuService := service.NewMenuService(menuRepo)

	loggerInstance := logrus.New()

	// 初始化处理器
	cloudAccountHandler := handler.NewCloudAccountHandler(cloudAccountService)
	databaseConfigHandler := handler.NewDatabaseConfigHandler(databaseConfigService)
	serverConfigHandler := handler.NewServerConfigHandler(serverConfigService, loggerInstance)
	k8sConfigHandler := handler.NewK8sConfigHandler(k8sConfigService)
	userHandler := handler.NewUserHandler(userService, jwtAuth)
	roleHandler := handler.NewRoleHandler(roleService)
	menuHandler := handler.NewMenuHandler(menuService)
	authHandler := handler.NewAuthHandler(userService, jwtAuth)

	// 初始化任务
	k8sSyncTask := task.NewK8sSyncTask(db, k8sConfigService, loggerInstance)
	go k8sSyncTask.Start(context.Background())

	// 初始化路由
	r := router.NewRouter(
		jwtAuth,
		cloudAccountHandler,
		databaseConfigHandler,
		serverConfigHandler,
		k8sConfigHandler,
		userHandler,
		roleHandler,
		menuHandler,
		authHandler,
	)

	// 启动服务器
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	if err := r.Run(addr); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
