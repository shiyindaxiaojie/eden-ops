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
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	// 加载配置
	log.Printf("正在加载配置...")
	cfg, err := config.LoadConfigFromYAML("configs/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
	log.Printf("配置加载成功")

	// 初始化日志
	log.Printf("正在初始化日志...")
	if err := logger.Init(cfg.Log); err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}

	// 获取日志实例
	logInstance := logger.GetLogger()
	logInstance.Info("日志初始化成功")

	// 初始化数据库
	logInstance.Info("正在初始化数据库...")
	db, err := initialize.InitDB(cfg.Database)
	if err != nil {
		logInstance.Fatalf("初始化数据库失败: %v", err)
	}
	logInstance.Info("数据库初始化成功")

	// 初始化JWT
	logInstance.Info("正在初始化JWT...")
	jwtAuth := auth.NewJWTAuth(cfg.JWT.Secret, cfg.JWT.Expire)
	logInstance.Info("JWT初始化成功")

	// 初始化仓库
	logInstance.Info("正在初始化仓库...")
	cloudAccountRepo := repository.NewCloudAccountRepository(db)
	databaseConfigRepo := repository.NewDatabaseConfigRepository(db)
	serverConfigRepo := repository.NewServerConfigRepository(db)
	k8sConfigRepo := repository.NewK8sConfigRepository(db)
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	menuRepo := repository.NewMenuRepository(db)
	logInstance.Info("仓库初始化成功")

	// 初始化服务
	logInstance.Info("正在初始化服务...")
	cloudAccountService := service.NewCloudAccountService(cloudAccountRepo)
	databaseConfigService := service.NewDatabaseConfigService(databaseConfigRepo)
	serverConfigService := service.NewServerConfigService(serverConfigRepo)
	k8sConfigService := service.NewK8sConfigService(k8sConfigRepo)
	userService := service.NewUserService(userRepo, jwtAuth)
	roleService := service.NewRoleService(roleRepo)
	menuService := service.NewMenuService(menuRepo)
	logInstance.Info("服务初始化成功")

	loggerInstance := logrus.New()

	// 初始化处理器
	logInstance.Info("正在初始化处理器...")
	cloudAccountHandler := handler.NewCloudAccountHandler(cloudAccountService)
	databaseConfigHandler := handler.NewDatabaseConfigHandler(databaseConfigService)
	serverConfigHandler := handler.NewServerConfigHandler(serverConfigService, loggerInstance)
	k8sConfigHandler := handler.NewK8sConfigHandler(k8sConfigService)
	userHandler := handler.NewUserHandler(userService, jwtAuth)
	roleHandler := handler.NewRoleHandler(roleService)
	menuHandler := handler.NewMenuHandler(menuService)
	authHandler := handler.NewAuthHandler(userService, jwtAuth)
	logInstance.Info("处理器初始化成功")

	// 初始化任务
	logInstance.Info("正在启动后台任务...")
	k8sSyncTask := task.NewK8sSyncTask(db, k8sConfigService, loggerInstance)
	go k8sSyncTask.Start(context.Background())
	logInstance.Info("后台任务启动成功")

	// 初始化路由
	logInstance.Info("正在初始化路由...")
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
	logInstance.Info("路由初始化成功")

	// 启动服务器
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logInstance.Infof("正在启动服务器，监听地址: %s", addr)

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// 在一个单独的goroutine中启动服务器
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logInstance.Fatalf("服务器启动失败: %v", err)
		}
	}()

	logInstance.Infof("Eden Ops 服务已成功启动，监听端口: %d", cfg.Server.Port)

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	// kill (无参数) 默认发送 syscall.SIGTERM
	// kill -2 发送 syscall.SIGINT
	// kill -9 发送 syscall.SIGKILL，但无法被捕获，所以不需要添加
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logInstance.Info("正在关闭服务器...")

	// 创建一个5秒的上下文用于超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logInstance.Fatalf("服务器强制关闭: %v", err)
	}

	logInstance.Info("服务器已优雅关闭")
}
